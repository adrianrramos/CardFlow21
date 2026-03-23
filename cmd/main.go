package main

import (
	"cardflow21/simulation"
	"cardflow21/strategy"
	"flag"
	"fmt"
	"os"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Rounds       *int     `yaml:"rounds"`
	Decks        *int     `yaml:"decks"`
	Penetration  *float64 `yaml:"penetration"`
	UseTrueCount *bool    `yaml:"use_true_count"`
	Strategy     *string  `yaml:"strategy"`
}

// Using no pointers
type FinalConfig struct {
	Rounds       int
	Decks        int
	Penetration  float64
	UseTrueCount bool
	Strategy     string
}

func main() {
	file := flag.String("file", "default", "config name")

	rounds := flag.Int("rounds", 100000, "number of rounds")
	decks := flag.Int("decks", 6, "number of decks in the shoe")
	penetration := flag.Float64("penetration", .85, "% of shoe before cut card comes out")
	use_true_count := flag.Bool("use_true_count", false, "use true count to determine wager")
	strategy_name := flag.String("strategy", "BJA", "strategy to use: BJA or BJ101App")

	flag.Parse()
	cfg := Config{
		Rounds:       rounds,
		Decks:        decks,
		Penetration:  penetration,
		UseTrueCount: use_true_count,
		Strategy:     strategy_name,
	}

	configPath := fmt.Sprintf(".configs/%s.yaml", strings.ToLower(*file))
	fileCfg, err := LoadConfig(configPath)
	if err == nil {
		Merge(&cfg, fileCfg)
	} else {
		fmt.Println("No config file loaded: ", err)
	}
	final := resolve(cfg)

	var stratName strategy.StrategyName
	switch strings.ToUpper(strings.TrimSpace(final.Strategy)) {
	case "BJA":
		stratName = strategy.BJA
	case "BJ101APP":
		stratName = strategy.BJ101App
	default:
		// Keep CLI behavior simple: default to BJA on unknown values.
		stratName = strategy.BJA
	}

	stats, statsTracker := simulation.RunSimulation(final.Rounds, final.Decks, final.Penetration, final.UseTrueCount, stratName)

	p := message.NewPrinter(language.English)

	fmt.Println("==== CardFlow21 ====")
	fmt.Printf("%-17s %v\n", "Profit: ", stats.Profit)
	fmt.Printf("%-17s %v\n", "EV/Wagered: ", stats.Mean)
	fmt.Printf("%-17s %v\n", "EV/Hand: ", stats.MeanHands)
	fmt.Printf("%-17s %v\n", "Std Dev: ", stats.StdDev())
	fmt.Printf("%-17s %v\n", "Hands: ", p.Sprintf("%d", statsTracker.TotalHands))
	fmt.Printf("%-17s %v\n", "Doubles Won: ", statsTracker.DoubleWin)
	fmt.Printf("%-17s %v\n", "Doubles Lost: ", statsTracker.DoubleLoss)
	fmt.Printf("%-17s %v\n", "Splits: ", statsTracker.SplitHands)
	fmt.Printf("%-17s %v\n", "Wagered: ", statsTracker.TotalWagered)
	fmt.Printf("%-17s %v\n", "Insurance Took: ", statsTracker.TookInsurance)
	fmt.Printf("%-17s %v\n", "Surrendered: ", statsTracker.Surrendered)
}

func LoadConfig(path string) (Config, error) {
	var cfg Config

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}

	err = yaml.Unmarshal(data, &cfg)
	return cfg, err
}

func Merge(base *Config, override Config) {
	if override.Rounds != nil {
		base.Rounds = override.Rounds
	}
	if override.Decks != nil {
		base.Decks = override.Decks
	}
	if override.Penetration != nil {
		base.Penetration = override.Penetration
	}
	if override.UseTrueCount != nil {
		base.UseTrueCount = override.UseTrueCount
	}
	if override.Strategy != nil {
		base.Strategy = override.Strategy
	}
}

func resolve(cfg Config) FinalConfig {
	return FinalConfig{
		Rounds:       getInt(cfg.Rounds, 100000),
		Decks:        getInt(cfg.Decks, 6),
		Penetration:  getFloat(cfg.Penetration, 0.85),
		UseTrueCount: getBool(cfg.UseTrueCount, false),
		Strategy:     getString(cfg.Strategy, "BJA"),
	}
}

// ---- HELPERS ----

func getInt(v *int, def int) int {
	if v != nil {
		return *v
	}
	return def
}

func getFloat(v *float64, def float64) float64 {
	if v != nil {
		return *v
	}
	return def
}

func getBool(v *bool, def bool) bool {
	if v != nil {
		return *v
	}
	return def
}

func getString(v *string, def string) string {
	if v != nil {
		return *v
	}
	return def
}
