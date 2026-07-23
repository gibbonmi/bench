# Session handoff

Repository: /home/mgibs/workspace/bench — branch `main`. Everything below is
executable from a cold start; no conversation history is needed.

## State

- **The drain is committed** (`dcb24b9`): FT88 removed from `ROADMAP.md` — the
  gate owner record, the PID-naming refusal, process-group teardown, the two-leg
  marker wait, the contract-harness reap, and self-attributing conformance diags
  all shipped and the spec retired. `IDEAS.md` and `.bench/learnings.md` are both
  empty; the journal's reproduction-economics entry became roadmap row FT112.
- **FT109 is built and uncommitted** — the handoff-shape row, taken on the light
  path with the reviewer's explicit OK (no spec). Three parts: the phase-close
  paragraph in `AGENTS.md` gained the rewrite-in-full and conflict rules, this
  file gained the shape section below, and `bench status` gained a `handoff` row
  (severity 11) reporting how many commits have landed since the handoff was last
  written. Unit tests live at the `internal/status` seam; two mutations were
  verified red.
- **One approved deviation from the FT109 row.** The row specified a `written-at:`
  line inside the handoff. That was dropped: a self-reported date is the same
  remembered-not-computed defect the row exists to close, and it cannot name its
  own commit anyway, since the handoff lands in the commit it describes. The age
  is read from git history instead — `git log -1 -- session-handoff.md` — so there
  is nothing to maintain and nothing that can lie.
- **`specs/` is empty.** No spec is staged and no build is mid-flight beyond the
  uncommitted FT109 work above.
- **Known advisory debt:** `bench structure` reports 10 violations (crowded
  `internal/adopt/` and `internal/contract/surface/`, plus seven over-length
  files). The gate is green with them and no roadmap row covers them.
- **Unpushed:** `main` is ahead of origin by several commits. Pushing is the
  reviewer's call.

## Next command

`bench commit -m "<message>" AGENTS.md ROADMAP.md session-handoff.md internal/status/`
to land FT109 on a green gate. Its roadmap row is already removed in the same
change, so nothing remains to reconcile afterward.

After that, the roadmap's recommended sequence leads with `/bench-write-spec` for
FT87 slice 3 (command-wide parser and security-evidence capability) in a fresh
mid-tier session — its decision map is already closed.

## Shape

Rewritten in full at every phase close, pruned rather than accreted: a fresh
session pays for every line it reads cold, so drop anything it would not act on.
Keep the three sections above — **State** (what is true now, including anything
uncommitted), **Next command** (the exact harness-native invocation, not a
description of it), and this one.

The handoff carries no date of its own. `bench status` computes its age from the
commit that last wrote this file and reports a `handoff` row once anything has
landed since. Where this document and the tree disagree, the tree wins.
