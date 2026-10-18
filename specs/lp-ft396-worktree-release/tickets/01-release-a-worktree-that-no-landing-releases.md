# Release a worktree that no landing releases

Blocked by: none
Writes: .bench/BENCH.md, ROADMAP.md, roadmap/FT396.md
Covers: none

## What to build

Roadmap row FT396.

A landing releases the worktree that it lands from.
A worktree that no landing releases stays active until a person asks for cleanup.
Five occurrences from 2026-10-06 to 2026-10-14 show this result.
Examples are a superseded phase worktree, an outcome worktree after a landing from another worktree, and an assessment worktree with no commit of its own.

Add one paragraph to the workflow rules of `.bench/BENCH.md`.
The paragraph states these rules:

- A phase that moves its work to a new worktree releases the superseded worktree with `bench worktree release`.
- A session that ends with no landing also releases its worktree.
- A read-only assessment that makes no commit of its own also releases its worktree.

Remove the FT396 row from `ROADMAP.md` and delete `roadmap/FT396.md`.
No committed outcome pins FT396 as a source.

## Acceptance

- [ ] `.bench/BENCH.md` states the three release rules.
- [ ] `ROADMAP.md` holds no FT396 row, and `roadmap/FT396.md` does not exist.
- [ ] The landing gate is green.
