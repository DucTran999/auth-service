package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/DucTran999/auth-service/internal/model"
	"github.com/DucTran999/auth-service/internal/usecase/dto"
	"github.com/DucTran999/auth-service/internal/usecase/port"
	"github.com/DucTran999/auth-service/internal/usecase/shared"
	"github.com/DucTran999/auth-service/pkg/cache"
	"github.com/DucTran999/auth-service/pkg/signer"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	AccessTokenLifetime  = 15 * time.Minute
	RefreshTokenLifetime = 7 * 24 * time.Hour
)

type tokenClaimsParams struct {
	signAt     time.Time
	exp        time.Duration
	jti        string
	account    *model.Account
	includeJTI bool
}

type authJWTUsecase struct {
	signer signer.TokenSigner
	cache  cache.Cache

	accountVerifier shared.AccountVerifier
}

func NewAuthJWTUsecase(
	signer signer.TokenSigner,
	cache cache.Cache,
	accountVerifier shared.AccountVerifier,
) port.AuthJWTUsecase {
	return &authJWTUsecase{
		signer:          signer,
		cache:           cache,
		accountVerifier: accountVerifier,
	}
}

// Login authenticates the user and returns a pair of JWT tokens (access + refresh).
func (uc *authJWTUsecase) Login(ctx context.Context, input dto.LoginJWTInput) (*dto.TokenPairs, error) {
	// Verify user credentials
	account, err := uc.accountVerifier.Verify(ctx, input.Email, input.Password)
	if err != nil {
		return nil, fmt.Errorf("auth: failed to verify credentials: %w", err)
	}

	// Prepare token metadata
	jti := uuid.NewString()
	issuedAt := time.Now()

	// Generate access and refresh tokens
	tokenPairs, err := uc.signTokenPairs(jti, issuedAt, account)
	if err != nil {
		return nil, fmt.Errorf("auth: failed to sign token pairs: %w", err)
	}

	// Store refresh token in cache/session store
	deviceSession := uc.newDeviceSession(jti, account.ID.String(), issuedAt, input)
	if err := uc.cacheRefreshToken(ctx, deviceSession); err != nil {
		return nil, fmt.Errorf("auth: failed to cache refresh token: %w", err)
	}

	return tokenPairs, nil
}

func (uc *authJWTUsecase) RefreshToken(ctx context.Context, refreshToken string) (*dto.TokenPairs, error) {
	if refreshToken == "" {
		return nil, model.ErrInvalidCredentials
	}

	claims := new(model.TokenClaims)
	if err := uc.signer.ParseInto(refreshToken, claims); err != nil {
		return nil, err
	}

	key := cache.KeyRefreshToken(claims.Subject, claims.ID)
	ok, err := uc.cache.Has(ctx, key)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, model.ErrInvalidCredentials
	}

	tokens, err := uc.resignTokenPairs(ctx, *claims)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (uc *authJWTUsecase) RevokeRefreshToken(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return model.ErrInvalidCredentials
	}

	claims := new(model.TokenClaims)
	if err := uc.signer.ParseInto(refreshToken, claims); err != nil {
		return err
	}

	key := cache.KeyRefreshToken(claims.Subject, claims.ID)
	_ = uc.cache.Del(ctx, key)

	return nil
}

func (uc *authJWTUsecase) signTokenPairs(
	jti string,
	signAt time.Time,
	account *model.Account,
) (*dto.TokenPairs, error) {
	// Access token claims
	accessClaims := uc.buildClaims(tokenClaimsParams{
		signAt:  signAt,
		exp:     AccessTokenLifetime,
		account: account,
	})

	accessToken, err := uc.signer.Sign(accessClaims)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// Refresh token claims
	refreshClaims := uc.buildClaims(tokenClaimsParams{
		signAt:     signAt,
		jti:        jti,
		exp:        RefreshTokenLifetime,
		account:    account,
		includeJTI: true,
	})

	refreshToken, err := uc.signer.Sign(refreshClaims)
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return &dto.TokenPairs{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (uc *authJWTUsecase) resignTokenPairs(ctx context.Context, old model.TokenClaims) (*dto.TokenPairs, error) {
	now := time.Now()
	jti := uuid.NewString()

	// Step 1: Revoke old refresh token
	oldKey := cache.KeyRefreshToken(old.Subject, old.ID)
	if err := uc.cache.Del(ctx, oldKey); err != nil {
		return nil, fmt.Errorf("failed to invalidate old refresh token: %w", err)
	}

	// Step 2: Build access token claims (no JTI)
	accessClaims := model.TokenClaims{
		Email: old.Email,
		Role:  old.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   old.Subject,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(AccessTokenLifetime)),
		},
	}

	accessToken, err := uc.signer.Sign(accessClaims)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// Step 3: Build refresh token claims (new JTI)
	refreshClaims := model.TokenClaims{
		Email: old.Email,
		Role:  old.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   old.ID,
			ID:        jti,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(RefreshTokenLifetime)),
		},
	}

	refreshToken, err := uc.signer.Sign(refreshClaims)
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	// Step 4: Store new refresh token session
	newKey := cache.KeyRefreshToken(old.Subject, jti)
	if err := uc.cache.Set(ctx, newKey, refreshClaims, RefreshTokenLifetime); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &dto.TokenPairs{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (uc *authJWTUsecase) cacheRefreshToken(ctx context.Context, dev model.SessionDevice) error {
	key := cache.KeyRefreshToken(dev.AccountID, dev.JTI)
	return uc.cache.Set(ctx, key, dev, RefreshTokenLifetime)
}

func (uc *authJWTUsecase) newDeviceSession(
	jti, accountID string,
	signAt time.Time,
	input dto.LoginJWTInput,
) model.SessionDevice {
	return model.SessionDevice{
		JTI:       jti,
		AccountID: accountID,
		UserAgent: input.UserAgent,
		IP:        input.IP,
		CreatedAt: signAt,
		ExpiresAt: signAt.Add(RefreshTokenLifetime),
	}
}

func (uc *authJWTUsecase) buildClaims(params tokenClaimsParams) model.TokenClaims {
	claims := model.TokenClaims{
		Email: params.account.Email,
		Role:  params.account.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   params.account.ID.String(),
			IssuedAt:  jwt.NewNumericDate(params.signAt),
			ExpiresAt: jwt.NewNumericDate(params.signAt.Add(params.exp)),
		},
	}

	if params.includeJTI {
		claims.ID = params.jti
	}

	return claims
}
