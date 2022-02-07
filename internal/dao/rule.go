package dao

import "github.com/EZ4BRUCE/athena-server/internal/model"

// dao层方法，接收特定参数执行model方法

func (d RuleDao) CreateRule(level, action, description string) error {
	rule := model.Rule{Level: level, Action: action, Description: description}
	return rule.Create(d.engine)
}

func (d RuleDao) DeleteRule(id uint32) error {
	rule := model.Rule{Id: id}
	return rule.Delete(d.engine)
}

func (d RuleDao) GetRule(id uint32) (model.Rule, error) {
	rule := model.Rule{Id: id}
	return rule.Get(d.engine)
}

func (d RuleDao) UpdateRule(id uint32, level, action, description string) error {
	rule := model.Rule{Id: id, Level: level, Action: action, Description: description}
	return rule.Update(d.engine)
}

func (d RuleDao) ListRules() ([]model.Rule, error) {
	rule := model.Rule{}
	return rule.List(d.engine)
}
