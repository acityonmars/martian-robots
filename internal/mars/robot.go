package mars

import (
	"fmt"
)

type Robot struct {
	ID          int
	Position    Position
	Orientation Orientation
}

type RobotPlan struct {
	Robot    Robot
	Commands []Command
}

func (rp RobotPlan) Execute() error {
	for _, command := range rp.Commands {
		fmt.Sprintf("%q executes %q", rp.Robot.ID, command)
	}
	return nil
}
