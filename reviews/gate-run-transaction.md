# gate-run-transaction review pickup

Reviewed base: `00290cf2f36a4fe6631165c66812d76caa8abb46`
Reviewed tip: `67faecb3438f09ca18ceea742caa00ee1a0cbd3a`

## Standards

Finding count: 2. Worst issue: Low.

- Low, auto-fix: `specs/gate-run-transaction/spec.md:8` carries dated implementation and reviewer-process history. Remove that narration and keep only present-tense contract rationale in the GC5 and GT5 rows, per `.bench/BENCH.md:73-75`.
- Low, auto-fix: `internal/gate/run_outcomes_test.go:19-20` carries coverage-provenance commentary. Delete it; the authoritative GC10 citation remains in `specs/gate-run-transaction/spec.md:165`, per `.agents/skills/bench-craft-comments/SKILL.md:32-38`.

## Spec

Finding count: 0. Worst issue: none.

## Coverage

Finding count: 2. Worst issue: High.

- High, auto-fix: cancellation after pending persistence is unproved. Add a public-seam test that cancels a fixture-held run and asserts exit 130, no ready/red verdict, and recoverable pending state. Relevant branch: `internal/gate/run_transaction.go:139-147`; public promise: `internal/gate/gate.go:190`; hostile lifecycle requirement: `projects/benchkit.md:159`.
- Medium, auto-fix: initial pending-persistence refusal is unproved. Make `bench-last-gate` a directory, assert operational refusal before oracle dispatch and a zero run count. Relevant branch: `internal/gate/run_transaction.go:140`; existing GT3-GT5 only cover terminal persistence failures in `specs/gate-run-transaction/spec.md:171-173`.
