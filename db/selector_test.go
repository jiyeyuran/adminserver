package db

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/gocraft/dbr/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

type People struct {
	Id        int       `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email,omitempty"`
	Location  Location  `json:"location,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

func (loc Location) Value() (driver.Value, error) {
	data, _ := json.Marshal(loc)

	return string(data), nil
}

func (loc *Location) Scan(src interface{}) error {
	var source []byte
	switch src.(type) {
	case string:
		source = []byte(src.(string))
	case []byte:
		source = src.([]byte)
	default:
		return errors.New("Incompatible type for Location")
	}

	return json.Unmarshal(source, loc)
}

type Location struct {
	Log int `json:"log,omitempty"`
	Lat int `json:"lat,omitempty"`
}

func createSession(dbConfig Config) *dbr.Session {
	return NewSQLDB(dbConfig, true)
}

func reset(t *testing.T, sess *dbr.Session) {
	err := CreateTable(sess, "people", People{})
	require.NoError(t, err)
}

func TestSelector(t *testing.T) {
	session := createSession(Config{
		Driver: "sqlite3",
		DSN:    ":memory:",
	})
	reset(t, session)

	n := 10

	for i := 0; i < n; i++ {
		p := People{
			Id:        i + 1,
			Name:      fmt.Sprintf("aaa_%d", i),
			Email:     fmt.Sprintf("email_%d", i),
			Location:  Location{Log: 1, Lat: i},
			CreatedAt: time.Now(),
		}
		_, err := session.InsertInto("people").
			Columns("id", "name", "email", "location", "created_at").Record(p).Exec()
		require.NoError(t, err)
	}

	items := []People{}
	sel := NewSelector(session)

	result, err := sel.From("people").LoadPage(&items)
	require.NoError(t, err)
	require.EqualValues(t, n, result.Count)
	require.Equal(t, items, result.Items)
	require.Len(t, items, n)

	items = []People{}
	sel = NewSelector(session)
	result, err = sel.From("people").Paginate(0, 5).LoadPage(&items)
	require.NoError(t, err)
	require.EqualValues(t, n, result.Count)
	require.Equal(t, items, result.Items)
	require.Len(t, result.Items, 5)

	data, _ := json.Marshal(result)
	t.Log(string(data))
}
