# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `57799ee`, 1 dirty path, 6 unpushed commits
Spec: `specs/implement-spec-full-run.md` (Status: staged)
Gate: green at `073cd6e` — stale, work tree `6e566d0`

## State

- **Two approved specs are queued; build FT91 first.** The reviewer's 2026-07-28
  sequencing stands: FT91 lands before `specs/implement-spec-full-run.md`, because
  every later build pays the current gate price otherwise. The two do not collide
  on files.
- **FT91 stage 2 (`specs/ft91-canary-compiled-bites.md`, approved) is the build.**
  Behavior-owned fixtures stop spawning nested gates: one `go test -c` per contract
  package group, then one invocation of that binary per fixture root with the
  contract subject root swapped to the materialized tree. The canary sweep still
  owns bite, did-not-bite, and the vacuity baseline. The gate's inner-mode contract
  narrowing, its environment variable, and their tests retire. 12 stories, 28
  coverage rows, three seams. Build inputs: the spec,
  `decisions/gate-critical-path.md` (Handoff, rewritten for stage 2), and
  `decisions/assets/gate-critical-path-timeline.md`.
- **Four deviations from the map were approved with the spec** and stay closed:
  compile from the swept root rather than the kit checkout; keep the `GOMAXPROCS`
  pin on bite invocations (the map called it moot — true of the nested gate, not of
  the test binary underneath); single-source `BENCH_CONTRACT_ROOT` across the
  contract helper, which the Handoff had said stays untouched; and drop the
  phase-manifest refusal for this family, whose mechanism no longer exists.
- **Stage 2's mechanism is prototyped, not assumed.** `go test -c
  ./internal/contract/axi` compiles in 0.5 s warm, and the binary run against the
  materialized `roadmap-regressed` tree reds in 0.49 s with its EXPECT present in
  test-level output. The remaining 32 fixtures are unverified until the sweep runs —
  an EXPECT observable only in gate framing is an ordinary did-not-bite red, fixed
  at the owning test's failure message (story 9).
- **The falsification pass is folded in.** A reviewer-directed Codex `gpt-5.6-sol`
  xhigh pass returned block on the first draft. Its load-bearing finding: every
  sweep assertion runs through an injected fake, so a change emitting correct call
  metadata while the real runner still spawned gates would have passed the whole
  map. There is now a row on the default runner's actual dispatch. Two findings were
  rejected on inspection — the linked-repo refusal claim, and a doc-authoring line
  bump for stories 10 and 11, whose deviations are named in the spec instead.
- **implement-spec-full-run** is unchanged and still staged: an opt-in `--full` mode
  for `/bench-implement-spec`, plus the fix-don't-park and point-of-use verify-hook
  shared-rule changes. Prose plus its observers only. Twelve decisions flagged in
  that spec's header were approved with it and stay closed.
- **Map state:** `gate-critical-path` #2 and #3 are the only open tickets and gate
  the artifact-hoist slice that follows stage 2, not either queued build. FT91's
  stop condition is a measured dev gate ≤60 s, judged after the artifact hoist.
- **`ROADMAP.md` is stale** — one row names a retired spec, and neither staged spec
  has a row. The next `/bench-what-next` reconcile owns the rewrite; the tree wins
  meanwhile.

## Next command

`/bench-implement-spec — FT152, specs/implement-spec-full-run.md (staged, amended; the build is the only remaining action)`

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
