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
	Use:  "update",
	Long: "Update data through regenerated game info",
	Run:  updateRun,
}

type updateCommandConfig struct {
	ID         int
	IDType     string
	GameInfoID string
}

var updateCmdcfx updateCommandConfig

func init() {
	updateCmd.Flags().IntVarP(&updateCmdcfx.ID, "platform-id", "p", 0, "platform id")
	updateCmd.Flags().StringVarP(&updateCmdcfx.IDType, "type", "t", "", "id type")
	updateCmd.Flags().StringVarP(&updateCmdcfx.GameInfoID, "id", "i", "", "game info id")
	RootCmd.AddCommand(updateCmd)
}

func updateRun(cmd *cobra.Command, args []string) {
	id, err := primitive.ObjectIDFromHex(updateCmdcfx.GameInfoID)
	if err != nil {
		log.Logger.Error("Failed to parse game info id", zap.Error(err))
		return
	}
	oldInfo, err := db.GetGameInfoByID(id)
	if err != nil {
		log.Logger.Error("Failed to get game info", zap.Error(err))
		return
	}
	newInfo, err := crawler.GenerateGameInfo(updateCmdcfx.IDType, updateCmdcfx.ID)
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
