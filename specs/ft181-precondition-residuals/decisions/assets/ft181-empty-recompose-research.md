# FT181 #4 — empty-candidate recomposition: shipped-state evidence (2026-08-03, delegated read)

Supports the corrected answer to map ticket #4 (the residual wedge is the
blanket recompose refusal of checkpoint and start; promote already
fast-forwards an empty run).

- The promote-path fix landed 2026-08-02 as `2874d94` ("fix: rebase empty
  spec-build recomposition"), before the roadmap drain (`17af7e4`) that still
  describes it as open. `"No valid patches"` no longer appears in any `.go`
  source — only in capture and decision prose.
- Pinned green: `TestPromoteRecomposesEmptyRunOnWorkingAdvance`
  (`promotion_recompose_test.go:37`) — zero checkpoints, tip advanced,
  promote succeeds and moves base, candidate tip, and durable ref; the
  bootstrap-failure variant (line 62) proves the path stays fail-closed.
- The residual wedge, confirmed in code: `preconditions` refuses with
  `errRecompose` for every op except abandon whenever
  `run.Base != subject.tip` (`precondition.go:82,106`), regardless of whether
  any checkpoint exists — so checkpoint and start still refuse on an empty
  run whose tip moved; the fast-forward lives only inside
  `recomposePromotion`.
- Shipping vehicle note: the fix's ticket lives in the still-open
  `specs/light-path-composed-green-recomposition/` spec
  (`tickets/rebase-empty-recomposition.md`, checkboxes still unchecked
  despite the landed code and tests), and no spec-retire commit exists for
  that folder — bookkeeping for the spec's own close, not for this map.

Not read: `promotion_recompose_test.go` beyond line ~101; `checkpoint.go`
and `assign.go` bodies beyond the precondition call sites; no fresh
`go test` run — test presence and diffs taken as evidence.
