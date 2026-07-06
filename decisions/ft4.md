# ft4 — harness task list in /bench-implement-spec

Source: ROADMAP.md FT4 (MED-LOW), "harness task list in `/bench-implement-spec` —
per-harness adapter (Claude hook + phase line; Codex native)". Grilled
2026-07-06. The grill collapsed the roadmap's proposed shape: no hook, no
per-harness adapter mechanism — the feature is one generic phase line in the
canonical command file plus a Claude-adapter README note. The roadmap's "Claude
hook" framing is rejected on a mechanical fact (a Claude hook can observe or
block a tool call but cannot *create* a TodoWrite entry, so it cannot populate a
task list).

## #1: Mechanism — how does the native task list get populated and stay honest?

Type: Grill

### Question
The roadmap proposed "Claude hook + phase line". Does the task list need an
enforcement hook, or is it guidance prose? A hook is the kit's enforcement layer;
a task-list display is ergonomic and has no gate-observable behavior.

### Answer
**Phase-line prose only; no hook.** The command file instructs the agent to seed
its harness's native task list from `bench coverage <spec>`. No hook: a Claude
hook cannot create todos (hooks only observe/block tool calls), task-list display
is not gate-observable, and both target harnesses already have a native list.
Rejected: a PostToolUse validation hook on TodoWrite — brittle, Claude-only, and
it guards an ergonomic surface the gate has no stake in.

## #2: Granularity — what does each task-list item represent?

Type: Grill

### Question
Is a task a user story, a coverage row, or a phase step?

### Answer
**One task per acceptance-coverage-map row.** Each row becomes a task, marked
in-progress → done as its vertical slice turns red-to-green. This ties the visible
list to the spec's fixed target and to a single source — `bench coverage <spec>`
output — so the list cannot drift from the map. Rejected: per-story (hides partial
progress inside a multi-row story) and per-phase-step (generic, says nothing about
which stories are done).

## #3: Harness parity — one generic line, or a per-harness adapter?

Type: Grill

### Question
The command file `.agents/commands/bench-implement-spec.md` is a single canonical
file every harness reads. Does FT4 keep the roadmap's "per-harness adapter", or
dissolve into one generic line?

### Answer
**One generic line in the canonical command file, plus a short per-harness note.**
The command file stays harness-neutral: "in your harness's native task list (Codex
plan, Claude todos), seeded from `bench coverage <spec>`". The Claude-specific note
(names the concrete surface, TodoWrite) lands in `.claude/README.md`, which already
documents Claude-adapter specifics. Rejected: naming TodoWrite in the canonical file
(breaks harness-neutrality — the file must read for Codex and other AGENTS.md
harnesses too) and building a real adapter file/hook per harness.

Codex-note home (veto at spec): `.codex/` holds only `hooks.json` — no README — and
the Codex adapter skill (`.agents/skills/bench-implement-spec/SKILL.md`) is
deliberately content-free (points at the canonical file; adding phase content there
violates one-source-per-fact). Recommendation: the generic line's own "Codex plan"
naming carries the Codex case; do **not** create a `.codex/README.md` for a single
line. Only Claude gets the extra README pointer, because that README already exists
and documents Claude-adapter behavior. The asymmetry is the leaner call.

## #4: Where does the line live in the command file?

Type: Proposal (not grilled — veto at spec approval)

### Question
Add a new bullet, or extend an existing one?

### Answer
**Extend the existing bullet.** The "Then build" section already has the
conditional bullet: "If the spec has an acceptance coverage map, each vertical
slice names the coverage row it is turning red-to-green before editing that slice."
One-source-per-fact: the task-list instruction extends that sentence
("...names the coverage row it is turning red-to-green **in your harness's native
task list**...") rather than adding a duplicate bullet that restates the coverage-
row-naming discipline. The bullet's existing "if the spec has an acceptance coverage
map" guard already covers the no-map case (small changes with an obvious seam get no
list). Rejected: a standalone task-list bullet — it would re-derive the coverage-row
step the existing bullet owns.

## Handoff

1. **Module boundaries.** Two files, prose only, no new units:
   - `.agents/commands/bench-implement-spec.md` — the canonical phase file (every
     harness reads it; `.claude/commands/` is a symlink, linked repos get it via
     `bench link`). One existing bullet in "Then build" is extended (#4).
   - `.claude/README.md` — gains one short line noting `/bench-implement-spec`
     seeds TodoWrite from `bench coverage`. Claude-adapter-local knowledge that
     cannot live in the harness-neutral command file.
   No CLI change: `bench coverage <spec>` already emits the rows and its slug
   fallback already resolves a spec argument (shipped, see `implement-spec-lean`).
2. **Contracts.** The extended bullet keeps its existing conditional guard ("if the
   spec has an acceptance coverage map") and the load-bearing anchor phrase
   "turning red-to-green" (a gate docs anchor — must survive verbatim). The new text
   names its source by the token the stale-reference sweep recognizes
   (`bench coverage`). The `.claude/README.md` addition is documentation only, no
   contract. No new observable behavior anywhere — the feature is instruction text.
3. **Deep vs thin.** No units, no seam of their own. `bench coverage` remains the
   single source of the row set; the phase line consumes it, it does not re-derive
   rows. Each fact the line adds (seed the list, one task per row, from
   `bench coverage`) has exactly one home.
4. **Black-box assertables.** Nothing new is gate-observable — task-list display is
   an agent behavior, not a command output. The observable surface is entirely
   existing: the stale-reference sweep (the `bench coverage` token in the new prose
   must resolve), the docs-contract anchors on the command file (the "turning
   red-to-green" anchor must still match after the edit), and structure budgets. No
   Go test, no new anchor needed (see #5).
5. **Gate attachment.** `bench gate` as-is: docs anchors + stale-reference sweep +
   structure budgets over the two edited files; no new check. Whether the agent
   actually seeds the list is **not gate-observable** — the spec records this row as
   not-TDD-able (cold-session legibility + reviewer judgment), same posture as the
   deduped-prose legibility row in `implement-spec-lean`. Flag for the spec: decide
   whether to add a docs anchor pinning the new task-list clause so a later edit
   can't silently drop it. Recommendation: **no anchor** — the clause is ergonomic,
   not a correctness contract, and the leverage override already routes the edit top
   + high; anchoring an ergonomic line over-fits the gate. Veto at spec.
6. **Hostile-input owners.** From the profile's shell-CLI checklist, only the prose-
   relevant classes apply (no code path changes):
   - spec has no coverage map (absent vs present) → the bullet's existing "if the
     spec has an acceptance coverage map" guard owns it; the line must stay under
     that guard.
   - harness with no native task list (a plain AGENTS.md harness) → the generic
     wording must degrade gracefully ("your harness's native task list, if it has
     one"), never assume TodoWrite/plan exists. Owned by the phrasing.
   - invocation through every shipped surface → n/a, no CLI/hook/adapter change; the
     canonical file reaches all harnesses unchanged by the existing symlink/link
     mechanism.
   The remaining checklist classes (spaces/globs in paths, trailing newline, PATH
   tool missing, symlink invocation, SIGINT mid-loop, re-run idempotency, cwd depth)
   are code-path classes with no surface here — a prose edit touches none.
7. **Uncertainty flags.** None open. Two items are veto-at-spec proposals, not
   uncertainty: the Codex-note home (#3, recommend no new file) and the anchor
   question (#5, recommend no anchor). Neither needs escalation above the mid tier
   for authoring; both are reviewer sign-offs on a stated recommendation.
8. **Rejected alternatives.** The enforcement hook (#1); per-story and per-phase-
   step granularity (#2); a real per-harness adapter file and naming TodoWrite in
   the canonical command file (#3); a `.codex/README.md` created for one line (#3);
   a standalone new bullet duplicating the coverage-row step (#4); a docs anchor on
   the ergonomic clause (#5).
9. **Domain watch-outs.** The command file is a leverage artifact: every
   `/bench-implement-spec` session loads it, so the edit routes **top model, high
   effort** per the profile's skill/command/doc-authoring leverage override — not the
   default mid-tier spec-authoring line, and not cheap despite being "just one
   sentence". The one-source-per-fact standard grades the diff: the new clause must
   not restate the coverage-row-naming discipline the existing bullet already owns
   (#4), and must name `bench coverage` as the row source rather than implying the
   agent hand-lists rows. The harness-neutral constraint on the canonical file is
   load-bearing: it is read by Codex and other AGENTS.md harnesses, so no
   Claude-specific tool name may leak into it — that knowledge lives in
   `.claude/README.md` only. Linked repos see the change on their next `bench link`.

Dependency order: n/a — single spec.
