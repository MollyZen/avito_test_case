package repository

import (
	"avito_test_case/internal/datastruct"
	"avito_test_case/pkg/logger"
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
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

func (p PostgresSegmentRepository) Create(ctx context.Context, segment datastruct.Segment) (datastruct.Segment, error) {
	q := `
		INSERT INTO segmenting.segment
			(name, creationdate)
		VALUES 
			($1, now())
		RETURNING id, name, creationdate
`
	var res datastruct.Segment
	if err := pgxscan.Get(ctx, p.db, &res, q, segment.Name); err != nil {
		p.log.Error(err)
		return datastruct.Segment{}, err
	}

	return res, nil
}

func (p PostgresSegmentRepository) GetForIds(ctx context.Context, slugs []string) ([]datastruct.Segment, error) {
	q := `
		SELECT * FROM segmenting.segment
		WHERE (name) = ANY(UNNEST($1::varchar[]))
		`
	var res []datastruct.Segment
	if err := pgxscan.Select(ctx, p.db, &res, q, slugs); err != nil {
		p.log.Error(err)
		return nil, err
	}

	return res, nil
}

func (p PostgresSegmentRepository) DeleteById(ctx context.Context, segmentId int64) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresSegmentRepository) DeleteBySlug(ctx context.Context, segmentSlug string) {
	//TODO implement me
	panic("implement me")
}
