package model

import "github.com/google/uuid"

type User struct {
	ID   uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	Name string    `json:"name" gorm:"column:name;"`
}
