package database

import (
	"context"
	"io"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TestGenericContainer is a test that uses the GenericContainer APIs to run a PostgreSQL container.
// The database is initialized with the SQL files in the testdata directory.
// It demonstrates how to use the GenericContainer APIs to run a container and interact with it:
// - Copy files from the host to the container.
// - Use wait strategies to wait for the container to be ready.
// - Use lifecycle hooks to run code before and after the different events of the container lifecycle.
// - Use the GenericContainer APIs to copy files from the container to the host.
func TestGenericContainer(t *testing.T) {
	var (
		dbName         = "pokemon_battles"
		dbPwd          = "postgres"
		dbUser         = "postgres"
		initSQLFile    = filepath.Join("testdata", "00-schema.sql")
		insertsSQLFile = filepath.Join("testdata", "01-inserts.sql")
	)

	req := testcontainers.ContainerRequest{
		Image: "postgres:latest",
		Env: map[string]string{
			"POSTGRES_USER":     dbUser,
			"POSTGRES_PASSWORD": dbPwd,
			"POSTGRES_DB":       dbName,
		},
		ExposedPorts: []string{"5432/tcp"},
		Cmd:          []string{"postgres", "-c", "fsync=off"},
		WaitingFor:   wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      initSQLFile,
				ContainerFilePath: "/docker-entrypoint-initdb.d/00-schema.sql",
				FileMode:          0o755,
			},
			{
				HostFilePath:      insertsSQLFile,
				ContainerFilePath: "/docker-entrypoint-initdb.d/01-inserts.sql",
				FileMode:          0o755,
			},
		},
		LifecycleHooks: []testcontainers.ContainerLifecycleHooks{
			{
				PreCreates: []testcontainers.ContainerRequestHook{
					func(ctx context.Context, req testcontainers.ContainerRequest) error {
						t.Log("🔧 PreCreates 1")
						return nil
					},
					func(ctx context.Context, req testcontainers.ContainerRequest) error {
						t.Log("🔧 PreCreates 2")
						return nil
					},
				},
				PreStarts: []testcontainers.ContainerHook{
					func(ctx context.Context, container testcontainers.Container) error {
						t.Log("👋 PreStarts 1")
						return nil
					},
					func(ctx context.Context, container testcontainers.Container) error {
						t.Log("👋 PreStarts 2")
						return nil
					},
				},
				PreTerminates: []testcontainers.ContainerHook{
					func(ctx context.Context, container testcontainers.Container) error {
						t.Log("💀 PreTerminates 1")
						return nil
					},
					func(ctx context.Context, container testcontainers.Container) error {
						t.Log("💀 PreTerminates 2")
						return nil
					},
				},
			},
		},
	}

	genericContainerReq := testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	}

	dbContainer, err := testcontainers.GenericContainer(context.Background(), genericContainerReq)

	testcontainers.CleanupContainer(t, dbContainer)
	require.NoError(t, err)

	dbHost, err := dbContainer.Host(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, dbHost)

	dbPort, err := dbContainer.MappedPort(context.Background(), "5432/tcp")
	require.NoError(t, err)
	require.NotEmpty(t, dbPort.Port())

	t.Run("copy-from-container", func(t *testing.T) {
		rc, err := dbContainer.CopyFileFromContainer(
			context.Background(), "/docker-entrypoint-initdb.d/00-schema.sql",
		)
		require.NoError(t, err)

		content, err := io.ReadAll(rc)
		require.NoError(t, err)

		require.Contains(t, string(content), `CREATE TABLE pokemons (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL,
    hp INT NOT NULL,
    attack INT NOT NULL,
    defense INT NOT NULL
);`)
	})
}
