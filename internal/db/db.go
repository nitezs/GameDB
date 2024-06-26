package db

import (
	"GameDB/internal/config"
	"GameDB/internal/log"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var MongoDB *mongo.Client
var GameDownloadCollection *mongo.Collection
var LanguageCollection *mongo.Collection
var GameInfoCollection *mongo.Collection

func InitDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI(fmt.Sprintf(
		"mongodb://%s:%s@%s:%v",
		config.Config.Database.User,
		config.Config.Database.Password,
		config.Config.Database.Host,
		config.Config.Database.Port,
	))
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Logger.Panic("Failed to connect to MongoDB", zap.Error(err))
	}
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Logger.Panic("Failed to ping MongoDB", zap.Error(err))
	}
	log.Logger.Info("Connected to MongoDB")
	MongoDB = client

	GameDownloadCollection = MongoDB.Database(config.Config.Database.Database).Collection("game_downloads")
	LanguageCollection = MongoDB.Database(config.Config.Database.Database).Collection("languages")
	GameInfoCollection = MongoDB.Database(config.Config.Database.Database).Collection("game_infos")

	gameDetailsGamesIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "games", Value: 1},
			{Key: "name", Value: 1},
		},
	}
	searchGameDetailsIndex := mongo.IndexModel{
		Keys: bson.D{{Key: "name", Value: "text"}, {Key: "aliases", Value: "text"}},
	}
	_, err = GameDownloadCollection.Indexes().CreateOne(context.TODO(), gameDetailsGamesIndex)
	if err != nil {
		log.Logger.Error("Failed to create index", zap.Error(err))
	}
	_, err = GameInfoCollection.Indexes().CreateOne(context.TODO(), searchGameDetailsIndex)
	if err != nil {
		log.Logger.Error("Failed to create index", zap.Error(err))
	}
}
