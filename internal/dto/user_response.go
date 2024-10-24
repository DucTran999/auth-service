package dto

import "github.com/google/uuid"

type CreateUserResp struct {
	ID uuid.UUID `json:"id"`
}
