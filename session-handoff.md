# Session handoff

Repository: `bench` (origin `github.com/gibbonmi/bench`) — branch `main`,
checked out at `~/workspace/bench` on the machine that wrote this. Everything
below is executable from a cold start; no conversation history is needed.

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
- **The next writer of this file should dogfood `bench handoff` instead.** The
  command exists in the tree and is what this build is for; this rewrite was
  done by hand deliberately, because ungated code should not be the thing that
  writes the reviewer's cold-start artifact. Once the gate is green, that
  reasoning expires.
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
- **Three commits are unpushed and the gate is stale** (gated tree `5b1cc38`,
  work tree `eae35a5`). The last green gate is still `4ea4880`; every commit
  since is doc-only and landed with plain `git commit`.
- **A retired spec cited as `specs/<slug>.md` in `ROADMAP.md` reads as a
  dangling row.** `roadmapReconcileCounts` treats any such path as a row's spec
  pointer, so prose evidence names retired specs by bare slug. Cost one
  follow-up commit on 2026-07-25.
- `bench structure` reports 17 issues, all pre-existing. Read
  `projects/benchkit.md`'s cold-session notes before touching `internal/canary` —
  the nested-run trap deadlocks rather than fails, and this file deliberately
  does not restate the commands.

## Next command

`/bench-implement-spec specs/session-handoff-emission.md` — resume the
in-flight build; the remaining work is verification, not authoring. Interactive
rather than `bench shift`, because nine stories route to mid and one to top,
which fails `craft-line`'s venue-routing test for an unattended loop. In Codex,
that is `$bench-implement-spec`.

Nothing else should touch this tree first. A second diff on top of an
uncommitted 605-line build makes the gate answer for two changes at once.

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
