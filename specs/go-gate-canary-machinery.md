# go-gate-canary-machinery

Status: staged

## Problem

The canary tripwire is still implemented as sourced shell in
`.bench/lib/canary-run.sh`. That runner is shipped to consumers, copied by
`bench init`, and sourced by both the benchkit gate and scaffolded consumer gates.
The Go rewrite map closed the direction: the canary sweep moves behind
`bench canary`, while existing sourced consumers keep working through a
compatibility shim.

This is the retire-rule machinery the rest of the gate port depends on. Before any
check-by-check conformance port can safely delete a shell twin, the kit needs a Go
canary runner that can prove a fixture still bites, a gate inner-run mode that
avoids recursion and behavior-contract cost, and a root-parameterized Go
conformance harness that future checks can attach to without putting
kit-conformance logic in the shipped binary.

## Solution

Add a generic `bench canary [root]` subcommand backed by a new `internal/canary`
package. It sweeps `<root>/tests/canary`, copies each `files/` tree to a
throwaway git repo, restores `dot-` path segments to dotfiles, runs
`<root>/.bench/gate.sh` from inside the throwaway repo with
`BENCH_CANARY_INNER=1`, and asserts non-zero exit plus the fixture's `EXPECT`
substring. It also preserves the absent-or-empty harness red and the vacuous
`EXPECT` baseline.

Shrink `.bench/lib/canary-run.sh` to a compatibility shim. It resolves a bench
wrapper, calls `bench canary "$root"`, and folds a non-zero exit into the sourcing
gate's `err`/`fail` posture. New `bench init` scaffolds call
`bench canary "$root"` directly instead of sourcing the lib, while still installing
the shim as a compatibility asset.

Update the benchkit gate to have explicit outer and inner modes. The outer mode
builds the dev binary, runs Go verification with the no-cache posture required by
the map, runs the existing shell contract/conformance layer for now, then calls
`bench canary`. The inner mode, entered only under `BENCH_CANARY_INNER=1`, skips
canary recursion and runs the rest of the gate — including the not-yet-ported
behavior-contract fragments, because fifteen existing fixtures bite only through
them — plus the new root-parameterized conformance harness. Behavior fragments
leave the inner run only when `go-gate-behavior-contracts` retires each family with
its fixtures still biting. This spec adds the harness and wiring only; it does not
port individual conformance checks or behavior contracts.

## User stories

1. As the gate, I want `bench canary [root]` implemented in the Go core behind
   `internal/canary` -- validating absent/empty fixture sets, missing `EXPECT`,
   missing `files/`, vacuous `EXPECT`, fixture copy, recursive `dot-` restore,
   inner gate execution with `BENCH_CANARY_INNER=1`, red exit, targeted substring,
   and per-fixture attribution -- so that the gate's tripwire is testable Go
   machinery instead of sourced shell. Line: gpt-5.5 / xhigh. The map makes this
   the gate-critical retire rule for every later conformance deletion, and a rotted
   canary runner lets the oracle prove nothing.

2. As an existing consumer whose gate already sources `.bench/lib/canary-run.sh`,
   I want that file to become a one-glance compatibility shim that resolves a bench
   wrapper, delegates to `bench canary "$root"`, and calls the sourcing gate's
   `err` on any failure, so that a re-link preserves existing gates while retiring
   the shell sweep implementation. Line: gpt-5.5 / xhigh. The map explicitly
   rejects removing the lib because that would turn already-linked consumer gates
   red until a human edits them, so compatibility is gate-critical.

3. As a newly initialized consumer repo, I want `bench init` to scaffold a gate
   that calls `bench canary "$root"` directly, keeps the red-until-configured
   sentinel, keeps the seeded example fixture, keeps absent/empty harness failures,
   and never clobbers an existing configured gate on a second init, so that new
   consumers use the permanent binary surface without losing the tripwire
   guarantees. Line: gpt-5.5 / xhigh. The scaffold is the first gate many
   consumers get, and a green-by-default or unguarded scaffold would violate the
   canary ADR's threat model.

4. As the benchkit gate, I want explicit outer and inner modes: outer mode runs the
   normal build/vet/test/contracts/canary path, while `BENCH_CANARY_INNER=1` never
   recurses into the canary sweep but keeps running the not-yet-ported
   behavior-contract fragments, because fifteen existing fixtures bite only through
   them, so that each fixture grades the intended root without the sweep
   multiplying itself or orphaning behavior-targeted fixtures. Line: gpt-5.5 /
   xhigh. The map flags per-fixture inner-run cost as uncertain, and the wrong
   inner mode can make the gate recursive, too slow, or silently strip fixtures of
   their bite.

5. As the future conformance port, I want a root-parameterized Go conformance
   harness that accepts `BENCH_CONFORMANCE_ROOT` with a default to the current repo
   root, keeps the real kit root separate from the graded root, and exposes shared
   helpers for root-relative reads, executable checks, substring assertions, and
   shell command probes, so that later Go conformance tests can grade canary
   fixtures without shipping kit checks in the binary. Line: gpt-5.5 / xhigh. The
   map requires conformance checks to live in `go test` rather than the shipped
   binary, and this harness is the uncertain seam future checks depend on.

6. As a reviewer watching the oracle's cost, I want the canary sweep measured and
   kept at or under today's wall-clock by invoking the new Go conformance
   machinery in inner runs only as a filtered or precompiled test run, never
   `go test ./...` for every fixture, so that the fixture set stays usable as a
   normal gate layer. Today's baseline already pays for the behavior fragments per
   fixture, so keeping them in inner mode does not move the bound. Line: gpt-5.5 / xhigh. The Handoff calls this the live uncertainty flag, and a
   technically correct runner that makes the gate impractical is still an oracle
   failure.

7. As the dispatcher and package surface, I want `cmd/bench`, `bin/bench.sh`,
   command help, package contents, and conformance references updated for the new
   `canary` subcommand with no dangling references to shell runner logic and no
   weakening of existing canary fixtures, so that the new surface is reachable
   everywhere the gate and consumers need it. Line: gpt-5.5 / xhigh. This touches
   the CLI-to-binary route and shipped compatibility files, which are gate-critical
   surfaces where stale references silently strand consumers.

## Implementation decisions

- **Package split.** Add `internal/canary` as the deep unit behind
  `bench canary`. It owns fixture discovery, fixture validation, copy/restore,
  baseline execution, inner gate execution, result attribution, and cleanup.
  `cmd/bench` gets a direct `canary` case, matching `gate-run` and `shift`, because
  the command streams stderr/stdout and returns an exit code rather than fitting the
  TOON query command map.

- **CLI contract.** `bench canary [root]` grades `root`, defaulting to the current
  git root. Fixtures live at `root/tests/canary`; the gate under test is
  `root/.bench/gate.sh`. The subcommand runs that gate from inside each throwaway
  fixture repo, not from `root`, so the gate's own `git rev-parse --show-toplevel`
  sees the fixture root exactly like the current shell runner.

- **Canary output contract.** Success is quiet and exits 0. Any failed fixture
  exits non-zero and names the fixture, the expected substring, and the observed
  exit code. Missing harness, empty harness, missing `EXPECT`, missing `files/`,
  vacuous `EXPECT`, copy failure, and gate-start failure are all red and
  attributed. The existing `canary '<name>' did not bite` wording carries unless a
  test pins a better exact string in the same diff.

- **Environment contract.** The canary child env sets `BENCH_CANARY_INNER=1` and
  strips wrapper-routing internals (`BENCH_KIT`, `BENCH_WRAPPER`) through one shared
  helper with `internal/gate`, so fixture gates do not accidentally resolve the live
  kit through leaked routing state. It preserves normal user env such as `PATH` and
  `BENCH_GATE` unless the gate entry explicitly overrides it.

- **Fixture materialization.** Copy the fixture `files/` tree into a fresh temp dir,
  initialize git there, then restore any path segment named `dot-*` to `.<name>`
  after copy. This generalizes the current top-level restore while preserving the
  fixture convention that prevents harnesses from loading fixture dotdirs as real
  config.

- **Compatibility shim.** `.bench/lib/canary-run.sh` keeps its sourced contract
  (`root`, `err`, and `fail` exist; no `set -e` in the caller) but contains no sweep
  logic. It resolves the bench wrapper using the existing resolver when available,
  falls back to `bench` on `PATH`, runs `canary "$root"`, and calls
  `err "canary sweep failed"` on non-zero. Existing consumers can keep sourcing it;
  new scaffolded gates do not.

- **Scaffold posture.** `internal/adopt/init.go` continues installing
  `.bench/lib/canary-run.sh` as a compatibility asset, but the generated
  `.bench/gate.sh` calls `bench canary "$root"` directly. The sentinel, example
  `DO-NOT-SHIP` check, seed fixture, no-`set -e` posture, and second-init
  no-clobber behavior carry.

- **Benchkit gate posture.** `.bench/gate.sh` keeps the `gate: green` /
  `gate: red` lines. The map flagged possible consumers of those lines as
  uncertain; preserving them avoids turning that uncertainty into surface churn.
  The final few-line gate is not completed in this spec because the individual
  shell checks and behavior contracts are not ported yet.

- **Inner-run mode.** Under `BENCH_CANARY_INNER=1`, the benchkit gate skips only
  the canary sweep. Behavior-contract fragments keep running in inner mode:
  fifteen fixtures (`doctor-foreign-clobbered`, `doctor-manager-dir-chosen`,
  `doctor-stale-silent`, `postinstall-guard-bypassed`, `postinstall-nonzero-exit`,
  `session-start-advice-dropped`, `wrapper-args-dropped`, `status-regressed`,
  `roadmap-regressed`, `unscaffolded-bench-file`, `toon-escaping-dropped`,
  `learnings-parse-broken`, `guards-aggregation-dropped`,
  `coverage-extraction-dropped`, `diff-recorded-base-dropped`) bite only through
  those fragments, and they leave the inner run only when
  `go-gate-behavior-contracts` retires each family with its fixtures still biting.
  Root-grading checks that remain in shell also run, plus the new conformance
  harness. As checks port in `go-gate-conformance-checks`, they move from the shell
  root-grading subset into the Go conformance package one at a time.

- **Root-parameterized conformance harness.** Add `internal/conformance` as
  test-only gate machinery. It resolves `BENCH_CONFORMANCE_ROOT` for the graded
  tree and the module root for shared kit assets, with table tests proving those
  roots stay distinct. This spec adds the harness and self-tests only; no real
  conformance check migrates here until `go-gate-conformance-checks`.

- **Cost control.** The inner gate must not run full module `go test ./...` per
  fixture. The implementation chooses either a filtered `go test -run` package
  invocation or a precompiled conformance test binary, but the measured before/after
  wall-clock for `bench gate` must be recorded in the implementation notes and stay
  at or under today's shell-runner timing on the same machine.

## Testing decisions

- **What a good test is here.** Acceptance drives the public surfaces:
  `bench canary` as a subprocess-style command, scaffolded gates through
  `bench init`, and the benchkit gate through `bench gate`. Unit tests attach only
  where shell could not: fixture copy/restore, baseline attribution, env
  construction, and conformance root resolution.

- **Seams.** Three seams get tested:
  - The **canary CLI seam**: `bench canary [root]` drives `internal/canary` against
    throwaway fixture repos and observes exit code plus stderr/stdout.
  - The **gate inner/conformance seam**: `.bench/gate.sh` chooses outer vs inner
    mode, and `internal/conformance` grades `BENCH_CONFORMANCE_ROOT`.
  - The **scaffold compatibility seam**: `bench init` output and
    `.bench/lib/canary-run.sh` are exercised as shipped files, not by reading Go
    internals.

- **Gate command:** `bench gate`. Done for this spec means `bench gate` green,
  `go test -count=1 ./...` green through the gate path, all existing canary
  fixtures still biting, and the measured canary wall-clock recorded.

### Seam diagram

Canary CLI seam:

    trigger: gate or user runs `bench canary [root]`
        |
        v
    root/tests/canary  --> [ internal/canary: discover + validate fixtures ] --> red on absent/empty/malformed
    EXPECT files       --> [ vacuous baseline against an empty temp repo    ] --> red if EXPECT matches baseline
    files/ tree        --> [ copy to temp git repo + restore dot-* segments ] --> throwaway fixture root
    root/.bench/gate.sh--> [ run from fixture cwd with BENCH_CANARY_INNER=1 ] --> require non-zero + EXPECT substring
                            ^
                            tests attach here: Go tests use fake gates/runners for attribution and copy semantics;
                            gate contracts run the real command.

Gate inner/conformance seam:

    trigger: `.bench/gate.sh` invoked normally or by `bench canary`
        |
        v
    normal env            --> [ benchkit gate outer mode ] --> build dev binary, vet, go test -count=1, shell contracts, bench canary
    BENCH_CANARY_INNER=1  --> [ benchkit gate inner mode ] --> no sweep; shell checks + behavior fragments + internal/conformance
    BENCH_CONFORMANCE_ROOT--> [ internal/conformance harness ] --> checks grade fixture root while reading helper code from kit root
                              ^
                              tests attach here: shell contract proves inner mode does not recurse; Go tests prove
                              root override and kit-root separation.

Scaffold compatibility seam:

    trigger: `bench init` or an existing gate sources `.bench/lib/canary-run.sh`
        |
        v
    new repo              --> [ internal/adopt.Init scaffold ] --> .bench/gate.sh with direct `bench canary "$root"`
    existing sourced gate --> [ .bench/lib/canary-run.sh shim ] --> resolved bench wrapper runs `canary "$root"`
    missing runner/bench  --> [ scaffold/lib failure posture  ] --> gate red, not silent green
                              ^
                              tests attach here: gate-link contracts assert scaffold behavior; shim tests source the lib
                              with fake err/root/fail.

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | `bench canary [root]` exists, routes through `bin/bench.sh` to the Go binary, and exits 0 for a fixture set whose fixtures bite | canary CLI | new command test before dispatch exists returns `bench: unknown subcommand: "canary"` / shell wrapper help lacks route | catches a subcommand added only in Go or only in shell, which would strand either the gate or direct users |
| 1 | absent `tests/canary`, empty `tests/canary`, missing `EXPECT`, and missing `files/` each exit non-zero with an attributed error | canary CLI | `internal/canary` table tests with temp roots; before the package exists they do not compile | catches the lazy escape of deleting the tripwire or creating inert fixture directories |
| 1 | vacuous `EXPECT` baseline runs once against an empty fixture repo and fails any fixture whose `EXPECT` appears there | canary CLI | `internal/canary` fake-gate test where baseline output contains the `EXPECT` substring | catches a canary that proves only a generic gate failure rather than its planted regression |
| 1 | each fixture copy restores `dot-*` path segments, initializes git, runs the root gate from the fixture cwd, requires non-zero exit, and requires the targeted substring | canary CLI | `internal/canary` fixture-copy test with `dot-bench`, spaced paths, and fake gate transcripts | catches fixture materialization drift and a runner that accepts "red for the wrong reason" |
| 2 | `.bench/lib/canary-run.sh` contains no sweep loop and delegates to resolved `bench canary "$root"`, turning any non-zero exit into the sourcing gate's `err` | scaffold compatibility | shim source test with fake `bench` returning 7; before the shim it never calls the fake command | catches a compatibility file that still owns a second canary implementation or fails open |
| 2 | deleting `.bench/lib/canary-run.sh` from a scaffolded existing-style gate is still red with `canary runner missing` | scaffold compatibility | already covered by `gate-link-contracts.sh`; keep the row green after the shim rewrite | preserves the source guard that the runner cannot enforce after it has been deleted |
| 3 | `bench init` scaffolds a direct `bench canary "$root"` gate, keeps the sentinel red, keeps the seed fixture, and becomes green after sentinel removal because the seed canary bites | scaffold compatibility | existing init scaffold contracts plus a new assertion that scaffolded `.bench/gate.sh` names `bench canary` and does not source `canary-run.sh` | catches a scaffold that keeps new consumers on the retired sourcing API or drops the tripwire |
| 3 | second `bench init` never clobbers a configured `.bench/gate.sh` | scaffold compatibility | already covered by `gate-link-contracts.sh` second-init no-clobber row | catches scaffold updates that rewrite reviewer-owned gate content |
| 4 | `BENCH_CANARY_INNER=1` skips canary recursion and absent/empty harness errors | gate inner/conformance | existing inner guard contract in `gate-link-contracts.sh` plus a new benchkit gate inner-mode contract | catches the recursive canary failure where every fixture run starts another sweep |
| 4 | inner mode keeps running the not-yet-ported behavior fragments so the fifteen behavior-targeted fixtures keep biting | gate inner/conformance | `bench gate` canary layer reds with `did not bite` naming those fixtures if inner mode skips their fragment | catches an inner mode that silently orphans behavior-targeted fixtures while chasing sweep cost |
| 5 | `internal/conformance` resolves `BENCH_CONFORMANCE_ROOT` as the graded root and keeps kit root separate | gate inner/conformance | Go table test with a temp graded root marker; before the harness exists it does not compile | catches future checks accidentally reading the real kit tree while claiming to grade a fixture |
| 5 | conformance harness defaults to the current git root when `BENCH_CONFORMANCE_ROOT` is unset | gate inner/conformance | Go table test in a temp git repo with env unset | catches normal `go test` runs losing the default kit-root grading path |
| 6 | canary sweep invokes one baseline plus one inner gate per fixture and does not run full module tests per fixture | canary CLI + gate inner/conformance | fake runner count test for invocation count; wall-clock is not TDD-able and must be recorded before/after on the same machine | catches accidental nested sweeps mechanically, while the timing record covers the map's machine-dependent cost uncertainty |
| 7 | command help/reference/package surface names `bench canary`; no docs or shell fragments still describe `.bench/lib/canary-run.sh` as the runner implementation | scaffold compatibility | stale-reference/conformance sweep over command docs and package contents | catches a shipped surface where users can install the shim but cannot invoke the binary command |
| 7 | all existing `tests/canary/*` fixtures still bite through `bench canary`, including Go fixtures and docs/line/coverage anchors | canary CLI | `bench gate` canary layer; failures name the fixture | catches a port that implements the runner but weakens fixture attribution or skips fixture classes |

### Edge inventory

- **paths and directories containing spaces or glob characters** -- covered by
  canary copy tests and a root path with spaces; all gate invocations use argv, not
  shell string concatenation.
- **hand-edited files whose last line lacks a trailing newline** -- covered by
  `EXPECT` read tests that compare the exact substring without requiring a trailing
  newline.
- **absent file vs present-but-empty file** -- covered for `tests/canary` absent vs
  empty, `EXPECT` missing vs empty/vacuous, and `files/` missing vs empty
  directory.
- **unquoted multi-word arguments** -- covered by `bench canary [root]` tests using
  a spaced root argument; the shell wrapper passes `"$root"` as one arg.
- **required tool missing from PATH** -- covered by scaffold/shim failure posture:
  missing resolved bench command is gate-red with an explanatory error, not a silent
  skip. Go toolchain absence for a graded root with `go.mod` remains a hard red in
  the gate entry.
- **invocation through a symlink rather than the real path** -- safe by construction
  through the existing wrapper resolver and `bin/bench.sh` symlink resolution; no
  new symlink resolver is introduced.
- **dot-dir fixtures** -- covered by recursive `dot-*` segment restore before the
  inner gate runs.
- **interrupted/partial canary sweep** -- covered by temp-dir-per-fixture cleanup;
  no repo state is written except ordinary gate output.
- **re-run idempotency** -- covered by temp dirs created fresh per fixture and by
  second-init no-clobber.
- **cwd deeper than the repo root** -- covered by root defaulting through git root
  and by explicit root argument tests.
- **leaked wrapper env (`BENCH_KIT`, `BENCH_WRAPPER`)** -- covered by canary env
  tests that assert those internals are stripped before the inner gate.
- **`EXPECT` substring coupled to message text** -- covered by the existing fixture
  rule: updating a check message requires updating `EXPECT` in the same diff;
  vacuous baseline prevents empty/generic expectations.
- **Won't handle: adversarial same-diff deletion of both gate logic and its canary
  fixture** -- outside the ADR threat model; pinning the gate outside the writable
  tree is a separate capability.
- **Won't handle: porting individual conformance checks in this spec** -- those are
  `go-gate-conformance-checks`; this spec provides the runner and harness they
  attach to.
- **Won't handle: porting behavior contracts in this spec** -- those are
  `go-gate-behavior-contracts`; this spec only keeps their fixtures biting in inner
  mode, and retiring the fragments is that spec's job.

## Out of scope

- **Check-by-check conformance port** -- separate capability because each shell
  root-grading check must move with its targeted fixture proving the Go check bites
  before the shell twin retires -- `35 edits, 12 gate runs`.

- **Behavior-contract port** -- separate capability because runtime, shift, link,
  doctor, AXI, Go routing, and package contracts are assertion-for-assertion parity
  work, not canary machinery -- `40 edits, 12 gate runs`.

- **Deleting every `gate-*.sh` fragment and `gate-contract-runner.sh`** -- separate
  retirement step completed across the conformance and behavior specs; this spec
  deletes only the canary sweep logic from `.bench/lib/canary-run.sh` --
  `18 edits, 6 gate runs`.

- **Pinning the gate outside the writable working tree** -- separate
  adversarial-defense capability explicitly left out by the canary tripwire ADR; the
  current threat model makes weakening loud, not impossible -- `16 edits, 4 gate
  runs`.
