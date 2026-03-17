package engine

// import "fmt"

// TODO: might want to move Strategy interface
type Strategy interface {
	Decide(player Hand, dealerUpCard Card) Action
	Name() string
}

type StatsTracker struct {
	TotalHands int
	DoubleWin  int
	DoubleLoss int
	SplitHands int
}

func PlayHandWithStrategy(shoe *Shoe, strat Strategy, statsTracker *StatsTracker, use_true_count bool) float64 {
	var wagered int
	if use_true_count && shoe.true_count > 0 {
		wagered = shoe.true_count * 2
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

	if player.IsBlackjack() && !dealer.IsBlackjack() {
		return float64(wagered) * 1.5
	} else if player.IsBlackjack() && dealer.IsBlackjack() {
		return 0
	}

	// TODO: place bet for insurance and resolve insurance bet here
	if dealer.OffersInsurance() {
		if shoe.true_count > 3 && dealer.IsBlackjack() {
			return 0
		}
	}
	if dealer.IsBlackjack() {
		return float64(-wagered)
	}

	hand_queue := []*Hand{player}
	// create a loop that executes moves for each hand in the list
	// if more hands are added to the list (ie, splitting or other players)
	// keep executing each hand
	finished_hands := 0
	for finished_hands < len(hand_queue) {
		q_length := len(hand_queue)
		for i := 0; i < q_length; i++ {
			current_hand := hand_queue[i]

			// Check if hand was split and needs to be dealt a card
			if len(current_hand.Cards) < 2 {
				current_hand.AddCard(shoe.Draw())
			}

			// Lookup action to catch any splits
			action := strat.Decide(*current_hand, dealer.Cards[0])
			if action == Split && current_hand.SplitCount < 4 && current_hand.IsPair() {
				statsTracker.SplitHands++
				statsTracker.TotalHands++
				new_hand := &Hand{}
				new_hand.AddCard(current_hand.Cards[1])
				current_hand.RemoveCard()

				nextSplitCount := current_hand.SplitCount + 1
				current_hand.SplitCount = nextSplitCount
				new_hand.SplitCount = nextSplitCount

				hand_queue = append(hand_queue, new_hand)
				hand_queue[i] = current_hand
				continue // Start over with the new hand
			}

			PlayOutHand(hand_queue[i], dealer, strat, shoe)
			finished_hands++
		}
	}

	// Hit 17 rules
	for dealer.Value() < 17 || (dealer.Value() <= 17 && dealer.IsSoft()) {
		dealer.AddCard(shoe.Draw())
	}
	dealer.CheckBust()


	// Evaluate all hands
	total_profit := 0
	for _, hand := range hand_queue {
		if hand.IsBust {
			total_profit -= wagered
			continue
		}
		// TODO: missing payouts for doubled hands (currently does not 2x the wager)
		if dealer.IsBust {
			if hand.IsDoubled {
				statsTracker.DoubleWin++
			}
			continue
		}
		if hand.Value() > dealer.Value() {
			if hand.IsDoubled {
				statsTracker.DoubleWin++
			}
			total_profit += wagered
		}
		if hand.Value() < dealer.Value() {
			if hand.IsDoubled {
				statsTracker.DoubleLoss++
			}
			total_profit -= wagered
		}
	}

	return float64(total_profit)
}

// Takes a valid hand that cant split and plays it out until it busts or stands
func PlayOutHand(player *Hand, dealer *Hand, strat Strategy, shoe *Shoe) {
	for {
		action := strat.Decide(*player, dealer.Cards[0])
		
		// TODO: Double and Split feature not ready 
		if action == Double {
			player.Doubled()
			player.AddCard(shoe.Draw())
			break
		}
		
		if action == Split {
			action = Hit
		}

		if action == Stand {
			return
		}

		if action == Hit {
			player.AddCard(shoe.Draw())
			if player.CheckBust() {
				return
			}
		}
	}

}
