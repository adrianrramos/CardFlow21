# CardFlow21 – Working Notes
## TODOs

### P1 – High Priority

- [ ] Add Bankroll features
    - [ ] User can set starting bankroll
    - [ ] Game engine ends play if bankroll hits 0 or can no longer afford next wager
    - [ ] User can define bet spread (start with simple truncated true count, and only for count >= -3 || count <= +7, like BJA)
- [ ] Read other libraries for how they handle EV and hand automation

### P2 – Medium Priority

- [ ] Add calculations for n₀ (number of hands until profit is ~1 SD from mean EV)  
- [ ] Consider leveraging Strategy chart for game decision making
    - [ ] ie. Current logic performs decisions and game resolutions OUTSIDE of the chart -> Surrender parameters are not fully
            gathered from chart ie. only checks chart if count is <17 or >14 for example
    - [ ] GOAL: to be able to port many charts of many different strategy and not tweak the game engine at all
    - [ ] Consider the future where deviations will take a role and change behaviors of low value hands etc.
- [ ] Add support for additional players at the table  
- [ ] Add BJA deviations (Illustrious 18 / Fab 4, etc.)
- [ ] Try another basic strategy chart and compare outcomes
- [ ] Build basic web UI on LocalHost
- [ ] Graph results in web UI
    - [ ] if HTMX is used make sure support for charting / graphs is easy to implement
- [ ] Make commands to setup Web UI simple for easy portability
- [ ] Add Some of Composer’s metrics library on `agent-fix-1` branch
- [ ] TODO: fix issue with using files and flags in the same run command

### P3 – Lower Priority / Future Work

- [ ] Add kelly betting
- [ ] Broader deviation set beyond BJA baseline  
- [ ] Support different count systems (e.g. KO, Zen, Hi-Opt II, CAC2)

## Current Focus

- **Engine correctness**
  - Finish core H17 blackjack rules
  - Stabilize splitting (including 10-valued pairs and Aces)
  - Make payouts and EV calculations reflect real wagers

- **Simulation quality**
  - Use more realistic betting schemes (true count, bet spreads, bankroll)
  - Experiment with alternative basic strategy charts and deviations



## Scratchpad

- Cut card comes out and switches cut card state, so engine knows to finish current hand but shuffle on the next

## Methods for calculating EV & N0

**Welford's Online Algorithm**

```Python
class WelfordVariance:
    def __init__(self):   # Comparison to ShiftDataVariance:
        self.mean = 0.0   # = K + Ex / n
        self.count = 0    # = n
        self.M2 = 0.0     # = Ex2 - (Ex)^2 / n

    def add_variable(self, x: float):
        self.count += 1
        old_mean = self.mean
        self.mean += (x - self.mean) / self.count
        self.M2 += (x - old_mean) * (x - self.mean)

    def remove_variable(self, x: float):
        self.count -= 1
        new_mean = self.mean
        self.mean -= (x - self.mean) / self.count
        self.M2 -= (x - new_mean) * (x - self.mean)

    def get_mean(self) -> float:
        return self.mean

    def get_variance(self) -> float:
        return self.M2 / self.count

    def get_sample_variance(self) -> float:
        return self.M2 / (self.count - 1)
```
