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
	"slices"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

var TwitchToken string

func _GetIGDBID(name string, prepare bool) (int, error) {
	if prepare {
		name = GetIDPrepared(name)
	}
	log.Logger.Debug("Get IGDB ID", zap.String("key", name))
	name = regexp.MustCompile(`[:\-]`).ReplaceAllString(name, " ")
	name = regexp.MustCompile(`\s+`).ReplaceAllString(name, " ")
	var err error
	if TwitchToken == "" {
		TwitchToken, err = LoginTwitch()
		if err != nil {
			log.Logger.Error("Failed to login", zap.Error(err))
			return 0, err
		}
	}
	resp, err := utils.Fetch(utils.FetchConfig{
		Url: constant.IGDBGameURL,
		Headers: map[string]string{
			"Client-ID":     config.Config.Twitch.ClientID,
			"Authorization": "Bearer " + TwitchToken,
			"User-Agent":    "",
			"Content-Type":  "text/plain",
		},
		Data:   fmt.Sprintf(`search "%s"; fields *; limit 40;`, name),
		Method: "POST",
	})
	if err != nil {
		log.Logger.Error("Failed to fetch", zap.Error(err))
		return 0, err
	}
	var data model.IGDBGameDetails
	if err = json.Unmarshal(resp.Data, &data); err != nil {
		log.Logger.Error("Failed to unmarshal", zap.Error(err))
		return 0, err
	}
	if len(data) == 1 {
		return data[0].ID, nil
	}
	maxSim := 0.0
	maxSimID := 0
	maxSimItem := &model.IGDBGameDetail{}
	for _, item := range data {
		if item.Platforms != nil && len(item.Platforms) > 0 {
			if !slices.Contains(item.Platforms, 6) &&
				!slices.Contains(item.Platforms, 130) {
				continue
			}
		}
		if item.Name == name {
			return item.ID, nil
		} else {
			sim := utils.Similarity(item.Name, name)
			if sim > 0.8 && sim > maxSim {
				maxSim = sim
				maxSimID = item.ID
				maxSimItem = item
			}
		}
	}
	if maxSimID == 0 {
		log.Logger.Error("IGDB ID not found", zap.String("key", name))
		return 0, errors.New("IGDB ID not found")
	}
	if maxSimItem.ParentGame != 0 {
		log.Logger.Info("Found IGDB ID", zap.Int("id", maxSimItem.ParentGame), zap.String("key", name))
		return maxSimItem.ParentGame, nil
	}
	log.Logger.Info("Found IGDB ID", zap.Int("id", maxSimID), zap.String("key", name))
	return maxSimID, nil
}

func GetIGDBID(name string) (int, error) {
	id, err := _GetIGDBID(name, false)
	if err == nil {
		return id, nil
	}
	id, err = _GetIGDBID(name, true)
	if err == nil {
		return id, nil
	}
	return 0, errors.New("IGDB ID not found")
}

func GetIGDBIDCache(name string) (int, error) {
	if config.Config.RedisAvaliable {
		key := fmt.Sprintf("igdb_id:%s", name)
		val, exist := cache.Redis.Get(key)
		if exist {
			id, err := strconv.Atoi(val)
			if err != nil {
				return 0, err
			}
			return id, nil
		} else {
			id, err := GetIGDBID(key)
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
		return GetIGDBID(name)
	}
}

func GetIGDBAppDetail(id int) (*model.IGDBGameDetail, error) {
	var err error
	if TwitchToken == "" {
		TwitchToken, err = LoginTwitch()
		if err != nil {
			return nil, err
		}
	}
	resp, err := utils.Fetch(utils.FetchConfig{
		Url: constant.IGDBGameURL,
		Headers: map[string]string{
			"Client-ID":     config.Config.Twitch.ClientID,
			"Authorization": "Bearer " + TwitchToken,
			"User-Agent":    "",
			"Content-Type":  "text/plain",
		},
		Data:   fmt.Sprintf(`where id=%v; fields *;`, id),
		Method: "POST",
	})
	if err != nil {
		return nil, err
	}
	var data model.IGDBGameDetails
	if err = json.Unmarshal(resp.Data, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errors.New("IGDB App not found")
	}
	if data[0].Name == "" {
		return GetIGDBAppDetail(id)
	}
	return data[0], nil
}

func GetIGDBAppDetailCache(id int) (*model.IGDBGameDetail, error) {
	if config.Config.RedisAvaliable {
		key := fmt.Sprintf("igdb_game:%v", id)
		val, exist := cache.Redis.Get(key)
		if exist {
			var data model.IGDBGameDetail
			if err := json.Unmarshal([]byte(val), &data); err != nil {
				return nil, err
			}
			return &data, nil
		} else {
			data, err := GetIGDBAppDetail(id)
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
		return GetIGDBAppDetail(id)
	}
}

func GetIGDBScreenshots(game int) ([]string, error) {
	var err error
	if TwitchToken == "" {
		TwitchToken, err = LoginTwitch()
		if err != nil {
			return nil, err
		}
	}
	resp, err := utils.Fetch(utils.FetchConfig{
		Url: constant.IGDBScreenshotsURL,
		Headers: map[string]string{
			"Client-ID":     config.Config.Twitch.ClientID,
			"Authorization": "Bearer " + TwitchToken,
			"User-Agent":    "",
			"Content-Type":  "text/plain",
		},
		Data:   fmt.Sprintf(`where game=%v; fields *;`, game),
		Method: "POST",
	})
	if err != nil {
		return nil, err
	}
	var data model.IGDBScreenshots
	if err = json.Unmarshal(resp.Data, &data); err != nil {
		return nil, err
	}
	var screenshots []string
	for _, item := range data {
		screenshots = append(screenshots, strings.Replace(item.URL, "t_thumb", "t_original", 1))
	}
	if len(screenshots) == 0 {
		return nil, errors.New("screenshots not found")
	}
	if screenshots[0] == "" {
		return GetIGDBScreenshots(game) // server sometimes return wrong data, just contains id
	}
	return screenshots, nil
}

func GetIGDBScreenshotsCache(game int) ([]string, error) {
	if config.Config.RedisAvaliable {
		key := fmt.Sprintf("igdb_screenshots:%v", game)
		val, exist := cache.Redis.Get(key)
		if exist {
			var data []string
			if err := json.Unmarshal([]byte(val), &data); err != nil {
				return nil, err
			}
			return data, nil
		} else {
			data, err := GetIGDBScreenshots(game)
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
		return GetIGDBScreenshots(game)
	}
}

func GetIGDBCovers(game int) (string, error) {
	var err error
	if TwitchToken == "" {
		TwitchToken, err = LoginTwitch()
		if err != nil {
			return "", err
		}
	}
	resp, err := utils.Fetch(utils.FetchConfig{
		Url: constant.IGDBCoversURL,
		Headers: map[string]string{
			"Client-ID":     config.Config.Twitch.ClientID,
			"Authorization": "Bearer " + TwitchToken,
			"User-Agent":    "",
			"Content-Type":  "text/plain",
		},
		Data:   fmt.Sprintf(`where game=%v; fields *;`, game),
		Method: "POST",
	})
	if err != nil {
		return "", err
	}
	var data model.IGDBCovers
	if err = json.Unmarshal(resp.Data, &data); err != nil {
		return "", err
	}
	for _, item := range data {
		return strings.Replace(item.URL, "t_thumb", "t_original", 1), nil
	}
	if len(data) == 0 {
		return "", errors.New("cover not found")
	}
	if data[0].URL == "" {
		return GetIGDBCovers(game)
	}
	return "", errors.New("cover not found")
}

func GetIGDBCoversCache(game int) (string, error) {
	if config.Config.RedisAvaliable {
		key := fmt.Sprintf("igdb_covers:%v", game)
		val, exist := cache.Redis.Get(key)
		if exist {
			return val, nil
		} else {
			val, err := GetIGDBCovers(game)
			if err != nil {
				return "", err
			}
			err = cache.Redis.Add(key, val)
			if err != nil {
				return "", err
			}
			return val, nil
		}
	} else {
		return GetIGDBCovers(game)
	}
}

func GetIGDBLanguagesID(game int) ([]*model.Language, error) {
	var err error
	if TwitchToken == "" {
		TwitchToken, err = LoginTwitch()
		if err != nil {
			return nil, err
		}
	}
	resp, err := utils.Fetch(utils.FetchConfig{
		Url: constant.IGDBLanguageSupportsURL,
		Headers: map[string]string{
			"Client-ID":     config.Config.Twitch.ClientID,
			"Authorization": "Bearer " + TwitchToken,
			"User-Agent":    "",
			"Content-Type":  "text/plain",
		},
		Data:   fmt.Sprintf(`where game=%v; fields *;`, game),
		Method: "POST",
	})
	if err != nil {
		return nil, err
	}
	var data model.IGDBLanguageSupports
	if err = json.Unmarshal(resp.Data, &data); err != nil {
		return nil, err
	}
	languages := []int{}
	if len(data) == 0 {
		return []*model.Language{}, nil
	}
	if data[0].Language == 0 {
		return GetIGDBLanguagesID(game)
	}
	for _, item := range data {
		languages = append(languages, item.Language)
	}
	languages = utils.Unique(languages)
	return db.GetLanguages(languages)
}

func GetIGDBLanguagesIDCache(game int) ([]*model.Language, error) {
	if config.Config.RedisAvaliable {
		key := fmt.Sprintf("igdb_languages:%v", game)
		val, exist := cache.Redis.Get(key)
		if exist {
			var data []*model.Language
			if err := json.Unmarshal([]byte(val), &data); err != nil {
				return nil, err
			}
			return data, nil
		} else {
			val, err := GetIGDBLanguagesID(game)
			if err != nil {
				return nil, err
			}
			data, err := json.Marshal(val)
			if err != nil {
				return nil, err
			}
			err = cache.Redis.Add(key, data)
			if err != nil {
				return nil, err
			}
			return val, nil
		}
	} else {
		return GetIGDBLanguagesID(game)
	}
}

func GetIGDBAliases(game int) ([]string, error) {
	var err error
	if TwitchToken == "" {
		TwitchToken, err = LoginTwitch()
		if err != nil {
			return nil, err
		}
	}
	resp, err := utils.Fetch(utils.FetchConfig{
		Url: constant.IGDBAlternativeNamesURL,
		Headers: map[string]string{
			"Client-ID":     config.Config.Twitch.ClientID,
			"Authorization": "Bearer " + TwitchToken,
			"User-Agent":    "",
			"Content-Type":  "text/plain",
		},
		Data:   fmt.Sprintf(`where game=%v; fields *;`, game),
		Method: "POST",
	})
	if err != nil {
		return nil, err
	}
	var data model.IGDBAlternativeNames
	if err = json.Unmarshal(resp.Data, &data); err != nil {
		return nil, err
	}
	aliases := []string{}
	for _, item := range data {
		aliases = append(aliases, item.Name)
	}
	if len(aliases) == 0 {
		return []string{}, errors.New("no aliases found")
	}
	utils.Unique(aliases)
	if aliases[0] == "" {
		return GetIGDBAliases(game)
	}
	return aliases, nil
}

func GetIGDBAliasesCache(game int) ([]string, error) {
	if config.Config.RedisAvaliable {
		key := fmt.Sprintf("igdb_aliases:%v", game)
		val, exist := cache.Redis.Get(key)
		if exist {
			var data []string
			err := json.Unmarshal([]byte(val), &data)
			if err != nil {
				return nil, err
			}
			return data, nil
		} else {
			data, err := GetIGDBAliases(game)
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
		return GetIGDBAliases(game)
	}
}

func SaveIGDBLanguagesList() error {
	var err error
	if TwitchToken == "" {
		TwitchToken, err = LoginTwitch()
		if err != nil {
			return err
		}
	}
	resp, err := utils.Fetch(utils.FetchConfig{
		Url: constant.IGDBLanguagesURL,
		Headers: map[string]string{
			"Client-ID":     config.Config.Twitch.ClientID,
			"Authorization": "Bearer " + TwitchToken,
			"User-Agent":    "",
			"Content-Type":  "text/plain",
		},
		Data:   `fields *; limit 100;`,
		Method: "POST",
	})
	if err != nil {
		return err
	}
	var data model.IGDBLanguages
	if err = json.Unmarshal(resp.Data, &data); err != nil {
		return err
	}
	for _, item := range data {
		err = db.SaveLanguage(&model.Language{
			LID:        item.ID,
			Name:       item.Name,
			NativeName: item.NativeName,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func LoginTwitch() (string, error) {
	baseURL, _ := url.Parse(constant.TwitchAuthURL)
	params := url.Values{}
	params.Add("client_id", config.Config.Twitch.ClientID)
	params.Add("client_secret", config.Config.Twitch.ClientSecret)
	params.Add("grant_type", "client_credentials")
	baseURL.RawQuery = params.Encode()
	resp, err := utils.Fetch(utils.FetchConfig{
		Url:    baseURL.String(),
		Method: "POST",
		Headers: map[string]string{
			"User-Agent": "",
		},
	})
	if err != nil {
		return "", err
	}
	data := struct {
		AccessToken string `json:"access_token"`
	}{}
	err = json.Unmarshal(resp.Data, &data)
	if err != nil {
		return "", err
	}
	return data.AccessToken, nil
}

func GetIGDBInvolvedCompanies(game int) (model.IGDBInvolvedCompanies, error) {
	var err error
	if TwitchToken == "" {
		TwitchToken, err = LoginTwitch()
		if err != nil {
			return nil, err
		}
	}
	resp, err := utils.Fetch(utils.FetchConfig{
		Url: constant.IGDBInvolvedCompaniesURL,
		Headers: map[string]string{
			"Client-ID":     config.Config.Twitch.ClientID,
			"Authorization": "Bearer " + TwitchToken,
			"User-Agent":    "",
			"Content-Type":  "text/plain",
		},
		Data:   fmt.Sprintf(`where game=%v; fields *;`, game),
		Method: "POST",
	})
	if err != nil {
		return nil, err
	}
	var data model.IGDBInvolvedCompanies
	if err = json.Unmarshal(resp.Data, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errors.New("no data found")
	}
	if data[0].Company == 0 {
		return GetIGDBInvolvedCompanies(game)
	}
	return data, nil
}

func GetIGDBInvolvedCompaniesCache(game int) (model.IGDBInvolvedCompanies, error) {
	if config.Config.RedisAvaliable {
		key := fmt.Sprintf("igdb_involved_companies:%v", game)
		val, exist := cache.Redis.Get(key)
		if exist {
			var data model.IGDBInvolvedCompanies
			err := json.Unmarshal([]byte(val), &data)
			if err != nil {
				return nil, err
			}
			return data, nil
		} else {
			data, err := GetIGDBInvolvedCompanies(game)
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
		return GetIGDBInvolvedCompanies(game)
	}
}

func GetIGDBCompany(id int) (string, error) {
	var err error
	if TwitchToken == "" {
		TwitchToken, err = LoginTwitch()
		if err != nil {
			return "", err
		}
	}
	resp, err := utils.Fetch(utils.FetchConfig{
		Url: constant.IGDBCompaniesURL,
		Headers: map[string]string{
			"Client-ID":     config.Config.Twitch.ClientID,
			"Authorization": "Bearer " + TwitchToken,
			"User-Agent":    "",
			"Content-Type":  "text/plain",
		},
		Data:   fmt.Sprintf(`where id=%v; fields *;`, id),
		Method: "POST",
	})
	if err != nil {
		return "", err
	}
	var data model.IGDBCompanies
	if err = json.Unmarshal(resp.Data, &data); err != nil {
		return "", err
	}
	if len(data) == 0 {
		return "", errors.New("Not found")
	}
	if data[0].Name == "" {
		return GetIGDBCompany(id)
	}
	return data[0].Name, nil
}

func GetIGDBCompanyCache(id int) (string, error) {
	if config.Config.RedisAvaliable {
		key := fmt.Sprintf("igdb_companies:%v", id)
		val, exist := cache.Redis.Get(key)
		if exist {
			return val, nil
		} else {
			data, err := GetIGDBCompany(id)
			if err != nil {
				return "", err
			}
			err = cache.Redis.Add(key, data)
			if err != nil {
				log.Logger.Warn("Failed to add cache", zap.Error(err))
				return data, nil
			}
			return data, nil
		}
	} else {
		return GetIGDBCompany(id)
	}
}

func GenerateIGDBGameInfo(id int) (*model.GameInfo, error) {
	item := &model.GameInfo{}
	detail, err := GetIGDBAppDetailCache(id)
	if err != nil {
		return nil, err
	}
	item.IGDBID = id
	item.Name = detail.Name
	item.Description = detail.Summary

	languages, err := GetIGDBLanguagesIDCache(id)
	if err != nil {
		log.Logger.Warn("Failed to get igdb languages", zap.Error(err))
	}
	for _, language := range languages {
		item.Languages = append(item.Languages, language.Name)
	}

	if len(detail.Screenshots) > 0 {
		screenshots, err := GetIGDBScreenshotsCache(id)
		if err != nil {
			log.Logger.Warn("Failed to get igdb screenshots", zap.Error(err))
		}
		item.Screenshots = append(item.Screenshots, screenshots...)
	}

	if len(detail.AlternativeNames) > 0 {
		aliases, err := GetIGDBAliasesCache(detail.ID)
		if err != nil {
			log.Logger.Warn("Failed to get igdb aliases", zap.Error(err))
		}
		item.Aliases = append(item.Aliases, aliases...)
	}

	item.Cover, err = GetIGDBCoversCache(id)
	if err != nil {
		log.Logger.Warn("Failed to get igdb cover", zap.Error(err))
	}

	if len(detail.InvolvedCompanies) > 0 {
		involvedCompanies, err := GetIGDBInvolvedCompaniesCache(id)
		if err != nil {
			log.Logger.Warn("Failed to get igdb involved companies", zap.Error(err))
		} else {
			for _, company := range involvedCompanies {
				if company.Developer || company.Publisher {
					companyName, err := GetIGDBCompanyCache(company.Company)
					if err != nil {
						log.Logger.Warn("Failed to get igdb company", zap.Error(err))
						continue
					}
					if company.Developer {
						item.Developers = append(item.Developers, companyName)
					} else if company.Publisher {
						item.Publishers = append(item.Publishers, companyName)
					}
				}
			}
		}
	}

	return item, nil
}

func ProcessGameWithIGDB(game *model.GameDownload) (*model.GameInfo, error) {
	id, err := GetIGDBID(game.Name)
	if err != nil {
		return nil, err
	}
	d, err := db.GetGameInfoByPlatformID("igdb", id)
	if err == nil {
		d.GameIDs = append(d.GameIDs, game.ID)
		d.GameIDs = utils.Unique(d.GameIDs)
		return d, nil
	}
	info, err := GenerateGameInfo("igdb", id)
	if err != nil {
		return nil, err
	}
	info.GameIDs = append(info.GameIDs, game.ID)
	info.GameIDs = utils.Unique(info.GameIDs)
	return info, nil
}
