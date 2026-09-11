# Harness measures and result interception

Recommendation: use explicit capability records and preserve unavailable measures as unknown.
Keep transcript observations separate from documented interfaces and installed behavior.
Scope: Codex, Claude Code, OpenCode, and the current Bench measurement owners.
Evidence status: source-verified contracts and one observed Codex transcript; no runtime interception compatibility claim.

Consumed by: `specs/session-context-efficiency/decisions/session-context-efficiency.md`, especially tickets 5, 6, and 9.
Drift: refresh after a harness release, transcript schema change, integration change, or change to the named Bench owners.
Retire when: the resulting specs absorb this evidence or a later capability assessment supersedes it.

## Question graph

1. Which surfaces expose tool-result text?
2. Which records distinguish outer calls from nested tool calls?
3. Which sources expose input and cached-input tokens?
4. Which sources expose completed compactions?
5. Which sources expose reasoning-output measures?
6. Which owner can measure or replace each result?

Questions 1–5 are independent factual inputs to question 6.
All upstream sources below were retrieved on 2026-09-11.

## Facts: tool-result text

| Harness | Source boundary | Consequence |
| --- | --- | --- |
| Codex | The assessed transcript contains matched result payloads. Hook payloads carry tool-specific JSON. | Transcript text is observable for the pinned schema; hook serialization is a separate measure. |
| Claude Code | `PostToolBatch.tool_response` contains the corresponding model-visible content. | A reader can count text fields without treating structured tool objects as equivalent text. |
| OpenCode | The V2 plugin API exposes successful results and failures. | Exact model-visible serialization remains unverified. |

Sources: the [preliminary assessment](preliminary-assessment.md),
[Codex hook fields](https://learn.chatgpt.com/docs/hooks#posttooluse),
[Claude batch fields](https://code.claude.com/docs/en/hooks#posttoolbatch), and
[OpenCode tool hooks](https://opencode.ai/v2/docs/build/plugins#tools).

## Facts: outer and nested tool calls

The assessed Codex transcript contains 323 outer calls with matching results.
It does not provide a separately verified count of nested tool calls.
Its raw slice and digest remain in the [preliminary assessment](preliminary-assessment.md).

Codex documents hooks for nested code-mode calls.
Hosted tools and specialized paths can bypass these hooks.
A completion poll can deliver the original execution event.
Hook counts therefore describe their covered boundary, not all outer calls.
Source: [Codex tool coverage](https://learn.chatgpt.com/docs/hooks#tool-coverage).

Claude supplies tool-use identifiers with batch results.
This establishes per-tool correlation, not a nested-call hierarchy inside an arbitrary wrapper.
Source: [Claude batch fields](https://code.claude.com/docs/en/hooks#posttoolbatch).

OpenCode's stats implementation counts tool parts by tool name.
That aggregation does not establish outer-versus-nested attribution.
Source: [OpenCode stats](https://github.com/anomalyco/opencode/blob/dev/packages/opencode/src/cli/cmd/stats.ts), tool-part loop and session aggregation.
The source is the mutable `dev` branch, not a pinned installed release.

## Facts: input and cached-input tokens

| Harness | Available evidence | Limitation |
| --- | --- | --- |
| Codex | Local usage records expose input and cached-input counters. App Server documents usage updates. | The opened App Server page does not establish the local transcript's leaf schema. |
| Claude Code | OpenTelemetry documents input, output, cache-read, and cache-creation token counters. | Export configuration and local delivery remain unverified. |
| OpenCode | Stats source aggregates input and cache-read/write counters. | The source defaults missing values to zero; that display cannot prove availability. |

Sources: the [preliminary assessment](preliminary-assessment.md),
[Codex App Server events](https://learn.chatgpt.com/docs/app-server#events),
[Claude token counter](https://code.claude.com/docs/en/monitoring-usage#token-counter), and
[OpenCode stats](https://github.com/anomalyco/opencode/blob/dev/packages/opencode/src/cli/cmd/stats.ts), session token fallback and aggregation.

## Facts: completed compactions

Codex exposes completed compaction items through App Server and a `PostCompact` hook.
The local transcript contains one compaction record.
Sources: [Codex App Server items](https://learn.chatgpt.com/docs/app-server#items) and the [preliminary assessment](preliminary-assessment.md).

Claude documents `PostCompact` with manual or automatic triggers.
Source: [Claude compaction hook](https://code.claude.com/docs/en/hooks#postcompact).

The inspected OpenCode plugin and stats sources do not establish a completed-compaction measure.
This measure remains unknown in this assessment.

## Facts: reasoning output

Codex's local usage record includes `reasoning_output_tokens`.
This counter does not attribute reasoning to individual tool results.
Source: the [preliminary assessment](preliminary-assessment.md) and its pinned transcript, line 2777.

OpenCode stats source carries a reasoning-token aggregate.
Its per-model output total adds reasoning tokens to output tokens.
Those totals require separate labels before comparison.
Source: [OpenCode stats](https://github.com/anomalyco/opencode/blob/dev/packages/opencode/src/cli/cmd/stats.ts), per-model and total-token aggregation.

The opened Claude monitoring contract does not establish a distinct reasoning-token counter.
Visible thinking text and provider reasoning-token counts remain different measures.
Source boundary: [Claude monitoring](https://code.claude.com/docs/en/monitoring-usage).

## Facts: measurement ownership and interception

Bench's existing census records selected raw shell calls that name its worktree pool.
It stores a timestamp and command head, not tool output.
Source: `internal/census/census.go:31-66`.
The capability registry keeps all four harness measures unknown pending its named reader.
Source: `internal/harnesses/harnesses.go:116-175`.
FT231 retains the Bench-versus-harness ownership split.
Source: `roadmap/FT231.md:61-69`.

| Harness | Documented result control | Consequence |
| --- | --- | --- |
| Codex | Post-tool feedback can replace a result. A blocking hook rejects a nested promise; `continue: false` has different semantics. | A generic wrapper must not treat these controls as interchangeable. |
| Claude Code | `updatedToolOutput` replaces output before model delivery and must retain the tool's output shape. | Invalid built-in replacement shapes can expose the original output. |
| OpenCode | V2 plugins can mutate successful result objects. | Model-input timing and complete-output access remain unverified. |

Sources: [Codex post-tool hooks](https://learn.chatgpt.com/docs/hooks#posttooluse),
[Claude replacement control](https://code.claude.com/docs/en/hooks#posttooluse-decision-control), and
[OpenCode tool hooks](https://opencode.ai/v2/docs/build/plugins#tools).

Codex explicitly calls its transcript format unstable.
It also documents that `updatedMCPToolOutput` and `suppressOutput` are unsupported hook fields.
Source: [Codex hooks](https://learn.chatgpt.com/docs/hooks).

## Tested results

The coordinator re-read the assessment's original transcript slice.
Its SHA-256 matches the recorded digest.
Schema inspection found 255 custom calls, 68 function calls, matching outputs, 334 usage events, and one compaction record.
The usage objects expose total and last-record counters, including input, cached-input, and reasoning-output keys.
No transcript content was copied into this report.

The coordinator re-opened the load-bearing upstream passages and the Bench owner definitions.
No live hook, plugin, telemetry exporter, or cross-harness byte-compatibility probe ran.

## Inferences and unresolved joins

A documented field does not prove that an installed integration captures it.
A hook's structured payload need not equal transcript text or complete producer stdout.
A result can already contain truncation before a post-tool hook receives it.
No inspected contract proves recovery of those omitted producer bytes.

Agent ancestry and nested tool-call ancestry are separate relations.
The current evidence does not establish a portable join for both.
No source attributes provider token costs to each tool result.

## Proposals

Record the harness version, integration, source boundary, and availability for each raw measure.
Keep text bytes, structured payload bytes, and producer-output bytes separate.
Use observed field presence before numeric aggregation.
Resolve enforcement and budget posture through tickets 6 and 9.
These proposals do not authorize a new policy.

## Verification record

- [x] The report states the recommendation, scope, and evidence status.
- [x] Facts, tested results, inferences, and proposals have separate sections.
- [x] Each question has a synthesized section and source boundary.
- [x] Comparison tables state the consequence of each capability.
- [x] Unknowns and incompatible meanings remain explicit.
- [x] The question graph is small enough to read without a diagram.
- [x] Local installation and byte compatibility remain unverified.
- [x] Load-bearing claims cite primary URLs or exact local locations.
- [x] The report ends with a validation plan.

## Validation plan

Pin each harness version before a source reader claims support.
Capture one Unicode text result, one structured result, one failure, and one absent measure.
Compare text bytes at the producer, hook, and model-visible boundaries.
Exercise nested calls and a completion poll to check count attribution.
Check that an absent source remains unknown through aggregation.

For an overflow implementation, verify complete artifact retention and a working detail route before enabling replacement.
Verify execution status, tool-specific shapes, and nested promise behavior at that same boundary.
