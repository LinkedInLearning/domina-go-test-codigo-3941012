package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"pokemon-battle/internal/models"
	"strconv"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

// Service represents a service that interacts with a database.
type Service interface {
	MustDB() *sql.DB

	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error
}

type service struct {
	db *sql.DB
}

var (
	database   = os.Getenv("BLUEPRINT_DB_DATABASE")
	password   = os.Getenv("BLUEPRINT_DB_PASSWORD")
	username   = os.Getenv("BLUEPRINT_DB_USERNAME")
	port       = os.Getenv("BLUEPRINT_DB_PORT")
	host       = os.Getenv("BLUEPRINT_DB_HOST")
	schema     = os.Getenv("BLUEPRINT_DB_SCHEMA")
	dbInstance *service
	sbMu       sync.Mutex
)

const errNilDB = "database connection is nil"

// validateDB checks if the database connection is nil
func validateDB(db *sql.DB) error {
	sbMu.Lock()
	defer sbMu.Unlock()

	if db == nil {
		return errors.New(errNilDB)
	}
	return nil
}

// New returns a new database service.
// If the database service is already initialized, it returns the same instance.
// If the database service is not initialized, it initializes a new one.
// Thread safe.
func New() Service {
	sbMu.Lock()
	defer sbMu.Unlock()

	// Reuse Connection if it's already initialized and healthy
	if dbInstance != nil && dbInstance.db != nil && dbInstance.db.Ping() == nil {
		return dbInstance
	}

	return NewService(username, password, host, port, database, schema)
}

// NewService creates a new database service with the given parameters
// Not thread safe, use New() instead.
func NewService(username, password, host, port, database, schema string) Service {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", username, password, host, port, database, schema)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return &service{
		db: db,
	}
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf("db down: %v", err) // Log the error and terminate the program
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := s.db.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

// MustDB returns the database connection.
// It panics if the database connection is nil.
func (s *service) MustDB() *sql.DB {
	if err := validateDB(s.db); err != nil {
		panic(err)
	}

	return s.db
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", database)
	return s.db.Close()
}

type PokemonCRUDService interface {
	Create(ctx context.Context, obj *models.Pokemon) error
	Delete(ctx context.Context, id int) error
	GetAll(ctx context.Context) ([]models.Pokemon, error)
	GetByID(ctx context.Context, id int) (models.Pokemon, error)
	Update(ctx context.Context, obj models.Pokemon) error
}

type BattleCRUDService interface {
	Create(ctx context.Context, obj *models.Battle) error
	Delete(ctx context.Context, id int) error
	GetAll(ctx context.Context) ([]models.Battle, error)
	GetByID(ctx context.Context, id int) (models.Battle, error)
	Update(ctx context.Context, obj models.Battle) error
}
