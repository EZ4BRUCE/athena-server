package model

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// 保存每一次运行正常的聚合结果到数据库
type NormalEvent struct {
	Reports []Report `bson:"reports"`
}

// 每一个已经触发了的警告事件还会绑定一个对应该警告的聚合器Values
type WarningEvent struct {
	Reports        []Report `bson:"reports"`
	AggregatorId   uint32   `bson:"aggregator_id"`
	AggregatorName string   `bson:"aggregator_name"`
	Metric         string   `bson:"metric"`
	Function       Function `bson:"function"`
	Rule           Rule     `bson:"rule"`
	AggregateValue float64  `bson:"aggregate_value"`
}

func (n NormalEvent) Create(db *mongo.Database) error {
	collection := db.Collection("event")
	result, err := collection.InsertOne(context.TODO(), n)
	if err != nil {
		log.Printf("normalEvent.Create err:%s", err)
		return err
	}
	id := result.InsertedID.(primitive.ObjectID)
	id.Hex()
	return nil

}

func (w WarningEvent) Create(db *mongo.Database) error {
	collection := db.Collection("event")
	result, err := collection.InsertOne(context.TODO(), w)
	if err != nil {
		log.Printf("normalEvent.Create err:%s", err)
		return err
	}
	id := result.InsertedID.(primitive.ObjectID)
	id.Hex()
	return nil

}
