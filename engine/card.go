package engine

// Remove Suits? Not needed for game logic
type Suit int
type Rank int

type Card struct {
	Suit Suit
	Rank Rank
}

var RankToIndex = map[int]int{
	1:  0, 2:  1, 3:  2, 4:  3, 5:  4, 6:  5, 7:  6, 8:  7, 9:  8, 10: 9, 11: 9, 12: 9, 13: 9,
}

func (c Card) Index() int {
	return RankToIndex[int(c.Rank)]
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
