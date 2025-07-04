package container

import (
	"github.com/DucTran999/auth-service/internal/repository"
	"github.com/DucTran999/auth-service/internal/usecase/port"
)

type repositories struct {
	account port.AccountRepo
	session port.SessionRepository
}

func (c *Container) initRepositories() {
	c.repositories = &repositories{
		account: repository.NewAccountRepo(c.AuthDB.DB()),
		session: repository.NewSessionRepository(c.AuthDB.DB()),
	}
}
