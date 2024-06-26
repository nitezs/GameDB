package main

import (
	"GameDB/internal/cache"
	"GameDB/internal/cmd"
	"GameDB/internal/config"
	"GameDB/internal/db"
	"GameDB/internal/log"

	"go.uber.org/zap"
)

func main() {
	config.InitConfig()
	log.InitLogger(config.Config.LogLevel)
	db.InitDB()
	if config.Config.RedisAvaliable {
		cache.InitRedis()
	}
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Logger.Error("main", zap.Error(err))
	}
}
