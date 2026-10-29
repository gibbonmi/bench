# Package listing drops the unread EmbedPatterns field

Blocked by: none
Writes: internal/testreport/selection.go
Covers: none

## What to build

The `FT463` build (`f9bdbc5c3`) retired the check-floor embed flag.
After that build, no code reads the `EmbedPatterns` field of `listedPackage` in `internal/testreport/selection.go`.
The only writer is the `go list -json` decode, and the decoder skips a JSON key that the struct does not declare.
The `byImport` comment still says that the package "holds its embed patterns".

Delete the `EmbedPatterns` field.
Remove the stale clause from the `byImport` comment.
Keep the behavior identical.

## Acceptance

- [ ] A repository-wide search for `EmbedPatterns` finds no match.
- [ ] The `byImport` comment states only that the map takes an import path to its package.
- [ ] The `internal/testreport` tests pass.
- [ ] The `bench commit` lane passes.
