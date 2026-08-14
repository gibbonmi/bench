# Lighten the light path to write-ticket-then-implement-inline

Blocked by: none
Writes: .bench/BENCH.md, .agents/skills/bench-craft-tdd/SKILL.md, .agents/skills/bench-craft-tickets/SKILL.md, internal/anchors/registry_data.go, internal/conformance/fixture_bite_test.go, tests/canary/workflow-guidance-anchors/

## What to build

An under-threshold change follows the new route: write the one ticket file
(`craft-tickets` owns the template), then implement it inline in the same
session — no breakdown-approval pause, no write-delegate, no worktree — and
gate and commit on green. `.bench/BENCH.md`'s right-size table keeps the
`Right-size the process` marker and both observables (`one independently-green
ticket`, `crosses no declared seam`) verbatim; only the route cell changes,
and the old route text is retired as a Forbid+Require pair; the new route
cell must avoid BENCH.md's whole-file Forbid on the bare substring
`thorough`. `craft-tdd`'s
light-path bullet replaces the stop-for-seam-confirmation with: the ticket
file names the test seam and the reviewer vetoes post-hoc — the retired stop
trips a Forbid, rewriting the `craft-tdd-light-path-seam-gate` fixture.
`craft-tickets` gains the light-path carve-out — the right-size table is the
ticket's standing approval and it implements inline — scoping its "only route
onto the frontier" and "one fresh write-delegate charge" rules to spec-backed
builds; the carve-out displaces lines (the file sits at its 100-line budget),
keeps `ticket-skill-contract-anchor` green (it pins the
one-write-delegate-charge row; the reviewer-approved-breakdown row has no
canary fixture — its only bite is the cadence case) and craft-tickets'
cadence-pinned substrings byte-identical (else the byte-exact strings in
`internal/conformance/fixture_bite_test.go` update in this ticket), and gets
its own new fixture.
`ticket-light-path-anchor` is retained; a new light-path-inline fixture pins
the new route. The three edited files are grouped by reviewer-chosen
bundling: no thinner cut strands a gate red, but each intermediate commit
would ship a BENCH.md route contradicting craft-tdd's stop or
craft-tickets' contract — flagged for the reviewer, who may split. Shares
`.bench/BENCH.md` and `craft-tickets` (with move-slicing-into-write-spec.md,
both displacing lines at its 100-line limit), plus `registry_data.go` and
the fixtures tree with siblings — those paths land serially across the
whole spec.

## Acceptance

- [ ] BENCH.md's light-path route reads write-ticket-then-implement-inline with
      no approval pause, delegate, or worktree; both observables and the
      right-size marker survive; the old route text trips a Forbid
      (covers WF1)
- [ ] craft-tdd's light-path bullet carries the ticket-names-the-seam post-hoc
      veto with the retired stop tripping a Forbid, and the rewritten fixture
      bites both halves (covers WF2)
- [ ] craft-tickets scopes breakdown approval and the one-delegate charge to
      spec-backed builds and its carve-out fixture bites both halves
      (covers WF16)
- [ ] with everything landed, BENCH.md is ≤ 180 lines, craft-tdd ≤ 120, and
      craft-tickets ≤ 100 (covers WF13)
