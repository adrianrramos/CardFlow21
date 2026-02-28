package simulation

import (
	"fmt"
	"math"
)

// TheoreticalBaseline contains expected values for DAS, H17, 6-deck blackjack
type TheoreticalBaseline struct {
	EV                    float64
	PlayerBlackjackFreq   float64
	DealerBlackjackFreq   float64
	DealerBustFreq        float64
	Dealer17Freq          float64
	Dealer18Freq          float64
	Dealer19Freq          float64
	Dealer20Freq          float64
	Dealer21Freq          float64
	WinFreq               float64
	LossFreq              float64
	PushFreq              float64
}

var BaselineDASH17 = TheoreticalBaseline{
	EV:                  -0.006,
	PlayerBlackjackFreq: 0.048,
	DealerBlackjackFreq: 0.048,
	DealerBustFreq:     0.28,
	Dealer17Freq:       0.14,
	Dealer18Freq:       0.14,
	Dealer19Freq:       0.13,
	Dealer20Freq:       0.18,
	Dealer21Freq:       0.07,
	WinFreq:            0.43,
	LossFreq:           0.49,
	PushFreq:           0.08,
}

type Deviation struct {
	Metric     string
	Expected   float64
	Actual     float64
	Difference float64
	PercentDiff float64
}

type Validator struct {
	Baseline TheoreticalBaseline
}

func NewValidator() *Validator {
	return &Validator{
		Baseline: BaselineDASH17,
	}
}

func (v *Validator) CompareMetrics(stats *Stats) []Deviation {
	deviations := []Deviation{}

	if stats.TotalHands == 0 {
		return deviations
	}

	// EV comparison - use HouseEdge (per-wager) for true comparison
	evActual := stats.Mean
	if stats.TotalWagered > 0 {
		evActual = stats.HouseEdge
	}
	evDiff := evActual - v.Baseline.EV
	deviations = append(deviations, Deviation{
		Metric:     "House Edge",
		Expected:   v.Baseline.EV,
		Actual:     evActual,
		Difference: evDiff,
		PercentDiff: (evDiff / math.Abs(v.Baseline.EV)) * 100,
	})

	// Player blackjack frequency
	pbjFreq := float64(stats.PlayerBlackjacks) / float64(stats.TotalHands)
	pbjDiff := pbjFreq - v.Baseline.PlayerBlackjackFreq
	deviations = append(deviations, Deviation{
		Metric:     "Player Blackjack Frequency",
		Expected:   v.Baseline.PlayerBlackjackFreq,
		Actual:     pbjFreq,
		Difference: pbjDiff,
		PercentDiff: (pbjDiff / v.Baseline.PlayerBlackjackFreq) * 100,
	})

	// Dealer blackjack frequency
	dbjFreq := float64(stats.DealerBlackjacks) / float64(stats.TotalHands)
	dbjDiff := dbjFreq - v.Baseline.DealerBlackjackFreq
	deviations = append(deviations, Deviation{
		Metric:     "Dealer Blackjack Frequency",
		Expected:   v.Baseline.DealerBlackjackFreq,
		Actual:     dbjFreq,
		Difference: dbjDiff,
		PercentDiff: (dbjDiff / v.Baseline.DealerBlackjackFreq) * 100,
	})

	// Dealer bust frequency
	dbustFreq := float64(stats.DealerBusts) / float64(stats.TotalHands)
	dbustDiff := dbustFreq - v.Baseline.DealerBustFreq
	deviations = append(deviations, Deviation{
		Metric:     "Dealer Bust Frequency",
		Expected:   v.Baseline.DealerBustFreq,
		Actual:     dbustFreq,
		Difference: dbustDiff,
		PercentDiff: (dbustDiff / v.Baseline.DealerBustFreq) * 100,
	})

	// Dealer outcome frequencies
	if stats.TotalHands > 0 {
		d17Freq := float64(stats.Dealer17) / float64(stats.TotalHands)
		deviations = append(deviations, Deviation{
			Metric:     "Dealer 17 Frequency",
			Expected:   v.Baseline.Dealer17Freq,
			Actual:     d17Freq,
			Difference: d17Freq - v.Baseline.Dealer17Freq,
			PercentDiff: ((d17Freq - v.Baseline.Dealer17Freq) / v.Baseline.Dealer17Freq) * 100,
		})

		d18Freq := float64(stats.Dealer18) / float64(stats.TotalHands)
		deviations = append(deviations, Deviation{
			Metric:     "Dealer 18 Frequency",
			Expected:   v.Baseline.Dealer18Freq,
			Actual:     d18Freq,
			Difference: d18Freq - v.Baseline.Dealer18Freq,
			PercentDiff: ((d18Freq - v.Baseline.Dealer18Freq) / v.Baseline.Dealer18Freq) * 100,
		})

		d19Freq := float64(stats.Dealer19) / float64(stats.TotalHands)
		deviations = append(deviations, Deviation{
			Metric:     "Dealer 19 Frequency",
			Expected:   v.Baseline.Dealer19Freq,
			Actual:     d19Freq,
			Difference: d19Freq - v.Baseline.Dealer19Freq,
			PercentDiff: ((d19Freq - v.Baseline.Dealer19Freq) / v.Baseline.Dealer19Freq) * 100,
		})

		d20Freq := float64(stats.Dealer20) / float64(stats.TotalHands)
		deviations = append(deviations, Deviation{
			Metric:     "Dealer 20 Frequency",
			Expected:   v.Baseline.Dealer20Freq,
			Actual:     d20Freq,
			Difference: d20Freq - v.Baseline.Dealer20Freq,
			PercentDiff: ((d20Freq - v.Baseline.Dealer20Freq) / v.Baseline.Dealer20Freq) * 100,
		})

		d21Freq := float64(stats.Dealer21) / float64(stats.TotalHands)
		deviations = append(deviations, Deviation{
			Metric:     "Dealer 21 Frequency",
			Expected:   v.Baseline.Dealer21Freq,
			Actual:     d21Freq,
			Difference: d21Freq - v.Baseline.Dealer21Freq,
			PercentDiff: ((d21Freq - v.Baseline.Dealer21Freq) / v.Baseline.Dealer21Freq) * 100,
		})
	}

	// Win/Loss/Push frequencies
	winFreq := float64(stats.Wins) / float64(stats.TotalHands)
	deviations = append(deviations, Deviation{
		Metric:     "Win Frequency",
		Expected:   v.Baseline.WinFreq,
		Actual:     winFreq,
		Difference: winFreq - v.Baseline.WinFreq,
		PercentDiff: ((winFreq - v.Baseline.WinFreq) / v.Baseline.WinFreq) * 100,
	})

	lossFreq := float64(stats.Losses) / float64(stats.TotalHands)
	deviations = append(deviations, Deviation{
		Metric:     "Loss Frequency",
		Expected:   v.Baseline.LossFreq,
		Actual:     lossFreq,
		Difference: lossFreq - v.Baseline.LossFreq,
		PercentDiff: ((lossFreq - v.Baseline.LossFreq) / v.Baseline.LossFreq) * 100,
	})

	pushFreq := float64(stats.Pushes) / float64(stats.TotalHands)
	deviations = append(deviations, Deviation{
		Metric:     "Push Frequency",
		Expected:   v.Baseline.PushFreq,
		Actual:     pushFreq,
		Difference: pushFreq - v.Baseline.PushFreq,
		PercentDiff: ((pushFreq - v.Baseline.PushFreq) / v.Baseline.PushFreq) * 100,
	})

	return deviations
}

func (v *Validator) FindLargestDeviation(stats *Stats) Deviation {
	deviations := v.CompareMetrics(stats)
	if len(deviations) == 0 {
		return Deviation{}
	}

	largest := deviations[0]
	maxAbsDiff := math.Abs(deviations[0].Difference)

	for _, d := range deviations[1:] {
		absDiff := math.Abs(d.Difference)
		if absDiff > maxAbsDiff {
			maxAbsDiff = absDiff
			largest = d
		}
	}

	return largest
}

func (v *Validator) PrintDeviations(stats *Stats) {
	deviations := v.CompareMetrics(stats)
	
	if len(deviations) == 0 {
		return
	}

	largest := v.FindLargestDeviation(stats)

	fmt.Println("\n--- Deviation Analysis ---")
	for _, d := range deviations {
		marker := ""
		if d.Metric == largest.Metric {
			marker = " <-- LARGEST"
		}
		if math.Abs(d.PercentDiff) > 5.0 {
			marker += " *** >5% DEVIATION"
		}
		fmt.Printf("%s: Expected %.6f, Actual %.6f, Diff %.6f (%.2f%%)%s\n",
			d.Metric, d.Expected, d.Actual, d.Difference, d.PercentDiff, marker)
	}
}

// ValidateBlackjackPayouts validates blackjack payout logic
func (v *Validator) ValidateBlackjackPayouts(traces []HandTracer) (bool, []HandTracer) {
	mismatches := []HandTracer{}
	
	for _, trace := range traces {
		if !trace.Match {
			mismatches = append(mismatches, trace)
		}
	}
	
	return len(mismatches) == 0, mismatches
}
