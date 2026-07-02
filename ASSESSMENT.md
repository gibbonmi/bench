# Tooling assessment — 2026-07-02

Full assessment of the Bench kit, run at the reviewer's request. Three audit passes
(ad-hoc validation calls, craft skills, workflow phases), compared against roadmap
items 3 (guard self-disclosure) and 4 (learnings counter bug), then prioritized.
Findings below are the complete record; the roadmap holds the deprioritized work.

## Part 1 — ad-hoc check/validate calls → CLI candidates

### The core observation

`bench status` already computes seven signals as structured `severity|signal|detail|action`
tuples — gate verdict + staleness, git dirty/ahead, active worktrees, open learnings,
structure violations, unresolved decision maps, roadmap count — then flattens them into
a 5-row text budget and discards the detail. Three phases (`/bench-integrate-learnings`,
`/bench-shape-idea`, `/bench-implement-spec`) re-derive that same detail with
hand-assembled greps. The data exists; the query surface doesn't. First-wave work is
exposure, not new checkers.

### Ad-hoc patterns found (where agents assemble shell checks by hand)

| # | Pattern | Instructed in | Status |
|---|---|---|---|
| A1 | Branch-diff pinning: resolve base, `git diff <base>...HEAD`, confirm non-empty | review-implementation | no helper; base resolution wrong for shift branches |
| A3 | Partial gate: typecheck + single test file "frequently as you go" | implement-spec | agent reconstructs the stack; can drift from the gate |
| A4 | Stale-reference sweep after rename: `/name`, `$name`, basename, `dir/name` forms | implement-spec | assembled multi-pattern grep |
| A5 | Acceptance-coverage table: hand-parse spec rows, re-emit per-row status | write-spec, implement-spec, review-implementation | pure prose discipline, no deterministic backstop |
| A6 | Working-tree state before commit (no blind `git add -A`) | implement-spec, BENCH.md | logic already exists internally (`shift_dirty_paths`) |
| A7 | Link/init wiring preflight | setup-repo | no read-only "am I linked / drifted" query |
| A8 | Stack/remote/profile discovery | setup-repo | duplicates `run_gate`'s auto-detect logic |
| A9 | Model discovery fallback (harness-scraping when no API key) | setup-repo, craft-line | fragile; no contract when nothing resolves |
| A10 | Roadmap promotion drain (hand-edit `ROADMAP.md` line removal) | shape-idea | gate-enforced but hand-performed |
| A11 | Decision-map placeholder scan (`— (open` / `GRILL DEFERRED`) | shape-idea, craft-grill | `status` runs the same grep, discards per-map detail |
| A12 | Learnings-journal entry counting/grouping | integrate-learnings | `status` counts; entry list is a manual read |
| A13 | Structure violations | craft-seams | `status` re-greps `structure`'s own text output |
| A14 | Kit staleness audit (drift, cross-refs, stale paths) | craft-synthesis | much already in gate-docs-contracts; skill says "grep" |

### Candidate subcommands

**First wave (prioritized — exposure of existing tuples):** `bench status --json` /
per-signal query; `bench learnings` (structured open-entry list); `bench maps`
(per-map unresolved tickets). Shared parsers, one code path, deletes ad-hoc greps
from four agent-facing files.

**Second wave (parked):** `bench diff` (A1), `bench refs <stem>` (A4),
`bench coverage <spec>` (A5), `bench doctor` (A7/A14), `bench detect` (A8);
extend `bench structure` with a structured mode, `bench roadmap remove`,
scoped `bench gate`, structured `bench models`.

### The decision this reopens

The kit currently exempts `bench` itself from AXI ("plain text, stderr errors,
documented exit codes" — closed in `craft-cli`'s scope clause and the benchkit
profile). Making the new query surfaces AXI/TOON reopens that. Recommended
resolution: hybrid — existing text surfaces keep their contract; new query
surfaces are AXI-conformant. Reviewer's call; belongs in the shaping session.

## Part 2 — craft-skill findings

| Skill | Finding | Severity |
|---|---|---|
| craft-skills | Model-invoked census omitted `craft-line` (meta-standard mis-counting itself) | High — **fixed 2026-07-02** |
| craft-cli | No contrastive pair despite governing an output surface (the meta-standard's own rule) | High — **fixed 2026-07-02** |
| craft-cli | Errors-on-stdout stated without the why; agents may "fix" it back to Unix convention | Med — **fixed 2026-07-02** |
| craft-design-system | No contrastive pair for the token rule | Med — **fixed 2026-07-02** |
| craft-design-system | Names Claude Design in a harness-agnostic skill; "stop" on missing token has no headless unblock route | Med — parked (polish batch) |
| craft-tdd / write-spec / review-implementation | 7-edge-class list triplicated verbatim; a change is a three-place edit | Med — parked (polish batch) |
| craft-seams | Fuses two jobs: seam placement plus structure-gate splitting; splitting guidance can't fire on a `bench structure` failure (triggers don't match) | Med — parked (polish batch) |
| craft-line | Token cap has no sizing method (vibes number; the skill warns against those elsewhere) | Med — parked (polish batch) |
| craft-synthesis | Dogfood loop has no proportionality carve-out — read literally, a typo fix demands a full gated shift | Med — parked (polish batch) |
| craft-grill | Good-only example (no bad pair); trigger collision with ecosystem grill skills | Low — parked |
| craft-adr | Cleanest of the nine; no material findings | — |

### Missing skills (ranked)

1. **craft-gate** — authoring a strong oracle. Invariant 1 rests entirely on the
   gate, yet gate-authoring guidance lives only in run-once setup-repo. Nothing
   fires when a check is added, tightened, or weakened. Should cover: oracle vs
   theater, red-by-construction avoidance, paired-delta harnesses, the canary
   pattern.
2. **craft-review** — the three-axis judgment review-implementation demands
   (hard-violation vs judgment-call sorting, adversarial coverage, aggregate
   don't merge) has no model-reachable home.
3. **craft-delegate** — delegate prompts, scope bounding, worktree isolation,
   done-claim verification; currently smeared across four artifacts.
4. **craft-structure** (or re-triggered craft-seams) — splitting under the
   structure gate.
5. **craft-learnings** — promotable-entry quality; could fold into craft-synthesis.

## Part 3 — workflow phase findings

| # | Finding | Severity |
|---|---|---|
| P1 | **Per-story lines assigned in write-spec are dropped by implement-spec** — it declares one blanket line, never reads the spec's per-story routing; `bench shift` carries a single `BENCH_MODEL` so mixed-line specs can't shift as one unit and nothing says to partition. Direct invariant-2 contradiction. | High — prioritized (spec path) |
| P2 | implement-spec called the review "two-axis"; it is three (Coverage silently dropped at the handoff) | High — **fixed 2026-07-02** |
| P3 | No phase routes a capped/unmet shift; review-implementation's base (default branch) is wrong for shift branches (pre-shift HEAD); no exit for retiring a superseded spec | High — parked (workflow-exits entry) |
| P4 | Coverage-row status is entirely self-assessed; the gate's coverage check is structural (anchor text), not semantic (rows → passing tests) | Med — feeds `bench coverage` candidate |
| P5 | setup-repo model discovery is harness-scraping with no failure contract | Med — parked |
| P6 | Batch approval (multi-spec roadmap run) lives only in BENCH.md prose; no phase surface | Low — parked |
| P7 | bench-debug: repro test must be committed into the gate *before* `bench shift`, else the first red-iteration rollback destroys it — implied, not stated | Low — parked (workflow-exits entry) |
| P8 | Headless shift iteration prompt drops the coverage-row discipline implement-spec mandates — interactive and headless discipline diverge for the same spec | Low — parked |
| — | Adapter skills (`$bench-*`) carry no drift; the Claude/Codex invocation split is internally consistent | clean |

## Roadmap comparison and priority order

- **Roadmap 3 (guard self-disclosure)** is the same family as Part 1: a
  machine-readable state surface injected at SessionStart. Shape together as one
  structured-query surface (`status --json`, `learnings`, `maps`, `guards`).
- **Roadmap 4 (learnings counter)** — **fixed 2026-07-02** (entry-heading match +
  fixture alignment + template-line regression check); the same parser seeds
  `bench learnings`.

Priorities:

1. ~~Trivial-fix batch~~ — done (P2, craft-skills census, roadmap 4, contrastive pairs).
2. **Shape the structured state surface** (`/bench-shape-idea`) — merges Part 1
   first wave with roadmap 3; closes the AXI-exemption question.
3. **Per-story line inheritance** (P1) — spec path; has design weight (partitioning
   mixed-line specs for shift).
4. Everything else parked on the roadmap: second-wave parsers, missing-skills trio
   (craft-gate first), skill-polish batch, workflow exits.
