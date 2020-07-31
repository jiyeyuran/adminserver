package server

import (
	"github.com/gin-gonic/gin"
	"jhmeeting.com/adminserver/app"
)

type RoomServer struct {
	*app.App
}

func NewRoomServer(app *app.App) *RoomServer {
	return &RoomServer{
		App: app,
	}
}

func (s RoomServer) Info(c *gin.Context) {

}

func (s RoomServer) Create(c *gin.Context) {

}

func (s RoomServer) Delete(c *gin.Context) {

}

func (s RoomServer) Modify(c *gin.Context) {

}

func (s RoomServer) List(c *gin.Context) {

}

func (s RoomServer) Token(c *gin.Context) {

}
