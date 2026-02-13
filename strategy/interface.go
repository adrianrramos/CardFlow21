package strategy

import "cardflow21/engine"

type Strategy interface {
	Decide(player engine.Hand, dealerUpCard engine.Card) engine.Action
	Name() string
}