package simulation

import (
	"cardflow21/engine"
	"cardflow21/strategy"
)

func RunSimulation(rounds int, decks int) Stats {
	stats := Stats{}
	shoe := engine.NewShoe(decks)
	strat := strategy.NewBasicStrategy()

	for i := 0; i < rounds; i++ {
		result := engine.PlayHandWithStrategy(shoe, strat)
		stats.Update(result)
	}

	return stats
}