package persistence

import (
	"context"
	"os"
	"testing"
	"time"
	"toy-note/logger"

	"github.com/stretchr/testify/require"
)

var mongoConn = MongoConn{
	host: "localhost",
	port: 27017,
	user: "root",
	pass: "secret",
}

func newMongoRepo() (MongoRepository, error) {
	if err := logger.Init("debug", logPath, true); err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return NewMongoRepository(ctx, logger.TNLogger, mongoConn)
}

func TestMongoConnection(t *testing.T) {

	_, err := newMongoRepo()

	require.NoError(t, err)
}

func TestUploadAndDownloadFile(t *testing.T) {

	r, err := newMongoRepo()
	require.NoError(t, err)

	reader, err := os.Open("test.log")
	require.NoError(t, err)

	id, err := r.UploadFile(reader, "test.log.bak")
	require.NoError(t, err)

	_, err = r.DownloadFile(id)
	require.NoError(t, err)
}
