# Prioritization assessment — 2026-07-04

All 13 open learnings entries and all 22 roadmap items, each verified against the
current tree (Go port state, canary fixture tree, doctor/link/postinstall, toon-go
adoption) and prioritized: fixes first, then features. The no-grill cleanup prunes
or rewrites the stale roadmap lines; this file keeps the rationale and the current
split between light cleanup and deeper grills. This file replaces the 2026-07-02
assessment, which resolves that assessment's own "ASSESSMENT.md is a stale history
artifact" finding.

## Outdated — prune or rewrite

| Item (date) | Why it no longer holds | Disposition |
|---|---|---|
| Roadmap: shell one-source-per-fact dups (07-02) | The Go port deleted every cited site — the query script, the tree-hash mirror in the stop hook, the shell link script are gone | Remove |
| Roadmap: canary coverage gaps (07-02) | The canary machinery landed: 56 fixtures under `tests/canary/` cover the named gaps (codex-hooks, roadmap seam, JSON, frontmatter, shared-rule); the cited gate.sh line refs are dead since checks moved to Go | Remove; keep only the residual "one canary per check" meta-check idea if still wanted |
| Roadmap: BENCH.md token diet (07-02) | Shipped — reference sections live on demand in `.bench/BENCH-reference.md` | Remove; residual: craft-skill frontmatter trim, only if still felt |
| Roadmap: ASSESSMENT.md stale artifact (07-02) | Resolved by this replacement | Remove |
| Roadmap: hermetic doctor canary (07-03) | Doctor contract suite already asserts the property family (manager-dir/nvm avoidance, stale-target, foreign-shim); the narrow absolute-path string assertion is mostly subsumed | Remove, or keep as LOW polish |
| Roadmap: parser candidates trio (07-02) | `bench doctor` shipped (doctor --fix + contract tests) — strike it from the parked list; refs/detect remain parked | Rewrite line |
| Roadmap: bench refs (07-04) | Duplicate of the refs candidate already parked in the 07-02 trio | Merge the two lines |
| Learning: by-path routing broke linked-repo CLI (07-03) | Defect fixed — `bench link` now distributes the binary to the linked kit dir (first resolution candidate) | Retire the defect half; keep the rule proposal (edge-inventory: which surfaces invoke a ported command, through which CLI) |
| Learning: hand-rolled TOON dialect (07-03) | Concrete issue resolved — toon-go is now the dependency | Retire the defect half; fold the rule proposal into the runnable-probe rule below |

## Triage split

| Bucket | Items | Why |
|---|---|---|
| No grill; light cleanup | stale roadmap pruning/rewrite, duplicate `bench refs` merge, `bench symbols` drop, README one-source pass, generated-index pointer, README drift guard | The assessment already states the current fact and the one-source-per-fact rule decides the edit. No product semantics change. |
| No grill; leave parked | dashboard, Sonnet 5 revisit, `bench refs`/`bench detect`/`bench doc`/`bench specs --retired` pending evidence | Already declared low, scheduled, or evidence-gated. No immediate edit improves them. |
| Needs deeper grill/spec | line-governance cluster + token-cap-to-iteration-cap, review-findings persistence, stale-gate status split, tests/canary bloat discipline, auto-repair download fallback, adversarial gate pinning, `bench spec implemented` + `bench commit`, harness task list, `bench outline` | These pick new semantics, new command surfaces, or new enforcement posture. They need reviewer choices before implementation. |

## Fixes, in priority order

**F1 (closed by no-grill cleanup) — README drift.** README now points to
`.bench/BENCH.md` for shared operating rules and to `.bench/BENCH-reference.md` for
the generated skills index instead of carrying hand-maintained copies. The
conformance layer now guards README against reintroducing the shared-rule sections.

**F2 (HIGH) — continue `/bench-integrate-learnings`.** The line-governance cluster
and token-cap-to-iteration-cap decision are promoted. The remaining journal backlog
falls into three clusters:

1. **Review-findings persistence (MED-HIGH):** three-axis findings still live
   only in chat (no `reviews/` surface, no spec heading exists) — a mid-fix
   disconnect loses verified findings.
2. **Status stale split (MED):** verified still unimplemented — the status
   signal reports raw `stale` with no benign-drift (capture-scratch) vs
   real-drift (committed code moved) classification.
3. **One-liner batch (LOW, cheap):** skill/command edits take effect next
   session, so dogfood shifts on trigger changes need a fresh session; review
   findings are verified in the live session, not a worktree; Codex-style
   no-subagent harnesses run the three review axes inline and flag it;
   byte/wire-compat claims about external libraries require a runnable probe plus
   an official-implementation conformance check in the spec edge inventory.

**F3 (closed) — iteration-cap line declaration.** The line now
uses iteration caps as the stop condition; token estimates are optional sizing
notes only.

**F4 (MED-LOW) — tests/ bloat review.** Audit growth and pasted fixture
harnesses; the canary tree grew fast (56 fixtures) since this was parked, which
strengthens the case for a lean-suite discipline.

*Live repo signals, outside the two lists but worth noting: `bench status`
reports 7 structure issues and 16 merged specs awaiting retirement
(promote-then-delete). The retirement pass is real pending work and may generate
the evidence the parked `bench specs --retired` idea is waiting for.*

## Features, in priority order

**FT1 (HIGH) — auto-repair download fallback for the missing platform binary.**
Fresh evidence from today's session: a missing core binary degraded the git guard
hook and hard-blocked every git command until a manual source build. For npm
consumers there is no source tree to build — the registry-fetch + hash-verify
fallback in the launcher is the fix-shaped feature. Well-scoped (~15 edits).

**FT2 (MED) — adversarial gate pinning.** Hash-verify the gate outside the
writable tree in pre-push. Distinct threat model from the lazy-agent tripwire;
small (~6 edits) and closes the "determined agent weakens the gate" hole.

**FT3 (MED-LOW) — `bench spec implemented` + `bench commit`.** Pair them: the
roadmap already notes commit could fold in the spec status flip. Replaces footgun
prose in the implement phase with two small wrappers over existing logic.

**FT4 (MED-LOW) — harness task list in `/bench-implement-spec`.** Per-harness
adapter (Claude hook + phase line; Codex native). Promote via the F2 pass since
it's already learnings-funnel material.

**FT5 (LOW) — `bench outline`.** Marginal for this repo, real as a kit
affordance for large/polyglot linked repos. Needs its grill (languages, on-demand
vs committed, prose anchors).

**FT6 (LOW, parked pending evidence — leave parked):** `bench refs`, `bench
detect`, `bench doc`, `bench specs --retired` (watch the 16-spec retirement pass
for its evidence). `bench symbols` is not carried; restore only if agents
demonstrably burn turns on symbol search.

**FT7 (LOW) — dashboard.** Low priority by declaration; unchanged.

**FT8 (scheduled, not actionable) — Sonnet 5 mid-tier revisit.** Time-boxed to
2026-09-01 or the next frontier shift; keep as is.

## Recommended sequence

1. Finish the remaining `/bench-integrate-learnings` clusters, starting with
   review-findings persistence.
2. FT1 (binary auto-repair) — today's outage is the justification.
3. FT2 (adversarial gate pinning) — only after the current gate diagnosis is settled.
4. FT3/FT4 by appetite once the findings-persistence and harness-fallback decisions
   are closed.
