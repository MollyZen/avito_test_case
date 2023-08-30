package repository

import (
	"avito_test_case/internal/datastruct"
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, user datastruct.User) (datastruct.User, error)
	Upsert(ctx context.Context, user datastruct.User) (datastruct.User, error)
}
