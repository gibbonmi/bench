# Repair compatibility evidence and review closure

Blocked by: repair-registry-and-deep-cwd-proof.md
Writes: `cmd/bench/`, `internal/learnings/`, `internal/maps/`, `internal/guards/`, `internal/coverage/`, `internal/worktree/`, `specs/axi-query-disclosure/spec.md`, `reviews/axi-query-disclosure.md`

## What to build

Check in the QD6 pre-disclosure half of every named old/new response pair, compare the candidate response as an exact appended-help delta, reconcile QD3's red-signal token with the shipped check name, and delete the resolved review pickup artifact in the same green commit.

## Acceptance

- [ ] [RE1] (covers QD6) every QD6 extraction, terminal, incomplete-scan, diagnostic, and early-refusal state has a checked-in pre-disclosure response fixture plus a candidate comparison. Disclosure states permit only the approved help append, named early refusals are byte-identical, and worktree additionally permits only the QD2 sole-help delta.
- [ ] [RE2] (covers QD2) the paired worktree argv fixture proves that only sole-argument `--help`, `-h`, and `help` differ and explicitly enumerates multi-token help, bare and value-bearing `--`, an empty token, one ordinary extra, and multiple ordinary extras with old/new stdout, stderr, and exit.
- [ ] [RE3] (covers QD3) QD3's recorded red probe names the live `axi-query-registry` check, and all review findings in `reviews/axi-query-disclosure.md` are resolved before that pickup file is deleted.
