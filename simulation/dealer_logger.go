package simulation

import (
	"cardflow21/engine"
	"fmt"
)

// DealerHandLog logs a single dealer hand progression
type DealerHandLog struct {
	HandNumber   int
	InitialCards []engine.Card
	CardsDrawn   []engine.Card
	Values       []int  // Value after each card
	IsSoft       []bool // Soft status after each card
	FinalValue   int
	FinalBust    bool
	HitOnSoft17  bool
	Correct      bool // Whether H17 rule was followed correctly
}

// DealerLogger logs dealer hands for H17 validation
type DealerLogger struct {
	logs []DealerHandLog
}

func NewDealerLogger() *DealerLogger {
	return &DealerLogger{
		logs: []DealerHandLog{},
	}
}

// LogDealerHand logs a dealer hand (called from game engine)
func (dl *DealerLogger) LogDealerHand(handNumber int, initialCards []engine.Card, cardsDrawn []engine.Card, finalValue int, finalBust bool) {
	log := DealerHandLog{
		HandNumber:   handNumber,
		InitialCards: make([]engine.Card, len(initialCards)),
		CardsDrawn:   make([]engine.Card, len(cardsDrawn)),
		FinalValue:   finalValue,
		FinalBust:    finalBust,
	}
	copy(log.InitialCards, initialCards)
	copy(log.CardsDrawn, cardsDrawn)
	
	// Reconstruct value progression
	hand := &engine.Hand{}
	for _, card := range initialCards {
		hand.AddCard(card)
		log.Values = append(log.Values, hand.Value())
		log.IsSoft = append(log.IsSoft, hand.IsSoft())
	}
	
	for _, card := range cardsDrawn {
		hand.AddCard(card)
		log.Values = append(log.Values, hand.Value())
		log.IsSoft = append(log.IsSoft, hand.IsSoft())
		
		// Check if hit on soft 17
		if hand.Value() == 17 && hand.IsSoft() {
			log.HitOnSoft17 = true
		}
	}
	
	// Validate H17 behavior
	log.Correct = dl.validateH17Behavior(log)
	
	dl.logs = append(dl.logs, log)
}

// validateH17Behavior checks if dealer correctly follows H17 rule
func (dl *DealerLogger) validateH17Behavior(log DealerHandLog) bool {
	// Check if dealer hit on soft 17
	for i, value := range log.Values {
		if value == 17 && log.IsSoft[i] {
			// Should have hit (unless it's the final value)
			if i < len(log.Values)-1 {
				// Hit occurred, which is correct
				return true
			} else {
				// Final value is soft 17 - should have hit, but didn't
				return false
			}
		}
	}
	
	// No soft 17 encountered, behavior is correct
	return true
}

// ValidateH17Behavior validates all logged dealer hands
func (dl *DealerLogger) ValidateH17Behavior() {
	fmt.Println("\n=== Dealer H17 Validation ===")
	fmt.Printf("Logged %d dealer hands\n\n", len(dl.logs))
	
	soft17Count := 0
	soft17HitCount := 0
	hard17StandCount := 0
	errors := 0
	
	for _, log := range dl.logs {
		// Check for soft 17
		for i, value := range log.Values {
			if value == 17 && log.IsSoft[i] {
				soft17Count++
				if i < len(log.Values)-1 {
					soft17HitCount++
				} else {
					// Final value is soft 17 - error!
					errors++
					fmt.Printf("Hand #%d: ERROR - Dealer stood on soft 17!\n", log.HandNumber)
					fmt.Printf("  Cards: ")
					for _, c := range log.InitialCards {
						fmt.Printf("%d ", c.Rank)
					}
					for _, c := range log.CardsDrawn {
						fmt.Printf("%d ", c.Rank)
					}
					fmt.Printf("\n  Final: %d (soft)\n", log.FinalValue)
				}
				break
			}
		}
		
		// Check for hard 17 stand
		if log.FinalValue == 17 && !log.FinalBust {
			// Check if it's hard 17
			finalHand := &engine.Hand{}
			for _, c := range log.InitialCards {
				finalHand.AddCard(c)
			}
			for _, c := range log.CardsDrawn {
				finalHand.AddCard(c)
			}
			if !finalHand.IsSoft() {
				hard17StandCount++
			}
		}
		
		if !log.Correct {
			errors++
		}
	}
	
	fmt.Printf("\nSummary:\n")
	fmt.Printf("  Soft 17 encountered: %d\n", soft17Count)
	fmt.Printf("  Soft 17 hit correctly: %d\n", soft17HitCount)
	fmt.Printf("  Hard 17 stood correctly: %d\n", hard17StandCount)
	fmt.Printf("  Errors: %d\n", errors)
	
	if errors == 0 {
		fmt.Println("\n✓ All dealer H17 decisions correct")
	} else {
		fmt.Printf("\n*** FOUND %d H17 RULE VIOLATIONS ***\n", errors)
	}
}

// PrintSampleHands prints a sample of dealer hands
func (dl *DealerLogger) PrintSampleHands(count int) {
	fmt.Printf("\n=== Sample Dealer Hands (first %d) ===\n", count)
	
	printed := 0
	for _, log := range dl.logs {
		if printed >= count {
			break
		}
		
		fmt.Printf("\nHand #%d:\n", log.HandNumber)
		fmt.Printf("  Initial: ")
		for _, c := range log.InitialCards {
			fmt.Printf("%d ", c.Rank)
		}
		fmt.Printf("\n  Drawn: ")
		for _, c := range log.CardsDrawn {
			fmt.Printf("%d ", c.Rank)
		}
		fmt.Printf("\n  Progression: ")
		for i, v := range log.Values {
			soft := ""
			if log.IsSoft[i] {
				soft = "s"
			}
			fmt.Printf("%d%s ", v, soft)
		}
		fmt.Printf("\n  Final: %d", log.FinalValue)
		if log.FinalBust {
			fmt.Printf(" (BUST)")
		}
		if log.HitOnSoft17 {
			fmt.Printf(" [Hit on soft 17]")
		}
		fmt.Println()
		
		printed++
	}
}
