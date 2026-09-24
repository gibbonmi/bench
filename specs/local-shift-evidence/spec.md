# Versioned local shift evidence

Status: staged

Roadmap: FT71

Decision source: the named reviewed artifact `roadmap/FT71.md`, the row body that the 2026-09-22 drain `d-465b789396f9` settled at `65ac3e2a` with `Next: spec`.

The row inherits two closed decisions. On 2026-08-29 the reviewer decided that FT71's events are span attributes in the FT274 seam record, not a second ledger. The FT274 crash test passed, so that condition holds. The 2026-08 capability audit closed L-04: no new canonical work record exists, and this evidence must not become one.

Verification log: 3 iteration(s) to accept — iteration 1 rejected on 2 blockers and 15 other findings (R1–R17). Iteration 2 rejected on 2 major and 5 minor findings (R18–R24), and the reviewer raised the cap to 3. Iteration 3 accepted with 4 minor folds (R25–R28) and no further review.

## Problem

Bench keeps no trustworthy local evidence of a shift. A shift writes an intent entry, iteration commits, and gate spans, but no record joins them. No record names the adapter, the line, the adapter result, the recovery reference, or the cleanup decision. A green shift deletes `.bench-notes.md` at teardown, so the memory of its iterations is gone.

The intent ledger also misleads. A shift that exits normally records its outcome, but its entry stays live while its pool worktree exists, and pool worktrees persist. A shift that crashes leaves an open entry that nothing ever resolves. No rule separates a crashed owner from a live one.

The seam record itself is not redacted. A probe on 2026-09-23 found `OTEL_RESOURCE_ATTRIBUTES` and `OTEL_SERVICE_NAME` values in every record line, because the OpenTelemetry SDK merges them into the default resource. `DATA_HANDLING.md` says that a line carries no environment value, so the claim is false today. The record also grows without a bound: the live record of this repository holds 196 MB after three weeks.

## Solution

The seam record becomes the versioned local evidence. Every line names its record schema and the Bench version that wrote it, and the encoder writes only declared keys. No environment value reaches a line. The record rotates into sealed segments of a bounded size and keeps a bounded count of them.

A shift writes one `shift` span for its whole life and one pass span for each iteration and refactor pass. Each gate run is the child of its pass, so the gate fingerprint and verdict join their iteration. The spans name the adapter, the declared tier, the iteration cap, and each adapter result and commit. They also name the shift outcome, a work state, the recovery reference, and the cleanup decision. The adapter's own call to `bench resolve-model` records the resolved model in the same trace.

Before any scratch cleanup, the shift retains its notes as a private memory file beside the record. The span references that file by digest, and no prompt reads it.

A shift that exits normally ends its intent. The shift records its lease identity, and a recovery pass resolves each crashed shift. A crashed shift is recovered only when its lease identity still holds and its owner is dead. The pass takes the dead lease through the pool's own takeover protocol. It locks a dirty worktree in place, releases a clean one, and finishes an interrupted recovery on a re-run. Otherwise the pass abandons the entry and leaves its worktree untouched.

The worktree shell records its session boundaries the same way. `DATA_HANDLING.md` documents the record as mutable local evidence, not a tamper-proof central audit system.

## User stories

Line: opus / medium.
Implementation-line reason: LE-C3 is the hardest chunk, because the identity claim races other acquirers and its tests kill real processes. LE-A carries rotation under concurrent writers. The rows are exact, the seams exist, and focused package tests red each wrong answer.
Harder chunks: LE-A, LE-C3.

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
44. As an operator, I want the recovery pass at session start and before `bench shift` acquires, so that stale intent needs no manual step.
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

Recovery and record safety:

57. As an operator, I want recovery to take the dead lease through the pool's takeover protocol, so that a competing claimant keeps its tree.
58. As an operator, I want a re-run to finish a recovery that a kill interrupted, so that no recovery stays half done.
59. As an operator, I want a re-run over a finished recovery to change nothing, so that the pass stays idempotent.
60. As an operator, I want sealed segments named by a lock-held sequence, so that no rotation overwrites one and no clock step misorders them.
61. As a retro reader, I want a failed memory write recorded as failed, so that the span never leaves the memory keys undefined.
62. As a record consumer, I want a process with no version to write no version key, so that no false version appears.
63. As an operator, I want an entry skipped while a live other owner holds its lease, so that no verdict races a claim.
64. As an operator, I want a recovered entry left alone under another lease identity, so that a resume never touches a new owner.

## Implementation decisions

### The record contract

- The FT274 seam record is the one record. The record package owns the address, the segments, the encoder, the memory files, and the trace handoff.
- The encoder writes its own resource block for every line. After the version is set, the block holds exactly `service.name` with the value `bench`, `service.version`, and `bench.record.schema` with the value `1`. The SDK resource never reaches a line, so no detector or environment value can enter.
- The command layer hands its stamped version to the record package once, at process start. A process that sets no version writes no `service.version` key, so its block holds exactly two keys.
- The declared attribute set becomes the redaction contract. The encoder drops each span attribute whose key is not declared. A later key ships only when its ticket declares it.
- A reader treats a line whose `bench.record.schema` is not `1` as a malformed line. A line with no schema attribute is a legacy line and reads as before.
- `service.version` is the relationship to release evidence. A release producer selects the lines whose version equals its envelope's package version. The producer itself is out of scope.
- The gate fingerprint is the existing gate span's `bench.subject.id`, and the gate verdict is its `bench.outcome`. This spec adds no oracle key.

### Rotation and retention

- The live segment stays `traces.jsonl`. When an append would take it past `bounds.RecordSegmentLimit`, the writer seals it as `traces-<sequence>.jsonl` and appends to a new live segment.
- Only the rotation step takes a lock: a non-blocking exclusive lock on one lock file in the record directory. A writer that cannot take the lock appends to the live segment and rotates on a later append. Each append stays one synchronous `O_APPEND` write.
- Under the lock, the writer checks the live size again and reads the sealed names. The next sequence is one more than the highest sequence present. The listing holds every sealed name, planted or not, and only a lock holder seals, so that name is free. Then the writer renames the live segment to that name. When the highest sequence is the largest one, the writer refuses the rotation and appends to the live segment.
- The sequence is a zero-padded 20-digit decimal, so name order equals sequence order. No clock enters the name, so two rotations in one instant never share a name and a clock step never changes the order.
- After a rotation, the writer removes the lowest sequences until `bounds.RecordSegmentsRetained` sealed segments remain.
- One helper in the record package formats and parses the sealed name. The writer, the prune, and both readers use that helper.
- The readers read the sealed segments in sequence order and then the live segment. Each segment passes the writer's grade, `gradeRecordPath`, before the open. `ReadSpans` has no grade today, so this spec adds it.
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
- A pass span carries `bench.adapter.result` (`exited` or `spawn-failed`), `bench.adapter.exit` when the adapter exited with a code, and the commit that the pass made. An adapter that a signal ends records `exited` with no exit key. Every exit path ends an open pass span before it ends the `shift` span.
- The shift hands the adapter the trace handoff for the current pass span. `bench resolve-model` joins that trace when the handoff is present and records a `line.resolve` span. That span carries the harness, the declared tier, the resolved model when it is a safe model token, and the outcome.

### Shift memory

- The shift retains its notes on every path that reached the first iteration: before the teardown's scratch cleanup, and on the retain path. It reads `.bench-notes.md` through the bounded no-follow producer read, so a symlink, a FIFO, an oversized file, or invalid UTF-8 reads as refused.
- The record package writes the bytes to `memory/<UTC stamp>-<trace id>.md` below the record directory with mode 0600. It keeps the newest `bounds.RecordMemoryRetained` files.
- The `shift` span carries `bench.memory.state`: `retained`, `absent`, `refused`, or `failed`. A retained file adds `bench.memory.bytes` and `bench.memory.digest`. `failed` means the store write failed, and the shift outcome does not change. The record never holds the notes text.
- No prompt reads the memory directory. A new shift writes an empty notes file, as it does today.

### Intent completion and crash recovery

- The liveness rule drops an entry that holds an outcome and a recovery of empty or `none`. An entry with a `worktree:` recovery stays live while its worktree exists, even when its branch reads as landed.
- The ledger entry gains an optional `lease` field. Right after the acquire, the shift records its lease line without the final newline.
- The recovery pass lives in the shift package, because it reuses the shift's scratch policy, memory retention, and preservation. It judges each shift entry that has no outcome, and it resumes each entry with the outcome `recovered` whose act did not finish.
- For an entry with no lease, the intent package parses the owner process from the key beside `NewEntry`. A dead owner abandons the entry. A live owner or an unparsable key skips it.
- For an entry with a lease, the pass first grades the lease file without a follow. An absent file abandons the entry. A file that is not regular, is unreadable, or holds a malformed lease line skips the entry. This grade comes before the comparison, so a malformed line never counts as another owner (LE75).
- A well-formed lease line that differs from the recorded line abandons the entry only when that line's owner is dead. A differing line with a live owner skips the entry, because a concurrent pass or a new acquirer holds the lease. The recorded line with a live owner skips the entry. The recorded line with a dead owner starts a recovery.
- A recovery first retains the memory, before any claim. A concurrent acquire of a clean tree runs `git clean -qfdx`, so the notes must leave the tree before the pass competes for it. A concede then leaves an orphan memory file, and the memory prune bounds it.
- The recovery then takes the lease through the pool's takeover protocol. The protocol renames the lease to `.stale.<pid>`, compares the moved bytes with the recorded line, and creates the pass's own lease with an exclusive create. The worktree package owns that protocol, so a new identity claim beside `RetainAndLock` runs `claimAt` with the recorded line as the one accepted judgment. If the claim loses, the pass concedes and changes nothing more.
- After a won claim, the pass writes the entry: its own lease line and the outcome `recovered`. A dirty worktree gets the `worktree:` pointer, and a clean worktree gets the recovery `none`. Last, the pass acts on the worktree.
- The act follows `preserveAndRecover`. A worktree with dirty paths beyond the scratch is locked in place, and the pass removes its own lease. The lock step is skipped when the worktree is already locked. A dirty tree is never an acquire candidate, so its lock holds. A clean worktree loses its scratch and is released to the pool.
- After the act, the pass writes its entry once more with the outcome `recovered` and its pointer. So a concurrent abandon that landed between the two writes loses, and a dirty tree keeps its live pointer.
- The reviewer decided on 2026-09-23 that a recovery releases a clean crashed worktree and locks only a dirty one. A lock does not keep a clean tree out of the pool. `acquireAt` picks a tree by its clean state and never reads a git lock.
- A resume arm judges each entry with the outcome `recovered`. When its lease file holds the entry's lease line with a dead owner, the act did not finish. The pass then runs only the act again. The resume retains no second memory file and writes one recovery span with no memory keys.
- Every other `recovered` entry stays as it is. No lease file means that the act finished, and a live owner means that a pass still acts. Another identity means that a new owner holds the tree. A malformed or special lease file gives no verdict.
- An abandon changes no worktree file, lease, or lock. The pass sets the entry outcome to the work-state word, `recovered` or `abandoned`, and records one `shift.recovery` span with the intent key. The span starts a new trace, and the intent key joins it to the crashed `shift` span.
- Two triggers run the pass: a session-inspect phase after the resume phase, and `bench shift` before its acquire. `worktree.Acquire` has that one production caller. The pass prints `bench shift recovery: recovered <n>, abandoned <m>` only when it acted.

### The worktree shell session

- `bench worktree shell` opens one `worktree.shell` span before it creates the assignment and ends it at every return. The end line carries the assignment id, the work state, the cleanup decision, and the outcome.
- A normal shell exit is `completed`. The cleanup is `released` when the release exits 0 and `retained` otherwise. A signal is `interrupted` with `retained`, because the lease stays for a reclaim. A shell that cannot start is `failed`.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| LE-A / `1-move-the-trace-handoff-into-the-record.md`, `2-redact-and-version-every-record-line.md`, `3-rotate-and-retain-the-record.md` | Every line is versioned and redacted, the record rotates and retains by sequence, and the trace handoff has one owner. | LE1, LE2, LE3, LE4, LE5, LE96, LE6, LE7, LE8, LE9, LE10, LE11, LE12, LE13, LE92, LE94, LE106, LE14, LE93, LE15, LE16, LE17, LE18, LE19, LE20, LE21 | `bench test --package ./internal/otelrecord`, `bench test --package ./internal/gate`, `bench test --package ./cmd/bench`, `bench test --check system` | yes |
| LE-B1 / `4-record-the-shift-boundaries.md` | A shift writes one `shift` span with its line, outcome, work state, commit, recovery reference, and cleanup decision on every in-process exit path. | LE22, LE23, LE24, LE25, LE26, LE27, LE28, LE29, LE30, LE31, LE32, LE33, LE34, LE35, LE36, LE107 | `bench test --package ./internal/shift`, `bench test --check kit-compliance` | no |
| LE-B2 / `5-record-each-pass-under-the-shift.md`, `6-record-the-resolved-line.md`, `7-retain-the-shift-memory.md` | Each pass, its gate, and the resolved line join the shift trace, an interrupted shift ends its spans, and the notes survive as a private memory file. | LE37, LE38, LE39, LE40, LE41, LE42, LE43, LE44, LE45, LE46, LE47, LE48, LE49, LE50, LE51, LE52, LE53, LE54, LE55, LE56, LE57, LE58, LE88, LE95, LE102, LE59 | `bench test --package ./internal/shift`, `bench test --package ./internal/otelrecord`, `bench test --package ./cmd/bench`, `bench test --check kit-compliance` | no |
| LE-C1 / `8-complete-the-shift-intent-on-exit.md` | A normal exit ends the shift intent, and each shift records its lease identity. | LE60, LE61, LE62, LE63 | `bench test --package ./internal/intent/...`, `bench test --package ./internal/shift` | no |
| LE-C2 / `9-abandon-lease-less-stale-intent.md` | A recovery pass with two triggers abandons each lease-less entry whose key owner is dead and prints one line when it acts. | LE72, LE73, LE74, LE76, LE77, LE78, LE79 | `bench test --package ./internal/shift`, `bench test --package ./internal/intent/...`, `bench test --package ./internal/sessioninspect` | no |
| LE-C3 / `10-recover-crashed-shifts-by-lease-identity.md` | The pass recovers or abandons each leased entry by its lease identity through the takeover protocol, locks a dirty tree or releases a clean one, and a re-run finishes or keeps a recovery. | LE64, LE65, LE91, LE105, LE66, LE67, LE68, LE69, LE70, LE89, LE103, LE71, LE75, LE90, LE97, LE98, LE99, LE101, LE100, LE104 | `bench test --package ./internal/shift`, `bench test --package ./internal/worktree --run TestClaimRecordedLease` | yes |
| LE-D / `11-record-the-worktree-shell-session.md`, `12-prove-the-shift-trace-through-the-built-binary.md`, `13-document-the-local-evidence-contract.md` | The shell session is on record, the built binary proves the whole trace and the redaction, and `DATA_HANDLING.md` states the contract. | LE80, LE81, LE82, LE83, LE84, LE85, LE86, LE87 | `bench test --package ./internal/worktree --run TestSubshell`, `bench test --package ./internal/worktree --run TestWorktreeSeamsMatchTheRegistry`, `bench test --check system`, `bench gate-prose . -- DATA_HANDLING.md` | no |

Review iteration 1 split the old LE-C2 into LE-C2 and LE-C3. The old ticket 9 became tickets 9 and 10, and the old tickets 10, 11, and 12 are now tickets 11, 12, and 13. LE37 moved from LE-B1 to LE-B2, beside LE44. Review iteration 2 added LE103 to LE105 to LE-C3.

## Testing decisions

- A good test drives the real entry with a private Bench home and reads the record lines back. The line content is the external behavior. `TestLoopRetainsAndLocksDirtyWorktree` and the branch collision fixture are the in-process shift precedents.
- The encoder, reader, rotation, handoff, and memory store get unit tests in the record package. The fixture spans of `encode_test.go` and the two-writer test of `writer_test.go` are the precedents.
- The shift tests drive `Loop` in process for normal exits. For SIGINT and SIGKILL, they re-exec the shift test binary into one role test of the shift package, because the checkpoint exits through `os.Exit`.
- That helper-process start is the Go re-exec idiom that each package authors for its own roles, as `internal/gocache/lock_test.go` and `internal/testreport/cancel_test.go` do. It copies nothing from the worktree package's `descendant`, whose knowledge is the worktree journey census record. So no shared start owner is planned.
- The identity claim gets a worktree test through the existing takeover gap seam, beside `TestRetainAndLockLocksDropsLeaseAndPreservesDirt`. The recovery pass reaches its concede through the shift package's step fault seam in process. It reaches its interrupted act in a helper role beside ticket 5's role. That role arms the act fault from a test-only variable, runs the pass, and exits. The test waits for it, so the first pass's lease line names a dead pid and not a zombie that `PIDAlive` reads as alive.
- The liveness rule gets a pure policy test beside `TestCompactionRuleReadsTypedLivenessFacts`. The resolve-model span gets a command test beside the hook seam test in `cmd/bench`.
- The system suite proves the stamped version and the whole trace through the built binary. `TestOtelCrashKeepsStartedPhaseLine` and the verb span tests are the precedents.
- The gate observes the feature through the `test` phase and the `system` phase. No new gate phase is added. The existing seam check in `kit-compliance` reds a registered seam whose symbol starts no span.
- A planned row names no `_test.go` citation, because a citation must resolve to a declared test today. The build adds each citation when its test lands. LE83 cites an existing test, and LE86 and LE87 carry the review-owned marker.

### Seam diagram

    trigger: `bench shift`, `bench worktree shell`, `bench resolve-model`, session start
        │
        ▼
    verb entry ──▶ [ shift record: shift, pass, and recovery spans; memory retention ] ──▶ span lines
                      ◀ tests attach here: run Loop or the recovery pass with a private
                        Bench home, then read the lines and the ledger back
        │
        ▼
    span ──▶ [ record package: redacting encoder, sequenced writer, memory store ] ──▶ <home>/otel/<key>/
                      ◀ tests attach here: encode fixture spans, append across a small
                        segment limit, read the segments and memory files back
        │
        ▼
    ledger entry ──▶ [ liveness rule, lease identity, takeover claim ] ──▶ live, recovered, or abandoned
                      ◀ tests attach here: pure policy facts; the takeover gap seam;
                        a killed helper shift; the step fault seam

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| LE1 | 1 | An encoded line carries the resource attribute `bench.record.schema` with the value `1` | `internal/otelrecord/encode_test.go` (`TestEncodeWritesTheRecordSchema`) | An encoder that keeps the SDK resource has no schema key, so the value read reds. |
| LE2 | 2 | After the version is set to `9.9.9-test`, an encoded line carries `service.version` with that value | `internal/otelrecord/encode_test.go` (`TestEncodeWritesTheSetVersion`) | An encoder that never reads the set version writes no such key, so the value read reds. |
| LE3 | 2 | A `bench gate` run through the built binary writes lines whose `service.version` equals the second field of the `bench version` line, such as `0.2.0` | `internal/systemtest/otel_gate_test.go` (`TestOtelGateRecordNamesTheStampedVersion`) | A binary that never hands its stamped version to the record writes no version, so the equality reds. |
| LE4 | 3 | With `OTEL_RESOURCE_ATTRIBUTES=bench.leak=PROBEMARK` and `OTEL_SERVICE_NAME=SVCMARK` set, no line that `Begin` writes holds either marker | `internal/otelrecord/processor_test.go` (`TestBeginWritesNoEnvironmentResource`) | The 2026-09-23 probe found both markers in the lines of today's encoder, so the byte search reds. |
| LE5 | 3 | After the version is set, the resource block of a line holds exactly the keys `service.name`, `service.version`, and `bench.record.schema` | `internal/otelrecord/encode_test.go` (`TestEncodeResourceHoldsExactlyTheBenchKeys`) | A filter that keeps a `telemetry.sdk.*` key or a detector key fails the exact-set comparison. |
| LE96 | 62 | With no version set, the resource block of a line holds exactly the keys `service.name` and `bench.record.schema` | `internal/otelrecord/encode_test.go` (`TestEncodeWithoutAVersionWritesNoVersionKey`) | An encoder that writes an empty or `dev` version claims a release it never had, so the exact-set comparison reds. |
| LE6 | 4 | A span with the undeclared attribute `bench.payload` encodes to a line with no `bench.payload` key | `internal/otelrecord/encode_test.go` (`TestEncodeDropsAnUndeclaredAttribute`) | Today's encoder copies every attribute, so the key survives and the absence read reds. |
| LE7 | 4 | A span that carries every declared attribute encodes to a line that keeps each declared key | `internal/otelrecord/encode_test.go` (`TestEncodeKeepsEveryDeclaredAttribute`) | A filter that drops too much loses a declared key, so the presence read reds. |
| LE8 | 5 | `ReadSelected` reports a line whose schema attribute is `2` as a malformed-line problem and returns no span from it | `internal/otelrecord/reader_test.go` (`TestReadSelectedReportsAnUnknownSchemaAsMalformed`) | A reader that ignores the resource returns the span, so the problem read reds. |
| LE9 | 5 | `ReadSpans` returns no span from a line whose schema attribute is `2` | `internal/otelrecord/reader_test.go` (`TestReadSpansReturnsNoSpanOfAnUnknownSchema`) | A reader that ignores the resource returns the span, so the count reds. |
| LE10 | 6 | `ReadSpans` returns the finished span of a legacy line that has no schema attribute | `internal/otelrecord/reader_test.go` (`TestReadSpansReadsALegacyLine`) | A reader that requires the attribute drops the legacy span, so the count reds. |
| LE11 | 7 | A context extracted from the handoff environment of a parent span starts a child whose trace id and parent span id are the parent's | `internal/otelrecord/processor_test.go` (`TestAHandedOffChildJoinsItsParentSpan`) | A handoff that omits the traceparent starts a new trace, so the id comparison reds. |
| LE12 | 7 | No non-test Go file in `internal/gate` holds the literal `BENCH_OTEL_` | `internal/gate/otel_env_test.go` (`TestTheGateSpellsNoHandoffVariable`) | A surviving copy of a handoff variable name in the gate package reds the literal search. |
| LE13 | 8 | With a 1 KiB segment limit, an append past the limit renames `traces.jsonl` to the sealed segment of the next sequence and writes the line to a new live segment | `internal/otelrecord/retention_test.go` (`TestAnAppendPastTheLimitSealsTheLiveSegment`) | A writer without rotation keeps one growing file, so the sealed-segment read reds. |
| LE92 | 60 | Two rotations in a row leave two sealed segments with consecutive sequences, and the first keeps every line it sealed | `internal/otelrecord/retention_test.go` (`TestTwoRotationsSealConsecutiveSequences`) | A name built from a coarse clock stamp repeats in one stamp unit, so the rename overwrites the first segment and the line count reds. |
| LE94 | 60 | A planted file at the next sequence name stays unchanged, and the rotation seals the live segment under the following sequence | `internal/otelrecord/retention_test.go` (`TestARotationSkipsAPlantedSequenceName`) | A rotation that ignores the present names renames over the planted file, so the byte comparison reds. |
| LE106 | 60 | A planted file at the largest sequence name stays unchanged, and an append past the limit returns an error, seals no segment, and keeps its line in the live segment | `internal/otelrecord/retention_test.go` (`TestARotationRefusesTheLastSequence`) | A rotation whose next sequence wraps to 0 seals under sequence 0 and renames over it later, so the error read reds. |
| LE14 | 9 | After the retained count plus two rotations, exactly the retained count of sealed segments remains, and the two lowest sequences are gone | `internal/otelrecord/retention_test.go` (`TestThePruneKeepsTheRetainedCount`) | A writer that never prunes keeps every segment, so the count reds. |
| LE93 | 60 | The prune removes the lowest sequences even when their modification times are the newest in the directory | `internal/otelrecord/retention_test.go` (`TestThePruneRemovesTheLowestSequenceWhateverItsTime`) | A prune by modification time removes the newest segment, so the survivor set reds. |
| LE15 | 10 | `ReadSpans` returns the spans of two sealed segments in sequence order and then the spans of the live segment | `internal/otelrecord/reader_test.go` (`TestReadSpansReadsTheSealedSegmentsInOrder`) | A reader of the live segment alone drops the sealed spans, so the ordered comparison reds. |
| LE16 | 10 | `NewestLanding` returns the stages of a landing whose phase spans sit in a sealed segment and whose landing span sits in the live segment | `internal/otelrecord/reader_test.go` (`TestNewestLandingReadsStagesFromASealedSegment`) | A live-only read finds the landing without its stages, so the stage count reds. |
| LE17 | 10 | `ReadSelected` returns a selected span that sits in a sealed segment | `internal/otelrecord/reader_test.go` (`TestReadSelectedReadsASealedSegment`) | A live-only read misses the span, so the selection reds. |
| LE18 | 11 | Two writers with separate handles that append across forced rotations leave every line whole, and the line total equals the append total | `internal/otelrecord/retention_test.go` (`TestTwoWritersAcrossRotationsLeaveOnlyWholeLines`) | A rotation that copies or truncates the live segment loses or splits a line, so the parse or the count reds. |
| LE19 | 12 | While another open file holds the rotation lock, an append past the limit writes its line to the live segment and creates no sealed segment | `internal/otelrecord/retention_test.go` (`TestAHeldRotationLockSkipsTheRotation`) | A blocking lock wait hangs past the test deadline, and a lock-free rename seals the segment, so the absence read reds. |
| LE20 | 13 | A FIFO at a sealed segment name makes `ReadSpans` return an error within the test deadline | `internal/otelrecord/reader_test.go` (`TestReadSpansRefusesAFIFOSegment`) | An open of the FIFO blocks the read, so the deadline reds. |
| LE21 | 13 | A symlink at a sealed segment name makes `ReadSpans` return an error | `internal/otelrecord/reader_test.go` (`TestReadSpansRefusesASymlinkedSegment`) | A reader that follows the link reads bytes outside the record, so the error read reds. |
| LE22 | 14 | A green one-iteration shift writes a start line and an end line for one `shift` span | `TestAGreenShiftWritesAStartAndAnEndLine` in `internal/shift` | A shift with no span writes no line, so the pair read reds. |
| LE23 | 14 | A shift that exhausts its branch-name retries ends its `shift` span with `bench.shift.outcome` `usage` and `bench.work.state` `failed` | `TestAShiftThatExhaustsItsBranchRetriesRecordsUsage` in `internal/shift` | A span that starts after the branch step never records this exit, so the end-line read reds. |
| LE24 | 15 | With `BENCH_AGENT` at `<tmp>/DIRMARK dir/agent`, the `shift` span carries `bench.agent` `agent`, and no line holds `DIRMARK` | `TestTheShiftSpanCarriesTheAdapterBaseName` in `internal/shift` | A span that copies the adapter path carries the directory, so the byte search reds. |
| LE25 | 16 | With `BENCH_MODEL=mid` and `BENCH_MAX_ITERS=2`, the `shift` span carries `bench.line.tier` `mid` and `bench.shift.cap` `2` | `TestTheShiftSpanCarriesTheTierAndTheCap` in `internal/shift` | A span without the declared line lacks both keys, so the value read reds. |
| LE26 | 16 | With `BENCH_MODEL=custom model`, the `shift` span carries no `bench.line.tier` key | `TestTheShiftSpanCarriesNoTierForAnOperatorModel` in `internal/shift` | A span that copies the raw value puts operator text in the record, so the absence read reds. |
| LE27 | 17 | A green one-iteration shift ends with `bench.shift.outcome` `complete`, `bench.work.state` `completed`, and `bench.outcome` `green` | `TestACompleteShiftRecordsCompletedWork` in `internal/shift` | A span that records only the exit lacks the two FT71 keys, so the value read reds. |
| LE28 | 17 | A shift whose first gate is red ends with `bench.shift.outcome` `failed`, `bench.work.state` `failed`, and `bench.outcome` `red` | `TestAFailedShiftRecordsFailedWork` in `internal/shift` | A mapping that reads every stop as completed reds the work-state read. |
| LE29 | 17 | A shift that commits and then reaches its iteration cap ends with `bench.shift.outcome` `incomplete` and `bench.work.state` `completed` | `TestAnIncompleteShiftRecordsCompletedWork` in `internal/shift` | A mapping that reads incomplete as failed reds the work-state read. |
| LE30 | 17 | A shift whose adapter makes no change ends with `bench.shift.outcome` `no-op` and `bench.work.state` `completed` | `TestANoOpShiftRecordsCompletedWork` in `internal/shift` | A mapping with no no-op arm leaves the work state empty, so the value read reds. |
| LE31 | 18 | A shift with one commit carries the branch head commit as `bench.subject.id` on its `shift` span | `TestACommittedShiftCarriesTheBranchHead` in `internal/shift` | A span without the commit reference lacks the key, so the equality reds. |
| LE32 | 18 | A no-op shift carries no `bench.subject.id` on its `shift` span | `TestANoOpShiftCarriesNoSubject` in `internal/shift` | A span that records the base commit claims work the shift never made, so the absence read reds. |
| LE33 | 19 | A red shift that retained its dirty worktree carries `bench.recovery.kind` `worktree` and a `bench.recovery.key` equal to the base name of the retained path | `TestARetainedShiftCarriesTheRecoveryBaseName` in `internal/shift` | A span without the reference lacks both keys, so the value read reds. |
| LE34 | 19 | No line of the retained shift holds the absolute path of the Bench home | `TestARetainedShiftRecordHoldsNoHomePath` in `internal/shift` | A span that copies the recovery pointer holds the home path, so the byte search reds. |
| LE35 | 20 | A green shift carries `bench.cleanup` `released` | `TestAGreenShiftRecordsReleasedCleanup` in `internal/shift` | A span with no cleanup key reds the value read. |
| LE36 | 20 | A red shift that retained its worktree carries `bench.cleanup` `retained` | `TestARetainedShiftRecordsRetainedCleanup` in `internal/shift` | A span that records released on every path reds the value read. |
| LE107 | 19 | A green shift that left no recovery pointer carries `bench.recovery.kind` `none` and no `bench.recovery.key` | `TestAGreenShiftRecordsNoRecoveryKind` in `internal/shift` | A span that writes the kind only for a retained worktree lacks `bench.recovery.kind`, so the value read reds. |
| LE37 | 24 | A shift that receives SIGINT during its adapter writes an end line for its `shift` span with `bench.work.state` `interrupted` | `TestAnInterruptedShiftRecordsInterruptedWork` in `internal/shift` | The checkpoint exits through `os.Exit`, so a deferred end writes no line and the read reds. |
| LE38 | 21 | A two-iteration shift writes exactly two `shift.iteration` spans whose parent is the `shift` span | `TestATwoIterationShiftWritesTwoPassSpans` in `internal/shift` | A loop with no pass span writes no pass line, so the count reds. |
| LE39 | 21 | A pass whose adapter exits 3 carries `bench.adapter.result` `exited` and `bench.adapter.exit` `3` | `TestAPassRecordsTheAdapterExit` in `internal/shift` | A pass span without the adapter result lacks both keys, so the value read reds. |
| LE40 | 21 | A pass whose adapter removed its own execute bit in iteration 1 carries `bench.adapter.result` `spawn-failed` and no `bench.adapter.exit` key in iteration 2 | `TestAPassRecordsASpawnFailure` in `internal/shift` | A span that records a spawn failure as an exit reds the result read. |
| LE41 | 21 | A refactor pass writes one `shift.refactor` span whose parent is the `shift` span | `TestARefactorPassWritesARefactorSpan` in `internal/shift` | A refactor loop with no pass span writes no refactor line, so the read reds. |
| LE42 | 22 | Each pass span has exactly one child `gate` span, and that child carries `bench.subject.id` and `bench.outcome` | `TestEachPassParentsOneGateSpan` in `internal/shift` | A gate run on a fresh context starts a new trace, so the parent comparison reds. |
| LE43 | 23 | A pass that committed carries a `bench.subject.id` equal to the commit that the branch gained in that pass | `TestACommittedPassCarriesItsCommit` in `internal/shift` | A pass without the commit key reds the equality. |
| LE44 | 24 | A shift that receives SIGINT during its adapter writes the end line of the open `shift.iteration` span before the end line of the `shift` span | `TestAnInterruptedShiftEndsThePassFirst` in `internal/shift` | An exit path that ends only the shift span leaves the pass open, so the order read reds. |
| LE45 | 25 | The adapter environment carries `BENCH_OTEL_ROOT` and a `BENCH_OTEL_TRACEPARENT` whose span id is the current pass span | `TestTheAdapterReceivesThePassHandoff` in `internal/shift` | An adapter launch without the handoff leaves the variables out, so the comparison reds. |
| LE46 | 25 | `bench resolve-model --harness claude` with the handoff, a routed binding, and `BENCH_MODEL=mid` writes a `line.resolve` span under the handoff span with harness `claude`, tier `mid`, and the bound model | `TestAHandedOffResolutionRecordsItsLine` in `cmd/bench` | An uninstrumented resolver writes no line span, so the read reds. |
| LE47 | 25 | `bench resolve-model` with no handoff environment writes a `line.resolve` span that has no parent | `TestAStandaloneResolutionRecordsARootSpan` in `cmd/bench` | A resolver that records only under a handoff misses a standalone call, so the read reds. |
| LE48 | 25 | A refused resolution in a routed repo with no `BENCH_MODEL` writes a `line.resolve` span with `bench.outcome` `red` and no `bench.line.model` key | `TestARefusedResolutionRecordsRed` in `cmd/bench` | A span that records the refusal as green reds the outcome read. |
| LE49 | 26 | After a green shift whose adapter appended `MEMMARK` to the notes, one memory file below the record directory holds the notes bytes | `TestAGreenShiftRetainsItsNotes` in `internal/shift` | A teardown that deletes the notes first leaves no memory file, so the read reds. |
| LE50 | 27 | After a red shift that retained its worktree, one memory file holds the notes bytes | `TestARedShiftRetainsItsNotes` in `internal/shift` | A retention on the release path only leaves no file here, so the read reds. |
| LE51 | 28 | The `shift` span carries `bench.memory.state` `retained`, the byte count, and a `bench.memory.digest` equal to the SHA-256 of the memory file | `TestAGreenShiftRetainsItsNotes` in `internal/shift` | A span without the reference lacks the keys, so the equality reds. |
| LE52 | 28 | No record line holds `MEMMARK` | `TestTheRecordHoldsNoNotesText` in `internal/shift` | A span event with the notes body puts the text in the record, so the byte search reds. |
| LE53 | 29 | When the adapter replaces the notes with a symlink to a file that holds `SECRETMARK`, the span carries `bench.memory.state` `refused`, and no memory file holds `SECRETMARK` | `TestEachNotesStateIsRecorded/symlink` in `internal/shift` | A retention that follows the link copies the secret, so the byte search reds. |
| LE54 | 29 | When the adapter replaces the notes with a FIFO, the shift exits within the test deadline, and the span carries `bench.memory.state` `refused` | `TestEachNotesStateIsRecorded/fifo` in `internal/shift` | A retention that opens the FIFO blocks the teardown, so the deadline reds. |
| LE88 | 29 | When the notes exceed the control-record limit, the span carries `bench.memory.state` `refused`, and no memory file is written | `TestEachNotesStateIsRecorded/oversized` in `internal/shift` | A retention that reads without the bound copies an unbounded file, so the read reds. |
| LE55 | 30 | When the adapter deletes the notes, the span carries `bench.memory.state` `absent`, and no memory file is written | `TestEachNotesStateIsRecorded/deleted` in `internal/shift` | A retention that writes an empty file for an absent source makes the two states equal, so the read reds. |
| LE56 | 30 | When the notes stay empty, the span carries `bench.memory.state` `retained` and `bench.memory.bytes` `0` | `TestEachNotesStateIsRecorded/empty` in `internal/shift` | A retention that skips empty notes records absent, so the read reds. |
| LE57 | 31 | After the retained count plus one memory writes, exactly the retained count of memory files remains, and the oldest is gone | `TestTheMemoryStoreKeepsTheRetainedCount` in `internal/otelrecord` | A store that never prunes keeps every file, so the count reds. |
| LE58 | 31 | A symlinked memory directory makes the memory write return an error and write no file | `TestTheMemoryStoreRefusesASymlinkedDirectory` in `internal/otelrecord` | A store that follows the link writes outside the Bench home, so the refusal read reds. |
| LE95 | 61 | With a symlinked memory directory, a green shift's span carries `bench.memory.state` `failed` | `TestAFailedMemoryWriteKeepsTheOutcome` in `internal/shift` | A shift that drops the store error leaves the memory keys undefined, so the value read reds. |
| LE102 | 61 | With a symlinked memory directory, a green shift still ends with `bench.shift.outcome` `complete` and exit 0 | `TestAFailedMemoryWriteKeepsTheOutcome` in `internal/shift` | A shift that lets a record failure change its outcome exits nonzero, so the exit read reds. |
| LE59 | 32 | The first adapter of a second shift reads an empty `.bench-notes.md`, and its stdin prompt holds no `MEMMARK` | `TestASecondShiftStartsWithEmptyNotes` in `internal/shift` | A shift that seeds its notes from the retained memory puts `MEMMARK` in the worktree, so the read reds. |
| LE60 | 33 | `Live` drops an entry with the outcome `complete` and the recovery `none` whose worktree exists | `TestLiveDropsAFinishedShiftWithNoRecovery` in `internal/intent/admissionpolicy` | Today's rule keeps every entry whose worktree exists, so the entry stays live and the read reds. |
| LE61 | 33 | After a green shift exits, `intent.Snapshot` holds no entry for its key | `TestAGreenShiftEndsItsIntent` in `internal/shift` | A shift whose outcome does not end its intent leaves a live entry, so the read reds. |
| LE62 | 34 | `Live` keeps an entry with the outcome `failed` and a `worktree:` recovery whose worktree exists | `TestLiveKeepsAShiftWithAWorktreeRecovery` in `internal/intent/admissionpolicy` | A rule that drops every terminal entry loses the recovery pointer, so the read reds. |
| LE63 | 35 | After the acquire, the ledger entry of the shift carries a `lease` value equal to the lease line that the acquire wrote | `TestAShiftRecordsItsLease` in `internal/shift` | A shift that never records the lease leaves the field empty, so the equality reds. |
| LE72 | 41 | An entry with no lease and no worktree whose key names a dead process gets the outcome `abandoned` | `TestRecoverAbandonsADeadOwnersEntry` in `internal/shift` | A pass that skips entries with no lease leaves the stale intent live, so the read reds. |
| LE73 | 41 | An entry with no lease whose key names the live test process stays unchanged | `TestRecoverKeepsALiveOrUnknownEntry` in `internal/shift` | A pass that abandons every entry with no lease drops a running shift, so the read reds. |
| LE74 | 41 | An entry whose key does not parse as a shift key stays unchanged | `TestRecoverKeepsALiveOrUnknownEntry` in `internal/shift` | A parse that reads an unknown key as a dead process abandons it, so the read reds. |
| LE76 | 44 | The session-inspect sequence runs the recovery pass after the resume phase, and a seeded dead-key entry is abandoned after `Inspect` | `TestInspectRecoversAfterTheResumePhase` in `internal/sessioninspect` | A sequence without the phase leaves the entry open, so the read reds. |
| LE77 | 44 | `bench shift` runs the recovery pass before its acquire, and a seeded dead-key entry is abandoned after the shift | `TestAShiftRecoversBeforeItsAcquire` in `internal/shift` | A shift that skips the pass leaves the entry open, so the read reds. |
| LE78 | 45 | A pass that abandoned one entry prints exactly `bench shift recovery: recovered 0, abandoned 1` | `TestRecoverAbandonsADeadOwnersEntry` in `internal/shift` | A pass that prints nothing or another form reds the exact comparison. |
| LE79 | 45 | A pass with no open shift entry prints nothing | `TestRecoverWithNothingToDoPrintsNothing` in `internal/shift` | A pass that always prints adds a line to every session start, so the empty read reds. |
| LE64 | 36, 43 | After a helper shift gets SIGKILL during its adapter, the pass writes a `shift.recovery` span with `bench.work.state` `recovered` and the intent key of the crashed `shift` span | `TestRecoverRecoversAKilledDirtyShift` in `internal/shift` | A pass that never judges leased entries writes no recovery span, so the read reds. |
| LE65 | 36 | The recovered entry carries the outcome `recovered` | `TestRecoverRecoversAKilledDirtyShift` in `internal/shift` | A pass that records only the span leaves the entry open, so the outcome read reds. |
| LE91 | 36 | After the recovery of a crashed shift whose adapter appended `MEMMARK` to the notes, one memory file holds the notes bytes | `TestRecoverRecoversAKilledDirtyShift` in `internal/shift` | A pass that acts before it retains the memory can lose the notes, so the read reds. |
| LE105 | 26 | When the identity claim concedes, one memory file still holds the crashed shift's notes | `TestRecoverConcedesAFaultedClaim` in `internal/shift` | A pass that retains the memory after the claim loses the notes to the winner's clean, so the read reds. |
| LE66 | 37 | When the crashed worktree holds a dirty file, the pass leaves the worktree locked with that file's bytes unchanged | `TestRecoverRecoversAKilledDirtyShift` in `internal/shift` | A pass that releases a dirty worktree resets the file, so the byte comparison reds. |
| LE67 | 37 | The recovered entry of a dirty worktree carries a `worktree:` recovery pointer and stays live | `TestRecoverRecoversAKilledDirtyShift` in `internal/shift` | A pass that drops the pointer hides the preserved work, so the read reds. |
| LE68 | 38 | When the crashed worktree holds only scratch files, the pass removes the lease and the scratch files and records `bench.cleanup` `released` | `TestRecoverReleasesAKilledCleanShift` in `internal/shift` | A pass that locks a clean worktree leaves a tree that the next acquire resets under a live entry, so the lease read reds. |
| LE69 | 39 | When another identity whose owner is dead holds the lease, the pass sets the outcome `abandoned` and writes a `shift.recovery` span with `bench.work.state` `abandoned` | `TestRecoverAbandonsAnotherDeadIdentity` in `internal/shift` | A pass that trusts the path alone recovers the checkout of another owner, so the state read reds. |
| LE70 | 39 | When another identity whose owner is dead holds the lease, the pass leaves the lease bytes, the worktree files, and the lock state unchanged | `TestRecoverAbandonsAnotherDeadIdentity` in `internal/shift` | A pass that acts on a mismatch removes the lease of the other owner, so the byte comparison reds. |
| LE89 | 39 | When the lease file is absent, the pass sets the outcome `abandoned` and leaves every worktree file unchanged | `TestRecoverAbandonsAnAbsentLease` in `internal/shift` | A pass that reads absence as a dead owner acts on the worktree, so the file comparison reds. |
| LE103 | 63 | When another identity whose owner is alive holds the lease, the pass leaves the entry unchanged and writes no recovery span | `TestRecoverKeepsAnUnprovenLease/live_other_owner` in `internal/shift` | A pass that abandons on any mismatch writes `abandoned` under a live claimant, so the entry read reds. |
| LE71 | 40 | While the helper shift is alive, the pass leaves its entry unchanged and writes no recovery span | `TestRecoverKeepsAnUnprovenLease/live_helper` in `internal/shift` | A pass that judges by the lease match alone recovers a running shift, so the read reds. |
| LE75 | 42 | An entry whose lease file holds a malformed line stays unchanged and gets no recovery span | `TestRecoverKeepsAnUnprovenLease/malformed` in `internal/shift` | A pass that reads a malformed lease as another owner abandons it with no proof, so the read reds. |
| LE90 | 42 | An entry whose lease path is a FIFO stays unchanged, and the pass returns within the test deadline | `TestRecoverKeepsAnUnprovenLease/fifo` in `internal/shift` | A pass that opens the FIFO blocks the session start, so the deadline reds. |
| LE97 | 57 | When another writer replaces the lease in the takeover gap, the identity claim returns false and the other writer's lease stays | `TestClaimRecordedLeaseConcedesToAWriterInTheGap` in `internal/worktree` | A claim that renames without the byte comparison takes the other writer's lease, so the lease read reds. |
| LE98 | 57 | A pass whose identity claim concedes leaves the entry, the lease, and the lock state unchanged and writes no recovery span | `TestRecoverConcedesAFaultedClaim` in `internal/shift` | A pass that acts after a lost claim locks the other owner's tree, so the read reds. |
| LE99 | 58 | A first pass runs in a helper role with the act fault armed and exits after its entry write, and then a second pass locks the dirty worktree and removes the first pass's lease | `TestRecoverFinishesAnInterruptedRecovery` in `internal/shift` | A pass that judges only open entries never finishes the act, so the lock read reds. |
| LE101 | 58 | After that second pass, exactly one memory file exists for the crashed shift | `TestRecoverFinishesAnInterruptedRecovery` in `internal/shift` | A resume that retains the memory again writes a second file, so the count reds. |
| LE100 | 59 | A second pass over a finished recovery changes no file, lease, lock, or entry and writes no recovery span | `TestRecoverRecoversAKilledDirtyShift` in `internal/shift` | A pass that re-judges a recovered entry rewrites it or retains the memory again, so the comparison reds. |
| LE104 | 64 | A `recovered` entry whose lease file holds another identity with a dead owner stays unchanged, and the pass changes no worktree file | `TestRecoverKeepsARecoveryUnderAnotherLease` in `internal/shift` | A resume that acts on any dead lease locks or releases the new owner's tree, so the comparison reds. |
| LE80 | 46, 47 | A normal shell exit writes a `worktree.shell` span whose end line carries the assignment id, `bench.work.state` `completed`, and `bench.cleanup` `released` | `TestSubshellNormalExitReleasesItsAssignment` in `internal/worktree` | A session without a span writes no line, so the read reds. |
| LE81 | 48 | A signalled session writes an end line with `bench.work.state` `interrupted` and `bench.cleanup` `retained` | `TestSubshellSignalsLeaveAReclaimableLease` in `internal/worktree` | A span that records every exit as completed reds the state read. |
| LE82 | 49 | A session whose shell path does not exist writes an end line with `bench.work.state` `failed` | `TestSubshellRecordsAShellThatCannotStart` in `internal/worktree` | A span that records only a started shell misses this exit, so the read reds. |
| LE83 | 46 | The worktree seam test holds `worktree.shell` in both the package seams and the registry rows | `internal/worktree/otel_seams_test.go` (`TestWorktreeSeamsMatchTheRegistry`) | A registry row with no package seam constant, or the reverse, fails the set comparison. |
| LE84 | 50 | A `bench shift` run through the built binary, with an adapter that calls `bench resolve-model`, writes one trace in which the `shift` span is an ancestor of the pass, `line.resolve`, and `gate` spans | `TestOtelShiftTraceJourney` in `internal/systemtest` | A broken handoff at any hop starts a second trace, so the ancestry read reds. |
| LE85 | 51 | That run, with `OBJMARK` in the objective, `DIRMARK` in the adapter directory, and `ENVMARK` in `OTEL_RESOURCE_ATTRIBUTES`, leaves no marker in any record segment | `TestOtelShiftTraceJourney` in `internal/systemtest` | A leak through any hop puts a marker in the record, so the byte search reds. |
| LE86 | 52 | `DATA_HANDLING.md` states that local records are mutable evidence inputs and not a tamper-proof central audit system | review-owned: prose has no mechanical seam here | Review grades the sentence against the source row. |
| LE87 | 53 | `DATA_HANDLING.md` names the sealed segments, the memory files, the three bound entries, and the three resource keys | review-owned: prose has no mechanical seam here | Review grades each named item against the tree. |

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
- A dangling or live symlink: the notes read refuses both forms through the no-follow producer read (LE53). A symlinked sealed segment (LE21) and a symlinked memory directory (LE58, LE95) are refused.
- State written by one process and read by a fresh one: LE64 kills a helper process and runs the pass in the test process. LE99 runs the first pass in a helper process that exits, so its lease line names a dead pid.
- Destructive worktree state: a dead other owner (LE69, LE70), an absent lease (LE89), and a live other owner (LE103) change no worktree. So do a live entry owner (LE71), a lost claim (LE97, LE98), and a `recovered` entry under another identity (LE104).
- A test that swaps a package variable: unexported constructors and one setter supply the limits, the counts, and the version in process. The system rows use the stamped binary and the real bounds.
- Concurrent writers: LE18 and LE19 cover two writers and a held rotation lock. LE92 and LE94 cover the sealed name, and LE97 covers a competing lease claimant.
- An interrupted act: LE99 and LE101 cover a pass stopped between its entry write and its act, and LE100 covers a re-run after the act.

**Won't handle** lines:

- Harness session boundaries — no harness session-end event reaches a Bench seam. The `hook.session-inspect` span already marks each harness session start, and the shell session rows LE80 to LE82 survive.
- The effort of a line — effort has no enforcement surface, so no Bench seam knows it. The declared tier (LE25) and the resolved model (LE46) survive.
- A tamper-proof central audit system — the record is a local file that its owner can change. LE86 documents that posture.
- An older Bench binary that reads a ledger holding the new `lease` key — the strict decode refuses the unknown key. `Assignment.CreatedAt` and `RequestToken` set this precedent, and the current binary reads both shapes (LE63).
- A reader in the moment between a rename and the next append — it sees no live segment and reports the record absent once. The next read succeeds, and LE15 survives.
- A sealed segment that the prune removes between a reader's list and its open — `ReadSpans` returns an error once, and the next read succeeds. LE15 survives.
- A rotation between a reader's list and its live open — the reader misses the newly sealed segment and reports no error. The next read sees that segment, and LE15 survives.
- The existing oversized live record — the first rotation seals it whole, and the count retention ages it out. LE14 survives.
- A planted file at the largest sequence name — the writer refuses each rotation, and the live segment grows until an operator removes the file. LE106 survives.
- An append or rename failure inside the span processor — the processor drops it, because the FT274 record never changes a verb outcome. The memory store failure keeps its own state, and LE95 survives.
- A span status description — no seam sets one. The declared-key filter covers every attribute, and LE6 survives.
- Free text inside a declared value, such as the lane diagnostic — review grades the value source of each declared key. LE7 survives.
- Hook spans inside an iteration — they keep their own traces. The one resolver joins the pass trace, and LE46 survives.
- The `claude-agent` and `worktree` intent kinds — they keep their current liveness rules. LE60 survives for the shift kind.
- A host path with a symlinked component, such as `/var` on macOS — the pass locks the recorded path, as the shift does today. LE66 survives.
- A reused process id — a crashed owner whose id a new process took reads as alive, so the entry stays open until that process exits. LE71 survives.
- A kill between the won claim and the entry write — the lease then holds the pass's own dead line. A re-run reads another identity and abandons the entry. The abandon changes nothing, the work stays in place, and LE99 survives.
- A kill between the memory write and the entry write, or a concede after the memory write — an orphan memory file remains. A re-run retains a second copy, and the memory prune bounds both. LE105 survives.
- A kill after the act and before the recovery span ends — the record holds no finished `shift.recovery` span for that verdict. The started span line holds the verdict, and for a dirty tree the entry holds it too. LE65 survives.
- A kill inside the release of a clean tree — the entry already holds the recovery `none`, so the compaction can drop it before a resume. The tree keeps a dead lease and no dirt, so the pool reclaims it through the ordinary acquire. LE68 survives.
- Two concurrent passes over one entry — one pass can read the vacant lease slot between the other's rename and create, and it records `abandoned`. The winner writes its entry again after its act, so `recovered` and its pointer win. The record holds both spans, and the abandoning pass changes no worktree. LE97 survives.

## Ownership fences

- `reviews/local-shift-evidence.md`
- `.bench/structure.budgets`
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
- `cmd/bench/main_test.go`
- `tests/canary/package-core-guard/unrouted-subcommand`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `internal/shift/loop.go`
- `internal/shift/session.go`
- `internal/shift/result.go`
- `internal/shift/fault.go`
- `internal/shift/fault_test.go`
- `internal/shift/record.go`
- `internal/shift/record_test.go`
- `internal/shift/pass_test.go`
- `internal/shift/shift_test.go`
- `internal/shift/recover.go`
- `internal/shift/recover_test.go`
- `internal/intent/intent.go`
- `internal/intent/ledger/ledger.go`
- `internal/intent/ledger_aliases.go`
- `internal/intent/admissionpolicy/admissionpolicy.go`
- `internal/intent/admissionpolicy/liveness_test.go`
- `internal/sessioninspect/sessioninspect.go`
- `internal/status/status.go`
- `tests/canary/docs-currency-token-diet/signal-vocabulary-drift`
- `internal/sessioninspect/sessioninspect_test.go`
- `internal/worktree/snapshot.go`
- `internal/worktree/snapshot_test.go`
- `internal/worktree/lifecycle.go`
- `internal/worktree/lifecycle_test.go`
- `internal/worktree/subshell.go`
- `internal/worktree/clean.go`
- `internal/worktree/verb_span.go`
- `internal/worktree/subshell_test.go`
- `internal/worktree/otel_seams_test.go`
- `internal/systemtest/otel_gate_test.go`
- `internal/systemtest/otel_verbs_test.go`
- `DATA_HANDLING.md`
- `tests/canary/data-handling-derivation/undocumented-passlist-var`

Reviewer disposition: approved on 2026-09-23. The reviewer approved the structure grant for the shift package with the other decisions under Further notes.

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
| `8-complete-the-shift-intent-on-exit.md` | `4-record-the-shift-boundaries.md`, `5-record-each-pass-under-the-shift.md`, `6-record-the-resolved-line.md`, `7-retain-the-shift-memory.md` | LE-C1 |
| `9-abandon-lease-less-stale-intent.md` | `4-record-the-shift-boundaries.md`, `8-complete-the-shift-intent-on-exit.md` | LE-C2 |
| `10-recover-crashed-shifts-by-lease-identity.md` | `7-retain-the-shift-memory.md`, `8-complete-the-shift-intent-on-exit.md`, `9-abandon-lease-less-stale-intent.md` | LE-C3 |
| `11-record-the-worktree-shell-session.md` | `4-record-the-shift-boundaries.md` | LE-D |
| `12-prove-the-shift-trace-through-the-built-binary.md` | `2-redact-and-version-every-record-line.md`, `5-record-each-pass-under-the-shift.md`, `6-record-the-resolved-line.md` | LE-D |
| `13-document-the-local-evidence-contract.md` | `3-rotate-and-retain-the-record.md`, `7-retain-the-shift-memory.md`, `10-recover-crashed-shifts-by-lease-identity.md`, `11-record-the-worktree-shell-session.md` | LE-D |

Ticket 4 creates `internal/shift/record.go` and `internal/shift/record_test.go`, and ticket 9 creates `internal/shift/recover.go` and `internal/shift/recover_test.go`. Each later ticket that writes one of those files marks it `(new)`, because the file is absent from the tree at spec time. Ticket 8 names every LE-B2 ticket in its `Blocked by:` line. Tickets 5, 6, and 7 also write `pass_test.go`, ticket 8 writes `record_test.go` after ticket 5, and tickets 5 and 8 both write `loop.go`. So a delegated frontier starts LE-C1 only after the LE-B2 checkpoint.

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
| Atomic append, rotation, and explicit retention are part of the repository contract. | LE13, LE14, LE18, LE19, LE57, LE92, LE93 |
| Records distinguish completed, failed, interrupted, recovered, and abandoned work; | LE27, LE28, LE37, LE64, LE69 |
| normal exit completes the matching intent | LE60, LE61 |
| and crash recovery requires the active lease identity. | LE63, LE64, LE69, LE70, LE71, LE97 |
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
- The SDK behavior was read in the module cache at `go.opentelemetry.io/otel/sdk@v1.46.0`. `NewProvider` builds the provider without `WithResource`, so `resource.Default` applies, and its `fromEnv` detector reads both variables. The 2026-09-23 probe confirmed both markers in a record line.
- This repository ignores `.bench-objective` and `.bench-notes.md` in `.gitignore`. So a crashed shift's worktree can read as clean, and an acquire can take its dead lease. The takeover claim of LE97 answers that race.

### Reader sweep and proof checklist

- Cited symbols: each symbol resolves at `ba9b8621`. The record package: `Encode`, `Begin`, `BeginIn`, `NewProvider`, `Writer.Append`, `gradeRecordPath`, `ReadSpans`, `ReadSelected`, `NewestLanding`, `DeclaredAttributes`, and `Registry`. The gate: `beginGateSpan`, `withGateSpanEnv`, `beginPhaseRecord`, `otelGateEnv`, and `RunAndRecordContext`. The shift: `finish`, `evidenceResult`, `preserveAndRecover`, `teardown`, `cleanupScratch`, `runAdapter`, `runGate`, `refactorPhase`, `exitPreserving`, and `shiftFault`. Others: `admissionpolicy.Live`, `intent.NewEntry`, `worktree.RetainAndLock`, `worktree.LeaseFile`, `worktree.ProbeLease`, `claimAt`, `joins.claimTakeoverGap`, `lifecyclepolicy.LeaseOwnerPID`, `bounds.ClassifyNoFollow`, `lines.Tiers`, `modelid.SafeToken`, `prepareProcessEnvironment`, and `subshellAt`.
- Import edges: the shift package gains `internal/otelrecord`, `internal/lines`, and `internal/bounds`, and none of the three imports the shift package. The session-inspect package gains `internal/shift`, and the shift package imports no session-inspect code. The shift package already imports `internal/worktree`.
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
  - `claimAt` has two call sites in `acquireAt` and one in `subshellAt`, and the identity claim becomes a fourth. Ticket 10 edits all four sites if the judgment becomes a parameter.
  - `admissionpolicy.Live` has two, `Snapshot` and `Compact`.
  - `finish` has every exit of `loop` and `exitPreserving`, and `subshellAt` has `Subshell` and its tests.
- Readers of the changed artifacts:
  - Production: the assessment collection checks the live path before `ReadSelected`, and the retro scaffold calls `NewestLanding`.
  - Tests that read the live path: `internal/gate/lane_record_test.go`, `internal/worktree/reset_plan_test.go`, `internal/roadmap/retro_scaffold_test.go`, `internal/assessment/collection_repair_test.go`, and the system helpers that glob `traces.jsonl`. None of these tests rotates, so each keeps reading the live segment.
  - Tests that assert every attribute is declared: `internal/systemtest/otel_gate_test.go`, `internal/systemtest/otel_hook_test.go`, `internal/systemtest/otel_worktree_test.go`, and `internal/systemtest/otel_verbs_test.go`. The encoder filter keeps each assertion true.
  - History: `decisions/software-factory/assets/ft303-cli-assessment.md` quotes the path and takes a named exclusion.
- Pinning consumers: `internal/conformance/bounds_policy_test.go` requires `internal/shift/loop.go` to keep its five iteration bound names. The implementation keeps them there.
- Copy survival: LE12 reds a surviving handoff name in the gate package. The sealed-name helper is the one format for the writer, the prune, and the readers.

### Structure and fence disposition

- `internal/shift/` holds 11 files and the lane caps a directory at 12. This spec adds `record.go`, `record_test.go`, `pass_test.go`, `recover.go`, and `recover_test.go`. Ticket 4 writes the grant line `internal/shift/ 15`, and ticket 5 raises it to 16.
- `internal/otelrecord/` gains one file, `retention_test.go`, and reaches the cap of 12.
- `cmd/bench/main.go`, `internal/gate/runner.go`, `internal/worktree/lifecycle.go`, and `internal/worktree/subshell.go` are over their line budgets. Tickets 1, 2, 6, 10, and 11 move code out or edit in place so that none of them grows. Ticket 11 moves the ignored-inventory function from `subshell.go` into `clean.go`.
- The command registry, its two tests, and the two conformance registry tests join the fence through the binding registry closure only. The `cmd/bench` and `internal/worktree` packages are bound packages, and the build expects no edit to those five files.
- The three canary fixture paths join the fence through the fixture closure only. Each pins a file that a ticket writes, and the build expects no edit to a fixture.
- `cmd/bench/`, `internal/gate/`, `internal/systemtest/`, and `internal/worktree/` have no directory headroom, so this spec adds no file there.

### Reviewer decisions

The reviewer closed these four decisions at sign-off on 2026-09-23:

1. The structure grant `internal/shift/ 15` is approved. The four new files hold one package's record, recovery, and their tests, and a split fragments that evidence.
2. The retention sizes are a 16 MiB segment, 8 sealed segments, and 64 memory files. The record then stays under 144 MiB for each repository.
3. This spec does not ship the `bank.ft71.local_event` producer. The producer is priced under Out of scope as its own capability, because FT88 closed with no producer and the release track stays NO-GO.
4. The retained notes text lives in a 0600 memory file beside the record, referenced by digest. A span event would put model prose in the record, and the FT274 declared set forbids a payload.

The reviewer closed three decisions at the LE-A review on 2026-09-23:

1. The listing of sealed names is the existence check, so the writer has no separate free-name step. LE94 grades a rotation that ignores the present names.
2. A reader can miss a segment that a rotation seals between its list and its live open. A Won't handle line records this window.
3. A planted file at the largest sequence has no successor, so the writer refuses the rotation. LE106 grades that refusal.

The reviewer closed one more decision on 2026-09-23: a recovery releases a clean crashed worktree and locks only a dirty one, as `preserveAndRecover` does. The Implementation decisions section records it.

The reviewer closed one decision at the LE-B2 build on 2026-09-23: the grant rises to `internal/shift/ 16`. The pass, line, and memory tests go in `pass_test.go`, so each test file stays under the line cap.

The reviewer closed one decision at the LE-B2 review on 2026-09-23. An adapter that a signal ends records `exited` with no `bench.adapter.exit` key, because it has no exit code.

### Flagged additions

- `service.version` on every line is the spec's reading of "release-evidence relationship". The source names the relationship but not its field.
- The declared-key filter in the encoder turns the FT274 review-owned rule into a mechanical one. It is the mechanism for "redacted".
- The schema read rules of LE8 to LE10 are the reader half of "versioned". The source names a versioned schema but no read rule.
- The handoff move of LE11 and LE12 is a prefactor that tickets 6 and 12 consume. The source does not name it.
- The segment refusals of LE20 and LE21 extend the writer's grade to the readers. The source names rotation but not the read path.
- The sealed-name sequence and the prune by sequence of LE92 to LE94 are mechanisms that the source does not name.
- The `line.resolve` span records the resolved model at its one resolver. The shift itself knows only the declared tier.
- The worktree shell is the "session" of "shift/session boundaries", because the source's C-05 names the subshell. Harness sessions take a Won't handle line.
- The recovery triggers, the printed line, the takeover claim, the resume rule, and the memory states are mechanisms that the source did not name.
- The release-when-clean rule of LE68 is a mechanism that the source did not name. The reviewer closed it on 2026-09-23, and it follows `preserveAndRecover`.
- The skip for a live other owner (LE103) and the resume arm for another identity (LE104) are guards that the source did not name. They keep a verdict out from under a live claimant or a new owner.
- The memory write before the claim (LE105) is an ordering that the source did not name. It keeps the notes when an acquire wins the tree.
- The retention sizes are proposals that the reviewer owns.

### Completion plan

```bench-completion-plan
{"version":1,"chunks":[{"id":"LE-A","tickets":["1-move-the-trace-handoff-into-the-record.md","2-redact-and-version-every-record-line.md","3-rotate-and-retain-the-record.md"],"verification":[{"id":"otelrecord","command":"bench test --package ./internal/otelrecord"},{"id":"gate","command":"bench test --package ./internal/gate"},{"id":"cmd","command":"bench test --package ./cmd/bench"},{"id":"system","command":"bench test --check system"}]},{"id":"LE-B1","tickets":["4-record-the-shift-boundaries.md"],"verification":[{"id":"shift","command":"bench test --package ./internal/shift"},{"id":"kit-compliance","command":"bench test --check kit-compliance"}]},{"id":"LE-B2","tickets":["5-record-each-pass-under-the-shift.md","6-record-the-resolved-line.md","7-retain-the-shift-memory.md"],"verification":[{"id":"shift","command":"bench test --package ./internal/shift"},{"id":"otelrecord","command":"bench test --package ./internal/otelrecord"},{"id":"cmd","command":"bench test --package ./cmd/bench"},{"id":"kit-compliance","command":"bench test --check kit-compliance"}]},{"id":"LE-C1","tickets":["8-complete-the-shift-intent-on-exit.md"],"verification":[{"id":"intent","command":"bench test --package ./internal/intent/..."},{"id":"shift","command":"bench test --package ./internal/shift"}]},{"id":"LE-C2","tickets":["9-abandon-lease-less-stale-intent.md"],"verification":[{"id":"shift","command":"bench test --package ./internal/shift"},{"id":"intent","command":"bench test --package ./internal/intent/..."},{"id":"sessioninspect","command":"bench test --package ./internal/sessioninspect"}]},{"id":"LE-C3","tickets":["10-recover-crashed-shifts-by-lease-identity.md"],"verification":[{"id":"shift","command":"bench test --package ./internal/shift"},{"id":"claim","command":"bench test --package ./internal/worktree --run TestClaimRecordedLease"}]},{"id":"LE-D","tickets":["11-record-the-worktree-shell-session.md","12-prove-the-shift-trace-through-the-built-binary.md","13-document-the-local-evidence-contract.md"],"verification":[{"id":"subshell","command":"bench test --package ./internal/worktree --run TestSubshell"},{"id":"seams","command":"bench test --package ./internal/worktree --run TestWorktreeSeamsMatchTheRegistry"},{"id":"system","command":"bench test --check system"},{"id":"prose","command":"bench gate-prose . -- DATA_HANDLING.md"}]}],"final_verification":[{"id":"coverage","command":"bench coverage --check specs/local-shift-evidence/spec.md"},{"id":"otelrecord","command":"bench test --package ./internal/otelrecord"},{"id":"shift","command":"bench test --package ./internal/shift"},{"id":"kit-compliance","command":"bench test --check kit-compliance"},{"id":"system","command":"bench test --check system"}]}
```
