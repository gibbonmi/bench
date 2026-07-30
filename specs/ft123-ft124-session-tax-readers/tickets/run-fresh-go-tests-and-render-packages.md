# Run fresh Go tests and render packages

Blocked by: none

## What to build

Sessions can invoke a fresh real Go test run from any repository depth through
the built `bench test` command and receive one stable package-status row for every
observed package.

## Acceptance

- [x] Zero package arguments run `go test -json -count=1 ./...`; one package
  replaces only `./...`; both run at the Git root from root or deep cwd.
- [x] Zero, one, and excess package arguments, empty input, help, unknown flags,
  and a dash-leading package after `--` follow the shared grammar while paths and
  names containing spaces or globs remain one argv.
- [x] Every terminal package appears once in sorted order with status `pass`,
  `fail`, or `no-tests`.
