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

— (open: reviewer decision on the authority ceiling before fingerprint
research)

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

— (open: research not run)

## #3: When is an otherwise matching verdict fresh?

Blocked by: #1, #2
Type: Grill

### Question

Choose the freshness contract: lifetime, future-clock and malformed-time
posture, whether a project may override the lifetime, and whether changing the
policy itself changes the oracle identity. The result must preserve the fast
gate-then-commit path without making an old external environment authoritative.

### Answer

— (open: reviewer decision after the oracle subject is known)

## #4: Which persistence failures fail which requested actions?

Blocked by: #2, #3
Type: Grill

### Question

Set the observable behavior for standalone `bench gate`, `bench commit`,
shift, and Stop when pre-run invalidation, final green/red replacement, file
sync, rename, or directory sync fails. Also decide how cancellation and
tree-oracle drift during a run are classified. The critical invariant is
already fixed: an unrecordable red must leave no reusable older green and must
fail the requested action.

### Answer

— (open: reviewer decision on the symmetric fail-closed posture and CLI
results)

## #5: What crash and concurrency protocol preserves the verdict invariant?

Blocked by: #4
Type: Prototype

### Question

Compare at least a versioned pending record with attempt identity and
compare-and-swap against an execution lock or remove-then-replace protocol.
Exercise overlapping green/red runs, crash after invalidation, cancellation,
write/sync/rename failures, and a superseded writer. Return the smallest state
machine that prevents an old or earlier green from becoming reusable.

### Answer

— (open: prototype not run)

## #6: Where does the gate observe every closure case?

Blocked by: #2, #4, #5
Type: Research

### Question

Map black-box contracts and biting canaries through the real CLI for
gate-command change, gate-script or executable-content change, auto-detected
oracle change, expiry, legacy and malformed cache, every durable-write failure,
mid-run drift, crash, and overlapping writers. Name the lowest extra seam only
where the CLI cannot make a failure observable.

### Answer

— (open: coverage research not run)

## #7: Which deep module interface and build slices carry the result?

Blocked by: #1, #2, #3, #4, #5, #6
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

— (open: reviewer decision after research and prototype evidence)

## Not yet specified

- Exact versioned record encoding, bounded-read limit, permissions, and safe
  diagnostic projection.
- Whether project-declared fingerprint inputs need a new configuration surface
  or reuse an existing one.
- The final Handoff seam table and dependency order, which depend on #1–#7.

## Out of scope

- Changing the existing gate-resolution precedence.
- Replacing the separately reviewed `.bench` gate pin used by pre-push.
- Guaranteeing equivalence of remote services, network responses, or other
  undeclared transitive dependencies.
- FT82 release-preflight evidence and FT88's full agent/gate environment
  minimization beyond inputs required to make this verdict honest.
- Redesigning status, the dashboard page, or roadmap context beyond consuming
  the authoritative typed verdict projection.
