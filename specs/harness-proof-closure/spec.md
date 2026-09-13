# Close the remaining harness proof gaps

Status: staged

Roadmap: FT120

Decision source: `roadmap/FT120.md` (reviewed roadmap item).

Verification log: 1 iteration(s) to accept — the Astra/high author pass added HP19 for focused cleanup. The reviewer approved the package on 2026-09-12.

## Problem

Five live harness gaps let focused or concurrent work report an answer that the
whole gate would reject. Gate fixtures can use undeclared host tools. A package
run can skip live-root conformance. The anchor query locates prose without
grading it, and one changed anchor file does not select its check. A focused
system test can pass and then fail the whole-suite ledger. Concurrent freshness
publishers can interleave one executable, seal, and broker manifest.

FT120 also retains older claims whose implementations now exist. Rebuilding
those paths would duplicate working policy and widen this item without adding a
red-capable outcome.

## Solution

Close each live gap at its existing owner. Give test repositories one gate
fixture writer whose command requests derive the manifest and a private fixture
path. Default the live conformance entry to the harness repository. Make anchor
path inspection evaluate the registry and make registered paths select the docs
check. Apply the system ledger only to an unfiltered suite. Serialize freshness
publication on the stable manifest directory for the full live-state
transaction.

Keep the already-implemented FT120 clauses as explicit premise cuts. The build
does not revive the retired canary phase or replace existing teardown,
preflight, doctor, bounds, or isolated-status behavior.

## User stories

Line: gpt-6-astra / high.
Implementation-line reason: Cross-process freshness publication is the hardest chunk. The roadmap fixes the required outcomes, and existing owners supply strong seams, but interruption and contention need bounded process evidence.
Harder chunks: HP-C5.
The work crosses concurrent publication and executable test-harness seams.

1. As a fixture author, I want one owner for gate scripts and their manifest, so that repeated bodies cannot drift from their closure.
2. As a fixture author, I want each command request to derive its script token and tool row, so that invocation cannot outlive declaration.
3. As a maintainer, I want affected fixtures to run with only their declared commands, so that a developer image cannot hide a CI-only dependency.
4. As a package-test caller, I want unset live conformance to use the current repository, so that package runs see gate defects.
5. As a gate caller, I want an explicit conformance root to keep precedence, so that prospective and fixture trees remain gradeable.
6. As an agent, I want a satisfied registered path query to locate and evaluate its anchors, so that one command proves the live file.
7. As an agent, I want an unsatisfied path query to fail with registry diagnostics, so that a location table cannot conceal a red rule.
8. As an agent, I want an unregistered path query to remain an empty success, so that inspection does not invent policy for ordinary files.
9. As a committer, I want a modified registered anchor file to select the docs-currency check, so that the lane grades the changed policy file.
10. As a committer, I want deletion of a registered anchor file to select the same check, so that removing the subject cannot evade its rule.
11. As a system-test caller, I want a passing `-run` selection to exit green, so that the selected case's verdict is usable.
12. As a system-test caller, I want a failing `-run` selection to retain its own failure, so that the ledger cannot replace the selected diagnostic.
13. As a release reviewer, I want the unfiltered system suite to retain its complete ledger check, so that focused-run support cannot weaken release evidence.
14. As a system-test caller, I want focused runs to keep owner cleanup, so that selecting one case cannot leak its repositories or processes.
15. As a concurrent builder, I want publishers for one broker-manifest directory serialized, so that two live-state transactions cannot interleave.
16. As a runtime caller, I want the final publication triple to describe one build, so that the wrapper cannot authenticate a mixed result.
17. As an interrupted builder, I want rollback to precede lock release, so that recovery cannot undo a waiting publisher's result.
18. As a later builder, I want a dead publisher to leave no persistent lock artifact, so that crash recovery cannot permanently block publication.
19. As a maintainer, I want an unusable manifest directory refused before live mutation, so that lock acquisition cannot weaken target validation.

## Implementation decisions

### Gate fixture ownership

`internal/testrepo` owns a reusable gate-fixture description. It writes the
ordinary script, the optional prospective script, and canonical
`.bench/gate-inputs.json` bytes. A caller can name one body for both gate
paths instead of copying it.

The fixture owner supplies ambient command tokens. Requesting a token records
that command in the manifest, deduplicated and sorted. The same declaration can
construct a private fixture `PATH` that exposes only those commands. Shell
built-ins and absolute test executables are not ambient command declarations.
The migration covers the current gate fixtures that invoke external commands
in the commit, gate, landing, shift, status, system, and worktree test families.
It does not add a shell-language parser.

### Live-root conformance

`TestRootConformance` creates the existing harness before it resolves the
graded root. A non-empty `BENCH_CONFORMANCE_ROOT` remains authoritative.
Otherwise the entry uses the harness root, which already resolves through the
current Git repository. The unset case no longer emits an environment skip.
The ordinary gate continues to provide its explicit root and tier.

### Anchor path evaluation

The anchors package gains one path evaluation over a single classified read.
It returns the registered rows, their physical locations, and the diagnostics
produced by the same matching rules used by group evaluation. Existing group
evaluation derives from this owner; the command does not implement another
anchor matcher.

`bench anchors <path>` retains the `anchors{kind,section,needle,line}`
table. For a registered path, any registry diagnostic follows the rows and
makes the command exit 1. A satisfied registered path exits 0. An unregistered
regular path keeps an empty table and exit 0. Existing wrong-type, unreadable,
out-of-repository, and usage refusals remain unchanged.

Lane selection derives an `anchor-registry` path class from the distinct
`File` values returned by `anchors.Entries()`. Modified and deleted paths
in that class select `docs-currency-workflow` in addition to their existing
classes. The anchor registry remains the only path inventory.

### Focused system verdicts

The system owner reads Go's parsed `test.run` flag after normal flag parsing.
A non-empty value identifies a focused run. `TestMain` always performs its
current cleanup, returns the selected tests' exit status for a focused run, and
does not call the full-suite ledger verifier. With no `test.run` selection,
the current ledger verification and its diagnostics remain mandatory.

This item does not infer focus from raw argument strings and does not change
`-skip`, `-list`, or package selection.

### Publication serialization

The stable serialization identity is the broker-manifest directory. Before
`Publish` reads, backs up, or moves any live executable, seal, or manifest
state, it opens that already-existing directory and takes an exclusive kernel
lock. The lock remains held through publication, rollback, temporary cleanup,
signal teardown, and close.

The directory handle is the lock record. It creates no lock file, and process
exit releases it. Target classification remains fail-closed: a missing,
wrong-type, or unreadable manifest directory fails before any live member
changes. Publishers with different manifest directories do not contend.

Contention tests use separate processes and bounded sentinel handshakes. They
observe entry into live mutation, not elapsed-time guesses. Existing
interruption and rollback helpers remain the precedent.

## Implementation chunks

Each ticket is an independent tracer. The retained implementation session lands
them as serial green checkpoints in the order below. No ticket needs another
ticket's interface.

| chunk / ticket | blocked by | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- | --- |
| HP-C1 / `1-own-gate-fixture-inputs.md` | none | Declared, host-independent gate fixtures | HP1, HP2, HP3 | TestGateFixtureDerivesInputs and migrated fixture packages | no |
| HP-C2 / `2-run-root-conformance-by-default.md` | none | Package runs grade their live repository by default | HP4, HP5 | TestRootConformanceDefaultsToHarnessRoot and TestHarnessDefaultsToCurrentGitRoot | no |
| HP-C3 / `3-evaluate-anchor-paths.md` | none | One-file anchor proof and lane attachment | HP6, HP7, HP8, HP9, HP10 | anchors package, TestAnchors reports cases, and TestSelectLaneByClass | no |
| HP-C4 / `4-preserve-focused-system-verdicts.md` | none | Focused system runs return selected verdicts | HP11, HP12, HP13, HP19 | TestFocusedRunUsesSelectedVerdict, TestFocusedRunStillCleansOwnerRoot, and the full system suite | no |
| HP-C5 / `5-serialize-freshness-publication.md` | none | Cross-process atomic publication | HP14, HP15, HP16, HP17, HP18 | TestPublishSerializesManifestDirectory and existing publication interruption cases | yes |

The coverage map supplies the complete predicates for each checkpoint.
After each chunk, review freezes its predecessor and current tips for Standards,
Spec, and Coverage review. The successor starts after accepted findings have
repair coverage. Execution-plan changes follow `.bench/BENCH.md`.

## Completion plan

The fenced plan names the focused checks and the omission each chunk must make
red. Final verification follows HP-C5 and precedes landing.

```bench-completion-plan
{
  "version": 2,
  "execution": {
    "mode": "delegate",
    "run_id": "ft120-20260912-01",
    "orchestrator_session": "/root",
    "author_limit": 3,
    "assignments": {
      "1-own-gate-fixture-inputs.md": [
        {
          "session": "/root/ft120_c1",
          "assignment": "322e8f5a65d7e6587716301507c563f7",
          "model": "gpt-6-astra",
          "effort": "high",
          "source": "8c6f68da3e159e93d6af733938373222dda3fd32",
          "native_ref": "collaboration:/root/ft120_c1"
        }
      ],
      "2-run-root-conformance-by-default.md": [
        {
          "session": "/root/ft120_c2",
          "assignment": "cf3bf7ea402a0409197271434d1d6eb5",
          "model": "gpt-6-astra",
          "effort": "high",
          "source": "8c6f68da3e159e93d6af733938373222dda3fd32",
          "native_ref": "collaboration:/root/ft120_c2"
        }
      ],
      "3-evaluate-anchor-paths.md": [
        {
          "session": "/root/ft120_c3",
          "assignment": "cd1dc7c8eded6027ea11742134e06957",
          "model": "gpt-6-astra",
          "effort": "high",
          "source": "8c6f68da3e159e93d6af733938373222dda3fd32",
          "native_ref": "collaboration:/root/ft120_c3"
        }
      ],
      "4-preserve-focused-system-verdicts.md": [
        {
          "session": "/root/ft120_c4",
          "assignment": "2705b151c3e164ab416e30dcd65e7217",
          "model": "gpt-6-astra",
          "effort": "high",
          "source": "8c6f68da3e159e93d6af733938373222dda3fd32",
          "native_ref": "collaboration:/root/ft120_c4"
        }
      ],
      "5-serialize-freshness-publication.md": [
        {
          "session": "/root/ft120_c5",
          "assignment": "6c03260aa82b132c1e986acf9d24c5ad",
          "model": "gpt-6-astra",
          "effort": "high",
          "source": "8c6f68da3e159e93d6af733938373222dda3fd32",
          "native_ref": "collaboration:/root/ft120_c5"
        }
      ]
    }
  },
  "chunks": [
    {
      "id": "HP-C1",
      "tickets": ["1-own-gate-fixture-inputs.md"],
      "verification": [
        {
          "id": "fixture-owner",
          "command": "bench test --package ./internal/testrepo",
          "probe": "write an ambient command token without recording it in the manifest",
          "ticket": "1-own-gate-fixture-inputs.md"
        },
        {
          "id": "fixture-consumers",
          "command": "bench test --package ./internal/...",
          "ticket": "1-own-gate-fixture-inputs.md"
        }
      ]
    },
    {
      "id": "HP-C2",
      "tickets": ["2-run-root-conformance-by-default.md"],
      "verification": [
        {
          "id": "root-entry",
          "command": "bench test --package ./internal/conformance --run 'TestRootConformance|TestHarnessDefaultsToCurrentGitRoot'",
          "probe": "restore the unset-root environment skip",
          "ticket": "2-run-root-conformance-by-default.md"
        }
      ]
    },
    {
      "id": "HP-C3",
      "tickets": ["3-evaluate-anchor-paths.md"],
      "verification": [
        {
          "id": "anchor-owner",
          "command": "bench test --package ./internal/anchors",
          "probe": "return locations without the path's registry diagnostics",
          "ticket": "3-evaluate-anchor-paths.md"
        },
        {
          "id": "anchor-command",
          "command": "bench test --package ./cmd/bench --run TestAnchors",
          "ticket": "3-evaluate-anchor-paths.md"
        },
        {
          "id": "anchor-lane",
          "command": "bench test --package ./internal/gate --run TestSelectLaneByClass",
          "ticket": "3-evaluate-anchor-paths.md"
        },
        {
          "id": "docs-check",
          "command": "bench test --check docs-currency-workflow",
          "ticket": "3-evaluate-anchor-paths.md"
        }
      ]
    },
    {
      "id": "HP-C4",
      "tickets": ["4-preserve-focused-system-verdicts.md"],
      "verification": [
        {
          "id": "focused-system",
          "command": "go test -trimpath -count=1 -tags=system ./internal/systemtest -run '^TestFocusedRun'",
          "probe": "call the whole-suite ledger verifier after a focused passing child",
          "ticket": "4-preserve-focused-system-verdicts.md"
        },
        {
          "id": "full-system",
          "command": "bench test --check system",
          "ticket": "4-preserve-focused-system-verdicts.md"
        }
      ]
    },
    {
      "id": "HP-C5",
      "tickets": ["5-serialize-freshness-publication.md"],
      "verification": [
        {
          "id": "publication",
          "command": "bench test --package ./internal/freshness",
          "probe": "release the directory lock before rollback and temporary cleanup finish",
          "ticket": "5-serialize-freshness-publication.md"
        }
      ]
    }
  ],
  "final_verification": [
    {
      "id": "coverage",
      "command": "bench coverage --check harness-proof-closure"
    },
    {
      "id": "acceptance",
      "command": "bench test --package ./..."
    },
    {
      "id": "integration",
      "command": "bench test --check system"
    }
  ]
}
```

## Testing decisions

Tests attach at the existing public owners. The test-repository helper is
in-process and local-substitutable. Root and anchor tests use isolated Git
repositories. Focused-system and publication contention evidence re-executes
the package test binary as a local process. No remote service enters the suite.

Planned names below identify required behavior. Existing named tests remain
valid evidence where the map cites them.

### Seam diagram

```text
fixture declaration -> [ testrepo gate owner ] -> scripts + manifest + private PATH
                              ^ tests omit one tool and run the fixture

package invocation -> [ conformance root entry ] -> evaluated repository diagnostics
                              ^ tests vary explicit and absent root environment

path / changed path -> [ anchors registry evaluation ] -> rows + diagnostics / lane checks
                              ^ tests mutate an isolated registered file

go test -run -> [ system TestMain owner ] -> selected verdict + cleanup
                              ^ child-process tests compare focused and full invocations

two publishers -> [ manifest-directory lock + publication ] -> one coherent live triple
                              ^ process tests hold, interrupt, release, and verify
```

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
| --- | --- | --- | --- | --- |
| HP1 | 1 | One fixture declaration writes the requested gate paths and one canonical input manifest | review-owned: planned `TestGateFixtureDerivesInputs` in the new gate-fixture test file | Separate writers can reproduce the script-with-stale-manifest defect |
| HP2 | 2 | Every ambient command token in a migrated script appears once in the sorted manifest tools | review-owned: planned `TestGateFixtureDerivesToolClosure` in the new gate-fixture test file | A hand-written tool list can omit the command while the script still runs on a rich host |
| HP3 | 3 | The migrated fixtures pass with a private path that exposes only declared commands | `internal/landing/completion_evidence_test.go` (`TestLandingCompletionEvidence`), `internal/worktree/land_journey_test.go` (`TestLandCommandRefusesPostGateUnknownIgnoredMutation`), and affected package tests | An undeclared `grep`, `git`, `sed`, `cat`, `cp`, `sleep`, `mkdir`, `env`, or `bash` cannot resolve from the developer's ambient path |
| HP4 | 4 | With no conformance-root environment, the entry evaluates the harness Git root and emits no unset-root skip | `internal/conformance/gate_entry_test.go` (`TestRootConformance`) | Restoring the current early skip leaves the injected defect unreported |
| HP5 | 5 | A non-empty conformance-root environment wins over the harness root | `internal/conformance/harness_test.go` (`TestHarnessUsesBenchConformanceRootAsGradedRoot`, `TestHarnessDefaultsToCurrentGitRoot`) | Always using the current repository would grade the wrong prospective fixture |
| HP6 | 6 | A satisfied registered path prints ordered locations and exits 0 | `cmd/bench/anchor_help_test.go` (`TestAnchorsReportsNeedleLines`) | A command that only evaluates and loses locations breaks the current query contract |
| HP7 | 7 | A missing required anchor or present forbidden anchor prints its registry diagnostic and exits 1 | `cmd/bench/anchor_help_test.go` (`TestAnchorsReportsAbsentNeedles`) | The current zero line with exit 0 passes this defect |
| HP8 | 8 | An unregistered regular path prints an empty anchors table and exits 0 | review-owned: planned `TestAnchorsLeavesUnregisteredPathEmpty` in the anchor command test file | Treating every file as governed invents failures outside the registry |
| HP9 | 9 | A modified registered path selects `docs-currency-workflow` beside its existing classes | `internal/gate/lane_select_test.go` (`TestSelectLaneByClass`) | Markdown-only classification currently selects prose and misses the anchor check |
| HP10 | 10 | Deletion of a registered path still selects `docs-currency-workflow` | `internal/gate/lane_select_test.go` (`TestSelectLaneByClass`) | Classifying only files that still exist lets deletion evade the check |
| HP11 | 11 | A child system invocation with a passing non-empty `test.run` exits 0 and prints no ledger diagnostic | review-owned: planned `TestFocusedRunUsesSelectedVerdict` in the owner-selection test file | The current selected PASS followed by owner.verify exit 1 fails this row |
| HP12 | 12 | A child system invocation with a failing non-empty `test.run` stays nonzero and retains the selected failure | review-owned: planned failing-helper branch of `TestFocusedRunUsesSelectedVerdict` | Replacing the failure with a missing-ledger trailer hides the chosen test's verdict |
| HP13 | 13 | An unfiltered system suite still verifies repositories, executable observations, and terminal outcomes | `internal/systemtest/owner_test.go` (`TestMain`) | Skipping the ledger for every invocation would make focused tests green by weakening release proof |
| HP19 | 14 | A focused child removes its owner root and reports any cleanup failure | review-owned: planned `TestFocusedRunStillCleansOwnerRoot` in the owner-selection test file | Returning immediately after the selected verdict leaks the focused run's owned repositories |
| HP14 | 15 | A second process targeting one manifest directory cannot enter live mutation while the first holds the lock | review-owned: planned `TestPublishSerializesManifestDirectory` in the new publication-lock test file | Two publishers that take independent or late locks can both cross the mutation sentinel |
| HP15 | 16 | After two successful contenders, executable bytes, seal digest, manifest digest, path, and version belong to one completed publisher | review-owned: final-state assertions in planned `TestPublishSerializesManifestDirectory` | A mixed triple can pass if the test checks only that both calls returned |
| HP16 | 17 | Interrupting the holder before seal promotion restores its prior triple before the waiter publishes and completes | `internal/freshness/freshness_publish_test.go` (`TestPublishRestoresPriorPairWhenInterruptedBeforeItsSealLands`) plus the planned waiter case | Releasing before rollback lets the holder undo the waiter's installed result |
| HP17 | 18 | Terminating a lock holder leaves no lock file and a later publisher acquires the directory lock | review-owned: planned `TestPublicationLockDiesWithProcess` in the new publication-lock test file | A persistent lock record can remain after the owner process is gone |
| HP18 | 19 | A missing, wrong-type, or unreadable manifest directory returns an error before executable, seal, or manifest bytes change | `internal/freshness/freshness_publish_test.go` (`TestPublishAndVerifyRefuseEmptyExecutables/no_manifest_directory`) | Lock acquisition after executable backup or install can damage the live set before refusal |

### Edge inventory

HP1–HP3 cover ordinary, prospective, and shared bodies. They also cover
duplicate tool requests, an empty tool set, missing host resolution, and the
current external-command fixture census. The private path is test-only. It does
not change the production gate environment.

HP4–HP5 cover an unset, empty, and non-empty conformance-root environment. A
package invocation outside a Git repository retains the harness's existing
repository refusal.

HP6–HP10 cover required, forbidden, and section-scoped anchors. They also cover
missing or duplicated sections, missing files, refused subjects, changed
paths, deleted paths, and an unregistered file.

HP11–HP13 and HP19 cover passing and failing focused selections, focused
cleanup, and an unfiltered run.
The parsed `test.run` value, not argument spelling, decides the mode.

HP14–HP18 cover two successful publishers, a holder interrupted before seal
promotion, a dead holder, missing and refused manifest directories, and
existing rollback residue. Sentinels have explicit test deadlines.

Won't handle: parsing arbitrary shell text to discover commands — migrated fixtures obtain ambient command tokens from the owner, which keeps the manifest single-sourced.
Won't handle: `-skip` or `-list` as focused-system modes — FT120 names `-run`, and the unfiltered release route uses neither flag.
Won't handle: serialization across different broker-manifest directories — the production wrapper has one manifest identity, and distinct directories publish distinct triples.
Won't handle: a portable non-Unix publication lock — the publisher and its existing signal/process tests already use Unix system calls.

## Ownership fences

- `internal/testrepo/gate_fixture.go` (new)
- `internal/testrepo/gate_fixture_test.go` (new)
- `internal/commit/format_test.go`
- `internal/gate/run_outcomes_test.go`
- `internal/landing/completion_evidence_test.go`
- `internal/shift/fault_test.go`
- `internal/shift/shift_test.go`
- `internal/status/status_producible_test.go`
- `internal/systemtest/owner_landing_fixture_test.go`
- `internal/worktree/delegated_integration_test.go`
- `internal/worktree/land_fixtures_test.go`
- `internal/worktree/land_journey_test.go`
- `internal/conformance/gate_entry_test.go`
- `internal/conformance/harness_test.go`
- `internal/conformance/registry/registry.go`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `internal/anchors/registry.go`
- `internal/anchors/anchor_harness_diagnostics_test.go`
- `cmd/bench/anchors_command.go`
- `cmd/bench/anchor_help_test.go`
- `internal/gate/lane_select.go`
- `internal/gate/lane_select_test.go`
- `projects/benchkit.md`
- `tests/canary/guidance-prose-budgets/over-budget-skill/`
- `tests/canary/line-routing/line-binding-prose-drift/`
- `tests/canary/skill-description-budgets/budget-table-missing/`
- `tests/canary/skill-description-budgets/description-folded/`
- `tests/canary/skill-description-budgets/description-missing/`
- `tests/canary/skill-description-budgets/over-budget-command/`
- `tests/canary/skill-description-budgets/over-budget-description/`
- `tests/canary/workflow-guidance-anchors/benchkit-hostile-input-heading/`
- `tests/canary/workflow-guidance-anchors/benchkit-review-round-owner/`
- `tests/canary/workflow-guidance-anchors/benchkit-review-round-routing/`
- `tests/canary/workflow-guidance-anchors/benchkit-spec-ownership/`
- `tests/canary/workflow-guidance-anchors/benchkit-system-suite-route/`
- `internal/systemtest/owner_test.go`
- `internal/systemtest/owner_selection_test.go`
- `internal/freshness/freshness_publish.go`
- `internal/freshness/publication_lock.go` (new)
- `internal/freshness/freshness_publish_test.go`
- `internal/freshness/publication_lock_test.go` (new)
- `specs/harness-proof-closure/`
- `reviews/harness-proof-closure.md`
- `roadmap/FT120.md`
- `CHANGELOG.md`
- `capture/retros/harness-proof-closure.md` (new)

Reviewer disposition: approved on 2026-09-12. The fence is the union of
ticket writes, the review pickup, and phase-close records. `.bench/BENCH.md`
governs approved in-scope execution-plan changes.

## Ticket graph

| ticket | blocked by | delivered outcome |
| --- | --- | --- |
| [1. Own gate fixture inputs](tickets/1-own-gate-fixture-inputs.md) | none | Declared, host-independent gate fixtures |
| [2. Run root conformance by default](tickets/2-run-root-conformance-by-default.md) | none | Package runs grade their live repository by default |
| [3. Evaluate anchor paths](tickets/3-evaluate-anchor-paths.md) | none | One-file anchor proof and lane attachment |
| [4. Preserve focused system verdicts](tickets/4-preserve-focused-system-verdicts.md) | none | Focused system runs return selected verdicts |
| [5. Serialize freshness publication](tickets/5-serialize-freshness-publication.md) | none | Cross-process atomic publication |

## Out of scope

- Enforce ambient-command declarations for arbitrary consumer gate scripts — 9 edits, 2 gate runs. This spec changes test fixture construction, not the consumer gate schema.
- Add focused semantics for Go test filters other than `-run` — 4 edits, 1 gate run. Each filter needs an explicit ledger policy.
- Add a cross-platform publication-lock abstraction — 7 edits, 2 gate runs. The shipped publisher currently targets the Unix process and signal model.
- Add a multi-path anchor query — 6 edits, 2 gate runs. One changed path already gives the lane and agent a complete proof.

## Further notes

The FT120 premise audit found these clauses already implemented or retired:

| roadmap clause | current tree evidence | disposition |
| --- | --- | --- |
| Bounded release waits and process-group teardown | marker-wait conformance plus the subprocess cancellation owner | Retain; no new row |
| Per-fixture `BENCH_CANARY_PHASE` isolation | the canary execution phase no longer exists | Retired; do not restore |
| Nested preflight special-file refusal | TestSpecialTicketEntryRefused and current classified reads | Retain; no new row |
| Doctor launch-versus-child attribution | the harness separates process-start errors from ExitError status | Retain; no new row |
| Typed `bounds.ClassifyNoFollow` classes | typed states and classification tests | Retain; no new row |
| Isolated live-status byte comparisons | isolated command and handoff fixture repositories | Retain; no new row |
| Ordinary gate live-root propagation | rootConformanceEnv supplies the root and tier | Retain; HP4 closes package-run default only |

The live clauses map to rows as follows:

| source clause or occurrence | rows |
| --- | --- |
| Live-tree pass defaults to the repository root when the environment is unset | HP4, HP5 |
| System `TestMain` must not false-red a `-run` subset | HP11, HP12, HP13, HP19 |
| Gate fixtures declare ambient tools and share one owner | HP1, HP2, HP3 |
| A non-gate anchor handle evaluates one live path and its lane reaches docs currency | HP6, HP7, HP8, HP9, HP10 |
| Concurrent freshness publishers cannot interleave the live manifest triple | HP14, HP15, HP16, HP17, HP18 |

Flagged additions beyond the decision source: HP3 supplies a private fixture
path, HP8 preserves unregistered queries, HP19 preserves cleanup, and HP17
excludes persistent lock artifacts. Each row blocks the cheapest implementation
that would satisfy the positive FT120 sentence while weakening current behavior.

Pre-review proof checklist:

- `Cited symbols`: `TestRootConformance`, `NewHarness`, `anchors.Entries`, `EvaluateGroup`, `anchorsCommand`, `TestMain`, `systemOwner.verify`, `freshness.Publish`, and `beginPublication` resolve in the tree.
- `Import edges`: `internal/gate` can import `internal/anchors`; anchors does not import gate. Test packages can import `internal/testrepo`.
- `Source-row clauses and occurrences`: every still-live FT120 clause maps above; every closed clause has a premise disposition.
- `Promised field labels`: the anchor table retains `kind`, `section`, `needle`, and `line`; the gate manifest retains `tools`.
- `Changed-function callers`: the production freshness caller remains `cmd/bench/freshness_publish.go`; anchor group callers derive through the shared evaluator.
- `Copy survival`: gate script and manifest copies in migrated fixtures are removed in favor of the fixture owner.
