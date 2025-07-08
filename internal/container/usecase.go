package container

import (
	"github.com/DucTran999/auth-service/internal/handler/background"
	"github.com/DucTran999/auth-service/internal/usecase"
	"github.com/DucTran999/auth-service/internal/usecase/auth"
	"github.com/DucTran999/auth-service/internal/usecase/port"
	"github.com/DucTran999/auth-service/internal/usecase/shared"
)

type useCases struct {
	jwtAuth           port.JWTAuthUsecase
	sessionAuth       port.AuthUseCase
	account           port.AccountUseCase
	restSession       port.SessionUsecase
	backgroundSession background.SessionUsecase
}

func (c *Container) initUseCases() {
	accountVerifier := shared.NewAccountVerifier(
		c.Hasher,
		c.repositories.account,
	)
	accountUC := usecase.NewAccountUseCase(
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

	sessionUC := usecase.NewSessionUC(c.Cache, c.repositories.session)

	c.useCases = &useCases{
		account:           accountUC,
		jwtAuth:           jwtAuthUC,
		sessionAuth:       sessionAuthUC,
		restSession:       sessionUC,
		backgroundSession: sessionUC,
	}
}
