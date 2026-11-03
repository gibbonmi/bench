# Grade an unpublished candidate commit

Blocked by: none
Writes: internal/preflight/candidate.go (new), internal/preflight/candidate_test.go (new), internal/preflight/candidate_fixture_test.go (new), internal/git/status.go, internal/preflight/gather.go, internal/preflight/gather_inputs.go, internal/preflight/command.go, internal/preflight/default_branch.go, internal/preflight/default_branch_test.go, internal/diff/range.go, internal/diff/review_base_test.go, cmd/bench/commitment_test.go
Covers: PP01, PP02, PP03, PP04, PP05, PP06, PP07, PP08, PP09, PP32, PP33, PP34

## What to build

Add the exported in-process grade `preflight.BuildVerdictAt(root, slug, base, candidate string) (string, int)`. It grades a candidate commit that no ref names yet, and it prints the verdict text that `bench preflight build <slug> [--base <base>] --source-tip <candidate>` prints after HEAD moves to that candidate.

The grade resolves the candidate as a commit. A candidate that does not resolve refuses with the kind `cannot resolve candidate`, after the `cannot resolve --source-tip` precedent in `internal/preflight/gather.go:179`.
The grade substitutes the candidate for HEAD in each read that the spec lists under "The candidate grade": the default-branch ancestry, the derived review base, the source tip, and the dirty-source test.
The grade pins the candidate as the source-tip pin, so the verdict holds a green `tip-current` row.
With an explicit base, the grade resolves the range from that base to the candidate with `diff.ResolveSourceRange(root, base, candidate)`.

Keep each rule in one source. Add `diff.ResolveReviewBaseAt(root, tip)`, and make `diff.ResolveReviewBase(root)` call it with `HEAD`, so `bench diff` keeps its behavior. Give the default-branch ancestry and the dirty-source test a tip or a tree argument in the same way, and keep the HEAD callers on HEAD. The dirty-source test calls `git.WorktreeDirty` (`internal/git/status.go:53`), which compares with HEAD. Add its tree form in that file, so the dirty rule keeps one source.
Render through the one verdict renderer of `verdictCommand`, and keep the `commitment` row rule. If the `Ready` call moves out of `verdictCommand`, update the route row in `cmd/bench/commitment_test.go` so that the inventory names the new declaration.

`internal/preflight/gather.go` has 3 lines of headroom under the 400-line limit. Put the new grade in `internal/preflight/candidate.go`. Do not extend `internal/preflight/explicit_base_test.go`, which is already over the limit.
Build the merge-shaped fixture of PP33 with `RunGit` and `commit-tree -p HEAD -p <default tip>`. Put the fixture builders in `internal/preflight/candidate_fixture_test.go`, so each test file stays under the line limit. The existing `preflighttest` helpers are sufficient.

The public `--source-tip` flag keeps its HEAD comparison. No check row is added, so the bare row counts in `internal/preflight/source_tip_test.go` stay as they are.

## Acceptance

- [ ] PP01: `BuildVerdictAt` on a green candidate that HEAD does not hold exits 0 with a green `tip-current` row.
- [ ] PP02: a candidate with an unreadable completion-plan fence over a HEAD with a valid fence exits 1 with `completion-plan` red.
- [ ] PP03: a candidate that deletes a path that a decision-map source cites exits 1 with `map-sources` red.
- [ ] PP04: for the PP01, PP02, and PP03 fixtures, the candidate verdict text equals the published `--source-tip` verdict text after `update-ref` moves HEAD to that candidate.
- [ ] PP05: a candidate whose coverage map holds a two-predicate row exits 1 with the `coverage map invalid` refusal.
- [ ] PP06: with an explicit base, the verdict prints a `source` table whose tip is the candidate.
- [ ] PP07: a candidate that names no commit exits 1 with the refusal kind `cannot resolve candidate`.
- [ ] PP08: a base that is not an ancestor of the candidate exits 1 with the `--base is not an ancestor` refusal.
- [ ] PP09: `TestSourceTipMismatchRendersRedRow` stays green.
- [ ] PP32: `TestDecideHasNoCheckFloorRow` stays green.
- [ ] PP33: a merge candidate whose second parent is the default-branch tip over a HEAD behind that tip gives the published verdict text with `base-current` green.
- [ ] PP34: with the candidate's edits uncommitted in the checkout and one unfenced changed path equal to the default-branch tip, the candidate verdict text equals the published verdict text with `paths-authorized` green.
