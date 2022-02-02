package service

import (
	"log"

	pb "github.com/EZ4BRUCE/athena-proto/proto"
	"github.com/EZ4BRUCE/athena-server/internal/model"
)

func (svc *RuleService) ExecuteFunc(function model.Function, rawReports []*pb.ReportReq) (float64, bool) {
	switch function.Type {
	case "MAX":
		return MAX(rawReports), MAX(rawReports) >= function.Threshold
	case "MIN":
		// 场景是检测他有没有满血干活,聚合数据的最小值不可小于阈值
		return MIN(rawReports), MIN(rawReports) < function.Threshold
	case "AVG":
		return AVG(rawReports), AVG(rawReports) >= function.Threshold
	case "SUM":
		return SUM(rawReports), SUM(rawReports) >= function.Threshold
	default:
		log.Printf("Function：找不到 %s 类型的聚合函数，返回0", function.Type)
		return 0.0, false
	}
}

func MAX(rawReports []*pb.ReportReq) float64 {
	maxValue := 0.0
	for _, r := range rawReports {
		if r.GetValue() > maxValue {
			maxValue = r.GetValue()
		}
	}
	return maxValue
}

func MIN(rawReports []*pb.ReportReq) float64 {
	minValue := rawReports[0].GetValue()
	for _, r := range rawReports {
		if r.GetValue() < minValue {
			minValue = r.GetValue()
		}
	}
	return minValue
}

func SUM(rawReports []*pb.ReportReq) float64 {
	sum := 0.0
	for _, r := range rawReports {
		sum += r.GetValue()
	}
	return sum
}

func AVG(rawReports []*pb.ReportReq) float64 {
	sum := 0.0
	for _, r := range rawReports {
		sum += r.GetValue()
	}
	return sum / float64(len(rawReports))
}
