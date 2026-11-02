# Name command substitution in the follow-on refusal

Blocked by: 01-mark-live-substitution.md, 04-wire-rg-guard-hook.md
Writes: internal/benchguard, CHANGELOG.md, .bench/commitment.json
Covers: R51, R52, R53, R54

## What to build

When the Bench follow-on guard refuses at an unquoted `$(`, its message names command substitution in place of the opaque `operator=(`. The guard refuses exactly the commands that it refuses now.

Add `SideSubstitution` after `SideUnlexed` in the `Side` enum (`internal/benchguard/benchguard.go:53-62`), so every existing value keeps its number.
`operatorFor` (`:245-265`) names the operator token of a refusal. Let it also report whether that token has the `Substitution` mark from `01-mark-live-substitution.md`. `refusal` (`:267-273`) selects `SideSubstitution` only when the named operator is a marked `(` ControlOperator. Every other refusal keeps the side and the message that it has now, including a refusal whose lead is a shaping head or an unmarked `(`.
`repairSentence` (`:78-93`) gets an explicit case for `SideSubstitution` whose sentence names command substitution. The operator cell prints `$(` in place of `(`. The segment cell keeps the projected words as they are now, so `bench gate $(x)` prints `segment=bench gate $`.

Put the new table in `internal/benchguard/substitution_test.go` as `TestClassifyNamesCommandSubstitution`, in `package benchguard` like `unlexed_test.go`. Reuse `blockedPrefixWant` and `notBenchResolver` from `benchguard_test.go`, and keep a local expectation constant for the new sentence. The expectation is independent, so record one omission red under the `AGENTS.md` exception: return the plain operator side for a marked opener and record the R51 red, the restore, and the green.

This ticket is the last ticket of the build. Add one `CHANGELOG.md` entry for the whole spec: the rg guard in both harnesses, the live-substitution mark, the follow-on message, and the retired `AGENTS.md` bullet.

`benchguard.go` must stay within its 400-line budget. It had 363 lines before `02-classify-rg-calls.md` shortened `isBench`; if the change does not fit, move a helper to a new file in the package.

## Acceptance

- [ ] `Classify("bench gate $(x)")` refuses, and its message names `command substitution` and ends `operator=$(` (R51).
- [ ] `Classify("echo $(bench status)")` refuses, and its message names `command substitution` (R52).
- [ ] `Classify("bench learning t --what \"$(cat f)\"")` stays unblocked (R53).
- [ ] `Classify("true; (bench status)")` keeps the message that ends `operator=(` (R54).
- [ ] `bench test --package ./internal/benchguard` passes with every existing message table unchanged.
- [ ] `bench test --check system --run 'GuardProjection|BenchFollowOn|GuardRefusal'` passes with the golden byte-equal.
- [ ] `CHANGELOG.md` holds one entry under `## [Unreleased]` for the whole spec, and the entry names no `specs/`, `decisions/`, or `tests/` path.
- [ ] `.bench/commitment.json` is in this ticket's fence for the binding that the spec-stage close commits.
