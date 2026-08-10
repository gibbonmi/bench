# Migrate status and dashboard empty classes

Blocked by: migrate-status-dashboard-aggregates.md
Ownership fence: `internal/status`, `internal/dashboard`
Integration surfaces: typed empty carrier→`internal/axi/empty.go` exercised by SDE1; prose-clean producer→`internal/status/status.go` exercised by SDE1; per-section empty renderer→`internal/dashboard/render.go` exercised by SDE1; not-in-repo refusal→`internal/dashboard/dashboard.go` exercised by SDE1; aggregate sibling blocker→migrate-status-dashboard-aggregates.md; registry empty declaration→declare-empty-dispositions.md; legacy carrier contraction→contract-aggregate-empty-routes.md
Contracts: the status and dashboard empty classifications cross `internal/status/status.go` and `internal/dashboard/render.go`→`internal/axi/empty.go`; the classes are prose-clean (a successful human line, exit 0), per-section prose empty inside a rendered page, and absent refusal (`toon.NotInRepo`, exit 1); order is the section order the page already renders; no class may be reached by defaulting from another, asserted by SDE1 against the real `status.render` and `dashboard.Render` producers
Closure: SDE1/status-prose-clean, SDE1/section-gate-empty, SDE1/section-signals-empty, SDE1/section-roadmap-empty, SDE1/section-ideas-empty, SDE1/section-worktrees-empty, SDE1/status-absent-refusal, SDE1/dashboard-absent-refusal, SDE1/route

## What to build

`bench status` on a clean repository keeps its exact prose success line, the dashboard keeps
one definitive prose empty per section, and both commands keep the not-in-repo refusal as a
refusal rather than an empty success. The migration routes each of those classifications
through the shared typed empty carrier so the class is declared and observable, without any
one class defaulting into another.

AE8 is an already-covered row for this family: `TestRenderClean` (`internal/status`) and
`TestRenderEmptyStates` (`internal/dashboard`) stay exactly as they are and remain the named
existing controls. This ticket adds the subject mutation the row lacks — a new
`TestStatusAndDashboardEmptyClassesReachTheTypedRoute` in `internal/dashboard` asserting each
empty was classified through the shared carrier rather than by a local literal.

Tree condition that must hold when this ticket is refreshed: `internal/axi/empty.go` exists
and declares the exported empty classification type `EmptyClass` with distinct constants for
the prose-clean class and the absent-refusal class. If that path or the symbol is absent,
stop and report rather than build — the prerequisite `axi-carriers-and-registry` build has
not landed.

## Acceptance

- [ ] [SDE1] (covers AE8) status renders its prose-clean success, the dashboard renders each per-section prose empty, and both keep the not-in-repo refusal, each through its own declared class on the shared typed empty route.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| SDE1/status-prose-clean | in `render`, return the empty string when `len(signals) == 0` instead of `bench: clean — nothing pending\n` | `TestRenderClean` (`internal/status`) | run `go test ./internal/status -run TestRenderClean -count=1 -timeout 180s`; the exact `out != "bench: clean — nothing pending\n"` comparison fails against `""`; the fixture is a `t.TempDir()` git repository with one committed file, so only local `add`/`commit`/`status` calls run |
| SDE1/section-gate-empty | in `render.go`, drop the `No gate cache` branch so an absent gate cache renders an empty section body | `TestRenderEmptyStates` (`internal/dashboard`) | run `go test ./internal/dashboard -run TestRenderEmptyStates -count=1 -timeout 180s`; the `No gate cache` containment assertion fails; `Render` is pure over `baseSnapshot()`, so no IO runs at all |
| SDE1/section-signals-empty | in `render.go`, drop the `No signals` branch for an empty signal ladder | `TestRenderEmptyStates` (`internal/dashboard`) | run `go test ./internal/dashboard -run TestRenderEmptyStates -count=1 -timeout 180s`; the `No signals` containment assertion fails; pure render over `baseSnapshot()` |
| SDE1/section-roadmap-empty | in `render.go`, render the absent-roadmap section with the parked-ideas empty text instead of `No ROADMAP.md` | `TestRenderEmptyStates` (`internal/dashboard`) | run `go test ./internal/dashboard -run TestRenderEmptyStates -count=1 -timeout 180s`; the `No ROADMAP.md` containment assertion fails against the borrowed text — the exact class-defaulting this row exists to catch; pure render over `baseSnapshot()` |
| SDE1/section-ideas-empty | in `render.go`, drop the `No parked ideas` branch for an empty ideas slice | `TestRenderEmptyStates` (`internal/dashboard`) | run `go test ./internal/dashboard -run TestRenderEmptyStates -count=1 -timeout 180s`; the `No parked ideas` containment assertion fails; pure render over `baseSnapshot()` |
| SDE1/section-worktrees-empty | in `render.go`, render `No out-of-pool` whenever the worktree section is empty, including when `WorktreesErr` is set | `TestRenderWorktreeClassifyFailureIsVisible` (`internal/dashboard`) | run `go test ./internal/dashboard -run TestRenderWorktreeClassifyFailureIsVisible -count=1 -timeout 180s`; the assertion that a classify failure renders as a visible error fails because the empty-pool text hides it; pure render over an injected snapshot |
| SDE1/status-absent-refusal | in `status.Command`, return the prose-clean line and exit 0 when `git.Root()` fails instead of `toon.NotInRepo()` and exit 1 | `TestStatusOutsideARepositoryRefuses` (`internal/status`, new) | run `go test ./internal/status -run TestStatusOutsideARepositoryRefuses -count=1 -timeout 180s`; the assertion pairing exit 1 with the `toon.NotInRepo()` line fails with exit 0 and the clean line; the test chdirs into a non-repository `t.TempDir()` and `git.Root()` returns without waiting on anything |
| SDE1/dashboard-absent-refusal | in `dashboard.Command`, return an empty page and exit 0 when `git.Root()` fails | `TestDashboardOutsideARepositoryRefuses` (`internal/dashboard`, new) | run `go test ./internal/dashboard -run TestDashboardOutsideARepositoryRefuses -count=1 -timeout 180s`; the assertion pairing exit 1 with the `toon.NotInRepo()` line fails with exit 0 and an empty page, and no `bench-dashboard.html` is written either way; the test chdirs into a non-repository `t.TempDir()` |
| SDE1/route | classify the clean and per-section empties with local string literals and never construct the shared typed empty value | `TestStatusAndDashboardEmptyClassesReachTheTypedRoute` (`internal/dashboard`, new) | run `go test ./internal/dashboard -run TestStatusAndDashboardEmptyClassesReachTheTypedRoute -count=1 -timeout 180s`; the assertion that the prose-clean, per-section, and absent-refusal cases each carried their distinct `axi.EmptyClass` constant fails with no classification observed, even though every rendered byte is unchanged; pure render plus one `t.TempDir()` chdir, no subprocess |
