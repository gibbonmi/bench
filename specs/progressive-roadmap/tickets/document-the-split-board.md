# Document the split board in the glossary, commands, file maps, and CHANGELOG

Blocked by: register-the-roadmap-detail-integrity-check.md
Writes: CONTEXT.md, .agents/commands/bench-what-next.md, .agents/commands/bench-write-spec.md, .agents/commands/bench-shape-idea.md, .bench/BENCH.md, .bench/BENCH-reference.md, projects/benchkit.md, README.md, CHANGELOG.md, internal/anchors/registry_data.go, tests/canary/workflow-guidance-anchors

## What to build

`CONTEXT.md`'s **roadmap** and **roadmap index / roadmap detail** entries name the
index file and the `roadmap/FT<n>.md` detail owner; `/bench-what-next` says bodies
and ledgers are edited in the row file, ordering in the index, and retirement deletes
both — that sentence is anchored with a canary that fails for its own planted
removal; `/bench-write-spec`'s promote-then-delete step names both files while
keeping its existing anchor needle; `/bench-shape-idea`'s cold entry reads the index
and fetches detail per row; the capture paragraph, profile, reference file map, and
README tree name `roadmap/`; CHANGELOG gains the entry; `bench skills-index --check`
is green. Edits to `.bench/BENCH.md` and `bench-write-spec.md` are line-neutral (both
sit at their prose budget) unless the reviewer grants a budget change in this ticket.
Coverage rows PR23, PR24, PR25, PR26.

## Acceptance

- [ ] The glossary entries and file maps name `roadmap/FT<n>.md` as the detail owner.
- [ ] The what-next detail-owner sentence is anchored and its canary reds for the planted removal.
- [ ] The write-spec anchor `removes the spec's ROADMAP.md row` still matches and the sentence also names the row file.
- [ ] `bench skills-index --check` exits zero.
- [ ] `.bench/BENCH.md` stays at 180 lines and `bench-write-spec.md` at 73, or the profile budget row changed in this ticket.
