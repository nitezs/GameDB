package cmd

import (
	"GameDB/internal/db"
	"GameDB/internal/log"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Long:  "list",
	Short: "list",
	Run:   listRun,
}

type listCommandConfig struct {
	unid bool
}

var listCmdCfg listCommandConfig

func init() {
	listCmd.Flags().BoolVarP(&listCmdCfg.unid, "unid", "u", false, "unid")
	RootCmd.AddCommand(listCmd)
}

func listRun(cmd *cobra.Command, args []string) {
	if listCmdCfg.unid {
		games, err := db.GetGameDownloadsNotInGameInfos(-1)
		if err != nil {
			log.Logger.Error("Failed to get games", zap.Error(err))
		}
		for _, game := range games {
			log.Logger.Info(
				"Game",
				zap.Any("game_id", game.ID),
				zap.String("raw_name", game.RawName),
				zap.String("name", game.Name),
			)
		}
	}
}
