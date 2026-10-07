# Shell argument ownership

Status: staged

Decision source: named reviewed artifact `decisions/architecture-primitives.md`, resolved tickets `3.md` and `4.md` (2026-10-06).

Verification log: 2 spec iterations and 2 ticket iterations to accept — spec acceptance preceded slicing and separate independent ticket review.

## Problem

AXI and freshness define shell quote algorithms beside the sanitize owner.
AXI prints safe values bare, while sanitize always quotes values.
The differing spellings reach direct command lines and AXI executable actions.
Sources: `internal/axi/action.go:284`, `internal/sanitize/sanitize.go:109`, and `internal/freshness/freshness_verify.go:124`.

AXI shellSafeToken accepts an empty string.
Its renderKnownArgument caller refuses empty first.
The empty-helper result is not a reproduced reachable command defect.
Sources: `internal/axi/action.go:291`, `internal/axi/action.go:298`, and `internal/axi/action.go:305`.

## Solution

Use the existing sanitize.ShellQuote owner and its always-quoted contract.
Migrate callers and remove the independent production algorithms.
Preserve each caller control-byte policy and action validity check.
Approve exact output changes through domain fixtures.
Keep POSIX-shell argv evidence separate from contractual output-byte evidence.

## User stories

Line: gpt-5.6-sol / high.
Implementation-line reason: AXI validation and rendered-output consumers make the quote-policy migration harder than a helper substitution.
Harder chunks: S-A, S-E.
Author line: gpt-6.1-sol / high / one draft plus at most two review repair rounds.

### Argument and caller behavior

1. As a command reader, I want to retain an empty argument, so that the shell receives the original argument count.
2. As a command reader, I want to retain an embedded quote, so that the shell receives the original bytes.
3. As a command reader, I want to retain shell metacharacters, so that printed arguments cannot execute substitutions.
4. As a command reader, I want to receive one quote policy, so that argument spelling has one owner.
5. As an AXI reader, I want to retain every known argument, so that follow-up commands preserve the declared invocation.
6. As an AXI reader, I want to retain future-input slots, so that the action distinguishes known and unknown values.
7. As an AXI reader, I want to retain empty-value refusal, so that consolidation preserves the action validity contract.
8. As an AXI reader, I want to retain control-byte refusal, so that unsupported disclosure values remain absent.
9. As an AXI reader, I want to retain accepted layout bytes, so that TOON escapes remain reversible.
10. As an AXI reader, I want to retain command validity checks, so that quote spelling cannot weaken action authority.
11. As an AXI reader, I want to retain action deduplication, so that duplicate actions do not multiply.
12. As a harness reader, I want to retain phase arguments, so that shell policy does not alter harness syntax.
13. As a spec-history reader, I want to replay the selected detail command, so that the recovered command reads the same history.
14. As a consumers reader, I want to replay the citation command, so that the citation preserves every supplied operand.
15. As a worktree reader, I want to replay the build repair, so that the quoted label reaches the correct assignment.
16. As a worktree reader, I want to replay a recovery command, so that the next line names the original worktree.
17. As an evidence reader, I want to replay the next page, so that the command retains its evidence operands.
18. As an operator, I want to replay the rebuild command, so that the command retains its original paths.
19. As an anchors reader, I want to retain the exact rendered invocation, so that the action reaches the same command.
20. As a coverage reader, I want to retain the exact rendered invocation, so that the action reaches the same command.
21. As a consumers query reader, I want to retain the exact rendered invocation, so that the action reaches the same command.
22. As a consumers blast reader, I want to retain the exact rendered invocation, so that the action reaches the same command.
23. As a guards reader, I want to retain the exact rendered invocation, so that the action reaches the same command.
24. As a maps reader, I want to retain the exact rendered invocation, so that the action reaches the same command.
25. As a gate retry reader, I want to retain the exact rendered invocation, so that the action reaches the same command.
26. As a preflight retry reader, I want to retain the exact rendered invocation, so that the action reaches the same command.
27. As a worktree list reader, I want to retain the exact rendered invocation, so that the action reaches the same command.
28. As a worktree selected list reader, I want to retain the exact rendered invocation, so that the action reaches the same command.
29. As a worktree explicit clean reader, I want to retain the exact rendered invocation, so that the action reaches the same command.
30. As a worktree clean refusal reader, I want to retain the exact rendered invocation, so that the action reaches the same command.
31. As a worktree landed clean reader, I want to retain the exact rendered invocation, so that the action reaches the same command.
32. As a worktree unclaimed clean reader, I want to retain the exact rendered invocation, so that the action reaches the same command.
33. As a worktree pool reclaim reader, I want to retain the exact rendered invocation, so that the action reaches the same command.
34. As a maintainer, I want to remove independent production algorithms, so that quote policy cannot drift between packages.
35. As a maintainer, I want to retain line-output refusals, so that quoting cannot grant sink permission.
36. As a maintainer, I want to retain unchanged contracts, so that the approved byte delta stays bounded.
37. As a maintainer, I want to use the owner in fixture scripts, so that test fixture quote policy has one source.
38. As a maintainer, I want to prove the caller composition, so that a shared-helper test cannot hide a bypass.
39. As a maintainer, I want to update affected output readers, so that quoting preserves their original behavioral assertions.

## Implementation decisions

The owner stays `sanitize.ShellQuote(value string) string`.
It always encloses a value in single quotes.
It spells an embedded quote by closing, escaping, and reopening the quoted value.
An empty value becomes an empty quoted pair.
Do not add an option for bare safe tokens.

Remove axi.ShellQuote after its direct callers and renderKnownArgument migrate.
Remove freshness.shellQuote after RebuildAction migrates.
Keep AXI shellSafeToken for executable-name validation only.
Validation compares the declared value before quote output changes.
In particular, the invalid first value bench remains refused.
The quote owner does not authorize a command.

The direct quote owner preserves control bytes as values.
AXI retains its current accepted tab, newline, and carriage-return values through TOON.
AXI refuses its other unsupported control runes.
AXI still refuses empty known values and angle-bracket placeholders.
Line-structured callers retain LineSafe checks or the existing whole-line escape.
These callers must not reuse AXI layout-byte permission as line permission.

A quote owner is distinct from a complete command renderer.
Existing literal command tokens and shell operators keep their renderer-owned spellings.
ExecutableInvocation quotes the accepted KnownArgument values it already feeds through renderKnownArgument.
InspectFull and RetryDiff retain their current validated token renderers.
HarnessPhase and HarnessPhaseOn retain their harness syntax.
No new argument, flag, next action, schema field, or exit code belongs to this migration.

### Current direct caller inventory

Sources pin the tree at `a382f4848d432815968f044820fad593e856e297`.
A caller or renderer change invalidates its inventory row.

| current function | call site | approved byte delta | row |
|---|---|---|---|
| axi.renderKnownArgument | `internal/axi/action.go:291` | Accepted known values become always quoted. | S05-S12 |
| spec.historyDetail | `internal/spec/history_selected.go:121` | The operand uses the selected always-quoted spelling. | S13 |
| consumers.citation.cmd | `internal/consumers/citation.go:58` | Each argv value becomes always quoted. | S14 |
| worktree.PrintTreeBuildRefusal | `internal/worktree/tree_target.go:78` | The assignment label becomes always quoted. | S15 |
| worktree.recoveryRoute.line | `internal/worktree/list.go:252` | The optional path becomes always quoted. | S16 |
| evidencecmd.evidenceInvocation | `internal/preflight/evidencecmd/evidence.go:233` | Each rendered argv value becomes always quoted. | S17 |
| freshness.RebuildAction | `internal/freshness/freshness_verify.go:121` | No approved byte delta. The owner changes. | S18 |

The AXI function also has test callers in identifier_operand_test.go and unlanded_route_test.go.
The static census includes those callers.
Existing sanitize callers retain their behavior because the owner implementation already exists.

### Rendered AXI producer inventory

This inventory includes executable actions beyond direct ShellQuote calls.
Each producer row requires exact approved command output and shell argv evidence.
For a future-input placeholder, compare its bytes without executing the placeholder.

| producer | source | current command family | row |
|---|---|---|---|
| anchors | `cmd/bench/anchors_command.go:86` | `anchors <file>` | S19 |
| coverage | `internal/coverage/coverage.go:587` | `coverage --check <spec>` | S20 |
| consumers query | `internal/consumers/query.go:141` | `consumers <qualified-symbol> [scope] [--full]` | S21 |
| consumers blast | `internal/consumers/blast.go:320` | `consumers [original-args] --full` | S22 |
| guards | `internal/guards/guards.go:323` | `link` | S23 |
| maps | `internal/maps/maps.go:258` | `maps --template` | S24 |
| gate retry | `internal/gate/gate.go:332` | `gate --fresh` | S25 |
| preflight retry | `internal/preflight/command.go:214` | `preflight [original-args]` | S26 |
| worktree list | `internal/worktree/list.go:152` | `worktree path, exec, release, or clean` | S27 |
| worktree selected list | `internal/worktree/list_selected.go:90` | `worktree list` | S28 |
| worktree explicit clean | `internal/worktree/clean_set.go:296` | `worktree clean [selection] --apply <fingerprint>` | S29 |
| worktree clean refusal | `internal/worktree/clean_set_apply.go:163` | `worktree clean [selection]` | S30 |
| worktree landed clean | `internal/worktree/clean_landed.go:248` | `worktree clean --landed or worktree exec [assignment] -- [command]` | S31 |
| worktree unclaimed clean | `internal/worktree/clean_unclaimed.go:145` | `worktree clean [branch selection] --apply <fingerprint>` | S32 |
| worktree pool reclaim | `internal/worktree/pool_reclaim.go:318` | `worktree reclaim [--apply <fingerprint>]` | S33 |

### Output readers and fixtures

The following files read changed rendered bytes or related output contracts.
The inventory separates exact literals from unaffected envelope and semantic assertions.
Keep every pre-existing assertion.
Update only the approved quote spellings in exact command expectations.

| fixture group | files | disposition |
|---|---|---|
| AXI owner | `internal/axi/action_test.go` | Update executable-action command literals. Keep refusal, harness, order, and count assertions. |
| Consumers | `citation_test.go`, `candidates_test.go`, `loader_test.go`, `blast_edges_test.go`, `query_test.go`, `command_test.go`, `blast_test.go`, `refuse_test.go` under internal/consumers | Update replay and executable-help command literals. Keep citation hash and query behavior. |
| Worktree | `list_actions_test.go`, `identifier_operand_test.go`, `landed_test.go`, `unlanded_route_test.go`, `path_identifier_test.go` under internal/worktree | Update executable help and direct recovery operands. Keep identity and clean policy. |
| Maps, guards, coverage | `internal/maps/maps_command_test.go`, `internal/maps/freshness_test.go`, `internal/guards/guards_test.go`, `internal/coverage/coverage_command_test.go`, `internal/coverage/coverage_schema_test.go` | Update executable-help command literals. Keep empty help and harness-phase assertions. |
| Command facade | `cmd/bench/anchor_help_test.go`, `cmd/bench/anchors_dir_test.go`, `cmd/bench/main_test.go`, `cmd/bench/command_registry_test.go` | Update changed anchors command literals only. Registry and help grammar remain unchanged. |
| Gate and preflight | `internal/gate/run_outcomes_test.go`, `internal/preflight/explicit_base_test.go` | Update retry command literals. Keep retries and refusal semantics. |
| Selected history | `internal/spec/history_selected_test.go` | Preserve the real-shell history route test and add exact quote fixtures. |
| Evidence | `internal/preflight/evidencecmd/evidence_command_test.go`, `evidence_modes_test.go`, `evidence_summary_test.go`, `evidence_review_test.go`, `evidence_cleanup_test.go` in that package | Update next-command fixtures and preserve evidence selection. |
| Freshness | `internal/freshness/freshness_verify_test.go` | Preserve the executable rebuild fixture and its exact observed arguments. |
| Conformance | `internal/conformance/axi_query_registry_test.go` | Read-only exclusion. It grades inventory and envelopes, not changed command bytes. |

### Changed producer-to-assertion closure

These eleven additional readers belong to the approved output change.
Each path is an exact fixture fence, with its current assertion source pinned below.
The first three read PrintTreeBuildRefusal directly or through the target command.
The remaining readers observe AXI KnownArgument rendering through clean or reclaim producers.
S39 requires all listed readers to retain their behavioral assertions.

| affected fixture and existing seam | producer / row | approved assertion update and preserved purpose |
|---|---|---|
| `internal/treetarget/build_test.go:95`, TestRunKitWorktreeBuild | PrintTreeBuildRefusal / S15 | Quote alpha in TT44-TT46 expected repair lines. Keep hostile-label bytes, child non-start, environment, argv, and cwd assertions. |
| `internal/worktree/target_refusal_test.go:19`, TestTargetRefusalEscapesOnlyAnUnsafeLine | PrintTreeBuildRefusal / S15 | Change the plain-label repair to `bench worktree build 'alpha'`. Keep every control escape and single-line expectation. |
| `internal/systemtest/tree_target_test.go:108`, TestTreeTargetRefusesStaleKitBuild | PrintTreeBuildRefusal / S15 | Quote alpha in the exact stderr repair. Keep exit, empty stdout, stale reason, and child-marker absence. |
| `internal/worktree/pool_reclaim_test.go:154`, reclaimApplyHead and TestReclaimCommandPlansOnlyTheProvablyDeadKeys | poolReclaimAction / S33 | Quote worktree, reclaim, --apply, and the fingerprint in the decoded action. Keep aggregate counts and dead-versus-live classification. |
| `internal/worktree/clean_set_wiring_test.go:53`, TestCleanSetApplyTimeStaleWiring | clean refusal / S30 | Update landed and unclaimed selector command literals. Keep stale refusal, race execution, retained refs, and rejected fingerprint. |
| `internal/worktree/clean_landed_hostile_test.go:30`, TestCleanLandedQuotesSpaceAndGlobPaths and TestCleanLandedControlBytePathRetained | landed clean / S31 | Quote every known verb, flag, path, and assignment operand. Keep safe removal, dirty retention, escaped controls, and the single safe exec pointer. |
| `internal/worktree/clean_unclaimed_test.go:302`, TestCleanUnclaimedPlanNamesApplyOnlyWhenARowRemoves and TestCleanUnclaimedStaleClassRefusesTheOldPlan | unclaimed clean and refusal / S32, S30 | Update AXI apply and stale re-plan command expectations. Preserve the raw retained-row remedy expectation unchanged. Keep applicability, class, stale rejection, and ref preservation. |
| `internal/worktree/clean_set_outcomes_test.go:327`, TestCleanSetStaleReplanAction, TestCleanSetUnclaimedStaleReplanAction, and TestCleanSetApplyTimeStaleRefusal | clean refusal / S30 | Update canonical selector spellings and exact rejected-action checks. Keep selector order, digest omission, aggregate refusal, and retained refs. |
| `internal/worktree/clean_set_command_test.go:72`, TestCleanSetDiscardModifiers and TestCleanSetHostileOperand | explicit clean / S29 | Update positive modifier actions and negative replay-command patterns. Keep hostile refusals, no fingerprint, no shell effect, and removal authority. |
| `internal/worktree/clean_set_test.go:28`, TestCleanExplicitSetPlan and TestCleanExplicitSetSelectionFailure | explicit clean / S29 | Update complete-set apply and forbidden replay patterns. Keep membership-bound fingerprints, per-target outcomes, no action on failure, and correct removal. |
| `internal/worktree/clean_landed_test.go:100`, TestCleanLandedPlanAdvertisesApplyAndRemedies and TestCleanLandedPlanApplyCarriesModifiers | landed clean / S31 | Update apply, dirty-path remedy, and modifier literals. Keep plan-only non-mutation, identity, applicability, and modifier order. |

All AXI known arguments in these commands gain single quotes, including dynamic identities and fingerprints.
For example, the apply prefix becomes `bench 'worktree' 'reclaim' '--apply'`.
The fingerprint follows as a quoted known value.
Raw tree-build refusal keeps its fixed tokens and quotes only the label.
Do not change producer selection, authority, or output schemas to accommodate a fixture.

Update every affected expectation within each named file, including shared literal builders and their callers.
The reclaimApplyHead consumers include stale-plan and apply checks beyond the first census hit.
Decode TOON command cells where encoded quoting would make a substring assertion ambiguous.
Preserve each exact command predicate with an independent approved byte expectation.
Recover argv separately through axitest.RecoverHelpCommandArgv.

A negative bare-command search must become a negative search for the new command spelling.
Leaving its old bare pattern would silently weaken the no-replay assertion.
Keep control-byte refusal and shell-sentinel assertions independently of the quote-byte expectations.
An omitted fixture update or returned former bare spelling must turn the affected assertion red.

Read-only exclusion: `internal/systemtest/status_route_converge_test.go:54` reads a different renderer.
Its command expectation receives no rewrite authority from S39.

The broader hidden search also reads tests/canary and .github consumers.
Their command guidance names command tokens, not the changed quote-owner return bytes.
They receive no rewrite authority from this spec.
Any later exact-byte hit requires a fence and coverage amendment before slicing.

## Implementation chunks

Each ticket is one green checkpoint on a retained integration source.
Each named chunk receives independent Standards, Spec, and Coverage review before its dependent successors.
Independent ticket graph review accepted the frozen source recorded in assets/spec-review.md.

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
|---|---|---|---|---|
| S-A / 01-quote-tree-build-arguments.md | Quote tree-build arguments through the existing owner | S01, S02, S03, S04, S15 | owner-tests, worktree-tests, target-tests, system-tests, axi-registry, routing | yes |
| S-B1 / 02-quote-selected-history.md | Quote selected-history command operands | S13 | caller-tests | no |
| S-B2 / 03-quote-citation-commands.md | Quote citation command arguments | S14 | caller-tests, axi-registry, routing | no |
| S-B3 / 04-quote-recovery-operands.md | Quote direct recovery-route operands | S16 | caller-tests, axi-registry, routing | no |
| S-B4 / 05-quote-evidence-commands.md | Quote evidence pagination command arguments | S17 | caller-tests | no |
| S-D1 / 06-share-rebuild-and-fixture-quotes.md | Share rebuild and fixture shell quoting | S18, S37 | rebuild-tests, diff-tests, release-fixtures, axi-registry, routing | no |
| S-E / 07-switch-axi-output-atomically.md | Switch AXI known arguments with their complete output closure | S05, S06, S07, S08, S09, S10, S11, S12, S19, S20, S21, S22, S23, S24, S25, S26, S27, S28, S29, S30, S31, S32, S33, S34, S35, S36, S38, S39 | ordinary-tests, axi-registry, routing, coverage | yes |

Stable chunk mapping: S-A becomes S-A and S-E.
S-B becomes S-B1 through S-B4, with its freshness caller in S-D1.
S-C becomes S-E, the complete global AXI output checkpoint.
S-D becomes S-D1 and the final source census in S-E.

S-A serves the existing quote owner through a real tree-build caller.
Its three exact reader fixtures land with that behavior.
S-E follows every direct migration and switches AXI output with all affected executable-action readers.
The renderer change and its fixture closure must remain one green ticket.
No temporary bare-versus-quoted presentation policy is permitted.

Shared registries order S-A, S-B2, S-B3, S-D1, and S-E.
The approved eight clean and reclaim reader paths belong to S-E.
All eleven reader files remain required by the final S39 reconciliation.

### Completion plan

The authored version 1 plan records future implementation evidence.
It claims no current implementation red, pass, native qualification, or benchmark.
The implementation orchestrator records author sessions in the required version 2 amendment before dispatch.

```bench-completion-plan
{
  "version": 1,
  "chunks": [
    {
      "id": "S-A",
      "tickets": [
        "01-quote-tree-build-arguments.md"
      ],
      "verification": [
        {
          "id": "owner-tests",
          "command": "bench test --package ./internal/sanitize"
        },
        {
          "id": "worktree-tests",
          "command": "bench test --package ./internal/worktree --run 'TreeBuildArgument|TargetRefusal'"
        },
        {
          "id": "target-tests",
          "command": "bench test --package ./internal/treetarget"
        },
        {
          "id": "system-tests",
          "command": "bench test --check system"
        },
        {
          "id": "axi-registry",
          "command": "bench test --check axi-query-registry"
        },
        {
          "id": "routing",
          "command": "bench test --check subcommand-routing"
        },
        {
          "id": "owner-omission-proof",
          "command": "bench test --package ./internal/sanitize",
          "probe": "Return bare alpha, split an embedded quote, or drop the empty owner argument. The exact quoted-byte and NUL-framed recovered-argv owner assertions must fail, then pass after source restoration."
        },
        {
          "id": "owner-preservation",
          "command": "bench test --package ./internal/sanitize",
          "probe": "Alter the existing always-quoted owner result for spaces, Unicode, or embedded quotes. The independent input-to-argv witness must fail, then pass after restoration."
        },
        {
          "id": "caller-bypass-proof",
          "command": "bench test --package ./internal/worktree --run 'TreeBuildArgument|TargetRefusal'",
          "probe": "Bypass sanitize.ShellQuote in PrintTreeBuildRefusal or restore its former bare alpha spelling. The actual tree-build producer fixture and recovered argv assertion must fail, then pass after source restoration."
        },
        {
          "id": "caller-preservation",
          "command": "bench test --package ./internal/worktree --run 'TreeBuildArgument|TargetRefusal'",
          "probe": "Drop the tree-build operand, change its exact quoted spelling, or bypass retained whole-line control escaping. The real producer or refusal assertion must fail, then pass after restoration."
        }
      ]
    },
    {
      "id": "S-B1",
      "tickets": [
        "02-quote-selected-history.md"
      ],
      "verification": [
        {
          "id": "caller-tests",
          "command": "bench test --package ./internal/spec"
        },
        {
          "id": "omission-proof",
          "command": "bench test --package ./internal/spec",
          "probe": "Bypass the owner with former bare-token output. The independent exact fixture must fail even if argv still round-trips."
        },
        {
          "id": "preservation",
          "command": "bench test --package ./internal/spec",
          "probe": "Compare real caller payload, mode, refusal, authority, and output assertions with the frozen baseline. Keep independent expected bytes and demonstrate the named omission red, restoration, and pass."
        }
      ]
    },
    {
      "id": "S-B2",
      "tickets": [
        "03-quote-citation-commands.md"
      ],
      "verification": [
        {
          "id": "caller-tests",
          "command": "bench test --package ./internal/consumers"
        },
        {
          "id": "axi-registry",
          "command": "bench test --check axi-query-registry"
        },
        {
          "id": "routing",
          "command": "bench test --check subcommand-routing"
        },
        {
          "id": "omission-proof",
          "command": "bench test --package ./internal/consumers",
          "probe": "Return the former bare spelling at citation.cmd or omit a supplied flag. The exact fixture or recovered argv must turn red."
        },
        {
          "id": "preservation",
          "command": "bench test --package ./internal/consumers",
          "probe": "Compare real caller payload, mode, refusal, authority, and output assertions with the frozen baseline. Keep independent expected bytes and demonstrate the named omission red, restoration, and pass."
        }
      ]
    },
    {
      "id": "S-B3",
      "tickets": [
        "04-quote-recovery-operands.md"
      ],
      "verification": [
        {
          "id": "caller-tests",
          "command": "bench test --package ./internal/worktree --run 'RecoveryRouteArgument|List|PathIdentifier'"
        },
        {
          "id": "axi-registry",
          "command": "bench test --check axi-query-registry"
        },
        {
          "id": "routing",
          "command": "bench test --check subcommand-routing"
        },
        {
          "id": "omission-proof",
          "command": "bench test --package ./internal/worktree --run 'RecoveryRouteArgument|List|PathIdentifier'",
          "probe": "Bypass the owner only at recoveryRoute.line. The raw command fixture must fail independently of AXI help tests."
        },
        {
          "id": "preservation",
          "command": "bench test --package ./internal/worktree --run 'RecoveryRouteArgument|List|PathIdentifier'",
          "probe": "Compare real caller payload, mode, refusal, authority, and output assertions with the frozen baseline. Keep independent expected bytes and demonstrate the named omission red, restoration, and pass."
        }
      ]
    },
    {
      "id": "S-B4",
      "tickets": [
        "05-quote-evidence-commands.md"
      ],
      "verification": [
        {
          "id": "caller-tests",
          "command": "bench test --package ./internal/preflight/evidencecmd"
        },
        {
          "id": "omission-proof",
          "command": "bench test --package ./internal/preflight/evidencecmd",
          "probe": "Return a bare cursor or bypass the owner at evidenceInvocation. The exact producer fixture and argv oracle must detect each omission."
        },
        {
          "id": "preservation",
          "command": "bench test --package ./internal/preflight/evidencecmd",
          "probe": "Compare real caller payload, mode, refusal, authority, and output assertions with the frozen baseline. Keep independent expected bytes and demonstrate the named omission red, restoration, and pass."
        }
      ]
    },
    {
      "id": "S-D1",
      "tickets": [
        "06-share-rebuild-and-fixture-quotes.md"
      ],
      "verification": [
        {
          "id": "rebuild-tests",
          "command": "bench test --package ./internal/freshness"
        },
        {
          "id": "diff-tests",
          "command": "bench test --package ./internal/diff"
        },
        {
          "id": "release-fixtures",
          "command": "bench test --package ./internal/preprelease"
        },
        {
          "id": "axi-registry",
          "command": "bench test --check axi-query-registry"
        },
        {
          "id": "routing",
          "command": "bench test --check subcommand-routing"
        },
        {
          "id": "preservation",
          "command": "bench test --package ./internal/freshness",
          "probe": "Compare real caller payload, mode, refusal, authority, and output assertions with the frozen baseline. Keep independent expected bytes and demonstrate the named omission red, restoration, and pass."
        },
        {
          "id": "rebuild-argument-proof",
          "command": "bench test --package ./internal/freshness",
          "probe": "Emit a malformed quoted rebuild operand or drop a supplied rebuild operand. Its real producer fixture and independent shell argv assertion must fail, then pass after source restoration."
        },
        {
          "id": "diff-fixture-argument-proof",
          "command": "bench test --package ./internal/diff",
          "probe": "Corrupt quote escaping for a hostile diff fixture operand. Its independent observed argv assertion must fail, then pass after source restoration."
        },
        {
          "id": "release-fixture-argument-proof",
          "command": "bench test --package ./internal/preprelease",
          "probe": "Corrupt quote escaping for a hostile preprelease fixture operand. Its independent observed argv assertion must fail, then pass after source restoration."
        }
      ]
    },
    {
      "id": "S-E",
      "tickets": [
        "07-switch-axi-output-atomically.md"
      ],
      "verification": [
        {
          "id": "ordinary-tests",
          "command": "bench test --package ./..."
        },
        {
          "id": "axi-registry",
          "command": "bench test --check axi-query-registry"
        },
        {
          "id": "routing",
          "command": "bench test --check subcommand-routing"
        },
        {
          "id": "coverage",
          "command": "bench coverage --check shell-argument-ownership"
        },
        {
          "id": "omission-proof",
          "command": "bench test --package ./...",
          "probe": "Return former bare-token spelling at each enumerated family, quote a future slot, compare the rendered first command, or omit a fixture update. Record each independent exact-byte red and restoration separately from shell argv evidence. Keep the accepted empty-helper observation separate from reachable AXI empty refusal. Reconcile all source definitions and fixture helpers without granting status_route_converge rewrite authority."
        },
        {
          "id": "preservation",
          "command": "bench test --package ./...",
          "probe": "Compare real caller payload, mode, refusal, authority, and output assertions with the frozen baseline. Keep independent expected bytes and demonstrate the named omission red, restoration, and pass."
        }
      ]
    }
  ],
  "final_verification": [
    {
      "id": "ordinary-integration",
      "command": "bench test --package ./..."
    },
    {
      "id": "coverage",
      "command": "bench coverage --check shell-argument-ownership"
    },
    {
      "id": "ticket-grammar",
      "command": "bench test --check ticket-grammar"
    },
    {
      "id": "system-integration",
      "command": "bench test --check system"
    }
  ]
}
```

## Testing decisions

Test purpose takes precedence over the number of assertions.
The shared owner tests recover NUL-framed argv through a local POSIX shell.
Fixtures cover empty, safe, Unicode, quote, space, glob, backslash, dollar, semicolon, and substitution-looking arguments.
They detect both argument count and argument bytes.
A sentinel file detects accidental command execution.
NUL is outside the shell argv domain.

Use the existing axitest.RecoverHelpCommandArgv seam for executable help without future-input slots.
Decode TOON before the shell probe so escaped layout bytes reach the command intact.
For multiple actions, decode and inspect each command cell.
For future inputs, pin exact placeholder bytes and probe the known prefix separately.
Keep byte-fixture expectations independent of the quote implementation.
Record the omission or bypass mutation that proves each independent expectation necessary.

Each direct caller and AXI producer has its own composition row.
Apply a compiling bypass mutation at the caller, not only at sanitize.ShellQuote.
The exact producer fixture must turn red when it returns the former bare-token policy.
The shell oracle must turn red when a value loses a quote or an argument.
Existing tests alone do not establish the approved output delta.

Prior art: `internal/axi/axitest/help.go:15` recovers argv through sh.
Prior art: `internal/consumers/citation_test.go:190` crosses the real citation producer.
Prior art: `internal/spec/history_selected_test.go:98` executes the printed history detail command.
Prior art: `internal/freshness/freshness_verify_test.go:246` executes the printed rebuild action.
Read enforcement: `internal/conformance/axi_query_registry_test.go:27` grades the declared query inventory.

The dev gate unit phase executes the affected package tests.
The whole-project gate remains the acceptance oracle.

### Seam diagram

    caller input and sink policy
        |
        v
    exact argument --> sanitize.ShellQuote --> rendered command --> decoded cell --> POSIX-shell argv
                            ^                       ^                                ^
                            |                       |                                |
                    owner quote fixtures     real producer fixtures           byte and count oracle

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| S01 | 1 | The shared owner renders an empty string as one empty POSIX-shell argument. | planned internal/sanitize/shellquote_test.go (TestShellQuoteRoundTrip) | A bare empty string loses its argument. |
| S02 | 2 | An embedded single quote survives a POSIX-shell round trip. | planned internal/sanitize/shellquote_test.go (TestShellQuoteRoundTrip) | The shell oracle detects a split or altered argument. |
| S03 | 3 | Spaces, glob characters, dollar signs, and shell operators round-trip literally. | planned internal/sanitize/shellquote_test.go (TestShellQuoteRoundTrip) | A sentinel file and NUL-framed argv detect expansion or execution. |
| S04 | 4 | A safe nonempty value still renders with surrounding single quotes. | planned internal/sanitize/shellquote_test.go (TestShellQuoteAlwaysQuoted) | A safe-token shortcut violates the exact output predicate. |
| S05 | 5 | AXI executable actions quote each accepted known argument through sanitize.ShellQuote. | planned internal/axi/action_test.go (TestRenderHelpAlwaysQuoted) | A round trip alone misses the approved quote-policy delta. |
| S06 | 6 | Future-input slots retain their exact placeholder spelling. | planned internal/axi/action_test.go (TestRenderHelpFutureInputs) | Quoting a future slot changes the action grammar. |
| S07 | 7 | A known empty AXI argument still produces the honest empty help envelope. | `internal/axi/action_test.go` (`TestRenderHelpUsesHonestEmptyForUnsupportedDisclosureValue`) | Accepting empty at this caller changes a pre-existing assertion. |
| S08 | 8 | Unsupported AXI control runes still produce the honest empty help envelope. | `internal/axi/action_test.go` (`TestKnownArgumentRefusesUnsupportedControls`) | The accepted owner cannot replace the caller refusal policy. |
| S09 | 9 | AXI tab, newline, and carriage-return arguments retain their exact shell values. | `internal/axi/action_test.go` (`TestKnownArgumentControlValuesRoundTripThroughTOONAndPOSIXShell`) | A line-safe filter would incorrectly remove accepted bytes. |
| S10 | 10 | An executable action whose first known value is bench remains refused. | planned internal/axi/action_test.go (TestRenderHelpRejectsQuotedCommandSpoof) | Comparing rendered bytes misses the new quoted spelling. |
| S11 | 11 | Identical actions still collapse without source-order changes. | `internal/axi/action_test.go` (`TestRenderHelpDeduplicatesExactTemplatesWithoutReordering`) | A formatting migration cannot change action cardinality. |
| S12 | 12 | Harness-phase actions retain their current rendered argument bytes. | `internal/axi/action_test.go` (`TestRenderHelpRendersReusableInvocationAndHarnessPhase`) | Applying shell quotes to a harness operand changes a distinct interface. |
| S13 | 13 | Selected-history detail operands use the surviving quote owner. | `internal/spec/history_selected_test.go` (`TestSelectedSpecDetailRoute`) | The command executes through a shell and reaches the selected history. |
| S14 | 14 | Citation commands preserve the supplied argv through the surviving owner. | `internal/consumers/citation_test.go` (`TestCitationCmdRoundTripsThroughAPOSIXShellSplit`) | An omitted flag or bare metacharacter changes recovered argv. |
| S15 | 15 | Tree-build refusal commands preserve the exact assignment label. | planned internal/worktree/shell_arguments_test.go (TestTreeBuildArgument) | A hostile label and shell argv oracle detect a bypass. |
| S16 | 16 | RecoveryRoute.line quotes its path through the surviving owner. | planned internal/worktree/shell_arguments_test.go (TestRecoveryRouteArgument) | The raw refusal line must preserve argv independently of help output. |
| S17 | 17 | Evidence commands preserve identity, source, and cursor arguments through the surviving owner. | planned internal/preflight/evidencecmd/shell_arguments_test.go (TestEvidenceArgumentCommands) | A producer-level test detects a direct helper bypass. |
| S18 | 18 | Freshness rebuild commands use sanitize.ShellQuote for each path. | `internal/freshness/freshness_verify_test.go` (`TestRefusalRebuildActionIsCopyPasteSafeForHostilePaths`) | The existing real-script fixture observes cwd and both arguments. |
| S19 | 19 | The anchors producer emits the approved quoted command shape. | planned cmd/bench/anchor_help_test.go (TestAnchorHelpArguments) | A reviewed exact fixture and recovered argv detect wrong spelling or altered arguments. |
| S20 | 20 | The coverage producer emits the approved quoted command shape. | planned internal/coverage/shell_arguments_test.go (TestCoverageHelpArguments) | A reviewed exact fixture and recovered argv detect wrong spelling or altered arguments. |
| S21 | 21 | The consumers query producer emits the approved quoted command shape. | planned internal/consumers/shell_arguments_test.go (TestConsumersQueryHelpArguments) | A reviewed exact fixture and recovered argv detect wrong spelling or altered arguments. |
| S22 | 22 | The consumers blast producer emits the approved quoted command shape. | planned internal/consumers/shell_arguments_test.go (TestConsumersBlastHelpArguments) | A reviewed exact fixture and recovered argv detect wrong spelling or altered arguments. |
| S23 | 23 | The guards producer emits the approved quoted command shape. | planned internal/guards/guards_test.go (TestGuardsHelpArguments) | A reviewed exact fixture and recovered argv detect wrong spelling or altered arguments. |
| S24 | 24 | The maps producer emits the approved quoted command shape. | planned internal/maps/maps_command_test.go (TestMapsHelpArguments) | A reviewed exact fixture and recovered argv detect wrong spelling or altered arguments. |
| S25 | 25 | The gate retry producer emits the approved quoted command shape. | planned internal/gate/run_outcomes_test.go (TestGateHelpArguments) | A reviewed exact fixture and recovered argv detect wrong spelling or altered arguments. |
| S26 | 26 | The preflight retry producer emits the approved quoted command shape. | planned internal/preflight/shell_arguments_test.go (TestPreflightHelpArguments) | A reviewed exact fixture and recovered argv detect wrong spelling or altered arguments. |
| S27 | 27 | The worktree list producer emits the approved quoted command shape. | planned internal/worktree/shell_arguments_test.go (TestListHelpArguments) | A reviewed exact fixture and recovered argv detect wrong spelling or altered arguments. |
| S28 | 28 | The worktree selected list producer emits the approved quoted command shape. | planned internal/worktree/shell_arguments_test.go (TestSelectedListHelpArguments) | A reviewed exact fixture and recovered argv detect wrong spelling or altered arguments. |
| S29 | 29 | The worktree explicit clean producer emits the approved quoted command shape. | planned internal/worktree/shell_arguments_test.go (TestExplicitCleanHelpArguments) | A reviewed exact fixture and recovered argv detect wrong spelling or altered arguments. |
| S30 | 30 | The worktree clean refusal producer emits the approved quoted command shape. | planned internal/worktree/shell_arguments_test.go (TestCleanRefusalHelpArguments) | A reviewed exact fixture and recovered argv detect wrong spelling or altered arguments. |
| S31 | 31 | The worktree landed clean producer emits the approved quoted command shape. | planned internal/worktree/shell_arguments_test.go (TestLandedCleanHelpArguments) | A reviewed exact fixture and recovered argv detect wrong spelling or altered arguments. |
| S32 | 32 | The worktree unclaimed clean producer emits the approved quoted command shape. | planned internal/worktree/shell_arguments_test.go (TestUnclaimedCleanHelpArguments) | A reviewed exact fixture and recovered argv detect wrong spelling or altered arguments. |
| S33 | 33 | The worktree pool reclaim producer emits the approved quoted command shape. | planned internal/worktree/shell_arguments_test.go (TestReclaimHelpArguments) | A reviewed exact fixture and recovered argv detect wrong spelling or altered arguments. |
| S34 | 34 | AXI and freshness no longer define independent shell quote algorithms. | review-owned definition and consumer census | A surviving algorithm violates the selected owner decision. |
| S35 | 35 | Direct callers retain their current refusal or escape policy for control-bearing operands. | planned direct caller tests from the inventory | A quoted newline must not bypass a line-structured refusal. |
| S36 | 36 | Differential fixtures change only the approved argument quote spellings. | planned domain fixtures from the output-reader inventory | The comparison detects changed schema, counts, reason text, exit codes, or command order. |
| S37 | 37 | Fixture shell arguments use sanitize.ShellQuote instead of independent helpers. | planned migrated fixture tests from the ownership fences | A source inspection plus fixture execution detects surviving helpers or altered fixture argv. |
| S38 | 38 | Each migrated caller family fails its declared quote-bypass mutation. | planned production caller tests from the inventory | Returning the former bare-token output turns the exact producer fixture red. |
| S39 | 39 | All eleven affected output-reader files preserve their existing behavioral assertions with only approved quote-byte updates. | planned existing fixtures named in the changed producer-to-assertion closure | Exact commands and negative replay checks detect omissions without weakening control or authority assertions. |

### Edge inventory

The audience includes this kit and linked repositories that receive its executable.
Empty helper input belongs to S01.
Empty AXI KnownArgument input retains refusal through S07.
These are different interfaces with different contracts.
The source does not establish a reachable empty-argument defect.

Accepted AXI layout bytes belong to S09.
Unsupported control bytes belong to S08.
Line-structured direct callers retain their policy through S35.
Spaces, Unicode, trailing spaces, glob characters, and embedded quotes belong to S02-S03 and each producer fixture.
Flags, end-of-options markers, and future slots retain their positions through S06 and S36.
Exact duplicate actions and source order belong to S11.

**Won't handle:** NUL in shell argv — the surviving sanitize owner serves shell-representable values.
**Won't handle:** a new control-byte policy — the surviving history caller keeps LineSafe before it renders the operand.
**Won't handle:** harness-phase shell syntax — the surviving AXI HarnessPhaseOn renderer retains its separate contract.
**Won't handle:** command authorization — the surviving preflight caller keeps validation before the quote owner.
**Won't handle:** a new empty AXI action feature — the surviving renderKnownArgument caller keeps its explicit refusal.
**Won't handle:** literal fixed command grammar — the surviving InspectFull and RetryDiff renderers keep their current validated tokens.

## Ownership fences

The union below equals the implementation ticket Writes union.
Each ticket also owns its exact required registry and fixture closure.
The review pickup remains phase-owned.
Registry co-ownership preserves command membership, help grammar, and routing.
The required closure grants no additional product behavior.
No implementation ticket may sweep unrelated specs or future ticket folders.

- `cmd/bench/anchor_help_test.go`
- `cmd/bench/anchors_dir_test.go`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `cmd/bench/main_test.go`
- `internal/axi/action.go`
- `internal/axi/action_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `internal/consumers/blast_edges_test.go`
- `internal/consumers/blast_test.go`
- `internal/consumers/candidates_test.go`
- `internal/consumers/citation.go`
- `internal/consumers/citation_test.go`
- `internal/consumers/command_test.go`
- `internal/consumers/loader_test.go`
- `internal/consumers/query_test.go`
- `internal/consumers/refuse_test.go`
- `internal/consumers/shell_arguments_test.go`
- `internal/coverage/coverage_command_test.go`
- `internal/coverage/coverage_schema_test.go`
- `internal/coverage/shell_arguments_test.go`
- `internal/diff/identity_test.go`
- `internal/freshness/freshness_verify.go`
- `internal/freshness/freshness_verify_test.go`
- `internal/gate/run_outcomes_test.go`
- `internal/guards/guards_test.go`
- `internal/maps/freshness_test.go`
- `internal/maps/maps_command_test.go`
- `internal/preflight/evidencecmd/evidence.go`
- `internal/preflight/evidencecmd/evidence_cleanup_test.go`
- `internal/preflight/evidencecmd/evidence_command_test.go`
- `internal/preflight/evidencecmd/evidence_modes_test.go`
- `internal/preflight/evidencecmd/evidence_review_test.go`
- `internal/preflight/evidencecmd/evidence_summary_test.go`
- `internal/preflight/evidencecmd/shell_arguments_test.go`
- `internal/preflight/explicit_base_test.go`
- `internal/preflight/shell_arguments_test.go`
- `internal/preprelease/preprelease_test.go`
- `internal/sanitize/shellquote_test.go`
- `internal/spec/history_selected.go`
- `internal/spec/history_selected_test.go`
- `internal/systemtest/tree_target_test.go`
- `internal/treetarget/build_test.go`
- `internal/worktree/clean_landed_hostile_test.go`
- `internal/worktree/clean_landed_test.go`
- `internal/worktree/clean_set_command_test.go`
- `internal/worktree/clean_set_outcomes_test.go`
- `internal/worktree/clean_set_test.go`
- `internal/worktree/clean_set_wiring_test.go`
- `internal/worktree/clean_unclaimed_test.go`
- `internal/worktree/identifier_operand_test.go`
- `internal/worktree/landed_test.go`
- `internal/worktree/list.go`
- `internal/worktree/list_actions_test.go`
- `internal/worktree/path_identifier_test.go`
- `internal/worktree/pool_reclaim_test.go`
- `internal/worktree/shell_arguments_test.go`
- `internal/worktree/target_refusal_test.go`
- `internal/worktree/tree_target.go`
- `internal/worktree/unlanded_route_test.go`
- `reviews/shell-argument-ownership.md`

## Out of scope

No implementation, runtime reproduction, benchmark, ticket, or roadmap edit belongs to this author phase.
FT354 retains the strict-JSON and durable-file outcomes.
This spec has no Roadmap retirement field.

A new empty AXI action contract is a separate capability: 2 edits, 1 gate run.
A new line-output control policy is a separate capability: 4 edits, 2 gate runs.
A complete command-renderer redesign is a separate capability: 8 edits, 2 gate runs.
These estimates describe future planning units, not implementation promises.

## Further notes

### Source-sentence coverage

| reviewed source clause | coverage |
|---|---|
| Use sanitize.ShellQuote and the always-quoted contract. | S01-S05, S13-S34 |
| Remove independent algorithms without a presentation option. | S04, S34, S37 |
| Preserve empty and embedded-quote arguments. | S01-S03 |
| Establish reachability before calling empty input a reproduced defect. | Current-code distinction and S07 |
| Retain caller control-byte refusal for line output. | S08-S09, S35 |
| Quoting grants no authority. | S10 and the explicit scope cuts |
| Enumerate every changed printed command. | Direct and AXI producer inventories, S13-S33, S39 |
| Test POSIX-shell round trips and exact contractual output. | S01-S06, S13-S33, S36 |
| Preserve error categories and existing refusals. | S07-S12, S35-S36, S39 |
| Attach a bypass mutation to each migrated family. | S38 and the per-producer rows |

### Pre-review proof checklist

- Cited symbols: definition reads and bench consumers resolve the direct quote algorithms and their callers.
- Import edges: go list resolves sanitize, axi, freshness, and existing producer packages. The new quote imports remain cycle checks before implementation.
- Source-row clauses and occurrences: the source table consumes tickets 3 and 4 without moving the shared parent map.
- Promised field labels: no field changes. Existing cmd, why, detail, next, and citation fields retain their schemas.
- Changed-function callers: the direct inventory includes renderKnownArgument and its executable-action producer closure.
- Copy survival: S34 and S37 grade remaining algorithms. S38 requires producer bypass reds.
- Rendered-shape readers: the output tables include the eleven reviewed tree-build, clean, and reclaim readers. S39 owns every approved fixture update.
- Pin operators: exact approved command bytes and exact recovered argv. Unchanged output assertions retain their current operators.
- Entry reads: none introduced.
- Derived expectations: fixture argv comes from supplied inputs, not from the quote implementation.
- Consolidated rules: the direct inventory gives the old owner and approved new spelling per site.
- Quantified obligations: each direct caller and each AXI executable producer receives a composition row.
- Workflow-step writes: none introduced.

### Fixture algorithm dispositions

The diff identity and preprelease fixtures define independent shell-token quote helpers.
The freshness test quoteForShell helper reproduces it too.
Migrate those fixture inputs to sanitize.ShellQuote.
Keep their observed argv expectations independent.
The systemtest shellQuoteJSON helper uses JSON string encoding for a tool envelope.
It remains outside the shell-token algorithm inventory.
Sources: `internal/diff/identity_test.go:272`, `internal/preprelease/preprelease_test.go:179`, and `internal/systemtest/bench_follow_on_test.go:340`.

### Evidence status and validation plan

Facts come from the pinned tree and the named reviewed artifacts.
The static consumer census includes production and test callers.
A hidden tree search includes the canary fixtures and workflow files.
No runtime reproduction or byte-compatibility probe ran during authoring.
The current empty helper observation remains a source-level fact.

Implementation must produce the separate shell and output-byte evidence before acceptance.
The spec-only checks are prose mechanics and acceptance coverage.
Independent review must approve the provisional fences and chunks before ticket creation.

### Review repair disposition

| finding | disposition | coverage and fence evidence |
|---|---|---|
| SAO-R1-OUTPUT-CLOSURE | Repaired for confirming review. | S39 names eleven exact reader paths, existing seams, producer rows, and quote-only fixture updates. All eleven paths join the Writes union. |

This is review repair round one.
No implementation or runtime reproduction supports this planning change.
The prior argv, control-byte, authority, and refusal assertions remain required.
Independent confirming review precedes any ticket slicing.

## Ticket approval

Slicing line: gpt-6.1-sol / high / one pass plus at most two bounded repairs.
Independent ticket review accepted the frozen graph before this spec-stage close.
Accepted spec hashes and closed findings are recorded in assets/spec-review.md.
The accepted product rows, exclusions, API, and failure guarantees remain unchanged.

| numbered ticket | Blocked by | delivered outcome |
|---|---|---|
| 1. 01-quote-tree-build-arguments.md | none | Quote tree-build arguments through the existing owner |
| 2. 02-quote-selected-history.md | 01-quote-tree-build-arguments.md | Quote selected-history command operands |
| 3. 03-quote-citation-commands.md | 01-quote-tree-build-arguments.md | Quote citation command arguments |
| 4. 04-quote-recovery-operands.md | 01-quote-tree-build-arguments.md, 03-quote-citation-commands.md | Quote direct recovery-route operands |
| 5. 05-quote-evidence-commands.md | 01-quote-tree-build-arguments.md | Quote evidence pagination command arguments |
| 6. 06-share-rebuild-and-fixture-quotes.md | 01-quote-tree-build-arguments.md, 04-quote-recovery-operands.md | Share rebuild and fixture shell quoting |
| 7. 07-switch-axi-output-atomically.md | 01-quote-tree-build-arguments.md, 02-quote-selected-history.md, 03-quote-citation-commands.md, 04-quote-recovery-operands.md, 05-quote-evidence-commands.md, 06-share-rebuild-and-fixture-quotes.md | Switch AXI known arguments with their complete output closure |

### Source headroom

S-E updates existing reader literals without growing over-budget fixture files.
Place its new worktree list-help tests in the co-owned shell_arguments_test.go file.
Keep the existing fixture helpers and their assertion purpose.
A required extraction must use an exact co-owned destination before the affected ticket starts.
No ticket may raise structure budgets or add an acceptance grant.

### Verification ownership

S-A separates owner evidence from real tree-build caller evidence.
S-D1 uses separate freshness, diff, and preprelease commands for their argument witnesses.
S34 and S37 retain independent source-census review of surviving quote algorithms.
Ticket grammar cannot replace that source inspection or a caller bypass probe.

### Accepted ticket-review source

Accepted graph commit: 3e7193118ec710e9e1e25c3cbe0255d92ab7b02a
Accepted spec SHA256 before closure metadata: 5a218cec5b4529b8563d7f09c61f5067cd1390c21dc340d149f2831e2160c87e
Reviewer: /root/cleanup_ticket_review
Reviewer line: gpt-6.1-sol / high
Judgment: Standards 0 blockers, Spec 0 blockers, Coverage 0 blockers
Review-axis confidence: 9 each, static planning only

The version 1 plan remains an authored implementation plan, not completed implementation evidence.
No runtime test, native qualification, manual gate, or benchmark supports this stage.
Ticket bytes, acceptance rows, ownership fences, and implementation decisions remain at the accepted source.
SAO-T1-RAW-REMEDY closed on independent confirmation with finding-closure confidence 10.
