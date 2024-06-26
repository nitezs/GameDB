package cmd

import (
	"GameDB/internal/crawler"
	"GameDB/internal/db"
	"GameDB/internal/log"
	"strings"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var formatCmd = &cobra.Command{
	Use:   "format",
	Long:  "Format games name",
	Short: "Format games name",
	Run:   FixRun,
}

type FormatCommandConfig struct {
	source string
}

var formatCmdCfg FormatCommandConfig

func init() {
	formatCmd.Flags().StringVarP(&formatCmdCfg.source, "source", "s", "", "source to fix (fitgirl/dodi/kaoskrew/freegog/xatab/onlinefix)")
	RootCmd.AddCommand(formatCmd)
}

func FixRun(cmd *cobra.Command, args []string) {
	formatCmdCfg.source = strings.ToLower(formatCmdCfg.source)
	switch formatCmdCfg.source {
	case "fitgirl":
		items, err := db.GetFitgirlAllGameDownloads()
		if err != nil {
			log.Logger.Error("Failed to get games", zap.Error(err))
			return
		}
		for _, item := range items {
			oldName := item.Name
			item.Name = crawler.FitgirlFormatter(item.RawName)
			if oldName != item.Name {
				log.Logger.Info("Fix name", zap.String("old", oldName), zap.String("raw", item.RawName), zap.String("name", item.Name))
				err := db.SaveGameDownload(item)
				if err != nil {
					log.Logger.Error("Failed to update item", zap.Error(err))
				}
			}
		}
	case "dodi":
		items, err := db.GetDODIAllGameDownloads()
		if err != nil {
			log.Logger.Error("Failed to get games", zap.Error(err))
			return
		}
		for _, item := range items {
			oldName := item.Name
			item.Name = crawler.DODIFormatter(item.RawName)
			if oldName != item.Name {
				log.Logger.Info("Fix name", zap.String("old", oldName), zap.String("raw", item.RawName), zap.String("name", item.Name))
				err := db.SaveGameDownload(item)
				if err != nil {
					log.Logger.Error("Failed to update item", zap.Error(err))
				}
			}
		}
	case "kaoskrew":
		items, err := db.GetKaOsKrewAllGameDownloads()
		if err != nil {
			log.Logger.Error("Failed to get games", zap.Error(err))
			return
		}
		for _, item := range items {
			oldName := item.Name
			item.Name = crawler.KaOsKrewFormatter(item.RawName)
			if oldName != item.Name {
				log.Logger.Info("Fix name", zap.String("old", oldName), zap.String("raw", item.RawName), zap.String("name", item.Name))
				err := db.SaveGameDownload(item)
				if err != nil {
					log.Logger.Error("Failed to update item", zap.Error(err))
				}
			}
		}
	case "freegog":
		items, err := db.GetAllFreeGOGGameDownloads()
		if err != nil {
			log.Logger.Error("Failed to get games", zap.Error(err))
			return
		}
		for _, item := range items {
			oldName := item.Name
			item.Name = crawler.FreeGOGFormatter(item.RawName)
			if oldName != item.Name {
				log.Logger.Info("Fix name", zap.String("old", oldName), zap.String("raw", item.RawName), zap.String("name", item.Name))
				err := db.SaveGameDownload(item)
				if err != nil {
					log.Logger.Error("Failed to update item", zap.Error(err))
				}
			}
		}
	case "xatab":
		items, err := db.GetXatabGameDownloads()
		if err != nil {
			log.Logger.Error("Failed to get games", zap.Error(err))
			return
		}
		for _, item := range items {
			oldName := item.Name
			item.Name = crawler.XatabFormatter(item.RawName)
			if oldName != item.Name {
				log.Logger.Info("Fix name", zap.String("old", oldName), zap.String("raw", item.RawName), zap.String("name", item.Name))
				err := db.SaveGameDownload(item)
				if err != nil {
					log.Logger.Error("Failed to update item", zap.Error(err))
				}
			}
		}
	case "onlinefix":
		items, err := db.GetOnlineFixGameDownloads()
		if err != nil {
			log.Logger.Error("Failed to get games", zap.Error(err))
			return
		}
		for _, item := range items {
			oldName := item.Name
			item.Name = crawler.OnlineFixFormatter(item.RawName)
			if oldName != item.Name {
				log.Logger.Info("Fix name", zap.String("old", oldName), zap.String("raw", item.RawName), zap.String("name", item.Name))
				err := db.SaveGameDownload(item)
				if err != nil {
					log.Logger.Error("Failed to update item", zap.Error(err))
				}
			}
		}
	}
}
