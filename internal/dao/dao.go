package dao

import (
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

// Dao封装一个特定db引擎实例用于执行model方法

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
