package repository

import (
	"avito_test_case/internal/datastruct"
	"avito_test_case/pkg/logger"
	"context"
	"database/sql"
	"fmt"
	"log"
)

type DatabaseSegmentRepository struct {
	db  *sql.DB
	log *logger.Logger
}

func NewDatabaseSegmentRepository(db *sql.DB, logger *logger.Logger) *DatabaseSegmentRepository {
	return &DatabaseSegmentRepository{
		db:  db,
		log: logger,
	}
}

func (d DatabaseSegmentRepository) Create(ctx context.Context, segment datastruct.Segment) (datastruct.Segment, error) {
	q := `
		INSERT INTO segmenting.segment
			(name, creationdate)
		VALUES 
			($1, now())
		RETURNING id, name, creationdate
`
	var res datastruct.Segment
	if err := d.db.QueryRowContext(ctx, q, segment.Name).Scan(&res.ID, &res.Name, &res.CreationDate); err != nil {
		log.Print(err)
		return datastruct.Segment{}, err
	}

	fmt.Println(res)

	return res, nil
}

func (d DatabaseSegmentRepository) DeleteById(ctx context.Context, segmentId int64) {
	//TODO implement me
	panic("implement me")
}

func (d DatabaseSegmentRepository) DeleteBySlug(ctx context.Context, segmentSlug string) {
	//TODO implement me
	panic("implement me")
}
