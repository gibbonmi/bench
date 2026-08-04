# Align the guidance prose with the reclaim verb

Blocked by: add-spec-build-reclaim.md
Ownership fence: `.bench/BENCH.md`, `.agents/commands/bench-implement-spec.md`, `projects/benchkit.md`
Contracts: the spec-build operation inventory crosses `internal/spec/build.go`→`.bench/BENCH.md`'s CLI Inventory and the two guidance documents that restate it, asserted by AL1 by enumerating the parser's operation set against every advertisement rather than checking one document
Assumptions: `add-spec-build-reclaim.md` has landed, so the parser's operation set is the fact these documents advertise; the harness-run lifecycle and the maintainer-run maintenance verbs are separate claims and the prose keeps them separate rather than collapsing to a single count; the top binding applies because guidance prose compounds through every session that loads it; claims re-derived from the tree at pickup

## What to build

`reclaim` makes the parser's operation set nine, and three always-loaded
documents advertise the old count as an eight-item mutation surface. Leaving them
stale is the one-source-per-fact defect the code standard names, and bumping the
number blindly is worse: the sentence in `/bench-implement-spec` that says "these
eight operations are the complete mutation surface" is making a claim about what
*the harness* runs during a build, and `reclaim` is not that. A maintainer runs
it, once, over a terminal run.

So the edit distinguishes rather than counts. The eight lifecycle operations
remain the complete surface a build harness drives; `reclaim` joins as a
maintainer-run maintenance verb alongside `abandon`'s plan/apply discipline, in
the same register the CLI Inventory already uses for maintainer-run entries such
as `bench prep-release`.

Three documents carry it: the CLI Inventory's spec-build line in `.bench/BENCH.md`,
the lifecycle sentence in `.agents/commands/bench-implement-spec.md`, and the
seam paragraph in `projects/benchkit.md` that says reviewed builds route through
all eight operations.

Keep the edits minimal and in each document's existing voice. This ticket adds no
behavior and changes no code.

## Acceptance

- [ ] [AL1] every document that advertises the spec-build operation set names `reclaim`, enumerated against the parser's own operation set rather than spot-checked in one file.
- [ ] [AL2] each of the three documents keeps the harness-run lifecycle and the maintainer-run verb as distinct claims rather than merging them into one count.
- [ ] [AL3] no shared-rule text is duplicated back into `AGENTS.md` and no code file is touched.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| AL1 | leave the `projects/benchkit.md` seam paragraph on the old eight-operation wording | the conformance docs-currency check | revert that paragraph, run `bench gate`, expect the stale-reference failure naming the file |
| AL2 | rewrite the `/bench-implement-spec` sentence as a flat nine-operation mutation surface | the reviewer, at review | make the flat edit, re-read the sentence against `cmd/bench/specbuild.go`'s dispatch, expect the claim that the harness drives `reclaim` to be false |
| AL3 | copy the reclaim description into `AGENTS.md` as well | the shared-rule single-sourcing conformance check | add the duplicate paragraph, run `bench gate`, expect the shared-rule duplication failure |
