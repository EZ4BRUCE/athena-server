package dao

import "github.com/EZ4BRUCE/athena-server/internal/model"

// dao层方法，接收特定参数执行model方法

func (rd RuleDao) CreateEmail(address string) error {
	email := model.Email{Address: address}
	return email.Create(rd.engine)
}

func (rd RuleDao) DeleteEmail(id uint32) error {
	email := model.Email{Id: id}
	return email.Delete(rd.engine)
}

func (rd RuleDao) GetEmail(id uint32) (model.Email, error) {
	email := model.Email{Id: id}
	return email.Get(rd.engine)
}

func (rd RuleDao) UpdateEmail(id uint32, address string) error {
	email := model.Email{Id: id, Address: address}
	return email.Update(rd.engine)
}

func (rd RuleDao) ListEmails() ([]model.Email, error) {
	email := model.Email{}
	return email.List(rd.engine)
}
