package datastruct

import "time"

type User struct {
	ID           int64     `db:"id"`
	CreationDate time.Time `db:"creationdate"`
}
