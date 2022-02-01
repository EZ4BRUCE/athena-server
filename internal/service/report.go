package service

import (
	pb "github.com/EZ4BRUCE/athena-proto/proto"
)

func (svc *ReportService) ExecuteWarnings() error {
	return nil
}

// Service的方法，svc作为接收者
func (svc *ReportService) CreateReport(r *pb.ReportReq) error {
	return svc.dao.CreateReport(r.GetUId(), r.GetTimestamp(), r.GetMetric(), r.GetDimensions(), r.GetValue())
}
