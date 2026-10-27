# Prove a stack tree in the gate and grade it whole

Blocked by: 01-share-the-stack-owner-rule.md
Writes: internal/gate/stack_completion.go (new), internal/gate/stack_completion_test.go (new), internal/gate/checkpoint.go, internal/gate/completion.go, internal/gate/evaluation.go, internal/gate/reader_plan.go
Covers: LS7, LS8, LS9, LS10, LS11, LS12, LS13, LS14, LS15, LS16, LS17

## What to build

Export `gate.StackMember{Source, Folder string}` and `gate.WithStack(ctx context.Context, members []gate.StackMember) context.Context`. `Source` is the member's source tip, and `Folder` is its tickets folder, which is empty for a spec-less member. Tickets 04 and 08 call this one signature. A stack context goes with the landing value of `gate.WithLanding`, which names the destination commit.

`checkpointEvaluation` reads the stack context beside the completion context, and its signature stays. So its four callers in `internal/gate/engine.go` and `internal/gate/gate.go` keep their calls, and the reuse path of `inspectProspective` reads the same context.

When a stack context is present, `applyCheckpoint` runs the stack proof in place of the single completion proof. The proof runs for every stack context, also a context of spec-less members only, so the early return for a run with no spec and no tickets folder must not skip it. The proof lives in `internal/gate/stack_completion.go`, because `checkpoint.go` has 11 lines of headroom.

The proof captures the destination tree and each member tree, and it calls `stackpolicy.Owners` over their entries. It then applies these rules to each path of the graded tree and the destination, in this order:

1. A path under the tickets folder of a member is absent. This rule wins over the next rules for the paths of that folder.
2. A path that no member changes keeps the destination's entry. Otherwise the proof faults with `stack composition changes <path>; review the destination delta`.
3. A path that one member changes keeps that member's entry.
4. A `capture/` path that two or more members change passes.
5. A path outside `capture/` that two or more members change faults, and the message names both member indexes.

The proof also applies the verified closure of each member that has a tickets folder, as the single proof does. `completionClosure` reads `e.completionSource` today (`internal/gate/completion.go:124`), and a stack context sets no such field. So `completionClosure` takes the source tip as a parameter, and its one caller in `completionTree` passes `e.completionSource`. The stack proof passes each member's own tip. A spec-less member skips the closure, because `Closure` of an empty path reaches `reviewrecord.RecordPath("")` (`internal/commitment/repository/closure.go:42`) and fails with a non-proof error. Each fault is a `completionProofError`, so the authorization keeps the `completion proof` cause.

The reader plan gets the case `landing stack` right after `no landing value`. A landing with a stack context then prints `gate: whole gate (reader set refused: landing stack)` and runs the whole gate.

Ticket 04 grades each stack through this context. A context of one member is a valid bisect prefix.

## Acceptance

- [ ] The stack proof accepts the exact tree of a two-member stack (LS7).
- [ ] A tree of a stack of two spec-less members whose bytes at a member's changed path differ from that member's bytes refuses, and the message names the path (LS8).
- [ ] A tree of a stack of two spec-less members that changes a path that no member changes refuses with `stack composition changes <path>; review the destination delta` (LS9).
- [ ] A stack tree that keeps an entry under the tickets folder of a member refuses with `completion keeps closed <path>` (LS10).
- [ ] A stack tree whose bytes at a `capture/` path that two members change match neither member passes (LS11).
- [ ] A stack tree with a path outside `capture/` that two members change refuses, and the message names both member indexes (LS12).
- [ ] A stack tree passes only when it lacks a destination path that one member deletes (LS13).
- [ ] A one-member stack context over a member that holds the destination accepts the exact prefix tree (LS14).
- [ ] After a green of the top tree under no stack context, a grade of that tree under a stack context whose proof fails refuses (LS15).
- [ ] A landing with a stack context prints exactly `gate: whole gate (reader set refused: landing stack)` and runs the whole gate (LS16).
- [ ] A stack landing over an uncovered path prints the `landing stack` case and not the `uncovered path` case (LS17).

The reader-plan rows reuse the `landingFixture` of `internal/gate/reader_plan_test.go` from the new test file, because that file is at its 400-line limit. The proof rows follow `gradeTicketsClose` in `internal/gate/commitment_completion_test.go`, which drives `checkpointEvaluation` and then `completionTree` on a crafted tree.

## Checkpoint verification

- `bench test --package ./internal/gate`
- `bench test --check root-conformance`

## Mutation proof

Skip rule 2 for every path that no member changes. The LS9 test must fail. Spell the probe against the code that the author writes, and require `bench probe` to report `bit` with a byte-exact restore.
