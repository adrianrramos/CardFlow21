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

	var doubled bool
	for {
		action := strat.Decide(*player, dealer.Cards[0])
		
		// TODO: Double and Split feature not ready 
		if action == Double {
			doubled = true
			break
		}
		
		if action == Split {
			action = Hit
		}

		if action == Stand {
			break
		}

		if action == Hit {
			player.AddCard(shoe.Draw())
			if player.IsBust() {
				return -1
			}
		}
	}

	for dealer.Value() < 17 {
		dealer.AddCard(shoe.Draw())
	}

	if dealer.IsBust() {
		return 1
	}

	if player.Value() > dealer.Value() {
		if doubled {
			return 2
		}
		return 1
	} else if player.Value() < dealer.Value() {
		if doubled {
			return -2
		}
		return -1
	}

	return 0
}