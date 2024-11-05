package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Username  string    `json:"username" gorm:"type:varchar(255);not null"`
	Password  string    `json:"password" gorm:"type:varchar(255);not null"`
	Address   string    `json:"address,omitempty" gorm:"type:varchar(255)"`
	FirstName string    `json:"first_name" gorm:"type:varchar(50);not null"`
	LastName  string    `json:"last_name" gorm:"type:varchar(50);not null"`
	Email     string    `json:"email" gorm:"type:varchar(50);unique;not null"`
	Gender    string    `json:"gender,omitempty" gorm:"type:varchar(50)"`
	Phone     string    `json:"phone,omitempty" gorm:"type:varchar(50)"`
	IsActive  int       `json:"is_active" gorm:"type:int;default:0"`
	IsDeleted int       `json:"is_deleted" gorm:"type:int;default:0"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (c *User) TableName() string { return "users" }
