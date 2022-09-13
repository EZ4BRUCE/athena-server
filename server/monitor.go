package server

import (
	"context"
	"time"

	pb "github.com/EZ4BRUCE/athena-proto/proto"
	"github.com/EZ4BRUCE/athena-server/global"
	"github.com/EZ4BRUCE/athena-server/internal/service"
)

// 定时检测某主机状态是否正常
func checkAlive(agent *Agent) {
	// 若在指定时间的2倍内没有上报数据，则视为连接异常
	var ticker *time.Ticker = time.NewTicker(time.Duration(agent.CheckAliveTime) * time.Second)
	for range ticker.C {
		if agent.CheckAliveStatus {
			agent.CheckAliveStatus = false
		} else {
			// 断开连接并释放资源
			global.Logger.Errorf("[连接断开] 主机 %s(%s) 出现连接异常，需要邮件通知\n", agent.UId, agent.Description)
			go sendOfflineEmail(agent)
			agent.IsDead = true
			release(agent)
			return
		}
	}
}

// 释放断开连接的主机的资源
func release(agent *Agent) {
	for k := range agent.MetricMap {
		// 遍历每一个指标的channel将其关闭，关闭过后对应的monitor协程会退出
		close(agent.MetricMap[k])
	}
	RegisterMap.Delete(agent.UId)
}

// 监控单个指标对应的channel
// 这里应该要传入一个agent，这样不用每次都去读sync.Map
func monitor(agent *Agent, metricChan chan *pb.ReportReq, aggregationTime int32) {
	// 聚合计数器，每当收到的report到达聚合时间的数目即需要进行聚合
	counter := aggregationTime
	var list []*pb.ReportReq
	for report := range metricChan {
		// 每次收到该agent上报，证明其连接正常，将检测变量设为true
		agent.CheckAliveStatus = true
		list = append(list, report)
		// 更新聚合计数器
		counter--
		if counter == 0 {
			// 执行一次聚合操作
			go doAggregation(agent, list, report)
			// 清空list
			list = make([]*pb.ReportReq, global.RPCSetting.AggregationTime)
			// 计数器设置回初值
			counter = aggregationTime
		}
	}
}

// 聚合操作函数
func doAggregation(agent *Agent, list []*pb.ReportReq, report *pb.ReportReq) {
	// 规则服务，用于从规则数据库中读取规则
	ruleSvc := service.NewRuleService(context.Background())
	// 告警服务，用于保存聚合信息以及执行告警动作
	reportSvc := service.NewReportService(context.Background())
	safe := true
	// 从数据库中取出该指标对应的所有聚合器(每个聚合器包含对应的聚合函数和告警规则)
	aggregators, err := ruleSvc.SearchAggregators(report.GetMetric())
	if err != nil {
		global.Logger.Errorf("[规则载入错误] ruleSvc.SearchAggregators err:%s", err)
		return
	}

	// 对该指标的每个聚合器进行告警判断
	for _, aggregator := range aggregators {
		result, danger := ruleSvc.ExecuteFunc(aggregator.Function, list)
		if danger {
			// 系统异常，需要告警
			safe = false
			global.Logger.Infof("[监控异常] Timestamp:%v 主机%s(%s)  指标 %s 出现异常，%s 型函数聚合值为 %v 告警等级为 %s ，需要执行 %s 动作\n", report.GetTimestamp(), agent.UId, agent.Description, aggregator.Metric, aggregator.Function.Type, result, aggregator.Rule.Level, aggregator.Rule.Action)
			// 执行告警动作(发邮件)可以开一个协程去做，不需要阻塞
			// go ruleSvc.ExecuteRule(report, &aggregator, result)
			ruleSvc.ExecuteRule(report, &aggregator, result, agent.Description)
			err = reportSvc.CreateWarningEvent(list, aggregator.Id, aggregator.Name, aggregator.Metric, aggregator.Function, aggregator.Rule, result)
			if err != nil {
				global.Logger.Errorf("[数据库错误] reportSvc.CreateWarningEvent err:%s", err)
			}
		}
	}
	if safe {
		// 系统正常
		global.Logger.Infof("[监控正常] Timestamp:%v 主机%s(%s) 指标 %s 正常，无需告警\n", report.GetTimestamp(), agent.UId, agent.Description, report.GetMetric())
		// 保存聚合事件
		err = reportSvc.CreateNormalEvent(list)
		if err != nil {
			global.Logger.Errorf("[数据库错误] reportSvc.CreateNormalEvent err:%s", err)
		}
	}
}
