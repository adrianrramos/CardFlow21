package engine

import "testing"

// --- Test Helpers ---

// c creates a Card with the given rank (1=Ace, 2-10, 11=Jack, 12=Queen, 13=King).
func c(rank int) Card {
	return Card{Suit: 0, Rank: Rank(rank)}
}

// newTestShoe creates a Shoe that deals cards in the given argument order.
// Shoe.Draw() pops from the end of the slice, so the first argument is written
// to the last index (drawn first), the second to the previous index, etc.
func newTestShoe(cards ...Card) *Shoe {
	// 1. Start from a real shoe
	s := NewShoe(6, 0.85) // or whatever decks/penetration you want for tests

	// 2. Ensure the shoe has room (it will)
	// 3. Overwrite the last len(cards) positions so the first argument is at
	//    len-1 (first Draw), the second at len-2, etc.
	for i, card := range cards {
		s.cards[len(s.cards)-1-i] = card
	}

	// Optionally reset count/true_count if you need them specific
	s.count = 0
	s.true_count = 0

	return s
}

// mockStrategy returns actions from a predetermined slice, defaulting to
// Stand once the slice is exhausted.
type mockStrategy struct {
	actions          []Action
	idx              int
	surrenderInitial bool // when true, surrender is taken before the round plays out
}

func (m *mockStrategy) Decide(_ Hand, _ Card) Action {
	if m.idx >= len(m.actions) {
		return Stand
	}
	a := m.actions[m.idx]
	m.idx++
	return a
}

func (m *mockStrategy) Name() string { return "mock" }

func (m *mockStrategy) CheckSurrenderChart(_ Hand, _ Card) bool {
	return m.surrenderInitial
}

// --- PlayHandWithStrategy Tests ---
//
// Deal order from newTestShoe(p1, d1, p2, d2, ...):
//   shoe.Draw() #1 → p1 (player card 1)
//   shoe.Draw() #2 → d1 (dealer card 1)
//   shoe.Draw() #3 → p2 (player card 2)
//   shoe.Draw() #4 → d2 (dealer card 2)
//   subsequent draws → extra cards for hits / dealer draws
//
// Action slices must account for TWO Decide calls per normal hand: one in
// the outer split-detection loop, then one (or more) inside PlayOutHand.

// TestPlayHandWithStrategy_PlayerBlackjack verifies the 3:2 payout when
// only the player has blackjack.
func TestPlayHandWithStrategy_PlayerBlackjack(t *testing.T) {
	// Player: Ace(1) + King(13) = 21 BJ; Dealer: 2 + 3 = 5 (no BJ)
	shoe := newTestShoe(c(1), c(2), c(13), c(3))
	got := PlayRound(shoe, &mockStrategy{}, &StatsTracker{}, false)
	if got != 1.5 {
		t.Errorf("player BJ, no dealer BJ: want 1.5, got %v", got)
	}
}

// TestPlayHandWithStrategy_BothBlackjack verifies current behavior when both
// player and dealer have blackjack.
func TestPlayHandWithStrategy_BothBlackjack(t *testing.T) {
	// Player: Ace + King = 21 BJ; Dealer: King + Ace = 21 BJ
	shoe := newTestShoe(c(1), c(13), c(13), c(1))
	got := PlayRound(shoe, &mockStrategy{}, &StatsTracker{}, false)
	if got != 0 {
		t.Errorf("both BJ: want 0 (dealer blackjack precedence), got %v", got)
	}
}

// TestPlayHandWithStrategy_DealerBlackjack verifies the player loses when
// the dealer has blackjack and the player does not, before any further play.
func TestPlayHandWithStrategy_DealerBlackjack(t *testing.T) {
	// Player: 5 + 6 = 11 (no BJ); Dealer: King + Ace = 21 BJ — only four cards dealt.
	shoe := newTestShoe(c(5), c(13), c(6), c(1), c(10))
	strat := &mockStrategy{actions: []Action{Stand, Hit, Stand}}
	got := PlayRound(shoe, strat, &StatsTracker{}, false)
	if got != -1.0 {
		t.Errorf("dealer BJ: want -1.0, got %v", got)
	}
}

// TestPlayHandWithStrategy_PlayerWins verifies a standard win.
func TestPlayHandWithStrategy_PlayerWins(t *testing.T) {
	// Player: King + Queen = 20; Dealer: 8 + 9 = 17
	shoe := newTestShoe(c(13), c(8), c(12), c(9))
	strat := &mockStrategy{actions: []Action{Stand, Stand}}
	got := PlayRound(shoe, strat, &StatsTracker{}, false)
	if got != 1.0 {
		t.Errorf("player 20 vs dealer 17: want 1.0, got %v", got)
	}
}

// TestPlayHandWithStrategy_PlayerLoses verifies a standard loss.
func TestPlayHandWithStrategy_PlayerLoses(t *testing.T) {
	// Player: 6 + 10 = 16; Dealer: 8 + 10 = 18
	shoe := newTestShoe(c(6), c(8), c(10), c(10))
	strat := &mockStrategy{actions: []Action{Stand, Stand}}
	got := PlayRound(shoe, strat, &StatsTracker{}, false)
	if got != -1.0 {
		t.Errorf("player 16 vs dealer 18: want -1.0, got %v", got)
	}
}

// TestPlayHandWithStrategy_Push verifies that equal hand values produce a push.
func TestPlayHandWithStrategy_Push(t *testing.T) {
	// Player: 8 + 10 = 18; Dealer: 8 + 10 = 18
	shoe := newTestShoe(c(8), c(8), c(10), c(10))
	strat := &mockStrategy{actions: []Action{Stand, Stand}}
	got := PlayRound(shoe, strat, &StatsTracker{}, false)
	if got != 0.0 {
		t.Errorf("player 18 vs dealer 18: want 0.0 (push), got %v", got)
	}
}

// TestPlayHandWithStrategy_PlayerBusts verifies that a busted player loses.
//
// Current implementation compares raw totals after bust because IsBust is not
// persisted on the original Hand.
func TestPlayHandWithStrategy_PlayerBusts(t *testing.T) {
	// Player: 7 + 9 = 16, hits 10 → 26 (bust); Dealer: 8 + 10 = 18
	shoe := newTestShoe(c(7), c(8), c(9), c(10), c(10))
	strat := &mockStrategy{actions: []Action{Stand, Hit}}
	got := PlayRound(shoe, strat, &StatsTracker{}, false)
	if got != -1.0 {
		t.Errorf("player busts (current behavior): want -1.0, got %v", got)
	}
}

// TestPlayHandWithStrategy_DealerBusts verifies the player wins when the
// dealer busts.
//
// Current implementation compares raw totals after dealer bust because IsBust
// is not persisted on the original Hand.
func TestPlayHandWithStrategy_DealerBusts(t *testing.T) {
	// Player: 10 + 8 = 18 (stands); Dealer: 6 + 10 = 16, draws 10 → 26 (bust)
	shoe := newTestShoe(c(10), c(6), c(8), c(10), c(10))
	strat := &mockStrategy{actions: []Action{Stand, Stand}}
	got := PlayRound(shoe, strat, &StatsTracker{}, false)
	if got != 1.0 {
		t.Errorf("dealer busts (current behavior): want 1.0, got %v", got)
	}
}

// TestPlayHandWithStrategy_DoubleWin verifies a doubled wager is paid out
// in full on a win.
func TestPlayHandWithStrategy_DoubleWin(t *testing.T) {
	// Player: 10 + 8 = 18 (doubles, gets 5th card); Dealer: 7 + 10 = 17 → +2
	// 5th card needed for the double; use 2 so player ends 18+2=20 (still wins)
	shoe := newTestShoe(c(8), c(6), c(3), c(10), c(12))

	var expected_win float64
	if shoe.true_count > 0 {
		expected_win = float64(shoe.true_count) * 2 * 2
	} else {
		expected_win = 2.0
	}
	strat := &mockStrategy{actions: []Action{Stand, Double}}
	got := PlayRound(shoe, strat, &StatsTracker{}, false)
	if got != expected_win {
		t.Errorf("double win: want %v, got %v", expected_win, got)
	}
}

// TestPlayHandWithStrategy_DoubleLoss verifies a doubled wager is lost in
// full on a loss.
func TestPlayHandWithStrategy_DoubleLoss(t *testing.T) {
	// Player: 5 + 7 = 12 (doubles, gets 5th card); Dealer: 8 + 10 = 18 → -2
	// 5th card needed for the double; use 5 so player ends 12+5=17 (still loses)
	shoe := newTestShoe(c(5), c(8), c(7), c(10), c(5))
	strat := &mockStrategy{actions: []Action{Stand, Double}}
	got := PlayRound(shoe, strat, &StatsTracker{}, false)
	if got != -2.0 {
		t.Errorf("double loss: want -2.0, got %v", got)
	}
}

// TestPlayHandWithStrategy_Split_BothWin verifies the combined profit when
// both split hands win.
//
// Shoe layout (deal order → split draw order → dealer draw):
//
//	p1=9, d1=6, p2=9, d2=Jack → player pair of 9s; dealer 6+10=16
//	5th card (hand 0 extra): 10 → [9,10]=19
//	6th card (hand 1 extra): 10 → [9,10]=19
//	7th card (dealer draw):  2  → dealer 16+2=18
//	Both 19 > 18 → +2
func TestPlayHandWithStrategy_Split_BothWin(t *testing.T) {
	shoe := newTestShoe(c(9), c(6), c(9), c(11), c(10), c(10), c(2))
	// Split (outer), Stand+Stand for hand 0, Stand+Stand for hand 1
	strat := &mockStrategy{actions: []Action{Split, Stand, Stand, Stand, Stand}}
	got := PlayRound(shoe, strat, &StatsTracker{}, false)
	if got != 2.0 {
		t.Errorf("split both win: want 2.0, got %v", got)
	}
}

// TestPlayHandWithStrategy_Split_BothLose verifies the combined loss when
// both split hands lose.
//
// Shoe layout:
//
//	p1=7, d1=8, p2=7, d2=10 → player pair of 7s; dealer 8+10=18 (stands)
//	5th card (hand 0 extra): 5 → [7,5]=12
//	6th card (hand 1 extra): 5 → [7,5]=12
//	Both 12 < 18 → -2
func TestPlayHandWithStrategy_Split_BothLose(t *testing.T) {
	shoe := newTestShoe(c(7), c(8), c(7), c(10), c(5), c(5))
	strat := &mockStrategy{actions: []Action{Split, Stand, Stand, Stand, Stand}}
	got := PlayRound(shoe, strat, &StatsTracker{}, false)
	if got != -2.0 {
		t.Errorf("split both lose: want -2.0, got %v", got)
	}
}

// ===== Game rules & payouts (not basic-strategy chart tests) =====

// --- Hand / ace scoring ---

func TestHand_AceHandling(t *testing.T) {
	tests := []struct {
		name string
		r    []int
		want int
	}{
		{"A+9=20", []int{1, 9}, 20},
		{"A+9+5=15", []int{1, 9, 5}, 15},
		{"A+A=12", []int{1, 1}, 12},
		{"A+A+9=21", []int{1, 1, 9}, 21},
		{"A+A+A+9=12", []int{1, 1, 1, 9}, 12},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var h Hand
			for _, rank := range tt.r {
				h.AddCard(c(rank))
			}
			if g := h.Value(); g != tt.want {
				t.Errorf("Value(): want %d, got %d", tt.want, g)
			}
		})
	}
}

func TestHand_BlackjackOnlyTwoCards(t *testing.T) {
	var bj Hand
	bj.AddCard(c(1))
	bj.AddCard(c(13))
	if !bj.IsBlackjack() {
		t.Fatal("two-card 21 should be blackjack")
	}
	bj.AddCard(c(9))
	if bj.IsBlackjack() {
		t.Error("three-card 21 must not count as blackjack")
	}
}

// --- Initial deal & dealer blackjack peek ---

func TestPlayRound_InitialDealConsumesFourCardsOnEarlyExit(t *testing.T) {
	shoe := newTestShoe(c(1), c(2), c(13), c(3))
	before := len(shoe.cards)
	PlayRound(shoe, &mockStrategy{}, &StatsTracker{}, false)
	if got := before - len(shoe.cards); got != 4 {
		t.Errorf("player BJ early exit: want 4 cards dealt, got %d", got)
	}
}

func TestPlayRound_DealerBlackjackNoExtraPlayerDraws(t *testing.T) {
	shoe := newTestShoe(c(5), c(13), c(6), c(1))
	before := len(shoe.cards)
	PlayRound(shoe, &mockStrategy{actions: []Action{Hit}}, &StatsTracker{}, false)
	if got := before - len(shoe.cards); got != 4 {
		t.Errorf("dealer BJ: want exactly 4 cards (no further play), got %d", got)
	}
}

// --- Dealer H17 ---

func TestPlayRound_DealerHitsSoft17(t *testing.T) {
	shoe := newTestShoe(c(10), c(1), c(10), c(6), c(2))
	strat := &mockStrategy{actions: []Action{Stand, Stand}}
	got := PlayRound(shoe, strat, &StatsTracker{}, false)
	if got != 1.0 {
		t.Errorf("dealer hits soft 17: want +1.0, got %v", got)
	}
}

func TestPlayRound_DealerStandsHard17(t *testing.T) {
	shoe := newTestShoe(c(13), c(10), c(12), c(7))
	strat := &mockStrategy{actions: []Action{Stand, Stand}}
	got := PlayRound(shoe, strat, &StatsTracker{}, false)
	if got != 1.0 {
		t.Errorf("player 20 vs dealer hard 17: want +1.0, got %v", got)
	}
}

// Current engine: dealer still draws after player busts (no peek/skip).
func TestPlayRound_DealerStillDrawsAfterPlayerBusts(t *testing.T) {
	shoe := newTestShoe(c(10), c(10), c(6), c(6), c(10), c(5))
	before := len(shoe.cards)
	strat := &mockStrategy{actions: []Action{Stand, Hit}}
	PlayRound(shoe, strat, &StatsTracker{}, false)
	if dealt := before - len(shoe.cards); dealt != 6 {
		t.Errorf("player bust then dealer hits 16: want 6 cards dealt, got %d", dealt)
	}
}

// --- Hit / stand ---

func TestPlayRound_HitAddsOneCardPerHit(t *testing.T) {
	shoe := newTestShoe(c(5), c(10), c(5), c(10), c(3), c(2))
	before := len(shoe.cards)
	strat := &mockStrategy{actions: []Action{Stand, Hit, Hit, Stand}}
	PlayRound(shoe, strat, &StatsTracker{}, false)
	if dealt := before - len(shoe.cards); dealt != 6 {
		t.Errorf("initial 4 + two hits: want 6 cards dealt, got %d", dealt)
	}
}

// --- Double (incl. push) & split + double ---

func TestPlayRound_DoublePush(t *testing.T) {
	shoe := newTestShoe(c(7), c(9), c(4), c(9), c(7))
	strat := &mockStrategy{actions: []Action{Stand, Double}}
	got := PlayRound(shoe, strat, &StatsTracker{}, false)
	if got != 0.0 {
		t.Errorf("double push: want 0, got %v", got)
	}
}

func TestPlayRound_DoubleAfterSplit_DAS(t *testing.T) {
	// 6,6 vs 5–10 (15); split; [6,5]=11 doubles to 19; [6,10]=16 stands; dealer 15+10 busts
	shoe := newTestShoe(c(6), c(5), c(6), c(10), c(5), c(10), c(8), c(10))
	strat := &mockStrategy{actions: []Action{Split, Stand, Double, Stand, Stand}}
	st := &StatsTracker{}
	got := PlayRound(shoe, strat, st, false)
	if got != 3.0 {
		t.Errorf("split+DAS: want +3.0 (2+1), got %v", got)
	}
	if st.SplitHands != 1 || st.DoubleWin != 1 {
		t.Errorf("stats: SplitHands=%d DoubleWin=%d", st.SplitHands, st.DoubleWin)
	}
}

// --- Split: basics, net zero, push+loss, two 21s pay 1:1 each ---

func TestPlayRound_Split88_OneWinOneLose_NetZero(t *testing.T) {
	// Tail order must match draw order: 5th=split fill hand0, 6th=hit, 7th=split fill hand1
	// (see newTestShoe: args map to consecutive draws from the shoe end).
	shoe := newTestShoe(c(8), c(9), c(8), c(10), c(2), c(10), c(6))
	strat := &mockStrategy{actions: []Action{Split, Stand, Hit, Stand, Stand, Stand}}
	got := PlayRound(shoe, strat, &StatsTracker{}, false)
	if got != 0.0 {
		t.Errorf("split 8,8 one win one lose: want 0, got %v", got)
	}
}

func TestPlayRound_Split_OnePushOneLose(t *testing.T) {
	shoe := newTestShoe(c(9), c(8), c(9), c(10), c(9), c(8), c(10))
	strat := &mockStrategy{actions: []Action{Split, Stand, Stand, Stand, Stand}}
	got := PlayRound(shoe, strat, &StatsTracker{}, false)
	if got != -1.0 {
		t.Errorf("split one push one lose: want -1.0, got %v", got)
	}
}

func TestPlayRound_SplitAcesWithTenEachPaysEvenMoneyNotThreeToTwo(t *testing.T) {
	// Each split hand is [A,K]=21 with two cards (IsBlackjack true on the hand), but
	// settlement still pays 1:1 per hand vs a non-BJ dealer total.
	shoe := newTestShoe(c(1), c(9), c(1), c(10), c(13), c(13))
	strat := &mockStrategy{actions: []Action{Split, Stand, Stand}}
	got := PlayRound(shoe, strat, &StatsTracker{}, false)
	if got != 2.0 {
		t.Errorf("split A,A to two 21s vs 19: want +2.0 (1:1 each), got %v", got)
	}
}

func TestPlayRound_SplitAces_OneCardEachNoFurtherPlay(t *testing.T) {
	shoe := newTestShoe(c(1), c(6), c(1), c(5), c(4), c(3), c(10))
	before := len(shoe.cards)
	strat := &mockStrategy{actions: []Action{Split, Stand, Stand}}
	st := &StatsTracker{}
	got := PlayRound(shoe, strat, st, false)
	if st.SplitHands != 1 {
		t.Errorf("SplitHands: want 1, got %d", st.SplitHands)
	}
	if got != -2.0 {
		t.Errorf("split aces vs dealer 21: want -2.0, got %v", got)
	}
	if dealt := before - len(shoe.cards); dealt != 7 {
		t.Errorf("split aces: want 7 cards dealt (no extra hit cards), got %d", dealt)
	}
}

func TestPlayRound_Resplit_UpToFourHands(t *testing.T) {
	// Four 8s → split three times → four hands of 8; dealer 7–10 (17); all stand on 8
	shoe := newTestShoe(c(8), c(7), c(8), c(10), c(8), c(8), c(8), c(8), c(8), c(8), c(8))
	strat := &mockStrategy{actions: []Action{
		Split, Split, Split,
		Stand, Stand, Stand, Stand, Stand, Stand, Stand, Stand, Stand, Stand,
	}}
	st := &StatsTracker{}
	got := PlayRound(shoe, strat, st, false)
	if st.SplitHands != 3 {
		t.Errorf("SplitHands: want 3 resplits, got %d", st.SplitHands)
	}
	if got != -4.0 {
		t.Errorf("four 8s vs 17: want -4.0, got %v", got)
	}
}

// --- Surrender & insurance ---

func TestPlayRound_SurrenderLosesHalf(t *testing.T) {
	shoe := newTestShoe(c(10), c(10), c(6), c(9))
	st := &StatsTracker{}
	got := PlayRound(shoe, &mockStrategy{surrenderInitial: true}, st, false)
	if got != -0.5 {
		t.Errorf("surrender: want -0.5, got %v", got)
	}
	if st.Surrendered != 1 || st.TotalHands != 1 {
		t.Errorf("stats Surrendered=%d TotalHands=%d", st.Surrendered, st.TotalHands)
	}
}

func TestPlayRound_SurrenderRunsBeforeDealerBlackjackCheck(t *testing.T) {
	shoe := newTestShoe(c(10), c(13), c(6), c(1))
	got := PlayRound(shoe, &mockStrategy{surrenderInitial: true}, &StatsTracker{}, false)
	if got != -0.5 {
		t.Errorf("surrender vs dealer BJ hand: engine applies surrender first, want -0.5, got %v", got)
	}
}

func TestPlayRound_InsuranceWhenTrueCountDealerBlackjack(t *testing.T) {
	shoe := newTestShoe(c(10), c(1), c(10), c(13))
	// PlayRound updates Hi-Lo on each Draw(); keep running count high so true_count
	// stays >= 3 after the initial four cards (see Shoe.updateHiLoCount).
	shoe.count = 19
	shoe.true_count = 3
	st := &StatsTracker{}
	got := PlayRound(shoe, &mockStrategy{actions: []Action{Stand, Stand}}, st, true)
	if got != 0.0 {
		t.Errorf("insurance + dealer BJ (wager from count): want 0, got %v", got)
	}
	if st.TookInsurance != 1 {
		t.Errorf("TookInsurance: want 1, got %d", st.TookInsurance)
	}
}

// --- Deck / shoe ---

func TestShoe_Draw_FromEmptyReshufflesWithoutPanic(t *testing.T) {
	s := NewShoe(1, 0.5)
	for len(s.cards) > 0 {
		_ = s.Draw()
	}
	_ = s.Draw()
}

// --- Invalid / unsupported enforcement (document current behavior) ---

func TestPlayRound_InvalidDoubleAfterHitNotRejectedByEngine(t *testing.T) {
	// Engine applies Double whenever strategy requests it; invalid doubles are not blocked here.
	shoe := newTestShoe(c(5), c(10), c(5), c(10), c(2), c(10))
	strat := &mockStrategy{actions: []Action{Stand, Hit, Double}}
	got := PlayRound(shoe, strat, &StatsTracker{}, false)
	if got != -2.0 {
		t.Errorf("double after hit (3 cards): want -2.0 loss, got %v (documents current behavior)", got)
	}
}

// --- Betting stats ---

func TestPlayRound_StatsTotalWagered_SplitAddsSecondBaseWager(t *testing.T) {
	shoe := newTestShoe(c(9), c(6), c(9), c(11), c(10), c(10), c(2))
	strat := &mockStrategy{actions: []Action{Split, Stand, Stand, Stand, Stand}}
	st := &StatsTracker{}
	PlayRound(shoe, strat, st, false)
	if st.TotalWagered != 2.0 {
		t.Errorf("split two base wagers: TotalWagered want 2.0, got %v", st.TotalWagered)
	}
}

// Simulation EV / state-machine tests omitted: use benchmarks or separate sim package.
