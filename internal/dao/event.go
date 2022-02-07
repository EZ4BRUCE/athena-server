package dao

import "github.com/EZ4BRUCE/athena-server/internal/model"

// dao层方法，接收特定参数执行model方法

func (d ReportDao) CreateNormalEvent(reports []model.Report) error {
	normalEvent := model.NormalEvent{Reports: reports}
	return normalEvent.Create(d.engine)
}

func (d ReportDao) CreateWarningEvent(reports []model.Report, aggregatorId uint32, aggregatorName, metric string, function model.Function, rule model.Rule, aggregateValue float64) error {
	warningEvent := model.WarningEvent{Reports: reports, AggregatorId: aggregatorId, AggregatorName: aggregatorName, Metric: metric, Function: function, Rule: rule, AggregateValue: aggregateValue}
	return warningEvent.Create(d.engine)
}
