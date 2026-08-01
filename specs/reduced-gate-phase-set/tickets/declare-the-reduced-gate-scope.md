# Declare the reduced gate scope

Blocked by: none

Ownership fence: `internal/gate/scope.go`, `internal/gate/scope_test.go`
Assumptions: `internal/status` already imports `internal/gate`, so the declaration can live in `internal/gate` and be consumed by every current holder of a private capture-path copy without a new package or an import cycle

## What to build

One declaration holding both halves of a single claim — *these paths are invisible
to these phases*. The path allowlist is the `capture/` and `specs/`
directories, `ROADMAP.md`, and
`.bench-notes.md`. The excludable phase set is `gofmt`, `vet`, `test`, `race`,
`contract`, `shellcheck`, and `canary`; the included set is `conformance` and
`conformance-suite`; `build` is unconditional in both modes because it produces the
binary the other phases exec.

With it, the predicate every later ticket routes through: is a changeset confined to
the allowlist? Membership is byte-exact against the two file entries and
containment under a declared directory. It is not a prefix test —
`ROADMAP.md.bak` and `capture-old/x.md` are not members — and it does not enumerate
the files inside `capture/` or `specs/`, because membership following location is
what the co-location migration bought and re-enumerating would hand the
exhaustiveness burden back. `specs/` is a declared directory on the reviewer's
2026-08-01 ruling: it is formatted documents whose graders are the included phases.

## Acceptance

- [ ] [R02] A newly added file under `capture/` is a member without its own declaration entry, and the same holds for `specs/`.
- [ ] [R03] A near-miss sibling of a declared file, including `ROADMAP.md.bak`, is not confined.
- [ ] [R04] Any descendant of a declared directory is a member, including `capture/retros/<slug>.md` and `specs/<slug>/tickets/<t>.md`, while `capture-old/x.md` and `specs-old/x.md` are not.

## Ruling on a spec ambiguity

The spec's R04 row reads "a nested or sibling-prefixed path is not" a member, which
contradicts its own R02 row requiring every co-located capture surface to be a
member — `capture/retros/<slug>.md` is nested and must be covered, and `specs/`
is deeper still. Read "nested" as the sibling-prefix case the row's red signal
actually names (`capture-old/x.md`). Membership is descendant containment under a
declared directory, with the path boundary respected so a sibling prefix never
matches. Flagged for reviewer veto.
- [ ] [R13] The excludable set contains the `contract` phase, so the declaration cannot be satisfied by excluding one trivial phase.
