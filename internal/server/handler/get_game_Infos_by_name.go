package handler

import (
	"GameDB/internal/db"
	"GameDB/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetGameInfosByNameRequest struct {
	Name string `uri:"name" binding:"required"`
}

type GetGameInfosByNameResponse struct {
	Status    string            `json:"status"`
	Message   string            `json:"message,omitempty"`
	GameInfos []*model.GameInfo `json:"game_infos,omitempty"`
}

func GetGameInfosByName(c *gin.Context) {
	var req GetGameInfosByNameRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, GetGameInfosByNameResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	games, err := db.GetGameInfosByName(req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, GetGameInfosByNameResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, GetGameInfosByNameResponse{
		Status:    "ok",
		GameInfos: games,
	})
}
