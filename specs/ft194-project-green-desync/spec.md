# ft194-project-green-desync

Status: staged

Decision source: `ROADMAP.md` FT194 (HIGH, reproduced defect), verdicted in the reviewed 2026-08-04 `/bench-what-next` drain (commit 40be04f).

## Problem

When an empty spec-build run fast-forwards its recorded base onto a moved
branch tip, the project-green marker stays at the prior base. Every later
marker consumer assumes the marker still sits at the run's recorded base, so
the sanctioned lifecycle wedges itself: promotion's recomposition refuses with
`project-green marker conflicts with another tip`, and a direct promotion of a
fast-forwarded run advances the working branch, then fails the marker
compare-and-swap — leaving a non-terminal run on an already-moved branch that
re-entry can never finish. Checkpoint, assign, and review refuse toward the
promote that refuses; a fresh gate cannot advance the marker; abandon discards
the reviewed composition. The transition owns no repair path.

## Solution

The lifecycle recognizes a lagging marker instead of calling it a conflict. A
marker that is an ancestor of both the expected lineage (the run base) and the
destination tip is the same marker the run started from, merely behind; every
marker advance accepts it and swaps from its actual position. Divergent
markers, markers ahead of the lineage, and unreadable markers keep refusing.
The sanctioned reproduction — start an empty run at a project-green tip,
advance the branch by an ordinary gated commit, let the run fast-forward, then
promote — completes to a terminal run with branch and marker at the promotion
commit.

## User stories

1. As the reviewer driving a spec build, after an ordinary gated commit moved
   the branch and the empty run fast-forwarded onto it, I complete
   assign → checkpoint → integrate → review → promote and the run lands
   terminal, with the working branch and the project-green marker both at the
   promotion commit. Line: opus / medium. Authorization-owner state machine —
   the profile's cached gate-logic routing holds oracle-adjacent code at mid.
2. As the reviewer whose fast-forwarded run gained checkpointed work before
   the tip moved again, `bench spec build promote` recomposes instead of
   refusing with `project-green marker conflicts with another tip`, and the
   retried promote completes. Line: opus / medium. Same authorization state
   machine, same cached routing.
3. As the reviewer recovering a promotion that died between the branch advance
   and the marker advance, re-entering promote completes the publication even
   though the marker lags the run base. Line: opus / medium. Crash-recovery
   half of the same publication policy.
4. As the reviewer trusting the marker's protective half, a marker that is
   *not* a recognized ancestor of the expected lineage — divergent, ahead of
   the lineage, or with undecidable ancestry — still refuses everywhere,
   including at promotion publication, mutating no refs and no run state.
   Line: opus / medium. These rows carry the degenerate-killing weight on the
   acceptance predicate, and the cached gate-logic routing holds oracle code
   at mid.

## Implementation decisions

- **Recognition over republication (the decision source's second option).**
  The fix is the acceptance rule in the marker owner, not a marker move at
  fast-forward time. The empty-run fast-forward stays evidence-free and
  marker-free: it is documented as running no gate precisely because promote
  is the lifecycle's only green transition, and moving the marker there would
  publish a green position with no composed-green validation and would race
  sibling runs sharing the branch marker. Rejecting the source's first option
  (advance base and marker together) is a decision of this spec.
- **The recognition rule.** An existing marker that differs from the caller's
  expected lineage is accepted only when it is an ancestor of both the
  expected lineage and the destination commit. The compare-and-swap then moves
  the marker from its actual read position, never from the assumed one — the
  swap remains the concurrency guard. A marker equal to the destination stays
  the idempotent early accept. Absent-marker and empty-expectation behavior is
  unchanged. A marker strictly between lineage and destination, a divergent
  marker, and a marker whose ancestry cannot be decided (the ancestry probe
  errors) all refuse. The recognized case keeps the existing guard that the
  caller's expectation is itself an ancestor of the destination.
- **One source for marker policy.** The recognition-and-advance rule lives in
  the gate-authorization owner. `Bootstrap` applies it at its conflict branch,
  and the promotion's two publication sites — the direct green-path publish
  and the crash-recovery publish — stop assuming the marker sits at the run
  base and route their marker advance through the `GateOwner` seam via one
  focused marker-advance operation implemented by the same owner. The
  lifecycle package never re-derives marker ancestry.
- **Publication does not re-run evidence validation.** The marker-advance
  operation validates lineage and swaps; it does not re-check composed green.
  Both publication sites run only after promotion's own prospective evidence
  is proven (the gate's green outcome on the direct path, retained-evidence
  validation on recovery), and coupling publication to the working-tree
  evidence cache's freshness window would add a refusal class promotion
  already guards against.
- No state-schema change: the run record keeps its current fields; the fix is
  entirely in the transition and its acceptance rule.

## Testing decisions

- A good test drives the public lifecycle — `Start`, `Assign`, `Checkpoint`,
  `Integrate`, `Review`, `Promote` — over a real Git fixture repository and
  observes refs (`refs/bench/green/<branch>`, the working branch, the
  candidate) and the run's public status, never internals.
- **Seam 1: the authorization owner's unit surface.** Prior art:
  `greenBootstrapRepo`-style fixtures that plant markers and assert
  advance-or-refusal plus the final ref. TDD here — the recognition rule's
  accept and refuse cases each start red.
- **Seam 2: the spec-build lifecycle package surface.** Prior art: the
  existing fast-forward, recomposition, and fault-injection lifecycle tests.
  TDD here for the wedge reproductions.
- **Rows graded on recognition run the real authorization owner.** The
  package's marker-writing fakes mirror only the compare-and-swap, so a
  fake-driven row could go green while the real owner still refuses — the
  composition degenerate. Fakes remain for rows not about marker recognition.
  Honest cost: prior art exists for a real-owner adapter and for a fixture
  with real gate evidence, but no existing test drives the full promotion
  against the real owner — the recognition rows need one new fixture
  combining the real-gate repository with the assignment/checkpoint/review
  machinery (today built on a gate-less fixture). Priced into the build:
  ~3 edits, ~1 gate run of the total below.
- **No new binary-level runtime contract.** The command adapter is a one-line
  delegation to the authorization owner, already exercised by the existing
  runtime lifecycle contracts; the composition this defect lives in is graded
  in-process against the real owner. Veto surface: adding a full wedge replay
  through the built binary prices at ~3 edits, ~3 gate runs.
- **Gate seam:** the dev gate's test, race, and conformance-suite phases
  observe both packages through the module test closure. The feature is fully
  gate-observable.
- **Quantifier enumerated — every site that writes or expects the marker
  position:** (1) fresh start and (2) start resume, both expecting the live
  marker (behavior unchanged, already covered); (3) promotion recomposition's
  bootstrap, expecting the run base; (4) the direct publication after the
  prospective gate; (5) the crash-recovery publication. No sixth site writes
  the marker. One further site *reads* the marker position without writing
  it: red-disposition attribution treats the subject as inherited-green only
  when the marker sits at the tip, so a fast-forwarded run's lagging marker
  degrades a candidate-caused prospective red's attribution to `inherited`.
  That degradation is pre-existing and this spec leaves it (edge inventory
  and out of scope carry the disposition).

### Seam diagram

    trigger: spec-build start, promote recomposition, promotion publication
        │
        ▼
    branch, destination commit,      [ gate-authorization marker owner:        refs/bench/green/<branch>
    lineage expectation         ──▶    recognize the existing marker,     ──▶  advanced to the destination,
    (run base, or live marker)         compare-and-swap from its actual        or a structured refusal
                                       position ]
                      ◀ tests attach here: fixtures plant the marker (lagging, divergent,
                        between, unreadable, absent) and assert the advance or the refusal
                        plus the final ref position

    trigger: bench spec build promote (direct green path; crash-recovery re-entry)
        │
        ▼
    reviewed candidate +             [ promote: prospective gate →             branch and marker at the
    prospective green evidence  ──▶    branch compare-and-swap →          ──▶  promotion commit; run terminal
                                       marker advance via the gate
                                       owner's seam ]
                      ◀ tests attach here: lifecycle fixtures with real gate evidence drive the
                        public operations and read refs plus run status, including after a
                        fresh-service reload and after injected crashes

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | The sanctioned reproduction (start empty at a green tip, ordinary gated commit, fast-forward, assign→checkpoint→integrate→review, promote at the unmoved tip) completes terminal | lifecycle surface driving the real authorization owner | to be observed red at build start — warranted by code reading (the publication compare-and-swap assumes the marker at the run base and fails after the branch has already advanced); the 2026-08-04 incident reproduced story 2's refusal, not this site | A fix that softens only bootstrap leaves the direct publication site assuming the marker at the run base, and this row still reds |
| 1 | After that promote, the marker, the working branch, and the run's terminal state agree on the promotion commit, and a fresh service reloads the run as terminal | lifecycle surface driving the real authorization owner, fresh-service reload | to be observed red at build start (same fixture, assertions on final refs and reloaded status) | Completing promote while leaving the marker behind would re-create the desynchronization for the next run; the reload half catches state that only lives in the first process |
| 2 | A lagging-marker run with checkpointed work whose tip moved again recomposes at promote — no `project-green marker conflicts with another tip` — and the retried promote completes | lifecycle surface driving the real authorization owner | reproduced 2026-08-04: this is the incident's exact refusal string | This is the bootstrap conflict branch itself; an always-refuse (current) implementation cannot pass, and the completion assertion rejects a recognition that advances the marker but corrupts the recomposed run |
| 3 | A promotion crashed between the branch advance and the marker advance (injected fault) completes on re-entry with the marker lagging the run base | lifecycle surface driving the real authorization owner, fault injection at the existing promote fault points | to be observed red at build start: today re-entry fails the recovery publication's compare-and-swap forever | Fixing only the direct path leaves the recovery publication assuming the run base; this row reds on that partial fix — and a fake-driven version would pass a naive advance the real owner refuses |
| 4 | A marker divergent from a non-empty expectation (sibling commit) refuses at bootstrap, marker and refs untouched | authorization unit surface | to be observed red against the degenerate always-recognize mutation of the acceptance predicate | The conflict branch is the line being softened; without this row, replacing the ancestor check with `true` passes every accept-side row |
| 4 | A marker strictly between the lineage and the destination (descendant of the run base, not at the tip) refuses | authorization unit surface | to be observed red at build start (outside the sanctioned recognition; no current test pins it) | Recognition scoped as "ancestor of the destination" alone would accept it; the source sanctions only ancestors of *both* the run base and the tip |
| 4 | A divergent marker planted after review, at promotion publication, refuses: the marker stays untouched and the run stays recoverable (re-entry still possible) | lifecycle surface driving the real authorization owner | to be observed red against a publication that swaps from wherever the marker sits without an ancestry guard | The accept-side rows alone are satisfiable by an unguarded read-and-swap; today the run-base compare-and-swap is the only thing preventing publication from clobbering a foreign marker, and this row is what keeps that protection |
| 4 | A marker that peels but whose ancestry cannot be decided (broken history behind it, so the ancestry probe errors) refuses, marker untouched | authorization unit surface | to be observed red at build start; the plain unreadable-marker case is already covered by existing fail-closed tests and runs unchanged | An ancestry check that reads a failed probe as "not divergent" would advance over a corrupt store |
| edge | The empty-run fast-forward still moves no marker and runs no gate | existing fast-forward lifecycle tests plus one new marker assertion | already covered for the gate half; the marker assertion is new and observed red only if the build wrongly adopts the rejected republication option | Pins the rejected first option: base and candidate move, the marker stays, promote remains the only green transition |
| edge | A refusal during recomposition still mutates nothing (state, candidate, branch, marker) | existing recomposition-refusal lifecycle tests, plus the story-4 rows' untouched-refs assertions | already covered for the lifecycle's handling of a refusal — the existing controls inject the refusal at the seam boundary and never reach the new predicate; the new predicate's own refuse-without-mutation is what the story-4 rows assert | Prevents the acceptance change from widening into partial writes on the refuse path, at both the lifecycle and the predicate |
| edge | Terminal re-entry and replayed marker advances stay idempotent | existing promote-twice and idempotent-bootstrap tests | already covered; run unchanged | The recognition rule adds an accept path; these controls red if it double-moves refs or replays operations |

### Edge inventory

- Error path — divergent-marker refusal row (story 4) and the
  refusal-without-mutation row above.
- Empty/absent input — absent marker with empty expectation, and absent marker
  with non-empty expectation, are already covered at the unit seam and run
  unchanged. **Won't handle:** automatic healing of a marker hand-deleted
  mid-promotion — a tampered ref outside the sanctioned lifecycle; publication
  keeps failing closed and promote re-entry stays available once the operator
  restores the marker.
- Boundary values — marker at the destination (already covered idempotent
  accept), marker equal to the expectation (already covered), marker between
  lineage and destination (story 4 row).
- Malformed input — abbreviated object IDs and ref names as expectations are
  already covered and run unchanged; the unreadable-object row covers the new
  ancestry read.
- Interrupted/partial state — the crash-between-branch-and-marker row
  (story 3), at the existing injected fault points.
- Re-run idempotency — the terminal re-entry row above.
- Process-boundary lifecycle — the fresh-service reload assertion in story 1's
  second row; existing fast-forward and recomposition reload tests run
  unchanged.
- Hostile environment (profile checklist) — the fixtures already exercise a
  spec slug containing a space; state serialized by one process and reloaded
  by a fresh one is the reload row. **Won't handle:** control-byte and
  quoting classes — this change adds no CLI text surface, TOON sink, or
  parser; every new value is a Git object ID or ref the existing surfaces
  already carry.
- Marker readers — **Won't handle:** red-disposition attribution during the
  lagging-marker window: a candidate-caused prospective red on a
  fast-forwarded run is attributed `inherited` because inherited-green
  recognizes only a marker at the tip. Pre-existing behavior that this fix
  makes routine rather than wedged; it degrades a diagnostic label, never a
  green/red verdict, and repairing it changes attribution semantics — a
  reviewer decision, priced in out of scope.

## Out of scope

- Surfacing a lagging or desynchronized marker in `bench status` / `bench
  doctor` diagnostics — a separate observability capability; ~4 edits, ~2 gate
  runs.
- Lagging-aware inherited-green attribution — teaching red-disposition
  attribution to recognize a lagging ancestor marker with retained tip
  evidence, so a fast-forwarded run's candidate reds stop degrading to
  `inherited`; ~3 edits, ~2 gate runs. Deferred not by size but because it
  changes what the marker means to attribution — flagged for the reviewer to
  pull into this build if wanted.
- An operator marker-resync verb — rejected by the decision source ("fix the
  state transition rather than add an operator rule"), not deferred.
