# Refuse production code that re-implements a standard-library function

Blocked by: none
Writes: internal/conformance/registry/checks.go, internal/conformance/checks_test.go, internal/conformance/stdlib_reimplementation_test.go (new), projects/benchkit.md, CHANGELOG.md, .bench/commitment.json, internal/adopt/doctor.go, internal/adopt/link.go, internal/adopt/link_hook.go, internal/assessment/collection.go, internal/benchguard/benchguard.go, internal/canary/mutation.go, internal/capability/capability.go, internal/compatibility/capabilities.go, internal/consumers/resolve.go, internal/coverage/citations.go, internal/gate/manifest.go, internal/gate/runner.go, internal/git/staged.go, internal/gitguard/verdict.go, internal/gitguard/verdict_push.go, internal/harnesses/harnesses.go, internal/landing/landing.go, internal/lines/lines.go, internal/maps/schema.go, internal/outline/outline.go, internal/preflight/preflighttest/reviewfiles.go, internal/publication/fixture_registry.go, internal/publication/npm_registry.go, internal/releaseevidence/release_requirements.go, internal/repairpilot/record.go, internal/retros/recommendations.go, internal/roadmap/roadmap.go, internal/sessioninspect/sessioninspect.go, internal/structure/budgets.go, internal/testreport/environment.go, internal/testreport/named_check.go, internal/worktree/clean.go, internal/worktree/land_refusal.go
Covers: none

## What to build

Production code still re-implements standard functions. FT373 keeps four categories of the `modernize` analyzers in `golang.org/x/tools`, which the module already lists as a direct dependency:

- `slicescontains`: `slices.Contains` and `slices.ContainsFunc`
- `mapsloop`: a `maps.Copy` loop
- `minmax`: a user-defined `min` or `max` function
- `stringscut` and `stringscutprefix`: `strings.Cut`, `strings.CutPrefix`, and `strings.CutSuffix`

On `main` at `216bd7c5`, `modernize -test=false -slicescontains -mapsloop -minmax -stringscut -stringscutprefix ./...` reports 46 sites in the 33 production files that the `Writes:` line names. First, apply `modernize -fix` to those categories and make sure that each rewrite keeps its behavior. Then add one Go-source conformance check that runs the same analyzers and refuses each production site. Register the check in the conformance registry, and add its row to the profile's input table.

The check leaves out `rangeint`, `SplitSeq`, `FieldsSeq`, `TypeFor`, `WaitGroup.Go`, and the efficiency hints, because they are idiom, not duplicated knowledge.

## Acceptance

- [ ] FT373.check: the new conformance check refuses a production file that re-implements one function of each kept category. A recorded red shows each refusal, and the check ignores test files and the excluded categories.
- [ ] FT373.sites: no production site in the kept categories remains, and the whole-project gate is green with the check on.
- [ ] The conformance registry and the profile's input table name the check once each.
- [ ] `CHANGELOG.md` has one entry for the new check.
