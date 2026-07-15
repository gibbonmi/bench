# Static, bounded guard discovery (FT80)

Status: implemented

## Problem

`bench guards` and the SessionStart inspection describe every guard by **executing**
each hook's `bash <hook> --describe` protocol. Executing a file just to read its
description is the collision this surface was meant to avoid: a hostile or slow hook
body runs (bounded only by a per-hook timeout and a foreign-pre-push carve-out), and
the one file already guarded against foreign execution — the git pre-push hook — is
the exception that proves the rule. `RR:R-07` / `RC:H-02` stay open until the guards
path spawns no process at all.

## Solution

Every guard describes itself from **static comment-header lines** at the top of its
script (`# name:`, `# boundary:`, `# denies:`, `# why:`), parsed as text. `bench
guards` and the SessionStart inspection read those headers and never execute a hook.
The `--describe` protocol is deleted everywhere. A new Go plumbing subcommand runs
the three SessionStart inspections under one aggregate deadline so a hung phase warns
and the session still opens. A runtime sentinel proves non-execution, and a recorded
mutation at TDD-red proves the sentinel bites.

**Reviewer-approved deviation from the map (decision #1).** The map said the header
carries each guard's `denies` verbatim. Three guards derive `denies` from live state
today — `block-dangerous-git` from the gitguard class list, `check-agent-line` from
the repo's `lines.env` binding, `pre-push` from live `origin/HEAD` + pin arming — so a
frozen header would ship the kit's own per-repo values into every linked repo and
duplicate a single-sourced fact. Per reviewer decision (2026-07-15), every `# denies:`
header is a **static generic** human string; the exact live values leave the guards
row and stay at their enforcing source. This trades away the row's current live
precision (exact classes / bound ids / protected branch / pin state) for a fully
static, per-repo-correct surface.

## User stories

1. As an agent, I want `bench guards` to describe every guard without executing any
   hook script, so a hostile or slow hook body can never run or stall aggregation.
   Line: gpt-5.6-terra / medium. This is the deep parser at the oracle-adjacent
   surface where non-execution correctness is the whole point.

2. As an agent, I want each hook's manifest read from a static comment header
   (`# name:`, `# boundary:`, `# denies:`, `# why:`) at the top of the script, so
   reading a guard's description can never execute code. Line: gpt-5.6-luna / low.
   Authoring the header lines on the shipped hooks is mechanical text editing at a
   fixed shape.

3. As an agent, I want the `--describe` branch deleted from every hook, the prepush
   template, and the now-orphaned core `guard-git --describe-classes` /
   `check-agent-line --describe-binding` subcommands, so exactly one description
   mechanism exists and no dead exec path remains. Line: gpt-5.6-luna / low. This is
   deletion verified by a reference sweep, not new design.

4. As an agent, I want a hook with an absent or incomplete header to render a
   definitive `no manifest` row rather than a silent omission, so a malformed or
   header-less hook is visible in the surface. Line: gpt-5.6-terra / medium. Same
   parser seam as story 1; the fallback classification is correctness-sensitive.

5. As an agent, I want the pre-push row classified managed / unmanaged / not-installed
   from the bench marker and a static read — never executing the pre-push body — so a
   foreign pre-push is never run in order to be described. Line: gpt-5.6-terra /
   medium. The non-execution guarantee is load-bearing for the one historically
   foreign-executed file.

6. As an agent, I want informational hooks (`session-start`, `worktree-lifecycle`,
   whose header reads `denies: nothing (informational)`) excluded from the guard
   rows, so the surface stays deny-capable-only. Line: gpt-5.6-luna / low. This
   preserves existing exclusion behavior across the header change.

7. As a session, I want SessionStart's three inspections (resume-clean, status,
   guards --brief) run by one Go plumbing subcommand under a single ~10s aggregate
   deadline, so a hung phase emits a warning and the session still opens (exit 0).
   Line: gpt-5.6-terra / medium. The never-block posture is a correctness contract,
   not plumbing.

8. As a session, I want `session-start.sh` reduced to resolve-and-exec plus the
   CLI-location line, so the orchestration and the deadline live in the testable Go
   core. Line: gpt-5.6-luna / low. Shrinking the shim to a pass-through is mechanical.

9. As the kit maintainer, I want the conformance manifest check to validate static
   headers instead of executing `--describe`, and the stale `--describe` references
   in `projects/benchkit.md` and `.bench/BENCH-reference.md` corrected in the same
   diff, so the gate proves the shipped hooks self-describe without running them and
   no doc points at a deleted protocol. Line: gpt-5.6-terra / medium. Conformance is
   oracle logic where a wrong check is the worst class of kit bug, and the doc edits
   are mechanical stale-reference corrections tracking that code change rather than the
   net-new guidance authoring the top-tier leverage override is meant for.

10. As the kit maintainer, I want the `tests/canary/package-core-guard/guard-describe-*`
    fixture family migrated to the header convention, so the gate still goes red when a
    shipped header drops a required key. Line: gpt-5.6-luna / low. Fixture-tree edits
    at a known EXPECT shape.

11. As the kit maintainer, I want a runtime contract sentinel proving `bench guards`
    and the inspection leave no evidence file — a hook that writes-on-exec is never
    run — so FT80's non-execution guarantee is gate-observable against the built
    `dist/bench`. Line: gpt-5.6-terra / medium. This test is the spec's closure proof;
    its correctness matters more than its speed.

12. As the kit maintainer, I want the non-execution rule's bite demonstrated by a
    recorded mutation at the TDD-red step (revert guards to executing the hook, watch
    the sentinel test go red, record it), so we hold evidence the sentinel catches a
    regression the gate cannot hold continuously. Line: gpt-5.6-terra / medium. This
    is a judged process step, not a standing gate check.

## Implementation decisions

- **`internal/guards` becomes a static header parser + row assembler.** Delete
  `guardDescribe`, `describeTimeout`, `waitGrace`, and all `exec`/`syscall`/`context`
  process-group machinery. `guardRow` reads the file, parses the header, and builds
  the row; `manifestField` becomes a header-line parser. `Rows`, `Command`,
  `wiredHarnesses`, and `withWired` keep their shapes — `wiredHarnesses` already reads
  harness configs as *data*, not by execution, and is unchanged. The package doc
  comment loses its `--describe`/time-bound narrative and states the static-header
  contract.

- **Header grammar (map deferred this).** The manifest is the set of `# <key>: <value>`
  lines found in the file's **leading comment block** — the run of lines from the top
  that are blank or begin with `#` (a `#!` shebang is skipped as a comment line).
  Parsing stops at the first line that is neither blank nor `#`-prefixed (in practice
  `set -uo pipefail`). Keys `name`, `boundary`, `denies`, `why` are order-independent;
  first occurrence wins; the value is the text after `# <key>: ` with trailing
  whitespace trimmed; an empty value counts as missing. Interleaved comment prose
  (block-dangerous-git's threat-model paragraphs) is tolerated. Any of the four keys
  absent or empty ⇒ `no manifest` row.

- **Static generic `denies` per guard** (the reviewer-approved deviation):
  `block-dangerous-git` → `destructive git operations`; `check-agent-line` → `Agent
  delegation off the bound tier`; `stop` → `stopping an armed shift while the gate is
  red`; `session-start` and `worktree-lifecycle` → `nothing (informational)`;
  `pre-push` template → `direct push to the protected branch; .bench drift when
  pinned`. Boundaries and names are unchanged from today's manifests.

- **Pre-push classification.** `prePushRow` reads the installed hook: marker present ⇒
  parse its static header (managed); marker absent ⇒ `unmanaged (no manifest)`; file
  absent ⇒ `not installed`. No branch of this path executes the file. The `git`
  wired-cell constant is unchanged.

- **Delete the `--describe` protocol everywhere.** Remove the `--describe` branch from
  all five `.bench/hooks/*.sh` and `internal/adopt/prepush.sh`. The core subcommands
  `guard-git --describe-classes` and `check-agent-line --describe-binding` orphan once
  no hook calls them (confirmed: referenced only by the hooks and their own tests);
  delete both handlers and their tests. The underlying gitguard class table stays — it
  feeds the actual verdict message, not just the deleted describe.

- **Retire the tests that assert a hook's `--describe` output** (in the same diff — the
  protocol they exercise is gone). This is a reviewer-visible disposition, not a
  discover-at-red: (a) `internal/contract/surface/prepush_test.go` — the
  `runPrePushDescribe` assertions that the pre-push `--describe` denies varies with pin
  state (`.bench drift from bench gate pin` present when pinned, `drift check disarmed`
  when not) are removed; the pre-push **enforcement** tests (a pinned hook blocks a
  drifted push) stay untouched. (b) `internal/contract/axi/axi_fail_closed_test.go` —
  the `denies: manifest unavailable (analyzer missing)` describe assertion is removed
  (static parsing has no analyzer dependency, so the describe fail-closed posture has
  no analogue); the block-dangerous-git **enforcement** fail-closed stays. (c) In
  `internal/contract/axi/axi_guards_test.go` — the describe sub-assertions
  (`testAXIGuardsDescribeTimeoutBound`, `testAXIBlockDangerousGitCoreUnreachableManifest`,
  and the pre-push manifest-key loop inside `testAXIGuardsAggregation`) are removed or
  rewritten to assert the static row instead. These retirements are deletions, not
  new red signals, so they carry no coverage row — they are disposed here explicitly.

- **Sentinel fixtures are write-on-exec across every header class.** The non-execution
  proof must defeat not just an execute-every-hook implementation but the cheapest
  static-parse-**with-exec-fallback** degenerate (parse headers, but `exec` the hook
  when the header is absent or incomplete). So the sentinel fixture family carries a
  write-on-exec body in *each* class — full header, incomplete header, absent header,
  and an informational one — and asserts no evidence file after `bench guards` + the
  inspection. A single full-header sentinel would pass a fallback impl and prove
  nothing about the header-less path.

- **New `bench session-inspect` plumbing subcommand.** A `func Inspect(ctx
  context.Context, w io.Writer, root string) int` runs the three phases (resume-clean,
  status, guards --brief) in order, writing their combined output; `Command` wraps it
  with a 10s deadline via `context.WithTimeout`. On deadline trip: emit a warning to
  stderr and return 0 — never block. `session-start.sh` reduces to resolve-the-wrapper,
  print the CLI-location line, then `exec <wrapper> session-inspect`. Added to the
  plumbing-subcommands list in `.bench/BENCH-reference.md` (hook-driven, never typed by
  sessions).

- **Deadline test seam (map uncertainty flag #7, spec-writer's pick).** The deadline is
  tested at the **unit-level context seam**, not with a real sleep or a production
  test-only env knob. The inspect package exposes its phase list as an unexported,
  test-overridable slice (`var phases = []phase{…}`) so a same-package unit can
  substitute a stub phase that blocks on `ctx.Done()`; the unit drives `Inspect` with
  that stub and an already-short/cancelled context, asserting the warning + return 0.
  `Command` owns the concrete 10s value; the seam owns the trip behavior. Because a
  `Command` that forgot `context.WithTimeout` would still pass that unit **and** the
  happy path, a second unit asserts the context `Command` hands to `Inspect` carries a
  deadline — the never-block wiring is verified, not assumed. The black-box happy-path
  (exit 0, phases in order, guards brief present) is asserted against the built binary
  through the existing SessionStart contract tests, repointed at `session-inspect`.

- **Conformance migration.** `checkGuardDescribeManifests` stops executing `bash
  <hook> --describe` and instead validates each shipped hook's static header (and the
  prepush template's), asserting `name`/`boundary`/`denies`/`why` present and
  `session-start` informational. It reads through the **`internal/guards` header
  parser** — the deep unit (Handoff item 3) — never a second grammar implementation in
  `internal/conformance`; the grammar has one source and cannot drift between the two
  readers. Rename to reflect static validation. Fix the two stale `--describe`
  references in `projects/benchkit.md` (lines ~36 and ~136) and the one in
  `.bench/BENCH-reference.md` (~131) in the same diff so the stale-command sweep stays
  green.

- **Canary migration.** `tests/canary/package-core-guard/guard-describe-boundary-dropped`
  (and any sibling `guard-describe-*`) move from an executable `--describe` that drops a
  key to a static header that omits `# boundary:`; the EXPECT substring (`manifest
  missing boundary`) is preserved so the family still bites.

## Testing decisions

- **What a good test is here:** drive `bench guards` / the inspection at the CLI seam
  against fixture hook trees and observe TOON rows and the evidence-file's absence —
  never read the parser's internals. The deadline trip is the one exception: it tests
  through the `Inspect` context seam because a black-box sleep is slow and flaky.
- **Seams tested:** (1) `internal/guards` `Rows`/`Command` via the AXI runtime contract
  and `guards_test.go`; (2) the conformance static-header check via a canary fixture +
  the conformance suite; (3) `Inspect(ctx,…)` via a Go unit for the deadline, and the
  built binary via the repointed SessionStart contract tests for the happy path.
- **Prior art:** `internal/contract/axi/axi_guards_test.go` (aggregation, brief,
  unmanaged-pre-push safety, path-with-spaces, subdirectory), the SessionStart
  injection/order/never-block tests in the same file, and
  `checkGuardDescribeManifests` in `internal/conformance/package_core_checks_test.go`.
- **Gate command:** the project gate, `.bench/gate.sh`.

### Seam diagram

**Seam 1 — `internal/guards` static parser + row assembler:**

    trigger: `bench guards [--brief]`, or `session-inspect` calling guards --brief
        │
        ▼
    .bench/hooks/*.sh headers ──▶  [ Rows: read file → parse       ]  ──▶  guards[N]{guard,
    installed pre-push header  ──▶  [ leading-comment header →       ]        boundary,denies,wired}
    .claude/.codex configs     ──▶  [ row; marker for pre-push;      ]        TOON table / --brief lines
    (read as data, never exec) ──▶  [ no manifest on absent/partial ]
                                        ◀ tests attach here: AXI runtime contract feeds a
                                          fixture hook tree (incl. a write-on-exec sentinel)
                                          and asserts rows + evidence-file absence

**Seam 2 — conformance static-header validation:**

    trigger: gate → TestRootConformance → checkGuard…Manifests(root)
        │
        ▼
    shipped hook + prepush ──▶  [ read leading-comment header →   ]  ──▶  [] diags (green)
    template headers            [ require name/boundary/denies/why ]       or "manifest missing <key>"
                                [ session-start informational      ]
                                    ◀ tests attach here: canary fixture drops `# boundary:`,
                                      gate goes red with "manifest missing boundary"

**Seam 3 — `Inspect` under one aggregate deadline:**

    trigger: session-start.sh → `exec <wrapper> session-inspect`
        │
        ▼
    ctx(10s), root ──▶  [ Inspect: resume-clean → status →  ]  ──▶  combined stdout,
                        [ guards --brief, in order;          ]       exit 0 (always)
                        [ deadline trip → warn(stderr), 0    ]
                            ◀ tests attach here: Go unit drives Inspect with a ctx-blocking
                              stub + short context (trip → warn+0); built-binary contract
                              asserts happy-path order/content

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1, 11 | guards describes every guard while executing nothing | Seam 1 (AXI contract) | new sentinel test: fixture hooks dir holds write-on-exec scripts in *every* header class — full, incomplete, absent, informational — each `touch`ing a distinct evidence file if run; after `bench guards` + `session-inspect`, `os.Stat` on each must error | an execute-every-hook impl and the cheaper static-parse-**with-exec-fallback** impl (execs only on absent/incomplete headers) both leave an evidence file → red; a single full-header sentinel would miss the fallback |
| 2, 4 | header parsed from leading comment block; absent/partial ⇒ `no manifest` | Seam 1 (`guards_test.go` + AXI) | fixture hook with no header, or missing one key, must render `<name>,"",no manifest,none`; a hook with a full header renders its parsed row | a parser that only handled `--describe` output, or dropped header-less hooks, would omit the row or mis-shape it |
| 5 | pre-push classified managed/unmanaged/not-installed, never executed | Seam 1 (`testAXIGuardsUnmanagedPrePushSafety`, kept) | foreign pre-push writes a sentinel on exec; row is `unmanaged (no manifest)` and the sentinel is absent | executing the foreign pre-push to describe it writes the sentinel → red |
| 6 | informational hooks excluded from rows | Seam 1 (AXI aggregation) | `session-start` and `worktree-lifecycle` must not appear as guard rows; the four deny-capable guards must | a header parser that ignored `denies: nothing (informational)` would leak an informational row |
| 3 | `--describe` protocol fully deleted | conformance stale-reference + build | grep/reference sweep: no `--describe` in hooks/prepush; orphaned core subcommands removed; build green | a leftover `--describe` branch or an orphaned subcommand keeps a dead exec path the gate flags |
| 7 | inspection runs under one deadline; trip warns and exits 0 | Seam 3 (Go unit) | `Inspect` with a ctx-blocking stub phase + short deadline must return 0 and emit the warning; without the deadline handling it hangs or returns nonzero | a missing/incorrect deadline branch blocks the session or propagates a nonzero exit |
| 7 | `Command` actually installs the aggregate deadline | Seam 3 (Go unit) | a unit asserts the context `Command` hands to `Inspect` has a deadline set (`ctx.Deadline()` ok) | a `Command` that forgot `context.WithTimeout` passes the trip unit and the happy path but silently drops the never-block guarantee |
| 7, 8 | happy-path inspection: resume → status → guards brief, in order, exit 0 | Seam 3 (built-binary contract, repointed) | existing SessionStart order/injection/never-block tests pass against `session-inspect`; wrong order or a swallowed phase reds | the current `testSessionStart*` assertions already pin order and content; repointing preserves the red signal |
| 9 | conformance validates headers, not `--describe`; docs corrected | Seam 2 (canary + conformance) | canary `guard-describe-boundary-dropped` reds with `manifest missing boundary`; stale-reference sweep green after doc edits | a check that still executed `--describe`, or a doc still naming `--describe`, would pass a header-broken tree or fail the sweep |
| 10 | canary family bites on a dropped header key | Seam 2 (canary) | the migrated fixture's EXPECT (`manifest missing boundary`) matches the conformance diagnostic | a fixture left on the old `--describe` shape would no longer exercise the migrated check |
| 12 | the sentinel actually bites | mutation demo (TDD-red, recorded) | at red step, revert `guardRow` to `exec`; the story-1/11 sentinel test goes red; record it | proves the non-execution assertion is not vacuous — a fixture tree cannot violate binary behavior, so the mutation is the only evidence |

### Edge inventory

Walked per the profile's shell-CLI hostile-input checklist and map item 6:

- **absent header vs incomplete header** — coverage rows (story 4): both render `no
  manifest`, asserted distinctly from a full-header row.
- **leading shebang + blank/prose comment lines before the keys** — covered by the
  grammar (skip shebang, scan the leading comment run); block-dangerous-git's long
  header is the live fixture.
- **missing trailing newline on the last header line** — the line-split parser treats a
  final unterminated `# key: value` as a valid line — **Won't handle** as a separate
  row: the shipped hooks end in newlines and the parser's `strings.Split` already
  yields the last line; no distinct behavior to assert beyond story 2.
- **control/binary bytes in header text** — a value with an ESC/BEL byte flows to
  `toon.Table`, which refuses control bytes and errors (existing render-error path);
  **Won't handle** a bespoke row: the TOON emitter already owns this and returns the
  structured render error the AXI contract covers elsewhere.
- **huge header file** — reading a hook file into memory to parse a leading comment run
  is bounded by the hook's own size (kilobytes); **Won't handle** a size cap — the
  input is kit-shipped or project-authored hook scripts, not untrusted network data.
- **non-`.sh` entries and subdirectories in `.bench/hooks/`** — existing skip logic in
  `Rows` (unchanged); asserted by `testAXIGuardsUsageSubdirectory` and aggregation.
- **foreign / marker-less pre-push** — coverage row (story 5).
- **pre-push absent** — coverage row (story 5): `not installed`.
- **unparseable harness config** — `wiredHarnesses` scans invalid JSON as not-wired
  (unchanged); **Won't handle** anew — JSON-validity conformance owns malformedness.
- **path with spaces / invocation from a subdirectory / via symlink** — existing AXI
  contract tests (`testAXIGuardsPathWithSpaces`, `testAXIGuardsUsageSubdirectory`, the
  `pwd -P` resolution) still hold; static parsing does not touch path resolution.
- **inspection deadline trip** — coverage row (story 7).
- **resume-clean failure inside the inspection** — the existing
  `testSessionStartSurfacesResumeFailure` warning must survive the move to
  `session-inspect`; covered by the repointed happy-path contract.
- **`session-inspect` outside a repo** — existing never-block-outside-repo posture:
  silent, exit 0; covered by the repointed `testSessionStartNeverBlocksOutsideRepo`.
- **required tool missing from PATH (no wrapper)** — `session-start.sh` keeps its
  resolve-or-exit-0 rim; the shim degrades to silence, never blocks — preserved from
  today.
- **project-added hooks that answered `--describe`** — on a kit update they read as
  `no manifest` until they carry static headers (Handoff item 9). **Won't handle** an
  auto-migration or warning: the header convention is a documented contract change; the
  `no manifest` row is the honest, visible signal, and story 4 already asserts it.

## Out of scope

- **Deleting the gitguard class table now unused by describe.** It still feeds the
  live verdict message, so it is not dead; no cut to make. (0 edits.)
- **A conformance cross-check that the generic `denies` strings match their enforcing
  source.** The reviewer chose generic-static denies precisely to *stop* advertising
  the live value, so there is nothing to cross-check; a drift gate here would
  re-introduce the coupling the decision removed. Separate capability if ever wanted.
  (~2 edits, 1 gate run.)
- **Making `session-inspect` a user-facing command.** It is hook-driven plumbing; a
  documented user surface is a separate capability with its own help/usage contract.
  (~3 edits, 1 gate run.)
