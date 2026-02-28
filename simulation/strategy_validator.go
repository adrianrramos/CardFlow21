package simulation

import (
	"cardflow21/engine"
	"cardflow21/strategy"
	"fmt"
)

// StrategyValidator validates strategy matrix implementation
type StrategyValidator struct {
	basicStrategy *strategy.BasicStrategy
}

func NewStrategyValidator() *StrategyValidator {
	return &StrategyValidator{
		basicStrategy: strategy.NewBasicStrategy(),
	}
}

// ValidateHardTotals tests hard total decisions
func (sv *StrategyValidator) ValidateHardTotals() {
	fmt.Println("\n=== Hard Totals Validation ===")
	
	testCases := []struct {
		cards    []int // Card ranks
		dealer   int   // Dealer upcard rank
		expected string // Expected action
		desc     string
	}{
		// Hard 16 vs 10 (should stand)
		{[]int{10, 6}, 10, "S", "Hard 16 vs 10"},
		// Hard 11 vs 10 (should double)
		{[]int{6, 5}, 10, "D", "Hard 11 vs 10"},
		// Hard 12 vs 3 (should stand)
		{[]int{7, 5}, 3, "S", "Hard 12 vs 3"},
		// Hard 9 vs 5 (should double)
		{[]int{5, 4}, 5, "D", "Hard 9 vs 5"},
		// Hard 10 vs 10 (should hit)
		{[]int{6, 4}, 10, "H", "Hard 10 vs 10"},
	}
	
	errors := 0
	for _, tc := range testCases {
		hand := &engine.Hand{}
		for _, rank := range tc.cards {
			hand.AddCard(engine.Card{Rank: engine.Rank(rank)})
		}
		
		dealerCard := engine.Card{Rank: engine.Rank(tc.dealer)}
		action := sv.basicStrategy.Decide(*hand, dealerCard)
		
		expectedAction := actionToString(tc.expected, hand)
		actualAction := actionToSymbol(action)
		
		match := (expectedAction == actualAction)
		status := "✓"
		if !match {
			status = "✗"
			errors++
		}
		
		fmt.Printf("%s %s: Expected %s, Got %s", status, tc.desc, expectedAction, actualAction)
		if !match {
			fmt.Printf(" *** MISMATCH ***")
		}
		fmt.Println()
	}
	
	if errors == 0 {
		fmt.Println("\n✓ All hard total tests passed")
	} else {
		fmt.Printf("\n*** FOUND %d HARD TOTAL MISMATCHES ***\n", errors)
	}
}

// ValidateSoftTotals tests soft total decisions
func (sv *StrategyValidator) ValidateSoftTotals() {
	fmt.Println("\n=== Soft Totals Validation ===")
	
	testCases := []struct {
		cards    []int
		dealer   int
		expected string
		desc     string
	}{
		// Soft 18 vs 9 (should hit)
		{[]int{1, 7}, 9, "H", "Soft 18 vs 9"},
		// Soft 18 vs 6 (should double)
		{[]int{1, 7}, 6, "Ds", "Soft 18 vs 6"},
		// Soft 19 vs 5 (should double)
		{[]int{1, 8}, 5, "Ds", "Soft 19 vs 5"},
		// Soft 20 vs 10 (should stand)
		{[]int{1, 9}, 10, "S", "Soft 20 vs 10"},
	}
	
	errors := 0
	for _, tc := range testCases {
		hand := &engine.Hand{}
		for _, rank := range tc.cards {
			hand.AddCard(engine.Card{Rank: engine.Rank(rank)})
		}
		
		dealerCard := engine.Card{Rank: engine.Rank(tc.dealer)}
		action := sv.basicStrategy.Decide(*hand, dealerCard)
		
		expectedAction := actionToString(tc.expected, hand)
		actualAction := actionToSymbol(action)
		
		match := (expectedAction == actualAction)
		status := "✓"
		if !match {
			status = "✗"
			errors++
		}
		
		fmt.Printf("%s %s: Expected %s, Got %s", status, tc.desc, expectedAction, actualAction)
		if !match {
			fmt.Printf(" *** MISMATCH ***")
		}
		fmt.Println()
	}
	
	if errors == 0 {
		fmt.Println("\n✓ All soft total tests passed")
	} else {
		fmt.Printf("\n*** FOUND %d SOFT TOTAL MISMATCHES ***\n", errors)
	}
}

// ValidatePairSplits tests pair splitting decisions
func (sv *StrategyValidator) ValidatePairSplits() {
	fmt.Println("\n=== Pair Split Validation ===")
	
	testCases := []struct {
		rank     int
		dealer   int
		expected bool // true = split, false = don't split
		desc     string
	}{
		// Pair 8 vs 10 (should split)
		{8, 10, true, "Pair 8s vs 10"},
		// Pair 10 vs 10 (should not split)
		{10, 10, false, "Pair 10s vs 10"},
		// Pair Aces vs 10 (should split)
		{1, 10, true, "Pair Aces vs 10"},
		// Pair 5s vs 10 (should not split)
		{5, 10, false, "Pair 5s vs 10"},
	}
	
	errors := 0
	for _, tc := range testCases {
		hand := &engine.Hand{}
		hand.AddCard(engine.Card{Rank: engine.Rank(tc.rank)})
		hand.AddCard(engine.Card{Rank: engine.Rank(tc.rank)})
		
		dealerCard := engine.Card{Rank: engine.Rank(tc.dealer)}
		action := sv.basicStrategy.Decide(*hand, dealerCard)
		
		actualSplit := (action == engine.Split)
		match := (tc.expected == actualSplit)
		status := "✓"
		if !match {
			status = "✗"
			errors++
		}
		
		expectedStr := "Split"
		if !tc.expected {
			expectedStr = "Don't Split"
		}
		actualStr := "Split"
		if !actualSplit {
			actualStr = "Don't Split"
		}
		
		fmt.Printf("%s %s: Expected %s, Got %s", status, tc.desc, expectedStr, actualStr)
		if !match {
			fmt.Printf(" *** MISMATCH ***")
		}
		fmt.Println()
	}
	
	if errors == 0 {
		fmt.Println("\n✓ All pair split tests passed")
	} else {
		fmt.Printf("\n*** FOUND %d PAIR SPLIT MISMATCHES ***\n", errors)
	}
}

// ValidateDoubleRules tests double-down rules
func (sv *StrategyValidator) ValidateDoubleRules() {
	fmt.Println("\n=== Double Rules Validation ===")
	
	// Test that doubles only allowed on 2-card hands
	hand2Cards := &engine.Hand{}
	hand2Cards.AddCard(engine.Card{Rank: 6})
	hand2Cards.AddCard(engine.Card{Rank: 5})
	
	hand3Cards := &engine.Hand{}
	hand3Cards.AddCard(engine.Card{Rank: 6})
	hand3Cards.AddCard(engine.Card{Rank: 5})
	hand3Cards.AddCard(engine.Card{Rank: 2})
	
	dealerCard := engine.Card{Rank: 10}
	
	action2Cards := sv.basicStrategy.Decide(*hand2Cards, dealerCard)
	action3Cards := sv.basicStrategy.Decide(*hand3Cards, dealerCard)
	
	fmt.Printf("Hard 11 (2 cards) vs 10: %s", actionToSymbol(action2Cards))
	if action2Cards == engine.Double {
		fmt.Println(" ✓ (can double)")
	} else {
		fmt.Println(" ✗ (should be able to double)")
	}
	
	fmt.Printf("Hard 13 (3 cards) vs 10: %s", actionToSymbol(action3Cards))
	if action3Cards != engine.Double {
		fmt.Println(" ✓ (cannot double)")
	} else {
		fmt.Println(" ✗ (should NOT be able to double)")
	}
}

func actionToString(expected string, hand *engine.Hand) string {
	switch expected {
	case "H":
		return "Hit"
	case "S":
		return "Stand"
	case "D":
		if hand.IsTwoCardTotal() {
			return "Double"
		}
		return "Hit"
	case "Ds":
		if hand.IsTwoCardTotal() {
			return "Double"
		}
		return "Stand"
	default:
		return expected
	}
}

func actionToSymbol(action engine.Action) string {
	switch action {
	case engine.Hit:
		return "Hit"
	case engine.Stand:
		return "Stand"
	case engine.Double:
		return "Double"
	case engine.Split:
		return "Split"
	default:
		return "Unknown"
	}
}
