package hasher

import "github.com/alexedwards/argon2id"

type Hasher interface {
	// HashPassword securely hashes a plain password
	HashPassword(password string) (string, error)

	// Verify password
	ComparePasswordAndHash(pass, hash string) (bool, error)
}

type hasher struct{}

func NewHasher() *hasher {
	return &hasher{}
}

// HashPassword securely hashes a plain password.
func (h *hasher) HashPassword(password string) (string, error) {
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}

func (h *hasher) ComparePasswordAndHash(pass, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(pass, hash)
}
