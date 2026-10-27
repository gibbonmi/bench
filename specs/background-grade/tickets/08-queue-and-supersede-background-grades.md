# Queue and supersede background grades

Blocked by: 07-start-a-background-grade-after-a-lane-pass.md
Writes: internal/gate/backgroundstate (new), internal/gate/background_grade.go (new), internal/gate/background_launch.go (new), internal/gate/background_queue_test.go (new), CHANGELOG.md
Covers: BG46, BG47, BG48, BG49, BG50, BG51, BG52, BG53, BG54, BG55, BG56, BG83

## What to build

Replace the single slot of ticket 07 with the machine queue of ticket 05.

- A lane pass queues nothing when its grade key has a fresh green, a live claim, or a queued entry.
- Each worktree keeps at most one queued grade. A newer tip with a new key replaces the queued grade of its worktree.
- A newer tip with a new key stops the running grade of its worktree.
- A grade whose key is the key of the latest lane-passed tip of another worktree is not stopped and not replaced.
- The queue serves the grades in the order of their lane passes.
- When a grade ends, its process starts the oldest queued grade and then exits. No resident process exists.
- A lane pass that finds no live running grade starts the oldest queued grade, which can be its own.
- A grade whose owner died records `aborted` when the next lane pass reads its claim.
- A stopped grade writes no green, and its background record names `stopped`.

## Acceptance

- [ ] Queued grades start in the order of their lane passes.
- [ ] A lane pass whose key holds a fresh green, a live claim, or a queued entry queues nothing.
- [ ] A newer tip with a new key replaces the queued grade and stops the running grade of its worktree.
- [ ] A running grade or a queued grade whose key is the latest tip key of another worktree keeps its place.
- [ ] An ended grade starts the oldest queued grade and then exits.
- [ ] A lane pass that finds no live running grade starts the oldest queued grade.
- [ ] A grade whose owner died records `aborted` at the next lane pass.
- [ ] A stopped grade writes no green, and its record names `stopped`.
