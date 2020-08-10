package server

import (
	"github.com/gocraft/dbr/v2"
	"net/http"

	"github.com/gin-gonic/gin"
	"jhmeeting.com/adminserver/app"
	"jhmeeting.com/adminserver/db"
)

type RecordServer struct {
	*app.App
}

func NewRecordServer(app *app.App) *RecordServer {
	return &RecordServer{
		App: app,
	}
}

func (s RecordServer) Info(c *gin.Context) {
	var param struct {
		ID int64
	}
	if c.BindJSON(&param) != nil {
		return
	}

	uid := c.GetInt64(app.UserID)

	record := app.RecordInfo{}
	err := s.DB().Select("*").From(app.RecordTableName).
		Where("id=? and conference_uid=?", param.ID, uid).LoadOneContext(c, &record)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	c.JSON(http.StatusOK, record)
}

// 删除
func (s RecordServer) Delete(c *gin.Context) {
	var param struct {
		ID int64
	}
	if c.BindJSON(&param) != nil {
		return
	}
	uid := c.GetInt64(app.UserID)

	_, err := s.DB().
		DeleteFrom(app.RecordTableName).Where("id=? and conference_uid=?", param.ID, uid).ExecContext(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

// 列表查看
func (s RecordServer) List(c *gin.Context) {
	var param db.Pagination
	if c.BindJSON(&param) != nil {
		return
	}

	uid := c.GetInt64(app.UserID)

	records := []app.RecordInfo{}

	result, _ := db.NewSelector(s.DB()).From(app.RecordTableName).
		Where(dbr.Eq(app.RecordConferenceUidCol, uid)).
		Paginate(param.Page, param.PerPage).
		OrderDesc(app.RecordIdCol).
		LoadPage(&records)
	c.JSON(http.StatusOK, result)
}
