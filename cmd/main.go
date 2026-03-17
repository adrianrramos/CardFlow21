package main

import (
	"cardflow21/simulation"
	"flag"
	"fmt"
)

func main() {
	rounds := flag.Int("rounds", 100000, "number of rounds")
	decks := flag.Int("decks", 6, "number of decks in the shoe")
	penetration := flag.Float64("penetration", .85, "% of shoe before cut card comes out")
	use_true_count := flag.Bool("use_true_count", false, "use true count to determine wager")
	flag.Parse()

	stats, statsTracker := simulation.RunSimulation(*rounds, *decks, *penetration, *use_true_count)

	fmt.Println("==== CardFlow21 ====")
	fmt.Println("Hands: ", statsTracker.TotalHands)
	fmt.Println("Profit: ", stats.Profit)
	fmt.Println("EV/Hand: ", stats.Mean)
	fmt.Println("Std Dev: ", stats.M2)
	fmt.Println("Doubles Won: ", statsTracker.DoubleWin)
	fmt.Println("Doubles Lost: ", statsTracker.DoubleLoss)
	fmt.Println("Splits: ", statsTracker.SplitHands)
}
