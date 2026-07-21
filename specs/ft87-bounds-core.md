# FT87 slice 1: bounded Go-side resource policy

Status: implemented

Map: `decisions/bounded-network-resource-cli.md` (tickets #1, #2, #4, #5,
#6, and the Handoff), closed at `abd4a1d`.

## Problem

Bench can initiate repair, Git refresh, provider discovery, gate execution,
guard inspection, repository reads, and shift iteration without one policy that
bounds or suppresses those operations. A worktree acquisition silently fetches,
model providers run serially, large responses and outlines are read without a
cap, and `BENCH_OFFLINE=1` has no production meaning. A hung or oversized
operation can therefore consume a session indefinitely, while skipped or
truncated work can look like success.

## Solution

Add one Go policy package for named resource bounds and result classification,
make `BENCH_OFFLINE=1` the master switch for the network operations in this
slice, and apply the policy at the existing worktree, model, outline, guard,
gate, and shift seams. Every suppression or bound hit remains observable:
offline, timeout, oversized, truncated, and incomplete are distinct states.
Online Git refresh becomes explicit, model discovery becomes concurrent,
default outline output becomes a bounded summary, and healthy unbounded agent
iterations remain governed by validated iteration caps rather than an implicit
wall timer.

## User stories

1. As a Bench maintainer, I want every Go-side duration, read limit, row limit,
   and iteration range in this slice owned by one policy package, so callers
   cannot drift on values or on timeout and cancellation classification.
   Line: gpt-5.6-terra / medium. This changes the shared policy used by the
   oracle, so the project profile's gate and conformance routing applies.

2. As an operator who sets `BENCH_OFFLINE=1`, I want binary repair, requested
   Git refresh, the Codex provider subprocess, and the OpenAI and Anthropic HTTP
   providers suppressed with explicit evidence, so Bench makes none of the
   network attempts this slice owns.
   Line: gpt-5.6-terra / medium. The contract crosses the bootstrap shell and
   Go runtime and its zero-egress proof is workflow-attached rather than fully
   observable from an ordinary local gate run.

3. As a worktree user, I want acquisition to use local refs by default and an
   explicit refresh to be bounded, noninteractive, non-recursive, and
   observable without making refresh failure fatal, so worktree startup never
   silently waits on or discards a Git network operation.
   Line: gpt-5.6-luna / medium. The worktree seam and failure contract are
   exact and gate-observable, while subprocess and hostile-repository cases
   justify medium effort.

4. As a reviewer choosing a line, I want Codex, OpenAI, and Anthropic model
   discovery to run concurrently under independent deadlines and read caps,
   so one slow or oversized provider cannot delay or erase the other results.
   Line: gpt-5.6-luna / medium. The provider set and output contract are exact
   and fixture-tested, while concurrency requires more than a mechanical low
   pass.

5. As an agent locating seams, I want default `bench outline` output to report
   repository and symbol totals, emit at most 200 symbol rows, and name every
   omitted or skipped class, while `--full` deliberately removes only the row
   cap, so a routine query is bounded without hiding incompleteness.
   Line: gpt-5.6-luna / low. This is a deterministic query-surface change at an
   existing AXI seam with direct large-fixture coverage.

6. As a session starting through a harness, I want guard inspection to finish
   within five seconds total and return partial rows plus an explicit
   incomplete status on timeout, so startup cannot hang or pretend it inspected
   every guard.
   Line: gpt-5.6-terra / medium. Time-bounded filesystem inspection has a known
   output seam but needs careful partial-result and goroutine cleanup behavior.

7. As a gate operator, I want gate execution to stop after 45 minutes with a
   distinct timeout verdict and nonzero exit while killing the gate subprocess
   tree, so a hung oracle is bounded red and can never become a reusable green.
   Line: gpt-5.6-terra / medium. Gate execution follows the profile's cached
   oracle routing because an incorrect timeout or cache state could authorize
   the wrong tree.

8. As a shift operator, I want both iteration caps accepted only as integers in
   the inclusive range 1 through 100, unset caps to retain their current
   defaults, and an unset `BENCH_MAX_WALL` to impose no wall ceiling, so invalid
   authority is rejected while long interactive iterations remain
   cancellation-driven.
   Line: gpt-5.6-luna / medium. Existing runtime contracts cover the validation
   range, and the default-wall correction needs a focused process-lifecycle
   regression pass.

## Implementation decisions

- Map-carried reviewer veto surface is locked, not left open: offline mode also
  suppresses the local Codex discovery subprocess, and the following policy
  bullet owns all numeric values including the default outline row count. A
  later veto changes the owning fact rather than creating a second override.
- `internal/bounds` is the single Go owner of policy values and of the helpers
  that apply a parent context, run a subprocess tree, bound a read, and classify
  `timeout`, `canceled`, `oversized`, ordinary exit, and start failure. Callers
  receive classified results instead of recreating `context.WithTimeout`,
  `exec.CommandContext`, `io.LimitReader`, or deadline-error tests.
- The production policy locks these values: provider discovery gets one
  10-second budget per provider (the Codex live and bundled attempts share that
  budget); Git refresh gets 30 seconds; guard inspection gets five seconds
  total; each HTTP model body is limited to 5 MiB; each outline file is limited
  to 2 MiB; default outline emits at most 200 symbol rows; gate execution gets
  45 minutes; both shift iteration caps use the inclusive range 1 through 100.
  The existing shift defaults of 12 main iterations and 4 refactor iterations
  remain unchanged.
- `BENCH_OFFLINE` is enabled only by the exact value `1`. The Go policy owns
  that interpretation for binary callers. The bootstrap wrapper must make the
  same exact check before a binary exists; this is the necessary cross-runtime
  copy, and conformance pins parity so the shell and Go meanings cannot drift.
  `BENCH_NO_REPAIR` remains the narrower repair-only lever, and offline implies
  it without changing its existing online behavior in this slice.
- The offline operation set is enumerated, not open-ended: wrapper binary
  repair, explicit worktree `git fetch`, `codex debug models` including its
  bundled fallback, the OpenAI models HTTP request, and the Anthropic models
  HTTP request. Generic project gate and agent subprocesses are not classified
  as Bench network clients; Bench cancels their process trees but does not claim
  to firewall code they run.
- Binary repair remains otherwise unchanged for slice 2, but the wrapper checks
  offline before resolving or invoking Node. Missing-binary output names
  `BENCH_OFFLINE=1` as the reason repair was suppressed. No repair subprocess is
  started under offline mode.
- The three user-invoked worktree acquisition forms accept explicit refresh:
  `bench worktree --refresh [objective...]`, `bench worktree create --refresh
  --request <opaque-id> --label <work-item>`, and `bench shift --refresh
  <objective...>`. Harness-created worktrees do not refresh because their event
  carries no explicit operator choice. `internal/worktree.Acquire` and
  `internal/worktree.Create` receive the resolved refresh choice; neither
  fetches by default.
- Requested refresh runs `git fetch -q --no-recurse-submodules origin` with
  `GIT_TERMINAL_PROMPT=0` and the 30-second policy context before the worktree
  start ref is resolved. It emits one
  `worktree_refresh[1]{status,detail}:` row with `refreshed`, `failed`, or
  `offline`; failure and offline skip both continue from local refs, and a
  failure's detail preserves the control-safe underlying Git error. No refresh
  flag means no fetch and no refresh row.
- Model discovery starts exactly three provider jobs together and joins them
  before rendering in the stable order Codex, OpenAI, Anthropic. Each provider
  keeps its own result when another fails. Source rows use distinct `status`
  values `available`, `unavailable`, `offline`, `timeout`, and `oversized`;
  hints name the bound or missing credential. HTTP reads use limit-plus-one so
  exactly 5 MiB is accepted and the next byte is classified oversized.
- `bench outline` preserves the existing first
  `outline[N]{file,line,kind,name}:` table. It then emits
  `outline_meta[1]{tracked_files,scanned_files,skipped_files,total_symbols,emitted_symbols,omitted_symbols,truncated}:`
  and `outline_skips[N]{file,reason}:`. Default mode emits the first 200 symbols
  in existing deterministic order. `--full` removes the 200-row cap but retains
  the 2 MiB per-file read cap. `truncated=true` whenever rows or files were
  omitted; skip reasons distinguish `oversized`, `unreadable`, `nonregular`,
  and `binary` instead of silently dropping those files.
- `bench guards` preserves its existing `guards` table and appends
  `guard_scan[1]{status,inspected,total,omitted,reason}:`. Status is `complete`
  with reason `none`, or `incomplete` with reason `timeout`; rows completed
  before the five-second total deadline remain visible. Candidate enumeration
  is inside the same deadline: after enumeration, total and omitted are exact;
  if enumeration itself times out, both cells read `unknown` instead of
  fabricating counts. `--brief` prints the same state as one control-safe line
  before its existing footer. An incomplete startup remains exit 0 because it
  is degraded evidence, not a false complete result or a generic command
  failure.
- Gate execution applies the 45-minute policy inside the shared
  resolve-run-record owner, not only in the shell adapter or shift caller. A
  deadline kills the process group, prints `gate: timeout`, exits 124, and
  records a non-reusable ready verdict with `status: timeout`; ordinary SIGINT
  remains exit 130 with an interrupted pending record. Cache parsing and the
  ambient dashboard recognize timeout explicitly, rank it with the existing
  red gate signal, and direct the operator to inspect the hang and re-run; no
  consumer treats it as red or green by inference.
- Agent subprocesses keep process-group signal propagation and gain no default
  wall ceiling. Unset or empty `BENCH_MAX_WALL` disables the timer; an explicitly
  supplied value retains the existing `(0,24h]` validation and opt-in deadline.
  Invalid `BENCH_MAX_ITERS`, `BENCH_REFACTOR_ITERS`, or `BENCH_MAX_WALL` remains
  a structured `shift_result` usage error before worktree acquisition.
- Existing runtime, AXI, gate, and native/offline contract packages own the
  behavioral tests. Conformance checks that named Go limits have one production
  owner and that the unavoidable shell offline check matches the Go semantics;
  it does not duplicate the numeric values in a second executable registry.

## Testing decisions

- A good test drives the production CLI or gate with hermetic fake Git,
  provider, file, and child-process fixtures and observes exit code plus output
  state. Policy unit tests use short injected durations and bounded readers so
  the gate never waits on production-length deadlines.
- New seam: the `internal/bounds` policy interface, pre-agreed in the map
  Handoff. Tests attach to its classified subprocess and read results; callers
  and tests use the same interface. Existing higher seams remain the `bench`
  CLI, the gate contract, and the native/offline workflow.
- Prior art: `internal/models/models_test.go` already injects command and HTTP
  providers; `internal/contract/axi/axi_outline_test.go` drives the built query;
  `internal/guards/guards_test.go` owns static inspection; gate cancellation
  lives in `internal/gate/story4_proof_test.go`; shift range and process outcomes
  live in `internal/contract/runtime/runtime_shift_outcomes_test.go`; and the
  offline sentinel lives in the native artifact workflow.
- Handoff assertable disposition is complete for this slice: the zero-attempt
  sentinel, hung provider, oversized model JSON, bounded/full outline, and gate
  timeout have rows below. Oversized repair tarballs and the stable evidence
  record belong to slice 2; trailing-garbage/help/`--` matrices, root-anchored
  coverage, and capability-skip rows belong to slice 3; version and package
  metadata belong to slice 2. Those are stated exceptions, not silently missing
  tests.
- Gate command: `.bench/gate.sh`.

### Seam diagram

    trigger: Go caller starts a bounded operation
        │
        ▼
    operation kind + parent context ──▶ [ internal/bounds policy ] ──▶ classified result
    command or reader              ──▶ [ timeout/read/process owner ]     + bounded bytes
                                          ◀ tests attach here: inject a fast policy,
                                            hung child, or limit+1 reader

    trigger: operator or SessionStart invokes bench
        │
        ▼
    env + CLI args + fixture repo ──▶ [ bench CLI: worktree/models/outline/guards/shift ]
                                      ──▶ exit code + TOON/human evidence
                                          ◀ tests attach here: run the built binary and
                                            assert rows, statuses, ordering, and markers

    trigger: bench gate or an in-process shift gate
        │
        ▼
    repo + resolved gate + 45m policy ──▶ [ gate resolve-run-record owner ]
                                         ──▶ child-tree termination + exit + cached verdict
                                             ◀ tests attach here: a hung gate under a fast
                                               injected policy and cache inspection

    trigger: native/offline workflow smoke
        │
        ▼
    BENCH_OFFLINE=1 + egress/process sentinels ──▶ [ shipped wrapper + built bench ]
                                                  ──▶ zero attempt markers + skip evidence
                                                      ◀ tests attach here: workflow record
                                                        consumed by the gate contract

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | all production Go values are the locked policy values and callers do not redeclare them | internal/bounds + conformance | new conformance mutation test, run first against `abd4a1d`, fails because `internal/bounds` and its registry do not exist | a scaffold with duplicated caller constants still fails the owner and duplicate scan |
| 1 | subprocess completion, timeout, parent cancellation, nonzero exit, start failure, exact-limit read, and limit-plus-one read are distinct results | internal/bounds | new table test `go test ./internal/bounds`, run first, fails to compile because the policy interface is absent | an always-error or always-timeout helper cannot satisfy the seven enumerated cases |
| 2 | exact `BENCH_OFFLINE=1` suppresses repair, requested refresh, Codex live and bundled subprocesses, OpenAI HTTP, and Anthropic HTTP | native/offline workflow + bench CLI | extended offline sentinel contract, run first, observes the repair or Codex marker on `abd4a1d` and goes red | the sentinel enumerates all five operation classes, so suppressing only HTTP or only Git cannot pass |
| 2 | every suppressed operation emits a flag-naming offline row or notice and offline mode remains idempotent on a second run | bench CLI | new runtime matrix for missing binary, refresh, and all three model rows, run first, fails because `BENCH_OFFLINE` is ignored | a silent no-op or authoritative empty model set lacks the required evidence on either run |
| edge of 2 | unset, empty, `0`, and `true` do not enable offline; only exact `1` does | internal/bounds + wrapper contract | new value matrix, run first, fails because neither runtime has the exact shared contract | truthy or nonempty parsing fails at least one enumerated non-`1` case |
| 3 | every acquisition form is local-only by default and invokes no fetch | bench CLI | fake-Git marker matrix over worktree subshell, worktree create, and shift, run first, records the current implicit `Acquire` fetch | leaving the fetch in any one caller writes that caller's marker |
| 3 | explicit refresh uses the 30-second bound, disables prompts and submodule recursion, and emits `refreshed` | bench CLI | fake Git asserts argv/env and a short injected deadline; run first, current surfaces reject or absorb `--refresh` | an unbounded command, prompt-capable command, or wrong argv fails before worktree creation is accepted |
| 3 | refresh failure and timeout preserve the underlying safe error, emit `failed`, and continue from local refs | bench CLI | hung and exit-23 fake origins, run first, either hang or emit no warning on `abd4a1d` | a discarded error or fatal refresh cannot produce both the worktree and the classified row |
| 3 | offline plus `--refresh` invokes no Git fetch and emits `offline` while continuing locally | bench CLI | fake-Git marker under `BENCH_OFFLINE=1`, run first, records the current implicit fetch | a notice-only implementation that still starts Git trips the marker |
| 4 | Codex, OpenAI, and Anthropic start concurrently, retain deterministic render order, and preserve successful peers | internal/models + bench CLI | barrier-based provider test, run first, times out waiting for peers because `abd4a1d` calls providers sequentially | a sequential port cannot release the barrier, while completion-order rendering fails the fixed row order |
| 4 | each provider owns one 10-second total deadline and reports `timeout` without collapsing peer rows | internal/models | hung command and endpoint tests under a fast policy, run first, Codex is unbounded and HTTP only reports generic unavailable | generic failure or one shared deadline cannot produce the required per-provider status matrix |
| 4 | model HTTP accepts exactly 5 MiB and reports `oversized` on byte 5 MiB + 1 | internal/models | limit and limit-plus-one response tests, run first, `io.ReadAll` consumes both and never emits `oversized` | an off-by-one limit or unbounded reader fails one of the paired cases |
| 4 | missing credentials, malformed JSON, non-2xx, empty catalogs, and unsafe IDs retain their existing unavailable/empty/filter behavior | internal/models + bench CLI | already covered by `internal/models/models_test.go`; extended assertions require the new statuses without weakening these cases | the existing table tests fail if bounded discovery conflates ordinary provider failures or emits hostile IDs |
| 5 | default outline emits at most 200 rows plus exact totals and `truncated=true` with omitted counts on a 201-symbol fixture | bench CLI | new AXI large-outline contract, run first, emits 201 rows and no metadata on `abd4a1d` | a count-only summary, off-by-one cap, or silent truncation fails the exact row and metadata assertions |
| 5 | `--full` emits every discovered symbol from eligible files while retaining metadata | bench CLI | same 201-symbol fixture with `--full`, run first, exits 2 because the flag is unknown | a default-only cap or a `--full` alias that still emits 200 fails the exact symbol count |
| 5 | exactly 2 MiB is scanned, 2 MiB + 1 is skipped as oversized, and unreadable, nonregular, and binary files each get named skip rows | bench CLI | size-boundary and hostile-file AXI matrix, run first, the oversized regular file is fully read and all skip classes vanish silently | every file class has its own expected reason, so a generic or missing skip cannot pass |
| 5 | empty repo, final line without newline, literal space/glob path, deep cwd, control-byte row, and rerun ordering keep their existing behavior | bench CLI | already covered by `internal/outline/outline_test.go` and `internal/contract/axi/axi_outline_test.go`; metadata is added to those fixtures | retaining the existing assertions prevents the bounded renderer from amputating supported callers or destabilizing order |
| 6 | guard inspection always emits complete metadata when all candidates finish | bench CLI | extended guards contract, run first, has no `guard_scan` table or brief status on `abd4a1d` | a partial-only implementation cannot claim complete without exact candidate counts |
| 6 | one candidate that blocks past the five-second total bound yields partial rows, `incomplete`, reason `timeout`, and exit 0; blocked enumeration reports unknown counts | internal/guards + bench CLI | injected blocking enumerator and inspector under a fast policy, run first, fail because `Rows` has no cancellable total-bound seam | waiting for every candidate hangs; fabricating counts, discarding completed rows, or returning generic red fails output and exit assertions |
| 6 | absent hook directory, empty directory, FIFO/device/socket candidates, malformed JSON wiring, and unreadable headers stay nonblocking and honestly classified | internal/guards + bench CLI | existing special-file tests plus new absent/empty/wrong-type matrix; the new metadata assertions start red | the candidate totals and manifest rows expose silent omission or an attempted special-file read |
| 7 | a gate deadline kills the leader and descendant process group, prints `gate: timeout`, exits 124, and records ready/timeout | gate contract | hung-gate process-tree test under a fast policy, run first, current context cancellation exits 130 and leaves pending | killing only the leader, recording red, or leaving pending fails the PID, exit, message, and cache assertions |
| 7 | a ready/timeout verdict is non-reusable and the ambient dashboard reports timeout as its own signal | gate contract + bench status | cache projection test, run first, rejects `timeout` as an invalid ready status | treating timeout as green, generic red, or invalid fails the exact inspection and dashboard state |
| 7 | SIGINT remains exit 130 with interrupted pending, and a healthy gate just below its bound records ordinary green/red | gate contract | existing cancellation proof plus paired fast-policy deadline boundary test | the pair prevents the new deadline from reclassifying human cancellation or racing a healthy gate |
| 8 | empty or unset main/refactor caps use 12 and 4; 1 and 100 are accepted; 0, 101, negative, and non-integer values exit with structured usage | bench CLI | already covered in `internal/contract/runtime/runtime_shift_outcomes_test.go`; add refactor-cap parity cases before moving policy | both cap names and both inclusive boundaries are enumerated, so validating only one variable or one edge fails |
| 8 | unset or empty `BENCH_MAX_WALL` starts no timer, while an explicit value in `(0,24h]` keeps the existing opt-in deadline | internal/shift + bench CLI | fast loop test, run first, observes `parseWallDuration` returning the current two-hour default | retaining an implicit timer or deleting the explicit opt-in fails one half of the paired contract |
| 8 | SIGINT continues to cancel adapter and gate process trees and preserves dirty recovery state | bench CLI | already covered by runtime shift interruption and gate cancellation contracts | the existing descendant-PID and recovery assertions fail if policy centralization drops propagation or cleanup |

Degenerate implementations checked: a constants-only package fails the
classified-result row; an offline notice that still launches a child fails the
sentinels; a sequential model port fails the barrier; an outline that prints
only counts fails the exact first-200 rows and `--full`; a guard that always
claims incomplete fails the complete case; an always-red gate timeout fails the
healthy boundary case; and cap parsing for only `BENCH_MAX_ITERS` fails the
refactor parity matrix.

### Edge inventory

- Error paths — the refresh-failure, provider-failure, outline-skip,
  guard-incomplete, gate-timeout, and cap-validation rows cover ordinary child
  failure, timeout, malformed response, unreadable input, and cache projection.
- Empty or absent input — the offline-evidence, provider-compatibility,
  outline-compatibility, guard-hostile-input, cap-validation, and wall-default
  rows cover missing keys, empty catalogs/repos/directories, and unset or empty
  resource variables.
- Boundary values — the exact-offline-value, model-read-limit,
  outline-default/size, gate-deadline, cap-validation, and wall-default rows
  cover byte limits, 200/201 rows, deadline edges, iteration 1/100, and explicit
  wall duration limits.
- Malformed input — the exact-offline-value, provider-compatibility,
  outline-skip, guard-hostile-input, and cap-validation rows cover malformed
  environment values, JSON, files, wiring, and cap strings.
- Interrupted or partial state — the refresh-failure, provider-deadline,
  guard-incomplete, gate-timeout/cancellation, and shift-cancellation rows cover
  canceled refresh/providers, partial guard scans, gate process trees/cache
  state, and shift recovery.
- Re-run idempotency — the offline-evidence, offline-refresh,
  provider-compatibility, outline-compatibility, guard-hostile-input,
  timeout-cache, and shift-cancellation rows assert repeatable behavior.
- Hostile environment — the offline-suppression, exact-offline-value,
  offline-refresh, provider-compatibility, outline-skip/compatibility,
  guard-hostile-input, gate-timeout, and shift-cancellation rows cover missing
  credentials/tools, fake or hung processes, hostile file types/names, deep
  cwd, control bytes, and descendant processes.
- Paths and directory names containing spaces or glob characters — the outline
  compatibility row.
- Control bytes in Git-sourced text — the refresh argv/error and outline
  compatibility rows require control-safe errors and retain `toon.Table`
  refusal/filtering.
- Files lacking a trailing newline — the outline compatibility row.
- Absent file versus present-but-empty file — the outline compatibility and
  guard hostile-input rows.
- Special files in discovery paths — the outline-skip and guard hostile-input
  rows.
- Required tool missing from PATH — the offline-suppression row covers repair,
  and the provider-compatibility row covers Codex; existing structured
  unavailable behavior remains.
- Invocation through a symlink — **Won't handle** — wrapper resolution is
  unchanged and already contract-tested; no new path resolution is introduced.
- Invocation through every shipped surface — the offline-suppression and
  default-local-acquisition rows enumerate wrapper, worktree subshell,
  worktree create, shift, and workflow entry points.
- Destructive worktree state — **Won't handle** — refresh precedes creation and
  never changes ownership, cleanup, recovery, or plan/apply policy from ADR
  0005; existing lifecycle contracts remain authoritative.
- Interrupt mid-loop — the gate-cancellation and shift-cancellation rows.
- Cwd deeper than the repository root — the outline compatibility row; other
  commands already resolve through their existing root seams.
- Host firewall, proxy allowlists, certificate policy, disk quotas, CPU limits,
  and process-count limits — **Won't handle** — these are platform controls and
  cannot be truthfully enforced by this repository policy.

## Out of scope

1. **FT87 slice 2: repair hardening, release identity/metadata, and the FT83
   offline/network-control evidence record** — this is a separate bootstrap and
   artifact-governance capability in the Node/wrapper runtime, not the Go bounds
   core; 12 edits, 4 gate runs.
2. **FT87 slice 3: shared argument grammar, root-anchored coverage, directory
   commit paths, capability-skip evidence, and deadline decoupling** — this is a
   separate parser and security-evidence capability with its own command-wide
   migration; 14 edits, 5 gate runs.
