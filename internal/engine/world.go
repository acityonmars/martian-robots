package engine

import (
	"log/slog"
)

type Orientation int

const (
	North Orientation = iota
	East
	South
	West
)

func (o Orientation) String() string {
	switch o {
	case North:
		return "N"
	case East:
		return "E"
	case South:
		return "S"
	case West:
		return "W"
	default:
		return "unknown orientation"
	}
}

type Position struct {
	X int
	Y int
}

// World represents Mars where robots roam around.
type World struct {
	maxX     int // maximum x coordinate of the world
	maxY     int // maximum y coordinate of the world
	scents   map[Position]struct{}
	robotPos map[int]Position // robot ID -> position

	logger *slog.Logger
}

func NewWorld(upperRight Position, logger *slog.Logger) *World {
	return &World{
		upperRight.X,
		upperRight.Y,
		map[Position]struct{}{},
		map[int]Position{},
		logger,
	}
}

// MoveTo decides what happens when a robot tries to move to a position.
//
//  1. to is on the grid: the robot moves successfully.
//  2. to is off the grid and the robot's current cell has a scent: the move is ignored and the robot stays where it is.
//  3. to is off the grid with no scent: a scent is left on the current cell, the robot is marked Lost, and its position
//     stays at the last on-grid cell (which is what gets printed with LOST).
func (w *World) MoveTo(robot *Robot, to Position) {
	if w.HasPosition(to) {
		// 1: the robot moves successfully.
		robot.Pos = to
		return
	}
	// "to" is off grid.

	if w.HasScent(robot.Pos) {
		// 2: the move is ignored and the robot stays where it is.
		return
	}

	// 3: the robot is LOST and leaves a scent on the current position.
	w.scents[robot.Pos] = struct{}{}
	robot.Lost = true
}

// HasScent checks if "pos" has a scent on it.
func (w *World) HasScent(pos Position) bool {
	_, ok := w.scents[pos]
	return ok
}

// HasPosition checks if "to" position is inside the world.
func (w *World) HasPosition(to Position) bool {
	return to.X >= 0 && to.X <= w.maxX && to.Y >= 0 && to.Y <= w.maxY
}
