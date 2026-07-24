# Session handoff

Repository: /home/mgibs/workspace/bench — branch `main`. Everything below is
executable from a cold start; no conversation history is needed.

## State

- **FT87 slice 3 is built, reviewed, and green.**
  `specs/cli-grammar-and-capability-evidence.md` is `Status: implemented`; all 13
  stories are in, and `bench coverage --check` is green on 30 rows (27 original
  plus three the review round added). The gate is green at HEAD.
- **Semantic review ran and its round-1 fixes landed** (`892eea0`). Closed: the
  `help`-in-positional regression that made `bench idea help <text>` print usage
  and park nothing; an empty positional resolving to the cwd and staging every
  changed file beneath it; three usage lines duplicated as literals beside their
  grammar's `Help` (already drifted — the commit error line had lost `[--]`); a
  routing conformance check a doc comment could satisfy; and `t.SkipNow()`
  walking through the skip-ownership guard.
- **The spec was amended to match the tree** (`9edfa2a`). Two decisions the build
  made against the spec are now spec text, not drift: capability evidence travels
  by `BENCH_SKIP_LOG` side-channel file rather than a tee of phase stdout (`go
  test` without `-v` discards a passing package's stdout), and a repeated flag is
  a usage error — kept rather than reverted, because no flag in this CLI
  accumulates a list, so last-one-wins only hides a mistyped invocation. The two
  rules round 1 introduced (sole-argument bare `help`, empty positional) are
  stated with coverage rows.
- **`reviews/cli-grammar-and-capability-evidence.md` is live pickup work: 16
  findings across three axes.** It is tracked, not drift. The three with real
  teeth:
  - `internal/gate/capability_skips.go:67-71` — `readSkipTally` swallows
    `os.ReadFile`'s error and returns an empty tally, so under
    `BENCH_REQUIRE_CAPABILITIES=1` an unreadable log reads exactly like a fully
    capable runner. That is silent de-enforcement in the release workflows.
  - `internal/gate/gate.go:150` — `gateEnv` strips `BENCH_SKIP_LOG` so a canary's
    inner run cannot contaminate the outer tally, but no test asserts it. Delete
    the clause and the suite stays green; the live cost is a canary fixture's
    skip turning a real release red.
  - `internal/conformance/subcommand_routing_test.go:78-89` — six adopt
    subcommands are exempted under a `whyNested` reason that is factually wrong
    (`internal/adopt/doctor.go:181-189` is a flat hand-rolled switch, and
    `bench doctor -h` exits 2). The exemption reason is free text nothing grades.
- **Spec retirement is deliberately deferred.** `bench status` will flag the spec
  for retirement now that it is `implemented`. Do not retire it while the pickup
  file is open — those findings cite spec line numbers, and retiring deletes the
  file they cite.
- **Two calls still open for reviewer veto.** Neither blocks anything.
  - The marker-wait conformance check grades only the *slow* deadline argument of
    package-qualified `WaitForTwoLegMarkers` calls; the fast leg is bounded by no
    named policy.
  - `capability.Capability`/`Environment` take a local `capability.TB` interface
    rather than `testing.TB`, so `internal/gate` can import the line shape
    without linking `testing` into `dist/bench`.
- **One history defect, unrepaired, reviewer's call.** Commit `c82ba1f` is
  labelled "capture: park the gate phase-timeout headroom idea" but also contains
  the entire stories 10–11 slice (649 insertions, 11 files). Cause: a bare `git
  commit` after `git merge --squash` commits the whole index. A later full gate
  ran green on that exact tree, so the content is verified; only the history is
  wrong. `main` is unpushed, so a split is risk-free.
- **Known advisory debt.** `bench structure` reports 15 issues; `internal/gate/`
  is over its 16-file budget from the collector's new files. Gate is green.
- **Seven ideas parked, zero open learnings.** Two carry real risk and neither is
  in the pickup file: **real data races in `guards.Scan`** that fail under
  `-race` on `main` today, which the gate never runs; and `waitForPIDFile`'s
  hardcoded 2s literal deadline — the same defect class story 13 fixed, at a call
  site the spec did not name.
- **Build gotcha.** A plain `go build -o dist/bench ./cmd/bench` stamps
  `version=dev` and fails two `internal/contract/surface` contracts. Hand-running
  that package needs
  `go build -ldflags "-X main.version=0.2.0" -o dist/bench ./cmd/bench`.
- **Unpushed:** `main` is well ahead of origin. Pushing is the reviewer's call.

## Next command

`/bench-what-next` in a fresh session — seven parked ideas and a 16-finding
pickup file are both waiting on a reconcile, and the drain is what turns them
into sequenced roadmap rows.

If you would rather keep closing FT87 instead, the alternative is
`/bench-implement-spec` against `reviews/cli-grammar-and-capability-evidence.md`,
starting with the two `internal/gate` findings above — they are the only ones
that change what the oracle reports, and both are small.

## Shape

Rewritten in full at every phase close, pruned rather than accreted: a fresh
session pays for every line it reads cold, so drop anything it would not act on.
Keep the three sections above — **State** (what is true now, including anything
uncommitted), **Next command** (the exact harness-native invocation, not a
description of it), and this one.

The handoff carries no date of its own. `bench status` computes its age from the
commit that last wrote this file and reports a `handoff` row once anything has
landed since. Where this document and the tree disagree, the tree wins.
