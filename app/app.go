package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/gocraft/dbr/v2"
	"github.com/spf13/viper"
	"jhmeeting.com/adminserver/db"
)

const CookieName = "rtcadmin"

type App struct {
	config     AppConfig
	httpClient *http.Client
	redisCli   redis.UniversalClient
	db         *dbr.Session
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

	if err := db.CreateDatabase(appConfig.DB); err != nil {
		panic(err)
	}

	sqlDB := db.NewSQLDB(appConfig.DB, gin.Mode() == gin.DebugMode)
	InitSqlDB(sqlDB)

	return &App{
		config: appConfig,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		redisCli: newRedis(appConfig.Redis),
		db:       sqlDB,
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

func (app App) CreateToken(claims jwt.StandardClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(app.config.Secret))
	if err != nil {
		panic(err)
	}
	return tokenString
}

func (app App) ParseToken(tokenString string) (claims *jwt.StandardClaims, err error) {
	_, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(app.config.Secret), nil
	})
	return
}

func (app App) HttpClient() *http.Client {
	return app.httpClient
}

func (app App) NewAPIRequest(path string, body interface{}) *http.Request {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	url := strings.TrimSuffix(app.config.API.URL, "/") + path
	if body == nil {
		body = gin.H{}
	}
	var data []byte
	if bytes, ok := body.([]byte); ok {
		data = bytes
	} else {
		bytes, err := json.Marshal(body)
		if err != nil {
			panic(err)
		}
		data = bytes
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", "Bearer "+app.config.API.Token)

	return req
}

func (app App) SendAPIRequest(path string, body interface{}) (data []byte, err error) {
	req := app.NewAPIRequest(path, body)
	resp, err := app.HttpClient().Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if resp.StatusCode >= 400 {
		var errResult struct {
			Message string
		}
		if err = json.Unmarshal(data, &errResult); err != nil {
			return
		}
		err = errors.New(errResult.Message)
	}
	return
}

func (app App) APIRoute(c *gin.Context, apiPath string) {
	remote, _ := url.Parse(app.Config().API.URL)

	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Director = func(req *http.Request) {
		req.Host = remote.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = apiPath

		req.Header.Set("Authorization", "Bearer "+app.config.API.Token)
		req.Header.Set("X-Uid", c.GetString("uid"))
	}
	proxy.ServeHTTP(c.Writer, c.Request)
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
