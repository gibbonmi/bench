# Classify destination ref update failures

Blocked by: none
Writes: internal/landing

## What to build

After a failed destination compare-and-swap, re-read the destination ref and
report movement only when its identity changed; preserve infrastructure failure
when the expected identity remains current.

## Acceptance

- [ ] Destination movement retains the compare-and-swap refusal.
- [ ] A held or failing ref lock with an unchanged destination retains infrastructure attribution.
