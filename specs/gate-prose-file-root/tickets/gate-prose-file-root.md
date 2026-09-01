# Refuse a file as the `bench gate-prose` root operand with a usage message

Blocked by: none
Writes: internal/gate/gate_prose.go, the internal/gate test file that covers parseGateProseArgs, tests/canary fixtures only if the command help changes

## What to build

Today a file as the root operand, as in `bench gate-prose README.md`, exits 1
with an unreadable exclusion-file error. The root operand must be a directory.
When the operand names an existing file, the command refuses with a usage
message on stderr. The message names the operand and says the root must be a
directory. It also points at the file form the usage line already documents.

The exit code is the usage exit code for a malformed argument. A missing path keeps
its current diagnostic.

The refusal happens before the exclusion file is read, so the exclusion-file
error never appears for a file root. Do not change the lane wrappers or the
verdict shape. Do not change the named-result rendering. Do not change the usage
line unless it is wrong. If the usage line changes, update the help canary owner
under `tests/canary` in the same diff.

## Acceptance

- [ ] `bench gate-prose <existing-file>` exits with the usage exit code and a stderr line that names the operand and says the root must be a directory.
- [ ] The stderr text for a file root does not contain `prose-exclusions`.
- [ ] `bench gate-prose <directory>` behavior is unchanged: an existing test for a directory root stays green.
- [ ] A test in `internal/gate` goes red when the file-root guard is removed.
- [ ] `go vet` and `gofmt` are clean on the edited files.
