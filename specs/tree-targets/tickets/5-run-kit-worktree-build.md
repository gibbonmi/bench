# 5. Run a kit worktree target on its own current build

Blocked by: 4-run-verbs-in-tree-target.md
Writes: internal/treetarget/ (new), internal/systemtest/
Covers: TT43, TT44, TT45, TT46, TT47, TT48, TT50

## What to build

Chunk: TT-C4.

A label target whose root declares Bench build inputs by `freshness.DeclaresBuildInputs` is a kit worktree target. Its child executable is `freshness.PublishedExecutable` of the target root. The executable runs only after `freshness.Verify` accepts it against that root.

A missing build refuses with `bench <name> --in: worktree build is missing`. Any other `Verify` refusal prints `bench <name> --in: worktree build does not match the tree`. Both refusals then print `next=bench worktree build <label>`, with the label through `sanitize.ShellQuote` when it needs quoting. Each refusal exits 1 before any child starts, and neither line prints the executable path.

The tests publish a marker script through `freshness.Publish`, so the seal is real. The stale case edits one listed build input after the publication. The changed case appends one byte to the published executable. Add one system row with `BENCH_KIT` set to the kit root, through the real executable.

## Acceptance

- [ ] A kit worktree target with a current build starts `<root>/dist/bench` as the child.
- [ ] No build exits 1 with `worktree build is missing` and `next=bench worktree build alpha`, and no marker exists.
- [ ] A changed build input exits 1 with `worktree build does not match the tree`, and no marker exists.
- [ ] One appended byte in the build exits 1 with the same lines, and no marker exists.
- [ ] The label `my alpha` gives `next=bench worktree build 'my alpha'`.
- [ ] The child sees no `BENCH_RUN_BINARY`, no `BENCH_KIT`, and `BENCH_WRAPPER` set to `<root>/bin/bench.sh`.
- [ ] In the system suite, a stale kit build refuses `bench status --in alpha` and starts no child.
