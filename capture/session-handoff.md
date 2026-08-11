# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `215a7fc`, 1 dirty path, 22 unpushed commits
Spec: `specs/axi-coherent-diff/spec.md` (Status: staged), `specs/axi-query-disclosure/spec.md` (Status: staged), `specs/pocock-guidance-doctrine/spec.md` (Status: staged), `specs/single-build-serial-gate/spec.md` (Status: staged)
Gate: green at `9841e3f` — stale, work tree `add9144`

## State

Pocock-alignment Spec C is staged at
`specs/pocock-guidance-doctrine/spec.md`. Its compiled decision map is the
byte-identical retained `pocock-alignment` source, limited to decisions #4, #5,
#6, #8, #9, and #10. The upstream comparison was refreshed through
`mattpocock/skills` `84fdeff`; prototype branch retention remains deliberately
unadopted in favor of the reviewed discard-after-verdict rule.

Fable medium with bypass permissions rejected the first draft on four omissions;
all were repaired. A fresh Terra xhigh falsification pass accepted the 24-row
candidate with zero Standards, Spec, or Coverage findings. Closed decisions stay
closed: gate authority and all three parallel review axes remain; ticket fences
become advisory `Writes:` notes; the one-build Opus `/bench-debug` override is not
standing; unowned `bench-link` convergence remains limited to canonical same-file
adapter targets, with divergent and byte-identical foreign targets refused. No
push is authorized. The other staged specs remain `axi-coherent-diff`,
`axi-query-disclosure`, and `single-build-serial-gate`.

## Next command

`$bench-implement-spec`

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
