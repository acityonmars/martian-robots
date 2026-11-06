package engine

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
)

// Mission is the ordered list of robot plans the engine runs.
type Mission struct {
	Plans []RobotPlan
}

type Engine struct {
	world   *World
	mission *Mission
	logger  *slog.Logger
}

// Run runs all mission plans, and prints out the result of each robot.
// Note that each plan is associated with a single robot and a series of commands that belong to the robot.
func (e Engine) Run(outputDest io.Writer) error {
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

	// Print out the result - the final state of each robot per line.
	// e.g.
	// 1 1 E
	// 3 3 N LOST
	// 2 3 S
	for _, robot := range robots {
		if _, err := fmt.Fprintln(outputDest, robot); err != nil {
			return fmt.Errorf("write result: %w", err)
		}
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
