package engine

import (
	"errors"
	"fmt"
	"log/slog"
)

// Robot represents a single robot
type Robot struct {
	ID   int
	Pos  Position
	Ori  Orientation // Ori is the direction it faces
	Lost bool        // Lost indicates if it is a lost robot
}

func (r Robot) String() string {
	line := fmt.Sprintf("%d %d %s", r.Pos.X, r.Pos.Y, r.Ori)
	if r.Lost {
		line += " LOST" // e.g. "3 3 N LOST"
	}
	return line
}

// RobotPlan represents a single robot's commands.
type RobotPlan struct {
	Robot    *Robot
	Commands []Command
	World    *World
}

// Execute runs the commands sequentially.
func (rp RobotPlan) Execute(logger *slog.Logger) (*Robot, error) {
	if rp.Robot == nil {
		return nil, errors.New("robot isn't set on the plan")
	}

	if rp.World == nil {
		return nil, errors.New("world isn't set on the plan")
	}

	logger.Debug(fmt.Sprintf("robot id: %d apply commands", rp.Robot.ID))
	for _, command := range rp.Commands {
		if rp.Robot.Lost {
			logger.Debug(fmt.Sprintf("robot id: %d LOST", rp.Robot.ID))
			break
		}
		command.ApplyTo(rp.Robot, rp.World)
	}
	return rp.Robot, nil
}
