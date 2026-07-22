# Session handoff

Repository: /home/mgibs/workspace/bench — branch `main` at `c701642`
(`spec-retire: repo-aware-setup`, on top of `365e5bb`, the FT76 landing).

## State

- **FT76 shipped.** `bench setup` (repo-aware one-command adoption) landed
  gate-green in `365e5bb`: wrapper route, adopt-package deep module
  (inspect → preview → confirm → FT84 transaction → doctor → next action),
  gate-inference table, per-harness doctor rows, packed-artifact offline leg,
  slimmed `/bench-setup-repo`, truthful README quickstart. Spec reviewed
  (three-axis, 11 findings fixed) and retired in `c701642`; roadmap row
  removed; decision map `decisions/repo-aware-bootstrap.md` deleted.
- **Decisions that stay closed:** the spec's four defaulted decisions
  (non-TTY posture, converge-and-report re-run, gate-inference table,
  porcelain output) were built as written. Empty pre-existing `CLAUDE.md`
  stays preserved-and-red (link-lifecycle contract wins over the spec's
  "converge unconditionally" sentence) — flagged for reviewer veto, default
  is keep. Profile seeds through a `"seed"` plan kind: atomic with the
  transaction, never manifest-tracked.
- **Known flake (open):** under full gate load on this WSL2 host, two
  pre-existing checks flake — conformance's inner core `go test` (diag
  discards output) and `binary_repair_hardened_test.go`'s 2s wall-clock sync
  deadline. Ledgered 2026-07-22 in `.bench/learnings.md` with three proposed
  gate-authoring changes (reviewer's call). Retry-once at the same tier is the
  working mitigation.
- **Unpushed:** 9 commits on `main`; push is the reviewer's.

## Next command

`/bench-what-next` — drains the parked idea (worktree-recovery discard gap),
the open learnings entry, and reconciles `ROADMAP.md` (structure row shows 8
split suggestions).
