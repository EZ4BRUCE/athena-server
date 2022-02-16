package model

import (
	"context"

	"github.com/EZ4BRUCE/athena-server/global"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// 正常事件：保存每一次运行正常的聚合结果到告警数据库
type NormalEvent struct {
	Reports []Report `bson:"reports"`
}

// 异常事件：保存触发了警告的事件
type WarningEvent struct {
	Reports        []Report `bson:"reports"`
	AggregatorId   uint32   `bson:"aggregator_id"`
	AggregatorName string   `bson:"aggregator_name"`
	Metric         string   `bson:"metric"`
	Function       Function `bson:"function"`
	Rule           Rule     `bson:"rule"`
	AggregateValue float64  `bson:"aggregate_value"`
}

// 创建正常事件至告警数据库
func (n NormalEvent) Create(db *mongo.Database) error {
	collection := db.Collection("event")
	result, err := collection.InsertOne(context.TODO(), n)
	if err != nil {
		global.Logger.Errorf("[数据库错误] normalEvent.Create err:%s", err)
		return err
	}
	id := result.InsertedID.(primitive.ObjectID)
	id.Hex()
	return nil
}

// 创建异常事件至告警数据库
func (w WarningEvent) Create(db *mongo.Database) error {
	collection := db.Collection("event")
	result, err := collection.InsertOne(context.TODO(), w)
	if err != nil {
		global.Logger.Errorf("[数据库错误] normalEvent.Create err:%s", err)
		return err
	}
	id := result.InsertedID.(primitive.ObjectID)
	id.Hex()
	return nil
}
