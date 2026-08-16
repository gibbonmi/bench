# Move the spec template into craft-spec

Blocked by: none
Writes: .agents/skills/bench-craft-spec/SKILL.md, .agents/commands/bench-write-spec.md, internal/anchors/registry_data.go, projects/benchkit.md, tests/canary/workflow-guidance-anchors

## What to build

A spec author reads one file for the artifact's shape and a short one for the
phase. The spec template moves out of `bench-write-spec.md` into `craft-spec`,
which already owns the authoring discipline, and the command keeps only its entry
contract, ownership, and exit handoff. The template it carries is the reduced
one: `| row | story | behavior | seam | why it catches the failure |`.

`projects/benchkit.md`'s guidance-prose-budget table gains a
`.agents/commands/bench-write-spec.md | 60` row, so the shrink is enforced rather
than a one-time diff that re-accretes. It also gains an exact
`.agents/skills/bench-craft-spec/SKILL.md` row: craft-spec is 71 lines and the
template block is 65, so taking the template puts it past the 120-line `*/SKILL.md`
glob it currently matches. Set the row to the landed line count plus modest
headroom (~150) and state the landed count in the ticket's evidence. The budget
table is the one source the check parses, so both are data edits and nothing else.

Every anchor whose needle lives inside the moved template block — the
acceptance-coverage vocabulary, the seam diagram, `tests attach here`,
`Won't handle`, `Status: staged`, `why it catches the failure`, the header line,
and the approval paragraph — retargets to `craft-spec` with its section named,
never deleted. Anchors whose needle stays in the thinned command stay put.
Leave the review-loop and `--reviewer` needles alone; the next ticket owns them.

## Acceptance

- [ ] `(covers SR18)` `bench-write-spec.md` is at most 60 lines and the profile's budget table
      carries its row; raising the file above 60 turns the budget check red.
- [ ] `(covers SR18)` `craft-spec` carries the spec template with the reduced five-column header.
- [ ] `(covers SR18)` The budget table carries an exact `craft-spec` row above its landed
      line count, and raising the file past it turns the budget check red.
- [ ] `(covers SR25, SR26)` Every template-block needle is retargeted to `craft-spec` with its section
      named, and each retargeted anchor's canary still reproduces its own
      `EXPECT` line.
- [ ] `(covers SR25)` No needle is deleted by this ticket.
- [ ] `(covers SR25)` `bench anchors` is green.
