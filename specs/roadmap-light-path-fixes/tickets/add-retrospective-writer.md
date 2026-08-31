# Add a sanctioned retrospective writer

Blocked by: make-bare-worktree-usage-safe.md
Writes: cmd/bench/main.go, cmd/bench/command_registry_test.go, internal/roadmap/retro.go (new), internal/roadmap/retro_test.go (new), internal/retros/retros.go, internal/retros/retros_test.go
Covers: LF12

## What to build

Add the missing retrospective writer and validate parser-owned grammar before
recording. Preserve the existing learning writer and all primary-local paths.

## Acceptance

- [ ] Malformed retrospective input changes no file.
- [ ] A successful write parses back as one eligible retrospective.
- [ ] Repeated writes preserve earlier capture.
- [ ] The destination remains ignored and primary-local.

