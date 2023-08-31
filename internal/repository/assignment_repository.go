package repository

import (
	"avito_test_case/internal/datastruct"
	"context"
)

type AssignmentRepository interface {
	Create(ctx context.Context, assignments []datastruct.Assignment) ([]datastruct.Assignment, error)
	Update(ctx context.Context, assignments []datastruct.Assignment) ([]datastruct.Assignment, error)
	Delete(ctx context.Context, assignments []datastruct.Assignment) ([]datastruct.Assignment, error)
	DeleteAllForSeg(ctx context.Context, segmentID int64) ([]datastruct.Assignment, error)
	GetAllForUser(ctx context.Context, userID int64) ([]datastruct.Assignment, error)
}
