package signer

import (
	"fmt"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
)

type TokenSigner interface {
	SignAccessToken(claims jwt.Claims) (string, error)
	SignRefreshToken(claims jwt.Claims) (string, error)
	// Parse(token string) (*jwt.Claims, error)
}

type tokenSigner struct {
	method jwt.SigningMethod
	key    any
}

func NewTokenSigner(alg SigningAlgorithm, keyIdentifier string) (TokenSigner, error) {
	keyIdentifier = strings.TrimSpace(keyIdentifier)
	if keyIdentifier == "" {
		return nil, ErrMissingKey
	}

	method, err := alg.ToJWTMethod()
	if err != nil {
		return nil, err
	}

	ts := &tokenSigner{
		method: method,
	}

	if err := ts.loadKey(keyIdentifier); err != nil {
		return nil, err
	}

	return ts, nil
}

func (ts *tokenSigner) SignAccessToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(ts.method, claims)
	return token.SignedString(ts.key)
}

func (ts *tokenSigner) SignRefreshToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(ts.method, claims)
	return token.SignedString(ts.key)
}

func (ts *tokenSigner) loadKey(keyIdentifier string) error {
	switch ts.method.Alg() {
	case jwt.SigningMethodHS256.Alg(),
		jwt.SigningMethodHS384.Alg(),
		jwt.SigningMethodHS512.Alg():
		// HMAC uses raw secret
		ts.key = []byte(keyIdentifier)
		return nil

	case jwt.SigningMethodRS256.Alg(),
		jwt.SigningMethodRS384.Alg(),
		jwt.SigningMethodRS512.Alg():
		return ts.loadRSAKey(keyIdentifier)

	case jwt.SigningMethodES256.Alg(),
		jwt.SigningMethodES384.Alg(),
		jwt.SigningMethodES512.Alg():
		return ts.loadECDSAKey(keyIdentifier)

	default:
		return ts.loadEdDSAKey(keyIdentifier)
	}
}

func (ts *tokenSigner) loadRSAKey(path string) error {
	pemData, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read RSA key: %w", err)
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(pemData)
	if err != nil {
		return fmt.Errorf("parse RSA key: %w", err)
	}
	ts.key = key
	return nil
}

func (ts *tokenSigner) loadECDSAKey(path string) error {
	pemData, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read ECDSA key: %w", err)
	}
	key, err := jwt.ParseECPrivateKeyFromPEM(pemData)
	if err != nil {
		return fmt.Errorf("parse ECDSA key: %w", err)
	}
	ts.key = key
	return nil
}

func (ts *tokenSigner) loadEdDSAKey(path string) error {
	keyData, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read EdDSA key: %w", err)
	}
	privateKey, err := jwt.ParseEdPrivateKeyFromPEM(keyData)
	if err != nil {
		return fmt.Errorf("parse EdDSA key: %w", err)
	}
	ts.key = privateKey
	return nil
}
