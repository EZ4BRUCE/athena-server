package model

import (
	"context"
	"log"

	"github.com/EZ4BRUCE/athena-server/global"
	"github.com/EZ4BRUCE/athena-server/pkg/setting"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewDBEngine(databaseSetting *setting.DatabaseSettingS) (*mongo.Database, error) {
	var err error
	clientOptions := options.Client().ApplyURI("mongodb://" + global.DatabaseSetting.Host)
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
	return client.Database(global.DatabaseSetting.DBName), nil

}
