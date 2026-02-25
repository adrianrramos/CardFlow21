package simulation

import (
	"cardflow21/engine"
	"cardflow21/strategy"
)

func RunSimulation(rounds int, decks int, penetration float64) Stats {
	stats := Stats{}
	shoe := engine.NewShoe(decks, penetration)
	strat := strategy.NewBasicStrategy()

	for i := 0; i < rounds; i++ {
		result := engine.PlayHandWithStrategy(shoe, strat)
		stats.Update(result)
	}

	return stats
}