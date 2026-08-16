# Distinguish missing Git at dispatch

Blocked by: none
Writes: internal/skillsindex/command.go, internal/skillsindex/command_test.go, cmd/bench/command_registry_test.go

## What to build

At the real `cmd/bench` registry dispatch, classify the underlying error already
returned by `git.Root`: executable-not-found produces exit 1 and names required tool
`git` as missing or non-executable, while an executed Git probe outside a work tree
retains the established not-in-repository line. Do not edit `internal/git` and do not
couple this command-only tracer to producer parsing or replacement cleanup.

The paired environments are one outcome partition; splitting them would allow a
missing-tool special case to regress the outside-repository contract while its ticket
remained green. This ticket is independently shippable from the first frontier; its
write overlap with the lifecycle ticket prevents parallel assignment but creates no
behavioral blocker.

## Acceptance

- [ ] `(covers HI11)` Empty `PATH` crosses real registry dispatch, exits 1, names
  required tool `git`, and never prints “not in a git repository.”
- [ ] `(covers HI11)` Git present outside a repository crosses the same dispatch,
  exits 1, and preserves the established diagnostic byte-for-byte.
