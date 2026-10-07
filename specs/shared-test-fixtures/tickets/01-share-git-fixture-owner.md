# Share generic Git commands and snapshots

Blocked by: none
Writes: internal/conformance/native_workflow_test.go, internal/conformance/release_probe_fixture_test.go, internal/gittest/commands.go (new), internal/gittest/commands_test.go (new), internal/gittest/gittest.go, internal/gittest/gittest_test.go, internal/gittest/identity.go (new), internal/gittest/snapshot.go (new), internal/gittest/snapshot_test.go (new), internal/reviewrecord/recordtest/fixture.go, internal/testrepo/working_tree.go
Covers: GF1, GF2, GF3, GF4, GF5, GF6, GF7, GF8, GF9, GF11, GF12, GF13, GF14, GF15, GF20, GF32, GF33, GF34, GF35, GF36, GF37

## What to build

Move the Git portion of testrepo into gittest and migrate every former export caller in this ticket. Keep GateFixture under testrepo.
Preserve exported error-returning Run and existing Output delegation. Add MustRun, OutputBytes, Commit, configuration, and identity-argument contracts from the accepted spec.
Preserve error-returning snapshot functions and the index-only Commit distinction. Use new files for headroom rather than enlarging gittest.go.

The existing Run definition is `internal/gittest/gittest.go:277`; Output delegates at line 268.
Keep callers that return Run errors, including `internal/commitment/commitmenttest/repo.go:107`, compatible without editing them.
Their error-returning APIs must not become fatal helper calls. Preserve existing root operands, including an empty root passed to Git's `-C` option.
New helper forms use the accepted current-directory behavior for an empty root.

Create a repository in a path with spaces. Stage one file and leave another unstaged. Commit only the index and inspect the message and identity. Then copy a visible tree with an executable, an ignored file, and a dangling link.

## Acceptance

- [ ] MustRun fails its calling test on a failed Git command (GF1).
- [ ] Output preserves trimmed combined output (GF2).
- [ ] OutputBytes preserves raw stdout with trailing whitespace and NUL bytes (GF3).
- [ ] A root with spaces and shell metacharacters reaches the intended repository (GF4).
- [ ] Missing Git makes MustRun fail the calling test (GF5).
- [ ] Commit preserves the supplied message and explicit empty-commit option (GF6).
- [ ] Commit leaves an unstaged sentinel outside the new commit (GF7).
- [ ] RepoOnBranch and snapshots use the same canonical identity (GF8).
- [ ] An explicit identity override survives migration (GF9).
- [ ] The snapshot preserves visible file bytes and modes (GF11).
- [ ] The snapshot preserves live and dangling symlink targets (GF12).
- [ ] The snapshot excludes ignored untracked files and source Git metadata (GF13).
- [ ] A special source entry returns an error before a blocking read (GF14).
- [ ] A failed snapshot Git command returns its process error (GF15).
- [ ] GateFixture keeps its script and input-manifest behavior (GF20).
- [ ] Commit without an empty-commit option fails on an unchanged index (GF32).
- [ ] A refused no-change Commit preserves HEAD (GF33).
- [ ] Run retains root operands and returns trimmed combined output on success without a testing handle (GF34).
- [ ] Run returns empty output and its wrapped command error to a calling goroutine (GF35).
- [ ] Missing Git makes Run return an error and empty output (GF36).
- [ ] Output delegates to Run and fails its calling test on nonzero exit or missing Git (GF37).

## Checkpoint verification

- `bench test --package ./internal/gittest`
- `bench test --package ./internal/reviewrecord/recordtest`
- `bench test --package ./internal/testrepo`
- `bench test --package ./internal/conformance`
- `bench diff`

The planned contract tests live in the already fenced `internal/gittest/commands_test.go`. Observe goroutine errors through a bounded channel result.
Keep fatal-wrapper observations separate from the returned-error assertions. Reuse the fixture owner and existing repository setup; do not add another generic runner.

Use the frozen baseline and the candidate on the same fixture inputs. Exclude commit timestamps and object IDs from comparisons unless the scenario grades them.
Record the observed facts and the named omission or swap in the review pickup. A changed outcome stops the migration.
