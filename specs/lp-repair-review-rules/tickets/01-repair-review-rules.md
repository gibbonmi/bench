# Repair guidance states the repair-cost rules and review guidance states the variant sweep

Blocked by: none
Writes: .bench/BENCH.md, .agents/skills/bench-craft-line/references/bounded-repair-policy.md, .agents/skills/bench-craft-delegate/references/delegation-discipline.md, .agents/skills/bench-craft-review/references/finding-discipline.md, .agents/commands/bench-review-implementation.md, CHANGELOG.md, ROADMAP.md, roadmap/FT381.md, roadmap/FT403.md, roadmap/FT447.md, docs/adr/0023-each-ticket-gets-a-fresh-author.md
Covers: none

## What to build

The reviewer decided the FT447 repair-cost rules on 2026-10-25 and 2026-10-29, and the FT403 review rules are decided kit edits.
Put each rule in its one owner, and point to an existing owner instead of a restatement.

- The bounded repair policy states the axis rule after a production repair and the severity route to optional advice.
- The repair-charge template states the probe rules of a repair charge.
- `.bench/BENCH.md` states the one repair session for the small findings of a chunk round. `/bench-review-implementation` and the delegation discipline point to that rule.
- The finding discipline states the variant sweep and the stop-finding citation, and it points to the `craft-delegate` registry rule.
- The review phase states the probe owner on a shared tree.

Retire `roadmap/FT447.md` and `roadmap/FT403.md` and their `ROADMAP.md` headings. Repoint the FT447 references in `roadmap/FT381.md` to the live owners.

## Acceptance

- [ ] After a production repair, each axis that issued a finding re-reviews in full, and each other axis gives a native reaffirmation that reads only the repair delta and stays a recorded axis result.
- [ ] An initial-round finding at severity 3 or lower that is not a correctness defect goes to optional advice and takes no repair cycle.
- [ ] A repair charge names only the probes of the plan, and the one coordinator probe carries the finding mutants.
- [ ] One fresh repair session takes every small accepted finding of a chunk round, routed to the ticket that holds most of their paths, and `bench spec widen` adds the other paths to that ticket.
- [ ] A repair charge never skips a named-check plan probe.
- [ ] After a held finding whose cause is a repeatable pattern, the axis searches for that pattern; an in-fence hit becomes its own finding, and an out-of-fence hit goes to `bench learning`.
- [ ] The Coverage member sets point to the `craft-delegate` registry rule for an extended family.
- [ ] A finding that asks for a new stop or refusal first cites the validator's current diagnostic for that input.
- [ ] Each review charge names the Coverage axis as the probing axis, and while an axis probes, the orchestrator reads code only through `git show` at the frozen tip.
- [ ] No live file names FT447 or FT403.
- [ ] The `bench commit` lane passes.
