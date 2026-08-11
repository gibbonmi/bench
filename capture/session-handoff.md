# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` at `ca2b0ee2` — tree clean; active `bench spec build` lifecycle state is durable, not a working-tree diff
Spec: `specs/ticket-bundle-refusal/spec.md` (staged, signed off), `specs/checkpoint-scoped-review/spec.md` (staged, signed off), `specs/axi-spec-build-complete/spec.md` (staged, run active), `specs/axi-coherent-diff/spec.md` (staged), `specs/axi-query-disclosure/spec.md` (staged), `specs/single-build-serial-gate/spec.md` (staged), `specs/axi-compatibility-oracle/spec.md` (staged, superseded — pending abandon)

## State

Two new specs landed this session, both Sol-falsified (two rounds and one
round respectively, all findings repaired) and reviewer-signed:

- `ticket-bundle-refusal` (`1d236fab`): assign refuses any ticket over
  rows > 5 (deflation-resistant count, non-`R` ranges by span) or closure
  tokens > 15 without a header-block `Bundle-approved:` line
  (reviewer-owned, inert at or after the first `##`); `craft-tickets` drops
  the author-asserted keep-together exception for lifecycle work;
  `craft-spec` + `/bench-write-spec` gain the single-ticket-landing
  disclosure. Four implied tickets: refusal core (TB1+TB3), override and
  anchoring (TB2+TB4, blocked by core), craft-tickets prose (TB5),
  craft-spec/write-spec prose (TB6). Closed decisions: bounds values, no
  fence-entry dimension, grammar-independence, header-block anchoring.
- `checkpoint-scoped-review` (`ca2b0ee2`): advisory three-axis review per
  assignment on the write delegate's return, before done-claim verification
  and checkpoint (ordering: return → review → dispositions → verify settled
  tree → checkpoint); dispositions closed and authority-split (fixed =
  write delegate; risk-accepted = reviewer; deferred); subject-bound
  evidence line per assignment in this file, consolidated into the retro.
  Prose-only across four command/skill files; FT184/FT200 seams untouched.

Implementation order is coupled: implement `checkpoint-scoped-review` first
(prose fences, no AXI collision), because its cadence is the dogfood oracle
for the `ticket-bundle-refusal` build that follows. The
`ticket-bundle-refusal` code tickets touch `internal/specbuild/`, which sits
inside the active AXI candidate's fence — prefer closing the AXI run first;
a main-tip move before then just forces routine recomposition at promote.

`axi-spec-build-complete`'s run is active, subject
`f1b951edc6c653a6f42f1e1f5f32a3eb70e2a377`, next operation `review` (the
prior candidate's review round left five accepted findings awaiting reviewer
triage — S1 test-repo tension (contestable), S2 comment provenance label,
P1 `git.Root()` before help/usage dispatch, P2 non-exact refresh-refusal
remedy, C1 `Runs` TOCTOU — and the subject has since moved, so the composed
review must run fresh against `f1b951ed` regardless). Build order unchanged:
`axi-spec-build-complete`, then `axi-coherent-diff`, then
`axi-query-disclosure`.

## Next command

`/bench-implement-spec specs/checkpoint-scoped-review/spec.md`

Then `/bench-implement-spec specs/ticket-bundle-refusal/spec.md` under the
new cadence. The AXI run resumes separately with
`bench spec build review axi-spec-build-complete` (fresh composed review of
`f1b951ed`, five prior findings to re-triage).

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
