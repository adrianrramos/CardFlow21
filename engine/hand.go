package engine

type Hand struct {
	Cards []Card
	SplitCount int
	IsDoubled bool
	IsBust bool
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

func (h Hand) IsPair() bool {
	// Use Index() to compare the rank of the cards instead of Rank directly
	return len(h.Cards) == 2 && h.Cards[0].Index() == h.Cards[1].Index()
}

func (h Hand) IsTwoCardTotal() bool {
	return len(h.Cards) == 2
}

func (h *Hand) CheckBust() bool {
	h.IsBust = h.Value() > 21
	return h.IsBust
}

func (h Hand) IsBlackjack() bool {
	return len(h.Cards) == 2 && h.Value() == 21
}

func (h *Hand) Doubled() {
	h.IsDoubled = true
}

// Dealer is showing an Ace
func (h Hand) OffersInsurance() bool {
	return len(h.Cards) == 2 && h.Cards[0].Rank == 1
}