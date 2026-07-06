# Roadmap

The working prioritization document: every row is open work, verified against
the tree; a row leaves when the work ships (spec-retire) or a
`/bench-what-next` reconcile removes it. Raw capture never lands here — it
goes to `IDEAS.md` and enters only through a reviewed drain.

## Features, in priority order

**FT4 (MED-LOW) — harness task list in `/bench-implement-spec`.** Per-harness
adapter (Claude hook + phase line; Codex native).

**FT9 (MED-LOW) — `bench diff --full` base-relative context bundle.** Grilled
and mapped (`decisions/ft9.md`): grow `bench diff` with a `--full` flag that
appends the base-relative diff body plus a `log[N]{sha,subject}` table, and
repoint `/bench-review-implementation`'s two-call git prose at it. The generic
status+diff+log+staged framing is dropped — no repeated call site. One
uncertainty flag for the spec-writer: how the TOON-stdout conformance assertion
exempts the raw diff-body block (needs `craft-gate`). Ready to spec.

**FT11 (MED-LOW, defect) — `bench commit` cannot stage a file deletion.**
`git add :(literal)<path>` exits 128 on a removed path, so the sanctioned commit
path can't complete a spec-retire (which always deletes a spec); the workaround
today is raw `git commit`, which forfeits the block-check + gate-order
guarantees. Teach `bench commit` staging to record deletions for named paths,
with a gate row driving a deleted path through the command.

**FT10 (LOW) — doctor installs the kit repo's pre-push guard.** `bench guards`
already reports the missing guard; `bench doctor` should detect it on the kit
repo itself and offer the install (consumer repos get it via `bench link`).

**FT5 (LOW) — `bench outline`.** Marginal for this repo, real as a kit
affordance for large/polyglot linked repos. Needs its grill (languages,
on-demand vs committed, prose anchors).

**FT6 (LOW, parked pending evidence — leave parked):** `bench refs`, `bench
detect`, `bench doc`, `bench specs --retired`, doctor binary-presence row,
`conformanceFamilies`-vs-dispatch reconcile meta-check, and a per-anchor
bite-proof meta-test (canaries prove one needle per family today; graduate on
observed anchor rot). `bench symbols` is not carried; restore only if agents
demonstrably burn turns on symbol search.

**FT7 (LOW) — dashboard.** Low priority by declaration.

**FT8 (scheduled, not actionable) — Sonnet 5 mid-tier revisit.** Time-boxed to
2026-09-01 or the next frontier shift.

## Watch

- `bench worktree concurrent-acquire` contract test failed once under
  full-gate load, then passed 3/3 in isolation and on rerun — likely a timing
  flake surfaced by gate phase concurrency. Journal it if it recurs.

## Recommended sequence

1. FT4 harness task list in `/bench-implement-spec` — /bench-write-spec
2. FT9 `bench diff --full` context bundle (mapped, ready) — /bench-write-spec
3. FT11 `bench commit` deletion-staging fix — /bench-write-spec
