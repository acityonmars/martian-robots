package engine

import (
	"errors"
	"fmt"
	"log/slog"
)

type Engine struct {
	world   *World
	mission *Mission
	logger  *slog.Logger
}

func (e Engine) Run() error {
	e.logger.Debug("Run missions")

	var errs []error
	var robots []*Robot
	for _, plan := range e.mission.Plans {
		result, err := plan.Execute()
		if err != nil {
			errs = append(errs, err)
		}
		if result != nil {
			robots = append(robots, result)
		}
	}

	for _, robot := range robots {
		fmt.Println(robot)
	}

	return errors.Join(errs...)
}

func NewEngine(world *World, mission *Mission, logger *slog.Logger) Engine {
	return Engine{
		world:   world,
		mission: mission,
		logger:  logger,
	}
}
