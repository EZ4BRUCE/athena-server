package service

import (
	pb "github.com/EZ4BRUCE/athena-server/proto"
)

// Service的方法，svc作为接收者
func (svc *ReportService) CreateReport(r *pb.ReportReq) error {
	return svc.dao.CreateReport(r.GetTimestamp(), r.GetMetric(), r.GetDimensions(), r.GetValue())
}
