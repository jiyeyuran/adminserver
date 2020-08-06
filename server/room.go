package server

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocraft/dbr/v2"
	"jhmeeting.com/adminserver/app"
	"jhmeeting.com/adminserver/db"
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
	var param struct {
		ID int64
	}
	if c.BindJSON(&param) != nil {
		return
	}
	uid := c.GetInt64("uid")
	room := app.RoomInfo{}
	err := s.DB().Select("*").From("room").
		Where("id=? and uid=?", param.ID, uid).LoadOneContext(c, &room)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, room)
}

func (s RoomServer) Create(c *gin.Context) {
	roomInfo := app.RoomInfo{}
	if c.BindJSON(&roomInfo) != nil {
		return
	}

	roomInfo.Uid = c.GetInt64("uid")
	roomInfo.Ctime = time.Now()

	room := app.RoomInfo{}
	err := s.DB().Select("*").From("room").
		Where("uid=? and name=?", roomInfo.Uid, roomInfo.RoomName).LoadOneContext(c, &room)
	if room.RoomName == roomInfo.RoomName {
		c.AbortWithError(http.StatusBadRequest, errors.New("会议名已存在！"))
		return
	}

	_, err = s.DB().InsertInto("room").
		Columns("uid", "participant_limits", "room_name", "allow_anonymous", "config", "ctime").
		Record(&roomInfo).ExecContext(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, gin.H{
		"id": roomInfo.Id,
	})
}

func (s RoomServer) Delete(c *gin.Context) {
	var param struct {
		ID int64
	}
	if c.BindJSON(&param) != nil {
		return
	}
	uid := c.GetInt64("uid")
	_, err := s.DB().DeleteFrom("room").Where("id=? and uid=?", param.ID, uid).ExecContext(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

func (s RoomServer) Modify(c *gin.Context) {
	roomInfo := app.RoomInfo{}
	if c.BindJSON(&roomInfo) != nil {
		return
	}
	uid := c.GetInt64("uid")
	_, err := s.DB().Update("room").
		Set("participant_limits", roomInfo.ParticipantLimits).
		Set("allow_anonymous", roomInfo.AllowAnonymous).
		Set("config", roomInfo.Config).
		Where("id=? and uid=?", roomInfo.Id, uid).
		ExecContext(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

func (s RoomServer) List(c *gin.Context) {
	var param db.Pagination
	if c.BindJSON(&param) != nil {
		return
	}

	uid := c.GetInt64("uid")
	rooms := []app.RoomInfo{}

	result, _ := db.NewSelector(s.DB()).From("room").
		Where(dbr.Eq("uid", uid)).
		Paginate(param.Page, param.PerPage).
		OrderDesc("id").
		LoadPage(&rooms)
	c.AbortWithStatusJSON(http.StatusOK, result)
}

func (s RoomServer) Token(c *gin.Context) {
	// JSON Body: see RoomTokenRequest
	s.APIRoute(c, "/api/conference/token")
}
