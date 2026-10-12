# Record a fold of main and accept it in the review chain

Blocked by: 2-route-a-moved-destination.md
Writes: cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/gate/authorization/authorization.go, internal/gate/review_checkpoint_commits_test.go, internal/git/admin_readers_test.go, internal/git/default_branch.go, internal/git/git.go, internal/landing/composition.go, internal/landing/composition_fold_test.go (new), internal/landing/composition_test.go, internal/landing/landing.go, internal/landing/merge.go, internal/reviewrecord/coverage.go, internal/reviewrecord/fold.go (new), internal/reviewrecord/fold_test.go (new), internal/reviewrecord/parse.go, internal/reviewrecord/record.go, internal/reviewrecord/recordcmd/command.go, internal/reviewrecord/recordcmd/fold.go (new), internal/reviewrecord/recordcmd/fold_test.go (new), internal/reviewrecord/recordcmd/ordering_test.go, internal/reviewrecord/recordcmd/refusal.go, internal/reviewrecord/recordtest/fixture.go, internal/reviewrecord/recordtest/fold.go (new), internal/reviewrecord/repair.go, internal/reviewrecord/repair_test.go, internal/reviewrecord/source_test.go, internal/reviewrecord/write.go, internal/testreport/floor.go, internal/worktree/completion_fixture_test.go, internal/worktree/merge.go
Covers: BW25, BW26, BW27, BW28, BW29, BW30, BW31, BW32, BW33, BW34, BW35, BW36, BW37, BW38, BW41

## What to build

`bench record fold <slug> --source <commit>` records a fold of `main` into the build source, and the review chain accepts a recorded fold between chunks and after the last chunk.
A clean fold has an empty resolution and needs no review.

`Record` gains `folds`, a list of `Fold{ID, After, From, Commit, Incoming, SourceDigest, Resolution, ResolutionDigest, TestPackages, Reviews}`.
`ID` counts `F1`, `F2`, and so on, and `After` names the chunk that the fold follows.
The form refuses, each before any write:

- a source that is not HEAD: `fold source <commit> is not HEAD; record the fold before the next commit`
- a commit without exactly two parents: `fold <commit> is not a two-parent merge`
- a second parent that is not an ancestor of the default branch tip: `fold <commit> does not merge a default-branch commit`
- a first parent whose digest is neither the covered digest nor a comment-only gap from it: `fold <commit> does not follow the reviewed source`

The covered digest is the source digest of the last chunk entry, or of the last fold after it.
The automatic tree is the landing composition of the two parents (`internal/landing/composition.go:54-90`).
`Compose` drops its tree on a conflict today, so expose the conflict-marker tree that `git merge-tree --write-tree` writes through the same owner.
`Resolution` lists, in byte order, each path where the fold commit's tree differs from the automatic tree, without the review record path.
Read paths through `-z` plumbing, so a path with a space or a double quote keeps its exact bytes.
`ResolutionDigest` is the digest of that diff.

`TestPackages` is the sorted package list that `testreport.WritesFloor` derives from the `_test.go` paths that differ between `From` and `Commit`.
`internal/reviewrecord` cannot import `internal/testreport` or `internal/landing`, because both depend on it.
So the record command in `internal/reviewrecord/recordcmd` computes the resolution and the packages and stores them.

`checkSource` (`internal/reviewrecord/coverage.go:19`) changes in two places in this ticket.
A later chunk base passes when the last fold after the predecessor covers it.
That fold covers the base when three conditions hold:

- The base digest equals the fold's `SourceDigest`.
- The fold commit is an ancestor of the base.
- The fold's `From` covers the predecessor's digest.

A source past the last chunk tip passes when fold records cover it.
The chain-gap refusal becomes a repair of form `fold`, and its text names `a default-branch merge lands before the first chunk or as a recorded fold`.
The stale-source refusal becomes a repair of form `fold` when the uncovered range holds a fold of a default-branch commit, and it keeps `stale reviewed source: no chunk review covers`.
Ticket 4 adds the review rule for a non-empty resolution.

Register the form in the record command registry and the help inventory.
Keep `internal/reviewrecord/repair.go`'s write order coherent, and update `internal/reviewrecord/recordcmd/ordering_test.go` if the fold joins that order.
`internal/reviewrecord/recordtest/fixture.go` has about 20 lines of headroom, so put the fold fixture builder in `recordtest/fold.go`.

## Acceptance

- [ ] `bench record fold` over a fold that `bench worktree merge` published writes entry `F1` with `After`, `From`, `Commit`, `Incoming`, and the source digest.
- [ ] That clean fold's `Resolution` is empty.
- [ ] A hand-resolved conflict fold lists exactly the conflicted path in `Resolution`.
- [ ] A clean fold with one extra hand edit lists exactly the edited path in `Resolution`.
- [ ] A resolution path that holds a space and a double quote is recorded with its exact bytes.
- [ ] A non-HEAD source, a one-parent commit, a sibling-tip merge, and a fold over unreviewed work each refuse before any write.
- [ ] A chunk-2 checkpoint passes when chunk 2's base is the record commit that follows a recorded clean fold over the chunk-1 tip.
- [ ] The same checkpoint without the fold entry refuses as a repair of form `fold`.
- [ ] The chain-gap refusal holds `a default-branch merge lands before the first chunk or as a recorded fold`.
- [ ] `bench record completion` passes at a source past the last chunk tip when a clean recorded fold covers it.
- [ ] Without the fold entry, that completion refuses with `stale reviewed source` as a repair of form `fold`.
- [ ] A fold that brings in `pkg/a/a_test.go` stores `TestPackages` equal to `testreport.WritesFloor` over that path.
