package global

import (
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

var (
	// 告警服务数据库操作实例
	ReportDBEngine *mongo.Database
	// 规则服务数据库操作实例
	RuleDBEngine *gorm.DB
)
