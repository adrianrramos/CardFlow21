# 🧪 Blackjack Engine Debug Checklist


# FINAL REVIEW
- [ ] Verify TotalWagered is being tracked properly

## 🟠 5. Accounting & EV Tracking

### Bet tracking
- [ ] Includes:
  - [ ] Base bets
  - [ ] Doubles
  - [ ] Splits
  - [ ] Insurance

### Profit calculation
- [ ] Profit = net outcome per hand
- [ ] No double-counting wins
- [ ] No missing losses

## 🟡 7. Strategy Execution

- [ ] Make sure that each OUTCOME is tested from unit tests
- [ ] Hard totals correct
- [ ] Soft totals correct
- [ ] Pair splitting correct
- [ ] Surrender decisions correct

---

## 🟢 8. Sanity Tests (Quick Validation)

### Test 1
- [ ] Player 20 vs Dealer 6 behaves realistically (~80% win)

### Test 2
- [ ] Disable splits/doubles/surrender → EV ≈ **−2% to −3%**

### Test 3
- [ ] Disable insurance → EV barely changes

---

## 🚨 High-Signal Bug Checks (FASTEST TO VERIFY)

- [ ] Dealer not hitting **soft 17**
- [ ] Player acts after dealer blackjack
- [ ] Split 21 treated as blackjack
- [ ] Double not increasing risk
- [ ] Insurance mispaid
- [ ] Hidden information leak
