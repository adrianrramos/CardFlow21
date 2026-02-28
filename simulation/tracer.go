package simulation

import (
	"cardflow21/engine"
	"fmt"
)

// HandTracer logs individual hand details for validation
type HandTracer struct {
	HandNumber      int
	PlayerCards     []engine.Card
	DealerCards     []engine.Card
	PlayerBlackjack bool
	DealerBlackjack bool
	ExpectedPayout  float64
	ActualPayout    float64
	Match           bool
	Notes           string
}

// TraceBlackjacks traces hands that involve blackjacks
func TraceBlackjacks(results []engine.HandResult, limit int) []HandTracer {
	traces := []HandTracer{}
	count := 0
	
	for i, result := range results {
		if count >= limit {
			break
		}
		
		// Only trace hands with blackjacks
		if result.PlayerBlackjack || result.DealerBlackjack {
			trace := HandTracer{
				HandNumber:      i + 1,
				PlayerBlackjack: result.PlayerBlackjack,
				DealerBlackjack: result.DealerBlackjack,
				ActualPayout:    result.Profit,
			}
			
			// Calculate expected payout
			if result.PlayerBlackjack && result.DealerBlackjack {
				trace.ExpectedPayout = 0.0
				trace.Notes = "Both blackjack - push"
			} else if result.PlayerBlackjack && !result.DealerBlackjack {
				// Player BJ pays 1.5x bet (assuming 1 unit bet for initial hand)
				trace.ExpectedPayout = 1.5
				trace.Notes = "Player blackjack - pays 3:2"
			} else if result.DealerBlackjack {
				trace.ExpectedPayout = -1.0
				trace.Notes = "Dealer blackjack - player loses"
			}
			
			trace.Match = (trace.ExpectedPayout == trace.ActualPayout)
			traces = append(traces, trace)
			count++
		}
	}
	
	return traces
}

// PrintBlackjackTraces prints traced blackjack hands
func PrintBlackjackTraces(traces []HandTracer) {
	fmt.Println("\n=== Blackjack Payout Trace ===")
	fmt.Printf("Traced %d blackjack hands\n\n", len(traces))
	
	mismatches := 0
	for _, trace := range traces {
		status := "✓"
		if !trace.Match {
			status = "✗ MISMATCH"
			mismatches++
		}
		
		fmt.Printf("Hand #%d: %s\n", trace.HandNumber, status)
		fmt.Printf("  Player BJ: %v, Dealer BJ: %v\n", trace.PlayerBlackjack, trace.DealerBlackjack)
		fmt.Printf("  Expected: %.2f, Actual: %.2f\n", trace.ExpectedPayout, trace.ActualPayout)
		fmt.Printf("  %s\n", trace.Notes)
		if !trace.Match {
			fmt.Printf("  *** DEVIATION: %.2f\n", trace.ActualPayout-trace.ExpectedPayout)
		}
		fmt.Println()
	}
	
	if mismatches > 0 {
		fmt.Printf("*** FOUND %d PAYOUT MISMATCHES ***\n", mismatches)
	} else {
		fmt.Println("All blackjack payouts correct!")
	}
}
