# Create the shared refusal-route registry

Blocked by: none
Writes: internal/refusalroute/registry.go (new), internal/refusalroute/route.go (new), internal/refusalroute/registry_test.go (new), internal/refusalroute/route_test.go (new)
Covers: RR04, RR05, RR06, RR07, RR13

## What to build

Create the leaf package `internal/refusalroute`.
This package is the seam that every later ticket consumes, so this ticket is its own small review chunk.
The package imports no write-verb package.

The package holds these parts:

- the face type with its five facts: the verb, the name, the sentence, the authority, and the route
- the step types: a command step, which is a template over named facts, and an instruction step, which is one imperative sentence
- the ordered face inventory, which stays empty in this ticket
- the one constructor, which takes the face name and the facts, and returns a typed refusal
- the typed refusal, which carries the face, the sentence, the paths, and the rendered route
- the route renderer
- the exported `next` field label

The renderer joins the steps with `; then ` in declared order.
A reviewer route starts with `reviewer: `.
A tree-scoped Bench step at a worktree renders `--in <label>` as the first argument after the verb.
A slot that the operator fills renders as its placeholder.
The renderer shell-quotes a line-safe value, and it prints the slot placeholder for a value that is not line-safe.

An unknown face name returns the refusal `refusal face <name> is unregistered` with a reviewer route.
It does not panic.

The uniqueness walk refuses two faces with the same name.
Its test runs the walk over the real inventory and over an injected list with a duplicate name.
The renderer tests use injected faces, because the inventory is empty until ticket 02.

## Acceptance

- [ ] An injected list with two faces of the same name makes the uniqueness walk fail, and the real inventory passes it.
- [ ] A rendered reviewer route starts with `reviewer: `.
- [ ] A rendered agent route with three steps joins them with `; then ` in declared order.
- [ ] A command step that runs a tree-scoped Bench verb at a worktree renders `--in <label>` after the verb.
- [ ] The constructor for an unknown name returns `refusal face <name> is unregistered` with a route that starts with `reviewer: `.
- [ ] `internal/refusalroute` imports no package under `internal/worktree`, `internal/landing`, `internal/commit`, `internal/gate`, or `internal/commitment`.
