# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `1d58ffd`, 2 dirty paths, 8 unpushed commits
Spec: `specs/ft123-ft124-session-tax-readers/spec.md` (Status: staged)
Gate: green at `0bd5adc` — stale, work tree `56e2972`

## State

- **Implementation and semantic review are complete.** Five ticket commits landed
  from `e4496cc` through `1893e16`. A fresh-context three-axis review found two
  spec defects and two missing coverage proofs, with no standards or contestable
  design findings.
- **All review findings are repaired.** `34ee835` enforces the exec separator
  boundary, refuses incomplete package streams, proves TOON render refusal, and
  preserves colon-ending skip prose. Its full repository gate passed every phase,
  and the review pickup artifact is cleared.
- **Cross-harness falsification is closed.** The approved `gpt-5.6-sol` / high
  Codex CLI pass ran read-only with approval disabled over `1c42a80..34ee835`.
  Its sole counterexample, exact flag-like worktree labels, is repaired in
  `1d58ffd`; that repair's full repository gate passed every phase.
- **Final-check is ready.** All five tickets are accepted, no review pickup
  artifact remains, and `bench coverage --check` validates all 17 rows. The
  spec-aware landing commit is next and will flip `Status: staged` to
  `Status: implemented`.
- **Uncommitted state is intentional.** The staged spec, completed ticket files,
  and this handoff await the final spec commit. `IDEAS.md` is an unrelated
  pre-existing edit and must remain untouched and outside every feature commit.

## Next command

`bench commit -m "spec-implemented: ft123-ft124-session-tax-readers" --spec ft123-ft124-session-tax-readers -- session-handoff.md specs/ft123-ft124-session-tax-readers`

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
