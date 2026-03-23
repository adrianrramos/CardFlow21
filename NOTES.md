# CardFlow21 – Working Notes
## TODOs

### P1 – High Priority

- [ ] Do a [SANITY_CHECK](SANITY_CHECK.md) walkthrough
- [ ] Double check logic for pair of 8's is not being pre-maturely surrendered 
- [ ] DO NOT MOVE ON to AP simming until
    - [ ] Baseline EV is established
    - [ ] Proper metrics tracking is handled
    - [ ] Library is well tested where it matters most
- [ ] Consider leveraging Strategy chart for game decision making
    - [ ] ie. Current logic performs decisions and game resolutions OUTSIDE of the chart -> Surrender parameters are not fully 
            gathered from chart ie. only checks chart if count is <17 or >14 for example
    - [ ] GOAL: to be able to port many charts of many different strategy and not tweak the game engine at all
    - [ ] Consider the future where deviations will take a role and change behaviors of low value hands etc. 
- [ ] Add Bankroll features
    - [ ] User can set starting bankroll
    - [ ] Game engine ends play if bankroll hits 0 or can no longer afford next wager
    - [ ] User can define bet spread (start with simple truncated true count, and only for count >= -3 || count <= +7, like BJA)

### P2 – Medium Priority

- [ ] Add calculations for n₀ (number of hands until profit is ~1 SD from mean EV)  

- [ ] Add support for additional players at the table  
- [ ] Add BJA deviations (Illustrious 18 / Fab 4, etc.)
- [ ] Try another basic strategy chart and compare outcomes
- [ ] Build basic web UI on LocalHost
- [ ] Graph results in web UI
    - [ ] if HTMX is used make sure support for charting / graphs is easy to implement
- [ ] Make commands to setup Web UI simple for easy portability
- [ ] Add Some of Composer’s metrics library on `agent-fix-1` branch

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
