package dbutil

import (
	"fmt"
	"testing"
	"time"

	"github.com/gocraft/dbr/v2"
	"github.com/gocraft/dbr/v2/dialect"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	"jhmeeting.com/adminserver/app"
)

type People struct {
	Id        int       `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

func createSession(dbConfig app.DBConfig) *dbr.Session {
	return app.NewDB(dbConfig, true)
}

func reset(t *testing.T, sess *dbr.Session) {
	var autoIncrementType string

	switch sess.Dialect {
	case dialect.MySQL, dialect.PostgreSQL:
		autoIncrementType = "SERIAL PRIMARY KEY"
	case dialect.SQLite3:
		autoIncrementType = "INTEGER PRIMARY KEY AUTOINCREMENT"
	}

	for _, v := range []string{
		`DROP TABLE IF EXISTS dbr_people`,
		fmt.Sprintf(`CREATE TABLE dbr_people (
			id %s,
			name varchar(255) NOT NULL,
			email varchar(255),
			created_at datetime NOT NULL
		)`, autoIncrementType),
	} {
		_, err := sess.InsertBySql(v).Exec()
		require.NoError(t, err)
	}
}

func TestSelect(t *testing.T) {
	session := createSession(app.DBConfig{
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
			CreatedAt: time.Now(),
		}
		_, err := session.InsertInto("dbr_people").
			Columns("id", "name", "email", "created_at").Record(p).Exec()
		require.NoError(t, err)
	}

	items := []People{}
	sel := NewSelect(session)

	result, err := sel.From("dbr_people").LoadPage(&items)
	require.NoError(t, err)
	require.EqualValues(t, n, result.Count)
	require.Equal(t, items, result.Items)
	require.Len(t, items, n)

	items = []People{}
	sel = NewSelect(session)
	result, err = sel.From("dbr_people").Paginate(0, 5).LoadPage(&items)

	require.NoError(t, err)
	require.EqualValues(t, n, result.Count)
	require.Equal(t, items, result.Items)
	require.Len(t, result.Items, 5)
}
