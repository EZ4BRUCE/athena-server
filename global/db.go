package global

import (
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

var (
	ReportDBEngine *mongo.Database
	RuleDBEngine   *gorm.DB
)
