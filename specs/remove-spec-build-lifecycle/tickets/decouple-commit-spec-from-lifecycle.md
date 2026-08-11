# Decouple bench commit --spec from the lifecycle

Blocked by: none

## What to build

`bench commit --spec <slug>` flips the named spec staged→implemented on its
green landing commit without consulting any spec-build run state: the
`specbuild.New(...).Status` consult and the active-build refusal in
`internal/commit` are deleted, and `bench commit --help` states that `--spec`
marks the named spec implemented. Lands green while the lifecycle still
exists, so the later deletion cannot leave a dangling consult. Covers RM6.

## Acceptance

- [ ] `internal/commit` has no specbuild import; the active-build refusal and
      its tests are gone, replaced by tests of the direct staged→implemented
      flip.
- [ ] `bench commit --help` names the `--spec` semantics (closes the open
      learnings entry's help gap).
- [ ] `go test ./internal/commit/...` green.
