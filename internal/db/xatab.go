package db

import (
	"GameDB/internal/model"
)

func GetXatabGameDownloads() ([]*model.GameDownload, error) {
	return GetAllGameDownloadsWithAuthor("xatab")
}

func IsXatabCrawled(flag string) bool {
	return IsGameCrawled(flag, "xatab")
}
