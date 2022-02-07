package service

import "github.com/EZ4BRUCE/athena-server/internal/model"

// service层的聚合器结构体
type Aggregator struct {
	Id       uint32
	Name     string
	Metric   string
	Function model.Function
	Rule     model.Rule
}

// service层方法，接收请求结构体或特定参数执行dao方法

// 返回特定指标的所有聚合器（包含聚合函数和告警规则）
func (svc *RuleService) SearchAggregators(metric string) ([]Aggregator, error) {
	all, err := svc.dao.SearchAggregators(metric)
	if err != nil {
		return nil, err
	}
	var finalResults []Aggregator
	for _, agg := range all {
		temp, err := svc.combine(&agg)
		if err != nil {
			return nil, err
		}
		finalResults = append(finalResults, temp)
	}
	return finalResults, nil

}

func (svc *RuleService) ListAggregators() ([]Aggregator, error) {
	all, err := svc.dao.ListAggregators()
	if err != nil {
		return nil, err
	}
	var finalResults []Aggregator
	for _, agg := range all {
		temp, err := svc.combine(&agg)
		if err != nil {
			return nil, err
		}
		finalResults = append(finalResults, temp)
	}
	return finalResults, nil
}

// 工具函数，用于把从数据库中取出的聚合器、函数、规则等信息融合起来
func (svc *RuleService) combine(a *model.Aggregator) (Aggregator, error) {
	f, err := svc.dao.GetFunction(a.FunctionId)
	if err != nil {
		return Aggregator{}, err
	}
	r, err := svc.dao.GetRule(a.RuleId)
	if err != nil {
		return Aggregator{}, err
	}
	return Aggregator{Id: a.Id, Name: a.Name, Metric: a.Metric, Function: f, Rule: r}, nil
}
