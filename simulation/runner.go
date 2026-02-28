package simulation

import (
	"cardflow21/engine"
	"cardflow21/strategy"
)

func RunSimulation(rounds int, decks int, penetration float64) (Stats, engine.StatsTracker) {
	stats := Stats{}
	shoe := engine.NewShoe(decks, penetration)
	strat := strategy.NewBasicStrategy()
	statsTracker := engine.StatsTracker{}
	for i := 0; i < rounds; i++ {
		result := engine.PlayHandWithStrategy(shoe, strat, &statsTracker)
		stats.Update(result)
	}


	return stats, statsTracker
}