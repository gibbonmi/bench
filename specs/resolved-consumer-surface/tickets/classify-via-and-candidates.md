# Classify via values and answer ambiguous names

Blocked by: add-consumers-command-surface.md
Writes: internal/consumers/

## What to build

The `via` classifier and the candidates path. A call-position use emits
`via=call` through the static callee test, and a value use stays
`via=reference`. A queried interface emits `via=implements` rows for its
satisfying fixture types. A bare name with several matches emits
`consumers_candidates[N]{qualified,file,line,kind}` at exit 0 and zero
consumer rows. The response ends with one literal re-query action per
candidate row through the `internal/axi` owner, because every argument is
known.

## Acceptance

- [ ] CS2: a call-position use emits `via=call` and a value use emits
      `via=reference`.
- [ ] CS3: an interface query emits `via=implements` rows for its satisfying
      fixture types.
- [ ] CS5: a bare name with two fixture matches emits the candidates table
      and zero consumer rows.
- [ ] CS17: each candidates row carries the exact qualified re-query
      spelling.
- [ ] CS18: the candidates response ends with one literal re-query action
      per candidate row.
