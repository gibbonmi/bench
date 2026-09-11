# Retained implementation workflow review pickup

Frozen delta: `935ca83cb65107897f51d11e118e6e71c5440817..1212e8ec029d6dacf5c4d04f48376a58e086ce21`.

## Standards

Finding count: 0. Worst issue: none.

## Spec

Finding count: 1. Worst issue: P2.

1. [P2] `auto-fix` — Preserve chunk identity across plan amendments. `.bench/BENCH.md:127` permits splitting and combining chunks but omits the approved requirement to record old-to-new chunk IDs. Add that requirement to the canonical policy and its coverage.

## Coverage

Finding count: 0. Worst issue: none.
