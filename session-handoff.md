# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `5b2955b`, 11 dirty paths, 4 unpushed commits
Spec: `specs/session-handoff-emission.md` (Status: staged)
Gate: green at `5b1cc38` — stale, work tree `53c6b1d`

## State

- **FT122's build is in the tree, uncommitted and ungated.** `internal/handoff/`
  (7 files, 605 lines) plus three runtime contract files, a conformance file,
  and edits to `.bench/BENCH.md`, `bin/bench.sh`, `cmd/bench/main.go`,
  `internal/status/handoff.go`, and `internal/conformance/subcommand_routing_test.go`.
  `go build ./...` passes and `go test ./internal/handoff/` is green; the
  contract and conformance layers have **not** been run, and neither has the
  gate. Test names cover all 26 stories, but coverage-by-name is not a verdict —
  the first `bench gate` on this tree is still the first real one.
- **`specs/session-handoff-emission.md` is still `Status: staged`.** Flip it
  with `bench commit --spec` on the green commit, not by hand.
- **This file was written by the in-tree `bench handoff`, and doing so found a
  defect.** The header block is correct and is now the single source for
  repository, path, branch, HEAD, dirty/unpushed counts, spec status, and gate
  verdict — do not restate those in `## State`. But the default `## Next
  command` derivation resolved to `commit on green / /bench-final-check / push`,
  the status board's leading signal, which is a *signal*, not the next action;
  the real one had to be supplied with `--next`. Story 8 derives that field from
  the same signals `bench status` ranks, and a generic dirty-tree row outranks
  the thing a cold session should actually do. Decide before the gate run
  whether story 8's derivation needs a rule that skips non-actionable signals,
  or whether `--next` is simply mandatory at a phase close.
- **Story 15 carries a reviewer-approved top-tier bump** — the scaffolded
  skeleton's guidance prose, `gpt-5.6-sol` / high under the leverage override,
  bounded to that prose. Story 14 (moving this file's `## Shape` text into the
  binary) stays cheap because it is transcription. Nine stories route to mid,
  the rest cheap.
- **The falsification pass is folded in** (`2f51fe9`); `decisions/session-handoff-emission.md`
  carries the **[veto]** marks. The finding worth carrying: a bare `gate: green`
  would have emitted a cached verdict whose tree had already moved, so the field
  renders verdict, cached tree, and staleness together.
- **The drain ran (`120fff9`, `66991a6`) — `IDEAS.md` is empty.** Four rows
  landed from it: FT123 (`bench worktree path`/`exec`, which must render
  `~`-relative paths), FT124 (`bench test`), FT125 (the two reader-slice
  candidates, merged), and FT126 (`bench roadmap --context` should report the
  reconcile's workload boundary and a discrepancies block rather than raw
  fields). FT122 got the row it had been missing.
- **The last green gate is `4ea4880`**; every commit since is doc-only and
  landed with plain `git commit`. (Counts and tree hashes are in the header.)
- **A retired spec cited as `specs/<slug>.md` in `ROADMAP.md` reads as a
  dangling row.** `roadmapReconcileCounts` treats any such path as a row's spec
  pointer, so prose evidence names retired specs by bare slug. Cost one
  follow-up commit on 2026-07-25.
- `bench structure` reports 17 issues, all pre-existing. Read
  `projects/benchkit.md`'s cold-session notes before touching `internal/canary` —
  the nested-run trap deadlocks rather than fails, and this file deliberately
  does not restate the commands.

## Next command

`/bench-implement-spec specs/session-handoff-emission.md`

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
