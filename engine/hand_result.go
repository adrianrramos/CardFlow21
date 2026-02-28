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
	// Pair opportunity tracking
	PairOpportunity bool
	PairRank        int  // Rank of pair (1=Ace, 2-9, 10=10/J/Q/K, 0=no pair)
	DealerUpCard    Card // Dealer upcard when pair opportunity occurred
	TotalWagered    float64 // Total amount wagered in this round (accounts for doubles/splits)
}
