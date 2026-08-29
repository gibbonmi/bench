# Review pickup: worktree-merge

Frozen pair: base `b6eccb4f`, reviewed tip `d9da7542`. Three axes ran as
opus / medium read-only delegates. Raw findings: 7. Repair targets: 6 code
repairs and 1 spec amendment.

## Standards

Count: 2 hard, 3 judgment (no-op). Worst: the status-clean predicate is
derived twice.

- S1 `auto-fix` — `mergeCheckoutClean` (`internal/worktree/merge.go`) repeats
  the block in `landingDestination` (`internal/worktree/land_identity.go`):
  same status call, same `ParsePorcelainZStrict`, same paths loop, same
  refusal. Rule: "Two derivations of the same fact must collapse into one
  source." Route `landingDestination` through the one predicate.
- S2 `auto-fix` (contestable; the axis said `ask-user`) —
  `mergeOnAssignmentBranch` is a third `symbolic-ref HEAD == a.Branch`
  spelling beside `landingSource` (`land_identity.go`) and
  `validateCreationBundle` (`ownership.go`). All three files sit inside the
  spec's `internal/worktree/` fence, so the collapse is in approved scope.

## Spec

Count: 2. Worst: WM17 is closed below the seam the row names.

- P1 `auto-fix` as a spec amendment — WM17's test drives `mergeTargetTip`
  directly, because a linked worktree's HEAD cannot disagree with the shared
  ref inside one repository. The row's seam column now names the helper and
  the reason. Reviewer veto surface: the amended row in the coverage map.
- P2 `auto-fix` — `mergeDefaultBranchCommit` folds a failed ancestry query,
  and an unresolved default branch, into `owned=false`, so the refusal
  misattributes the value. The spec's discipline: "a failed query refuses
  instead of classifying." Raise a refusal that names the failed query.

## Coverage

Count: 3. Worst: an ambiguous `--from` prefix is swallowed and a colliding
commit merges at exit 0.

- C1 `auto-fix` — `mergeSiblingTip` discards every `selectAssignment` error
  as "no sibling", including `target is ambiguous`. Probe: `--from shared-p`
  with two `shared-prefix-*` assignments and a branch `shared-p` exited 0
  with `kind=merge`. New row WM35.
- C2 `auto-fix` — `LaneAuthority.namedMarkdown` splits `git diff --name-only`
  on `\n` under default `core.quotepath`; `café-notes.md` reaches the prose
  check C-quoted and is never graded. Use `-z` and split on NUL. New row
  WM36. The same shape exists in `internal/git/tree.go` and
  `internal/worktree/land.go`; those sites are outside this spec and stay as
  they are (reviewer decision).
- C3 `auto-fix` (contestable; the axis said `ask-user`) — the spec's lookup
  runs "over active, bundle-valid assignments", but `mergeSiblingTip` matches
  any state. Filter non-active assignments out of the sibling lookup, so the
  spelling falls through to the commit lookup. A matched active sibling whose
  bundle fails still refuses by component, because the bundle is the
  bootstrap proof.
