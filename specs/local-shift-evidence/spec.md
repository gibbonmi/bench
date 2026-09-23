# Versioned local shift evidence

Status: staged

Roadmap: FT71

Decision source: the named reviewed artifact `roadmap/FT71.md`, the row body that the 2026-09-22 drain `d-465b789396f9` settled at `65ac3e2a` with `Next: spec`.

The row inherits two closed decisions. On 2026-08-29 the reviewer decided that FT71's events are span attributes in the FT274 seam record, not a second ledger. The FT274 crash test passed, so that condition holds. The 2026-08 capability audit closed L-04: no new canonical work record exists, and this evidence must not become one.

Verification log: 0 iteration(s) to accept — pending the independent review round at reviewer sign-off.

## Problem

Bench keeps no trustworthy local evidence of a shift. A shift writes an intent entry, iteration commits, and gate spans, but no record joins them. No record names the adapter, the line, the adapter result, the recovery reference, or the cleanup decision. A green shift deletes `.bench-notes.md` at teardown, so the memory of its iterations is gone.

The intent ledger also misleads. A shift that exits normally records its outcome, but its entry stays live while its pool worktree exists, and pool worktrees persist. A shift that crashes leaves an open entry that nothing ever resolves. No rule separates a crashed owner from a live one.

The seam record itself is not redacted. A probe on 2026-09-23 found `OTEL_RESOURCE_ATTRIBUTES` and `OTEL_SERVICE_NAME` values in every record line, because the OpenTelemetry SDK merges them into the default resource. `DATA_HANDLING.md` says that a line carries no environment value, so the claim is false today. The record also grows without a bound: the live record of this repository holds 196 MB after three weeks.

## Solution

The seam record becomes the versioned local evidence. Every line names its record schema and the Bench version that wrote it, and the encoder writes only declared keys. No environment value reaches a line. The record rotates into sealed segments of a bounded size and keeps a bounded count of them.

A shift writes one `shift` span for its whole life and one pass span for each iteration and refactor pass. Each gate run is the child of its pass, so the gate fingerprint and verdict join their iteration. The spans name the adapter, the declared tier, the iteration cap, and each adapter result and commit. They also name the shift outcome, a work state, the recovery reference, and the cleanup decision. The adapter's own call to `bench resolve-model` records the resolved model in the same trace.

Before any scratch cleanup, the shift retains its notes as a private memory file beside the record. The span references that file by digest, and no prompt reads it.

A shift that exits normally ends its intent. The shift records its lease identity, and a recovery pass resolves each crashed shift. A crashed shift is recovered only when its lease identity still holds and its owner is dead. Otherwise the pass abandons it and leaves its worktree untouched. The worktree shell records its session boundaries the same way. `DATA_HANDLING.md` documents the record as mutable local evidence, not a tamper-proof central audit system.

## User stories

Line: opus / medium.
Implementation-line reason: LE-C2 is the hardest chunk, because crash recovery acts on worktrees at session start and its tests kill real processes. LE-A carries rotation under concurrent writers. The rows are exact, the seams exist, and focused package tests red each wrong answer.
Harder chunks: LE-A, LE-C2.

One versioned, redacted record:

1. As an auditor, I want each record line to carry the schema version, so that a reader knows its schema.
2. As a release reviewer, I want each line to carry the Bench version, so that release evidence selects its own events.
3. As an operator, I want no environment value in a record line, so that an OTel variable cannot leak a secret.
4. As a reviewer, I want the encoder to drop each undeclared attribute, so that no payload attribute ships by accident.
5. As a record consumer, I want an unknown schema version read as a malformed line, so that I never misread it.
6. As a record consumer, I want a line with no schema attribute to read as before, so that old records stay readable.
7. As a maintainer, I want one owner for the trace handoff to a child process, so that the handoff cannot drift.

Rotation and retention:

8. As an operator, I want the record to rotate into sealed segments at a bounded size, so that it cannot fill the disk.
9. As an operator, I want a bounded count of sealed segments kept, so that the retention is explicit.
10. As a record consumer, I want the readers to read every segment in order, so that a rotation hides no span.
11. As a concurrent verb, I want appends across a rotation to leave only whole lines, so that no line is lost.
12. As a concurrent verb, I want an append to skip a rotation that another writer holds, so that no verb blocks.
13. As an operator, I want a reader to refuse a symlinked or FIFO segment, so that a planted path cannot misdirect a read.

Shift boundaries and the line:

14. As an auditor, I want one shift span from the intent to every exit, so that the record holds the shift boundaries.
15. As an auditor, I want the shift span to name the adapter by base name, so that no host path enters.
16. As an auditor, I want the shift span to carry the declared tier and iteration cap, so that the declared line shows.
17. As an auditor, I want the shift span to carry the outcome and a work state, so that the end states differ.
18. As an auditor, I want the shift span to carry the branch head commit, so that the record references the shift's work.
19. As an auditor, I want the recovery reference without an absolute path, so that preserved work is findable and host paths stay out.
20. As an auditor, I want the shift span to carry the cleanup decision, so that the record shows released or retained.
21. As an auditor, I want one pass span for each iteration and refactor pass, so that each adapter result is on record.
22. As an auditor, I want each gate run recorded under its pass span, so that the gate fingerprint and verdict join their iteration.
23. As an auditor, I want the pass span to carry its commit, so that the record references each commit.
24. As an auditor, I want an interrupted shift to end its open spans before exit, so that an interruption leaves complete boundaries.
25. As an auditor, I want a `line.resolve` span with the resolved model, so that the one resolver records the resolved line.

Shift memory:

26. As a retro reader, I want a shift to retain its notes before the scratch cleanup, so that its memory survives.
27. As a retro reader, I want a retaining shift to retain its notes too, so that a later worktree cleanup loses no memory.
28. As an auditor, I want the span to reference the memory by digest and size, so that the record holds no notes text.
29. As an operator, I want notes that are not a bounded regular file refused, so that an agent cannot leak or block.
30. As a retro reader, I want absent notes and empty notes recorded differently, so that I can tell the two apart.
31. As an operator, I want a bounded count of private memory files kept, so that the memory retention is explicit.
32. As a shift author, I want each new shift to start with empty notes, so that retained evidence never re-enters a prompt.

Intent completion and crash recovery:

33. As an operator, I want a normal shift exit to leave no live intent, so that status stops showing finished work.
34. As an operator, I want a retaining shift to keep its live intent and recovery pointer, so that I find the preserved work.
35. As an operator, I want each shift to record its lease identity, so that crash recovery can prove the owner.
36. As an operator, I want a crashed shift recovered when its lease still holds and its owner is dead, so that its intent completes.
37. As an operator, I want a recovered dirty worktree retained and locked in place, so that the crash loses no work.
38. As an operator, I want a recovered clean worktree released, so that the pool gets the worktree back.
39. As an operator, I want a crashed shift abandoned untouched when its lease no longer holds, so that another owner's checkout stays safe.
40. As an operator, I want an open entry left alone while its owner lives, so that recovery never judges a running shift.
41. As an operator, I want an entry with no lease abandoned only for a dead key owner, so that stale intent resolves safely.
42. As an operator, I want an unreadable or malformed lease to stop its verdict, so that recovery fails closed.
43. As an auditor, I want each recovery verdict recorded with the intent key, so that it joins the crashed shift's span.
44. As an operator, I want the recovery pass at session start and before each acquire, so that stale intent needs no manual step.
45. As an operator, I want the pass to print one line only when it acted, so that session start stays quiet.

Worktree shell session:

46. As an auditor, I want a `worktree.shell` span for each shell session, so that the record holds the session boundaries.
47. As an auditor, I want a normal shell exit recorded as completed with its release, so that subshell completion is on record.
48. As an auditor, I want a signalled session recorded as interrupted and retained, so that the interruption is on record.
49. As an auditor, I want a session whose shell cannot start recorded as failed, so that a failed session is on record.

Composition and documentation:

50. As a reviewer, I want one built-binary shift to join its shift, pass, line, and gate spans in one trace, so that composition holds.
51. As an auditor, I want a built-binary shift to leave no marker text in the record, so that the redaction holds end to end.
52. As an auditor, I want `DATA_HANDLING.md` to call local records mutable evidence inputs, so that no reader overtrusts them.
53. As an auditor, I want `DATA_HANDLING.md` to name the segments, memory files, bounds, and resource keys, so that its claims match.

Reviewed exclusions:

54. As a release operator, I want the `bank.ft71.local_event` producer left to its own capability, so that this spec stays local.
55. As an auditor, I want harness session boundaries left to the hook spans, so that no seam claims an unseen event.
56. As an auditor, I want tamper-proof central storage left out, so that the local record stays an honest evidence input.

## Implementation decisions

### The record contract

- The FT274 seam record is the one record. The record package owns the address, the segments, the encoder, the memory files, and the trace handoff.
- The encoder writes its own resource block for every line. The block holds exactly `service.name` with the value `bench`, `service.version` with the running Bench version, and `bench.record.schema` with the value `1`. The SDK resource never reaches a line, so no detector or environment value can enter.
- The command layer hands its stamped version to the record package once, at process start. A process that sets no version writes no `service.version` key.
- The declared attribute set becomes the redaction contract. The encoder drops each span attribute whose key is not declared. A later key ships only when its ticket declares it.
- A reader treats a line whose `bench.record.schema` is not `1` as a malformed line. A line with no schema attribute is a legacy line and reads as before.
- `service.version` is the relationship to release evidence. A release producer selects the lines whose version equals its envelope's package version. The producer itself is out of scope.
- The gate fingerprint is the existing gate span's `bench.subject.id`, and the gate verdict is its `bench.outcome`. This spec adds no oracle key.

### Rotation and retention

- The live segment stays `traces.jsonl`. When an append would take it past `bounds.RecordSegmentLimit`, the writer renames it to a sealed segment `traces-<UTC stamp>-<pid>.jsonl` and appends to a new live segment.
- Only the rotation step takes a lock: a non-blocking exclusive lock on one lock file in the record directory. A writer that cannot take the lock appends to the live segment and rotates on a later append. Each append stays one synchronous `O_APPEND` write.
- After a rotation, the writer removes the oldest sealed segments until `bounds.RecordSegmentsRetained` remain.
- The readers read the sealed segments in name order and then the live segment. Each segment passes the existing symlink and regular-file grade before the open.
- The values are reviewer-owned sizes, as the gate log retention is. This spec proposes a 16 MiB segment, 8 sealed segments, and 64 memory files.

### The trace handoff

- The two handoff variables, `BENCH_OTEL_ROOT` and `BENCH_OTEL_TRACEPARENT`, move from the gate to the record package with their inject and extract functions. The gate, the shift, and `bench resolve-model` call that one owner. The gate still strips inherited values before it composes its phase child.

### The shift record

- The `shift` span starts after the intent entry persists, and `finish` ends it on every exit path, the checkpoint exit included. The start line carries the intent key, the adapter base name, the declared tier, and the iteration cap.
- The declared tier is `BENCH_MODEL` only when it names a tier of the line binding. Any other value records no tier key.
- The end line carries the shift outcome, the work state, and the green or red outcome of the exit code. It also carries the branch head commit when one or more passes committed, the recovery reference, the cleanup decision, and the memory reference.
- No span carries a committed count or a pass index. A reader counts the pass spans that carry a commit and orders them by start time. The FT274 record already derives the gate iteration count the same way.
- The work-state vocabulary is `completed`, `failed`, `interrupted`, `recovered`, and `abandoned`. The shift maps `complete`, `no-op`, and `incomplete` to `completed`; `failed` and `usage` to `failed`; and `interrupted` to `interrupted`.
- The cleanup vocabulary is `released`, `retained`, and `none`. The recovery reference is `bench.recovery.kind` (`none` or `worktree`) and `bench.recovery.key`, the base name of the retained path. No absolute path enters the record.
- Each main iteration opens a `shift.iteration` span, and each refactor pass opens a `shift.refactor` span, under the `shift` span. Each pass gives its own context to the gate, so the gate span is the pass span's child.
- A pass span carries `bench.adapter.result` (`exited` or `spawn-failed`), `bench.adapter.exit` when the adapter exited, and the commit that the pass made. Every exit path ends an open pass span before it ends the `shift` span.
- The shift hands the adapter the trace handoff for the current pass span. `bench resolve-model` joins that trace when the handoff is present and records a `line.resolve` span. That span carries the harness, the declared tier, the resolved model when it is a safe model token, and the outcome.

### Shift memory

- The shift retains its notes on every path that reached the first iteration: before the teardown's scratch cleanup, and on the retain path. It reads `.bench-notes.md` through the bounded no-follow producer read, so a symlink, a FIFO, an oversized file, or invalid UTF-8 reads as refused.
- The record package writes the bytes to `memory/<UTC stamp>-<trace id>.md` below the record directory with mode 0600. It keeps the newest `bounds.RecordMemoryRetained` files.
- The `shift` span carries `bench.memory.state` (`retained`, `absent`, or `refused`), `bench.memory.bytes`, and `bench.memory.digest`. The record never holds the notes text.
- No prompt reads the memory directory. A new shift writes an empty notes file, as it does today.

### Intent completion and crash recovery

- The liveness rule drops an entry that holds an outcome and a recovery of empty or `none`. An entry with a `worktree:` recovery stays live while its worktree exists.
- The ledger entry gains an optional `lease` field. Right after the acquire, the shift records its lease line without the final newline.
- The recovery pass lives in the shift package, because it reuses the shift's scratch policy, memory retention, and preservation. It judges each shift entry that has no outcome.
- For an entry with a lease, the pass first grades the lease file without a follow. An absent file abandons the entry. A file that is not regular, is unreadable, or holds a malformed lease line skips the entry. This grade comes before the comparison, so a malformed line never counts as another owner (LE75).
- A well-formed lease line that differs from the recorded line abandons the entry. The recorded line with a live owner skips the entry. The recorded line with a dead owner recovers the entry.
- A recovery first retains the memory. A worktree with dirty paths beyond the scratch is retained and locked, and the entry gets its `worktree:` pointer. A clean worktree loses its scratch and is released. An abandon changes no worktree file, lease, or lock.
- For an entry with no lease, the intent package parses the owner process from the key beside `NewEntry`. A dead owner abandons the entry. A live owner or an unparsable key skips it.
- The pass sets the entry outcome to the work-state word, `recovered` or `abandoned`, and records one `shift.recovery` span with the intent key. The span starts a new trace, and the intent key joins it to the crashed `shift` span.
- Two triggers run the pass: a session-inspect phase after the resume phase, and `bench shift` before its acquire. The pass prints `bench shift recovery: recovered <n>, abandoned <m>` only when it acted.

### The worktree shell session

- `bench worktree shell` opens one `worktree.shell` span before it creates the assignment and ends it at every return. The end line carries the assignment id, the work state, the cleanup decision, and the outcome.
- A normal shell exit is `completed`. The cleanup is `released` when the release exits 0 and `retained` otherwise. A signal is `interrupted` with `retained`, because the lease stays for a reclaim. A shell that cannot start is `failed`.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| LE-A / `1-move-the-trace-handoff-into-the-record.md`, `2-redact-and-version-every-record-line.md`, `3-rotate-and-retain-the-record.md` | Every line is versioned and redacted, the record rotates and retains, and the trace handoff has one owner. | LE1, LE2, LE3, LE4, LE5, LE6, LE7, LE8, LE9, LE10, LE11, LE12, LE13, LE14, LE15, LE16, LE17, LE18, LE19, LE20, LE21 | `bench test --package ./internal/otelrecord`, `bench test --package ./internal/gate`, `bench test --package ./cmd/bench`, `bench test --check system` | yes |
| LE-B1 / `4-record-the-shift-boundaries.md` | A shift writes one `shift` span with its line, outcome, work state, commit, recovery reference, and cleanup decision on every exit path. | LE22, LE23, LE24, LE25, LE26, LE27, LE28, LE29, LE30, LE31, LE32, LE33, LE34, LE35, LE36, LE37 | `bench test --package ./internal/shift`, `bench test --check kit-compliance` | no |
| LE-B2 / `5-record-each-pass-under-the-shift.md`, `6-record-the-resolved-line.md`, `7-retain-the-shift-memory.md` | Each pass, its gate, and the resolved line join the shift trace, and the notes survive as a private memory file. | LE38, LE39, LE40, LE41, LE42, LE43, LE44, LE45, LE46, LE47, LE48, LE49, LE50, LE51, LE52, LE53, LE54, LE55, LE56, LE57, LE58, LE59, LE88 | `bench test --package ./internal/shift`, `bench test --package ./internal/otelrecord`, `bench test --package ./cmd/bench`, `bench test --check kit-compliance` | no |
| LE-C1 / `8-complete-the-shift-intent-on-exit.md` | A normal exit ends the shift intent, and each shift records its lease identity. | LE60, LE61, LE62, LE63 | `bench test --package ./internal/intent/...`, `bench test --package ./internal/shift` | no |
| LE-C2 / `9-recover-crashed-shifts-by-lease-identity.md` | A recovery pass resolves each crashed shift by its lease identity and records the verdict. | LE64, LE65, LE66, LE67, LE68, LE69, LE70, LE71, LE72, LE73, LE74, LE75, LE76, LE77, LE78, LE79, LE89, LE90, LE91 | `bench test --package ./internal/shift`, `bench test --package ./internal/intent/...`, `bench test --package ./internal/sessioninspect` | yes |
| LE-D / `10-record-the-worktree-shell-session.md`, `11-prove-the-shift-trace-through-the-built-binary.md`, `12-document-the-local-evidence-contract.md` | The shell session is on record, the built binary proves the whole trace and the redaction, and `DATA_HANDLING.md` states the contract. | LE80, LE81, LE82, LE83, LE84, LE85, LE86, LE87 | `bench test --package ./internal/worktree --run TestSubshell`, `bench test --package ./internal/worktree --run TestWorktreeSeamsMatchTheRegistry`, `bench test --check system`, `bench gate-prose . -- DATA_HANDLING.md` | no |

## Testing decisions

- A good test drives the real entry with a private Bench home and reads the record lines back. The line content is the external behavior. `TestLoopRetainsAndLocksDirtyWorktree` and the branch collision fixture are the shift precedents. `TestSubshellSignalsLeaveAReclaimableLease` is the helper-process precedent for a signal or a kill.
- The encoder, reader, rotation, handoff, and memory store get unit tests in the record package. The fixture spans of `encode_test.go` and the two-writer test of `writer_test.go` are the precedents.
- The shift tests drive `Loop` in process for normal exits. They drive `Loop` in a helper process for SIGINT and SIGKILL, because the checkpoint exits through `os.Exit`.
- The liveness rule gets a pure policy test beside `TestCompactionRuleReadsTypedLivenessFacts`. The resolve-model span gets a command test beside the hook seam test in `cmd/bench`.
- The system suite proves the stamped version and the whole trace through the built binary. `TestOtelCrashKeepsStartedPhaseLine` and the verb span tests are the precedents.
- The gate observes the feature through the `test` phase and the `system` phase. No new gate phase is added. The existing seam check in `kit-compliance` reds a registered seam whose symbol starts no span.

### Seam diagram

    trigger: `bench shift`, `bench worktree shell`, `bench resolve-model`, session start
        │
        ▼
    verb entry ──▶ [ shift record: shift, pass, and recovery spans; memory retention ] ──▶ span lines
                      ◀ tests attach here: run Loop or the recovery pass with a private
                        Bench home, then read the lines and the ledger back
        │
        ▼
    span ──▶ [ record package: redacting encoder, rotating writer, memory store ] ──▶ <home>/otel/<key>/
                      ◀ tests attach here: encode fixture spans, append across a small
                        segment limit, read the segments and memory files back
        │
        ▼
    ledger entry ──▶ [ liveness rule, lease identity ] ──▶ live intent or done
                      ◀ tests attach here: pure policy facts; a killed helper shift

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| LE1 | 1 | An encoded line carries the resource attribute `bench.record.schema` with the value `1` | planned encoder test in `internal/otelrecord` | An encoder that keeps the SDK resource has no schema key, so the value read reds. |
| LE2 | 2 | After the version is set to `9.9.9-test`, an encoded line carries `service.version` with that value | planned encoder test in `internal/otelrecord` | An encoder that never reads the set version writes no such key, so the value read reds. |
| LE3 | 2 | A `bench gate` run through the built binary writes lines whose `service.version` equals the version that `bench version` prints | planned system test in `internal/systemtest` | A binary that never hands its stamped version to the record writes no version, so the equality reds. |
| LE4 | 3 | With `OTEL_RESOURCE_ATTRIBUTES=bench.leak=PROBEMARK` and `OTEL_SERVICE_NAME=SVCMARK` set, no line that `Begin` writes holds either marker | planned provider test in `internal/otelrecord` | The 2026-09-23 probe found both markers in the lines of today's encoder, so the byte search reds. |
| LE5 | 3 | The resource block of a line holds exactly the keys `service.name`, `service.version`, and `bench.record.schema` | planned encoder test in `internal/otelrecord` | A filter that keeps a `telemetry.sdk.*` key or a detector key fails the exact-set comparison. |
| LE6 | 4 | A span with the undeclared attribute `bench.payload` encodes to a line with no `bench.payload` key | planned encoder test in `internal/otelrecord` | Today's encoder copies every attribute, so the key survives and the absence read reds. |
| LE7 | 4 | A span that carries every declared attribute encodes to a line that keeps each declared key | planned encoder test in `internal/otelrecord` | A filter that drops too much loses a declared key, so the presence read reds. |
| LE8 | 5 | `ReadSelected` reports a line whose schema attribute is `2` as a malformed-line problem and returns no span from it | planned reader test in `internal/otelrecord` | A reader that ignores the resource returns the span, so the problem read reds. |
| LE9 | 5 | `ReadSpans` returns no span from a line whose schema attribute is `2` | planned reader test in `internal/otelrecord` | A reader that ignores the resource returns the span, so the count reds. |
| LE10 | 6 | `ReadSpans` returns the finished span of a legacy line that has no schema attribute | planned reader test in `internal/otelrecord` | A reader that requires the attribute drops the legacy span, so the count reds. |
| LE11 | 7 | A context extracted from the handoff environment of a parent span starts a child whose trace id and parent span id are the parent's | planned handoff test in `internal/otelrecord` | A handoff that omits the traceparent starts a new trace, so the id comparison reds. |
| LE12 | 7 | The gate package names no handoff variable and calls the record package for the phase handoff | review-owned | The move changes no behavior, so `TestGateEnvStripsTheRecordVariables` and the phase-span system test stay green either way, and review grades the single owner. |
| LE13 | 8 | With a 1 KiB segment limit, an append past the limit renames `traces.jsonl` to one sealed `traces-<stamp>-<pid>.jsonl` segment and writes the line to a new live segment | planned retention test in `internal/otelrecord` | A writer without rotation keeps one growing file, so the sealed-segment read reds. |
| LE14 | 9 | After the retained count plus two rotations, exactly the retained count of sealed segments remains, and the two oldest names are gone | planned retention test in `internal/otelrecord` | A writer that never prunes keeps every segment, so the count reds. |
| LE15 | 10 | `ReadSpans` returns the spans of two sealed segments in name order and then the spans of the live segment | planned reader test in `internal/otelrecord` | A reader of the live segment alone drops the sealed spans, so the ordered comparison reds. |
| LE16 | 10 | `NewestLanding` returns the stages of a landing whose phase spans sit in a sealed segment and whose landing span sits in the live segment | planned reader test in `internal/otelrecord` | A live-only read finds the landing without its stages, so the stage count reds. |
| LE17 | 10 | `ReadSelected` returns a selected span that sits in a sealed segment | planned reader test in `internal/otelrecord` | A live-only read misses the span, so the selection reds. |
| LE18 | 11 | Two writers with separate handles that append across forced rotations leave every line whole, and the line total equals the append total | planned retention test in `internal/otelrecord` | A rotation that copies or truncates the live segment loses or splits a line, so the parse or the count reds. |
| LE19 | 12 | While another open file holds the rotation lock, an append past the limit writes its line to the live segment and creates no sealed segment | planned retention test in `internal/otelrecord` | A blocking lock wait hangs past the test deadline, and a lock-free rename seals the segment, so the absence read reds. |
| LE20 | 13 | A FIFO at a sealed segment name makes `ReadSpans` return an error within the test deadline | planned reader test in `internal/otelrecord` | An open of the FIFO blocks the read, so the deadline reds. |
| LE21 | 13 | A symlink at a sealed segment name makes `ReadSpans` return an error | planned reader test in `internal/otelrecord` | A reader that follows the link reads bytes outside the record, so the error read reds. |
| LE22 | 14 | A green one-iteration shift writes a start line and an end line for one `shift` span | planned shift test in `internal/shift` through `Loop` with a private Bench home | A shift with no span writes no line, so the pair read reds. |
| LE23 | 14 | A shift that exhausts its branch-name retries ends its `shift` span with `bench.shift.outcome` `usage` and `bench.work.state` `failed` | planned shift test in `internal/shift` with the branch collision fixture | A span that starts after the branch step never records this exit, so the end-line read reds. |
| LE24 | 15 | With `BENCH_AGENT` at `<tmp>/DIRMARK dir/agent`, the `shift` span carries `bench.agent` `agent`, and no line holds `DIRMARK` | planned shift test in `internal/shift` | A span that copies the adapter path carries the directory, so the byte search reds. |
| LE25 | 16 | With `BENCH_MODEL=mid` and `BENCH_MAX_ITERS=2`, the `shift` span carries `bench.line.tier` `mid` and `bench.shift.cap` `2` | planned shift test in `internal/shift` | A span without the declared line lacks both keys, so the value read reds. |
| LE26 | 16 | With `BENCH_MODEL=custom model`, the `shift` span carries no `bench.line.tier` key | planned shift test in `internal/shift` | A span that copies the raw value puts operator text in the record, so the absence read reds. |
| LE27 | 17 | A green one-iteration shift ends with `bench.shift.outcome` `complete`, `bench.work.state` `completed`, and `bench.outcome` `green` | planned shift test in `internal/shift` | A span that records only the exit lacks the two FT71 keys, so the value read reds. |
| LE28 | 17 | A shift whose first gate is red ends with `bench.shift.outcome` `failed`, `bench.work.state` `failed`, and `bench.outcome` `red` | planned shift test in `internal/shift` | A mapping that reads every stop as completed reds the work-state read. |
| LE29 | 17 | A shift that commits and then reaches its iteration cap ends with `bench.shift.outcome` `incomplete` and `bench.work.state` `completed` | planned shift test in `internal/shift` | A mapping that reads incomplete as failed reds the work-state read. |
| LE30 | 17 | A shift whose adapter makes no change ends with `bench.shift.outcome` `no-op` and `bench.work.state` `completed` | planned shift test in `internal/shift` | A mapping with no no-op arm leaves the work state empty, so the value read reds. |
| LE31 | 18 | A shift with one commit carries the branch head commit as `bench.subject.id` on its `shift` span | planned shift test in `internal/shift` | A span without the commit reference lacks the key, so the equality reds. |
| LE32 | 18 | A no-op shift carries no `bench.subject.id` on its `shift` span | planned shift test in `internal/shift` | A span that records the base commit claims work the shift never made, so the absence read reds. |
| LE33 | 19 | A red shift that retained its dirty worktree carries `bench.recovery.kind` `worktree` and a `bench.recovery.key` equal to the base name of the retained path | planned shift test in `internal/shift` with the retain fixture | A span without the reference lacks both keys, so the value read reds. |
| LE34 | 19 | No line of the retained shift holds the absolute path of the Bench home | planned shift test in `internal/shift` with the retain fixture | A span that copies the recovery pointer holds the home path, so the byte search reds. |
| LE35 | 20 | A green shift carries `bench.cleanup` `released` | planned shift test in `internal/shift` | A span with no cleanup key reds the value read. |
| LE36 | 20 | A red shift that retained its worktree carries `bench.cleanup` `retained` | planned shift test in `internal/shift` with the retain fixture | A span that records released on every path reds the value read. |
| LE37 | 24 | A shift that receives SIGINT during its adapter writes an end line for its `shift` span with `bench.work.state` `interrupted` | planned shift test in `internal/shift` that runs `Loop` in a helper process | The checkpoint exits through `os.Exit`, so a deferred end writes no line and the read reds. |
| LE38 | 21 | A two-iteration shift writes exactly two `shift.iteration` spans whose parent is the `shift` span | planned shift test in `internal/shift` | A loop with no pass span writes no pass line, so the count reds. |
| LE39 | 21 | A pass whose adapter exits 3 carries `bench.adapter.result` `exited` and `bench.adapter.exit` `3` | planned shift test in `internal/shift` | A pass span without the adapter result lacks both keys, so the value read reds. |
| LE40 | 21 | A pass whose adapter file lost its execute bit carries `bench.adapter.result` `spawn-failed` and no `bench.adapter.exit` key | planned shift test in `internal/shift` | A span that records a spawn failure as an exit reds the result read. |
| LE41 | 21 | A refactor pass writes one `shift.refactor` span whose parent is the `shift` span | planned shift test in `internal/shift` with an over-budget touched file | A refactor loop with no pass span writes no refactor line, so the read reds. |
| LE42 | 22 | Each pass span has exactly one child `gate` span, and that child carries `bench.subject.id` and `bench.outcome` | planned shift test in `internal/shift` | A gate run on a fresh context starts a new trace, so the parent comparison reds. |
| LE43 | 23 | A pass that committed carries a `bench.subject.id` equal to the commit that the branch gained in that pass | planned shift test in `internal/shift` | A pass without the commit key reds the equality. |
| LE44 | 24 | A shift that receives SIGINT during its adapter writes the end line of the open `shift.iteration` span before the end line of the `shift` span | planned shift test in `internal/shift` that runs `Loop` in a helper process | An exit path that ends only the shift span leaves the pass open, so the order read reds. |
| LE45 | 25 | The adapter environment carries `BENCH_OTEL_ROOT` and a `BENCH_OTEL_TRACEPARENT` whose span id is the current pass span | planned shift test in `internal/shift` with an adapter that writes its environment | An adapter launch without the handoff leaves the variables out, so the comparison reds. |
| LE46 | 25 | `bench resolve-model --harness claude` with the handoff, a routed binding, and `BENCH_MODEL=mid` writes a `line.resolve` span under the handoff span with harness `claude`, tier `mid`, and the bound model | planned command test in `cmd/bench` | An uninstrumented resolver writes no line span, so the read reds. |
| LE47 | 25 | `bench resolve-model` with no handoff environment writes a `line.resolve` span that has no parent | planned command test in `cmd/bench` | A resolver that records only under a handoff misses a standalone call, so the read reds. |
| LE48 | 25 | A refused resolution in a routed repo with no `BENCH_MODEL` writes a `line.resolve` span with `bench.outcome` `red` and no `bench.line.model` key | planned command test in `cmd/bench` | A span that records the refusal as green reds the outcome read. |
| LE49 | 26 | After a green shift whose adapter appended `MEMMARK` to the notes, one memory file below the record directory holds the notes bytes | planned shift test in `internal/shift` | A teardown that deletes the notes first leaves no memory file, so the read reds. |
| LE50 | 27 | After a red shift that retained its worktree, one memory file holds the notes bytes | planned shift test in `internal/shift` with the retain fixture | A retention on the release path only leaves no file here, so the read reds. |
| LE51 | 28 | The `shift` span carries `bench.memory.state` `retained`, the byte count, and a `bench.memory.digest` equal to the SHA-256 of the memory file | planned shift test in `internal/shift` | A span without the reference lacks the keys, so the equality reds. |
| LE52 | 28 | No record line holds `MEMMARK` | planned shift test in `internal/shift` | A span event with the notes body puts the text in the record, so the byte search reds. |
| LE53 | 29 | When the adapter replaces the notes with a symlink to a file that holds `SECRETMARK`, the span carries `bench.memory.state` `refused`, and no memory file holds `SECRETMARK` | planned shift test in `internal/shift` | A retention that follows the link copies the secret, so the byte search reds. |
| LE54 | 29 | When the adapter replaces the notes with a FIFO, the shift exits within the test deadline, and the span carries `bench.memory.state` `refused` | planned shift test in `internal/shift` | A retention that opens the FIFO blocks the teardown, so the deadline reds. |
| LE55 | 30 | When the adapter deletes the notes, the span carries `bench.memory.state` `absent`, and no memory file is written | planned shift test in `internal/shift` | A retention that writes an empty file for an absent source makes the two states equal, so the read reds. |
| LE56 | 30 | When the notes stay empty, the span carries `bench.memory.state` `retained` and `bench.memory.bytes` `0` | planned shift test in `internal/shift` | A retention that skips empty notes records absent, so the read reds. |
| LE57 | 31 | After the retained count plus one memory writes, exactly the retained count of memory files remains, and the oldest is gone | planned retention test in `internal/otelrecord` | A store that never prunes keeps every file, so the count reds. |
| LE58 | 31 | A symlinked memory directory makes the memory write return an error and write no file | planned retention test in `internal/otelrecord` | A store that follows the link writes outside the Bench home, so the refusal read reds. |
| LE88 | 29 | When the notes exceed the control-record limit, the span carries `bench.memory.state` `refused`, and no memory file is written | planned shift test in `internal/shift` | A retention that reads without the bound copies an unbounded file, so the read reds. |
| LE59 | 32 | The first adapter of a second shift reads an empty `.bench-notes.md`, and its stdin prompt holds no `MEMMARK` | planned shift test in `internal/shift` that runs two shifts | A shift that seeds its notes from the retained memory puts `MEMMARK` in the worktree, so the read reds. |
| LE60 | 33 | `Live` drops an entry with the outcome `complete` and the recovery `none` whose worktree exists | planned policy test in `internal/intent/admissionpolicy` | Today's rule keeps every entry whose worktree exists, so the entry stays live and the read reds. |
| LE61 | 33 | After a green shift exits, `intent.Snapshot` holds no entry for its key | planned shift test in `internal/shift` | A shift whose outcome does not end its intent leaves a live entry, so the read reds. |
| LE62 | 34 | `Live` keeps an entry with the outcome `failed` and a `worktree:` recovery whose worktree exists | planned policy test in `internal/intent/admissionpolicy` | A rule that drops every terminal entry loses the recovery pointer, so the read reds. |
| LE63 | 35 | After the acquire, the ledger entry of the shift carries a `lease` value equal to the lease line that the acquire wrote | planned shift test in `internal/shift` with an adapter that copies the lease file | A shift that never records the lease leaves the field empty, so the equality reds. |
| LE64 | 36, 43 | After a helper shift gets SIGKILL during its adapter, the pass writes a `shift.recovery` span with `bench.work.state` `recovered` and the intent key of the crashed `shift` span | planned recovery test in `internal/shift` with a killed helper process | A pass that never judges open entries writes no recovery span, so the read reds. |
| LE65 | 36 | The recovered entry carries the outcome `recovered` | planned recovery test in `internal/shift` with a killed helper process | A pass that records only the span leaves the entry open, so the outcome read reds. |
| LE91 | 36 | After the recovery of a crashed shift whose adapter appended `MEMMARK` to the notes, one memory file holds the notes bytes | planned recovery test in `internal/shift` with a killed helper process | A pass that recovers before it retains the memory loses the notes on the release path, so the read reds. |
| LE66 | 37 | When the crashed worktree holds a dirty file, the pass leaves the worktree locked with that file's bytes unchanged | planned recovery test in `internal/shift` with a killed helper process | A pass that releases a dirty worktree resets the file, so the byte comparison reds. |
| LE67 | 37 | The recovered entry of a dirty worktree carries a `worktree:` recovery pointer and stays live | planned recovery test in `internal/shift` with a killed helper process | A pass that drops the pointer hides the preserved work, so the read reds. |
| LE68 | 38 | When the crashed worktree holds only scratch files, the pass removes the lease and the scratch files and records `bench.cleanup` `released` | planned recovery test in `internal/shift` with a killed helper process | A pass that retains a clean worktree keeps the lease, so the lease read reds. |
| LE69 | 39 | When another identity holds the lease, the pass sets the outcome `abandoned` and writes a `shift.recovery` span with `bench.work.state` `abandoned` | planned recovery test in `internal/shift` with a rewritten lease | A pass that trusts the path alone recovers the checkout of another owner, so the state read reds. |
| LE70 | 39 | When another identity holds the lease, the pass leaves the lease bytes, the worktree files, and the lock state unchanged | planned recovery test in `internal/shift` with a rewritten lease | A pass that releases on a mismatch removes the lease of the other owner, so the byte comparison reds. |
| LE89 | 39 | When the lease file is absent, the pass sets the outcome `abandoned` and leaves every worktree file unchanged | planned recovery test in `internal/shift` with a removed lease | A pass that reads absence as a dead owner releases the worktree, so the file comparison reds. |
| LE71 | 40 | While the helper shift is alive, the pass leaves its entry unchanged and writes no recovery span | planned recovery test in `internal/shift` with a live helper process | A pass that judges by the lease match alone recovers a running shift, so the read reds. |
| LE72 | 41 | An entry with no lease and no worktree whose key names a dead process gets the outcome `abandoned` | planned recovery test in `internal/shift` with a seeded ledger | A pass that skips entries with no lease leaves the stale intent live, so the read reds. |
| LE73 | 41 | An entry with no lease whose key names the live test process stays unchanged | planned recovery test in `internal/shift` with a seeded ledger | A pass that abandons every entry with no lease drops a running shift, so the read reds. |
| LE74 | 41 | An entry whose key does not parse as a shift key stays unchanged | planned recovery test in `internal/shift` with a seeded ledger | A parse that reads an unknown key as a dead process abandons it, so the read reds. |
| LE75 | 42 | An entry whose lease file holds a malformed line stays unchanged and gets no recovery span | planned recovery test in `internal/shift` with a seeded lease | A pass that reads a malformed lease as dead recovers with no proof, so the read reds. |
| LE90 | 42 | An entry whose lease path is a FIFO stays unchanged, and the pass returns within the test deadline | planned recovery test in `internal/shift` with a FIFO lease | A pass that opens the FIFO blocks the session start, so the deadline reds. |
| LE76 | 44 | The session-inspect sequence runs the recovery pass after the resume phase, and a seeded dead-key entry is abandoned after `Inspect` | planned session-inspect test in `internal/sessioninspect` | A sequence without the phase leaves the entry open, so the read reds. |
| LE77 | 44 | `bench shift` runs the recovery pass before its acquire, and a seeded dead-key entry is abandoned after the shift | planned recovery test in `internal/shift` | A shift that skips the pass leaves the entry open, so the read reds. |
| LE78 | 45 | A pass that abandoned one entry prints exactly `bench shift recovery: recovered 0, abandoned 1` | planned recovery test in `internal/shift` | A pass that prints nothing or another form reds the exact comparison. |
| LE79 | 45 | A pass with no open shift entry prints nothing | planned recovery test in `internal/shift` | A pass that always prints adds a line to every session start, so the empty read reds. |
| LE80 | 46, 47 | A normal shell exit writes a `worktree.shell` span whose end line carries the assignment id, `bench.work.state` `completed`, and `bench.cleanup` `released` | planned subshell test in `internal/worktree` with a private Bench home | A session without a span writes no line, so the read reds. |
| LE81 | 48 | A signalled session writes an end line with `bench.work.state` `interrupted` and `bench.cleanup` `retained` | planned subshell test in `internal/worktree` with the signal helper | A span that records every exit as completed reds the state read. |
| LE82 | 49 | A session whose shell path does not exist writes an end line with `bench.work.state` `failed` | planned subshell test in `internal/worktree` | A span that records only a started shell misses this exit, so the read reds. |
| LE83 | 46 | `TestWorktreeSeamsMatchTheRegistry` holds `worktree.shell` in both the package seams and the registry rows | existing seam test in `internal/worktree` | A registry row with no package seam constant, or the reverse, fails the set comparison. |
| LE84 | 50 | A `bench shift` run through the built binary, with an adapter that calls `bench resolve-model`, writes one trace in which the `shift` span is an ancestor of the pass, `line.resolve`, and `gate` spans | planned system test in `internal/systemtest` | A broken handoff at any hop starts a second trace, so the ancestry read reds. |
| LE85 | 51 | That run, with `OBJMARK` in the objective, `DIRMARK` in the adapter directory, and `ENVMARK` in `OTEL_RESOURCE_ATTRIBUTES`, leaves no marker in any record segment | planned system test in `internal/systemtest` | A leak through any hop puts a marker in the record, so the byte search reds. |
| LE86 | 52 | `DATA_HANDLING.md` states that local records are mutable evidence inputs and not a tamper-proof central audit system | review-owned | Prose has no mechanical seam here, so review grades the sentence against the source row. |
| LE87 | 53 | `DATA_HANDLING.md` names the sealed segments, the memory files, the three bound entries, and the three resource keys | review-owned | Prose has no mechanical seam here, so review grades each named item against the tree. |

Not covered: story 54 — the reviewed exclusion changes no behavior, and the Out of scope section prices it.
Not covered: story 55 — the reviewed exclusion changes no behavior, and the Won't handle line below records it.
Not covered: story 56 — the reviewed exclusion changes no behavior, and LE86 documents the posture.

### Edge inventory

The canonical edge classes and the profile's hostile-input checklist, walked at each seam:

- Paths with spaces and glob characters: LE24 puts a space in the adapter directory. No path reaches a shell, because the record carries base names only.
- Control bytes in git-sourced text: no commit subject, branch name, or path enters the record. Commit ids and base names are the only git-derived values.
- A command whose write changes a fact it reports: the recovery pass counts its own verdicts after it writes them (LE78).
- A path read out of a file: the pass reads the worktree path from the ledger entry. The shift wrote that path as the absolute pool path. The pass compares lease bytes and resolves no relative path.
- A last line with no newline: the lease parse requires one final newline. So a lease with no newline is malformed, and the pass skips its entry (LE75).
- Absent versus empty: the notes file (LE55, LE56) and the lease file (LE89) each take a row for absence.
- Special files: a FIFO at a sealed segment (LE20), at the notes path (LE54), and at the lease path (LE90) is refused before any open.
- A dangling or live symlink: the notes read refuses both forms through the no-follow producer read (LE53). A symlinked sealed segment (LE21) and a symlinked memory directory (LE58) are refused.
- State written by one process and read by a fresh one: LE64 kills a helper process and runs the pass in the test process.
- Destructive worktree state: an identity mismatch (LE69, LE70), an absent lease (LE89), and a live owner (LE71) all change no worktree.
- A test that swaps a package variable: unexported constructors and one setter supply the limits, the counts, and the version in process. The system rows use the stamped binary and the real bounds.
- Concurrent writers: LE18 and LE19 cover two writers and a held rotation lock.

**Won't handle** lines:

- Harness session boundaries — no harness session-end event reaches a Bench seam. The `hook.session-inspect` span already marks each harness session start, and the shell session rows LE80 to LE82 survive.
- The effort of a line — effort has no enforcement surface, so no Bench seam knows it. The declared tier (LE25) and the resolved model (LE46) survive.
- A tamper-proof central audit system — the record is a local file that its owner can change. LE86 documents that posture.
- An older Bench binary that reads a ledger holding the new `lease` key — the strict decode refuses the unknown key. `Assignment.CreatedAt` and `RequestToken` set this precedent, and the current binary reads both shapes (LE63).
- A reader in the moment between a rename and the next append — it sees no live segment and reports the record absent once. The next read succeeds, and LE15 survives.
- The existing oversized live record — the first rotation seals it whole, and the count retention ages it out. LE14 survives.
- A span status description — no seam sets one. The declared-key filter covers every attribute, and LE6 survives.
- Free text inside a declared value, such as the lane diagnostic — review grades the value source of each declared key. LE7 survives.
- Hook spans inside an iteration — they keep their own traces. The one resolver joins the pass trace, and LE46 survives.
- The `claude-agent` and `worktree` intent kinds — they keep their current liveness rules. LE60 survives for the shift kind.
- A host path with a symlinked component, such as `/var` on macOS — the pass locks the recorded path, as the shift does today. LE66 survives.
- A backward clock step between two rotations — a newer sealed segment can sort before an older one, but the readers still return every span. LE15 survives.

## Ownership fences

- `reviews/local-shift-evidence.md`
- `internal/otelrecord/provider.go`
- `internal/otelrecord/processor_test.go`
- `internal/otelrecord/encode.go`
- `internal/otelrecord/encode_test.go`
- `internal/otelrecord/attributes.go`
- `internal/otelrecord/registry.go`
- `internal/otelrecord/reader.go`
- `internal/otelrecord/reader_test.go`
- `internal/otelrecord/writer.go`
- `internal/otelrecord/retention_test.go`
- `internal/gate/telemetry.go`
- `internal/gate/runner.go`
- `internal/gate/gate.go`
- `internal/gate/otel_env_test.go`
- `internal/bounds/bounds.go`
- `tests/canary/package-core-guard/bounds-duplicate-owner`
- `cmd/bench/process_env.go`
- `cmd/bench/main.go`
- `cmd/bench/guards.go`
- `cmd/bench/otel_hook_seams_test.go`
- `tests/canary/package-core-guard/unrouted-subcommand`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `internal/shift/loop.go`
- `internal/shift/session.go`
- `internal/shift/result.go`
- `internal/shift/record.go`
- `internal/shift/record_test.go`
- `internal/shift/recover.go`
- `internal/shift/recover_test.go`
- `internal/intent/intent.go`
- `internal/intent/ledger/ledger.go`
- `internal/intent/admissionpolicy/admissionpolicy.go`
- `internal/intent/admissionpolicy/liveness_test.go`
- `internal/sessioninspect/sessioninspect.go`
- `internal/sessioninspect/sessioninspect_test.go`
- `internal/worktree/subshell.go`
- `internal/worktree/clean.go`
- `internal/worktree/verb_span.go`
- `internal/worktree/subshell_test.go`
- `internal/worktree/otel_seams_test.go`
- `internal/systemtest/otel_gate_test.go`
- `internal/systemtest/otel_verbs_test.go`
- `DATA_HANDLING.md`
- `tests/canary/data-handling-derivation/undocumented-passlist-var`

Reviewer disposition: pending. The fence also needs one structure grant for the shift package, which the Further notes section states.

## Ticket graph

| ticket | Blocked by | chunk |
| --- | --- | --- |
| `1-move-the-trace-handoff-into-the-record.md` | none | LE-A |
| `2-redact-and-version-every-record-line.md` | none | LE-A |
| `3-rotate-and-retain-the-record.md` | none | LE-A |
| `4-record-the-shift-boundaries.md` | `2-redact-and-version-every-record-line.md` | LE-B1 |
| `5-record-each-pass-under-the-shift.md` | `4-record-the-shift-boundaries.md` | LE-B2 |
| `6-record-the-resolved-line.md` | `1-move-the-trace-handoff-into-the-record.md`, `5-record-each-pass-under-the-shift.md` | LE-B2 |
| `7-retain-the-shift-memory.md` | `3-rotate-and-retain-the-record.md`, `4-record-the-shift-boundaries.md` | LE-B2 |
| `8-complete-the-shift-intent-on-exit.md` | none | LE-C1 |
| `9-recover-crashed-shifts-by-lease-identity.md` | `4-record-the-shift-boundaries.md`, `7-retain-the-shift-memory.md`, `8-complete-the-shift-intent-on-exit.md` | LE-C2 |
| `10-record-the-worktree-shell-session.md` | `4-record-the-shift-boundaries.md` | LE-D |
| `11-prove-the-shift-trace-through-the-built-binary.md` | `2-redact-and-version-every-record-line.md`, `5-record-each-pass-under-the-shift.md`, `6-record-the-resolved-line.md` | LE-D |
| `12-document-the-local-evidence-contract.md` | `3-rotate-and-retain-the-record.md`, `7-retain-the-shift-memory.md`, `9-recover-crashed-shifts-by-lease-identity.md`, `10-record-the-worktree-shell-session.md` | LE-D |

## Out of scope

- The `bank.ft71.local_event` producer — about 7 edits and 3 gate runs. It needs a fixture-journey script, a native workflow step, the release runbook, and the workflow conformance pin. FT88 closed with no producer for its own record, so the tree has no producer precedent to follow. The retirement drain moves the registry owner to a new roadmap row.
- A hook span inside an iteration that joins the pass trace — about 4 edits and 1 gate run.
- An operator setting for the retention sizes — about 3 edits and 1 gate run.

## Further notes

### Source trace

| source sentence | rows or exclusion |
| --- | --- |
| Add a redacted, append-only schema for local events. | LE1, LE4, LE5, LE6, LE7, LE85 |
| It records shift/session boundaries, | LE22, LE23, LE37, LE80, LE81, LE82 |
| the resolved agent and line, | LE24, LE25, LE26, LE46, LE47, LE48 |
| and the gate fingerprint and verdict. | LE42 |
| It also records the adapter result, | LE39, LE40 |
| commit or recovery reference, | LE31, LE33, LE43, LE67 |
| cleanup decision, | LE35, LE36, LE68, LE80 |
| and release-evidence relationship. | LE2, LE3; the producer is out of scope |
| Atomic append, rotation, and explicit retention are part of the repository contract. | LE13, LE14, LE18, LE19, LE57 |
| Records distinguish completed, failed, interrupted, recovered, and abandoned work; | LE27, LE28, LE37, LE64, LE69 |
| normal exit completes the matching intent | LE60, LE61 |
| and crash recovery requires the active lease identity. | LE63, LE64, LE69, LE70, LE71 |
| Local records are documented as mutable evidence inputs, not a tamper-proof central audit system. | LE86 |
| Closure includes success, failure, interruption, recovery, normal subshell completion, stale intent, and redaction fixtures. | LE27, LE28, LE37, LE64, LE80, LE72, LE85 |
| The shift record retains the per-iteration memory, the notes or a later bounded state, before the scratch cleanup. | LE49, LE50, LE51 |
| Retained evidence never re-enters the iteration prompt. | LE59 |

The source sentence "The repository-controlled bank evidence requirement makes this row active" states a priority, not a behavior.

### Sources re-read

- `RR:C-05` is row C-05 of the 2026-07-11 release-readiness assessment, read at `af19d71c:ASSESSMENT.md`. The current `ASSESSMENT.md` is a later refresh that no longer carries it.
- `RC:H-03` is finding H-03 of the 2026-07-11 compliance assessment, read at `62f92b2e^:COMPLIANCE_ASSESSMENT.md`, the last version before its removal.
- The C-05 subshell half is fixed today. `TestSubshellNormalExitReleasesItsAssignment` proves that a normal exit releases the assignment. The shift half is still open, and LE60 and LE61 own it.
- `capture/audits/skill-state-shift-memory.md` was re-read. The paper behind it was not re-read.
- The FT274 decision was read at `868ff55c:roadmap/FT274.md` and in the retired spec at `8ae8b80d^:specs/otel-seam-record/spec.md`. That spec priced FT71's rotation, retention, and absorption.
- The SDK behavior was read in the module cache at `go.opentelemetry.io/otel/sdk@v1.46.0`: `resource.Default` runs the `fromEnv` detector, and `WithResource` merges `resource.Environment`. The 2026-09-23 probe confirmed both markers in a record line.

### Reader sweep and proof checklist

- Cited symbols: each symbol resolves at `ba9b8621`. The record package: `Encode`, `Begin`, `BeginIn`, `NewProvider`, `Writer.Append`, `gradeRecordPath`, `ReadSpans`, `ReadSelected`, `NewestLanding`, `DeclaredAttributes`, and `Registry`. The gate: `beginGateSpan`, `withGateSpanEnv`, `beginPhaseRecord`, `otelGateEnv`, and `RunAndRecordContext`. The shift: `finish`, `evidenceResult`, `preserveAndRecover`, `teardown`, `cleanupScratch`, `runAdapter`, `runGate`, `refactorPhase`, and `exitPreserving`. Others: `admissionpolicy.Live`, `intent.NewEntry`, `worktree.RetainAndLock`, `worktree.Release`, `worktree.LeaseFile`, `worktree.ProbeLease`, `lifecyclepolicy.LeaseOwnerPID`, `bounds.ClassifyNoFollow`, `lines.Tiers`, `modelid.SafeToken`, `prepareProcessEnvironment`, and `subshellAt`.
- Import edges: the shift package gains `internal/otelrecord`, `internal/lines`, and `internal/bounds`, and none of the three imports the shift package. The session-inspect package gains `internal/shift`, and the shift package imports no session-inspect code.
- Source-row clauses and occurrences: each clause of the row maps in the source trace above. The one occurrence names the 2026-08-30 assessment; `bench status --all` on 2026-09-23 still shows no intent row, and the ledger holds zero entries and 45 assignments.
- Promised field labels, in four groups:
  - Resource keys: `bench.record.schema`, `service.name`, and `service.version`.
  - Shift keys: `bench.intent.key`, `bench.agent`, `bench.shift.cap`, `bench.shift.outcome`, `bench.work.state`, `bench.cleanup`, `bench.recovery.kind`, and `bench.recovery.key`.
  - Pass, line, and memory keys: `bench.adapter.result`, `bench.adapter.exit`, `bench.line.harness`, `bench.line.tier`, `bench.line.model`, `bench.memory.state`, `bench.memory.bytes`, and `bench.memory.digest`.
  - Ledger key: `lease`.
- Changed-function callers:
  - `Encode` and `Writer.Append` each have one production caller, the processor.
  - `ReadSpans` has one, `NewestLanding`, and `NewestLanding` has one, the retro scaffold.
  - `ReadSelected` has one, the assessment collection.
  - `withGateSpanEnv` has one caller in `engine.go`, and `beginPhaseRecord` has one in `runner.go`.
  - `admissionpolicy.Live` has two, `Snapshot` and `Compact`.
  - `finish` has every exit of `loop` and `exitPreserving`, and `subshellAt` has `Subshell` and its tests.
- Readers of the changed artifacts: the assessment collection checks the live path before `ReadSelected`, and the retro scaffold calls `NewestLanding`. Tests read the live path in `internal/gate/lane_record_test.go`, `internal/worktree/reset_plan_test.go`, `internal/roadmap/retro_scaffold_test.go`, and the system helpers that glob `traces.jsonl`. None of these tests rotates, so each keeps reading the live segment. `decisions/software-factory/assets/ft303-cli-assessment.md` quotes the path as history and takes a named exclusion. `internal/systemtest/otel_worktree_test.go` asserts that every attribute is declared, and the filter keeps that assertion true.
- Pinning consumers: `internal/conformance/bounds_policy_test.go` requires `internal/shift/loop.go` to keep its five iteration bound names. The implementation keeps them there.
- Copy survival: LE12 grades the one handoff owner. No other copy of a replaced fact survives.

### Structure and fence disposition

- `internal/shift/` holds 11 files and the lane caps a directory at 12. This spec adds `record.go`, `record_test.go`, `recover.go`, and `recover_test.go`. The reviewer decides the grant `internal/shift/ 15` in `.bench/structure.budgets` before ticket 4 starts. The alternative is a split of the shift tests, which fragments one package's evidence.
- `internal/otelrecord/` gains one file, `retention_test.go`, and reaches the cap of 12.
- `cmd/bench/main.go`, `internal/gate/runner.go`, and `internal/worktree/subshell.go` are over their line budgets. Tickets 1, 2, 6, and 10 move code out so that none of them grows. Ticket 10 moves the ignored-inventory function from `subshell.go` into `clean.go`.
- The command registry, its two tests, and the two conformance registry tests join the fence through the binding registry closure only. The `cmd/bench` and `internal/worktree` packages are bound packages, and the build expects no edit to those five files.
- The three canary fixture paths join the fence through the fixture closure only. Each pins a file that a ticket writes, and the build expects no edit to a fixture.
- `cmd/bench/`, `internal/gate/`, `internal/systemtest/`, and `internal/worktree/` have no directory headroom, so this spec adds no file there.

### Late questions for the reviewer

1. Does this spec ship the `bank.ft71.local_event` producer? Recommended answer: no. The producer is priced under Out of scope as its own capability, because FT88 closed with no producer and the release track stays NO-GO.
2. Where does the retained notes text live? Recommended answer: in a 0600 memory file beside the record, referenced by digest. A span event would put model prose in the record, and the FT274 declared set forbids a payload.

### Flagged additions

- `service.version` on every line is the spec's reading of "release-evidence relationship". The source names the relationship but not its field.
- The declared-key filter in the encoder turns the FT274 review-owned rule into a mechanical one. It is the mechanism for "redacted".
- The `line.resolve` span records the resolved model at its one resolver. The shift itself knows only the declared tier.
- The worktree shell is the "session" of "shift/session boundaries", because the source's C-05 names the subshell. Harness sessions take a Won't handle line.
- The two recovery triggers, the printed line, and the memory state vocabulary are mechanisms that the source did not name.
- The retention sizes are proposals that the reviewer owns.

### Completion plan

```bench-completion-plan
{"version":1,"chunks":[{"id":"LE-A","tickets":["1-move-the-trace-handoff-into-the-record.md","2-redact-and-version-every-record-line.md","3-rotate-and-retain-the-record.md"],"verification":[{"id":"otelrecord","command":"bench test --package ./internal/otelrecord"},{"id":"gate","command":"bench test --package ./internal/gate"},{"id":"cmd","command":"bench test --package ./cmd/bench"},{"id":"system","command":"bench test --check system"}]},{"id":"LE-B1","tickets":["4-record-the-shift-boundaries.md"],"verification":[{"id":"shift","command":"bench test --package ./internal/shift"},{"id":"kit-compliance","command":"bench test --check kit-compliance"}]},{"id":"LE-B2","tickets":["5-record-each-pass-under-the-shift.md","6-record-the-resolved-line.md","7-retain-the-shift-memory.md"],"verification":[{"id":"shift","command":"bench test --package ./internal/shift"},{"id":"otelrecord","command":"bench test --package ./internal/otelrecord"},{"id":"cmd","command":"bench test --package ./cmd/bench"},{"id":"kit-compliance","command":"bench test --check kit-compliance"}]},{"id":"LE-C1","tickets":["8-complete-the-shift-intent-on-exit.md"],"verification":[{"id":"intent","command":"bench test --package ./internal/intent/..."},{"id":"shift","command":"bench test --package ./internal/shift"}]},{"id":"LE-C2","tickets":["9-recover-crashed-shifts-by-lease-identity.md"],"verification":[{"id":"shift","command":"bench test --package ./internal/shift"},{"id":"intent","command":"bench test --package ./internal/intent/..."},{"id":"sessioninspect","command":"bench test --package ./internal/sessioninspect"}]},{"id":"LE-D","tickets":["10-record-the-worktree-shell-session.md","11-prove-the-shift-trace-through-the-built-binary.md","12-document-the-local-evidence-contract.md"],"verification":[{"id":"subshell","command":"bench test --package ./internal/worktree --run TestSubshell"},{"id":"seams","command":"bench test --package ./internal/worktree --run TestWorktreeSeamsMatchTheRegistry"},{"id":"system","command":"bench test --check system"},{"id":"prose","command":"bench gate-prose . -- DATA_HANDLING.md"}]}],"final_verification":[{"id":"coverage","command":"bench coverage --check specs/local-shift-evidence/spec.md"},{"id":"otelrecord","command":"bench test --package ./internal/otelrecord"},{"id":"shift","command":"bench test --package ./internal/shift"},{"id":"kit-compliance","command":"bench test --check kit-compliance"},{"id":"system","command":"bench test --check system"}]}
```
