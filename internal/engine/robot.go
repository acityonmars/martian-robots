package engine

import (
	"errors"
	"fmt"
	"io"
)

type Robot struct {
	ID   int
	Pos  Position
	Ori  Orientation // Ori is the direction it faces
	Lost bool        // Lost indicates if it is a lost robot
}

func (r Robot) Print(stdout io.Writer) {
	fmt.Println("test")
}

func (r Robot) String() string {
	line := fmt.Sprintf("%d %d %s", r.Pos.X, r.Pos.Y, r.Ori)
	if r.Lost {
		line += " LOST"
	}
	return line
}

type RobotPlan struct {
	Robot    *Robot
	Commands []Command
	World    *World
}

func (rp RobotPlan) Execute() (*Robot, error) {
	if rp.Robot == nil {
		return nil, errors.New("robot isn't set on the plan")
	}

	if rp.World == nil {
		return nil, errors.New("world isn't set on the plan")
	}

	for _, command := range rp.Commands {
		if rp.Robot.Lost {
			break
		}
		command.ApplyTo(rp.Robot, rp.World)
	}
	return rp.Robot, nil
}
