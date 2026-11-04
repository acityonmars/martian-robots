package engine

var commandsByLetter = map[rune]Command{
	'L': TurnLeft{},
	'R': TurnRight{},
	'F': MoveForward{},
}

func CommandFor(letter rune) (Command, bool) {
	cmd, ok := commandsByLetter[letter]
	return cmd, ok
}

type Command interface {
	ApplyTo(robot *Robot, world *World)
}

type TurnLeft struct{}

func (l TurnLeft) ApplyTo(robot *Robot, _ *World) {
	switch robot.Ori {
	case North:
		robot.Ori = West
	case East:
		robot.Ori = North
	case South:
		robot.Ori = East
	case West:
		robot.Ori = South
	}
}

type TurnRight struct{}

func (r TurnRight) ApplyTo(robot *Robot, _ *World) {
	switch robot.Ori {
	case North:
		robot.Ori = East
	case East:
		robot.Ori = South
	case South:
		robot.Ori = West
	case West:
		robot.Ori = North
	}
}

type MoveForward struct{}

func (m MoveForward) ApplyTo(robot *Robot, world *World) {
	nextPos := robot.Pos

	switch robot.Ori {
	case North:
		nextPos.Y += 1
	case East:
		nextPos.X += 1
	case South:
		nextPos.Y -= 1
	case West:
		nextPos.X -= 1
	}

	world.MoveTo(robot, nextPos)
}
