# Pin the broker manifest symlink resolution

Blocked by: none
Writes: internal/brokermanifest/replacement_test.go
Covers: none

## What to build

`Write` resolves symbolic links in the broker executable path before it digests and binds the broker. `replacement_test.go` passes a pre-resolved regular file path, so no test in `internal/brokermanifest` witnesses that resolution. Add a test that passes a symbolic link as the broker path and checks the bound target path.

## Acceptance

- [ ] `Write` resolves a symbolic-link broker path and returns the resolved target path.
- [ ] `Write` records the resolved target path in the published manifest while the symbolic link stays intact.
