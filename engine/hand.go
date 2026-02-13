package engine

type Hand struct {
	Cards []Card
}

func (h *Hand) AddCard(c Card) {
	h.Cards = append(h.Cards, c)
}

func (h Hand) Value() int {
	total := 0
	aces := 0
	for _, card := range h.Cards {
		total += card.Value()
		if card.Rank == 1 {
			aces++
		}
	}
	for total > 21 && aces > 0 {
		total -= 10
		aces--
	}
	return total
}

func (h Hand) IsBust() bool {
	return h.Value() > 21
}

func (h Hand) IsBlackjack() bool {
	return len(h.Cards) == 2 && h.Value() == 21
}