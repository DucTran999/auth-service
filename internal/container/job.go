package container

import "github.com/DucTran999/auth-service/internal/handler/background"

func (c *Container) initJobs() {
	c.CleanupSessionHandler = background.NewSessionCleaner(c.Logger, c.useCases.backgroundSession)
}
