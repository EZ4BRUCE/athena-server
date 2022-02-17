package service

import (
	"github.com/EZ4BRUCE/athena-server/internal/model"
)

// service层方法，接收请求结构体或特定参数执行dao方法

func (svc *RuleService) GetAllSmtps() ([]model.Smtp, error) {
	return svc.dao.ListSmtps()
}
