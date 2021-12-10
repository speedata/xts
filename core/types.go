package core

type positioning uint8

const (
	positioningUnknown positioning = iota
	positioningAbsolute
	positioningGrid
)

func (p positioning) String() string {
	switch p {
	case positioningUnknown:
		return "unknown positioning"
	case positioningAbsolute:
		return "absolute positioning"
	case positioningGrid:
		return "grid positioning"
	}
	return ""
}
