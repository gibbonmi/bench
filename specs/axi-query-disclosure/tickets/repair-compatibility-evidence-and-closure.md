# Repair compatibility evidence and review closure

Blocked by: repair-registry-and-deep-cwd-proof.md
Writes: `cmd/bench/`, `internal/learnings/`, `internal/maps/`, `internal/guards/`, `internal/coverage/`, `internal/worktree/`, `specs/axi-query-disclosure/spec.md`, `reviews/axi-query-disclosure.md`

## What to build

Check in the QD6 pre-disclosure half of every named old/new response pair, compare the candidate response as an exact appended-help delta, reconcile QD3's red-signal token with the shipped check name, and delete the resolved review pickup artifact in the same green commit.

## Acceptance

- [ ] [RE1] (covers QD6) every changed surface and named state has a checked-in pre-disclosure response fixture plus a candidate comparison proving primary bytes, streams, exits, and argv are unchanged except for the approved help append and worktree help-spelling delta.
- [ ] [RE2] (covers QD2) the paired worktree argv fixture names the old and new grammars and proves that only `--help`, `-h`, and bare `help` differ.
- [ ] [RE3] (covers QD3) QD3's recorded red probe names the live `axi-query-registry` check, and all review findings in `reviews/axi-query-disclosure.md` are resolved before that pickup file is deleted.
