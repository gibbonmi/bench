# FT173 helper and byte-compatibility census

Evidence asset for `decisions/byte-preserving-axi-foundation/ft173-axi-contract.md` ticket #3. Inspected
2026-08-09 at `974020e4af8de5ed75098c4c5934a8907952bb2b`. This is a
read-only source census: no command behavior was changed and no test or gate command
was run.

## Compatibility boundary

The published [AXI contract](https://axi.md/) requires bounded content with a size
hint and `--full`, plus pre-computed totals and explicit empty results. It also
requires structured stdout errors with honest exits, live content at query homes,
contextual `help[]`, and consistent leaf help. Bench's current local contract
applies AXI to selected surfaces
and gives operational commands a different stderr/exit contract
(`.agents/skills/bench-craft-cli/SKILL.md:7-16,45-75`). FT173 therefore cannot use
"AXI conformance" as permission to change an existing operational stream. Ticket #4
is the explicit exception that makes the complete spec-build family an AXI surface.

For this census, **byte-preserving** means preserving the complete observable
contract for every previously accepted and rejected argv:

- exact stdout and stderr bytes, including block order, field order, quoting,
  whitespace, and final newlines;
- exit status and the distinction between help, usage, refusal, success, and no-op;
- accepted spellings, flag repetition, `--` behavior, and default versus `--full`
  behavior; and
- each truncation policy's cap, counting unit, suffix, metadata, sorting, and safety
  posture.

An internal type or call-site move can satisfy that boundary. Adding a flag is not
strictly byte-preserving over the argv domain: the new spelling previously failed
with usage. Adding a field, `help[]` block, family home, structured refusal, or newly
included patch body necessarily changes public output even when existing fields keep
their values.

## Current owners and consumers

| Fact | Current derivation owner | Current rendering and consumers | Compatibility classification | Focused proof that must bite |
|---|---|---|---|---|
| Flat structured success | `internal/toon.table` owns the exact flat-table TOON adapter, including schema-bearing zero rows, cell typing/escaping, count header, indentation, and trailing newline. Command packages own schemas, row order, and block composition (`internal/toon/toon.go:22-80`). | `outputCommand` only forwards the returned bytes and code (`cmd/bench/main.go:120-125`). There are 43 non-test `Table`/`TableTyped` calls in 21 files; spec build renders typed lifecycle values in `cmd/bench/specbuild.go:128-238`. | Routing existing renderers through a shared result envelope is byte-preserving only while the same command-local schemas, row order, block order, and TOON adapter remain. Registry-declared schemas, minimal-field defaults, `--fields`, or a new envelope are output-changing. | `TestTableCellEscaping`, `TestTableUnrepresentableCellErrors`, `TestRepresentableMatchesEncoder`, `TestTable`, and `TestTableTyped`. Mutate the empty special case, final newline, encoder options, string-versus-number type, or block ordering; the focused owner or command test must fail. |
| Content truncation | Four owners: `sanitize.Preview`, `roadmap.limited`, `worktree.inventoryIgnored`, and outline's row projection. Their policies are not interchangeable; the detailed ledger below records their distinct units and metadata. | Eleven `sanitize.Preview` calls span status, test report, shift, worktree refresh, and the compiled wrapper. Nine roadmap parser call sites consume `limited`. Worktree cleanup consumes the ignored inventory for both display and safety/fingerprints. Outline emits bounded rows plus metadata. | A parameterized bounded-projection mechanic can be byte-preserving. Replacing all policies with `sanitize.Preview`, double-truncating capture-bounded values, changing units/caps, or changing suffix/metadata is not. Adding a preview where content is currently omitted is output-changing. | The four owner tests in the truncation ledger, plus a paired before/after fixture for every current call-site family. Mutations must independently cover cap, unit, UTF-8 boundary, `--full`, total size, omitted count, and safety scan. |
| Aggregates | `toon.table` derives emitted row count only. Total, inspected, omitted, unknown, status, and lifecycle counts are command-local typed facts: guard scan, outline metadata, roadmap context/drain, ignored inventory, spec-build abandon/reclaim, publication state, status, and dashboard. | Renderers serialize owner-supplied values. `status.Signals`, `roadmap.DrainCounts`, spec-build plans/status, and publication records are consumed as Go values before rendering; no generic renderer can infer their meaning from displayed rows. | A shared aggregate carrier/renderer can preserve bytes if it accepts owner-derived facts and retains names/order. Recomputing total from emitted rows, treating unknown as zero, or defining completeness in the renderer changes meaning even if the block still parses. | `TestCommandAlwaysEmitsCompleteGuardScanMetadata`, `TestScanTimeoutPreservesPartialRowsAndHonestCounts`, `TestScanEnumerationTimeoutUsesUnknownCounts`, `TestCommandBoundsRowsAndFullRetainsMetadata`, roadmap context tests, and `TestReclaimReceiptReportsEveryClassIncludingTheEmptyOnes`. Mutate total to `len(rows)`, force incomplete scans complete, coerce unknown to zero, or drop an empty class; one focused test must fail per mutation. |
| Structured errors | `toon.Errorf`, `RenderError`, `Usage`, `MissingArg`, `NotInRepo`, and `RecordError` own common line shapes (`internal/toon/toon.go:110-150`). `buildError` sanitizes spec-build refusal and hint text into one stdout line (`cmd/bench/specbuild.go:124-126`). Other direct/system handlers still own their sinks and prose. | Query commands return stdout/code pairs. Spec build returns sanitized stdout/1 for lifecycle refusal. The compiled root dispatcher emits absent/unknown subcommands on stderr/2 (`cmd/bench/command_registry.go:37-60`). Publication prints a render failure but its `printRecord` helper cannot return a changed exit (`internal/publication/command.go:320-336`). | Consolidating identical line formatting and returning a typed error before rendering can preserve bytes. Moving an existing stderr error to stdout, replacing prose with a TOON block, appending contextual help, or changing the hint is output-changing and belongs to an explicitly migrated surface. | `TestErrorfUsage`, `TestSpecBuildErrorCannotSplitStructuredOutput`, `TestRunUnknownExits2`, and command-level outside-repo/usage tests. Mutate stdout/stderr, exit 1/2, kind, hint, control escaping, or the single-line guarantee. A formatting-only unit test is insufficient for sink and exit; retain a command-level probe. |
| Empty state | `toon.table` owns `name[0]{fields}:`. Commands independently choose a zero-row table, prose, or a one-row state object. Spec build deliberately represents no run as `spec_build[1]` with `state=empty`, and `--full` adds empty assignments/review blocks (`internal/specbuild/state.go:554-566`; `cmd/bench/specbuild.go:148-179`). Status separately emits `bench: clean — nothing pending`. | Agents and tests read the rendered result. Internal composition does not infer absence from the stdout bytes: spec-build and status state exist as typed values first. | Reusing the TOON empty renderer preserves bytes for existing zero-row tables. Replacing a one-row/prose empty with a zero-row table, or making a formerly silent command explicit, is output-changing even though it may improve AXI conformance. | `TestTable`, `TestSpecBuildStatusRendersDefinitiveEmptyProjection`, `TestStatusHasDefinitiveEmptyAndActiveProjections`, `TestRenderClean`, and `TestRenderEmptyStates`. Mutate nil rows to blank output, drop the schema, omit spec-build's empty row, or collapse unknown/unreadable to empty. |
| Argument grammar and exits | `usage.Parse` owns shared help spellings, declared flags, repeated-flag refusal, required values, arity, empty positional, `--`, and success/help/usage codes for 23 call sites (`internal/usage/parse.go:47-149`). Nested/manual parsers and the root wrapper retain separate grammars. `Command.Run` is the final exit dispatcher. | Callers print the non-empty parser line and return its code. No caller reparses that text. Root, nested, and shell dispatchers consume argv directly, so byte compatibility includes the accepted-language boundary as well as output. | Moving an existing grammar unchanged into `usage.Parse` can preserve bytes only after a paired accepted/rejected-argv delta proves parity. Adding `--fields`, bare help, a family home, or changing missing/unknown classification is behavior-changing. | `TestParseHelpSpellings`, `TestParseBareHelpOnlyWhenSole`, `TestParseMissingFlagValueDistinctFromUnknownFlag`, `TestParseArityUnmetIsMissingArg`, `TestParseExcessPositionalNamesFirstExcess`, `TestParseSuccess`, repeated-flag tests, and root/spec-build dispatcher tests. Mutate one spelling, repeat handling, terminator behavior, sink, or code. |
| Contextual action/help | There is no production `help[]` owner or emitter. Status `action`, spec-build `next`, publication `next_action`, error hints, and wrapper examples are separate prose derivations. Spec-build derives `Status.Next` once, but values include both invokable commands and orchestration labels such as `release assignment`, `delegate assignment`, and `retry promote` (`internal/specbuild/state.go:74-114`). | Handoff consumes typed status signals and filters invocable commands with `IsInvocable`; it does not parse rendered `bench status` (`internal/handoff/facts.go:159-205`). Publication and spec build render their own strings. Agents are the public consumer. | Introducing an internal typed action with fixed arguments versus open placeholders can be byte-preserving before it is rendered. Emitting `help[]`, replacing `action`/`next`/`next_action`, or changing an error hint is necessarily output-changing under the roadmap's explicit exception (`ROADMAP.md:735-746`). Existing prose must not be copied blindly. | The future helper needs an exact table test plus command-level tests from `ft173-command-help-inventory.md`. Mutate by deleting a useful row, replacing a known slug/id/fingerprint with a placeholder, guessing an unknown runtime value, dropping a carried flag, or converting a non-command label into an advertised command. Each mutation must fail independently. |

## Truncation policy ledger

The roadmap requires four derivations to route through one shared call without changing
bytes (`ROADMAP.md:735-746`). That is possible only if the call is a policy-parameterized
mechanic, not a fifth universal policy. It must return already-bounded content and the
owner's total/emitted/omitted facts without choosing semantic units on the owner's
behalf.

| Owner | Current policy tuple | Semantic consumers | Byte-preserving move | Exact independent mutation probe |
|---|---|---|---|---|
| `sanitize.Preview` | Escape controls/backslashes; cap at 120 Unicode code points; append `… (N bytes)` using the original byte length; no `--full` argument (`internal/sanitize/sanitize.go:41-56`). `sanitize.Controls` is the uncapped sibling. | Status details, shift objective/result, worktree-refresh detail, compiled prompt observation, and `bench test` default diagnostics. `bench test --full` explicitly switches to uncapped `Controls` (`internal/testreport/testreport.go:337-342`). | A shared projection may supply the first 120 runes, original byte count, and truncated bit while the current escaping/suffix renderer remains. It must not count bytes for the cap or re-truncate a value already bounded upstream. | `TestPreviewBoundariesAndControls` and `TestControlsEscapesWithoutCapping`. Mutate 120 to 119/121, runes to bytes, original bytes to visible bytes, remove control/backslash escaping, or make full diagnostics call `Preview`; each must turn a focused test red. |
| `roadmap.limited` | Cap raw content at 4096 bytes; back off to valid UTF-8; return original byte count and truncated bit; `--full` returns the complete raw string (`internal/roadmap/context_types.go:9,117-127`). It does not sanitize or append a suffix; renderers expose `*_bytes` and `truncated`. | Nine capture/context projections for roadmap rows, ideas, learnings, retros, parse failures, and source bodies. | Route the cap/count/full computation through the shared mechanic while preserving raw bytes and existing per-row metadata. Do not pass its result through `sanitize.Preview`. | `TestContextBodyLimitBoundaries` plus `TestBuildContextCarriesRetrosAndDegradedEvidence`. Mutate 4096, remove UTF-8 backoff, count runes, clear `truncated`, falsify original bytes, or make `--full` remain capped. |
| `worktree.inventoryIgnored` | Enumerate and sort ignored paths; reject unsafe/escaping/stat-raced entries; count entries and bytes while building fingerprint material; stop beyond 1000 entries and mark at-least/over-limit; default displays 20 paths, full up to 1000; `Truncated` reflects undisplayed or at-least state (`internal/worktree/subshell.go:389-439`; `internal/worktree/ownership.go:24-25`). | Worktree release/cleanup safety, preservation verdict, classification rendering, recovery fingerprint, and resume summary. This is authority-bearing inventory, not merely display truncation. | Only the final visible-count projection may use a generic bound. Enumeration, safety checks, total/at-least meaning, byte cap, canonical parts, and digest remain with worktree ownership. | `TestIgnoredInventoryEntryAndByteBoundaries`, `TestIgnoredInventoryStatRaceRetains`, and cleanup-policy tests. Mutate default 20, full 1000, byte limit, sort, at-least state, unsafe-path refusal, stat-failure retention, or digest inputs. A display-only golden cannot replace these safety probes. |
| `outline` | Discover all eligible symbols; default emits at most `bounds.OutlineRowLimit` (200), full emits all; metadata retains tracked/scanned/skipped/total/emitted/omitted/truncated; any skip also makes `truncated=true` (`internal/outline/outline.go:230-264`; `internal/bounds/bounds.go:34`). | `bench outline` agents consume bounded rows and exact completeness metadata. | Route row slicing and count result through the shared projection after outline has derived total symbols and skips. The helper must accept total independently from `len(rows)` because unrepresentable rows and skips affect completeness. | `TestCommandBoundsRowsAndFullRetainsMetadata` and skip tests. Mutate 200, derive total from emitted rows, ignore skips when setting `truncated`, lose omitted count, or keep bounding under `--full`. |

## Aggregate and empty-state boundary

The common abstraction may carry facts, but it must not become their semantic owner.
Three counterexamples show why:

1. `guard_scan` may know `total=unknown`; converting that to zero for a numeric helper
   fabricates completeness (`internal/guards/guards_test.go:63-75`).
2. Outline's `total_symbols` is not the table header count: bounds, skipped files, and
   unrepresentable symbols can separate the two (`internal/outline/outline.go:230-264`).
3. Spec-build reclamation reports every disposition class, including zero-count ones,
   because absence of a class is part of the receipt (`cmd/bench/specbuild.go:217-238`).

The same distinction governs empty state. A TOON zero-row table says "this collection
was read successfully and contains no rows." A one-row `state=empty` says "the named
lifecycle exists as a query target but has no run." Prose clean state is a human board
contract. A shared renderer must preserve those meanings rather than normalize all
three to one spelling.

## Runtime consumer census

Source inspection found no production TOON decoder and no production process that
invokes a Bench query command and parses its stdout back into semantic state. The only
non-test `toon-go` import is the encoder adapter in `internal/toon`. Current in-repo
composition follows typed or pass-through seams:

- `sessioninspect.statusPhase` calls `status.Command(nil)` in process and forwards its
  bytes; it does not split or decode them (`internal/sessioninspect/sessioninspect.go:73-78`).
- Handoff calls `status.SignalsWith`, selects an invocable typed action, and renders its
  own result (`internal/handoff/facts.go:159-205`).
- Dashboard gathers typed status, roadmap, and worktree projections
  (`internal/dashboard/dashboard.go:87-104`).
- Spec-build receives typed lifecycle status/plans from the service and renders them at
  the command boundary (`cmd/bench/specbuild.go:35-121`).
- Session hooks and wrapper cases execute or forward commands; they do not parse query
  TOON. Tests and external agents are therefore the current byte consumers. Raw patch
  output from `bench diff --full` is also intentionally consumable as a patch even
  though it has no in-repo parser.

This absence lowers internal migration risk but does not relax compatibility. Exact
tests, session logs, shell users, and agents are consumers. A decoder added later must
bind to a versioned schema rather than becoming an undocumented reason old bytes can
never move.

## Migration classification

### Byte-preserving foundation

The following can land without an output exception, provided a paired-delta run proves
the whole old argv matrix is identical:

- introduce a policy-parameterized bounded-projection result and route the four current
  derivations through it while retaining their policy tuples and owner metadata;
- introduce a shared aggregate carrier/renderer that serializes facts supplied by the
  current semantic owners, without recomputing totals or completeness;
- introduce typed success, refusal, empty, and action values behind current renderers;
- move matching manual grammars into `usage.Parse` only where accepted/rejected argv,
  help bytes, streams, and exits remain exact; and
- enrich the command registry with internal metadata that is not yet emitted.

The foundation is not allowed to consolidate by flattening distinctions. In particular,
it cannot replace roadmap/worktree/outline behavior with `sanitize.Preview`, turn
unknown into zero, or turn spec-build's empty state into an empty table. It also
cannot move operational stderr or advertise prose `next` values as commands.

### Necessarily output-changing slices

These require their own explicit compatibility contract and paired before/after
expectations:

- emitting `help[]` or replacing `action`, `next`, `next_action`, or refusal hints;
- changing an existing sink, error shape, exit, no-op result, or root/family no-argument
  behavior;
- adding `--fields`, reducing or widening a default schema, or changing block order;
- adding sized previews where content is currently omitted or unbounded;
- extending `bench diff` to a coherent live snapshot and including untracked regular-file
  bodies under `--full`; and
- migrating every spec-build success, empty, refusal, recovery, and family-home result to
  the full ten-principle envelope decided by ticket #4.

## Paired-delta and mutation contract

The compatibility harness for the foundation should execute the same selected Bench
executable against frozen fixtures twice—baseline and candidate—and compare stdout,
stderr, and exit exactly. It is a focused contract harness, not a whole-project gate.
Its matrix must include:

| Seam | Required paired cases | Independent mutation that must be observed |
|---|---|---|
| TOON success/empty | zero/one/many rows; numeric-looking string versus typed number; every quoting/control class; multiple concatenated blocks | Remove empty schema, final newline, quoting, representability refusal, row count, or block order. |
| Grammar/exit | success, no-op where defined, all help spellings, unknown/repeated/missing flag, missing/excess/empty positional, `--`, outside repo, lifecycle refusal | Swap 0/1/2; move sink; accept an unknown/repeated flag; collapse missing into unknown; change help bytes. |
| Each truncation policy | below/exactly/above cap; multibyte boundary; default/full; zero and over-limit; skipped/unrepresentable/unsafe source | Change cap/unit/total, double-truncate, make full incomplete, or remove the authority-bearing worktree refusal. |
| Aggregates | complete, partial, unknown, zero class, emitted less than total, skipped input | Derive total from visible rows, unknown as zero, complete from success exit, or omit zero classes. |
| Contextual actions | success, empty, stale, refusal, plan/apply, terminal, and recovery examples from the exhaustive command inventory | Delete action; fail to carry a known fixed value; guess an unknown value; carry a stale fingerprint; advertise an orchestration label as executable. |

For byte-changing slices, comparison is not equality. Each fixture needs an approved
old-to-new delta naming the exact added/removed blocks, stream/exit change, and action
template. Everything outside that approved delta remains byte-equal. A semantic review
can judge whether the new suggestion is useful. But only the executable paired-delta
and mutation probes prove the formatter, parser, and exit contract bite.

## Resulting seam judgment

Bench already has strong reusable mechanics for TOON bytes and common argv grammar,
plus typed semantic owners in status, roadmap, spec build, publication, and worktree.
It does not have one reusable truncation policy, aggregate semantics, empty-state type,
or contextual-help owner. The safe foundation is therefore a set of parameterized
carriers behind existing renderers—not a universal renderer that infers meaning.

That foundation can preserve bytes. Contextual disclosure, coherent Git inspection,
and the complete spec-build AXI migration cannot; they need separately approved deltas.
This is the compatibility boundary tickets #7 and #8 must carry forward.

## Sources

- `https://axi.md/` — published ten-principle AXI contract, especially principles 3
  through 6 and 8 through 10; checked 2026-08-09.
- `.agents/skills/bench-craft-cli/SKILL.md` — current local AXI scope, seven-principle
  guidance, output/error rules, and paired-delta expectation.
- `ROADMAP.md` — FT173 byte-preservation, contextual-disclosure exception, Git-owner,
  and no-double-truncation constraints.
- `decisions/byte-preserving-axi-foundation/assets/ft173-axi-surface-census.md` — command/call-site counts and
  ten-principle implementation census.
- `decisions/byte-preserving-axi-foundation/assets/ft173-command-help-inventory.md` — exhaustive action surface and
  exact known-value/placeholder requirements.
- `internal/toon`, `internal/usage`, `internal/sanitize`, `internal/roadmap`,
  `internal/worktree`, `internal/outline`, `internal/status`, `internal/handoff`,
  `internal/specbuild`, `internal/publication`, and `cmd/bench` — current derivations,
  renderers, consumers, and focused contract tests cited above.
