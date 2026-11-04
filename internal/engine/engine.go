package engine

import (
	"log/slog"

	"github.com/acityonmars/martian-robots/internal/display"
	"github.com/acityonmars/martian-robots/internal/mars"
)

type Engine struct {
	logger *slog.Logger
}

func (e Engine) Run(mission mars.Mission) error {
	e.logger.Debug("Run missions")
	for _, plan := range mission.Plans {
		if err := plan.Execute(); err != nil {
			// TODO:
		}
	}

	return nil
}

func (e Engine) Init(disp display.Display) error {
	e.logger.Debug("Initialising engine")
	return nil
}

func NewEngine(logger *slog.Logger) Engine {
	return Engine{logger: logger}
}
