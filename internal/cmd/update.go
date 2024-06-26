package cmd

import (
	"GameDB/internal/crawler"
	"GameDB/internal/db"
	"GameDB/internal/log"

	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Long:  "Update Game Info",
	Short: "Update Game Info",
	Run:   updateRun,
}

type updateCommandConfig struct {
	id         int
	idtype     string
	gameInfoID string
}

var updateCmdcfx updateCommandConfig

func init() {
	updateCmd.Flags().IntVarP(&updateCmdcfx.id, "id", "i", 0, "platform id")
	updateCmd.Flags().StringVarP(&updateCmdcfx.idtype, "type", "t", "", "id type")
	updateCmd.Flags().StringVarP(&updateCmdcfx.gameInfoID, "info", "g", "", "game info id")
	RootCmd.AddCommand(updateCmd)
}

func updateRun(cmd *cobra.Command, args []string) {
	id, err := primitive.ObjectIDFromHex(updateCmdcfx.gameInfoID)
	if err != nil {
		log.Logger.Error("Failed to parse game info id", zap.Error(err))
		return
	}
	oldInfo, err := db.GetGameInfoByID(id)
	if err != nil {
		log.Logger.Error("Failed to get game info", zap.Error(err))
		return
	}
	newInfo, err := crawler.GenerateGameInfo(updateCmdcfx.idtype, updateCmdcfx.id)
	if err != nil {
		log.Logger.Error("Failed to generate game info", zap.Error(err))
		return
	}
	newInfo.ID = id
	newInfo.GameIDs = oldInfo.GameIDs
	err = db.SaveGameInfo(newInfo)
	if err != nil {
		log.Logger.Error("Failed to save game info", zap.Error(err))
	}
}
