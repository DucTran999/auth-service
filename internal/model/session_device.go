package model

import "time"

type SessionDevice struct {
	JTI       string
	AccountID string

	UserAgent string
	IP        string

	CreatedAt time.Time
	ExpiresAt time.Time
}
