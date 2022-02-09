package server

import (
	pb "github.com/EZ4BRUCE/athena-proto/proto"
)

// 全局定义每个agent主机的结构体
type Agent struct {
	UId              string
	CheckAliveTime   int32
	CheckAliveStatus bool
	IsDead           bool
	// 记录主机的最近上报信息（用于聚合数据）
	// key:metric value:存储对应metric的最近收到的report的channel
	MetricMap map[string]chan *pb.ReportReq
}

var (
	// RegisterMap记录每个已注册的主机
	RegisterMap map[string]*Agent
)
