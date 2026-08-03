# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `5b4e240`, 1 dirty path, 13 unpushed commits
Spec: `specs/pre-push-guard-visibility/spec.md` (Status: staged)
Gate: green at `94889f6` — stale, work tree `396fbd6`

## State

- **FT164 is implemented and terminal.** Promoted at `83c630e` (candidate
  `cbdcf7d`, run terminal). All 30 coverage rows closed; 27 section-scoped
  needles each with a mutation-table row; the example-agreement check is live
  in the conformance registry with its canary family. Retro at
  `capture/retros/ft164-ticket-contracts.md`, including the open reviewer-veto
  surface (duplicate re-derivation sentence, `###` fence-heading workaround,
  ParseTicket doc comment, inventory grouping label, held probe-kind entry,
  anchor-file structure-budget overage).
- **Reviewer's next action: run the Codex hotfix** (prompt already delivered in
  conversation; hotfix-class, `$bench-debug`, light-path tickets): relax
  `bench spec build start`'s exact-green precondition to accept composed
  green (reduced/partial verdict + intact inherited evidence), fix its
  remediation text, and make zero-checkpoint recomposition trivial instead of
  feeding git an empty patch. Both defects are parked in `capture/IDEAS.md`.
  Nothing now holds `internal/specbuild` or the tip — safe to run.
- Two open learnings entries (gate-subject mutation; tree-freeze during active
  runs) and one pending retro drain await `/bench-what-next`.
- FT173's row carries the reviewer's per-command AXI acceptance conditions
  (recorded 2026-08-02); FT174 is unblocked by FT164's identifier decisions.
- Nothing pushed; push is the reviewer's call.

## Next command

`/bench-what-next`

## Shape

Rewritten in full at every phase close, pruned rather than accreted: a fresh
session pays for every line it reads cold, so drop anything it would not act on.
Operational gotchas are placed by lifetime, not copied here: one that recurs across
phases belongs in `projects/benchkit.md`'s cold-session notes, and one scoped to a
build belongs in that spec's coverage rows. This file names at most when you'll hit
one, never the command — a second copy drifts from the source.
Keep the three sections above — **State** (what is true now, including anything
uncommitted), **Next command** (the exact harness-native invocation, not a
description of it), and this one.

The handoff carries no date of its own. `bench status` computes its age from the
commit that last wrote this file and reports a `handoff` row once anything has
landed since. Where this document and the tree disagree, the tree wins.
