package mars

type Orientation int

const (
	North Orientation = iota
	East
	South
	West
)

type Position struct {
	X int
	Y int
}
type World struct {
	maxX   int
	maxY   int
	scents map[Position]struct{}
}
