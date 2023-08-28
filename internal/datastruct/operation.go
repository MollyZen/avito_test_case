package datastruct

import "time"

type Operation struct {
	Name         string
	Description  string
	CreationDate time.Time
	UpdateDate   time.Time
	IsActive     bool
}
