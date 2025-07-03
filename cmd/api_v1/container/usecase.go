package container

import (
	"github.com/DucTran999/auth-service/internal/v1/handler/background"
	"github.com/DucTran999/auth-service/internal/v1/usecase"
	"github.com/DucTran999/auth-service/internal/v1/usecase/port"
)

type useCases struct {
	auth              port.AuthUseCase
	account           port.AccountUseCase
	restSession       port.SessionUsecase
	backgroundSession background.SessionUsecase
}

func (c *container) initUseCases() {
	accountUC := usecase.NewAccountUseCase(
		c.hasher,
		c.repositories.account,
	)

	authUC := usecase.NewAuthUseCase(
		c.hasher,
		c.cache,
		c.repositories.account,
		c.repositories.session,
	)

	sessionUC := usecase.NewSessionUC(c.cache, c.repositories.session)

	c.useCases = &useCases{
		account:           accountUC,
		auth:              authUC,
		restSession:       sessionUC,
		backgroundSession: sessionUC,
	}
}
