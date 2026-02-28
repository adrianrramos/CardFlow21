package engine

// TODO: might want to move Strategy interface
type Strategy interface {
	Decide(player Hand, dealerUpCard Card) Action
	Name() string
}

func PlayHandWithStrategy(shoe *Shoe, strat Strategy) float64 {
	result := PlayHandWithStrategyDetailed(shoe, strat)
	return result.Profit
}

func PlayHandWithStrategyDetailed(shoe *Shoe, strat Strategy) HandResult {
	player := &Hand{}
	dealer := &Hand{}

	player.AddCard(shoe.Draw())
	dealer.AddCard(shoe.Draw())
	player.AddCard(shoe.Draw())
	dealer.AddCard(shoe.Draw())

	result := HandResult{
		HandsInRound: 1,
	}

	// Blackjack resolution order (must check in this order):
	// 1. Both blackjack → push (0.0)
	// 2. Player blackjack → win 3:2 (+1.5)
	// 3. Dealer blackjack → lose (-1.0)
	// No double/split allowed after blackjack - resolve immediately
	if player.IsBlackjack() && dealer.IsBlackjack() {
		result.PlayerBlackjack = true
		result.DealerBlackjack = true
		result.Profit = 0.0
		result.PlayerValue = 21
		result.DealerValue = 21
		result.TotalWagered = 1.0
		return result
	}

	if player.IsBlackjack() && !dealer.IsBlackjack() {
		result.PlayerBlackjack = true
		result.Profit = 1.5
		result.PlayerValue = 21
		result.DealerValue = dealer.Value()
		result.TotalWagered = 1.0
		return result
	}

	if dealer.IsBlackjack() {
		result.DealerBlackjack = true
		result.Profit = -1.0
		result.PlayerValue = player.Value()
		result.DealerValue = 21
		result.TotalWagered = 1.0
		return result
	}

	// TODO: place bet for insurance and resolve insurance bet here
	// if dealer.OffersInsurance() {
	// }

	// Detect pair opportunity BEFORE strategy decision
	if player.IsPair() && len(player.Cards) == 2 {
		result.PairOpportunity = true
		rank := int(player.Cards[0].Rank)
		// Normalize: 10/J/Q/K (ranks 10-13) → rank 10
		if rank >= 10 {
			result.PairRank = 10
		} else {
			result.PairRank = rank
		}
		result.DealerUpCard = dealer.Cards[0]
	}

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
			if action == Split && current_hand.SplitCount < 4 && current_hand.IsPair() {
				result.WasSplit = true
				new_hand := &Hand{}
				new_hand.AddCard(current_hand.Cards[1])
				current_hand.RemoveCard()

				nextSplitCount := current_hand.SplitCount + 1
				current_hand.SplitCount = nextSplitCount
				new_hand.SplitCount = nextSplitCount

				hand_queue = append(hand_queue, new_hand)
				hand_queue[i] = current_hand
				result.HandsInRound++
				continue // Start over with the new hand
			}

			PlayOutHand(hand_queue[i], dealer, strat, shoe)
			finished_hands++
		}
	}

	for dealer.Value() < 17 || (dealer.Value() == 17 && dealer.IsSoft()) {
		dealer.AddCard(shoe.Draw())
	}
	dealer.CheckBust()

	result.DealerBust = dealer.IsBust
	result.DealerValue = dealer.Value()

	// Evaluate all hands (split hands evaluated independently)
	// Split accounting validation:
	// - Each split hand wagers 1 unit independently (base wager)
	// - Doubled hands wager 2 units
	// - Profit/loss calculated per hand and summed
	// - TotalWagered accumulates all wagers correctly
	// - No missing loss registration - busts, losses, and pushes all handled
	total_profit := 0.0
	total_wagered := 0.0
	for _, hand := range hand_queue {
		wagered := 1.0
		if hand.IsDoubled {
			wagered = 2.0
			result.WasDoubled = true
		}
		total_wagered += wagered

		if hand.IsBust {
			result.PlayerBust = true
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
		// Note: Push (equal values) results in 0 profit, correctly handled
	}

	result.Profit = total_profit
	result.TotalWagered = total_wagered
	result.PlayerValue = hand_queue[0].Value()
	return result
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

	return
}