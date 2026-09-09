package tui

type ConfigSide int

const (
	LeftSide	ConfigSide = iota
	RightSide
)

type side[T any] struct {
	left	T
	right	T
}

func (s side[T]) Get(side ConfigSide) T {
	if side == LeftSide {
		return s.left
	}
	return s.right
}

func (s *side[T]) Set(side ConfigSide, v T) {
	if side == LeftSide {
		s.left = v
	} else {
		s.right = v
	}
}