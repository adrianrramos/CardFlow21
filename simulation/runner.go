package simulation

import (
	"cardflow21/engine"
	"cardflow21/strategy"
)

func RunSimulation(rounds int, decks int, penetration float64) Stats {
	return RunSimulationWithStrategy(rounds, decks, penetration, strategy.NewBasicStrategy())
}

func RunSimulationWithStrategy(rounds int, decks int, penetration float64, strat engine.Strategy) Stats {
	stats := Stats{}
	shoe := engine.NewShoe(decks, penetration)

	for i := 0; i < rounds; i++ {
		result := engine.PlayHandWithStrategyDetailed(shoe, strat)
		stats.UpdateDetailed(result)
	}

	return stats
}