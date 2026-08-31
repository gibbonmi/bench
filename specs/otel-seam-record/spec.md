# Bench OpenTelemetry seam record

Status: staged

Roadmap: FT274

Decision source: roadmap/FT274.md (named reviewed artifact, reviewer decisions dated 2026-08-29; the spans-only scope and the full seam set were confirmed on 2026-08-30)

Verification log: 2 iterations to accept — nine blocking findings (a home-read import cycle, subject text on a span, an unpriced rotation promise) folded after round one; the acceptance pass returned two nits, folded into acceptance

## Problem

FT231 needs measured evidence, and no shared record supplies it. The gate writes its own run log, the lane record drops the failing check, and the census counts raw calls, but no record joins them. Every scorecard cell is hand-authored. An interrupted run keeps no evidence of the phase it started.

## Solution

Bench starts an OpenTelemetry span at each declared seam entry. It appends each span as one OTLP-JSON line to a local file under the Bench home. The processor writes a line at span start and a second line at span end, so a killed run keeps its started seam. The encoder is hand-written, because every official OTLP-JSON encoder requires a protobuf module and the reviewer excluded that footprint. FT231, FT232, FT204, FT241, and the scorecards read this one record. This spec ships spans only.

## User stories

Group 1 — the record package.
Line: opus / medium. The encoder and the start-writing processor are novel work against an exact upstream format.

1. As a record consumer, I want each span as one OTLP-JSON line under the key `resourceSpans`, so that a collector receiver ingests the record.
2. As a record consumer, I want hex ids, integer enums, lowerCamelCase names, and quoted 64-bit integers, so that the OTLP-JSON deviations hold.
3. As a retro reader, I want a record line at span start, so that an interrupted run keeps its started seam.
4. As a record consumer, I want the start line marked with a start attribute, so that a consumer can filter unfinished spans.
5. As a record consumer, I want the end line to carry the complete span, so that elapsed time derives from the record.
6. As an operator, I want the record at `otel/<repository key>/traces.jsonl` under the Bench home, so that records key by repository as census records do.
7. As a concurrent verb, I want each record appended in one write, so that parallel writers keep every line intact.
8. As an operator, I want a symlinked record directory refused, so that a redirected record cannot write outside the Bench home.
9. As a verb caller, I want a failed record write to leave the outcome unchanged, so that the record stays evidence, never a condition.
10. As a reviewer, I want only official-org modules with no protobuf and no gRPC, so that the footprint obeys the dependency standard.

Group 2 — the instrumented seams.
Line: opus / medium. The instrumentation crosses packages at reviewed seams under a covering gate.

11. As a retro reader, I want a root gate span with the subject id, mode, and exit, so that the run lands in the record.
12. As a retro reader, I want a phase span per executed phase with name and exit, so that per-phase timing derives from the record.
13. As a retro reader, I want a skipped phase's span to name its blocker, so that a cascade skip stays attributed.
14. As the FT232 tripwire, I want the lane span to carry the failing check and its diagnostic, so that the advisory reads the lane's red.
15. As a retro reader, I want a commit span with the subject digest and the outcome, so that commit cost lands in the record.
16. As a retro reader, I want `bench worktree land` to write a landing span covering composition and publication, so that landing cost lands in the record.
17. As a retro reader, I want a span per worktree verb that names the assignment, so that worktree traffic lands in the record.
    The instrumented set is `create`, `exec`, `merge`, `release`, `land`, `build`, and `reauthorize`. Each acts on one named assignment.
    `clean` and `reclaim` stay out, because a bulk verb acts on many assignments and names no single subject. The read-only verbs stay out too.
18. As a retro reader, I want a span per hook plumbing verb with its exit, so that the harness events join the record.
19. As a reviewer, I want span attributes limited to seam name, subject id, outcome, and measures, so that no payload enters the record.

Group 3 — the Bench-owned measures.
Line: opus / medium. The measures ride on the group 2 spans.

20. As a scorecard author, I want the landing span to carry the census raw-call count, so that the count survives the release.
21. As a scorecard author, I want the composed path count on the commit and landing spans, so that write-set size lands in the record.
22. As a record consumer, I want the gate iteration count derived by counting a subject's gate spans, so that no second counter duplicates the record.

Group 4 — the capability record.
Line: opus / low. The cells land at a known seam under an existing independent expectation.

23. As a measure consumer, I want harness-record cells for tokens, tool calls, Read paths, and turns, so that the record declares each measure's supplier.
24. As an operator, I want `bench harnesses <harness>` to print the measure cells, so that the declaration is readable without the source.

Group 5 — durability, conformance, and docs.
Line: opus / medium. A kill-based system test and a `go/ast` check need care.

25. As a retro reader, I want a gate killed mid-phase to keep the started phase line, so that recovery evidence survives the crash.
26. As a reviewer, I want a conformance check that reds a registered seam with no span, so that no seam stays silently uninstrumented.
27. As an operator, I want `DATA_HANDLING.md` to name the local record file, so that the no-upload claim stays true and complete.
28. As the reviewer, I want the build to stop and report if story 25's test cannot pass, so that FT71 stays a separate ledger.

## Implementation decisions

- One new package, `internal/otelrecord`, is the only package that knows the record's address, file layout, encoder, and processor. The intent ledger sets this single-owner precedent.
- A new stdlib-only package, `internal/benchhome`, owns the one `BENCH_HOME` read, and `internal/worktree` delegates its `Home()` to it. `internal/worktree` imports `internal/gate`, so the gate cannot import the read where it lives today. A prefactor ticket moves the read before any span ships.
- The dependencies are `go.opentelemetry.io/otel`, `go.opentelemetry.io/otel/trace`, and `go.opentelemetry.io/otel/sdk`, pinned at v1.46.0 under Apache-2.0. The transitive set is logr, stdr, xxhash, uuid, x/sys, and the passive `auto/sdk` shim. `sdk/metric` stays out with the metric signal.
- The encoder is hand-written and mirrors the struct schema of upstream's own internal `otlpjson` package. Every official OTLP-JSON encoder requires `google.golang.org/protobuf`, so the hand-written encoder is the only route through the reviewer's no-protobuf decision.
- One span processor writes at `OnStart` and at `OnEnd`. Each record is one synchronous `O_APPEND` write, census-style, with no buffer and no background goroutine. A kill therefore loses at most the line in flight.
- The start line carries the attribute `bench.record=start` and no end time. The end line carries the complete span. No consumer merges the pair; a reader filters by the marker. This disposition is open to reviewer veto.
- A span carries the subject digest — the composed tree or commit id — and never subject text. A commit subject carries objective text by design, and the record must not copy it into a third durable place.
- The write-set size is the composed path count at each publication seam. The commit span counts the named attributed paths, and the landing span counts the reviewed name-only diff. Each seam derives its own count where it already holds the list.
- A verb boundary resolves the Bench home once through the `internal/benchhome` read and constructs the provider there. The gate threads its run log through `context.Context` today, and the tracer follows the same pattern.
- `internal/otelrecord` exports the seam registry: each entry names the seam and the Go symbol that starts its span. The `go/ast` conformance check enumerates this registry, so the registry is the one source for the instrumented set.
- The hook plumbing verbs are instrumented at their one dispatch seam in `cmd/bench`, not in each adapter package.
- The record file is append-only. Rotation and retention stay with FT71 and are priced under Out of scope.
- The record directory is `<home>/otel/<poolkey.Key(root)>/`, a sibling of the census directory.

## Testing decisions

- A good test drives a real verb with `BENCH_HOME` set to a temporary directory and reads the JSON lines back. The line content is the external behavior.
- The encoder and writer tests live in `internal/otelrecord` with fixture spans. The lane span test joins the existing lane tests in `internal/gate`. The verb-level spans and the crash test live in `internal/systemtest` against the built binary.
- The system phase runs only when the graded root is the kit checkout. Each system-seam ticket therefore names the focused hand-run, with the system build tag and the run-binary environment, as its ticket-time observation.
- The capability record cells get an independent want-list expectation in `internal/harnesses`, beside the mechanics expectation that already reds a dropped name. The seam-instrumentation check is a new `go/ast` test in `internal/conformance`, with `internal/consumers` as the loader prior art.
- The gate observes the feature through the `test` and `system` phases. No new gate phase is added.

### Seam diagram

    trigger: a Bench verb (gate, commit, land, worktree, hook plumbing)
        │
        ▼
    seam entry  ──▶  [ internal/otelrecord provider + processor ]  ──▶  <home>/otel/<key>/traces.jsonl
                      ◀ tests attach here: run the verb with a temp BENCH_HOME, then parse the lines

### Acceptance coverage map
| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| OT1 | 1 | A record line parses as one JSON object with the `resourceSpans` key | new encoder test in `internal/otelrecord` | a Go-struct encoding has no `resourceSpans` key and the parse assertion reds |
| OT2 | 2 | The encoder writes trace and span ids as lowercase hex strings | encoder fixture test in `internal/otelrecord` | stock proto3 JSON emits base64 ids and the hex assertion reds |
| OT3 | 2 | The encoder writes span kind and status code as integers | encoder fixture test in `internal/otelrecord` | an enum name string reds the integer assertion |
| OT4 | 2 | The encoder writes `startTimeUnixNano` as a quoted decimal string | encoder fixture test in `internal/otelrecord` | an unquoted number reds the string-type assertion |
| OT5 | 3 | A record line exists before the span ends | processor test in `internal/otelrecord` | an end-only processor writes nothing before the end and the read reds |
| OT6 | 4 | The start line carries the `bench.record=start` attribute | processor test in `internal/otelrecord` | an unmarked start line reads as a finished span and the attribute read reds |
| OT7 | 5 | The end line carries `endTimeUnixNano` | processor test in `internal/otelrecord` | a start-only writer never writes the end time and the read reds |
| OT8 | 6 | The record file sits at `otel/<repository key>/traces.jsonl` under the resolved home | package test with `BENCH_HOME` set to a temp directory | a write to a hard-coded home misses the temp home and the exact-path read reds |
| OT9 | 7 | Two writers with separate file handles leave only intact JSON lines | writer test with two independent openers | a multi-call writer interleaves across handles and a line fails to parse |
| OT10 | 8 | A symlinked record directory refuses the write | writer test with a symlinked directory | a follow-through write lands outside the home and the refusal assertion reds |
| OT11 | 9 | The writer returns an error on an unwritable record directory | package test with an unwritable directory | a swallowed error hides the failed write and the error assertion reds |
| OT13 | 10 | go.mod names no protobuf module and no gRPC module | new dependency-footprint test | an OTLP exporter import drags `google.golang.org/protobuf` in and the check reds |
| OT14 | 11 | `bench gate` writes a root span carrying the subject id and the run's exit | system test against the built binary | a span without the subject id cannot group iterations and the attribute read reds |
| OT15 | 12 | Each executed phase writes a phase span with its name and exit | system test against the built binary | a root-only tracer leaves no phase line and the per-phase read reds |
| OT16 | 13 | A skipped phase's span names its blocker | system test with a red need | a skip recorded without its blocker reds the attribute read |
| OT17 | 14 | The lane span carries the first failing check with its first diagnostic line | lane test in `internal/gate` | a phase-and-exit-only record drops the check and the attribute read reds |
| OT18 | 15, 21 | `bench commit` writes a commit span carrying the subject digest and the composed path count | system test against the built binary | a commit span without the count reds the attribute read |
| OT19 | 16, 20 | `bench worktree land` writes a landing span carrying the census raw-call count | system test against the built binary | a landing span without the count reds the attribute read |
| OT20 | 17 | Each worktree verb writes a span naming the verb and the assignment | system test against the built binary | an uninstrumented verb leaves no line for its own invocation |
| OT21 | 18 | A hook plumbing verb writes a span naming the verb and its exit | system test piping a hook envelope | an uninstrumented dispatch leaves no hook line and the read reds |
| OT22 | 19 | Every span attribute stays inside the declared attribute set | review-owned | a payload attribute survives every mechanical check, so the review grades it |
| OT23 | 23 | The harness record carries the four measure cells | record test in `internal/harnesses` | an independent want-list expectation reds a record without the cells |
| OT24 | 24 | `bench harnesses claude` prints the four measure cells | command test in `internal/harnesses` | a projection that skips the cells reds the output read |
| OT25 | 25 | A gate killed mid-phase keeps the started phase span line | system test that kills the gate process mid-phase | an end-only or buffered writer loses the start line and the read reds |
| OT26 | 26 | The conformance check reds on a registered seam whose symbol starts no span | new `go/ast` check test in `internal/conformance` | an uninstrumented seam otherwise stays silently green |
| OT27 | 27 | `DATA_HANDLING.md` names the record file under the Bench home | review-owned | prose carries no mechanical seam, so the review grades the claim |
| OT28 | 9 | `bench gate` keeps its exit code with an unwritable record directory | system test with an unwritable record directory | a propagated record error flips the verb exit and the row reds |
| OT29 | 8 | A non-regular file at the record path refuses the write | writer test with a FIFO at the record path | an open of a FIFO blocks the first span and every recorded verb hangs |

Not covered: story 22 — the iteration count is a consumer derivation, and no reader ships in this spec.
Not covered: story 28 — the FT71 condition is a build-exit rule, and the workflow's material-shortfall rule owns it.

### Edge inventory

- An unset `BENCH_HOME` falls back to `~/.bench` through the one `benchhome` read; OT8 drives the env path.
- A kill between the span start and its append loses that one line; the crash test bounds the exposure to it.

- **Won't handle:** the metric signal — every Bench-owned measure rides on a span, and the cut is priced under Out of scope.
- **Won't handle:** rotation and retention of the record file — FT71 owns the retention contract, and the cut is priced under Out of scope.
- **Won't handle:** live collector ingest verification — the fixture pins the documented line shape, and no collector binary enters the gate.
- **Won't handle:** harness-owned measures (tokens, tool calls, Read paths, turns) — FT204 owns the transcript reader, and the cells stay Unknown.
- **Won't handle:** a merge of the start and end lines — the start marker lets any consumer filter, and no consumer ships here.
- **Won't handle:** span parenting across a process the gate did not start — an external propagation format enters no Bench seam. The subject id joins those spans for a reader. The gate hands its own phases child this run's traceparent, so the phase spans do join the run's root span.
- **Won't handle:** an operator disable switch — the record is a local file under the operator's own home, and a future spec can add one.
- **Won't handle:** control characters in a subject id — `encoding/json` escapes every control rune, and `internal/sanitize` owns the intake policy.
- **Won't handle:** a record line near the receiver's 1 MiB split bound — the no-payload rule keeps every line far under it.

## Ownership fences

- `specs/otel-seam-record/`
- `reviews/otel-seam-record.md`
- `internal/otelrecord/`
- `internal/benchhome/`
- `go.mod`
- `go.sum`
- `internal/gate/gate.go`
- `internal/gate/engine.go`
- `internal/gate/runner.go`
- `internal/gate/lane.go`
- `internal/gate/lane_test.go`
- `internal/gate/lane_record_test.go`
- `internal/gate/otel_env_test.go`
- `internal/commit/commit.go`
- `internal/landing/landing.go`
- `internal/landing/composition.go`
- `internal/worktree/worktree.go`
- `internal/worktree/build.go`
- `internal/worktree/build_test.go`
- `internal/worktree/reauthorize.go`
- `internal/worktree/exec_test.go`
- `internal/worktree/otel_seams_test.go`
- `internal/worktree/effects.go`
- `internal/worktree/exec.go`
- `internal/worktree/merge.go`
- `internal/worktree/land.go`
- `internal/worktree/main_test.go`
- `internal/harnesses/harnesses.go`
- `internal/harnesses/harnesses_test.go`
- `internal/harnesses/command.go`
- `internal/harnesses/command_test.go`
- `internal/conformance/`
- `cmd/bench/main.go`
- `cmd/bench/command_registry.go`
- `cmd/bench/otel_hook_seams_test.go`
- `DATA_HANDLING.md`
- `internal/systemtest/`

## Out of scope

- The metric signal export — about 10 edits, 3 gate runs. It adds the metricdata mapping, a periodic reader with a shutdown flush, and a second file.
- FT71 rotation, retention, and absorption — about 6 edits, 2 gate runs, decided after the crash-test verdict.
- The FT204 record query and transcript reader — a separate spec on its own roadmap row.
- The FT275 profile conformance declaration for linked projects — about 3 edits, 1 gate run, owned by FT275's own spec.
- An operator disable switch — about 2 edits, 1 gate run.

## Further notes

The upstream claims rest on sources outside the tree. The research pass of 2026-08-30 read the published `.mod` files at v1.46.0, the `otlpjsonfilereceiver` README, the OTLP file-exporter specification, and the opentelemetry-go exporter sources. The OpenTelemetry Go tree itself was not vendored or read in full. No published start-writing exporter exists, so the processor is novel work on a sanctioned SDK hook.

If the OT25 crash test cannot pass, the build stops and reports. FT71 then stays a separate ledger, and its recommended dependency on this row is withdrawn.
