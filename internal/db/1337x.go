package db

import (
	"GameDB/internal/model"
)

func GetFitgirlAllGameDownloads() ([]*model.GameDownload, error) {
	return GetAllGameDownloadsWithAuthor("fitgirl")
}

func GetDODIAllGameDownloads() ([]*model.GameDownload, error) {
	return GetAllGameDownloadsWithAuthor("dodi")
}

func GetKaOsKrewAllGameDownloads() ([]*model.GameDownload, error) {
	return GetAllGameDownloadsWithAuthor("kaoskrew")
}
