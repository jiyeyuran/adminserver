package app

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/gocraft/dbr/v2"
	"github.com/spf13/viper"
	"jhmeeting.com/adminserver/db"
)

type App struct {
	config   AppConfig
	redisCli redis.UniversalClient
	db       *dbr.Session
}

type AppConfig struct {
	Port   int         `json:"port,omitempty"`
	Secret string      `json:"secret,omitempty"`
	API    APIConfig   `json:"api,omitempty"`
	Redis  RedisConfig `json:"redis,omitempty"`
	DB     db.Config   `json:"db,omitempty"`
}

type APIConfig struct {
	URL   string `json:"url,omitempty"`
	Token string `json:"token,omitempty"`
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

	debug := gin.Mode() == gin.DebugMode

	if err := db.CreateDatabase(appConfig.DB); err != nil {
		panic(err)
	}

	return &App{
		config:   appConfig,
		redisCli: newRedis(appConfig.Redis),
		db:       db.NewSQLDB(appConfig.DB, debug),
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

func (app App) CreateToken(claims jwt.MapClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(app.config.Secret)
	if err != nil {
		panic(err)
	}
	return tokenString
}

func (app App) ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(app.config.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims.(jwt.MapClaims), nil
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
