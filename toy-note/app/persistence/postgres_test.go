package persistence

import (
	"testing"
	"toy-note/logger"

	"github.com/stretchr/testify/require"
)

const logPath = "test.log"
const sqlUri = "host=localhost user=root password=secret dbname=dev port=5432 sslmode=disable TimeZone=Asia/Shanghai"

func TestConnectionAndDataMigration(t *testing.T) {
	if err := logger.Init("debug", logPath, true); err != nil {
		panic(err)
	} else {
		defer logger.TNLogger.Sync()
	}

	r, err := NewPgRepository(logger.TNLogger, sqlUri)

	require.NoError(t, err)

	err = r.AutoMigrate()

	require.NoError(t, err)
}
