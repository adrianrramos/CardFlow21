package engine

// HandResult contains detailed information about a single hand outcome
type HandResult struct {
	PlayerBlackjack bool
	DealerBlackjack bool
	PlayerBust      bool
	DealerBust      bool
	PlayerValue     int
	DealerValue     int
	Profit         float64
	WasDoubled     bool
	WasSplit       bool
	HandsInRound   int
}
