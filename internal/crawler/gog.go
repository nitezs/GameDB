package crawler

import (
	"GameDB/internal/cache"
	"GameDB/internal/config"
	"GameDB/internal/constant"
	"GameDB/internal/db"
	"GameDB/internal/log"
	"GameDB/internal/model"
	"GameDB/internal/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

func _GetGOGID(name string, prepare bool) (int, error) {
	log.Logger.Debug("Get GOG ID", zap.String("key", name))

	if prepare {
		name = GetIDPrepared(name)
	}

	baseURL, _ := url.Parse(constant.GOGSearchURL)
	params := url.Values{}
	params.Add("mediaType", "game")
	params.Add("search", name)
	baseURL.RawQuery = params.Encode()
	resp, err := utils.Fetch(utils.FetchConfig{
		Url: baseURL.String(),
	})
	if err != nil {
		log.Logger.Error("Failed to fetch", zap.String("url", baseURL.String()), zap.Error(err))
		return 0, err
	}
	data := model.GOGSearch{}
	err = json.Unmarshal(resp.Data, &data)
	if err != nil {
		log.Logger.Error("Failed to unmarshal JSON", zap.Error(err))
		return 0, err
	}
	maxSim := 0.0
	maxSimID := 0
	for _, item := range data.Products {
		sim := utils.Similarity(item.Title, name)
		if sim > maxSim && sim >= 0.8 {
			maxSim = sim
			maxSimID = item.ID
		}
	}
	if maxSimID == 0 {
		log.Logger.Warn("GOG ID not found", zap.String("key", name))
		return 0, errors.New("GOG ID not found")
	}
	return maxSimID, nil
}

func GetGOGID(name string) (int, error) {
	id, err := _GetGOGID(name, false)
	if err == nil {
		return id, nil
	}
	id, err = _GetGOGID(name, true)
	if err == nil {
		return id, nil
	}
	return 0, errors.New("GOG ID not found")
}

func GetGOGIDCache(name string) (int, error) {
	if config.Config.RedisAvaliable {
		key := fmt.Sprintf("gog_id:%s", name)
		val, exist := cache.Redis.Get(key)
		if exist {
			id, err := strconv.Atoi(val)
			if err != nil {
				return 0, err
			}
			return id, nil
		} else {
			id, err := GetGOGID(key)
			if err != nil {
				return 0, err
			}
			err = cache.Redis.Add(key, id)
			if err != nil {
				log.Logger.Warn("Failed to add cache", zap.Error(err))
			}
			return id, nil
		}
	} else {
		return GetGOGID(name)
	}
}

func GetGOGAppDetail(id int) (*model.GOGAppDetail, error) {
	baseURL, _ := url.Parse(fmt.Sprintf("%s/%v", constant.GOGDetailsURL, id))
	params := url.Values{}
	params.Add("expand", "downloads,expanded_dlcs,description,screenshots,videos,related_products,changelog")
	params.Add("locale", "zh")
	baseURL.RawQuery = params.Encode()
	resp, err := utils.Fetch(utils.FetchConfig{
		Url: baseURL.String(),
	})
	if err != nil {
		return nil, err
	}
	data := model.GOGAppDetail{}
	err = json.Unmarshal(resp.Data, &data)
	if err != nil {
		log.Logger.Error("JSON unmarshal failed on GOG data", zap.Error(err))
		return nil, err
	}
	if data.ID == 0 {
		return nil, errors.New("GOG App not found")
	}
	return &data, nil
}

func GetGOGAppDetailCache(id int) (*model.GOGAppDetail, error) {
	if config.Config.RedisAvaliable {
		key := fmt.Sprintf("gog_app:%d", id)
		val, exist := cache.Redis.Get(key)
		if exist {
			data := model.GOGAppDetail{}
			err := json.Unmarshal([]byte(val), &data)
			if err != nil {
				log.Logger.Error("JSON unmarshal failed on GOG data", zap.Error(err))
				return nil, err
			}
			return &data, nil
		} else {
			data, err := GetGOGAppDetail(id)
			if err != nil {
				return nil, err
			}
			dataBytes, err := json.Marshal(data)
			if err != nil {
				return nil, err
			}
			err = cache.Redis.Add(key, dataBytes)
			if err != nil {
				log.Logger.Warn("Failed to add cache", zap.Error(err))
				return data, nil
			}
			return data, nil
		}
	} else {
		return GetGOGAppDetail(id)
	}
}

func GenerateGOGGameInfo(id int) (*model.GameInfo, error) {
	item := &model.GameInfo{}
	info, err := GetGOGAppDetailCache(id)
	if err != nil {
		return nil, err
	}
	item.GOGID = id
	item.Name = info.Title
	item.Description = info.Description.Full
	item.Cover = info.Images.Logo2X
	for _, language := range info.Languages {
		item.Languages = append(item.Languages, language)
	}
	screenshots := []string{}
	for _, screenshot := range info.Screenshots {
		screenshots = append(screenshots, strings.Replace(screenshot.FormatterTemplateURL, "{formatter}", "ggvgl_2x", -1))
	}
	item.Screenshots = screenshots
	return item, nil
}

func ProcessGameWithGOG(game *model.GameDownload) (*model.GameInfo, error) {
	id, err := GetGOGID(game.Name)
	if err != nil {
		return nil, err
	}
	d, err := db.GetGameInfoByPlatformID("gog", id)
	if err == nil {
		d.GameIDs = append(d.GameIDs, game.ID)
		d.GameIDs = utils.Unique(d.GameIDs)
		return d, nil
	}
	detail, err := GenerateGameInfo("gog", id)
	if err != nil {
		return nil, err
	}
	detail.GameIDs = append(detail.GameIDs, game.ID)
	detail.GameIDs = utils.Unique(detail.GameIDs)
	return detail, nil
}
