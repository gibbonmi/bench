# FT181 #1 — restart previousGreen: tree evidence (2026-08-03, delegated read)

Supports the resolved answer to map ticket #1 (prefer the live marker).

- `refs/bench/green/<branch>` is per-branch: any run reaching `Promote`
  advances it via CAS from that run's own `run.Base`
  (`internal/specbuild/assign.go:392,401`; recovery branch `assign.go:324`).
  A sibling run sharing the same start tip can therefore move the marker while
  another run is terminal; nothing in abandon touches it.
- `authorization.Bootstrap` (`internal/gate/authorization/authorization.go:35-70`)
  refuses "project-green marker conflicts with another tip" when the live
  marker exists, differs from `expected`, and differs from `tip` (lines 50-52).
  The terminal-restart branch (`assign.go:242`) passes `run.Base` as
  `previousGreen`, skipping the live-marker fallback (`assign.go:266-268`), so
  a legitimate sibling promotion is indistinguishable from tampering there. A
  fresh start reads the marker live and can never hit this refusal.
- Tradeoff the decision accepts: with the live marker as `expected`, the
  ancestor check (`authorization.go:56-59`) becomes self-referential — the
  restart no longer asserts what the run itself observed as green. The map's
  answer keeps `run.Base` recorded as evidence, not as the comparison operand.
- Coverage gap for the spec: no test pins the restart `previousGreen` value.
  The marker tests in `start_test.go:81-128` exercise only the fresh path;
  `abandon_test.go` never constructs sibling-promotion-then-restart.

Not read: `internal/gate`'s engine internals beyond `authorization.go`;
`precondition.go`'s subject derivation in full.
