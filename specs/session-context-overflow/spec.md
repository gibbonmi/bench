# Preserve complete output before verified overflow projection

Status: staged

Decision source: `specs/session-context-efficiency/decisions/session-context-efficiency.md` (ready compiled map).

Verification log: 2 iteration(s) to accept — Sol/high accepted the suite. Trace-only partials and review-state bookkeeping are folded.

Coordinator: [Session context efficiency](../session-context-efficiency/spec.md)

## Problem

A wrapper can expose an oversized result after an upstream layer already discarded producer bytes.
A hook with the wrong replacement semantics can also hide failure status or interrupt nested tool execution.

## Solution

Verify each installed harness path before it can replace a result.
Preserve the complete result first, then return approved diagnostics, status, true size, and an exact detail route.
Unsupported paths retain the original result and report the limitation.
No efficiency denial or automatic retry joins the operation.

## User stories

Line: gpt-5.6-terra / high.
The work crosses uncertain runtime or lifecycle boundaries.

1. As a reviewer, I want installed-runtime proof, so that documented hooks do not become assumed capabilities.
2. As an agent, I want complete output preserved before replacement, so that projection cannot destroy evidence.
3. As an agent, I want the producer status retained, so that a failed command cannot look successful.
4. As an agent, I want bounded diagnostics, so that the result still supports the next decision.
5. As an agent, I want the true byte count, so that the summary describes the full artifact.
6. As an agent, I want a working complete-detail route, so that I can retrieve omitted output.
7. As an agent, I want unsupported shapes preserved, so that an unrecognized callback cannot corrupt evidence.
8. As an agent, I want failed preservation reported without denial, so that efficiency work cannot block a completed operation.
9. As a reviewer, I want completed operations invoked once, so that overflow recovery cannot repeat a mutation.
10. As an agent, I want upstream truncation identified, so that a partial source cannot be called complete.
11. As a reviewer, I want runtime drift to disable unverified replacement, so that old proof cannot authorize a changed path.
12. As a reviewer, I want measured budget approval, so that the adapter cannot choose its own cap.
13. As a reviewer, I want existing guards retained, so that projection cannot change safety authority.
14. As an agent, I want duplicate callbacks handled consistently, so that artifacts cannot race or multiply the producer.
15. As an agent, I want private artifact ownership, so that hostile identifiers cannot redirect writes.
16. As an agent, I want empty results handled correctly, so that known absence is not an unsupported source.
17. As an agent, I want small results preserved, so that projection only changes results that need it.
18. As an agent, I want the verified installed path exercised, so that package tests cannot conceal missing wiring.

## Implementation decisions

The first ticket records runtime capability evidence before adapter implementation.
Its inventory comes from the compiled harness record and the source matrix's named tool paths.
Each path records harness version, event shape, producer boundary, callback timing, replacement semantics, and retrieval result.
Codex local tools, nested tool calls, completion polls, and hosted-tool exceptions receive separate dispositions.
Claude structured callbacks and model-visible batch serialization receive separate dispositions.
OpenCode remains unverified until its installed runtime proves the required boundary.

The runtime probe uses a deterministic producer with known complete bytes and a separately observed invocation count.
It tests success, failure, empty output, large output, multibyte diagnostics, and nested completion where supported.
The probe compares the retrieved artifact with the independent producer record.
A simulated callback cannot replace this installed-runtime evidence.
A path without complete producer bytes remains unsupported even if its callback accepts replacement.

The second ticket is blocked by runtime evidence and the measurement child's approved budget decision.
Before it starts, spec authoring records the exact eligible adapters, callback mappings, byte values, and detail-route grammar.
The author also closes adoption and configuration fences for those verified adapters.
If several adapters qualify, the author splits the final ticket into one vertical ticket per adapter.

The current candidate fence is not authority to guess an unverified adapter.
A missing proof or policy stops the ticket before product writes.
The child remains staged when no path qualifies.

A single typed overflow owner receives a completed producer result through a verified adapter.
It does not execute the producer.
The existing executable dispatch and hook launcher remain the launch chain.
This work adds no executable authenticity claim and does not replace their trust roots.
The adapter uses independently observed runtime evidence as a compatibility condition, not a self-authenticating binary claim.

Artifacts use private directories and exclusive creation under a Bench-owned local artifact root.
A generated invocation identity selects the artifact name; raw callback identifiers never form a path.
The owner preserves bytes before publishing a detail route.
A completed artifact remains available after the callback and across session compaction.
Its route names the file and its full byte count without embedding complete output into the summary.

At a write, close, or readback failure, the adapter preserves the original result.
It reports that complete-output preservation is unavailable.
It never repeats the operation to reconstruct missing bytes.
If the harness cannot express this posture, that path remains unsupported.
The existing safety guards retain their behavior and execution order.

## Testing decisions

Tests drive the production owner through controlled inputs and its existing injected boundaries.
Ordinary package tests execute through the existing Go test phase.
Hook integration tests execute through `bench test --check system` where the row requires a real subprocess.
Live harness evidence remains review-owned and cannot be replaced by a simulated callback.
Planned test names below identify required tests, not tests that already exist.

### Seam diagram

```text
caller -> domain entrypoint -> verification -> existing operation -> complete result
                ^                  ^                  ^
          command tests       stale/fault cases   durable evidence
```

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
| --- | --- | --- | --- | --- |
| OV1 | 1 | Each eligible harness path has installed-runtime evidence of model-visible replacement | review-owned: runtime-evidence.md with pinned runtime and tool path | A documentation citation alone cannot establish replacement semantics |
| OV2 | 2 | A replacement points to a complete artifact that exists before the replacement returns | planned TestOverflowPreservesBeforeReturn in internal/harnessoverflow | An omitted or shortened artifact fails the independent byte comparison |
| OV3 | 3 | A nonzero producer status remains nonzero after projection | planned TestOverflowFailureStatus in internal/harnessoverflow | Returning the hook success as tool success hides the producer failure |
| OV4 | 4 | Projected results retain the approved bounded diagnostic content | planned TestOverflowDiagnosticProjection in internal/harnessoverflow | Only printing an artifact path drops the failure evidence |
| OV5 | 5 | The result reports the full artifact byte count | planned TestOverflowTrueBytes in internal/harnessoverflow | Counting the projected text understates the preserved output |
| OV6 | 6 | The detail route retrieves the exact preserved bytes | planned TestOverflowDetailRoute in internal/systemtest | A stale or inaccessible path cannot satisfy the retrieval comparison |
| OV7 | 7 | An unsupported callback retains its original tool result | planned TestOverflowUnsupportedShape in internal/harnessoverflow | Guessing a callback mapping changes unsupported output |
| OV8 | 8 | A failed artifact write retains the original tool result | planned TestOverflowSpillFailure in internal/harnessoverflow | Replacing after a failed spill loses the only complete result |
| OV9 | 8 | A spill failure reports its capability limitation without denying the operation | planned TestOverflowLimitation in internal/harnessoverflow | A new efficiency refusal violates the advisory posture |
| OV10 | 9 | The adapter never reruns the completed producer operation | planned TestOverflowNoReplay in internal/systemtest | A producer sentinel count above one proves a repeated mutation |
| OV11 | 10 | Unverified upstream completeness prevents automatic replacement | planned TestOverflowIncompleteProducer in internal/harnessoverflow | Saving an already-truncated result cannot certify complete producer preservation |
| OV12 | 11 | Unknown or changed harness versions remain in the unsupported posture | planned TestOverflowVersionBoundary in internal/harnessoverflow | A prior runtime proof cannot silently authorize a changed callback contract |
| OV13 | 12 | Automatic replacement requires an approved per-surface byte budget | review-owned: ticket-entry evidence and policy record | A build-selected cap cannot satisfy reviewer approval |
| OV14 | 13 | The adapter preserves existing guard outcomes | planned TestOverflowPreservesGuards in internal/systemtest | A new callback cannot suppress a prior safety refusal |
| OV15 | 14 | Repeated completion delivery resolves to the same complete artifact | planned TestOverflowDuplicateCompletion in internal/harnessoverflow | Concurrent callback delivery cannot overwrite another invocation artifact |
| OV16 | 15 | The artifact path rejects traversal and symlink redirection | planned TestOverflowArtifactPath in internal/harnessoverflow | A hostile result identifier cannot redirect the write outside its private directory |
| OV17 | 16 | An empty supported result retains its status and empty body | planned TestOverflowEmptyResult in internal/harnessoverflow | Treating empty as missing creates a false limitation |
| OV18 | 17 | Supported results at or below the approved budget retain their original body | planned TestOverflowWithinBudget in internal/harnessoverflow | Unnecessary replacement changes a result that needs no projection |
| OV19 | 18 | Installed configuration exposes only the verified harness paths | planned TestOverflowInstalledWiring in internal/systemtest | A supported core without its actual adapter cannot satisfy the end-to-end result |
| OV20 | 2 | Each eligible runtime path proves complete producer bytes before its model-visible replacement | review-owned: runtime-evidence.md compares the independent producer record with the stored artifact | A replacement-only demonstration cannot certify complete-output preservation |
| OV21 | 3 | Each eligible runtime path preserves a nonzero producer status at the model boundary | review-owned: runtime-evidence.md records the failing producer and observed result | A successful callback cannot stand in for the failed producer status |
| OV22 | 6 | Each eligible runtime path retrieves the exact complete artifact through its reported detail route | review-owned: runtime-evidence.md records actual retrieval and byte comparison | A plausible path without a successful retrieval cannot authorize enablement |
| OV23 | 8 | An artifact-creation failure retains the original tool result | planned TestOverflowCreateFailure in internal/harnessoverflow | Replacement after failed creation loses the only complete result |
| OV24 | 8 | An artifact-close failure retains the original tool result | planned TestOverflowCloseFailure in internal/harnessoverflow | A successful write cannot certify a failed close |
| OV25 | 8 | An artifact-readback failure retains the original tool result | planned TestOverflowReadbackFailure in internal/harnessoverflow | A published path cannot certify unreadable preserved bytes |

### Edge inventory

OV2 checks preservation before replacement.
OV8 and OV23–OV25 separately cover write, creation, close, and readback failures.
OV7 and OV11 cover malformed callbacks, missing producer bytes, and upstream truncation.
OV12 covers absent and changed runtime identities.

OV14 and OV19 cover the actual installed path rather than a test-only package variable.

OV15 covers repeated delivery and concurrent distinct invocations.
OV16 covers absent artifact directories, existing empty directories, path traversal, symlinks, and name collisions.
OV17 and OV18 cover empty, below-budget, exact-budget, and oversized supported results.
These behaviors serve the kit and every linked repository with a verified installed path.

Won't handle: unsupported replacement mechanisms — the adapter keeps the original result and reports its limitation.
Won't handle: reconstruction of upstream-discarded bytes — the original tool path remains the caller.
Won't handle: automatic producer retry — the agent retains the existing explicit operation boundary.

## Ownership fences

- `.bench/hooks/result-overflow.sh`
- `.claude/settings.json`
- `.codex/hooks.json`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `cmd/bench/main.go`
- `cmd/bench/otel_hook_seams_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/harness_record_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `internal/harnesses/harnesses.go`
- `internal/harnesses/harnesses_test.go`
- `internal/harnessoverflow`
- `internal/otelrecord/registry.go`
- `internal/systemtest/harness_overflow_test.go`
- `reviews/session-context-overflow.md`
- `specs/session-context-overflow/assets/runtime-evidence.md`
- `tests/canary/line-routing/agent-hook-unwired`
- `tests/canary/line-routing/stop-hook-unwired`
- `tests/canary/load-validity-metadata/codex-hooks-broken`
- `tests/canary/load-validity-metadata/codex-hooks-timeout`
- `tests/canary/load-validity-metadata/codex-hooks-timeout-typed`
- `tests/canary/package-core-guard/unrouted-subcommand`

Reviewer disposition: Sol/high review accepted; user spec and ticket sign-off remains pending.
The fence is the union of ticket writes and the review pickup.
A build cannot change this spec, its acceptance rows, or its tickets.

## Ticket graph

| Ticket | Blocked by | Delivered coverage |
| --- | --- | --- |
| [1. Prove installed overflow paths](tickets/1-prove-runtime-paths.md) | none | OV1, OV20, OV21, OV22 |
| [2. Install verified overflow projection](tickets/2-install-verified-overflow.md) | 1-prove-runtime-paths.md | OV2, OV3, OV4, OV5, OV6, OV7, OV8, OV9, OV10, OV11, OV12, OV13, OV14, OV15, OV16, OV17, OV18, OV19, OV23, OV24, OV25 |

## Out of scope

Provider-token attribution is a separate measurement capability: approximately 4 edits, 1 gate run.
A general artifact browser is a separate capability: approximately 6 edits, 1 gate run.
Automatic artifact retention policy is separate: approximately 4 edits, 1 gate run.

## Further notes

### Source trace

| Source clause | Coverage |
| --- | --- |
| Ticket 3: complete artifact, bounded diagnostics, status, true size, detail route | OV2–OV6 |
| Ticket 5: unresolved runtime and upstream boundaries | OV1, OV7, OV11, OV12, OV19–OV22 |
| Ticket 6: preserve the original, report limitations, add no denial | OV7–OV10, OV14, OV23–OV25 |
| Ticket 9: measured budget approval | OV13 |

### Reader sweep and proof checklist

[Seam evidence](../session-context-efficiency/assets/seam-evidence.md) records the inspected hook and harness owners.
Cited symbols: the command registry's hook flag feeds `beginHookSpan` and the OTEL hook inventory.
Import edges: the internal hook command imports the new overflow owner after the checkpoint fixes its grammar.
Source-row clauses and occurrences: the sole compiled map owns the clauses above.

Promised field labels: no harness wire field is promised before runtime verification.
Changed-function callers: the actual hook configuration and adoption readers require closure at the checkpoint.
Copy survival: OV2 compares the entire artifact against an independently known producer record.
The deterministic omission case removes the final artifact chunk while retaining the summary.
OV2 must fail that case.

### Flagged additions

Private artifact ownership and duplicate-delivery behavior are necessary preservation details under OV15 and OV16.
The staged adapter ticket is not executable until its runtime, budget, and adoption-fence checkpoint completes.
This limitation retains capability-dependent scope instead of claiming support from documentation.
