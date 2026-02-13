package engine

import (
	"math/rand"
	"time"
)

type Shoe struct {
	cards []Card
	decks int
}

func NewShoe(decks int) *Shoe {
	s := &Shoe{decks: decks}
	s.init()
	return s
}

func (s *Shoe) init() {
	rand.Seed(time.Now().UnixNano())
	s.cards = []Card{}
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
}

func (s *Shoe) Draw() Card {
	if len(s.cards) == 0 {
		s.init()
	}

	card := s.cards[len(s.cards)-1]
	s.cards = s.cards[:len(s.cards)-1]
	return card
}