package persistence

import (
	"testing"
	"toy-note/logger"

	"github.com/stretchr/testify/require"
)

const logPath = "test.log"
const sqlUri = `host=localhost
				user=root
				password=secret
				dbname=dev
				port=5432
				sslmode=disable
				TimeZone=Asia/Shanghai`

func newPgRepo() (PgRepository, error) {
	if err := logger.Init("debug", logPath, true); err != nil {
		panic(err)
	}

	return NewPgRepository(logger.TNLogger, sqlUri)
}

func TestConnectionAndDataMigration(t *testing.T) {

	r, err := newPgRepo()
	require.NoError(t, err)

	err = r.AutoMigrate()
	require.NoError(t, err)
}
