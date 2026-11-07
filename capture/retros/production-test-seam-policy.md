## Outcome

FT343 landed at `9bfb04c` with a green whole gate and `Status: implemented`.
The injected-port owner moved to `internal/conformance/injectedports/` behind one `Check(root, Policy)` boundary.
It now derives named ports across every production package, plus variable and setter seams, test-runtime branches, and environment controls that a test set reaches.
The one registry admits each seam with an owner, a consumer, a purpose, and isolation evidence.
ADR 0034 states the decided rules.

The reviewer made four decisions during the build. The reviewer retired `releaseevidence.indexEncoder`, accepted TS3-C2 as a Won't handle limitation, and widened named-port derivation as spec line 114 states.
The landing release stays incomplete, because a dead protected user retains the build worktree.

## Cost

- landing: 9bfb04c7d320f97ded7940de764119ece6b31c3c, range 9bfb04c7d320f97ded7940de764119ece6b31c3c^1..9bfb04c7d320f97ded7940de764119ece6b31c3c^2
- wall: 1771 min
- commits: 55 total, 46 bookkeeping, 9 change
- graph: 4 tickets, chain 4, 4 chunks, 4 serial checkpoints
- dispatch: 4 authors, 5 repairs
- trace: commit 59 spans in 460423 ms, lane 0 spans in 0 ms, landing 1 spans in 321479 ms
- wall per ticket: 442.9 min against 60 min, over budget
- bookkeeping per ticket: 11.5 against 7.1, over budget

## Gate-stage timings

- landing: commit 9bfb04c7d320f97ded7940de764119ece6b31c3c, trace 6cbc5ab3bc5c06dbca9e522d5e611903
- gofmt: 209 ms
- vet: 1741 ms
- test: 238132 ms
- race: 2509 ms
- system: 50501 ms
- shellcheck: 565 ms

## Ticket-versus-spec-slice and delegate performance

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| TS2 review | TS2-P1 same-package field replacement is not derived | confirmed | 0.8 | 1 | opus / medium / reviewer |
| TS2 review | TS2-P2 validator stubs have no hook-active refusal witness | confirmed | 0.5 | 1 | opus / medium / reviewer |
| TS2 review | TS2-C1 no fixture covers a field or index replacement | confirmed | 0.8 | 1 | opus / medium / reviewer |
| TS3 review | TS3-P1 two environment rows name shell-only consumers | confirmed | 0.8 | 1 | opus / medium / reviewer |
| TS3 review | TS3-C1 an assertion string counts as an environment set | confirmed | 0.7 | 1 | opus / medium / reviewer |
| TS3 review | TS3-C2 reach fails open through function values | confirmed | 0.6 | 1 | opus / medium / reviewer |
| TS4 review | TS4-P1 the fault check covers only the owner | confirmed | 0.8 | 1 | opus / medium / reviewer |
| TS4 review | TS4-C1 no fixture has a second fault producer | confirmed | 0.8 | 1 | opus / medium / reviewer |
| TS4 review | TS4-C2 a var-declared local hides a real error | confirmed | 0.7 | 1 | opus / medium / reviewer |

Brier mean: 0.088 over 9 pairs, with 0 abstentions.

The four authors ran opus / high, one fresh author for each ticket.
Ticket 01 landed in one pass with zero findings.
Tickets 02, 03, and 04 each stopped before their commit on a decision that the spec did not settle, and each resumed on the same assignment.
Ticket 03 took four stops and two repair cycles, and its chunk took about 4 hours 30 minutes.
Ticket 04 grew from a census close into a 28-port admission after the reviewer widened it.

## Coordinator catches

- The staleness pass found that the plan grade refuses plain chunk commands, so the build-entry amendment moved every plain command to one final `./...` run.
- The final plan named `bench test --check release-evidence-probe`, which `bench test` refuses by design, so the plan moved to the ship-tier conformance route.
- The staleness pass found four deleted profile fixtures in the ticket 03 fence and a stale DATE citation.
- The orchestrator refused a proposed consumer that never runs its owner, and it asked for a real in-process consumer.
- The orchestrator found that the recorded consumer-omission probe tested the orphan rule, and it asked for the named mutation.
- The checkpoint refused verifications by the orchestrator, because the owed performer is the latest assigned session.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 01-relocate-the-existing-audit.md | 0 | none |
| 02-admit-variable-and-setter-callers.md | 1 | edge-found |
| 03-admit-runtime-and-environment-controls.md | 2 | spec-row, edge-found |
| 04-close-census-and-omission-proof.md | 2 | delegate-error, edge-found |

## Agent-experience improvements

### Bench CLI

- Add a `--witness <file>` option to `bench record ticket`, so that the author's record commit carries its witness table and the orchestrator writes no separate witness commit.
  Feeds: new
- Add a plan-edit form for one completion-plan requirement command, so that a build-entry amendment needs no hand edit of the plan fence.
  Feeds: new
- Let `bench record assignment` accept a dirty spec when the only change is a `bench spec widen` change, so that a widen and its assignment share one plan commit.
  Feeds: new
- Let `bench preflight build --charge` accept a dirty tree whose changes sit inside the ticket's fence, so that a mid-ticket widen can regenerate the author's charge.
  Feeds: new
- Let the chunk checkpoint accept a verification by any session in the ticket's assignment history at the current pair, so that an evidence rerun needs no resume of the latest session.
  Feeds: new
- Let a disposition-only spec line ride in a review-record commit without a move of the chunk subject, so that a reviewer decision needs no amendment, refreeze, and three reaffirmations.
  Feeds: new
- Give `bench probe` a route to a ship-tier conformance check, so that a ship-tier witness such as PS34 can show a demonstrated omission red.
  Feeds: new
- Make the protection owner treat a user whose actor and process group are dead on the current boot as released, so that a landing release completes.
  Feeds: new
- Count sanctioned read heads such as `rg` apart in the assignment census, because 160 of 168 raw calls were `rg` and they hide the other classes.
  Feeds: new

### Skills

- Make `craft-spec` author chunk inventories as probes only and the plain package commands as one final `./...` run.
  Feeds: new
- Make `craft-spec` state, for a spec that defines an eligibility rule and also lists an observed census, which of the two bounds the admission set.
  Feeds: new

### Process

- Land a reviewer decision that adds spec prose in a plan commit before the chunk freeze, as TS4 did, and not after it, as TS3 did.
  Feeds: new
- Make each author and repair charge require the exact failing assertion and the source commit for each witness red.
  Feeds: new