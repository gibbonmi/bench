# Repeat the package flag in one focused run

Blocked by: none
Writes: internal/testreport/command.go, internal/testreport/selection_facts.go, internal/testreport/selection_facts_test.go, internal/testreport/check_test.go, internal/testreport/testreport_test.go, internal/testreport/checks_face_test.go, internal/preflight/check_floor.go, internal/preflight/check_floor_test.go, cmd/bench/help_inventory_test.go, cmd/bench/command_registry_test.go, CHANGELOG.md
Covers: FTPKG1, FTPKG2, FTPKG3, FTPKG4

## What to build

`bench test --package <expr>` takes one Go package expression. A caller that must run several consumer packages runs one call for each package, or uses `--changed`. Two build authors asked for one call that runs several packages.

Make `--package` repeatable through the shared parser's repeatable flag. `bench test --package ./internal/a --package ./internal/b` passes both expressions to one `go test` child, in argv order, and prints one TOON report. Each value keeps the current package-pattern rule. One `--package`, the legacy positional package, and `--changed` keep their current behavior. A repeated `--package` together with `--changed`, a legacy positional, or `--check` stays a usage refusal.

The `bench test` usage line shows the repeated flag, and the help note states the repeat as the remedy for several packages. The `bench probe` grammar keeps one `--package`, so its help note keeps the one-call-per-package remedy. Both notes share one sentence for the one-expression rule.

The completion-plan check floor reads each package expression of a prepared selection. So a plan command with a repeated `--package` satisfies each element that one of its expressions covers.

## Acceptance

- [ ] FTPKG1: `bench test --package <a> --package <b>` runs both packages in one focused run and prints one report that names both packages.
- [ ] FTPKG2: a repeated `--package` with `--changed`, a legacy positional, or `--check` exits 2 with the usage line and starts no Go child.
- [ ] FTPKG3: the `bench test` usage line shows `[--package <expr>]...`, and the help note names the repeated flag. The `bench probe` help note is unchanged.
- [ ] FTPKG4: a plan command with a repeated `--package` satisfies the floor element of each named package.
- [ ] `CHANGELOG.md` has one entry for the repeated flag.
