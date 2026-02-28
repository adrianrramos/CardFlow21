package strategy

import "cardflow21/engine"

// NoAdvancedStrategy wraps a strategy and disables doubles and splits
type NoAdvancedStrategy struct {
	base engine.Strategy
}

func NewNoAdvancedStrategy(base engine.Strategy) *NoAdvancedStrategy {
	return &NoAdvancedStrategy{base: base}
}

func (n *NoAdvancedStrategy) Decide(player engine.Hand, dealerUpCard engine.Card) engine.Action {
	action := n.base.Decide(player, dealerUpCard)
	
	// Disable doubles and splits
	if action == engine.Double {
		// For "D" actions, hit instead
		// For "Ds" actions, stand instead
		// We'll default to hit for simplicity
		return engine.Hit
	}
	
	if action == engine.Split {
		return engine.Stand
	}
	
	return action
}

func (n *NoAdvancedStrategy) Name() string {
	return "No Advanced (" + n.base.Name() + ")"
}
