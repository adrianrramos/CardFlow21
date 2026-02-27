package simulation

import "math"

type Stats struct {
	TotalHands    int
	Profit        float64
	Mean          float64
	M2            float64
	DoubleWin     int
	DoubleLoss    int
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

func (s *Stats) Variance() float64 {
	if s.TotalHands < 2 {
		return 0
	}
	return s.M2 / float64(s.TotalHands-1)
}

func (s *Stats) StdDev() float64 {
	return math.Sqrt(s.Variance())
}
