package service

import (
	pb "github.com/EZ4BRUCE/athena-proto/proto"
	"github.com/EZ4BRUCE/athena-server/internal/model"
)

// Service的方法，svc作为接收者
func (svc *ReportService) CreateNormalEvent(rawReports []*pb.ReportReq) error {
	var reports []model.Report
	for _, r := range rawReports {
		temp := model.Report{UId: r.GetUId(), Timestamp: r.GetTimestamp(), Metric: r.GetMetric(), Dimensions: r.GetDimensions(), Value: r.GetValue()}
		reports = append(reports, temp)
	}
	return svc.dao.CreateNormalEvent(reports)
}

// Service的方法，svc作为接收者
func (svc *ReportService) CreateWarningEvent(rawReports []*pb.ReportReq, aggregatorId uint32, aggregatorName string, metric string, function model.Function, rule model.Rule, aggregateValue float64) error {
	var reports []model.Report
	for _, r := range rawReports {
		temp := model.Report{UId: r.GetUId(), Timestamp: r.GetTimestamp(), Metric: r.GetMetric(), Dimensions: r.GetDimensions(), Value: r.GetValue()}
		reports = append(reports, temp)
	}
	return svc.dao.CreateWarningEvent(reports, aggregatorId, aggregatorName, metric, function, rule, aggregateValue)
}
