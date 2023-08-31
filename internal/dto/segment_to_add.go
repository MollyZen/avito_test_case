package dto

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type SegmentToAdd struct {
	Slug      string             `json:"slug" validate:"required"`
	UntilDate pgtype.Timestamptz `json:"untilDate"`
}
