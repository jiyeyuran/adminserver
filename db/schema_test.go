package db

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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
