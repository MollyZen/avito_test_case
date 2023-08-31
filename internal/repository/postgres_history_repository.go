package repository

import (
	"avito_test_case/internal/datastruct"
	"avito_test_case/pkg/logger"
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type PostgresHistoryRepository struct {
	db  *pgxpool.Pool
	log logger.Logger
}

func NewPostgresHistoryRepository(db *pgxpool.Pool, logger logger.Logger) *PostgresHistoryRepository {
	return &PostgresHistoryRepository{
		db:  db,
		log: logger,
	}
}

func (p PostgresHistoryRepository) CreateAllWithConn(ctx context.Context, history []datastruct.History, conn *pgx.Conn) error {
	q := `
		INSERT INTO segmenting.history
			(userid, segmentid, operationid)
		(SELECT * FROM UNNEST($1::bigint[], $2::bigint[], $3::bigint[]))
		RETURNING id, userid, segmentid, operationid, timestamp
		`
	var res []datastruct.History
	var userIds = make([]int64, len(history))
	var segmentIds = make([]int64, len(history))
	var operationIds = make([]int64, len(history))
	for i, val := range history {
		userIds[i] = val.UserID
		segmentIds[i] = val.SegmentID
		operationIds[i] = val.OperationID
	}
	if err := pgxscan.Select(ctx, conn, &res, q, userIds, segmentIds, operationIds); err != nil {
		return err
	}
	return nil
}

func (p PostgresHistoryRepository) CreateAll(ctx context.Context, history []datastruct.History) error {
	conn, err := p.db.Acquire(ctx)
	if err != nil {
		return err
	}
	return p.CreateAllWithConn(ctx, history, conn.Conn())
}

func (p PostgresHistoryRepository) GetAllForUserPeriod(ctx context.Context, userID int64, start, end time.Time) ([]datastruct.History, error) {
	q := `
		SELECT * from segmenting.history
		WHERE userid = $1
		AND timestamp >= $2
		AND timestamp < $3	
		`
	var res []datastruct.History
	if err := pgxscan.Select(ctx, p.db, &res, q, userID, start, end); err != nil {
		return nil, err
	}
	return res, nil
}
