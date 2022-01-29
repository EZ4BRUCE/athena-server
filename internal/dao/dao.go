package dao

import (
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type ReportDao struct {
	engine *mongo.Database
}

type RuleDao struct {
	engine *gorm.DB
}

func NewReportDao(engine *mongo.Database) *ReportDao {
	return &ReportDao{engine: engine}
}

func NewRuleDao(engine *gorm.DB) *RuleDao {
	return &RuleDao{engine: engine}
}
