# Record the repair effort, review round, and propose-writes slicing rules

Blocked by: none
Writes: .agents/skills/bench-craft-line/SKILL.md, .agents/skills/bench-craft-tickets/references/slicing-checks.md, internal/anchors, tests/canary/workflow-guidance-anchors
Covers: none

## What to build

The `internal/anchors` and `tests/canary/workflow-guidance-anchors` entries are the anchor and fixture closure of the two guidance files. This ticket changes a closure file only when a pinned sentence must change.

The kit guidance states three decided rules, and each rule has one source.

A post-review repair runs at the declared effort of its ticket and keeps the declared model of its ticket. The leverage override still applies, so a guidance repair runs at high effort.
A first review round runs each axis on the conditional review line. Only a confirming round or a later round can run an axis on the top binding, and only at the reviewer's direction.
Before the ticket graph goes to approval, and before each plan-expansion commit, the slicer runs the `--propose-writes` form of build preflight for each affected ticket. The slicer adds each listed closure file to the `Writes:` line of that ticket.

## Acceptance

- [ ] `craft-line` states that a post-review repair keeps the declared model and effort of its ticket, and no guidance file states the old low-effort repair rule.
- [ ] `craft-line` states that a first review round uses the conditional review line. Only a confirming or later round can use the top binding, at the reviewer's direction.
- [ ] The three review-line sentences that `internal/anchors` pins stay byte for byte.
- [ ] `slicing-checks.md` states the `bench preflight build <slug> --propose-writes --ticket <basename> --base <commit> --source-tip <commit>` run before ticket approval and before each plan-expansion commit, and the fold of each listed closure file into `Writes:`.
