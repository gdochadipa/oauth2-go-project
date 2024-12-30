package repository

import (
	"database/sql"
	"errors"
)

var (
	ErrDuplicate    = errors.New("record already exists")
	ErrNotExist     = errors.New("row does not exist")
	ErrUpdateFailed = errors.New("update failed")
	ErrDeleteFailed = errors.New("delete failed")
)

type Repository interface {
	Close()
	ItemRepository
}

type dbRepository struct {
	db *sql.DB
}

func NewDBRepository(db *sql.DB) Repository {
	return &dbRepository{db}
}

func (r *dbRepository) Close() {
	r.db.Close()
}

func (r *dbRepository) Ping() error {
	return r.db.Ping()
}
