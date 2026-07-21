# Bounded network, resource, and CLI behavior (FT87)

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

## Handoff

1. **Module boundaries.** `internal/bounds` (policy constants + bounded-exec
   helpers) — owns every Go-side bound; `bin/bench-repair-binary.mjs` — owns
   repair bounds, pin verification, promotion, prune; a shared arg-grammar
   helper — owns arity/`--`/help semantics for all Go subcommands;
   `internal/worktree.Acquire` — loses fetch, gains explicit refresh;
   `internal/models` — concurrent bounded discovery; `internal/outline` —
   bounded default rendering; a capability-skip helper package for tests +
   gate aggregation.
2. **Contracts.** `BENCH_OFFLINE=1` → zero Bench-initiated network attempts,
   each suppressed operation an explicit notice/row; bound hit → distinct
   timeout/truncated status (exit and output state it), never silent success;
   repair: no implicit run, pinned-manifest digest mismatch → hard red, exit
   127 path names `bench repair`; usage errors exit 2, help exits 0; gate
   timeout → distinct `timeout` verdict, non-zero.
3. **Deep vs thin.** `internal/bounds` and the repair script are deep (hide
   timeout/size/cancellation mechanics); the arg-grammar helper is deep for
   parsing semantics; per-command `Command` funcs stay thin pass-throughs;
   `Acquire`'s refresh becomes a thin call into bounded-exec.
4. **Black-box assertables.** Offline contract: sentinel proves zero egress
   under `BENCH_OFFLINE=1` (exists as smoke harness; now exercises the real
   flag). Hung-endpoint fixtures prove deadlines fire with distinct status.
   Oversized fixture tarball/JSON proves size caps red before exhaustion.
   `bench outline` default on a large fixture: bounded bytes + `truncated`
   metadata; `--full` unbounded. Trailing-garbage/`--`/help matrices assert
   exit codes. `bench coverage <slug>` from a subdirectory equals from root.
   Gate output contains `capability-skips` rows on a capability-poor fixture.
   Version string and package metadata asserted in artifact contracts.
5. **Gate attachment.** Existing contract/conformance phases host all of the
   above (runtime contracts for CLI behavior, artifact contracts for
   packages, conformance for single-sourcing checks). The offline egress
   proof runs in the native/offline workflow like today's offline smokes —
   the gate sees its recorded result, not live network absence; that seam
   stays workflow-attached, flagged as such.
6. **Hostile-input owners.** Leading-dash/space/glob paths → arg-grammar
   helper + commit path grammar (#7); required tool missing (`node` absent)
   → repair's structured failure (#3); cwd deeper than root → root-anchored
   slug resolution (#7); interrupt mid-loop → cancellation propagation (#2);
   special files/oversize files → outline skip metadata (#6); absent vs
   empty → pinned repair manifest fails closed on both (#3); control bytes →
   existing `toon.Table` refusal (unchanged).
7. **Uncertainty flags.** (a) All numeric bound values are worker
   recommendations — spec locks them, reviewer may adjust cheaply. (b)
   Canonical `repository`/`homepage`/`bugs` URLs are reviewer-supplied at
   spec time. (c) The gate-timeout ceiling interacts with FT91's wall-clock
   work — if FT91 later shrinks the gate, the ceiling can tighten; keep the
   constant single-sourced so that is a one-line move.
8. **Rejected alternatives.** Per-operation env flags instead of one
   `BENCH_OFFLINE` (sprawl; unprovable as a single record). Config-file
   tunable bounds (config surface for constants nobody should tune).
   Registry `dist.integrity` as the digest source (trusts the transport it
   defends against). Adopting a CLI framework (cobra et al.) for the grammar
   (dependency and diff far larger than the hand-rolled helper the toon
   helpers already anchor). Silent capability skips (a skipped security test
   indistinguishable from a passed one). Automatic repair-cache GC (deletion
   policy hidden inside a fetch path).
9. **Domain watch-outs.** The repair script runs where no Go binary exists —
   it cannot share Go-side constants; treat its bounds as its own facts or
   the one-source rule will be "fixed" into a broken import. A gate timeout
   must never be reachable by a healthy whole-tree run — it exists to end
   hangs, not to race FT91. `GIT_TERMINAL_PROMPT=0` without a deadline still
   hangs on non-credential stalls; both are required.

Dependency order: three slices — (1) `internal/bounds` + `BENCH_OFFLINE` +
Go-side bounding (refresh, models, outline, guards, caps); (2) repair
hardening + identity/metadata (wrapper/Node side, produces the FT83 evidence
record); (3) CLI argument grammar + capability-skip evidence + deadline
decoupling. Slice 1 first — it defines the policy the others cite; 2 and 3
are independent of each other. Slicing stays the reviewer's call.
