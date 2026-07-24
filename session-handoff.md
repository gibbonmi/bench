# Session handoff

Repository: /home/mgibs/workspace/bench — branch `main`. Everything below is
executable from a cold start; no conversation history is needed.

## State

- **FT87 slice 3 is built, reviewed, and its defects are closed.**
  `specs/cli-grammar-and-capability-evidence.md` is `Status: implemented`, 13
  stories in, `bench coverage --check` green on 30 rows. The gate is green at
  HEAD — one full run over both round-3 commits.
- **Round 3 closed the 11 remaining defects** (`6b44e26` docs, `6164844` fixes).
  Four changed behavior:
  - `readSkipTally` now returns an error. An absent log is still an empty tally,
    but any other read failure is diagnosed on stderr and turns the run red under
    `BENCH_REQUIRE_CAPABILITIES=1`. Before, an unreadable log read exactly like a
    fully capable runner — silent de-enforcement in the release workflows.
  - `gateEnv`'s stripping of `BENCH_SKIP_LOG` is now asserted, so a canary's inner
    run can no longer contaminate the outer tally undetected.
  - `bench commit -m x <dir>` works on a *deleted* directory. `isDir` falls back to
    the index when the working tree has nothing to Lstat, so `rm -r sub` then
    naming `sub` commits its removals instead of reporting every child as an
    unexplained offender. The sibling-`subdir` and outside-file guards still hold.
  - `bounds.TestDeadline` clamps negatives to the floor and saturates at
    `math.MaxInt64` instead of wrapping to a negative deadline.
  The rest were comment, import-grouping, profile, and spec-row corrections. The
  profile now documents the `BENCH_REQUIRE_CAPABILITIES=1` knob and the three
  conformance checks the slice added.
- **`reviews/cli-grammar-and-capability-evidence.md` is now six reviewer
  decisions, not defects.** Nothing in it is a bug to fix; each entry is a call
  only the reviewer can make. Two are one decision: three hand-rolled help parsers
  (`internal/spec/spec.go:232`, `internal/worktree/list.go:18`,
  `cmd/bench/main.go:328`) and `specArg`'s missing `--` case both resolve by
  routing `spec` through the grammar or blessing the exemption. The sharpest of
  the rest is the routing registry's `whyNested` exemption for six adopt
  subcommands, whose stated reason is factually wrong for `doctor`.
- **Spec retirement stays deferred.** `bench status` flags it. Do not retire while
  the pickup file is open — its entries cite spec line numbers.
- **Two calls still open for reviewer veto.** Neither blocks anything. The
  marker-wait conformance check grades only the *slow* deadline argument, leaving
  the fast leg unbounded by named policy; and `capability.TB` is a local interface
  rather than `testing.TB`, so `internal/gate` can import the line shape without
  linking `testing` into `dist/bench`.
- **One history defect, unrepaired, reviewer's call.** Commit `c82ba1f` is
  labelled "capture: park the gate phase-timeout headroom idea" but also contains
  the entire stories 10–11 slice (649 insertions, 11 files) — a bare `git commit`
  after `git merge --squash` took the whole index. A later full gate ran green on
  that tree, so only the history is wrong. `main` is unpushed, so a split is
  risk-free.
- **Known advisory debt.** `bench structure` reports 15 issues; `internal/gate/` is
  over its 16-file budget. Gate is green regardless.
- **Seven ideas parked, zero open learnings.** Two carry real risk and neither is
  in the pickup file: **real data races in `guards.Scan`** that fail under `-race`
  on `main` today, which the gate never runs; and `waitForPIDFile`'s hardcoded 2s
  literal deadline — the defect class story 13 fixed, at a call site the spec did
  not name.
- **Build gotcha.** A plain `go build -o dist/bench ./cmd/bench` stamps
  `version=dev` and fails two `internal/contract/surface` contracts. Hand-running
  that package needs
  `go build -ldflags "-X main.version=0.2.0" -o dist/bench ./cmd/bench`.
- **Unpushed:** `main` is well ahead of origin. Pushing is the reviewer's call.

## Next command

`/bench-what-next` in a fresh session. FT87 has no defects left, so the drain is
the real next move: seven parked ideas and six reviewer decisions are both waiting
on a reconcile, and two of the parked ideas (the `guards.Scan` races, the
`waitForPIDFile` literal) outrank anything in the pickup file.

If you would rather settle FT87 first, the alternative is to work the pickup
file's six decisions directly — they are questions to answer, not code to write,
so they belong in conversation rather than in `/bench-implement-spec`.

## Shape

Rewritten in full at every phase close, pruned rather than accreted: a fresh
session pays for every line it reads cold, so drop anything it would not act on.
Keep the three sections above — **State** (what is true now, including anything
uncommitted), **Next command** (the exact harness-native invocation, not a
description of it), and this one.

The handoff carries no date of its own. `bench status` computes its age from the
commit that last wrote this file and reports a `handoff` row once anything has
landed since. Where this document and the tree disagree, the tree wins.
