package task

import (
	"GameDB/internal/crawler"
	"GameDB/internal/db"
	"GameDB/internal/log"
	"GameDB/internal/model"

	"go.uber.org/zap"
)

func Crawl() {
	var err error
	var games []*model.GameDownload
	var g []*model.GameDownload
	g, err = crawler.CrawlFitgirlMulti([]int{1, 2, 3})
	if err == nil {
		games = append(games, g...)
	}
	g, err = crawler.CrawlDODIMulti([]int{1, 2, 3})
	if err == nil {
		games = append(games, g...)
	}
	g, err = crawler.CrawlKaOsKrewMulti([]int{1, 2, 3})
	if err == nil {
		games = append(games, g...)
	}
	g, err = crawler.CrawlXatabMulti([]int{1, 2, 3})
	if err == nil {
		games = append(games, g...)
	}
	g, err = crawler.CrawlFreeGOGAll()
	if err == nil {
		games = append(games, g...)
	}
	g, err = crawler.CrawlOnlineFixMulti([]int{1, 2, 3})
	if err == nil {
		games = append(games, g...)
	}

	for _, game := range games {
		gameInfo, err := crawler.ProcessGameWithIGDB(game)
		if err == nil {
			err = db.SaveGameInfo(gameInfo)
			if err != nil {
				log.Logger.Error("Failed to save game info", zap.Error(err))
			}
			continue
		}
		gameInfo, err = crawler.ProcessGameWithSteam(game)
		if err == nil {
			err = db.SaveGameInfo(gameInfo)
			if err != nil {
				log.Logger.Error("Failed to save game info", zap.Error(err))
			}
			continue
		}
		gameInfo, err = crawler.ProcessGameWithGOG(game)
		if err == nil {
			err = db.SaveGameInfo(gameInfo)
			if err != nil {
				log.Logger.Error("Failed to save game info", zap.Error(err))
			}
			continue
		}
		log.Logger.Error("Failed to process game", zap.Error(err))
	}
}
