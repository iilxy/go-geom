package graph

type LabelPos int

const (
	OnLabel LabelPos = iota
	LeftOfLabel
	RightOfLabel
)

func (pos LabelPos) opposite() LabelPos {
	switch pos {
	case LeftOfLabel:
		return RightOfLabel
	case RightOfLabel:
		return LeftOfLabel
	default:
		return pos
	}
}
