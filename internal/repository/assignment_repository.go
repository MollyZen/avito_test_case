package repository

import (
	"avito_test_case/internal/datastruct"
	"context"
)

type AssignmentRepository interface {
	Assign(ctx context.Context, userID int64, segments []datastruct.Segment)
	Remove(ctx context.Context, userID int64, segments []datastruct.Segment)
	GetAllForUser(ctx context.Context, userID int64) ([]datastruct.Segment, error)
}
