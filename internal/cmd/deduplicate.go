package cmd

import (
	"GameDB/internal/db"
	"GameDB/internal/log"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var deduplicateCmd = &cobra.Command{
	Use:  "deduplicate",
	Long: "Remove duplicate games caused by incorrect crawling",
	Run: func(cmd *cobra.Command, args []string) {
		err := db.DeduplicateGames()
		if err != nil {
			log.Logger.Error("Failed to deduplicate games", zap.Error(err))
		}
	},
}

func init() {
	RootCmd.AddCommand(deduplicateCmd)
}
