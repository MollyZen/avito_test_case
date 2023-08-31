package repository

import (
	"avito_test_case/internal/datastruct"
	"context"
	"time"
)

type HistoryRepository interface {
	CreateAll(ctx context.Context, history []datastruct.History) error
	GetAllForUserPeriod(ctx context.Context, userID int64, start, end time.Time) ([]datastruct.History, error)
}
