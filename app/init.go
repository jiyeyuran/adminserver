package app

import (
	"github.com/gocraft/dbr/v2"
	"jhmeeting.com/adminserver/db"
)

var DBTables = map[string]interface{}{
	UserTableName:        User{},
	RoomTableName:        RoomInfo{},
	ConfereneceTableName: ConferenceInfo{},
	RecordTableName:      RecordInfo{},
}

func InitSqlDB(session *dbr.Session) {
	for table, schema := range DBTables {
		if err := db.CreateTable(session, table, schema); err != nil {
			panic(err)
		}
	}
}
