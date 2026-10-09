# Known issues

This file is the one owner of the known issues in this repository and their workarounds.
A coordinator reads it before it writes a delegate charge, and a delegate reads it during its work.
Each row names a retire condition: an obligation in the body of the owning roadmap row, or "no owner".
When the retire condition holds, the change that makes it hold removes the row.

| Issue | Workaround | Retire condition |
|---|---|---|
| After a ticket commit, the `dist/bench` binary of the worktree is stale, and the build preflight shows `binary-seal` red. | Run `bench worktree build <target>`. | FT406: a Go commit in a build worktree leaves no stale `dist/bench` for the next preflight. |
| A `bench probe` mutant that leaves an import unused does not compile. The verdict is then `invalid`, and the probe proves nothing. | Keep each import in use in the mutated text. | no owner |
| `bench probe` links the package test binary twice for each call, once for the baseline run and once for the mutant. In a large package such as `internal/adopt`, each call costs about one minute. A charge that asks for one probe per independent expectation multiplies that minute by the expectation count. The author of FT217 ticket 01 ran 54 serial probes after it finished its code in 13 minutes. | Cap the witnesses in a charge at one per coverage row plus one per production branch that the ticket adds. Write fewer independent expectations, for example one table-driven test whose rows share one derivation. The one-source rule in `AGENTS.md` still owes each independent expectation a demonstrated red. Run the probes after the floor checks, so that an invalid mutant does not block the commit. | FT168: one probe takes several mutations in one private build. |
