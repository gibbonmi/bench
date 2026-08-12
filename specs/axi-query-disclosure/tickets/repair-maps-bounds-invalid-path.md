# Repair maps bounds-invalid path

Blocked by: none
Writes: `internal/maps/` (tests only)

## What to build

Pin the carried path on the bounds-classified invalid branch (review finding
R10): `maps.go:58` passes `candidate.Path` into the invalid-map disclosure,
but every existing test exercises only the `ValidateDecisionMap` diagnostic
branch (`maps.go:63`), so mutating line 58 to pass `mapName` stays green while
the disclosed repair path is wrong. Add a public-command case for a
bounds-invalid map file.

## Acceptance

- [ ] [MB1] (covers QD1) a public `maps` command test with a bounds-invalid
  map file (empty file, exercising the `bounds.Classify` non-parsed branch)
  asserts the disclosed action carries the full `decisions/<name>.md` path —
  the `bench maps --template` invocation plus the named path — and the
  bounds-derived diagnostic with exit 1.
- [ ] [MB2] (covers QD1) the mutation `candidate.Path` → `mapName` at
  `maps.go:58` is demonstrated red against the new test, then restored green,
  recorded in the ticket evidence.
