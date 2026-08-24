# Add the path-form example to the commit help

Blocked by: none
Writes: internal/commit/commit.go, internal/commit/landing_test.go

## What to build

`bench commit --help` shows one concrete example line that demonstrates the
trailing `-- <path>...` form, so a caller does not learn the shape by tripping
the usage line.

## Acceptance

- [ ] The `--help` output contains one example invocation with `-m`, `--`, and a path.
- [ ] A grammar error still prints only the one-line usage, not the example.
