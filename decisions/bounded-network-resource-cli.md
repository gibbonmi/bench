# Bounded network, resource, and CLI behavior (FT87)

Status: shaping

## Destination

One explicit policy bounds every Bench-initiated network attempt, subprocess,
read, and output; `BENCH_OFFLINE=1` provably prevents all of them; repair is
explicit and manifest-pinned; and the CLI has one argument grammar. Sources:
`RR:A-14`, `RR:C-06`, `RR:C-07`, `RR:C-09`–`RR:C-12`; `RC:H-04`, `RC:M-02`.

Closed 2026-07-21 under the reviewer's blanket approval of the worker's
recommendations; contestable calls are marked **[veto]** for post-hoc review.

## #1: What does `BENCH_OFFLINE=1` mean, exactly?

Blocked by: none
Type: Grill

### Question

Which operations count as "Bench-initiated network attempts", what does each do
instead when offline, and how does the flag relate to `BENCH_NO_REPAIR`?

### Answer

`BENCH_OFFLINE=1` is the master no-network switch. Under it: binary repair
never runs (resolution failure exits with a structured error naming the flag);
worktree Git refresh is skipped with an observable skip notice; model
discovery skips HTTP providers *and* provider subprocesses (`codex`), emitting
explicit `offline` status rows rather than empty results **[veto — the codex
subprocess is local, but its own egress is not Bench's to vouch for, so it is
excluded too]**. `BENCH_NO_REPAIR` remains as the narrower lever; `BENCH_OFFLINE`
implies it. Every skip is explicit evidence in output, never a silent pass.

FT87 also produces the offline/network-control evidence record the FT83
requirement registry demands for the `public` release profile: a contract run
proving zero repository-initiated network attempts under `BENCH_OFFLINE=1`
(the existing offline sentinel/smoke machinery is the proving harness; the
production flag is what it now exercises).

## #2: Where does the one bounds policy live and what are its defaults?

Blocked by: #1
Type: Grill

### Question

One timeout/size/cancellation/iteration-cap policy must govern agent and gate
execution, repair, Git refresh, model discovery, guard startup, reads, and
large output. What owns it, and with what values?

### Answer

A single Go policy package (`internal/bounds` or equivalent) single-sources
every Go-side named bound; callers replicate the existing
`internal/git/git.go` pattern (`context.WithTimeout` + `exec.CommandContext`,
documented fail-safe default per caller). Hitting a bound yields a distinct
timeout/truncated/incomplete status, never a silent success or a generic red.
Cancellation (SIGINT) propagates through every subprocess tree.

Recommended defaults (spec locks the numerals; constants, not config)
**[veto — values]**: provider subprocess and HTTP deadlines 10s (existing);
Git refresh 30s; guard startup inspection 5s total with degraded-but-honest
output on timeout; HTTP reads via `io.LimitReader` at 5 MB; outline skips
files over 2 MB with skip metadata; shift iteration caps validated as
positive integers ≤ 100, invalid caps a structured error. Gate execution gets
a generous ceiling (45 min) with a distinct `timeout` verdict — a bounded-red,
never a fake green. Agent execution gets cancellation, not a default wall
ceiling: interactive iteration length is legitimately unbounded, and the
iteration cap is the resource control there.

The Node repair script cannot import the Go package and runs before any
binary exists, so it owns its two bounds (total fetch deadline 60s; download
and decompression caps, #3) as its own facts — a documented, deliberate
second runtime, not a drifted copy.

## #3: What is the hardened repair posture?

Blocked by: #1, #2
Type: Grill

### Question

Repair today is implicit-by-default, digest-trusts the transport it fetches
from, unbounded, and never prunes its cache. What ships instead?

### Answer

Repair becomes explicit: it never runs as a silent side effect of resolution
failure. The failure message names the explicit action (`bench repair`, a
wrapper-routed subcommand); `BENCH_REPAIR=1` permits implicit repair for
automation environments **[veto — default flip is user-visible]**. The
expected digest moves out of registry `dist.integrity` into a
manifest shipped inside the wrapper package, pinned independently of
transport metadata; a mismatch is a hard red. Fetch and decompression are
deadline- and size-bounded (recommended: 100 MB download / 200 MB
decompressed **[veto — values]**). Promotion stays write-to-unique-temp +
fsync + atomic rename (already correct); failure cleanup removes only the
process's own temp files, never another process's installed target. Cache
pruning is an explicit `bench repair --prune`, never automatic.

## #4: What replaces the implicit worktree `git fetch origin`?

Blocked by: #1, #2
Type: Grill

### Question

`Acquire` silently runs an unbounded, possibly interactive fetch and discards
the result. Explicit for whom, and what does failure look like?

### Answer

`Acquire` drops the implicit fetch. Refresh is an explicit opt-in
(`--refresh` on the worktree-acquiring surfaces; the shift loop passes it by
its own explicit configuration, not by default). When it runs it is bounded
(#2), noninteractive (`GIT_TERMINAL_PROMPT=0` and `--no-recurse-submodules`
posture), and observable: failure is a structured non-fatal warning with the
underlying git error, never a discarded error. `BENCH_OFFLINE=1` skips it
with a notice (#1).

## #5: How is model discovery bounded and parallelized?

Blocked by: #2
Type: Grill

### Question

Deadline, concurrency, and read-bound shape for `bench models`.

### Answer

The three providers query concurrently (goroutines joined by a WaitGroup);
`runCommand` moves to `exec.CommandContext` under the policy deadline; HTTP
bodies read through `io.LimitReader`. A provider that times out or oversizes
returns a distinct per-provider error row; the others still report. Offline
behavior per #1.

## #6: What does default `bench outline` emit?

Blocked by: #2
Type: Grill

### Question

Default outline emitted 169 KB with no cap. What is the bounded default?

### Answer

Default output is a bounded summary: total file and symbol counts plus the
first N symbol rows (recommended N=200 **[veto — value]**) and explicit
truncation metadata (`truncated=true`, counts of omitted rows/files) in the
TOON output. `--full` restores the complete listing as a deliberate act.
Oversize and unreadable files appear in skip metadata rather than vanishing.

## #7: What is the one CLI argument grammar?

Blocked by: none
Type: Grill

### Question

Parsing is per-command and inconsistent: trailing garbage variously accepted,
help sometimes exits 2, no `--` for leading-dash commit paths, directory
commit paths don't authorize children, coverage resolves slugs CWD-relative.

### Answer

A small shared parsing helper (hand-rolled, beside the existing `toon.Usage`
helpers; no third-party CLI framework) owns the grammar, and every Go
subcommand routes through it: exact arity with trailing garbage rejected as a
usage error (exit 2); `help`/`--help`/`-h` always exit 0; `--` ends flag
parsing so leading-dash paths are expressible; spec-slug fallbacks resolve
from the repository root everywhere (`coverage` adopts the `spec` package's
root anchoring). `bench commit` path arguments gain the conventional
directory grammar: naming a directory authorizes its changed children
**[veto — widens what a commit path stages]**; hostile filenames (dashes,
spaces, globs) are contract-tested.

## #8: How do capability skips become evidence, and deadlines decouple?

Blocked by: none
Type: Grill

### Question

Capability-dependent security tests (`symlink`, `mkfifo`, pid semantics,
multi-CPU) silently `t.Skip`, and fixture ages are multiples of the same
one-minute `staleAfter` an outer timeout can collide with.

### Answer

A shared capability-skip helper replaces bare `t.Skip` in
security-relevant tests, emitting a recognizable structured skip line; the
gate aggregates these into explicit `capability-skips` evidence rows in its
output. On the release/native workflow the skip count must be zero — a
skipped security class there is red. On dev machines skips stay non-fatal but
visible. Deadlines decouple by construction: lease-staleness windows and test
fixture ages derive from one injectable constant, and security-test deadlines
are set independently of (and provably larger than) the subprocess bounds
they contain.

## #9: What is the one user-facing identity and complete package metadata?

Blocked by: none
Type: Grill

### Question

Three names ship today — `bench` (command), `redbench` (npm), `benchkit`
(version output and prose) — and `package.json` lacks repository, homepage,
bugs, and author fields.

### Answer

The user-facing identity is **Bench**: `bench` is the command and product
name in every user-facing string; `redbench` remains the distribution
(registry) name per ADR 0004, appearing only in install/registry contexts;
`benchkit` is retired from user-facing output (it survives only as this
repo's internal project-profile name). `bench version` prints
`bench <version> (<os>/<arch>)` **[veto — identity choice]**. The root and
every platform `package.json` gain `repository`, `homepage`, `bugs`, and
`author`; the canonical URLs are reviewer-supplied at spec time (see
Handoff uncertainty flags).

## Not yet specified

- Whether the shift loop should ever pass `--refresh` by default in CI-like
  environments (today: never implicit).
- Whether repair's shipped digest manifest and FT83's release evidence index
  can later share a generator (single-source candidate once both exist).

## Out of scope

- Host firewalls, egress enforcement, IAM, endpoint controls — outside the
  repository-controlled scope.
- Shift failure/evidence semantics (`RC:H-05`, FT71) and objective data
  exposure (`RR:C-08`, shipped).
- Consumer/maintainer capability separation (FT85) and transactional link
  lifecycle (FT84).
- Gate wall-clock reduction (FT91) — the gate timeout here is a ceiling, not
  a speed fix.
- Reopening the npm distribution identity (`redbench`, ADR 0004).

## Spec-writer discretion

## Sources
