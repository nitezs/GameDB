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
	"regexp"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

func GetSteamIDFromSearchPage(name string) (int, error) {
	log.Logger.Debug("Get Steam ID", zap.String("key", name))
	baseURL, _ := url.Parse(constant.SteamSearchURL)
	params := url.Values{}
	params.Add("term", name)
	baseURL.RawQuery = params.Encode()

	log.Logger.Debug("Get Steam ID", zap.String("url", baseURL.String()))
	resp, err := utils.Fetch(utils.FetchConfig{
		Url: baseURL.String(),
	})
	if err != nil {
		log.Logger.Error("Failed to fetch", zap.String("url", baseURL.String()), zap.Error(err))
		return 0, err
	}
	idRegex := regexp.MustCompile(`data-ds-appid="(.*?)"`)
	nameRegex := regexp.MustCompile(`<span class="title">(.*?)</span>`)
	idRegexRes := idRegex.FindAllStringSubmatch(string(resp.Data), -1)
	nameRegexRes := nameRegex.FindAllStringSubmatch(string(resp.Data), -1)

	if len(idRegexRes) == 0 {
		return 0, errors.New("Steam ID not found")
	}

	maxSim := 0.0
	maxSimID := 0
	for i, id := range idRegexRes {
		idStr := id[1]
		nameStr := nameRegexRes[i][1]
		if index := strings.Index(idStr, ","); index != -1 {
			idStr = idStr[:index]
		}
		if strings.EqualFold(strings.TrimSpace(nameStr), strings.TrimSpace(name)) {
			log.Logger.Info("Steam ID found", zap.String("key", name), zap.String("id", idStr))
			return strconv.Atoi(idStr)
		} else {
			sim := utils.Similarity(nameStr, name)
			log.Logger.Debug("Similarity", zap.String("str1", name), zap.String("str2", nameStr), zap.Float64("similarity", sim))
			if sim >= 0.8 && sim > maxSim {
				maxSim = sim
				maxSimID, _ = strconv.Atoi(idStr)
			}
		}
	}
	if maxSimID != 0 {
		log.Logger.Info("Steam ID found", zap.String("key", name), zap.String("id", strconv.Itoa(maxSimID)))
		return maxSimID, nil
	}
	log.Logger.Info("Steam ID not found", zap.String("key", name))
	return 0, errors.New("Steam ID not found")
}

func GetSteamIDFromSteamDB(name string) (int, error) {
	log.Logger.Debug("Get Steam ID", zap.String("key", name))
	baseURL, _ := url.Parse(constant.GoogleSearchURL)
	params := url.Values{}
	params.Add("q", fmt.Sprintf("%s site:steamdb.info/app", name))
	baseURL.RawQuery = params.Encode()
	log.Logger.Debug("Get Steam ID", zap.String("url", baseURL.String()))
	// dataStr, err := utils.FlareSolverr(baseURL.String())
	urls, err := utils.BingSearch(fmt.Sprintf("%s site:steamdb.info/app", name))
	dataStr := strings.Join(urls, "\n")
	if err != nil {
		log.Logger.Error("Failed to fetch", zap.String("url", baseURL.String()), zap.Error(err))
		return 0, err
	}
	urlRegex := regexp.MustCompile(`(?i)https://steamdb.info/app/(\d*)`)
	urlRegexRes := urlRegex.FindAllStringSubmatch(dataStr, -1)
	if len(urlRegexRes) == 0 {
		log.Logger.Warn("Steam ID not found", zap.String("key", name))
		return 0, errors.New("Steam ID not found")
	}
	maxSim := 0.0
	maxSimID := 0
	ids := []int{}
	for _, url := range urlRegexRes {
		idStr := url[1]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Logger.Warn("Failed to convert Steam ID", zap.String("key", name), zap.String("id", idStr), zap.Error(err))
			continue
		}
		ids = append(ids, id)
	}
	ids = utils.Unique(ids)
	details, err := GetSteamAppDetailsCache(ids)
	if err != nil {
		return 0, err
	}
	for id, detail := range details {
		sim := utils.Similarity(detail.Data.Name, name)
		log.Logger.Debug("Similarity", zap.String("str1", name), zap.String("str2", detail.Data.Name), zap.Float64("similarity", sim))
		if sim >= 0.9 && sim > maxSim {
			maxSim = sim
			maxSimID, _ = strconv.Atoi(id)
		}
	}
	if maxSimID != 0 {
		log.Logger.Info("Steam ID found", zap.String("key", name), zap.String("id", strconv.Itoa(maxSimID)))
		return maxSimID, nil
	}
	log.Logger.Warn("Steam ID not found", zap.String("key", name))
	return 0, errors.New("Steam ID not found")
}

func GetIDPrepared(key string) string {
	key = regexp.MustCompile(`(?i)\b(\w+)\s+(Edition|Vision|Collection)\b`).ReplaceAllString(key, " ")
	key = regexp.MustCompile(`(?i)GOTY`).ReplaceAllString(key, "Game of the year")
	key = regexp.MustCompile(`(?i)nsw for pc`).ReplaceAllString(key, "")
	key = regexp.MustCompile(`\s+`).ReplaceAllString(key, " ")
	return strings.TrimSpace(key)
}

func _GetSteamID(name string, prepare bool) (int, error) {
	if prepare {
		name = GetIDPrepared(name)
	}
	id, err := GetSteamIDFromSearchPage(name)
	if err == nil {
		return int(id), nil
	}
	// id, err = GetSteamIDFromSteamDB(key)
	// if err == nil {
	// 	return int(id), nil
	// }
	return 0, errors.New("Steam ID not found")
}

func GetSteamID(key string) (int, error) {
	key = GetIDPrepared(key)
	id, err := _GetSteamID(key, false)
	if err == nil {
		return id, nil
	}
	id, err = _GetSteamID(key, true)
	if err == nil {
		return id, nil
	}
	return 0, errors.New("Steam ID not found")
}

func GetSteamIDCache(key string) (int, error) {
	if config.Config.RedisAvaliable {
		key := fmt.Sprintf("steam_id:%s", key)
		val, exist := cache.Redis.Get(key)
		if exist {
			id, err := strconv.Atoi(val)
			if err != nil {
				return 0, err
			}
			return id, nil
		} else {
			id, err := GetSteamID(key)
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
		return GetSteamID(key)
	}
}

func GetSteamAppDetail(id int) (*model.SteamAppDetail, error) {
	baseURL, _ := url.Parse(constant.SteamAppDetailURL)
	params := url.Values{}
	params.Add("appids", strconv.Itoa(id))
	// params.Add("l", "schinese")
	baseURL.RawQuery = params.Encode()
	log.Logger.Debug("Get Steam App Detail", zap.String("url", baseURL.String()))
	resp, err := utils.Fetch(utils.FetchConfig{
		Url: baseURL.String(),
		Headers: map[string]string{
			"User-Agent": "",
		},
	})
	if err != nil {
		log.Logger.Error("Failed to fetch", zap.String("url", baseURL.String()), zap.Error(err))
		return nil, err
	}
	var detail map[string]*model.SteamAppDetail
	if err = json.Unmarshal(resp.Data, &detail); err != nil {
		log.Logger.Error("Failed to unmarshal JSON", zap.Error(err))
		return nil, err
	}
	if _, ok := detail[strconv.Itoa(id)]; !ok {
		return nil, errors.New("Steam App not found")
	}
	if detail[strconv.Itoa(id)] == nil {
		return nil, errors.New("Steam App not found")
	}
	return detail[strconv.Itoa(id)], nil
}

func GetSteamAppDetailCache(id int) (*model.SteamAppDetail, error) {
	if config.Config.RedisAvaliable {
		key := fmt.Sprintf("steam_game:%d", id)
		val, exist := cache.Redis.Get(key)
		if exist {
			var detail model.SteamAppDetail
			if err := json.Unmarshal([]byte(val), &detail); err != nil {
				return nil, err
			}
			return &detail, nil
		} else {
			data, err := GetSteamAppDetail(id)
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
		return GetSteamAppDetail(id)
	}
}

func GetSteamAppDetailsCache(ids []int) (map[string]model.SteamAppDetail, error) {
	res := make(map[string]model.SteamAppDetail)
	for _, id := range ids {
		detail, err := GetSteamAppDetail(id)
		if err != nil {
			return nil, err
		}
		res[strconv.Itoa(id)] = *detail
	}
	return res, nil
}

func GenerateSteamGameInfo(id int) (*model.GameInfo, error) {
	item := &model.GameInfo{}
	detail, err := GetSteamAppDetailCache(id)
	if err != nil {
		return nil, err
	}
	item.SteamID = id
	item.Name = detail.Data.Name
	item.Description = detail.Data.ShortDescription
	item.Cover = fmt.Sprintf("https://shared.cloudflare.steamstatic.com/store_item_assets/steam/apps/%v/library_600x900_2x.jpg", id)
	item.Developers = detail.Data.Developers
	item.Publishers = detail.Data.Publishers
	screenshots := []string{}
	for _, screenshot := range detail.Data.Screenshots {
		screenshots = append(screenshots, screenshot.PathFull)
	}
	item.Screenshots = screenshots
	return item, nil
}

func ProcessGameWithSteam(game *model.GameDownload) (*model.GameInfo, error) {
	id, err := GetSteamID(game.Name)
	if err != nil {
		return nil, err
	}
	d, err := db.GetGameInfoByPlatformID("steam", id)
	if err == nil {
		d.GameIDs = append(d.GameIDs, game.ID)
		d.GameIDs = utils.Unique(d.GameIDs)
		return d, nil
	}
	detail, err := GenerateGameInfo("steam", id)
	if err != nil {
		return nil, err
	}
	detail.GameIDs = append(detail.GameIDs, game.ID)
	detail.GameIDs = utils.Unique(detail.GameIDs)
	return detail, nil
}
