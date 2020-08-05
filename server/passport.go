package server

import (
    "errors"
    "net/http"
    "time"

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
        CaptchaId   string `json:"captcha_id,omitempty"`
        CaptchaCode string `json:"captcha_code,omitempty"`
    }
    if c.BindJSON(&param) != nil {
        return
    }

    /*if !captcha.VerifyString(param.CaptchaId, param.CaptchaCode) {
    	c.AbortWithError(http.StatusBadRequest, errors.New("图片验证码错误"))
    	return
    }*/

    var err error
    param.Password, err = app.HashPassword(param.Password)
    if err != nil {
        c.AbortWithError(http.StatusInternalServerError, err)
        return
    }

    user := app.User{}
    err = s.DB().Select("*").From("users").
        Where("name=?", param.Name).LoadOneContext(c, &user)
    if user.Name == param.Name {
        c.AbortWithError(http.StatusBadRequest, errors.New("用户名已存在！"))
        return
    }

    param.Ctime = time.Now()
    ctx := c.Request.Context()

    _, err = s.DB().InsertInto("users").
        Columns("name", "password", "ctime").
        Record(&param.User).ExecContext(ctx)
    if err != nil {
        c.AbortWithError(http.StatusInternalServerError, err)
        return
    }
    token := s.CreateToken(param.Id)
    param.Password = ""
    c.SetCookie(app.CookieName, token, 0, "/", "", true, true)
    c.AbortWithStatusJSON(http.StatusOK, gin.H{
        "id": param.Id,
    })
}

func (s PassportServer) Login(c *gin.Context) {
    var param struct {
        app.User
        CaptchaId   string `json:"captcha_id,omitempty"`
        CaptchaCode string `json:"captcha_code,omitempty"`
    }
    if c.BindJSON(&param) != nil {
        return
    }

    /*if !captcha.VerifyString(param.CaptchaId, param.CaptchaCode) {
    	c.AbortWithError(http.StatusBadRequest, errors.New("图片验证码错误"))
    	return
    }*/

    errMsg := "用户名或密码错误"
    user := app.User{}
    err := s.DB().Select("*").From("users").
        Where("name=?", param.Name).LoadOneContext(c, &user)
    if err != nil {
        c.AbortWithError(http.StatusBadRequest, errors.New(errMsg))
        return
    }

    pass := app.CheckPasswordHash(param.Password, user.Password)

    if !pass {
        c.AbortWithError(http.StatusBadRequest, errors.New(errMsg))
        return
    }
    token := s.CreateToken(user.Id)
    c.SetCookie(app.CookieName, token, 0, "/", "", true, true)
}

func (s PassportServer) Logout(c *gin.Context) {
    c.SetCookie(app.CookieName, "", -1, "/", "", true, true)
}
