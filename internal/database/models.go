// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0

package database

import (
	"time"
)

type Feed struct {
	ID            int64
	Name          string
	LastPostTitle string
	Url           string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
