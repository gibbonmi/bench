# Ignore the capture inboxes and write ideas to the primary checkout

Blocked by: none
Writes: .gitignore, internal/roadmap/roadmap.go, internal/roadmap/roadmap_test.go, .bench/BENCH.md

## What to build

Git ignores `capture/IDEAS.md` and `capture/learnings.md` in this repo. The
files stay on disk, and the drain still reads them. When the ideas file is
ignored, `bench idea` writes to the file in the primary checkout from any
checkout. When the file is tracked, the verb keeps the current refusal and
the worktree-local write, so a linked repo keeps today's behavior.

## Acceptance

- [ ] `git check-ignore` exits zero for `capture/IDEAS.md` and for `capture/learnings.md`, and neither file is in the index.
- [ ] With the ideas file ignored, `bench idea` on the primary checkout appends to `capture/IDEAS.md` and exits zero.
- [ ] With the ideas file ignored, `bench idea` in a linked worktree appends to the primary checkout's `capture/IDEAS.md`, not to the worktree copy.
- [ ] With the ideas file tracked, the verb still refuses the primary checkout and writes the worktree file.
