package datastruct

import "time"

type Segment struct {
	ID           int64
	Slug         string
	IsActive     bool      `db:"isactive"`
	CreationDate time.Time `db:"creationdate"`
}
