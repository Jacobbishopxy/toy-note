package persistence

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
	"toy-note/logger"

	"github.com/stretchr/testify/require"
)

const mongoUri = "mongodb://admin:password@localhost:27017"

func TestMongoConnection(t *testing.T) {
	if err := logger.Init("debug", logPath, true); err != nil {
		panic(err)
	} else {
		defer logger.TNLogger.Sync()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := NewMongoRepository(ctx, logger.TNLogger, mongoUri)

	require.NoError(t, err)
}

func TestUploadAndDownloadFile(t *testing.T) {
	if err := logger.Init("debug", logPath, true); err != nil {
		panic(err)
	} else {
		defer logger.TNLogger.Sync()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	r, err := NewMongoRepository(ctx, logger.TNLogger, mongoUri)

	require.NoError(t, err)

	reader, err := os.Open("test.log")

	require.NoError(t, err)

	id, err := r.UploadFile(reader, "test.log.bak")
	require.NoError(t, err)

	fmt.Println(id)

	data, err := r.DownloadFile("61bb4c3ca9406a0f51fbd5ca")
	require.NoError(t, err)

	fmt.Println(data)
}
