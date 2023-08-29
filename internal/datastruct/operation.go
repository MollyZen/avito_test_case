package datastruct

import "time"

type Operation struct {
	ID           int
	Name         string
	Description  string
	CreationDate time.Time
	UpdateDate   time.Time
	IsActive     bool
}
