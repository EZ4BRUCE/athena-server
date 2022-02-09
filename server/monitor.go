package server

import (
	"context"
	"log"
	"time"

	pb "github.com/EZ4BRUCE/athena-proto/proto"
	"github.com/EZ4BRUCE/athena-server/global"
	"github.com/EZ4BRUCE/athena-server/internal/service"
)

func checkAlive(agent *Agent) {
	log.Println(agent)
	var ticker *time.Ticker = time.NewTicker(time.Duration(agent.CheckAliveTime) * time.Second)
	for range ticker.C {
		if agent.CheckAliveStatus {
			agent.CheckAliveStatus = false
		} else {
			log.Printf("[连接断开] 主机 %s 出现连接异常，需要邮件通知\n", agent.UId)
			go sendOfflineEmail(agent)
			agent.IsDead = true
			release(agent)
			return
		}
	}
}

// 释放断开连接的资源
func release(agent *Agent) {
	for k := range agent.MetricMap {
		// 遍历每一个指标的channel将其关闭，关闭过后对应的monitor协程会退出
		close(agent.MetricMap[k])
	}
	delete(RegisterMap, agent.UId)
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
		// 收到上报，将检测变量设为true
		RegisterMap[report.GetUId()].CheckAliveStatus = true
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
