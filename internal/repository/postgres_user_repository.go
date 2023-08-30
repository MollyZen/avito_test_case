package repository

import (
	"avito_test_case/internal/datastruct"
	"avito_test_case/pkg/logger"
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepository struct {
	db  *pgxpool.Pool
	log logger.Logger
}

func NewPostgresUserRepository(db *pgxpool.Pool, logger logger.Logger) *PostgresUserRepository {
	return &PostgresUserRepository{
		db:  db,
		log: logger,
	}
}

func (p PostgresUserRepository) Create(ctx context.Context, user datastruct.User) (datastruct.User, error) {
	q := `
		INSERT INTO segmenting.user
			(id, creationdate)
		VALUES 
			($1, now())
		RETURNING creationdate
		`
	var res datastruct.User
	if err := pgxscan.Get(ctx, p.db, &res, q, user.ID); err != nil {
		p.log.Error(err)
		return datastruct.User{}, err
	}

	return res, nil
}

func (p PostgresUserRepository) Upsert(ctx context.Context, user datastruct.User) (datastruct.User, error) {
	q := `
		INSERT INTO segmenting.user
			(id, creationdate)
		VALUES 
			($1, now())
		ON CONFLICT DO NOTHING 
		RETURNING creationdate
		`
	var res datastruct.User
	if err := pgxscan.Get(ctx, p.db, &res, q, user.ID); err != nil {
		p.log.Error(err)
		return datastruct.User{}, err
	}

	return res, nil
}
