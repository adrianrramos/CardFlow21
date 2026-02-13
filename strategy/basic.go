package strategy

import "cardflow21/engine"

type BasicStrategy struct{}

func NewBasicStrategy() *BasicStrategy {
	return &BasicStrategy{}
}

func (b *BasicStrategy) Name() string {
	return "Basic Strategy (Simplified)"
}

func (b *BasicStrategy) Decide(player engine.Hand, dealerUpCard engine.Card) engine.Action {
	// TODO: replace with BJA H17 4+ Deck BS
	if player.Value() < 17 {
		return engine.Hit
	}
	return engine.Stand
}