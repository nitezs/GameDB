package cmd

import (
	"GameDB/internal/crawler"
	"GameDB/internal/db"
	"GameDB/internal/log"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var organizeCmd = &cobra.Command{
	Use:  "organize",
	Long: "Organize game info by repack game's name",
	Run:  organizeRun,
}

type organizeCommandConfig struct {
	Num int
}

var organizeCmdCfg organizeCommandConfig

func init() {
	organizeCmd.Flags().IntVarP(&organizeCmdCfg.Num, "num", "n", 1, "number of items to process")
	RootCmd.AddCommand(organizeCmd)
}

func organizeRun(cmd *cobra.Command, args []string) {
	games, err := db.GetGameDownloadsNotInGameInfos(organizeCmdCfg.Num)
	if err != nil {
		log.Logger.Error("Failed to get games", zap.Error(err))
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
