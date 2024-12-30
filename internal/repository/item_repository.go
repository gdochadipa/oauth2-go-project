package repository

import (
	"context"
	"errors"

	"github.com/gdochadipa/oauth2-go-project/internal/entity"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

type ItemRepository interface {
	GetListItem(ctx context.Context, skip uint32, limitPage uint32) ([]entity.Item, error)
	UpdateItem(ctx context.Context, id uuid.UUID, item entity.Item) (*entity.Item, error)
	InsertItem(ctx context.Context, item entity.Item) (*entity.Item, error)
	GetItem(ctx context.Context, id uuid.UUID) (*entity.Item, error)
	DeleteItem(ctx context.Context, id uuid.UUID) error
}

func (r *dbRepository) GetListItem(ctx context.Context, skip uint32, limitPage uint32) ([]entity.Item, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, description FROM items ORDER BY id DESC OFFSET $1 LIMIT $2", skip, limitPage)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	items := []entity.Item{}

	for rows.Next() {
		item := &entity.Item{}
		if err = rows.Scan(&item.Id, &item.Name, &item.Description); err == nil {
			items = append(items, *item)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *dbRepository) UpdateItem(ctx context.Context, id uuid.UUID, item entity.Item) (*entity.Item, error) {
	res, err := r.db.ExecContext(ctx, "UPDATE items SET name=$1, description=$2 WHERE id = $3", item.Name, item.Description, id)
	if err != nil {
		var pgxE *pgconn.PgError
		if errors.As(err, &pgxE) {
			if pgxE.Code == "23505" {
				return nil, ErrDuplicate
			}
		}
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, ErrUpdateFailed
	}
	return &item, nil
}

func (r *dbRepository) InsertItem(ctx context.Context, item entity.Item) (*entity.Item, error) {
	var id uuid.UUID
	err := r.db.QueryRowContext(ctx, "INSERT INTO items(name, description) values ($1, $2) RETURNING id", item.Name, item.Description).Scan(&id)
	if err != nil {
		var pgxE *pgconn.PgError
		if errors.As(err, &pgxE) {
			if pgxE.Code == "23505" {
				return nil, ErrDuplicate
			}
		}
	}

	item.Id = id
	return &item, err
}

func (r *dbRepository) GetItem(ctx context.Context, id uuid.UUID) (*entity.Item, error) {
	data := r.db.QueryRowContext(ctx, "SELECT * from items where id = $1", id)

	item := &entity.Item{}
	if err := data.Scan(&item.Id, &item.Name, &item.Description); err != nil {
		return nil, err
	}

	return item, nil
}

func (r *dbRepository) DeleteItem(ctx context.Context, id uuid.UUID) error {
	data, err := r.db.ExecContext(ctx, "DELETE FROM items where id = $id", id)

	if err != nil {
		return err
	}

	affected, err := data.RowsAffected()
	if affected == 0 {
		return ErrDeleteFailed
	}
	return err
}
