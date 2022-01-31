package server

import (
	"context"
	"fmt"
	"log"

	"github.com/EZ4BRUCE/athena-server/global"
	"github.com/EZ4BRUCE/athena-server/internal/service"
	pb "github.com/EZ4BRUCE/athena-server/proto"
)

type ReportServerServer struct{}

func NewReportServer() *ReportServerServer {
	return &ReportServerServer{}
}

// 每个连接单独处理要在注册的时候用，因为注册只会注册一次
// 防止重复注册
func (s *ReportServerServer) Register(ctx context.Context, r *pb.RegisterReq) (*pb.RegisterRsp, error) {
	return nil, nil
}

func (s *ReportServerServer) Report(ctx context.Context, r *pb.ReportReq) (*pb.ReportRsp, error) {
	_, ok := global.RegisterMap[r.GetUId()]
	if !ok {
		log.Printf("收到未注册Agent上报！time:%v, metric:%s, dimensions:%v, value:%v\n", r.GetTimestamp(), r.GetMetric(), r.GetDimensions(), r.GetValue())
		return &pb.ReportRsp{Code: 10000001, Msg: "Agent未注册！"}, nil
	}
	global.RegisterReports[r.GetUId()] <- r

	repoetSvc := service.NewReportService(ctx)

	err := repoetSvc.CreateReport(r)
	if err != nil {
		return &pb.ReportRsp{Code: 10000001, Msg: "插入数据库失败"}, err
	}
	fmt.Printf("report: time:%v, metric:%s, dimensions:%v, value:%v \n", r.GetTimestamp(), r.GetMetric(), r.GetDimensions(), r.GetValue())
	return &pb.ReportRsp{Code: 10000001, Msg: "hello"}, nil
}
