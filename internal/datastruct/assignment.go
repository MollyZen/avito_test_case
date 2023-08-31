package datastruct

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Assignment struct {
	UserID    int64              `db:"userid"`
	SegmentID int64              `db:"segmentid"`
	UntilDate pgtype.Timestamptz `db:"untildate"`
}
