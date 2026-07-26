---
name: craft-delegate
description: The delegation discipline — when to spawn a subagent, how to charge and scope it, when it needs an isolated worktree, and how a done-claim is verified. Use whenever spawning any delegate (an axis review, a scoped build, a fan-out search) or accepting one's result. The model/effort half of the decision lives in craft-line.
index: spawning a delegate / verifying a delegate's done-claim
---

# Delegating without losing the plot

A delegate buys two things: parallelism, and isolation of a heavy read-set from
your context. It costs one thing that matters more when misjudged: the work
happens where you can't see it. This skill keeps the cost bounded.

## Delegate or inline

The coordinator scopes, routes, and verifies work; a write-delegate authors code.
The inline allowance is exactly one source-line insertion, one source-line
deletion, or one source-line replacement in code or tests. A replacement counts
as one correction even though a unified diff renders it as one deleted line plus
one added line. This allowance spans the current reviewer request — the current
user objective — and does not reset when work is split into tasks, slices,
delegates, or verification rounds. A no-spec change admitted by the
lighter-path threshold in `.bench/BENCH.md`'s "Right-size the process"
paragraph may also remain inline — the only allowance beyond the
one-source-line rule. All other code authorship runs as a
write-delegation in an isolated worktree, including mechanical work, atomic
diffs, and fixes the coordinator has already diagnosed. Read-only coordination
stays inline. Never delegate a decision the reviewer owns (what ships, spec
content, an irreversible choice); a delegate inherits your authority ceiling,
not the reviewer's.

This policy is capability-aware: a harness that cannot spawn a write subagent
never falls back inline for work beyond the allowances above. It stops before
editing and emits one executable resume handoff to a subagent-capable harness —
the repository path, the working branch or worktree, the spec or change name,
the destination harness, and that harness's exact invocation. One route, not a
menu; the phase being resumed supplies its own harness-native invocation forms.

## The charge

A delegate has no conversation memory: everything it needs must be in the
prompt. A complete charge names the objective, the inputs by path, the seam it
works at, the return shape (raw data or a structured report — its final message
is the deliverable), and its budget. Route its model and effort with
`craft-line` — the line is declared by you, not chosen by the delegate — and
carry both explicitly: the model goes on the call itself as the bound alias
(never omitted — an Agent call without a model inherits *your* model, which is
silent escalation when you run top-tier), and the effort and iteration cap go in the
charge text, because the Agent tool has no effort parameter — effort rides in
the charge or it rides nowhere.

A context-inheriting delegate (a harness's fork type, which clones this
conversation) is the one exception to the contextless premise — and it always
runs the parent model, ignoring any model override, so in a top-tier session
every fork is a top-tier delegate. A fork is legal only where `craft-line`'s
table already routes top, and it is declared like any line. Fork because the
read-set genuinely is this conversation, never to skip writing the charge;
everything else spawns fresh on the bound alias.

Prefer compressed inputs over inherited context: when a decision map has a
Handoff, give the delegate that digest plus line-ranged excerpts it must quote,
not the orchestrator's whole read list.

Name exemplar files to mirror, and say so explicitly when no exemplar exists.
A convention stated as prose degrades as the tree grows: "follow the repo's
error idiom" hands the delegate a judgment it cannot make, because the idiom
lives in files it will not read; "mirror the error shape in
`internal/x/foo.go`" survives translation to a low-context delegate, because
a path resolves the same way at any context size.

A write-delegation from a spec carries its stories' coverage rows — behavior,
seam, red signal — in the charge, every time, and requires the delegate to show
each row red before the edit and green after. That is what makes the done-claim
verifiable by running the gate instead of re-reading the work; a charge without
its rows buys a diff you must read line-by-line to trust.

The charge also names the gate layer that owns each artifact class the
delegate touches — workflows and `.bench/` content to canary, gate output
shape to canary, skills and commands to conformance — and states the converse
to the delegate in the same breath: the named list is a floor, not a ceiling.
Omit the mapping and the delegate breaks a layer nothing in its charged
verification list can see; omit the converse and it treats the list as
permission to check nothing else.

A worktree-isolated delegate's charge opens with the stale-base check: run
`git merge --ff-only main`, verify HEAD equals main, stop and report if the
merge is denied or diverges. Worktrees get cut behind a moving main, and a
delegate that builds on a stale base re-fights landed work. The orchestrator
holds the other end of the contract: a blocked worktree is fast-forwarded by
the orchestrating session, which then resumes the same delegate. Read-only
delegates are unaffected. A fix-pass charge against a repository snapshot
names a commit-specific sentinel — a function or test introduced by the
commit under fix — and requires the delegate to verify it before editing,
then stop and report if it is absent.

```
Implement story 3 of specs/retry-backoff.md in this worktree. Open with the
stale-base check: run `git merge --ff-only main`, verify HEAD equals main,
stop and report if denied. Coverage rows: [the story's rows]. Effort: medium,
~3 iterations. Stop at diff ready with focused tests green; return the
red→green log per row.
```
Good — a write-delegation whose opener rides in the charge and whose rows make
the done-claim verifiable.

```
Review this diff on the Standards axis only. Base: run `bench diff` for the
changed files; read AGENTS.md and .bench/BENCH.md for the conventions. Charge:
.agents/skills/bench-craft-review/SKILL.md. Effort: medium, ~1 iteration — one
pass, no fix iteration. Return findings under ## Standards, each citing the rule,
under 400 words. Do not edit any file.
```
Good — self-contained: inputs by path, the charge by path, the return shape and
a write prohibition stated.

```
Review the changes we discussed against our standards, like last time.
```
Bad — "we discussed", "our", "last time": every referent lives in a context the
delegate does not have; it will guess all three.

## Scope

One delegate, one coherent unit — one axis, one story, one search question. A
delegate charged with several loosely-joined jobs returns a summary that blurs
them, and you can't verify what you can't attribute.

## Isolation

A write-delegation runs in an isolated worktree (`bench worktree`), so stray
edits can't land in reviewer-owned files — the delegate gets a checkout, not
your checkout. Concurrent delegates get *separate* worktrees, one each: two
writers in one checkout collide, and a mixed `git status` makes both
done-claims unverifiable. The harness's own `isolation: worktree` cannot cut
the second one — it derives its request ID from the harness session ID alone,
so a second concurrent request collides with the first assignment and is
refused (`worktree create request conflicts with its existing assignment`).
The coordinator cuts the worktrees instead: run
`bench worktree create --request <opaque-id> --label <work-item>` once per
delegate — a distinct request id each — and hand each delegate its returned
root path. Share a worktree — or run sequentially — only when one delegate's
work genuinely depends on another's output. A charge that shares an existing
worktree names its root and pins every file-tool path to that root; shell CWD
does not retarget file tools.

The whole-tree gate is a serialized resource: concurrent `bench commit` gates
flake load-sensitive tests that pass serially — a red that answers for machine
load rather than for any diff. A write-delegate stops at "diff ready, focused
tests green"; the coordinator runs `bench commit` per worktree, one at a time.
When `bench commit` reports nothing to commit beside a visibly modified file,
the coordinator diagnoses a CWD/tree mismatch before treating the command as
defective.

A worktree isolates the working tree, not repository-global git surfaces —
the stash stack above all. Two delegates in separate worktrees share one
stash stack, and stashing cross-applies their in-flight edits until neither
diff can be attributed. A charge bans `git stash` — the destructive-git guard
refuses it, and the guard's deny table owns which verbs — and names the
per-worktree substitute instead: copy the working file aside with `cp`,
restore the committed version with `git show HEAD:<path> > <path>`, run the
test, then copy the working file back.

Read-only delegations need no worktree; say "do not edit any file" in the
charge and mean it. Review delegates return findings only. The coordinator
verifies the repair in the checkout that owns the diff. Repairs beyond the
allowance under Delegate or inline are re-charged to a write-delegate in an
isolated worktree; that delegate receives the finding and a commit-specific
sentinel for the diff under repair.

### The shared-checkout exception

Some work no worktree can hold: a large uncommitted build the gate is red on
cannot be committed first, and a worktree branched from HEAD would not contain
the code under repair. A delegate may then run in the main checkout under
exactly four conditions — one writer at a time; a named file allowlist in the
charge; no commit authority; a `git status` check verified on return. All
four, every time: this is the one loosening of the isolation rule, and the
conditions are what carry the safety.

## Verifying the done-claim

A delegate's done-claim is a claim, not a result. Before accepting one:

- run the gate against the delegate's work — its green report is not your green;
- check every coverage row in the charge went red-then-green — a missed row is
  a missed case, found now instead of at review;
- run `git status` in the worktree it used — files touched outside the charge
  are a finding, not a footnote;
- resolve every identifier in an absence, exclusion, or withholding claim to a
  real thing before accepting it — a misspelled identifier passes its contract
  by asserting the absence of something that never existed;
- probe at least one accepted behavior independently of the delegate's own
  tests — through the built binary or a fixture the delegate did not author.
  Delegates write the tests that pin their work, so gate-green alone cannot
  tell a correct build from a self-consistent wrong one; keep the probe
  constant across a batch, not front-loaded;
- spot-check the citations in any summary before folding it into your own
  report — a delegate's confident paraphrase inherits none of the source's
  authority.

Accepting a done-claim unverified is grading your own work at one remove —
invariant #1 with extra steps.

Report every verification round in one line, like a ladder move: accepted, or
what was missed and where the fix went. Reject a miss and re-charge its concrete
repair to a write-delegate when it exceeds the canonical allowance under
Delegate or inline. Recurring misses across delegates are a charge defect, not a
delegate defect — tighten the rows in the charge before re-sending it.
