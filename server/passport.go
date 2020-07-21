package server

import (
	"github.com/gin-gonic/gin"
	"jhmeeting.com/adminserver/app"
)

type PassportServer struct {
	*app.App
}

func (s PassportServer) Signup(c *gin.Context) {

}

func (s PassportServer) Login(c *gin.Context) {

}

func (s PassportServer) Logout(c *gin.Context) {

}
