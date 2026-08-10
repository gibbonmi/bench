# Declare AXI metadata for every worktree nested member

Blocked by: declare-axi-query-root-metadata.md, declare-axi-operational-root-metadata.md
Ownership fence: `cmd/bench/command_registry.go`, `cmd/bench/command_registry_test.go`, `internal/usage`
Integration surfaces: declaration type and required-field validator→declare-axi-query-root-metadata.md; worktree grammar owner→`internal/usage/worktree.go` (`worktreeCommands` and the default subshell form, already the single source behind `WorktreeUsage`); parent root `worktree` declaration→declare-axi-operational-root-metadata.md; completeness and inertness validation→validate-axi-registry.md
Contracts: the exported nested-member list crosses `internal/usage`→`cmd/bench/command_registry.go`, membership is the eight worktree surfaces — the seven entries of `worktreeCommands` plus the default subshell form `bench worktree [--refresh] [objective...]` that `WorktreeUsage` renders as its own first line — order is subshell form first then that slice's order, and an absent member is an error rather than an empty list, asserted by NW1 against the real exported list rather than a transcribed copy
Closure: NW1/worktree-subshell, NW1/worktree-list, NW1/worktree-path, NW1/worktree-exec, NW1/worktree-create, NW1/worktree-release, NW1/worktree-clean, NW1/worktree-recovery

## What to build

`internal/usage` exports the eight worktree surfaces it already owns in
`internal/usage/worktree.go` — the seven operations in `worktreeCommands` plus
the default subshell form `bench worktree [--refresh] [objective...]` that
`WorktreeUsage` prints as its own first line — and `cmd/bench` declares one AXI
record per surface. `bench worktree list` is declared AXI; the other seven are
declared explicitly non-AXI. Each record carries its grammar owner, renderer
family, empty class, detail modes, and shared routes.

The subshell form is a member rather than a flag spelling of the root: the FT173
census grades it as its own surface with its own known-state and contextual-action
row (`decisions/byte-preserving-axi-foundation/assets/ft173-command-help-inventory.md`
line 104), distinct from the `bench worktree list` row beneath it.

Exporting the subshell form as a named constant that `WorktreeUsage` renders is
what keeps one source: the usage line and the declaration cannot disagree, and a
ninth surface added to the usage owner makes the declaration red until it is
declared. Exactness of the approved AXI set is a later ticket's claim, and
whole-registry totality and byte-inertness stay with `validate-axi-registry.md`,
which is what lets this ticket land green alone.

## Acceptance

- [ ] [NW1] (covers CR5) each of the eight worktree surfaces exported by `internal/usage` has exactly one declaration record carrying its complete metadata.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| NW1/worktree-subshell | delete the `bench worktree [--refresh] [objective...]` record from the nested declaration table | `TestAXIDeclarationCoversEveryWorktreeSurface` in `cmd/bench` (this ticket authors it) | run `go test ./cmd/bench -run TestAXIDeclarationCoversEveryWorktreeSurface -timeout 120s`; expect the set-difference failure `axi declaration missing nested member "bench worktree [--refresh] [objective...]"` against the exported subshell-form constant in `internal/usage/worktree.go`; bound is the `-timeout 120s` binary deadline, and the test reads the exported member list in process without spawning a command |
| NW1/worktree-list | delete the `bench worktree list` record from the nested declaration table | `TestAXIDeclarationCoversEveryWorktreeSurface` in `cmd/bench` (this ticket authors it) | run `go test ./cmd/bench -run TestAXIDeclarationCoversEveryWorktreeSurface -timeout 120s`; expect the set-difference failure `axi declaration missing nested member "bench worktree list"` against `worktreeCommands` in `internal/usage/worktree.go`; bound is the `-timeout 120s` binary deadline, and the test reads the exported member list in process without spawning a command |
| NW1/worktree-path | delete the `bench worktree path` record from the nested declaration table | `TestAXIDeclarationCoversEveryWorktreeSurface` in `cmd/bench` (this ticket authors it) | run `go test ./cmd/bench -run TestAXIDeclarationCoversEveryWorktreeSurface -timeout 120s`; expect the set-difference failure `axi declaration missing nested member "bench worktree path"` against `worktreeCommands` in `internal/usage/worktree.go`; bound is the `-timeout 120s` binary deadline, and the test reads the exported member list in process without spawning a command |
| NW1/worktree-exec | delete the `bench worktree exec` record from the nested declaration table | `TestAXIDeclarationCoversEveryWorktreeSurface` in `cmd/bench` (this ticket authors it) | run `go test ./cmd/bench -run TestAXIDeclarationCoversEveryWorktreeSurface -timeout 120s`; expect the set-difference failure `axi declaration missing nested member "bench worktree exec"` against `worktreeCommands` in `internal/usage/worktree.go`; bound is the `-timeout 120s` binary deadline, and the test reads the exported member list in process without spawning a command |
| NW1/worktree-create | delete the `bench worktree create` record from the nested declaration table | `TestAXIDeclarationCoversEveryWorktreeSurface` in `cmd/bench` (this ticket authors it) | run `go test ./cmd/bench -run TestAXIDeclarationCoversEveryWorktreeSurface -timeout 120s`; expect the set-difference failure `axi declaration missing nested member "bench worktree create"` against `worktreeCommands` in `internal/usage/worktree.go`; bound is the `-timeout 120s` binary deadline, and the test reads the exported member list in process without spawning a command |
| NW1/worktree-release | delete the `bench worktree release` record from the nested declaration table | `TestAXIDeclarationCoversEveryWorktreeSurface` in `cmd/bench` (this ticket authors it) | run `go test ./cmd/bench -run TestAXIDeclarationCoversEveryWorktreeSurface -timeout 120s`; expect the set-difference failure `axi declaration missing nested member "bench worktree release"` against `worktreeCommands` in `internal/usage/worktree.go`; bound is the `-timeout 120s` binary deadline, and the test reads the exported member list in process without spawning a command |
| NW1/worktree-clean | delete the `bench worktree clean` record from the nested declaration table | `TestAXIDeclarationCoversEveryWorktreeSurface` in `cmd/bench` (this ticket authors it) | run `go test ./cmd/bench -run TestAXIDeclarationCoversEveryWorktreeSurface -timeout 120s`; expect the set-difference failure `axi declaration missing nested member "bench worktree clean"` against `worktreeCommands` in `internal/usage/worktree.go`; bound is the `-timeout 120s` binary deadline, and the test reads the exported member list in process without spawning a command |
| NW1/worktree-recovery | delete the `bench worktree recovery` record from the nested declaration table | `TestAXIDeclarationCoversEveryWorktreeSurface` in `cmd/bench` (this ticket authors it) | run `go test ./cmd/bench -run TestAXIDeclarationCoversEveryWorktreeSurface -timeout 120s`; expect the set-difference failure `axi declaration missing nested member "bench worktree recovery"` against `worktreeCommands` in `internal/usage/worktree.go`; bound is the `-timeout 120s` binary deadline, and the test reads the exported member list in process without spawning a command |
