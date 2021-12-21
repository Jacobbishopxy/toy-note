package service

import (
	"testing"
	"toy-note/api/persistence"
	"toy-note/logger"

	"github.com/stretchr/testify/require"
)

const logPath = "test.log"

var sqlConn = persistence.PgConn{
	Host:    "localhost",
	Port:    5432,
	User:    "root",
	Pass:    "secret",
	Db:      "dev",
	Sslmode: "disable",
}

var mongoConn = persistence.MongoConn{
	Host: "localhost",
	Port: 27017,
	User: "root",
	Pass: "secret",
}

func newService() (*ToyNoteService, error) {
	if err := logger.Init("debug", logPath, true); err != nil {
		panic(err)
	}

	return NewToyNoteService(logger.TNLogger, sqlConn, mongoConn)
}

func TestNewService(t *testing.T) {

	_, err := newService()

	require.NoError(t, err)
}
