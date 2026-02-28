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
	}
	if result.WasSplit {
		s.SplitsAttempted++
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

func (s *Stats) PrintMetrics() {
	fmt.Println("==== CardFlow21 Metrics ====")
	fmt.Printf("Total Rounds: %d\n", s.TotalHands)
	fmt.Printf("Total Hands: %d\n", s.TotalHandsInRounds)
	fmt.Printf("Profit: %.2f\n", s.Profit)
	fmt.Printf("EV/Hand: %.6f\n", s.Mean)
	fmt.Printf("Std Dev: %.6f\n", s.StdDev())
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
	fmt.Printf("Splits Attempted: %d\n", s.SplitsAttempted)
}
