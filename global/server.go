package global

import pb "github.com/EZ4BRUCE/athena-proto/proto"

var (
	// 记录当前已注册的主机id
	RegisterMap map[string]struct{}
	// 记录每个已注册的主机的最近上报信息（用于聚合数据）
	ReportMap map[string]chan *pb.ReportReq
)
