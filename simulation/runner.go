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

// RunSimulationWithTracing runs simulation and collects hand results for tracing
func RunSimulationWithTracing(rounds int, decks int, penetration float64, strat engine.Strategy) (Stats, []engine.HandResult) {
	stats := Stats{}
	results := []engine.HandResult{}
	shoe := engine.NewShoe(decks, penetration)

	for i := 0; i < rounds; i++ {
		result := engine.PlayHandWithStrategyDetailed(shoe, strat)
		stats.UpdateDetailed(result)
		results = append(results, result)
	}

	return stats, results
}