package engine

import (
	"math"
	"math/rand"
	"time"
)

type Shoe struct {
	cards []Card
	cut_card_position int

	decks int
	penetration float64
	count int
	true_count int
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewShoe(decks int, penetration float64) *Shoe {
	if penetration < 0 {
		penetration = 0
	}
	if penetration > 1 {
		penetration = 1
	}

	s := &Shoe{decks: decks, penetration: penetration}
	s.init()
	return s
}

func (s *Shoe) init() {
	s.cards = make([]Card, 0, s.decks*52)
	for d := 0; d < s.decks; d++ {
		for suit := 0; suit < 4; suit++ {
			for rank := 1; rank <= 13; rank++ {
				s.cards = append(s.cards, Card{Suit: Suit(suit), Rank: Rank(rank)})
			}
		}
	}
	rand.Shuffle(len(s.cards), func(i, j int) {
		s.cards[i], s.cards[j] = s.cards[j], s.cards[i]
	})

	cut := int(math.Round(float64(len(s.cards)) * s.penetration))
	s.cut_card_position = len(s.cards) - cut
	s.count = 0
}

func (s *Shoe) Draw() Card {
	// TODO: shuffle the deck on the next hand after the cut card comes out
	if len(s.cards) == 0 || len(s.cards) <= s.cut_card_position  {
		s.init()
	}

	card := s.cards[len(s.cards)-1]
	s.cards = s.cards[:len(s.cards)-1]
	s.updateHiLoCount(card)
	return card
}

func (s *Shoe) updateHiLoCount(card Card) {
	decks_remaining := len(s.cards) / 52
	if decks_remaining == 0 {
		decks_remaining = 1
	}
	if card.Rank >= 2 && card.Rank <= 6 {
		s.count++
		s.true_count = s.count / decks_remaining
	} else if card.Rank >= 10 || card.Rank == 1 {
		s.count--
		s.true_count = s.count / decks_remaining
	}
}