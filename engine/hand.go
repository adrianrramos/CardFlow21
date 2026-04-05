package engine

type Hand struct {
	Cards      []Card
	SplitCount int
	IsDoubled  bool
}

func (h *Hand) AddCard(c Card) {
	h.Cards = append(h.Cards, c)
}

func (h *Hand) RemoveCard() {
	if len(h.Cards) > 0 {
		h.Cards = h.Cards[:len(h.Cards)-1]
	}
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

func (h Hand) IsSoft() bool {
	total := 0
	aces := 0
	for _, card := range h.Cards {
		total += card.Value()
		if card.Rank == 1 {
			aces++
		}
	}
	return aces > 0 && total <= 21
}

func (h Hand) IsSplitHandWithAce() bool {
	return h.IsSoft() && h.SplitCount >= 1
}

func (h Hand) CanSplit() bool {
	return h.SplitCount < 4 && h.IsPair()
}

func (h Hand) CheckBust() bool {
	return h.Value() > 21
}
func (h *Hand) Doubled() {
	h.IsDoubled = true
}

func (h Hand) IsPair() bool {
	// Use Index() to compare the rank of the cards instead of Rank directly
	return len(h.Cards) == 2 && h.Cards[0].Value() == h.Cards[1].Value()
}

func (h Hand) IsTwoCardTotal() bool {
	return len(h.Cards) == 2
}

func (h Hand) IsBlackjack() bool {
	return len(h.Cards) == 2 && h.Value() == 21
}

// Dealer is showing an Ace
func (h Hand) IsShowingAce() bool {
	return len(h.Cards) == 2 && h.Cards[0].Rank == 1
}
