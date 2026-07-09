# CLI hygiene batch (FT51)

Status: staged

## Problem

The `bench` CLI has a cluster of small exit-code, feedback, and wording defects
that make it lie to a wrapping script or misdirect a user:

- An unknown subcommand (a typo like `bench frobnicate`) prints the help text to
  **stdout and exits 0** — indistinguishable from success to any script that
  checks the exit code, and the lone outlier against the exit-2 norm every real
  subcommand already follows for bad args.
- `--version` and `--help` fall into the same catch-all `*)` case, so `--version`
  prints help instead of the version line and `--help` shares the fate of a typo.
- `bench canary` — an oracle command — prints **nothing** on a passing run at exit
  0, giving no "did it actually run?" feedback, where its sibling `bench structure`
  prints `structure ok`.
- `bench coverage` with no argument reports an *unknown argument* for a *missing*
  one.
- `bench link` outside a git repo tells the user to "run inside a Bench-linked
  repo" — nonsensical, since link's whole job is to *create* that linkage.

Separately, a posture question is open and unrecorded: CLI strings that print
`/bench-*` phase names have no decided stance on whether they stay canonical or go
harness-neutral.

## Solution

A batch of small, independent fixes at the shell dispatcher and three Go command
functions, plus one decision-record edit:

1. Unknown subcommand → print usage to **stderr, exit 2**; an explicit help
   request (bare `bench`, `bench help`, `--help`, `-h`) still prints help to
   stdout at exit 0.
2. `--version` routes to the same implementation as `bench version`.
3. `bench canary` success prints one definitive line (`canary ok …`) at exit 0,
   modeled on `structure ok`.
4. `bench coverage` with no argument says the argument is missing/required (not
   "unknown argument"), still exit 2.
5. `bench link` outside a git repo tells the user to run inside a git repository
   (e.g. run `git init` first), exit 1.
6. Record the harness-form posture: CLI-emitted `/bench-*` strings keep the
   canonical `/bench-*` form; the relaying agent translates at the point of
   recommendation. This extends `.bench/BENCH.md`'s existing communication rule to
   explicitly cover CLI-emitted strings.

## User stories

1. As a script wrapping `bench`, I want an unknown subcommand to exit non-zero
   (2) with its message on stderr, so that a typo is distinguishable from success
   and matches the exit-2 norm every real subcommand follows.
   Line: claude-sonnet-5 / low. Shell-plumbing at the known `*)` dispatch seam in
   `bin/bench.sh` with surface-contract coverage, per the profile's cached CLI
   routing.

2. As a user, I want `bench --help`, `bench -h`, bare `bench`, and `bench help` to
   keep printing the help text to stdout at exit 0, so that fixing the typo case
   does not break the legitimate help request.
   Line: claude-sonnet-5 / low. Same dispatch seam and the same contract-tested
   surface as story 1.

3. As a user, I want `bench --version` to print the same version line as
   `bench version`, so that the conventional flag reaches the real implementation
   instead of the help fallthrough.
   Line: claude-sonnet-5 / low. The dispatcher routes the flag to the existing
   `version` path — no new resolver — and the version-routing contract already
   exercises this surface.

4. As an agent running the oracle, I want `bench canary` to print one `canary ok`
   line at exit 0 on success, so that a passing run gives the same definitive
   feedback `bench structure` already does.
   Line: claude-sonnet-5 / low. A one-line addition to `canary.Run` mirroring the
   `structure ok` precedent, tested at the same package seam.

5. As a spec author, I want `bench coverage` with no argument to tell me the
   argument is missing/required rather than "unknown argument", so that the
   message names the actual error, still exiting 2.
   Line: claude-sonnet-5 / low. A dedicated missing-argument message in
   `coverage.Command`, at the seam its command test already drives.

6. As someone adopting Bench in a fresh directory, I want `bench link` outside a
   git repo to tell me to run inside a git repository (run `git init` first), so
   that the remedy matches link's job of creating the linkage.
   Line: claude-sonnet-5 / low. A link-specific message replacing the shared
   AXI `NotInRepo()` phrasing, at the adopt/link seam its surface contract drives.

7. As a teammate reading the communication rules, I want the harness-form posture
   for CLI-emitted `/bench-*` strings recorded, so that the canonical-form-plus-
   translate stance is a decided rule rather than a per-session guess.
   Line: claude-opus-4-8 / high. Guidance-prose edit to `.bench/BENCH.md`: the
   `craft-line` leverage override would route top, but the reviewer capped this
   batch below the top tier, so it routes claude-opus-4-8 / high and the
   orchestrator's top-tier review compensates.

## Implementation decisions

- **Dispatcher (`bin/bench.sh`, the `*)` catch-all).** Stories 1–3 all land in the
  one dispatch seam. The catch-all currently serves both the help/bare-invocation
  path and the unknown-token path at exit 0; split them: an explicit help request
  (`help`, `--help`, `-h`, and the `${1:-help}` bare default) prints help to
  stdout at exit 0; add a `--version`/`-v` case that routes to the existing
  `version` implementation (the same path `bench version` takes, not a second
  resolver); every other unrecognized token prints its message to stderr and
  exits 2. Keep the change inside the existing `case` block — no new resolver,
  no second dispatch table.
- **`bench canary` success (`internal/canary/canary.go`, `Run`).** On a clean
  sweep, print one line beginning `canary ok` to the passed stdout writer before
  returning 0, following `structure ok`'s `<noun> ok (<parenthetical>)` shape
  (e.g. a fixture count). `Run` already takes `stdout`/`stderr` writers, so the
  line is assertable at the package seam without a live gate run.
- **`bench coverage` missing arg (`internal/coverage/coverage.go`, `Command`).**
  The no-arg branch currently returns `toon.Usage("bench coverage", "<spec.md> is
  required")`, which renders the contradictory `usage: bench coverage (unknown
  argument: <spec.md> is required)`. Replace it with a dedicated
  missing/required-argument message that does not use the `unknown argument`
  template, still exit 2. Do not change the flag-shaped and extra-positional
  branches, which are correctly "unknown argument" cases.
- **`bench link` non-git (`internal/adopt/link.go`, `Link`).** Line 37 emits the
  shared `toon.NotInRepo()` ("run inside a Bench-linked repo"). link is an
  adoption command, not an AXI query command, so it gets its own message pointing
  at a git repository / `git init`; leave the shared `toon.NotInRepo()` untouched
  so every genuine AXI query command keeps the shared phrasing (one source stays
  one source). Exit code stays 1.
- **Harness-form posture (`.bench/BENCH.md`).** The record edit extends the
  existing communication clause about `/bench-*` phase forms to state explicitly
  that CLI-emitted strings keep the canonical form and the relaying agent
  translates. It lands in `.bench/BENCH.md` (the canonical source), never copied
  into `AGENTS.md`, so the shared-rule single-sourcing conformance check stays
  green.

### Defaulted decisions (batch-drain: veto here)

This spec is compiled from a reviewer-directed assessment drain (ASSESSMENT.md §4
+ §2 nits + backlog item 10), which substitutes for a decision map. Every call
below was pre-made by the orchestrator; each is post-hoc veto surface:

- Unknown subcommand → usage to **stderr, exit 2** (matches the exit-2 norm).
- `--version` → same impl as `bench version`; `--help`/`-h`/bare/`help` → help,
  stdout, exit 0.
- Canary success wording `canary ok` (modeled on `structure ok`).
- Coverage no-arg → missing/required wording, exit 2 (verified **not** already
  fixed: current code still renders "unknown argument").
- link non-git → git-repository / `git init` wording, exit 1 unchanged.
- CLI `/bench-*` strings keep canonical form; agent translates at point of use.

## Testing decisions

- **What a good test is here.** Drive each surface as a black box and observe its
  exit code, stdout, and stderr — never assert internals. The dispatcher stories
  are driven through the built `dist/bench` wrapper (`bin/bench.sh`); the three Go
  message/feedback stories are driven at their command functions, which already
  return or write their full output.
- **Seams and prior art.**
  - Dispatcher (stories 1–3): the surface contract at
    `internal/contract/surface/go_routing_test.go`, which already runs
    `bin/bench.sh` (via `Fixture.Bench` and `Fixture.Run("bash", …/bench.sh, …)`)
    and asserts `Probe.ExitCode` / `Stdout` / `Stderr` — the exact prior art for
    exit-code and routing assertions against the wrapper. No existing test pins
    the current unknown-subcommand exit-0 behavior, which is why the defect
    slipped in; these rows add that pin.
  - Canary success (story 4): `internal/canary/canary_test.go`, mirroring
    `internal/structure/structure_test.go`'s `structure ok` assertions.
  - Coverage missing arg (story 5): `internal/coverage/coverage_test.go`, which
    already has a "flag-shaped argument stays a usage error, exit 2" case to
    extend.
  - link non-git (story 6): `internal/contract/surface/link_test.go` (surface) or
    the adopt package test — drive `bench link` in a non-git temp dir and assert
    the message.
  - Posture record (story 7): no test seam — a guidance-prose edit graded by the
    conformance doc layer (shared-rule single-sourcing, stale-reference sweep)
    and review, with the orchestrator's top-tier read as the compensating control.
- **Gate command.** The project gate: `.bench/gate.sh`.

### Seam diagram

Seam A — the shell dispatcher (stories 1–3):

    trigger: a user or wrapping script runs `bench <token>`
        │
        ▼
    "frobnicate" ──▶ [ bin/bench.sh case "$1" ] ──▶ usage → stderr, exit 2
    "--version"  ──▶ [   (the *) catch-all +   ] ──▶ version line → stdout, exit 0
    "--help"/""  ──▶ [    --version/help arms   ] ──▶ help → stdout, exit 0
                        ◀ tests attach here: surface/go_routing_test.go drives
                          bin/bench.sh and asserts Probe.ExitCode + Stdout/Stderr

Seam B — the three Go command functions (stories 4–6):

    trigger: `bench canary` / `bench coverage` / `bench link`
        │
        ▼
    (no args)   ──▶ [ canary.Run(stdout,stderr) ] ──▶ "canary ok …"    exit 0
    (no spec)   ──▶ [ coverage.Command(args)     ] ──▶ "…required…"     exit 2
    (non-git)   ──▶ [ adopt.Link(…, stderr)      ] ──▶ "…git repo…"     exit 1
                        ◀ tests attach here: each package's own test file
                          (canary_test / coverage_test / link_test) asserts
                          the returned/written output and exit code

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | Unknown subcommand → message on stderr, exit 2 | `surface/go_routing_test.go` | `bench frobnicate` asserting `ExitCode==2` and non-empty `Stderr` — starts red: `bin/bench.sh` `*)` prints help to stdout and exits 0 (bench.sh:232–262) | A degenerate "still exit 0 / stdout" impl fails the exit-2 and stderr assertions |
| 2 | `--help`/`-h`/bare/`help` → help on stdout, exit 0 (preserved) | `surface/go_routing_test.go` | `bench --help` asserting `ExitCode==0` and help text in `Stdout` — starts red once story 1 lands if help is naively swept into the exit-2 arm | Pins the split so the typo fix cannot cannibalize the legitimate help request |
| 3 | `--version` → identical output to `bench version` | `surface/go_routing_test.go` | `bench --version` asserting `Stdout` equals `bench version`'s and `ExitCode==0` — starts red: `--version` currently hits `*)` and prints help | A fallthrough-to-help impl fails the version-line equality |
| 4 | `bench canary` success → one `canary ok` line, exit 0 | `internal/canary/canary_test.go` | `Run` on a clean sweep asserting `Stdout` contains `canary ok` and returns 0 — starts red: `Run` returns 0 with no stdout today (canary.go:60–64) | The always-silent success stub fails the stdout-contains assertion, exactly as `structure ok` pins `structure` |
| 5 | `bench coverage` no arg → missing/required wording, exit 2 | `internal/coverage/coverage_test.go` | No-arg `Command` asserting output says required/missing and does **not** contain "unknown argument", code 2 — starts red: current branch renders `unknown argument: <spec.md> is required` (coverage.go:218–219) | The unchanged `toon.Usage` template fails the "not unknown argument" assertion |
| 6 | `bench link` non-git → git-repository / `git init` wording, exit 1 | `surface/link_test.go` | `bench link` in a non-git temp dir asserting `Stderr` names git repo / `git init` and not "Bench-linked repo", exit 1 — starts red: link.go:37 emits shared `toon.NotInRepo()` ("run inside a Bench-linked repo") | Keeping the shared AXI message fails the link-specific wording assertion |
| 7 | Harness-form posture recorded for CLI `/bench-*` strings | — (`.bench/BENCH.md` prose) | not TDD-able — guidance-prose edit graded by the conformance doc layer (single-sourcing, stale-reference sweep) and review; orchestrator top-tier review is the compensating control | No red command; a wrong or absent record is caught by review, not a test |

Degenerate-implementation check: for each of rows 1–6 the cheapest wrong
implementation is the *current* code (exit 0 for a typo, help for `--version`,
silent canary, "unknown argument" for a missing one, "Bench-linked repo" for
link), and each row's assertion is confirmed red against that code by the cited
source line. Row 7 is honestly classified as not-TDD-able.

### Edge inventory

Walked against the profile's shell-CLI hostile-input checklist for each mapped
behavior:

- **Unquoted multi-word argument (`$*` vs `$1`)** — covered: an unknown quoted
  token (`bench "foo bar"`) is a single `$1`, unknown → exit 2 via row 1; no new
  behavior.
- **Invocation through a symlink** — covered by reuse: `--version` routes to the
  existing `version` path, which `go_routing_test.go` already exercises through a
  symlinked wrapper; no new resolver is introduced.
- **Invocation through every shipped surface** — covered: all fixes live in
  `bin/bench.sh` and the shared Go binary, so the real kit CLI, the linked-repo
  by-path CLI, hooks, and adapters all reach the identical routed behavior; the
  surface contract drives both `Fixture.Bench` and the by-path wrapper.
- **Error path / empty-or-absent input** — covered: rows 4 (empty canary success
  path), 5 (absent coverage arg), and 6 (absent git repo) are the empty/absent
  cases for their commands.
- **Re-run idempotency** — Won't handle: every changed path is a pure function of
  its input (message text and exit code), idempotent by construction; no state to
  re-run against.
- **Paths/dir names with spaces or glob characters** — Won't handle: this batch
  adds no new path argument; message and exit-code changes touch no path parsing.
- **Control bytes in git-sourced text** — Won't handle: no new git-sourced text is
  rendered; `toon.Table`'s existing refusal is unchanged.
- **Hand-edited file without trailing newline** — Won't handle: no file parsing is
  added or altered by this batch.
- **Required tool missing from PATH / no `readlink -f`** — Won't handle: binary
  and wrapper resolution is unchanged; the dispatcher edit runs after resolution.
- **Interrupt (SIGINT) mid-loop** — Won't handle: none of these commands hold a
  loop, lease, or scratch state to leave behind.
- **cwd deeper than repo root** — Won't handle: `canary` and `link` resolve via
  `git.Root()` from anywhere in the tree; that resolution is untouched.

## Out of scope

- **Symlink-loop cap in `resolve_script_path`** — a separate hardening capability
  (a circular symlink hangs the CLI; `readlink -f` would detect it), tracked under
  FT53's test/hardening batch, not a wording/exit-code fix. ~1 edit, 1 gate run.
- **The BENCH.md communication-clause *guidance* rewrite beyond the CLI-string
  extension** — any broader reshaping of the harness-form rule is its own
  synthesis-discipline edit, distinct from recording this one posture. ~1 edit,
  1 gate run.
