package repository

import (
	"avito_test_case/internal/datastruct"
	"avito_test_case/pkg/logger"
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresSegmentRepository struct {
	db  *pgxpool.Pool
	log logger.Logger
}

func NewPostgresSegmentRepository(db *pgxpool.Pool, logger logger.Logger) *PostgresSegmentRepository {
	return &PostgresSegmentRepository{
		db:  db,
		log: logger,
	}
}

func (p PostgresSegmentRepository) UpsertWithConn(ctx context.Context, segment datastruct.Segment, conn *pgx.Conn) (datastruct.Segment, error) {
	q := `
		INSERT INTO segmenting.segment
			(slug, creationdate)
		VALUES 
			($1, now())
		ON CONFLICT (slug)
			DO UPDATE 
			SET isactive = TRUE
		RETURNING id, slug, isactive, creationdate`
	var res datastruct.Segment
	if err := pgxscan.Get(ctx, conn, &res, q, segment.Slug); err != nil {
		return datastruct.Segment{}, err
	}

	return res, nil
}

func (p PostgresSegmentRepository) Upsert(ctx context.Context, segment datastruct.Segment) (datastruct.Segment, error) {
	conn, err := p.db.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return datastruct.Segment{}, err
	}
	return p.UpsertWithConn(ctx, segment, conn.Conn())
}

func (p PostgresSegmentRepository) GetAllBySlugWithConn(ctx context.Context, slugs []string, conn *pgx.Conn) ([]datastruct.Segment, error) {
	q := `
		SELECT * FROM segmenting.segment
		WHERE (slug) = ANY($1::varchar[])
		`
	var res []datastruct.Segment
	if err := pgxscan.Select(ctx, conn, &res, q, slugs); err != nil {
		return nil, err
	}

	return res, nil
}

func (p PostgresSegmentRepository) GetAllBySlug(ctx context.Context, slugs []string) ([]datastruct.Segment, error) {
	conn, err := p.db.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return nil, err
	}
	return p.GetAllBySlugWithConn(ctx, slugs, conn.Conn())
}

func (p PostgresSegmentRepository) GetAllByIdWithConn(ctx context.Context, ids []int64, conn *pgx.Conn) ([]datastruct.Segment, error) {
	q := `
		SELECT * FROM segmenting.segment
		WHERE (id) = ANY($1::bigint[])
		`
	var res []datastruct.Segment
	if err := pgxscan.Select(ctx, conn, &res, q, ids); err != nil {
		return nil, err
	}

	return res, nil
}

func (p PostgresSegmentRepository) GetAllById(ctx context.Context, ids []int64) ([]datastruct.Segment, error) {
	conn, err := p.db.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return nil, err
	}
	return p.GetAllByIdWithConn(ctx, ids, conn.Conn())
}

func (p PostgresSegmentRepository) DeleteByIdWithConn(ctx context.Context, segmentId int64, conn *pgx.Conn) (datastruct.Segment, error) {
	q := `
		UPDATE segmenting.segment 
		SET isactive = FALSE
		WHERE id = $1
		RETURNING id, slug, isactive, creationdate`
	var res datastruct.Segment
	if err := pgxscan.Get(ctx, conn, &res, q, segmentId); err != nil {
		return datastruct.Segment{}, err
	}

	return res, nil
}

func (p PostgresSegmentRepository) DeleteById(ctx context.Context, segmentId int64) (datastruct.Segment, error) {
	conn, err := p.db.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return datastruct.Segment{}, err
	}
	return p.DeleteByIdWithConn(ctx, segmentId, conn.Conn())
}

func (p PostgresSegmentRepository) DeleteBySlugWithConn(ctx context.Context, segmentSlug string, conn *pgx.Conn) (datastruct.Segment, error) {
	q := `
		UPDATE segmenting.segment 
		SET isactive = FALSE
		WHERE slug = $1
		RETURNING id, slug, isactive, creationdate`
	var res datastruct.Segment
	if err := pgxscan.Get(ctx, conn, &res, q, segmentSlug); err != nil {
		return datastruct.Segment{}, err
	}

	return res, nil
}

func (p PostgresSegmentRepository) DeleteBySlug(ctx context.Context, segmentSlug string) (datastruct.Segment, error) {
	conn, err := p.db.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return datastruct.Segment{}, err
	}
	return p.DeleteBySlugWithConn(ctx, segmentSlug, conn.Conn())
}

func (p PostgresSegmentRepository) AddToPercentOfUsersWithConn(ctx context.Context, segmentID int64, percent float64, untilDate pgtype.Timestamptz, conn *pgx.Conn) ([]datastruct.Assignment, error) {
	q := `
		WITH users AS (SELECT id FROM segmenting.user TABLESAMPLE BERNOULLI($1)),
		segs AS (SELECT $2::bigint AS segmentid, $3::timestamptz AS untildate),
		ljoin AS (SELECT u.id AS userid, s.segmentid, s.untildate FROM users u LEFT JOIN segs s on TRUE)
		INSERT INTO segmenting.assignment
		SELECT userid, segmentid, untildate FROM ljoin
		RETURNING userid
	`
	var ids []int64
	if err := pgxscan.Select(ctx, conn, &ids, q, percent, segmentID, untilDate); err != nil {
		return nil, err
	}

	res := make([]datastruct.Assignment, len(ids))
	for i, v := range ids {
		res[i].UserID = v
	}

	return res, nil
}

func (p PostgresSegmentRepository) AddToPercentOfUsers(ctx context.Context, segmentID int64, percent float64, untilDate pgtype.Timestamptz) ([]datastruct.Assignment, error) {
	conn, err := p.db.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return nil, err
	}
	return p.AddToPercentOfUsersWithConn(ctx, segmentID, percent, untilDate, conn.Conn())
}
