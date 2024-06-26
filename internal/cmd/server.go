package cmd

import (
	"GameDB/internal/server"

	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:  "server",
	Long: "Start api server",
	Run:  ServerRun,
}

type serverCommandConfig struct {
	Addr      string
	AutoCrawl bool
}

var serverCmdCfg serverCommandConfig

func init() {
	serverCmd.Flags().StringVarP(&serverCmdCfg.Addr, "addr", "a", ":8080", "server address")
	serverCmd.Flags().BoolVarP(&serverCmdCfg.AutoCrawl, "auto-crawl", "c", true, "enable auto crawl")
	RootCmd.AddCommand(serverCmd)
}

func ServerRun(cmd *cobra.Command, args []string) {
	server.Run(serverCmdCfg.Addr, serverCmdCfg.AutoCrawl)
}
