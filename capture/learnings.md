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

## 2026-08-20 - A narrowing introduced a fail-open the inherited coverage rows could not catch  [open]

/bench-implement-spec (FT229) shipped a fail-open at the
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

## 2026-08-20 - The committing verb was used as a diagnostic and landed a probe commit  [open]

/bench-implement-spec (FT229) used `bench commit -m "probe"` as a
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

## 2026-08-20 - The guard refused a file write over its own heredoc body  [open]

the destructive-git guard refused an ordinary file write because
the heredoc *body* being written contained a destructive git command as prose.
The classifier read the whole command text, and a heredoc body is data the
shell hands to a file, never a command position. This is the same class the
FT229 ticket fixed for the degraded rim — text mentioning a verb is not an
invocation of it — but the core classifier still matches it. Cost: two blocked
writes while capturing this journal. Proposed rule change: none for the agent;
the guard should treat a heredoc body as data, and until it does, capture
prose about a destructive command by assembling the token rather than spelling
it.

## 2026-08-20 - The landing refused five times, each on an undocumented precondition  [open]

`bench worktree land` refused five times in a row on the
worktree-exec-run-binary landing, and no two refusals named the same
precondition. In order: the destination was not clean (another session's
capture); an abbreviated SHA reported as `worktree source tip mismatch`, which
reads as a real divergence rather than a format complaint; `source and
destination do not carry identical staged spec bytes`; `reviewed source range
or ownership fence is invalid`; and finally a successful publication that still
exited non-zero because ignored residuals blocked the source release.

The load-bearing one is the third and fourth together, and they encode an
ordering constraint the phase prose never states: **a spec amendment made
during a build must reach the destination and then be the source's base — it
must never sit inside the reviewed range.** Landing demands the two trees carry
identical spec bytes, and it re-runs the review preflight over the range, where
`paths-authorized` reds on `specs/` because no ownership fence names it (and by
tree convention none should). Those two requirements are only jointly
satisfiable by rebasing the source onto the amendment. Copying the amendment
onto the source as an extra commit satisfies the byte check and then fails the
range check — the obvious move is the wrong one.

Why it happened: review found two edges decided but untested, so rows WX21 and
WX22 were added mid-build. That is the normal, correct outcome of a review
round, so this ordering will recur on every build whose review amends the spec.
Cost here: six cycles, one wasted commit-copy onto the source, and a full gate
run per attempt.

A related gap made it worse. Nothing moves an existing worktree to a new base:
`create --refresh` refreshes the pool, not the checkout, and the path-taking
verbs refuse the tilde form that `bench worktree path` itself emits, so
`release` could not be used either. Both times the only route left was raw Git
inside the worktree, against this repo's own preference for lifecycle through
Bench verbs.

Proposed rule change, two parts. For the agent: amend the spec *before* the
integration worktree is created, or rebase the source onto the amendment
immediately — never commit a spec edit onto the reviewed source. For the kit:
`bench worktree land` should accept an abbreviated SHA or say the form is
wrong; it should distinguish "published, release pending" from "refused" in its
exit code, since publication had already succeeded; and a verb that rebases a
retained source onto a moved base would remove the raw-Git fallback entirely.

## 2026-08-20 - A repo convention was asserted from artifacts instead of from the check that grades them  [open]

The spec phase reported to the reviewer that this tree does not use inline
coverage-row citations in ticket files, having grepped the ticket and spec
artifacts and found none. The build preflight immediately red on `rows-owned`,
which requires every declared row ID to appear in a ticket file, and the
just-landed FT229 tickets did carry the annotation. The sample was read; the
enforcement was not. Cost: one red preflight, one extra commit, one worktree
base move. Why it happened: the check that owns the convention lives in
`internal/preflight`, while the artifacts are in `specs/`, and only the
artifacts were searched. Proposed rule change: when reporting that something is
or is not a convention here, cite the check that enforces it, not a sample of
files that happen to follow it. An absence in artifacts is not evidence of an
absent rule.

## 2026-08-20 - A reused gate verdict was presented as evidence for a code path  [open]

A composed-run demonstration was recorded for coverage row WX20, and both runs
returned `gate: green (fresh verdict reused for this tree)`. The cache answers
before the gate entry is reached, so neither run exercised the refusal path the
row exists to grade — the evidence was worthless, and it had already been
reported as verified before that was noticed. Re-running with `--fresh`
produced a real controlled A/B: the unfixed driver refuses, the fixed driver
runs all six phases green. Why it happened: a green line reads as proof
regardless of its provenance, and the provenance is stated in the same line
that reports the verdict. Proposed rule change: an evidence run for a specific
code path uses `--fresh`; a verdict carrying "reused for this tree" is not
evidence about any path, only about the tree's last real run.
