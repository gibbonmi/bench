# Slim BENCH.md to 150 and land the AGENTS.md shell conventions

Blocked by: 06-review-rederive-doctrine.md
Ownership fence: `.bench/BENCH.md`, `AGENTS.md`, `internal/anchors/registry_data.go`, `tests/canary/workflow-guidance-anchors`, `.bench/BENCH-reference.md`
Integration surfaces: BENCH anchors→`internal/anchors/registry_data.go`; anchor canaries→`tests/canary/workflow-guidance-anchors`; marker-phrase gate→`AGENTS.md` (shared-rule literal markers must not reappear); reference lookup→`.bench/BENCH-reference.md`
Contracts: the shared-rule marker phrases crossing `.bench/BENCH.md`→`AGENTS.md`, asserted by BA4 against the real gate's substring check
Closure: BA1/bench-150, BA2/contradiction-and-convergence, BA3/agents-shell-rules, BA4/gate-and-anchors-green, BA5/source-before-call, BA6/acceptance-shortfall

## What to build

Rewrite `.bench/BENCH.md` to ≤150 lines. Keep: roles, the CLI inventory, the
four invariants, the communication rules, the workflow with proportionality
and light-path table, capture, and the phase-close handoff pointer — each in
tighter form, moving branch-only detail to `.bench/BENCH-reference.md` where
it already lives. Add the surviving operational predicates this file owns:
no API or function is called before its definition has been read this session
(strengthening "read the surrounding code"); a non-behavioral spec
contradiction follows the current tree convention and is flagged for veto
while behaviorally different readings stop; a material acceptance shortfall
exits and reports rather than silently landing; only diff-owned reds count
toward fix-loop convergence (pointer — `craft-line` owns the classification).
Preserve the 17 BENCH anchors' surviving obligations (worktree
release/clean, upgrade, IDEAS append format, independently-green rule,
final-check retro and what-next drain pointers, "NEVER assume, always
verify"); migrate or retire pins with their canaries where wording moves. In
`AGENTS.md`, extend the shell-conventions section with the four FT107 rules:
wait on a PID or sentinel, never a self-matching pattern; destructive scripts
run plan-before-apply with exact sampled targets; repository-wide sweeps use
`rg --hidden` excluding `.git`; discover Bench verbs non-interactively via
`bench commands --brief` or source, never a bare unknown verb, with
non-interactive stdin. Do not let any shared-rule literal marker phrase from
BENCH.md reappear in AGENTS.md.

## Acceptance

- [ ] [BA1] (covers local) `.bench/BENCH.md` is ≤150 lines with roles, invariants, communication rules, workflow, CLI inventory, and capture intact (PG13's composed verdict is owned by ticket 09's PB1).
- [ ] [BA2] (covers PG20) the non-behavioral-contradiction predicate (tree-consistent reading flagged for veto, behaviorally different readings stop) and the owned-red-convergence pointer are present.
- [ ] [BA3] (covers PG22) `AGENTS.md` states all four shell rules.
- [ ] [BA5] (covers local) the source-before-call predicate is present in BENCH.md (PG19's composed verdict is owned by ticket 08's SD2).
- [ ] [BA6] (covers local) the material acceptance-shortfall exit is present in BENCH.md (PG21's composed verdict is owned by ticket 08's SD1).
- [ ] [BA4] (covers local) the gate's marker-phrase check, `go test ./internal/anchors ./internal/conformance`, and `.bench/skills-index.sh --check` are all green.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| BA1/bench-150 | restore a removed ceremony section | budget check (ticket 09) / review | until ticket 09 lands, review cites the budget target; after, the conformance red |
| BA2/contradiction-and-convergence | delete the non-behavioral-contradiction predicate | anchors check | remove it, run the docs-currency check, expect the new owning anchor's red |
| BA5/source-before-call | delete the source-before-call predicate | anchors check | remove it, run the docs-currency check, expect the new owning anchor's red |
| BA6/acceptance-shortfall | delete the shortfall exit | anchors check | remove it, run the docs-currency check, expect the new owning anchor's red |
| BA3/agents-shell-rules | drop the hidden-sweep rule | semantic review reread | the PG22 planted counterexample (missed dot-path) goes uncaught; review cites it |
| BA4/gate-and-anchors-green | copy a marker phrase into AGENTS.md | gate marker check | run `bench gate` fragment, expect the substring red |
