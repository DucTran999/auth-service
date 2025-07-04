package container

import "github.com/DucTran999/auth-service/internal/handler/background"

type jobs struct {
	SessionCleaner background.SessionCleaner
}

func (c *container) initJobs() {
	c.jobs = &jobs{
		SessionCleaner: background.NewSessionCleaner(c.logger, c.useCases.backgroundSession),
	}
}
