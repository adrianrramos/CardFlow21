package simulation

import (
	"cardflow21/engine"
	"cardflow21/strategy"
)

func RunSimulation(rounds int, decks int, penetration float64, use_true_count bool, strategy_name strategy.StrategyName) (Stats, engine.StatsTracker) {
	stats := Stats{}
	shoe := engine.NewShoe(decks, penetration)
	strat := strategy.NewBasicStrategy(strategy_name)
	statsTracker := engine.StatsTracker{}
	for i := 0; i < rounds; i++ {
		result := engine.PlayRound(shoe, strat, &statsTracker, use_true_count)
		stats.Update(result)
	}


	return stats, statsTracker
}