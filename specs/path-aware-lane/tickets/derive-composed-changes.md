# Derive the lane's change list from the composed tree

Blocked by: none
Writes: internal/gate/lane_select.go (new), internal/gate/lane_select_test.go (new), internal/gate/lane.go, internal/gate/authorization/authorization.go, internal/gate/authorization/lane_test.go, internal/commit/commit.go, internal/commit/lane_test.go, internal/worktree/merge.go, internal/worktree/merge_test.go, internal/worktree/land.go, internal/status/status_test.go

## What to build

The lane grades the diff between the expected base tree and the composed tree.
A commit that names the directory `docs` therefore reaches the prose check with
`docs/note.md`, and a spec named by its folder reaches prose.

The package `internal/gate` owns the derivation. A new file
`internal/gate/lane_select.go` declares `ComposedChange` with the status letter,
the source mode, the destination mode, and the path. `ComposedChanges(root,
base, tree)` runs `git diff --raw --no-renames -z <base>^{tree} <tree>` through
`benchgit.Raw` and parses the NUL frames. The `-z` framing keeps a path's own
bytes, so a path with a space or a non-ASCII byte stays verbatim. The
`--no-renames` flag makes a rename one `D` entry and one `A` entry, so no `R`
entry reaches a reader.

The package `internal/gate/authorization` owns the authority. `LaneAuthority`
gains `Base` and loses both `NamedMarkdown` and `PreviousTip`. `Authorize` calls
`ComposedChanges` once. It resolves the prose placeholder to the changed
Markdown whose destination mode is a regular file, so a deleted Markdown file
contributes no subject. `LaneRequest` loses `NamedMarkdown` and gains `Changes`.
The authority passes its one derivation there, `RunLane` resolves the prose
placeholder from it, and the authority keeps today's lane line.

The package `internal/gate` also gains the `Lane` value that `LaneForCommit`
returns: the declared checks and the source root, in place of today's two
positional results. `select-kit-lane-checks.md` extends that same type with a
`Selective` field. The `Lane` type, `LaneForCommit`, and `ComposedChanges` live
in `internal/gate/lane_select.go`, because `internal/gate/lane.go` is already
over the structure budget and must not grow. The join field `mergeLane` in
`internal/worktree/land.go` follows the new signature, and `mergeOwner` in
`internal/worktree/merge.go` reads the value.

The package `internal/commit` passes the expected base commit it already reads
from `HEAD^{commit}`. The package `internal/worktree` passes the merge's
previous tip. The attribution owner still runs before the composition, so a
special or unreadable named path refuses with the existing message and no lane
run starts.

The derivation tests attach to `internal/gate/lane_select_test.go`, with
fixtures by the precedent of `laneFixture` and `commitFiles` in
`internal/gate/authorization/lane_test.go`. The commit tests attach to
`internal/commit/lane_test.go` beside `TestLaneProseGradesOnlyTheNamedMarkdown`,
with the manifest fixture that `laneRepo` builds. The merge row is the existing
`TestMergeResolvesTheProsePlaceholderToTheIncomingMarkdown`, which the build
session reads before it renames the field.

## Acceptance

- [ ] PL1: `ComposedChanges` over a base that holds `docs/a.md` and `docs/b.go`, and a tree that changes both, lists exactly those two paths with status `M`.
- [ ] PL2: a tree that renames `docs/old.md` to `docs/new.md` lists `docs/new.md` with status `A` and `docs/old.md` with status `D`, and no `R` entry.
- [ ] PL3: a tree that adds a symbolic link `link.go` lists it with destination mode `120000`.
- [ ] PL4: the base holds `kept.md`, `changed.md`, and `gone.md`. A tree that changes `changed.md`, adds `added.md`, and deletes `gone.md` hands the prose check exactly `added.md changed.md`.
- [ ] PL5: a tree that adds `café notes.md` hands the prose check that path's own bytes as one argument.
- [ ] PL6: `bench commit -m m fifo` where `fifo` is a FIFO exits 1, prints `special file "fifo" is not attributable`, and leaves the lane's tally file absent.
- [ ] PL40: a named path with mode `000` refuses with `unreadable path`, and an empty change list refuses with `nothing to commit`, each before the authority runs.
- [ ] PL7: a merge whose incoming side changes `incoming.md` hands the prose check exactly the incoming Markdown.
- [ ] PL8: `LaneAuthority` declares `Base` and declares neither `NamedMarkdown` nor `PreviousTip`.
- [ ] PL25: the manifest fixture commit that names the directory `docs` holding a 27-word sentence in `docs/note.md` exits 1 and prints `docs/note.md:3:`.
- [ ] the gate `test` phase stays green for the whole `internal/worktree` package.
