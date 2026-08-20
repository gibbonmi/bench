# Wire registry-adapter selection through bench release

Blocked by: none
Writes: internal/publication/command.go, internal/publication/npm_registry.go, internal/publication/*_test.go, cmd/bench/main.go, cmd/bench/main_test.go

## What to build

`bench release submit/promote/rollback` select their registry adapter from
`--adapter npm|fixture` (default `fixture`, no environment twin) instead of
hardcoding the fixture. The npm adapter publishes with `--access public` under
the public profile and appends `--provenance` only when the new `--provenance`
flag was passed. `--path staged --adapter npm` refuses up front — before the
release lock and before any registry call — naming the unimplemented staged
capability and the `--path first` alternative. The usage line, the `cmd/bench`
inventory suffix, and the help golden all gain the two flags in this same
slice. This slice also adds the record-level ordering assertion over the state
machine's package plan (all four platform transitions before the wrapper), so
the ordering has a live owner before T2 retires the workflow step-name check.
The contract T2 consumes: the exact invocation grammar `release submit
--version <v> --profile public --path first --adapter npm --provenance
--registry <url>` exists and is exercised by tests, and the ordering assertion
already stands.

## Acceptance

- [ ] `submit --adapter npm` against a stub `npm` on PATH invokes `npm publish` with the approved tarball; without the flag no npm process spawns and the fixture path behaves exactly as before (spec rows R1, R2).
- [ ] `--adapter bogus` exits 2 with a usage line naming `--adapter` (R3).
- [ ] The recorded publish argv carries `--access public` under the public profile and `--provenance` only when flagged (R4, R5).
- [ ] `submit --path staged --adapter npm` refuses before the lock and spawns no process (R6).
- [ ] `promote --adapter npm` and `rollback --adapter npm` drive the stub npm; without the flag they drive the fixture (R7).
- [ ] A new test in `internal/publication` asserts the package plan orders all four platform packages before the wrapper (R13).
