package main

import (
	"cardflow21/engine"
	"cardflow21/simulation"
	"cardflow21/strategy"
	"flag"
	"fmt"
	"math"
)

func main() {
	rounds := flag.Int("rounds", 100000, "number of rounds")
	decks := flag.Int("decks", 6, "number of decks in the shoe")
	penetration := flag.Float64("penetration", .85, "% of shoe before cut card comes out")
	noAdvanced := flag.Bool("no-advanced", false, "disable doubles and splits for baseline testing")
	flag.Parse()

	var strat engine.Strategy = strategy.NewBasicStrategy()
	if *noAdvanced {
		strat = strategy.NewNoAdvancedStrategy(strat)
	}

	stats := simulation.RunSimulationWithStrategy(*rounds, *decks, *penetration, strat)

	stats.PrintMetrics()

	if !*noAdvanced {
		fmt.Println("\n--- Validation vs Theoretical Baseline ---")
		validator := simulation.NewValidator()
		deviations := validator.CompareMetrics(&stats)
		largest := validator.FindLargestDeviation(&stats)
		
		fmt.Printf("\nLargest Deviation: %s\n", largest.Metric)
		fmt.Printf("  Expected: %.6f\n", largest.Expected)
		fmt.Printf("  Actual:   %.6f\n", largest.Actual)
		fmt.Printf("  Diff:     %.6f (%.2f%%)\n", largest.Difference, largest.PercentDiff)
		
		fmt.Println("\nAll Deviations (>1% shown):")
		for _, d := range deviations {
			if d.Metric == largest.Metric {
				fmt.Printf("  *** %s: Expected %.6f, Actual %.6f, Diff %.6f (%.2f%%)\n",
					d.Metric, d.Expected, d.Actual, d.Difference, d.PercentDiff)
			} else if d.Metric == "EV" || math.Abs(d.PercentDiff) > 1.0 {
				fmt.Printf("  %s: Expected %.6f, Actual %.6f, Diff %.6f (%.2f%%)\n",
					d.Metric, d.Expected, d.Actual, d.Difference, d.PercentDiff)
			}
		}
	}
}
