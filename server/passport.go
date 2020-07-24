package server

import (
	"errors"
	"time"

	"github.com/dchest/captcha"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"jhmeeting.com/adminserver/app"
)

type PassportServer struct {
	*app.App
}

func NewPassportServer(app *app.App) *PassportServer {
	return &PassportServer{
		App: app,
	}
}

func (s PassportServer) Signup(c *gin.Context) {
	var param struct {
		app.User
		CaptchaId   string
		CaptchaCode string
	}
	if c.BindJSON(&param) != nil {
		return
	}

	if !captcha.VerifyString(param.CaptchaId, param.CaptchaCode) {
		c.AbortWithError(400, errors.New("图片验证码错误"))
		return
	}

	var err error
	param.Password, err = app.HashPassword(param.Password)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	param.Ctime = time.Now()
	ctx := c.Request.Context()

	_, err = s.DB().InsertInto("users").
		Columns("name", "mobile", "password", "ctime").
		Record(param.User).ExecContext(ctx)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	param.Password = ""
	token := s.CreateToken(jwt.MapClaims{
		"id": param.Id,
	})
	c.SetCookie("rtcadmin", token, 0, "/", "", true, true)
	c.JSON(200, param.User)
}

func (s PassportServer) Login(c *gin.Context) {
	var param struct {
		app.User
		CaptchaId   string
		CaptchaCode string
	}
	if c.BindJSON(&param) != nil {
		return
	}

	if !captcha.VerifyString(param.CaptchaId, param.CaptchaCode) {
		c.AbortWithError(400, errors.New("图片验证码错误"))
		return
	}

	password, err := app.HashPassword(param.Password)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	param.Password = ""

	param.Id, _ = s.DB().Select("id").From("users").
		Where("name=? and password=?", param.Name, password).ReturnInt64()

	if param.Id == 0 {
		c.AbortWithError(400, errors.New("用户名或密码错误"))
		return
	}
	token := s.CreateToken(jwt.MapClaims{
		"id": param.Id,
	})
	c.SetCookie("rtcadmin", token, 0, "/", "", true, true)
	c.JSON(200, param.User)
}

func (s PassportServer) Logout(c *gin.Context) {
	c.SetCookie("rtcadmin", "", -1, "/", "", true, true)
}
