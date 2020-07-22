package app

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/gocraft/dbr/v2"
	"github.com/spf13/viper"
)

type App struct {
	config   AppConfig
	redisCli redis.UniversalClient
	db       *dbr.Session
}

type AppConfig struct {
	Port  int         `json:"port,omitempty"`
	API   APIConfig   `json:"api,omitempty"`
	Redis RedisConfig `json:"redis,omitempty"`
	DB    DBConfig    `json:"db,omitempty"`
}

type APIConfig struct {
	URL   string `json:"url,omitempty"`
	Token string `json:"token,omitempty"`
}

type DBConfig struct {
	Driver   string `json:"driver,omitempty"`
	DSN      string `json:"dsn,omitempty"`
	Timezone string `json:"timezone,omitempty"`
}

type RedisConfig struct {
	Addr     []string `json:"addr,omitempty"`
	Password string   `json:"password,omitempty"`
	DB       int      `json:"db,omitempty"`
}

func NewApp(configLocations ...string) *App {
	viper.SetConfigName("config") // name of config file (without extension)

	for _, configLocation := range configLocations {
		viper.AddConfigPath(configLocation)
	}

	// 默认配置目录
	viper.AddConfigPath("/config")
	viper.AddConfigPath(".")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	appConfig := AppConfig{
		Port: 8004,
	}

	if err := viper.Unmarshal(&appConfig); err != nil {
		panic(err)
	}

	if gin.Mode() != gin.TestMode {
		data, _ := json.MarshalIndent(viper.AllSettings(), "", "\t")

		log.Printf("config: %s", data)
	}

	if len(appConfig.DB.Driver) == 0 {
		if parts := strings.SplitN(appConfig.DB.DSN, "://", 2); len(parts) == 2 {
			appConfig.DB.Driver = parts[0]
		}
	}

	time.Local = time.UTC
	debug := gin.Mode() == gin.DebugMode

	if err := initDB(appConfig.DB, NewDBEventReceiver(debug)); err != nil {
		panic(err)
	}

	return &App{
		config:   appConfig,
		redisCli: newRedis(appConfig.Redis),
		db:       NewDB(appConfig.DB, debug),
	}
}

func (app App) Config() AppConfig {
	return app.config
}

func (app App) RedisCli() redis.UniversalClient {
	return app.redisCli
}

func (app App) DB() *dbr.Session {
	return app.db
}

func newRedis(config RedisConfig) redis.UniversalClient {
	redisCli := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})

	if err := redisCli.Ping().Err(); err != nil {
		panic(err)
	}

	return redisCli
}
