# implement-spec-full-run

## Destination

An opt-in `--full` mode on `/bench-implement-spec` that carries a spec from build to
push-ready (implement → review → final-check, debug on demand) with per-phase status
persisted for cold resumption — plus two shared-rule changes: fix-don't-park for small
discovered defects, and point-of-use reinforcement of the assume/verify rule.

## #1: Where does each phase of a full run execute?

Type: Grill

### Question
Same-session phases are token-cheap (cached context) but review becomes self-review;
subagents buy fresh context at re-read cost. Which phases get which?

### Answer
Implement runs inline. Review runs as a fresh-context subagent — independence is a
quality requirement, not just a token tradeoff: the implementing context inherits the
assumptions that produced the bugs. The subagent reads only spec + diff (~10–20k
tokens), cheaper than dragging the implementation context through review.
Final-check runs inline — mechanical (gate + commit), loaded context helps.
Debug is invoked inline when an issue needs deep analysis; it needs the failing context.

## #2: Is orchestration the default invocation or opt-in, and do standalones survive?

Type: Grill

### Answer
Opt-in via a flag (recommended spelling: `--full` in the command's argument
convention). Plain `/bench-implement-spec` keeps its current implement-only
semantics. `/bench-review-implementation`, `/bench-final-check`, and `/bench-debug`
remain standalone commands for strict phased use and mid-run resumption.

## #3: Where does per-phase status persist for cold resumption?

Type: Grill

### Answer
`session-handoff.md`. It is already the designated cold-resume artifact with an
existing rewrite-in-full discipline and `bench handoff` support. Each phase boundary
rewrites it with phase reached, spec path, commit, and the exact next harness-native
command. No new file; no collision with the spec status line that `bench commit
--spec` owns.

## #4: What does a full run do with review findings?

Type: Grill

### Answer
Fix concrete defects (bugs, spec misses, missing coverage) and re-gate without
stopping. Contestable design/judgment findings are flagged in the exit report for
reviewer veto, mirroring the existing AFK batch rule. Full-autonomy and
stop-and-report were both rejected.

## #5: What threshold decides fix-inline vs park for defects discovered mid-work?

Type: Grill

### Answer
Decision-based, not size-based: a discovered fix lands in the active workflow — as
its own commit, keeping gate attribution clean per invariant 4 — unless it needs a
reviewer decision, a new seam, or spec-level design. Parking to `IDEAS.md` /
learnings is reserved for those cases only. The rule lives in `.bench/BENCH.md`
(shared platform prose) so linked repos inherit it.

## #6: How is "NEVER assume, always verify" made effective?

Type: Grill

### Question
The sentence already sits in `.bench/BENCH.md` Roles — the distributed platform file,
not AGENTS.md, so placement matches the requirement. Is one ambient sentence enough?

### Answer
No — add point-of-use reinforcement: a one-clause verify hook in each phase where the
failure was observed (grill: look up facts in the tree before asking; implement:
verify claims against the tree, not memory; review: cite evidence, don't recall).
Keep the Roles sentence. Add a canary anchor so the prose can't silently drift out.

## #7: Which tier runs the review subagent?

Type: Grill

### Answer
Mid (Terra) by default. When the orchestrator judges the diff's scope large enough
that mid could miss important bugs, it prompts the reviewer — as a structured
decision list (AskUserQuestion in Claude Code) with a recommendation for the current
run — offering: (1) continue with mid, (2) top tier via Claude CLI (fable),
(3) top tier via Codex CLI (gpt-5.6-sol, yolo). No silent escalation.

## Not yet specified

- Exact prose shape of the escalation trigger ("scope large enough") in the command file.
- Mechanics of invoking Codex CLI as the option-3 review runner from a Claude Code session.
- Whether `bench handoff` needs a flag or the orchestrator rewrites `session-handoff.md` directly.

## Out of scope

- Making orchestration the default invocation (rejected in #2).
- A new `/bench-build` wrapper command (rejected in #2).
- A dedicated run-state file or per-spec status block (rejected in #3).
- Auto-applying design-judgment review findings (rejected in #4).

## Handoff

1. **Module boundaries.** `bench-implement-spec.md` command file (gains the `--full`
   orchestration section) — mirrored in `.agents/commands/` and `.claude/commands/`
   per the kit's existing sync; `bench-review-implementation.md`, `bench-final-check.md`,
   `bench-debug.md` unchanged as standalones; `.bench/BENCH.md` (fix-don't-park rule
   in the Workflow section; Roles sentence stays); the phase command files gaining
   one-clause verify hooks; a canary under `tests/canary/workflow-guidance-anchors/`.
2. **Contracts.** `--full` argument convention on invocation; each phase boundary
   rewrites `session-handoff.md` in full (phase, spec path, commit, next command);
   review subagent input = spec path + diff, output = findings split into concrete
   defects vs judgment calls; escalation prompt = structured decision list with a
   per-run recommendation, three fixed options (#7).
3. **Deep vs thin.** The orchestration prose is thin coordination over existing
   phases — no new logic seam; the review subagent is the deep unit (fresh-context
   three-axis review per `craft-review`).
4. **Black-box assertables.** Canary anchors assert prose presence: the Roles
   sentence, the `--full` section, the fix-don't-park rule, the verify hooks.
   `session-handoff.md` content after a phase boundary is observable file state.
5. **Gate attachment.** The workflow-guidance-anchors canary family observes the
   prose. The orchestration *behavior* (does a session actually spawn the review
   subagent, actually rewrite the handoff) is prose-driven and gate-invisible —
   verify by dogfood on the next real spec build.
6. **Hostile-input owners.** Interrupted run mid-phase → `session-handoff.md` +
   tree-wins rule (`bench status` commit count) own recovery; review subagent
   done-claim → invariant 1 owns verification against gate and `git status`;
   `--full` with no/unknown spec argument → command entry contract refuses.
7. **Uncertainty flags.** Codex-CLI-as-reviewer invocation mechanics (option 3 in
   #7) are unverified from a Claude Code session — spec-writer should treat as a
   research item or scope it to a stub that reports unavailability. The escalation
   trigger is judgment prose, not a measurable threshold — accepted.
8. **Rejected alternatives.** Default-on orchestration; new wrapper command;
   all-inline (self-review); all-phases-in-subagents; new run-state file; size-based
   fix threshold; Roles-sentence-only and strengthened-Roles-prose for #6;
   session-inherited and always-top review tiers.
9. **Domain watch-outs.** Command prose ships to both harness dirs; the shared-rule
   drift canaries fail if `.bench/BENCH.md` content is restated in `AGENTS.md` —
   the fix-don't-park rule must land in BENCH.md only.

Dependency order: slice A — `.bench/BENCH.md` fix-don't-park rule + verify hooks +
canary anchors (small prose changes, likely lighter-path); slice B — the `--full`
orchestration mode with status persistence and the tier-escalation prompt (the
spec-worthy slice, builds on A's rules being in place). Slicing stays the
reviewer's call.
