# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — tree clean; the repair round below is durable `bench spec build` lifecycle state, not a working-tree diff
Spec: `specs/axi-spec-build-complete/spec.md` (Status: staged), `specs/axi-coherent-diff/spec.md` (Status: staged), `specs/axi-query-disclosure/spec.md` (Status: staged), `specs/single-build-serial-gate/spec.md` (Status: staged), `specs/axi-compatibility-oracle/spec.md` (staged, superseded — pending abandon)

## State

`axi-spec-build-complete`'s spec-build run (`b6f79889...`) is active, candidate
`123edbab919d2462c73f206dd6623aba5feccc9f`. The previously accepted SB4 finding
(P1-nondigest-state-name / C1-nondigest-state-name: a valid retained record
under a non-digest `*.json` filename rendered healthy instead of its own
diagnostic) is repaired and integrated — ticket
`specs/axi-spec-build-complete/tickets/repair-nondigest-state-finding.md`,
assignment `73324cb44cd716dcf9f4b2600ce46c41`. `Service.Runs` now diagnoses
that case as `nondigest_name`.

A fresh Terra/xhigh review of the whole composed candidate (receipt digest
`925a9ec921f57225b529def2885158427ba9858ecacf152b64402751357e250b`) confirmed
that repair and found one new pair of `accepted` findings not yet acted on: a
partial-reclaim failure renders two `help[]` blocks — `RenderReclamation`
then `RenderRefusal` each append one, violating the spec's single-envelope
rule (`internal/specbuild/render.go:251-287`, `spec.md:24`) — and the existing
partial-reclaim test never renders that path, so the duplicate goes unpinned
(`internal/specbuild/reclaim_test.go:435-469`). This round's repair authority
was scoped to exactly the nondigest finding, so promotion was withheld rather
than widened: `bench spec build status axi-spec-build-complete --full` shows
`next: bench spec build assign axi-spec-build-complete` — the run is waiting
on a reviewer decision to authorize a repair ticket for the new finding (or a
different route).

Build order unchanged: `axi-spec-build-complete` next, then
`axi-coherent-diff`, then `axi-query-disclosure`. FT185 composition stays
closed: promote's gate payload is composed when FT185 exists, never
re-derived. `axi-compatibility-oracle`'s reclaim is still deferred per the
prior handoff.

## Next command

Reviewer decision on the new duplicate-`help[]` finding, then
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
