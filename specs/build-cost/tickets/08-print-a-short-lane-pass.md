# Print a short lane pass

Blocked by: none
Writes: internal/gate/lane.go, internal/gate/lane_output.go, internal/gate/lane_report_test.go (new), internal/gate/lane_run_test.go, internal/gate/authorization/lane_test.go, internal/commit/lane_test.go, internal/worktree/merge_test.go, internal/worktree/land_journey_test.go, internal/worktree/merge_target_lane_test.go, .bench/BENCH-reference.md
Covers: BC49, BC50, BC51, BC52

## What to build

The lane writes each check's prefixed stream to `bench-lane.log` in the git directory of its worktree, and it truncates that log at the start of each run. The schedule call in `runLane` reads `schedule(runCtx, checkout, checks, laneWriters(log, log, diagnostics))`, because the probe-08 command swaps that exact text.
On a lane pass, the response holds no check line. On a lane red, the lane copies the failing check's own lines from the log to the response before the `lane{...}` line.
The `lane{...}` line, the background line, and the `committed{...}` line keep their shapes. `bench worktree merge` gets the same short pass, because it runs the same lane.

The lane paragraph of the command notes in `.bench/BENCH-reference.md` names the lane log and its truncation.

## Acceptance

- [ ] A lane pass over the prose and docs-currency-workflow checks prints no line that starts with `[prose]` or `[docs-currency-workflow]` (BC49).
- [ ] After that pass, `bench-lane.log` holds those lines (BC50).
- [ ] A lane red prints the failing check's own lines (BC51).
- [ ] A second run truncates the log (BC52).
