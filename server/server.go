package server

import (
	"context"
	"fmt"
	"log"

	"github.com/EZ4BRUCE/athena-server/global"
	"github.com/EZ4BRUCE/athena-server/internal/service"
	pb "github.com/EZ4BRUCE/athena-server/proto"
	"github.com/google/uuid"
)

type ReportServerServer struct{}

func NewReportServer() *ReportServerServer {
	return &ReportServerServer{}
}

// 每个连接单独处理要在注册的时候用，因为注册只会注册一次
// 防止重复注册
func (s *ReportServerServer) Register(ctx context.Context, r *pb.RegisterReq) (*pb.RegisterRsp, error) {
	// 分配UId
	uId := uuid.New().String()
	// 将该agent注册到服务端
	global.RegisterMap[uId] = struct{}{}
	global.ReportMap[uId] = make(chan *pb.ReportReq, global.RPCSetting.AggregationTime*2)
	go monitor(global.ReportMap[uId])
	return &pb.RegisterRsp{Code: 10000001, UId: uId, Msg: "注册成功"}, nil
}

func monitor(reportMap chan *pb.ReportReq) {
	counter := global.RPCSetting.AggregationTime
	ruleSvc := service.NewRuleService(context.Background())
	reportSvc := service.NewReportService(context.Background())
	var list []*pb.ReportReq
	// 这里记得开一个for循环

	for report := range reportMap {
		list = append(list, report)
		counter--
		fmt.Println(counter)
		if counter == 0 {
			// 从数据库中取出该指标对应的所有聚合器(每个聚合器包含对应的聚合函数和告警规则)
			aggregators, err := ruleSvc.SearchAggregators(report.GetMetric())
			if err != nil {
				log.Printf("ruleSvc.SearchAggregators err:%s", err)
			}
			fmt.Println(report.GetMetric())
			for _, aggregator := range aggregators {
				result, danger := ruleSvc.ExecuteFunc(aggregator.Function, list)
				if danger {
					// 系统异常，需要告警
					log.Printf("指标 %s 出现异常，%s 型函数聚合值为 %v 告警等级为 %s ，需要执行 %s 动作\n", aggregator.Metric, aggregator.Function.Type, result, aggregator.Rule.Level, aggregator.Rule.Action)
					err = reportSvc.CreateWarningEvent(list, aggregator.Id, aggregator.Name, aggregator.Metric, aggregator.Function, aggregator.Rule, result)
					if err != nil {
						log.Printf("reportSvc.CreateWarningEvent err:%s", err)
					}
				} else {
					// 系统正常
					err = reportSvc.CreateNormalEvent(list)
					if err != nil {
						log.Printf("reportSvc.CreateNormalEvent err:%s", err)
					}
				}
			}

			list = make([]*pb.ReportReq, global.RPCSetting.AggregationTime)
			counter = global.RPCSetting.AggregationTime
		}
	}

}

func (s *ReportServerServer) Report(ctx context.Context, r *pb.ReportReq) (*pb.ReportRsp, error) {
	_, ok := global.RegisterMap[r.GetUId()]
	if !ok {
		log.Printf("收到未注册Agent上报！time:%v, metric:%s, dimensions:%v, value:%v\n", r.GetTimestamp(), r.GetMetric(), r.GetDimensions(), r.GetValue())
		return &pb.ReportRsp{Code: 10000001, Msg: "Agent未注册！"}, nil
	}
	global.ReportMap[r.GetUId()] <- r
	repoetSvc := service.NewReportService(ctx)
	err := repoetSvc.CreateReport(r)
	if err != nil {
		return &pb.ReportRsp{Code: 10000001, Msg: "插入数据库失败"}, err
	}
	// fmt.Printf("report: time:%v, metric:%s, dimensions:%v, value:%v \n", r.GetTimestamp(), r.GetMetric(), r.GetDimensions(), r.GetValue())
	return &pb.ReportRsp{Code: 10000001, Msg: "hello"}, nil
}
