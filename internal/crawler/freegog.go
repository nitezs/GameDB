package crawler

import (
	"GameDB/internal/constant"
	"GameDB/internal/db"
	"GameDB/internal/log"
	"GameDB/internal/model"
	"GameDB/internal/utils"
	"bytes"
	"encoding/base64"
	"html"
	"math"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
)

func CrawlFreeGOG(num int) ([]*model.GameDownload, error) {
	haveCrawled := 0
	resp, err := utils.Fetch(utils.FetchConfig{
		Url: constant.FreeGOGListURL,
	})
	if err != nil {
		log.Logger.Error("Failed to fetch", zap.Error(err))
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Data))
	if err != nil {
		log.Logger.Error("Failed to parse HTML", zap.Error(err))
		return nil, err
	}

	urls := []string{}
	updateFlags := []string{} //rawName+link
	doc.Find(".items-outer li a").Each(func(i int, s *goquery.Selection) {
		urls = append(urls, s.AttrOr("href", ""))
		updateFlags = append(updateFlags, s.Text()+s.AttrOr("href", ""))
	})

	res := []*model.GameDownload{}
	for i := 0; i < len(urls); i++ {
		if haveCrawled >= num {
			break
		}
		if db.IsFreeGOGCrawled(updateFlags[i]) {
			log.Logger.Info("Skipping already crawled item", zap.String("URL", urls[i]))
			continue
		}
		log.Logger.Info("Crawling item", zap.String("URL", urls[i]))
		resp, err := utils.Fetch(utils.FetchConfig{
			Url: urls[i],
		})
		if err != nil {
			log.Logger.Warn("Failed to fetch", zap.Error(err), zap.String("URL", urls[i]))
			continue
		}
		item, err := db.GetGameDownloadByUrl(urls[i])
		if err != nil {
			log.Logger.Error("Failed to get game item", zap.Error(err))
			return nil, err
		}
		item.Url = urls[i]
		item.UpdateFlag = updateFlags[i]
		rawTitleRegex := regexp.MustCompile(`(?i)<h1 class="entry-title">(.*?)</h1>`)
		rawTitleRegexRes := rawTitleRegex.FindStringSubmatch(string(resp.Data))
		if len(rawTitleRegexRes) > 1 {
			rawName := rawTitleRegexRes[1]
			rawName = html.UnescapeString(rawName)
			rawName = strings.Replace(rawName, "–", "-", -1)
			item.RawName = rawName
		} else {
			log.Logger.Warn("Failed to get title", zap.String("url", urls[i]))
			continue
		}
		item.Name = FreeGOGFormatter(item.RawName)
		sizeRegex := regexp.MustCompile(`(?i)>Size:\s?(.*?)<`)
		sizeRegexRes := sizeRegex.FindStringSubmatch(string(resp.Data))
		if len(sizeRegexRes) > 1 {
			item.Size = sizeRegexRes[1]
		}
		magnetRegex := regexp.MustCompile(`<a class="download-btn" href="https://gdl.freegogpcgames.xyz/download-gen\.php\?url=(.*?)"`)
		magnetRegexRes := magnetRegex.FindStringSubmatch(string(resp.Data))
		if len(magnetRegexRes) > 1 {
			magnet, err := base64.StdEncoding.DecodeString(magnetRegexRes[1])
			if err != nil {
				log.Logger.Warn("Failed to decode magnet", zap.String("url", urls[i]))
				continue
			}
			item.Magnet = string(magnet)
		} else {
			log.Logger.Warn("Failed to get magnet", zap.String("url", urls[i]))
			continue
		}
		item.Author = "FreeGOG"
		res = append(res, item)
		haveCrawled++
		err = db.SaveGameDownload(item)
		if err != nil {
			log.Logger.Error("Failed to save game item", zap.Error(err))
		}
	}
	return res, nil
}

func CrawlFreeGOGAll() ([]*model.GameDownload, error) {
	return CrawlFreeGOG(math.MaxInt)
}

var freeGOGRegexps = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\(.*\)`),
}

func FreeGOGFormatter(name string) string {
	for _, re := range freeGOGRegexps {
		name = re.ReplaceAllString(name, "")
	}

	reg1 := regexp.MustCompile(`(?i)v\d+(\.\d+)*`)
	if index := reg1.FindIndex([]byte(name)); index != nil {
		name = name[:index[0]]
	}
	if index := strings.Index(name, "+"); index != -1 {
		name = name[:index]
	}

	reg2 := regexp.MustCompile(`(?i):\sgoty`)
	name = reg2.ReplaceAllString(name, ": Game Of The Year")

	return strings.TrimSpace(name)
}
