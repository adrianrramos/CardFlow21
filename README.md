# CardFlow21

CardFlow21 is a blackjack engine and simulation toolkit focused on analyzing expected value (EV) and betting strategies for multi-deck H17 games. It provides a rules-accurate game engine, a basic strategy module, and a simulation runner for exploring long‑run performance of different play and bet schemes.

## Features

- **Blackjack engine**  
  - Multi-deck shoe with configurable deck count and penetration  
  - Dealer H17 rules (dealer hits soft 17)  
  - Standard actions: hit, stand, double, split, (future: surrender, insurance)  

- **Strategy module**  
  - Basic strategy for DAS, H17, 4+ deck games  
  - Pluggable strategy interface for custom decision logic

- **Simulation runner**  
  - Command-line entrypoint to run large numbers of hands  
  - Hi-Lo counting with true count support for bet sizing experiments  
  - Basic statistics: total profit, EV per hand, variance-related measures

## Usage

From the project root:

```bash
go run ./cmd \
  -rounds 100000 \
  -decks 6 \
  -penetration 0.85 \
  -use_true_count
```

### Using .config YAML files
From project root go to `./.configs/` and create a new YAML file. Below is an
an example of `benchmark.yaml` with the following parameters:
```yaml
rounds: 100000000
decks: 6
penetration: 0.85
use_true_count: false
strategy: BJA
detailed_logs: true
```
Now to run this file simply run this command from project root:
```bash
go run ./cmd --file benchmark
```
