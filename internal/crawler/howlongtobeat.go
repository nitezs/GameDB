package crawler

import (
	"GameDB/internal/constant"
	"GameDB/internal/log"
	"GameDB/internal/utils"
	"encoding/json"
	"errors"
	"strings"

	"go.uber.org/zap"
)

type hltbSearchRequest struct {
	SearchPage  int      `json:"searchPage"`
	SearchTerms []string `json:"searchTerms"`
	SearchType  string   `json:"searchType"`
	Size        int      `json:"size"`
}

type hltbSearch struct {
	Color           string           `json:"color"`
	Title           string           `json:"title"`
	Category        string           `json:"category"`
	Count           int              `json:"count"`
	PageCurrent     int              `json:"pageCurrent"`
	PageTotal       int              `json:"pageTotal"`
	PageSize        int              `json:"pageSize"`
	Data            []hltbSearchData `json:"data"`
	UserData        []interface{}    `json:"userData"`
	DisplayModifier interface{}      `json:"displayModifier"`
}

type hltbSearchData struct {
	GameID          int    `json:"game_id"`
	GameName        string `json:"game_name"`
	GameNameDate    int    `json:"game_name_date"`
	GameAlias       string `json:"game_alias"`
	GameType        string `json:"game_type"`
	GameImage       string `json:"game_image"`
	CompLvlCombine  int    `json:"comp_lvl_combine"`
	CompLvlSp       int    `json:"comp_lvl_sp"`
	CompLvlCo       int    `json:"comp_lvl_co"`
	CompLvlMp       int    `json:"comp_lvl_mp"`
	CompLvlSpd      int    `json:"comp_lvl_spd"`
	CompMain        int    `json:"comp_main"`
	CompPlus        int    `json:"comp_plus"`
	Comp100         int    `json:"comp_100"`
	CompAll         int    `json:"comp_all"`
	CompMainCount   int    `json:"comp_main_count"`
	CompPlusCount   int    `json:"comp_plus_count"`
	Comp100Count    int    `json:"comp_100_count"`
	CompAllCount    int    `json:"comp_all_count"`
	InvestedCo      int    `json:"invested_co"`
	InvestedMp      int    `json:"invested_mp"`
	InvestedCoCount int    `json:"invested_co_count"`
	InvestedMpCount int    `json:"invested_mp_count"`
	CountComp       int    `json:"count_comp"`
	CountSpeedrun   int    `json:"count_speedrun"`
	CountBacklog    int    `json:"count_backlog"`
	CountReview     int    `json:"count_review"`
	ReviewScore     int    `json:"review_score"`
	CountPlaying    int    `json:"count_playing"`
	CountRetired    int    `json:"count_retired"`
	ProfileDev      string `json:"profile_dev"`
	ProfilePopular  int    `json:"profile_popular"`
	ProfileSteam    int    `json:"profile_steam"`
	ProfilePlatform string `json:"profile_platform"`
	ReleaseWorld    int    `json:"release_world"`
}

func SearchHowLongToBeat(key string) (*hltbSearchData, error) {
	log.Logger.Info("Search HowLongToBeat", zap.String("key", key))
	resp, err := utils.Fetch(utils.FetchConfig{
		Url:    constant.HowLongToBeatSearchURL,
		Method: "POST",
		Headers: map[string]string{
			"Referer": constant.HowLongToBeatURL,
		},
		Data: hltbSearchRequest{
			SearchPage:  1,
			SearchTerms: strings.Split(key, " "),
			SearchType:  "games",
			Size:        100,
		},
	})
	if err != nil {
		log.Logger.Error("Failed to fetch", zap.Error(err))
		return nil, err
	}
	res := &hltbSearch{}
	if err = json.Unmarshal(resp.Data, res); err != nil {
		log.Logger.Error("Failed to unmarshal JSON", zap.Error(err))
		return nil, err
	}
	maxSim := 0.0
	maxSimItem := &hltbSearchData{}
	for _, item := range res.Data {
		if item.GameName == key {
			return &item, nil
		} else {
			sim := utils.Similarity(key, item.GameName)
			if sim >= 0.8 && sim > maxSim {
				maxSim = sim
				maxSimItem = &item
			}
		}
	}
	if maxSim >= 0.8 {
		return maxSimItem, nil
	}
	log.Logger.Warn("Failed to find", zap.String("key", key))
	return nil, errors.New("Not found")
}
