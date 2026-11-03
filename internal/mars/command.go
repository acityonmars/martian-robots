package mars

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
	ApplyTo(robot Robot, world *World)
}

type CommandLetter int

type TurnLeft struct{}

func (l TurnLeft) ApplyTo(robot Robot, world *World) {}

type TurnRight struct{}

func (r TurnRight) ApplyTo(robot Robot, world *World) {}

type MoveForward struct{}

func (m MoveForward) ApplyTo(robot Robot, world *World) {}
