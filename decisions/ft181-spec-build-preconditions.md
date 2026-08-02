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

## Not yet specified

## Spec-writer discretion

- Error-message wording and the exact classification helper shape for #2/#3,
  provided the liveness/ownership boundary decided above is preserved.

## Out of scope

- Any change to promotion's gate cadence or the compare-and-swap integrate
  contract; the four faces are precondition and ownership classification only.

## Sources
