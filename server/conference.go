package server

import (
	"errors"
	"log"
	"net/http"
	"time"

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
	var param struct {
		ID int64 `json:"id,omitempty"`
	}
	if c.BindJSON(&param) != nil {
		return
	}
	info := app.ConferenceInfo{}
	err := s.DB().Select("*").From("conference").
		Where("id=? and uid=?", param.ID, c.GetInt64(app.UserID)).LoadOneContext(c, &info)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	c.JSON(http.StatusOK, info)
}

// Runing 获取正在进行的会议室列表
func (s ConferenceServer) Runing(c *gin.Context) {
	items := []app.ConferenceInfo{}
	result, err := db.NewSelector(s.DB()).From("conference").Where(
		dbr.Eq("uid", c.GetInt64(app.UserID)),
		dbr.Eq("etime", nil),
	).LoadPage(&items)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, result)
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
		Val: c.GetInt64(app.UserID),
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

	if len(req.Room) == 0 {
		c.AbortWithError(http.StatusBadRequest, errors.New("room is invalid"))
		return
	}

	switch req.Action {
	case MUC_ROOM_INFO:
		var roomInfo app.RoomInfo
		err := s.DB().Select("*").From("room").Where("room_name=?", req.Room).LoadOne(&roomInfo)
		if err != nil {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}
		c.JSON(200, roomInfo)

	case MUC_ROOM_PRE_CREATE:
		uid, _ := s.DB().Select("id").From("room").Where("room_name=?", req.Room).ReturnInt64()
		if uid == 0 {
			c.AbortWithError(http.StatusNotFound, errors.New("房间不存在"))
			return
		}
		confereceInfo := app.ConferenceInfo{
			Uid:        uid,
			RoomName:   req.Room,
			ApiEnabled: req.ApiEnabled,
			Ctime:      time.Now(),
		}
		_, err := s.DB().InsertInto("conference").
			Columns("uid", "room_name", "api_enabled", "ctime").
			Record(&confereceInfo).ExecContext(c)
		if err == nil {
			c.AbortWithError(http.StatusNotFound, errors.New("房间不存在"))
			return
		}
		c.JSON(200, confereceInfo)

	case MUC_ROOM_CREATED:
		log.Printf("room: %s created", req.Room)

	case MUC_OCCUPANT_PRE_JOIN:
		participantLimits, _ := s.DB().Select("participant_limits").From("room").Where("id=?", req.ConferenceId).ReturnInt64()
		if participantLimits > 0 && req.Participants >= int(participantLimits) {
			c.AbortWithError(http.StatusServiceUnavailable, errors.New("会议室人数已达上限"))
			return
		}

	case MUC_OCCUPANT_JOINED:
		s.DB().Update("conference").Set("participants", req.Participants).Where("id=?", req.ConferenceId).ExecContext(c)
		s.DB().Update("conference").Set("max_participants", req.Participants).
			Where("id=? and max_participants<?", req.ConferenceId, req.Participants).ExecContext(c)
		// TODO: 数据库记录参会者

	case MUC_OCCUPANT_LEFT:
		s.DB().Update("conference").Set("participants", req.Participants).Where("id=?", req.ConferenceId).ExecContext(c)
		// TODO: 数据库更新参会者

	case MUC_ROOM_DESTROYED:
		s.DB().Update("conference").
			Set("etime", time.Now()).
			Set("is_recording", false).
			Set("is_streaming", false).
			Where("id=?", req.ConferenceId).ExecContext(c)

	case MUC_ROOM_SECRET:
		s.DB().Update("conference").Set("lock_password", req.Secret).Where("id=?", req.ConferenceId).ExecContext(c)
	case MUC_ROOM_RECORDING_START:
		if recording := req.Recording; recording != nil {
			isStreaming := len(recording.Streaming) > 0

			s.DB().Update("conference").
				Set("is_recording", !isStreaming).
				Set("is_streaming", isStreaming).
				Where("id=?", req.ConferenceId).ExecContext(c)

			if isStreaming {
				uid, _ := s.DB().Select("uid").From("room").Where("room_name=?", req.Room).ReturnInt64()
				recordInfo := app.RecordInfo{
					ConferenceUid: uid,
					ConferenceId:  req.ConferenceId,
					RoomName:      req.Room,
					StreamingUrl:  recording.Streaming,
					Ctime:         time.Now(),
				}
				s.DB().InsertInto(app.RecordTableName).
					Columns(app.RecordConferenceUidCol, app.RecordConferenceIdCol, app.RecordRoomNameCol, app.RecordStreamUrlCol).
					Record(&recordInfo).ExecContext(c)
			}
		}

	case MUC_ROOM_RECORDING_STOP:
		if recording := req.Recording; recording != nil {
			s.DB().Update("conference").
				Set("is_recording", false).
				Set("is_streaming", false).
				Where("id=?", req.ConferenceId).ExecContext(c)

			if len(recording.Streaming) > 0 {
				s.DB().Update(app.RecordTableName).
					Set(app.RecordDurationCol, recording.Duration).
					Set(app.RecordSizeCol, recording.Size).
					Where("conference_id=? and streaming_url=? and duration=0", req.ConferenceId, recording.Streaming).
					ExecContext(c)
			} else {
				uid, _ := s.DB().Select("uid").From("room").Where("room_name=?", req.Room).ReturnInt64()
				recordInfo := app.RecordInfo{
					ConferenceUid: uid,
					ConferenceId:  req.ConferenceId,
					RoomName:      req.Room,
					Duration:      recording.Duration,
					Size:          recording.Size,
					DownloadUrl:   recording.ObjectKey,
					StreamingUrl:  recording.Streaming,
					Ctime:         time.Now(),
				}
				s.DB().InsertInto(app.RecordTableName).
					Columns(app.RecordConferenceUidCol, app.RecordConferenceIdCol, app.RecordRoomNameCol,
						app.RecordDurationCol, app.RecordSizeCol, app.RecordDownUrlCol, app.RecordStreamUrlCol, app.RecordCtimeCol).
					Record(&recordInfo).ExecContext(c)
			}
		}
	}
}
