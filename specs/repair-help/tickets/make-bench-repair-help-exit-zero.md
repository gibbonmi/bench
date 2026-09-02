# Make bench repair --help answer on stdout

Blocked by: none
Writes: bin/bench.sh, cmd/bench/main_test.go, CHANGELOG.md
Covers: none

## What to build

`bench repair` is a wrapper-only verb. The wrapper accepts no argument and
`--prune`, and it sends all other argument lists to the usage-error arm. So
`bench repair --help` prints the usage line on stderr and exits 2. Each
Go-implemented verb prints its usage line on stdout and exits 0 for `--help`.
The wrapper verb therefore breaks the help contract that an agent reads across
the command surface.

Give `repair_command` a help arm before its usage-error arm. The arm accepts
`--help`, `-h`, and the bare word `help`, which is the shape `gate_command`
already uses. The arm prints `usage: bench repair [--prune]` on stdout and
returns 0. The usage string stays in one place: a local usage function prints it
for the help arm and for the usage-error arm.

An unknown argument keeps its current behavior. The usage line goes to stderr,
and the wrapper exits 2.

## Acceptance

- [ ] `bench repair --help` prints `usage: bench repair [--prune]` on stdout and exits 0.
- [ ] `bench repair --bogus` prints `usage: bench repair [--prune]` on stderr and exits 2.
- [ ] `shellcheck bin/bench.sh` reports no finding.
