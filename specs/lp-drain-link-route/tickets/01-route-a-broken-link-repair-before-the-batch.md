# Route a broken link repair before the batch

Blocked by: none
Writes: .agents/commands/bench-drain.md
Covers: none

## What to build

A drain can retire a roadmap row or a spec that a file outside `ROADMAP.md`, `roadmap/`, and `capture/` links to.
The batch commit refuses such a file, because the file is a production path with no delivery binding.
The drain command tells the drain to repair that link through an implement-now light-path fix.
The existing implement-now rule then lands the fix before the batch worktree opens.

## Acceptance

- [ ] The `## 1. Reconcile first` section names a link that a row or spec retirement breaks outside `ROADMAP.md`, `roadmap/`, and `capture/`.
- [ ] That section routes the link repair through an implement-now light-path fix that lands before the batch worktree opens.
- [ ] The new text adds no more than two sentences, and no other drain rule changes.
