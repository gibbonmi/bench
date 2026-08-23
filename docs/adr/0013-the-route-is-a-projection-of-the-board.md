# The next action is a projection of the existing board, never a second derivation

Every session — human or agent — begins by asking what to do next. The kit answers with a single routed command derived entirely from signals it already computes. The board's severity ladder is the one ranking. The route ranks that ladder, takes the highest row carrying an invocable command, and emits exactly one next row of state, why, and command. It also emits one runners-up line naming every later row that also carries a command.

There is no router state file, no compiler, and no signal computed a second time for routing's benefit. The handoff surface consumes the same selection owner, rather than keeping its own copy. So the board, the route, and the handoff cannot disagree about what is pending or what to run.

The defect this closes is that nothing routed from observed state to the next action. The board's actions were partly prose — "fix before commit", "split", "resume interrupted work" — which no harness can invoke. Any consumer had to skip past them to find something runnable. The workflow's own middle was invisible. A repository with a staged spec and a build due reported nothing, and an unadopted repository reported that it was clean. Each harness therefore began a session by guessing.

Two rules follow, and they are the substance of the decision. First, **every board action is an invocable command or it is empty** — never prose. A row that genuinely has no command renders with an empty command cell, which is honest and machine-readable, where prose was neither. The fallback command applies only to a genuinely empty board, never to a board whose rows merely lack commands.

Second, **absence of adoption and presence of staged work are signals like any other**, ranked in the same ladder rather than special-cased ahead of it. An unadopted tree leads with setup, because its severity says so. A staged spec sits at its own rung, between the phases that produce and consume it.

Harness translation is explicit, never inferred. The route takes a named harness and renders each phase in that harness's invocation form, from one shared prefix table. It does not sniff its environment to decide what it is running under. A surface a harness lacks is a dead key. So naming the harness is the caller's job, and the adapters pass it through.

The route names commands; it does not execute them. No trusted-execution or refusal-before-execution claim is introduced, and the routed command carries exactly the authority the reviewer already granted whatever it names.

## Considered options

- **A dedicated router state file recording the next action.** Rejected. It becomes a second derivation of every signal it summarizes, and the moment the tree moves it is stale in a way nothing detects. The board is already recomputed on demand; projecting it costs nothing and cannot drift.
- **Keeping prose actions and having each consumer parse them.** Rejected: that is the status quo the decision removes, and it puts the same fragile translation in every consumer. An empty command is strictly better than prose a caller must interpret.
- **Auto-detecting the harness from the environment.** Rejected. Detection is guesswork that fails silently, by emitting an invocation form the caller cannot run. The kit already had an explicit harness argument on the neighbouring surface. Explicitness costs one argument and removes a class of dead recommendations.
- **Routing across worktrees** — reading another checkout's spec or review state to route the session sitting in this one. Rejected as out of scope. The route reads the tree it runs in, plus the shared ledger. The one scenario that motivated it, a reviewed source awaiting landing, is answered by the worktree listing instead. Landing needs a request token only the assignment record holds, and a clean review leaves no artifact to observe.
- **Leaving the inventory in the wrapper.** Rejected, in favour of a single registry. That registry classifies every command public or internal and owns its help rows, with the binary projecting them and the wrapper routing to it. The accepted cost is that help now requires the binary; the wrapper reports its install message when the binary is absent. Two hand-maintained inventories could not be kept in agreement. That is the one-source rule applied to the kit's own front page.

## Consequences

- The bare invocation of both the wrapper and the binary is the route, and they return its exit code. The two front doors that previously disagreed — one printing a long inventory, the other reporting no subcommand — now give the same answer. The inventory moved to an explicit help verb.
- The phase that reconciles the roadmap and drains capture is named for what it does, rather than for the question the route now answers. So the router-sounding name no longer points at roadmap maintenance.
- An empty roadmap file still counts as present for the fallback. Reporting its emptiness belongs to the roadmap surface, not to the route. This keeps the route out of the business of judging the content of what it points at.
- Because the route reuses the ladder rather than re-ranking, adding a signal to the board adds it to the route for free. Any new row must therefore arrive with an invocable command or a deliberate empty one.
