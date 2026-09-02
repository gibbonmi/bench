# Roadmap light-path fixes, batch 3

Status: staged

Roadmap: FT276, FT279, FT285, FT202, FT264, FT201

Decision source: the six named roadmap rows, each carrying its "Decided 2026-09-02" paragraph from the reviewer's grill of that date

Verification log: 2 iteration(s) to accept — the first round blocked on a fourth census copy, an unreachable fixture, a headroom chain, a wrong cache seam, a swallowed error outside the fence, and an ungraded owner; the second pass found three more derivation copies, folded into the owner

## Problem

Six roadmap rows carry a decided fix that fits one ticket each. Three
derivations of the canonical repository path drift apart because no package
below them owns one. `bench link` installs a launcher in the kit source
checkout that the shim then prefers. A wrapped decision-map field with a colon
reaches the wrong refusal. Three purity-census harnesses repeat one scan with
no bite test. Three callers run a Go build unlocked when the cache lock fails.

Two production signal registrations miss a cancel signal, and a builder group
survives its owner's death.

Shipped one at a time, every row pays a full gate. Shipped together under one
landing, they pay one.

## Solution

Keep one reviewed umbrella for scheduling and retirement. Split every row at
its executable seam. Each ticket owns one small behavior, focused evidence, and
an exact write fence. Every ticket first verifies its row's premise against the
code it names. The ticket graph carries each `Writes:` overlap as a blocker
edge, and the rest of the frontier runs in parallel.

This spec builds after the `roadmap-light-path-fixes-2` landing, so its base is
that landing's commit or a descendant.

## User stories

### Group A — one canonical path owner

Line: opus / low. The three bodies are equivalent, and two tests already pin
the symlink case.

1. As a maintainer, I want one dependency-free package to derive the canonical repository path, so that three copies cannot drift.
2. As a worktree caller, I want the worktree, preflight, and canary derivations to call that owner, so that a symlinked root resolves the same everywhere.
3. As a maintainer, I want the new package to carry the same purity census the policy packages carry, so that it stays a leaf.
4. As a preflight caller, I want a path that does not exist yet to keep its absolute spelling, so that a not-yet-created destination still resolves.

### Group B — link refuses the kit checkout

Line: opus / low. The predicate exists, and the refusal shape is one line.

5. As a kit maintainer, I want `bench link` to exit 1 in the kit source checkout, so that no tracked launcher copy shadows the shim.
6. As a consumer, I want `bench link` in an already-linked consumer repo to stay green, so that a relink still works.
7. As a caller outside a git repository, I want the git-repository refusal to come first, so that the kit predicate never reads an empty root.
8. As a cold session, I want the project profile to state the rule in one sentence, so that the refusal is not a surprise.

### Group C — the spaced field name

Line: opus / low. The guard is one branch between two existing ones.

9. As a map author, I want a continuation with a spaced field name to get the one-physical-line refusal, so that the message names the wrap.
10. As a map author, I want an unknown field name with no space to keep the unexpected-field message, so that a typo stays a typo.
11. As a gate operator, I want a canary fixture for the colon-bearing wrapped line, so that the guard is graded.
12. As a maintainer, I want the decision-map fixture inventory to match the fixture directories, so that a stale count cannot hide a missing fixture.

### Group D — one purity census

Line: opus / medium. The helper is new, and the census has no bite test today.

13. As a maintainer, I want one shared helper to own the purity scan, so that three harnesses do not repeat one policy.
14. As a policy-package author, I want a one-line wrapper in each package, so that the scan keeps the package's own directory.
15. As a gate operator, I want the helper to red a forbidden import, an ambient effect, and a `t.Parallel` call, so that it is graded.
16. As a gate operator, I want each wrapper to prove the scanned set holds the package's own source file, so that a wrong-directory helper reds.
17. As a reviewer, I want a process-backed git fixture to count as an ambient effect, so that a shell call buys no exemption.
18. As a maintainer, I want the helper's own package scanned under the same policy, so that it is not a blind spot.

### Group E — the build cache

Line: opus / low. Three call sites and one validator, all gate-covered.

19. As a gate operator, I want a `gocache.Hold` error to refuse the gate, the lane, and the focused run, so that no build runs unlocked.
20. As a gate operator, I want the refusal to print the hold error and the cache path, so that the next action is visible.
21. As a caller, I want an empty `GOCACHE` to stay absent, so that the home derivation still applies.
22. As a caller, I want a relative `GOCACHE` to refuse and print the value, so that a cache path cannot move with the working directory.
23. As a caller, I want the applied cache directory absolute and writable before a build, so that a bad cache refuses early.

### Group F — cancel signals and the builder group

Line: opus / medium for the check, opus / low for the attribute. The check is a
new AST sweep at a known seam, and the attribute needs a platform split.

24. As a maintainer, I want a dev-tier check to red a production signal registration that does not use `subprocess.CancelSignals`, so that one source stays one.
25. As a maintainer, I want the two non-conforming registrations migrated, so that the check lands green.
26. As a gate operator, I want the check to inspect call sites, not file bytes, so that a comment token does not red it.
27. As a test author, I want test files excluded from the check, so that a fixture can register any signal.
28. As a gate operator, I want a builder child to die when SIGKILL takes its owner on Linux, so that no orphaned group survives.
29. As a release maintainer, I want the darwin build to keep compiling, so that the Linux-only attribute lives behind a build tag.

## Implementation decisions

- A new package `internal/canonicalpath` exports one function that returns a path's absolute, symlink-resolved, cleaned form. When the path does not exist, it keeps the absolute spelling. The six existing production functions become one-line wrappers, and the package carries a purity census. The owner check grades production files only.
- `bench link` reads `kitSourceCheckout` after the git-root refusal and before the plan builds. The refusal uses the adopt `toon.Errorf` shape and exit 1. The rule sentence sits under the profile's cold-session notes.
- The Sources guard tests the field name alone for a space, between the no-separator branch and the unknown-field branch. It reuses the `expected` slice; it does not restate the field names.
- A new package `internal/puritycensus` owns the scan policy: forbidden-import patterns, the ambient-effect set, the `t.Parallel` ban, and the self-exempt file name. `exec.Command` and `exec.CommandContext` stay in the ambient set, so a process-backed git fixture counts. Each policy package keeps a one-line wrapper. The helper scans its own directory under the same policy. The census scope stays the three policy packages and the two new leaf packages.
- `leasedRepo` stays as declared residue on FT202's successor row, if any.
- A `Hold` error refuses at all three call sites with `toon.Errorf` and exit 1, printing the error and the sanitized cache path. `FromEnv` refuses a relative inbound `GOCACHE` value, because `Apply` never passes the inbound value to a child. An unwritable derived directory fails the lock-file open inside `Hold`, so no second check exists. A missing directory is created by `HoldDir`, as today; the "exists" leg is satisfied by that create. The gate subject closure propagates an `Apply` error.
- A `Writes:` file that preflight adds for fixture or registry closure is authorization headroom. Only a file that a ticket's `What to build` names creates a blocker edge between tickets.
- The cancel-signal check mirrors the bounds-policy walk: production files only, the owning package exempt, an AST call-site inspection. It registers at the dev tier with Go-source inputs. `internal/skillsindex` and `internal/worktree/subshell.go` migrate to `subprocess.CancelSignals`.
- `Pdeathsig` is set on the builder child in `internal/runbinary` only, in a Linux build-tagged file with darwin and other legs, mirroring the release-evidence exchange files. `Setpgid` stays.

## Testing decisions

- Path derivation uses the two existing symlink tests plus a new package test.
- The link refusal uses one new adopt test and the two existing link journeys as the over-broad guard.
- The Sources guard uses the record-shape table and one new canary fixture.
- The purity census uses a helper-level bite test over in-memory sources and a wrapper-level scanned-set assertion.
- The cache refusals use one new refusal test per call site with an unwritable cache directory.
- The signal check uses a canary fixture with the token in a call, a comment, and a string.
- `Pdeathsig` uses the parking-builder fixture with a SIGKILL to the owner.

### Seam diagram

    roadmap row
        |
        v
    focused ticket -> package, CLI, or check seam -> focused evidence
        |
        v
    retained integration source -> serial gate -> reviewed landing

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| LQ1 | 1, 2 | a symlinked root resolves to the ledger's canonical spelling through the shared owner | `TestGatherAssignmentTarget` and `TestRestoreMutationFixtureRefusesSymlinkedRootSpelling` | a `Clean`-only owner leaves the link spelling |
| LQ2 | 4 | a destination path that does not exist resolves to its absolute spelling | `TestRestoreMutationFixtureReinstatesBaseAndRemovesOverlay` | an owner that propagates the symlink error breaks the restore |
| LQ3 | 3 | the new package imports nothing under `internal/` | `TestPurePackageSourceCensus` in the new package | an internal import reds the census |
| LQ4 | 5 | `bench link` in the kit source checkout exits 1 and names the kit checkout | new adopt test with `BENCH_KIT` at the repo root | the old path installs the launcher and exits 0 |
| LQ5 | 6 | `link copy` in an already-linked consumer exits 0 | `TestWrapperInstallFreshnessAndReloadJourneys` | a marker-based predicate refuses the relink |
| LQ6 | 7 | outside a git repository the git-repository message prints first | `TestLinkOutsideGitRepoNamesGitRepository` | a predicate before the root check reads an empty root |
| LQ7 | 8 | the profile's cold-session notes hold the rule sentence | review-owned: the sentence is present under that heading | an omitted sentence leaves the rule unstated |
| LQ8 | 9 | a continuation whose field name holds a space yields the one-physical-line message | `TestMapSourcesRequireExactRecordShape` new row | a guard after the unknown-field branch is unreachable |
| LQ9 | 10 | an unknown field name with no space yields the unexpected-field message | `TestMapSourcesRequireExactRecordShape` rows unknown field and second path | an over-broad guard rewrites those rows |
| LQ10 | 11 | the colon-bearing wrapped fixture bites | new fixture under `tests/canary/decision-map-integrity` | a fixture whose mutation misses the guard passes silently |
| LQ11 | 12 | the fixture inventory names every fixture directory and the count matches | `TestDecisionMapIntegrityFixtureInventoryRejectsDeletion` | a hand count that lags the tree hides a missing fixture |
| LQ12 | 13, 14 | each policy package's census is a one-line wrapper over the helper with its own directory | review-owned: the three wrappers name `"."` | a helper that resolves the directory itself scans one tree three times |
| LQ13 | 15 | in-memory sources with a forbidden import, an ambient effect, and a `t.Parallel` call yield three diagnostics with file and line | new helper test in `internal/puritycensus` | a helper that returns no diagnostics stays green |
| LQ14 | 16 | each wrapper asserts the scanned set holds the package's own source file | the three wrappers | a wrong-directory helper reds on the absent file |
| LQ15 | 17 | a source that calls `exec.Command` reds the census | new helper test | an ambient set without exec lets a git fixture through |
| LQ16 | 18 | the helper package's own census passes under the same policy | `TestPurePackageSourceCensus` in the helper package | an unscanned helper is a blind spot |
| LQ17 | 19, 20 | an unwritable cache directory refuses the gate run, the lane, and the focused run with the error and the path | new refusal test per call site | an `err == nil` guard runs the build unlocked |
| LQ18 | 21 | `GOCACHE=` falls through to the home derivation | `TestFromEnvFallsBackToTheHomeDerivation` | a refusal on empty breaks the derivation |
| LQ19 | 22 | `FromEnv` with `GOCACHE=cache` returns a refusal that names `cache`, and the footprint report prints it | new `FromEnv` test row | a verbatim return reports a path that moves with the caller |
| LQ20 | 23 | an absent derived cache directory is created by `HoldDir`, and an unwritable one fails the lock-file open | `TestHolderCreatesTheDirectoryAndTheLockFile` and the LQ17 refusal tests | a stat-based check refuses every first run |
| LQ21 | 24, 25 | the live tree passes the cancel-signal check after the two migrations | the new check over the live root | an unmigrated registration reds |
| LQ22 | 26 | a fixture with the token in a call reds and the same token in a comment and a string does not | new fixture under `tests/canary` | a byte grep reds the comment |
| LQ23 | 27 | a `_test.go` registration of `os.Interrupt` alone does not red | new check test | a walk over test files reds the fixtures |
| LQ24 | 28 | a builder child dies after SIGKILL to its owner on Linux | new `runbinary` test on the parking fixture | without `Pdeathsig` the child survives the drain that never runs |
| LQ25 | 29 | `GOOS=darwin go build ./...` compiles | `TestResidualCheckCallsCrossCompileMatrix` | an untagged `Pdeathsig` field fails the darwin build |
| LQ26 | 28 | `Setpgid` stays set on the builder child | `TestCanonicalBuilderDrainsDescendantsBeforeReturningSelection` | a dropped `Setpgid` addresses the owner's own group |
| LQ27 | 1 | no production file outside `internal/canonicalpath` calls both `filepath.Abs` and `filepath.EvalSymlinks` | new dev-tier owner check and its fixture | a new package beside six untouched copies stays green |
| LQ28 | 19 | the gate subject closure returns the `Apply` error | new `subject.go` test | a dropped entry changes the oracle hash in silence |

### Edge inventory

- Error paths: a hold failure, a relative cache path, an unwritable cache, a spaced field name, and the kit checkout each refuse by name.
- Empty input: `GOCACHE=` is absent; a Sources continuation with no colon keeps the no-separator refusal.
- Boundary values: a field name with one space is wrapped; a name with none and outside the set is unknown.
- Interrupted state: a SIGKILLed owner leaves no builder child on Linux.
- Re-run idempotency: `link` in a linked consumer stays green on a second run.
- Hostile paths: a symlinked root and a not-yet-existing destination both resolve.
- Partial implementation: a package-internal change with the three swallow sites untouched reds LQ17.

**Won't handle** — `resolvedPath` in `internal/adopt` — the row names three derivations, and the adopt copy has no `Abs` step; the reviewer decides its fold.

**Won't handle** — `Pdeathsig` on the fourteen other `Setpgid` sites — the row names builder children, and the drain paths cover the rest.

**Won't handle** — widening the purity census to `internal/worktree` itself — the journey fixtures spawn git by design, and the row's scope is the policy packages.

**Won't handle** — a missing cache directory refusing — `HoldDir` creates it, and a first run must lock rather than refuse.

**Won't handle** — `leasedRepo` folded into `gittest` — it carries a lease fact `gittest` has no business knowing.

## Ownership fences

- `specs/roadmap-light-path-fixes-3/`
- `reviews/roadmap-light-path-fixes-3.md`
- `internal/canonicalpath/`
- `internal/worktree/subshell.go`
- `internal/preflight/gather.go`
- `internal/canary/mutation.go`
- `internal/adopt/link.go`
- `internal/adopt/adopt_test.go`
- `projects/benchkit.md`
- `internal/maps/validation.go`
- `internal/maps/maps_graph_test.go`
- `internal/conformance/decision_map_integrity_test.go`
- `tests/canary/decision-map-integrity/`
- `internal/puritycensus/`
- `internal/worktree/lifecyclepolicy/purity_census_test.go`
- `internal/worktree/reclaimpolicy/purity_census_test.go`
- `internal/worktree/landingpolicy/purity_census_test.go`
- `internal/gocache/gocache.go`
- `internal/gocache/gocache_test.go`
- `internal/gate/run_transaction.go`
- `internal/gate/lane.go`
- `internal/gate/subject.go`
- `internal/gate/report.go`
- `internal/testreport/command.go`
- `internal/conformance/canonical_path_owner_test.go`
- `tests/canary/canonical-path-owner/`
- `internal/skillsindex/command.go`
- `internal/conformance/cancel_signal_registrations_test.go`
- `internal/conformance/checks_test.go`
- `internal/conformance/registry/registry.go`
- `tests/canary/cancel-signal-registrations/`
- `internal/runbinary/`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/main_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_test.go`
- `tests/canary/guidance-prose-budgets/over-budget-skill`
- `tests/canary/line-routing/line-binding-prose-drift`
- `tests/canary/workflow-guidance-anchors/benchkit-hostile-input-heading`
- `tests/canary/workflow-guidance-anchors/benchkit-review-round-owner`
- `tests/canary/workflow-guidance-anchors/benchkit-review-round-routing`
- `tests/canary/workflow-guidance-anchors/benchkit-spec-ownership`

The fence is the union of the seven tickets' `Writes:` lines, closed by
`bench preflight build` over the fixture and registry pins. The derive ticket
lands first, because it rewrites six functions across the packages five
siblings edit. Those five follow it in parallel, and the link ticket runs
free.

## Out of scope

- The adopt `resolvedPath` fold: 2 edits, 1 gate run.
- `Pdeathsig` on the other fourteen `Setpgid` sites: 14 edits, 1 gate run.
- A repository-wide import-layering check: a spec, 2 gate runs.

## Further notes

Flagged additions beyond the decision source:

- The purity census bite test and the wrapper scanned-set assertion. The rows ask for a shared helper; without these the cheapest wrong helper stays green.
- The "exists" leg of the cache validation reads as "created by `HoldDir`", not "refused". The code creates the directory today, and a first run must lock.
- The canary fixture for the cancel-signal check and the fixture-inventory repair for decision maps. Both are what the existing fixture-bite sweeps demand of a new check or fixture.
- The canonical-path owner check. Without it, a new package beside three untouched copies stays green.
- The gate subject closure propagates an `Apply` error. The caller sits outside the row, and the refusal would otherwise vanish there.
- The shared forbidden-import list carries `internal/bounds` for `landingpolicy` too. That package imports no bounds today, so the tree stays green.

Two review rounds disagreed on the five registry closure files. The rule above
settles it: headroom files create no edge. A reviewer who wants the stricter
rule reopens it at sign-off.

Source-sentence-to-row table:

| source sentence | rows |
|---|---|
| FT276: one dependency-free package derives the path; three importers; no other derivation survives | LQ1, LQ2, LQ3, LQ27 |
| FT279: `bench link` exits non-zero through the predicate; one profile sentence | LQ4, LQ5, LQ6, LQ7 |
| FT285: a spaced field name is wrapped; the guard precedes the unknown-field branch | LQ8, LQ9, LQ10, LQ11 |
| FT202: one shared helper; package-local wrappers; process-backed fixtures count | LQ12 to LQ16 |
| FT264: `Hold` error refuses with error and path; empty absent; relative refuses; exists, absolute, writable | LQ17 to LQ20, LQ28 |
| FT201: a dev-tier check on `CancelSignals`; `Pdeathsig` on Linux | LQ21 to LQ26 |

Reviewer decisions closed on 2026-09-02 stand as the rows record them. Every
subagent runs `opus` at low or medium effort.
