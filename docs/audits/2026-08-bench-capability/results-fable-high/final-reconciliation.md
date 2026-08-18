# Bench audit reconciliation — final report

Adjudicator: Claude Code, **Fable 5** (`claude-fable-5`), effort high. Subject:
`58d966e2f92f7f37eba07b6215e8eef45371b72d` — the same commit both prior audits examined
(EXACT MATCH; `baseline.md`). Inputs: Sol (Codex, GPT-5.6 Sol xhigh) `report.md` /
`evidence.md` / `questions.md`; Opus (Claude Code, Opus 5 xhigh) `opus5-xhigh.md.md` +
three evidence logs. Companion files: `reconciliation-ledger.md` (L-ids),
`bug-triage.md` (B-ids), `roadmap-inventory.md`, `roadmap-dispositions.yaml`,
`proposed-roadmap.md`, `principle-control-matrix.md`, `action-items.yaml` (A-ids),
`next-ticket.md`.

**Model continuity.** Recorded at start: the harness system prompt declares Fable 5 /
`claude-fable-5`. Re-checked at close: the same declaration; no fallback, downgrade, or
model-switch message was surfaced by the harness at any point. No runtime surface exposes
a live build id beyond that declaration, so this is what can honestly be recorded.
`INVALID-RUN.md` is therefore not written.

---

## A. Executive verdict

**Bench is directionally correct.** Its deterministic core — the content-addressed gate
verdict keyed on `sha256(tree‖oracle)`, prospective path-scoped commits, ownership fences
and `bench preflight`, `bench coverage --check`, the git and agent-line guards, pooled
worktrees, `bench setup` — held under every probe in this run and in both prior audits.
The prose layer is disciplined (0 vague verbs in 27,758 words) and unmeasured.

**Continue incrementally, with a bounded consolidation.** Every defect reproduced here is
local: an environment variable the ordinary test driver does not carry, an `if` ordered
before a tree comparison, a missing status row, a scaffold file `bench setup` does not
write, a test helper that forgets `t.Setenv("BENCH_HOME")`. No observed failure needs a
new spine, a work-state store, a context compiler, or a claim graph; each of those would
add a second reader of facts Bench already single-sources.

**Strongest demonstrated mechanism:** the tree-keyed gate verdict — reuse refused on
drift, red recorded at its own tree, green correctly reused after revert (verified). Close
behind it: `bench coverage --check` and `bench preflight` (both audits' hostile inputs
caught), and — as practitioner evidence rather than a controlled result — the
repro-loop-first debugging discipline.

**Largest architectural weakness:** the oracle does not grade the kit's own contracts.
The 29-check conformance registry (the thing that makes Bench more than a prose kit)
reaches the live tree only under `BENCH_CONFORMANCE_ROOT`, which the dev gate stopped
setting on 2026-08-09; ten checks are red at HEAD while `bench gate` is green, and the
footer hides the loss as `environment=1`. Second: there is no front door — no surface
routes from observed state, `bench status` says "clean" in an un-adopted repo, a staged
spec produces no row, and `what-next` is maintenance wearing a router's name.

**What happens next:** ship `next-ticket.md` (live-root conformance in the dev gate;
environment skip inside the oracle is red; skips named), then the verdict-reader fix,
then the router. Measure before cutting prose.

## B. Audit relationship

**Agreed (and verified here):** no front door; handoff `## State` stale at the current
commit; 34 tickets-only spec folders; invalid decision map + 62 structure issues coexist
with a green gate; the deterministic substrate is the demonstrated value; `/bench-debug`
preserves the upstream mechanism and must not be compressed further without a benchmark;
Codex/Claude asymmetry in phase invocation and the agent-line hook; guard blocks
read-only text when the core is missing; nothing measures the prose layer.

**Compatible framings:** context volume (Opus 2.3k auto-loaded vs Sol ~6k mandated cold —
both true, different definitions); failure inheritance (Sol: no attempt schema; Opus:
represented in prose shapes with no carrier); repair-loop trigger (Sol: attempt tripwire
in state; Opus: projection over gate records); Codex phase visibility (Sol observed the
effect of the `allow_implicit_invocation: false` policy Opus located).

**Material disagreements:** strategy (Sol strangler vs Opus incremental); work state (Sol
new canonical record vs Opus consolidate); context compiler (Sol yes vs Opus no); claims
(Sol typed graph vs Opus none needed); final-check (Sol merge vs Opus keep); prose quality
(Sol overload vs Opus disciplined). Each is resolved in § E on fresh evidence; every one
resolves toward Opus's position, with the caveat that Opus's *date* for the conformance
loss was wrong (9 days, not 43).

**Where one had stronger evidence:** Opus ran the discriminating experiments for the
conformance gap, the verdict-reader contradiction, the adoption chain, and the un-adopted
status — all reproduced here. Sol ran the gate reuse/red/timeout/moved-subject tests,
the handoff-rewrite clones, and the wrapper/binary probes, and alone caught the release
workflow bypass and the test-fixture leak (whose root cause this run located and found
to be far larger than Sol's framing).

**Where both lacked evidence:** any outcome effect of any prose artifact; review
independence under a withheld commit log; delegate-count value; the actual token/turn
cost of a Bench session. Neither ran model trials (correctly — self-benchmarking is
invalid). Neither noticed FT212 was already fixed or that FT190/FT165 were satisfied.

## C. Reconciliation matrix

| Topic | Sol | Opus | Repository evidence (this run) | Classification | Final disposition |
|---|---|---|---|---|---|
| Conformance in the dev gate | silent | P0: unwired 43 days, 10 red | REPRODUCED: 10 red in 3.2 s; phase table lacks it; removed 2026-08-09 by 3701c4a0 (profile: separate driver retired, coverage lost accidentally) | REPRODUCED (+date correction) | **A1 / next ticket** |
| Red verdicts drift-checked | silent | P0 | REPRODUCED: gate green (reused) vs status "gate red" same tree | REPRODUCED | **A2** |
| Front door / router | logical `bench` + work record + compiler | `bench status --route` + thin `/bench` | REPRODUCED gap on every surface; ask-matt unported | REPRODUCED gap; Sol's extra layers INFERRED | **ADOPT /bench** as adapter over `--route` (A3) |
| Work state | introduce canonical store | consolidate; only repro missing | owners verified; State stale; no repro storage | Sol PARTIALLY SUPPORTED; Opus COMPATIBLE | **MINIMALLY EXTEND** (A8) |
| Context compiler | needed | not needed | 2.3k auto / 7.7k mandated / 27.8k estate; outline 24.6 KB; no observed failure | Sol RESEARCH-INSPIRED; Opus COMPATIBLE | no compiler; measure (A11); fix outline (A9) |
| Claims/evidence | typed graph | none needed | gate = control; structure/maps/conformance = monitoring | Sol UNSUPPORTED; Opus COMPATIBLE | convert the specific monitoring surfaces to control (A1) |
| Strategy | strangler | incremental + one consolidation | all reproduced defects local | Opus COMPATIBLE | **A — incremental** |
| Final check | merge | keep | invokes gate/commit; writes retro (repair table) | Sol UNSUPPORTED | keep |
| Prose quality | overload | disciplined | 0 vague verbs; 17 negations/2.3k | Sol CONTRADICTED | none |
| Prose volume | 31k estate | 3.6× upstream per feature | measured both | COMPATIBLE | measure before FT100 |
| `/bench-debug` fidelity | preserved, diluted | preserved + 3 compressions | textual diff confirms Opus's list | REPRODUCED | **A5** |
| Codex phase visibility | parity failure | policy; debug not self-invocable | config verified | COMPATIBLE | reviewer decision in A5 |
| Adoption gate | not run | cannot go green | REPRODUCED end to end | REPRODUCED | **A4** |
| Test fixture leak | subject polluted (E006) | silent | 759 orphans in ~/.bench/worktrees, +10/gate run; owner located | REPRODUCED (extended) | **A6** |
| Release publish path | P0 bypass | not run | static: raw `npm publish` vs runbook | REPRODUCED (static) | **A7** |
| Wrapper vs binary root | schema failure | help consistent | REPRODUCED (plumbing surface) | REPRODUCED, low impact | with A3 |
| Review independence | unverifiable | real; commit log leaks | `diff --full` appends log | PARTIALLY SUPPORTED | EXPERIMENT G (A11) |
| Delegation counts | ceremony | justified | unmeasured | UNINVESTIGATED | ablate (A11) |
| Repair-loop trigger | attempt ledger | advisory row | per-phase exits exist in gate log | HYPOTHESIS over OBSERVED data | EXPERIMENT (A10) |
| Handoff State stale | E011 | J.3 | REPRODUCED | REPRODUCED | A8 |
| Orphan ticket dirs | E016 | I.6 | 34 dirs; light path has no close step | REPRODUCED | A9 |
| Guard rim / gaps | narrow parser | tokenize; 3 gaps | REPRODUCED (session start; probes) | REPRODUCED | A9 |
| Skip disclosure | "7 skips" | hides environment class | REPRODUCED; reasons retained, unprinted | REPRODUCED | A1 |
| FT212 stale row | — | — | `--landed` works | STALE | archive |

## D. Highest-confidence findings (evidence-backed only)

1. `bench gate` at HEAD is green and the live-root conformance suite is red on 10 checks;
   the dev gate has not graded the live root since `3701c4a0` (2026-08-09).
2. A red gate record is never drift-checked; `bench status` and the handoff pin can assert
   a red the oracle denies for the same tree.
3. `bench setup`'s scaffolded gate cannot go green in a new repo (no `gate-inputs.json`
   → unbound `HOME`; then empty canary inventory).
4. `bench status` cannot distinguish un-adopted from clean, has no staged-spec signal,
   and its action column is half prose; no surface routes from state; `what-next` is
   maintenance and hidden.
5. The kit's own tests leak ~10 pool worktrees per gate run into the operator's real
   `BENCH_HOME` (759 at last count).
6. `/bench-debug` preserves the four load-bearing `diagnosing-bugs` constraints and
   dropped the loop-construction menu, two hard gates, the checkbox completion form, and
   "tighten the loop"; upstream is switched off in this repo; `$bench-debug` cannot
   self-invoke on Codex.
7. The gate's evidence store is a real control (commit refuses on red/stale; reuse
   refused on drift; green correctly reused after revert); `structure`, `maps`, and
   conformance are monitoring only.
8. Prose is disciplined (0 vague verbs; low negation density; enforced budgets) and
   entirely unmeasured; 27,758 words agent-facing; ~2.3k auto-loaded, ~7.7k mandated cold.
9. 34 tickets-only spec folders exist because the light-path route has no close step
   (e.g. `light-path-duplicate-acceptance-ids` landed 2026-08-03; folder remains).
10. Three roadmap rows are stale or satisfied (FT212, FT190, FT165); FT100 sits at rank 1
    with no measurement behind it and its own text says anchor consolidation and a
    benchmark route come first.

## E. Resolved disputes

**1. Strategy — strangler vs incremental.** *Claims:* Sol: new logical entry + work
record + compiler + typed evidence in front of existing modules; Opus: incremental with
one consolidation. *Evidence needed:* is any reproduced failure non-local? *Obtained:*
L-01…L-30 — every reproduced defect is a small local fix; the deterministic core held.
*Experiment:* the sum of the probes (gate reuse/red/revert; adoption chain; status probes;
conformance run). *Result:* no failure requires a spine; a spine would duplicate readers.
*Conclusion:* **A, incremental**, with A2/A3/A8 as the bounded consolidations.

**2. Work state — new store vs consolidate.** *Claims:* as above. *Evidence:* inventory
of owners in the tree; handoff at HEAD; search for repro storage. *Result:* goal, scope,
verified state, open questions all have mechanical owners; the weak spots are prose
`## State` and the missing repro/next-discriminator carrier; three "Next" sources.
*Conclusion:* **MINIMALLY EXTEND EXISTING STATE** (A8) — constrain State, add `Repro:`,
router owns Next.

**3. Front door — `/bench` where?** *Claims:* both want a router; Sol behind a new work
record/compiler; Opus as `bench status --route`. *Evidence:* status signals and their
severity ordinals exist; two signals missing (staged spec, un-adopted); action column
mixed. *Result:* the projection is buildable from existing facts. *Conclusion:*
**ADOPT /bench** as a thin harness adapter (Claude model-invocable; Codex explicit per
policy) over `bench status --route`; rename `what-next` → `drain`.

**4. Context compiler.** *Claims:* Sol needed; Opus not. *Evidence:* word counts;
outline bytes; absence of any observed context-caused failure in three audits.
*Conclusion:* not justified now; instrument and measure (A11); fix `outline` bare form.

**5. Gate vs final-check.** *Claims:* Sol duplicative → merge; Opus complementary.
*Evidence:* `bench final-check` is not a command; the phase invokes `bench gate`/`bench
commit`, reports retained land evidence, writes the retro with the repair-attribution
table. *Conclusion:* complementary; keep.

**6. Claims — control or reporting.** *Evidence:* gate green/red/stale each change
behavior (commit, Stop hook, reuse); structure/maps/conformance reds change nothing.
*Conclusion:* the gate verdict is a control system; the advisory reds are monitoring —
fix by wiring the mechanical ones into the gate (A1), not by adding a claim graph.

**7. Pocock — wrapped vs absorbed vs preserved.** *Evidence:* clause-by-clause diff of
`diagnosing-bugs`; provenance of the other skills as both audits describe (spot-checked:
craft-tdd/craft-review/craft-tickets carry the upstream invariants plus checkable
additions). *Conclusion:* preserved and, in three places, converted into checkable
predicates (coverage rows, citation standard, probe rule); `/bench-debug` compressed —
restore (A5); nothing to absorb or replace.

**8. Failure inheritance — missing vs indirect.** *Evidence:* blocked-delegate return
shape and `craft-line` ladder are prose; gate log carries per-phase exits; no durable
repro carrier. *Conclusion:* represented indirectly with no mechanical carrier; add the
smallest carrier (`Repro:` line, A8) and try the tripwire (A10).

**9. Prose overload.** *Evidence:* rg counts. *Conclusion:* Sol's characterization does
not survive measurement; volume is the only open question and needs A11.

**10. Conformance-loss dating.** *Claims:* Opus "43 days since 72c037a1". *Evidence:*
`git show 72c037a1` added a Go conformance phase; `git show 3701c4a0` removed it
(2026-08-09); profile text. *Conclusion:* nine days; the phase removal was deliberate
(branch-native architecture), the coverage loss was not (profile still claims enforcement;
later occurrences misread the skip). This changes the fix shape (carry the env in the
ordinary driver, don't resurrect a phase), not the priority.

## F. Unresolved hypotheses

| Hypothesis | Why it matters | Cheapest resolving experiment |
|---|---|---|
| The prose layer changes outcomes at all | FT100 and every skill edit hinge on it | A11 arm B vs A on 10–15 retired-spec tasks, ≥5 trials |
| `/bench-debug` transfers `diagnosing-bugs`' repair-loop-breaking effect | protects the one practitioner-evidenced mechanism | A11 arms D/E/F on `Occurrence:`-sourced bugs |
| Withholding the commit log changes review findings | review independence | A11 arm G — one variable, same diff |
| Fixed delegate counts buy anything | orchestration cost | leave-one-out over craft skills; fan-out 1 vs 3 review axes on seeded defects |
| A phase-level repair-loop tripwire reduces repair rounds | trigger for the debug discipline | A10 advisory row vs retro repair table over the next N specs |
| Codex live behavior matches its config (phase adapters hidden; no agent guard) | harness parity | one Codex session opening this repo and listing skills; one `$bench-debug` invocation |
| Bench's cost per outcome vs plain instructions | SWE-Effi's warning | A11 token/turn instrumentation |

## G. "Run the dang test" assessment

*Where Bench already forces observation:* the gate on every landing and shift iteration;
`bench coverage --check`, `bench preflight`, `bench maps` when invoked; `/bench-debug`
Phase 1 (one command already run before any hypothesis); `craft-review` "refute before
you report"; `craft-gate` "prove it bites"; `craft-tdd`'s vacuity and `ok`-hides-a-skip
warnings; `craft-delegate`'s independent probe at a different site and kind.

*Where `diagnosing-bugs` changes the trajectory:* the four constraints that force
contact with reality at exactly the point the model would substitute confidence — a
red-capable loop before theorizing, the exact symptom, one command already run, and the
original-loop rerun after the fix. Practitioner evidence says these break repair loops;
the mechanism is enforced sequencing plus executable feedback, not wording — which is why
Bench's compressions of the *menu* and *hard gates* matter (they are how the sequence is
entered and held) while dropping "be aggressive" does not.

*Where speculation still continues:* the oracle itself, which read `ok
internal/conformance` and stopped (the exact trap `craft-tdd` names); the ambient board,
which asserted a red nobody re-derived; the 2026-08-13 self-assessment, which enumerated
skips and did not name the fifth; ordinary implementation and `craft-line`'s ladder,
which allow a second patch after a red without a repro; review, which forbids tests even
where one would discriminate. And this run: the guard's substring rim cost turns before a
`git status` could be run.

*What generalizes:* one two-sentence doctrine pointer (Opus I.10's wording is good) beside
the four invariants; the mechanical slice — the gate names every skip and treats an
environment skip inside the oracle as red (A1); the advisory tripwire (A10). *What stays
diagnosis-specific:* minimisation, hypothesis ranking, instrumentation tags, seam-of-the-
regression rules. *How to measure:* observation-opportunity delay, speculation-after-
discriminator, blank-retry rate, time-to-red-capable-loop, root-cause-before-fix rate
(A11 § "diagnosing-bugs").

*The candidate principle* ("when a bounded safe executable observation can distinguish
active hypotheses, execute it before continuing speculative repair") is a PROCEDURAL
DEFAULT taught by `/bench-debug`/`craft-review`/`craft-tdd`, pointed to from core doctrine
in two sentences, enforced mechanically only where a check can see the loop (A1's
skip rule; A10's tripwire, advisory first). Not a skill, not a work-state field, not a
gate rule on its own.

## H. Pocock integration verdict (by capability)

| Capability | Bench artifact | Verdict | Why |
|---|---|---|---|
| Diagnosing bugs | `/bench-debug` | **PRESERVE + STRENGTHEN** | four load-bearing constraints intact; restore menu/gates/tighten; settle Codex trigger; benchmark before any further cut |
| TDD | `craft-tdd` | PRESERVE | upstream invariants + vacuity/`ok`-hides-skip additions; the oracle should obey its own warning (A1) |
| Grilling / shaping | `craft-grill`, `/bench-shape-idea` | PRESERVE | facts-vs-decisions and confirmation gate are upstream bug fixes adopted; over-shaping is a routing question the router addresses |
| Spec creation | `craft-spec`, `/bench-write-spec` | PRESERVE; coverage map **WRAP WITH DETERMINISTIC CONTROL** (already exists — wire it: A1) | `bench coverage --check` bites; it is a conformance check that does not run |
| Ticket decomposition | `craft-tickets` | PRESERVE; delegate-per-ticket → BENCHMARK | vertical slicing faithful; mandatory fresh writer unmeasured |
| Implementation | `/bench-implement-spec` | SIMPLIFY (restore its structural headings; label `--full` as prose orchestration) | Bench-native; the file lost its contract headings and landed green |
| Code review | `craft-review`, `/bench-review-implementation` | PRESERVE; BENCHMARK arm G | citation standard and refute-first are real additions; commit-log leak is the one open independence question |
| Delegation / handoff | `craft-delegate`, handoff | PRESERVE; minimal state extension (A8); BENCHMARK counts | probe-differs-in-kind-and-site is the best anti-vacuity clause in the repo |
| Domain modeling | `craft-domain` | PRESERVE | wired into shape-idea (FT165 satisfied) |
| Router (`ask-matt`) | none | **REPLACE** with a state-driven `/bench` | the one upstream capability lost in adaptation; Bench can route from state, not prose |
| Prototype | `prototype` | PRESERVE | direct empirical escape |
| Skill authoring / synthesis / update-kit | kit-maintainer skills | PRESERVE, keep out of linked-repo cold context (already excluded from the payload) | — |

Consistency gains survive across harnesses only for the deterministic layer (same binary,
same checks); the prose triggers differ by policy (`allow_implicit_invocation: false` on
Codex; six model-invocable phases on Claude) and by hook coverage (agent-line guard
Claude-only). The gains are preserved and Bench simplified by *wiring what exists* (A1),
*projecting what exists* (A3), and *measuring before cutting* (A11).

## I. Entry-point verdict

**ADOPT /bench** — as a thin adapter, not a new subsystem.

```
state inspection   bench status (existing signals + `specs: staged` + `setup`)
router             bench status --route → next[1]{state,why,command} (+ runners-up, one line)
expert side doors  /bench-shape-idea · /bench-write-spec · /bench-debug · /bench-implement-spec
                   /bench-review-implementation · /bench-final-check (unchanged, direct)
maintenance        /bench-drain (was what-next) · assess · setup-repo · update-kit
Claude adapter     .agents/commands/bench.md, model-invocable; loads only the routed phase
Codex adapter      $bench, allow_implicit_invocation: false (policy), same route text
shell              bare `bench` prints the route; `bench help` prints the inventory;
                   binary no-arg/help aligned with the wrapper
```

State is inspected before any question; resume reads the `intent` row and the handoff
`## Next command`; ambiguity is answered by the highest-severity actionable row and named
runners-up, never by a question the tree can answer. Context load: one row plus one phase
file. Duplication: none — the router is a projection, and Sol's "expose why it chose the
route" is the `why` column.

## J. Target architecture

- **Engineering-practice skills** (judgment): craft-* and the phase files; unchanged in
  role; `/bench-debug` restored; measured by A11 before any cut.
- **Deterministic Bench core**: gate (now grading the kit's contracts too), commit, land,
  worktree/leases, preflight, coverage, guards, adopt/setup, status, roadmap/capture
  readers, release state machine.
- **Current work state**: git + spec `Status:`/coverage rows + fences + decision maps +
  intent ledger + handoff pin block; extended by a constrained `## State` and a `Repro:`
  line; Next projected once by `status --route`.
- **Claims/evidence**: the content-addressed gate record (control) — no claim graph;
  advisory surfaces either become gate checks (conformance) or stay explicitly
  on-demand (structure — reviewer decision).
- **Context selection**: progressive disclosure as today; outline bare form fixed;
  instrumentation added; compiler explicitly not built.
- **Gates**: dev gate = gofmt · vet · test (with live-root conformance) · race · system ·
  shellcheck; environment skips inside the oracle red; skips named; ship = prep-release.
- **Semantic review**: three fresh axes, citation standard; commit-log question settled
  by arm G.
- **Claude adapter**: settings hooks (session-start, stop, git guard, agent-line),
  `.claude/skills` symlinks, model-invocable phase policy mirrored by a conformance check.
- **Codex adapter**: `.codex/hooks.json`, openai.yaml policy (gate-enforced once A1
  runs), `$bench`; agent-line parity parked (FT24).
- **Model-specific profiles**: `.bench/lines.env` + profile `Lines`; effort/cap stay
  profile prose until measured.
- **Benchmark suite**: A11 harness + A10 tripwire measurement; task corpus from retired
  specs and `Occurrence:` bugs.

## K. Roadmap verdict

Remains active (as the twelve items): A1–A12 in `proposed-roadmap.md`. Removed from the
active board: FT212, FT190, FT165 (done); FT100 demoted behind measurement (A11) with
FT170/FT138 merged into A11. Rewritten before working: FT98, FT169, FT164, FT158, FT144,
FT199/FT206, FT197, FT174, FT130, FT172, FT141 (name collision). Fixed now inside the
portfolio: the ten conformance diagnostics, the verdict reader, the test leak, the
adoption scaffold. Explicitly not worth fixing: wrapper/binary root disagreement on its
own (align with the router), prose "quality" (not a defect), Sol's E006 as stated
(superseded by the real leak). Everything else: `roadmap-dispositions.yaml` — deferred
and batched by file so each lands under one gate.

## L. Principle / control verdict

Seven core invariants (`principle-control-matrix.md`): the existing four (gate is the
oracle; declare the line; document current state; one small change / one source per
fact) plus three made explicit — **C5** a rule a program can check is checked by a program
in the dev gate (owner: conformance registry in the gate, coverage/preflight, guards);
**C6** work state lives in the tree and survives reset/harness/model swap (owner: git +
spec status + handoff pin block + router Next + lines.env); **C7** verification claims
disclose what they did not grade (owner: gate skip disclosure by name; final-check
report; review citations). Empirical escape, failure inheritance, review independence,
smallest context, and "more agents is not better" stay procedural defaults/heuristics with
named teaching skills; the tripwire and the measurement harness are experiments; the
work-state file, compiler, and claim graph are rejected doctrine pending A11 evidence.

## M. Keep / Change / Kill

**Keep (evidence: held under probes here and in both audits):** tree-keyed gate
evidence and reuse; `bench commit`/`bench worktree land`; ownership fences + preflight;
`bench coverage --check`; git guard live surface (incl. `restore --staged` allowed);
`check-agent-line` (denies omitted model); `bench setup` plan-first UX; pooled worktrees
and plan-then-apply cleanup; `/bench-debug`'s four constraints; `craft-review` citation
standard; `craft-delegate` probe rule; retro repair-attribution table; final-check as
close-out; TOON query plane; tier binding design.

**Change (evidence: reproduced defects):** gate test phase carries the conformance root
and names skips (A1); verdict staleness rule (A2); status → route, two signals, invocable
actions, `/bench`, what-next → drain (A3); setup scaffold + wrapper HOME message + canary
line (A4); `/bench-debug` restorations + Codex trigger + Claude parity check (A5); worktree
test isolation (A6); release publish path (A7); handoff `## State` constraint + `Repro:`
(A8); hygiene batch (A9).

**Kill (evidence: stale, inert, or unsupported):** the inert `disable-model-invocation`
key on Codex adapters; the 34 tickets-only folders (after review) and the light path's
missing close step; four stale structure-accept rows; the stale "built-in conformance
phase" prose; `bench outline`'s bare 200-row dump; the substring guard rim; the idea of a
canonical work-state file / context compiler / claim graph / review runner / strangler
spine (A12); FT100's rank-1 position pending measurement.

## N. Final action set

A1 P0 conformance graded in the dev gate · A2 P0 verdict staleness rule · A3 P1 `/bench`
router · A4 P1 adoption smoke · A5 P1 `/bench-debug` restore + trigger + parity · A6 P1
test-isolation leak · A7 P1 governed publish path · A8 P2 minimal handoff/work-state
extension · A9 P2 hygiene batch · A10 EXPERIMENT repair-loop tripwire · A11 EXPERIMENT
measurement harness (blocks FT100) · A12 DELETE the spine/store/compiler/graph proposals.
Full fields in `action-items.yaml`; identical to `proposed-roadmap.md`.

## O. Immediate next ticket

**A1 — the dev gate grades the kit's own contracts** (`next-ticket.md`). It outranks A2
because it changes what the oracle sees rather than what a reader reports (and would have
caught A2's doc-drift class); outranks A3 because the router's new phase file, rename,
and adapters are exactly what the registry grades — land it first or land the router
ungraded; outranks Sol's CR-001 tracer because that is INFERRED, large, and duplicates
readers; outranks FT100 because a prose cut with no measurement is the highest-variance
action on the board. It is small (env plumbing, skip rendering, ten dispositions), red
first at HEAD, mutation-probed in two kinds, and it unblocks four other items.

## P. Benchmark plan

Arms: **A** model + `AGENTS.md` only; **B** current Bench at `58d966e2`; **C** Bench after
A1–A9. Debug arms: **D/E** current Bench without/with `/bench-debug`; **F** upstream
`diagnosing-bugs` vs `/bench-debug`. Review arm **G**: commit log withheld vs supplied.
Same commit, task set, model, effort, tool permissions, budgets; ≥5 trials per arm per
task; both harnesses where a Codex run is possible; medians and spread; pre-registered
criteria committed before the first scored run (B beats A on requirement satisfaction and
epistemic quality at ≤3× tokens; E beats D on first-fix success and time-to-root-cause;
C beats B on tokens at equal outcome; G decides the log). Tasks: 10–15 reconstructed from
retired specs (`git log --diff-filter=D -- specs/`, parent-of-landing checkout, coverage
rows as hidden requirements) plus bugs from `roadmap/FT*.md` `Occurrence:` lines.
Measures, per the brief's list: outcome (task success, requirement satisfaction, hidden
defects, CI, reviewer findings); epistemic (unsupported assertions, stale evidence, missed
validation, false completion); context (tokens, files read, duplicate reads, repeated
searches, irrelevant reads); execution (duplicate commands, repair-loop count, blank
retries, observation-opportunity delay, speculation after discriminator, failure
inheritance); recovery (context-reset and model-swap resume, checkpoint completeness,
repeated dead ends); orchestration (subagent count, duplicate discovery, overlapping
modifications, conflicts, orchestrator overhead); entry (ambiguity, steps to correct
workflow, unnecessary questions, incorrect routing); `diagnosing-bugs` (steps to first
repro and first discriminating observation, repair-loop count, root-cause accuracy,
first-fix success, regression validation, tokens/tool calls). Leave-one-out over the 16
craft skills. Instrumentation lands with A11 (FT138); the retro repair table is already a
channel.

## Q. Ten hard conclusions

1. `bench gate` is green at HEAD while ten registered contract checks are red; the dev
   gate has not graded the live root since 2026-08-09, and the footer hides it as
   `environment=1`. (REPRODUCED)
2. A Bench build removed a phase file's required headings on 2026-08-11 and landed green;
   the check that forbids it exists and did not run. (REPRODUCED)
3. `bench gate` and `bench status` disagree about the same tree in the same second because
   red records are never drift-checked; single-sourcing propagated the wrong rule to the
   handoff faithfully. (REPRODUCED)
4. A newly adopted repo cannot reach a green gate for two kit-caused reasons no surface
   names. (REPRODUCED)
5. The kit's own tests write ~10 orphan worktrees into the operator's real `BENCH_HOME`
   per gate run; 759 were present. (OBSERVED)
6. There is no front door: bare `bench` is a 44-line inventory, an un-adopted repo reads
   "clean", a staged spec has no row, and the router upstream ships (`ask-matt`) was never
   ported. (REPRODUCED)
7. `/bench-debug` keeps every load-bearing `diagnosing-bugs` constraint and dropped the
   menu, two hard gates, and "tighten the loop"; on Codex it cannot self-invoke; upstream
   is off in this repo. (REPRODUCED textually; PRACTITIONER EVIDENCE for value)
8. The prose is disciplined and unmeasured: 27,758 words, zero vague verbs, ~3.6× the
   upstream pipeline's weight per feature, no A/B evidence — and the rank-1 roadmap item
   proposes cutting it by judgment. (OBSERVED / UNINVESTIGATED)
9. Every disagreement between the two audits that could change architecture resolves, on
   fresh evidence, toward the smaller change: incremental over strangler, extend over new
   store, project over compile, wire over claim-graph. (OBSERVED)
10. Bench's problem is not that it built too much; it is that the most valuable thing it
    built — mechanical enforcement of its own guidance — is switched off, and no surface
    says so. Wiring it is one small ticket, and it is the next one. (REPRODUCED)

---

## Run close — state and cleanup

```
$ git status --short
 M .claude/README.md      (pre-existing, untouched — see baseline.md)
 M .claude/settings.json  (pre-existing, untouched)
?? audit/                 (inputs + this run's outputs under audit/fable-high/)
```

No tracked file was modified by this run; nothing was committed, pushed, or released.
Temporary experiments and their side effects were removed: the gofmt mutation of
`internal/toon/toon.go` was restored from a cp-aside copy (`git diff -- internal/` empty);
the throwaway `git init` repo under the session scratchpad was deleted; the stale red
`bench-last-gate` record left by the red/revert experiment was reconciled with one final
`bench gate --fresh` (green); the 30 orphan pool fixtures that this run's four gate runs
wrote into `~/.bench/worktrees` (ledger L-06 — the leak itself is a Bench defect, item A6)
were removed, returning the pool to its pre-run 739 entries. Remaining artifacts outside
the tree: `dist/bench` (gitignored local build, needed for the CLI to resolve in this
worktree) and four gitignored `.logs/gate-*.jsonl` progress logs from this run's gate runs.

Model identity at close: unchanged — Fable 5 (`claude-fable-5`), no fallback reported.
