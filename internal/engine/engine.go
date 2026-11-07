package engine

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
)

type Engine struct {
	world      *World
	robotPlans []RobotPlan // executed in order
	logger     *slog.Logger
}

// Run runs all robot plans in order, and prints out the result of each robot.
// Note that each plan is associated with a single robot and a series of commands that belong to the robot.
func (e Engine) Run(outputDest io.Writer) error {
	e.logger.Debug("run robot plans")

	var errs []error
	var robots []*Robot
	for _, plan := range e.robotPlans {
		result, err := plan.Execute(e.logger)
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
	e.logger.Debug("print final robot states")
	fmt.Fprintln(outputDest, "")
	for _, robot := range robots {
		if _, err := fmt.Fprintln(outputDest, robot); err != nil {
			return fmt.Errorf("write result: %w", err)
		}
	}

	return errors.Join(errs...)
}

func NewEngine(world *World, robotPlans []RobotPlan, logger *slog.Logger) Engine {
	return Engine{
		world:      world,
		robotPlans: robotPlans,
		logger:     logger,
	}
}
