# Introduce the bounded projection carrier

Blocked by: none
Ownership fence: `internal/axi`
Integration surfaces: projection API→`internal/axi`; registry declarations→declare-production-axi-registry.md
Contracts: selected content, integer total/emitted/omitted, boolean truncated, completeness enum, and counting-unit enum cross caller→`internal/axi`, owner domain is declared policy, field order is fixed, and unknown completeness is explicit, asserted by PC1
Closure: PC1/content, PC1/total, PC1/emitted, PC1/omitted, PC1/truncated, PC1/completeness, PC1/unit

## What to build

projections preserve selected content and every owner-supplied size, completeness, and unit fact without inference.

## Acceptance

- [ ] [PC1] (covers CR3) projections preserve selected content and every owner-supplied size, completeness, and unit fact without inference.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| PC1/content | replace selected content with a recomputed prefix | projection test | supply content and require identity |
| PC1/total | derive total from emitted | projection counterexample test | supply unequal values and require both |
| PC1/emitted | replace emitted with content length under the wrong unit | projection counterexample test | supply multibyte content and require owner value |
| PC1/omitted | derive omitted from total minus emitted | projection counterexample test | supply owner-specific omission and require it |
| PC1/truncated | derive truncation from omitted | projection counterexample test | supply independent facts and require both |
| PC1/completeness | coerce unknown completeness to false | projection test | require explicit unknown |
| PC1/unit | normalize bytes and code points | projection test | construct both and require distinct enums |

