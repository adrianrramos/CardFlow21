package main

import (
	"cardflow21/simulation"
	"fmt"
	"flag"
)

func main() {
	rounds := flag.Int("rounds", 100000, "number of rounds")
	decks := flag.Int("decks", 6, "number of decks in the shoe")
	flag.Parse()

	stats := simulation.RunSimulation(*rounds, *decks)

	fmt.Println("==== CardFlow21 ====")
	fmt.Println("Hands: ", stats.TotalHands)
	fmt.Println("Profit: ", stats.Profit)
	fmt.Println("EV/Hand: ", stats.Mean)
	fmt.Println("Std Dev: ", stats.M2)
}