# Add the known-issue table and point delegate charges at it

Blocked by: none
Writes: docs/known-issues.md (new), .agents/skills/bench-craft-delegate/SKILL.md, ROADMAP.md, roadmap/FT378.md, roadmap/FT369.md
Covers: none

## What to build

The reviewer decided on 2026-10-08 that one tracked file holds the known issues of this repository.
`docs/known-issues.md` holds one table with an issue, a workaround, and a retire condition in each row.
A retire condition cites an obligation in the body of the owning roadmap row, or it states "no owner".
Linked repositories get no copy of the file.
The charge section of `craft-delegate` names the file by its path, so a charge copies no row.
This change delivers FT378, so the FT378 index line and its detail file leave the roadmap.

## Acceptance

- [ ] `docs/known-issues.md` holds a purpose paragraph and one table with the columns issue, workaround, and retire condition.
- [ ] The table holds four rows: the stale plan digest, the probe exit code, the stale worktree binary, and the unused import in a mutant.
- [ ] Each retire condition quotes an obligation that the cited roadmap row body holds, or it states "no owner".
- [ ] The charge section of `craft-delegate` names `docs/known-issues.md` and tells the charge to copy no row.
- [ ] `ROADMAP.md` holds no FT378 index line, and `roadmap/FT378.md` is absent.
- [ ] `bench test --check roadmap-detail-integrity` passes after the removal.
