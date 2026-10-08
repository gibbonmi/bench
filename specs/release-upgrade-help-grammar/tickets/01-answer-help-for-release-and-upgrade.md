# Answer help for release and upgrade

Blocked by: none
Writes: internal/publication/command.go, internal/publication/npm_registry_test.go, internal/adopt/upgrade.go, internal/adopt/compatibility_test.go
Covers: none

## What to build

`bench release --help` and `bench upgrade --help` print the grammar of the verb on stdout and exit 0.
Each verb prints the usage line that its package already owns, so the change adds no second copy of the grammar.

The upgrade parser uses `usage.Parse` with a `usage.Grammar`, as `bench doctor` does.
The shared parser therefore decides which arguments are help requests.
When that parser refuses the arguments, upgrade prints its current usage line on stderr and exits 2.

The release parser accepts a help spelling only when that spelling is the sole argument.
That argument is in the subcommand position.
The release parser uses the same shared parser to identify the help spelling.

All other arguments keep their current refusal and exit code.
`bench repair --help` and the help contract for all verbs are not in the scope of this ticket.

The test seam is the command entry point of each package: `publication.Command` and `adopt.Upgrade`.

## Acceptance

- [ ] `bench release --help`, `bench release -h`, and `bench release help` each print the release usage line on stdout and exit 0.
- [ ] `bench upgrade --help`, `bench upgrade -h`, `bench upgrade help`, and `bench upgrade --check --help` each print the upgrade usage line on stdout, write nothing to stderr, and exit 0.
- [ ] `bench release --bogus` still prints `usage: bench release (unknown argument: --bogus)` on stdout and exits 2.
- [ ] `bench upgrade --bogus` still prints the upgrade usage line on stderr, writes nothing to stdout, and exits 2.
- [ ] `bench upgrade --check --check` prints the upgrade usage line on stderr, writes nothing to stdout, and exits 2. A repeated flag is a usage error, as the shared parser decides for each verb that uses it.
