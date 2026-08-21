# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — clean once the spec-retire commit lands; commits unpushed.
Spec: none staged. `land-executable-freshness` landed and is retired.
Gate: green at the landing commit.

## State

FT242 shipped at `84b7c4b0`: `bench worktree land` proves its own executable
through `freshness.Verify` before any repository proof, only where the
repository declares Go build inputs, with `--resume` exempt. FT242 is closed on
the board and `roadmap/FT242.md` is deleted. Durable content promoted to
`projects/benchkit.md`'s cold-session notes: the proof is self-attestation, so
the gate's private exact-source build and the operator's sanctioned rebuild stay
the independent roots.

Folded into that landing at the reviewer's explicit direction: the `regroup`
example profile is retired across the packed-asset list, the consumer payload,
the npm files list, the profile note, and its gitignore line;
`COMPLIANCE_ASSESSMENT.md` is removed; `ui_example/` is untracked but kept on
disk; the README lost the Regroup walkthrough, kept the design-system oracle it
documented, and gained a dev-build freshness note.

Pending capture for the next drain: `capture/retros/land-executable-freshness.md`,
a refreshed `capture/agent-performance/claude-models.md`, and two open
`capture/learnings.md` entries (the spec-map ordering row, and `bench handoff`
overwriting a mid-phase next command). The retro names two further CLI defects
worth roadmap items — `bench worktree release` refusing with a `--discard-ignored`
flag its own help does not list, and `bench diff` rejecting the `--source-tip`
that `bench preflight` accepts.

`ASSESSMENT.md:99` now cites a README example this landing removed; left alone
because that document is a dated record, not live guidance.

## Next command

`/bench-drain` — verdict the retro, the scorecard, and the two learnings
entries into roadmap items.

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
