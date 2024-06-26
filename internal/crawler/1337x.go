package crawler

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"GameDB/internal/constant"
	"GameDB/internal/db"
	"GameDB/internal/log"
	"GameDB/internal/model"
	"GameDB/internal/utils"

	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
)

type Formatter func(string) string

func Crawl1337x(source string, page int, formatter Formatter) ([]*model.GameDownload, error) {
	var resp *utils.FetchResponse
	var doc *goquery.Document
	var err error
	requestUrl := fmt.Sprintf("%s/%s/%d/", constant.C1337xBaseURL, source, page)
	log.Logger.Info("Crawling item", zap.String("url", requestUrl))
	resp, err = utils.Fetch(utils.FetchConfig{
		Url: requestUrl,
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
	trSelection := doc.Find("tbody>tr")
	urls := []string{}
	trSelection.Each(func(i int, trNode *goquery.Selection) {
		nameSelection := trNode.Find(".name").First()
		if aNode := nameSelection.Find("a").Eq(1); aNode.Length() > 0 {
			url, _ := aNode.Attr("href")
			urls = append(urls, url)
		}
	})
	var res []*model.GameDownload
	for i := 0; i < len(urls); i++ {
		urls[i] = fmt.Sprintf("%s%s", constant.C1337xBaseURL, urls[i])
		if db.IsGameCrawledByURL(urls[i]) {
			log.Logger.Info("Skipping already crawled item", zap.String("url", urls[i]))
			continue
		}
		requestUrl = urls[i]
		log.Logger.Info("Crawling item", zap.String("url", requestUrl))
		resp, err = utils.Fetch(utils.FetchConfig{
			Url: requestUrl,
		})
		if err != nil {
			log.Logger.Warn("Failed to fetch", zap.String("url", requestUrl), zap.Error(err))
			continue
		}
		var game = &model.GameDownload{}
		game.Url = urls[i]
		doc, err = goquery.NewDocumentFromReader(bytes.NewReader(resp.Data))
		if err != nil {
			return nil, err
		}
		selection := doc.Find(".torrent-detail-page ul.list>li")
		info := make(map[string]string)
		selection.Each(func(i int, item *goquery.Selection) {
			info[strings.TrimSpace(item.Find("strong").Text())] = strings.TrimSpace(item.Find("span").Text())
		})
		magnetRegex := regexp.MustCompile(`magnet:\?[^"]*`)
		magnetRegexRes := magnetRegex.FindStringSubmatch(string(resp.Data))
		game.Size = info["Total size"]
		game.RawName = doc.Find("title").Text()
		game.RawName = strings.Replace(game.RawName, "Download ", "", 1)
		game.RawName = strings.TrimSpace(strings.Replace(game.RawName, "Torrent | 1337x", " ", 1))
		game.Name = formatter(game.RawName)
		game.Magnet = magnetRegexRes[0]
		game.Author = strings.Replace(source, "-torrents", "", -1)
		err = db.SaveGameDownload(game)
		if err != nil {
			log.Logger.Warn("Failed to save", zap.Error(err))
			continue
		}
		res = append(res, game)
	}
	if err != nil {
		log.Logger.Error("Failed to add game item", zap.Error(err))
	}
	return res, nil
}

func Crawl1337xMulti(source string, pages []int, formatter Formatter) (res []*model.GameDownload, err error) {
	var items []*model.GameDownload
	totalPageNum, err := Get1337xTotalPageNum(source)
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		if page > totalPageNum {
			log.Logger.Warn("Current page exceed total page", zap.Int("page", page))
			continue
		}
		items, err = Crawl1337x(source, page, formatter)
		res = append(res, items...)
		if err != nil {
			return nil, err
		}
	}
	log.Logger.Info("Crawled finished", zap.Int("num", len(res)))
	return res, nil
}

func Crawl1337xAll(source string, formatter Formatter) (res []*model.GameDownload, err error) {
	totalPageNum, err := Get1337xTotalPageNum(source)
	if err != nil {
		return nil, err
	}
	var items []*model.GameDownload
	for i := 1; i <= totalPageNum; i++ {
		items, err = Crawl1337x(source, i, formatter)
		res = append(res, items...)
		if err != nil {
			return nil, err
		}
	}
	log.Logger.Info("Crawled finished", zap.Int("num", len(res)))
	return res, nil
}

func Get1337xTotalPageNum(source string) (int, error) {
	var resp *utils.FetchResponse
	var doc *goquery.Document
	var err error

	requestUrl := fmt.Sprintf("%s/%s/%d/", constant.C1337xBaseURL, source, 1)
	resp, err = utils.Fetch(utils.FetchConfig{
		Url: requestUrl,
	})
	if err != nil {
		return 0, err
	}
	doc, _ = goquery.NewDocumentFromReader(bytes.NewReader(resp.Data))
	selection := doc.Find(".last")
	pageStr, exist := selection.Find("a").Attr("href")
	if !exist {
		return 0, errors.New("total page num not found")
	}
	pageStr = strings.ReplaceAll(pageStr, source, "")
	pageStr = strings.ReplaceAll(pageStr, "/", "")
	totalPageNum, err := strconv.Atoi(pageStr)
	if err != nil {
		return 0, err
	}
	return totalPageNum, nil
}
