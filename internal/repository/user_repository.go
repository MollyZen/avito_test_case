package repository

import (
	"avito_test_case/internal/datastruct"
	"context"
)

type UserRepository interface {
	Get(ctx context.Context, useID int64) (datastruct.User, error)
	Create(ctx context.Context, user datastruct.User) (datastruct.User, error)
	Upsert(ctx context.Context, user datastruct.User) (datastruct.User, error)
}
