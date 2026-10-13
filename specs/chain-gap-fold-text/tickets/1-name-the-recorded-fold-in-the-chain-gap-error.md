# Name the recorded fold in the chain-gap error

Blocked by: none
Writes: internal/reviewrecord/coverage.go, internal/gate/review_checkpoint_commits_test.go
Covers: none

## What to build

The review chain-gap refusal names recorded folds beside record commits and comment-only corrections after a chunk tip. The checkpoint test verifies the updated chain rule in the refusal output.

## Acceptance

- [ ] The chain-gap refusal states that only record commits, recorded folds, and comment-only corrections follow a chunk tip.
