package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"nba-remake/internal/config"
)

// NewMongoDBClient 创建一个新的MongoDBClient实例
func NewMongoDBClient(conf *config.MongoDBConfig) *mongo.Client {
	// 设置MongoDB连接选项
	clientOptions := options.Client().ApplyURI(conf.URI)

	// 连接到MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(err)
	}
	return client
}
