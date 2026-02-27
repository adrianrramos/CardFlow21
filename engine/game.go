package engine

// TODO: might want to move Strategy interface
type Strategy interface {
	Decide(player Hand, dealerUpCard Card) Action
	Name() string
}

func PlayHandWithStrategy(shoe *Shoe, strat Strategy) float64 {
	player := &Hand{}
	dealer := &Hand{}

	player.AddCard(shoe.Draw())
	dealer.AddCard(shoe.Draw())
	player.AddCard(shoe.Draw())
	dealer.AddCard(shoe.Draw())

	if player.IsBlackjack() && !dealer.IsBlackjack() {
		return 1.5
	}

	// TODO: place bet for insurance and resolve insurance bet here
	// if dealer.OffersInsurance() {
	// }

	hand_queue := []*Hand{ player } 
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
			if action == Split && current_hand.SplitCount < 4 {
				new_hand := &Hand{}
				new_hand.AddCard(current_hand.Cards[1])
				current_hand.RemoveCard()

				current_hand.SplitCount++
				new_hand.SplitCount++

				hand_queue = append(hand_queue, new_hand)
				hand_queue[i] = current_hand
				continue // Start over with the new hand
			}

			PlayOutHand(hand_queue[i], dealer, strat, shoe)
			finished_hands++
		}
	}

	for dealer.Value() < 17 {
		dealer.AddCard(shoe.Draw())
	}
	dealer.CheckBust()

	// Evaluate all hands
	total_profit := 0
	for _, hand := range hand_queue {
		wagered := 1
		if hand.IsDoubled {
			wagered *= 2
		}

		if hand.IsBust {
			total_profit -= wagered
			continue
		}
		if dealer.IsBust {
			total_profit += wagered
			continue
		}
		if hand.Value() > dealer.Value() {
			total_profit += wagered
		}
		if hand.Value() < dealer.Value() {
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

	return
}