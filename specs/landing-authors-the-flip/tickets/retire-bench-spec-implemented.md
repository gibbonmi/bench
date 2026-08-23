# Retire bench spec implemented

Blocked by: retire-the-commit-route-flip-and-close.md
Writes: internal/spec/spec.go, internal/spec/spec_test.go, cmd/bench/main.go, cmd/bench/main_test.go, cmd/bench/command_registry_test.go

## What to build

`bench spec implemented` leaves the executable. The subcommand switch, the
working-tree flip function, and the shared locate helper leave; the helper's
last caller is gone. The `bench spec` usage text, the `bench help` row, and
the registry help case drop it. The derived-implemented-bytes function stays as the one
source of the flip for the landing verb.

`bench spec retire` and `bench spec history` keep their grammar and behavior. Delete the flip tests with the
behavior. The blocker is a same-file edit of `internal/spec/spec.go` and the
help golden, not a behavior dependency.

## Acceptance

- [ ] FA3: `bench spec implemented x` exits 2 as an unknown subcommand and leaves `specs/x/spec.md` byte-identical.
- [ ] FA9: the `bench help` spec rows show `retire` and `history` only.
- [ ] FA4: `bench spec retire` on a merged-implemented spec still deletes the folder and exits 0.
