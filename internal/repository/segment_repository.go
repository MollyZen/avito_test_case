package repository

import (
	"avito_test_case/internal/datastruct"
	"context"
)

type SegmentRepository interface {
	Create(ctx context.Context, segment datastruct.Segment) (datastruct.Segment, error)
	GetAllBySlug(ctx context.Context, slugs []string) ([]datastruct.Segment, error)
	DeleteById(ctx context.Context, segmentId int64) (datastruct.Segment, error)
	DeleteBySlug(ctx context.Context, segmentSlug string) (datastruct.Segment, error)
}
