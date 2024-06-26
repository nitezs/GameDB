package crawler

import (
	"GameDB/internal/cache"
	"GameDB/internal/config"
	"GameDB/internal/constant"
	"GameDB/internal/log"
	"GameDB/internal/model"
	"GameDB/internal/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
)

func GetSteam250(url string) ([]model.Steam250Item, error) {
	resp, err := utils.Fetch(utils.FetchConfig{
		Url: url,
	})
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Data))
	if err != nil {
		return nil, err
	}
	var res []model.Steam250Item
	var item model.Steam250Item
	doc.Find(".appline").Each(func(i int, s *goquery.Selection) {
		if s.Find(".free").Length() > 0 {
			return
		}
		item.Name = s.Find(".title>a").First().Text()
		idStr := s.Find(".store").AttrOr("href", "")
		idSlice := regexp.MustCompile(`app/(\d+)/`).FindStringSubmatch(idStr)
		if len(idSlice) < 2 {
			return
		}
		item.SteamID, _ = strconv.Atoi(idSlice[1])
		res = append(res, item)
	})
	return res[:10], nil
}

func GetSteam250Top250() ([]model.Steam250Item, error) {
	return GetSteam250(constant.Steam250Top250URL)
}

func GetSteam250Top250Cache() ([]model.Steam250Item, error) {
	return GetSteam250Cache("top250", GetSteam250Top250)
}

func GetSteam250BestOfTheYear() ([]model.Steam250Item, error) {
	return GetSteam250(fmt.Sprintf(constant.Steam250BestOfTheYearURL, time.Now().UTC().Year()))
}

func GetSteam250BestOfTheYearCache() ([]model.Steam250Item, error) {
	return GetSteam250Cache(fmt.Sprintf("bestoftheyear:%v", time.Now().UTC().Year()), GetSteam250BestOfTheYear)
}

func GetSteam250WeekTop50() ([]model.Steam250Item, error) {
	return GetSteam250(constant.Steam250WeekTop50URL)
}

func GetSteam250WeekTop50Cache() ([]model.Steam250Item, error) {
	return GetSteam250Cache("weektop50", GetSteam250WeekTop50)
}

func GetSteam250MostPlayed() ([]model.Steam250Item, error) {
	return GetSteam250(constant.Steam250MostPlayedURL)
}

func GetSteam250MostPlayedCache() ([]model.Steam250Item, error) {
	return GetSteam250Cache("mostplayed", GetSteam250MostPlayed)
}

func GetSteam250Cache(k string, f func() ([]model.Steam250Item, error)) ([]model.Steam250Item, error) {
	if config.Config.RedisAvaliable {
		key := k
		val, exist := cache.Redis.Get(key)
		if exist {
			var res []model.Steam250Item
			err := json.Unmarshal([]byte(val), &res)
			if err != nil {
				return nil, err
			}
			return res, nil
		} else {
			data, err := f()
			if err != nil {
				return nil, err
			}
			dataBytes, err := json.Marshal(data)
			if err != nil {
				log.Logger.Warn("JSON marshal failed", zap.Error(err))
				return data, nil
			}
			err = cache.Redis.AddWithExpire(key, dataBytes, 24*time.Hour)
			if err != nil {
				log.Logger.Warn("Failed to add cache", zap.Error(err))
				return data, nil
			}
			return data, nil
		}
	} else {
		return f()
	}
}
