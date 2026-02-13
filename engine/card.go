package engine

type Suit int
type Rank int

type Card struct {
	Suit Suit
	Rank Rank
}

func (c Card) Value() int {
	if c.Rank >= 10 {
		return 10
	}
	if c.Rank == 1 {
		return 11
	}
	return int(c.Rank)
}
