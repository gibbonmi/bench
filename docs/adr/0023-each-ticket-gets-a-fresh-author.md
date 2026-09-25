# Each ticket of a spec-backed build gets a fresh author

Status: accepted

A long implementation session reads its full context again on each model request, so the context length drives the cost of a build.
Every spec-backed build therefore gives each ticket a fresh author session on the declared line, with or without a full run.
The author starts from a narrow charge: the ticket, its coverage rows, its expected writes, the evidence identity, and the line.
The build records its authors in the version 2 delegate plan with an author limit of 1.
So the existing record checks name the author of each ticket, and the record schema does not change.

## Consequences

- The authors work serially in one integration worktree, in dependency order. Only an explicit delegated run adds concurrent authors and the full tier range.
- Outside a delegated run, a tier move of a fresh author asks the reviewer first.
- The orchestrator reads manifests, returns, and verdicts, not code. It refreshes the session handoff at each chunk checkpoint, and it reconciles the final acceptance and integration.
- A post-review repair goes to a fresh repair session for each affected ticket. The plan records that session as a new assignment with the user-directed trigger, and the session runs the verification of its ticket again.
- The orchestrator does not repair at final reconciliation. A finding there goes to a fresh repair session for the ticket whose expected writes hold the path. A finding on a path that no ticket holds is a material acceptance shortfall.
- Each ticket takes its own session, because the review record refuses one session for two tickets.
- The one-ticket light path stays in the main session, because a dispatch costs more than its context. The operating guide names one exception for a captured learning.

## Considered options

- One retained author for the full ticket graph: rejected, because its context grows through each ticket, repair, and reconciliation.
- A serial execution mode in the review record: rejected, because the delegate plan with an author limit of 1 records the same facts.
