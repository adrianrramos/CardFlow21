package engine

import "fmt"

type Strategy interface {
	Decide(player Hand, dealerUpCard Card) Action
	CheckSurrenderChart(player Hand, dealerUpCard Card) bool
	Name() string
}

type StatsTracker struct {
	TotalHands    int
	DoubleWin     int
	DoubleLoss    int
	SplitHands    int
	Surrendered   int
	TookInsurance int
	TotalWagered  float64
}

func PlayRound(shoe *Shoe, strat Strategy, statsTracker *StatsTracker, use_true_count bool) float64 {
	total_profit := 0.0
	var wagered float64
	if use_true_count && shoe.true_count > 0 {
		wagered = float64(shoe.true_count) * 2
	} else {
		wagered = 1
	}

	statsTracker.TotalHands++
	player := &Hand{}
	dealer := &Hand{}

	player.AddCard(shoe.Draw())
	dealer.AddCard(shoe.Draw())
	player.AddCard(shoe.Draw())
	dealer.AddCard(shoe.Draw())

	if shouldSurrender := strat.CheckSurrenderChart(*player, dealer.Cards[0]); shouldSurrender {
		statsTracker.Surrendered++
		return (wagered / 2) * -1.0 // surrendering causes player to LOSE half their bet
	}

	if player.IsBlackjack() && !dealer.IsBlackjack() {
		return wagered * 1.5
	} else if player.IsBlackjack() && dealer.IsBlackjack() {
		return 0
	}

	// Taking Insurance
	if use_true_count && shoe.true_count >= 3 && dealer.OffersInsurance() {
		statsTracker.TookInsurance++
		if dealer.IsBlackjack() {
			// Hand is dead, insurance pays 2:1 covering your original bet
			return 0
		}

		insurance_ammount := float64(wagered) / 2
		total_profit -= insurance_ammount
	}

	if dealer.IsBlackjack() {
		return float64(-wagered)
	}

	hands_stack := []*Hand{player}
	// create a loop that executes moves for each hand in the list
	// if more hands are added to the list (ie, splitting or other players)
	// keep executing each hand
	finished_hands := 0
	for finished_hands < len(hands_stack) {
		q_length := len(hands_stack)
		for i := 0; i < q_length; i++ {
			current_hand := hands_stack[i]

			// Check if hand was split and needs to be dealt a second card
			if len(current_hand.Cards) < 2 {
				current_hand.AddCard(shoe.Draw())
			}

			// Lookup action to catch any splits
			action := strat.Decide(*current_hand, dealer.Cards[0])

			// Handle Splitting Action
			if action == Split && current_hand.CanSplit() {
				statsTracker.SplitHands++
				new_hand := &Hand{}
				new_hand.AddCard(current_hand.Cards[1])
				current_hand.RemoveCard()

				nextSplitCount := current_hand.SplitCount + 1
				current_hand.SplitCount, new_hand.SplitCount = nextSplitCount, nextSplitCount

				hands_stack = append(hands_stack, new_hand)
				continue // Start over with the new hand
			}

			PlayOutHand(hands_stack[i], dealer, strat, shoe)
			finished_hands++
		}
	}

	// Hit 17 rules
	for dealer.Value() < 17 || (dealer.Value() <= 17 && dealer.IsSoft()) {
		dealer.AddCard(shoe.Draw())
	}
	dealer.CheckBust()

	// Evaluate all hands
	for _, hand := range hands_stack {
		var hand_wager float64
		if hand.IsDoubled {
			hand_wager = wagered * 2
		} else {
			hand_wager = wagered
		}
		statsTracker.TotalWagered += hand_wager

		if hand.CheckBust() {
			total_profit -= hand_wager
			continue
		}
		if dealer.CheckBust() {
			if hand.IsDoubled {
				statsTracker.DoubleWin++
			}
			total_profit += hand_wager
			continue
		}
		if hand.Value() > dealer.Value() {
			if hand.IsDoubled {
				statsTracker.DoubleWin++
			}
			total_profit += hand_wager
		} else if hand.Value() < dealer.Value() {
			if hand.IsDoubled {
				statsTracker.DoubleLoss++
			}
			total_profit -= hand_wager
		}
		// Push if values are equal
	}

	// Shuffle the shoe if the cut card has been dealt
	if shoe.cut_card_out {
		shoe.ShuffleShoe()
	}

	return float64(total_profit)
}

/*
PlayOutHand
player: Hand to play out
dealer: Dealer's hand to access up card and value
strat: Strategy to use
shoe: Shoe to use

This method when given a valid hand will play it out until it busts or stands
The state of the hand is set on the hand object itself and can be accessed by
the caller to determine if the hand was busted or not.

Only hitting, standing, and doubling are supported; because splitting and surrendering
effect the state of the game, so they are handled by the caller.
*/
func PlayOutHand(player *Hand, dealer *Hand, strat Strategy, shoe *Shoe) {
	if player.IsAceSplit() {
		return
	}

	for {
		action := strat.Decide(*player, dealer.Cards[0])

		switch action {
		case Double:
			player.Doubled()
			player.AddCard(shoe.Draw())
			return
		case Hit, Split: // Splittable hand should be treated as hard value
			player.AddCard(shoe.Draw())
			if player.CheckBust() {
				return
			}
		case Stand:
			return
		default:
			fmt.Printf("Action of: %v is not recognized\n", action)
			return
		}
	}
}
