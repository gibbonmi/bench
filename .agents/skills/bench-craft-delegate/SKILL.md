---
name: craft-delegate
description: The delegation discipline — when to spawn a subagent, how to charge and scope it, when it needs an isolated worktree, and how a done-claim is verified. Use whenever spawning any delegate (an axis review, a scoped build, a fan-out search) or accepting one's result. The model/effort half of the decision lives in craft-line.
index: spawning a delegate / verifying a delegate's done-claim
---

# Delegating without losing the plot

A delegate buys parallelism and isolates a heavy read-set; misjudged, it
costs the thing that matters: the work happens unseen.

## Delegate or inline

The coordinator scopes, routes, and verifies work; a write-delegate authors
code. The inline allowance is exactly one source-line insertion, one
source-line deletion, or one source-line replacement in code or tests. A
replacement counts as one correction. This allowance spans the current
reviewer request and does not reset when work is split into tasks, slices,
delegates, or verification rounds. A no-spec change admitted by the
lighter-path threshold in `.bench/BENCH.md` may also remain inline — the
only other allowance. All other code authorship runs as a write-delegation
in an isolated worktree; read-only coordination stays inline. Never
delegate a decision the reviewer owns.
This policy is capability- and posture-aware: a harness that cannot spawn a
write subagent, or a reviewer who has prohibited delegation for this work,
never falls back inline beyond the allowances above — a spec-doc-only
correction is not a silent exception to either stop. Either posture stops
before editing and emits one executable resume handoff to a subagent-capable
harness — the repository path, the working branch or worktree, the spec or
change name, the destination harness, and that harness's exact invocation.
Surface before spawn any delegation that changes who performs requested work.

## The charge

A delegate has no conversation memory: everything it needs is in the prompt —
objective, inputs by path, seam, return shape, budget. Route model and effort with
`craft-line`: name the resolved bound model id on every call (omission inherits your
model and may silently escalate); put effort and iteration cap in the charge.
An own-family reviewer uses the harness's native agent surface, never that family's CLI.
Cross-family reviews and the no-native-surface fallback use the exact recipes in
`references/cross-harness-reviewers.md`. Prefer compressed inputs — the named
decision source, exact passages, coverage rows, and the fence's fixture-and-seam
inventory — so the delegate uses prior art instead of re-deriving it.
Name exemplar files to mirror when one exists. A charge that extends an
enumerated family names every registry the family already appears in, traced
from one existing sibling through the tree; a registry the charge does not
name is one the delegate will miss.

A write-delegation from a spec carries its stories' coverage rows — behavior,
seam, why it catches the failure — every time, and requires the delegate to
show each row red before the edit and green after; check each slice against
`craft-spec`'s "Slicing a build for delegates" first. Name the mutation that
breaks the change's central property, and require the
delegate to apply it to its own finished work, report the observed result, and
add the missing row when the mutation comes back silently green.
A delegate blocked by a defect outside its fence stops and reports rather
than fixing out of fence. A new worktree charge opens with the stale-base check
(`git merge --ff-only main`, verify HEAD equals main, stop if denied); dependent
tickets in a reviewed spec chain share the retained integration source and
verify its expected tip. A fix-pass charge names a commit-specific sentinel.
A ticket delegate returns focused evidence and its own mutation probe from
its worktree; it does not land the diff. The coordinator probes the exact
returned tree independently before landing it. The
coordinator probe's mutation kind differs from the delegate author's mutation
kind. It also differs in site from every probe the delegate ran: a second probe
at the same site is vacuous, and a vacuous probe is indistinguishable from a
pass. A repeat site is not independent evidence.
```
Implement story 3 of specs/retry-backoff/spec.md. Stale-base check first.
Coverage rows: [rows]. Effort: medium, ~3 iterations. Stop at diff ready;
return the red→green log per row. Self-probe: apply the central-property
mutation; report the observed result and the mutation's kind (omission or swap).
```
Good — rows make the done-claim verifiable, and the self-probe names its kind.

## Scope

One delegate, one coherent unit: one axis, one story, one search question.

## Isolation

A write-delegation runs in an isolated worktree (`bench worktree`), so stray edits can't
land in reviewer-owned files. Concurrent delegates get separate worktrees, one each — the
harness's own `isolation: worktree` cannot cut a second one (its request ID derives from the session ID alone),
so the coordinator runs `bench worktree create --request <opaque-id> --label <work-item>`
once per delegate, then addresses it by label (`bench worktree exec "<label>" -- <command>`,
`bench worktree path "<label>"`), never a cached path.
Share a worktree only when a delegate's work depends on another's output; reviewed dependent
tickets share one retained integration source, and each charge names its root and expected tip.
The whole-tree gate is serialized: a write-delegate stops at diff-ready with focused tests green;
the coordinator runs `bench commit` per worktree, one at a time.
A worktree isolates the working tree, not the repo-global stash stack a
concurrent delegate shares. A charge bans `git stash` — the destructive-git
guard refuses it — and names the substitute: `cp` the working file aside,
restore the committed version with `git show HEAD:<path> > <path>`, test,
then copy it back. The copy lives inside the delegate's own worktree
under a unique name, and every restore names exact files, never a
glob.
Read-only delegations need no worktree; say "do not edit any file" and mean
it. A large uncommitted build no worktree can hold may run in the main
checkout under exactly four conditions: one writer, a named file allowlist,
no commit authority, a `git status` check verified on return.

## Verifying the done-claim

A delegate's done-claim is a claim, not a result. Before accepting one: run
the gate against its work; check every coverage row went red-then-green;
run `git status` in the worktree it used; resolve every identifier in an
absence or exclusion claim to a real thing; probe one accepted behavior
independently of the delegate's own tests, kept constant across a batch;
spot-check citations before folding a summary in.
Report every verification round in one line: accepted, or what was missed
and where the fix went. Repairs beyond the allowance under Delegate or
inline continue the authoring delegate for its own slice when the harness can
resume it; otherwise a fresh charge in an isolated worktree carries the
finding and a sentinel. The coordinator verifies the repair in the checkout
that owns the diff.
Acceptance closes an independent worktree after its slice lands: the coordinator runs `bench worktree release --request <opaque-id> <path>` for it. A reviewed
dependent chain remains retained through explicit source review; only
`bench worktree land` releases it after the reviewed source publishes.
