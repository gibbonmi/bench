# Repair landing-refusal review findings

Blocked by: distinguish-published-incomplete-exit.md, enrich-refusals-through-one-emitter.md, name-reauthorize-recovery-for-lost-request.md
Writes: internal/worktree/worktree.go, internal/worktree/ownership.go, internal/worktree/reauthorize.go, internal/worktree/reauthorize_test.go, internal/worktree/land.go, internal/worktree/land_test.go, reviews/landing-refusal-diagnostics.md

## What to build

Close the accepted landing-diagnostic review findings at the existing command
seams. Reauthorization proves the addressed assignment before diagnosing a
non-ancestor review base with the assignment's recorded start. A resume with a
lost request token finds the same unique active assignment recovery as a first
landing attempt while completed resumes still authenticate through their
terminal receipt. Retained-release continuations remain safe for every
target/request hostile-input pairing.

## Acceptance

- [ ] Reauthorizing with a non-ancestor `--base` and otherwise-valid assignment evidence exits 1 and names the recorded assignment start as `wanted`.
- [ ] Resuming a published landing with an unknown request and exactly one active assignment exits 1 with that assignment id and a reauthorize continuation.
- [ ] A completed resume with no active assignment still replays its terminal receipt.
- [ ] A retained release with a safe target and control-bearing request exits 1 on stderr with `<request>`, no caller token, and no raw control bytes.
