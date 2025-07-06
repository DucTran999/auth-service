package dto

// LoginInput represents the input required to authenticate a user using email and password.
// It also includes optional request metadata for logging, auditing, or session management.
type LoginInput struct {
	// CurrentSessionID is an optional field used to associate or trace login attempts
	CurrentSessionID string

	Email     string `json:"email"`    // User's email address
	Password  string `json:"password"` // Plain-text password from the login form
	IP        string `json:"-"`        // Client IP address (injected by handler, not from JSON)
	UserAgent string `json:"-"`        // User-Agent header string (injected by handler)
}

type LoginJWTInput struct {
	Email     string
	Password  string
	IP        string
	UserAgent string
}

type TokenPairs struct {
	AccessToken  string
	RefreshToken string
}
