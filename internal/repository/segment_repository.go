package repository

import (
	"avito_test_case/internal/datastruct"
	"context"
)

type SegmentRepository interface {
	Create(ctx context.Context, segment datastruct.Segment)
	DeleteById(ctx context.Context, segmentId int64)
	DeleteBySlug(ctx context.Context, segmentSlug string)
}
