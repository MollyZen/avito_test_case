package repository

import (
	"avito_test_case/internal/datastruct"
	"avito_test_case/pkg/logger"
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
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

func (p PostgresHistoryRepository) CreateAll(ctx context.Context, history []datastruct.History) error {
	q := `
		INSERT INTO segmenting.history
			(userid, segmentid, operationid)
		(SELECT * FROM UNNEST($1::bigint[], $2::bigint[], $3::timestamptz[]))
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
	if err := pgxscan.Select(ctx, p.db, &res, q, userIds, segmentIds, operationIds); err != nil {
		return err
	}
	return nil
}

func (p PostgresHistoryRepository) GetAllForUserPeriod(ctx context.Context, userID int64, start, end time.Time) ([]datastruct.History, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresHistoryRepository) GetForUserMonth(ctx context.Context, userID int64, month, year int) {
	//TODO implement me
	panic("implement me")
}
