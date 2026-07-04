# go-gate-conformance-checks

Status: staged

## Problem

The gate still grades most kit conformance through shell in `.bench/gate.sh` and
root-scanning fragments. Those checks are the oracle for the kit content surface:
they catch stale command references, broken skill indexes, drifted line bindings,
missing package assets, invalid JSON, lost workflow anchors, and other defects that
do not show up as normal product behavior.

The parent map split gate content into three specs. `go-gate-canary-machinery`
supplies the canary runner and the root-parameterized Go conformance harness. This
spec is the middle slice: port the root-grading conformance checks into Go tests
without weakening the canary net. The hard part is not translating greps; it is
preserving attribution. A check may leave shell only after its existing fixture goes
red under the Go check with the targeted `EXPECT` substring, and the shell twin
deletes in that same diff.

## Solution

Add a Go conformance test suite that grades a root path, defaulting to the kit repo
and accepting the fixture root through `BENCH_CONFORMANCE_ROOT` from the canary
machinery spec. The suite ports the root-grading checks from `.bench/gate.sh`,
`.bench/gate-docs-contracts.sh`, `.bench/gate-line-contracts.sh`, and the root-only
parts of the package, Go, and guard-manifest fragments. Behavior contracts stay
shell until `go-gate-behavior-contracts`.

Each port moves one check family at a time. The Go check first produces the same
targeted substring for its canary fixture, then the shell twin is deleted in the
same diff. During the transition the gate may run both Go and not-yet-ported shell
families, but no individual check may be live in both once its fixture bites in Go.

## User stories

1. As the gate, I want a named Go conformance registry that classifies every
   canary fixture under `tests/canary/` -- assigned to a root-grading check family
   or exempted by name as behavior-owned -- so that the port cannot silently skip a
   shell check or orphan a fixture. Line: gpt-5.5 / xhigh. The map makes omission the central failure mode
   for this slice, and an unported check that still looks green weakens the oracle.

2. As the gate, I want the inline load and validity checks ported to Go tests:
   shell syntax for remaining shell files, executable git modes for command-path
   scripts and adapters, extensionless `.bench/gate`/`.bench/done` references,
   required JSON validity, Codex hook wiring, skill frontmatter, craft visible
   names, Claude skill mirroring, and shared-rule single-sourcing, so that the first
   root-grading layer no longer lives in bash. Line: gpt-5.5 / xhigh. These are
   gate-critical checks whose canary messages are precise today, and the map
   requires each one to bite under Go before its shell twin retires.

3. As the gate, I want kit-content sync checks ported: generated skills index
   equality, missing or stale skill index entries, command references in the
   operating guide, Codex adapter files and metadata, roadmap-promotion anchors, and
   command-adapter documentation, so that disk, docs, and harness adapters keep one
   source of truth after the shell checks delete. Line: gpt-5.5 / xhigh. The
   one-source-per-fact rule is the project standard, and a false green here lets
   guidance drift across every future session.

4. As the gate, I want doc-conformance and workflow-anchor checks ported: stale
   slash/Codex references, cold-pickup CLI command lists, BENCH-reference
   token-diet placement, README command-first layout, command Entry/Exit sections,
   acceptance-coverage anchors, edge-inventory anchors, oracle-authoring anchors,
   review/delegation anchors, spec-lifecycle anchors, spec-sourcing anchors, README
   file-layout anchors, skills-index generate/verify, and `bench coverage --check`
   over mapped specs, so that structural guidance contracts stay enforced by Go.
   Line: gpt-5.5 / xhigh. These checks grade prose the normal test suite cannot
   understand, and the map flags gate-criticality over speed for conformance logic.

5. As the gate, I want line-routing conformance ported: `.bench/lines.env` tier
   shape through the binary's own binding view, alias shape, profile prose sync,
   Claude Agent hook wiring, `check-agent-line.sh` allow/deny/degraded cases,
   adapter fail-closed/pass-through cases, incomplete-binding fallbacks, and the
   line prose anchors, so that invariant #2 enforcement keeps its bite after the
   shell fragment retires. Line: gpt-5.5 / xhigh. The line binding crosses hooks,
   adapters, Go, and project prose, and a false green silently disables the routing
   invariant the map treats as oracle-grade.

6. As the gate, I want root-only package, compiled-core, release-structure, and
   guard-manifest conformance ported: `package.json` `files[]` resolution, npm
   dry-run required/forbidden assets, repo-only package-claim sweep, Go
   fmt/build/vet/test/cross-compile checks where `go.mod` exists, release workflow
   structure assertions, and hook `--describe` manifest key checks, while leaving
   version routing, package-generator behavior, AXI, doctor, link, status, roadmap,
   postinstall, and shift behavior contracts to `go-gate-behavior-contracts`. Line: gpt-5.5 / xhigh. This row contains the map's ambiguous edge between root grading
   and behavior contracts, so it needs the strongest review line to avoid deleting
   shell coverage from the wrong side.

7. As a reviewer, I want the retire rule enforced per check: a fixture must fail
   under the Go conformance check with its targeted `EXPECT` substring before the
   matching shell assertion is deleted, and the deletion happens in the same diff,
   so that there is never a window with no biting check and never a lasting second
   implementation. Line: gpt-5.5 / xhigh. The map names this as the slice's core
   safety rule, and violating it turns the canary layer into a false assurance.

8. As the gate entry, I want `.bench/gate.sh` to keep only the entry work and the
   not-yet-ported behavior fragments after each conformance family retires, with no
   dangling shell source, stale `EXPECT` ownership, or duplicated check message, so
   that the strangler window keeps shrinking without changing the gate exit
   contract. Line: gpt-5.5 / xhigh. The gate is the oracle, and the map's
   shell-or-Go rule makes duplicate or dangling checks a gate-critical migration
   risk.

## Implementation decisions

- **Conformance lives in Go test code, not the shipped binary.** Put the ported
  root-grading checks in a test-only conformance package. Helper code should be
  `_test.go` unless a helper is already a production parser, such as
  `internal/coverage` or `internal/lines`. No `bench selfcheck` subcommand is
  introduced.

- **One root parameter.** The conformance helper resolves the graded root from
  `BENCH_CONFORMANCE_ROOT`, defaulting to the kit repo. That is the single root
  override shared with `go-gate-canary-machinery`; do not introduce a second env
  name.

- **Diagnostics are the contract.** Each check returns stable diagnostics whose text
  preserves the existing fixture `EXPECT` substring unless the fixture updates in the
  same diff. Ordering is deterministic so one failing fixture is attributable.

- **One check registry for retirement.** The Go suite carries a registry that
  classifies every fixture under `tests/canary/`: each is either assigned to a
  ported check family (with its shell twin location) or exempted by name as
  behavior-owned, retiring with `go-gate-behavior-contracts`. A registry test fails
  when any fixture is unclassified -- including fixtures added after this port --
  or when a retired fixture's targeted substring still appears in
  `.bench/gate*.sh`.

- **Inventory to port.**

| family | shell source today | root-grading canary fixtures |
|---|---|---|
| load/validity/metadata | `.bench/gate.sh` inline checks | `invalid-json`, `codex-hooks-broken`, `bad-frontmatter`, `claude-skills-unmirrored`, `extensionless-gate-ref`, `shared-rule-drift` |
| skills index and command adapter sync | `.bench/gate.sh`, `.bench/skills-index.sh` checks | `dangling-index`, `missing-index-field`, `stale-index-wording`, `unindexed-skill` |
| docs currency and token-diet placement | `.bench/gate-docs-contracts.sh` | `stale-command-reference`, `stale-codex-adapter-reference`, `stale-cli-doc-reference`, `historical-marker-prose`, `benchref-missing`, `benchref-pointer-dropped`, `benchref-imported`, `benchref-section-duplicated`, `readme-command-first` |
| workflow and guidance anchors | `.bench/gate-docs-contracts.sh`, `.bench/gate-line-contracts.sh` | `acceptance-coverage-anchor`, `coverage-axis-anchor`, `command-handoff-anchor`, `debug-archaeology-anchor`, `edge-inventory-anchor`, `implement-spec-status-flip-anchor`, `shape-idea-bypass`, `shape-idea-bypass-wrapped`, `shape-idea-handoff-anchor`, `story-line-anchor-missing`, `write-spec-handoff-anchor`, `write-spec-map-required`, `line-anchor-missing` |
| coverage-map validation | `.bench/gate-docs-contracts.sh` via `bench coverage --check` | `broken-coverage-map` |
| line routing | `.bench/gate-line-contracts.sh` | `line-binding-prose-drift`, `agent-hook-unwired`, `agent-hook-broken`, `adapter-line-broken` |
| package/core/guard root checks | `.bench/gate-package-contracts.sh`, root-only parts of `.bench/gate-go-contracts.sh`, guard-manifest preflight in `.bench/gate-axi-contracts.sh` | `missing-files-entry`, `go-build-broken`, `go-test-failing`, `guard-describe-boundary-dropped` |

- **Only root-grading fixtures are assigned here.** The fifteen fixtures outside
  the inventory table (`doctor-foreign-clobbered`, `doctor-manager-dir-chosen`,
  `doctor-stale-silent`, `postinstall-guard-bypassed`, `postinstall-nonzero-exit`,
  `session-start-advice-dropped`, `wrapper-args-dropped`, `status-regressed`,
  `roadmap-regressed`, `unscaffolded-bench-file`, `toon-escaping-dropped`,
  `learnings-parse-broken`, `guards-aggregation-dropped`,
  `coverage-extraction-dropped`, `diff-recorded-base-dropped`) bite through
  behavior-contract fragments and are exempted by name in the registry as
  behavior-owned; `go-gate-behavior-contracts` keeps them biting when those
  fragments retire. They are not conformance work here.

- **Reuse production parsers where they exist.** Coverage-map validation calls the
  Go coverage parser, line binding reads through `internal/lines`, and guard
  manifest checks use the hook `--describe` surface. Do not re-derive the same fact
  with a second ad hoc parser when the Go core already owns it.

- **Shell deletion is check-sized.** Each commit can port a family or sub-family,
  but the smallest unit is one fixture-backed check: Go diagnostic green on the real
  root, fixture red with targeted substring, shell twin deleted, registry updated.

## Testing decisions

- **What a good test is here.** Drive the conformance suite at a root path and
  observe diagnostics and exit status. For canary proof, point the Go test at the
  fixture root and assert the targeted substring. Do not test private helper
  internals except for pure parsers already exposed as package seams.

- **Seams.** Three seams get tested:
  - The **root-grading conformance seam**: a root path enters the Go checks and
    diagnostics come out.
  - The **retire-rule seam**: each fixture's `files/` tree plus `EXPECT` proves the
    Go check bites before the shell twin disappears.
  - The **gate-entry transition seam**: `.bench/gate.sh` still exits with the same
    `gate: green` / `gate: red` contract while conformance families leave shell and
    behavior fragments remain.

- **Gate command:** `bench gate`, plus focused `go test -count=1 ./...` while this
  slice is landing. Done means the project gate is green, every ported root-grading
  fixture bites through Go, and the shell twin sweep has no retired check left.

### Seam diagram

Root-grading conformance:

    trigger: gate runs `go test -count=1` or canary inner run points tests at a fixture root
        |
        v
    graded root path --> [ Go conformance registry + check families ] --> diagnostics, exit 0/1
    repo files       --> [ parsers and root scanners, reusing Go cores where they exist ]
                         ^
                         tests attach here: table/fixture tests set the graded root and assert exact
                         targeted substrings for each family in the inventory.

Retire rule:

    trigger: implementation ports one check family
        |
        v
    tests/canary/<fixture>/files --> [ canary fixture provision + Go conformance root ] --> red + EXPECT substring
    tests/canary/<fixture>/EXPECT--> [ registry maps fixture to check + shell twin      ]
                                      ^
                                      tests attach here: fixture fails before the shell twin deletes; a registry
                                      sweep fails if the retired shell message still appears in `.bench/gate*.sh`.

Gate-entry transition:

    trigger: reviewer or shift runs `bench gate`
        |
        v
    repo root --> [ thin .bench/gate.sh entry ] --> build/vet/go test, remaining shell behavior fragments, canary
              --> [ ported Go conformance tests ]
                  ^
                  tests attach here: normal gate run plus fixture-targeted Go tests verify exit/status
                  continuity while individual shell conformance checks retire.

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | every fixture under `tests/canary/` is classified: inventory fixtures assigned to exactly one Go check family, the fifteen behavior-owned fixtures exempted by name | root-grading conformance | registry test fails on the first unclassified fixture, including any fixture added after this port | prevents a cheap port that leaves canary coverage orphaned and forces every future fixture to declare an owner |
| 2 | inline load/validity fixtures `invalid-json`, `codex-hooks-broken`, `bad-frontmatter`, `claude-skills-unmirrored`, `extensionless-gate-ref`, and `shared-rule-drift` each red under Go with their current `EXPECT` substring | root-grading conformance + retire rule | fixture-targeted Go test for each listed fixture fails until the corresponding Go diagnostic exists | catches a port that validates the real root but no longer catches the planted regression |
| 3 | skills-index and command-adapter fixtures `dangling-index`, `missing-index-field`, `stale-index-wording`, and `unindexed-skill` each red under Go; shell index/adapter twins delete after bite | root-grading conformance + retire rule | fixture-targeted Go test fails before index equality and frontmatter validation are ported | catches drift between generated index, disk skills, and adapter metadata |
| 4 | docs currency/token/workflow fixtures listed in the docs and workflow rows of the inventory each red under Go, including wrapped bypass text and stale Codex adapter references | root-grading conformance + retire rule | fixture-targeted Go tests fail until stale-reference scanning, token-diet checks, README checks, Entry/Exit checks, and anchor greps are ported | catches a prose-surface regression that normal unit tests cannot grade |
| 4 | coverage-map validation fixture `broken-coverage-map` reds through the Go conformance suite using the same coverage parser as `bench coverage --check` | root-grading conformance | fixture-targeted Go test fails until the conformance suite calls the coverage parser and forwards the parser diagnostic | catches a second parser or missing validation that lets malformed acceptance maps through |
| 5 | line-routing fixtures `line-binding-prose-drift`, `agent-hook-unwired`, `agent-hook-broken`, and `adapter-line-broken` each red under Go with their current `EXPECT` substring | root-grading conformance + retire rule | fixture-targeted Go test fails until binding shape, profile sync, hook wiring, hook behavior, and adapter cases are ported | catches a port that preserves docs anchors but silently disables model-line enforcement |
| 6 | package/core/guard fixtures `missing-files-entry`, `go-build-broken`, `go-test-failing`, and `guard-describe-boundary-dropped` red under Go; non-root fixtures stay exempt as behavior-owned | root-grading conformance + retire rule | fixture-targeted Go tests fail until package files, Go build/test, and guard manifest checks move; registry fails if an exempt fixture is claimed by a conformance family | catches deleting root checks from shell while accidentally pulling non-conformance work into this spec |
| 7 | for each ported check, shell twin text is absent from `.bench/gate*.sh` in the same diff where the Go fixture bite is green | retire rule | retired-shell sweep fails while a ported fixture's `EXPECT` substring or shell check label still appears in shell fragments | enforces one source per check and prevents long-lived duplicate implementations |
| 8 | `.bench/gate.sh` still returns the same exit contract and sources only not-yet-ported behavior fragments after conformance retirement | gate-entry transition | `bench gate` plus a stale-reference sweep fails if a deleted fragment is still sourced or if the gate omits the Go test invocation | catches a strangler break where the entry stops grading conformance or references deleted shell |

### Edge inventory

- **paths and dirs with spaces or glob characters** -- covered by the conformance
  helper taking a root path as data, not interpolating shell strings; fixture roots
  can be provided under the space-path canary harness from
  `go-gate-canary-machinery`.
- **hand-edited files with no trailing newline** -- covered by reading whole files
  as bytes and splitting without requiring a final `\n`; `EXPECT` files are trimmed
  only for the comparison newline, not for meaningful substring content.
- **absent vs present-empty** -- covered: missing JSON differs from invalid JSON,
  absent optional surfaces keep today's guarded-skip posture, empty `lines.env`
  values remain malformed/incomplete, and coverage maps distinguish no-map from a
  malformed map.
- **malformed input** -- covered by fixtures for invalid JSON, malformed coverage
  rows, stale index entries, missing frontmatter/index fields, bad hook manifests,
  and malformed line bindings.
- **interrupted or partial state** -- safe by construction: conformance checks are
  read-only except subprocess probes that write to `t.TempDir`; cleanup uses
  `t.Cleanup`.
- **re-run idempotency** -- covered by deterministic diagnostics and sorted
  filesystem walks; running the suite twice against the same root yields the same
  set of messages.
- **hostile environment / missing tools** -- covered for required Go toolchain when
  `go.mod` exists; shellcheck remains best-effort if retained by the entry, and
  npm/go subprocess checks report a diagnostic rather than silently skipping when
  the checked surface exists.
- **cwd deeper than repo root** -- covered by passing the graded root explicitly;
  tests do not infer root from process cwd once the env root is set.
- **symlink/dot-dir fixtures** -- covered by the canary fixture provision and
  `dot-` restore; this spec consumes the restored root only.
- **Won't handle: semantic quality of guidance prose** -- the gate checks
  structural anchors and stale references; review still judges whether the prose is
  good.
- **Won't handle: optional shellcheck policy changes** -- this spec preserves the
  current best-effort posture if it remains; deciding to make shellcheck hard or
  remove it is a separate gate-policy change.

## Out of scope

- **Canary machinery and root-parameterized harness** -- prerequisite, not this
  capability; it owns `bench canary`, the shim, fixture provisioning, inner-run
  mode, vacuous `EXPECT`, and absent/empty harness checks. Estimate:
  `18-25 edits, 6-8 gate runs`.

- **Behavior contract fragments** -- separate capability covering runtime/shift/
  link/doctor/AXI/go-routing/package behavior contracts by assertion-for-assertion
  parity, with no canary fixture ownership in this spec. Estimate:
  `35-55 edits, 10-16 gate runs`.

- **New conformance rules beyond the port** -- separate future checks need their own
  fixtures and review of the policy, not a hidden expansion while migrating shell to
  Go. Estimate per small family: `4-8 edits, 2-3 gate runs`.

- **Changing shellcheck from best-effort to hard enforcement** -- a gate-policy
  decision, not part of preserving current root-grading behavior. Estimate:
  `2-4 edits, 1-2 gate runs`.
