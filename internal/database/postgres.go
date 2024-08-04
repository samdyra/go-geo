package database

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/samdyra/go-geo/internal/config"
)

func NewDB(cfg *config.Config) *sqlx.DB {
    dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

    db, err := sqlx.Connect("postgres", dsn)
    if err != nil {
        log.Fatalf("Could not connect to database: %v", err)
    }

    return db
}