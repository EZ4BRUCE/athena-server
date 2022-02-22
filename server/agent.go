package server

import (
	"sync"

	pb "github.com/EZ4BRUCE/athena-proto/proto"
)

// 根据用户注册提供的参数定义每个agent主机的结构体
type Agent struct {
	UId              string                        // 用户注册后发放的标识UId
	CheckAliveTime   int32                         // 连接诊断时间
	AggregationTime  int32                         // 聚合时间(/单位聚合粒度)
	CheckAliveStatus bool                          // 当前连接存活状态
	IsDead           bool                          // 标识该Agent是否已断开
	MetricMap        map[string]chan *pb.ReportReq // 最近上报信息(用于聚合数据)
	Description      string                        //主机描述
}

var (
	RegisterMap sync.Map // RegisterMap记录每个已注册的主机
)
