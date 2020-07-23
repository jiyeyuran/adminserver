package db

import (
	"testing"

	"github.com/gocraft/dbr/v2"
	"github.com/gocraft/dbr/v2/dialect"
	"github.com/stretchr/testify/require"
)

type TestTable struct {
	A string
	B int
	C struct{}
	D dbr.NullTime
	E *string
	F string `sql:"type:text length:15 default:'aa'"`
}

type TestTable2 struct {
	ID int
	A  string `sql:"index:a_b,unique"`
	B  string `sql:"index:a_b"`
	C  string `sql:"index:c"`
}

var sqlite3Session = createSession(Config{
	Driver: "sqlite3",
	DSN:    ":memory:",
})

func TestCreateTable(t *testing.T) {
	err := CreateTable(sqlite3Session, "test_table", TestTable{})
	require.NoError(t, err)

	err = CreateTable(sqlite3Session, "test_table", TestTable{})
	require.NoError(t, err)

	err = CreateTable(sqlite3Session, "test_table2", TestTable2{})
	require.NoError(t, err)
}

func TestSplitDSN(t *testing.T) {
	dbName, dsn, err := splitDSN("mysql", "user:password@tcp(localhost:5555)/dbname?tls=skip-verify&autocommit=true")
	require.NoError(t, err)
	require.Equal(t, "dbname", dbName)
	require.Equal(t, "user:password@tcp(localhost:5555)/?tls=skip-verify&autocommit=true", dsn)

	dbName, dsn, err = splitDSN("postgres", "postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full")
	require.NoError(t, err)
	require.Equal(t, "pqgotest", dbName)
	require.Equal(t, "postgres://pqgotest:password@localhost/?sslmode=verify-full", dsn)
}

func TestSchema2CreateTableSQL(t *testing.T) {
	testCases := []struct {
		Driver  string
		Dialect dbr.Dialect
		Want    string
	}{
		{
			Driver:  "sqlite3",
			Dialect: dialect.SQLite3,
			Want: `CREATE TABLE IF NOT EXISTS "test_table"(
"a" VARCHAR(255) NOT NULL DEFAULT '',
"b" INT(11) NOT NULL DEFAULT '0',
"c" TEXT NOT NULL,
"d" DATETIME NOT NULL,
"e" VARCHAR(255) NULL,
"f" TEXT(15) NOT NULL DEFAULT 'aa'
)`,
		},
		{
			Driver:  "mysql",
			Dialect: dialect.MySQL,
			Want: "CREATE TABLE IF NOT EXISTS `test_table`(\n" +
				"`a` VARCHAR(255) NOT NULL DEFAULT '',\n" +
				"`b` INT(11) NOT NULL DEFAULT '0',\n" +
				"`c` TEXT NOT NULL,\n" +
				"`d` DATETIME NOT NULL,\n" +
				"`e` VARCHAR(255) NULL,\n" +
				"`f` TEXT(15) NOT NULL DEFAULT 'aa'\n)",
		},
		{
			Driver:  "postgres",
			Dialect: dialect.PostgreSQL,
			Want: `CREATE TABLE IF NOT EXISTS "test_table"(
"a" VARCHAR(255) NOT NULL DEFAULT '',
"b" INT(11) NOT NULL DEFAULT '0',
"c" TEXT NOT NULL,
"d" TIMESTAMP NOT NULL,
"e" VARCHAR(255) NULL,
"f" TEXT(15) NOT NULL DEFAULT 'aa'
)`,
		},
	}

	for _, c := range testCases {
		sqlstr := schema2CreateTableSQL(c.Dialect, "test_table", TestTable{})
		require.Equal(t, c.Want, sqlstr, c.Driver)
	}
}

func TestListTableIndexSQL(t *testing.T) {
	type Test struct {
		A string `sql:"index:a_b,unique"`
		B string `sql:"index:a_b"`
	}
	sqls := listTableIndexSQL(dialect.MySQL, "test", Test{})
	require.Len(t, sqls, 1)

	require.Equal(t, "CREATE UNIQUE INDEX `a_b` ON `test` (`a`,`b`)", sqls[0])
}
