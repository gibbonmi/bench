# Known issues

This file is the one owner of the known issues in this repository and their workarounds.
A coordinator reads it before it writes a delegate charge, and a delegate reads it during its work.
Each row names a retire condition: an obligation in the body of the owning roadmap row, or "no owner".
When the retire condition holds, the change that makes it hold removes the row.

| Issue | Workaround | Retire condition |
|---|---|---|
| After a ticket commit, the `dist/bench` binary of the worktree is stale, and the build preflight shows `binary-seal` red. | Run `bench worktree build <target>`. | FT406: a Go commit in a build worktree leaves no stale `dist/bench` for the next preflight. |
| A `bench probe` mutant that leaves an import unused does not compile. The verdict is then `invalid`, and the probe proves nothing. | Keep each import in use in the mutated text. | no owner |
