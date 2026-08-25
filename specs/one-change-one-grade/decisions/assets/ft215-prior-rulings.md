# FT215 prior rulings and current gate shape

Read on 2026-08-25 from the decision maps, the profile, ADR 0016, and the gate code.
This asset supports `specs/one-change-one-grade/decisions/one-change-one-grade.md`.

## Closed rulings

- Diff-scoped gating is unsound. `gate-budget`, `gate-critical-path`,
  `gate-concurrency`, and `cost-follows-project-size` #7 each record the ruling.
  The reason is the same in each map: contract and canary checks are behavior
  contracts with no file-to-test map.
- A scoped gate is legitimate only on a reviewer-declared boundary, never on a
  boundary derived from a diff (`cost-follows-project-size` #7). Wall-clock is
  never the justification. A scoped verdict is explicit evidence, never a silent skip.
- `projects/benchkit.md` pins the phase table and states that the gate has no
  per-package loop and no component partition.
- The Markdown-only reduced lane shipped on 2026-08-01 (`6f3486af`) and FT183 #1
  retired it on 2026-08-03 (`b41b4d20`) as unreachable.
- ADR 0016 (accepted 2026-08-25): green evidence keys to the prospective tree and
  the baseline runner identity. Candidate code cannot authorize its own publication.
- `spec-build-review-gate-cadence` #2 and #5: no ticket-local check runs the
  whole-project gate, and `promote` alone runs the one whole-project gate. That
  map waits on FT173 and FT130 for its ticket-local evidence machinery.

## What the tree does today

- The `test` phase runs `go test -count=1 ./...` (`internal/gate/phases.go:89`).
  Go schedules packages in parallel, so the wall is the slowest package plus the
  build. `internal/worktree` at 51–56 s sets that wall.
- `bench commit` rewrites the named Go files with `go/format` before it composes
  the tree (`internal/commit/commit.go:69-82`). No prose check runs before the gate.
  `prose-mechanics` runs only inside the `test` phase.
- A green verdict reuses in 0.57 s when the tree, the oracle key, and a 60-minute
  window match (`internal/gate/engine.go:100-127`). The evidence store is the
  common Git dir, so worktrees share it.
- The oracle key frames the identity root path and `PATH`
  (`internal/gate/subject.go:71-73`). A prospective subject also frames the
  baseline runner digest (`internal/gate/prospective.go:49-57,95`). A `bench commit`
  in a worktree keys to the worktree; `bench worktree land` keys to `main`.
- The reviewed landing composes the source onto the destination. It then rewrites
  the spec's `Status: staged` line to `Status: implemented`, removes the
  tickets-only folder, and authorizes that tree (`internal/landing/landing.go:203-222`).
- No conformance check in `internal/conformance` reads a `Status:` line or a
  tickets-only folder.
