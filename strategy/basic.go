package strategy

import (
	"cardflow21/engine"
)

type BasicStrategy struct{}

func NewBasicStrategy() *BasicStrategy {
	return &BasicStrategy{}
}

func (b *BasicStrategy) Name() string {
	return "Basic Strategy"
}

type ChartType int

const (
	Hard ChartType = iota
	Soft
	Splits
)

type BasicStrategyChart map[ChartType]map[int][]string
type UncommonChart map[ChartType]map[string]string

// ┌─────────┬───────────────────────────────────────────────┐
// │ Symbol  │ Meaning                                       │
// ├─────────┼───────────────────────────────────────────────┤
// │   H     │ Hit                                           │
// │   S     │ Stand                                         │
// │   D     │ Double if allowed, otherwise Hit              │
// │  Ds     │ Double if allowed, otherwise Stand            │
// │  Sur    │ Surrender if allowed, otherwise Hit           │
// │   Y     │ Yes (split)                                   │
// │   N     │ No (do not split)                             │
// └─────────┴───────────────────────────────────────────────┘

// BS for DAS, H17, 4+ Decks
var Chart BasicStrategyChart = BasicStrategyChart{
	// 2,  3,   4,   5,   6,   7,   8,   9,   10,  A
	Hard: {
		17: {"S", "S", "S", "S", "S", "S", "S", "S", "S", "S"},
		16: {"S", "S", "S", "S", "S", "H", "H", "H", "H", "H"},
		15: {"S", "S", "S", "S", "S", "H", "H", "H", "H", "H"},
		14: {"S", "S", "S", "S", "S", "H", "H", "H", "H", "H"},
		13: {"S", "S", "S", "S", "S", "H", "H", "H", "H", "H"},
		12: {"H", "H", "S", "S", "S", "H", "H", "H", "H", "H"},
		11: {"D", "D", "D", "D", "D", "D", "D", "D", "D", "D"},
		10: {"D", "D", "D", "D", "D", "D", "D", "D", "H", "H"},
		9:  {"H", "D", "D", "D", "D", "H", "H", "H", "H", "H"},
	},

	Soft: {
		20: {"S", "S", "S", "S", "S", "S", "S", "S", "S", "S"},      // A,9 (soft 20)
		19: {"S", "S", "S", "S", "Ds", "S", "S", "S", "S", "S"},     // A,8 (soft 19)
		18: {"Ds", "Ds", "Ds", "Ds", "Ds", "S", "S", "H", "H", "H"}, // A,7 (soft 18)
		17: {"H", "D", "D", "D", "D", "H", "H", "H", "H", "H"},      // A,6 (soft 17)
		16: {"H", "H", "D", "D", "D", "H", "H", "H", "H", "H"},      // A,5 (soft 16)
		15: {"H", "H", "D", "D", "D", "H", "H", "H", "H", "H"},      // A,4 (soft 15)
		14: {"H", "H", "H", "D", "D", "H", "H", "H", "H", "H"},      // A,3 (soft 14)
		13: {"H", "H", "H", "D", "D", "H", "H", "H", "H", "H"},      // A,2 (soft 13)
	},

	Splits: {
		11: {"Y", "Y", "Y", "Y", "Y", "Y", "Y", "Y", "Y", "Y"},
		10: {"N", "N", "N", "N", "N", "N", "N", "N", "N", "N"},
		9:  {"Y", "Y", "Y", "Y", "Y", "N", "Y", "Y", "N", "N"},
		8:  {"Y", "Y", "Y", "Y", "Y", "Y", "Y", "Y", "Y", "Y"},
		7:  {"Y", "Y", "Y", "Y", "Y", "Y", "N", "N", "N", "N"},
		6:  {"Y", "Y", "Y", "Y", "Y", "N", "N", "N", "N", "N"},
		5:  {"N", "N", "N", "N", "N", "N", "N", "N", "N", "N"},
		4:  {"N", "N", "N", "Y", "Y", "N", "N", "N", "N", "N"},
		3:  {"Y", "Y", "Y", "Y", "Y", "Y", "N", "N", "N", "N"},
		2:  {"Y", "Y", "Y", "Y", "Y", "Y", "N", "N", "N", "N"},
	},
	// TODO: add surrender table
}

func (b *BasicStrategy) Decide(player engine.Hand, dealerUpCard engine.Card) engine.Action {
	if player.Value() <= 8 {
		return engine.Hit
	}
	if player.Value() >= 21 {
		return engine.Stand
	}
	if !player.IsSoft() && player.Value() > 17 {
		return engine.Stand
	}

	// TODO: Feature not ready
	// PAIR SPLITTING
	if player.IsPair() && Chart[Splits][player.Cards[0].Value()][dealerUpCard.Index()] == "Y" {
		return engine.Split
	}

	var action string
	// SOFT HANDS
	if player.IsSoft() {
		action = Chart[Soft][player.Value()][dealerUpCard.Index()]
	} else {
		action = Chart[Hard][player.Value()][dealerUpCard.Index()]
	}

	switch action {
	case "H":
		return engine.Hit
	case "S":
		return engine.Stand
	case "D":
		if player.IsTwoCardTotal() {
			return engine.Double
		}
		return engine.Hit
	case "Ds":
		if player.IsTwoCardTotal() {
			return engine.Double
		}
		return engine.Stand
	default:
		panic("Invalid action in chart")
	}
}
