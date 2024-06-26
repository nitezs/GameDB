package server

import (
	"GameDB/internal/server/handler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func initRoute(app *gin.Engine) {
	app.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
	}))

	app.GET("/raw/:id", handler.GetGameDownload)
	app.GET("/game/search", handler.SearchGames)
	app.GET("/game/:id", handler.GetGameInfo)
	app.GET("/game/name/:name", handler.GetGameInfosByName)
	app.GET("/ranking/:type", handler.GetSteam250)
}
