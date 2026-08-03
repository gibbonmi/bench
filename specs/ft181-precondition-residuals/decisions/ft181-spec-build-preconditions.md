# FT181 spec-build precondition residuals

Status: ready

## Destination

Fix directions for the four reviewer-call faces left open by the shipped
spec-build-lifecycle-preconditions build, all in `internal/specbuild`'s
precondition and ownership layer. This map is the decision source for the
FT181 spec; seams, tests, and coverage belong to `/bench-write-spec`.

## #1: Restart-after-terminal refuses on a sibling-advanced green marker

Blocked by: none
Type: Grill

### Question

Restart passes `run.Base` as `previousGreen`, so a green marker advanced by a
sibling build still refuses restart of an abandoned run — the same
benign-marker refusal story 4 removed, surviving on the sibling path. Prefer
the live marker, or keep `run.Base` as anti-tamper evidence?

### Answer

Resolved 2026-08-02: prefer the live green marker on restart. Restart reads
the current marker instead of `run.Base`, so a sibling build's advance stops
refusing an abandoned run's restart. `run.Base` stays recorded as evidence,
not as the comparison operand.

## #2: A husk or dangling symlink at an assignment path is fatal even for abandon

Blocked by: none
Type: Grill

### Question

Liveness classifies by `Lstat` alone: a directory whose `.git` is gone, or a
dangling symlink, reads as `errOwnership` — fatal even for `abandon`. Does
present-but-not-a-checkout soften to liveness?

### Answer

Resolved 2026-08-02: yes. A husk (directory present, `.git` gone) or a
dangling symlink classifies as liveness, so `abandon` can proceed. Identity
mismatches remain fatal.

Amended 2026-08-03 after doc review: the specbuild classification alone does
not deliver the outcome — abandon routes any `Lstat`-present target through
`internal/worktree`'s explicit planner, which rejects a non-checkout. The
worktree-layer change is authorized: the abandon planner treats
present-but-not-a-checkout like the already-handled removed-checkout case, so
the map's scope includes `internal/worktree`'s abandon planning path.

Amended 2026-08-03 at spec falsification (reviewer decision): the bytes of a
present-but-not-a-checkout path are preserved. Abandon releases the
registration and intent entry through a new non-deleting plan action that
names the leftover path; it never inherits the removed-checkout path's
force-removal. Disposal of the leftover bytes routes through the existing
size-bounded clean surface.

## #3: The prepared-abandon exemption softens identity errors

Blocked by: none
Type: Grill

### Question

The prepared-abandon exemption softens ANY ownership error, identity
included, once a prepared abandon operation exists. Narrow it?

### Answer

Resolved 2026-08-02: narrow it by classifying "registration no longer found"
as liveness rather than ownership, and remove the blanket swallow of any
ownership error under a prepared abandon — both moves together are what
makes the exemption never reach identity errors; identity stays fatal
everywhere.

## #4: A run whose tip moves before any checkpoint wedges

Blocked by: none
Type: Grill

### Question

Recomposition cannot replay an empty candidate (`No valid patches in input`),
so a run whose tip moves before any checkpoint wedges. **Post-approval
correction, flagged:** the source capture said checkpoint, start, and abandon
all refuse, but the shipped tree already exempts abandon from the
recomposition refusal (`precondition.go:84`), so the wedge is checkpoint and
start only. Fast-forward the base, or keep the standing rule of sequencing
capture and `main` commits outside run windows?

### Answer

Resolved 2026-08-02: fast-forward the base rather than replaying zero
patches, unwedging checkpoint and start (abandon is already exempt). When
this face ships, the standing run-window commit-sequencing rule retires with
the roadmap row that carries it.

Corrected 2026-08-03 after doc review: the "cannot replay an empty candidate"
mechanism is stale — promote's recomposition already fast-forwards an empty
run directly onto the advanced tip. The residual wedge is the blanket
recompose refusal of checkpoint and start in the mutation precondition; the
decision stands as extending that same fast-forward treatment to checkpoint
and start instead of refusing them.

Amended 2026-08-03 at spec falsification (reviewer decision): the
fast-forward's operation set is closed — checkpoint and start, on
non-terminal runs only. `assign` keeps today's refusal (start first, then
assign); promote keeps its recomposition route and Bootstrap; review keeps
refusing without mutation; terminal runs keep the restart path, which
depends on the recompose refusal firing.

## Not yet specified

## Spec-writer discretion

- Error-message wording and the exact classification helper shape for #2/#3,
  provided the liveness/ownership boundary decided above is preserved.

## Out of scope

- Any change to promotion's gate cadence or the compare-and-swap integrate
  contract; the four faces are precondition and ownership classification,
  plus the `internal/worktree` abandon-planning change #2's amendment
  authorizes.

## Sources

- Path: `specs/ft181-precondition-residuals/decisions/assets/ft181-restart-marker-research.md`
  Supports: #1 — marker/Bootstrap mechanics, the exact refusal, and the untested restart path; a 2026-08-03 delegated read confirming the answer.
  Drift: re-verify if `internal/gate/authorization/authorization.go` or the terminal-restart branch in `internal/specbuild/assign.go` changes before the spec reads this map.
- Path: `specs/ft181-precondition-residuals/decisions/assets/ft181-liveness-classification-research.md`
  Supports: #2 and #3 — the classification table, why #2's worktree amendment is load-bearing, the exemption's reach, and the untested states; a 2026-08-03 delegated read confirming both answers.
  Drift: re-verify if `internal/specbuild/precondition.go` or `internal/worktree/resume.go`'s abandon planning changes before the spec reads this map.
- Path: `specs/ft181-precondition-residuals/decisions/assets/ft181-empty-recompose-research.md`
  Supports: #4's correction — commit evidence that promote's fast-forward shipped (`2874d94`) and that the residual wedge is checkpoint/start only.
  Drift: re-verify if `internal/specbuild/recompose.go` or the `errRecompose` refusal in `precondition.go` changes before the spec reads this map.
