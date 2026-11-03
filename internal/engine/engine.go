package engine

import (
	"github.com/acityonmars/martian-robots/internal/display"
	"github.com/acityonmars/martian-robots/internal/mars"
)

type Engine struct{}

func (e Engine) Run(mission mars.Mission) error {
	for _, plan := range mission.Plans {
		if err := plan.Execute(); err != nil {
			// TODO:
		}
	}

	return nil
}

func (e Engine) Init(disp display.Display) error {
	return nil
}

func NewEngine() Engine {
	return Engine{}
}
