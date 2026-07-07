# Enforcement Postures ADR (FT31)

## #1: Which accepted postures does the record cover, and where does it live?

Blocked by: —
Type: Grill

### Question
Several deliberate enforcement postures are defensible but citable nowhere —
they live in a drained assessment's disposition table, which is now history.
Every future "is this a hole?" question re-derives the answer.

### Answer
One ADR (next number in docs/adr/) records the accepted postures as current
decided state: interactive commit-on-red remains possible (raw git isn't
fully fenced — the gate wrapper and shift loop are the enforced paths);
non-shift done-claims are honor-system (invariant 1 governs behavior, hooks
enforce only what they can see); "declare the line" has no enforcement
surface for effort (declaration discipline, membership enforced); the
remaining fail-open guard rims after [FT27] tightens the routed-repo
missing-model branch (unrouted repos, malformed envelopes); recompute-always
gating (the verdict cache is advisory for the dashboard, never a
short-circuit — a stale cache can't lie green); family-level canary
granularity (one needle per family, per-check needles only on observed anchor
rot). Each entry: the posture, why accepted, and what evidence would reopen
it. Per invariant 3: no file paths, no code snippets. Rejected: scattering
the postures into per-package doc comments (not citable as decisions);
recording in the profile (profiles describe the repo, not platform
decisions).

## Handoff

1. **Module boundaries.** One new file in docs/adr/; craft-adr skill governs
   form; no code.
2. **Contracts.** ADR states six postures as accepted decisions with reopen
   conditions; status "accepted".
3. **Deep vs thin.** Thin — the decisions exist; this writes them down.
4. **Black-box assertables.** Conformance docs scan sees the file (stale CLI
   references would red); nothing behavioral.
5. **Gate attachment.** Docs conformance only; content is review-verified.
6. **Hostile-input owners.** n/a — prose.
7. **Uncertainty flags.** The FT27 outcome changes the fail-open entry's
   wording — write this ADR after FT27 lands (dependency, not uncertainty).
8. **Rejected alternatives.** Per-package comments; profile notes; leaving
   the postures in assessment history.
9. **Domain watch-outs.** An ADR that promises enforcement it doesn't have
   would be worse than silence — each entry must say "accepted risk", not
   "mitigated".

Dependency order: after FT27 (the fail-open entry records the posture FT27
leaves behind).
