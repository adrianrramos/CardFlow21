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

/*
The Welford Online Algorithm as described in
https://en.wikipedia.org/wiki/Algorithms_for_calculating_variance
*/
type WelfordVariance struct {
	Mean  float64
	Count int64
	M2    float64
}

/*
Performs an update to overall stats from the payout
of a new hand
*/
func (w *WelfordVariance) AddVariable(payout float64) {
	w.Count++
	oldMean := w.Mean
	w.Mean += (payout - w.Mean) / float64(w.Count)
	w.M2 += (payout - oldMean) * (payout - w.Mean)
}

/*
Reverses an update to overall stats with a given payout
*/
func (w *WelfordVariance) RemoveVariable(payout float64) {
	w.Count--
	newMean := w.Mean
	w.Mean -= (payout - w.Mean) / float64(w.Count)
	w.M2 -= (payout - newMean) * (payout - w.Mean)
}

func (w *WelfordVariance) GetMean() float64 {
	return w.Mean
}

func (w *WelfordVariance) GetVariance() float64 {
	return w.M2 / float64(w.Count)
}

func (w *WelfordVariance) GetSampleVariance() float64 {
	return w.M2 / float64(w.Count-1)
}

func (w *WelfordVariance) GetStdDev() float64 {
	return math.Sqrt(w.GetVariance())
}
func (w *WelfordVariance) GetSampleStdDev() float64 {
	return math.Sqrt(w.GetSampleVariance())
}
