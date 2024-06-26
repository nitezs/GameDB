package cmd

import (
	"GameDB/internal/crawler"
	"GameDB/internal/db"
	"GameDB/internal/log"
	"GameDB/internal/utils"
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Long:  "Manually add game info",
	Short: "Manually add game info",
	Run:   addRun,
}

type AddCommandConfig struct {
	GameID string `json:"game_id"`
	IDtype string `json:"id_type"`
	ID     int    `json:"id"`
	Config string
}

var addCmdCfg AddCommandConfig

func init() {
	addCmd.Flags().StringVarP(&addCmdCfg.GameID, "game", "g", "", "game id")
	addCmd.Flags().StringVarP(&addCmdCfg.IDtype, "type", "t", "", "id type")
	addCmd.Flags().StringVarP(&addCmdCfg.Config, "config", "c", "", "config path")
	addCmd.Flags().IntVarP(&addCmdCfg.ID, "id", "i", 0, "platform id")
	RootCmd.AddCommand(addCmd)
}

func addRun(cmd *cobra.Command, args []string) {
	c := []*AddCommandConfig{}
	if addCmdCfg.Config != "" {
		data, err := os.ReadFile(addCmdCfg.Config)
		if err != nil {
			log.Logger.Error("Failed to read config file", zap.Error(err))
			return
		}
		if err = json.Unmarshal(data, &c); err != nil {
			log.Logger.Error("Failed to unmarshal config file", zap.Error(err))
			return
		}
	} else {
		c = append(c, &addCmdCfg)
	}
	for _, v := range c {
		objID, err := primitive.ObjectIDFromHex(v.GameID)
		if err != nil {
			log.Logger.Error("Failed to parse game id", zap.Error(err))
			continue
		}
		info, err := db.GetGameInfoByPlatformID(v.IDtype, v.ID)
		if err == nil {
			info.GameIDs = append(info.GameIDs, objID)
			info.GameIDs = utils.Unique(info.GameIDs)
			err = db.SaveGameInfo(info)
			if err != nil {
				log.Logger.Error("Failed to save game info", zap.Error(err))
			}
			log.Logger.Info("Updated game info", zap.String("game_id", v.GameID), zap.String("id_type", v.IDtype), zap.Int("id", v.ID))
			continue
		}
		err = crawler.AddGameInfoManually(objID, v.IDtype, v.ID)
		if err != nil {
			log.Logger.Error("Failed to add game info", zap.Error(err))
		}
		log.Logger.Info("Added game info", zap.String("game_id", v.GameID), zap.String("id_type", v.IDtype), zap.Int("id", v.ID))
	}
}
