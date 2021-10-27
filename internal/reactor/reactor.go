package reactor

import (
	"github.com/unknowntpo/todos/internal/logger"
)

type Reactor struct {
	Logger logger.Logger
}

func NewReactor(logger logger.Logger) *Reactor {
	return &Reactor{Logger: logger}
}
