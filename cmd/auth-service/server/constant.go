package server

type RunningEnvironment int

const (
	ProductionEnv RunningEnvironment = iota
	DevelopmentEnv
)

func (r RunningEnvironment) String() string {
	switch r {
	case DevelopmentEnv:
		return "dev"
	case ProductionEnv:
		return "prod"
	// Set to default value dev if env invalid
	default:
		return "dev"
	}
}
