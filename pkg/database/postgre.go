package database

import (
	"database/sql"
	"fmt"

	"github.com/gdochadipa/oauth2-go-project/pkg/configs"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgressConnection(cfg configs.DatabaseConfig) (*sql.DB, error) {

	connStr := fmt.Sprintf("user='%s' password=%s host=%s dbname='%s'", cfg.User, cfg.Password, cfg.Host, cfg.Name)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	// Test the connection
	if err = db.Ping(); err != nil {
		return nil, err
	}

	fmt.Println("Successfully connected")
	return db, nil
}

func GormDB(cfg configs.DatabaseConfig) (*gorm.DB, error) {
	connStr := fmt.Sprintf("user='%s' password=%s host=%s dbname='%s'", cfg.User, cfg.Password, cfg.Host, cfg.Name)
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, err

}
