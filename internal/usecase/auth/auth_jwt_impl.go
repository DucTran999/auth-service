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
	"github.com/google/uuid"
)

const (
	AccessTokenLifetime  = 15 * time.Minute
	RefreshTokenLifetime = 7 * 24 * time.Hour
)

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

func (uc *authJWTUsecase) Login(ctx context.Context, input dto.LoginJWTInput) (*dto.TokenPairs, error) {
	account, err := uc.accountVerifier.Verify(ctx, input.Email, input.Password)
	if err != nil {
		return nil, err
	}

	jti := uuid.NewString()
	signAt := time.Now()

	tokens, err := uc.signTokenPairs(jti, signAt, account)
	if err != nil {
		return nil, err
	}

	sessionDevice := uc.newDeviceSession(jti, account.ID.String(), signAt, input)
	if err := uc.cacheRefreshToken(ctx, sessionDevice); err != nil {
		return nil, err
	}

	return tokens, nil
}

func (uc *authJWTUsecase) RefreshToken(ctx context.Context, refreshToken string) (*dto.TokenPairs, error) {
	if refreshToken == "" {
		return nil, model.ErrInvalidCredentials
	}

	claims, err := uc.signer.Parse(refreshToken)
	if err != nil {
		return nil, err
	}

	tokenClaim, err := model.MapClaimsToTokenClaims(*claims)
	if err != nil {
		return nil, err
	}

	key := cache.KeyRefreshToken(tokenClaim.ID.String(), tokenClaim.JTI)
	ok, err := uc.cache.Has(ctx, key)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, model.ErrInvalidCredentials
	}

	tokens, err := uc.resignTokenPairs(ctx, *tokenClaim)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (uc *authJWTUsecase) RevokeRefreshToken(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return model.ErrInvalidCredentials
	}

	claims, err := uc.signer.Parse(refreshToken)
	if err != nil {
		return err
	}

	tokenClaim, err := model.MapClaimsToTokenClaims(*claims)
	if err != nil {
		return err
	}

	key := cache.KeyRefreshToken(tokenClaim.ID.String(), tokenClaim.JTI)
	_ = uc.cache.Del(ctx, key)

	return nil
}

func (uc *authJWTUsecase) signTokenPairs(
	jti string,
	signAt time.Time,
	account *model.Account,
) (*dto.TokenPairs, error) {
	// Access token claims
	accessClaims := model.TokenClaims{
		ID:        account.ID,
		Email:     account.Email,
		Role:      account.Role,
		IssuedAt:  signAt.Unix(),
		ExpiresAt: signAt.Add(AccessTokenLifetime).Unix(),
	}

	accessToken, err := uc.signer.SignAccessToken(accessClaims.ToMapClaims())
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// Refresh token claims
	refreshClaims := model.TokenClaims{
		ID:        account.ID,
		Email:     account.Email,
		Role:      account.Role,
		IssuedAt:  signAt.Unix(),
		ExpiresAt: signAt.Add(RefreshTokenLifetime).Unix(),
		JTI:       jti,
	}

	refreshToken, err := uc.signer.SignRefreshToken(refreshClaims.ToMapClaims())
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return &dto.TokenPairs{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (uc *authJWTUsecase) resignTokenPairs(ctx context.Context, oldClaims model.TokenClaims) (*dto.TokenPairs, error) {
	jti := uuid.NewString()
	now := time.Now()

	// Access token claims
	accessClaims := oldClaims
	accessClaims.JTI = ""
	accessClaims.IssuedAt = now.Unix()
	accessClaims.ExpiresAt = now.Add(AccessTokenLifetime).Unix()

	accessToken, err := uc.signer.SignAccessToken(accessClaims.ToMapClaims())
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// Refresh token claims (rotate JTI and extend lifetime)
	refreshClaims := oldClaims
	refreshClaims.JTI = jti
	refreshClaims.IssuedAt = now.Unix()
	refreshClaims.ExpiresAt = now.Add(RefreshTokenLifetime).Unix()

	refreshToken, err := uc.signer.SignRefreshToken(refreshClaims.ToMapClaims())
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	key := cache.KeyRefreshToken(refreshClaims.ID.String(), jti)
	if err = uc.cache.Set(ctx, key, refreshClaims, RefreshTokenLifetime); err != nil {
		return nil, err
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
