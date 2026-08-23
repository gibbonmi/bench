# The landing verb is the one author of the spec flip and the tickets-only close

`bench worktree land --spec <slug>` is the only command that turns a spec's `Status: staged` line into `Status: implemented`, and the only command that closes a tickets-only folder. The ordinary commit verb has no spec flag. It gates the named paths and publishes them on green, and a spec folder named as a path publishes like any other path. No verb flips a spec in the working tree.

The commit verb reports its publication boundary with exit 3. Exit 1 is a refusal before publication, and exit 2 is a grammar error. Exit 3 means the commit is published and the checkout did not reconcile. The record names the published commit, the path that did not reconcile, and one restore command over every named path as the repair. A path that is not line-safe takes a placeholder in the repair, so the record stays pasteable.

A spec retirement names the board remainder it leaves. When the spec carries a `Roadmap: FT<n>` line, the retire verb names the board row `FT<n>` and, when it exists, the row's detail file. Bench removes neither; the spec-retire commit owns the board. Without a valid `Roadmap:` line, the verb names the row and the detail file generically.

## Considered options

- **Keep three flip authors.** Rejected. They race: an early commit flips the spec before the build ends, and a working-tree flip makes the landing refuse.
- **A retry verb for an unreconciled checkout.** Rejected. The named restore is the repair, and a re-run of the commit verb reports nothing to commit.
- **A Bench write that removes the board row and the detail file.** Deferred. The drain owns the board shape; a retirement names the remainder instead.
- **A hook that refuses the commit verb on the default branch.** Deferred. It needs a reviewer enforcement decision.

## Consequences

- A green commit reads green: a reconcile failure after publication is a named repair, not a red-looking error.
- A spec's `Roadmap:` line now has a consumer, so a spec author states it.
- The anchor canary fixtures still teach the retired spec flag; they are fixtures, not guidance.
