package service

import (
	pb "github.com/EZ4BRUCE/athena-proto/proto"
	"github.com/EZ4BRUCE/athena-server/global"
	"github.com/EZ4BRUCE/athena-server/internal/model"
)

// service层方法，接收请求结构体或特定参数执行dao方法

// 执行聚合函数，返回聚合结果以及是否异常：异常为true，正常为false
func (svc *RuleService) ExecuteFunc(function model.Function, rawReports []*pb.ReportReq) (float64, bool) {
	switch function.Type {
	case global.MAXFunction:
		return MAX(rawReports), MAX(rawReports) >= function.Threshold
	case global.MINFunction:
		// 场景是检测他有没有满血干活,聚合数据的最小值不可小于阈值
		return MIN(rawReports), MIN(rawReports) < function.Threshold
	case global.AVGFunction:
		return AVG(rawReports), AVG(rawReports) >= function.Threshold
	case global.SUMFunction:
		return SUM(rawReports), SUM(rawReports) >= function.Threshold
	case global.LOGICFunction:
		return LOGIC(rawReports), LOGIC(rawReports) >= function.Threshold
	default:
		global.Logger.Errorf("[参数错误] Function：找不到 %s 类型的聚合函数，返回0", function.Type)
		return 0.0, false
	}
}

// 定义每种函数类型的处理逻辑

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

func LOGIC(rawReports []*pb.ReportReq) float64 {
	sum := 0.0
	for _, r := range rawReports {
		sum += r.GetValue()
	}
	return sum / float64(len(rawReports))
}
