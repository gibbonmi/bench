# Assessment owner

Status: implemented

## Problem

The platform assessment has run three consecutive times (`ASSESSMENT.md` §6's recurring `low/recurring` row, and the closing operational note) with **no owner**. Each run re-derives the whole method from the prior file: verify the last drain landed against the tree, fan out read-only area sweeps on the mid tier, synthesize adversarially on the top tier with ✓ re-verification of load-bearing claims, grade findings on a high/med/low severity grammar, and produce a ranked backlog sized in agent-time plus a verification-notes section with known limits. The method is established fact but lives only as the shape of the last output — a cold session must reconstruct it from the predecessor. `ASSESSMENT.md` is both the decision doc for this row and the exemplar the drill must codify.

## Solution

A new deliberately-invoked phase command, **`/bench-assess`** (`.agents/commands/bench-assess.md`), codifies the drill. It joins `/bench-setup-repo`, `/bench-update-kit`, and `/bench-what-next` as a periodic phase the reviewer runs on cadence or on demand — not a craft skill, and not part of the build pipeline. The method lives in the command file. The change ships with every registration surface a new command must touch, so the gate's conformance layer stays green.

**Defaulted decisions (batch-drain substitute for a decision map — flagged for post-hoc veto):**

- **D1 — owner is a command phase, not a craft skill.** `/bench-assess` codifies the method in the command file; there is no separate `craft-assess` skill. Precedent: setup-repo / update-kit / what-next are also periodic phases invoked deliberately, with their method in the command file. *Veto surface: if you want mid-work autonomous assessment guidance the model fires on its own, that is a craft skill instead — say so.*
- **D2 — routing cap for this batch.** No story routes `claude-fable-5`; the command-authoring story routes `claude-opus-4-8 / high` even though `craft-line`'s leverage override normally sends doc/command authoring to the top tier. The top-tier orchestrator performs the compensating review pass before merge. *Veto surface: lift the cap and route the authoring story to `claude-fable-5` if you'd rather spend the top tier inline.*
- **D3 — the method describes the tier split; the binding stays in the profile.** The command says *mid sweeps, top synthesizes*; the concrete tier→model binding stays in `projects/<name>.md` per invariant 2, not baked into the command.

## User stories

Routing note: `cheap = claude-sonnet-5`, `mid = claude-opus-4-8`. Per D2 the command-authoring stories route `claude-opus-4-8 / high` rather than the `craft-line` leverage-override top tier; the top-tier orchestrator does the compensating review before merge. Stories 1–6 all edit the one command file (one authoring pass, shared line); 7–11 are the wiring.

1. As a reviewer, I want `/bench-assess` to state its **entry contract** — run it periodically or on my explicit ask, as a deliberately-invoked phase — so a cold session knows when the phase applies and never fires it autonomously. `Line: claude-opus-4-8 / high.` This is the guidance-prose leverage surface capped to mid per D2, so the whole compounds through every future assessment and warrants the higher effort under the batch cap.

2. As a reviewer, I want the command to require **verifying the previous assessment's backlog landed** against the tree before any new sweep, so the new file builds on a reconciled baseline instead of re-listing shipped work (mirrors `ASSESSMENT.md` §1). `Line: claude-opus-4-8 / high.` Same authoring pass; the precondition is the load-bearing guard that keeps the drill honest.

3. As a reviewer, I want the command to name the **sweep/synthesis tier split** — mid-tier read-only area sweeps fanned out, adversarial top-tier synthesis with ✓ re-verification of load-bearing claims — while pointing at `projects/<name>.md` for the actual tier binding, so the method is portable and the binding stays reviewer-owned (D3). `Line: claude-opus-4-8 / high.` Same authoring pass.

4. As a reviewer, I want the command to fix the **output contract** — a dated `ASSESSMENT.md` at repo root that replaces its predecessor (which lives in git history), the high/med/low severity grammar, a ranked backlog sized in agent-time, and a verification-notes section including known coverage limits — so every run produces the same shape without re-deriving it. `Line: claude-opus-4-8 / high.` Same authoring pass.

5. As a reviewer, I want the command to require **parked/known items be re-confirmed rather than re-filed**, so the backlog doesn't churn already-tracked rows (mirrors `ASSESSMENT.md`'s parked-item handling). `Line: claude-opus-4-8 / high.` Same authoring pass.

6. As a reviewer, I want the command's **exit handoff** to route findings by kind — operational items to `/bench-what-next`, backlog items into `IDEAS.md`/`ROADMAP.md` through the drain — recommended in this harness's invocation form, so the assessment feeds the workflow instead of dead-ending. `Line: claude-opus-4-8 / high.` Same authoring pass.

7. As a Codex user, I want a **command-adapter skill** at `.agents/skills/bench-assess/` (SKILL.md + `agents/openai.yaml`) so `$bench-assess` reaches the same phase file, wired like the existing adapters. `Line: claude-sonnet-5 / low.` Mechanical mirror of the `bench-what-next` adapter at a known seam.

8. As a Claude Code user, I want `.claude/commands/bench-assess.md` mirrored so `/bench-assess` is invocable in this harness. `Line: claude-sonnet-5 / low.` A copy of the command file at a known seam.

9. As a cold session, I want `/bench-assess` **named in `.bench/BENCH.md`**'s "How the pieces fit" Commands bullet so the phase is discoverable and the adapter clears the operating-guide-documentation conformance check. `Line: claude-opus-4-8 / medium.` Touches shared canonical prose under the one-source discipline, so it earns a tier over pure plumbing.

10. As the gate, I want the **skills index regenerated** (`.bench/skills-index.sh --write`) so the generated-index equality check stays green after the new skill dir lands. `Line: claude-sonnet-5 / low.` Mechanical regeneration.

11. As the gate, I want a few **load-bearing anchors** for `bench-assess.md` pinned in the `workflow-guidance-anchors` conformance check, so the drill's non-negotiable clauses (verify-last-drain, tier split, replaces-the-prior, exit-to-what-next) carry a real red signal. `Line: claude-opus-4-8 / medium.` Choosing which clauses are load-bearing is a judgment call, not rote.

## Implementation decisions

- **New periodic phase command** `.agents/commands/bench-assess.md`. Frontmatter carries `description:` plus `disable-model-invocation: true` — matching `bench-what-next.md` and `bench-update-kit.md`, the other deliberately-invoked maintenance phases. Structure follows their prior art: `## Entry orientation`, `## Exit handoff`, then numbered method sections.
- **No `craft-assess` skill** (D1). The method is command-file content, single-sourced there; nothing model-invoked duplicates it.
- **The command describes the tier split; the binding stays in `projects/<name>.md`** (D2/D3). The command says *mid sweeps, top synthesizes with ✓ re-verification*; it does not name models.
- **`ASSESSMENT.md` is a living root doc replaced in place**, not a `specs/*.md` retire. Each run's dated file supersedes the prior one, whose text lives in git history — the command states this replace-not-append rule. No spec-retire commit and no archive folder; this is distinct from the promote-then-delete spec lifecycle.
- **Registration surfaces, all in one diff** (D per orchestrator): the Codex adapter skill, the Claude command mirror, the BENCH.md Commands-bullet mention, the regenerated skills index, and the pinned workflow anchors. The command adapter's skill dir is **index-excluded** (the skills-index generator skips any `.agents/skills/<name>/` that has a matching `.agents/commands/<name>.md`), so the generated index stays byte-identical; regeneration confirms no drift rather than adding a line.
- **ROADMAP FT54 row removal** rides the merge, not this spec's authoring — row presence is status.

## Testing decisions

- **What a good test is here.** The command's *wiring* is gate-observable through the kit-content conformance layer and gets real red signals. The command's *content quality* — whether the prose actually captures the drill well — is **not TDD-able**: it is guidance prose, compensated by the top-tier orchestrator's review pass before merge (D2). The anchor pins give a real red signal that a few load-bearing clauses are *present*, not that the surrounding prose is good.
- **Seam.** The kit-content surface conformance layer (`internal/conformance`, run as the gate's Go-root conformance phase). Relevant families: `skills-index-command-adapters` (adapter presence, frontmatter-name match, command-file reference, `agents/openai.yaml` explicit-invocation metadata, implicit-invocation disabled, adapter documented in the operating guide, generated-index equality) and `workflow-guidance-anchors` (per-command required phrases). `docs-currency-token-diet` catches stale `/bench-*` references and always-loaded token placement. Prior art: the `bench-what-next` adapter + its anchor `require()` block in `internal/conformance/docs_workflow_helpers_test.go`.
- **Gate command:** `.bench/gate.sh` (the project gate).

### Seam diagram

    trigger: bench gate → Go-root conformance phase (go test ./internal/conformance)
        │
        ▼
    .agents/commands/bench-assess.md ──▶ [ checkSkillsIndexAndCommandAdapters ] ──▶ diag[] (empty = green)
    .agents/skills/bench-assess/*     ──▶ [ checkWorkflowAnchors             ]
    .bench/BENCH.md Commands bullet   ──▶ [ checkDocsCurrency (stale refs)   ]
    .bench/skills-index.sh --write    ──▶ [ checkSkillsIndex (equality)      ]
        ◀ tests attach here: feed the conformance suite the tree under grade;
          a missing/mis-wired surface emits a named diag → phase non-zero.
    tests/canary/skills-index-command-adapters, tests/canary/workflow-guidance-anchors
        ◀ the meta-gate: deliberately-broken fixtures assert each family still bites.

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 7 | Codex adapter exists, frontmatter `name: bench-assess`, references `.agents/commands/bench-assess.md`, ships `agents/openai.yaml` with implicit invocation disabled | `skills-index-command-adapters` family (`checkCodexCommandAdapters`) | Run conformance with the adapter dir absent / name-mismatched / openai.yaml missing → named diag, phase non-zero; canary fixture `tests/canary/skills-index-command-adapters` asserts the family bites | The check enumerates every `.agents/commands/*.md` and demands a matching, correctly-wired adapter; a missing or malformed adapter emits its specific diag |
| 9 | `/bench-assess` documented in the operating guide (BENCH.md Commands bullet) | same family (the "adapter documented in `.bench/BENCH.md` or `.bench/BENCH-reference.md`" clause) | Remove the BENCH.md mention → "not documented in the operating guide" diag | The clause requires the adapter name to appear in the guide; an undocumented phase fails closed |
| 6, 9 | No dead `/bench-*` references introduced by the new prose | `docs-currency-token-diet` family (stale command references) | A reference to a non-existent phase → stale-reference diag | The check cross-lists referenced phases against the command set on disk |
| 10 | Generated skills index stays byte-identical (adapter is command-excluded) | `checkSkillsIndex` equality | Already covered — not a new red signal: a stray index line for the excluded dir would fail equality, but the exclusion rule keeps it green | The generator skips command-named skill dirs; regeneration confirms no drift rather than adding coverage |
| 11 | Load-bearing drill clauses present in `bench-assess.md` (verify-last-drain, tier split, replaces-the-prior, exit-to-what-next) | `workflow-guidance-anchors` family (`checkWorkflowAnchors` `require()` entries) | Drop a pinned phrase from the command → missing-anchor diag; canary fixture `tests/canary/workflow-guidance-anchors` asserts the family bites | Each pinned phrase is a `require(file, needle)`; deleting the clause makes the needle absent and the check red |
| 1–6 | The command captures the drill *well* (prose quality) | command file | **Not TDD-able** — guidance prose; compensated by the top-tier orchestrator's compensating review pass before merge (D2) | No gate observes prose adequacy; the anchor rows pin presence only, and review covers the rest |
| 8 | Claude mirror `.claude/commands/bench-assess.md` present and invocable | dogfood invocation + review | **Not gate-observable** — no conformance parity check exists for `.claude/commands` copies; caught by invoking `/bench-assess` in this harness during the dogfood run and by review | The mirror is a hand-maintained copy with no automated parity check; the synthesis dogfood loop is its red signal |

### Edge inventory

Walked per behavior; the applicable classes land as rows above, the shell-CLI classes resolve as **Won't handle** because this change adds **no `bench` subcommand or `bin/bench.sh` route** — `/bench-assess` is a markdown phase invoked as a slash command / `$bench` skill, with no executable surface.

- **stale-reference (dead `/bench-*`)** — covered (row: `docs-currency-token-diet`).
- **index-regeneration drift** — covered (row: `checkSkillsIndex` equality, green by the command-exclusion rule).
- **adapter-metadata malformity** (name mismatch, missing `openai.yaml`, implicit invocation not disabled) — covered (row: `skills-index-command-adapters`).
- **missing workflow anchor** — covered (row: `workflow-guidance-anchors`).
- **Won't handle — Claude command-mirror parity:** no conformance check compares `.claude/commands` to `.agents/commands`; caught by the synthesis dogfood invocation and review, not the gate.
- **Won't handle — all shell hostile-input classes** (paths with spaces/globs, control bytes, missing-PATH tool, symlink invocation, SIGINT mid-run, re-run idempotency, cwd below root): no CLI surface is added, so none is reachable.
- **Won't handle — absent vs empty prior `ASSESSMENT.md`:** the command instructs the human/agent to treat a missing predecessor as a first run; there is no code path to guard.

## Out of scope

- **A `craft-assess` craft skill** — a *separate capability* only if mid-work autonomous assessment guidance is wanted (a model-invoked skill the agent fires on its own, distinct from a deliberately-invoked phase). D1 defaults against it. Later cost if reversed: ~3 edits (SKILL.md + skills-index regen + BENCH.md index line), 2 gate runs.
- **A `bench assess` CLI subcommand or automation of the sweeps/drain-verification** — a *separate capability*: it would add a `bin/bench.sh` route, a Go core package, and its own contract fixtures, turning a guidance phase into tooling. Later cost: ~8 edits, 4+ gate runs (new route + core + AXI contract + canary).
- **Tooling that auto-verifies the last drain landed** — the command instructs the manual reconcile (story 2); automating it against the tree is the same separate-capability class as the CLI subcommand above and folds into that estimate.
