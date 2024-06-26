package server

import (
	"GameDB/internal/config"
	"GameDB/internal/log"
	"GameDB/internal/server/middleware"
	"GameDB/internal/task"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

func Run(addr string, autoCrawl bool) {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard
	app := gin.New()
	app.Use(middleware.Logger())
	app.Use(middleware.Recovery())
	initRoute(app)
	log.Logger.Info("Server running", zap.String("addr", addr))
	if config.Config.AutoCrawl || autoCrawl {
		go func() {
			c := cron.New()
			_, err := c.AddFunc("0 0 * * *", task.Crawl)
			if err != nil {
				log.Logger.Error("Error adding cron job", zap.Error(err))
			}
			c.Start()
		}()
	}
	err := app.Run(addr)
	if err != nil {
		log.Logger.Panic("Failed to run server", zap.Error(err))
	}
}
