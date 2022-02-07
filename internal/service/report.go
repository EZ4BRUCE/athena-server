package service

import (
	pb "github.com/EZ4BRUCE/athena-proto/proto"
)

// service层方法，接收请求结构体或特定参数执行dao方法

func (svc *ReportService) CreateReport(r *pb.ReportReq) error {
	return svc.dao.CreateReport(r.GetUId(), r.GetTimestamp(), r.GetMetric(), r.GetDimensions(), r.GetValue())
}
