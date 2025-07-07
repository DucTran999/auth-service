package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/DucTran999/auth-service/internal/model"
	"github.com/DucTran999/auth-service/internal/usecase/dto"
	"github.com/DucTran999/auth-service/internal/usecase/port"
	"github.com/DucTran999/auth-service/pkg/cache"
	"github.com/DucTran999/auth-service/pkg/hasher"
	"github.com/DucTran999/auth-service/pkg/signer"
	"github.com/google/uuid"
)

const (
	sessionDuration      = 60 * time.Minute
	AccessTokenLifetime  = 15 * time.Minute
	RefreshTokenLifetime = 7 * 24 * time.Hour
)

type AuthUseCase struct {
	hasher      hasher.Hasher
	signer      signer.TokenSigner
	cache       cache.Cache
	accountRepo port.AccountRepo
	sessionRepo port.SessionRepository
}

func NewAuthUseCase(
	hasher hasher.Hasher,
	signer signer.TokenSigner,
	cache cache.Cache,
	accountRepo port.AccountRepo,
	sessionRepo port.SessionRepository,
) *AuthUseCase {
	return &AuthUseCase{
		hasher:      hasher,
		signer:      signer,
		cache:       cache,
		accountRepo: accountRepo,
		sessionRepo: sessionRepo,
	}
}

// Login authenticates a user using email and password.
// It verifies credentials, checks account status, and creates a new session on success.
func (uc *AuthUseCase) Login(ctx context.Context, input dto.LoginInput) (*model.Session, error) {
	session, err := uc.tryReuseSession(ctx, input.CurrentSessionID)
	if err != nil {
		return nil, err
	}
	if session != nil {
		return session, nil
	}

	account, err := uc.findAccountByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	if err := uc.checkAccountActive(account); err != nil {
		return nil, err
	}

	if err := uc.verifyPassword(input.Password, account.PasswordHash); err != nil {
		return nil, err
	}

	return uc.createSession(ctx, account, input)
}

func (uc *AuthUseCase) Logout(ctx context.Context, sessionID string) error {
	// Fast check sessionID must uuid
	if _, err := uuid.Parse(sessionID); err != nil {
		return fmt.Errorf("logout: %w session=%s error=%w", model.ErrInvalidSessionID, sessionID, err)
	}

	// Remove session in cache
	_ = uc.cache.Del(ctx, cache.KeyFromSessionID(sessionID))

	// expire the session in db
	err := uc.sessionRepo.UpdateExpiresAt(ctx, sessionID, time.Now())
	if err != nil {
		return fmt.Errorf("logout: failed to update expires_at in DB for session=%s, error=%w", sessionID, err)
	}

	return nil
}

func (uc *AuthUseCase) LoginJWT(ctx context.Context, input dto.LoginJWTInput) (*dto.TokenPairs, error) {
	account, err := uc.findAccountByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	if err := uc.checkAccountActive(account); err != nil {
		return nil, err
	}

	if err := uc.verifyPassword(input.Password, account.PasswordHash); err != nil {
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

func (uc *AuthUseCase) RefreshToken(ctx context.Context, refreshToken string) (*dto.TokenPairs, error) {
	if refreshToken == "" {
		return nil, ErrInvalidCredentials
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
		return nil, ErrInvalidCredentials
	}

	tokens, err := uc.resignTokenPairs(ctx, *tokenClaim)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

// tryReuseSession checks if the current session is valid and updates its expiration.
// Returns the updated session if reusable; otherwise, returns nil.
func (uc *AuthUseCase) tryReuseSession(ctx context.Context, sessionID string) (*model.Session, error) {
	if sessionID == "" || sessionID == uuid.Nil.String() {
		return nil, nil
	}

	// Try get session from cache
	sessionKey := cache.KeyFromSessionID(sessionID)
	cachedSession := uc.getSessionFromCache(ctx, sessionKey)
	if cachedSession != nil {
		ttl, _ := uc.cache.TTL(ctx, sessionKey)

		// Example: Refresh TTL if less than 30% of sessionDuration remains
		if ttl < int64(sessionDuration.Seconds())/3 {
			_ = uc.cache.Expire(ctx, sessionKey, sessionDuration)
		}
		return cachedSession, nil
	}

	// Cache missing try to retrieve from DB
	session, err := uc.sessionRepo.FindByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if session == nil || session.IsExpired() {
		return nil, nil
	}

	// Re-cache the session after extend ttl
	_ = uc.cache.Set(ctx, sessionKey, session, sessionDuration)
	return session, nil
}

func (uc *AuthUseCase) findAccountByEmail(
	ctx context.Context,
	email string,
) (*model.Account, error) {
	account, err := uc.accountRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, ErrInvalidCredentials
	}

	return account, nil
}

func (uc *AuthUseCase) checkAccountActive(account *model.Account) error {
	if !account.IsActive {
		return model.ErrAccountDisabled
	}
	return nil
}

func (uc *AuthUseCase) verifyPassword(plain, hashed string) error {
	match, err := uc.hasher.ComparePasswordAndHash(plain, hashed)
	if err != nil {
		return err
	}
	if !match {
		return ErrInvalidCredentials
	}
	return nil
}

func (uc *AuthUseCase) createSession(
	ctx context.Context,
	account *model.Account,
	input dto.LoginInput,
) (*model.Session, error) {
	session := &model.Session{
		AccountID: account.ID,
		Account: model.Account{
			ID:       account.ID,
			Email:    account.Email,
			Role:     account.Role,
			IsActive: account.IsActive,
		},
		IPAddress: input.IP,
		UserAgent: input.UserAgent,
	}

	if err := uc.sessionRepo.Create(ctx, session); err != nil {
		return nil, err
	}

	_ = uc.cache.Set(ctx, cache.KeyFromSessionID(session.ID.String()), session, sessionDuration)

	return session, nil
}

func (uc *AuthUseCase) getSessionFromCache(
	ctx context.Context,
	sessionKey string,
) *model.Session {
	var session model.Session
	if err := uc.cache.GetInto(ctx, sessionKey, &session); err != nil {
		// If get session from cache got error just pass it.
		// Already has fallback form DB
		return nil
	}

	return &session
}

func (uc *AuthUseCase) signTokenPairs(
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

func (uc *AuthUseCase) resignTokenPairs(ctx context.Context, oldClaims model.TokenClaims) (*dto.TokenPairs, error) {
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

func (uc *AuthUseCase) cacheRefreshToken(ctx context.Context, dev model.SessionDevice) error {
	key := cache.KeyRefreshToken(dev.AccountID, dev.JTI)
	return uc.cache.Set(ctx, key, dev, RefreshTokenLifetime)
}

func (uc *AuthUseCase) newDeviceSession(
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
