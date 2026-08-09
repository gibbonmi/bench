# Migrate worktree aggregates

Blocked by: none
Ownership fence: `internal/worktree`
Integration surfaces: shared aggregate/empty carriers→implemented prerequisite `axi-carriers-and-registry`; final contraction→contract-aggregate-empty-routes.md
Contracts: worktree inventory integers/booleans and fingerprint authority cross `internal/worktree`→shared aggregate; domain is complete owner inventory; order is current renderer order; absence is empty safe inventory, asserted by WA1
Closure: WA1/count, WA1/bytes, WA1/shown, WA1/truncated, WA1/fingerprint-authority, WA1/order, WA1/route

## What to build

Migrate worktree aggregates through the shared carriers without changing owner semantics or public bytes.

## Acceptance

- [ ] [WA1] (covers AE4) migrate worktree aggregates preserve count, bytes, shown, truncated, fingerprint-authority, order, route.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| WA1/count | derive count from shown | the independent producer route or compatibility test | apply the subject mutation, run the real public producer fixture under its named bound, and require the specific red |
| WA1/bytes | derive bytes from visible | the independent producer route or compatibility test | apply the subject mutation, run the real public producer fixture under its named bound, and require the specific red |
| WA1/shown | report total as shown | the independent worktree inventory test | cross the visible cap and require owner shown count |
| WA1/truncated | derive truncation from shown alone | the independent worktree inventory test | exercise authority ceiling and display cap separately |
| WA1/fingerprint-authority | move fingerprint derivation into the carrier | the independent ownership test | mutate one authority input and require fingerprint mismatch |
| WA1/order | reorder fields | the independent exact renderer test | render inventory facts and require current order |
| WA1/route | bypass shared aggregate | the independent route test | invoke worktree list and require the missing route marker |
