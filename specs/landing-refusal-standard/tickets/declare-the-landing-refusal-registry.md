# Declare the landing refusal registry

Blocked by: none
Writes: internal/worktree/land_refusal.go, internal/worktree/land.go, internal/worktree/identity_component.go, internal/worktree/identity_component_test.go, internal/worktree/land_surface_test.go, internal/worktree/land_reauthorization_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: LRS1, LRS2, LRS3, LRS20, LRS21

## What to build

A declared registry holds the landing's refusal faces. Each entry names one
face and one route builder. The registry follows the shape that
`identityComponents` sets in `internal/worktree/identity_component.go`. One
constructor turns an entry into a `refusalError`. The constructor takes the
route as a required argument, so no call site can omit the argument.

`TestIdentityComponentRegistryHasAProducingFixture` extends to the new
registry. The walk drives one producing fixture for each declared face. The
same walk asserts that each fixture prints a `next=` field that is not empty.
An empty route string still compiles, so this walk is the assertion that
catches it.

Each landing-preflight route ends with the caller's own re-run. The
face-specific repair comes first. The re-run repeats the `--request`, `--base`,
`--source-tip`, and `--spec` values the caller passed. A source path that is
not line-safe takes the existing `atSourceWorktree` pointer form. A refusal
that fires before the assignment resolves names the operator's own worktree
path.

Later tickets add their own face to this registry. The constructor is the
contract those tickets share, so no later ticket composes a route of its own.

## Acceptance

- [ ] LRS1 — the registry walk drives one producing fixture for each declared
      landing refusal face.
- [ ] LRS2 — the registry walk asserts that each face's fixture prints a
      `next=` field that is not empty.
- [ ] LRS3 — each landing-preflight refusal ends `next=` with the caller's own
      `bench worktree land` re-run.
- [ ] LRS20 — a refusal before the assignment resolves names the operator's
      worktree path in its re-run.
- [ ] LRS21 — a source path that is not line-safe produces `next=` in the
      `bench worktree exec` pointer form.
- [ ] A face declared with no producing fixture reds the walk.

## Delegate charge

You work in the Bench repo on the `landing-refusal-standard` spec. Line: opus /
medium. Effort: medium, at most 3 iterations.

Read `specs/landing-refusal-standard/spec.md` first. Read
`internal/worktree/land_refusal.go`, `internal/worktree/classifier.go`, and
`internal/worktree/identity_component.go` in full. Read the preflight faces in
`internal/worktree/land.go`. Read `identityComponentFixtures` and
`TestIdentityComponentRegistryHasAProducingFixture` in
`internal/worktree/identity_component_test.go`.

Coverage rows: LRS1, LRS2, LRS3, LRS20, LRS21. Show each row red before your
edit. Show each row green after. Return the red-to-green log per row.

Mirror `identityComponents` for the declared data. Keep the `refusal` struct's
optional `next` field, because the verbs outside this spec still use it. The
five `cmd/bench` and `internal/conformance` entries are the registry closure
for the `internal/worktree` package. Edit them only if your change reaches
them.

Self-probe with a behavioral mutation. A deleted route argument does not
compile, and a probe that fails to compile proves nothing. Point one face's
route builder at a wrong command, such as `bench worktree release`, or pass
that face an empty route string. Show LRS2's registry-walk assertion red. Show
LRS3's re-run-tail assertion red. Report both observed results, and add the
missing assertion if either row stays green.

Run `bench worktree exec "<label>" -- go test -parallel 2 ./internal/worktree/`.
The exec form is the only command form. Do not use `cd`. Do not commit. Do not
edit the spec.
