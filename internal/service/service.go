package service

import (
	"context"

	"github.com/EZ4BRUCE/athena-server/global"
	"github.com/EZ4BRUCE/athena-server/internal/dao"
)

type ReportService struct {
	ctx context.Context
	dao *dao.ReportDao
}

type RuleService struct {
	ctx context.Context
	dao *dao.RuleDao
}

func NewReportService(ctx context.Context) ReportService {
	svc := ReportService{ctx: ctx}
	svc.dao = dao.NewReportDao(global.ReportDBEngine)
	return svc
}

func NewRuleService(ctx context.Context) RuleService {
	svc := RuleService{ctx: ctx}
	svc.dao = dao.NewRuleDao(global.RuleDBEngine)
	return svc
}
