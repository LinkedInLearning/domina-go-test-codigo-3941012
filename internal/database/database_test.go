package database

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	srv := MustNewWithDatabase(t)
	require.NotNil(t, srv)
}

func TestMustDB(t *testing.T) {
	srv := service{}

	require.Panics(t, func() {
		srv.MustDB()
	})
}

func TestHealth(t *testing.T) {
	srv := MustNewWithDatabase(t)

	stats := srv.Health()

	require.Equal(t, "up", stats["status"])
	require.Empty(t, stats["error"])
	require.Equal(t, "It's healthy", stats["message"])
}

func TestClose(t *testing.T) {
	srv := MustNewWithDatabase(t)
	require.NoError(t, srv.Close())
}
