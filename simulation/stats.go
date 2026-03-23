package simulation

import "math"

type Stats struct {
	TotalHands   int
	TotalWagered float64
	Profit       float64
	Mean         float64
	MeanHands    float64
	M2           float64
	DoubleWin    int
	DoubleLoss   int
}

func (s *Stats) Update(result float64) {
	s.TotalHands++
	s.Profit += result
	// TODO:
	// This is a real bug in the EV/Wagered calculation. In simulation/stats.go, you accumulate
	// TotalWagered using math.Abs(result) (abs profit), not the amount actually wagered:
	s.TotalWagered += math.Abs(result)

	delta := result - s.Mean
	deltaHands := result - s.MeanHands
	s.Mean += delta / s.TotalWagered
	s.MeanHands += deltaHands / float64(s.TotalHands)
	s.M2 += delta * (result - s.Mean)
}

func (s *Stats) Variance() float64 {
	if s.TotalHands < 2 {
		return 0
	}
	return s.M2 / float64(s.TotalHands-1)
}

// TODO: this is not being used
func (s *Stats) StdDev() float64 {
	return math.Sqrt(s.Variance())
}
