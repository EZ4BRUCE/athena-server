package global

import (
	pb "github.com/EZ4BRUCE/athena-proto/proto"
)

type Agent struct {
	UId       string
	MetricMap map[string]chan *pb.ReportReq
}

var (
	// 记录每个已注册的主机的最近上报信息（用于聚合数据）
	RegisterMap map[string]Agent
)
