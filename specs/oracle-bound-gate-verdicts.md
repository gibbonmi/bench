# Oracle-bound gate verdicts

Status: implemented

## Problem

Bench currently authorizes `bench commit` from a three-field gate-cache line keyed
only to the working-tree hash. The cache does not identify the gate command, its
launcher or tool content, ignored inputs, inherited environment, schema, or a real
freshness policy. A byte-identical working tree can therefore reuse green from a
different oracle.

The current writer is also not part of the gate's success contract. It overwrites the
cache directly, ignores write failure, and records only after the gate finishes. An
older green can survive a red or otherwise unrecordable run. Gate executions have no
single-owner protocol, so overlapping writers can publish results out of order. Stop
runs the standalone gate path and then writes the verdict again.

Disposable throwaway-repository probes against the real wrapper reproduced the
defects in the current tree:

- changing `BENCH_GATE` from exit 0 to exit 23 after a green gate let `bench commit`
  exit 0 and create `oracle-probe` without running the new oracle;
- changing an inherited `ORACLE_INPUT` from `green` to `red`, without changing the
  command or tree, likewise let commit create `env-probe`;
- changing an ignored gate input from `green` to `red` let commit create
  `ignored-probe`;
- planting a matching green dated `2000-01-01` let a red gate command be skipped and
  commit create `freshness-probe`;
- making the cache and Git directory unwritable left the planted green intact, still
  ran the red gate, and returned its exit 31 instead of failing before execution;
- while one gate was blocked, a second gate exited 0 and ran concurrently;
- a fake standalone wrapper that wrote its own cache record was overwritten by Stop's
  second writer; and
- a planted schema-1 pending record was not classified as interrupted-pending by
  status or roadmap context, while the dashboard page exposed unrelated raw text
  containing the word `pending`.

Oracle A green can therefore authorize an action under oracle B, and an unrecordable
result can leave the prior authorization in place. That violates the gate-is-the-
oracle invariant at the exact point where Bench mechanizes commit-on-green.

## Solution

Make the gate module the single deep owner of oracle resolution, closed-subject
construction, strict verdict inspection, exclusive execution, and durable verdict
replacement. A read-only inspection exposes the only reusable-green predicate. An
execution snapshots the exact tree and oracle, acquires the repository's gate lock,
durably replaces any older record with pending, runs only the captured oracle, rejects
subject drift, and durably publishes the final green or red.

Introduce strict repository-owned gate input configuration. Gate processes receive
only `PATH` plus explicitly declared environment variables. The canonical subject
digest binds the resolution kind and exact command, repository-root working directory,
tree hash, launcher and declared tool identities, declared ignored/local inputs,
passed environment values, schema policy, and freshness policy. Missing, malformed,
remote, unreadable, escaped, or otherwise incomplete inputs leave the subject open:
the gate still runs and its latest result remains inspectable, but no cached green can
authorize reuse.

Replace the legacy cache with strict schema-1 JSON in the Git directory. The sole
writer uses a same-directory mode-0600 temporary file, file sync, atomic rename, and
directory sync for pending and final records. One nonblocking OS-held lock spans the
entire pending-through-final transaction. A crash leaves non-reusable pending evidence;
a later owner may replace it and run.

Route standalone gate, commit, shift, Stop, status, the dashboard page, and roadmap
context through that owner. Commit alone may reuse a matching fresh green. The other
gate-running actions execute; the read-only consumers inspect without running or
repairing. Runtime contracts and two behavior-owned canaries prove both independent
invariants: oracle identity is compared, and old green is durably invalidated before
execution.

## User stories

1. As a repository maintainer, I want a cached green to identify one closed working-
   tree-and-oracle subject and one freshness policy, so that changing any input that
   can change the gate verdict forces the real gate to run again. Line:
   `gpt-5.6-terra` / medium. This is gate authorization logic at a settled seam, and
   the project profile routes gate and conformance changes to the mid line because a
   false green is the kit's highest-cost defect.

2. As a developer invoking gate, commit, shift, or an armed Stop, I want one durable
   gate execution owner to invalidate older authorization before running and publish
   a result only after durable finalization, so that crashes, contention, drift, and
   persistence failures fail closed without staging or losing work. Line:
   `gpt-5.6-terra` / medium. The state machine is fully decided and observable, but
   cross-process ordering and persistence remain oracle code rather than cheap shell
   plumbing.

3. As a reviewer reading status, the dashboard page, or roadmap context, I want every
   surface to project the same typed verdict inspection without running or reparsing
   the gate, so that absent, reusable, red, stale, pending, invalid, and unavailable
   evidence cannot disagree. Line: `gpt-5.6-luna` / low. The deep owner fixes the
   semantics, leaving a mechanical consumer migration with literal runtime coverage.

4. As a kit maintainer, I want real-wrapper contracts, deterministic injected-fault
   tests, hostile-input coverage, and biting canaries for verdict authorization and
   invalidation, so that weakening either invariant makes the gate red with its own
   attributed message. Line: `gpt-5.6-terra` / medium. This story authors the oracle's
   regression layer and canary tripwire, which the project profile keeps on the mid
   line even though every assertion is executable.

The four stories land as one atomic implementation. The new schema, writer, reuse
predicate, and consumers cannot form independently green authorization slices without
a temporary second source of cache knowledge. The implementation author therefore
runs at the highest story line: `gpt-5.6-terra` / medium.

## Implementation decisions

### One deep gate owner

The gate module owns resolution, subject construction, input configuration, cache
codec, freshness, inspection, locking, execution, and durable replacement. The Git
module continues to own repository-root, absolute-Git-directory, and working-tree-hash
resolution. Action and presentation modules consume typed gate results; none parses
cache bytes, reconstructs a subject, or decides freshness.

The production interface is concrete and small:

- `Inspect(root)` is read-only and returns the cache state, ready status, cached and
  current tree identities, recorded time, reason, and the sole `ReusableGreen`
  predicate.
- `Execute(ctx, root, stdout, stderr)` returns the original gate exit, the requested
  action exit, and the resulting inspection. It is the only operation that changes
  verdict state.

Resolver, manifest collectors, clock, filesystem, and lock dependencies stay behind
an unexported engine used by tests. There is one production implementation, so no
public Resolver, Store, Oracle, clock, filesystem, or lock interface is introduced.
The deletion test for this module is decisive: without it, command resolution,
fingerprinting, parsing, freshness, durability, and authorization would reappear in
every caller.

`Inspect` classifies exactly these states:

- `absent` — no cache record;
- `ready` — a valid schema-1 green or red record;
- `pending` — a valid schema-1 in-progress record, projected as `locked-pending` when
  the OS lock has a live owner and `interrupted-pending` when it does not;
- `invalid` — present bytes or metadata that violate the strict record contract; and
- `unavailable` — Git directory, tree, cache metadata, or required read could not be
  inspected safely.

`ReusableGreen` is true only for `ready` green when the current subject is closed, the
cached and current trees and oracle digests match, and the record is within the fixed
freshness window. No caller may rebuild that predicate from fields.

### Closed oracle subjects and gate input configuration

The repository-owned `.bench/gate-inputs.json` is strict JSON with this complete
semantic shape:

```json
{
  "schema": 1,
  "closure": "local",
  "environment": ["NAME"],
  "paths": ["repo-relative/path"],
  "tools": ["executable-name-or-path"]
}
```

All five fields are required. `closure` is either `local` or `remote`; `remote`
explicitly makes every result non-reusable. The three arrays may be empty, are
deduplicated and sorted for canonical hashing, and do not derive execution policy
from source order. Environment names use the portable shell-name grammar. Declared
paths are slash-formed, repository-relative, traversal-free inputs; absolute paths,
escape through a symlink, missing entries, unsupported file types, and control bytes
make the subject open. Declared tools may be command names resolved through the exact
passed `PATH`, repository-relative executables, or absolute executables. A missing,
non-executable, cyclic, or unresolvable tool makes the subject open.

The manifest is capped at 16 KiB and rejects duplicate JSON names, unknown fields,
trailing tokens, unsupported schemas, wrong types, invalid UTF-8/control bytes, and
an absent final object. An absent, empty, or malformed manifest never blocks gate
execution; it supplies a typed non-reuse reason. Readers do not repair it.

Gate execution receives only inherited `PATH` plus the variables named in
`environment`; wrapper-routing variables and every other ambient variable are absent.
The exact values actually passed join the subject digest. A declared but absent
variable makes the subject open and remains absent in the child. Environment names or
values never enter the verdict cache or diagnostics.

Every subject is a SHA-256 digest over a versioned, length-framed canonical envelope:

1. verdict schema, fingerprint-policy version, and freshness-policy version;
2. resolution kind, exact executable/argv or exact single `BENCH_GATE` shell string,
   and canonical repository-root working directory;
3. working-tree hash;
4. the complete recognized launcher closure, including path, symlink disposition,
   executable mode, and streamed content identity;
5. the canonical manifest identity and sorted environment, path, and tool entry
   identities; and
6. the exact passed `PATH` and declared environment values.

Launcher closure follows direct executable paths, shebang interpreters, and
`/usr/bin/env` selection. The known collectors are enumerated by kind: gate script
uses its shebang closure; `BENCH_GATE` uses Bash; pnpm uses Bash and pnpm; npm uses
Bash, npm, and npm's interpreter closure; Python uses Bash, mypy, pytest, ruff, and
their interpreter closures; Cargo uses Bash, cargo, rustc/toolchain, and clippy's
resolved executable closure. The manifest names any further executable reached by an
opaque script or package command.

Repository paths may name files, symlinks, or directories. Directory identity walks
entries in sorted slash-path order without following symlinks and hashes relative
name, type, mode, link target, and file content. Subject collection streams content
and is capped at 100,000 entries and 1 GiB aggregate bytes across declared paths and
launcher/tool content. Crossing either cap makes the subject open and runs the gate;
it never substitutes metadata for content identity. Symlink cycles, more than 64
link hops, or a declared repository path escaping the root have the same posture.

The fixed resolver precedence remains gate script, `BENCH_GATE`, pnpm, npm, Python,
Cargo, then no gate. Selection-file bytes are already in the tree when unignored; an
ignored selection file must be declared. The configuration is an assertion by the
project that the listed local dependencies complete an otherwise opaque oracle.
Bench verifies every declared entry but cannot infer an omitted shell dependency or
remote service.

### Strict versioned verdict records

The cache is one regular, non-symlink mode-0600 file in the repository's absolute Git
directory. It accepts at most 16,384 bytes including optional final newline; a
16,385th byte is invalid. A ready record contains exactly:

- integer `schema` equal to 1;
- `state` equal to `ready`;
- `status` equal to `green` or `red`;
- lowercase hexadecimal `tree` and 64-character lowercase hexadecimal `oracle`;
- UTC RFC3339 `recorded_at` with whole-second precision.

A pending record contains exactly schema 1, state `pending`, the pre-run `tree` and
`oracle`, whole-second UTC `started_at`, and positive integer diagnostic `owner_pid`.
PID and time are evidence for a human, never lock ownership authority.

Unknown or duplicate fields, trailing JSON, unsupported schema, bad enum, missing or
wrong-type fields, malformed hashes or timestamps, future timestamps, legacy text,
wrong file type, symlinks, broader permissions, truncation, and oversize bytes are
invalid and non-reusable. Readers classify but never rewrite. They never echo hostile
bytes. Command strings, manifest entries, environment names or values, declared path
content, tool output, and gate output never enter the record.

A matching green is fresh strictly before ten minutes after `recorded_at`. At ten
minutes or later it is stale. Missing, malformed, or future time is invalid. There is
no project override in schema 1. Changing the freshness-policy version changes the
subject digest.

### Exclusive, durable execution

One nonblocking OS-held exclusive lock per absolute Git directory is acquired before
pending replacement and held through final directory sync. Lock contention fails
immediately with action exit 1. Different Git directories remain independent. The OS
lock is the sole live-owner authority; PID, timestamp, lock-file bytes, and age never
authorize stealing or waiting.

Execution order is fixed:

1. resolve and snapshot the current subject, including its closed/open reason;
2. acquire the repository gate lock;
3. re-resolve the subject under the lock and refuse plan drift;
4. durably install pending through a same-directory mode-0600 temporary file, file
   sync, atomic rename, and directory sync;
5. run that exact captured oracle from the repository root with only the gate
   passlist, streaming its output;
6. rebuild the subject and reject tree, resolution, command, manifest, environment,
   declared input, tool, or launcher drift; and
7. durably install ready green or red with the same write protocol before releasing
   the lock.

Failure to establish and sync pending prevents the oracle from starting. Once pending
is durable, a crash or cancellation can never reveal the older green. A process crash
releases the OS lock and leaves interrupted-pending; the next owner may replace it and
run normally. Cancellation kills the gate process group, leaves durable pending, and
returns the interruption result. A final write, file sync, rename, directory sync,
subject recheck, or other operational failure returns action exit 1 and never claims
green. If rename succeeds but directory sync fails, the new bytes may be visible, but
the previously synced pending remains the crash-recovery authorization floor.

The result preserves both exits. Durable green returns gate 0/action 0. Durable red
returns the gate's original nonzero exit as both the gate and action verdict. No oracle
returns 3 and records nothing. Cancellation returns the process's interruption result.
Lock, subject, persistence, and other operational failures return action exit 1 while
retaining the original gate exit and streamed output as non-secret diagnostics.

Open subjects still execute and durably record their latest result. Inspection names
the closure reason, and `ReusableGreen` remains false even for fresh matching green.
Re-running an interrupted, invalid, expired, open, or red state always executes.

### Thin action and projection consumers

Standalone `bench gate` always executes. `bench commit` performs its existing
block-check and stage plan, inspects, and executes unless `ReusableGreen` is true; any
nonzero action result refuses before spec flip or staging. Shift executes on every
iteration and preserves its work on any red or operational result. Armed Stop invokes
the standalone wrapper path once, translates nonzero to its blocking response, and
does not record a second verdict. Unarmed and already-active Stop remain non-writing.

Status, the dashboard page, and `bench roadmap --context` inspect only. They preserve
their existing presentation roles and derive them from the same typed state:

- absent and reusable green remain non-signals in the ambient dashboard;
- red remains the immediate fix-before-commit signal;
- stale/expired and invalid evidence request a gate rerun;
- locked-pending reports a live gate owner and does not invite another run;
- interrupted-pending requests a new gate run;
- unavailable reports that gate state cannot be inspected safely; and
- the dashboard page and roadmap context expose the same stable state names and
  fields while retaining their HTML and TOON contracts.

No projection reads cache bytes, computes a tree, tests freshness, probes the lock, or
runs the gate independently. Raw record content and non-sensitive reason details are
not rendered where control bytes or secrets could escape.

### Atomic sequence and structure

Implementation proceeds in four ordered story slices on one branch: subject,
configuration, strict codec, and `Inspect`; lock plus durable `Execute`; consumer
migration and deletion of the old parser and second writer; then full runtime,
injected-fault, hostile-input, and canary proof. Green is promised at the final
integrate-and-verify step because retaining the old schema or writer beside the new
authorization path would violate the one-source standard.

The gate package currently has directory headroom. New responsibilities split into
cohesive files inside that package rather than new packages or one oversized source.
The already-granted runtime-contract directory remains one shared fixture package;
tests join the existing command-family files or replace obsolete expectations rather
than creating a second harness.

## Testing decisions

- The highest seam is the real `bench` wrapper in throwaway Git repositories. Tests
  assert exit code, stdout/stderr, run count, cache bytes and metadata, HEAD/index,
  retained shift work, and consumer output as applicable.
- The only lower seam is the unexported gate engine with injected clock, filesystem,
  and lock operations. It exists solely for exact byte limits, future-clock
  evaluation, lock-acquisition error, file sync, rename, after-rename, and directory-
  sync failures that cannot be induced portably through the CLI.
- Runtime contracts extend the existing gate, commit, status, shift, and roadmap
  context families. A second cache parser in fixtures is not evidence; tests either
  drive typed inspection or compare strict literal bytes where the cache encoding is
  the contract.
- Every new contract calls `NoteContractFailure` with its unique named message and
  continues the family. The gate runs the real path and aggregates failures.
- Two `behavior-owned` canaries sabotage independent invariants. The oracle-binding
  canary trusts tree/status while ignoring the oracle digest and must trigger
  `oracle-bound gate verdict contract failed`. The invalidation canary runs before
  pending and permits old green after final-write failure and must trigger
  `fail-closed gate verdict persistence contract failed`.
- The feature must pass `.bench/gate.sh`; focused red/green work uses the relevant Go
  package or runtime-contract command before the full gate.

### Seam diagram

Real CLI/action/projection seam:

    trigger: developer, shift loop, Stop hook, SessionStart, or roadmap reader
        │
        ▼
    command + repo state ──▶ [ bench wrapper ──▶ gate owner ] ──▶ exit/output + cache/action state
    cache + repo state   ──▶ [ Inspect ──▶ typed projections     ] ──▶ status/page/context
                                  ◀ tests attach here: invoke real shipped wrapper in throwaway repos

Private deterministic-fault seam:

    trigger: internal gate table test
        │
        ▼
    fixed clock + ordered FS/lock outcomes ──▶ [ unexported gate engine ] ──▶ typed inspection/result
                                                       ◀ tests attach here: inject one impossible-to-portably-induce fault

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | A complete local manifest plus an unchanged closed subject lets gate-then-commit run the gate once; absent, empty, malformed, wrong-schema, remote, missing-variable, missing-path, missing-tool, escaped-symlink, over-limit, and control-byte manifests run on every request with a non-reuse reason. | real CLI/action/projection seam | Observed red: the disposable changed-command, changed-environment, and changed-ignored-input probes all let commit exit 0 without any manifest or closure proof. | The positive run-count control prevents an always-rerun implementation, while every enumerated closure failure prevents best-effort identity from authorizing commit. |
| 1 | Changing the exact `BENCH_GATE` string, gate script/interpreter or declared tool content/mode/target, declared environment value, declared file/directory content, ignored input, `PATH`, repository root, or auto-detected kind reruns; the new red refuses commit. | real CLI/action/projection seam | Observed red: disposable command, `ORACLE_INPUT`, and ignored-file mutations each kept the tree constant, skipped the new red oracle, and created a commit. | Each mutation changes one subject component while holding the tree fixed, so tree-only reuse and partial collectors cannot pass. |
| 1 | Gate execution receives only `PATH` and the enumerated manifest variables; `BENCH_KIT`, `BENCH_WRAPPER`, `HOME`, `CI`, locale variables, and arbitrary inherited names are absent unless individually declared. | real CLI/action/projection seam | Observed red: the `ORACLE_INPUT` probe proved an undeclared inherited variable currently reaches the gate and can change its result while commit reuses green. | A gate that prints presence bits for the enumerated names catches both ambient leakage and a passlist that accidentally strips declared input. |
| 1 | Gate-script, `BENCH_GATE`, pnpm, npm, Python, and Cargo resolution preserve the existing precedence and bind the enumerated launcher/tool closures; no-gate remains exit 3 with no record. | real CLI/action/projection seam | Observed red: the changed-command probe committed under a different resolved command, while the existing resolution-order contract already supplies the positive precedence control. | Combining mutation with the established positive matrix rejects both a reordered resolver and a digest that omits kind, command, or launcher identity. |
| 1 | Strict inspection distinguishes absent, ready green/red, locked pending, interrupted pending, invalid, and unavailable for zero-byte, no-newline, trailing-token, duplicate/unknown-field, wrong-type/schema/enum/hash/time, legacy, truncated, exactly-16-KiB, over-16-KiB, wrong-mode, symlink, directory, and unreadable records. | real CLI/action/projection seam plus private deterministic-fault seam | Observed red: a planted schema-1 pending record was not classified as interrupted-pending by status or roadmap context, and the expired legacy probe was accepted for commit. | The paired state and metadata matrix forces strict decoding and file checks; exact boundary cases catch prefix readers and off-by-one limits. |
| 1 | A ready matching green reuses only before ten minutes; exactly ten minutes, expired, future, malformed, or policy-version-mismatched evidence reruns, while readers never rewrite it. | real CLI/action/projection seam plus private deterministic-fault seam | Observed red: a matching green dated `2000-01-01` skipped a red gate and created `freshness-probe`. | Fixed-clock boundary pairs reject the current timestamp-ignored parser and prevent a project override or inclusive ten-minute edge. |
| 1 | Cache and diagnostics contain no command string, environment name/value, manifest path, input content, tool output, or gate output, including secrets and unsafe controls. | real CLI/action/projection seam | Observed red: the current cache cannot encode an oracle at all, so the changed-environment probe reused green without a secret-safe aggregate identity. | Sentinel secrets across every enumerated source must affect the digest yet remain absent from cache and output, catching both omission and raw-value persistence. |
| 1 | Deep cwd, spaces and glob characters in repository and declared paths, no trailing newline, symlink invocation, missing global `bench`, executable mode changes, symlink chains, and external launcher targets reach the same subject rules without raw control-byte output. | real CLI/action/projection seam | Observed red: the ignored-input probe changed a repository-local dependency without changing the subject; existing deep-cwd and symlink-wrapper contracts are the positive controls. | Driving the shipped wrapper with hostile paths and launch forms exposes quoting, root-resolution, mode, and target-identity bugs hidden below the CLI. |
| 2 | Before a marker inside the gate can appear, any prior green is durably replaced by pending using temp mode 0600, file sync, rename, and directory sync; failure at any pre-run step returns action exit 1 and the marker stays absent. | real CLI/action/projection seam plus private deterministic-fault seam | Observed red: with an unwritable cache/Git directory, the red gate still ran, exited 31, and left the planted green intact. | Marker ordering proves execution cannot race ahead of invalidation; injected calls distinguish file-write, file-sync, rename, and directory-sync failures. |
| 2 | Durable green returns 0, durable red preserves the gate's original nonzero exit, and green/red final-write, file-sync, rename, after-rename directory-sync, or subject-recheck failure returns action exit 1 without exposing older green. | real CLI/action/projection seam plus private deterministic-fault seam | Observed red: the unwritable-cache probe retained old green and returned the gate's exit instead of an operational persistence failure. | Rereading bytes and inspection after each ordered fault rejects ignored write errors, asymmetric green success, and false finalization claims. |
| 2 | Tree, command, manifest, environment, declared path/tool, launcher, or auto-detect drift during a run returns non-success and leaves pending; cancellation kills the process group, returns interruption, and leaves pending. | real CLI/action/projection seam | Observed red: the current runner computes its cache tree only after execution and has no pre/post oracle comparison; the disposable changed-input probes show those components are not captured at all. | Blocking gates mutate each enumerated component between start and release, so a post-run-only snapshot or tree-only drift check fails. |
| 2 | While one gate owns the Git-directory lock, a second standalone gate, commit, shift, and armed Stop each fail immediately without running, staging, or discarding work; different Git directories still run independently. | real CLI/action/projection seam | Observed red: while one disposable gate blocked, a second gate exited 0 and created its execution marker. | One blocking owner plus four consumers catches missing, per-command, or waiting locks; a parallel second-repository control prevents a process-global lock. |
| 2 | Killing the owner leaves interrupted-pending; the released lock lets the next run replace pending and finish, while PID and age never block or authorize recovery. | real CLI/action/projection seam | Observed red: the concurrent-gate probe demonstrated there is no live-owner authority or pending crash evidence in the current implementation. | A killed child with both live-looking and dead-looking PIDs rejects PID/timeout ownership and proves recovery is tied only to the OS lock. |
| 2 | Commit reuses only `ReusableGreen`; every other state executes, and any red or operational result leaves HEAD and index unchanged and never flips or stages the spec. | real CLI/action/projection seam | Observed red: command, environment, ignored-input, and expired-legacy probes all exited 0 and moved HEAD under a red current oracle. | HEAD, index, run-count, and spec-byte assertions attach directly to the authorization consumer and reject a second predicate in commit. |
| 2 | Shift executes every iteration and preserves its worktree, branch, intent, and uncommitted changes on red, lock, persistence, drift, or cancellation failure. | real CLI/action/projection seam | Observed red: the common unwritable-cache probe returned the gate result after failed recording, so the current shared runner cannot give shift an operational failure to preserve against. | Injecting each common Execute result through a real one-iteration shift proves the loop does not reinterpret or clean up a failed oracle transaction. |
| 2 | Armed Stop invokes the standalone wrapper once, blocks on every nonzero gate/action result, and never records again; unarmed and active Stop do not inspect, run, or write. | real CLI/action/projection seam | Observed red: a fake wrapper wrote `wrapper-owned`, returned 0, and current Stop overwrote it with its own green record. | A wrapper-owned sentinel plus run/write counts catches the second writer without reimplementing the real gate; red and operational cases prove blocking translation. |
| 2 | Repeated interrupted or invalid runs are idempotent at the state-machine level, never resurrect older green, and never leave same-directory temporary files as reusable evidence. | real CLI/action/projection seam plus private deterministic-fault seam | Observed red: current execution has no pending state or lock, and the unwritable-cache probe left the older green as the only durable evidence. | Repeating every partial boundary tests recovery from durable truth instead of assuming a destructive write step ran exactly once. |
| 3 | Status, the dashboard page, and roadmap context consume one inspection and consistently project absent, reusable-green, red, stale, locked-pending, interrupted-pending, invalid, and unavailable without parsing cache bytes. | real CLI/action/projection seam | Observed red: the planted pending probe produced no pending state in status or roadmap context while the dashboard page matched unrelated raw text. | A literal three-surface matrix catches divergent parsers, raw-byte leakage, and projection vocabulary drift for every enumerated state. |
| 3 | Read-only consumers never run the gate, acquire a long-lived execution lock, rewrite invalid/legacy/pending bytes, repair permissions, or change timestamps. | real CLI/action/projection seam | Observed red: the current consumers accept legacy cache knowledge owned outside the gate module, and the expired legacy probe showed that knowledge can authorize an action. | Gate-run markers plus before/after cache metadata prove inspection is observational and that all policy moved behind the deep owner. |
| 3 | Ambient status preserves its show-only-on-signal and severity roles, while the dashboard page stays valid self-contained HTML and roadmap context stays AXI TOON with stable typed gate state. | real CLI/action/projection seam | Observed red: the pending projection probe failed to expose one consistent typed state; existing literal status, dashboard, and roadmap contracts provide positive format controls. | Format controls prevent consumer migration from fixing semantics by redesigning or bypassing the established public surfaces. |
| 4 | Private table tests force future-clock, exact byte-bound, lock-acquisition, temporary-write, file-sync, rename, after-rename, and directory-sync outcomes and assert ordered calls plus the resulting typed state. | private deterministic-fault seam | Observed red: the unwritable-cache probe showed the production runner ignores persistence failure and exposes no injectable ordered durability seam. | Each single-fault call trace distinguishes operations that a broad CLI permission failure would collapse and makes reordered durability visible. |
| 4 | The oracle-binding behavior-owned canary substitutes a wrapper that trusts matching tree/status while ignoring the oracle digest and makes the gate red with `oracle-bound gate verdict contract failed`. | real CLI/action/projection seam through behavior-owned canary | Observed red: `rg --files tests/canary/behavior-owned` found no `gate-verdict-oracle-binding-bypassed` fixture. | The planted cheapest wrong implementation passes tree-only tests but must fail the real runtime contract with the attributed message. |
| 4 | The durable-invalidation behavior-owned canary substitutes a wrapper that runs before pending and permits old green after final-write failure and makes the gate red with `fail-closed gate verdict persistence contract failed`. | real CLI/action/projection seam through behavior-owned canary | Observed red: `rg --files tests/canary/behavior-owned` found no `gate-verdict-invalidation-bypassed` fixture. | This independent needle catches removal or weakening of pre-run invalidation even if oracle comparison remains correct. |
| 4 | Focused gate, commit, shift, Stop, status, dashboard, and roadmap-context families all exercise the real wrapper; the two canaries are registered exactly once as `behavior-owned`, and the empty baseline still rejects vacuous EXPECT text. | real CLI/action/projection seam through runtime contracts and canary registry | Observed red: both named fixture searches exit 1, while current runtime probes demonstrate the contract families do not yet reject the planted behaviors. | Registry ownership, real-path execution, and baseline-vacuity controls prevent a fixture from passing through wrong routing, duplicate ownership, or an expectation satisfied by ambient noise. |

**Degenerate-implementation check.** The cheapest wrong implementation keeps the old
three-field cache and adds only command text, hashes the full inherited environment,
writes pending without syncing, waits or races on a PID file, lets each consumer parse
JSON, or simply reruns every gate forever. The positive exact-reuse row rejects
always-rerun. Command/environment/ignored/tool mutations reject tree-only and partial
subjects. Pre-run marker and injected sync rows reject non-durable invalidation. The
blocking-owner matrix rejects PID, per-command, and waiting locks. The three-surface
projection row rejects copied parsers. The two independent canaries ensure neither an
always-green authorization check nor an always-green persistence check survives.

### Edge inventory

- **Error paths:** covered by manifest closure, strict inspection, pre-run
  invalidation, finalization, lock, drift, action-consumer, projection, and injected-
  fault rows. Git-root/tree lookup; manifest/cache open, read, decode, metadata,
  temporary creation, write, file sync, rename, directory sync; lock open/acquire;
  executable resolution/hash/start; and post-run inspection each have an attributable
  fail-closed result.
- **Empty and absent input:** covered by manifest, cache-state, no-oracle, and
  projection rows. Absent versus zero-byte manifest/cache, empty arrays, absent
  declared variables/paths/tools, empty gate output, no gate, and a definitive empty
  projection are distinct cases.
- **Boundary values:** covered by strict inspection, freshness, subject collection,
  and concurrency rows. These include 16,384/16,385 cache and manifest bytes; just
  before/exactly at/after ten minutes; zero/one/multiple declarations; 64/65 symlink
  hops; 100,000/100,001 entries; just below/at/above 1 GiB; one/multiple repositories;
  and gate exits 0, 1, 2, 3, 127, signal interruption, and a larger arbitrary nonzero.
- **Malformed input:** covered by manifest and cache matrices. Invalid UTF-8, decoded
  control bytes, no final newline, trailing tokens, duplicate/unknown fields, wrong
  types, unsupported schemas, bad enums/hashes/times/names/paths, legacy text,
  truncation, oversize, wrong mode/type, symlink, and partial writes never authorize
  reuse or leak raw bytes.
- **Interrupted and partial state:** covered by pre-run, finalization, drift,
  cancellation, crash recovery, and idempotency rows. Every interruption before and
  after temporary write, file sync, rename, directory sync, child start/exit, drift
  check, and final replacement leaves the last durable non-reusable truth.
- **Re-run and concurrency:** covered by the owner, crash, and idempotency rows.
  Same-Git-dir actions fail immediately behind one live owner; a released crash lock
  permits recovery; separate Git dirs proceed; repeated partial-state runs never
  resurrect prior green.
- **Hostile environment:** covered by passlist, subject, and hostile-wrapper rows.
  Spaces and glob characters, deep cwd, symlinked launch, missing global tools,
  declared/missing variables, hostile PATH, ignored files, executable modes, shebangs,
  external launcher targets, and remote-marked or omitted dependencies all have an
  explicit posture.
- **External remote equivalence:** **Won't handle** — a project may mark a gate remote
  and get run-every-time behavior, but Bench cannot prove two network or service
  responses equivalent from repository-controlled evidence.
- **Malicious local repository owner:** **Won't handle** — a user who can replace the
  binary, Git directory, lock substrate, and cache can forge local evidence; this
  contract closes accidental and ordinary cross-process drift, not a hostile machine
  owner.

## Out of scope

- **Changing gate-resolution precedence** — a separate oracle-selection policy; this
  feature fingerprints and preserves the existing order. Estimated later cost:
  `4 edits, 3 gate runs`.
- **Replacing or folding the push-time gate pin into verdict identity** — a separate
  reviewer-trust capability outside the writable-tree cache; the existing pin remains
  independently enforced. Estimated later cost: `6 edits, 4 gate runs`.
- **Remote-state equivalence or reusable remote attestations** — a separate trust and
  protocol capability; schema 1 marks such gates non-reusable. Estimated later cost:
  `10 edits, 6 gate runs`.
- **FT82 release-preflight evidence** — a separate release decision capability that
  consumes gate truth but does not define verdict identity or durability. Estimated
  later cost: `10 edits, 5 gate runs`.
- **FT88 full agent and process environment minimization** — a separate platform-wide
  data-inventory capability; this feature narrows only the gate environment required
  for an honest subject. Estimated later cost: `12 edits, 6 gate runs`.
- **Redesigning ambient status, the dashboard page, or roadmap context** — separate
  presentation capabilities; this feature only replaces their gate-state source and
  adds the decided typed projections. Estimated later cost: `8 edits, 4 gate runs`.
