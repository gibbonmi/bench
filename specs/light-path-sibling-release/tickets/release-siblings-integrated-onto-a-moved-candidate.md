# Release siblings integrated onto a moved candidate

Blocked by: none

Ownership fence: `internal/worktree`, `internal/contract/runtime/runtime_worktree_sibling_release_test.go`
Assumptions: checkpoint refs remain retained and immutable; the spec-build lifecycle's own operations are unchanged

## What to build

`bench spec build integrate` succeeds and then fails to release the assignment's
checkout, reporting `provisional release evidence drifted; checkout retained`. The
integration is real — the candidate advances — but the assignment stays unreleased, and
`promote` refuses while any assignment is unreleased. A build whose tickets integrate in
sequence therefore cannot promote at all. Observed on a ten-ticket build where six of ten
assignments could not release.

The guard requires `checkpointTree == integratedTree`. That equality holds only for the
first ticket integrated at a given candidate base; every sibling after it legitimately
replays onto a moved candidate and produces a different integrated tree. The property the
guard exists for — releasing this checkout loses nothing — is carried by `liveTree ==
checkpointTree` together with the retained checkpoint ref. The `integratedTree` conjunct
adds nothing about the checkout and is false by construction for sibling composition.

Drop that conjunct; keep every other condition exactly as it is. This is a safety check,
so the change must narrow nothing else: a genuinely drifted checkout, a missing checkpoint
ref, a checkpoint not based at the assignment base, and unsaved work in the checkout must
all still refuse.

## Acceptance

- [ ] [RS1] Two siblings integrated in sequence from the same candidate base both release, and the second's release is not refused for having a different integrated tree.
- [ ] [RS2] A checkout carrying work absent from its checkpoint still refuses to release and is retained.
- [ ] [RS3] A checkpoint whose parent is not the assignment base, or whose ref no longer resolves, still refuses to release.
