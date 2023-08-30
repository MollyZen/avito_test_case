package repository

import (
	"avito_test_case/internal/datastruct"
	"avito_test_case/pkg/logger"
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type PostgresAssignmentRepository struct {
	db  *pgxpool.Pool
	log logger.Logger
}

func NewPostgresAssignmentRepository(db *pgxpool.Pool, logger logger.Logger) *PostgresAssignmentRepository {
	return &PostgresAssignmentRepository{
		db:  db,
		log: logger,
	}
}

func (p PostgresAssignmentRepository) Create(ctx context.Context, assignments []datastruct.Assignment) ([]datastruct.Assignment, error) {
	q := `
		INSERT INTO segmenting.assignment
			(userid, segmentid, untildate)
		(SELECT * FROM UNNEST($1::bigint[], $2::bigint[], $3::timestamptz[]))
		ON CONFLICT (userid, segmentid) DO NOTHING
		RETURNING userid, segmentid, untildate
		`
	var res []datastruct.Assignment
	var userIds = make([]int64, len(assignments))
	var segmentIds = make([]int64, len(assignments))
	var untilDates = make([]time.Time, len(assignments))
	for i, val := range assignments {
		userIds[i] = val.UserID
		segmentIds[i] = val.SegmentID
		untilDates[i] = val.UntilDate
	}
	if err := pgxscan.Select(ctx, p.db, &res, q, userIds, segmentIds, untilDates); err != nil {
		return nil, err
	}
	return res, nil
}

func (p PostgresAssignmentRepository) Update(ctx context.Context, assignments []datastruct.Assignment) ([]datastruct.Assignment, error) {
	q := `
		UPDATE segmenting.assignment
		SET untildate =  c.untildate
		FROM (
			(SELECT * FROM UNNEST($1::bigint[], $2::bigint[], $3::timestamptz[]))
		) as c(userid, segmentid, untildate)
		WHERE userid = c.userid AND segmentid = c.segmentid
		RETURNING userid, segmentid, untildate
		`
	var res []datastruct.Assignment
	var userIds = make([]int64, len(assignments))
	var segmentIds = make([]int64, len(assignments))
	var untilDates = make([]time.Time, len(assignments))
	for i, val := range assignments {
		userIds[i] = val.UserID
		segmentIds[i] = val.SegmentID
		untilDates[i] = val.UntilDate
	}
	if err := pgxscan.Select(ctx, p.db, &res, q, userIds, segmentIds, untilDates); err != nil {
		return nil, err
	}
	return res, nil
}

func (p PostgresAssignmentRepository) Delete(ctx context.Context, assignments []datastruct.Assignment) ([]datastruct.Assignment, error) {
	q := `
		DELETE FROM segmenting.assignment
		WHERE (userid, segmentid) = ANY(SELECT * FROM UNNEST($1::bigint[], $2::bigint[]))
		RETURNING userid, segmentid, untildate
		`
	var res []datastruct.Assignment
	var userIds = make([]int64, len(assignments))
	var segmentIds = make([]int64, len(assignments))
	for i, val := range assignments {
		userIds[i] = val.UserID
		segmentIds[i] = val.SegmentID
	}
	if err := pgxscan.Select(ctx, p.db, &res, q, userIds, segmentIds); err != nil {
		return nil, err
	}
	return res, nil
}

func (p PostgresAssignmentRepository) GetAllForUser(ctx context.Context, userID int64) ([]datastruct.Assignment, error) {
	q := `
		SELECT * from segmenting.assignment
		WHERE userid = $1
		`
	var res []datastruct.Assignment
	if err := pgxscan.Select(ctx, p.db, &res, q, userID); err != nil {
		return nil, err
	}
	return res, nil
}
