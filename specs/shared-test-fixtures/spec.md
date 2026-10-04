# Share Git fixtures and capability reasons

Status: staged

Roadmap: FT360

Decision source: `roadmap/FT360.md`, the reviewed quality-survey artifact from drain `d-0bca6e72fedd`.

Verification log: 2 iteration(s) to accept the spec — GPT-6.1 Sol/high accepted the caller-class and no-change commit repairs. Slice review is pending.

## Problem

Generic fixture helpers repeat Git execution and default identities. The current `gittest` package exports `Output`, while its run helper remains private.
The Git portion of `testrepo` has another command loop and another identity. Callers also write their own reasons for unavailable FIFOs and symlinks.
The source reads and candidate inventory below identify these copies. The old survey counts are historical evidence, not acceptance targets.

## Solution

Tests use `gittest` for generic Git execution, commit creation, and the default fixture identity. The Git portion of `testrepo` moves into that owner.
The capability owner formats unavailable-capability reasons. Callers retain the operation detail and the existing skip classification.

This spec prepares work for the approved quality successor milestone. It does not start implementation or change the active delivery commitment.
The implementation rechecks this inventory after the roadmap-delivery-commitment build lands. It includes new generic fixtures within the same approved outcome.

## User stories

Line: gpt-5.6-sol / high.
Implementation-line reason: The migration spans shared fixtures, raw Git output, process supervision, and conformance checks. Exact behavior and existing tests constrain the work.
Harder chunks: GF-C1, GF-C5C, GF-C5D, GF-C6, GF-C7, GF-C8, GF-C13.

### Generic Git fixtures

1. As a test author, I want one Git run helper, so that failure handling does not drift.
2. As a test author, I want the existing trimmed output contract, so that query assertions remain valid.
3. As a test author, I want raw stdout when required, so that paths and binary delimiters survive.
4. As a test author, I want arguments passed separately, so that unusual paths cannot become shell commands.
5. As a test author, I want missing Git to fail, so that absent evidence cannot look green.
6. As a test author, I want one commit helper, so that fixture commits share execution behavior.
7. As a test author, I want explicit staging, so that a commit does not absorb unrelated fixture files.
8. As a test author, I want one default identity, so that fixtures do not depend on personal Git configuration.
9. As a test author, I want intentional identity overrides, so that identity-sensitive tests retain their subject.

### Snapshot and caller migration

10. As a test author, I want the snapshot helper under the Git fixture owner, so that initialization and commits share their policy.
11. As a test author, I want visible snapshot entries preserved, so that copied fixtures represent the source tree.
12. As a test author, I want excluded snapshot entries absent, so that private state stays outside the fixture.
13. As a test author, I want snapshot errors returned, so that callers retain their error-handling contract.
14. As a maintainer, I want generic helper copies removed, so that a future correction has one owner.
15. As a maintainer, I want process supervisors preserved, so that the refactor cannot abandon child cleanup.
16. As a maintainer, I want specialized Git probes preserved, so that tests retain environment and failure evidence.
17. As a test author, I want gate fixtures preserved, so that their command manifests remain valid.

### Capability reasons and closure

18. As a test author, I want one unavailable-capability prefix, so that callers state only the failed operation detail.
19. As a reviewer, I want the skip class preserved, so that strict capability grading remains accurate.
20. As a reviewer, I want the skip record written first, so that an unavailable assertion leaves evidence.
21. As a maintainer, I want leftover generic copies detected, so that a partial migration cannot claim completion.
22. As a maintainer, I want negative controls on ownership checks, so that source text examples remain valid.
23. As a reviewer, I want before-and-after fixture evidence, so that the refactor proves its preserved behavior.
24. As a maintainer, I want each migration batch independently green, so that later batches need no broken intermediate state.

## Implementation decisions

### Terms and boundaries

A **generic Git fixture** prepares repository state without grading the process transport. Avoid: Git adapter, production runner.
A **fixture identity** is the default author and committer name and email. Avoid: operator identity, authentication identity.

A **specialized probe** grades a process environment, failure, stream, or lifecycle. Avoid: generic helper exemption.
An **unavailable capability** is an existing capability-class skip with an operation detail. Avoid: environment skip, successful assertion.

The production Git adapter remains unchanged. Test fixtures do not import it merely to create repositories.
The existing `gittest` dependency on `testrepo` ends when snapshot functions move. The remaining `testrepo` gate fixtures do not import `gittest`.
This direction keeps `internal/git` tests able to import `gittest` without a cycle.

### Git helper contract

Expose `Run(t, root, args...)` beside the existing `Output(t, root, args...)`.
`Run` fails the test on a nonzero exit. `Output` preserves trimmed combined output and the existing error behavior.
Provide `OutputBytes(t, root, args...)` for untrimmed stdout. It does not combine stderr with successful stdout.

One private execution owner implements these forms. Each public test helper calls `t.Helper()`.

Expose `Commit(t, root, message, options...)`. It commits the current index and passes each additional commit option as its own argument.
It does not stage files. Callers retain their existing `add` paths and flags, including the distinction between `-A` and forced additions.
Empty commits require the caller's explicit option. Ordinary no-change commits continue to fail.

Retain `bench test` and `bench@example.invalid` as the canonical default identity. Store each value once inside `gittest`.
Expose configuration and per-command identity arguments derived from those values. Callers must not restate those literals.
`RepoOnBranch` and moved snapshot initialization use that identity. Generated fixture scripts obtain their default identity from the same owner.

The helpers preserve ambient environment and current-directory behavior unless a caller already supplies an explicit root or environment.
An empty root means the current working directory. A nonempty root uses an argv operand, never a shell interpolation.
Do not add global environment mutation, implicit configuration isolation, or production security guarantees.
Tests that deliberately alter author, committer, timestamp, or Git configuration retain those inputs.

### Snapshot move

Move `CommitWorkingTree` and `CommitAll` from `testrepo` to `gittest`. Preserve their signatures and error-returning behavior.
`KitCopy` calls the moved function directly. Migrate every former export caller in the same chunk, then remove the old definitions.

The error-returning snapshot functions and fatal test wrappers share the private Git execution owner.
`CommitAll` retains its force-add behavior. Its semantics remain distinct from the index-only `Commit` helper.

Preserve the current snapshot membership rule. The source includes tracked files and visible untracked files, including tracked files matched by ignore rules.
Preserve regular-file bytes and modes, live symlinks, and dangling symlinks. Ignore source entries deleted after enumeration, as the current function does.

Keep `.git` and ignored untracked files outside the snapshot. Refuse other special source entries before reading their bytes.
The copy remains a private repository with one initial commit and no pending changes.
No transactional rollback promise is added for a failed copy.

### Migration boundaries

`inventory.json` records candidate sites against the pinned base. Its patterns discover work; they do not decide semantics.
Before editing a candidate, classify the enclosing operation as a generic fixture, an intentional subject, or a specialized probe.
Every generic command wrapper, default identity, and snapshot caller migrates. Domain-specific file setup can remain as composition around `gittest`.

A wrapper that only renames `gittest.Run` or `Output` leaves unless it preserves an existing public fixture API.
An API-preserving wrapper contains no execution or identity policy of its own.

Worktree `descendant` and the system `owner.runAt` remain the process authorities. Their Git wrappers retain transport, cleanup, and census behavior.
They use the shared default identity, while their existing stdout and stderr handling remains intact.

The preflight fast-import fixture retains its input stream and timestamps. Environment-policy probes retain their explicit child environments.
Git failure probes retain their expected nonzero status. Queries that authenticate the tested repository remain production-subject observations.

These are operation-level exceptions. A file containing one exception does not exempt its unrelated fixture setup.

The ownership check records each retained specialized site by path and enclosing symbol, with its preserved contract.
A new generic runner beside an allowed specialized symbol still fails. String literals containing example source are data.
The build records the final site classifications in the review pickup. It does not copy historical counts into a new policy table.

### Capability reasons

Add `Unavailable(t, class, detail)` in `capability`. It derives the prefix from the existing class token.
The reason shape is `<class> capability unavailable: <detail>`. The existing `Capability` helper writes the record and ends the test.
Callers pass the underlying operation error or the existing host limitation detail. They do not repeat a second capability prefix.

Migrate the FIFO and symlink caller family, including existing socket and device fixtures classified as FIFO.
Preserve their current class rather than creating a new capability vocabulary.
The new helper does not probe the host, change the skip condition, or change strict-mode policy.

Keep `Kind`, `Class`, `Name`, `Reason`, `BENCH_SKIP_LOG`, `Render`, and `ParseLine` unchanged.

Record each migrated call site and its baseline class in the independent class-preservation fixture.
Use path, enclosing symbol, and a stable case label to distinguish multiple calls within one function.
The checker compares real call sites with that fixture. Missing sites and changed classes fail.

This independent expectation qualifies under ADR 0006 only after the named class-change mutation produces a recorded red.
The mutation changes the socket case in `clean_landed_hostile_test.go` from FIFO to symlink. The implementation also records restoration.

The owner tests still exercise the low-level `Capability` API. Other capability families retain their existing reasons in this feature.

### Ownership enforcement

Extend the existing Git and skip ownership checks after the relevant caller batches migrate.
The Git check adds a fixture-specific rule without narrowing the current production administration-flag rule.
Reuse the existing module inventory and Go AST helpers. Do not build another general source scanner.

The fixture rule rejects generic Git process wrappers and default identity definitions outside `gittest`.
It reads real syntax, resolves import aliases, and uses exact specialized-site exceptions.
The skip rule rejects low-level FIFO and symlink `Capability` calls outside the capability owner after migration.
Literal fixtures and comments remain permitted. Existing environment skips and other capability classes remain permitted.

Keep existing registry names and execution routes. Add focused helper files when an existing checker lacks line-budget headroom.
The terminal migration enables these checks only after their complete consumer sets pass.
Before enablement, intermediate chunks retain all existing checks and add their own focused behavior evidence.

## Implementation chunks

Each ticket is one serial green checkpoint and one review chunk. Its predecessor supplies an accepted owner or the completed preceding migration batch.
Shared conformance registries make the listed order explicit. A later ticket cannot supply a missing helper to an earlier checkpoint.
The author rechecks the candidate sites before implementation. Source drift expands only the affected fences through the approved amendment route.

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| GF-C1 / `01-share-git-fixture-owner.md` | Share generic Git commands and snapshots | GF1, GF2, GF3, GF4, GF5, GF6, GF7, GF8, GF9, GF11, GF12, GF13, GF14, GF15, GF20, GF32, GF33 | Completion-plan checks | yes |
| GF-C2 / `02-migrate-raw-git-fixtures.md` | Migrate raw Git fixture callers | GF16, GF29, GF30 | Completion-plan checks | no |
| GF-C3 / `03-migrate-leaf-fixtures.md` | Migrate leaf command fixtures | GF16, GF29, GF30 | Completion-plan checks | no |
| GF-C4 / `04-migrate-workflow-readers.md` | Migrate workflow reader fixtures | GF16, GF29, GF30 | Completion-plan checks | no |
| GF-C5A / `05a-migrate-cli-adopt-fixtures.md` | Migrate CLI and adoption fixtures | GF16, GF19, GF29, GF30 | Completion-plan checks | no |
| GF-C5B / `05b-migrate-commit-gate-fixtures.md` | Migrate commit and gate fixtures | GF16, GF19, GF29, GF30 | Completion-plan checks | no |
| GF-C5C / `05c-migrate-preflight-conformance-fixtures.md` | Migrate preflight and conformance fixtures | GF16, GF19, GF29, GF30 | Completion-plan checks | yes |
| GF-C5D / `05d-migrate-workflow-execution-fixtures.md` | Migrate workflow execution fixtures | GF16, GF19, GF29, GF30 | Completion-plan checks | yes |
| GF-C6 / `06-share-landing-fixture-identity.md` | Share the landing fixture identity | GF17, GF29, GF30 | Completion-plan checks | yes |
| GF-C7 / `07-share-worktree-fixture-identity.md` | Share the remaining worktree fixture identity | GF17, GF29, GF30 | Completion-plan checks | yes |
| GF-C8 / `08-share-system-fixture-identity.md` | Share the system fixture identity | GF18, GF29, GF30 | Completion-plan checks | yes |
| GF-C9 / `09-own-capability-reasons.md` | Own unavailable-capability reasons | GF21, GF22, GF23, GF30 | Completion-plan checks | no |
| GF-C10 / `10-migrate-capability-leaf-callers.md` | Migrate capability callers in leaf packages | GF22, GF29, GF30 | Completion-plan checks | no |
| GF-C11 / `11-migrate-capability-workflow-callers.md` | Migrate capability callers in workflow packages | GF22, GF29, GF30 | Completion-plan checks | no |
| GF-C12 / `12-migrate-capability-conformance-callers.md` | Migrate capability callers in conformance and CLI tests | GF22, GF29, GF30 | Completion-plan checks | no |
| GF-C13 / `13-enforce-fixture-ownership.md` | Enforce completed fixture ownership | GF10, GF24, GF25, GF26, GF27, GF28, GF31, GF30 | Completion-plan checks | yes |

```bench-completion-plan
{"version":1,"chunks":[{"id":"GF-C1","tickets":["01-share-git-fixture-owner.md"],"verification":[{"id":"pkg-internal-gittest","command":"bench test --package ./internal/gittest"},{"id":"pkg-internal-reviewrecord-recordtest","command":"bench test --package ./internal/reviewrecord/recordtest"},{"id":"pkg-internal-testrepo","command":"bench test --package ./internal/testrepo"},{"id":"conformance","command":"bench test --package ./internal/conformance"},{"id":"preservation","command":"bench diff","probe":"Compare the frozen baseline and this ticket over its named fixture family. Record facts, specialized exceptions, and unchanged skip classes in the review pickup."}]},{"id":"GF-C2","tickets":["02-migrate-raw-git-fixtures.md"],"verification":[{"id":"pkg-internal-diff","command":"bench test --package ./internal/diff"},{"id":"pkg-internal-git","command":"bench test --package ./internal/git"},{"id":"preservation","command":"bench diff","probe":"Compare the frozen baseline and this ticket over its named fixture family. Record facts, specialized exceptions, and unchanged skip classes in the review pickup."}]},{"id":"GF-C3","tickets":["03-migrate-leaf-fixtures.md"],"verification":[{"id":"pkg-internal-consumers","command":"bench test --package ./internal/consumers"},{"id":"pkg-internal-gitguard","command":"bench test --package ./internal/gitguard"},{"id":"pkg-internal-guards","command":"bench test --package ./internal/guards"},{"id":"pkg-internal-outline","command":"bench test --package ./internal/outline"},{"id":"pkg-internal-poolkey","command":"bench test --package ./internal/poolkey"},{"id":"pkg-internal-probe","command":"bench test --package ./internal/probe"},{"id":"pkg-internal-skillsindex","command":"bench test --package ./internal/skillsindex"},{"id":"pkg-internal-spec","command":"bench test --package ./internal/spec"},{"id":"pkg-internal-stophook","command":"bench test --package ./internal/stophook"},{"id":"pkg-internal-structure","command":"bench test --package ./internal/structure"},{"id":"pkg-internal-treetarget","command":"bench test --package ./internal/treetarget"},{"id":"preservation","command":"bench diff","probe":"Compare the frozen baseline and this ticket over its named fixture family. Record facts, specialized exceptions, and unchanged skip classes in the review pickup."}]},{"id":"GF-C4","tickets":["04-migrate-workflow-readers.md"],"verification":[{"id":"pkg-internal-handoff","command":"bench test --package ./internal/handoff"},{"id":"pkg-internal-maps","command":"bench test --package ./internal/maps"},{"id":"pkg-internal-roadmap","command":"bench test --package ./internal/roadmap"},{"id":"pkg-internal-roadmapflow","command":"bench test --package ./internal/roadmapflow"},{"id":"pkg-internal-status","command":"bench test --package ./internal/status"},{"id":"preservation","command":"bench diff","probe":"Compare the frozen baseline and this ticket over its named fixture family. Record facts, specialized exceptions, and unchanged skip classes in the review pickup."}]},{"id":"GF-C5A","tickets":["05a-migrate-cli-adopt-fixtures.md"],"verification":[{"id":"pkg-cmd-bench","command":"bench test --package ./cmd/bench"},{"id":"pkg-internal-adopt","command":"bench test --package ./internal/adopt"},{"id":"pkg-internal-adopt-repairtest","command":"bench test --package ./internal/adopt/repairtest"},{"id":"preservation","command":"bench diff","probe":"Compare the frozen baseline and this ticket over its named fixture family. Record facts, specialized exceptions, and unchanged skip classes in the review pickup."}]},{"id":"GF-C5B","tickets":["05b-migrate-commit-gate-fixtures.md"],"verification":[{"id":"pkg-internal-commit","command":"bench test --package ./internal/commit"},{"id":"pkg-internal-gate","command":"bench test --package ./internal/gate"},{"id":"pkg-internal-gate-authorization","command":"bench test --package ./internal/gate/authorization"},{"id":"pkg-internal-gate-greenmarker","command":"bench test --package ./internal/gate/greenmarker"},{"id":"pkg-internal-gate-prospectiveartifact","command":"bench test --package ./internal/gate/prospectiveartifact"},{"id":"preservation","command":"bench diff","probe":"Compare the frozen baseline and this ticket over its named fixture family. Record facts, specialized exceptions, and unchanged skip classes in the review pickup."}]},{"id":"GF-C5C","tickets":["05c-migrate-preflight-conformance-fixtures.md"],"verification":[{"id":"pkg-internal-conformance","command":"bench test --package ./internal/conformance"},{"id":"pkg-internal-env","command":"bench test --package ./internal/env"},{"id":"pkg-internal-preflight-preflighttest","command":"bench test --package ./internal/preflight/preflighttest"},{"id":"preservation","command":"bench diff","probe":"Compare the frozen baseline and this ticket over its named fixture family. Record facts, specialized exceptions, and unchanged skip classes in the review pickup."}]},{"id":"GF-C5D","tickets":["05d-migrate-workflow-execution-fixtures.md"],"verification":[{"id":"pkg-internal-intent","command":"bench test --package ./internal/intent"},{"id":"pkg-internal-landing","command":"bench test --package ./internal/landing"},{"id":"pkg-internal-responsebound-responseboundtest","command":"bench test --package ./internal/responsebound/responseboundtest"},{"id":"pkg-internal-sessioninspect","command":"bench test --package ./internal/sessioninspect"},{"id":"pkg-internal-shift","command":"bench test --package ./internal/shift"},{"id":"pkg-internal-testreport","command":"bench test --package ./internal/testreport"},{"id":"preservation","command":"bench diff","probe":"Compare the frozen baseline and this ticket over its named fixture family. Record facts, specialized exceptions, and unchanged skip classes in the review pickup."}]},{"id":"GF-C6","tickets":["06-share-landing-fixture-identity.md"],"verification":[{"id":"pkg-internal-worktree","command":"bench test --package ./internal/worktree"},{"id":"preservation","command":"bench diff","probe":"Compare the frozen baseline and this ticket over its named fixture family. Record facts, specialized exceptions, and unchanged skip classes in the review pickup."}]},{"id":"GF-C7","tickets":["07-share-worktree-fixture-identity.md"],"verification":[{"id":"pkg-internal-worktree","command":"bench test --package ./internal/worktree"},{"id":"preservation","command":"bench diff","probe":"Compare the frozen baseline and this ticket over its named fixture family. Record facts, specialized exceptions, and unchanged skip classes in the review pickup."}]},{"id":"GF-C8","tickets":["08-share-system-fixture-identity.md"],"verification":[{"id":"system","command":"bench test --check system"},{"id":"preservation","command":"bench diff","probe":"Compare the frozen baseline and this ticket over its named fixture family. Record facts, specialized exceptions, and unchanged skip classes in the review pickup."}]},{"id":"GF-C9","tickets":["09-own-capability-reasons.md"],"verification":[{"id":"pkg-internal-bounds","command":"bench test --package ./internal/bounds"},{"id":"pkg-internal-capability","command":"bench test --package ./internal/capability"},{"id":"pkg-internal-gittest","command":"bench test --package ./internal/gittest"},{"id":"preservation","command":"bench diff","probe":"Compare the frozen baseline and this ticket over its named fixture family. Record facts, specialized exceptions, and unchanged skip classes in the review pickup."}]},{"id":"GF-C10","tickets":["10-migrate-capability-leaf-callers.md"],"verification":[{"id":"pkg-internal-canary","command":"bench test --package ./internal/canary"},{"id":"pkg-internal-census","command":"bench test --package ./internal/census"},{"id":"pkg-internal-env","command":"bench test --package ./internal/env"},{"id":"pkg-internal-git","command":"bench test --package ./internal/git"},{"id":"pkg-internal-guards","command":"bench test --package ./internal/guards"},{"id":"pkg-internal-learnings","command":"bench test --package ./internal/learnings"},{"id":"pkg-internal-outline","command":"bench test --package ./internal/outline"},{"id":"pkg-internal-prose","command":"bench test --package ./internal/prose"},{"id":"pkg-internal-skillsindex","command":"bench test --package ./internal/skillsindex"},{"id":"pkg-internal-spec","command":"bench test --package ./internal/spec"},{"id":"preservation","command":"bench diff","probe":"Compare the frozen baseline and this ticket over its named fixture family. Record facts, specialized exceptions, and unchanged skip classes in the review pickup."}]},{"id":"GF-C11","tickets":["11-migrate-capability-workflow-callers.md"],"verification":[{"id":"pkg-internal-commit","command":"bench test --package ./internal/commit"},{"id":"pkg-internal-gate","command":"bench test --package ./internal/gate"},{"id":"pkg-internal-landing","command":"bench test --package ./internal/landing"},{"id":"pkg-internal-reviewrecord-recordcmd","command":"bench test --package ./internal/reviewrecord/recordcmd"},{"id":"pkg-internal-status","command":"bench test --package ./internal/status"},{"id":"pkg-internal-worktree","command":"bench test --package ./internal/worktree"},{"id":"preservation","command":"bench diff","probe":"Compare the frozen baseline and this ticket over its named fixture family. Record facts, specialized exceptions, and unchanged skip classes in the review pickup."}]},{"id":"GF-C12","tickets":["12-migrate-capability-conformance-callers.md"],"verification":[{"id":"pkg-cmd-bench","command":"bench test --package ./cmd/bench"},{"id":"conformance","command":"bench test --package ./internal/conformance"},{"id":"preservation","command":"bench diff","probe":"Compare the frozen baseline and this ticket over its named fixture family. Record facts, specialized exceptions, and unchanged skip classes in the review pickup."}]},{"id":"GF-C13","tickets":["13-enforce-fixture-ownership.md"],"verification":[{"id":"git-owner","command":"bench test --check git-plumbing-owner"},{"id":"skip-owner","command":"bench test --check skip-ownership"},{"id":"census","command":"bench test --check ordinary-build-census"},{"id":"conformance","command":"bench test --package ./internal/conformance"},{"id":"preservation","command":"bench diff","probe":"Compare the frozen baseline and this ticket over its named fixture family. Record facts, specialized exceptions, and unchanged skip classes in the review pickup."}]}],"final_verification":[{"id":"coverage","command":"bench coverage --check specs/shared-test-fixtures/spec.md"},{"id":"git-owner","command":"bench test --check git-plumbing-owner"},{"id":"skip-owner","command":"bench test --check skip-ownership"},{"id":"system","command":"bench test --check system"},{"id":"census","command":"bench test --check ordinary-build-census"}]}
```

## Testing decisions

Use real temporary Git repositories for Git fixture behavior. Existing prior art is `internal/gittest/gittest_test.go` and the Git adapter fixture helpers.
Use the capability owner's existing fake or captured skip transport for reason tests. Preserve the write-before-skip tests.
Use the existing conformance fixture builders for ownership checks. Their synthetic Go sources drive the actual registered checker.

Plan evidence now, without implementing the feature to obtain a red. The build records one meaningful red and green per acceptance row.
Shared helper tests attach at their exported contract. Consumer tests assert the caller's scenario rather than repeat helper implementation details.
The completion run includes the whole-project gate. Run the tagged system suite through `bench test --check system` with its supplied `BENCH_KIT`.

### Seam diagram

    test setup -> gittest run/output/commit -> Git argv -> repository facts
    snapshot caller -> gittest snapshot -> file membership and private commit
    unavailable host operation -> capability reason -> record -> skip
    source tree -> existing conformance owners -> duplicate-site diagnostic

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
| --- | --- | --- | --- | --- |
| GF1 | 1 | A failed Git command fails its calling test | planned TestRunFailure in internal/gittest/commands_test.go | Ignoring the process error leaves the negative fixture green. |
| GF2 | 2 | Output preserves trimmed combined output | planned TestOutputContract in internal/gittest/commands_test.go | A changed stream or trim rule changes the observed bytes. |
| GF3 | 3 | OutputBytes preserves raw stdout with trailing whitespace and NUL bytes | planned TestOutputBytesContract in internal/gittest/commands_test.go | A TrimSpace or combined-output substitution changes the expected bytes. |
| GF4 | 4 | A root with spaces and shell metacharacters reaches the intended repository | planned TestCommandArgv in internal/gittest/commands_test.go | Shell interpolation changes the target or creates a sentinel. |
| GF5 | 5 | Missing Git fails the calling test | planned TestMissingGit in internal/gittest/commands_test.go | A skip or empty success cannot satisfy the failure observation. |
| GF6 | 6 | Commit preserves the supplied message and explicit empty-commit option | planned TestCommitContract in internal/gittest/commands_test.go | Dropping an argv element changes the repository result. |
| GF7 | 7 | Commit leaves an unstaged sentinel outside the new commit | planned TestCommitDoesNotStage in internal/gittest/commands_test.go | Implicit add absorbs the sentinel. |
| GF8 | 8 | RepoOnBranch and snapshots use the same canonical identity | planned TestFixtureIdentity in internal/gittest/commands_test.go | A second identity or personal global fallback differs. |
| GF9 | 9 | An explicit identity override survives migration | planned TestFixtureIdentityOverride in internal/gittest/commands_test.go | An unconditional helper override changes the commit author. |
| GF10 | 10 | Every former testrepo Git export caller reaches gittest | planned TestGitFixtureOwnership in internal/conformance/git_fixture_owner_test.go | A retained definition or import edge triggers a diagnostic. |
| GF11 | 11 | The snapshot preserves visible file bytes and modes | `internal/gittest/gittest_test.go` (`TestKitCopyPreservesTheVisibleWorkingTree`) | Omitting the untracked file or executable mode fails the existing fixture. |
| GF12 | 11 | The snapshot preserves live and dangling symlink targets | `internal/gittest/gittest_test.go` (`TestKitCopyPreservesTheVisibleWorkingTree`) | Dereferencing either link changes the link assertion. |
| GF13 | 12 | The snapshot excludes ignored untracked files and source Git metadata | `internal/gittest/gittest_test.go` (`TestKitCopyPreservesTheVisibleWorkingTree`) | Copying the entire source introduces forbidden entries. |
| GF14 | 13 | A special source entry returns an error before a blocking read | planned TestSnapshotSpecialEntry in internal/gittest/snapshot_test.go | A tracked regular file replaced by a FIFO reaches the refusal before a blocking read. |
| GF15 | 13 | A failed snapshot Git command returns its process error | planned TestSnapshotCommandFailure in internal/gittest/snapshot_test.go | Swallowing the error or calling Fatal changes the caller contract. |
| GF16 | 14 | Generic helper callers retain their fixture results | review-owned: each migration ticket compares its named fixture family before and after the change | The same input family produces different repository facts after migration. |
| GF17 | 15 | Worktree Git children retain descendant cleanup and census effects | review-owned: existing worktree journey harness checks plus focused lifecycle evidence | Replacing descendant with an ordinary process loses its recorded effect or cleanup. |
| GF18 | 15 | System Git children retain owner.runAt and its repository budget | review-owned: bench test --check system and ordinary-build-census | A direct process bypass fails the owner or census assertions. |
| GF19 | 16 | Specialized Git probes retain their input, environment, and exit contracts | review-owned: compare each retained inventory symbol and its package test evidence | A convenience wrapper cannot erase an expected failure or custom input. |
| GF20 | 17 | GateFixture keeps its script and input-manifest behavior | review-owned: bench test --package ./internal/testrepo before and after the move | Removing the whole package loses an unrelated fixture API. |
| GF21 | 18 | Unavailable derives the canonical prefix from the class | planned TestUnavailableReason in internal/capability/unavailable_test.go | A caller-specific prefix differs from the contract. |
| GF22 | 19 | Migrated FIFO and symlink sites preserve kind and class | planned TestUnavailableClassification in internal/capability/unavailable_test.go | An environment skip or changed class changes the parsed record. |
| GF23 | 20 | Unavailable writes its record before it skips | planned TestUnavailableWritesBeforeSkip in internal/capability/unavailable_test.go | Moving the write after Skip leaves the transport empty. |
| GF24 | 21 | A planted generic Git runner outside gittest fails the registered check | planned TestGitFixtureOwnershipBites in internal/conformance/git_fixture_owner_test.go | Leaving one private runner cannot pass terminal verification. |
| GF25 | 21 | A planted default identity outside gittest fails the registered check | planned TestGitFixtureIdentityBites in internal/conformance/git_fixture_owner_test.go | A second source remains visible even when the runner delegates. |
| GF26 | 21 | A planted direct FIFO or symlink skip caller fails the registered check | planned TestUnavailableOwnershipBites in internal/conformance/skip_reason_owner_test.go | One unconverted caller cannot pass terminal verification. |
| GF27 | 22 | Source examples in literals do not trigger ownership diagnostics | planned TestFixtureOwnershipLiteralData in internal/conformance/git_fixture_owner_test.go | A text-only search rejects the permitted example. |
| GF28 | 22 | A specialized-site exception does not permit a generic sibling wrapper | planned TestFixtureOwnershipExceptionScope in internal/conformance/git_fixture_owner_test.go | A file-wide allowlist accepts the planted sibling. |
| GF29 | 23 | The same fixture input family has equivalent repository facts before and after migration | review-owned: frozen baseline and candidate differential evidence | Package green alone cannot prove facts the old suite never asserted. |
| GF30 | 24 | Each migration chunk passes without its successor implementation | review-owned: frozen chunk verification plan and three-axis review | A missing future helper is observable at the chunk checkpoint. |
| GF31 | 19 | Every migrated capability caller retains its baseline class | planned TestUnavailableCallerClasses in internal/conformance/skip_reason_owner_test.go | Changing a socket caller from FIFO to symlink fails the independent site expectation. |
| GF32 | 6 | Commit without an empty-commit option fails on an unchanged index | planned TestCommitNoChanges in internal/gittest/commands_test.go | An injected allow-empty option makes the negative fixture succeed. |
| GF33 | 7 | A refused no-change Commit preserves HEAD | planned TestCommitNoChanges in internal/gittest/commands_test.go | Creating an empty commit moves the observed ref. |

### Edge inventory

| edge | disposition and evidence |
| --- | --- |
| Spaces, quotes, glob characters, and shell operators in paths | GF4 uses separate argv operands and a sentinel. |
| Newlines, NUL delimiters, and trailing spaces in output | GF3 observes untrimmed stdout. GF2 retains the existing trimmed API. |
| Missing executable or nonzero exit | GF1 and GF5 require a failed test, not a skip. |
| Empty versus absent file | GF11 and GF29 compare both snapshot inputs. |
| Deleted source entry after enumeration | GF29 preserves the current skip of an absent source entry. |
| FIFO, device, and socket snapshot entries | GF14 refuses nonregular entries before a read. |
| Live and dangling symlinks | GF12 preserves link targets without reading their referents. |
| Personal configuration and deliberate fixture identity | GF8 and GF9 distinguish the default from an explicit subject. |
| Child environment, custom stdin, and expected failure | GF17 to GF19 preserve each specialized probe. |
| Concurrent skip transport | The existing append and collector contracts remain unchanged. GF23 exercises the owner route. |
| Nonroot current directory | GF4 drives an explicit root from a nested current directory. |
| Renamed imports, local variables, and literal examples | GF24 to GF28 grade syntax and exception boundaries. |
| Production lifecycle, ref cleanup, publication, and hooks | Won't handle: existing production callers retain their current owners. |
| Prompting, JSON envelopes, Markdown grammars, and CLI field schemas | Won't handle: this feature changes no such entry surface. |
| Other capability families | Won't handle: existing PID, CPU, privilege, signal, and tool callers retain Capability. |

## Ownership fences

- `cmd/bench/anchor_help_test.go`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `cmd/bench/main_test.go`
- `cmd/bench/otel_hook_seams_test.go`
- `internal/adopt/adopt_test.go`
- `internal/adopt/broker_test.go`
- `internal/adopt/link_hook_test.go`
- `internal/adopt/link_transaction_test.go`
- `internal/adopt/repairtest/session_test.go`
- `internal/adopt/setup_prompt_test.go`
- `internal/bounds/classify_test.go`
- `internal/canary/mutation_test.go`
- `internal/capability/capability_test.go`
- `internal/capability/unavailable.go` (new)
- `internal/capability/unavailable_test.go` (new)
- `internal/census/census_fifo_test.go`
- `internal/census/census_test.go`
- `internal/commit/chain_grammar_test.go`
- `internal/commit/deletion_test.go`
- `internal/commit/dry_run_test.go`
- `internal/commit/landing_test.go`
- `internal/commit/lane_structure_test.go`
- `internal/commit/lane_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/bounds_policy_test.go`
- `internal/conformance/build_contracts_test.go`
- `internal/conformance/checks_test.go`
- `internal/conformance/claude_agent_definitions_test.go`
- `internal/conformance/docs_workflow_checks_test.go`
- `internal/conformance/fixture_bite_test.go`
- `internal/conformance/gate_entry_test.go`
- `internal/conformance/git_fixture_owner_test.go` (new)
- `internal/conformance/git_plumbing_owner_test.go`
- `internal/conformance/guidance_token_sweep_test.go`
- `internal/conformance/handoff_single_source_test.go`
- `internal/conformance/harness_record_test.go`
- `internal/conformance/harness_test.go`
- `internal/conformance/native_workflow_test.go`
- `internal/conformance/ordinary_build_census_test.go`
- `internal/conformance/package_core_checks_test.go`
- `internal/conformance/package_core_diagnostics_test.go`
- `internal/conformance/prose_budget_test.go`
- `internal/conformance/registry/checks.go`
- `internal/conformance/release_probe_fixture_test.go`
- `internal/conformance/skill_description_budgets_test.go`
- `internal/conformance/skip_ownership_test.go`
- `internal/conformance/skip_reason_owner_test.go` (new)
- `internal/conformance/subcommand_routing_table_test.go`
- `internal/conformance/testdata/fixture-skip-classes.json` (new)
- `internal/conformance/tier_test.go`
- `internal/conformance/validity_checks_test.go`
- `internal/consumers/blast_edges_test.go`
- `internal/consumers/blast_test.go`
- `internal/consumers/citation_test.go`
- `internal/consumers/refuse_test.go`
- `internal/diff/command_test.go`
- `internal/diff/compatibility_test.go`
- `internal/diff/explicit_base_test.go`
- `internal/diff/identity_test.go`
- `internal/diff/matrix_test.go`
- `internal/diff/review_base_test.go`
- `internal/diff/source_tip_pair_test.go`
- `internal/env/git_policy_test.go`
- `internal/env/kit_run_test.go`
- `internal/gate/authorization/lane_test.go`
- `internal/gate/cache_env_test.go`
- `internal/gate/gate_prose_staged_root_test.go`
- `internal/gate/gate_prose_staged_test.go`
- `internal/gate/greenmarker/greenmarker_test.go`
- `internal/gate/lane_record_test.go`
- `internal/gate/lane_run_test.go`
- `internal/gate/lane_select_test.go`
- `internal/gate/prospective_owner_test.go`
- `internal/gate/prospectiveartifact/prospectiveartifact_test.go`
- `internal/gate/run_failure_outcomes_test.go`
- `internal/gate/run_outcomes_test.go`
- `internal/gate/verdict_registry_guard_test.go`
- `internal/git/admin_readers_test.go`
- `internal/git/checkout_test.go`
- `internal/git/facts_test.go`
- `internal/git/localnote_test.go`
- `internal/git/push_destination_test.go`
- `internal/git/refs_test.go`
- `internal/git/staged_test.go`
- `internal/git/testhelpers_test.go`
- `internal/git/worktree_admin_enum_test.go`
- `internal/git/worktree_admin_hostile_test.go`
- `internal/gitguard/checker_junction_test.go`
- `internal/gittest/commands.go` (new)
- `internal/gittest/commands_test.go` (new)
- `internal/gittest/gittest.go`
- `internal/gittest/gittest_test.go`
- `internal/gittest/identity.go` (new)
- `internal/gittest/snapshot.go` (new)
- `internal/gittest/snapshot_test.go` (new)
- `internal/guards/guards_test.go`
- `internal/handoff/render_test.go`
- `internal/handoff/sections_test.go`
- `internal/handoff/state_file_test.go`
- `internal/handoff/state_scan_test.go`
- `internal/intent/assignment_lookup_test.go`
- `internal/intent/intent_test.go`
- `internal/intent/worktree_owner_test.go`
- `internal/landing/close_test.go`
- `internal/landing/composition_test.go`
- `internal/landing/landing_helpers_test.go`
- `internal/landing/merge_test.go`
- `internal/landing/state_test.go`
- `internal/learnings/learnings_test.go`
- `internal/maps/freshness_test.go`
- `internal/outline/outline_test.go`
- `internal/poolkey/poolkey_test.go`
- `internal/preflight/preflighttest/fixture.go`
- `internal/preflight/preflighttest/reviewfiles.go`
- `internal/probe/probe_test.go`
- `internal/prose/walk_test.go`
- `internal/responsebound/responseboundtest/checkout.go`
- `internal/reviewrecord/recordcmd/excerpt_test.go`
- `internal/reviewrecord/recordcmd/refusal_test.go`
- `internal/reviewrecord/recordtest/fixture.go`
- `internal/roadmap/learning_test.go`
- `internal/roadmap/retro_test.go`
- `internal/roadmap/roadmap_test.go`
- `internal/roadmapflow/flow_test.go`
- `internal/sessioninspect/sessioninspect_test.go`
- `internal/shift/fault_test.go`
- `internal/shift/refresh_test.go`
- `internal/shift/shift_test.go`
- `internal/skillsindex/command_test.go`
- `internal/skillsindex/skillsindex_reference_test.go`
- `internal/spec/history_command_test.go`
- `internal/spec/history_selected_test.go`
- `internal/spec/spec_test.go`
- `internal/status/handoff_test.go`
- `internal/status/route_test.go`
- `internal/status/status_command_test.go`
- `internal/status/status_counters_test.go`
- `internal/status/status_drain_learnings_test.go`
- `internal/status/status_fixtures_test.go`
- `internal/status/status_gatecache_test.go`
- `internal/status/status_producible_test.go`
- `internal/status/status_render_test.go`
- `internal/status/status_signals_test.go`
- `internal/status/status_spec_count_test.go`
- `internal/status/status_test.go`
- `internal/stophook/stophook_test.go`
- `internal/structure/structure_test.go`
- `internal/systemtest/bench_follow_on_test.go`
- `internal/systemtest/charge_evidence_test.go`
- `internal/systemtest/land_route_test.go`
- `internal/systemtest/otel_gate_test.go`
- `internal/systemtest/otel_verbs_test.go`
- `internal/systemtest/owner_artifact_recovery_test.go`
- `internal/systemtest/owner_land_race_test.go`
- `internal/systemtest/owner_landing_fixture_test.go`
- `internal/systemtest/owner_selection_test.go`
- `internal/systemtest/owner_test.go`
- `internal/systemtest/status_route_converge_test.go`
- `internal/testrepo/working_tree.go`
- `internal/testreport/check_test.go`
- `internal/testreport/selection_test.go`
- `internal/treetarget/identify_test.go`
- `internal/treetarget/run_test.go`
- `internal/worktree/build_test.go`
- `internal/worktree/classifier_shape_test.go`
- `internal/worktree/clean_branch_test.go`
- `internal/worktree/clean_discard_test.go`
- `internal/worktree/clean_discard_transaction_test.go`
- `internal/worktree/clean_landed_hostile_test.go`
- `internal/worktree/clean_operand_test.go`
- `internal/worktree/clean_set_apply_test.go`
- `internal/worktree/clean_set_test.go`
- `internal/worktree/completion_fixture_test.go`
- `internal/worktree/delegated_integration_test.go`
- `internal/worktree/eligibility_test.go`
- `internal/worktree/journey_test.go`
- `internal/worktree/land_bench_home_test.go`
- `internal/worktree/land_broker_notice_test.go`
- `internal/worktree/land_effects_cleanup_test.go`
- `internal/worktree/land_effects_test.go`
- `internal/worktree/land_facts_test.go`
- `internal/worktree/land_fixtures_test.go`
- `internal/worktree/land_flags_test.go`
- `internal/worktree/land_folded_base_test.go`
- `internal/worktree/land_freshness_test.go`
- `internal/worktree/land_journey_test.go`
- `internal/worktree/land_local_capture_test.go`
- `internal/worktree/land_prunes_landed_siblings_test.go`
- `internal/worktree/land_release_refusal_test.go`
- `internal/worktree/land_resume_refusal_test.go`
- `internal/worktree/land_resume_test.go`
- `internal/worktree/land_specless_test.go`
- `internal/worktree/land_surface_test.go`
- `internal/worktree/land_tickets_only_test.go`
- `internal/worktree/lifecycle_acquire_test.go`
- `internal/worktree/lifecycle_facts_test.go`
- `internal/worktree/lifecycle_policy_test.go`
- `internal/worktree/lifecycle_test.go`
- `internal/worktree/live_binary_test.go`
- `internal/worktree/merge_test.go`
- `internal/worktree/ownership_test.go`
- `internal/worktree/parallel_census_test.go`
- `internal/worktree/pool_reclaim_facts_test.go`
- `internal/worktree/reauthorize_test.go`
- `internal/worktree/recovery_retry_test.go`
- `internal/worktree/resume_test.go`
- `internal/worktree/show_test.go`
- `internal/worktree/test_run_test.go`
- `internal/worktree/worktree_test.go`
- `tests/canary/injected-ports/unregistered-port`
- `tests/canary/package-core-guard/reintroduced-bare-skip`
- `reviews/shared-test-fixtures.md` (new)

Only fixture and check changes described above are authorized within these files. Registry co-names satisfy existing fence closure.
Implementation excludes all other specs and their tickets. It does not change roadmap ordering or the active commitment.

## Out of scope

No new production Git adapter, process supervisor, or host-capability vocabulary is introduced.
The future process-group consolidation belongs to FT362. The source-scanner consolidation belongs to FT365.
These are existing separate outcomes, not deferred fragments of FT360. This plan assigns them zero edits and zero gate runs.
No scope cut removes a generic Git fixture or a FIFO or symlink reason from the reviewed source.

## Further notes

### Source verification

The author verified the roadmap premise at `ca339ea83ef401f5edea17a012ba6dc891ec2995` on 2026-10-04.
The author read `gittest.go`, `gittest_test.go`, `testrepo/working_tree.go`, and `testrepo/gate_fixture.go` under `internal/`.
The author also read the capability API, its transport tests, and the Git and skip ownership checks.
The fixture precedents include preflight, reviewrecord, responsebound, worktree journey, and system owner helpers.

`go list` verified the existing gittest, testrepo, and capability import edges. The planned snapshot move removes an edge rather than reversing it.
The author observed snapshot enumeration with spaces, a trailing space, an empty file, and a dangling link.
An untracked FIFO was absent from that list. Replacing a tracked regular file with a FIFO retained its enumerated path.
GF14 uses this reachable input.

No new production import or external dependency is required. ADR 0006 constrains independent expectations to demonstrated omission checks.

### Reader and writer sweep

`inventory.json` is the candidate-site inventory. It names exact paths, enclosing symbols, source lines, and discovery categories.
The sweep covers module Go files under `cmd` and `internal`. A repository-wide hidden-file search also covered docs, scripts, workflows, and canary inputs.
Literal Go source in conformance fixtures remains data. Production Git calls are outside the fixture migration.

The moved export callers are `gittest.KitCopy`, the review-record fixture, and the native-workflow and release-probe conformance fixtures.
Snapshot membership assertions remain in `TestKitCopyPreservesTheVisibleWorkingTree`. No public CLI output schema changes.
The old capability prefixes are caller-authored `FIFOs unavailable`, `symlinks unavailable`, and their singular and filesystem-qualified variants.
The candidate inventory includes direct error-only reasons, device details, and socket details in the same classified family.

The ongoing roadmap-delivery-commitment build writes some future migration consumers. This phase writes none of those consumers.
Before implementation, refresh the inventory and expand affected ticket fences through the normal plan-amendment route.
Shared scorecard changes require composition at phase close. Neither session may overwrite the other's observations.

### Source clauses and coverage

| source clause | coverage |
| --- | --- |
| gittest exports one run, output, and commit form with the identity | GF1 to GF9, GF16, GF24, GF25, GF32, GF33 |
| The git half of testrepo folds into it | GF10 to GF15, GF20 |
| The capability helper owns one skip message for each capability | GF21 to GF23, GF26, GF31, with the source's FIFO and symlink family |
| Private runners and identities duplicate fixture knowledge | GF16 to GF19, GF24, GF25, GF28, GF29 |

### Pre-review proof checklist

- Cited symbols: Current definitions were read. New symbols and tests carry the planned marker in their coverage rows.
- Import edges: go list confirmed the current graph. Remove gittest to testrepo before the move completes.
- Source-row clauses and occurrences: The source-clause table covers FT360's outcome and its quality-survey occurrence.
- Promised field labels: none. The existing capability fields retain their grammar.
- Changed-function callers: inventory.json identifies candidates. The moved-export caller set is listed above.
- Copy survival: GF24, GF25, and GF26 plant a surviving copy through the registered check.
- Rendered-shape readers: The skip parser reads the remainder as Reason. Exact owner transport tests remain valid for the low-level API.
- Restore and copy omissions: GF11 omits a visible untracked file. GF12 omits a dangling link. GF13 plants an excluded entry.
- Test-only cross-package helpers: none. Tests exercise public package exports or their local fixture owners.
- Compared report runs: none. GF29 compares normalized repository facts, excluding commit timestamps and object IDs.
- Executable authority: none changes. Fixture commands use the test host's existing Git authority.

### Flagged additions and delegated decisions

The user delegated routine questions on 2026-10-04. The author keeps the existing seams and preserves the specialized process owners.
OutputBytes is necessary for existing raw-output callers. It is not a new public CLI feature.

The caller-class fixture is an independent preservation expectation. Its recorded class-change red is required before acceptance.

The terminal ownership checks make the source's copy-removal outcome observable. They reuse existing registered check routes.
The spec review precedes slicing, and the slice review follows it, as the user requested.
No implementation author is dispatched during this phase.

### Slice accounting

Each ticket names exact source files and their local helper callers. The inventory also records candidate helper definitions and callers.
The broad source search includes intentional subjects. A candidate requires classification before an edit, and every retained exception needs its exact site.
Checkpoint verification reads only the changed helper and its callers before it broadens to the named package checks.
The final chunk owns the absence of generic copies and the complete caller-class oracle.

The slice review split GF-C5 into GF-C5A, GF-C5B, GF-C5C, and GF-C5D. Each new chunk has its own preservation evidence.
The split preserves the ownership union and acceptance rows. Shared registry writes keep these chunks serial.
