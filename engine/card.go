package engine

// Remove Suits? Not needed for game logic
type Suit int
type Rank int

type Card struct {
	Suit Suit
	Rank Rank
}

// Takes the card rank and converts to index for strategy chart column
// ie. 2, 3, 4, ... , 10, A -> 0, 1, 2, ... 8, 9
var RankToIndex = map[int]int{
	2: 0, 3: 1, 4: 2, 5: 3, 6: 4, 7: 5, 8: 6, 9: 7, 10: 8, 11: 8, 12: 8, 13: 8, 1: 9,
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
