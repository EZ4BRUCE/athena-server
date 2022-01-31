package service

import "github.com/EZ4BRUCE/athena-server/internal/model"

// 定义业务逻辑需要的聚合器，聚合函数与告警规则不再以id出现
type Aggregator struct {
	Id       uint32
	Name     string
	Metric   string
	Function model.Function
	Rule     model.Rule
}

// 取出数据库的规则等信息，返回该指标对应的所有聚合器
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
