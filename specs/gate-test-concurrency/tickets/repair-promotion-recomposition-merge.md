# Repair promotion recomposition of already-landed changes

Blocked by: none
Ownership fence: `internal/specbuild/recompose.go`, `internal/specbuild/promotion_recompose_test.go`
Integration surfaces: `internal/specbuild/assign.go` promote recomposition call→`internal/specbuild/recompose.go` + PRM1; `internal/specbuild/integrate.go` assignment replay→unchanged `replayCheckpoint` + PRM1; `internal/specbuild/refresh.go` preservation replay→unchanged `replayCheckpoint` + PRM1
Contracts: recorded base, working tip, and old candidate cross `internal/specbuild/recompose.go` as three immutable commit IDs to Git's three-way tree merge; success returns exactly one verified tree ID; conflict or malformed output reaches no Bootstrap, ref, state, HEAD, or green-marker mutation, asserted by PRM1 against the real promotion path
Closure: PRM1/identical-change, PRM1/compatible-change, PRM1/divergent-conflict, PRM1/explicit-base, PRM1/tree-only, PRM1/validate-before-commit, PRM1/empty-output, PRM1/extra-output, PRM1/non-tree-output, PRM1/command-failure, PRM1/cancellation, PRM1/replay-separation

## What to build

Promotion recomposition merges the complete working tip and old provisional
candidate against the run's recorded base. Use Git's tree-only three-way merge
mode with the base supplied explicitly; do not infer it and do not touch an
index or worktree. A clean merge yields one tree ID, which must resolve as a
tree before candidate commit construction. Conflict, command failure,
cancellation, extra output, or a non-tree ID refuses before Bootstrap or any
protected mutation.

Keep `replayCheckpoint` unchanged for integration and refresh. Those operations
apply one attributed assignment delta and must continue to conflict when an
overlapping non-empty payload is already present. Only promotion combines two
complete descendants where an identical change on both sides is already
satisfied rather than conflicting.

## Acceptance

- [ ] [PRM1] (covers local) promotion cleanly recomposes identical and compatible same-path changes into their exact union, still refuses divergent changes and malformed merge results without protected mutation, and integration/refresh retain byte-exact assignment replay.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| PRM1/identical-change | restore cumulative patch application for promotion | the identical-same-path recomposition regression | apply, run focused, expect the already-landed hunk conflict red |
| PRM1/compatible-change | reduce the merge to whole-file selection | the existing compatible-same-path regression | apply, run focused, expect one side's independent line change missing or conflict red |
| PRM1/divergent-conflict | accept a conflicted merge tree | the existing divergent-same-path refusal | apply, run focused, expect Bootstrap/protected-state mutation or unexpected green red |
| PRM1/explicit-base | omit the recorded base or let Git infer it | the injected runner argv audit | apply, run focused, expect the missing exact `--merge-base=<run.Base>` argument red |
| PRM1/tree-only | replace `merge-tree --write-tree` with index- or worktree-backed merge plumbing | the argv audit plus index/worktree snapshots | apply, run focused, expect the command-shape or checkout snapshot red |
| PRM1/validate-before-commit | construct the candidate commit before verifying the merge result resolves as a tree | the injected runner call-order audit | apply, run focused, expect `commit-tree` before successful `cat-file` validation red |
| PRM1/empty-output | accept empty successful merge output | the injected empty-output runner case | apply, run focused, expect refusal before Bootstrap and protected mutation |
| PRM1/extra-output | accept a tree ID plus an extra output line | the injected extra-output runner case | apply, run focused, expect refusal before Bootstrap and protected mutation |
| PRM1/non-tree-output | accept one well-shaped object ID that does not resolve as a tree | the injected non-tree runner case | apply, run focused, expect refusal before Bootstrap and protected mutation |
| PRM1/command-failure | treat a merge command error as a clean result | the injected failing-runner case | apply, run focused, expect the error to propagate before Bootstrap and protected mutation |
| PRM1/cancellation | swallow context cancellation from the merge command | the injected cancellation case | apply, run focused, expect `context.Canceled` before Bootstrap and protected mutation |
| PRM1/replay-separation | route integration or refresh through promotion's whole-candidate merge | existing non-empty overlapping-payload conflict tests | apply, run focused integration and refresh cases, expect their attributed conflict refusal to disappear |
