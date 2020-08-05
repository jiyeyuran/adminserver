package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
		RoomName string `json:"room_name,omitempty"`
		Range    struct {
			StartTime db.NullTime `json:"start_time,omitempty"`
			EndTime   db.NullTime `json:"end_time,omitempty"`
		} `json:"range,omitempty"`
		Page    uint64 `json:"page,omitempty"`
		PerPage uint64 `json:"per_page,omitempty"`
	}

	if c.BindJSON(&param) != nil {
		return
	}

	selector := db.NewSelector(s.DB())

	selector.Conditions = append(selector.Conditions, db.Condition{
		Col: "uid",
		Cmp: "eq",
		Val: c.GetInt64("uid"),
	})

	if len(param.RoomName) > 0 {
		selector.Conditions = append(selector.Conditions, db.Condition{
			Col: "room_name",
			Cmp: "eq",
			Val: param.RoomName,
		})
	}
	if param.Range.StartTime.Valid {
		selector.Conditions = append(selector.Conditions, db.Condition{
			Col: "ctime",
			Cmp: "gte",
			Val: param.Range.StartTime,
		})
	}
	if param.Range.EndTime.Valid {
		selector.Conditions = append(selector.Conditions, db.Condition{
			Col: "ctime",
			Cmp: "lte",
			Val: param.Range.EndTime,
		})
	}

	selector.Orders = []db.Order{
		{Col: "id"},
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
	req := ActionRequest{}
	if c.BindJSON(&req) != nil {
		return
	}

	switch req.Action {
	case MUC_ROOM_PRE_CREATE:
	case MUC_ROOM_CREATED:
	case MUC_OCCUPANT_PRE_JOIN:
	case MUC_OCCUPANT_JOINED:
	case MUC_OCCUPANT_LEFT:
	case MUC_ROOM_DESTROYED:
	case MUC_ROOM_SECRET:
	case MUC_ROOM_INFO:
		var roomInfo app.RoomInfo
		err := s.DB().Select("*").From("room").Where("room_name=?", req.Room).LoadOne(&roomInfo)
		if err != nil {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}
		c.JSON(200, roomInfo)
	case MUC_ROOM_RECORDING_START:
	case MUC_ROOM_RECORDING_STOP:
		// TODO: 保存req.RecordingFile
	}
}
