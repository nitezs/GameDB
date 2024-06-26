package handler

import (
	"GameDB/internal/db"
	"GameDB/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SearchGamesRequest struct {
	Keyword  string `form:"keyword" json:"keyword" binding:"required,min=4,max=64"`
	Page     int    `form:"page" json:"page"`
	PageSize int    `form:"page_size" json:"page_size"`
}

type SearchGamesResponse struct {
	Status    string            `json:"status"`
	Message   string            `json:"message,omitempty"`
	TotalPage int               `json:"total_page,omitempty"`
	GameInfos []*model.GameInfo `json:"game_infos,omitempty"`
}

func SearchGames(c *gin.Context) {
	var req SearchGamesRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, SearchGamesResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	if req.Page == 0 || req.Page < 0 {
		req.Page = 1
	}
	if req.PageSize == 0 || req.PageSize < 0 {
		req.PageSize = 10
	}
	if req.PageSize > 10 {
		req.PageSize = 10
	}
	items, totalPage, err := db.SearchGameInfos(req.Keyword, req.Page, req.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, SearchGamesResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	if len(items) == 0 {
		c.JSON(http.StatusOK, SearchGamesResponse{
			Status:  "ok",
			Message: "No results found",
		})
		return
	}
	c.JSON(http.StatusOK, SearchGamesResponse{
		Status:    "ok",
		TotalPage: totalPage,
		GameInfos: items,
	})
}
