package domain

type ChangePasswordInput struct {
	AccountID   string
	OldPassword string
	NewPassword string
}

type RegisterInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
