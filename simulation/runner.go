package simulation

import (
	"cardflow21/engine"
	"cardflow21/strategy"
	"fmt"
	"os"
	"strings"
)

func RunSimulation(rounds int, decks int, penetration float64, use_true_count bool, strategy_name strategy.StrategyName) (Stats, engine.StatsTracker) {
	stats := Stats{}
	shoe := engine.NewShoe(decks, penetration)
	strat := strategy.NewBasicStrategy(strategy_name)
	statsTracker := engine.StatsTracker{}
	print_progress := false // for dev
	const segments = 10
	lastBucket := -1
	for i := 0; i < rounds; i++ {
		result := engine.PlayRound(shoe, strat, &statsTracker, use_true_count)
		stats.Update(result)

		if print_progress {
			printProgress(segments, rounds, lastBucket, i)
		}
	}

	return stats, statsTracker
}

func printProgress(segments, rounds, lastBucket, i int) {

	done := i + 1
	bucket := done * segments / rounds
	if bucket <= lastBucket {
		return
	}
	lastBucket = bucket
	filled := bucket
	if filled > segments {
		filled = segments
	}
	bar := strings.Repeat("█", filled) + strings.Repeat("░", segments-filled)
	pct := bucket * 10
	if bucket == segments {
		pct = 100
	}
	fmt.Fprintf(os.Stderr, "\r[%s] %3d%%  %d / %d rounds", bar, pct, done, rounds)
	if bucket == segments {
		fmt.Fprintln(os.Stderr)
	}
}
