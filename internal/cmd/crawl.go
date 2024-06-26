package cmd

import (
	"GameDB/internal/crawler"
	"GameDB/internal/log"
	"GameDB/internal/model"
	"GameDB/internal/utils"
	"slices"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var crawlCmd = &cobra.Command{
	Use:  "crawl",
	Long: "Allow you to crawl games from specific platforms",
	Run:  crawlRun,
}

type CrawlCommandConfig struct {
	Source string
	Page   string
	All    bool
	Num    int
}

var crawlCmdCfg CrawlCommandConfig

func init() {
	crawlCmd.Flags().StringVarP(&crawlCmdCfg.Source, "source", "s", "", "source to crawl (fitgirl/dodi/kaoskrew/freegog/xatab/onlinefix)")
	crawlCmd.Flags().StringVarP(&crawlCmdCfg.Page, "pages", "p", "1", "pages to crawl (1,2,3 or 1-3), only available for fitgirl/dodi/kaoskrew/xatab/onlinefix")
	crawlCmd.Flags().BoolVarP(&crawlCmdCfg.All, "all", "a", false, "crawl all page, ignore pages")
	crawlCmd.Flags().IntVarP(&crawlCmdCfg.Num, "num", "n", 1, "number of items to crawl, only available for freegog")
	RootCmd.AddCommand(crawlCmd)
}

func crawlRun(cmd *cobra.Command, args []string) {
	crawlCmdCfg.Source = strings.ToLower(crawlCmdCfg.Source)
	if slices.Contains([]string{"fitgirl", "dodi", "kaoskrew"}, crawlCmdCfg.Source) {
		crawl1337x()
	} else if crawlCmdCfg.Source == "freegog" {
		crawlFreeGOG()
	} else if crawlCmdCfg.Source == "xatab" {
		crawlXatab()
	} else if crawlCmdCfg.Source == "onlinefix" {
		crwalOnlineFix()
	} else {
		log.Logger.Error("Invalid source", zap.String("source", crawlCmdCfg.Source))
	}
}

func pagination(pageStr string) ([]int, error) {
	var pages []int
	pageSlice := strings.Split(pageStr, ",")
	for i := 0; i < len(pageSlice); i++ {
		if strings.Contains(pageSlice[i], "-") {
			pageRange := strings.Split(pageSlice[i], "-")
			start, err := strconv.Atoi(pageRange[0])
			if err != nil {
				return nil, err
			}
			end, err := strconv.Atoi(pageRange[1])
			if err != nil {
				return nil, err
			}
			if start > end {
				return nil, err
			}
			for j := start; j <= end; j++ {
				pages = append(pages, j)
			}
		} else {
			p, err := strconv.Atoi(pageSlice[i])
			if err != nil {
				log.Logger.Error("Invalid page", zap.String("page", pageSlice[i]))
				return nil, err
			}
			pages = append(pages, p)
		}
	}
	return utils.Unique(pages), nil
}

func crawl1337x() {
	var crawl1337xMulti func(pages []int) ([]*model.GameDownload, error)
	var crawl1337xAll func() ([]*model.GameDownload, error)
	switch crawlCmdCfg.Source {
	case "fitgirl":
		crawl1337xMulti = crawler.CrawlFitgirlMulti
		crawl1337xAll = crawler.CrawlFitgirlAll
	case "dodi":
		crawl1337xMulti = crawler.CrawlDODIMulti
		crawl1337xAll = crawler.CrawlDODIAll
	case "kaoskrew":
		crawl1337xMulti = crawler.CrawlKaOsKrewMulti
		crawl1337xAll = crawler.CrawlKaOsKrewAll
	}

	if crawlCmdCfg.All {
		_, err := crawl1337xAll()
		if err != nil {
			return
		}
	} else {
		pages, err := pagination(crawlCmdCfg.Page)
		if err != nil {
			log.Logger.Error("Invalid page", zap.String("page", crawlCmdCfg.Page))
			return
		}
		_, err = crawl1337xMulti(pages)
		if err != nil {
			return
		}
	}
}

func crawlFreeGOG() {
	if crawlCmdCfg.Num <= 0 {
		log.Logger.Error("Invalid num", zap.Int("num", crawlCmdCfg.Num))
		return
	}
	var err error
	if crawlCmdCfg.All {
		_, err = crawler.CrawlFreeGOGAll()
	} else {
		_, err = crawler.CrawlFreeGOG(crawlCmdCfg.Num)
	}
	if err != nil {
		return
	}
}

func crawlXatab() {
	if crawlCmdCfg.All {
		_, err := crawler.CrawlXatabAll()
		if err != nil {
			return
		}
	} else {
		pages, err := pagination(crawlCmdCfg.Page)
		if err != nil {
			log.Logger.Error("Invalid page", zap.String("page", crawlCmdCfg.Page))
			return
		}
		_, err = crawler.CrawlXatabMulti(pages)
		if err != nil {
			return
		}
	}
}

func crwalOnlineFix() {
	if crawlCmdCfg.All {
		_, err := crawler.CrawlOnlineFixAll()
		if err != nil {
			return
		}
	} else {
		pages, err := pagination(crawlCmdCfg.Page)
		if err != nil {
			log.Logger.Error("Invalid page", zap.String("page", crawlCmdCfg.Page))
			return
		}
		_, err = crawler.CrawlOnlineFixMulti(pages)
		if err != nil {
			return
		}
	}
}
