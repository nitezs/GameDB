package handler

import (
	"GameDB/internal/crawler"
	"GameDB/internal/log"
	"GameDB/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type GetSteam250Response struct {
	Status  string            `json:"status"`
	Message string            `json:"message,omitempty"`
	Games   []*model.GameInfo `json:"games,omitempty"`
}

func GetSteam250(c *gin.Context) {
	rankingType, exist := c.Params.Get("type")
	if !exist {
		c.JSON(http.StatusBadRequest, GetSteam250Response{
			Status:  "error",
			Message: "Missing ranking type",
		})
	}
	var f func() ([]model.Steam250Item, error)
	switch rankingType {
	case "top":
		f = crawler.GetSteam250Top250Cache
	case "week-top":
		f = crawler.GetSteam250WeekTop50Cache
	case "best-of-the-year":
		f = crawler.GetSteam250BestOfTheYearCache
	case "most-played":
		f = crawler.GetSteam250MostPlayedCache
	default:
		c.JSON(http.StatusBadRequest, GetSteam250Response{
			Status:  "error",
			Message: "Invalid ranking type",
		})
		return
	}
	m, err := f()
	if err != nil {
		c.JSON(http.StatusInternalServerError, GetSteam250Response{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	var infos []*model.GameInfo
	for _, item := range m {
		info, err := crawler.GenerateSteamGameInfo(item.SteamID)
		if err != nil {
			log.Logger.Warn("Failed to generate game info", zap.Error(err))
			continue
		}
		infos = append(infos, info)
	}

	c.JSON(http.StatusOK, GetSteam250Response{
		Status: "ok",
		Games:  infos,
	})
}
