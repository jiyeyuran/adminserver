package server

import (
	"github.com/gin-gonic/gin"
	"github.com/gocraft/dbr/v2"
	"jhmeeting.com/adminserver/app"
	"jhmeeting.com/adminserver/db"
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
	var param struct {
		RoomName string
		Range struct {
			StartTime db.NullTime
			EndTime db.NullTime
		}
		Page uint64
		PerPage uint64
	}

	if c.BindJSON(&param) != nil {
		return
	}

	selector := db.NewSelector(s.DB())

	if len(param.RoomName) > 0 {
		selector.Conditions = append(selector.Conditions, []db.Condition{
			Col: "room_name",
			Cmp: "eq",
			Val: param.RoomName,
		})
	}	
	if param.Range.StartTime.Valid {
		selector.Conditions = append(selector.Conditions, []db.Condition{
			Col: "ctime",
			Cmp: "gte",
			Val: param.Range.StartTime,
		})
	}
	if param.Range.EndTime.Valid {
		selector.Conditions = append(selector.Conditions, []db.Condition{
			Col: "ctime",
			Cmp: "lte",
			Val: param.Range.EndTime,
		})
	}
	
	confereces := []app.ConferenceInfo{}
	result, err := selector.Paginate(param.Page, param.PerPage).LoadPage(&confereces)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.JSON(200, result)
}

//Action 会议室事件
func (s ConferenceServer) Action(c *gin.Context) {
	// TODO: 监听事件
	
}
