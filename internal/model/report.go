package model

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// 监控信息Report对应的告警数据库字段信息
type Report struct {
	UId        string            `bson:"uId"`
	Timestamp  int64             `bson:"timestamp"`
	Metric     string            `bson:"metric"`
	Dimensions map[string]string `bson:"dimensions"`
	Value      float64           `bson:"value"`
}

// 创建监控信息至告警数据库
func (r Report) Create(db *mongo.Database) error {
	collection := db.Collection("report")
	result, err := collection.InsertOne(context.TODO(), r)
	if err != nil {
		return err
	}
	id := result.InsertedID.(primitive.ObjectID)
	id.Hex()
	return nil
}
