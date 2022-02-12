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

// RPC:agent注册函数，每个agent连接应只执行一次
func (s *ReportServerServer) Register(ctx context.Context, r *pb.RegisterReq) (*pb.RegisterRsp, error) {

	// 简单参数校验
	if r.GetCheckAliveTime() <= 0 || len(r.GetMetrics()) == 0 || r.GetAggregationTime() < 0 {
		log.Printf("[拒绝连接] 新增主机参数错误，拒绝连接")
		return &pb.RegisterRsp{Code: 10000003, UId: "", Msg: "注册失败：参数错误"}, nil
	}

	// 分配UId
	uId := uuid.New().String()
	// 获取聚合时间
	var aggregationTime int32
	if r.GetAggregationTime() <= 0 {
		aggregationTime = global.RPCSetting.AggregationTime
	} else {
		aggregationTime = r.GetAggregationTime()
	}

	// 将该agent注册到服务端用户注册表中，并为用户指标信息表分配内存
	RegisterMap[uId] = &Agent{
		UId:              uId,
		CheckAliveTime:   r.GetCheckAliveTime(),
		CheckAliveStatus: false, // 刚注册未发送设为false
		AggregationTime:  aggregationTime,
		IsDead:           false,
		MetricMap:        make(map[string]chan *pb.ReportReq, len(r.Metrics)),
	}

	// 针对该agent要监控的每一个指标创建指标channel（存储聚合数据）
	for _, metric := range r.Metrics {
		RegisterMap[uId].MetricMap[metric] = make(chan *pb.ReportReq, global.RPCSetting.AggregationTime*2)
		// 对每一个指标channel使用一个协程监控并处理
		go monitor(RegisterMap[uId].MetricMap[metric], RegisterMap[uId].AggregationTime)
	}
	// 发送新增主机通知邮件
	go sendLoginEmail(r, uId)
	log.Printf("[新增主机] 新注册 Agent: %v , 连接诊断频率: %d s , 聚合时间 %d 个单位上报粒度时间 ", uId, r.GetCheckAliveTime(), r.GetAggregationTime())
	// 新建协程监控该主机连接状态
	go checkAlive(RegisterMap[uId])

	return &pb.RegisterRsp{Code: 10000000, UId: uId, Msg: "注册成功"}, nil
}

// RPC:agent上报函数
func (s *ReportServerServer) Report(ctx context.Context, r *pb.ReportReq) (*pb.ReportRsp, error) {
	_, ok := RegisterMap[r.GetUId()]
	if !ok || r.GetUId() == "" {
		log.Printf("[拒绝接收] 收到未注册Agent上报！time:%v, metric:%s, dimensions:%v, value:%v\n", r.GetTimestamp(), r.GetMetric(), r.GetDimensions(), r.GetValue())
		return &pb.ReportRsp{Code: 10000001, Msg: "Agent 未注册"}, nil
	}
	// 将收到的上报信息放入当前agent的指标表MetricMap的相应指标的channel中，由monitor处理
	RegisterMap[r.GetUId()].MetricMap[r.GetMetric()] <- r
	repoetSvc := service.NewReportService(ctx)
	// 保存上报结果至数据库
	err := repoetSvc.SaveReport(r)
	if err != nil {
		log.Printf("[数据库错误] Create Report err:%s\n", err)
		return &pb.ReportRsp{Code: 10000002, Msg: "插入数据库失败"}, err
	}
	log.Printf("[接收成功] report: agent:%s, time:%v, metric:%s, dimensions:%v, value:%v \n", r.GetUId(), r.GetTimestamp(), r.GetMetric(), r.GetDimensions(), r.GetValue())
	return &pb.ReportRsp{Code: 10000000, Msg: "发送成功"}, nil
}
