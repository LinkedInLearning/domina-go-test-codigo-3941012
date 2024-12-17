package server

import (
	"context"
	_ "embed"
	"path/filepath"
	"pokemon-battle/internal/database"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func MustNewWithDatabase() database.Service {
	baseDir := filepath.Join("..", "database", "testdata")

	var (
		dbName         = "pokemon_battles"
		dbPwd          = "postgres"
		dbUser         = "postgres"
		dbSchema       = "public"
		initSQLFile    = filepath.Join(baseDir, "00-schema.sql")
		insertsSQLFile = filepath.Join(baseDir, "01-inserts.sql")
	)

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

	return database.NewService(dbUser, dbPwd, dbHost, dbPort.Port(), dbName, dbSchema)
}
