package server

import (
    "github.com/gin-gonic/gin"
    "jhmeeting.com/adminserver/app"
    "jhmeeting.com/adminserver/db"
    "net/http"
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

    record := app.RecordInfo{}
    err := s.DB().Select("*").From(app.RecordTableName).
        Where(app.WhereRecordID, param.ID).LoadOneContext(c, &record)
    if err != nil {
        c.AbortWithError(http.StatusNotFound, err)
        return
    }
    c.AbortWithStatusJSON(http.StatusOK, record)
}

// 删除
func (s RecordServer) Delete(c *gin.Context) {
    var param struct {
        ID int64
    }
    if c.BindJSON(&param) != nil {
        return
    }
    _, err := s.DB().DeleteFrom(app.RecordTableName).Where(app.WhereRecordID, param.ID).ExecContext(c)
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

    records := []app.RecordInfo{}

    result, _ := db.NewSelector(s.DB()).From(app.RecordTableName).
        Paginate(param.Page, param.PerPage).
        OrderDesc(app.RecordIdCol).
        LoadPage(&records)
    c.AbortWithStatusJSON(http.StatusOK, result)
}
