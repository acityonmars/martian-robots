package display

import (
	"github.com/acityonmars/martian-robots/internal/mars"
)

type Display interface {
	Draw()
	Init(position mars.Position) error
}

type ConsoleDisplay struct {
}

func NewConsoleDisplay() ConsoleDisplay {
	return ConsoleDisplay{}
}

func (c ConsoleDisplay) Init(position mars.Position) error {
	return nil
}

func (c ConsoleDisplay) Draw() {

}
