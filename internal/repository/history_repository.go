package repository

import (
	"avito_test_case/internal/datastruct"
	"context"
	"time"
)

type HistoryRepository interface {
	AddUserHistory(ctx context.Context, history []datastruct.History)
	GetUserHistory(ctx context.Context, userID int64, start, end time.Time) ([]datastruct.History, error)
}
