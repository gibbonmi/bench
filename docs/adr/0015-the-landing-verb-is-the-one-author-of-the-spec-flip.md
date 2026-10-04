# The landing verb is the one author of the spec flip, the tickets-only close, and the verified roadmap closure

`bench worktree land --spec <slug>` is the only command that turns a spec's `Status: staged` line into `Status: implemented`, and the only command that closes a tickets-only folder. The ordinary commit verb has no spec flag. It gates the named paths and publishes them on green, and a spec folder named as a path publishes like any other path. No verb flips a spec in the working tree.

The commit verb reports its publication boundary with exit 3. Exit 1 is a refusal before publication, and exit 2 is a grammar error. Exit 3 means the commit is published and the checkout did not reconcile. The record names the published commit, the path that did not reconcile, and one restore command over every named path as the repair. A path that is not line-safe takes a placeholder in the repair, so the record stays pasteable.

The same landing is the one author of the verified roadmap closure. The active commitment can approve a spec as the complete delivery of named obligations. A green landing of that spec then publishes the spec flip and the closure in one commit. The closure records the delivery in the commitment policy. It also removes each satisfied row, the detail file of that row, and the sequence entry of each delivered outcome.

The closure also removes each satisfied row from the board dependency tables. A table row whose dependent is a satisfied row goes. A satisfied row leaves each dependency list, and a table row whose list becomes empty goes. A closed row therefore never keeps dependent work blocked.

A partial delivery closes only its named obligations, and every other row stays. A `Roadmap:` line alone closes nothing.

An explicitly listed legacy run has no delivery binding. When its approved scope lists the landed spec, the landing closes that spec's delivery in the same way. The landing releases the legacy run only when every approved deliverable in its scope is delivered. A partly delivered scope stays open.

The commitment owner derives one exact closure from the reviewed source. The landing applies that closure, and the completion gate compares the candidate with the same closure. A list of permitted paths is not sufficient. A kept sequence entry, a missing row removal, or an extra removal makes the gate refuse. The delivery fact names the reviewed source commit and its completion record, and never the publication that carries the fact.

A failure before publication changes neither the default branch nor the roadmap. After publication, the landing releases the local claims that the published fact satisfies. When that local step fails, the resume completes it and does not publish again. The published fact keeps the obligation closed while the local step waits.

A spec retirement names the board remainder it leaves. When the spec carries a `Roadmap: FT<n>` line, the retire verb names the board row `FT<n>` and, when it exists, the row's detail file. The retire verb removes neither. Without a valid `Roadmap:` line, the verb names the row and the detail file generically.

## Considered options

- **Keep three flip authors.** Rejected. They race: an early commit flips the spec before the build ends, and a working-tree flip makes the landing refuse.
- **A retry verb for an unreconciled checkout.** Rejected. The named restore is the repair, and a re-run of the commit verb reports nothing to commit.
- **A separate closure commit or a later drain.** Rejected. A separate write can succeed alone, so a red delivery could close open work or a green delivery could leave stale rows.
- **A hook that refuses the commit verb on the default branch.** Deferred. It needs a reviewer enforcement decision.

## Consequences

- A green commit reads green: a reconcile failure after publication is a named repair, not a red-looking error.
- A spec's `Roadmap:` line now has a consumer, so a spec author states it.
- The anchor canary fixtures still teach the retired spec flag; they are fixtures, not guidance.
