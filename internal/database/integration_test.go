package database

import (
	"context"
	_ "embed"
	"path/filepath"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// MustNewWithDatabase creates a new database service and returns it.
// It's using Testcontainers to run a PostgreSQL container, returning a new Service.
// The database is initialized with the SQL files in the testdata directory.
// Use this function in integration tests to obtain a new database.
func MustNewWithDatabase(t *testing.T) Service {
	var (
		dbName         = "pokemon_battles"
		dbPwd          = "postgres"
		dbUser         = "postgres"
		dbSchema       = "public"
		initSQLFile    = filepath.Join("testdata", "00-schema.sql")
		insertsSQLFile = filepath.Join("testdata", "01-inserts.sql")
	)

	/*
		req := testcontainers.ContainerRequest{
			Image: "postgres:latest",
			Env: map[string]string{
				"POSTGRES_USER":     "postgres",
				"POSTGRES_PASSWORD": "postgres",
				"POSTGRES_DB":       "postgres",
			},
			ExposedPorts: []string{"5432/tcp"},
			Cmd:          []string{"postgres", "-c", "fsync=off"},
		}

		genericContainerReq := testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		}

		dbContainer, err := testcontainers.GenericContainer(ctx, genericContainerReq)
	*/

	dbContainer, err := postgres.Run(
		context.Background(),
		"postgres:latest",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPwd),
		postgres.WithInitScripts(initSQLFile, insertsSQLFile),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	testcontainers.CleanupContainer(t, dbContainer)
	if err != nil {
		panic(err)
	}

	dbHost, err := dbContainer.Host(context.Background())
	if err != nil {
		panic(err)
	}

	dbPort, err := dbContainer.MappedPort(context.Background(), "5432/tcp")
	if err != nil {
		panic(err)
	}

	return NewService(dbUser, dbPwd, dbHost, dbPort.Port(), dbName, dbSchema)
}
