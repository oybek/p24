package database

import (
	"database/sql"
	"fmt"
	"log"

	"embed"

	_ "github.com/lib/pq"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

type Database struct {
	Conn *sql.DB
}

func Initialize(cfg Config) (Database, error) {
	db := Database{}
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Pass, cfg.Name)
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return db, err
	}

	db.Conn = conn
	err = db.Conn.Ping()
	if err != nil {
		return db, err
	}
	log.Println("Database connection established")
	return db, nil
}

//go:embed migrations/*.sql
var fs embed.FS

// Migrate - runs migrations against db
func Migrate(cfg Config) {
	driver, err := iofs.New(fs, "migrations")
	if err != nil {
		log.Fatalf("Migration driver initialization error: %s", err)
	}

	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Name,
	)
	m, err := migrate.NewWithSourceInstance("iofs", driver, dbURL)
	if err != nil {
		log.Fatalf("Migration initialization error: %s", err)
	}

	err = m.Up()
	if err != nil {
		log.Printf("Migration result: %s\n", err)
	}
}
