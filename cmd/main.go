package main

import (
	"cardflow21/simulation"
	"cardflow21/strategy"
	"flag"
	"fmt"
	"strings"
)

func main() {
	rounds := flag.Int("rounds", 100000, "number of rounds")
	decks := flag.Int("decks", 6, "number of decks in the shoe")
	penetration := flag.Float64("penetration", .85, "% of shoe before cut card comes out")
	use_true_count := flag.Bool("use_true_count", false, "use true count to determine wager")
	strategy_name := flag.String("strategy", "BJA", "strategy to use: BJA or BJ101App")
	flag.Parse()

	var stratName strategy.StrategyName
	switch strings.ToUpper(strings.TrimSpace(*strategy_name)) {
	case "BJA":
		stratName = strategy.BJA
	case "BJ101APP":
		stratName = strategy.BJ101App
	default:
		// Keep CLI behavior simple: default to BJA on unknown values.
		stratName = strategy.BJA
	}

	stats, statsTracker := simulation.RunSimulation(*rounds, *decks, *penetration, *use_true_count, stratName)

	fmt.Println("==== CardFlow21 ====")
	fmt.Println("Hands: ", statsTracker.TotalHands)
	fmt.Println("Profit: ", stats.Profit)
	fmt.Println("EV/Hand: ", stats.Mean)
	fmt.Println("Std Dev: ", stats.M2)
	fmt.Println("Doubles Won: ", statsTracker.DoubleWin)
	fmt.Println("Doubles Lost: ", statsTracker.DoubleLoss)
	fmt.Println("Splits: ", statsTracker.SplitHands)
	fmt.Println("Total Wagered: ", statsTracker.TotalWagered)
}
