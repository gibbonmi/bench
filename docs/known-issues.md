# Known issues

This file is the one owner of the known issues in this repository and their workarounds.
A coordinator reads it before it writes a delegate charge, and a delegate reads it during its work.
Each row names a retire condition: an obligation in the body of the owning roadmap row, or "no owner".
When the retire condition holds, the change that makes it hold removes the row.

| Issue | Workaround | Retire condition |
|---|---|---|
| After a ticket commit, the `dist/bench` binary of the worktree is stale, and the build preflight shows `binary-seal` red. | Run `bench worktree build <target>`. | FT406: a Go commit in a build worktree leaves no stale `dist/bench` for the next preflight. |
| A `bench probe` mutant that leaves an import unused or a variable unused does not compile. The verdict is then `invalid`, and the probe proves nothing. | Keep each import and each variable in use in the mutated text. | no owner |
| `bench probe` links the package test binary twice for each call, once for the baseline run and once for the mutant. In a large package such as `internal/adopt`, each call costs about one minute. A charge that asks for one probe per independent expectation multiplies that minute by the expectation count. The author of FT217 ticket 01 ran 54 serial probes after it finished its code in 13 minutes. | Apply the witness cap and the batched `bench test` call that `/bench-implement-spec` states. Run the probes after the floor checks, so that an invalid mutant does not block the commit. | FT168: one probe takes several mutations in one private build. |
| An adoption writer now removes its lock file when it releases its lease. An older Bench binary does not check that the lock file path still names the file that it locked. When an older binary and a newer binary write to the same destination at the same time, both writers can hold that destination. | Use one Bench version at a time for adoption writes (Link, Setup, Upgrade, Unlink, and compatibility repair). Finish or stop each session on the older binary before you use the newer binary. | no owner: every session uses a Bench binary that is built from the adoption lock fix. |
| After `main` moves during a chunked build, a bare `bench preflight build` grades `base-current` red between chunks. The row names a `main` merge that the review chain refuses. | Pass an explicit `--base` and `--source-tip` to `bench preflight build`. | FT342: a bare `bench preflight build` grades `base-current` green between chunks after `main` moves. |
| `bench record verification` reads each excerpt from a file. A response that Bench does not spill leaves no file. | Write the excerpt to a file first, then pass that file to `--excerpt`. | FT405: the record verb reads an excerpt from the last response of the worktree. |
