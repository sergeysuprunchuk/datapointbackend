package repository

import (
	"context"
	"datapointbackend/internal/entity"
	"datapointbackend/pkg/database"
)

type SourceRepository struct {
	db *database.Database
}

func NewSourceRepository(db *database.Database) *SourceRepository {
	return &SourceRepository{db: db}
}

func (r *SourceRepository) GetAll(ctx context.Context) ([]entity.Source, error) {
	rows, err := r.db.Builder.
		Select("id", "name", "host", "port", "username", "password", "database_name", "driver").
		From("source").
		QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var sl []entity.Source
	for rows.Next() {
		var s entity.Source
		if err = rows.Scan(&s.Id, &s.Name, &s.Host, &s.Port, &s.Username, &s.Password, &s.DatabaseName, &s.Driver); err != nil {
			return nil, err
		}
		sl = append(sl, s)
	}

	return sl, nil
}

func (r *SourceRepository) GetOne(ctx context.Context, id string) (entity.Source, error) {
	var s entity.Source
	return s, r.db.Builder.
		Select("id", "name", "host", "port", "username", "password", "database_name", "driver").
		From("source").
		Where("id = ?", id).
		QueryRowContext(ctx).
		Scan(&s.Id, &s.Name, &s.Host, &s.Port, &s.Username, &s.Password, &s.DatabaseName, &s.Driver)
}

func (r *SourceRepository) Edit(ctx context.Context, s entity.Source) error {
	_, err := r.db.Builder.
		Update("source").
		Set("name", s.Name).
		Set("host", s.Host).
		Set("port", s.Port).
		Set("username", s.Username).
		Set("password", s.Password).
		Set("database_name", s.DatabaseName).
		Set("driver", s.Driver).
		Where("id = ?", s.Id).
		ExecContext(ctx)
	return err
}

func (r *SourceRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Builder.
		Delete("source").
		Where("id = ?", id).
		ExecContext(ctx)
	return err
}

func (r *SourceRepository) Create(ctx context.Context, s entity.Source) (string, error) {
	var id string
	return id, r.db.Builder.
		Insert("source").
		Columns("name", "host", "port", "username", "password", "database_name", "driver").
		Values(s.Name, s.Host, s.Port, s.Username, s.Password, s.DatabaseName, s.Driver).
		Suffix("RETURNING id").
		QueryRowContext(ctx).
		Scan(&id)
}
