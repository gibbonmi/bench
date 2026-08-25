# Print the spec subcommand inventory

Blocked by: none
Writes: internal/spec/spec.go, internal/spec/spec_test.go

## What to build

`bench spec --help` prints the public subcommand inventory instead of an unknown-argument error.

## Acceptance

- [ ] `bench spec --help` exits 0 and names `retire` and `history`.
- [ ] `bench spec -h` prints the same inventory and exits 0.
- [ ] An unknown subcommand still prints an unknown-argument error and exits 2.
