# Land a green stack with bench worktree land --stack

Blocked by: 03-admit-a-stack-publication.md, 04-compose-and-publish-a-stack.md, 06-read-a-stack-file.md
Writes: internal/usage/worktree.go, internal/worktree/land.go, internal/worktree/joins.go, internal/worktree/land_stack.go (new), internal/worktree/land_stack_rerun.go (new), internal/worktree/land_stack_test.go (new), internal/worktree/land_stack_rerun_test.go (new), internal/worktree/land_fixtures_test.go, internal/worktree/land_flags_test.go, internal/worktree/land_refusal.go, internal/worktree/land_effects.go
Covers: LS40, LS41, LS42, LS43, LS44, LS45, LS50, LS51, LS52, LS53, LS63, LS86, LS87, LS94, LS95, LS101, LS104, LS105

## What to build

Ship `bench worktree land --stack <file> [--source-tip <request>=<commit>]...` for a green stack.

The usage owner declares the grammar as `WorktreeLandStack` beside `WorktreeLand`, and the stack grammar's usage line reads it. Ticket 11 lists the form in `worktreeCommands` and in `bench help`, together with the resume form, so this ticket changes no `cmd/bench` file.

### The dispatch

`landAttributed` dispatches to the stack flow when `--stack` is a flag, and it checks this before `--resume`. `usage.FlagPresent` matches only the flags that its grammar declares (`internal/usage/parse.go:230`). So `land_stack.go` declares one dispatch grammar that holds every flag of `landGrammar` and `resumeLandGrammar`, `--stack`, and the repeatable `--source-tip` of the stack form. Then a `-m` value of `--stack` stays a single landing. `internal/worktree/land.go` has 3 lines of headroom, so the dispatch takes at most three lines there. The flow runs inside the landing span of `landWith`, and the span names the assignment of the first member.

### The flow

1. It calls the stack reader of ticket 06 with the override values. A reader refusal exits 2 before any ledger read.
2. It resolves each member path the way `landAttributed` does. Two members that resolve to one canonical path exit 2 and name both indexes.
3. It runs the destination proofs once. It runs the assignment and source proofs of one landing for each member, over the member's effective tip. One preflight prints every refusal of these groups, and each refusal names its member index and label.
4. It takes the one lease with `leaseLanding` for the first member.
5. It prints the `landing source{...}` line of each member on stderr. It then calls `LandStack` of ticket 04 through a new `landStack` seam in `joins`. The stack admission implements `StackAdmission`: `CheckMember` calls `AdmitPublication`, and `PublishStack` calls `PublishStackAdmitted` of ticket 03.
6. For the published members, it runs the green marker advance, the commitment reconcile, the destination checkout reconcile, and the branch prune once. It reads the census count of each member and releases its assignment, in file order. It runs the two landing effects once.
7. It prints the `members[N]{member,label,source_tip,published_commit,outcome,census}` table after the lease and the `stack{destination_base,published_commit,landed,worktree}` record last.

For step 6, split the effects run from the `landed{...}` record in `land_effects.go` with a new helper. Keep the signature of `landedAfterEffects`, so its two callers in `internal/worktree/land_resume.go` stay as they are.

The flow handles each result of `LandStack`: k published members and an optional refusal. Until ticket 09, it refuses a red member or an unattributed stop with `land-handback`. Until ticket 11, it ends a stack whose post-publication step stopped with `land-handback`. Until ticket 10, the stack preflight proofs and the window decision of each member before the lease are absent. The gate proof and the commitment decision still guard the destination.

### The rerun composer

The rerun of each stack-raised refusal follows the seven rules of the spec section "The refusal faces". Build them in one composer in the new file `land_stack_rerun.go`, so `land_stack.go` keeps room for tickets 09, 10, and 11. The composer:

- starts from the parsed overrides that the reader of ticket 06 returns;
- reads `repairsSource` and `retargetSourceTip` from `land_refusal.go` instead of a copy;
- takes the members that a route removes or resets as an argument, which tickets 09 and 10 supply.

This ticket adds no refusal face and no registry entry, so its chunk runs no `./internal/refusalroute`. A widening that adds a face adds that run.

## Acceptance

- [ ] Two tickets-only members from one destination land with one stack command, and the fixture gate tally reads one run (LS40).
- [ ] After a green two-member stack, both member assignments are released (LS41).
- [ ] After a green two-member stack, the destination holds neither tickets folder (LS42).
- [ ] After a green two-member stack, the green marker names the stack tip (LS43).
- [ ] The `members[2]` table rows read `landed` and carry each member's census count (LS44).
- [ ] The last stdout line is the `stack{...}` record with the stack tip as `published_commit` and `landed` 2 (LS45).
- [ ] The stack prints one `landing source{...}` line on stderr for each member, in file order (LS101).
- [ ] A stack whose members 1 and 2 both fail a source proof prints two refusals, each with its member index (LS50), and the ledger then holds no lease (LS51).
- [ ] While another landing holds the lease, the stack refuses with `land-lease-held`, and the tally reads 0 (LS52).
- [ ] A single landing that starts during the grade of a stack refuses with `land-lease-held` and names the assignment of member 1 (LS53).
- [ ] Two members that name one worktree exit 2, name both indexes, and leave no lease (LS63).
- [ ] A stack file that the reader refuses exits 2, and the ledger holds no lease (LS104).
- [ ] A single landing whose `-m` value is `--stack` stays a single landing (LS105).
- [ ] In a three-member stack, the last step of a `source-not-clean` route for member 2 is `bench worktree land --stack <file> --source-tip <request>=<repaired-tip>` with the request of member 2 (LS86).
- [ ] A member whose branch tip differs from its `source_tip` refuses with `source-tip-mismatch`, which names the observed tip (LS87).
- [ ] A call with an override for member 1 that refuses member 2 prints a rerun that keeps the member 1 override and adds the slot of member 2 (LS94).
- [ ] An override replaces the file `source_tip` of its member, and the stack publishes that member's merge with the override tip as second parent (LS95).

The journeys and the lease rows live in `land_stack_test.go`, and the rerun rows LS86, LS87, LS94, and LS95 live in `land_stack_rerun_test.go`, so neither file holds more than six tests. The journeys build two or three assignment worktrees from one review base with the landing fixture builders, and they count the tally with `landingGateTally`. The lease row follows `barrierLandingFixture` in `internal/worktree/land_lease_test.go`. Each expected `members` row comes from the same `toon.Table` call that the producer makes.

## Checkpoint verification

- `bench test --package ./internal/worktree`
- `bench test --package ./internal/usage`
- `bench test --package ./cmd/bench`
- `bench test --check root-conformance`

## Mutation proof

Make the release step release only the first published member. The LS41 test must fail. Spell the probe against the code that the author writes, and require `bench probe` to report `bit` with a byte-exact restore.
