# Prove the ticket-grammar canary family bites

Blocked by: register-ticket-grammar-sweep.md
Writes: tests/canary/ticket-grammar (new), internal/conformance/registry/registry.go
Covers: TG29

## What to build

A canary family named `ticket-grammar` proves each new diagnostic class through
a real gate run. Each fixture plants a synthetic staged spec under its `files/`
overlay. Each fixture carries the one `EXPECT` line the check emits for it.

The family binding and the fixtures land together in this ticket. The meta
check reds a bound family with no fixture directory, so the binding cannot land
before the fixtures. The generic proof over the family then observes each red
and each restored green.

Nine fixtures cover these diagnostic classes:

- a missing required field
- a duplicate field line
- an unterminated fence
- a dangling blocker
- a blocker cycle
- a self-edge
- a phantom `Covers:` token
- an uncovered declared row
- a binding row that names an absent registry file

## Acceptance

- [ ] TG29 — every fixture in the family reds a real gate run with its `EXPECT` line.
- [ ] Each fixture restores its subject and the red disappears.
- [ ] The family has exactly one registered check owner.
- [ ] `bench canary` validates the family with no diagnostic.

## Delegate charge

You work in the Bench repo on the `ticket-grammar` spec. Read
`specs/ticket-grammar/spec.md` first. Then read
`internal/conformance/fixture_bite_test.go` and
`internal/conformance/registry/registry.go` in full. Read
`tests/canary/prose-mechanics/long-sentence/` for the overlay shape. Read
`internal/conformance/ticket_grammar_test.go`, which a sibling ticket landed.

Create nine fixture directories under `tests/canary/ticket-grammar/`. Give each
one a `files/` overlay that plants a synthetic staged spec and its tickets
directory. Encode a leading dot in a directory name the way the existing
fixtures encode it.

Cover one diagnostic class per fixture. The classes are a missing required
field, a duplicate field line, an unterminated fence, a dangling blocker, and a
blocker cycle. The remaining classes are a self-edge, a phantom `Covers:`
token, an uncovered declared row, and a binding row that names an absent
registry file.

Write one `EXPECT` file per fixture. Copy the exact diagnostic text the check
emits. Do not paraphrase it.

Add the `ticket-grammar` family row to `familyChecks` in
`internal/conformance/registry/registry.go`. Bind it to the `ticket-grammar`
check. Add nothing else to that file.

Run `bench worktree exec ft174-ticket-grammar -- go test ./internal/conformance/...`.
Then run `bench worktree exec ft174-ticket-grammar -- go test ./internal/canary/`.
Do not commit. Do not edit the spec.
