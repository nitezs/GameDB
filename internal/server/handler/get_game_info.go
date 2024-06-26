package handler

import (
	"GameDB/internal/db"
	"GameDB/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetGameInfoRequest struct {
	ID string `uri:"id" binding:"required"`
}

type GetGameInfoResponse struct {
	Status   string          `json:"status"`
	Message  string          `json:"message,omitempty"`
	GameInfo *model.GameInfo `json:"game_info,omitempty"`
}

func GetGameInfo(c *gin.Context) {
	var req GetGameDownloadRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, GetGameInfoResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	id, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, GetGameInfoResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	gameInfo, err := db.GetGameInfoByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, GetGameInfoResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	gameInfo.Games, err = db.GetGameDownloadsByIDs(gameInfo.GameIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, GetGameInfoResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, GetGameInfoResponse{
		Status:   "ok",
		GameInfo: gameInfo,
	})
}
