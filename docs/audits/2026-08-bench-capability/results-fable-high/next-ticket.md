# Next ticket — the dev gate grades the kit's own contracts

**Title.** Live-root conformance runs inside `bench gate`'s ordinary test phase; an
environment-class skip inside the oracle is red; every skip is named.

**Blocked by:** none.
**Writes (proposed fence):** `internal/gate/phases.go`, `internal/gate/capability_skips.go`
(+ tests), and the files the ten current diagnostics name (`.agents/commands/bench-implement-spec.md`,
`.agents/commands/bench-final-check.md`, `.bench/BENCH.md`, `.bench/BENCH-reference.md`,
`decisions/spec-build-review-gate-cadence.md`, `internal/canary/canary.go`) — the last group
subject to the reviewer's per-diagnostic call (fix the tree vs amend the contract).

## Problem

`internal/conformance` registers 29 checks that enforce the kit's guidance/doc/adapter
contracts (`docs-currency-workflow`, `skills-index-command-adapters`,
`handoff-shape-single-source`, `decision-map-integrity`, `coverage-map-validation`,
`guidance-prose-budgets`, `axi-query-registry`, …). They reach the live tree only through
`TestRootConformance`, which opens with `capability.Environment(t, "BENCH_CONFORMANCE_ROOT
not set")`. `bench gate`'s phase table is gofmt · vet · test · race · system · shellcheck,
its test phase runs `go test -count=1 ./...` with no such variable, and the only
non-test setter in the repo is `bench prep-release` (ship tier). So between releases the
kit's own contracts are unenforced, the gate footer folds the loss into
`capability-skips: 7 (capability=6 environment=1)` without naming it, and at HEAD ten
diagnostics are red while the gate is green — including a phase file
(`bench-implement-spec.md`) that lost its required `## Entry orientation` / `## Exit
handoff` headings on 2026-08-11 (`fa4e1f02`) and landed green.

History matters for the fix shape: `72c037a1` (2026-07-05) moved the shell conformance
run into a Go `conformance` phase; `3701c4a0` (2026-08-09, "adopt branch-native test
architecture") removed that phase. `projects/benchkit.md` records the decided state:
*"Go owns package scheduling inside the one ordinary test driver. There is no separate
contract or conformance dev driver …"* — while the same profile still says *"The gate's
conformance layer enforces those contracts so disk, docs, and adapters do not drift."*
The decided shape is one ordinary driver **that grades the live root**; what shipped is
one ordinary driver that skips it.

## Evidence

- `BENCH_CONFORMANCE_ROOT=$PWD go test -count=1 ./internal/conformance -run
  '^TestRootConformance$'` → FAIL, 10 diagnostics, 3.24 s (this run; identical list to
  Opus E6). `bench gate` → green, 64.8 s, `environment=1` unnamed (this run).
- `rg BENCH_CONFORMANCE_ROOT --glob '!*_test.go'` → only `internal/preprelease/preprelease.go:100`.
- `internal/gate/capability_skips.go` already retains `environmentReasons`; `skipRows`
  never prints them.
- Roadmap already knows the trap: FT133 (occurrences 2026-07-26 and "FT126 recurrence"),
  FT120 (2026-08-14 "ordinary `go test ./...` skipped TestRootConformance"), FT213
  (2026-08-16 "two delegates read one TestRootConformance environment skip two different
  ways") — none names that the phase was removed.
- Ledger L-01, L-21; bug triage B1, B10, B19.

## Architectural boundary

Owner: the gate (Bench core) and the conformance registry. The change stays inside the
profile's decided shape — no separate driver, no new phase name, no change to
`prep-release`'s ship tier, no change to the evidence key. Judgment checks stay prose;
this ticket only makes the *registered mechanical* checks run where the doctrine already
says they run.

## Scope

1. In `BenchkitPhases`/`toolchainPhases`, give the ordinary `test` phase
   `BENCH_CONFORMANCE_ROOT=<graded root>` and `BENCH_CONFORMANCE_TIER=dev` in its `Env`
   (or make `TestRootConformance` default the root to the kit module when the graded root
   *is* the kit — pick the one that keeps linked repos unaffected; the env route is
   simplest and mirrors what `prep-release` does at ship tier).
2. In `capability_skips.go`: render `class=environment: N (TestName: reason)` for every
   environment skip; when the run is the oracle (not `bench test`), an environment-class
   skip is red with a message naming the test and reason.
3. Run the gate at HEAD → red with the ten diagnostics. Disposition each with the
   reviewer: restore the two headings in `bench-implement-spec.md`; fix the two stale
   `$bench-finalize-spec` references and the two dangling `Sources` paths in
   `decisions/spec-build-review-gate-cadence.md` (this also greens `bench maps`); remove
   the `spec build` token at `BENCH-reference.md:106`; restore the implementation-retro
   owner lines in `bench-final-check.md`/`.bench/BENCH.md` *or* retire that contract if
   the drop was intentional; make `canary.go` consume `bounds.CanaryInnerWidth` *or*
   retire that binding. Each amendment to a contract is a reviewer decision recorded in
   the commit.
4. Update `projects/benchkit.md` (line 206 region) and `.bench/BENCH-reference.md:131`
   so the profile describes the shipped shape (ordinary driver grades the live root;
   no separate phase).
5. Close or annotate the three roadmap occurrences (FT133, FT120, FT213).

## Non-goals

No new gate phase; no router; no prose rewrite beyond the ten diagnostics; no change to
capability-class (fifo/privilege) skip posture; no change to `bench test`'s skip
reporting semantics beyond naming; no re-litigation of the branch-native decision.

## Acceptance criteria

- [ ] `bench gate` at `58d966e2` exits non-zero and names all ten current diagnostics
      with their existing messages (red first — proves the phase is connected to the
      live root, not to fixtures).
- [ ] After the dispositions land, `bench gate` is green and the test-phase output shows
      `TestRootConformance` **ran** (`--- PASS`/`ok` with the test executed), not skipped.
- [ ] Deleting `## Entry orientation` from any `.agents/commands/*.md` reds the gate with
      `docs-currency-workflow`'s own message; restoring greens it.
- [ ] Removing the env assignment reds the gate with an environment-skip-inside-the-oracle
      message naming `TestRootConformance`, not a footer count.
- [ ] Skip line format: `capability-skips class=environment: 1 (TestRootConformance:
      BENCH_CONFORMANCE_ROOT not set)` when it does happen (e.g. under `bench test` on a
      linked repo).
- [ ] `bench maps` exits 0.
- [ ] `bench prep-release`'s ship-tier conformance still runs once (no double-run of the
      dev set) — verified by its existing test (`preprelease_test.go:54`).
- [ ] Added gate wall-clock reported in landing evidence (expected ≈ +3 s on 65 s).
- [ ] FT133/FT120/FT213 skip occurrences closed or annotated.

## Validation

- Red-first at HEAD (above).
- Mutation probe of omission kind: drop the env assignment → red on the environment rule.
  Restore.
- Mutation probe of a different kind and site (per `craft-delegate`): plant a heading
  deletion in a phase file → red with the check's message. Restore.
- Green after dispositions; `bench canary` still reports 233 bindings; race/system phases
  unchanged.
- Cost line in the landing report.

## Migration

None for linked repos: the env is set only when the graded root carries the conformance
package (kit-only), exactly as the race and system phases already materialize.

## Dependencies

None. Unblocks A3 (router phase file and rename are graded), A5 (debug phase file
budget/contract), A8 (handoff-shape check), A9 (prose fixes), A10 (structured skip
output).

## Rollback

Revert the two Go hunks; the gate returns to today's behavior (green with an unnamed
environment skip). Contract/tree dispositions are independent commits and can stay.

## Why this outranks the alternatives

- **A2 (verdict reader)** is equally small and equally P0, but it changes what the
  dashboard *reports*; A1 changes what the oracle *sees*, and A1 is what would have
  caught A2's class of drift in docs. Ship A2 immediately after — same session if the
  reviewer allows two light-path tickets.
- **A3 (/bench router)** is the user-visible win, but its new phase file, rename, and
  adapters are exactly the artifacts the conformance registry grades — landing it before
  A1 means landing it ungraded.
- **Sol's CR-001 (work-state/router/compiler tracer)** is INFERRED, large, and would
  duplicate readers Bench's own standard forbids; nothing observed in three audits
  requires it.
- **FT100 (prose cut, current rank 1)** is the highest-variance action available with no
  measurement behind it; A1 first makes any later cut graded, and A11 makes it measured.
