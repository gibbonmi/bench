# Spec comparison: gate-phase-concurrency

Two candidate specs for the same feature, both compiled from the closed map
`decisions/gate-phase-concurrency.md`:

- **Spec A** — `specs/gate-phase-concurrency.md` (9 stories, mid-tier routing)
- **Spec B** — `specs/gate-phase-concurrency-sonnet.md` (7 stories, all-cheap routing)

## Verdict

**Spec A is the better result.** Keep it; retire Spec B after porting one
insight from it (noted below). Spec A wins on the three things that matter most
for this feature: its line routing follows the profile, its coverage map
actually pins the headline behavior (concurrency), and its edge inventory
follows the template's row-or-won't-handle discipline. Spec B is leaner and has
one genuinely sharper piece of analysis, but it under-routes the oracle's own
logic and its coverage map would pass a sequential implementation.

## Decisive differences

### 1. Line routing — Spec A follows the profile, Spec B contradicts it

The profile (`projects/benchkit.md`, Lines) caches a routing: **"Gate /
conformance logic → mid effort. Correctness of the oracle matters more than
speed."** This feature rewrites the oracle's own orchestration.

- Spec A routes the executor, multiplexer, aggregation, inner-mode, and signal
  stories to `claude-opus-4-8` (mid) and only the mechanical dispatch plumbing
  (story 9) to `claude-sonnet-5` — consistent with the cached routing and with
  `craft-line`'s up-bias row (buy insurance where the gate can't cheaply catch
  a confidently wrong answer; concurrency and line-multiplexing bugs are
  exactly that class — nondeterministic, only partially observable by tests).
- Spec B routes **every** story `claude-sonnet-5 / low`, justified as "the gate
  fully observes the contract." That routing also flips the venue: under
  `craft-line`, an all-cheap fully-observable spec qualifies for a headless
  `bench shift` on Sonnet — so the mis-routing isn't cosmetic, it changes where
  and by what the oracle gets rewritten. Under-escalation is the expensive
  error (`craft-line`), and this is the one seam the whole kit stands on.

### 2. Coverage of the headline behavior — Spec B never asserts concurrency

- Spec A, row 1: inject four fake phases (marker + sleep), assert all markers
  present **and wall time below the serial sum** — red until a concurrent
  executor exists, green milliseconds thereafter.
- Spec B, row 1: assert `bench gate-phases` exists and an all-green run exits 0.
  That red signal proves the subcommand exists, not that anything runs
  concurrently. **A fully sequential Go port would pass Spec B's entire
  coverage map.** Spec B leans on the one-off build-time wall-clock measurement
  (correctly not a gate assertion per Handoff #5), but that leaves the feature's
  defining behavior with no permanent seam at all.

Spec B's row 1 also implies running the real all-green gate inside a test —
which recurses (the canary-sweep phase runs 61 inner gates) and costs ~35s per
run. Spec A's injectable `[]Phase` seam with `bash -c` stubs avoids both.

### 3. Handoff fidelity — Spec A covers all the black-box assertables

Handoff #4 lists the assertables. Spec A has a coverage row for each, including
the **verbatim `gate: green` / `gate: red` final line and exit 0/1/3** (its
story 3 / row 3, plus exit 130 on SIGINT). Spec B mentions the verbatim line in
prose but has **no coverage row asserting it** — and that line is load-bearing
surface (verdict cache, Stop hook, pre-push, humans; Handoff #9 flags it
explicitly).

### 4. Edge inventory discipline

The template requires every edge class to land as a coverage row or a one-line
**Won't handle**. Spec A does exactly that (root-with-space and
shellcheck-absent are explicit rows with concrete red signals; the n/a classes
are Won't-handle lines with one-clause justifications). Spec B has several
"Covered by story N" claims with no backing row: go-toolchain-missing is
"covered by story 3," but story 3's test injects a stub failing phase and never
exercises a missing toolchain; paths-with-spaces is "covered by story 7's"
test, which is itself the map's vaguest row ("drives the gate through more than
one surface"). Claimed-but-unbacked coverage is the failure the edge inventory
exists to prevent.

## Where Spec B is better

- **The canary-can't-catch-it analysis.** Spec B argues precisely that the
  canary sweep's `EXPECT` check is a raw substring match, so inner-mode output
  gaining a `[phase] ` prefix or reordering would **still pass** the sweep —
  which is why the direct inner-mode byte-shape test is the *only* real pin
  (its rows 4a/4b split shape vs. bite). Spec A includes the same direct test
  but its row-5 why-clause claims the opposite ("any prefix … breaks EXPECT
  substring matching"), which is likely factually wrong for mid-line
  substrings. No operational difference — both specs mandate the direct test —
  but Spec B's reasoning is correct and Spec A's justification should be fixed.
- **Measurement transcription.** Spec B carries the map's per-phase timings and
  the cores-scaling caveat into the Problem section verbatim — better cold-read
  context than Spec A's single aggregate number.
- **Three separate seam diagrams** (runner / phase table / gate.sh entry line)
  are marginally easier to veto one at a time than Spec A's two.

## Resolution

1. `specs/gate-phase-concurrency.md` (Spec A) is the spec of record.
2. Ported into it from Spec B: the corrected why-clause on the inner-mode row
   (canary EXPECT substring matching does *not* catch a prefix — that is why
   the direct byte-shape test is the primary pin), and the per-phase timing
   breakdown plus scaling caveat in Problem.
3. `specs/gate-phase-concurrency-sonnet.md` is kept in-repo as the comparison
   record (reviewer decision); it is not a spec of record and is not staged for
   build.

## Scorecard

| Criterion | Spec A | Spec B |
|---|---|---|
| Line routing per profile + craft-line | ✅ mid for oracle logic, cheap for plumbing | ❌ all-cheap; contradicts the profile's cached routing |
| Concurrency pinned by a permanent seam | ✅ fake-phase timing/marker test | ❌ only the one-off build measurement |
| All Handoff #4 assertables mapped | ✅ incl. verbatim `gate:` line, exit 130 | ❌ verbatim final line has no row |
| Edge inventory: row or Won't-handle, nothing hand-waved | ✅ | ❌ several unbacked "covered by" claims |
| Test economics (no real-gate recursion) | ✅ injected `[]Phase` stubs | ❌ row 1 runs the real gate in a test |
| Canary-EXPECT weakness correctly analyzed | ❌ why-clause likely wrong | ✅ sharp, correct |
| Story-line justification sentences read plainly | ✅ mostly | ➖ story 5 is a stacked-clause pile-up |
| Template structure (sections, diagrams, status line) | ✅ | ✅ |
