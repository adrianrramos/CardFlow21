package engine

import "testing"

// --- Test Helpers ---

// c creates a Card with the given rank (1=Ace, 2-10, 11=Jack, 12=Queen, 13=King).
func c(rank int) Card {
	return Card{Suit: 0, Rank: Rank(rank)}
}

// newTestShoe creates a Shoe that deals cards in the given argument order.
// Shoe.Draw() pops from the end of the internal slice, so we reverse the
// arguments so that the first argument is the first card dealt.
func newTestShoe(cards ...Card) *Shoe {
	// 1. Start from a real shoe
	s := NewShoe(6, 0.85) // or whatever decks/penetration you want for tests

	// 2. Ensure the shoe has room (it will)
	// 3. Overwrite the LAST len(cards) elements with your custom cards,
	//    in the order they should be drawn.
	for i, card := range cards {
		s.cards[len(s.cards)-len(cards)+i] = card
	}

	// Optionally reset count/true_count if you need them specific
	s.count = 0
	s.true_count = 0

	return s
}

// mockStrategy returns actions from a predetermined slice, defaulting to
// Stand once the slice is exhausted.
type mockStrategy struct {
	actions []Action
	idx     int
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

// CheckSurrenderChart implements the Strategy interface.
// The test suite does not cover surrender behavior, so default to "never surrender".
func (m *mockStrategy) CheckSurrenderChart(_ Hand, _ Card) bool { return false }

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
// the dealer has blackjack and the player does not.
//
// Current implementation resolves dealer blackjack as an immediate loss.
func TestPlayHandWithStrategy_DealerBlackjack(t *testing.T) {
	// Player: 5 + 6 = 11 (no BJ); Dealer: King + Ace = 21 BJ
	// Player hits a 10 to reach 21 → should still lose to dealer BJ, but
	// the missing early-exit produces a push instead.
	shoe := newTestShoe(c(5), c(13), c(6), c(1), c(10))
	strat := &mockStrategy{actions: []Action{Stand, Hit, Stand}}
	got := PlayRound(shoe, strat, &StatsTracker{}, false)
	if got != -1.0 {
		t.Errorf("dealer BJ: want -1.0, got %v (Bug: no early-exit for dealer blackjack)", got)
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
