package service

import pb "github.com/EZ4BRUCE/athena-server/proto"

func (svc *ReportService) MAX(rawReports []*pb.ReportReq) float64 {
	maxValue := 0.0
	for _, r := range rawReports {
		if r.GetValue() > maxValue {
			maxValue = r.GetValue()
		}
	}
	return maxValue
}

func (svc *ReportService) MIN(rawReports []*pb.ReportReq) float64 {
	minValue := rawReports[0].GetValue()
	for _, r := range rawReports {
		if r.GetValue() < minValue {
			minValue = r.GetValue()
		}
	}
	return minValue
}

func (svc *ReportService) SUM(rawReports []*pb.ReportReq) float64 {
	sum := 0.0
	for _, r := range rawReports {
		sum += r.GetValue()
	}
	return sum
}

func (svc *ReportService) AVG(rawReports []*pb.ReportReq) float64 {
	sum := 0.0
	for _, r := range rawReports {
		sum += r.GetValue()
	}
	return sum / float64(len(rawReports))
}
