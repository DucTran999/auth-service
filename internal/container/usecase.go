package container

import (
	"github.com/DucTran999/auth-service/internal/handler/background"
	"github.com/DucTran999/auth-service/internal/usecase"
	"github.com/DucTran999/auth-service/internal/usecase/port"
)

type useCases struct {
	auth              port.AuthUseCase
	account           port.AccountUseCase
	restSession       port.SessionUsecase
	backgroundSession background.SessionUsecase
}

func (c *Container) initUseCases() {
	accountUC := usecase.NewAccountUseCase(
		c.Hasher,
		c.repositories.account,
	)

	authUC := usecase.NewAuthUseCase(
		c.Hasher,
		c.Cache,
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
