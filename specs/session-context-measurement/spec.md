# Measure session context and propose result budgets

Status: staged

Decision source: `specs/session-context-efficiency/decisions/session-context-efficiency.md` (ready compiled map).

Verification log: 2 iteration(s) to accept — Sol/high accepted the suite. Trace-only partials and review-state bookkeeping are folded.

Coordinator: [Session context efficiency](../session-context-efficiency/spec.md)

## Problem

Current evidence mixes directly counted result text with harness-dependent counters.
A new output budget needs comparable task outcomes and a clear account of unavailable facts.

## Solution

Add one opt-in record inspection mode to the existing harness command.
Use its observations to produce a cited budget report for FT231.
The reader supplies FT204 facts without introducing a general session-query namespace.
The report proposes budgets for later reviewer approval.

## User stories

Line: gpt-5.6-terra / high.
The source fixes the outcome, but the seams require careful evidence and compatibility work.

1. As an agent, I want direct result-text measures, so that optimization starts with observed bytes.
2. As an agent, I want UTF-8 byte counts, so that character counts cannot misstate a byte budget.
3. As an analyst, I want separate invocation counts, so that nested execution and completion polls cannot inflate the census.
4. As an analyst, I want absent, empty, and unsupported inputs distinguished, so that unknown facts cannot look like zero.
5. As an analyst, I want incomplete evidence identified, so that malformed events cannot produce false precision.
6. As an analyst, I want native usage semantics retained, so that cumulative counters cannot become additive events.
7. As a reviewer, I want token attribution limits stated, so that bytes cannot become invented provider costs.
8. As an analyst, I want observed compactions, so that context effects remain separate from call volume.
9. As an agent, I want existing harness views preserved, so that optional observation cannot alter compiled declarations.
10. As an analyst, I want provenance with each measure, so that another reader can reproduce the observation.
11. As a reviewer, I want task success and recovery costs, so that a smaller result must still support the task.
12. As a reviewer, I want separate budget approval, so that measurement cannot select new defaults.
13. As an analyst, I want existing Bench metrics reused, so that one fact retains one owner.
14. As an agent, I want records treated as data, so that an inspected transcript cannot execute instructions.
15. As an analyst, I want missing result boundaries exposed, so that transcript visibility cannot imply complete producer capture.

16. As an analyst, I want explicit turn observations, so that unavailable turn boundaries cannot become invented counts.
17. As an analyst, I want explicit read-path observations, so that observed reads remain separate from shell-text guesses.

## Implementation decisions

Proposed grammar: `bench harnesses <harness> --record <path> --format <source-id>`.
Both flags require one named harness.
Existing zero-flag invocations remain compiled, terminal reads.
The opt-in view reads exactly one explicit regular file and never discovers session directories.
It does not change the compiled capability record from one local observation.

The first source ID is `codex-rollout-2026-09-11`.
Its schema derives from the pinned assessment transcript and minimized fixtures from that source.
A source ID pins event shapes, text boundaries, invocation identities, and native counter semantics.
Claude and OpenCode inputs remain explicitly unsupported until separate pinned source evidence supplies their mappings.
An unsupported format returns exit 1 with unknown observations.
Usage errors return exit 2.

The reader returns typed observations, not pre-rendered transcript excerpts.
Its metric inventory owns the rendered row count and names.
Each row has `metric`, `value`, `unit`, `availability`, `source`, and `boundary` fields.
The record identity includes its digest, source ID, harness, and observed interval.
Availability distinguishes observed, incomplete, and unknown.
Only an observed zero has a numeric zero value.

Required metrics include result-text bytes, result-text characters, result-text lines, outer calls, observed nested calls, and unmatched calls.
Native token dimensions, compactions, turns, and explicit read paths retain their source-specific boundaries.
Each token dimension has its own observation, including input, cached input, output, and reasoning.
A missing dimension stays unknown.

Shell text is not a reliable read-path census and receives no speculative parser.
Cumulative snapshots use their documented identity and interval semantics rather than a sum.
The command never estimates provider dollars or tokens from byte counts.

The report uses the existing Bench census and span readers for Bench facts.
It records those facts separately from harness observations.
The reader does not alter the existing OTEL span format or the SessionStart inspector.

The second ticket produces budget evidence through actual representative tasks.
Its case inventory covers raw file, Git, test, and shell reads, Bench queries, and supported wrapper results.
Each case names its task, input digest, harness version, baseline, candidate bytes, task verdict, and recovery calls.
Large histories, long diagnostics, multibyte text, empty output, and failed commands join the case inventory.
Unavailable cases remain unavailable and cannot justify a default.
The report names each proposed surface owner and retains current policy until approval.

## Testing decisions

The command tests exercise the real owner with controlled files and records.
The ordinary Go test phase executes new package tests through its existing package census.
The conformance phase checks command inventories and ticket ownership.
No test-only helper crosses a package boundary.
Planned test names below identify required tests, not tests that already exist.

### Seam diagram

```text
named input -> existing command owner -> typed producer -> projected result
                       ^                       ^
                 command tests          controlled records
```

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
| --- | --- | --- | --- | --- |
| ME1 | 1, 2 | The record view counts UTF-8 bytes from model-visible text fields only | planned TestObservedTextBoundary in internal/harnesses | Counting serialized metadata or character length changes the multibyte fixture total |
| ME2 | 3 | The report separates outer calls from observed nested calls | planned TestObservedCallAncestry in internal/harnesses | A nested-only event cannot increase the outer-call count |
| ME3 | 3 | A repeated completion event contributes once per invocation identity | planned TestObservedDuplicateCompletion in internal/harnesses | Counting polling completions twice changes the known invocation total |
| ME4 | 4 | Unsupported source formats report unknown measures at exit 1 | planned TestObservedUnsupportedFormat in internal/harnesses | Guessing a parser produces unsupported numeric facts |
| ME5 | 4 | A valid empty record reports zero observed calls and bytes | planned TestObservedEmptyRecord in internal/harnesses | Conflating empty with absent hides the known empty observation |
| ME6 | 4 | An absent record reports a read error at exit 1 | planned TestObservedAbsentRecord in internal/harnesses | An empty success would falsely certify a missing input |
| ME7 | 5 | Malformed relevant events mark affected measures incomplete | planned TestObservedMalformedEvent in internal/harnesses | Silently skipped events would certify an understated complete total |
| ME8 | 6 | Provider counters retain their source semantics and availability | planned TestObservedNativeUsage in internal/harnesses | Summing cumulative snapshots inflates the canned token total |
| ME9 | 7 | The report leaves per-result token attribution unknown | planned TestObservedNoTokenAttribution in internal/harnesses | A bytes-to-token estimate would appear as an observed measure |
| ME10 | 8 | Compaction counts come only from identified compaction events | planned TestObservedCompactions in internal/harnesses | A short context window alone cannot increment the count |
| ME11 | 9 | Existing compiled harness views match the baseline for every registered harness | planned TestObservedPreservesCompiledViews in internal/harnesses | A differential run catches changes to the current terminal views |
| ME12 | 10 | Each reported measure includes its source boundary and availability | planned TestObservedMeasureProvenance in internal/harnesses | An omitted provenance cell makes unknown and observed facts indistinguishable |
| ME13 | 11 | The budget report compares task success on matched representative cases | review-owned: Spec axis checks budget-evidence.md against the case records | Smaller output alone cannot satisfy a failed task case |
| ME14 | 11 | The budget report counts full-detail recovery calls for each candidate | review-owned: Spec axis checks recorded follow-on invocations | A candidate that hides required detail cannot conceal its recovery cost |
| ME15 | 12 | The budget report proposes per-surface bytes without approving any default | review-owned: budget evidence and reviewer decision record | Treating the historical diagnostic cut as policy fails the approval boundary |
| ME16 | 13 | Bench-owned elapsed and census facts remain separately sourced | review-owned: evidence provenance inspection | Parsing harness text as a replacement Bench metric creates a second owner |
| ME17 | 14 | Record contents never execute as shell commands | planned TestObservedHostileRecord in internal/harnesses | A command-shaped string produces data without creating its sentinel file |
| ME18 | 15 | The report identifies unobserved tool results without claiming complete producer output | planned TestObservedMissingResult in internal/harnesses | An unmatched call or upstream truncation cannot become a complete-byte claim |
| ME19 | 14 | Non-regular or symlink record inputs refuse before the first content read | planned TestObservedRegularFileBoundary in internal/harnesses | A FIFO or linked record cannot block the reader or supply unowned bytes |
| ME20 | 6 | Input-token observations match the pinned source field availability and value | planned TestObservedInputTokens in internal/harnesses | Always returning unknown fails a fixture with an explicit input counter |
| ME21 | 6 | Cached-input-token observations match the pinned source field availability and value | planned TestObservedCachedInputTokens in internal/harnesses | Omitting cache usage fails a fixture with an explicit cached-input counter |
| ME22 | 6 | Output-token observations match the pinned source field availability and value | planned TestObservedOutputTokens in internal/harnesses | Omitting output usage fails a fixture with an explicit output counter |
| ME23 | 6 | Reasoning-token observations match the pinned source field availability and value | planned TestObservedReasoningTokens in internal/harnesses | Treating reasoning as universally absent fails a source that supplies it |
| ME24 | 16 | Turn observations match the pinned source field availability and value | planned TestObservedTurns in internal/harnesses | Omitting turn availability fails the pinned event inventory |
| ME25 | 17 | Explicit read-path observations match the pinned source field availability and value | planned TestObservedReadPaths in internal/harnesses | Inferring shell paths or dropping observed read paths fails the pinned tool inventory |

### Edge inventory

ME4–ME7 cover unsupported, empty, absent, and malformed records for kit and linked-repository callers.
ME1 covers multibyte text and newline accounting.
ME2, ME3, and ME18 cover nested calls, repeated polls, unmatched events, and upstream truncation.
ME8 covers repeated snapshots, missing counters, and source resets within the observed interval.
ME17 covers control-bearing paths and command-shaped record fields as inert data.
The command uses no package-variable substitution across a subprocess boundary.

Won't handle: automatic session discovery — the agent supplies the explicit record path.
Won't handle: unsupported harness decoding — the record view reports unknown until a pinned adapter has reviewed evidence.
Won't handle: provider pricing — FT231 can consume independent pricing evidence in a separate capability.

## Ownership fences

- `.agents/skills/bench-craft-cli/SKILL.md`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `cmd/bench/main.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `internal/harnesses/command.go`
- `internal/harnesses/command_test.go`
- `internal/harnesses/observed_test.go`
- `internal/harnesstranscript`
- `reviews/session-context-measurement.md`
- `specs/session-context-measurement/assets/budget-evidence.md`
- `tests/canary/package-core-guard/unrouted-subcommand`

Reviewer disposition: Sol/high review accepted; user spec and ticket sign-off remains pending.
The fence is the union of ticket writes and the review pickup.
A build cannot change this spec, its acceptance rows, or its tickets.

## Ticket graph

| Ticket | Blocked by | Delivered coverage |
| --- | --- | --- |
| [1. Inspect one pinned harness record](tickets/1-inspect-record.md) | none | ME1, ME2, ME3, ME4, ME5, ME6, ME7, ME8, ME9, ME10, ME11, ME12, ME17, ME18, ME19, ME20, ME21, ME22, ME23, ME24, ME25 |
| [2. Compare representative budget cases](tickets/2-compare-budget-cases.md) | 1-inspect-record.md | ME13, ME14, ME15, ME16 |

## Out of scope

Automatic session search is a separate FT204 capability: approximately 5 edits, 1 gate run.
Provider pricing is a separate FT231 capability: approximately 3 edits, 1 gate run.
Neither capability blocks the explicit record view.

## Further notes

### Source trace

| Source clause | Coverage |
| --- | --- |
| Ticket 1: measured tool-result census | ME1, ME2, ME3, ME18 |
| Ticket 2: result text before call consolidation | ME1, ME13, ME14 |
| Ticket 5: documented, observed, and unknown measures | ME4–ME12, ME16, ME20–ME25 |
| Ticket 6: measure and advise | ME15, ME17 |
| Ticket 9: representative cases, task success, and recovery calls | ME13–ME15 |

### Reader sweep and proof checklist

[Seam evidence](../session-context-efficiency/assets/seam-evidence.md) records the definition reads and enforcement closure.
Cited symbols: the existing `harnesses.Command` is the production entrypoint.
Import edges: the command imports the new typed reader within the same Go module.
Source-row clauses and occurrences: the sole compiled map owns the quoted clauses above.

Promised field labels: `metric`, `value`, `unit`, `availability`, `source`, and `boundary` belong to the new observed view only.
Changed-function callers: the command registry wraps `harnesses.Command`; its command tests and AXI inventory consume that result.
Copy survival: no existing reader moves or disappears.

The compiled record consumers in guards, lines, and status keep their existing API and values.
Their exclusion follows ME9's differential contract.
The existing OTEL reader and SessionStart inspector remain unchanged under ME16.

### Flagged additions

The opt-in `harnesses` grammar is a proposed command-owner extension for spec approval.
It preserves the existing compiled views and does not replace FT173 policy.
ME1–ME12, ME17–ME25 grade the new mode.
The evidence ticket is an independent, review-owned deliverable after the reader ships.
