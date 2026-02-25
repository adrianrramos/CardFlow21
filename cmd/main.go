package main

import (
	"cardflow21/simulation"
	"fmt"
	"flag"
)

func main() {
	rounds := flag.Int("rounds", 100000, "number of rounds")
	decks := flag.Int("decks", 6, "number of decks in the shoe")
	penetration := flag.Float64("penetration", .85, "% of shoe before cut card comes out") 
	flag.Parse()

	stats := simulation.RunSimulation(*rounds, *decks, *penetration)

	fmt.Println("==== CardFlow21 ====")
	fmt.Println("Hands: ", stats.TotalHands)
	fmt.Println("Profit: ", stats.Profit)
	fmt.Println("EV/Hand: ", stats.Mean)
	fmt.Println("Std Dev: ", stats.M2)
}