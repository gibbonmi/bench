# Derive the retro cost section

Blocked by: 02-refuse-plain-chunk-requirements.md
Writes: internal/retros/retros.go, internal/retros/retros_test.go, internal/retros/recommendations_test.go, internal/retros/testdata/eligible.md, internal/roadmap/retro_scaffold.go, internal/roadmap/retro_cost.go (new), internal/roadmap/retro_cost_test.go (new), internal/roadmap/retro_test.go, internal/roadmap/retro_scaffold_test.go, internal/spec/landing.go (new), internal/spec/landing_test.go (new), internal/reviewrecord/plan.go
Covers: BC1, BC2, BC72, BC3, BC4, BC5, BC6, BC7, BC8, BC9, BC10, BC11, BC12, BC13, BC15, BC16, BC17, BC18

## What to build

`bench retro <slug> --scaffold` prints a `## Cost` section directly after `## Outcome`. The constant `CostHeading = "## Cost"` sits in `internal/retros/retros.go`, and the canonical heading list names it on its own line as `CostHeading,`, because the probe-04 command omits that exact text. That text occurs once in the file.
The retro parser then refuses a body without the heading. The marker grade of pending retros reads no heading, so a pending retro without `## Cost` stays valid.

One function in `internal/spec` answers the landing commit of a slug. It is the newest first-parent commit of the default branch that changes `specs/<slug>/spec.md` to `Status: implemented`: its tree holds that status, and the tree of its first parent does not. It reuses the status parser of `internal/spec`.

The cost renderer in `internal/roadmap/retro_cost.go` prints these lines from the landing L:

- the landing commit and the range `L^1..L^2`
- the wall minutes, from the committer time of the oldest non-merge commit in the range to the committer time of L
- the commit count of the range, split into bookkeeping commits and change commits by `git diff-tree -z` paths
- the ticket count, the longest `Blocked by:` chain, the chunk count, and the longest chunk chain, read from the tree of L
- the author and repair dispatch counts of a version 2 plan, or `unknown` for a version 1 plan
- the wall minutes per ticket and the bookkeeping commits per ticket, each with its budget from one table

A bookkeeping commit changes only paths below `specs/<slug>/` or the path `reviews/<slug>.md`. A per-ticket value whose printed one-decimal form is above its budget ends its line with `over budget`. The budget table holds 60 minutes and 7.1 commits per ticket.
The renderer reads the plan of L through a decode-only reader that `internal/reviewrecord` exports beside `checkPlan`, so a plan that the current grade refuses still gives its facts. A slug without a landing prints `unknown` in each derived line.

Each valid retro body in the tests and the testdata gains the `## Cost` heading.

## Acceptance

- [ ] The scaffold prints `## Cost` after `## Outcome`, and the draft passes the parser (BC1, BC16).
- [ ] The first cost line names L and its range, also when a later landing of another slug comes before the retirement (BC2, BC72).
- [ ] The wall, commits, graph, and dispatch lines match their fixture arithmetic (BC3, BC4, BC5, BC6, BC7).
- [ ] A version 1 plan, a slug without a landing, and a retired spec folder print as the rows state (BC8, BC9, BC10).
- [ ] The budget lines mark only a value above its budget (BC11, BC12, BC13).
- [ ] `bench retro <slug> --body` refuses a body without `## Cost` and writes no file (BC15).
- [ ] The marker grade passes a pending retro without `## Cost` and a retro with `over budget` lines (BC17, BC18).
