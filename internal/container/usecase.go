package container

import (
	"github.com/DucTran999/auth-service/internal/handler/background"
	"github.com/DucTran999/auth-service/internal/usecase/account"
	"github.com/DucTran999/auth-service/internal/usecase/auth"
	"github.com/DucTran999/auth-service/internal/usecase/port"
	"github.com/DucTran999/auth-service/internal/usecase/session"
	"github.com/DucTran999/auth-service/internal/usecase/shared"
)

type useCases struct {
	jwtAuth     port.JWTAuthUsecase
	sessionAuth port.SessionAuthUseCase

	account     port.AccountUseCase
	restSession port.SessionUsecase

	backgroundSession background.SessionUsecase
}

func (c *Container) initUseCases() {
	accountVerifier := shared.NewAccountVerifier(
		c.Hasher,
		c.repositories.account,
	)
	accountUC := account.NewAccountUseCase(
		c.Hasher,
		c.repositories.account,
	)

	sessionAuthUC := auth.NewAuthUseCase(
		c.Hasher,
		c.Cache,
		accountVerifier,
		c.repositories.account,
		c.repositories.session,
	)
	jwtAuthUC := auth.NewAuthJWTUsecase(c.Signer, c.Cache, accountVerifier)

	sessionUC := session.NewSessionUC(c.Cache, c.repositories.session)

	c.useCases = &useCases{
		account:           accountUC,
		jwtAuth:           jwtAuthUC,
		sessionAuth:       sessionAuthUC,
		restSession:       sessionUC,
		backgroundSession: sessionUC,
	}
}
