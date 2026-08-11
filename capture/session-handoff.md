# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — tree clean; the repair round below is durable `bench spec build` lifecycle state, not a working-tree diff
Spec: `specs/axi-spec-build-complete/spec.md` (Status: staged), `specs/axi-coherent-diff/spec.md` (Status: staged), `specs/axi-query-disclosure/spec.md` (Status: staged), `specs/single-build-serial-gate/spec.md` (Status: staged), `specs/axi-compatibility-oracle/spec.md` (staged, superseded — pending abandon)

## State

`axi-spec-build-complete`'s spec-build run (`b6f79889...`) is active, candidate
`399dca908c7b1e1a4162eb7625e497dfb6786750`. Two reviewer-authorized repair
rounds are integrated: the nondigest-state-name diagnostic
(`repair-nondigest-state-finding.md`) and the duplicate partial-reclaim
`help[]` envelope (`repair-duplicate-help-partial-reclaim.md`). A fresh
Terra/xhigh review of the whole composed candidate (receipt digest not yet
reviewer-triaged, full text at
`/tmp/.../scratchpad/terra-review-output-2.md` — not durable, re-run the
review if that scratch file is gone) confirms both repairs and independently
found five new `accepted` findings, none yet authorized for repair:

- Standards S1: `disclosure_observation.go`/`disclosure_test.go` create real
  git repos and run OS processes, which the reviewer read as violating
  `projects/benchkit.md:210-212`'s "ordinary tests create no repositories"
  rule — contestable: the disclosure fixture harness is the spec's SB2
  real-service observation seam, and an earlier Fable round's ticket already
  says "Preserve S3 and S5 exactly as risk-accepted review judgments" over
  what may be this same tension under different IDs. Needs reviewer judgment,
  not an assumed repair.
- Standards S2: `internal/specbuild/reclaim_test.go:481-482`'s new comment
  keeps the ticket-provenance label `(DH1)`, against `bench-craft-comments`'
  timeless-current-state rule. Concrete, trivial one-line fix.
- Spec P1: `cmd/bench/specbuild.go`'s `specBuildCommand` calls `git.Root()`
  before dispatching help spellings or malformed argv, so `help`/`--help`/`-h`
  and bad argv outside a git checkout return error/1 instead of the required
  catalog/0 or usage/2 (`spec.md:28,33-34`).
- Spec P2: the stale/spent-refresh refusal remedy in
  `internal/specbuild/disclosure.go`'s `RefusalForClass` emits a generic
  `assign` action with open placeholders instead of the exact original
  ticket, request, and `--refresh` receipt — not the class-exact remedy
  story 2 requires.
- Coverage C1: `Service.Runs` in `internal/specbuild/state.go` classifies a
  directory entry from the `os.ReadDir` snapshot, then reopens it by path for
  `os.ReadFile` — a TOCTOU window a special file could exploit between
  classification and read; no test drives that race.

Promotion was withheld (review not clean). `bench spec build status
axi-spec-build-complete --full` shows `next: bench spec build assign
axi-spec-build-complete` — the run is waiting on a reviewer decision on which
(if any) of the five new findings to authorize as repair tickets.

Build order unchanged: `axi-spec-build-complete` next, then
`axi-coherent-diff`, then `axi-query-disclosure`. FT185 composition stays
closed: promote's gate payload is composed when FT185 exists, never
re-derived. `axi-compatibility-oracle`'s reclaim is still deferred per the
prior handoff.

## Next command

Reviewer decision on the five new findings, then
`/bench-implement-spec --full specs/axi-spec-build-complete/spec.md` resumes
the same repair → checkpoint → integrate → review → promote cycle.

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
