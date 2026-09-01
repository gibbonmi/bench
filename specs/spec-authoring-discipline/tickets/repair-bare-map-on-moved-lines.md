# Repair: write the two-word term on the reflowed shape-idea line

Blocked by: define-the-map-glossary-terms.md
Writes: .agents/commands/bench-shape-idea.md, reviews/spec-authoring-discipline.md
Covers: SAD42

## What to build

Row SAD42 grades every sentence the build adds or moves. The review found one bare
"map" noun on a line the build reflowed. `.agents/commands/bench-shape-idea.md`
line 48 writes "moves a ready map and its owned assets". The noun becomes "decision
map". No anchor needle and no canary fixture pins that line.

The repair also deletes `reviews/spec-authoring-discipline.md`, because this commit
closes its one repair target.

## Acceptance

- [ ] SAD42 — line 48 of the shape-idea phase file writes "ready decision map".
- [ ] Every `workflow-guidance-anchors` fixture still bites.
- [ ] `reviews/spec-authoring-discipline.md` is deleted.

## Delegate charge

You work in the Bench repo on the `spec-authoring-discipline` spec. Line: opus /
medium. Effort: medium, at most 2 iterations.

Read `.agents/commands/bench-shape-idea.md` at line 48. Replace the bare "map" noun
on that line with "decision map". Change no other line. Run `bench gate-prose` on
the file.

Run `bench worktree exec "<label>" -- go test -parallel 2 ./internal/anchors/ ./internal/conformance/`.
The exec form is the only command form. Do not use `cd`. Do not commit. Do not
edit the spec.
