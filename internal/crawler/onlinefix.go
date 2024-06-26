package crawler

import (
	"GameDB/internal/config"
	"GameDB/internal/constant"
	"GameDB/internal/db"
	"GameDB/internal/log"
	"GameDB/internal/model"
	"GameDB/internal/utils"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
)

var cookies map[string]string

func init() {
	cookies = make(map[string]string)
}

func CrawlOnlineFix(page int) ([]*model.GameDownload, error) {
	if !config.Config.OnlineFixAvaliable {
		log.Logger.Error("Need Online Fix account")
		return nil, errors.New("Online Fix is not available")
	}
	if len(cookies) == 0 {
		err := LoginOnlineFix()
		if err != nil {
			log.Logger.Error("Failed to login", zap.Error(err))
			return nil, err
		}
	}
	requestURL := fmt.Sprintf("%s/page/%d/", constant.OnlineFixURL, page)
	log.Logger.Info("Crawling item", zap.String("url", requestURL))
	resp, err := utils.Fetch(utils.FetchConfig{
		Url:     requestURL,
		Cookies: cookies,
		Headers: map[string]string{
			"Referer": constant.OnlineFixURL,
		},
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
	doc.Find("article.news").Each(func(i int, s *goquery.Selection) {
		urls = append(urls, s.Find(".big-link").First().AttrOr("href", ""))
		updateFlags = append(
			updateFlags,
			s.Find(".big-link").First().AttrOr("href", "")+
				s.Find("time").Text(),
		)
	})

	var res []*model.GameDownload
	for i, u := range urls {
		if db.IsOnlineFixCrawled(updateFlags[i]) {
			log.Logger.Info("Skipping already crawled item", zap.String("url", u))
			continue
		}
		log.Logger.Info("Crawling item", zap.String("URL", u))
		resp, err = utils.Fetch(utils.FetchConfig{
			Url:     u,
			Cookies: cookies,
			Headers: map[string]string{
				"Referer": constant.OnlineFixURL,
			},
		})
		if err != nil {
			log.Logger.Error("Failed to fetch", zap.Error(err))
			continue
		}
		titleRegex := regexp.MustCompile(`(?i)<h1.*?>(.*?)</h1>`)
		titleRegexRes := titleRegex.FindAllStringSubmatch(string(resp.Data), -1)
		if len(titleRegexRes) == 0 {
			log.Logger.Error("Failed to find title", zap.Error(err))
			continue
		}
		downloadRegex := regexp.MustCompile(`(?i)<a[^>]*\bhref="([^"]+)"[^>]*>(Скачать Torrent|Скачать торрент)</a>`)
		downloadRegexRes := downloadRegex.FindAllStringSubmatch(string(resp.Data), -1)
		if len(downloadRegexRes) == 0 {
			log.Logger.Error("Failed to find download button", zap.Error(err))
			continue
		}
		item, err := db.GetGameDownloadByUrl(u)
		if err != nil {
			log.Logger.Error("Failed to get game", zap.Error(err))
			continue
		}
		item.UpdateFlag = updateFlags[i]
		item.RawName = titleRegexRes[0][1]
		item.Name = OnlineFixFormatter(item.RawName)
		item.Url = u
		item.Author = "OnlineFix"
		item.Size = "0"
		resp, err = utils.Fetch(utils.FetchConfig{
			Url:     downloadRegexRes[0][1],
			Cookies: cookies,
			Headers: map[string]string{
				"Referer": u,
			},
		})
		if err != nil {
			log.Logger.Error("Failed to fetch", zap.Error(err))
			continue
		}
		if strings.Contains(downloadRegexRes[0][1], "uploads.online-fix.me") {
			magnetRegex := regexp.MustCompile(`(?i)"(.*?).torrent"`)
			magnetRegexRes := magnetRegex.FindAllStringSubmatch(string(resp.Data), -1)
			if len(magnetRegexRes) == 0 {
				log.Logger.Error("Failed to find magnet", zap.Error(err))
				continue
			}
			log.Logger.Info("Found magnet", zap.String("magnet", downloadRegexRes[0][1]+strings.Trim(magnetRegexRes[0][0], "\"")))
			resp, err = utils.Fetch(utils.FetchConfig{
				Url:     downloadRegexRes[0][1] + strings.Trim(magnetRegexRes[0][0], "\""),
				Cookies: cookies,
				Headers: map[string]string{
					"Referer": u,
				},
			})
			if err != nil {
				log.Logger.Error("Failed to fetch", zap.Error(err))
				continue
			}
			item.Magnet, item.Size, err = utils.ConvertTorrentToMagnet(resp.Data)
			if err != nil {
				log.Logger.Error("Failed to convert torrent to magnet", zap.Error(err))
				continue
			}
		} else if strings.Contains(downloadRegexRes[0][1], "online-fix.me/ext") {
			if strings.Contains(string(resp.Data), "mega.nz") {
				if !config.Config.MegaAvaliable {
					log.Logger.Error("Mega is not avaliable")
					continue
				}
				megaRegex := regexp.MustCompile(`(?i)location.href=\\'([^\\']*)\\'`)
				megaRegexRes := megaRegex.FindAllStringSubmatch(string(resp.Data), -1)
				if len(megaRegexRes) == 0 {
					log.Logger.Error("Failed to find download link")
					continue
				}
				log.Logger.Info("Downloading torrent", zap.String("URL", megaRegexRes[0][1]))
				path, files, err := utils.MegaDownload(megaRegexRes[0][1], "torrent")
				if err != nil {
					log.Logger.Error("Failed to download torrent", zap.Error(err))
					continue
				}
				torrent := ""
				for _, file := range files {
					if strings.HasSuffix(file, ".torrent") {
						torrent = file
						break
					}
				}
				dataBytes, err := os.ReadFile(torrent)
				if err != nil {
					log.Logger.Error("Failed to read torrent", zap.Error(err))
					continue
				}
				item.Magnet, item.Size, err = utils.ConvertTorrentToMagnet(dataBytes)
				if err != nil {
					log.Logger.Error("Failed to convert torrent to magnet", zap.Error(err))
					continue
				}
				err = os.RemoveAll(path)
				if err != nil {
					log.Logger.Error("Failed to remove torrent", zap.Error(err))
				}
			} else {
				log.Logger.Error("Failed to find download link")
				continue
			}
		} else {
			log.Logger.Error("Failed to find download link")
			continue
		}
		res = append(res, item)
		err = db.SaveGameDownload(item)
		if err != nil {
			log.Logger.Error("Failed to save game", zap.Error(err))
			continue
		}
	}
	return res, nil
}

func CrawlOnlineFixMulti(pages []int) ([]*model.GameDownload, error) {
	var res []*model.GameDownload
	for _, page := range pages {
		items, err := CrawlOnlineFix(page)
		if err != nil {
			return nil, err
		}
		res = append(res, items...)
	}
	log.Logger.Info("Crawled finished", zap.Int("num", len(res)))
	return res, nil
}

func CrawlOnlineFixAll() ([]*model.GameDownload, error) {
	var res []*model.GameDownload
	totalPageNum, err := GetOnlineFixTotalPageNum()
	if err != nil {
		return nil, err
	}
	for i := 1; i <= totalPageNum; i++ {
		items, err := CrawlOnlineFix(i)
		if err != nil {
			return nil, err
		}
		res = append(res, items...)
	}
	log.Logger.Info("Crawled finished", zap.Int("num", len(res)))
	return res, nil
}

func GetOnlineFixTotalPageNum() (int, error) {
	resp, err := utils.Fetch(utils.FetchConfig{
		Url: constant.OnlineFixURL,
		Headers: map[string]string{
			"Referer": constant.OnlineFixURL,
		},
	})
	if err != nil {
		log.Logger.Error("Failed to fetch", zap.Error(err))
		return 0, err
	}
	pageRegex := regexp.MustCompile(`(?i)<a href="https://online-fix.me/page/(\d+)/">.*?</a>`)
	pageRegexRes := pageRegex.FindAllStringSubmatch(string(resp.Data), -1)
	if len(pageRegexRes) == 0 {
		log.Logger.Error("Failed to find total page num", zap.Error(err))
		return 0, err
	}
	totalPageNum, err := strconv.Atoi(pageRegexRes[len(pageRegexRes)-2][1])
	if err != nil {
		log.Logger.Error("Failed to parse total page num", zap.Error(err))
		return 0, err
	}
	return totalPageNum, nil
}

type csrf struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

func LoginOnlineFix() error {
	resp, err := utils.Fetch(utils.FetchConfig{
		Url: constant.OnlineFixCSRFURL,
		Headers: map[string]string{
			"X-Requested-With": "XMLHttpRequest",
			"Referer":          constant.OnlineFixURL,
		},
	})
	if err != nil {
		log.Logger.Error("Failed to fetch", zap.Error(err))
		return err
	}
	var csrf csrf
	if err = json.Unmarshal(resp.Data, &csrf); err != nil {
		log.Logger.Error("Failed to unmarshal JSON", zap.Error(err))
		return err
	}

	for _, cookie := range resp.Cookie {
		cookies[cookie.Name] = cookie.Value
	}
	params := url.Values{}
	params.Add("login_name", config.Config.OnlineFix.User)
	params.Add("login_password", config.Config.OnlineFix.Password)
	params.Add(csrf.Field, csrf.Value)
	params.Add("login", "submit")
	resp, err = utils.Fetch(utils.FetchConfig{
		Url:     constant.OnlineFixURL,
		Method:  "POST",
		Cookies: cookies,
		Headers: map[string]string{
			"Origin":       constant.OnlineFixURL,
			"Content-Type": "application/x-www-form-urlencoded",
			"Referer":      constant.OnlineFixURL,
		},
		Data: params,
	})
	if err != nil {
		log.Logger.Error("Failed to fetch", zap.Error(err))
		return err
	}
	for _, cookie := range resp.Cookie {
		cookies[cookie.Name] = cookie.Value
	}
	return nil
}

func OnlineFixFormatter(name string) string {
	name = strings.Replace(name, "по сети", "", -1)
	reg1 := regexp.MustCompile(`(?i)\(.*?\)`)
	name = reg1.ReplaceAllString(name, "")
	return strings.TrimSpace(name)
}
