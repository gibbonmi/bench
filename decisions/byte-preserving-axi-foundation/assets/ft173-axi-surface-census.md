# FT173 current AXI surface census

Evidence asset for `decisions/byte-preserving-axi-foundation/ft173-axi-contract.md` ticket #1. Inspected
2026-08-09 at `974020e4af8de5ed75098c4c5934a8907952bb2b`. This is a
source census only: no command behavior was changed and no gate or test command was
run.

## Contract boundary

The published [AXI contract](https://axi.md/) currently names ten principles:
token-efficient output, minimal default schemas, content truncation, pre-computed
aggregates, definitive empty states, structured errors and exit codes, ambient
context, content first, contextual disclosure, and a consistent way to get help.

Bench has three different statements of that boundary today:

- `craft-cli` documents only principles 1 through 7 and declares AXI for the
  query surfaces rather than the whole binary
  (`.agents/skills/bench-craft-cli/SKILL.md:7-16,18-67`).
- The project profile names seven query commands as the conforming surface:
  `anchors`, `learnings`, `maps`, `guards`, `diff`, `coverage`, and `worktree
  list` (`projects/benchkit.md:40-52`).
- The FT173 map now requires all ten principles for every `bench spec build`
  operation while leaving other operational families outside that automatic
  widening (`decisions/byte-preserving-axi-foundation/ft173-axi-contract.md:124-148`).

The production registry contains 48 root names, but a command definition carries
only `Name`, `Attachment`, and `Run`; it cannot declare AXI scope, output schema,
empty-state behavior, no-argument behavior, or help/disclosure metadata
(`cmd/bench/main.go:66-118`; `cmd/bench/command_registry.go:21-27`). AXI is
therefore a call-site convention, not a registry-derived contract. The current
tree has 43 non-test `toon.Table`/`TableTyped` call sites across 21 files and 23
non-test `usage.Parse` calls across 20 files. Neither set is the declared seven-
command boundary, and neither proves all ten principles.

## Ten-principle census

| # | Principle | Current documented owner | Current implementation and independent derivations | Uncovered surface | Compatibility or gate constraint |
|---:|---|---|---|---|---|
| 1 | Token-efficient output | `craft-cli` #1; project-profile AXI query seam | `internal/toon.table` is the shared flat-table TOON boundary over `toon-go`; `Table` and `TableTyped` provide string and typed rows (`internal/toon/toon.go:22-80`). `outputCommand` prints a returned string to stdout but does not validate its shape (`cmd/bench/main.go:120-125`). `internal/worktree/refresh` independently hand-renders a fallback `worktree_refresh[1]` block (`internal/worktree/refresh/refresh.go:59-63`). `bench diff --full` deliberately appends a raw patch outside TOON (`internal/diff/diff.go:120-129,253-260`). | Registry entries do not say which commands must emit TOON. AXI-shaped output has spread beyond the seven documented query commands into `test`, `outline`, `roadmap --context`, `spec history`, spec-build lifecycle, release, worktree mutation receipts, and `shift_result`, while direct/system handlers may retain plain output. | `TestTableCellEscaping`, `TestTableUnrepresentableCellErrors`, `TestRepresentableMatchesEncoder`, `TestTable`, and `TestTableTyped` pin the shared encoder boundary (`internal/toon/toon_test.go`). The dev gate includes `go test -count=1 ./...` (`internal/gate/phases.go:83-89`). No check derives the AXI surface from the command registry. |
| 2 | Minimal default schemas | `craft-cli` #2 | Schemas are literal slices at each of the 43 table call sites. The declared query defaults are 3 fields for `anchors`, 2 for `learnings`, 5 for `maps`, 7 plus a 5-field scan row for `guards`, 2 for `diff`, 3 for `coverage`, and 8 for `worktree list` (`cmd/bench/main.go:162`; `internal/learnings/learnings.go:216,234`; `internal/maps/maps.go:158`; `internal/guards/guards.go:283-290`; `internal/diff/diff.go:239`; `internal/coverage/coverage.go:434-440`; `internal/worktree/list.go:13-56`). Spec-build default status has 4 fields; `--full` assignments have 9 (`cmd/bench/specbuild.go:128-179`). | There is no schema owner in the registry and no implemented `--fields` flag anywhere under non-test `cmd/` or `internal/`; the only occurrence is the recommendation in `craft-cli`. Several default tables exceed the published 3-4-field list target, and no conformance check distinguishes justified detail schemas from wide defaults. | Exact field names are pinned piecemeal, for example `TestSpecBuildStatusRendersDefinitiveEmptyProjection`, `TestCommandBoundsRowsAndFullRetainsMetadata`, and the TOON tests. There is no generic schema-width or `--fields` assertion. Consolidation must preserve the existing field order and bytes until ticket #8 chooses a byte-changing boundary. |
| 3 | Content truncation | `craft-cli` #3 | Four independent policies remain. (1) `sanitize.Preview`: first 120 code points plus original byte count, 11 non-test calls across 6 files (`internal/sanitize/sanitize.go:41-56`). (2) `roadmap.limited`: 4096 UTF-8-safe bytes plus `*_bytes`/`truncated`, 9 calls in `internal/roadmap` (`internal/roadmap/context_types.go:9,117-127`; `internal/roadmap/context_parse.go`). (3) `worktree.inventoryIgnored`: 20 displayed paths by default, an entry-limit view under `--full`, plus count/bytes/shown/truncated (`internal/worktree/subshell.go:389-439`). (4) `outline`: `bounds.OutlineRowLimit` rows by default, all rows under `--full`, plus total/emitted/omitted metadata (`internal/outline/outline.go:245-264`). `bench test` already composes `sanitize.Preview` for diagnostics and `sanitize.Controls` under `--full` (`internal/testreport/testreport.go:337-342`). | No common truncation result carries preview, total size, truncation state, and escape hatch across surfaces. `bench diff` omits the default log/body rather than emitting a sized preview, and spec-build `buildError` uses unbounded `sanitize.Controls`, not a bounded detail helper (`internal/diff/diff.go:120-129`; `cmd/bench/specbuild.go:124-126`). | `TestPreviewBoundariesAndControls`, `TestContextBodyLimitBoundaries`, `TestIgnoredInventoryEntryAndByteBoundaries`, and `TestCommandBoundsRowsAndFullRetainsMetadata` pin the four policies independently. Routing them through one helper is byte-preserving only if those exact caps, counting units, suffixes, and `--full` behavior remain unchanged. |
| 4 | Pre-computed aggregates | `craft-cli` #4 | `internal/toon.table` derives only the emitted row count in every header (`internal/toon/toon.go:31-53`). Total-versus-emitted counts and derived status remain command-local: `outline_meta`; `guard_scan`; roadmap body bytes, occurrence counts, source totals, and drain counts; worktree ignored count/bytes/shown/truncated; spec-build abandonment/reclamation counts and `Status.Next`; release `next_action`; and status/dashboard signal summaries (`internal/outline/outline.go:245-264`; `internal/guards/guards.go:283-290`; `internal/roadmap/context_render.go:52-116`; `internal/worktree/classifier.go:252-259`; `cmd/bench/specbuild.go:182-238`; `internal/publication/command.go:320-335`; `internal/status/status.go`). | There is no shared aggregate result type or renderer beyond the table row count. Commands independently decide names, counting units, incomplete/unknown posture, and whether an aggregate describes the page or the whole source. | `TestTable`, `TestCommandAlwaysEmitsCompleteGuardScanMetadata`, `TestScanTimeoutPreservesPartialRowsAndHonestCounts`, `TestCommandBoundsRowsAndFullRetainsMetadata`, roadmap context tests, and spec-build renderer tests pin individual derivations. Any shared mechanic must consume owner-derived facts rather than recalculate them. |
| 5 | Definitive empty states | `craft-cli` #5; project-profile AXI seam | `internal/toon.table` renders `name[0]{fields}:` and is the shared empty table owner (`internal/toon/toon.go:31-35,56-62`). Query commands reach it independently: absent learnings, no active maps, a clean diff, no worktrees, no outline symbols, and no spec history. Human renderers separately own prose empties: status clean, roadmap missing/recommended-sequence, and dashboard section empties. Spec-build represents no run as one `spec_build` row with `state=empty`, and `--full` adds zero-row assignments/review tables (`internal/specbuild/state.go:560-566`; `cmd/bench/specbuild.go:148-179`). | The registry cannot require an empty result, and there is no all-command empty-state matrix. A zero table, a one-row `state=empty`, and prose such as `No gate cache` are unrelated contracts rather than variants of one declared empty-state type. | `TestTable`, `TestSpecBuildStatusRendersDefinitiveEmptyProjection`, `TestStatusHasDefinitiveEmptyAndActiveProjections`, `TestRenderClean`, and `TestRenderEmptyStates` pin representative variants. The gate has no ten-principle empty-state check. |
| 6 | Structured errors and honest exit codes | `craft-cli` #6; project-profile hybrid stdout/stderr rule | `toon.Errorf`, `RenderError`, `Usage`, `MissingArg`, `NotInRepo`, and `RecordError` own common line shapes (`internal/toon/toon.go:110-150`). `usage.Parse` owns flat-command help, unknown/repeated flags, missing values, arity, `--`, and exit 0/2 (`internal/usage/parse.go:47-149`). Query functions return `(stdout, code)` through `outputCommand`. Spec-build composes `ParseBuild`, `buildError`, and the lifecycle service, so usage is stdout/2 and service refusal is sanitized stdout/1 (`internal/spec/build.go:62-89`; `cmd/bench/specbuild.go:19-32,118-133`). Lifecycle request journals separately own retry/replay idempotency. | Direct/system/nested commands still choose their own sinks; the compiled dispatcher writes missing/unknown root subcommands to stderr/2 (`cmd/bench/command_registry.go:37-60`). Operational idempotency and no-op meaning are state-machine facts, not represented by a shared result. Spec-build errors share a generic retry hint even when the service knows a more exact command. | `internal/usage/parse_test.go`, `TestErrorfUsage`, `TestSpecBuildErrorCannotSplitStructuredOutput`, and `internal/spec/build_test.go` pin the common pieces. The `subcommand-routing` dev check is intended to require `usage.Parse`, but its scanner still looks for the former `commands` map and `run` switch (`internal/conformance/subcommand_routing_test.go:15-31,50-108,153-215`) while production now uses `commandRegistry` and `Command.Run`; source inspection therefore shows that this check no longer derives the live dispatcher. |
| 7 | Ambient context | `craft-cli` #7; `.bench/BENCH.md` and project-profile status seam | `.bench/hooks/session-start.sh` resolves and prints the executable path, then runs `session-inspect`; `sessioninspect` serially runs resume, `status`, and `guards --brief` under one deadline (`.bench/hooks/session-start.sh:14-30`; `internal/sessioninspect/sessioninspect.go:20-36,48-86`). Harness wiring and the informational guard posture are independently declared in adapter configuration and conformance registries. | The local guidance omits the published contract's explicit-setup-first and generated on-demand-skill clauses. No skill is generated from the same home-view facts. `bench status` has no spec-build import, so an active spec build is not a first-class ambient projection. | `TestCommandInstallsTenSecondDeadline`, `TestInspectDeadlineWarnsAndReturnsZero`, `checkCodexHooks`, `checkClaudeHookWiring`, `TestClaudeHookWiringBites`, and the package-core session-start manifest check protect the current hook. They prove wiring and fail-open posture, not that the dashboard covers every AXI family. |
| 8 | Content first | Published AXI only; absent from `craft-cli` | Individual query commands with optional selectors generally render live data on empty argv (`learnings`, `maps`, `guards`, `diff`, `status`, `roadmap`, `models`, `outline`). Required-selector leaves such as `anchors` and `coverage` return usage. The public wrapper defaults root empty argv to its long help block (`bin/bench.sh:284-333`), and `bench spec build` with no operation returns a missing-argument usage line (`internal/spec/build.go:62-67`). | There is no registry-owned home view, executable/description metadata, or rule distinguishing a family home from a mutating leaf. Spec build has no live family home; operation-less invocation exits 2. The session-start executable line partially supplies the published bin-path fact, but not through no-argument command output. | `TestParseBuildNoArgDiagnosticListsExactlyTheGrammarOperations` and `TestParseBuildNoArgDiagnosticIsOrderStable` pin the current spec-build usage bytes. Changing the family root to a live home is necessarily byte- and exit-changing. The old roadmap constraint that content-first is query-scoped prevents interpreting this principle as permission to run a mutation without its required intent (`ROADMAP.md:776-780`). |
| 9 | Contextual disclosure | Published AXI only; absent from `craft-cli` | There are zero production `help[` emitters under `cmd/` and `internal/`. Existing substitutes are independently shaped: status rows have an `action`, spec-build status has `next`, publication has `next_action`, structured errors have free-text hints, and root help advertises examples. Spec-build's `Status.Next` is derived once in `record.status`, but mixes executable Bench commands with orchestration prose such as `release assignment`, `delegate assignment`, and gate diagnosis labels (`internal/specbuild/state.go:74-114`). | No shared `help[]` renderer, command template, placeholder rule, or carry-forward owner exists. Spec-build success does not expose help, and refusals do not consistently carry the known slug, run, candidate, assignment, and exact retry command required by map ticket #4. `Status.Next` cannot be copied blindly into `help[]` because several values are not invokable commands. | There is no contextual-disclosure gate assertion. Adding `help[]` changes pinned stdout by definition; `ROADMAP.md:735-746` grants that byte-changing exception only to the contextual-disclosure slice. Existing status/spec-build/release tests pin the ad hoc fields that a migration must deliberately replace or compose. |
| 10 | Consistent way to get help | Published AXI only; absent from `craft-cli` | `usage.Parse` recognizes bare `help`, `--help`, and `-h` and returns each grammar's declared help (`internal/usage/parse.go:47-65,86-103`). Spec-build derives all nine operation grammars and help lines from `buildOperationOrder`/`buildOperations` (`internal/spec/build.go:16-59`). The wrapper separately owns root help and gate help; nested dispatchers separately own worktree, spec, release, adoption, and other leaf help. | Command registry entries carry no grammar or help text, so root help, `.bench/BENCH.md` inventory, `commands --brief`, `subcommandRouting`, and leaf grammars are independent command advertisements. Help spellings are inconsistent: `specArg` and `worktree list` accept `-h/--help` but not bare `help`, while the shared parser accepts all three (`internal/spec/spec.go:307-329`; `internal/worktree/list.go:15-22`). No exhaustive help test walks the production registry. | `TestParseHelpSpellings`, `TestParseBareHelpOnlyWhenSole`, `TestParseBuildExposesEveryGrammarOperation`, and `TestParseBuildNoArgDiagnosticListsExactlyTheGrammarOperations` pin the shared/parser-local behavior. The stale `subcommand-routing` scanner currently cannot close the registry-wide coverage gap. |

## Independent owner ledger

The reusable-fact problem is broader than the ten rows make easy to scan:

| Fact | Current owners | Census result |
|---|---|---|
| Root command names | `commandRegistry`; wrapper dispatch/help; `.bench/BENCH.md` inventory; `subcommandRouting`; `commands --brief` | Five advertisements/derivations. Only `commandRegistry` executes the Go command, and none of the other four is generated from it. |
| Flat TOON bytes | `internal/toon`; one manual worktree-refresh fallback; raw/plain exceptions | One strong shared encoder with one production re-derivation and multiple explicitly non-TOON contracts. |
| Argument grammar/help | `usage.Parse` plus nested manual parsers and wrapper cases | Shared for 23 calls, not registry-complete. Spec build has a good family-local grammar table. |
| Truncation | `sanitize.Preview`; `roadmap.limited`; `worktree.inventoryIgnored`; `outline` row bound | Four independent policies with different caps/counting units and four separate test owners. |
| Aggregates | TOON header count plus command-local meta/status/count renderers | Only emitted row count is shared. Total, page, incomplete, and next-action meaning are local. |
| Empty state | TOON zero table; command-local prose; spec-build one-row empty state | No declared cross-command type or registry assertion. |
| Next action | status `action`; spec-build `Status.Next`; publication `next_action`; error hints; wrapper help | No shared executable-template type. Spec-build's one state derivation includes non-command prose. |

## Spec-build migration baseline

The spec-build family is already ahead of most operational commands on principles
1, 5, 6, and 10:

- one ordered operation table derives the nine public operation grammars (the eight
  lifecycle operations plus maintainer-only `reclaim`);
- success and detail projections use TOON;
- no-run status is explicit;
- usage/refusal exits are 2/1 on stdout; and
- `record.status` derives lifecycle state and the current next step once.

It is not yet full AXI. It has no family home or ambient projection, no `help[]`, no
bounded error/detail policy, no registry-declared schemas, and no exact executable
next-action type. The current `next` cell is evidence to reuse, not proof of
contextual disclosure. Promotion's gate-result payload also remains outside this
census under FT185 ownership.

## Gate and compatibility finding

The current oracle includes the package tests because the dev test phase runs
`go test -count=1 ./...`, but there is no standalone ten-principle AXI suite in
the tree. Enforcement is fragmented among:

- pure encoder/parser tests in `internal/toon` and `internal/usage`;
- command-local exact or substring output tests;
- ambient hook wiring checks;
- `docs-currency-workflow`'s two literal profile anchors for `bench diff` and
  `bench coverage` (`internal/conformance/docs_workflow_checks_test.go:312-324`);
  and
- the `subcommand-routing` registry/check, whose dispatcher scanner is stale as
  described above.

Therefore the existing gate can protect a byte-preserving helper consolidation only
where a named local test already pins the behavior. It cannot currently certify
complete AXI conformance, registry-complete help, contextual disclosure, or even the
declared seven-command scope as a set. Ticket #9 must choose proportional contract
rows and independent mutations after tickets #2, #3, #7, and #8 settle the public
bytes and scope.

## Sources

- [AXI: the ten principles](https://axi.md/), read 2026-08-09.
- `.agents/skills/bench-craft-cli/SKILL.md`
- `projects/benchkit.md`
- `cmd/bench/main.go`, `cmd/bench/command_registry.go`, and
  `cmd/bench/specbuild.go`
- `internal/toon`, `internal/usage`, `internal/spec`, `internal/specbuild`,
  `internal/sessioninspect`, and the cited command packages
- `internal/conformance` registry, routing, docs-workflow, hook-wiring, and
  package-core checks
- `ROADMAP.md` FT173
