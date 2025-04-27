package database

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/cenkalti/backoff/v5"
	"github.com/gdochadipa/oauth2-go-project/internal/entity"
	"github.com/gdochadipa/oauth2-go-project/internal/util"
	"github.com/gdochadipa/oauth2-go-project/pkg/configs"
	"github.com/golang/glog"
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

func GormDB(cfg configs.DatabaseConfig) (*sql.DB, error) {
	connStr := fmt.Sprintf("user='%s' password=%s host=%s dbname='%s'", cfg.User, cfg.Password, cfg.Host, cfg.Name)
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db.DB()

}

func CreatePostgreSQLConfig(user, password, postgresHost, dbName string, postgresPort uint16,
) string {
	var b bytes.Buffer
	if dbName != "" {
		fmt.Fprintf(&b, "database=%s ", dbName)
	}
	if user != "" {
		fmt.Fprintf(&b, "user=%s ", user)
	}
	if password != "" {
		fmt.Fprintf(&b, "password=%s ", password)
	}
	if postgresHost != "" {
		fmt.Fprintf(&b, "host=%s ", postgresHost)
	}
	if postgresPort != 0 {
		fmt.Fprintf(&b, "port=%d ", postgresPort)
	}
	fmt.Fprint(&b, "sslmode=disable")

	return b.String()
}

func initDBDriver(cfg configs.DatabaseConfig) string {
	var sqlConfig string

	dbName := "OAuthProject"
	sqlConfig = CreatePostgreSQLConfig(cfg.User, cfg.Password, cfg.Host, "postgres", 5432)

	var db *sql.DB
	var err error

	operation := func() (*sql.DB, error) {
		db, err = sql.Open("pg", sqlConfig)
		if err != nil {
			return nil, err
		}
		return db, nil
	}
	/**
		using exponential backoof for trying connect postgre server
		interval trying is 3s
	*/
	b := backoff.NewExponentialBackOff()
	b.MaxInterval = 3000

	db, err = backoff.Retry(context.TODO(), operation, backoff.WithBackOff(b))

	util.TerminateIfError(err)

	defer db.Close()

	opr := func() (*string, error) {
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
		if ignoreAlreadyExistError(err) != nil {
			return nil, err
		}

		return nil, nil
	}
	b = backoff.NewExponentialBackOff()

	_, err = backoff.Retry(context.TODO(), opr, backoff.WithBackOff(b))

	util.TerminateIfError(err)

	// Note: postgreSQL does not have the option `ClientFoundRows`
	sqlConfig = CreatePostgreSQLConfig(cfg.User, cfg.Password, cfg.Host, dbName, 5432)

	_, err = db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	if ignoreAlreadyExistError(err) != nil {
		util.TerminateIfError(err)
	}


	return sqlConfig

}

func InitDBClient(cfg configs.DatabaseConfig) (*sql.DB, error) {
	arg :=initDBDriver(cfg)

	db, err := gorm.Open(postgres.Open(arg), &gorm.Config{})
	util.TerminateIfError(err)


	response := db.AutoMigrate(
		&entity.OAuthClient{},
		&entity.OAuthCode{},
		&entity.OAuthScope{},
		&entity.OAuthToken{},
		&entity.OAuthUser{},
	)

	if ignoreAlreadyExistError(response)  != nil {
	 glog.Fatalf("Failed to initialize the databases. Error: %s", response.Error)
	}


	return db.DB()
}

func ignoreAlreadyExistError(err error) error {
	if strings.Contains(err.Error(), util.PGX_EXIST_ERROR) {
		return nil
	}
	return err
}
