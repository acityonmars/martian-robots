package display

import (
	"log/slog"

	"github.com/acityonmars/martian-robots/internal/mars"
)

type Display interface {
	Draw()
	Init(position mars.Position) error
}

type ConsoleDisplay struct {
	logger *slog.Logger
}

func NewConsoleDisplay(logger *slog.Logger) ConsoleDisplay {
	return ConsoleDisplay{logger: logger}
}

func (c ConsoleDisplay) Init(position mars.Position) error {
	c.logger.Debug("Initialising display")
	return nil
}

func (c ConsoleDisplay) Draw() {
}
