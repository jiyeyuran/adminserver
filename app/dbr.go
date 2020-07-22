package app

import (
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr/v2"
	"github.com/pkg/errors"
)

type localDialect struct {
	dbr.Dialect
	loc *time.Location
}

func newLocalDialect(parentDialect dbr.Dialect, loc *time.Location) dbr.Dialect {
	return &localDialect{
		Dialect: parentDialect,
		loc:     loc,
	}
}

func (d localDialect) EncodeTime(t time.Time) string {
	return `'` + t.In(d.loc).Format("2006-01-02 15:04:05") + `'`
}

func loadLocation(dbConfig DBConfig) (loc *time.Location, err error) {
	loc = time.Local

	if len(dbConfig.Timezone) > 0 {
		loc, err = time.LoadLocation(dbConfig.Timezone)
	}

	return
}

func initDB(dbConfig DBConfig, eventReceiver dbr.EventReceiver) (err error) {
	switch dbConfig.Driver {
	case "mysql", "postgres":
		cfg, err := mysql.ParseDSN(dbConfig.DSN)
		if err != nil {
			return errors.WithStack(err)
		}

		dbName := cfg.DBName
		cfg.DBName = ""

		conn, err := dbr.Open(dbConfig.Driver, cfg.FormatDSN(), eventReceiver)
		if err != nil {
			return errors.WithStack(err)
		}
		defer conn.Close()

		_, err = conn.Exec("CREATE DATABASE IF NOT EXISTS " + conn.QuoteIdent(dbName))
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return
}

func NewDB(dbConfig DBConfig, debug bool) (session *dbr.Session) {
	eventReceiver := NewDBEventReceiver(debug)
	loc, err := loadLocation(dbConfig)
	if err != nil {
		panic(err)
	}

	if err := initDB(dbConfig, eventReceiver); err != nil {
		panic(err)
	}

	conn, err := dbr.Open(dbConfig.Driver, dbConfig.DSN, eventReceiver)
	if err != nil {
		panic(err)
	}
	conn.Dialect = newLocalDialect(conn.Dialect, loc)

	return conn.NewSession(nil)
}

type DBEventReceiver struct {
	dbr.NullEventReceiver
	debug bool
}

func NewDBEventReceiver(debug bool) dbr.EventReceiver {
	return &DBEventReceiver{debug: debug}
}

// Event receives a simple notification when various events occur.
func (n *DBEventReceiver) Event(eventName string) {
	if n.debug {
		log.Printf("event: %s", eventName)
	}
}

// EventKv receives a notification when various events occur along with
// optional key/value data.
func (n *DBEventReceiver) EventKv(eventName string, kvs map[string]string) {
	if n.debug {
		log.Printf("event: %s, sql: %v", eventName, kvs["sql"])
	}
}

// EventErr receives a notification of an error if one occurs.
func (n *DBEventReceiver) EventErr(eventName string, err error) error {
	log.Printf("event: %s, err: %s", eventName, err)
	return err
}

// EventErrKv receives a notification of an error if one occurs along with
// optional key/value data.
func (n *DBEventReceiver) EventErrKv(eventName string, err error, kvs map[string]string) error {
	log.Printf("event: %s, err: %s, sql: %v", eventName, err, kvs["sql"])
	return err
}

// Timing receives the time an event took to happen.
func (n *DBEventReceiver) Timing(eventName string, nanoseconds int64) {
	if n.debug {
		log.Printf("event: %s, timing: %s", eventName, time.Duration(nanoseconds))
	}
}

// TimingKv receives the time an event took to happen along with optional key/value data.
func (n *DBEventReceiver) TimingKv(eventName string, nanoseconds int64, kvs map[string]string) {
	if n.debug {
		log.Printf("event: %s, timing: %s, sql: %v", eventName, time.Duration(nanoseconds), kvs["sql"])
	}
}
