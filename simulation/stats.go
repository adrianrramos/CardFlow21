package simulation

import (
	"cardflow21/engine"
	"fmt"
	"math"
)

type Stats struct {
	TotalHands    int
	Profit        float64
	Mean          float64
	HouseEdge     float64 // True house edge: Profit / TotalWagered
	M2            float64
	DoubleWin     int
	DoubleLoss    int

	// Comprehensive metrics
	PlayerBlackjacks int
	DealerBlackjacks int
	PlayerBusts      int
	DealerBusts      int
	Wins             int
	Losses           int
	Pushes           int
	DoublesAttempted int
	SplitsAttempted  int

	// Dealer outcome frequencies
	Dealer17 int
	Dealer18 int
	Dealer19 int
	Dealer20 int
	Dealer21 int

	// Hand counts
	TotalHandsInRounds int

	// Pair opportunity tracking
	PairOpportunities      int
	PairOpportunitiesByRank map[int]int
	SplitsByRank           map[int]int
	TotalWagered           float64
	SplitHandsTotal        int
	SplitRounds            int
	SplitProfit            float64
	DoublesAfterSplit      int
}

func (s *Stats) Update(result float64) {
	s.TotalHands++
	s.Profit += result

	if int(result) == 2 {
		s.DoubleWin++
	}
	if int(result) == -2 {
		s.DoubleLoss++
	}

	delta := result - s.Mean
	s.Mean += delta / float64(s.TotalHands)
	s.M2 += delta * (result - s.Mean)
}

func (s *Stats) UpdateDetailed(result engine.HandResult) {
	s.TotalHands++
	s.Profit += result.Profit
	s.TotalHandsInRounds += result.HandsInRound
	s.TotalWagered += result.TotalWagered

	// Track pair opportunities (only count initial pair, not resplits)
	if result.PairOpportunity && result.HandsInRound == 1 {
		s.PairOpportunities++
		if s.PairOpportunitiesByRank == nil {
			s.PairOpportunitiesByRank = make(map[int]int)
		}
		s.PairOpportunitiesByRank[result.PairRank]++
	}

	// Track splits
	if result.WasSplit {
		s.SplitsAttempted++
		if s.SplitsByRank == nil {
			s.SplitsByRank = make(map[int]int)
		}
		if result.PairRank > 0 {
			s.SplitsByRank[result.PairRank]++
		}
		s.SplitRounds++
		s.SplitHandsTotal += result.HandsInRound
		s.SplitProfit += result.Profit
	}

	if result.PlayerBlackjack {
		s.PlayerBlackjacks++
	}
	if result.DealerBlackjack {
		s.DealerBlackjacks++
	}
	if result.PlayerBust {
		s.PlayerBusts++
	}
	if result.DealerBust {
		s.DealerBusts++
	}
	if result.WasDoubled {
		s.DoublesAttempted++
		// Track doubles after split
		if result.WasSplit {
			s.DoublesAfterSplit++
		}
	}

	// Track dealer outcomes
	switch result.DealerValue {
	case 17:
		s.Dealer17++
	case 18:
		s.Dealer18++
	case 19:
		s.Dealer19++
	case 20:
		s.Dealer20++
	case 21:
		if !result.DealerBlackjack {
			s.Dealer21++
		}
	}

	// Track win/loss/push
	if result.Profit > 0 {
		s.Wins++
	} else if result.Profit < 0 {
		s.Losses++
	} else {
		s.Pushes++
	}

	// Track double wins/losses
	if result.WasDoubled {
		if result.Profit > 0 {
			s.DoubleWin++
		} else if result.Profit < 0 {
			s.DoubleLoss++
		}
	}

	delta := result.Profit - s.Mean
	s.Mean += delta / float64(s.TotalHands)
	s.M2 += delta * (result.Profit - s.Mean)
	
	// Calculate true house edge: Profit / TotalWagered
	if s.TotalWagered > 0 {
		s.HouseEdge = s.Profit / s.TotalWagered
	}
}

func (s *Stats) Variance() float64 {
	if s.TotalHands < 2 {
		return 0
	}
	return s.M2 / float64(s.TotalHands-1)
}

func (s *Stats) StdDev() float64 {
	return math.Sqrt(s.Variance())
}

// PerHandVariance calculates variance using direct formula (requires storing profits)
// This is a diagnostic method - for large simulations, use Variance() instead
func (s *Stats) PerHandVariance(profits []float64) float64 {
	if len(profits) < 2 {
		return 0
	}
	
	sumSquares := 0.0
	mean := s.Mean
	
	for _, profit := range profits {
		diff := profit - mean
		sumSquares += diff * diff
	}
	
	return sumSquares / float64(len(profits)-1)
}

// ValidateVariance compares Welford's algorithm to expected range
func (s *Stats) ValidateVariance() (bool, string) {
	if s.TotalHands < 2 {
		return true, "Insufficient data"
	}
	
	stdDev := s.StdDev()
	
	expectedMin := 1.15
	expectedMax := 1.30
	
	if stdDev >= expectedMin && stdDev <= expectedMax {
		return true, fmt.Sprintf("Std Dev %.6f within expected range [%.2f, %.2f]", stdDev, expectedMin, expectedMax)
	}
	
	return false, fmt.Sprintf("Std Dev %.6f OUTSIDE expected range [%.2f, %.2f]", stdDev, expectedMin, expectedMax)
}

func (s *Stats) PrintMetrics() {
	fmt.Println("==== CardFlow21 Metrics ====")
	fmt.Printf("Total Rounds: %d\n", s.TotalHands)
	fmt.Printf("Total Hands: %d\n", s.TotalHandsInRounds)
	fmt.Printf("Profit: %.2f\n", s.Profit)
	fmt.Printf("Total Wagered: %.2f\n", s.TotalWagered)
	fmt.Printf("EV/Hand: %.6f\n", s.Mean)
	if s.TotalWagered > 0 {
		fmt.Printf("House Edge: %.6f (%.4f%%)\n", s.HouseEdge, s.HouseEdge*100.0)
	}
	fmt.Printf("Std Dev: %.6f\n", s.StdDev())
	varianceValid, varianceMsg := s.ValidateVariance()
	if varianceValid {
		fmt.Printf("Variance Status: ✓ %s\n", varianceMsg)
	} else {
		fmt.Printf("Variance Status: ✗ %s\n", varianceMsg)
	}
	fmt.Println("\n--- Outcomes ---")
	fmt.Printf("Wins: %d (%.2f%%)\n", s.Wins, 100.0*float64(s.Wins)/float64(s.TotalHands))
	fmt.Printf("Losses: %d (%.2f%%)\n", s.Losses, 100.0*float64(s.Losses)/float64(s.TotalHands))
	fmt.Printf("Pushes: %d (%.2f%%)\n", s.Pushes, 100.0*float64(s.Pushes)/float64(s.TotalHands))
	fmt.Println("\n--- Blackjacks ---")
	fmt.Printf("Player BJ: %d (%.2f%%)\n", s.PlayerBlackjacks, 100.0*float64(s.PlayerBlackjacks)/float64(s.TotalHands))
	fmt.Printf("Dealer BJ: %d (%.2f%%)\n", s.DealerBlackjacks, 100.0*float64(s.DealerBlackjacks)/float64(s.TotalHands))
	fmt.Println("\n--- Busts ---")
	fmt.Printf("Player Busts: %d (%.2f%%)\n", s.PlayerBusts, 100.0*float64(s.PlayerBusts)/float64(s.TotalHandsInRounds))
	fmt.Printf("Dealer Busts: %d (%.2f%%)\n", s.DealerBusts, 100.0*float64(s.DealerBusts)/float64(s.TotalHands))
	fmt.Println("\n--- Dealer Outcomes ---")
	if s.TotalHands > 0 {
		fmt.Printf("17: %d (%.2f%%)\n", s.Dealer17, 100.0*float64(s.Dealer17)/float64(s.TotalHands))
		fmt.Printf("18: %d (%.2f%%)\n", s.Dealer18, 100.0*float64(s.Dealer18)/float64(s.TotalHands))
		fmt.Printf("19: %d (%.2f%%)\n", s.Dealer19, 100.0*float64(s.Dealer19)/float64(s.TotalHands))
		fmt.Printf("20: %d (%.2f%%)\n", s.Dealer20, 100.0*float64(s.Dealer20)/float64(s.TotalHands))
		fmt.Printf("21: %d (%.2f%%)\n", s.Dealer21, 100.0*float64(s.Dealer21)/float64(s.TotalHands))
		fmt.Printf("Bust: %d (%.2f%%)\n", s.DealerBusts, 100.0*float64(s.DealerBusts)/float64(s.TotalHands))
	}
	fmt.Println("\n--- Advanced Actions ---")
	fmt.Printf("Doubles Attempted: %d\n", s.DoublesAttempted)
	fmt.Printf("Doubles Won: %d\n", s.DoubleWin)
	fmt.Printf("Doubles Lost: %d\n", s.DoubleLoss)
	fmt.Printf("Doubles After Split: %d\n", s.DoublesAfterSplit)
	fmt.Printf("Splits Attempted: %d\n", s.SplitsAttempted)
	fmt.Printf("Split Rounds: %d (%.2f%%)\n", s.SplitRounds, 100.0*float64(s.SplitRounds)/float64(s.TotalHands))
	fmt.Printf("Split Hands Total: %d\n", s.SplitHandsTotal)
	fmt.Printf("Split Profit: %.2f\n", s.SplitProfit)
	if s.SplitRounds > 0 {
		fmt.Printf("Split EV Contribution: %.6f\n", s.SplitProfit/float64(s.SplitRounds))
	}
	fmt.Println("\n--- Pair Opportunities ---")
	fmt.Printf("Pair Opportunities: %d (%.2f%%)\n", s.PairOpportunities, 100.0*float64(s.PairOpportunities)/float64(s.TotalHands))
	fmt.Printf("Splits Taken: %d (%.2f%% of pairs)\n", s.SplitsAttempted, 100.0*float64(s.SplitsAttempted)/float64(max(s.PairOpportunities, 1)))
	if s.PairOpportunitiesByRank != nil {
		fmt.Println("Pairs by Rank:")
		for rank := 1; rank <= 10; rank++ {
			if rank == 1 {
				fmt.Printf("  Aces: %d opportunities, %d splits\n", s.PairOpportunitiesByRank[1], s.SplitsByRank[1])
			} else if rank == 10 {
				fmt.Printf("  10s: %d opportunities, %d splits\n", s.PairOpportunitiesByRank[10], s.SplitsByRank[10])
			} else {
				fmt.Printf("  %ds: %d opportunities, %d splits\n", rank, s.PairOpportunitiesByRank[rank], s.SplitsByRank[rank])
			}
		}
	}
	fmt.Println("\n--- Betting ---")
	if s.TotalHands > 0 {
		fmt.Printf("Average Bet Per Hand: %.4f\n", s.TotalWagered/float64(s.TotalHandsInRounds))
		fmt.Printf("Total Wagered: %.2f\n", s.TotalWagered)
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
