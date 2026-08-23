# Route every phase through a worktree

Blocked by: make-spec-optional-on-the-landing.md, union-merge-the-phase-owned-journals.md, name-the-source-repair-in-the-conflict-refusal.md
Writes: .bench/BENCH.md, .bench/BENCH-reference.md, .agents/commands/bench-debug.md, AGENTS.md, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go
Line: fable / high — kit guidance prose; the leverage override applies.

## What to build

The workflow guidance says every phase runs in a bench worktree and lands
through `bench worktree land`, spec-less when the phase has no spec. `main`
receives writes only through landings. The right-size table's light-path row
now says the light path runs in a worktree too and lands through the verb with
its tickets-only `--spec`. If the reviewer vetoes the tickets-only close, the
row instead states the light path's exemption. The phase-close handoff rule
says the handoff is written in the worktree and lands with it, so the handoff
rule and the merge rule agree. The debug skill's isolation rule points at the
spec-less landing.

The landing paragraph in `.bench/BENCH-reference.md` names the optional spec,
the rule table verbs, and the repair the conflict refusal prints. The guidance
keeps merge composition as the primitive and rejects rebase. It states that the
worktree rule is guidance, not a hook.

The `.bench/BENCH.md` budget is 180 lines. Its Workflow needles and the
light-path needle are registered in `internal/anchors/registry_data.go`, and
the registry's bite tests pin them. This ticket moves each needle it edits and
keeps the bite red-capable. `.agents/commands/bench-debug.md` is 168 lines
against a 170-line budget, so the isolation-rule edit replaces its current
sentences rather than adding lines.

## Acceptance

- [ ] `.bench/BENCH.md`'s workflow names the worktree rule, the spec-less landing, and the no-rebase decision (review-owned, stories 18 and 21).
- [ ] `.bench/BENCH.md`'s right-size table's light-path row names the worktree and the verb, or names the exemption if the close is vetoed (review-owned, story 18).
- [ ] `AGENTS.md`'s phase-close handoff rule says the handoff is written in the worktree and lands with it (review-owned, story 19).
- [ ] `.agents/commands/bench-debug.md`'s isolation rule names the spec-less landing, and the file stays within its budget (review-owned, story 20).
- [ ] `.bench/BENCH-reference.md`'s landing paragraph names the optional spec, the three rule verbs, and the conflict repair (review-owned, story 9).
- [ ] The guidance states that no hook refuses a commit on the default branch (review-owned, story 22).
- [ ] Every anchor needle the edits move is re-registered, and its bite test still reds on the needle's removal (edge: the anchors registry).
