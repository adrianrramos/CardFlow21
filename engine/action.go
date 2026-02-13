package engine

type Action int

const (
	Stand Action = iota
	Hit
	Double
	Split
	Surrender
)

func (a Action) String() string {
	switch a {
	case Hit:
		return "Hit"
	case Stand:
		return "Stand"
	case Double:
		return "Double"
	case Split:
		return "Split"
	case Surrender:
		return "Surrender"
	default:
		return "Unknown"
	}
}