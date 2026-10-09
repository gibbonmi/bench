# Known issues

This file is the one owner of the known issues in this repository and their workarounds.
A coordinator reads it before it writes a delegate charge, and a delegate reads it during its work.
Each row names a retire condition: an obligation in the body of the owning roadmap row, or "no owner".
When the retire condition holds, the change that makes it hold removes the row.

| Issue | Workaround | Retire condition |
|---|---|---|
| `bench record chunk` accepts a plan digest that differs from the record plan digest. A later checkpoint then refuses with "stale plan digest". | After each plan commit, run `bench record amendment <slug> --source <commit>` before `bench record chunk`. Do the same after each ticket commit that rewrites spec seam cells. | FT318: `bench record chunk` refuses when the record plan digest differs from the plan at the tip. |
| `bench record verification --probe-exit-code` takes the exit code of the mutated test run, not the exit code of `bench probe`. The run exit is nonzero when the probe bites, and `bench probe` does not print it. | When the probe verdict is `bit` with the cause `failed`, record the probe exit code as 1. | FT168: The probe result reports the mutated test's exit code separately from the probe command's exit. |
| After a ticket commit, the `dist/bench` binary of the worktree is stale, and the build preflight shows `binary-seal` red. | Run `bench worktree build <target>`. | no owner |
| A `bench probe` mutant that leaves an import unused does not compile. The verdict is then `invalid`, and the probe proves nothing. | Keep each import in use in the mutated text. | no owner |
