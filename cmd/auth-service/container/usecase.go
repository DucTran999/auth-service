package container

import (
	"github.com/DucTran999/auth-service/internal/handler/background"
	"github.com/DucTran999/auth-service/internal/handler/rest"
	"github.com/DucTran999/auth-service/internal/usecase"
)

type useCases struct {
	auth              rest.AuthUseCase
	account           rest.AccountUseCase
	restSession       rest.SessionUsecase
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
