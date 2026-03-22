package simulation

type BankRoll struct {
	Balance float64
}

func (bank *BankRoll) SetBankRoll(amount float64) {
	bank.Balance = amount
}

// Takes positive / negative values
func (bank *BankRoll) AdjustBank(amount float64) {
	bank.Balance += amount
}



// TODO: adjust this according blackjack gameplay and spread ie. 1 full kelly, 2 full kelly based on true count
// KellyCriterion calculates the optimal percentage to wager
func (bank *BankRoll) KellyCriterion(winProb float64, netOdds float64) float64 {
    // f* = (bp - q) / b
    // winProb = p
    // netOdds = b
    // q = 1 - p
    
    q := 1 - winProb
    fraction := (winProb*netOdds - q) / netOdds
    
    if fraction < 0 {
        return 0 // No edge, do not bet
    }

    return fraction * bank.Balance
}

