# Report the board flow from git history

Blocked by: none
Writes: internal/roadmapflow/, cmd/bench/main.go, cmd/bench/main_test.go, cmd/bench/command_registry_test.go, internal/roadmap/roadmaptest

## What to build

`bench roadmap --flow` prints one TOON block with the board's flow. The new
package `internal/roadmapflow` runs one `git log --name-status --no-renames`
query over `roadmap/`, selects the window that spans the last three drain
commits, and derives opened, fed, retired, net, open mass, the target, and the
two boundary commits. The open mass comes from `roadmap.ParseDocument` over the
live board. The bare verb and `--context` keep their current output. Spec group
A, rows RF1, RF2, RF3, RF4, RF5, RF6, RF7, RF8, RF9, RF10.

## Acceptance

- [ ] A history whose window holds 4 added, 9 modified, and 6 deleted detail files reports `opened=4`, `fed=9`, `retired=6`, `net=-2`.
- [ ] A net delta of 0 reports `target_met=true`; a net delta of 1 reports `target_met=false`.
- [ ] The report names the third-newest drain commit and the newest commit as the window boundaries; a commit whose subject says `drain` but adds no detail file is not a boundary.
- [ ] A history with 2 drain commits reports `drains=2` and sums the events it found; a history with no flow event prints the block and exits 0.
- [ ] Outside a repository the command prints the structured error and exits 1; a detail directory that cannot be listed reports `open_mass=unknown`.
- [ ] `bench roadmap` and `bench roadmap --context` emit their current blocks unchanged.
