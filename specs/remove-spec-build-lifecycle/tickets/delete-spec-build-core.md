# Delete the spec-build lifecycle core

Blocked by: decouple-commit-spec-from-lifecycle.md

## What to build

`bench spec build <anything>` answers with the CLI's standard
unknown-subcommand structured error. Delete `internal/specbuild` whole,
`cmd/bench/specbuild.go`, the specbuild pins in `cmd/bench/wiring_pins.go`,
and the `spec build` grammar and help block in `bin/bench.sh`
(`spec implemented|retire|history` survive). Retire the specbuild rows in the
conformance registries (injected ports, ordinary-build census, fixture
registries), the example-agreement check with the `ParseTicket` parser it
exercises, the 33 spec-build anchor rows in
`internal/anchors/registry_data.go` with their
`tests/canary/workflow-guidance-anchors/spec-build-*` fixtures, and the
covers-annotation comment in `internal/coverage`. Add the enumerated
keep-list check: every kept verb (`worktree create/path/exec/release/clean`,
`gate`, `commit`, `status`, `guards`, `idea`, `roadmap`,
`spec implemented/retire/history`) answers `--help` exit 0. Covers RM1, RM2,
RM7, RM11.

## Acceptance

- [ ] `bench spec build start x` and `bench spec build --help` return the
      standard unknown-subcommand structured error.
- [ ] `rg -l specbuild` over the tree is empty; `go test ./...` green.
- [ ] The keep-list check exists, enumerates the kept verbs, and its red was
      demonstrated during authoring by temporarily dropping one kept route.
- [ ] No surviving conformance check or anchor's predicate is weakened —
      rows whose subject is deleted are removed, nothing else edited.
