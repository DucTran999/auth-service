package container

import (
	"github.com/DucTran999/auth-service/internal/usecase/account"
	"github.com/DucTran999/auth-service/internal/usecase/auth"
	"github.com/DucTran999/auth-service/internal/usecase/port"
	"github.com/DucTran999/auth-service/internal/usecase/session"
	"github.com/DucTran999/auth-service/internal/usecase/shared"
)

type useCases struct {
	jwtAuth     port.AuthJWTUsecase
	sessionAuth port.AuthSessionUsecase

	account port.AccountUsecase
	session port.SessionUsecase

	backgroundSession port.SessionMaintenanceUsecase
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

	sessionAuthUC := auth.NewAuthSessionUsecase(
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
		session:           sessionUC,
		backgroundSession: sessionUC,
	}
}
