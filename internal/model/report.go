package model

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Report struct {
	UId        string            `bson:"uId"`
	Timestamp  int64             `bson:"timestamp"`
	Metric     string            `bson:"metric"`
	Dimensions map[string]string `bson:"dimensions"`
	Value      float64           `bson:"value"`
}

func (r Report) Create(db *mongo.Database) error {
	collection := db.Collection("report")
	result, err := collection.InsertOne(context.TODO(), r)
	if err != nil {
		fmt.Print(err)
		return err
	}
	id := result.InsertedID.(primitive.ObjectID)
	fmt.Println("Auto Increasing ID:", id.Hex())
	return nil

}
