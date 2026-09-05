# Add the settle policy child

Blocked by: 01-forbid-the-new-parent-adapters-in-the-policy-census.md
Writes: internal/landing/settlepolicy/ (new)
Covers: LS20, LS21, LS23, LS26

## What to build

Verify the premise first. Read `Conflict`, `ConflictError`, `stageRecord`,
`unionStage`, `CaptureSide`, `resolveCaptureConflict`, `unionStages`,
`parseConflict`, and `contentConflictKind` in internal/landing/composition.go.
Read `TestComposeSettlesPhaseOwnedConflictsByRule` and
`TestComposeRefusesConflictsTheRuleTableCannotSettle` in
internal/landing/composition_test.go. Read
internal/worktree/landingpolicy/landingpolicy.go and
internal/worktree/landingpolicy/landingpolicy_test.go as the pure-policy shape.
Read `PolicyPackage` in internal/puritycensus/census.go.

Create the package `internal/landing/settlepolicy` with the import path
`github.com/gibbonmi/bench/internal/landing/settlepolicy`. Ticket
07-surface-the-settle-refusal-reason.md imports that exact path, so do not
rename it. This ticket writes no file inside `internal/landing/` itself,
because that directory already holds 16 source files against a 12-file budget.

Hold the settle decision in this package. It owns its own stage-record type,
which carries the mode, the object, the stage, and the path. It owns the
capture rule table, the mode-agreement rule, the mode-list conflict classifier,
and the settle verdict for each path. It answers a verdict of side, union,
remove, or refuse, and a refusal carries its reason.

The capture rule table answers `source` for `capture/session-handoff.md`. It
answers `union` for `capture/learnings.md` and for `capture/IDEAS.md`. It
answers `destination` for every other path under the prefix `capture/`. It
answers a refusal for a path outside the table, and `capture.md` is the prefix
boundary. The policy answers a removal for a path whose winning stage is
absent.

The four refusal reasons are `path outside the capture table`, `non-regular
mode`, `mode disagreement`, and `union content not text`. One reason wins. The
policy scans the stage records in merge-tree order and answers the first
refusal it finds. Inside one record the table check precedes the mode check.
The record scan precedes the settle scan.

The mode-list classifier answers `gitlink` for a mode list that holds `160000`.
It answers `symlink` for one that holds `120000`. It answers `mode` for two
unequal ordinary modes, and `textual` for one repeated mode. It checks the two
special modes before the ordinary modes.

The package reads no Git. It performs no text merge, because the census forbids
the `git merge-file` effect and the adapter keeps that merge. It reads neither
`internal/git` nor `internal/landing`.

Write the census wrapper internal/landing/settlepolicy/purity_census_test.go.
It runs `puritycensus.Scan` with `PolicyPackage()` and names this package's own
source through `MustHold`.

Write the two table tests in the new package. Name the first
`TestSettlePolicyAnswersTheCaptureRule`. It drives the policy with a literal
stage-record slice and observes the verdict for each path. It keeps the
existing `path-with-a-space` case.

Name the second table test `TestConflictKindClassifiesModeLists`. It drives the
classifier with a literal mode list. Neither table test builds a repository.

## Acceptance

- [ ] The settle policy answers `source` for `capture/session-handoff.md`.
- [ ] The settle policy answers `union` for `capture/learnings.md` and for `capture/IDEAS.md`.
- [ ] The settle policy answers `destination` for another path under `capture/`.
- [ ] The settle policy answers a refusal for `capture.md` with the reason `path outside the capture table`.
- [ ] The settle policy answers a removal for a path whose winning stage is absent.
- [ ] The classifier answers `gitlink`, `symlink`, `mode`, and `textual` for the four mode lists.
- [ ] The policy table tests build no repository.
- [ ] The policy census reports no import of `internal/git`, `internal/landing`, `internal/worktree`, `internal/bounds`, `os/exec`, or `syscall`, no ambient effect, and no parallel call.
- [ ] The policy census wrapper names this package's own source.
- [ ] Self-probe: check the ordinary modes before the special modes, and report the classifier test red.
