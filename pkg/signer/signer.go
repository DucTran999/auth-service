package signer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type TokenSigner interface {
	SignAccessToken(claims jwt.Claims) (string, error)
	SignRefreshToken(claims jwt.Claims) (string, error)
	Parse(token string) (*jwt.MapClaims, error)
}

type Config struct {
	Alg     SigningAlgorithm
	PrivPem string
	PubPem  string
	HashKey string
}

type tokenSigner struct {
	method  jwt.SigningMethod
	privKey any
	pubKey  any
}

func NewTokenSigner(config Config) (TokenSigner, error) {
	method, err := config.Alg.ToJWTMethod()
	if err != nil {
		return nil, err
	}

	ts := &tokenSigner{
		method: method,
	}

	if err := ts.loadKey(config); err != nil {
		return nil, err
	}

	return ts, nil
}

func (ts *tokenSigner) SignAccessToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(ts.method, claims)
	return token.SignedString(ts.privKey)
}

func (ts *tokenSigner) SignRefreshToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(ts.method, claims)
	return token.SignedString(ts.privKey)
}

func (ts *tokenSigner) Parse(tokenStr string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		// Validate that the token's algorithm matches our signer's algorithm
		if t.Method.Alg() != ts.method.Alg() {
			return nil, fmt.Errorf("unexpected signing method: got %v, expected %v", t.Header["alg"], ts.method.Alg())
		}
		if strings.HasPrefix(ts.method.Alg(), "HS") {
			return ts.privKey, nil
		}
		return ts.pubKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidTokenClaimType
	}

	if !token.Valid {
		return nil, ErrInvalidTokenSignature
	}

	return &claims, nil
}

func (ts *tokenSigner) loadKey(config Config) error {
	switch ts.method.Alg() {
	case jwt.SigningMethodHS256.Alg(),
		jwt.SigningMethodHS384.Alg(),
		jwt.SigningMethodHS512.Alg():
		// HMAC uses raw secret
		ts.privKey = []byte(config.HashKey)
		return nil

	case jwt.SigningMethodRS256.Alg(),
		jwt.SigningMethodRS384.Alg(),
		jwt.SigningMethodRS512.Alg():
		return ts.loadRSAKey(config.PrivPem, config.PubPem)

	case jwt.SigningMethodES256.Alg(),
		jwt.SigningMethodES384.Alg(),
		jwt.SigningMethodES512.Alg():
		return ts.loadECDSAKey(config.HashKey)
	case jwt.SigningMethodEdDSA.Alg():
		return ts.loadEdDSAKey(config.HashKey)

	default:
		return fmt.Errorf("unsupported signing algorithm: %s", ts.method.Alg())
	}
}

func (ts *tokenSigner) loadRSAKey(privPemFilename, pubPemFilename string) error {
	privPemPath, err := ts.buildPath(privPemFilename)
	if err != nil {
		return err
	}

	// #nosec G304 -- path is sanitized and restricted to baseDir
	privPem, err := os.ReadFile(privPemPath)
	if err != nil {
		return fmt.Errorf("read private RSA key: %w", err)
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(privPem)
	if err != nil {
		return fmt.Errorf("parse private RSA key: %w", err)
	}
	ts.privKey = privKey

	pubPemPath, err := ts.buildPath(pubPemFilename)
	if err != nil {
		return err
	}

	// #nosec G304 -- path is sanitized and restricted to baseDir
	pubPem, err := os.ReadFile(pubPemPath)
	if err != nil {
		return fmt.Errorf("read public RSA key: %w", err)
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubPem)
	if err != nil {
		return fmt.Errorf("parse public RSA key: %w", err)
	}

	ts.pubKey = pubKey

	return nil
}

func (ts *tokenSigner) loadECDSAKey(path string) error {
	absPath, err := ts.buildPath(path)
	if err != nil {
		return err
	}

	// #nosec G304 -- path is sanitized and restricted to baseDir
	pemData, err := os.ReadFile(absPath)
	if err != nil {
		return fmt.Errorf("read ECDSA key: %w", err)
	}
	key, err := jwt.ParseECPrivateKeyFromPEM(pemData)
	if err != nil {
		return fmt.Errorf("parse ECDSA key: %w", err)
	}
	ts.privKey = key
	return nil
}

func (ts *tokenSigner) loadEdDSAKey(path string) error {
	absPath, err := ts.buildPath(path)
	if err != nil {
		return err
	}

	// #nosec G304 -- path is sanitized and restricted to baseDir
	keyData, err := os.ReadFile(absPath)
	if err != nil {
		return fmt.Errorf("read EdDSA key: %w", err)
	}
	privateKey, err := jwt.ParseEdPrivateKeyFromPEM(keyData)
	if err != nil {
		return fmt.Errorf("parse EdDSA key: %w", err)
	}
	ts.privKey = privateKey
	return nil
}

func (ts *tokenSigner) buildPath(path string) (string, error) {
	const baseDir = "./keys/"

	cleanPath := filepath.Clean(path)
	absBase, err := filepath.Abs(baseDir)
	if err != nil {
		return "", fmt.Errorf("failed to resolve base dir: %w", err)
	}
	joinedPath := filepath.Join(baseDir, cleanPath)
	absPath, err := filepath.Abs(joinedPath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute path: %w", err)
	}
	if !strings.HasPrefix(absPath, absBase) {
		return "", fmt.Errorf("unauthorized path access: %s", path)
	}

	return absPath, nil
}
