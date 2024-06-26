package crawler

import (
	"GameDB/internal/constant"
	"GameDB/internal/db"
	"GameDB/internal/log"
	"GameDB/internal/model"
	"GameDB/internal/utils"
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
)

func CrawlXatab(page int) ([]*model.GameDownload, error) {
	requestURL := fmt.Sprintf("%s/page/%v", constant.XatabBaseURL, page)
	log.Logger.Info("Crawling item", zap.String("url", requestURL))
	resp, err := utils.Fetch(utils.FetchConfig{
		Url: requestURL,
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
	updateFlags := []string{} //link+date
	doc.Find(".entry").Each(func(i int, s *goquery.Selection) {
		urls = append(urls, s.Find(".entry__title.h2 a").AttrOr("href", ""))
		updateFlags = append(
			updateFlags,
			s.Find(".entry__title.h2 a").AttrOr("href", "")+
				s.Find(".entry__info-categories").Text(),
		)
	})
	var res []*model.GameDownload
	for i := 0; i < len(urls); i++ {
		if db.IsXatabCrawled(updateFlags[i]) {
			log.Logger.Info("Skipping already crawled item", zap.String("URL", urls[i]))
			continue
		}
		log.Logger.Info("Crawling item", zap.String("URL", urls[i]))
		resp, err := utils.Fetch(utils.FetchConfig{
			Url: urls[i],
		})
		if err != nil {
			log.Logger.Error("Failed to fetch", zap.Error(err))
			return nil, err
		}
		doc, err = goquery.NewDocumentFromReader(bytes.NewReader(resp.Data))
		if err != nil {
			log.Logger.Error("Failed to parse HTML", zap.Error(err))
			return nil, err
		}
		item, err := db.GetGameDownloadByUrl(urls[i])
		if err != nil {
			log.Logger.Error("Failed to get game item", zap.Error(err))
			return nil, err
		}
		item.Url = urls[i]
		item.UpdateFlag = updateFlags[i]
		item.RawName = doc.Find(".inner-entry__title").First().Text()
		item.Name = XatabFormatter(item.RawName)
		item.Author = "Xatab"
		downloadURL := doc.Find("#download>a").First().AttrOr("href", "")
		if downloadURL == "" {
			log.Logger.Error("Failed to find download URL", zap.String("item", item.Name))
			continue
		}
		resp, err = utils.Fetch(utils.FetchConfig{
			Url: downloadURL,
		})
		if err != nil {
			log.Logger.Error("Failed to fetch", zap.Error(err))
			continue
		}
		magnet, size, err := utils.ConvertTorrentToMagnet(resp.Data)
		if err != nil {
			log.Logger.Error("Failed to convert torrent to magnet", zap.Error(err))
			continue
		}
		item.Size = size
		item.Magnet = magnet
		res = append(res, item)
		err = db.SaveGameDownload(item)
		if err != nil {
			log.Logger.Error("Failed to save game item", zap.Error(err))
		}
	}
	return res, nil
}

func CrawlXatabMulti(pages []int) ([]*model.GameDownload, error) {
	totalPageNum, err := GetXatabTotalPageNum()
	if err != nil {
		return nil, err
	}
	var res []*model.GameDownload
	for _, page := range pages {
		if page > totalPageNum {
			log.Logger.Warn("Current page exceed total page", zap.Int("page", page))
			continue
		}
		items, err := CrawlXatab(page)
		if err != nil {
			return nil, err
		}
		res = append(res, items...)
	}
	return res, nil
}

func CrawlXatabAll() ([]*model.GameDownload, error) {
	totalPageNum, err := GetXatabTotalPageNum()
	if err != nil {
		return nil, err
	}
	var res []*model.GameDownload
	for i := 1; i <= totalPageNum; i++ {
		items, err := CrawlXatab(i)
		if err != nil {
			return nil, err
		}
		res = append(res, items...)
	}
	return res, nil
}

func GetXatabTotalPageNum() (int, error) {
	resp, err := utils.Fetch(utils.FetchConfig{
		Url: constant.XatabBaseURL,
	})
	if err != nil {
		log.Logger.Error("Failed to fetch", zap.Error(err))
		return 0, err
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Data))
	if err != nil {
		log.Logger.Error("Failed to parse HTML", zap.Error(err))
		return 0, err
	}
	pageStr := doc.Find(".pagination>a").Last().Text()
	totalPageNum, err := strconv.Atoi(pageStr)
	if err != nil {
		log.Logger.Error("Failed to parse total page num", zap.Error(err))
		return 0, err
	}
	return totalPageNum, nil
}

var xatabRegexps = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\sPC$`),
}

func XatabFormatter(name string) string {
	reg1 := regexp.MustCompile(`(?i)v(er)?\s?(\.)?\d+(\.\d+)*`)
	if index := reg1.FindIndex([]byte(name)); index != nil {
		name = name[:index[0]]
	}
	if index := strings.Index(name, "["); index != -1 {
		name = name[:index]
	}
	if index := strings.Index(name, "("); index != -1 {
		name = name[:index]
	}
	if index := strings.Index(name, "{"); index != -1 {
		name = name[:index]
	}
	if index := strings.Index(name, "+"); index != -1 {
		name = name[:index]
	}
	name = strings.TrimSpace(name)
	for _, re := range xatabRegexps {
		name = re.ReplaceAllString(name, "")
	}

	if index := strings.Index(name, "/"); index != -1 {
		names := strings.Split(name, "/")
		longestLength := 0
		longestName := ""
		for _, n := range names {
			if !utils.ContainsRussian(n) && len(n) > longestLength {
				longestLength = len(n)
				longestName = n
			}
		}
		name = longestName
	}

	return strings.TrimSpace(name)
}
