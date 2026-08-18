# Principle → control matrix

Rule applied: a CORE INVARIANT must have a concrete enforcement owner that exists (or is
in the twelve-item portfolio); a PROCEDURAL DEFAULT names the skill that teaches it;
harness rules stay in adapters; model behavior stays in the profile; unproven ideas stay
experiments. Bench's four existing invariants are kept verbatim as the base — they already
have owners that held under every probe in this run — and the ten candidate principles
from the brief are classified against them rather than layered on top.

## Classification of the ten candidates

| # | Candidate | Classification | Where it lands |
|---|---|---|---|
| 1 | Uncertainty may remain unresolved, but must remain visible | PROCEDURAL DEFAULT + one CLI owner | `craft-review` (dispositions), `craft-delegate` (blocked-report shape), handoff `## State`; mechanical slice: the gate names every skip (P0 ticket) |
| 2 | Observation and interpretation are different artifacts | Already CORE — this is invariant 1's content ("a done-claim is a claim, not a result"; the gate record is the observation) | gate evidence record; `bench commit`; Stop hook |
| 3 | When a safe executable observation can distinguish hypotheses, prefer it over speculation | PROCEDURAL DEFAULT (+ EXPERIMENT for a mechanical trigger) | `/bench-debug` Phase 1, `craft-review` "refute before you report", `craft-tdd` vacuity check, `craft-gate` "prove it bites"; two-sentence doctrine pointer in `.bench/BENCH.md`; tripwire = A10 |
| 4 | Deterministic invariants belong in CLI/schema/hook/gate, not model memory | CORE INVARIANT (new C5) | conformance registry running in the dev gate (A1); guards; `bench coverage --check`, `bench preflight` |
| 5 | Active context should be the smallest sufficient working set | HEURISTIC (measure before enforcing) | progressive disclosure already in `.claude/settings.json` skill policy + `craft-delegate` compressed inputs; `bench outline` fix; A11 measures |
| 6 | Failed attempts carry evidence into later attempts | PROCEDURAL DEFAULT (+ EXPERIMENT for a carrier) | `craft-line` ladder (classify reds, name the changed variable), blocked-delegate return shape, retro repair table; carrier = A8 `Repro:` line + A10 tripwire |
| 7 | Work survives context reset and model replacement | CORE INVARIANT (new C6) | git + spec `Status:` + coverage rows + `lines.env` tier binding + handoff pin block; enforced by `handoff-shape-single-source` conformance (once A1 runs it) and the status handoff-age row |
| 8 | Completion claims state and support their verification coverage | CORE INVARIANT (new C7) — the coverage half of invariant 1 | gate skip disclosure by name (A1), `bench worktree land` gating the composed tree, final-check report shape, `craft-review` citation standard |
| 9 | Independent review stays genuinely independent | PROCEDURAL DEFAULT + EXPERIMENT | `craft-review` parallel fresh axes / re-derive-then-compare; arm G tests the commit-log leak |
| 10 | More agents are not inherently better | HEURISTIC | `craft-delegate` (each boundary must buy isolation, rediscovery, or bounded research); leave-one-out in A11 |

Rejected doctrine (from the audit inputs): a J-Space-style Goal/Core/Verified/Open/Next
file as a new source of truth; a general claim/evidence graph; a fixed agent-count rule
either way. None has an observed failure behind it (ledger L-04, L-24, L-25).

## Core doctrine — seven durable invariants with owners

```yaml
- id: C1
  statement: The gate is the oracle; "done" means the external gate exits zero on the exact tree, and no agent grades its own work or weakens a check to pass.
  classification: CORE INVARIANT (existing invariant 1)
  scope: every landing, every shift iteration, every delegate done-claim
  failure_mode_prevented: self-certified completion; stale evidence authorizing a commit; a green that grades a different tree
  owner: Bench core (internal/gate evidence store keyed sha256(tree‖oracle); bench commit; bench worktree land) + Stop hook under BENCH_SHIFT
  enforcement: mechanical — commit refuses on red/stale; reuse refused on drift (verified this run); Stop hook blocks
  evidence: OBSERVED — reuse/red/revert sequence; both audits' attack tables
  regression_test_or_benchmark: internal/gate verdict tests; the three-surface agreement test from A2
  exceptions: none; `--fresh` only escalates

- id: C2
  statement: Declare the line before a long run — model tier, effort, iteration cap — and tiers bind to opaque ids per harness; an unavailable model returns to the binding, never a substitute.
  classification: CORE INVARIANT (existing invariant 2)
  scope: every delegate spawn, every headless shift
  failure_mode_prevented: silent tier escalation; a fork inheriting an unbound model unnoticed
  owner: check-agent-line hook (Claude PreToolUse:Agent) + shift adapters; profile `Lines` table; `.bench/lines.env`
  enforcement: mechanical on Claude (denies unbound/omitted model — verified by Opus, 10 probes); Codex parked (FT24); effort/cap remain prose (MODEL PROFILE, not enforced)
  evidence: OBSERVED (guard); iteration cap/effort HYPOTHESIS
  regression_test_or_benchmark: hook envelope probes; A11 records tier per trial
  exceptions: fork delegations (warned, cannot be denied)

- id: C3
  statement: Documents state the current decided state for a reader with no memory of how it got there; docs, ADRs, phase files, and adapters keep their structural contracts.
  classification: CORE INVARIANT (existing invariant 3)
  scope: AGENTS/BENCH, phase files, ADRs, profile, handoff
  failure_mode_prevented: stale prose steering a cold session (observed at HEAD: handoff State; BENCH-reference "conformance phase"; review doc)
  owner: conformance registry (docs-currency-workflow, handoff-shape-single-source, skills-index-command-adapters, decision-map-integrity, guidance-prose-budgets)
  enforcement: mechanical ONLY once the registry runs in the dev gate (A1) — today prose-only, and 10 red
  evidence: OBSERVED (10 diagnostics; drift landed green in fa4e1f02)
  regression_test_or_benchmark: planted heading deletion reds the gate (A1 acceptance)
  exceptions: reviewer-approved contract changes update the check in the same commit

- id: C4
  statement: One small change at a time on a green tree; one source per fact; read before calling; compose an existing seam.
  classification: CORE INVARIANT (existing invariant 4)
  scope: every diff
  failure_mode_prevented: drift between two derivations (the verdict-reader bug shows the cost when the one source is wrong — it is still the right rule); destructive git by reflex
  owner: gate + bench commit (path-scoped, atomic) + block-dangerous-git guard + craft-seams/craft-tickets (judgment half)
  enforcement: mechanical for landing and destructive git; prose for seam choice
  evidence: OBSERVED (guard probes; prospective commit path — both audits)
  regression_test_or_benchmark: guard corpus; A2's regression test as the one-source example
  exceptions: independently authored test expectations (documented in AGENTS.md)

- id: C5
  statement: A rule that a program can check is checked by a program in the dev gate, not remembered by the model — the kit's own registered checks included.
  classification: CORE INVARIANT (candidate 4, made explicit)
  scope: coverage-map validity, decision-map integrity, doc/adapter contracts, prose budgets, AXI registry, ownership fences, guard surface
  failure_mode_prevented: monitoring-without-control (structure/maps/conformance red while the gate is green — observed at HEAD)
  owner: bench gate phase table (A1) for the registry; bench coverage --check / bench preflight for spec entry; guards for git/model
  enforcement: mechanical after A1; preflight remains prose-invoked until FT200 decides the chokepoint
  evidence: OBSERVED (L-01, L-18, L-24)
  regression_test_or_benchmark: A1's environment-skip-is-red rule; conformance suite red-then-green at HEAD
  exceptions: judgment checks (seam choice, story breadth, finding severity) stay prose by design

- id: C6
  statement: Work state lives in the tree, not the transcript: goal (spec/ticket), scope (fences), verified state (gate record for the exact tree), open questions (decision maps), and next action (one projection) survive a context reset, a harness swap, and a model swap.
  classification: CORE INVARIANT (candidate 7)
  scope: every phase close and every delegate return
  failure_mode_prevented: a resuming session re-deriving or, worse, trusting stale prose; a bug's repro loop rebuilt from scratch
  owner: git + spec Status/coverage + handoff pin block (regenerated) + status/router Next projection (A3) + `.bench/lines.env`
  enforcement: mechanical for pins and Status; handoff-shape conformance (after A1); `## State` constrained to tree-contradictable facts + `Repro:` line (A8) — the one extension
  evidence: OBSERVED (pins recover; State drifts; repro not stored — L-04/L-07)
  regression_test_or_benchmark: A11 recovery measures (checkpoint completeness, cross-model resume)
  exceptions: conversation transcript, implementer rationale, delegate internal reasoning deliberately do not survive

- id: C7
  statement: A completion or verification claim discloses what it did and did not grade — tier, skipped checks by name, capability limits — so a local green is never read as global completion.
  classification: CORE INVARIANT (candidates 1 + 8, the coverage half of C1)
  scope: gate output, final-check report, delegate done-claims, review findings
  failure_mode_prevented: the 43-day/9-day blind spot: an "environment=1" footer that hid a whole suite; a delegate's focused green read as done
  owner: bench gate skip disclosure (A1: name test + reason; environment class inside the oracle = red); final-check report shape; craft-review citation standard; craft-delegate six-item verification
  enforcement: mechanical for gate skips (A1); prose for review/delegate claims
  evidence: OBSERVED (footer at HEAD; the 2026-08-13 assess walked past it)
  regression_test_or_benchmark: gate output test naming the skipping test; A11 "false completion" measure
  exceptions: capability skips (fifo/privilege) remain footnotes — they are host facts, not invocation bugs
```

## Procedural defaults, heuristics, adapter rules, profile, experiments

| Principle | Class | Owner (teaching skill / layer) | Enforcement today | Note |
|---|---|---|---|---|
| Empirical escape: build the red-capable loop before theorizing; retry only with new evidence | PROCEDURAL DEFAULT | `/bench-debug`, `craft-review`, `craft-tdd`, `craft-gate`; two-sentence pointer in `.bench/BENCH.md` | none mechanical | Do not add a skill; restore `/bench-debug`'s hard gates (A5); tripwire is A10 EXPERIMENT |
| Failed attempts inherit command/observation/ruled-out/changed-variable | PROCEDURAL DEFAULT | `craft-line`, `craft-delegate` blocked-report shape | none | carrier: A8 `Repro:` line |
| Smallest sufficient context | HEURISTIC | `craft-delegate` (compressed inputs), harness skill policy, `bench roadmap --row` | prose budgets (conformance, after A1) | measure (A11) before FT100 cuts |
| Independent review = fresh axes re-deriving from primary sources | PROCEDURAL DEFAULT | `craft-review` | none | arm G decides whether to withhold the commit log |
| Delegate only for isolation, independent rediscovery, or bounded research | HEURISTIC | `craft-delegate`, `craft-line` | fan-out/effort/cap unenforced | ablate in A11 |
| Phases are reviewer-chosen entry points on Codex; six are model-invocable on Claude | HARNESS ADAPTER RULE | `agents/openai.yaml` (checked), `.agents/commands` frontmatter (unchecked) | Codex side conformance (after A1) | add the Claude-side mirror check; decide `$bench-debug` (A5) |
| Git guard classifies destructive spellings; fail-closed when the core is missing | HARNESS ADAPTER RULE | `.bench/hooks/block-dangerous-git.sh` + `internal/gitguard` | mechanical | narrow the degraded rim (A9) |
| Tier ids per harness; effort/cap per work class | MODEL PROFILE | `projects/benchkit.md` Lines, `.bench/lines.env` | tier: hook; effort/cap: none | keep out of universal doctrine |
| Repair-loop tripwire from gate records | EXPERIMENT | A10 | — | advisory first |
| Skill/prose marginal value; per-skill ablation | EXPERIMENT | A11 | — | blocks FT100 |
| Canonical work-state file; context compiler; claim graph; executable review runner | REJECTED DOCTRINE (for now) | — | — | no observed failure; revisit only on A11 evidence |
