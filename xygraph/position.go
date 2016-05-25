package xygraph

type position int

const (
	ON position = iota
	LEFT
	RIGHT
)

func (p position) opposite() position {
	switch p {
	case LEFT:
		return RIGHT
	case RIGHT:
		return LEFT
	default:
		return p
	}
}
