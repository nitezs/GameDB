package cmd

import (
	"GameDB/internal/log"
	"GameDB/internal/task"

	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var cronCmd = &cobra.Command{
	Use:  "cron",
	Long: "Execute scheduled tasks to crawl games",
	Run: func(cmd *cobra.Command, args []string) {
		task.Crawl()
		c := cron.New()
		_, err := c.AddFunc("0 0 * * *", task.Crawl)
		if err != nil {
			log.Logger.Error("Error adding cron job", zap.Error(err))
		}
		c.Start()
		select {}
	},
}

func init() {
	RootCmd.AddCommand(cronCmd)
}
