package server

import (
	"context"
	"log"

	pb "github.com/EZ4BRUCE/athena-proto/proto"
	"github.com/EZ4BRUCE/athena-server/global"
	"github.com/EZ4BRUCE/athena-server/internal/service"
	"github.com/google/uuid"
)

type ReportServerServer struct{}

func NewReportServer() *ReportServerServer {
	return &ReportServerServer{}
}

// 每个连接单独处理要在注册的时候用，因为注册只会注册一次
func (s *ReportServerServer) Register(ctx context.Context, r *pb.RegisterReq) (*pb.RegisterRsp, error) {
	// 分配UId
	uId := uuid.New().String()
	// 将该agent注册到服务端
	global.RegisterMap[uId] = global.Agent{UId: uId, MetricMap: make(map[string]chan *pb.ReportReq, len(r.Metrics))}
	for _, metric := range r.Metrics {
		global.RegisterMap[uId].MetricMap[metric] = make(chan *pb.ReportReq, global.RPCSetting.AggregationTime*2)
		go monitor(global.RegisterMap[uId].MetricMap[metric])
	}
	log.Printf("New Agent: %v", uId)
	return &pb.RegisterRsp{Code: 10000001, UId: uId, Msg: "注册成功"}, nil
}

// 监控单个指标对应的channel
func monitor(reportMap chan *pb.ReportReq) {
	counter := global.RPCSetting.AggregationTime
	ruleSvc := service.NewRuleService(context.Background())
	reportSvc := service.NewReportService(context.Background())
	var list []*pb.ReportReq
	// 这里记得开一个for循环
	for report := range reportMap {
		list = append(list, report)
		counter--
		if counter == 0 {
			safe := true
			// 从数据库中取出该指标对应的所有聚合器(每个聚合器包含对应的聚合函数和告警规则)
			aggregators, err := ruleSvc.SearchAggregators(report.GetMetric())
			if err != nil {
				log.Printf("ruleSvc.SearchAggregators err:%s", err)
			}
			for _, aggregator := range aggregators {
				result, danger := ruleSvc.ExecuteFunc(aggregator.Function, list)
				if danger {
					safe = false
					// 系统异常，需要告警
					log.Printf("Timestamp:%v 指标 %s 出现异常，%s 型函数聚合值为 %v 告警等级为 %s ，需要执行 %s 动作\n", report.GetTimestamp(), aggregator.Metric, aggregator.Function.Type, result, aggregator.Rule.Level, aggregator.Rule.Action)
					// 执行告警动作(发邮件)可以开一个协程去做，不需要阻塞
					go ruleSvc.ExecuteRule(report, &aggregator, result)
					err = reportSvc.CreateWarningEvent(list, aggregator.Id, aggregator.Name, aggregator.Metric, aggregator.Function, aggregator.Rule, result)
					if err != nil {
						log.Printf("reportSvc.CreateWarningEvent err:%s", err)
					}
				}
			}
			if safe {
				// 系统正常
				log.Printf("Timestamp:%v 指标 %s 正常，无需告警\n", report.GetTimestamp(), report.GetMetric())
				err = reportSvc.CreateNormalEvent(list)
				if err != nil {
					log.Printf("reportSvc.CreateNormalEvent err:%s", err)
				}
			}
			list = make([]*pb.ReportReq, global.RPCSetting.AggregationTime)
			counter = global.RPCSetting.AggregationTime
		}
	}
}

func (s *ReportServerServer) Report(ctx context.Context, r *pb.ReportReq) (*pb.ReportRsp, error) {
	_, ok := global.RegisterMap[r.GetUId()]
	if !ok || r.GetUId() == "" {
		log.Printf("收到未注册Agent上报！time:%v, metric:%s, dimensions:%v, value:%v\n", r.GetTimestamp(), r.GetMetric(), r.GetDimensions(), r.GetValue())
		return &pb.ReportRsp{Code: 10000001, Msg: "Agent未注册！"}, nil
	}
	global.RegisterMap[r.GetUId()].MetricMap[r.GetMetric()] <- r
	repoetSvc := service.NewReportService(ctx)
	err := repoetSvc.CreateReport(r)
	if err != nil {
		return &pb.ReportRsp{Code: 10000002, Msg: "插入数据库失败"}, err
	}
	log.Printf("report: time:%v, metric:%s, dimensions:%v, value:%v \n", r.GetTimestamp(), r.GetMetric(), r.GetDimensions(), r.GetValue())
	return &pb.ReportRsp{Code: 10000000, Msg: "发送成功！"}, nil
}
