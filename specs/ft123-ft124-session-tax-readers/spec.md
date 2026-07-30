# FT123 + FT124 — portable worktree handles and structured test triage

Status: implemented

Compiled from the FT123 and FT124 evidence rows in `ROADMAP.md` under the
reviewer-directed batch-drain override: the reviewer explicitly invoked
`$bench-write-spec for FT123 + FT124`, and both rows are reviewed findings from
the week-of-2026-07-19 transcript evidence. No per-feature decision map exists.
Every decision supplied by this spec rather than by those evidence rows is
marked **[defaulted]** for post-hoc veto.

## Problem

A session that needs its owned worktree still copies a machine-specific,
110-character absolute path into every command. The label already recorded with
the assignment is visible in `bench worktree list`, but no command resolves it
back to the tree or runs a command there.

The same session also repeatedly rebuilds lossy shell filters around `go test`.
Those filters hide package structure, discard failure context, and make honest
skips indistinguishable from deleted assertions. A bare package-level `ok` is
especially misleading when a capability or environment skip was swallowed by
non-verbose output.

## Solution

Add `bench worktree path <target>` and
`bench worktree exec <target> -- <command> [args...]`. Both resolve the same
Bench-owned active assignment by exact label or by an absolute/`~`-relative
path. `path` emits the resolved tree with the current home compacted to `~`;
`exec` runs an argv directly from that tree while preserving the child streams
and exit status. Ambiguous labels, stale assignments, foreign registrations,
and ownership mismatches fail closed.

Add `bench test [package]`, defaulting to `./...`. It runs the real
`go test -json -count=1` interface from the repository root and emits three
stable TOON blocks:

```text
packages[N]{package,status}:
failures[N]{package,test,line}:
skips[N]{package,test,reason}:
```

The empty `failures[0]` and `skips[0]` headers are load-bearing evidence. They
distinguish a fully observed pass from the old bare `ok` without inventing a
sub-second timing heuristic. Test failures remain a nonzero result; the gate
remains the only oracle that calls work done.

## User stories

1. As a session holding a Bench-owned worktree label, I want
   `bench worktree path <label>` to resolve exactly one active, present,
   registered, ownership-matched assignment, so I can recover my tree without
   copying its hashed absolute path. Line: `gpt-5.6-luna` / low. The assignment
   reader and worktree command seam already exist, and a runtime contract fully
   grades the resolver.
2. As a session moving commands between machines, I want both worktree
   subcommands to accept exact labels, absolute assignment paths, and
   `~`/`~/...` assignment paths, while `path` compacts a tree under my current
   home to `~` form, so transcripts and committed handoffs do not embed
   `/home/<user>`. Line: `gpt-5.6-luna` / low. Portable path normalization is
   deterministic plumbing at the cached CLI line.
3. As a session facing duplicate labels or stale worktree state, I want
   resolution to reject zero matches, multiple matches, non-active assignment
   states, missing or foreign registrations, and ownership-marker mismatches
   without mutating lifecycle state, so a convenient handle never selects the
   wrong checkout. Line: `gpt-5.6-terra` / medium. Fail-closed ownership joins
   are correctness-sensitive and a wrong selection could execute arbitrary
   work despite a known package seam.
4. As a session running work in its tree, I want
   `bench worktree exec <target> -- <command> [args...]` to invoke the argv
   directly with that worktree as cwd, inherited stdin and environment, exact
   child stdout/stderr, and the child's exit status, so I can avoid a shell
   string and an explicit `cd`. Line: `gpt-5.6-luna` / low. The direct-process
   contract is exact, familiar, and black-box observable.
5. As a session discovering the worktree-local measurement route, I want
   general help and `bench worktree --help` to advertise `path`, `exec`, and
   `bash bin/bench.sh gate --fresh` as the command that exercises the current
   worktree's gate, so I do not accidentally measure the main checkout's
   installed wrapper. Line: `gpt-5.6-luna` / low. This is a small help-surface
   addition whose exact text is gate-observable.
6. As a Go developer, I want `bench test` to run
   `go test -json -count=1 ./...` from the repository root and
   `bench test <package>` to replace only the package expression, so one command
   works identically from the root or a deeper cwd and never reports a cached
   run as fresh evidence. Line: `gpt-5.6-luna` / low. The invocation contract is
   exact and a real-Go fixture grades its argv and cwd.
7. As a session triaging a test run, I want one deterministic package row for
   every observed package with status `pass`, `fail`, or `no-tests`, so I can
   see the tested population without re-parsing raw Go output. Line:
   `gpt-5.6-terra` / medium. Streaming and aggregating concurrent Go JSON events
   is a new deep parser, although its output seam is fully specified.
8. As a session triaging failures, I want one row per directly failing test
   leaf with its first non-runner diagnostic line, plus a package-level row for
   build or setup failures that have no test name, so the first actionable
   failure survives without dumping the raw stream. Line: `gpt-5.6-terra` /
   medium. Go's event attribution and parent/subtest failure collapse need
   careful semantics even though fixtures make them observable.
9. As a reviewer checking whether green is honest, I want every skipped test
   listed with its reason at default verbosity, preferring the structured
   capability/environment reason when present and otherwise retaining the
   test's emitted skip diagnostic, so a skip cannot look like a passing
   assertion. Line: `gpt-5.6-terra` / medium. The load-bearing distinction
   crosses the Go event stream and the kit's capability protocol.
10. As an agent consuming `bench test`, I want help and usage to follow the
    shared grammar, successful and failing runs to emit only the three TOON
    blocks on stdout, usage to exit 2, unsatisfiable execution to emit a
    structured stdout error and exit 1, and test failures to emit their tables
    and exit 1, so the command is AXI-conformant and scriptable. Line:
    `gpt-5.6-terra` / medium. AXI stream and exit-code semantics are exact but
    span every result class.
11. As a session reading hostile or very long diagnostics, I want control
    characters escaped and each default `line` or `reason` preview capped by
    the shared 120-code-point policy with its original byte count, while
    `--full` emits the complete escaped selected line, so terminal control and
    token floods do not corrupt triage or silently erase evidence. Line:
    `gpt-5.6-luna` / low. The shared sanitizer and TOON emitter already own the
    policy and its boundary tests.
12. As a kit maintainer, I want both commands registered once in the compiled
    CLI, reachable through the real kit and linked-repository wrappers, covered
    by the runtime contract and command-routing registry, and included in the
    canonical CLI inventory, so no shipped surface or advertisement drifts.
    Line: `gpt-5.6-terra` / medium. Cross-surface routing and single-source
    advertisement are conformance concerns at the profile's gate-logic line.

## Implementation decisions

- **[defaulted] One combined spec, two independent seams.** FT123 and FT124
  share the measured session-tax theme but no production module. Worktree
  resolution stays in `internal/worktree`; structured test execution lives in
  a new `internal/testreport` package. Their build slices may land separately
  and green because neither imports the other.
- **[defaulted] Resolver contract.** One read-only resolver serves both `path`
  and `exec`. A non-path target is an exact, byte-for-byte assignment label. An
  absolute target, `~`, or a target beginning `~/` is expanded and matched as a
  path. Other `~user` forms and relative path spellings are rejected rather
  than guessed. Path syntax takes precedence over a label that happens to look
  like a path.
- **[defaulted] Ownership join.** A usable result has exactly one validated
  assignment in `active` state, an existing exact Git worktree registration,
  and the immutable ownership marker that matches the assignment's owner and
  identity. The resolver composes the same intent reader, Git worktree facts,
  and marker validation as lifecycle code; it does not parse `bench worktree
  list` output or create a second assignment schema. A duplicate label fails
  with a pointer to `bench worktree list`; its exact path is the
  disambiguator.
- **[defaulted] Portable rendering.** Home compaction is path-component aware:
  the home itself becomes `~`, a descendant becomes `~/relative/path`, and a
  path outside home remains absolute. `path` writes exactly that value plus a
  newline on stdout. It refuses a value containing any control rune because a
  plain line cannot safely carry it; spaces and glob characters remain
  verbatim. `exec` accepts the same path forms as one argv value and does not
  depend on the rendered line.
- **[defaulted] Operational rather than AXI worktree faces.** `path` is the
  one-line resolver used by a human or shell; `exec` is a transparent command
  runner. Their usage and errors follow the existing operational worktree
  stderr contract rather than wrapping child output in TOON. `worktree list`
  remains the AXI population query.
- **[defaulted] Exact exec grammar.** The shared `internal/usage` parser records
  whether `--` ended Bench flag parsing so `exec` can require the separator,
  one target before it, and at least one child argv after it without a local
  parser. The child is started directly, never through `sh -c`; a caller who
  wants shell syntax names `sh -c` explicitly. Start errors are Bench errors;
  after start, stdout, stderr, stdin, and the process exit status pass through
  exactly. No assignment or cleanup record is changed.
- **[defaulted] Test grammar and root.** `bench test` accepts `--full` and zero
  or one package positional through `internal/usage`; `--` permits a
  dash-leading package spelling. It resolves the Git root and runs there.
  Zero positionals means `./...`; one is passed verbatim as the sole package
  expression. It exposes no arbitrary `go test` flag forwarding.
- **[defaulted] Real fresh Go interface.** The subprocess argv is exactly
  `go test -json -count=1 <package>`. It is an argv invocation, not a shell
  pipeline. The child inherits the ordinary environment except
  `BENCH_SKIP_LOG`, which is stripped so the capability helper emits its
  reason into the JSON-observed stdout instead of an unrelated gate side
  channel. Interrupt cancellation reaches the process group and returns one
  structured interrupted error instead of a partial table.
- **[defaulted] Streaming event owner.** `internal/testreport` owns the Go JSON
  event schema, subprocess lifecycle, aggregation, stable sorting, and the
  three-block render. It decodes incrementally rather than retaining raw test
  output. Package rows sort by package; failure and skip rows sort by package
  then test. The shared `toon` package remains the only encoder.
- **[defaulted] Package status.** A package terminal `pass` becomes `pass`, a
  terminal `fail` or a nonzero run affecting that package becomes `fail`, and
  Go's explicit `[no test files]` result becomes `no-tests`. A successful
  process that reports no package terminal event is not an authoritative empty
  state; it exits 1 with a structured error. The always-present
  `skips[0]{...}` block, not elapsed time, marks a fully observed no-skip pass.
- **[defaulted] Failure selection.** For each test, the parser retains its
  first non-empty output line that is not a Go runner marker, package summary,
  or structured skip line. A failed test with failed descendants and no direct
  diagnostic is suppressed as a parent aggregate; a failed test with a direct
  diagnostic, or with no failed descendant, is emitted. A leaf that emitted no
  diagnostic says `no diagnostic emitted`. Package/build failures use the
  first package-scoped non-runner diagnostic with an empty `test` cell. This
  defines “first failure line” at the information boundary Go exposes: its
  JSON stream does not distinguish `t.Log` from `t.Error`.
- **[defaulted] Skip selection.** A `bench-skip` line is parsed through
  `internal/capability`, and the reason cell retains its kind and, for a
  capability, class. Otherwise the last non-runner diagnostic attached before
  the test's `skip` action is the reason. If Go emits none, the row says
  `reason not emitted`; a skip is never dropped for lacking prose.
- **[defaulted] Preview and error posture.** Default diagnostic cells use
  `sanitize.Preview`; `--full` uses `sanitize.Controls`. Malformed/truncated
  JSON, failure to start `go`, no matched package, and interrupted execution
  produce `toon.Errorf` on stdout and exit 1. Valid test/build failures still
  render all trustworthy package/failure/skip rows and exit 1. Usage/help exit
  2/0. Stderr is empty for the AXI command.
- The gate does not call `bench test`, and `bench test` never reads or writes
  the gate cache. It is a focused reader for triage; only `bench gate` can
  declare the tree shippable.

## Testing decisions

- **[defaulted]** A good FT123 test drives the built CLI in a throwaway Git repository, creates
  real owned assignments, and observes path/exec behavior at the public
  worktree command seam. Prior art is
  `TestRuntimeWorktreeContracts`, especially its create, ownership, hostile
  surface, list-row, and wrapper subtests.
- **[defaulted]** A good FT124 test drives the built CLI against a throwaway Go module. A fake
  `go` executable supplies malformed, interrupted, and exact JSON edge streams;
  the installed real Go tool supplies the compatibility probe for pass, direct
  failure, build failure, parent/subtest failure, generic skip, structured
  capability skip, and no-test packages. Prior art is the runtime query
  contracts, `internal/capability`'s writer/reader round trip, and the shared
  TOON/sanitizer boundary tests.
- **[defaulted] Red-probe record.** On 2026-07-30 a throwaway real-Git fixture
  created an owned assignment labeled `alpha`. A black-box `path alpha`
  assertion required its exact home-compacted path; a black-box `exec alpha --
  sh -c ...` assertion required the child's spaced/glob argv, cwd, environment,
  stderr, and exit 37. Both assertions exited 1 against the current CLI. A
  second fake-Go fixture ran from a deep cwd and asserted the exact mixed
  package/failure tables and root-relative argv, a structured skip with
  `BENCH_SKIP_LOG` set, and malformed JSON under `--full -- ./...`; all three
  assertions exited 1 while the absent command returned 2. Real-Go probes
  separately observed the conformance environment-skip event and the deliberate
  canary failure that those assertions must preserve. These are behavior
  assertions, not help-presence probes; implementation promotes them into the
  named runtime tests before writing production code.
- **[defaulted]** The cheapest wrong FT123 implementation searches rendered list text and
  chooses the first label, then executes a shell string. Duplicate-label,
  stale/foreign/marker-mismatch, spaced-argv, metacharacter, and exact-exit rows
  all go red on it.
- **[defaulted]** The cheapest wrong FT124 implementation pipes ordinary `go test` through the
  recurring grep expressions and prints an always-green package row. The real
  JSON compatibility fixture, failing exit, build failure, generic and
  structured skips, malformed stream, and zero-package rows all go red on it.
- **[defaulted]** Focused contracts:
  `go test -count=1 ./internal/worktree ./internal/testreport
  ./internal/contract/runtime ./internal/conformance`.
- Gate: `.bench/gate.sh` through `bench gate --fresh`.

### Seam diagram — owned worktree command

    trigger: session invokes bench worktree path|exec
        │
        ▼
    label / absolute path / ~/path
        ──▶ [ assignment + registration + ownership resolver ]
                  ├─▶ portable path line
                  └─▶ direct child argv at resolved cwd ──▶ child streams + exit
                       ◀ tests attach here: built CLI in a real worktree fixture

### Seam diagram — structured Go test reader

    trigger: session invokes bench test [package]
        │
        ▼
    repo root + package ──▶ [ go test -json -count=1 ]
                                      │
                                      ▼
                              [ event aggregator ]
                                      │
                                      └─▶ packages + failures + skips TOON / exit
    fake and real Go fixtures ─────────◀ tests attach here through the built CLI

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | an exact unique active label resolves its owned present tree | built `bench worktree path` in the runtime worktree fixture | observed red 2026-07-30: the FT123 black-box path probe created real label `alpha`, required its exact resolved path, and exited 1 | the assertion crosses the real assignment and rejects a help-only or constant-path stub |
| 2 | `~`, `~/...`, absolute paths, and labels select the same assignment; home-exact, descendant, and prefix-sibling outputs compact only on path components | built worktree path/exec commands | observed red 2026-07-30: the FT123 target-form black-box table reused the real assignment and failed on its first label form; implementation promotes the complete enumerated table before production code | a label-only resolver or substring home replacement fails at least one exact form, including the prefix sibling |
| edge of 2 | spaces and glob characters survive as literal bytes; control-bearing output, whitespace-only target, `~user`, and relative path are rejected | built worktree path/exec commands | observed red 2026-07-30: the FT123 hostile-target table failed before any accepted or refused form could meet its expected line/code tuple | unquoted splitting, glob expansion, raw control output, or guessed path syntax violates a distinct tuple |
| 3 | zero, duplicate, inactive, missing-registration, foreign-registration, and ownership-mismatch states each fail without ledger, Git, or filesystem mutation | read-only resolver at the built CLI | observed red 2026-07-30: the FT123 state-table black-box probe could not reach a resolver and failed its first result/snapshot tuple; implementation promotes all six cases before production code | a first-match resolver, stale-record trust, or cleanup side effect changes an asserted outcome or snapshot |
| 4 | exec requires `--`, invokes argv without a shell, preserves stdin/environment/cwd and exact stdout/stderr/exit | built `bench worktree exec` | observed red 2026-07-30: the FT123 exec probe required spaced/glob argv, cwd, environment, stderr, and exit 37 and exited 1 | a shell-string implementation reinterprets metacharacters, while an exit-mapping wrapper loses 37 |
| edge of 4 | non-TTY stdin passes unchanged; child exits 0, 1, 2, 37, and 130 pass through; SIGINT leaves no descendant or worktree-state mutation | built worktree exec with signal-aware child | observed red 2026-07-30: the FT123 exec result matrix failed before the current CLI could start its sentinel child; implementation promotes every enumerated exit/signal case before production code | mapping exits to one Bench code, dropping stdin, or signaling only the leader fails an exact result or liveness assertion |
| 5 | both help surfaces name path, exec, and the worktree-local `bash bin/bench.sh gate --fresh` route | wrapper and worktree usage renderers | observed red 2026-07-30: exact `rg -q` assertions for `bench worktree path` and `bench worktree exec` each exited 1 | exact help assertions fail if any discoverability face remains absent |
| 6 | no package runs `go test -json -count=1 ./...`; one package replaces only `./...`; both run at the Git root from root and deep cwd | built `bench test` against a real/fake Go fixture | observed red 2026-07-30: the fake-Go mixed probe ran from `deep/cwd`, required root cwd and exact argv, and exited 1 while the absent command returned 2 | a cwd-relative runner, cached run, or argument-appending implementation produces different captured argv/cwd |
| edge of 6 | zero/one/excess package args, empty arg, help, unknown flag, and a dash-leading package after `--` follow the shared grammar; repository roots and package names with spaces/globs remain one argv | built bench test grammar/runtime fixture | observed red 2026-07-30: the fake-Go grammar matrix's `--full -- ./...` case exited 1 before invoking Go; implementation promotes every enumerated tuple before production code | an ad hoc parser, shell expansion, or root-only fixture cannot match the complete code/stdout/argv matrix |
| 7 | every terminal package appears once, sorted, with `pass`, `fail`, or `no-tests` | `packages` TOON block from the built command | observed red 2026-07-30: the fake-Go mixed package probe required exact sorted package rows and exited 1; the real-Go compatibility probe supplies the no-test form | grep summaries omit no-test packages and an event-per-line renderer duplicates packages |
| 8 | direct failing leaves carry their first diagnostic, parent-only aggregates are suppressed, and build failures use an empty test cell | `failures` TOON block from real Go JSON | observed red 2026-07-30: the fake-Go child-failure probe required `example/fail,TestChild,...boom` and exited 1; direct, nested, and compile-failure real-Go cases are the promoted table | an always-green stub, raw `--- FAIL` parser, or parent/child duplicate cannot match all three cases |
| 9 | generic and structured capability/environment skips appear at default verbosity with reasons, while a no-skip run prints `skips[0]` | `skips` TOON block from real Go JSON | observed red 2026-07-30: the real-Go conformance probe exposed `BENCH_CONFORMANCE_ROOT not set`, and the black-box assertion requiring its package/failure/skip blocks exited 1 | both skip protocols and the explicit empty block are asserted, so swallowing output or keying on elapsed time fails |
| 9 | inherited `BENCH_SKIP_LOG` cannot divert a skip reason away from the reader | built command with a sentinel skip-log environment | observed red 2026-07-30: the fake-Go structured-skip probe set `BENCH_SKIP_LOG`, required `environment: host absent` plus an untouched sentinel, and exited 1 | an unfiltered environment makes the row reasonless or mutates the unrelated side channel |
| 10 | pass exits 0; test/build failure renders all three tables and exits 1; usage/help and start/no-package/malformed errors use their declared stdout/code posture | built AXI command | observed red 2026-07-30: mixed-failure, structured-skip, and malformed-JSON fake-Go assertions expected codes 1/0/1 but each exited 1 because the current command returned 2 | a pass-through `go test`, stderr error, partial table, or always-zero wrapper violates at least one exact tuple |
| edge of 10 | SIGINT during fake Go kills its process group, emits one structured interrupted error on stdout, exits 1, and emits no partial TOON; zero terminal package events are a structured error | built AXI command with signal-aware fake Go | observed red 2026-07-30: the interruption/empty-stream result matrix failed before current dispatch could start its marker process; implementation promotes both cases before production code | killing only the leader leaves the descendant marker live, while streaming partial tables or accepting empty success violates exact output |
| 11 | diagnostics at 120/121 code points and with ESC, BEL, newline, tab, backslash, and invalid UTF-8 remain one representable bounded cell; `--full` restores the complete escaped selected line | TOON diagnostic cells | observed red 2026-07-30: the fake-Go diagnostic matrix's `--full` invocation exited 1 before it could satisfy any exact preview/full cell; implementation promotes all byte/rune cases before production code | raw output makes TOON refuse or forges layout; silent truncation loses the byte-count/full escape hatch |
| 12 | compiled routing, real-kit and linked-repository wrappers, help, runtime contract registration, CLI inventory, and subcommand registry agree on the new surfaces | conformance command-routing and kit-content checks | observed red 2026-07-30: one five-part source/runtime probe reported `compiled-route=1 wrapper-help=1 cli-inventory=1 runtime-contract=1 routing-registry=1`; the runtime contract is the linked-wrapper assertion while the shipped wrapper source is the real-kit assertion | adding only package code, one wrapper face, or one advertisement leaves its named independent probe red |

### Edge inventory

- **Error path:** missing/ambiguous/stale/unowned worktree targets, child start
  and nonzero exit, missing `go`, Go build/setup failure, malformed JSON, no
  matched package, and TOON render refusal all have coverage rows above.
- **Empty or absent input:** worktree path lacks a target; exec lacks a target,
  separator, or child; test has a deliberate zero-positional default; an empty
  Go stream is an error; a no-test package is a visible `no-tests` row.
- **Boundary values:** zero/one/two label matches, home exactly versus one
  descendant versus a prefix sibling, zero/one/many packages, zero/one/many
  failures and skips, 120/121 code points, and child exits 0/1/2/37/130 are
  enumerated in the coverage rows.
- **Malformed input:** whitespace-only target, unsupported `~user`, relative
  path pretending to be a target, flag-looking package before and after `--`,
  truncated JSON, invalid UTF-8, and control-bearing diagnostics are refused or
  rendered according to the rows above.
- **Interrupted or partial state:** SIGINT while exec's child or Go's process
  group is running must leave no child and no partial TOON block. Worktree
  resolution is read-only and must leave assignment/registration bytes
  unchanged.
- **Re-run idempotency:** path resolution and help are byte-idempotent; exec
  deliberately re-runs the requested process; test deliberately uses
  `-count=1` and re-runs the selected tests. **Won't handle:** idempotency of
  reviewer-selected child commands or repository tests belongs to those
  programs, not the runner.
- **Paths and names with spaces or glob characters:** passed as exact argv and
  exercised under `BENCH_HOME` and repository roots containing both.
- **Control bytes in line sinks:** `path` refuses a control-bearing result
  before writing; operational errors escape hostile values; test diagnostics
  use the shared preview/control sanitizer before TOON.
- **Missing trailing newline:** fake-Go JSON's final event is accepted at EOF;
  a selected diagnostic without a newline is retained.
- **Absent, empty, special, or dangling control records:** the shared intent
  reader owns classification before the resolver; no new ledger read exists.
  Ownership-marker wrong type or dangling link fails the ownership join before
  exec.
- **Required tools missing:** missing Git facts fail worktree resolution; a
  missing `go` executable is a structured test error. Neither surface fetches
  or repairs implicitly.
- **Invocation through every shipped surface:** real kit CLI and linked-repo
  by-path CLI are both runtime-tested; the compiled registry and wrapper help
  are conformance-tested.
- **Deep cwd:** both commands resolve repository/intent state from a deep cwd;
  test package expressions are always evaluated at the repository root.
- **Non-TTY stdin:** exec passes it through unchanged; test does not prompt.
- **Destructive worktree state:** foreign, identity-mismatched, reused,
  missing, and non-active trees are rejected. Neither new command cleans,
  unlocks, releases, or creates a tree.
- **Compatibility probe:** the installed Go tool is the official producer for
  `go test -json`; the real-Go fixture covers its pass/fail/skip/no-test/build
  outputs in addition to fake-stream parser faults.
- **Won't handle:** unassigned pooled or foreign worktrees — FT123 explicitly
  resolves a session's owned assignment, and admitting foreign registrations
  would bypass the ownership contract.
- **Won't handle:** `~other-user` expansion — it is machine- and account-
  specific, which is the portability failure this feature removes.
- **Won't handle:** arbitrary `go test` flag or multi-expression forwarding —
  the evidence calls for an optional single package reader; a caller needing a
  bespoke invocation still uses `go test` directly.

## Out of scope

- **FT125 section/symbol slice readers** — a separate reader capability over
  specs and source symbols, with no shared production module or output schema;
  about 10 edits, 4 gate runs.
