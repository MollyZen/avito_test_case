package repository

import (
	"avito_test_case/internal/datastruct"
	"avito_test_case/pkg/logger"
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
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

func (p PostgresUserRepository) GetWithConn(ctx context.Context, userID int64, conn *pgx.Conn) (datastruct.User, error) {
	q := `
		SELECT * FROM segmenting."user"
		WHERE id = $1
		`
	var res datastruct.User
	if err := pgxscan.Get(ctx, conn, &res, q, userID); err != nil {
		p.log.Error(err)
		return datastruct.User{}, err
	}

	return res, nil
}

func (p PostgresUserRepository) Get(ctx context.Context, userID int64) (datastruct.User, error) {
	conn, err := p.db.Acquire(context.TODO())
	defer conn.Release()
	if err != nil {
		return datastruct.User{}, err
	}
	return p.GetWithConn(context.TODO(), userID, conn.Conn())
}

func (p PostgresUserRepository) CreateWithConn(ctx context.Context, user datastruct.User, conn *pgx.Conn) (datastruct.User, error) {
	q := `
		INSERT INTO segmenting."user"
			(id, creationdate)
		VALUES 
			($1, now())
		RETURNING creationdate
		`
	var res datastruct.User
	if err := pgxscan.Get(ctx, conn, &res, q, user.ID); err != nil {
		p.log.Error(err)
		return datastruct.User{}, err
	}

	return res, nil
}

func (p PostgresUserRepository) Create(ctx context.Context, user datastruct.User) (datastruct.User, error) {
	conn, err := p.db.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return datastruct.User{}, err
	}
	return p.CreateWithConn(ctx, user, conn.Conn())
}

func (p PostgresUserRepository) UpsertWithConn(ctx context.Context, user datastruct.User, conn *pgx.Conn) (datastruct.User, error) {
	q := `
		with new_row as (INSERT INTO segmenting."user"
			(id, creationdate)
		VALUES 
			($1, now())
		ON CONFLICT DO NOTHING
		RETURNING id, creationdate)
		SELECT id, creationdate FROM new_row
		UNION
		SELECT id, creationdate FROM segmenting."user" WHERE id = $1
		`
	var res datastruct.User
	if err := pgxscan.Get(ctx, conn, &res, q, user.ID); err != nil {
		p.log.Error(err)
		return datastruct.User{}, err
	}

	return res, nil
}

func (p PostgresUserRepository) Upsert(ctx context.Context, user datastruct.User) (datastruct.User, error) {
	conn, err := p.db.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return datastruct.User{}, err
	}
	return p.UpsertWithConn(ctx, user, conn.Conn())
}
