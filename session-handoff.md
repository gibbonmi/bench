# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `1b04333`, 1 dirty path, 14 unpushed commits
Spec: `specs/implement-spec-full-run.md` (Status: staged)
Gate: green at `56eedcc` — stale, work tree `a266821`

## State

- **Two approved specs are queued; build FT91 first.** The reviewer agreed
  2026-07-28 that `specs/ft91-canary-contract-scoping.md` lands before
  `specs/implement-spec-full-run.md`, because FT91 cuts the dev gate from
  267 s toward a ≤100 s solo-canary target and every later build pays the
  current price otherwise. The two specs do not collide on files.
- **FT91 (stage 1 of the eighth arm)** scopes each `behavior-owned` canary
  fixture's nested gate to the one contract package owning its EXPECT via
  subfamily directories, with `BENCH_CANARY_CONTRACT_PACKAGE` narrowing the
  inner contract argv, loud reds for every failure mode, per-group scoped
  vacuity baselines, a kit-only guard pinning the flat set to two named
  relocations, and a ≤100 s acceptance threshold on the ship evidence. Build
  inputs: the spec, `decisions/gate-critical-path.md` (Handoff),
  `decisions/assets/gate-critical-path-timeline.md`, and
  `decisions/assets/behavior-owned-package-bindings.md`.
- **implement-spec-full-run** adds an opt-in `--full` mode to
  `/bench-implement-spec` (implement inline → review in a fresh-context
  delegate → final-check inline, debug on demand, `session-handoff.md`
  rewritten at every phase boundary), plus two shared-rule changes:
  fix-don't-park in `.bench/BENCH.md`, and point-of-use verify hooks in the
  grill, implement, and review phases. 10 stories, 28 coverage rows. It is
  prose plus its observers only — no production Go, no new `bench`
  subcommand. Build input: the spec and `decisions/implement-spec-full-run.md`
  (closed). Twelve decisions are flagged in that spec's header for post-hoc
  veto; they were approved with the spec and stay closed unless reopened.
- **Both specs' falsification passes ran and are folded in.** For
  implement-spec-full-run, a reviewer-directed Codex `gpt-5.6-sol` xhigh pass
  returned BLOCK on the first draft; six of seven findings were verified
  against the tree and folded (section-scoped anchors replacing whole-file
  presence, require/`forbid` pairs on contradiction-prone facts, the
  Handoff item 6 done-claim obligation that had fallen through, corrected
  "already covered" claims). One was rejected in part: `bench handoff` needs
  no phase field, because `render.go` preserves the State section and the
  orchestrator writes the phase there.
- **Map state:** `gate-critical-path` tickets #2 and #3 are the only open
  tickets and gate the artifact-hoist slice, not either queued build. FT91's
  stop condition is a measured dev gate ≤60 s.
- **`ROADMAP.md` is stale** — one row names a retired spec, and neither
  staged spec has a row. The next `/bench-what-next` reconcile owns the
  rewrite; the tree wins meanwhile.

## Next command

`/bench-shape-idea` — the board's leading invocable signal (`decisions`).

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
