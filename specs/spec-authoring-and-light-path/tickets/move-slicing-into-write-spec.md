# Move ticket slicing into the write-spec phase

Blocked by: none
Writes: .agents/commands/bench-write-spec.md, .agents/commands/bench-implement-spec.md, .agents/skills/bench-craft-tickets/SKILL.md, .bench/BENCH.md, docs/field-guide.html, internal/anchors/registry_data.go, tests/canary/workflow-guidance-anchors/

## What to build

`/bench-write-spec` charges `craft-tickets` after the spec's falsification
review completes — worded against the file's current step 9, whose verdict
is advisory, so this commit references no step or verdict the file lacks;
verification-loops.md re-points the trigger to loop 1 later — writes `specs/<slug>/tickets/`, and carries
the numbered title/`Blocked by:`/outcome breakdown into the same approval
table as the spec (new slicing-in-write-spec fixture; the charge and
tickets-path registry rows move here from the implement file). BENCH.md's
workflow list step 2 gains one clause — "and slice the tickets" — staying
within the 180-line budget and avoiding the file's whole-file Forbid on the
bare substring `thorough`. In `/bench-implement-spec`, remove
only the two slicing clauses — the `Charge \`craft-tickets\` ...` breakdown
charge and its reviewer-approved-breakdown / AFK-carve-out continuation,
which share one grammatical sentence in the live file — and
keep the three anchored delegation sentences verbatim — the write-subagent
assignment, the read-only-helper disqualifier, the `craft-delegate`
incapable-harness clause — and reword the integration-worktree sentence only
as far as its now-orphaned "After approval" referent, naming the
write-spec-phase approval. That sentence has no registry row today: it gains
a first-time Require row (diagnostic outside the `workflow integration
source: ` prefix) and fixture. The three delegation anchors keep their
existing rows — no duplicates — and gain first-time fixtures, so
over-deleting any of the four sentences is observable in the ordinary gate. The retained
`implement-spec-mandatory-delegation-anchor` and
`implement-spec-inline-exception` fixtures stay green untouched. Its entry
validates that `specs/<slug>/tickets/` exists and routes a ticketless
spec-backed run back to `/bench-write-spec` (preflight already accepts
present tickets — no Go change); a Forbid pins the old charge out of the file
and a Require pins the route-back, with a new entry-validation fixture.
write-spec's Exit handoff sentence is rewritten — sign-off, then a fresh
mid-tier build session on one retained integration source, slicing no longer
routed through `/bench-implement-spec` — as a Forbid+Require pair on the
anchored post-slicing-handoff needle, its backing fixture updated or created;
the Forbid's diagnostic must not begin with `workflow integration source: `
(a bespoke helper counts that family at an exact size — single-sourced in
the test — and fatals on a second row per file; its test file is unfenced
and stays untouched).
`craft-tickets`' trigger text moves from build entry to spec authoring with
its `index:` frontmatter left alone (so `.bench/BENCH-reference.md` needs no
edit); spec-backed breakdown, template, and frontier rules unchanged, and the
cadence-pinned substrings in `craft-tickets` stay byte-identical. field-guide
passages that say slicing happens at implement entry are corrected, and the
edit near the retained frontier-ticket sentence leaves that anchored sentence
byte-identical. implement-spec stays ≤ 60 lines, craft-tickets ≤ 100. Shares
BENCH.md, `craft-tickets` (with lighten-light-path.md, both displacing lines
at its 100-line limit), `registry_data.go`, and the fixtures tree with
siblings — those paths land serially across the whole spec.

## Acceptance

- [ ] bench-write-spec.md charges `craft-tickets` after the spec's
      falsification review completes (worded against the current advisory
      verdict), names `specs/<slug>/tickets/`, folds the breakdown into the
      approval table, and its Exit handoff no longer routes slicing through
      `/bench-implement-spec` (covers WF3)
- [ ] BENCH.md workflow step 2 carries "and slice the tickets" (covers WF3)
- [ ] docs/field-guide.html no longer says tickets are sliced at implement
      entry, and the adjacent anchored frontier-ticket sentence is
      byte-identical (covers WF3)
- [ ] bench-implement-spec.md drops exactly the two slicing clauses, keeps
      the four sentences above (three verbatim; the worktree sentence's
      approval referent updated), validates tickets at entry, and routes a
      ticketless spec-backed run to `/bench-write-spec` (covers WF4)
- [ ] the three delegation anchors and the reworded worktree sentence gain
      first-time fixtures whose bite and restore halves both prove out, so
      over-deleting any of the four sentences reds the ordinary gate
      (covers WF4)
- [ ] with everything landed, BENCH.md is ≤ 180 lines, implement-spec ≤ 60,
      and craft-tickets ≤ 100 (covers WF13)
