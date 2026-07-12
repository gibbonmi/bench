## Destination

Make a cached gate verdict reusable only for the exact working tree and oracle
that produced it, with durable fail-closed replacement so an older green can
never authorize work after a different or unrecordable result.

## #1: What is the reuse ceiling for an oracle with opaque dependencies?

Blocked by: none
Type: Grill

### Question

Decide whether reuse requires a closed, declared oracle subject. A
`BENCH_GATE` shell command, an auto-detected package script, or an executable
reached through `PATH` can depend on environment, ignored files, installed
tools, or remote state that Bench cannot infer exhaustively. The choice is
between disabling reuse when the subject is not closed, allowing projects to
declare extra fingerprint inputs, or accepting a documented best-effort
ceiling.

### Answer

Reuse requires a closed, declared oracle subject. Bench may close its owned
`.bench/gate.sh` automatically when every repository and executable input is
fingerprintable. An opaque `BENCH_GATE` command or auto-detected gate must
declare the additional local inputs that complete its subject; without that
declaration Bench reruns the gate instead of reusing a verdict. A gate whose
result intentionally depends on remote state is non-reusable. Best-effort
identity never authorizes a commit.

## #2: What exact material identifies each supported oracle?

Blocked by: #1
Type: Research

### Question

Inventory the canonical fingerprint inputs for `.bench/gate.sh`,
`BENCH_GATE`, and every auto-detected gate kind under the reuse ceiling chosen
in #1. Cover resolution kind and command, executable realpath/mode/content,
script and configuration content, selected environment, ignored inputs,
symlinks, and the working-tree hash. Prove claimed sensitivity with runnable
mutation probes against the real resolver and record the cited result in a
short linked research asset.

### Answer

The canonical inventory and runnable mutation probes are recorded in
[the fingerprint research](oracle-bound-gate-verdicts-fingerprint.md). Every
kind shares a fixed subject envelope: schema/policy version, resolved kind and
exact command/cwd, working-tree hash, recognized launcher identity, declared
local inputs, and freshness-policy identity. Kind-specific declarations close
opaque commands, tools, ignored files, and external targets. Missing or
unreadable required material disables reuse. The current full inherited
environment prevents honest closure until #3 decides its posture.

## #3: How does the gate environment become part of a closed subject?

Blocked by: #1, #2
Type: Grill

### Question

The runner currently passes every inherited variable except two wrapper
internals. Choose whether FT78 introduces the minimal gate-environment passlist
needed for a closed subject, fingerprints the entire inherited environment, or
treats any undeclared inherited environment as making the gate non-reusable.
This decides how much of FT88's later environment-minimization work moves into
the oracle fix now.

### Answer

FT78 introduces a gate-specific environment passlist. Gate processes receive a
small required baseline plus project-declared variable names; the values
actually passed join the oracle subject digest. Undeclared inherited variables
do not reach the gate. This moves only the gate-environment slice needed for an
honest verdict into FT78; agent environment minimization and the broader data
inventory remain in FT88. Projects declare additional environment names,
ignored paths, and tool executables in a versioned, repository-owned
`.bench/gate-inputs.json`; the manifest stores names only, while current
values/content feed one aggregate subject digest and never enter the cache.
Missing declared material disables reuse and runs the gate again.

## #4: When is an otherwise matching verdict fresh?

Blocked by: #1, #2, #3
Type: Grill

### Question

Choose the freshness contract: lifetime, future-clock and malformed-time
posture, whether a project may override the lifetime, and whether changing the
policy itself changes the oracle identity. The result must preserve the fast
gate-then-commit path without making an old external environment authoritative.

### Answer

A matching green verdict is reusable for 10 minutes after `recorded_at`, with
no project override in v1. The record stores only `recorded_at`; expiry derives
from the single policy constant. At or after expiry, or when the timestamp is
missing, malformed, or in the future, the caller reruns the gate rather than
failing for staleness alone. Red, pending, and invalid records never authorize
reuse. The freshness-policy version joins the oracle subject, so a policy
change invalidates earlier records.

## #5: Which persistence failures fail which requested actions?

Blocked by: #2, #3, #4
Type: Grill

### Question

Set the observable behavior for standalone `bench gate`, `bench commit`,
shift, and Stop when pre-run invalidation, final green/red replacement, file
sync, rename, or directory sync fails. Also decide how cancellation and
tree-oracle drift during a run are classified. The critical invariant is
already fixed: an unrecordable red must leave no reusable older green and must
fail the requested action.

### Answer

Durable verdict persistence is part of success for every gate-running action.
The runner snapshots the subject, durably replaces any prior record with
`pending`, executes that captured oracle, rejects post-run subject drift, and
atomically installs and syncs the final green or red. Failure to establish
`pending` prevents execution. Persistence failure, subject drift, or an
unrecordable result makes standalone gate, commit, shift, and armed Stop fail
their requested action; commit never stages, shift preserves work, and Stop
blocks. Cancellation leaves durable `pending` and preserves its interrupted
result. A durable red preserves the gate's nonzero exit; other operational
failures use action-level exit 1 while retaining the original gate exit and
output in diagnostics.

## #6: What crash and concurrency protocol preserves the verdict invariant?

Blocked by: #5
Type: Prototype

### Question

Compare at least a versioned pending record with attempt identity and
compare-and-swap against an execution lock or remove-then-replace protocol.
Exercise overlapping green/red runs, crash after invalidation, cancellation,
write/sync/rename failures, and a superseded writer. Return the smallest state
machine that prevents an old or earlier green from becoming reusable.

### Answer

The prototype selects one exclusive process lock per Git directory, held from
before durable `pending` replacement through final verdict sync. Gate
executions for one repository cannot overlap; different repositories remain
independent. A process crash releases the OS lock while leaving `pending`
non-reusable. This is smaller than a short lock plus attempt-ID compare-and-swap
and avoids superseding useful in-flight work. Unlocked remove-then-replace is
rejected because the prototype demonstrated that an older green can overwrite
a newer red. Lock contention fails immediately with operational exit 1; commit
does not stage, shift preserves work, and armed Stop blocks. A lock-free
`pending` record is abandoned or interrupted evidence, so the next lock owner
may durably replace it and run normally. The OS lock is the sole live-owner
authority; PID and age are diagnostic only. Status distinguishes locked
`pending` from lock-free interrupted/incomplete state, and neither authorizes
reuse.

## #7: Where does the gate observe every closure case?

Blocked by: #2, #5, #6
Type: Research

### Question

Map black-box contracts and biting canaries through the real CLI for
gate-command change, gate-script or executable-content change, auto-detected
oracle change, expiry, legacy and malformed cache, every durable-write failure,
mid-run drift, crash, and overlapping writers. Name the lowest extra seam only
where the CLI cannot make a failure observable.

### Answer

The black-box matrix, lowest injected seam, canary attachment, and hostile-input
ownership are recorded in
[the coverage research](oracle-bound-gate-verdicts-coverage.md). Real-CLI
runtime contracts own authorization, closure, freshness, durable invalidation,
execution ownership, consumer projection, and the commit/shift/Stop effects.
An injected clock/filesystem/lock seam inside `internal/gate` is limited to
individual sync/rename/clock/lock faults that cannot be reproduced portably.
Two behavior-owned canaries independently sabotage oracle comparison and
pre-run durable invalidation.

## #8: What is the persisted record and diagnostic contract?

Blocked by: #1, #2, #3, #4, #5, #6, #7
Type: Grill

### Question

Choose the strict versioned encoding, byte bound, file mode, pending/ready
fields, and which diagnostic material may be stored. The cache must distinguish
locked and abandoned pending state without making PID or age authoritative,
reject legacy/malformed content, and avoid persisting command strings,
environment values, or declared-input content.

### Answer

The cache is a strict schema-1 JSON object capped at 16 KiB and stored as a
regular, non-symlink `0600` file. Ready records contain state, green/red status,
tree hash, aggregate oracle digest, and `recorded_at`; pending records contain
state, pre-run tree and oracle digest, `started_at`, and diagnostic
`owner_pid`. PID and time never decide lock ownership. Unknown or duplicate
fields, trailing content, unsupported schema, bad enums/hashes/times, oversize,
wrong file type, and legacy text are invalid and non-reusable. Readers classify
but never repair or echo hostile bytes; a later gate run replaces invalid state
through the normal pending protocol. Command strings, environment names or
values, declared paths/content, and tool output never enter the cache. Writes
use a same-directory `0600` temporary file, file sync, atomic rename, and
directory sync.

## #9: Which deep module interface and build slices carry the result?

Blocked by: #1, #2, #3, #4, #5, #6, #7, #8
Type: Grill

### Question

Choose the public shape and spec slicing after the behavior is fixed. The
design comparison converged on `internal/gate` owning resolution,
fingerprinting, strict parsing, reuse classification, execution, and durable
replacement; `status`, dashboard, roadmap context, commit, shift, and Stop are
consumers. The remaining interface choice is a minimal concrete inspect/run
API, a zero-configuration concrete oracle object, or an extensible verdict
facade with injectable resolver and store ports. Decide whether persistence,
consumer migration, and concurrency land as one atomic spec or ordered slices.

### Answer

`internal/gate` exposes one concrete deep interface: `Inspect(root)` returns a
typed inspection with the cache state, ready status, cached/current trees,
recorded time, reason, and the sole `ReusableGreen` predicate; `Execute(ctx,
root, stdout, stderr)` returns both the original gate exit and the requested
action's final exit plus its resulting inspection. Commit inspects then executes
unless reuse is authorized; standalone gate and shift execute; status,
dashboard, and roadmap context inspect only; Stop invokes the standalone path
and no longer records a second verdict. Resolver, manifest collectors, clock,
filesystem, and lock stay behind an unexported engine for tests. No public
Resolver, Store, or Oracle interface is added for the single production
implementation. FT78 ships as one atomic spec with four ordered stories:
subject/configuration and `Inspect`; locked durable execution and `Execute`;
consumer migration and deletion of the old parser/second writer; then the
real-CLI, injected-fault, hostile-input, and canary proof. Separate landed specs
are rejected because the new schema, writer, reuse predicate, and consumers are
one authorization invariant; an intermediate compatibility layer would
duplicate cache knowledge without independently closing the defect.

## Not yet specified

n/a — all build-shaping fog is resolved in the tickets and Handoff.

## Out of scope

- Changing the existing gate-resolution precedence.
- Replacing the separately reviewed `.bench` gate pin used by pre-push.
- Guaranteeing equivalence of remote services, network responses, or other
  undeclared transitive dependencies.
- FT82 release-preflight evidence and FT88's full agent/gate environment
  minimization beyond inputs required to make this verdict honest.
- Redesigning status, the dashboard page, or roadmap context beyond consuming
  the authoritative typed verdict projection.

## Handoff

1. **Module boundaries.** `internal/gate` is the single deep owner of gate
   resolution, subject construction, `.bench/gate-inputs.json`, the gate
   environment passlist, strict cache parsing, freshness, inspection, execution
   locking, pending invalidation, and atomic synced finalization. Keep those
   responsibilities in one package, split into cohesive files rather than new
   packages. `internal/git` continues to own the working-tree hash and Git-dir
   resolution. Commit, shift, standalone gate, and Stop are action consumers;
   status, dashboard, and roadmap context are read-only projection consumers.
   The manifest is project-owned runtime configuration, not prose and not a
   second resolver. Gate pinning remains separate.
2. **Contracts.** `Inspect(root)` is read-only: it resolves the current subject,
   reads at most 16 KiB, classifies absent/ready/pending/invalid/unavailable,
   distinguishes locked from abandoned pending through the OS lock, and alone
   answers `ReusableGreen`. A ready green reuses only when schema, closed
   subject, tree, oracle digest, and fixed ten-minute freshness all match.
   `Execute(ctx, root, stdout, stderr)` snapshots the subject, fails immediately
   on lock contention, durably installs pending, runs that exact oracle from the
   repository root with only the baseline plus declared environment, rejects
   drift, and durably installs green/red. It returns both gate and action exits:
   0 only for durable green, the original nonzero for durable red, 3 for no
   oracle, the interruption result for cancellation, and 1 for lock,
   persistence, subject, or operational failure. Open-subject gates still run
   and record their latest result for diagnostics but never authorize reuse.
   Commit does not stage on failure; shift preserves work; armed Stop blocks and
   never writes a second verdict.
3. **Deep vs thin.** Subject collectors and the codec/store/lock engine hide
   every fingerprint, schema, filesystem-ordering, and concurrency detail behind
   `Inspect` and `Execute`; their clock/filesystem/lock injections are private
   test seams. `Inspection.ReusableGreen` is the only authorization predicate.
   CLI and status consumers translate typed results to their existing outputs
   and add no parsing or policy. There is one production implementation, so no
   public Resolver, Store, or Oracle interface exists.
4. **Black-box assertables.** Through the real wrapper, assert exact-subject
   reuse by gate-run count and commit state; rerun/refusal after gate command,
   declared environment, ignored input, tool, symlink target, or auto-detect
   changes; no reuse for an open subject; ten-minute, future, legacy, malformed,
   truncated, oversized, and wrong-type behavior; pre-run pending before a gate
   marker; no run on pending-write failure; no old green after finalization
   failure; drift/cancellation/crash state; immediate contention failure;
   lock-free pending recovery; no staging or lost shift work; and one typed
   projection across status, dashboard, and roadmap context. The full matrix and
   named red signals are in
   [the coverage research](oracle-bound-gate-verdicts-coverage.md).
5. **Gate attachment.** Attach observable behavior to
   `internal/contract/runtime` at the CLI seam, extending the existing gate,
   commit, status, and shift families. Use one injected `internal/gate`
   clock/filesystem/lock seam only for deterministic file-sync, directory-sync,
   after-rename, clock, strict-byte-bound, and lock-acquisition faults the CLI
   cannot portably induce. Add behavior-owned oracle-comparison and durable-
   invalidation canaries with the two named contract messages. The gate can see
   every repository-controlled seam; equivalence of declared remote state is
   intentionally outside the claim.
6. **Hostile-input owners.** The manifest/cache codec owns absent versus empty,
   present JSON with or without a final newline, second-token/trailing content,
   16 KiB edges, wrong type/schema, control bytes, and path strings containing
   spaces or glob characters; decoded control bytes are rejected rather than
   rendered. Subject construction owns executable modes,
   symlink chains and external targets, ignored files, missing tools/variables,
   the exact single-argument `BENCH_GATE` string, and deep-CWD root resolution.
   Execution owns signals, crashes, repeated runs, same-Git-dir contention, and
   plan/run drift; separate worktree Git dirs remain independent. Consumer
   contracts reach standalone gate, commit, shift, Stop, status, dashboard, and
   roadmap context. Destructive worktree classification remains outside this
   feature.
7. **Uncertainty flags.** None. The spec must inventory the smallest executable
   baseline and kind-specific collectors from the supported command paths, but
   that is factual research constrained by the closed-subject policy, not a new
   discretion to widen environment inheritance or accept best-effort reuse.
8. **Rejected alternatives.** Rejected best-effort identity; reuse for remote-
   dependent gates; full inherited-environment fingerprinting; ambient or prose
   input declarations; project freshness overrides; asymmetric green-write
   success; unbounded lock waits; PID/age ownership; concurrent attempt-ID CAS;
   unlocked remove/replace; raw commands, environment values, or input content
   in the cache; reader repair of hostile records; exported hypothetical ports;
   a second cache parser/writer; and separately landed core/consumer specs with
   a temporary compatibility layer.
9. **Domain watch-outs.** The working-tree hash excludes ignored files and
   external symlink-target content, so declared collectors are load-bearing.
   The current runner inherits almost the whole environment and Stop records a
   second verdict; both must disappear atomically with consumer migration. A
   synced pending replacement before execution is what prevents an old green
   from surviving a crash or unrecordable red. Final rename followed by a failed
   directory sync may expose the correct new record in-process, but the already-
   synced pending entry is the crash-recovery floor. A valid manifest is the
   project's assertion that an opaque command's dependency list is complete;
   Bench verifies every declared input but cannot infer an omitted shell or
   remote dependency. Inspection may stream executable content to hash it but
   never executes a tool; large-input latency must remain bounded without
   weakening content identity. Cache timestamps and PIDs are diagnostic
   evidence, never ownership authority.

Dependency order: n/a — single spec. Within it: subject/configuration and
`Inspect`; locked durable execution and `Execute`; consumer migration and old-
owner deletion; real-CLI, injected-fault, hostile-input, and canary proof.
