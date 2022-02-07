package dao

import "github.com/EZ4BRUCE/athena-server/internal/model"

// dao层方法，接收特定参数执行model方法

func (d RuleDao) CreateAggregator(name, metric string, function_id, rule_id uint32) error {
	aggregator := model.Aggregator{Name: name, Metric: metric, FunctionId: function_id, RuleId: rule_id}
	return aggregator.Create(d.engine)
}

func (d RuleDao) DeleteAggregator(id uint32) error {
	aggregator := model.Aggregator{Id: id}
	return aggregator.Delete(d.engine)
}

func (d RuleDao) GetAggregator(id uint32) (model.Aggregator, error) {
	aggregator := model.Aggregator{Id: id}
	return aggregator.Get(d.engine)
}

func (d RuleDao) UpdateAggregator(id uint32, name, metric string, function_id, rule_id uint32) error {
	aggregator := model.Aggregator{Id: id, Name: name, Metric: metric, FunctionId: function_id, RuleId: rule_id}
	return aggregator.Update(d.engine)
}

func (d RuleDao) SearchAggregators(metric string) ([]model.Aggregator, error) {
	aggregator := model.Aggregator{}
	return aggregator.Search(d.engine, metric)
}

func (d RuleDao) ListAggregators() ([]model.Aggregator, error) {
	aggregator := model.Aggregator{}
	return aggregator.List(d.engine)
}
