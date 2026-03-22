# CardFlow21 – Working Notes

## Current Focus

- **Engine correctness**
  - Finish core H17 blackjack rules
  - Stabilize splitting (including 10-valued pairs and Aces)
  - Make payouts and EV calculations reflect real wagers

- **Simulation quality**
  - Use more realistic betting schemes (true count, bet spreads, bankroll)
  - Experiment with alternative basic strategy charts and deviations

## TODOs

### P1 – High Priority

- [ ] Do a [SANITY_CHECK](SANITY_CHECK.md) walkthrough
- [ ] Add Some of Composer’s metrics library on `agent-fix-1` branch
- [ ] Try another basic strategy chart and compare outcomes

### P2 – Medium Priority

- [ ] Add proper wagers, payouts, and bankroll tracking  
- [ ] Add kelly bets
- [ ] Add calculations for n₀ (hands required to be ~1 SD from EV)  
- [ ] Add support for additional players at the table  
- [ ] Implement bet spreads (bet ramp by true count)  
- [ ] Add BJA deviations (Illustrious 18 / Fab 4, etc.)

### P3 – Lower Priority / Future Work

- [ ] Broader deviation set beyond BJA baseline  
- [ ] Support different count systems (e.g. KO, Zen, Hi-Opt II, CAC2)


## Scratchpad

- Cut card comes out and switches cut card state, so engine knows to finish current hand but shuffle on the next
- **CURRENT** EV: -1.5% (not good)