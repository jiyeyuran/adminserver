package server

import (
	"github.com/gin-gonic/gin"
	"jhmeeting.com/adminserver/app"
)

// ConferenceServer 会议室服务
type ConferenceServer struct {
	*app.App
}

func NewConferenceServer(app *app.App) *ConferenceServer {
	return &ConferenceServer{
		App: app,
	}
}

// Info 获取会议室信息
func (s ConferenceServer) Info(c *gin.Context) {
	s.APIRoute(c, "/api/conference/status")
}

// Runing 获取正在进行的会议室列表
func (s ConferenceServer) Runing(c *gin.Context) {
	s.APIRoute(c, "/api/conference/runing")
}

// Dispose 解散会议室
func (s ConferenceServer) Dispose(c *gin.Context) {
	s.APIRoute(c, "/api/conference/dispose")
}

// Lock 锁定会议室
func (s ConferenceServer) Lock(c *gin.Context) {
	s.APIRoute(c, "/api/conference/lock")
}

//Unlock 解锁会议室
func (s ConferenceServer) Unlock(c *gin.Context) {
	s.APIRoute(c, "/api/conference/unlock")
}

//History 会议室历史记录
func (s ConferenceServer) History(c *gin.Context) {
	// TODO: 从数据库获取
}

//Destoryed 会议室解散会议室通知
func (s ConferenceServer) Destroyed(c *gin.Context) {
	
}
