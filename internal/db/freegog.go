package db

import (
	"GameDB/internal/model"
)

func GetAllFreeGOGGameDownloads() ([]*model.GameDownload, error) {
	return GetAllGameDownloadsWithAuthor("freegog")
}
func IsFreeGOGCrawled(flag string) bool {
	return IsGameCrawled(flag, "freegog")
}
