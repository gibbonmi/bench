# Probe only a tree with no author checks

Blocked by: none
Writes: .agents/commands/bench-implement-spec.md, ROADMAP.md, roadmap/FT436.md
Covers: none

## What to build

Roadmap row FT436.

`bench probe` writes a mutation into its subject until it restores the subject.
An author check that runs in the same tree during that window can read the mutant.
In the 2026-10-17 anchor-admission-prune build, one coordinator probe ran while another author ran checks in the same tree.

Add one paragraph to the implement-spec phase guidance.
The paragraph states that the orchestrator runs `bench probe` only in a tree where no author runs checks.
Remove the FT436 row from `ROADMAP.md` and delete `roadmap/FT436.md`.
No committed outcome pins FT436 as a source.

## Acceptance

- [ ] `.agents/commands/bench-implement-spec.md` states that the orchestrator runs `bench probe` only in a tree where no author runs checks, and it gives the reason.
- [ ] `ROADMAP.md` holds no FT436 row, and `roadmap/FT436.md` does not exist.
- [ ] The landing gate is green.
