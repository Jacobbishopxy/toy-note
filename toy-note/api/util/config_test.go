package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDevConfig(t *testing.T) {
	cfg, err := LoadConfig(false, "../../env")
	require.NoError(t, err)

	require.Equal(t, cfg.PG_HOST, "localhost")
	require.Equal(t, cfg.PG_PORT, 5432)
	require.Equal(t, cfg.PG_USER, "root")
	require.Equal(t, cfg.PG_PASS, "secret")
	require.Equal(t, cfg.PG_DB, "dev")
	require.Equal(t, cfg.MONGO_HOST, "localhost")
	require.Equal(t, cfg.MONGO_PORT, 27017)
	require.Equal(t, cfg.MONGO_USER, "root")
	require.Equal(t, cfg.MONGO_PASS, "secret")
	require.Equal(t, cfg.MONGO_DB, "dev")
}

func TestProdConfig(t *testing.T) {
	cfg, err := LoadConfig(true, "../../env")
	require.NoError(t, err)

	require.Equal(t, cfg.PG_HOST, "127.0.0.1")
	require.Equal(t, cfg.PG_PORT, 5432)
	require.Equal(t, cfg.PG_USER, "root")
	require.Equal(t, cfg.PG_PASS, "secret")
	require.Equal(t, cfg.PG_DB, "dev")
	require.Equal(t, cfg.MONGO_HOST, "127.0.0.1")
	require.Equal(t, cfg.MONGO_PORT, 27017)
	require.Equal(t, cfg.MONGO_USER, "root")
	require.Equal(t, cfg.MONGO_PASS, "secret")
	require.Equal(t, cfg.MONGO_DB, "dev")
}
