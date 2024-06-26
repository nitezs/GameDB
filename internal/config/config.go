package config

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"strconv"
)

type SConfig struct {
	LogLevel              string       `json:"log_level"`
	Database              Database     `json:"database"`
	Redis                 Redis        `json:"redis"`
	FlareSolverr          FlareSolverr `json:"flaresolverr"`
	OnlineFix             OnlineFix    `json:"online_fix"`
	Twitch                Twitch       `json:"twitch"`
	FlareSolverrAvaliable bool
	OnlineFixAvaliable    bool
	MegaAvaliable         bool
	RedisAvaliable        bool
	AutoCrawl             bool
}

type Database struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type Twitch struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type Redis struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	DBIndex  int    `json:"db_index"`
}

type FlareSolverr struct {
	Url string `json:"url"`
}

type OnlineFix struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

var Config SConfig

func InitConfig() {
	Config = SConfig{
		LogLevel:      "info",
		Database:      Database{},
		FlareSolverr:  FlareSolverr{},
		MegaAvaliable: TestMega(),
	}
	if _, err := os.Stat("config.json"); err == nil {
		configData, err := os.ReadFile("config.json")
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(configData, &Config)
		if err != nil {
			panic(err)
		}
	}
	loadEnvVariables()
}

func loadEnvVariables() {
	if env := os.Getenv("LOG_LEVEL"); env != "" {
		Config.LogLevel = env
	}
	if env := os.Getenv("DB_HOST"); env != "" {
		Config.Database.Host = env
	}
	if env := os.Getenv("DB_PORT"); env != "" {
		Config.Database.Port, _ = strconv.Atoi(env) // Note: error handling can be improved
	}
	if env := os.Getenv("DB_USER"); env != "" {
		Config.Database.User = env
	}
	if env := os.Getenv("DB_PASSWORD"); env != "" {
		Config.Database.Password = env
	}
	if env := os.Getenv("DB_NAME"); env != "" {
		Config.Database.Database = env
	}
	if env := os.Getenv("TWITCH_CLIENT_ID"); env != "" {
		Config.Twitch.ClientID = env
	}
	if env := os.Getenv("TWITCH_CLIENT_SECRET"); env != "" {
		Config.Twitch.ClientSecret = env
	}
	if env := os.Getenv("REDIS_HOST"); env != "" {
		Config.Redis.Host = env
	}
	if env := os.Getenv("REDIS_PORT"); env != "" {
		Config.Redis.Port, _ = strconv.Atoi(env)
	}
	if env := os.Getenv("REDIS_PASSWORD"); env != "" {
		Config.Redis.Password = env
	}
	if env := os.Getenv("REDIS_DB_INDEX"); env != "" {
		Config.Redis.DBIndex, _ = strconv.Atoi(env)
	}
	if env := os.Getenv("FLARESOLVERR_URL"); env != "" {
		Config.FlareSolverr.Url = env
	}
	if env := os.Getenv("ONLINE_FIX_USER"); env != "" {
		Config.OnlineFix.User = env
	}
	if env := os.Getenv("ONLINE_FIX_PASSWORD"); env != "" {
		Config.OnlineFix.Password = env
	}
	if env := os.Getenv("AUTO_CRAWL"); env != "" {
		Config.AutoCrawl, _ = strconv.ParseBool(env)
	}
	if Config.FlareSolverr.Url != "" {
		Config.FlareSolverrAvaliable = true
	}
	if Config.OnlineFix.User != "" && Config.OnlineFix.Password != "" {
		Config.OnlineFixAvaliable = true
	}
	if Config.Redis.Host != "" && Config.Redis.Port != 0 {
		Config.RedisAvaliable = true
	}
	if Config.Twitch.ClientID == "" || Config.Twitch.ClientSecret == "" {
		panic("Missing twitch credentials")
	}
}

func TestMega() bool {
	cmd := exec.Command("mega-get", "--help")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return err == nil
}
