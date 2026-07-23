# Session handoff

Repository: /home/mgibs/workspace/bench — branch `main`. Written 2026-07-23 at
the close of a `/bench-what-next` drain. Everything below is executable from a
cold start; no conversation history is needed.

## State

- **The gate-trustworthiness work is done and retired.** The spec
  `trustworthy-gate-under-load` was implemented (`d3cf129`), retired
  (`f36c729`), and hardened (`dcbba47`): the gate owner record, the
  already-in-progress refusal naming the owner PID and its liveness,
  `gate-run` process-group teardown on signal, the two-leg marker wait, the
  `.bench-contract-env` cleanup reap, and self-attributing conformance diags
  all shipped. Its roadmap row is removed in this drain.
- **`specs/` is empty.** No spec is staged; no build is mid-flight.
- **Capture is drained.** `IDEAS.md` is empty and `.bench/learnings.md` holds
  no open entries — the 2026-07-23 reproduction-economics entry became roadmap
  row FT112.
- **This drain is uncommitted** unless the reviewer has approved it: it edits
  `ROADMAP.md`, `.bench/learnings.md`, and this file only. Doc-only, so it
  takes a plain `git commit`, not `bench commit`.
- **Known advisory debt:** `bench structure` reports 10 violations (crowded
  `internal/adopt/` and `internal/contract/surface/`, plus seven over-length
  files). The gate is green with them, and no roadmap row covers them.
- **Unpushed:** `main` is ahead of origin by several commits. Pushing is the
  reviewer's call.

## Next command

`/bench-write-spec` — FT87 slice 3, the command-wide parser and
security-evidence capability. It is the top row of `ROADMAP.md`'s recommended
sequence and its decision map (FT87 tickets #7 and #8) is already closed, so
the spec compiles from a reviewed source rather than the batch-drain override.
Run it in a fresh **mid-tier** session, the profile's spec default.

Alternatives in the same sequence, in order: `/bench-shape-idea` for FT91's
first arm (core-count-aware gate/phase concurrency — a grill, because the cut
line between speed and oracle authority is a reviewer decision), then
`/bench-write-spec` for FT86.
