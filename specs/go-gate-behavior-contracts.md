# go-gate-behavior-contracts

Status: implemented

## Problem

The Go rewrite has moved the kit's operational behavior into the binary, but the
gate still proves that behavior through large sourced shell fragments and
`gate-contract-runner.sh`. Those fragments are now test code, not product code, yet
they still carry shell fixture setup, implicit `$root`/`$tmp` state, heredoc
`eval`, and one-off assertion styles.

This is the last behavior-contract step in `decisions/go-gate-port.md`. Canary
machinery and root-grading conformance checks are prerequisites from the prior two
gate specs. This spec ports only the behavior contract fragments: runtime,
runtime-git, shift, link, doctor, AXI core, AXI guards, AXI wave 2, Go routing, and
package/installable-surface behavior. Fifteen canary fixtures bite through these
fragments (seven doctor/postinstall, two runtime, two AXI core, one AXI guards, two
AXI wave 2, one link), so the retire rule is twofold: assertion-for-assertion
parity for every row before deleting each shell twin, and for the fixture-covered
families each fixture must bite through the ported Go tests in the same diff the
fragment retires.

## Solution

Create a shared Go contract helper package for integration tests that replaces
`gate-contract-runner.sh`: fixture provisioning, optional space-containing paths,
optional no-repo fixtures, git setup, command execution, output/exit capture, and
cleanup live in one Go helper. One Go test file mirrors each former behavior
fragment and drives the same public seam the shell did: `bin/bench.sh`, hook shims,
generated wrappers, fixture repos, npm pack/generator behavior, and fabricated kit
layouts. For the fixture-covered families, the retiring diff also swaps the
inner-mode gate from the shell fragment to a root-parameterized run of the ported
Go tests (filtered or precompiled, per the canary spec's cost rule), so the fifteen
behavior-targeted fixtures never stop biting.

Each former shell assertion gets a named Go subtest or table row. The implementation
may group setup more cleanly, but review must be able to map every old contract
label to its Go assertion before the shell fragment retires. A shell fragment is
deleted in the same diff that its Go parity tests pass. `gate-contract-runner.sh`
deletes last, after no behavior fragment sources it.

After this spec, behavior contracts run under `go test -count=1`. `.bench/gate.sh`
no longer sources behavior fragments; it keeps the stable gate entry and delegates
to the already-ported gate stack from the canary and conformance specs.

## User stories

1. As the gate maintainer, I want a shared Go contract helper package to replace
   `gate-contract-runner.sh` with the same fixture contract: fresh temp fixture,
   optional parent path containing a space, optional no-repo mode, git init when
   requested, kit-root discovery, wrapper execution, stdout/stderr/exit capture,
   environment isolation, and cleanup on success and failure, so that every ported
   behavior contract uses one provisioning/reporting source. Line: gpt-5.5 / xhigh. The map calls behavior parity gate-blind, and a helper bug would make
   every migrated assertion look green while the old shell runner is gone.

2. As a reviewer of the runtime surface, I want `gate-runtime-contracts.sh` ported
   into Go integration tests: idea/roadmap, gate cwd and `BENCH_GATE` cwd, gate
   resolution order, status clean/footer/stale/fresh/decision/budget/warm-pool rows,
   gate-cache write and no-cache fail-safes, missing-core fail-safes, retirement and
   learnings signals, structure budgets and path-with-spaces checks, worktree
   lease/reuse/hardening/concurrency, symlinked kit-dir portability,
   `stop_hook_active`, and missing-bench fail-open, so that the runtime contract no
   longer depends on sourced shell. Line: gpt-5.5 / xhigh. Only `status-regressed`
   and `roadmap-regressed` have canary bite-proof here, so the dominant risk is
   parity drift that the gate cannot see once the shell twin is deleted.

3. As the shift-loop maintainer, I want `gate-runtime-shift-contracts.sh` ported
   into Go integration tests: green gated loop commits and benchBase, worktree
   isolation, touched-path staging, red rollback, commit-failure reporting,
   touched-scope refactor, no-op refactor, interrupt cleanup, gate-interrupt
   cleanup, `.bench/done.sh` early completion, scratch survival, refactor prompt
   scope, adapter preflight, adapter single-argument invocation, and reference
   adapter file shape, so that the shift behavior net runs in `go test`. Line: gpt-5.5 / xhigh. Shift has the map's most expensive gate-blind parity risk
   because a missing or weakened assertion can allow commit-on-red, leaked leases,
   or split prompts without a canary catching the deleted shell contract.

4. As the destructive-git guard maintainer, I want
   `gate-runtime-git-contracts.sh` ported into Go integration tests with the full
   allow/block matrix and its exit and `BLOCKED:` output checks intact, so that the
   guard's behavior contract is tested through the hook seam without the shell
   fragment. Line: gpt-5.5 / xhigh. The map shorthand calls this runtime behavior,
   and losing one matrix case during the port is gate-blind once the shell matrix is
   retired.

5. As the adoption-surface maintainer, I want the non-canary behavior in
   `gate-link-contracts.sh` ported into Go integration tests: safe fresh link,
   relink, existing `AGENTS.md`, malformed markers, conflict without manifest,
   modified managed files, linked local CLI hooks, metachar kit paths, linked
   worktrees, `core.hooksPath`, default-branch fallback, fenced marker examples,
   hooksPath conflict, and unclosed-fence refusal, so that safe link parity survives
   without the shell fragment. Line: gpt-5.5 / xhigh. The canary scaffolding rows
   belong to the canary machinery spec, and the remaining link rows are, except the
   fixture-backed `unscaffolded-bench-file`, gate-blind parity checks over
   reviewer-owned file mutation.

6. As the doctor/postinstall maintainer, I want `gate-doctor-contracts.sh` ported
   into Go integration tests: doctor report states, manifest skew, `--fix` write,
   spaced-target quoting, idempotency, foreign-file refusal, fallback dir selection,
   PATH notice, stale generated shim, argument passthrough, postinstall fail-open
   behavior, and session-start advisory, so that PATH-shim behavior is preserved in
   Go tests. Line: gpt-5.5 / xhigh. Doctor reads and writes user PATH state, and
   seven of the fifteen behavior-targeted fixtures bite through this fragment, so
   retirement here is graded by both parity review and fixture bite.

7. As the AXI query-surface maintainer, I want `gate-axi-contracts.sh` ported into
   Go integration tests: learnings rows, empty/template output, maps unresolved rows,
   TOON escaping, usage exit 2, subdirectory root resolution, path-with-spaces, maps
   parser hardening, CRLF stripping, no-Type tickets, ASCII-separator learnings
   titles, Handoff close-readiness, and maps `--count`, so that the core AXI surface
   keeps its hybrid stdout/exit contract. Line: gpt-5.5 / xhigh. AXI output is
   machine-read by agents, and a missing parity row can silently change the protocol
   after the shell assertion disappears.

8. As the guard-aggregation maintainer, I want
   `gate-axi-guards-contracts.sh` ported into Go integration tests: guards
   aggregation, `--brief`, usage/subdirectory, path-with-spaces, `--describe`
   timeout bound, unmanaged pre-push safety, core-unreachable manifest,
   linked-worktree classification, session-start guard injection, and
   never-blocks-outside-repo, so that guard advertisement and hook injection stay
   contract-tested. Line: gpt-5.5 / xhigh. This surface joins enforcement and
   advertisement, and the map's gate-blind parity warning applies directly to any
   row that is dropped during the port.

9. As the wave-2 AXI maintainer, I want `gate-axi-wave2-contracts.sh` ported into
   Go integration tests: diff recorded-base, diff fallback/shape, diff error
   posture, coverage extraction, coverage state/error, and coverage `--check`
   validation, so that review-base and acceptance-coverage behavior no longer
   depend on shell fixtures. Line: gpt-5.5 / xhigh. These contracts protect later
   workflow phases, and parity drift is gate-blind because the old shell fragment is
   the only current executable oracle for several edge cases.

10. As the packaging and router maintainer, I want the behavior portions of
    `gate-go-contracts.sh` and the installable-surface contracts ported into Go
    integration tests: `bench version` output and outside-repo execution,
    fabricated version-routing precedence, missing/empty/non-executable binary
    rims, symlink invocation, off-matrix platform errors, linked-worktree binary
    resolution, platform-package generator idempotence and output, and npm pack
    behavior not already moved as root conformance, so that the installable surface
    is tested under `go test`. Line: gpt-5.5 / xhigh. The compiled-core and
    root-conformance portions are prerequisites, while the named behavior rows are
    canary-blind parity risks if their shell fragment is deleted first.

11. As the gate maintainer, I want each ported behavior fragment deleted only after
    its Go parity tests pass and its behavior-targeted fixtures (where the family
    has them) bite through the ported Go tests in the same diff, `.bench/gate.sh`
    stopped from sourcing those fragments, `gate-contract-runner.sh` deleted after
    the last user is gone, all dangling references fixed in the same change, and
    `bench gate` still emitting `gate: green` / `gate: red` at the final verdict
    and exiting 3 outside a git repo, so that the gate has one behavior-contract
    runner and no stale shell test surface. Line: gpt-5.5 / xhigh. Retirement is
    assertion-for-assertion and review-graded except where a fixture bites, and the
    Handoff keeps the gate exit and stdout contract across each flip, so deletion
    discipline is the control that keeps the migration honest.

## Implementation decisions

- **Helper package path.** Add `internal/contract` as the shared contract-test
  package. It is Go test infrastructure in the kit module, not a binary subcommand
  and not shipped product behavior.

- **Test file layout mirrors the old fragments.** Use one Go test file per retired
  behavior fragment: `runtime_test.go`, `runtime_git_test.go`,
  `runtime_shift_test.go`, `link_test.go`, `doctor_test.go`, `axi_test.go`,
  `axi_guards_test.go`, `axi_wave2_test.go`, `go_routing_test.go`, and
  `package_test.go`. This keeps review parity mechanical.

- **Subtest names preserve old labels.** Each former shell contract label becomes a
  Go `t.Run` name or table label. When a shell fragment had multiple inline
  assertions under one `err` label, the Go test may split them into clearer
  subtests, but the old label must appear once as the group name.

- **The helper drives the public seam.** Tests execute `bash <kit>/bin/bench.sh`,
  hook shims, generated shims, git, npm, and fixture gate scripts as subprocesses.
  They do not call internal production functions to prove behavior contracts. Unit
  tests that already exist beside production packages remain separate.

- **The helper owns fixture options.** It provides the Go equivalents of
  `--space-path` and `--no-repo`, plus common helpers for git init, commits, writing
  executable files, checking output substrings, checking exact lines, and reading
  filesystem results. It sets isolated `HOME`, `XDG_*`, `PATH`, npm cache, and Bench
  env as each contract requires.

- **Fixture-covered families stay root-parameterized.** The contract helper
  resolves the graded kit root through the same `BENCH_CONFORMANCE_ROOT` override
  the conformance harness uses, defaulting to the real kit root. When a family with
  behavior-targeted fixtures ports (runtime, link, doctor, AXI core, AXI guards,
  AXI wave 2), the retiring diff points the inner-mode gate at its Go tests so
  `bench canary` still reds with each fixture's targeted `EXPECT`. Families with no
  fixtures (runtime-git, shift, Go routing, package) need no inner-mode wiring.

- **No heredoc-shell transplant.** Port assertions to Go fixtures and Go subprocess
  calls. Short fixture scripts are fine when the product behavior itself requires a
  shell script, but wrapping old shell contract bodies in `bash -c` would keep the
  shell runner in disguise.

- **Binary build source stays single.** The contract helper uses the existing
  `scripts/go-build.sh` when it needs a fresh `dist/bench`; it must not duplicate
  ldflags, platform mapping, version stamping, or binary-resolution knowledge.

- **Retirement is family-scoped.** A fragment can retire when its whole Go test file
  has assertion parity and passes. `gate-contract-runner.sh` retires only after all
  fragment sources are gone. There is no long-lived period where the same behavior is
  asserted in both shell and Go.

- **Prerequisite split is respected.** Canary machinery and root-grading
  conformance checks are not redesigned here. If one of the named behavior fragments
  still contains prerequisite-owned checks when this spec starts, move those checks
  under the prerequisite spec's result first, then port only the remaining behavior
  rows here.

## Testing decisions

- **What a good test is here.** A good test drives the same public entry the shell
  contract drove, asserts exit code, stdout/stderr, filesystem, and git state, and
  fails with a label that maps to the retired shell contract. It tests observable
  behavior, not internal implementation.

- **Seams.**
  - The **Go behavior contract seam**: `internal/contract` provisions a fixture and
    runs the real CLI/hook/package surface as a subprocess.
  - The **gate retirement seam**: `.bench/gate.sh` reaches behavior contracts only
    through `go test -count=1`, with no sourced behavior fragments or shell runner.

- **Gate command:** `bench gate`. The focused implementation loop should also run
  `go test -count=1 ./internal/contract` while porting each family, then the full
  gate before deleting each shell fragment.

### Seam diagram

Go behavior contract seam:

    trigger: go test runs a former-fragment test file
        |
        v
    test options --> [ internal/contract fixture helper ] --> temp repo / no-repo / space path
    kit root     --> [ bash bin/bench.sh or hook shim     ] --> stdout, stderr, exit code
    fixture files--> [ real git/npm/filesystem behavior   ] --> resulting files, commits, cache, manifests
                    ^
                    tests attach here: Go assertions named after old shell contract labels.

Gate retirement seam:

    trigger: bench gate
        |
        v
    .bench/gate.sh --> [ build / vet / go test -count=1 from prior specs ] --> Go behavior contract tests
    old fragments  --> [ no sourced gate-*-behavior fragments remain     ] --> no gate-contract-runner.sh
                      ^
                      tests attach here: gate load, stale-reference sweep, gate stdout checks, and the Go contract package.

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | helper provisions git repo, no-repo, and space-path fixtures; captures stdout/stderr/exit; cleans up; isolates env | Go behavior contract | new `internal/contract` helper tests fail to compile before the helper exists | catches a helper that cannot replace `gate-contract-runner.sh` faithfully |
| 2 | every `gate-runtime-contracts.sh` label has a Go subtest and preserves the same runtime assertion | Go behavior contract | not TDD-able as product behavior: shell contracts are already green; parity is review-graded by old-label-to-Go-subtest mapping before deletion | catches dropped runtime assertions during shell fragment retirement |
| 3 | every `gate-runtime-shift-contracts.sh` label has a Go subtest, including interrupt and adapter cases | Go behavior contract | not TDD-able as product behavior: current shell fragment is the green parity net; Go subtests become the net before shell deletion | catches shift-loop regressions that would otherwise lose their only behavior assertion |
| 4 | full destructive-git allow/block matrix ports from `gate-runtime-git-contracts.sh` | Go behavior contract | not TDD-able as product behavior: existing matrix passes; review requires each matrix row in Go before deleting the shell matrix | catches missing allow/block cases and fail-closed drift |
| 5 | non-canary `gate-link-contracts.sh` rows port with safe-link mutation and conflict assertions intact | Go behavior contract | not TDD-able as product behavior: existing link shell rows are green; Go parity rows must exist before fragment deletion | catches link preflight or project-owned file mutation drift |
| 6 | doctor, generated shim, postinstall, and session-start advisory rows port from `gate-doctor-contracts.sh` | Go behavior contract | not TDD-able as product behavior: existing shell rows are green; Go parity rows must replace them label-for-label | catches PATH shim state, skew, quoting, fallback, and fail-open regressions |
| 7 | AXI core rows port from `gate-axi-contracts.sh` with stdout, exit-code, TOON, parser, and count assertions intact | Go behavior contract | not TDD-able as product behavior: existing shell rows are green; Go parity rows become the protocol pin before deletion | catches machine-output drift in learnings/maps/TOON behavior |
| 8 | AXI guards rows port from `gate-axi-guards-contracts.sh`, including timeout, unmanaged pre-push, and session-start rows | Go behavior contract | not TDD-able as product behavior: shell rows are current oracle; parity mapping is required before deletion | catches enforcement/advertisement drift and missing guard injection |
| 9 | AXI wave-2 rows port from `gate-axi-wave2-contracts.sh`, including diff and coverage validation cases | Go behavior contract | not TDD-able as product behavior: shell rows already pass; Go rows must mirror each label before retirement | catches review-base and acceptance-coverage parser drift |
| 10 | Go routing and installable-surface behavior rows port while compiled-core and root conformance remain prerequisite-owned | Go behavior contract | mixed: routing/package behavior is currently covered by shell rows; helper/package tests fail until Go assertions call npm, fabricated layouts, and wrapper routing correctly | catches binary-resolution, package-generator, npm-pack, and repo-only package-claim regressions |
| 11 | behavior fragments and `gate-contract-runner.sh` delete only after Go parity; `.bench/gate.sh` no longer sources them; references are fixed | gate retirement | gate load or stale-reference sweep goes red if a deleted fragment is still sourced or named | catches half-retired shell behavior tests and dangling references |
| 11 | each fixture-covered family retires with its behavior-targeted fixtures biting through the root-parameterized Go tests in the same diff | gate retirement | `bench gate` canary layer reds with `did not bite` naming the fixture if a retiring diff drops its bite | catches a family flip that orphans the fifteen behavior-targeted fixtures |
| 11 | after the last behavior fragment retires, `bench gate` still prints `gate: green` on a passing tree and `gate: red` on a failing tree | gate retirement | focused gate-entry test runs one green and one red root after fragment deletion and fails if stdout/stderr lacks the expected verdict line | catches the final shrink accidentally dropping the Handoff's preserved gate stdout contract |
| 11 | after the last behavior fragment retires, `.bench/gate.sh` still exits 3 outside a git repo | gate retirement | focused gate-entry test runs the gate from a non-repo cwd and fails unless the exit code is 3 | catches the final shrink dropping the map's carried outside-repo exit contract |

### Edge inventory

- **paths and directory names containing spaces or glob characters** -- covered by
  the helper `SpacePath` option and existing rows in link metachar paths, doctor
  spaced targets, AXI path-with-spaces, structure path-with-spaces, and shift
  touched-path staging.
- **hand-edited files whose last line lacks a trailing newline** -- covered by the
  existing link marker, roadmap/idea, structure budget, manifest, and parser rows
  that must port with parity.
- **absent file vs present-but-empty file** -- covered by no-gate/runtime rows,
  doctor missing/empty shim behavior, manifest rows, and Go routing's empty binary
  case.
- **unquoted multi-word arguments** -- covered by shift objective assembly, adapter
  single-argument prompt, and doctor shim arg-passthrough.
- **required tool missing from PATH** -- covered where the old behavior contract
  covers it: missing core binary, missing bench wrapper, binary-unreachable guard
  behavior, and no-repo command posture. Missing `git`, `go`, `node`, or `npm` for
  the contract runner itself is a hard prerequisite failure, not product behavior.
- **invocation through a symlink** -- covered by version-routing symlink invocation
  and symlinked kit-dir portability.
- **interrupt or partial state** -- covered by shift interrupt cleanup and
  gate-interrupt cleanup rows.
- **re-run idempotency** -- covered by relink, doctor `--fix` idempotency, worktree
  reuse, package generator idempotency, and repeated fixture setup.
- **cwd deeper than repo root** -- covered by AXI subdirectory root-resolution, gate
  repo-root cwd, `BENCH_GATE` cwd, and runtime root resolution rows.
- **Won't handle: adding new behavior coverage beyond the old shell fragments** --
  safe to skip here because this spec is a parity migration; new product contracts
  need their own map/spec so parity review stays mechanical.
- **Won't handle: true live npm publish or GitHub release execution** -- safe to
  skip because the package behavior contract checks generated packages, npm dry-run
  contents, and workflow structure, while actual release execution is a separate
  release-smoke capability.
- **Won't handle: new behavior-contract canary fixtures** -- the fifteen existing
  behavior-targeted fixtures keep biting through the ported Go tests, but adding
  bite-proof for currently canary-blind rows is a separate gate-strengthening
  feature.

## Out of scope

- **Canary machinery port** -- separate prerequisite: `bench canary`, the lib shim,
  scaffold updates, fixture format, and inner-run mode. Estimate to build
  separately: `26 edits, 8 gate runs`.

- **Root-grading conformance check port** -- separate prerequisite: gate inline
  checks plus docs, line, coverage, anchor, package/core root checks, and
  root-parameterized conformance tests. Estimate to build separately:
  `34 edits, 10 gate runs`.

- **New behavior contracts beyond shell parity** -- a separate gate-strengthening
  capability because this spec's review surface is assertion-for-assertion migration,
  not contract expansion. Estimate to build separately: `10 edits, 3 gate runs`.

- **Behavior canaries for these contracts** -- a distinct capability that would add
  bite-proof fixtures for currently canary-blind behavior rows. Estimate to build
  separately: `48 edits, 14 gate runs`.

- **Live release or global npm install smoke** -- a separate release validation
  capability, not part of moving package behavior contracts from shell to Go.
  Estimate to build separately: `6 edits, 2 gate runs`.
