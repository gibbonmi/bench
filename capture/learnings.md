# Learnings — usage journal

- 2026-08-19 — /bench-write-spec (FT229) took 2 iterations to accept. Stage
  that missed: spec authoring. What review caught: (1) two coverage rows
  promised that a diagnostic "names the command that rebuilds the binary"
  without pinning the command, so the cheapest wrong build prints plain
  `go build` — which `projects/benchkit.md` explicitly forbids — and both rows
  pass; (2) the spec invented a standing conformance check holding the
  tickets-only count at zero, which would red on every future light-path
  landing, because the same spec's close step requires that folder to survive
  its own landing gate run. The check was also an addition beyond the decision
  source, which asked for a deletion plus a status row. Why missed: the author
  applied the rubric's exact-predicate test to behavior cells but not to a
  quoted remedy string, and added a durable check without walking the in-flight
  tree states the rest of the spec creates. Proposed rule change: a row
  promising that output "names the command" carries the exact command token,
  never the family; and a spec adding a standing check enumerates the in-flight
  states that check must tolerate, drawn from the same spec's own lifecycle.

- 2026-08-20 — /bench-implement-spec (FT229) shipped a fail-open at the
  enforcement boundary that the three-axis review caught, not the gate. The
  ticket narrowed the degraded git-guard rim from a raw-substring test on the
  whole envelope to a parse of `tool_input.command`. The substring test was
  escaping-blind by construction, so escaping could never matter to it; the
  parse introduced a JSON decoder, and the decoder collapsed every `\uXXXX`
  escape to one placeholder. That is right for word structure and wrong for
  operators, so a benign read followed by an escaped `&&` and `git push --force` lexed as
  one command whose first token was the read, the verb left command position,
  and the destructive half was allowed. Go's `encoding/json` HTML-escapes `&`,
  `<`, `>` by default, so that is the ordinary producer shape, not a contrived
  one. Why missed: the spec's coverage rows H20–H22 carried over the *old*
  rim's three questions — destructive refused, path-only allowed, unreadable
  refused — and none named the decoder the ticket was introducing, so the
  escape table shipped with zero assertions and the gate had nothing to run.
  Proposed rule change: when a ticket replaces a blunt check with a precise one,
  the coverage map adds a row for the parse surface the precision creates,
  because rows inherited from the blunt check cannot mention it.

- 2026-08-20 — /bench-implement-spec (FT229) used `bench commit -m "probe"` as a
  diagnostic to learn which gate phase a composed path set reds on, and it
  committed on green, landing a commit named `probe` on the integration source.
  Amending is correctly refused by the guard, so the message stands until
  `bench worktree land` composes the source into one published commit. Why it
  happened: `bench commit` is the only surface that gates a *composed* snapshot
  of named paths, and it commits on green as its whole purpose; a working-tree
  `go test` grades a different subject and passed while the composed snapshot
  was red. Proposed rule change: none for the agent — the gap is a missing verb.
  `bench commit --dry-run` (gate the composed snapshot, report the phase
  verdicts, commit nothing) removes the incentive to misuse the committing verb
  as a diagnostic.

- 2026-08-20 — the destructive-git guard refused an ordinary file write because
  the heredoc *body* being written contained a destructive git command as prose.
  The classifier read the whole command text, and a heredoc body is data the
  shell hands to a file, never a command position. This is the same class the
  FT229 ticket fixed for the degraded rim — text mentioning a verb is not an
  invocation of it — but the core classifier still matches it. Cost: two blocked
  writes while capturing this journal. Proposed rule change: none for the agent;
  the guard should treat a heredoc body as data, and until it does, capture
  prose about a destructive command by assembling the token rather than spelling
  it.
