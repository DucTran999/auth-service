package container

import (
	"github.com/DucTran999/auth-service/internal/handler/background"
	"github.com/DucTran999/auth-service/internal/usecase"
	"github.com/DucTran999/auth-service/internal/usecase/port"
	"github.com/DucTran999/auth-service/internal/usecase/shared"
)

type useCases struct {
	auth              port.AuthUseCase
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

	authUC := usecase.NewAuthUseCase(
		c.Hasher,
		c.Signer,
		c.Cache,
		accountVerifier,
		c.repositories.account,
		c.repositories.session,
	)

	sessionUC := usecase.NewSessionUC(c.Cache, c.repositories.session)

	c.useCases = &useCases{
		account:           accountUC,
		auth:              authUC,
		restSession:       sessionUC,
		backgroundSession: sessionUC,
	}
}
