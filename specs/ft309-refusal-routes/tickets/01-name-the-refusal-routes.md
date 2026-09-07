# Name the route in the fold, landing, and growth refusals

Blocked by: none
Writes: internal/worktree/merge.go, internal/worktree/merge_from_sha_test.go (new), internal/worktree/land.go, internal/worktree/land_identity.go, internal/worktree/land_identity_test.go, internal/structure/structure.go, internal/structure/growth_test.go, ROADMAP.md, roadmap/FT309.md (deleted)
Covers: none

## What to build

An operator who reads the handoff pins a sibling's tip sha. So
`bench worktree merge --from <sha> <target>` accepts the sha of an active
sibling's branch tip and composes the same fold that the sibling's label
composes. The sibling checks stay in force: the sibling is on its branch, its
checkout is at the tip, and the checkout is clean. A sha that is no sibling
tip and is outside the default branch's history still refuses, and that
refusal names both routes: a default-branch commit, or a sibling's label or
tip.

The label-and-commit ambiguity refusal stays as it is.

`bench worktree land --base <commit>` refuses a base outside the
destination's history. Its detail now names the distinction between the two
bases. The review reads the fold commit as its frozen base. The landing
takes the default-branch tip that the source folded. The detail text is:
`review base is not an ancestor of the landing destination: --base takes the landing base, the default-branch tip the source folded, not the fold commit the review read`.
The observed and wanted fields stay as they are.

`internal/worktree/land.go`
is five lines under its budget, so the new detail constant lives in
`internal/worktree/land_identity.go` beside `identityRefusal`, and
`land.go` only references it.

The growth ratchet's red summary names its base commit. The lane grades
growth against the tip, and a spec grades it against the spec base.
The red summary line becomes
`structure growth: N file(s) grew past budget since <base>. Split along responsibility (see the craft-seams skill), or record a reviewer grant in .bench/structure-accept.`
The green line already names the base and does not change.

`internal/worktree/merge_test.go` is over its budget. So the new fold tests
live in a new file `internal/worktree/merge_from_sha_test.go`. That file
reuses the fixtures `merge_test.go` already exports to the package.

The landing retires the FT309 row: the index line and the sequence entry
leave `ROADMAP.md`, the sequence renumbers, and `roadmap/FT309.md` leaves the
tree.

## Acceptance

- [ ] `bench worktree merge --from <sibling tip sha> <target>` composes the sibling fold with the same kind, subject, and publication as `--from <sibling label>`.
- [ ] `bench worktree merge --from <sibling tip sha> <target>` still refuses when the sibling checkout is dirty, detached, or behind its tip, with the existing sibling refusal details.
- [ ] `bench worktree merge --from <sha> <target>` with a sha that is no sibling tip and is outside the default branch's history refuses, and the detail names both routes.
- [ ] `bench worktree land --base <fold commit>` refuses with the detail text stated above, `observed=` the base, and `wanted=` the destination.
- [ ] `bench structure --growth <base>` prints the red summary line with `since <base>` when a file grew past budget, and `internal/structure/growth_test.go` reds when the base is missing from that line.
- [ ] `ROADMAP.md` carries no FT309 line, its sequence is renumbered, and `roadmap/FT309.md` is absent.
