package db

import (
	"GameDB/internal/model"
)

func GetOnlineFixGameDownloads() ([]*model.GameDownload, error) {
	return GetAllGameDownloadsWithAuthor("onlinefix")
}

func IsOnlineFixCrawled(flag string) bool {
	return IsGameCrawled(flag, "onlinefix")
}
