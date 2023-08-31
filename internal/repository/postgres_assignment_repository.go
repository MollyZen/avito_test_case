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

func (p PostgresAssignmentRepository) CreateWithConn(ctx context.Context, assignments []datastruct.Assignment, conn *pgx.Conn) ([]datastruct.Assignment, error) {
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
	var untilDates = make([]pgtype.Timestamptz, len(assignments))
	for i, val := range assignments {
		userIds[i] = val.UserID
		segmentIds[i] = val.SegmentID
		untilDates[i] = val.UntilDate
	}
	if err := pgxscan.Select(ctx, conn, &res, q, userIds, segmentIds, untilDates); err != nil {
		return nil, err
	}
	return res, nil
}

func (p PostgresAssignmentRepository) Create(ctx context.Context, assignments []datastruct.Assignment) ([]datastruct.Assignment, error) {
	conn, err := p.db.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return nil, err
	}
	return p.CreateWithConn(ctx, assignments, conn.Conn())
}

func (p PostgresAssignmentRepository) UpdateWithConn(ctx context.Context, assignments []datastruct.Assignment, conn *pgx.Conn) ([]datastruct.Assignment, error) {
	q := `
		UPDATE segmenting.assignment
		SET untildate =  c.untildate
		FROM (
			(SELECT * FROM UNNEST($1::bigint[], $2::bigint[], $3::timestamptz[]))
		) as c(userid, segmentid, untildate)
		WHERE assignment.userid = c.userid AND assignment.segmentid = c.segmentid
		RETURNING c.userid, c.segmentid, c.untildate
		`
	var res []datastruct.Assignment
	var userIds = make([]int64, len(assignments))
	var segmentIds = make([]int64, len(assignments))
	var untilDates = make([]pgtype.Timestamptz, len(assignments))
	for i, val := range assignments {
		userIds[i] = val.UserID
		segmentIds[i] = val.SegmentID
		untilDates[i] = val.UntilDate
	}
	if err := pgxscan.Select(ctx, conn, &res, q, userIds, segmentIds, untilDates); err != nil {
		return nil, err
	}
	return res, nil
}

func (p PostgresAssignmentRepository) Update(ctx context.Context, assignments []datastruct.Assignment) ([]datastruct.Assignment, error) {
	conn, err := p.db.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return nil, err
	}
	return p.UpdateWithConn(ctx, assignments, conn.Conn())
}

func (p PostgresAssignmentRepository) DeleteWithConn(ctx context.Context, assignments []datastruct.Assignment, conn *pgx.Conn) ([]datastruct.Assignment, error) {
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
	if err := pgxscan.Select(ctx, conn, &res, q, userIds, segmentIds); err != nil {
		return nil, err
	}
	return res, nil
}

func (p PostgresAssignmentRepository) Delete(ctx context.Context, assignments []datastruct.Assignment) ([]datastruct.Assignment, error) {
	conn, err := p.db.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return nil, err
	}
	return p.DeleteWithConn(ctx, assignments, conn.Conn())
}

func (p PostgresAssignmentRepository) DeleteAllForSegWithConn(ctx context.Context, segmentID int64, conn *pgx.Conn) ([]datastruct.Assignment, error) {
	q := `
		DELETE FROM segmenting.assignment
		WHERE segmentid = $1
		RETURNING userid, segmentid, untildate
		`
	var res []datastruct.Assignment
	if err := pgxscan.Select(ctx, conn, &res, q, segmentID); err != nil {
		return nil, err
	}
	return res, nil
}

func (p PostgresAssignmentRepository) DeleteAllForSeg(ctx context.Context, segmentID int64) ([]datastruct.Assignment, error) {
	conn, err := p.db.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return nil, err
	}
	return p.DeleteAllForSegWithConn(ctx, segmentID, conn.Conn())
}

func (p PostgresAssignmentRepository) GetAllForUserWithConn(ctx context.Context, userID int64, conn *pgx.Conn) ([]datastruct.Assignment, error) {
	q := `
		SELECT * from segmenting.assignment
		WHERE userid = $1
		AND (untildate > now() OR untildate IS NULL)
		`
	var res []datastruct.Assignment
	if err := pgxscan.Select(ctx, conn, &res, q, userID); err != nil {
		return nil, err
	}
	return res, nil
}

func (p PostgresAssignmentRepository) GetAllForUser(ctx context.Context, userID int64) ([]datastruct.Assignment, error) {
	conn, err := p.db.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return nil, err
	}
	return p.GetAllForUserWithConn(ctx, userID, conn.Conn())
}

func (p PostgresAssignmentRepository) DeleteExpired(ctx context.Context) error {
	q := `
		DELETE FROM segmenting.assignment
		WHERE untildate < now()
		`
	if _, err := p.db.Exec(context.TODO(), q); err != nil {
		return err
	}
	return nil
}
