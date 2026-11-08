# Start no background lane pass for a Markdown-only commit or merge

Blocked by: none
Writes: internal/gate/background_launch.go, internal/gate/background_launch_test.go, internal/gate/background_queue_test.go, internal/gate/lane_select.go, internal/commit/background.go, internal/commit/background_test.go, internal/worktree/merge.go, internal/worktree/merge_background_test.go, internal/worktree/parallel_census_test.go, .bench/BENCH-reference.md, .bench/commitment.json, CHANGELOG.md, internal/systemtest/markdown_merge_background_test.go (new), internal/systemtest/parallel_census_test.go, internal/systemtest/parallel_census_cases_test.go, internal/systemtest/serial_census_test.go, internal/systemtest/serial_table_test.go
Covers: SP21, SP22, SP23, SP24, SP25, SP26, SP27, SP28, SP29, SP30, SP31, SP32, SP36

## What to build

Review chunk: SP-C3.

In `gate.StartBackground` (`internal/gate/background_launch.go:34-41`), after the existing worktree check, read the change list of `<tip>^1` to `<tip>` through `benchgit.TreeChanges` (`internal/git/tree.go:301`). Start no pass, and print nothing, when all three hold:

- the list holds at least one change;
- `prose.IsSubjectName` (`internal/prose/walk.go:27`) accepts every changed path;
- `freshness.DiffChangesBuildInputs(root, tip+"^1", tip)` answers false.

An empty list, a root commit, and a list that `TreeChanges` cannot read each start the pass as today. Both callers reach this one launch (`internal/commit/background.go:77`, `internal/worktree/merge.go:137`), so neither caller gets its own copy. Keep `startGrade` in `background_launch.go` under its name, because `internal/conformance/process_lifetime_test.go:53` pins it.

Change `proseSubject` (`internal/gate/lane_select.go:41`) and the match of the lane's `markdown` class (`internal/gate/lane_select.go:147-151`) to call `prose.IsSubjectName`. Keep the class name and its checks, so the rendered profile lane table does not change. Import `internal/freshness` as `benchfreshness`, as `internal/gate/prospective.go:10` does, because the gate package already declares a `freshness` identifier (`internal/gate/verdict.go:25`).

The new merge test calls `t.Parallel()`, and the worktree package pins its test count exactly (`internal/worktree/parallel_census_test.go:926`), so raise that pin by the number of new top-level tests. Write the launch tests with `newLaunchFixture` and `gittest`, not with `exec.Command`. The SP27 test writes its own minimal module: `go.mod`, a `cmd/bench` main package, and a `scripts/go-build.inputs` line that names the Markdown file.

Write SP23 as a system test in `internal/systemtest/markdown_merge_background_test.go`, because a Markdown-only merge selects the kit lane, which the BG6 merge fixture cannot run. The test obeys the serial and parallel census rules of `internal/systemtest`.

Add the exemption to `.bench/BENCH-reference.md` `### Background grade`, beside the sentence that each lane pass detaches one background pass. Add one `CHANGELOG.md` entry. This is the last ticket of the spec, so it alone names `.bench/commitment.json`.

## Acceptance

- [ ] A kit-worktree commit of one Markdown file starts no lane pass (SP21).
- [ ] A kit-worktree commit of one Markdown file and one text file starts one lane pass (SP22).
- [ ] A non-current merge whose fold brings only one Markdown file starts no lane pass (SP23).
- [ ] A launch of a tip that only deletes a Markdown file starts no pass (SP24).
- [ ] A launch of a tip that renames `notes.txt` to `notes.md` starts one pass (SP25).
- [ ] A launch of a tip whose only change is `a "b".md` starts no pass (SP26).
- [ ] A launch of a tip whose only change is a Markdown build input starts one pass (SP27).
- [ ] A launch of a root commit and a launch of a tip that names no commit each start one pass (SP28, SP29).
- [ ] `TestSecondGradeStartsNoSecondOracle`, `TestCommitStartsABackgroundGrade`, and `TestMergeStartsABackgroundGrade` stay green (SP30, SP31, SP32).
- [ ] `### Background grade` states the exemption (SP36).

## Verification

Run `bench test --package ./internal/gate`, `bench test --package ./internal/commit`, and `bench test --package ./internal/worktree` whole. This ticket writes system-tagged test files, so `BENCH_KIT` must name the candidate kit for each system run; `bench test --check system` supplies it. Run `bench test --check system --run TestMarkdownOnlyMergeStartsNoGrade`, then `bench test --check system` once. Observe each planned row red before its production edit. Self-probe: read the change list with rename detection on, and report that SP25 goes red.
