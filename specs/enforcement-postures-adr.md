# enforcement-postures-adr — the accepted enforcement postures are recorded as decided state

Status: staged

Source: closed decision map `decisions/enforcement-postures-adr.md` (FT31); its
Handoff is the pre-agreed seam set, taken as-is. Depends on FT27
(`decisions/enforcement-verification.md`): the fail-open entry records the
residual rims FT27 leaves behind, so this spec is authored **after FT27 lands** —
a dependency on the shipped wording, not an open uncertainty.

## Problem

Six deliberate enforcement postures were accepted during an assessment drain, but
they are citable nowhere. They live in that assessment's disposition table, which
is now git history — so every future "is this a hole?" question re-derives the
answer from scratch. Worse, a cold session that mistakes an accepted risk for an
unaddressed bug can "fix" a posture the reviewer deliberately chose, or re-litigate
an acceptance that was already settled. There is no current-state record a session
can read to learn that these gaps are decisions, not oversights.

## Solution

One ADR at the next sequential number in `docs/adr/`, status accepted, records the
six postures as current decided state. Each entry states three things: the posture,
why it is accepted, and what evidence would reopen it — framed as **accepted risk,
never as mitigated**, because a record that promises enforcement it does not have is
worse than silence. Placing the postures in per-package doc comments (not citable as
decisions) or in the project profile (which describes the repo, not platform
decisions) was considered and rejected; the ADR records those rejections so nobody
re-litigates the placement.

## User stories

1. As a teammate who just walked in, I want the ADR to exist at the next number in
   `docs/adr/` with status accepted, every entry framed as accepted risk with a
   named reopen condition, no file paths and no code snippets, and the rejected
   placements (per-package comments, the profile) recorded — so the record reads as
   the current decided state, addressed to someone with no memory of the assessment
   that produced it.
   Line: claude-fable-5 / high. The record's framing — accepted risk, never
   mitigated — is the load-bearing judgment the profile's doc-authoring leverage
   override exists to buy the top tier.

2. As a session asking whether raw git needs fencing, I want the ADR to record that
   interactive commit-on-red remains possible by design — the gate wrapper and the
   shift loop are the enforced paths, so an interactive session can still commit on a
   red gate — accepted because fencing every raw git invocation is out of proportion
   and invariant 4 governs the behavior on the paths that matter, reopened by evidence
   that interactive commit-on-red happens often enough to warrant fencing raw git.
   Line: claude-fable-5 / high. Whether raw git stays deliberately unfenced is a
   subtle posture a future reader must not misread as a bug, so the override's top
   tier authors it.

3. As a session weighing a done-claim made outside a shift, I want the ADR to record
   that non-shift done-claims are honor-system — invariant 1 governs the behavior, but
   hooks enforce only what they can observe, so a claim made off the shift path is not
   gate-verified by a hook — accepted because a hook can only enforce what it sees and
   the invariant covers the rest by discipline, reopened by evidence of unverified
   done-claims landing off the shift path.
   Line: claude-fable-5 / high. The honor-system boundary is easy to overstate as
   enforced, and getting that line exactly right is doc-authoring leverage the
   override routes to the top tier.

4. As a session reading a declared line, I want the ADR to record that "declare the
   line" has no enforcement surface for effort — the model-tier binding (membership)
   is enforced by the agent-line guard, but the effort level is declaration discipline
   with no observable surface a hook can check — accepted because effort is not
   observable to a hook the way a model id is and the declaration invariant governs
   it, reopened by a mechanism that makes effort enforceable or by evidence that
   effort-declaration drift causes real cost.
   Line: claude-fable-5 / high. Effort-declaration having no enforcement surface is a
   precise claim that a careless phrasing turns into a false promise, so it earns the
   override's top tier.

5. As a session auditing the agent-line guard, I want the ADR to record the fail-open
   guard rims that remain after FT27 tightens the routed-repo missing-model branch —
   unrouted repos (no binding to enforce) and malformed or absent envelopes (degraded
   hook input) stay fail-open by design — accepted because an unrouted repo has no
   binding to enforce and the residual degraded branches are not the silent-escalation
   attack path the guard exists for, reopened by evidence that a residual rim is being
   exploited as a silent-escalation path.
   Line: claude-fable-5 / high. The residual rims must be stated as accepted risk
   after FT27's tightening without implying the rims are closed, a distinction the
   override buys with the top tier.

6. As a session worried a stale verdict could lie green, I want the ADR to record
   recompute-always gating — the gate always recomputes, and the verdict cache is
   advisory for the dashboard, never a short-circuit, so a stale cache cannot lie
   green because the gate never trusts it to decide "done" — accepted because
   correctness of the oracle outranks speed and a cache that could short-circuit the
   gate would be the worst class of bug in a kit whose premise is "the gate is the
   oracle," reopened by a gate-runtime cost high enough that a verified,
   non-short-circuiting cache becomes worth the risk.
   Line: claude-fable-5 / high. Recompute-always is a correctness posture whose "a
   stale cache can't lie green" reasoning must be exact, which the doc-authoring
   leverage override routes to the top tier.

7. As a session reading the canary, I want the ADR to record family-level canary
   granularity — one needle per check family, with per-check needles minted only when
   anchor rot is actually observed rather than one needle per check up front — accepted
   because per-check needles up front over-fit the canary and cost maintenance while
   family granularity still catches the always-pass rot the tripwire exists for,
   reopened by observed anchor rot inside a family that the family-level needle failed
   to catch.
   Line: claude-fable-5 / high. Family-level canary granularity is a deliberate
   under-fitting a reader could mistake for a gap, so the override authors the entry at
   the top tier.

## Implementation decisions

- **One new file in `docs/adr/`, next sequential number.** Scan the folder and
  increment — 0002 today, since 0001 is the only existing record. The file carries
  status accepted. `craft-adr` governs the form.

- **`craft-adr` form is binding: current decided state, no history.** Record that
  each posture is accepted and why; do not narrate the assessment path that led here
  (git holds that). No file paths, no code snippets — each posture is named by its
  role and mechanism (the gate wrapper, the shift loop, the agent-line guard, the
  verdict cache, the canary), not by path. The record reads to a session that has
  never seen the assessment.

- **Every entry is framed as accepted risk with a reopen condition — never
  "mitigated."** This is the load-bearing content decision (map Handoff item 9): an
  ADR that promises enforcement it does not have is worse than silence. Each of the
  six entries states the posture, why it is accepted, and what evidence would reopen
  it — the reopen condition is what keeps the record honest as an accepted risk rather
  than a claim of coverage.

- **The fail-open entry records what FT27 leaves behind, and re-decides nothing.**
  FT27 tightens the one branch that is also the attack path (routed repo, Agent call,
  no model → deny); the ADR entry names the rims that remain fail-open by design after
  that tightening, without implying they are now closed. Because the entry's exact
  wording depends on FT27's shipped posture, this spec is built after FT27 merges — a
  sequencing dependency, not an uncertainty flag, so no top-tier falsification pass is
  triggered on that ground.

- **Rejected placements are recorded in the ADR so they are not re-litigated.**
  Per-package doc comments were rejected (a comment is not a citable decision); the
  project profile was rejected (profiles describe the repo, not platform decisions).
  A short considered-options note in the ADR carries these, per `craft-adr`'s rule
  for rejected alternatives worth remembering.

- **No behavioral change and no gate-check addition.** This spec adds prose only. It
  wires no new conformance check; the file is graded by the existing docs-conformance
  scan (stale references) and by reviewer cold-read.

## Testing decisions

- **What a good test is here:** almost none to write — the record's correctness is
  semantic. The gate guards regressions structurally (the docs-conformance
  stale-reference sweep runs over the new ADR and reds on any stale CLI reference),
  and the reviewer enforces the content at cold-read. The rows below classify honestly
  rather than dressing prose review as TDD — the docs-drift spec set this precedent.
- **The cheapest wrong implementation passes every gate check.** Recording three of
  the six postures, framing an entry as "mitigated," or stating a wrong reopen
  condition all leave the gate green — which is precisely why these rows are
  review-verified and say so openly. The six postures are enumerated as stories 2–7 so
  "record the postures" cannot collapse to the cheapest subset.
- **Seams:** the docs-conformance stale-reference sweep, and reviewer cold-read.
- **Gate:** the project gate, `bench gate`.

### Seam diagram

    trigger: `bench gate` (docs-conformance stale-reference sweep) + reviewer cold-read
        │
        ▼
    new ADR in docs/adr/  ──▶  [ docs-conformance stale-reference sweep ]  ──▶  red on any stale
    (six posture entries, ──▶  [   (structural: sees the new file)      ]        CLI reference in
     accepted, reopen     ──▶  [                                        ]        the added prose
     conditions, no paths)──▶  [                                        ]
                      ◀ tests attach here: run `bench gate` after authoring — the sweep
                        catches stale references structurally. Posture content, the
                        accepted-risk framing, and the reopen conditions attach at
                        reviewer cold-read — stated openly, not disguised as TDD.

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | ADR exists at the next number, status accepted, framed accepted-risk throughout, no file paths, rejected placements recorded | stale-reference sweep + reviewer cold-read | the stale-reference sweep reds on any stale CLI reference in the new file (structural half); status, numbering, no-paths, and the accepted-risk framing are not TDD-able — reviewer cold-read | a "mitigated"-framed, history-shaped, or path-bearing record passes every structural check, so the reviewer is the enforcement — stated openly |
| 2 | interactive commit-on-red recorded as accepted with its reopen condition | reviewer cold-read | not TDD-able — posture content is semantic; no gate check reads posture prose | a record omitting this posture or stating a wrong reopen condition stays green; the reviewer confirms the entry is present and framed as accepted risk |
| 3 | non-shift done-claims recorded as honor-system with its reopen condition | reviewer cold-read | not TDD-able — semantic prose | same honest posture as story 2 |
| 4 | declare-line effort recorded as having no enforcement surface, with its reopen condition | reviewer cold-read | not TDD-able — semantic prose | same |
| 5 | residual fail-open rims after FT27 recorded as accepted, with their reopen condition | reviewer cold-read | not TDD-able — semantic prose; the FT27-dependent wording is verified against the shipped posture at cold-read | a record that implies the rims are closed, or that predates FT27's wording, passes every check; the reviewer verifies it matches the shipped posture |
| 6 | recompute-always gating recorded as accepted, with its reopen condition | reviewer cold-read | not TDD-able — semantic prose | same honest posture as story 2 |
| 7 | family-level canary granularity recorded as accepted, with its reopen condition | reviewer cold-read | not TDD-able — semantic prose | same |

### Edge inventory

This spec touches no code path; the runtime edge classes (error path, empty/absent
input, boundary values, malformed input, interrupted/partial state, re-run
idempotency, hostile environment) are **Won't handle — no runtime surface** (the
docs-drift precedent). The documentation-specific edges walked:

- **a posture recorded as "mitigated" rather than "accepted risk"** → covered by the
  story-1 row; this is the load-bearing domain watch-out and attaches at reviewer
  cold-read.
- **fewer than six postures recorded** → covered: stories 2–7 enumerate all six, and
  the reviewer confirms every entry is present.
- **a wrong or missing reopen condition on an entry** → covered by the per-posture
  rows (2–7); reviewer cold-read.
- **a file-path or code-snippet token in the ADR body** → the stale-reference sweep
  reds on stale CLI references structurally; any other path token is `craft-adr` /
  invariant 3 and is reviewer-checked.
- **the fail-open entry authored before FT27's wording is final** → not a runtime edge
  but a build-sequencing constraint: the build runs after FT27 lands, and the reviewer
  verifies the entry against the shipped posture (story-5 row).
- **ADR number collision with a concurrently added record** → the author scans
  `docs/adr/` and increments at write time; low risk under single-writer, reviewer
  confirms the number.

## Out of scope

- **A conformance family for ADR-posture currency** — detecting when a posture this
  ADR calls "accepted" has since been closed by shipped enforcement (for example, a
  future guard that fences raw git, retiring the commit-on-red posture). Semantic
  currency detection is a separate capability and probably a review-time duty, not a
  gate check. No estimate — needs a shape decision first.
- **Auditing the rest of the kit for other decided-but-unrecorded postures** — a
  separate sweep with its own findings; this ADR records the six the assessment
  verified, and a future accepted posture is its own future decision appended when it
  is made. Estimate: ~1 session of read-only sweep, then its own drain.
