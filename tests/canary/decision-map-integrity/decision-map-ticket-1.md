# Data model — one canonical family, or symmetric peers?

Blocked by: none
Type: Grill

### Question

`BENCH_TIER_*` is named for the abstract tier, but it holds Codex ids. Fix the
presentation only, or restructure the data model so that no family is canonical?

### Answer

**Symmetric peers.** The tier becomes the only canonical identity, and each
harness family is a peer binding. Presentation-only was rejected, because the
prose would still call one family the tier.
