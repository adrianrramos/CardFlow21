# 🧪 Blackjack Engine Debug Checklist

## 🔴 1. Payout & Resolution (CRITICAL)

# FINAL REVIEW
- [ ] Verify TotalWagered is being tracked properly

## 🟠 4. Splits & Advanced Actions

### Splitting
- [ ] Resplit rules (RSA) enforced correctly
- [ ] resplitting Aces vs only serving 1 card to split aces

### Split Aces
- [ ] Only **one card drawn** (if rule applies)
- [ ] No further hits allowed (if rule applies)

### Blackjack after split
- [ ] 21 after split is **NOT blackjack**
- [ ] Pays **1:1**, not 3:2

### Double after split (DAS)
- [ ] Allowed after split
- [ ] Works like normal double

---

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

### Hand resolution
- [ ] Each split hand resolved independently
- [ ] Totals aggregated correctly

---

## 🟡 6. Shoe / Card Mechanics

### Deck integrity
- [ ] Cards removed after dealing
- [ ] No duplicate cards
- [ ] No missing cards

### Penetration (85%)
- [ ] Reshuffle at correct depth
- [ ] No reshuffle mid-hand

### Shuffle randomness
- [ ] Shuffle is unbiased
- [ ] No deterministic patterns

---

## 🟡 7. Strategy Execution

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

- [ ] Blackjack paying **2:1 instead of 3:2**
- [ ] Surrender returning **full bet**
- [ ] Dealer not hitting **soft 17**
- [ ] Player acts after dealer blackjack
- [ ] Split 21 treated as blackjack
- [ ] Double not increasing risk
- [ ] Insurance mispaid
- [ ] Hidden information leak
