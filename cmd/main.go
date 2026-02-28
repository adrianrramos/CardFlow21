package main

import (
	"cardflow21/engine"
	"cardflow21/simulation"
	"cardflow21/strategy"
	"flag"
	"fmt"
	"math"
	"strings"
)

func main() {
	rounds := flag.Int("rounds", 100000, "number of rounds")
	decks := flag.Int("decks", 6, "number of decks in the shoe")
	penetration := flag.Float64("penetration", .85, "% of shoe before cut card comes out")
	noAdvanced := flag.Bool("no-advanced", false, "disable doubles and splits for baseline testing")
	traceBlackjacks := flag.Bool("trace-blackjacks", false, "trace 50 blackjack hands for payout validation")
	diagnostic := flag.Bool("diagnostic", false, "enable diagnostic mode with phase-by-phase validation")
	flag.Parse()

	var strat engine.Strategy = strategy.NewBasicStrategy()
	if *noAdvanced {
		strat = strategy.NewNoAdvancedStrategy(strat)
	}

	var stats simulation.Stats
	var handResults []engine.HandResult
	
	if *traceBlackjacks {
		stats, handResults = simulation.RunSimulationWithTracing(*rounds, *decks, *penetration, strat)
	} else {
		stats = simulation.RunSimulationWithStrategy(*rounds, *decks, *penetration, strat)
	}

	stats.PrintMetrics()

	// Phase 1: Core Statistics Validation Table
	if *diagnostic {
		fmt.Println("\n=== Phase 1: Core Statistics Validation ===")
		validator := simulation.NewValidator()
		deviations := validator.CompareMetrics(&stats)
		
		fmt.Println("\nMetric Comparison Table:")
		fmt.Printf("%-30s %12s %12s %12s %10s\n", "Metric", "Expected", "Actual", "Difference", "Status")
		fmt.Println(strings.Repeat("-", 80))
		for _, d := range deviations {
			status := "✓"
			if math.Abs(d.PercentDiff) > 2.0 {
				status = "✗ >2%"
			}
			fmt.Printf("%-30s %12.6f %12.6f %12.6f %10s\n", 
				d.Metric, d.Expected, d.Actual, d.Difference, status)
		}
	}

	// Phase 2: Blackjack payout validation
	if *traceBlackjacks {
		traces := simulation.TraceBlackjacks(handResults, 50)
		simulation.PrintBlackjackTraces(traces)
		
		validator := simulation.NewValidator()
		valid, mismatches := validator.ValidateBlackjackPayouts(traces)
		if !valid {
			fmt.Printf("\n*** BLACKJACK PAYOUT VALIDATION FAILED: %d mismatches found ***\n", len(mismatches))
		} else {
			fmt.Println("\n✓ Blackjack payout validation passed")
		}
	}

	// Phase 3: Split logic validation (diagnostic mode)
	if *diagnostic && !*noAdvanced {
		splitValidator := simulation.NewSplitValidator(&stats)
		splitValidator.ValidateSplitDecisions()
		splitValidator.CalculateSplitEV()
		splitValidator.ValidateDAS()
		splitValidator.ValidateResplitRules()
	}

	// Phase 5: Strategy matrix validation (diagnostic mode)
	if *diagnostic {
		strategyValidator := simulation.NewStrategyValidator()
		strategyValidator.ValidateHardTotals()
		strategyValidator.ValidateSoftTotals()
		strategyValidator.ValidatePairSplits()
		strategyValidator.ValidateDoubleRules()
	}

	// Phase 6: Variance validation (always shown in metrics)
	// Variance validation is included in PrintMetrics() output

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
