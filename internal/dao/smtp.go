package dao

import "github.com/EZ4BRUCE/athena-server/internal/model"

// dao层方法，接收特定参数执行model方法

func (rd RuleDao) CreateSmtp(host string, port int, isSSL bool, userName, password, from string) error {
	rule := model.Smtp{
		Host:     host,
		Port:     port,
		IsSSL:    isSSL,
		UserName: userName,
		Password: password,
		From:     from,
	}
	return rule.Create(rd.engine)
}

func (rd RuleDao) DeleteSmtp(id uint32) error {
	rule := model.Smtp{Id: id}
	return rule.Delete(rd.engine)
}

func (rd RuleDao) GetSmtp(id uint32) (model.Smtp, error) {
	rule := model.Smtp{Id: id}
	return rule.Get(rd.engine)
}

func (rd RuleDao) UpdateSmtp(id uint32, host string, port int, isSSL bool, userName, password, from string) error {
	rule := model.Smtp{
		Id:       id,
		Host:     host,
		Port:     port,
		IsSSL:    isSSL,
		UserName: userName,
		Password: password,
		From:     from,
	}
	return rule.Update(rd.engine)
}

func (rd RuleDao) ListSmtps() ([]model.Smtp, error) {
	rule := model.Smtp{}
	return rule.List(rd.engine)
}
