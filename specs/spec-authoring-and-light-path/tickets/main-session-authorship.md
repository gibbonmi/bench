# State main-session spec authorship

Blocked by: move-slicing-into-write-spec.md
Writes: .agents/commands/bench-write-spec.md, .agents/commands/bench-shape-idea.md, projects/benchkit.md, docs/field-guide.html, internal/anchors/registry_data.go, tests/canary/workflow-guidance-anchors/

## What to build

bench-write-spec.md's "Who runs this phase" says the session holding the
decision source authors the spec and tickets at whatever tier it runs; the
retired wording is that section's mid-tier-authors sentence ("carries their
decision context through reviewed ticket slicing"), landing as a
Forbid+Require pair — do not touch the Exit handoff's retained
fresh-mid-tier *build* recommendation, whose Require row belongs to the
exact-size `workflow integration source: ` family. The stale `new session on
the mid tier` registry row is retired here and a new authorship fixture
(none exists today) bites both halves. benchkit.md's spec-authoring Lines
row is rewritten to main-session authorship with the harness×tier binding
table untouched, and this ticket owns the `benchkit-spec-ownership` fixture
rewrite: its MUTATE old-string sits inside that spec-authoring bullet, and
the rewrite deliberately supersedes the pinned derivation clause so the
mutation string must be re-derived from the new text — the fixture red is by
construction, not contingent on incidental reflow. bench-shape-idea.md's exit
recommendation drops the fresh-mid-tier routing and recommends
`/bench-write-spec` from the session holding the ready decision source, as a
Forbid+Require pair with its fixture updated or created. field-guide passages
naming a fresh mid-tier authoring session are corrected (field-guide is
shared with move-slicing-into-write-spec.md — its passages interleave with
that ticket's slicing-at-implement targets; the blocker edge already forces
the order). Blocked by
move-slicing-into-write-spec.md because the successor clause says the session
"authors the spec and tickets" — a write-spec slicing step that does not
exist until that ticket lands. The four edited files are grouped by
reviewer-chosen bundling: no thinner cut strands a gate red, but each
intermediate commit would ship guidance whose authorship story contradicts
benchkit's Lines row or shape-idea's routing — flagged for the reviewer, who
may split. Shares bench-write-spec.md, `registry_data.go`, and the fixtures
tree with siblings — those paths land serially across the whole spec.

## Acceptance

- [ ] bench-write-spec.md carries the main-session authorship clause, the old
      wording trips a Forbid, and the new fixture bites both halves
      (covers WF5)
- [ ] benchkit.md's spec-authoring row reads main-session authorship,
      line-routing stays green, and the rewritten `benchkit-spec-ownership`
      fixture bites on its re-derived mutation string (covers WF10)
- [ ] the retired `new session on the mid tier` row no longer appears in the
      root-conformance sweep output (covers WF12)
- [ ] docs/field-guide.html no longer names a fresh mid-tier authoring
      session in its spec-authoring passages (covers WF5)
- [ ] bench-shape-idea.md's exit recommends `/bench-write-spec` from the
      session holding the ready source, the retired routing trips a Forbid,
      and the fixture bites both halves (covers WF18)
