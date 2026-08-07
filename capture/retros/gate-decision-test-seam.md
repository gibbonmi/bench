# Retro: gate-decision-test-seam

## Outcome

Landed. Promotion commit `3cce8aa` published candidate `f7f0ea8c` (one file,
`internal/gate/check_slots_test.go`, +300/−44): the exhaustive public-document
mapping matrix moved to the read-only decision seam with the full-engine
composition, fail-closed, and evaluation controls retained. Three tickets, one
repair cycle (SPEC-OP1-001), one forced recomposition after a concurrent
capture-only commit (`ff3818b`) advanced `main` mid-build. Spec retired the same
session.

## Gate-stage timings

All runs at `GOMAXPROCS=2`. Exact-tip recovery gate after the interrupted-pending
record: conformance-only partial, ~2.5 min (recorded 04:46:24Z), every other
component inherited from content-addressed slots authored 00:06:21Z. Promote's
prospective full-inventory gate: ~9 min (verdict 04:59:20Z, tree `0e064c7d`).
Spec-retire commit gate (specs/ + ROADMAP.md, reduced scope): a few minutes,
green with one privilege capability skip. For contrast, the mistakenly-started
`bench gate --fresh` in the prior session was interrupted rather than paid.

## Ticket-versus-spec-slice and delegate performance

Ticket-sized charges performed well: the repair delegate's receipt shows focused
matrix green in 19.4s and the full `internal/gate` package in 139.9s, with three
named mutations red (`captured 2 generations want one`). The re-review of the
recomposed candidate cost three read-only mid-tier axis delegates (~2–4 min
each, parallel) — cheap because the whole composition was one test file.

## Coordinator catches

- The Spec axis did not take the repair's word: overlay mutation runs proved the
  re-added raw capture alias is red against the repaired candidate and green
  against pre-repair `4c40ca47`, confirming SPEC-OP1-001 was real and is closed.
- Standards caught 13 sites of duplicated expectation knowledge
  (projection strings restating the check-name literals); Coverage caught that
  the decision-seam matrix leaves reload coverage to a single retained
  full-engine control. All dispositions risk-accepted under the batch approval
  and parked; receipt digest `80069545…` is the veto surface.
- Checkpoint's `--evidence` flag silently requires an absolute path; the relative
  path produced only `invalid spec build receipt`.

## Agent-experience improvements

### Bench CLI

- The recompose/start refusal hardcodes `run bench gate --fresh` as advice. For
  the common causes (interrupted-pending record, stale or mismatched record) the
  cheap, correct recovery is plain `bench gate`, demonstrated this session;
  `--fresh` is only load-bearing when a reusable green verdict exists but will
  not compose. Advice should name the cause and route accordingly — parked for
  reviewer decision.
- An interrupted `bench gate` overwrites `.git/bench-last-gate` with a pending
  record, destroying the prior verdict that diagnosis needs. A preserved
  last-good record (or naming the overwritten verdict in the pending record)
  would have made the original refusal provable instead of permanently lost.
- Receipt-path validation should say "receipt path must be absolute" instead of
  the generic invalid-receipt refusal.

### Skills

- The `$bench-debug` handoff prompt from the interrupted session carried enough
  pinned state (run, candidate, assignment, receipt path, refusal text) to
  resume cold with zero re-derivation. Keep that shape.

### Process

- A concurrent capture-only commit on `main` during an active spec build forces
  recomposition, which discards the review and re-prices a full three-axis pass
  plus the prospective gate. Cheap here; on a wide candidate it would not be.
  Worth deciding whether capture commits should pause while a build is between
  review and promote.
