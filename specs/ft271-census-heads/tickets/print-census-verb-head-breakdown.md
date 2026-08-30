# Print the census verb-head breakdown at the landing

Blocked by: none
Writes: internal/census/census.go, internal/census/census_test.go, internal/worktree/land.go, internal/worktree/land_resume.go

## What to build

The landing deletes an assignment's census records at the release step. The exit
then has no per-head counts, and the retro reconstructs them from memory. Before
the release, the landing prints one evidence line with the raw-call count for
each verb head.

The package `internal/census` owns the per-head reader. Add a function that
returns the count for each verb head in one assignment's record file. Give the
reader the same posture as `Counts`: an absent directory, an absent file, and an
unreadable file each return an empty result. The census is evidence beside the
landing, never a condition on it.

The landing prints the breakdown where it reads the record count today: in
`landWith` in `internal/worktree/land.go`, and in the resume path in
`internal/worktree/land_resume.go`. Print the line to stderr beside the
`landing source{...}` evidence line, in the shape `census heads{<head>=<count>,...}`.
Sort the heads by count, largest first. Sort a tie by the head name. If the
assignment has zero records, print nothing.

Pass the head text through `sanitize.Controls` before the print. The recorder
sanitizes at the write, but the print must stay safe for a record file that a
foreign writer changed.

## Acceptance

- [ ] A landing over an assignment with recorded raw calls prints one `census heads{...}` line to stderr before the release step. The per-head counts agree with the record file.
- [ ] The heads print sorted by count, largest first, with a tie sorted by the head name.
- [ ] A landing over an assignment with no census records prints no heads line. The `landed{...,census=0}` record keeps its shape.
- [ ] The resume path prints the same line from the retained record file.
- [ ] The reader counts heads from the second tab field of each record line. An unreadable census gives an empty breakdown, not a refusal.
