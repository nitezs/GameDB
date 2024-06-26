package handler

import (
	"GameDB/internal/db"
	"GameDB/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetGameDownloadRequest struct {
	ID string `uri:"id" binding:"required"`
}

type GetGameDownloadResponse struct {
	Status  string              `json:"status"`
	Message string              `json:"message,omitempty"`
	Game    *model.GameDownload `json:"game,omitempty"`
}

func GetGameDownload(c *gin.Context) {
	var req GetGameDownloadRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, GetGameDownloadResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	id, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, GetGameDownloadResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	game, err := db.GetGameDownloadByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, GetGameDownloadResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, GetGameDownloadResponse{
		Status: "ok",
		Game:   game,
	})
}
