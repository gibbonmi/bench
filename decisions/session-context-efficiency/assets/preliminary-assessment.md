# Preliminary session context-efficiency assessment

Recommendation: make returned tool-result text the primary control surface and
require projection before aggregation. Treat call count as a separate,
secondary measure. Start design work with bounded read/query results and a
set-wide worktree cleanup operation; do not spend the first build on polling or
quiet archive calls.

Scope: the root Codex transcript from the `$bench-drain` request through the
last tool result before the reviewer requested this assessment. The slice is
JSONL lines 11–2786 in
`/home/mgibs/.codex/sessions/2026/09/11/rollout-2026-09-11T03-26-10-01a08f5c-118f-7620-816b-adc8596406c5.jsonl`.
The SHA-256 of raw lines 11–2786 is
`eb92f9af6148368ec018c98908704bb72bfcbe1e5e0f23424b6664a51bbfc89e`.
Evidence status: observed for counts and call shapes in that slice; inferred for
their design implications; proposed for interventions. This is one long Codex
session, not a cross-harness baseline.

Consumed by: `decisions/session-context-efficiency.md` and its tickets.
Drift: refresh when the transcript schema, tool wrapper, or measured surfaces change.
Retire when: the map's resulting specs absorb this evidence or a later census supersedes it.

## Question graph

1. **How concentrated was tool-result text?** Independent; answered below.
2. **Which repeated shapes represented one domain intent?** Independent;
   preliminarily ranked below.
3. **Which existing artifacts own output, measurement, and guidance?**
   Independent; answered below.
4. **Which projection and aggregation changes should ship?** Depends on 1–3;
   left to map tickets 3 and 4.
5. **What can every harness measure?** Extends 1 across harnesses; left to map
   ticket 5.
6. **What denies, advises, or only measures?** Depends on 4 and 5; left to map
   ticket 6.

## Facts: the tool-result-text distribution

The census joins every tool call to its matching output. It counts the output
text once and stops immediately before the request that initiated this map. It
does not count reasoning or assistant prose as tool output.

### Census method

- Parse JSONL lines 11–2786, inclusive, in file order.
- Index `custom_tool_call` and `function_call` payloads by `call_id`.
- Join the matching `custom_tool_call_output` or `function_call_output` payload.
- For a list output, concatenate each dictionary block's `text` in array order.
  For another output type, use its string representation.
- Count characters with Python `len(text)` and bytes with
  `len(text.encode("utf-8"))`.
- Count result lines as `text.count("\n") + bool(text)`.
- Sort per-result character counts. The 90th percentile is the value at
  `floor(0.9 * (N - 1))`, which is index 289 for 323 results.
- Mark a result truncated when its text contains `truncated output`, ignoring
  case.
- Match diagnostic command families against the serialized call input or
  arguments. Families overlap and are not additive.

| Signal | Observed value |
| --- | ---: |
| Outer tool calls | 323 |
| `exec` calls | 255 |
| Collaboration calls | 68 |
| Tool-result text | 1,050,384 characters / 1,052,380 UTF-8 bytes |
| Returned lines | 11,965 |
| Median result | 320 characters |
| 90th-percentile result | 8,893 characters |
| Results at least 10,000 characters | 31 calls, 768,697 characters (73.2%) |
| Results at least 20,000 characters | 17 calls, 563,733 characters (53.7%) |
| Results at least 40,000 characters | 7 calls, 280,729 characters (26.7%) |
| Results carrying a truncation warning | 9 calls, 332,835 characters (31.7%) |

The distribution is top-heavy. Cutting many small calls cannot match the gain
from preventing a few 30,000–40,000-character results. Seven near-cap results
alone account for more than one quarter of all tool-result text. They include
whole-rule reads, batched source reads, one historical spec, and a roadmap batch
(transcript lines 18, 43, 59, 383, 399, 455, 1335, and 1374).

The transcript's last usage record reports 43,760,211 cumulative input tokens.
Of those, 43,290,624 are cached. It also reports 91,546 output tokens and 39,099
reasoning-output tokens (transcript line 2777). Those are provider counters
accumulated across model responses.

They do not attribute input tokens to tool
results. This record alone cannot yield a tool-output cost. One compaction
occurred in the measured slice (transcript line 1160).

## Facts: recurring output patterns

These diagnostic categories overlap and therefore do not sum to the total.

| Pattern | Calls | Returned characters | What happened |
| --- | ---: | ---: | --- |
| Whole-file or whole-history reads | 20 | 356,004 | `perl ... print`, empty-pattern `rg`, and complete `git show <rev>:<path>` acquired bodies before selecting needed evidence. |
| Git status/diff snapshots | 31 | 129,809 | Repeated snapshots often included full or multi-worktree detail when only a verdict or changed-path set was needed. |
| `bench roadmap --context` family | 4 | 123,811 | Two broad snapshots and later row selections returned complete large bodies; row selection reduced scope but did not impose a result budget. |
| `bench worktree list` family | 7 | 34,512 | The nine-field population table and up to two help actions per active row scale with assignment count. |
| Aggregating exec wrappers | 17 | 229,446 | Several outer calls correctly grouped independent commands but concatenated every raw child result. |
| Polling (`wait_agent` plus `write_stdin`) | 66 | 12,585 | The largest repeated call family produced about 1.2% of tool-result text. |
| One-at-a-time bundle/log preservation | 8 | 1,207 | Four bundles and four log copies were independent, visible repetitions with negligible text volume. |

Five patterns explain the failure:

1. **Fetch then filter.** The caller retrieved complete skills, source files,
   historical specs, diffs, or roadmap bodies and selected evidence afterward.
2. **Line limits substituted for byte limits.** A bounded line count still
   admitted very long roadmap, idea, table, and diff lines. The current research
   reproduced this failure with a 180-line search result that was still
   truncated.
3. **Aggregation lacked projection.** Parallel fan-in lowered the number of
   visible calls while returning the concatenation of all child outputs. This
   is why call grouping and context reduction must be measured separately.
4. **Snapshots replayed unchanged fields.** Population, roadmap, status, and
   diff queries returned complete snapshots across several checkpoints. Some
   repeated reads were necessary after mutations; their full schemas were not.
5. **High caps became a target.** Nine truncations still returned 332,835
   characters. A transport cap prevents unlimited output but does not supply the
   summary, total, omitted count, or detail route needed for the next decision.

These findings agree with AXI's existing requirement for compact output,
minimal default schemas, sized previews, and pre-computed aggregates
(`.agents/skills/bench-craft-cli/SKILL.md:14-47`). The session behavior failed
to consume that design consistently, and several raw shell/file surfaces sit
outside the approved AXI query set.

## Facts: repeated call sequences

The clearest orchestration-only miss was archive preservation: four independent
`git bundle create` calls followed by four independent `.logs` copies
(transcript lines 2161–2176 and 2190–2205). One orchestrator call could have run
each group concurrently and returned a compact eight-row verdict. This saves
calls but little text.

The stronger product-level candidate is bulk worktree cleanup. One wrapper
resolved 26 paths and planned target cleanup in parallel. Applying the first
plan changed repository-wide state. That invalidated later per-target
fingerprints, so the run re-planned and applied targets sequentially
(transcript lines 2230, 2262, 2294, 2326, and 2360). This is one user intent
whose invariant spans the set.

`bench worktree reclaim` already demonstrates a set-wide plan, one
aggregate fingerprint, per-key verdicts, and one apply action
(`internal/worktree/pool_reclaim.go:289-352`). A bulk cleanup design should
reuse that ownership shape rather than wrap N independent transactions.

Two smaller query candidates also appeared. The coordinator resolved 26
worktree paths through child calls inside one wrapper. It queried three spec
histories the same way (transcript lines 2230 and 2441). A selector or
multi-target query could return the needed columns once.

The worktree population query owns nine fields. It derives two actions for each active row
(`internal/worktree/list.go:15-17,89-131`), so a target/path projection must be
designed with that owner rather than parsed from its stdout.

Build-then-run is a possible fourth candidate. The shell follow-on guard
correctly rejected a chained build and probe, after which the run used separate
calls. A domain command could own that dependency, but the current session does
not show enough repetition to rank it above cleanup or query selection.

Polling, gate-to-land transitions, and approval boundaries are poor first
targets. Polling is frequent but small, and the other sequences cross changing
state or authority. A lower call count would not justify hiding their separate
outcomes.

## Current owners and overlap

FT231 already owns performance instrumentation, the split between Bench and
harness measures, a deep raw-call census, and the session context-cost audit
(`roadmap/FT231.md:21-25,61-69,84-94`). This map supplies the requested concrete
session evidence; it should amend or feed FT231 rather than create another token
ledger.

FT173 owns AXI query output and records that its four truncation policies have
different caps, units, metadata, and `--full` behavior. Shared mechanics may be
parameterized, while the policies stay with command domains
(`decisions/byte-preserving-axi-foundation/ft173-axi-contract.md:118-140`).
`bench worktree list` is explicitly in FT173's scope
(`decisions/byte-preserving-axi-foundation/ft173-axi-contract.md:16-22`). Any
default-schema or help-row change belongs there or in an explicit amendment.

FT100 owns always-loaded guidance weight and already requires evidence before
cuts (`roadmap/FT100.md:1-25,30-55`). FT89 supplies the placement rule: a
deterministic operation belongs in the CLI, while judgment remains with the
agent (`roadmap/FT89.md:31-45`). These owners cover guidance reduction and new
commands; this map decides which observed gaps they receive.

## Tested results

- An independent parser replay reproduced 323 calls, 1,050,384 characters,
  11,965 lines, and the 73.2% upper-tail share from the cited transcript slice.
- `bench maps` accepts the map and resolves its blocker graph.
- `bench gate-prose` accepts the map, report, tickets, and glossary changes.
- No byte-compatibility or cross-harness behavior has been tested.

## Inferences

- Tool-result text is a concentration problem. The first useful control should
  target the upper tail, not an average-call budget.
- Call consolidation improves context only when the owner returns a projection
  or aggregate. Fan-in of raw child results can make the primary cost worse.
- A transport-level truncation marker is too late. The caller has already paid
  for a large result and may need another call because the useful tail was cut.
- Repository query changes alone cannot bound raw shell, Git, and file reads.
  The final design needs separate CLI, harness-wrapper, and operating-guidance
  decisions with one metric owner.
- Bulk cleanup has higher leverage than archive batching because it removes
  both repeated calls and a stale-plan failure mode.

## Preliminary proposals

1. **Bound the return contract.** For any potentially large result, return the
   requested projection, exact total size/count, omitted count, and one explicit
   detail route. A wrapper-level overflow should spill complete output to a
   named artifact and return a structured summary, not the truncated prefix.
2. **Make selectors cheap.** Add or use file ranges, failure-only diagnostics,
   changed-path/status projections, roadmap row selection with bounded bodies,
   and target/field selection for population queries.
3. **Aggregate facts, not text.** An orchestrator batch should emit one row per
   child with status and a bounded diagnostic. Complete child stdout stays in a
   referenced artifact unless explicitly requested.
4. **Give set-wide intent one transaction.** Explore bulk worktree cleanup with
   a set fingerprint and per-target verdicts, using `worktree reclaim` as the
   precedent. Batch archive preservation and multi-target lookups are secondary
   candidates.
5. **Measure before denying.** Record tool-result bytes and upper-tail results
   first. Deny only mechanically unambiguous forms after ticket 5 establishes
   reliable harness sources and ticket 6 sets the false-positive posture.

No numeric output budget is proposed from one session. The observed 10,000-
character concentration is a diagnostic cut, not a policy threshold.

## Contradictions and unknowns

- The transcript exposes result text and cumulative token records, but it does
  not attribute cached or uncached input tokens to individual tool results.
- The JSONL schema is harness-owned and may change; no stable cross-harness
  capability record for these fields was assessed.
- The census counts outer tool calls. One `exec` wrapper may issue several child
  commands, so a child-call metric needs explicit instrumentation rather than a
  regex over JavaScript input.
- Parallel user steering and one compaction make phase attribution ambiguous.
- Some repeated snapshots followed real state changes. The defect is their
  unneeded payload, not the fact that state was re-read.
- The prior drained idea listed eighteen candidate behaviors. They included
  unfiltered test output, redundant traversals, repeated Git snapshots,
  unbounded delegates, stack traces, formatting calls, ANSI output, and hidden
  payload attribution. FT231 retains the measurement request. This census has
  not tested all eighteen behaviors or any external candidate tool.

## Verification record

- [x] The report opens with recommendation, scope, and evidence status.
- [x] Facts, tested results, inferences, and proposals have separate sections.
- [x] Questions 1–3 have synthesized sections; questions 4–6 are explicitly deferred.
- [x] Comparison tables record each observed option's consequence. No product option is resolved here.
- [x] Contradictions and residual unknowns are retained.
- [x] No diagram is needed; the numbered dependency graph remains legible.
- [x] Unverified measures and harnesses are named.
- [x] Exact local sources accompany the load-bearing claims.
- [x] The report ends with a validation plan.

## Validation plan

- Run the same field-based census over one comparable phase from each supported
  harness. Preserve absent fields as absent.
- Replay the top five large-result shapes with candidate projections. Compare
  total tool-result bytes and the 90th percentile first, then call count, task
  success, and extra recovery calls.
- Prove a bounded-result implementation with a fixture whose useful evidence is
  beyond the cap: the summary must give the true total, omitted amount, and a
  working detail route.
- Prove bulk cleanup with mixed clean, ignored, stale, and retained targets. One
  plan fingerprint must cover the full set; a stale apply must remove nothing
  and name one re-plan operation.
- Keep full outputs as evidence artifacts during evaluation so a smaller default
  cannot hide missed failures.
