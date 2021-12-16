package persistence

import (
	"context"
	"os"
	"testing"
	"time"
	"toy-note/logger"

	"github.com/stretchr/testify/require"
)

const mongoUri = "mongodb://root:secret@localhost:27017"

func newMongoRepo() (MongoRepository, error) {
	if err := logger.Init("debug", logPath, true); err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return NewMongoRepository(ctx, logger.TNLogger, mongoUri)
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
