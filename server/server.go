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

// agent注册函数，每个agent连接只会执行一次
func (s *ReportServerServer) Register(ctx context.Context, r *pb.RegisterReq) (*pb.RegisterRsp, error) {
	// 分配UId
	uId := uuid.New().String()
	// 将该agent注册到服务端用户注册表中，并为用户指标信息表分配内存
	global.RegisterMap[uId] = global.Agent{UId: uId, MetricMap: make(map[string]chan *pb.ReportReq, len(r.Metrics))}
	// 针对agent要监控的每一个指标创建聚合需要的存储channel
	for _, metric := range r.Metrics {
		global.RegisterMap[uId].MetricMap[metric] = make(chan *pb.ReportReq, global.RPCSetting.AggregationTime*2)
		// 对每一个指标使用一个协程监控并处理
		go monitor(global.RegisterMap[uId].MetricMap[metric])
	}
	log.Printf("New Agent: %v", uId)
	return &pb.RegisterRsp{Code: 10000000, UId: uId, Msg: "注册成功"}, nil
}

// 监控单个指标对应的channel
func monitor(reportMap chan *pb.ReportReq) {
	// 聚合计数器，每当收到的report到达聚合时间的数目即需要进行聚合
	counter := global.RPCSetting.AggregationTime
	// 规则服务，用于从规则数据库中读取规则
	ruleSvc := service.NewRuleService(context.Background())
	// 告警服务，用于保存聚合信息以及执行告警动作
	reportSvc := service.NewReportService(context.Background())
	var list []*pb.ReportReq
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
			// 对该指标的每个聚合器进行告警判断
			for _, aggregator := range aggregators {
				result, danger := ruleSvc.ExecuteFunc(aggregator.Function, list)
				if danger {
					// 系统异常，需要告警
					safe = false
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
			// 清空list
			list = make([]*pb.ReportReq, global.RPCSetting.AggregationTime)
			// 计数器设置回初值
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
	// 将收到的上报信息放入当前agent的指标表MetricMap的相应指标的channel中，由monitor处理
	global.RegisterMap[r.GetUId()].MetricMap[r.GetMetric()] <- r
	repoetSvc := service.NewReportService(ctx)
	err := repoetSvc.CreateReport(r)
	if err != nil {
		log.Printf("Create Report err:%s\n", err)
		return &pb.ReportRsp{Code: 10000002, Msg: "插入数据库失败"}, err
	}
	log.Printf("report: agent:%s, time:%v, metric:%s, dimensions:%v, value:%v \n", r.GetUId(), r.GetTimestamp(), r.GetMetric(), r.GetDimensions(), r.GetValue())
	return &pb.ReportRsp{Code: 10000000, Msg: "发送成功！"}, nil
}
