package server

import (
	"context"
	"fmt"

	pb "github.com/EZ4BRUCE/athena-server/proto"
)

type ReportServerServer struct{}

func NewReportServer() *ReportServerServer {
	return &ReportServerServer{}
}

func (s *ReportServerServer) Report(ctx context.Context, r *pb.ReportReq) (*pb.ReportRsp, error) {
	fmt.Printf("report: time:%v, metric:%s, dimensions:%v, value:%v \n", r.GetTimestamp(), r.GetMetric(), r.GetDimensions(), r.GetValue())
	return &pb.ReportRsp{Code: 10000001, Msg: "hello"}, nil
}
