package simulation

import (
	"fmt"
	"math"
)

// SplitValidator validates split logic and decisions
type SplitValidator struct {
	stats *Stats
}

func NewSplitValidator(stats *Stats) *SplitValidator {
	return &SplitValidator{stats: stats}
}

// ValidateSplitDecisions compares pair opportunities to splits taken
func (sv *SplitValidator) ValidateSplitDecisions() {
	fmt.Println("\n=== Split Decision Validation ===")
	
	if sv.stats.PairOpportunities == 0 {
		fmt.Println("No pair opportunities recorded")
		return
	}
	
	splitRate := float64(sv.stats.SplitsAttempted) / float64(sv.stats.PairOpportunities) * 100.0
	expectedRate := 0.30 // Rough estimate - should be ~30% of pairs split (varies by rank)
	
	fmt.Printf("Pair Opportunities: %d\n", sv.stats.PairOpportunities)
	fmt.Printf("Splits Taken: %d\n", sv.stats.SplitsAttempted)
	fmt.Printf("Split Rate: %.2f%% (expected ~30%%)\n", splitRate)
	
	if math.Abs(splitRate-expectedRate) > 5.0 {
		fmt.Printf("*** WARNING: Split rate deviates by %.2f%% from expected ***\n", math.Abs(splitRate-expectedRate))
	}
	
	// Validate by rank
	if sv.stats.PairOpportunitiesByRank != nil {
		fmt.Println("\nSplit Decisions by Rank:")
		for rank := 1; rank <= 10; rank++ {
			opps := sv.stats.PairOpportunitiesByRank[rank]
			splits := sv.stats.SplitsByRank[rank]
			if opps > 0 {
				rate := float64(splits) / float64(opps) * 100.0
				rankName := getRankName(rank)
				fmt.Printf("  %s: %d opportunities, %d splits (%.1f%%)\n", rankName, opps, splits, rate)
			}
		}
	}
}

// CalculateSplitEV computes split contribution to overall EV
func (sv *SplitValidator) CalculateSplitEV() {
	fmt.Println("\n=== Split EV Analysis ===")
	
	if sv.stats.SplitRounds == 0 {
		fmt.Println("No split rounds recorded")
		return
	}
	
	splitEV := sv.stats.SplitProfit / float64(sv.stats.SplitRounds)
	totalEV := sv.stats.Mean
	
	fmt.Printf("Split Rounds: %d\n", sv.stats.SplitRounds)
	fmt.Printf("Split Hands Total: %d\n", sv.stats.SplitHandsTotal)
	fmt.Printf("Split Profit: %.2f\n", sv.stats.SplitProfit)
	fmt.Printf("Split EV per Round: %.6f\n", splitEV)
	fmt.Printf("Overall EV: %.6f\n", totalEV)
	
	// Split contribution percentage
	if sv.stats.TotalHands > 0 {
		splitContribution := (sv.stats.SplitProfit / float64(sv.stats.TotalHands)) / totalEV * 100.0
		fmt.Printf("Split Contribution to EV: %.2f%%\n", splitContribution)
	}
}

// ValidateDAS checks if Double After Split is working
func (sv *SplitValidator) ValidateDAS() {
	fmt.Println("\n=== DAS (Double After Split) Validation ===")
	
	if sv.stats.SplitRounds == 0 {
		fmt.Println("No split rounds to validate DAS")
		return
	}
	
	dasRate := float64(sv.stats.DoublesAfterSplit) / float64(sv.stats.SplitHandsTotal) * 100.0
	
	fmt.Printf("Split Hands Total: %d\n", sv.stats.SplitHandsTotal)
	fmt.Printf("Doubles After Split: %d\n", sv.stats.DoublesAfterSplit)
	fmt.Printf("DAS Rate: %.2f%%\n", dasRate)
	
	// Expected DAS rate is typically 5-10% of split hands
	if dasRate < 1.0 {
		fmt.Println("*** WARNING: DAS rate very low - may indicate DAS not enabled ***")
	} else if dasRate > 15.0 {
		fmt.Println("*** WARNING: DAS rate very high - may indicate accounting error ***")
	} else {
		fmt.Println("✓ DAS rate within expected range")
	}
}

// ValidateResplitRules checks resplit implementation
func (sv *SplitValidator) ValidateResplitRules() {
	fmt.Println("\n=== Resplit Rules Validation ===")
	
	// Calculate average hands per split round
	if sv.stats.SplitRounds > 0 {
		avgHandsPerSplit := float64(sv.stats.SplitHandsTotal) / float64(sv.stats.SplitRounds)
		fmt.Printf("Average Hands per Split Round: %.2f\n", avgHandsPerSplit)
		
		// Expected: 2.0 for single splits, 2.1-2.2 if resplits occur
		if avgHandsPerSplit < 1.9 {
			fmt.Println("*** WARNING: Average hands per split too low - resplits may not be working ***")
		} else if avgHandsPerSplit > 2.5 {
			fmt.Println("*** WARNING: Average hands per split too high - may indicate accounting error ***")
		} else {
			fmt.Println("✓ Resplit behavior appears normal")
		}
	}
}

func getRankName(rank int) string {
	switch rank {
	case 1:
		return "Aces"
	case 10:
		return "10s"
	default:
		return fmt.Sprintf("%ds", rank)
	}
}
