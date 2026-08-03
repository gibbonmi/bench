# Retro — ft181-precondition-residuals

## Outcome

Promoted terminal at candidate a25e2803, published as e11df18 on main. All four
stories landed: restart reads the live green marker (recorded base retained as
history evidence only); decayed checkouts classify as liveness by shape with
identity faults fatal everywhere; the blanket prepared-abandon exemption is
deleted with the enumerated identity table pinning it; empty non-terminal runs
fast-forward on checkpoint and start behind a closed operation list. The
27-row coverage map is fully dispositioned (13 review findings: 10 resolved, 3
endorsed). Three review rounds ran: round 1 found 8 blocking findings (4 real
production gaps among them), round 2 confirmed closure and found 1 fresh
blocker in the repair itself, round 3 was waived by reviewer directive after
the prescribed fix was independently probed.

## Gate-stage timings

Promotion-stage timings are not retained in the terminal record (known gap,
already on the roadmap via the check-level-conformance-scoping retro). Session
wall-clock observations: ticket-landing gate (reduced scope, conformance only)
~1 min; full-package `go test ./internal/specbuild` ~37 s and
`./internal/worktree` ~21 s per delegate verification; the two capture/prose
commits with `.agents/` inputs (contract + canary components fresh) ~4–6 min
each; promote's composed gate ~5 min. Two sol review passes: ~313k and ~250k
tokens, several minutes each.

## Ticket-versus-spec-slice and delegate performance

Four planned tickets plus three repair assignments, all charged as tickets with
fences and row lists — no whole-spec-slice charge was used. Every build
delegate returned first-pass green on its focused checks; no delegate needed a
re-charge. The repair delegates (two opus, one sonnet) also landed first-pass;
the sonnet repair (special .git guard) confirmed the per-ticket cheap-tier
re-route works when the finding is precise. Charges carrying the pre-digested
fixture inventory (helper names, file:line prior art) are the plausible cause
of the first-pass rate; that input is now codified in craft-delegate.

## Coordinator catches

- Vacuous probe detected twice: a perl pattern that didn't match (probe "ran"
  against unmutated code) and a probe at an unpinned branch (ShapeUnknown) —
  both caught by checking the mutation actually applied and re-probing at a
  pinned site.
- A finished delegate's stale re-notification claimed its integrated work was
  lost; verified against the candidate (ClassifyPathShape present) instead of
  acting on the claim.
- The review-receipt disposition vocabulary ("accepted" = open defect) caught
  only at promote refusal; cleared via recomposition onto the prose-fix commit.
- Missed catches, recorded in capture/learnings.md: two skipped craft-tickets
  steps at derivation (contract discovery, prefactor) that the review then
  found as P1/C1/S1; findings auto-repaired without triage; delegate caveat
  about the fake owner relayed but not acted on.

## Agent-experience improvements

### Bench CLI

- Receipt-skeleton plumbing helper (parked in IDEAS): coordinators currently
  reverse-engineer checkpoint/review receipt schemas (~30k tokens) and can
  mis-use the disposition vocabulary; a skeleton emitter with enumerated
  dispositions removes both.
- Fixture-and-seam inventory generator (parked): `bench outline` extension so
  the charge input craft-delegate now requires costs nothing to assemble.
- Terminal record retains no promotion-stage timings (existing roadmap row);
  this retro's timing section is session memory, not evidence.

### Skills

- Landed this session (reviewer-directed): craft-tickets `Contracts:` field;
  craft-spec composition degenerate + existing-control edge rule; craft-line
  per-ticket ceiling; craft-delegate fixture inventory + focused checks;
  implement-spec research scaling. Expected effect: the P1-class composition
  gap becomes visible at derivation, the P2-class stale-control claim becomes
  a row, and tier/research over-provisioning stops being the default.

### Process

- Review-round tail: the terminal repair-pass bound was applied one round too
  late; a prescribed-fix delta with an independent probe should not re-open a
  full review round. Candidate rule parked in learnings (triage threshold).
- The Claude worktree hook's one-assignment-per-session limit did not bite:
  `bench spec build assign` cut every worktree; the FT176 manual-create
  workaround was never needed in this run.
