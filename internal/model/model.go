package model

import (
	"context"
	"fmt"
	"log"

	"github.com/EZ4BRUCE/athena-server/global"
	"github.com/EZ4BRUCE/athena-server/pkg/setting"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func NewReportDBEngine(reportDBSetting *setting.ReportDBSettingS) (*mongo.Database, error) {
	var err error
	clientOptions := options.Client().ApplyURI("mongodb://" + global.ReportDBSetting.Host)
	// 连接到MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	// 检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return client.Database(global.ReportDBSetting.DBName), nil

}

func NewRuleDBEngine(ruleDBSetting *setting.RuleDBSettingS) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local",
		ruleDBSetting.UserName,
		ruleDBSetting.Password,
		ruleDBSetting.Host,
		ruleDBSetting.DBName,
		ruleDBSetting.Charset,
		ruleDBSetting.ParseTime,
	)

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
		Logger:         logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}

	db.Logger.LogMode(logger.Info)
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(ruleDBSetting.MaxIdleConns)
	sqlDB.SetMaxOpenConns(ruleDBSetting.MaxOpenConns)
	return db, nil

}
