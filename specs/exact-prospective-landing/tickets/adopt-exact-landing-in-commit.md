# Adopt exact landing in bench commit

Blocked by: build-exact-landing-owner.md
Ownership fence: `internal/commit`, `internal/contract/runtime/runtime_commit_test.go`, `CHANGELOG.md`
Integration surfaces: landing request/result→build-exact-landing-owner.md; compiled subcommand dispatch→existing `cmd/bench/main.go` plus AC1; shell route→existing `bin/bench.sh` plus AC1; runtime contract registry→existing `TestRuntimeCommitContracts`; contract gate phase→existing `internal/gate/phases.go`; shipped behavior notice→`CHANGELOG.md`
Contracts: parsed message, optional spec slug, and literal caller paths cross `internal/commit`→`internal/landing`, asserted by AC1-AC8 against the real compiled command; landing outcome crosses `internal/landing`→`internal/commit`, asserted by AC1-AC8 through stdout, stderr, ref, index, and filesystem observations

## What to build

Replace the ambient cleanliness, gate, stage, and commit sequence in `bench commit` with the exact landing owner while preserving the command grammar, exit posture, literal path semantics, fail-fast spec validation, and branch-agnostic role. Extend the existing real-binary runtime family and record the user-visible concurrency behavior in the changelog.

## Acceptance

- [ ] [AC1] A named file lands through the real command while foreign tracked-modified and untracked files remain byte-identical and uncommitted; foreign staged content stays staged and absent from the commit while named staged-plus-unstaged content lands working-copy bytes.
- [ ] [AC2] A clean named path returns `nothing to commit` with zero gate runs and no publication even when unrelated staged content exists.
- [ ] [AC3] Existing literal path coverage remains green for spaces, glob characters, dash-led names, additions, deletions, staged deletions, rename halves, named and deleted directory descendants, duplicate arguments, and whole-segment sibling exclusion.
- [ ] [AC4] The public gate script observes the exact prospective tree with named changes and base bytes but no foreign staged, unstaged, untracked, or ignored content.
- [ ] [AC5] Gate red leaves the destination ref, real index, and all working-copy bytes unchanged, while invalid or already-implemented `--spec` still refuses before any gate run.
- [ ] [AC6] With `--spec`, the gate, published commit, and invoking checkout contain identical implemented bytes; gate red or compare-and-swap loss leaves the checkout spec staged.
- [ ] [AC7] A two-process winner/loser journey preserves the winning tip on CAS loss, then a fresh rerun recomposes and lands a second commit containing both independently attributed changes without mixed authorship.
- [ ] [AC8] After success every named path is clean and every unnamed staged blob plus unstaged, untracked, and ignored byte and classification matches its pre-command fingerprint.
- [ ] [AC9] `bench commit` retains its existing help, usage, operational exit codes, stdout/stderr posture, and errors that identify attribution, prospective authorization, destination CAS, or landed-but-incomplete reconciliation failures without exposing internal receipts.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| AC1 | restore the global unexplained-file refusal or seed the prospective tree from the real index | public runtime commit fixture | prepare named and foreign staged, modified, and untracked work, run the real command, expect only named bytes committed and all foreign state preserved |
| AC2 | retain the real-index `nothingStaged` predicate | public runtime empty control | stage a foreign file, name a clean file, run the command, expect refusal, zero gate tally, unchanged ref, and unchanged index |
| AC3 | replace literal pathspec handling with glob or raw-prefix semantics | existing hostile-path runtime cases | rerun every enumerated path case through the real command and expect exact names only |
| AC4 | call the ordinary ambient-checkout gate instead of exact-tree authorization | gate-checkout inspection fixture | prepare foreign state, run the command, have the gate assert named/base bytes and foreign absence, expect green only for the prospective checkout |
| AC5 | mutate the real index before authorization or accept gate red | runtime red-gate and spec fail-fast controls | fingerprint ref, index, and files, invoke each refusal, expect exact equality and the required gate tally |
| AC6 | move the spec flip after authorization | public runtime spec fixture | have the gate read the spec, run `--spec`, compare gate, commit, and checkout bytes, then repeat with red and CAS loss |
| AC7 | publish with unconditional update or reuse the losing process's stale plan | two-process runtime journey | advance destination between authorization and publish, expect winner preserved, then rerun fresh and inspect the two commits' content and authorship |
| AC8 | reset the whole index or checkout after publication | public runtime fingerprint test | capture index blobs and file classifications, publish named work, expect named cleanliness and byte-identical unnamed state |
| AC9 | print an internal receipt or collapse failure-stage diagnostics | existing command grammar/error runtime family | invoke help, usage, attribution, gate, CAS, and reconciliation outcomes, expect stable public posture and stage-specific bounded diagnostics |
