# Resume a published stack

Blocked by: 10-refuse-a-stack-member.md
Writes: internal/usage/worktree.go, cmd/bench/worktree_leaves.go, cmd/bench/help_inventory_test.go, internal/worktree/land_stack.go (new), internal/worktree/land_stack_resume.go (new), internal/worktree/land_stack_resume_test.go (new)
Covers: LS78, LS79, LS80, LS81, LS82, LS100

## What to build

Ship `bench worktree land --resume <published-commit> --stack <file> [--source-tip <request>=<commit>]...`. The usage owner declares the grammar as `WorktreeLandStackResume`, and `worktreeCommands` lists it. The `worktree` help rows in `cmd/bench/worktree_leaves.go` gain one row: "resume incomplete post-publication stack work".

The stack flow of ticket 07 dispatches the stack form with `--resume` to the resume flow in the new file `land_stack_resume.go`. The stack flow also replaces its interim `land-handback` for a stopped post-publication step with `land-incomplete`, whose `next=` names this resume with each override of the call.

The resume flow does these steps:

1. It reads the stack file and the overrides with the stack reader of ticket 06.
2. It proves that the published commit resolves exactly and that the destination reaches it.
3. The second parent of the published commit must be the tip of exactly one member, and that member's index is the published count k.
4. It walks k first parents back from the published commit. Each merge must have the tip of the next member, in reverse file order, as its second parent, and the last first parent is the destination base of the stack.
5. For each published member, it proves the merge with `publicationFacts`, as one resume does, so the reviewed range must authenticate.
6. It finishes the post-publication steps of ticket 07 once for the published merge: the marker, the commitment reconcile, the destination checkout reconcile with its residue guard, and the prune.
7. It releases each published member that is still active, and it skips each released member that has a terminal receipt. It then runs the two landing effects once.

Each member after k stays active. A second resume of a finished stack exits 0 and writes nothing.

## Acceptance

- [ ] A stack called with one override whose release of member 2 fails exits 3, and its `next=` names `bench worktree land --resume <tip> --stack <file>` with that override (LS78).
- [ ] The stack resume refuses a published commit whose second parent is not the tip of the listed member, before any step (LS79).
- [ ] The stack resume releases the active published member and exits 0 (LS80).
- [ ] A second stack resume of the same stack exits 0 and changes no ref, no ledger record, and no checkout (LS81).
- [ ] The resume of a stack that published members 1 and 2 of four leaves members 3 and 4 active (LS82).
- [ ] `bench help` prints the stack resume grammar as one row (LS100).

The release fault follows the `unknown-ignored` case of `TestLandCommandPublicRealGitJourney` in `internal/worktree/land_journey_test.go`, which leaves residue that the release refuses.

## Checkpoint verification

- `bench test --package ./internal/worktree`
- `bench test --package ./cmd/bench`
- `bench test --check root-conformance`

## Mutation proof

Make the resume skip the second-parent proof of each walked merge. The LS79 test must fail. Spell the probe against the code that the author writes, and require `bench probe` to report `bit` with a byte-exact restore.
