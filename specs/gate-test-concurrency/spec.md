# gate-test-concurrency

Status: staged

Decision source: `decisions/gate-budget.md` decision #22 and ticket #23, reviewer-resolved 2026-08-07 (named reviewed artifact; the map stays top-level because decisions #8 and #24–#26 remain open).

## Problem

The dev gate's destination is a whole run under 120 seconds, but `internal/gate`
alone is 147–160 s focused and 230 s inside a gate. Decision #21 proved the
package is a strictly serial chain of subprocess waits — 241 top-level tests
whose elapsed sum equals the package wall at CPU/wall 0.68 — and that the sized
serial cuts project only a 105–125 s floor. The remaining lever is overlapping
tests inside the package run, and it is blocked by process-global state:
production `kitRoot` reads ambient `BENCH_KIT`, four `t.Setenv("BENCH_KIT")`
pins guard it — two of them inside the shared fixture constructors reached from
roughly 52 construction sites, which is what serializes most of the package —
and a handful of tests mutate other process-global state (working directory,
package-level variables, a shared lock registry).

## Solution

The kit root becomes an explicit input below the gate's public entries, so a
test injects it through the same seam production resolution uses instead of
pinning the process environment. With the constructor pins retired, every test
that touches no process-global state adopts `t.Parallel`, so the package's wall
is bounded by its longest overlapping chains instead of the serial sum of every
test's subprocess waits. Observable behavior is unchanged: a wrapper-routed
entry still honors `BENCH_KIT`, a direct entry still falls back to the graded
root, and the oracle's verdicts, evidence, and scoping decisions are
byte-identical.

One deviation from the decision source's wording, flagged for sign-off: #22's
recorded answer says the two chdir helpers retire. They cannot — those tests
exercise cwd-derived behavior (resolution from outside or beneath the graded
root), so the chdir *is* the input under test. They keep their chdir and stay
serial; on approval the map's answer is amended to match.

## User stories

1. As a kit developer, my `internal/gate` test run holds regardless of what an
   enclosing gate exported: evaluation resolves the kit root from an explicit
   injected input, and ambient `BENCH_KIT` below the public entries is inert
   for every consumer — phase-table resolution, component identity, and
   component scoping alike.
   Line: `opus` / medium. Production oracle seam change — the profile's cached
   gate-logic routing, with the gate covering behavior preservation.
2. As a kit developer, a focused `internal/gate` run overlaps its eligible
   tests under `t.Parallel`: the median of three focused `-count=1` runs on the
   12-CPU reference host lands at or below 90 s (from the 147–160 s serial
   floor), with every process-global-pinning test still serial and correct.
   Line: `opus` / medium. The eligibility walk is judgment over shared state,
   not purely mechanical; tickets may re-route down at charge time per
   `craft-line`.

## Implementation decisions

- **One resolution rule, resolved once.** The rule stays exactly today's —
  `BENCH_KIT` when set and non-empty, else the graded root — but it is applied
  only at the exported gate entries (the gate run command, the execute and
  tree-execute paths, the phases plumbing command, and the composed-green
  query `internal/status` consults; their signatures are unchanged, so no
  consumer outside the package moves). Below those entries the
  kit root travels as an explicit value; no production code in `internal/gate`
  reads the variable. The single-derivation property `kitRoot`'s comment
  guards ("every resolver of a phase table shares this rule") is preserved by
  construction: one entry-time resolution feeds every consumer. All three
  current consumers move — phase-table resolution, component-identity
  resolution, and component scoping — because a consumer left on the old
  fallback returns the right answer whenever kit equals root, which is every
  fixture; only the hostile-ambient guard can tell the difference.
- **The carrier is implementer discretion.** Whether the resolved kit rides as
  a parameter or as a field on an evaluation-scoped value is a reversible
  internal choice, bounded by the constraint above and by story 1's guard row.
- **Fixture constructors inject.** The kit-shaped and routed fixtures claim kit
  identity by injection rather than `t.Setenv`; the eligibility predicate
  component scoping applies (kit and root are the same directory) continues to
  hold for them by construction.
- **A kit-taking boundary beneath the exported entries.** Many fixture tests
  drive the exported entries directly (`RunCommand --fresh`, the reuse
  execute, the composed-green query), and entry-time env resolution would hand
  them a hostile ambient kit once the pins retire. Story 1 therefore exposes
  the package-internal execution boundary that takes the resolved kit
  explicitly; story 1's build migrates the non-entry-subject fixture tests
  onto it, while a small representative set keeps driving each exported entry
  with its own explicit env pin (serial — env is the input at that seam),
  mirroring decision #18's representative-control pattern.
- **Eligibility for `t.Parallel` is structural, over all process-global
  state.** A test adopts it exactly when it mutates none of: the environment
  (`t.Setenv`, `os.Setenv`/`os.Unsetenv`), the working directory
  (`t.Chdir`/`os.Chdir`), PATH, a package-level production variable (the gate
  timeout, the builtin-table stub), a shared global registry (the execution
  lock owners map), or process lifecycle (the `os.Exit` re-exec helpers).
  Go's runner fences only the `t.Setenv`/`t.Chdir` cases; the variable swaps
  restore via `t.Cleanup` and would corrupt concurrent tests silently, so the
  predicate — not the runner — is the authority. The build records the
  enumerated serial list with each test's pinning reason, and review grades
  that list against this predicate.
- **No test requires real overlap.** Adopted tests stay correct at
  `-parallel=1` and `GOMAXPROCS=2` (the canary-inner width), so the adoption
  cannot deadlock a narrow host; no cross-test synchronization is introduced.
  The three pre-existing `t.Parallel` calls simply gain company.
- **No width policy here.** `-p`, `-parallel`, and `GOMAXPROCS` pinning for the
  gate's own phases stay with `gate-budget` #8/#11; this spec changes the
  workload's shape, not its scheduling.

## Testing decisions

- The behavior a good test exercises: evaluation over a constructed fixture
  root resolves that root's phase table, component identities, and scoping
  verdicts — driven through the same resolution seam production entries use,
  never by patching the environment around it.
- Seams receiving tests: the package-internal evaluation seam, whose prior art
  is the kit-shaped fixture and the routed fixture helpers with their marker
  observability. One new serial guard test owns ambient inertness and is the
  durable holder of the property after the build; a reintroduced env pin on a
  parallel test additionally panics under Go's own runner.
- Gate seam observing the feature: the dev gate's `test` phase runs the
  package; a broken injection or a parallel-safety defect reds it. The
  promotion gate over the composed candidate is the in-gate witness that the
  parallelized package holds under the seven-phase overlap #20 measured —
  the FT171 sample already showed deadline-bounded FIFO tests exhausting
  under load, so a deadline that cannot survive overlap surfaces there, and
  serializing that test or raising its bound is a recorded build decision,
  never a silent weakening. The race registry is unchanged; a focused
  `go test -race ./internal/gate` run and one `GOMAXPROCS=2` run are ticket
  evidence, not new gate phases.
- Wall-clock is measured, not gated: the build records three focused
  `-count=1` repetitions before and after on the same host against story 2's
  ≤ 90 s median target, and `gate-budget` #26's census owns the authoritative
  re-measure.

### Seam diagram

    trigger: an exported gate entry (wrapper-routed or direct), or a package test
        │
        ▼
    graded root + explicit kit root  ──▶  [ phase-table & component resolution ]  ──▶  resolved table, identities, scoping, verdict
                      ◀ tests attach here: construct a fixture root, inject it as the kit,
                        evaluate, and read the table/identities/markers back

### Acceptance coverage map

| row | story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|---|
| KC1 | 1 | With ambient `BENCH_KIT` naming a foreign path and no per-test pin, evaluation over a fixture root resolves the fixture's own phase table, component identities, and component scoping (eligible, runner root = fixture) | serial ambient-guard test at the evaluation seam, asserting all three consumers | new test, red before the seam lands: each consumer's `kitRoot` call returns the ambient path, so the table fails to resolve and scoping goes ineligible | a consumer left on the old fallback is correct whenever kit equals root — only hostile ambient distinguishes it, so the guard must name every consumer or the missed one ships silently |
| KC2 | 1 | A wrapper-routed entry with `BENCH_KIT` set resolves the kit from the environment at entry — including kit ≠ root, the linked-repo case — and set-but-empty falls back to the graded root | existing entry-level tests that pin `BENCH_KIT`, including the kit≠root manifest-resolution test, plus one new serial set-but-empty entry assertion | mostly already covered — the existing entry tests red if env stops being honored or kit≠root routing breaks; the set-but-empty clause is new coverage (no test pins an empty `BENCH_KIT` today) and goes red against an entry that accepts the empty string as a kit | the entry rule is the one place env may decide, so its regression cover must survive the seam and the emptiness clause needs a real owner |
| KC3 | 1 | Fixture constructions carry no `BENCH_KIT` pin and the whole package is green with a hostile `BENCH_KIT` exported for the entire run | package run under `BENCH_KIT=<foreign> go test -count=1 ./internal/gate` as build evidence | observed red at build time: removing a pin before the seam lands reds the affected tests under the exported variable | proves pin retirement is real and injection carries every construction; the durable holder afterward is KC1's guard test |
| KC4 | 1 | The closed subject environment still strips `BENCH_KIT`/`BENCH_WRAPPER` from the gate script's env, per gate-budget #3 | existing subject-env stripping test | already covered — the existing test pins hostile ambient values and asserts the subject env excludes them | subject-input hygiene is a separate consumer of the same variable and must not regress while its reader moves |
| KC5 | 2 | Every test satisfying the eligibility predicate calls `t.Parallel`; the ineligible tests are enumerated with their pinning reasons; the package is green under repeated `-count=1`, focused `-race`, and `GOMAXPROCS=2`; the three-run focused median lands ≤ 90 s on the reference host | the package run, the three evidence commands, and the build's recorded serial-list enumeration | not TDD-able as a breadth assertion — a shared-state collision reds the package or `-race` run; breadth is review-graded against the predicate with the recorded enumeration and the ≤ 90 s median as its checkable witnesses | a partial adoption cannot meet the wall target, an unjustified serial entry fails the enumeration review, and any unsafe adoption is a red the run itself produces |

Degenerate-implementation check: story 1's cheapest wrong implementation leaves
one `kitRoot` consumer on the old fallback — correct in every kit==root test —
and KC1 is red on exactly that because it asserts all three consumers under
hostile ambient. Story 2's is `t.Parallel` on a handful of tests and a
done-claim — the ≤ 90 s median and the serial-list enumeration in KC5 both
fail it checkably.

### Edge inventory

- Hostile environment (foreign, empty, unset `BENCH_KIT`): KC1, KC2, KC3.
- Boundary values kit == root (every fixture, direct entry) and kit ≠ root
  (wrapper-routed linked repo): KC1/KC3 and KC2 respectively.
- Re-run idempotency: KC5's repeated `-count=1` evidence.
- Hostile width (`GOMAXPROCS=2`, `-parallel=1`): KC5's narrow-width evidence
  run; no adopted test requires real overlap.
- Interrupted/partial state under overlap: deadline-bounded tests that
  exhausted under load in the FT171 sample — the promotion gate is the in-gate
  witness, and any test that cannot hold its deadline under overlap is
  serialized or re-bounded as a recorded build decision (Testing decisions).
- Process-boundary lifecycle: a child process (the prospective shell path, the
  built binary's own entries) still receives the kit through the environment
  and resolves it at its own entry — KC2's seam is that entry, exercised
  cross-process by the existing built-command tests; the `os.Exit` re-exec
  helper tests stay serial under the predicate.
- **Won't handle:** retiring the chdir helpers — cwd is the input those tests
  exercise, so they keep their chdir and stay serial; flagged above as the one
  deviation from #22's recorded wording, amended in the map on approval.
- **Won't handle:** validating that `BENCH_KIT` names an existing kit-shaped
  directory at entry — today's trust posture is unchanged, and a wrong path
  already fails resolution loudly rather than silently.
- **Won't handle:** interruption semantics — no production concurrency is
  added, the runner's process-group teardown is untouched, and the existing
  cancellation tests keep their live coverage.
- **Won't handle:** a conformance check forbidding ambient `BENCH_KIT` reads in
  `internal/gate` — KC1's serial guard test owns the regression durably, and a
  registry check plus canary family for one package-local constraint outprices
  the defect it prevents. Reviewer may veto.

## Out of scope

- `internal/specbuild` test-only `t.Parallel` — `gate-budget` #24, reviewer-
  decided separate light path landing alongside this build (~2 edits, 1 gate
  run).
- The three sized serial cuts (grace observation, synthesized matrix
  generations, closure memoization) — `gate-budget` #25, after this spec
  retires (~3 light-path tickets, each 1–3 edits and 1 gate run).
- Width and reserve pricing, the token pool, and any `-p`/`GOMAXPROCS` pinning
  of gate phases — `gate-budget` #8/#11 after #26's census re-run.
- Migrating the package-variable seams (gate timeout, builtin-table stub) to
  injection so their tests could parallelize too — a separate small capability
  (~4 edits, 1 gate run) that the serial-list enumeration will price precisely;
  the tests stay serial here.
