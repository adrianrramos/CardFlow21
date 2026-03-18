package engine

import "fmt"

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
	TotalWagered int
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

	if dealer.OffersInsurance() {
		if shoe.true_count > 3 && dealer.IsBlackjack() {
			return 0
		} 
		// TODO: need to handle alternative where player loses insurance bet
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
		if hand.IsDoubled {
			wagered *= 2
		}

		if hand.IsBust {
			total_profit -= wagered
			continue
		}
		if dealer.IsBust {
			if hand.IsDoubled {
				statsTracker.DoubleWin++
			}
			total_profit += wagered
			continue
		}
		if hand.Value() > dealer.Value() {
			if hand.IsDoubled {
				statsTracker.DoubleWin++
			}
			total_profit += wagered
		} else if hand.Value() < dealer.Value() {
			if hand.IsDoubled {
				statsTracker.DoubleLoss++
			}
			total_profit -= wagered
		}
		// Push if values are equal
	}

	// Shuffle the shoe if the cut card has been dealt
	if shoe.cut_card_out {
		shoe.ShuffleShoe()
	}

	statsTracker.TotalWagered += wagered
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
		case Surrender:
			fmt.Println("Surrender not supported yet")
			return
		}
	}
}
