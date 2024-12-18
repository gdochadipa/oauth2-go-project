package database

import (
	"database/sql"
	"fmt"

	"github.com/gdochadipa/oauth2-go-project/pkg/configs",
	_ "github.com/lib/pq"
)

func NewPostgressConnection(cfg configs.DatabaseConfig) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)

	db, err = sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	// Test the connection
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}