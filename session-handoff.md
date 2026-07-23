# Session handoff

Repository: /home/mgibs/workspace/bench — branch `main`. Everything below is
executable from a cold start; no conversation history is needed.

## State

- **The roadmap drain is committed** (`dcb24b9`): FT88 removed — the gate owner
  record, the PID-naming refusal, process-group teardown, the two-leg marker
  wait, the contract-harness reap, and self-attributing conformance diags all
  shipped and the spec retired. `IDEAS.md` and `.bench/learnings.md` are both
  empty; the journal's reproduction-economics entry became roadmap row FT112.
- **FT109 is shipped** (`d0263c2`, gate green, row removed). `bench status` now
  carries a `handoff` row (severity 11) reporting how many commits have landed
  since this file was last written; it is silent on an absent, untracked, or
  mid-rewrite handoff, and ranks last so it leads a quiet cold-pickup board
  without displacing a red gate or a dirty tree. `AGENTS.md`'s phase-close
  paragraph carries the rewrite-in-full and conflict rules; the Shape section
  below is this file's own template. Built on the light path with the reviewer's
  explicit OK — no spec.
- **One approved deviation, now closed.** The FT109 row specified a `written-at:`
  line inside the handoff. It was dropped: a self-reported date is the same
  remembered-not-computed defect the row exists to close, and it cannot name its
  own commit anyway, since the handoff lands in the commit it describes. The age
  comes from `git log -1 -- session-handoff.md` instead.
- **`specs/` is empty**; nothing is mid-flight and the tree is clean.
- **Known advisory debt:** `bench structure` reports 10 violations (crowded
  `internal/adopt/` and `internal/contract/surface/`, plus seven over-length
  files). The gate is green with them and no roadmap row covers them.
- **Unpushed:** `main` is ahead of origin. Pushing is the reviewer's call.

## Next command

`/bench-write-spec` for FT87 slice 3 — the command-wide parser and
security-evidence capability — in a fresh **mid-tier** session, the profile's
spec default. It leads the roadmap's recommended sequence and its decision map
(FT87 tickets #7 and #8) is already closed, so the spec compiles from a reviewed
source rather than the batch-drain override.

Fresh evidence for that row from this session: `bench commit` rejects a
directory path (`internal/status/`) and requires each file named, which is one
of slice 3's listed gaps.

## Shape

Rewritten in full at every phase close, pruned rather than accreted: a fresh
session pays for every line it reads cold, so drop anything it would not act on.
Keep the three sections above — **State** (what is true now, including anything
uncommitted), **Next command** (the exact harness-native invocation, not a
description of it), and this one.

The handoff carries no date of its own. `bench status` computes its age from the
commit that last wrote this file and reports a `handoff` row once anything has
landed since. Where this document and the tree disagree, the tree wins.
