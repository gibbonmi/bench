# Every write-verb refusal names a typed recovery route

Status: staged

Roadmap: FT393

Decision source: named reviewed artifact `roadmap/FT393.md`, opened by drain d-6259cc0d8421, which folded FT330 into it, and placed second in quality-1 by the reviewer on 2026-10-07

Audience: every repository that links the kit

Verification log: 2 iteration(s) to accept — the review reached its cap. A later read-only consultation found four folds, and the reviewer approved them.

## Problem

A write verb refuses, and the operator must find the way out alone.
Rules that are correct alone block the recovery of each other.
On 2026-10-07 a red source could not fold `main`, because the merge refused the inherited red and named no next action.
A dirty fold then refused with no next action.
The destructive-git guard denies two repairs that Bench prints: a raw `git merge` after a composition conflict, and `git restore --worktree` after a `bench commit` exit 3.
A checkpoint refusal prints one fixed help row, `bench gate --fresh, retry after restoring repository write access`, for every cause.

Only the landing declares its refusal faces in a registry today.
The merge, the reset, the commit, the gate checkpoint, and the commitment verbs each compose their recovery text by hand, and several print none.

## Solution

One refusal-route registry serves every write verb: commit, merge, reset, land, the gate checkpoint, and commitment.
Each registered face declares its verb, its name, its authority, and a typed route.
The authority is `agent` or `reviewer`.
An agent route lists the steps that an agent runs.
A reviewer route names the decision or the hand step that the reviewer owns.

Each verb prints the route of the face that refused, in the output shape that the verb uses today.
A conformance check proves that every agent route passes the wired guards.
A producing fixture for each face drives the refusal, follows the printed route, and shows that the face no longer prints.
A static check proves that no write-verb source composes a route outside the registry.
`bench recovery` renders the recovery matrix from the registry.

## User stories

Line: `opus` / high.
Implementation-line reason: the hardest chunk moves the landing's registry into a shared package and changes route text that many worktree tests pin. The spec fixes each face and its route, the seams follow the landing precedent, and each face has a red-capable producing fixture.
Harder chunks: RR-C1b, RR-C3.

### One registry declares every face

1. As an operator, I want every refusal of a write verb to come from one registry, so that no verb invents its own recovery text.
2. As an operator, I want each face to declare an agent or reviewer authority, so that I know whether I may run the route myself.
3. As an operator, I want a reviewer route to start with a fixed marker, so that an agent stops and hands back without guesswork.
4. As an operator, I want an agent route to list its steps in run order, so that I can follow it without archaeology.
5. As a maintainer, I want a face with no route step to turn the gate red, so that no face prints an empty route.
6. As a maintainer, I want two faces with the same name to turn the gate red, so that a name identifies one face.
7. As a maintainer, I want an unregistered face name to refuse with a reviewer route, so that a bookkeeping fault never aborts the session.

### Every agent route passes the wired guards

8. As an agent, I want every agent route to pass the destructive-git guard, so that the guard never denies a repair that Bench printed.
9. As an agent, I want every agent route to pass the follow-on guard, so that no route chains Bench calls or enters the pool.
10. As an agent, I want no agent route to run a Bench child through `bench worktree exec`, so that the route survives FT341's exec refusal.
11. As a maintainer, I want the guard check to turn red on a denied step, so that the check bites.
12. As an agent, I want a tree-scoped Bench step to address a worktree through `--in <label>`, so that the route uses the FT341 form.

### Each face has a producing fixture that follows its route

13. As a maintainer, I want each face to have exactly one producing fixture, so that a face added without proof turns the gate red.
14. As a maintainer, I want each fixture to follow the printed route and rerun the command, so that a wrong route turns the gate red.

### The landing keeps its faces in the shared registry

15. As an orchestrator, I want each landing face to keep its route, except for its authority or a guard-safe form, so that nothing else moves.
16. As an orchestrator, I want a composition conflict to print a reviewer route, so that the agent hands back the merge that the guard reserves.
17. As an orchestrator, I want a dirty landing destination to print a reviewer route, so that the agent never discards the reviewer's uncommitted work.
18. As an orchestrator, I want a dirty reviewed source to print a `bench commit --in <label>` route, so that the agent commits its own work through Bench.
19. As an orchestrator, I want each landing preflight route to keep ending with my own re-run, so that one paste finishes the recovery.

### The merge names a way out

20. As an orchestrator, I want a fold refusal for a red target to name a route, so that a red source can fold `main`.
21. As an orchestrator, I want that route to repair the red, commit it, and rerun the fold, so that the gate stays the oracle.
22. As an orchestrator, I want a fold into a dirty target to name the commit route, so that I can recover from a dirty fold.
23. As an orchestrator, I want a fold to separate a target red from a fold red, so that I hand back only a fold red.
24. As an orchestrator, I want a merge conflict to print the same reviewer route as the landing conflict, so that both verbs name one repair.
25. As an orchestrator, I want a dirty sibling to name `bench commit --in <label>`, so that the route uses no exec form.
26. As an operator, I want the authorization refusal sentence to drop its inline action, so that the route prints once, from the face.

### The reset names a way out

27. As an orchestrator, I want each reset refusal to name its face's route, so that a failed reset plan is recoverable.

### The commit names a way out

28. As an orchestrator, I want a commit exit 3 to name the reset plan at the published commit, so that the guard allows the repair.
29. As an orchestrator, I want that reset route to make the checkout match the published commit, so that the route repairs exit 3.
30. As an orchestrator, I want each commit refusal to name its face's route, so that a refused commit is recoverable.

### The gate checkpoint names the route of its cause

31. As an orchestrator, I want a checkpoint refusal to name the route of its cause, so that an evidence refusal never names write access.
32. As an orchestrator, I want a completion-evidence refusal to name `bench preflight review <slug>`, so that one read shows what the evidence lacks.
33. As an orchestrator, I want the FT392 dirty-checkout refusal to name the commit route, so that I can recover from it.
34. As an orchestrator, I want a subject-capture fault to name `bench doctor` and a fresh rerun, so that the fallback fits its faults.
35. As an orchestrator, I want the fixed write-access help row gone from the checkpoint refusal, so that one route prints per refusal.

### The commitment verb names a route of the right authority

36. As an orchestrator, I want a commitment refusal that only the reviewer clears to print a reviewer route, so that I stop and hand back.
37. As an orchestrator, I want a commitment refusal that I can clear to print an agent route, so that I keep my own fix.

### No verb composes a route outside the registry

38. As a maintainer, I want a static check to refuse a route literal in write-verb source, so that no later edit bypasses the registry.
39. As a maintainer, I want the static check to turn red on a planted route literal, so that the check bites.

### The recovery matrix renders from the registry

40. As an operator, I want `bench recovery` to list every face with its verb, authority, and route, so that I can read the matrix early.
41. As a maintainer, I want the matrix to have one row per registered face, so that the matrix and the registry cannot drift.
42. As a maintainer, I want the reference guide to name `bench recovery` as the recovery matrix, so that a new teammate finds it.

### Reviewed exclusions

43. As a reviewer, I want the moved-`main` decision to stay with FT342, so that this registry routes to that decision and does not make it.
44. As a reviewer, I want the exit-2 grammar refusals to keep their usage line, so that the usage grammar keeps its one owner.

## Implementation decisions

### The registry

A new leaf package `internal/refusalroute` owns the registry.
It imports no write-verb package, so each verb package imports it with no cycle.
It holds the one ordered face inventory, the route step types, the route renderer, and the `bench recovery` command.
Each verb declares its faces in its own file of the package, and the registry composes them into the one ordered inventory.

A face declares five facts:

- the verb: `commit`, `merge`, `reset`, `land`, `gate`, or `commitment`
- the name, unique across the registry
- the sentence, or none when a policy owns the sentence at the refusal
- the authority: `agent` or `reviewer`
- the route: an ordered list of one or more steps

A step is a command or an instruction.
A command step is a template over named facts that the raising site supplies, such as the re-run, the label, or the commit.
An instruction step is one imperative sentence for an action that the agent does with its own tools.
The renderer joins the steps with `; then `.

A reviewer route renders with the prefix `reviewer: ` before its steps.
A slot that the operator fills, such as `<msg>`, renders as its placeholder.
Only the production renderer prints a placeholder; the guard check fills every slot with a sample value.

The one constructor takes the face name and the facts.
It returns a typed refusal that carries the face, the sentence, the paths, and the rendered route.
An unknown face name returns the refusal `refusal face <name> is unregistered`, with a reviewer route.
This keeps the landing's fail-soft rule for a bookkeeping fault.

A raising site derives the `<label>` fact through `intent.AssignmentsOwning` over its own root.
The landing and the merge already hold their assignment, and the commit and the gate resolve theirs at the refusal.
When no assignment owns the root, the slot prints `<label>`.
A tree-scoped step names its tree target, because `--in` counts only as the first argument after the verb.
The `resume-marker` face has reviewer authority, because only a landing on the primary checkout advances the green marker.

### Authority

A face has agent authority when an agent can clear its cause with Bench verbs and its own file edits, inside its own worktree.
A face has reviewer authority when the clear needs one of these:

- a raw merge
- a change to the primary checkout
- a commitment change
- a decision that FT342 owns
- an edit of a gate check, such as the lane declaration

The destructive-git guard already states that the merge and any history rewrite are the reviewer's, so the conflict faces take reviewer authority.

### Output shape

Each verb keeps its output shape.
The landing, the merge, and the reset print `next=<route>` in the `refused{...}` record.
The commit prints `next=<route>` on its own stderr line after the refusal sentence, and its exit 3 record keeps `committed{published_commit=…,path=…,next=…}`.
The gate checkpoint prints `next=<route>` on stderr after its reason, and it no longer prints the fixed `help[1]{cmd,why}` row.
The commitment verb keeps its `next[1]{command}` table, and the cell holds the face's route.

The landing's skipped-proof sentence, `later proofs in this group did not run`, stays the first step of the route, after the authority marker.
So a short-circuited reviewer route starts with `reviewer: later proofs in this group did not run; `.

### The face inventory

The registry is the authoritative inventory.
The build sweeps every refusal site of the six verbs and maps each site to one face.
Several sites can share one face when they share one route.
Each verb has one reviewer face, `<verb>-handback`, for a cause outside agent authority that has no reviewer step of its own.

These faces and routes are required.

| verb | face | authority | route |
|---|---|---|---|
| land | `destination-not-clean` | reviewer | commit or discard the destination's uncommitted work; then the re-run |
| land | `destination-collision` | reviewer | move the `refusal_paths` entries out of the landing checkout; then the re-run |
| land | `source-tip-mismatch` | agent | the re-run, re-pointed at the source tip that the tree holds |
| land | `source-not-clean` | agent | `bench commit --in <label> -m <msg> -- <path>...`; then `/bench-review-implementation`; then the re-run at the repaired source tip |
| land | `source-not-fenced` | agent | the current fence instruction; then the re-run at the repaired source tip |
| land | `composition-conflict` | reviewer | the hand merge of the destination commit; then `bench commit`; then `/bench-review-implementation`; then the re-run |
| land | `composition-conflict-pending` | reviewer | finish the merge in progress; then `/bench-review-implementation`; then the re-run |
| land | `resume-destination-residue` | reviewer | commit or discard the destination's uncommitted work; then the resume |
| land | `resume-marker` | reviewer | land a green landing on main that covers the published commit, or restore main to it; then the resume |
| land | `land-red` | agent | repair each failure that the gate reports in `<label>`; then `bench commit --in <label> -m <msg> -- <path>...`; then `/bench-review-implementation`; then the re-run |
| land | `land-infrastructure` | agent | `bench doctor`; then the re-run |
| merge | `merge-target-red` | agent | repair each failing check in `<label>`; then `bench commit --in <label> -m <msg> -- <path>...`; then the re-run |
| merge | `merge-fold-red` | reviewer | the fold of `<from>` adds a red to `<label>`; decide who repairs a red that the fold of a moved `main` adds, the open FT342 question |
| merge | `merge-target-not-clean` | agent | `bench commit --in <label> -m <msg> -- <path>...`; then the re-run |
| merge | `merge-sibling-not-clean` | agent | `bench commit --in <sibling-label> -m <msg> -- <path>...`; then the re-run |
| merge | `merge-conflict` | reviewer | the hand merge of the incoming commit; then `bench commit`; then the re-run |
| merge | `merge-infrastructure` | agent | `bench doctor`; then the re-run |
| reset | `reset-checkout-conflicted` | agent | `bench worktree clean <id>` |
| reset | `reset-plan-stale` | agent | the current plan command |
| reset | `reset-tree-missing` | agent | the current missing-tree route |
| merge | `merge-published-unreconciled` | agent | the current reset plan at the published tip |
| land | `land-incomplete` | agent | the current resume command |
| commit | `commit-published-unreconciled` | agent | `bench worktree reset --to <published-commit> <checkout>` |
| commit | `commit-primary-checkout` | agent | `bench worktree create --request <opaque-id> --label <work-item>` |
| commit | `commit-red` | agent | repair each failure that the run reports; then the re-run |
| commit | `commit-infrastructure` | agent | `bench doctor`; then the re-run |
| commit | `commit-named-path` | agent | correct the paths or the files that the refusal names; then the re-run |
| commit | `commit-tip-moved` | agent | the re-run |
| gate | `checkpoint-completion-evidence` | agent | `bench preflight review <slug>` |
| gate | `checkpoint-dirty-checkout` | agent | `bench commit --in <label> -m <msg> -- <path>...`; then the re-run |
| gate | `checkpoint-composition` | reviewer | the delivery closure of the spec does not compose; hand back |
| gate | `checkpoint-tip-moved` | agent | the re-run |
| gate | `checkpoint-subject-unavailable` | agent | `bench doctor`; then the re-run with `--fresh` |
| commitment | `commitment-needs-assignment` | agent | `bench worktree create --request <request> --label <label>` |
| commitment | `commitment-plan-input` | agent | correct the input file; then the re-run |
| commitment | `commitment-verify-evidence` | agent | `bench commitment verify --milestone <id> --evidence <file>` |
| commitment | `commitment-decision` | reviewer | `bench commitment plan --input <file>` |
| commitment | `commitment-unbound` | agent | `bench commitment start --outcome <id> --request <request> --deliverable <path>` |
| commitment | `commitment-light-path-outside` | agent | add the path to the `Writes:` line of `<ticket>`; then the re-run with `<ticket>` among its paths |
| commitment | `commitment-run-unknown` | agent | `bench commitment inventory` |

The `<verb>-handback` reviewer faces join this list for each verb.
On the resume path, a dirty source takes `land-handback`, because a commit there moves the source away from the published tip.
The `commitment-decision` face covers each start, block, plan, and approve refusal whose clear changes the active commitment.

### The commitment route tails

The commitment policy errors carry their route today as an inline tail, such as `; run bench commitment plan --input <file>`.
Each such error becomes a typed error that names its face, and its sentence drops the tail.
The verb that prints the error renders the face's route, so `bench commit` and `bench commitment` print one route from one source.
No reader parses the tail text: `commitcmd.Outlook` replaces any projection error with `commitment.Unreadable()`.
`bench preflight` prints the sentence beside its own recovery column, and that column does not change.

These sites move to faces:

- `commitment-needs-assignment`: the owned-assignment errors in `candidate.go`, `readiness.go`, and `publication.go`
- `commitment-decision`: the adoption, protected-commitment, policy-approval, sequence, continuation, legacy-scope, and obligation errors in `admission.go`, `candidate.go`, `continuation.go`, `publication.go`, and `delivery.go`
- `commitment-unbound`: `errUnbound` in `readiness.go`
- `commitment-light-path-outside`: the light-path fence error in `light_path.go`
- `commitment-run-unknown`: the unknown-run continuation error in `continuation.go`
- `commitment-verify-evidence`: the missing verification receipt in `verification.go`

A face's verb names the verb whose rule raises it, and another write verb can print it.
So `bench commit` prints a commitment face when the commitment policy refuses its candidate.
The reset table keeps each current route of the reset verb unchanged.
The missing-tree route clears its cause: `bench worktree clean --landed` retires a landed assignment whose tree is missing, and `bench worktree release` releases an unlanded one (reviewer decision 2026-10-07).

### The red-source fold (collision 5a)

The decision source closes the merge policy: the gate stays the oracle, and a fold publishes only on a green or a lane pass.
So a fold red picks its face by its cause.
When the target tip alone is red, the `merge-target-red` face routes the agent to repair the target first.
The repair commit runs the lane on the repaired tree, and the fold then composes onto a green target.
When the target tip alone is green, the fold adds the red.
Then the `merge-fold-red` face routes to the reviewer, because FT342 owns the moved-`main` decision.

Only the `candidate` gate kind carries this attribution, and it takes `merge-fold-red`.
The `inherited` kind means only that the target has no reusable green marker, so a never-graded target can hide a fold red.
The `lane fail` kind carries no attribution, and the collision 5a repro refuses with it.
So on an `inherited` or a `lane fail` fold, the merge grades the target tip alone, before it picks the face.
That run uses the same lane, or the whole gate when the project declares no lane.
It runs on the refusal path only, and it changes no publish rule.

### The commit exit 3 (collision 8b)

The commit publishes, and then the checkout fails to reconcile.
The route names the reset plan at the published commit: `bench worktree reset --to <published-commit> <checkout>`.
The reset verb keeps the dirty layer in a recoverable envelope and then prints its own `--apply` command.
`<checkout>` is the commit's root path, shell-quoted, or the placeholder `<checkout>` when the path is not line-safe.

### The conflict repair (collision 8a)

The conflict faces of the landing and the merge take reviewer authority.
Their route keeps the hand merge text that `conflictRepairPrefix` and `conflictContinuePrefix` print today, behind the `reviewer: ` prefix.
The guard check grades agent routes only, so a reviewer step that names `git merge` passes the check.
The reference guide already states that the reviewer does this merge with raw Git, so the authority matches the decided state.

### The gate checkpoint (FT330)

Each cause in the pre-oracle funnel maps to its own face.
The completion-evidence face keeps the route that `routedRefusal` prints today.
The complete checkpoint raises three refusals of its own in `executeCompleteCheckpoint`, through `operational`, outside the pre-oracle funnel:

- the dirty-checkout refusal takes `checkpoint-dirty-checkout`
- the `complete checkpoint source unavailable` refusals take `checkpoint-subject-unavailable`
- the `complete checkpoint unavailable` refusal, the `published.Tree` compose error, takes `checkpoint-composition`

The constant `cleanCheckoutRefusal` stays the one source of the dirty-checkout sentence.
So the `checkpoint-dirty-checkout` face declares no sentence; it adds the route.
A subject-capture fault that no other face claims takes `checkpoint-subject-unavailable`.
The checkpoint no longer prints the fixed write-access help row.

### The authorization sentence

`landing.refusalMessage` keeps its prefix, its kind, and its explanation.
It drops its inline action, because the face now carries the route.
The merge's single infrastructure retry matches the typed kind, not the old sentence.
The landing's own authorization red, raised in `LandReviewed`, takes `land-red` for a red kind and `land-infrastructure` for an infrastructure kind.
So no landing red loses its route when the sentence drops its action.

### The FT392 dependency

This spec depends on the landed FT392 complete checkpoint.
The spec's base, `120e2715`, contains that landing, and the gate faces route the refusals that FT392 added.
The fence names `internal/gate/complete_checkpoint.go` and `internal/gate/complete_checkpoint_test.go`, because the dirty-checkout face changes them.

### The guard check

A conformance test renders each agent route with sample facts.
It fills every slot with a sample value, operator slots included, because the guard lexer reads `<msg>` as a redirection.
A path fact takes a sample under the pool prefix, so a `cd` or a `git -C` form turns red.
The test classifies each command step through `gitguard.Classify`, `benchguard.PoolReference`, and `benchguard.Classify`.
It also refuses a step that runs `bench worktree exec <target> -- bench ...`, under FT341's closed decision.
The production guards do not change.

### The static bypass check

The registry owns the `next=` field label and the commitment `next` table, so no write-verb printer spells either one.
The exit-3 remainders of the commit, the merge, and the landing print their routes through faces too.
The reset plan's `--apply` pointer renders through the same field function.

The registry exports the one `next` label.
The worktree `refusal` record printer, the commit help line, and the commitment verify success table compose their label from it.

A conformance test scans the string literals in the production Go files of the write-verb packages.
It refuses a literal outside `internal/refusalroute` that holds `next=`, `next[`, `run bench `, or the route joiner `; then `, or that equals `"next"`.
The scanned packages include `internal/commitment` and `internal/commitment/repository`.
The test holds one reviewed allowlist of files whose routes serve a non-write verb.
The allowlist holds `internal/worktree/path.go`, `internal/worktree/build.go`, and `internal/worktree/tree_target.go`, and each other entry needs a non-write caller.

### The recovery matrix

`bench recovery` is a repository-scoped read verb.
It prints `recovery[N]{verb,face,authority,route}` with one row per registered face, in registry order.
It reads only the compiled registry.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| RR-C1a / 01-create-the-shared-refusal-route-registry.md | The shared registry seam: the face and step types, the renderer, and the one constructor | RR04-RR07, RR13 | `internal/refusalroute` | no |
| RR-C1b / 02-move-the-landing-faces-into-the-shared-registry.md, 03-prove-each-agent-route-passes-the-wired-guards.md, 04-render-the-recovery-matrix-from-the-registry.md | The landing faces on the registry, the guard check, and `bench recovery` | RR01-RR03, RR08-RR12, RR14-RR20, RR51-RR54, RR58, RR60 | `internal/refusalroute`, `internal/conformance`, `internal/worktree` landing tests, `cmd/bench` | yes |
| RR-C2 / 05-route-each-merge-refusal-through-the-registry.md, 06-give-the-red-source-fold-an-exit.md, 07-route-each-reset-refusal-through-the-registry.md | The merge and the reset print registry routes, and the red-source fold has an exit | RR21-RR31, RR55-RR57, RR59, RR67 | `internal/worktree` merge and reset tests, `internal/landing` | no |
| RR-C3 / 08-route-the-commit-exit-3-to-the-reset-plan.md, 09-route-each-commit-refusal-through-the-registry.md | The commit prints registry routes, and the exit 3 route is the reset plan | RR32-RR38 | `internal/commit`, `internal/worktree` | yes |
| RR-C4 / 10-route-each-checkpoint-refusal-by-its-cause.md | The gate checkpoint prints the route of its cause | RR39-RR44, RR66, RR68 | `internal/gate` | no |
| RR-C5 / 11-route-the-commitment-policy-refusals-through-faces.md, 12-print-the-commitment-verb-routes-from-the-registry.md | The commitment verb prints routes of the right authority | RR45-RR48, RR61-RR64 | `internal/commit`, `internal/commitment`, `internal/commitment/commitcmd`, `cmd/bench` | no |
| RR-C6 / 13-refuse-a-route-literal-outside-the-registry.md | No write-verb source composes a route outside the registry | RR49-RR50, RR65 | `internal/conformance` | no |

The slice split RR-C1 into RR-C1a and RR-C1b.
RR-C1a holds ticket 01, which creates the seam that every later ticket consumes, so it is its own small review chunk.
RR-C1b holds tickets 02 to 04.
The other chunk IDs did not change.

RR-C1a comes first, because every later chunk raises its faces through the shared constructor.
Its chunk review closes before ticket 02 starts.
RR-C2 changes the authorization sentence that RR-C3 also prints.

Ticket 11 in RR-C5 runs after RR-C3, because it adds its rows to the commit test file that ticket 08 creates.
Each verb ticket adds its faces to the one registry file, so the tickets run in serial order.
RR-C6 comes last, because it turns red on any route that an earlier chunk has not moved.

### Completion plan

The version 2 plan records future implementation evidence.
It claims no current implementation pass, red, or probe result.
The orchestrator records each author session before that author's dispatch.

```bench-completion-plan
{
  "version": 2,
  "chunks": [
    {
      "id": "RR-C1a",
      "tickets": [
        "01-create-the-shared-refusal-route-registry.md"
      ],
      "verification": [
        {
          "id": "t1-registry",
          "command": "bench test --package ./internal/refusalroute",
          "ticket": "01-create-the-shared-refusal-route-registry.md"
        },
        {
          "id": "t1-renderer-proof",
          "command": "bench test --package ./internal/refusalroute --run 'TestRouteRendering'",
          "probe": "Make the renderer join the route steps with `, ` instead of `; then `. TestRouteRendering must fail, then pass after source restoration.",
          "ticket": "01-create-the-shared-refusal-route-registry.md"
        }
      ]
    },
    {
      "id": "RR-C1b",
      "tickets": [
        "02-move-the-landing-faces-into-the-shared-registry.md",
        "03-prove-each-agent-route-passes-the-wired-guards.md",
        "04-render-the-recovery-matrix-from-the-registry.md"
      ],
      "verification": [
        {
          "id": "t2-registry",
          "command": "bench test --package ./internal/refusalroute",
          "ticket": "02-move-the-landing-faces-into-the-shared-registry.md"
        },
        {
          "id": "t2-worktree-package",
          "command": "bench test --package ./internal/worktree",
          "ticket": "02-move-the-landing-faces-into-the-shared-registry.md"
        },
        {
          "id": "t2-landing-faces",
          "command": "bench test --package ./internal/worktree --run 'TestLandingRefusalRegistryHasAProducingFixture|TestConflictRepairIsAReviewerRoute|TestLandCommandReportsEveryRefusalInOnePreflight'",
          "ticket": "02-move-the-landing-faces-into-the-shared-registry.md"
        },
        {
          "id": "t2-conflict-proof",
          "command": "bench test --package ./internal/worktree --run 'TestConflictRepairIsAReviewerRoute'",
          "probe": "Declare the `composition-conflict` face with the agent authority. TestConflictRepairIsAReviewerRoute must fail, then pass after source restoration.",
          "ticket": "02-move-the-landing-faces-into-the-shared-registry.md"
        },
        {
          "id": "t3-registry",
          "command": "bench test --package ./internal/refusalroute",
          "ticket": "03-prove-each-agent-route-passes-the-wired-guards.md"
        },
        {
          "id": "t3-guard-check",
          "command": "bench test --package ./internal/conformance --run 'TestAgentRoutesPassTheWiredGuards|TestAgentRouteGuardCheckBites'",
          "ticket": "03-prove-each-agent-route-passes-the-wired-guards.md"
        },
        {
          "id": "t3-landing-faces",
          "command": "bench test --package ./internal/worktree --run 'TestLandingRefusalRegistryHasAProducingFixture|TestLandingFacesFollowTheirRoutes|TestConflictRepairIsAReviewerRoute|TestUnsafePathRouteUsesThePlaceholder|TestLandCommandReportsEveryRefusalInOnePreflight'",
          "ticket": "03-prove-each-agent-route-passes-the-wired-guards.md"
        },
        {
          "id": "t3-worktree-package",
          "command": "bench test --package ./internal/worktree",
          "ticket": "03-prove-each-agent-route-passes-the-wired-guards.md"
        },
        {
          "id": "t3-guard-proof",
          "command": "bench test --package ./internal/conformance --run 'TestAgentRoutesPassTheWiredGuards'",
          "probe": "Declare the first route step of the `source-not-clean` face as the command `git merge <commit>`. TestAgentRoutesPassTheWiredGuards must fail, then pass after source restoration.",
          "ticket": "03-prove-each-agent-route-passes-the-wired-guards.md"
        },
        {
          "id": "t4-registry",
          "command": "bench test --package ./internal/refusalroute",
          "ticket": "04-render-the-recovery-matrix-from-the-registry.md"
        },
        {
          "id": "t4-guard-check",
          "command": "bench test --package ./internal/conformance --run 'TestAgentRoutesPassTheWiredGuards|TestAgentRouteGuardCheckBites'",
          "ticket": "04-render-the-recovery-matrix-from-the-registry.md"
        },
        {
          "id": "t4-landing-faces",
          "command": "bench test --package ./internal/worktree --run 'TestLandingRefusalRegistryHasAProducingFixture|TestLandingFacesFollowTheirRoutes|TestConflictRepairIsAReviewerRoute|TestUnsafePathRouteUsesThePlaceholder|TestLandCommandReportsEveryRefusalInOnePreflight'",
          "ticket": "04-render-the-recovery-matrix-from-the-registry.md"
        },
        {
          "id": "t4-worktree-package",
          "command": "bench test --package ./internal/worktree",
          "ticket": "04-render-the-recovery-matrix-from-the-registry.md"
        },
        {
          "id": "t4-recovery-verb",
          "command": "bench test --package ./cmd/bench --run 'TestHelpInventoryIsComplete'",
          "ticket": "04-render-the-recovery-matrix-from-the-registry.md"
        },
        {
          "id": "t4-recovery-proof",
          "command": "bench test --package ./internal/refusalroute --run 'TestRecoveryListsEveryFace'",
          "probe": "Make `bench recovery` skip the last registered face. TestRecoveryListsEveryFace must fail, then pass after source restoration.",
          "ticket": "04-render-the-recovery-matrix-from-the-registry.md"
        }
      ]
    },
    {
      "id": "RR-C2",
      "tickets": [
        "05-route-each-merge-refusal-through-the-registry.md",
        "06-give-the-red-source-fold-an-exit.md",
        "07-route-each-reset-refusal-through-the-registry.md"
      ],
      "verification": [
        {
          "id": "t5-merge-routes",
          "command": "bench test --package ./internal/worktree --run 'TestMergeFacesFollowTheirRoutes|TestConflictRepairIsAReviewerRoute|TestMergeRetriesOnlyAVerifiedEmptyReasonInfrastructureRefusal|TestLandingRedRouteNamesTheRepair'",
          "ticket": "05-route-each-merge-refusal-through-the-registry.md"
        },
        {
          "id": "t5-worktree-package",
          "command": "bench test --package ./internal/worktree",
          "ticket": "05-route-each-merge-refusal-through-the-registry.md"
        },
        {
          "id": "t5-landing-package",
          "command": "bench test --package ./internal/landing",
          "ticket": "05-route-each-merge-refusal-through-the-registry.md"
        },
        {
          "id": "t5-conflict-proof",
          "command": "bench test --package ./internal/worktree --run 'TestConflictRepairIsAReviewerRoute'",
          "probe": "Declare the `merge-conflict` face with the agent authority. TestConflictRepairIsAReviewerRoute must fail, then pass after source restoration.",
          "ticket": "05-route-each-merge-refusal-through-the-registry.md"
        },
        {
          "id": "t6-merge-routes",
          "command": "bench test --package ./internal/worktree --run 'TestMergeFacesFollowTheirRoutes|TestRedSourceFoldNamesAnExit|TestConflictRepairIsAReviewerRoute|TestMergeRetriesOnlyAVerifiedEmptyReasonInfrastructureRefusal'",
          "ticket": "06-give-the-red-source-fold-an-exit.md"
        },
        {
          "id": "t6-worktree-package",
          "command": "bench test --package ./internal/worktree",
          "ticket": "06-give-the-red-source-fold-an-exit.md"
        },
        {
          "id": "t6-landing-package",
          "command": "bench test --package ./internal/landing",
          "ticket": "06-give-the-red-source-fold-an-exit.md"
        },
        {
          "id": "t6-target-alone-proof",
          "command": "bench test --package ./internal/worktree --run 'TestRedSourceFoldNamesAnExit'",
          "probe": "Make an `inherited` fold take the `merge-fold-red` face without the target-alone grade. TestRedSourceFoldNamesAnExit must fail, then pass after source restoration.",
          "ticket": "06-give-the-red-source-fold-an-exit.md"
        },
        {
          "id": "t7-merge-routes",
          "command": "bench test --package ./internal/worktree --run 'TestMergeFacesFollowTheirRoutes|TestRedSourceFoldNamesAnExit|TestConflictRepairIsAReviewerRoute|TestMergeRetriesOnlyAVerifiedEmptyReasonInfrastructureRefusal'",
          "ticket": "07-route-each-reset-refusal-through-the-registry.md"
        },
        {
          "id": "t7-worktree-package",
          "command": "bench test --package ./internal/worktree",
          "ticket": "07-route-each-reset-refusal-through-the-registry.md"
        },
        {
          "id": "t7-landing-package",
          "command": "bench test --package ./internal/landing",
          "ticket": "07-route-each-reset-refusal-through-the-registry.md"
        },
        {
          "id": "t7-walk-proof",
          "command": "bench test --package ./internal/worktree --run 'TestMergeFacesFollowTheirRoutes'",
          "probe": "Register one more reset face that no fixture produces. TestMergeFacesFollowTheirRoutes must fail, then pass after source restoration.",
          "ticket": "07-route-each-reset-refusal-through-the-registry.md"
        }
      ]
    },
    {
      "id": "RR-C3",
      "tickets": [
        "08-route-the-commit-exit-3-to-the-reset-plan.md",
        "09-route-each-commit-refusal-through-the-registry.md"
      ],
      "verification": [
        {
          "id": "t8-commit-package",
          "command": "bench test --package ./internal/commit",
          "ticket": "08-route-the-commit-exit-3-to-the-reset-plan.md"
        },
        {
          "id": "t8-commit-exit-three",
          "command": "bench test --package ./internal/worktree --run 'TestCommitExitThreeRouteReconcilesTheCheckout'",
          "ticket": "08-route-the-commit-exit-3-to-the-reset-plan.md"
        },
        {
          "id": "t8-reset-route-proof",
          "command": "bench test --package ./internal/commit --run 'TestPublishedUnreconciledRouteIsTheResetPlan'",
          "probe": "Give the `commit-published-unreconciled` face the old `git restore` route. TestPublishedUnreconciledRouteIsTheResetPlan must fail, then pass after source restoration.",
          "ticket": "08-route-the-commit-exit-3-to-the-reset-plan.md"
        },
        {
          "id": "t9-commit-package",
          "command": "bench test --package ./internal/commit",
          "ticket": "09-route-each-commit-refusal-through-the-registry.md"
        },
        {
          "id": "t9-commit-exit-three",
          "command": "bench test --package ./internal/worktree --run 'TestCommitExitThreeRouteReconcilesTheCheckout'",
          "ticket": "09-route-each-commit-refusal-through-the-registry.md"
        },
        {
          "id": "t9-red-route-proof",
          "command": "bench test --package ./internal/commit --run 'TestCommitFacesFollowTheirRoutes'",
          "probe": "Drop the re-run step from the `commit-red` route. TestCommitFacesFollowTheirRoutes must fail, then pass after source restoration.",
          "ticket": "09-route-each-commit-refusal-through-the-registry.md"
        }
      ]
    },
    {
      "id": "RR-C4",
      "tickets": [
        "10-route-each-checkpoint-refusal-by-its-cause.md"
      ],
      "verification": [
        {
          "id": "t10-gate-package",
          "command": "bench test --package ./internal/gate",
          "ticket": "10-route-each-checkpoint-refusal-by-its-cause.md"
        },
        {
          "id": "t10-checkpoint-routes",
          "command": "bench test --package ./internal/gate --run 'TestCheckpointFacesFollowTheirRoutes|TestReviewCheckpointRefusalRoute|TestGateRunRetainsSubjectConstructionCause'",
          "ticket": "10-route-each-checkpoint-refusal-by-its-cause.md"
        },
        {
          "id": "t10-help-row-proof",
          "command": "bench test --package ./internal/gate --run 'TestReviewCheckpointRefusalRoute|TestGateRunRetainsSubjectConstructionCause'",
          "probe": "Print the fixed write-access `help[1]{cmd,why}` row on the checkpoint refusal again. TestReviewCheckpointRefusalRoute must fail, then pass after source restoration.",
          "ticket": "10-route-each-checkpoint-refusal-by-its-cause.md"
        }
      ]
    },
    {
      "id": "RR-C5",
      "tickets": [
        "11-route-the-commitment-policy-refusals-through-faces.md",
        "12-print-the-commitment-verb-routes-from-the-registry.md"
      ],
      "verification": [
        {
          "id": "t11-commitment-packages",
          "command": "bench test --package ./internal/commitment/...",
          "ticket": "11-route-the-commitment-policy-refusals-through-faces.md"
        },
        {
          "id": "t11-commit-package",
          "command": "bench test --package ./internal/commit",
          "ticket": "11-route-the-commitment-policy-refusals-through-faces.md"
        },
        {
          "id": "t11-authority-proof",
          "command": "bench test --package ./internal/commit --run 'TestCommitFacesFollowTheirRoutes'",
          "probe": "Give the protected-commitment refusal the `commitment-unbound` face. TestCommitFacesFollowTheirRoutes must fail, then pass after source restoration.",
          "ticket": "11-route-the-commitment-policy-refusals-through-faces.md"
        },
        {
          "id": "t12-commitment-packages",
          "command": "bench test --package ./internal/commitment/...",
          "ticket": "12-print-the-commitment-verb-routes-from-the-registry.md"
        },
        {
          "id": "t12-commit-package",
          "command": "bench test --package ./internal/commit",
          "ticket": "12-print-the-commitment-verb-routes-from-the-registry.md"
        },
        {
          "id": "t12-commitment-verb",
          "command": "bench test --package ./cmd/bench --run 'TestCommitment'",
          "ticket": "12-print-the-commitment-verb-routes-from-the-registry.md"
        },
        {
          "id": "t12-authority-proof",
          "command": "bench test --package ./internal/commitment/commitcmd --run 'TestCommitmentFacesFollowTheirRoutes'",
          "probe": "Give the outside-milestone start refusal the `commitment-needs-assignment` face. TestCommitmentFacesFollowTheirRoutes must fail, then pass after source restoration.",
          "ticket": "12-print-the-commitment-verb-routes-from-the-registry.md"
        }
      ]
    },
    {
      "id": "RR-C6",
      "tickets": [
        "13-refuse-a-route-literal-outside-the-registry.md"
      ],
      "verification": [
        {
          "id": "t13-bypass-check",
          "command": "bench test --package ./internal/conformance --run 'TestNoWriteVerbComposesARouteOutsideTheRegistry|TestRouteBypassCheckBites'",
          "ticket": "13-refuse-a-route-literal-outside-the-registry.md"
        },
        {
          "id": "t13-conformance-package",
          "command": "bench test --package ./internal/conformance",
          "ticket": "13-refuse-a-route-literal-outside-the-registry.md"
        },
        {
          "id": "t13-scan-proof",
          "command": "bench test --package ./internal/conformance --run 'TestRouteBypassCheckBites'",
          "probe": "Remove `internal/commitment/repository` from the scanned packages of the bypass check. TestRouteBypassCheckBites must fail, then pass after source restoration.",
          "ticket": "13-refuse-a-route-literal-outside-the-registry.md"
        }
      ]
    }
  ],
  "final_verification": [
    {
      "id": "ordinary-integration",
      "command": "bench test --package ./..."
    },
    {
      "id": "coverage",
      "command": "bench coverage --check refusal-route-registry"
    }
  ],
  "execution": {
    "mode": "delegate",
    "run_id": "ft393-build-20261007",
    "orchestrator_session": "claude:ft393-orchestrator-20261007",
    "author_limit": 1,
    "assignments": {
      "01-create-the-shared-refusal-route-registry.md": [
        {
          "session": "claude:ft393_t1",
          "assignment": "86659b1a8e60fc398e93bcaf62691549",
          "model": "opus",
          "effort": "high",
          "source": "2474d5c9285ff9026ec821210209356f496ba84a",
          "native_ref": "claude-agent:ft393_t1"
        },
        {
          "session": "claude:ft393_t1_repair1",
          "assignment": "86659b1a8e60fc398e93bcaf62691549",
          "model": "opus",
          "effort": "high",
          "source": "e66a113731b4f70de8e262972d5c04a9297a06c8",
          "native_ref": "claude-agent:ft393_t1_repair1",
          "predecessor": "claude:ft393_t1",
          "trigger": "user-directed",
          "stopped": "the ticket 01 author session reported completion of the record commit 62bda9ae and has no live child",
          "preserved": "e66a113731b4f70de8e262972d5c04a9297a06c8"
        }
      ],
      "02-move-the-landing-faces-into-the-shared-registry.md": [
        {
          "session": "claude:ft393_t2",
          "assignment": "86659b1a8e60fc398e93bcaf62691549",
          "model": "opus",
          "effort": "high",
          "source": "90fa563633fcf6114ed714dfa480777cd36a176d",
          "native_ref": "claude-agent:ft393_t2"
        },
        {
          "session": "claude:ft393_t2_repair1",
          "assignment": "86659b1a8e60fc398e93bcaf62691549",
          "model": "opus",
          "effort": "high",
          "source": "5a0684de8441bd7b2b82e4e52bbef2d3e0e04fc5",
          "native_ref": "claude-agent:ft393_t2_repair1",
          "predecessor": "claude:ft393_t2",
          "trigger": "user-directed",
          "stopped": "the ticket 02 author session reported completion of its RR-C1b verification records and has no live child",
          "preserved": "5a0684de8441bd7b2b82e4e52bbef2d3e0e04fc5"
        }
      ],
      "03-prove-each-agent-route-passes-the-wired-guards.md": [
        {
          "session": "claude:ft393_t3",
          "assignment": "86659b1a8e60fc398e93bcaf62691549",
          "model": "opus",
          "effort": "high",
          "source": "1d22d3209f8d20710da4da8c7bdf1378a871969a",
          "native_ref": "claude-agent:ft393_t3"
        },
        {
          "session": "claude:ft393_t3_repair1",
          "assignment": "86659b1a8e60fc398e93bcaf62691549",
          "model": "opus",
          "effort": "high",
          "source": "5a0a4f783b1237efba18757e1d6887a7887327e3",
          "native_ref": "claude-agent:ft393_t3_repair1",
          "predecessor": "claude:ft393_t3",
          "trigger": "user-directed",
          "stopped": "the ticket 03 author session reported completion of its RR-C1b verification records and has no live child",
          "preserved": "5a0a4f783b1237efba18757e1d6887a7887327e3"
        },
        {
          "session": "claude:ft393_t3_repair2",
          "assignment": "86659b1a8e60fc398e93bcaf62691549",
          "model": "opus",
          "effort": "high",
          "source": "1603d883702d9a90ad8559f4dbdbf4d27d686eae",
          "native_ref": "claude-agent:ft393_t3_repair2",
          "predecessor": "claude:ft393_t3_repair1",
          "trigger": "user-directed",
          "stopped": "the first ticket 03 repair session reported completion of its RR-C1b verification records and has no live child",
          "preserved": "1603d883702d9a90ad8559f4dbdbf4d27d686eae"
        }
      ],
      "04-render-the-recovery-matrix-from-the-registry.md": [
        {
          "session": "claude:ft393_t4",
          "assignment": "86659b1a8e60fc398e93bcaf62691549",
          "model": "opus",
          "effort": "high",
          "source": "f99f7bcd61e976c1434781e33913029eb31f1eb6",
          "native_ref": "claude-agent:ft393_t4"
        },
        {
          "session": "claude:ft393_t4_repair1",
          "assignment": "86659b1a8e60fc398e93bcaf62691549",
          "model": "opus",
          "effort": "high",
          "source": "c359944cda758e4b555be502836e29cb530511fd",
          "native_ref": "claude-agent:ft393_t4_repair1",
          "predecessor": "claude:ft393_t4",
          "trigger": "user-directed",
          "stopped": "the ticket 04 author session reported completion of its RR-C1b verification record commit 2dd0ec7e and has no live child",
          "preserved": "c359944cda758e4b555be502836e29cb530511fd"
        }
      ],
      "05-route-each-merge-refusal-through-the-registry.md": [
        {
          "session": "claude:ft393_t5",
          "assignment": "86659b1a8e60fc398e93bcaf62691549",
          "model": "opus",
          "effort": "high",
          "source": "6affd5e5e85046b4b5bed7331cf92effb42f5eb7",
          "native_ref": "claude-agent:ft393_t5"
        },
        {
          "session": "claude:ft393_t5b",
          "assignment": "86659b1a8e60fc398e93bcaf62691549",
          "model": "opus",
          "effort": "high",
          "source": "01a59319d73960148284929697ff4fc542f9f898",
          "native_ref": "claude-agent:ft393_t5b",
          "predecessor": "claude:ft393_t5",
          "trigger": "user-directed",
          "stopped": "the ticket 05 author session committed 1fcfeef8, has no live child, and carries about 419k tokens of context; the reviewer's standing rule transfers work past about 300k tokens",
          "preserved": "1fcfeef8aacc9b37ce51cc0566c053974337d2af"
        },
        {
          "session": "claude:ft393_t5_repair1",
          "assignment": "86659b1a8e60fc398e93bcaf62691549",
          "model": "opus",
          "effort": "high",
          "source": "f932db29f82ba07271a8e11c8ac5ea011b0aa3ab",
          "native_ref": "claude-agent:ft393_t5_repair1",
          "predecessor": "claude:ft393_t5b",
          "trigger": "user-directed",
          "stopped": "the ticket 05 verification session reported its RR-C2 records and has no live child",
          "preserved": "f932db29f82ba07271a8e11c8ac5ea011b0aa3ab"
        }
      ],
      "06-give-the-red-source-fold-an-exit.md": [
        {
          "session": "claude:ft393_t6",
          "assignment": "86659b1a8e60fc398e93bcaf62691549",
          "model": "opus",
          "effort": "high",
          "source": "1fcfeef8aacc9b37ce51cc0566c053974337d2af",
          "native_ref": "claude-agent:ft393_t6"
        },
        {
          "session": "claude:ft393_t6_repair1",
          "assignment": "86659b1a8e60fc398e93bcaf62691549",
          "model": "opus",
          "effort": "high",
          "source": "eadc1b23a9247fef6e2ef3bd567eff24050b6726",
          "native_ref": "claude-agent:ft393_t6_repair1",
          "predecessor": "claude:ft393_t6",
          "trigger": "user-directed",
          "stopped": "the ticket 06 author reached the context limit and has no live child",
          "preserved": "f932db29f82ba07271a8e11c8ac5ea011b0aa3ab"
        }
      ],
      "07-route-each-reset-refusal-through-the-registry.md": [
        {
          "session": "claude:ft393_t7",
          "assignment": "86659b1a8e60fc398e93bcaf62691549",
          "model": "opus",
          "effort": "high",
          "source": "4812b3923e4c247a28e9e09f743c0c404515bdf6",
          "native_ref": "claude-agent:ft393_t7"
        },
        {
          "session": "claude:ft393_t7b",
          "assignment": "86659b1a8e60fc398e93bcaf62691549",
          "model": "opus",
          "effort": "high",
          "source": "a54fe8fca1910d7834bbc6c47c65d09ba6318e81",
          "native_ref": "claude-agent:ft393_t7b",
          "predecessor": "claude:ft393_t7",
          "trigger": "user-directed",
          "stopped": "the ticket 07 author session returned a blocked report at about 318k tokens of context and has no live child; the reviewer directed a fresh author",
          "preserved": "a54fe8fca1910d7834bbc6c47c65d09ba6318e81"
        },
        {
          "session": "claude:ft393_t7_repair1",
          "assignment": "86659b1a8e60fc398e93bcaf62691549",
          "model": "opus",
          "effort": "high",
          "source": "e686bbda27b1db9256b38d1eb8502ecefc21b2a2",
          "native_ref": "claude-agent:ft393_t7_repair1",
          "predecessor": "claude:ft393_t7b",
          "trigger": "user-directed",
          "stopped": "the ticket 07 successor reported its RR-C2 records and has no live child",
          "preserved": "f932db29f82ba07271a8e11c8ac5ea011b0aa3ab"
        },
        {
          "session": "claude:ft393_t7_repair1b",
          "assignment": "86659b1a8e60fc398e93bcaf62691549",
          "model": "opus",
          "effort": "high",
          "source": "d0c8c4955b07af76bcbc4e007a9b94c3714e6e96",
          "native_ref": "claude-agent:ft393_t7_repair1b",
          "predecessor": "claude:ft393_t7_repair1",
          "trigger": "user-directed",
          "stopped": "the ticket 07 repair session neared the 300k context rule after its commit and has no live child",
          "preserved": "d0c8c4955b07af76bcbc4e007a9b94c3714e6e96"
        }
      ],
      "08-route-the-commit-exit-3-to-the-reset-plan.md": [
        {
          "session": "claude:ft393_t8",
          "assignment": "86659b1a8e60fc398e93bcaf62691549",
          "model": "opus",
          "effort": "high",
          "source": "f97dcb703cb5e304115e9b475385b6aae99a405b",
          "native_ref": "claude-agent:ft393_t8"
        },
        {
          "session": "claude:ft393_t8_repair1",
          "assignment": "86659b1a8e60fc398e93bcaf62691549",
          "model": "opus",
          "effort": "high",
          "source": "ba5d8921e4de62d64f4046ad1fef241f51b446b8",
          "native_ref": "claude-agent:ft393_t8_repair1",
          "predecessor": "claude:ft393_t8",
          "trigger": "user-directed",
          "stopped": "the ticket 08 author reported its RR-C3 records and has no live child",
          "preserved": "ba5d8921e4de62d64f4046ad1fef241f51b446b8"
        }
      ],
      "09-route-each-commit-refusal-through-the-registry.md": [
        {
          "session": "claude:ft393_t9",
          "assignment": "86659b1a8e60fc398e93bcaf62691549",
          "model": "opus",
          "effort": "high",
          "source": "a62b2929153f8a1cf00a3c7cccba2b2fece0d861",
          "native_ref": "claude-agent:ft393_t9"
        },
        {
          "session": "claude:ft393_t9b",
          "assignment": "86659b1a8e60fc398e93bcaf62691549",
          "model": "opus",
          "effort": "high",
          "source": "7d86d4697932a09f683f0ecf33ecacadda1c10b6",
          "native_ref": "claude-agent:ft393_t9b",
          "predecessor": "claude:ft393_t9",
          "trigger": "user-directed",
          "stopped": "the ticket 09 author neared the 300k context rule after its commit and has no live child",
          "preserved": "7d86d4697932a09f683f0ecf33ecacadda1c10b6"
        },
        {
          "session": "claude:ft393_t9_repair1",
          "assignment": "86659b1a8e60fc398e93bcaf62691549",
          "model": "opus",
          "effort": "high",
          "source": "12c022023b8fe3062920e03adfdba79bb0dd30a1",
          "native_ref": "claude-agent:ft393_t9_repair1",
          "predecessor": "claude:ft393_t9b",
          "trigger": "user-directed",
          "stopped": "the ticket 09 verification session reported its RR-C3 records and has no live child",
          "preserved": "ba5d8921e4de62d64f4046ad1fef241f51b446b8"
        },
        {
          "session": "claude:ft393_t9_repair2",
          "assignment": "86659b1a8e60fc398e93bcaf62691549",
          "model": "opus",
          "effort": "high",
          "source": "4ddf054e5d72a88d3c75df034de47e7629f9e07f",
          "native_ref": "claude-agent:ft393_t9_repair2",
          "predecessor": "claude:ft393_t9_repair1",
          "trigger": "user-directed",
          "stopped": "the ticket 09 repair session reported its RR-C3 records and has no live child",
          "preserved": "eb02867a4f5d139439fa9532ff05ae5cb889a537"
        }
      ],
      "10-route-each-checkpoint-refusal-by-its-cause.md": [
        {
          "session": "claude:ft393_t10",
          "assignment": "86659b1a8e60fc398e93bcaf62691549",
          "model": "opus",
          "effort": "high",
          "source": "4abc3cc84ba09b72fc56796fa2e20ea42ea4dbcb",
          "native_ref": "claude-agent:ft393_t10"
        },
        {
          "session": "claude:ft393_t10b",
          "assignment": "86659b1a8e60fc398e93bcaf62691549",
          "model": "opus",
          "effort": "high",
          "source": "c53a2da583101247595f4d6806f6df5d68be9609",
          "native_ref": "claude-agent:ft393_t10b",
          "predecessor": "claude:ft393_t10",
          "trigger": "user-directed",
          "stopped": "the ticket 10 author passed the 300k context rule after its commit and has no live child",
          "preserved": "c53a2da583101247595f4d6806f6df5d68be9609"
        },
        {
          "session": "claude:ft393_t10_repair1",
          "assignment": "86659b1a8e60fc398e93bcaf62691549",
          "model": "opus",
          "effort": "high",
          "source": "61952037e487cb3aeea87b5055d4ac0193fb824a",
          "native_ref": "claude-agent:ft393_t10_repair1",
          "predecessor": "claude:ft393_t10b",
          "trigger": "user-directed",
          "stopped": "the ticket 10 verification session reported its RR-C4 records and has no live child",
          "preserved": "429a3d8af4ac5aa227cd705085ef138c58247aac"
        }
      ],
      "11-route-the-commitment-policy-refusals-through-faces.md": [
        {
          "session": "claude:ft393_t11",
          "assignment": "86659b1a8e60fc398e93bcaf62691549",
          "model": "opus",
          "effort": "high",
          "source": "e2991ccf285f97e8bdc05f30eb2c7419ce8b63fb",
          "native_ref": "claude-agent:ft393_t11"
        }
      ],
      "12-print-the-commitment-verb-routes-from-the-registry.md": [
        {
          "session": "claude:ft393_t12",
          "assignment": "86659b1a8e60fc398e93bcaf62691549",
          "model": "opus",
          "effort": "high",
          "source": "4cf6e4721b89fbcc06fdff9254442f3d95f18434",
          "native_ref": "claude-agent:ft393_t12"
        }
      ],
      "13-refuse-a-route-literal-outside-the-registry.md": []
    }
  }
}
```

## Testing decisions

The external behavior is the exit code and the printed route of each refused verb.
Each verb's tests drive the verb in process, as the landing registry test does today.
The guard check and the bypass check are conformance tests over the compiled registry and the source tree.

The precedent is `TestLandingRefusalRegistryHasAProducingFixture`.
It walks the registry, requires one producing fixture per face, and drives each fixture through the verb.
Each verb gets the same walk over the faces of its verb.
Each fixture also declares a follow step: it carries out the printed route, runs each printed agent command step verbatim, and reruns the refused command.

The three collision repros become ordinary tests with no build tag.
`TestCollisionRedSourceHasNoExitWhenMainMoves` becomes the merge red-source rows.
`TestCollisionGuardDeniesThePrintedConflictRepair` becomes the conflict authority rows and the guard check.
`TestCollisionGuardDeniesThePrintedRestore` becomes the commit exit 3 rows.

Each new expectation derives from the fixture inputs and the registry's declared faces, never from the renderer under test.

### Seam diagram

    trigger: a write verb refuses (commit, merge, reset, land, gate checkpoint, commitment)
        │
        ▼
    refusal site ──▶ [ refusalroute: face + facts → typed refusal ] ──▶ verb printer ──▶ next= / next[1] cell
                          │
                          ├── guard check: each agent step → gitguard, benchguard ◀ conformance
                          ├── bypass check: route literals outside the registry ◀ conformance
                          └── bench recovery ──▶ recovery[N]{verb,face,authority,route}
        ◀ tests attach at each verb's in-process runner: produce the face, follow the route, rerun

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| RR01 | 1, 2 | Every registered face declares a verb from the six write verbs | planned TestRegistryFacesAreComplete in internal/refusalroute/registry_test.go | A face with an empty or unknown verb escapes every per-verb walk |
| RR02 | 2 | Every registered face declares the authority `agent` or `reviewer` | planned TestRegistryFacesAreComplete in internal/refusalroute/registry_test.go | A face with no authority leaves the agent to guess |
| RR03 | 5 | Every registered face declares at least one route step | planned TestRegistryFacesAreComplete in internal/refusalroute/registry_test.go | A face with no step prints an empty route |
| RR04 | 6 | The registry walk fails when two faces share one name | planned TestRegistryFaceNamesAreUnique in internal/refusalroute/registry_test.go | A duplicate name lets the lookup return the wrong face |
| RR05 | 3 | A rendered reviewer route starts with `reviewer: ` | planned TestRouteRendering in internal/refusalroute/route_test.go | A reviewer route without the marker reads as an agent command |
| RR06 | 4 | A rendered agent route joins its steps with `; then ` in declared order | planned TestRouteRendering in internal/refusalroute/route_test.go | A renderer that sorts or drops steps prints a route in the wrong order |
| RR07 | 7 | The constructor for an unknown face name returns `refusal face <name> is unregistered` with a reviewer route | planned TestUnregisteredFaceRefusesWithAReviewerRoute in internal/refusalroute/registry_test.go | A panic or an empty route aborts the operator's session |
| RR08 | 8 | Every agent route, rendered with sample facts, gets the empty label from `gitguard.Classify` for each command step | planned TestAgentRoutesPassTheWiredGuards in internal/conformance/refusal_route_guard_test.go | A printed `git merge` or `git restore` path step passes unseen |
| RR09 | 9 | Every agent route command step, with every slot filled by a sample value, gets `Blocked == false` from `benchguard.Classify` | planned TestAgentRoutesPassTheWiredGuards in internal/conformance/refusal_route_guard_test.go | A step that chains a second Bench call passes unseen |
| RR10 | 9 | Every agent route command step gets the empty string from `benchguard.PoolReference` with a pool-path sample | planned TestAgentRoutesPassTheWiredGuards in internal/conformance/refusal_route_guard_test.go | A `cd` or `git -C` into the pool passes unseen |
| RR11 | 10 | The guard check refuses an agent step that runs `bench worktree exec <target> -- bench <verb>` | planned TestAgentRoutesPassTheWiredGuards in internal/conformance/refusal_route_guard_test.go | An exec route survives today and breaks when FT341 refuses it |
| RR12 | 11 | The guard check reds an injected agent face whose route step is `git merge <commit>` | planned TestAgentRouteGuardCheckBites in internal/conformance/refusal_route_guard_test.go | A check that never classifies passes every route |
| RR13 | 12 | A route step that runs a tree-scoped Bench verb at a worktree renders `--in <label>` | planned TestRouteRendering in internal/refusalroute/route_test.go | A renderer that keeps the exec form fails RR11 at the first real face |
| RR14 | 13, 15 | Each land face in the registry has exactly one producing fixture, and each fixture produces a registered land face | `internal/worktree/identity_component_test.go` (`TestLandingRefusalRegistryHasAProducingFixture`) | A face moved without its fixture reaches an operator unproven |
| RR15 | 14 | Each land fixture follows the printed route, reruns the landing, and the face's sentence no longer prints | planned TestLandingFacesFollowTheirRoutes in internal/worktree/refusal_route_follow_test.go | A route that names the wrong repair leaves the face in place |
| RR16 | 16 | The landing composition-conflict refusal prints a `next=` value that starts with `reviewer: ` | planned TestConflictRepairIsAReviewerRoute in internal/worktree/refusal_route_test.go | An agent route that prints `git merge` is the collision 8 denial |
| RR17 | 16 | The landing composition-conflict-pending refusal prints a `next=` value that starts with `reviewer: ` | planned TestConflictRepairIsAReviewerRoute in internal/worktree/refusal_route_test.go | The pending-merge arm prints `git merge --continue`, which the guard denies |
| RR18 | 17 | The landing destination-not-clean refusal prints a `next=` value that starts with `reviewer: ` | planned TestReviewerLandFacesOpenWithTheMarker in internal/worktree/refusal_route_test.go | An agent route lets the agent discard the reviewer's primary-checkout work |
| RR19 | 18 | The landing source-not-clean refusal prints a route that contains `bench commit --in ` and the source label | planned TestLandingFacesFollowTheirRoutes in internal/worktree/refusal_route_follow_test.go | A route that names no commit command leaves the agent at a raw commit, which the guard denies |
| RR20 | 19 | Each landing preflight route ends with the caller's own re-run | `internal/worktree/land_surface_test.go` (`TestLandCommandReportsEveryRefusalInOnePreflight`) | A face that drops the re-run leaves a second lookup |
| RR21 | 20 | A fold of `main` into a target whose committed tip fails its lane prints a `refused{` record that contains `next=` | planned TestRedSourceFoldNamesAnExit in internal/worktree/merge_route_test.go | This is the collision 5a repro: the `lane fail` refusal prints no `next=` today |
| RR22 | 21 | The target-red fold route contains `bench commit --in ` and the target label | planned TestRedSourceFoldNamesAnExit in internal/worktree/merge_route_test.go | A route that reruns the fold alone loops on the same red |
| RR23 | 21 | The target-red fold route ends with `bench worktree merge --from ` and the incoming spelling and the target | planned TestRedSourceFoldNamesAnExit in internal/worktree/merge_route_test.go | A route that stops at the commit leaves the fold undone |
| RR24 | 21 | After the fixture removes the red file and runs the printed commit step, the printed fold step exits 0 | planned TestRedSourceFoldNamesAnExit in internal/worktree/merge_route_test.go | A route whose steps do not clear the red keeps the deadlock |
| RR25 | 22 | A fold into a dirty target prints a `refused{` record whose `next=` contains `bench commit --in ` | planned TestRedSourceFoldNamesAnExit in internal/worktree/merge_route_test.go | This is the second refusal of the collision 5a repro, with no `next=` today |
| RR26 | 23 | A `lane fail` fold whose incoming commit adds the red to a lane-green target prints a `next=` value that starts with `reviewer: ` | planned TestMergeFacesFollowTheirRoutes in internal/worktree/merge_route_test.go | A merge that never grades the target alone gives every lane red the agent route, which loops |
| RR55 | 23 | A `lane fail` fold whose target tip alone fails the lane prints a `next=` value that does not start with `reviewer: ` | planned TestRedSourceFoldNamesAnExit in internal/worktree/merge_route_test.go | A merge that sends every lane red to the reviewer hands back a repair that is the agent's |
| RR57 | 23 | An `inherited` fold whose target tip alone grades green prints a `next=` value that starts with `reviewer: ` | planned TestMergeFacesFollowTheirRoutes in internal/worktree/merge_route_test.go | A face map that reads `inherited` as a target red gives a fold red the agent route, which loops |
| RR60 | 10 | An incomplete landing whose source path is not line-safe prints a resume route that contains `<checkout>` and no `bench worktree exec` | planned TestUnsafePathRouteUsesThePlaceholder in internal/worktree/land_release_refusal_test.go | The current resume pointer form runs a Bench child through exec, which FT341 refuses |
| RR58 | 10 | A landing refusal whose source path is not line-safe prints an agent route that contains `<checkout>` and no `bench worktree exec` | planned TestUnsafePathRouteUsesThePlaceholder in internal/worktree/land_release_refusal_test.go | The current pointer form runs a Bench child through exec, which FT341 refuses, and the sample-only guard check misses it |
| RR59 | 23 | An `inherited` fold whose target tip alone grades red prints a `next=` value that contains `bench commit --in ` and the target label | planned TestRedSourceFoldNamesAnExit in internal/worktree/merge_route_test.go | A merge that grades the target only on `lane fail` sends a lane-less project's target red to the reviewer, and the collision 5a deadlock returns |
| RR56 | 23 | A `candidate` gate-kind fold red prints a `next=` value that starts with `reviewer: ` | planned TestMergeFacesFollowTheirRoutes in internal/worktree/merge_route_test.go | A face map that ignores the gate kind's attribution gives a fold red the agent route |
| RR27 | 24 | A merge composition conflict prints a `next=` value that starts with `reviewer: ` | planned TestConflictRepairIsAReviewerRoute in internal/worktree/refusal_route_test.go | The merge keeps printing the guard-denied `git merge` as an agent step |
| RR28 | 25 | A fold with a dirty sibling prints a route that contains `bench commit --in ` and no `bench worktree exec` | planned TestMergeFacesFollowTheirRoutes in internal/worktree/merge_route_test.go | The current route prints the exec form that FT341 refuses |
| RR29 | 13, 27 | Each merge face and each reset face has exactly one producing fixture for each declared cause, and each fixture follows its route out of the face | planned TestMergeFacesFollowTheirRoutes in internal/worktree/merge_route_test.go | A merge or reset face with no fixture reaches an operator unproven |
| RR30 | 26 | The `inherited` authorization refusal sentence equals `prospective authorization refused: inherited (the gate ran red on the composed tree and no green baseline attributes the red to this diff)` | `internal/landing/landing_reviewed_test.go` (`TestRefusalMessageNamesTheOperatorActionAndTheOpenReason`) | An inline action prints a second route beside the face's route |
| RR67 | 26 | A landing whose composed tree grades red prints a `refused{` record whose `next=` contains `bench commit --in ` and ends with the caller's re-run | planned TestLandingRedRouteNamesTheRepair in internal/worktree/refusal_route_test.go | After the sentence drops its action, a landing red prints no route |
| RR31 | 26 | The merge retries an empty-reason infrastructure refusal exactly once | `internal/worktree/merge_test.go` (`TestMergeRetriesOnlyAVerifiedEmptyReasonInfrastructureRefusal`) | A retry that matches the old sentence never fires after the sentence changes |
| RR32 | 28 | A commit exit 3 prints `next=` with the value `bench worktree reset --to <published-commit> <checkout>` for its published commit and its root | planned TestPublishedUnreconciledRouteIsTheResetPlan in internal/commit/refusal_route_test.go | This is the collision 8 restore repro: the current route is `git restore`, which the guard denies |
| RR33 | 28 | A commit exit 3 `next=` value does not contain `git restore` | planned TestPublishedUnreconciledRouteIsTheResetPlan in internal/commit/refusal_route_test.go | A route that keeps the restore beside the reset still prints a denied step |
| RR34 | 28 | A commit exit 3 on a root path that is not line-safe prints the placeholder `<checkout>` | planned TestPublishedUnreconciledRouteIsTheResetPlan in internal/commit/refusal_route_test.go | A raw control byte reaches the line-structured record |
| RR35 | 29 | After a commit exit 3 in an assignment worktree, the printed reset plan and its `--apply` leave `git status --porcelain` empty at the published commit | planned TestCommitExitThreeRouteReconcilesTheCheckout in internal/worktree/commit_route_test.go | A reset route that does not reconcile leaves the exit 3 state |
| RR36 | 30 | Each commit face has exactly one producing fixture for each declared cause, and each fixture follows its route out of the face | planned TestCommitFacesFollowTheirRoutes in internal/commit/refusal_route_test.go | A commit face with no fixture reaches an operator unproven |
| RR37 | 30 | A red commit refusal prints a `next=` line on stderr whose route ends with the caller's commit command | planned TestCommitFacesFollowTheirRoutes in internal/commit/refusal_route_test.go | The current `inherited` refusal prints only `run bench gate --fresh`, which re-grades the same red |
| RR38 | 30 | A commit in the primary checkout prints a `next=` route that contains `bench worktree create --request` | planned TestCommitFacesFollowTheirRoutes in internal/commit/refusal_route_test.go | A refusal that keeps its route only inside the sentence bypasses the registry |
| RR39 | 31 | A checkpoint refusal for a missing completion record prints no `help[1]{cmd,why}` row | `internal/gate/review_checkpoint_test.go` (`TestReviewCheckpointRefusalRoute`) | This is the FT330 defect: the fixed write-access row prints for every cause |
| RR40 | 32 | A checkpoint refusal for a missing completion record prints `next=bench preflight review 'example'` | `internal/gate/review_checkpoint_test.go` (`TestReviewCheckpointRefusalRoute`) | The move to the registry drops the current evidence route |
| RR41 | 33 | A complete checkpoint with an uncommitted tracked edit prints `cleanCheckoutRefusal` and a `next=` route that contains `bench commit --in ` | `internal/gate/complete_checkpoint_test.go` (`TestCompleteCheckpointRefusesADirtyCheckout`) | The landed refusal names its recovery only in prose and prints no `next=` |
| RR68 | 33 | A complete checkpoint on a dirty assignment checkout prints a `next=` route that contains the assignment label after `--in ` | planned TestCheckpointFacesFollowTheirRoutes in internal/gate/refusal_route_test.go | A route that keeps the `<label>` placeholder where an assignment owns the root makes the agent look up its own label |
| RR66 | 31 | A complete checkpoint on a spec with no `Status: staged` line prints a `next=` value that starts with `reviewer: ` | `internal/gate/complete_checkpoint_test.go` (`TestCompleteCheckpointRefusesAnUntransformableSpec`) | The compose refusal prints through `operational` with no route today |
| RR42 | 34 | A subject-capture fault prints `next=` with a route that contains `bench doctor` and `--fresh` | `internal/gate/run_outcomes_test.go` (`TestGateRunRetainsSubjectConstructionCause`) | The fallback keeps the write-access text for a fault that is not a write-access fault |
| RR43 | 35 | A subject-capture fault prints no `help[1]{cmd,why}` row | `internal/gate/run_outcomes_test.go` (`TestGateRunRetainsSubjectConstructionCause`) | A route printed beside the old row gives two answers |
| RR44 | 13, 14 | Each gate face has exactly one producing fixture that follows its route out of the face | planned TestCheckpointFacesFollowTheirRoutes in internal/gate/refusal_route_test.go | A gate face with no fixture reaches an operator unproven |
| RR45 | 36 | A `bench commitment start` refusal for an outcome outside the active milestone prints a `next` cell that starts with `reviewer: ` | planned TestCommitmentFacesFollowTheirRoutes in internal/commitment/commitcmd/refusal_route_test.go | The current cell tells the agent to re-plan a commitment that only the reviewer changes |
| RR46 | 37 | A `bench commitment start` refusal without an owned assignment prints a `next` cell that contains `bench worktree create --request` | planned TestCommitmentFacesFollowTheirRoutes in internal/commitment/commitcmd/refusal_route_test.go | The current cell prints the plan command for a cause that the agent clears |
| RR47 | 37 | A `bench commitment verify` refusal keeps its `bench commitment verify --milestone` cell | `internal/commitment/verification_test.go` (`TestCommitmentUnmetCriterion`) | The move to the registry drops the current agent route |
| RR61 | 37 | A commit from an assignment with no delivery binding prints a `next=` route that contains `bench commitment start --outcome` | planned TestCommitFacesFollowTheirRoutes in internal/commit/refusal_route_test.go | The route stays only inside the sentence tail, outside the registry |
| RR62 | 36 | A commit whose candidate changes the protected commitment prints a `next=` value that starts with `reviewer: ` | planned TestCommitFacesFollowTheirRoutes in internal/commit/refusal_route_test.go | The tail tells the agent to re-plan a commitment that only the reviewer changes |
| RR63 | 37 | A light-path commit with a path outside its ticket's `Writes:` line prints an agent route that names the ticket and not `bench commitment start` | planned TestCommitFacesFollowTheirRoutes in internal/commit/refusal_route_test.go | The tail offers a commitment start for a fix that is one `Writes:` edit |
| RR64 | 26 | The `assignment has no current delivery binding` refusal sentence does not contain `run bench` | `internal/commitment/repository/publication_test.go` (`TestAdmitPublicationClosureAuthority`) | A sentence that keeps its tail prints a second route beside the face's route |
| RR65 | 39 | The bypass check reds a planted `; run bench` tail in an `internal/commitment/repository` production file | planned TestRouteBypassCheckBites in internal/conformance/refusal_route_bypass_test.go | A scan that skips the commitment packages lets a tail return |
| RR48 | 13, 14 | Each commitment face has exactly one producing fixture that follows its route out of the face | planned TestCommitmentFacesFollowTheirRoutes in internal/commitment/commitcmd/refusal_route_test.go | A commitment face with no fixture reaches an operator unproven |
| RR49 | 38 | The bypass check passes on the tree after every chunk lands | planned TestNoWriteVerbComposesARouteOutsideTheRegistry in internal/conformance/refusal_route_bypass_test.go | A hand-composed route left in a verb shows as a red |
| RR50 | 39 | The bypass check reds a planted `next=` literal, and a planted `next[1]:` literal, in a write-verb production file | planned TestRouteBypassCheckBites in internal/conformance/refusal_route_bypass_test.go | A check that scans no file passes every tree |
| RR51 | 40 | `bench recovery` prints `recovery[N]{verb,face,authority,route}` and exits 0 | planned TestRecoveryListsEveryFace in internal/refusalroute/command_test.go | A missing header breaks the TOON contract |
| RR52 | 41 | `bench recovery` prints exactly one row for each registered face, in registry order | planned TestRecoveryListsEveryFace in internal/refusalroute/command_test.go | A hand-kept matrix drifts from the registry |
| RR53 | 40 | `bench help` lists `bench recovery` | `cmd/bench/help_inventory_test.go` (`TestHelpInventoryIsComplete`) | An unlisted verb is a dead key |
| RR54 | 42 | The reference guide names `bench recovery` as the recovery matrix | review-owned | No gate check reads this prose |

Not covered: story 43 — FT342 owns the moved-`main` decision, and RR26 routes to it.
Not covered: story 44 — the exit-2 grammar refusals keep their behavior, and the usage package owns their usage line.

### Edge inventory

The hostile-input checklist applies to the facts that a route renders.
A path fact can hold a space, a glob byte, or a control byte.

An agent route for a value that is not line-safe prints that value's placeholder and never the exec pointer form.
For a path, the route prints `<checkout>` after a `bench worktree path <id>` step.
A missing tree is the exception: its `bench worktree path <id>` step refuses, so its route prints the command with `<checkout>` and no path step.
`landingResumeNext` and `atSourceWorktree` take this form, so the `land-incomplete` resume and each preflight re-run obey it.
The renderer shell-quotes a line-safe value and prints the slot placeholder for a value that is not line-safe; RR34 grades the placeholder.
A label fact comes from the assignment ledger, and the `--in` resolver refuses an ambiguous label under FT341.

The guard check runs over sample facts, and the producing fixtures run over real values.
A route can pass the check with samples and fail it with a real path.
So RR15 and the other follow rows run the real printed steps.

The fail-closed refusals are the unregistered face (RR07) and each verb's `<verb>-handback` face.

- **Won't handle**: the moved-`main` decision — FT342 owns it, and `merge-fold-red` routes to it.
- **Won't handle**: a fold that admits an inherited red — the decision source closes the merge policy, and `merge-target-red` routes the repair before the fold.
- **Won't handle**: the exit-2 grammar refusals — the usage package owns their usage line, and `usage.Parse` keeps it.
- **Won't handle**: the gate run-state refusals outside the checkpoint funnel — they serve every gate run, and `operational` keeps their reason line.
- **Won't handle**: the non-write worktree verbs (`list`, `path`, `clean`, `release`, `build`) — the decision source names six write verbs, and `recoveryRoute` keeps their routes.
- **Won't handle**: the exec arm of a reviewer conflict route for an unsafe path — the reviewer runs that step, and the check grades agent routes.
- **Won't handle**: a change to either wired guard — the closed decision forbids it, and the guard check reads the guards unchanged.

## Ownership fences

The prospective build owns these exact paths:

- `internal/refusalroute/registry.go`
- `internal/refusalroute/faces_land.go`
- `internal/refusalroute/faces_merge.go`
- `internal/refusalroute/faces_reset.go`
- `internal/refusalroute/faces_commit.go`
- `internal/refusalroute/faces_gate.go`
- `internal/refusalroute/faces_commitment.go`
- `internal/refusalroute/route.go`
- `internal/refusalroute/command.go`
- `internal/refusalroute/registry_test.go`
- `internal/refusalroute/route_test.go`
- `internal/refusalroute/command_test.go`
- `internal/conformance/refusal_route_guard_test.go`
- `internal/conformance/refusal_route_bypass_test.go`
- `internal/worktree/land_refusal.go`
- `internal/worktree/land.go`
- `internal/worktree/land_identity.go`
- `internal/worktree/land_resume.go`
- `internal/worktree/land_rerun.go`
- `internal/worktree/merge.go`
- `internal/worktree/merge_refusal.go`
- `internal/worktree/land_refusal_fixture_test.go`
- `internal/worktree/build.go`
- `internal/worktree/reset.go`
- `internal/worktree/reset_apply.go`
- `internal/worktree/reset_restore.go`
- `internal/worktree/ownership.go`
- `internal/worktree/reset_restore_refusal_test.go`
- `internal/worktree/reset_apply_test.go`
- `internal/worktree/missing_tree_recovery_test.go`
- `internal/worktree/merge_target_grade_test.go`
- `internal/worktree/merge_caller_root_test.go`
- `internal/worktree/list.go`
- `internal/worktree/list_actions_test.go`
- `internal/landing/landing_test.go`
- `internal/refusalroute/routetest/routetest.go`
- `internal/gate/completion.go`
- `internal/gate/engine.go`
- `internal/gate/authorization/authorization.go`
- `internal/commitment/repository/light_path_test.go`
- `internal/commit/commitment_test.go`
- `internal/worktree/commitment_light_landing_test.go`
- `internal/commit/assessment_span_test.go`
- `internal/commit/commit_test.go`
- `internal/commit/commitment_route_test.go`
- `internal/refusalroute/routetest/routetest_test.go`
- `internal/landing/attribution.go`
- `internal/landing/gitexec.go`
- `internal/worktree/lifecycle_test.go`
- `internal/worktree/worktree_test.go`
- `internal/worktree/lifecycle.go`
- `internal/worktree/worktree.go`
- `internal/worktree/clean_landed_test.go`
- `internal/worktree/clean_landed.go`
- `internal/worktree/path.go`
- `internal/worktree/classifier.go`
- `internal/worktree/identity_component.go`
- `internal/worktree/refusal_route_test.go`
- `internal/worktree/refusal_route_follow_test.go`
- `internal/worktree/merge_route_test.go`
- `internal/worktree/commit_route_test.go`
- `internal/worktree/identity_component_test.go`
- `internal/worktree/land_surface_test.go`
- `internal/worktree/land_tickets_only_test.go`
- `internal/worktree/land_release_refusal_test.go`
- `internal/worktree/land_folded_base_test.go`
- `internal/worktree/land_identity_test.go`
- `internal/worktree/land_resume_refusal_test.go`
- `internal/worktree/land_journey_test.go`
- `internal/worktree/parallel_census_test.go`
- `internal/worktree/land_resume_test.go`
- `internal/worktree/land_reauthorization_test.go`
- `internal/worktree/merge_test.go`
- `internal/landing/landing.go`
- `internal/landing/merge.go`
- `internal/landing/landing_reviewed_test.go`
- `internal/commit/commit.go`
- `internal/commit/landing_test.go`
- `internal/commit/dry_run_test.go`
- `internal/commit/refusal_route_test.go`
- `internal/gate/checkpoint.go`
- `internal/gate/run_transaction.go`
- `internal/gate/gate.go`
- `internal/gate/run_outcomes_test.go`
- `internal/gate/review_checkpoint_test.go`
- `internal/gate/refusal_route_test.go`
- `internal/commitment/commitcmd/command.go`
- `internal/commitment/commitcmd/admission.go`
- `internal/commitment/commitcmd/refusal_route_test.go`
- `internal/commitment/verification_test.go`
- `cmd/bench/main.go`
- `cmd/bench/commitment_test.go`
- `cmd/bench/help_inventory_test.go`
- `.bench/BENCH-reference.md`
- `CHANGELOG.md`
- `reviews/refusal-route-registry.md`
- `internal/gate/complete_checkpoint.go`
- `internal/gate/complete_checkpoint_test.go`
- `tests/canary/package-core-guard/unrouted-subcommand`
- `tests/canary/docs-currency-token-diet/benchref-imported`
- `tests/canary/docs-currency-token-diet/benchref-pointer-dropped`
- `tests/canary/docs-currency-token-diet/benchref-section-duplicated`
- `tests/canary/skills-index-command-adapters/adapter-inert-invocation-key`
- `tests/canary/skills-index-command-adapters/command-invocation-disabled-against-policy`
- `tests/canary/skills-index-command-adapters/dangling-index`
- `tests/canary/skills-index-command-adapters/debug-implicit-invocation-reverted`
- `tests/canary/skills-index-command-adapters/missing-index-field`
- `tests/canary/skills-index-command-adapters/stale-index-wording`
- `tests/canary/skills-index-command-adapters/unindexed-skill`
- `tests/canary/workflow-guidance-anchors/agents-handoff-section-rule`
- `tests/canary/workflow-guidance-anchors/reference-agent-push-rule`
- `tests/canary/workflow-guidance-anchors/reference-bench-operational-layer`
- `tests/canary/workflow-guidance-anchors/reference-category-context`
- `tests/canary/workflow-guidance-anchors/reference-category-oracle`
- `tests/canary/workflow-guidance-anchors/reference-category-setup`
- `tests/canary/workflow-guidance-anchors/reference-category-work`
- `tests/canary/workflow-guidance-anchors/reference-gate-authority`
- `tests/canary/workflow-guidance-anchors/reference-kit-only-ship`
- `tests/canary/workflow-guidance-anchors/reference-no-path-fallback`
- `tests/canary/workflow-guidance-anchors/reference-progressive-loading-term`
- `tests/canary/workflow-guidance-anchors/reference-refusal-route-shape`
- `tests/canary/workflow-guidance-anchors/reference-retro-capture-owner`
- `tests/canary/workflow-guidance-anchors/reference-retro-drain-owner`
- `tests/canary/workflow-guidance-anchors/reference-skills-guidance`
- `tests/canary/workflow-guidance-anchors/reference-upgrade-route`
- `tests/canary/workflow-guidance-anchors/changelog-reduced-schema-columns`
- `tests/canary/workflow-guidance-anchors/changelog-ticket-vocabulary`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry_data_test.go`
- `internal/anchors/registry_retained_workflow.go`
- `internal/worktree/merge_from_sha_test.go`
- `internal/worktree/land_bench_home_test.go`
- `internal/commitment/delivery.go`
- `internal/commitment/repository/admission.go`
- `internal/commitment/repository/candidate.go`
- `internal/commitment/repository/continuation.go`
- `internal/commitment/repository/publication.go`
- `internal/commitment/repository/light_path.go`
- `internal/commitment/repository/verification.go`
- `internal/commitment/repository/readiness.go`
- `internal/commitment/repository/publication_test.go`

This fence is the closure that the build preflight write proposal confirmed on the draft Writes line.
It also holds two rendered-text readers that the reader sweep found.
The two FT392 gate paths joined after the first proposal, and a second proposal at the folded tip listed no further path.
After the fold of main at 120e2715, a third proposal listed no further path.
The review pickup is `reviews/refusal-route-registry.md`.
The planning author owns this spec folder during the spec phase.

## Out of scope

- The gate run-state refusals outside the checkpoint funnel: 8 edits, 2 gate runs.
- The routes of the non-write worktree verbs: 6 edits, 1 gate run.

These are planning estimates, not approved priorities.

## Further notes

Evidence status: source-backed proposal. No probe ran, and the phase commits no code.

### Flagged additions

- The `bench recovery` read verb. The decision source requires that the matrix renders from the registry, and a read verb renders it with no second copy. A generated Markdown file with a drift check is the alternative.
- The `--in` rule for route steps (RR11, RR13). FT341 owns the `--in` forms, and its closed decision refuses an exec route for a Bench child.
- The static bypass check (RR49, RR50). The decision source requires proof that every refusal goes through the registry; the per-face fixtures prove only the faces that exist.
- The reviewer authority of the conflict faces (RR16, RR17, RR27). The guard states that the merge belongs to the reviewer, so no agent route exists for a conflict.
- The reset route for the commit exit 3 (RR32-RR35).
- The `merge-published-unreconciled` face. Its producing fixture is `TestMergeExitsThreeWhenTheReconcileFails` in `internal/worktree/merge_test.go`.
- The `land-incomplete` face. Its producing fixture is the interrupted landing in `internal/worktree/land_effects_test.go`.
- The registry's ownership of the `next` label and the commitment `next` table. Without it, the bypass check cannot tell a hand-spelled route from a registry route.
- The target-alone grade on an `inherited` or a `lane fail` fold (RR26, RR55, RR57). Neither kind tells a target red from a fold red, so without this grade a fold red cannot pick its authority.

### Posture change

- `TestPublicationRemainderNextRestoresEveryNamedPathShellQuoted` and `TestPublicationRemainderSanitizesTheFailedPathAndPointsAtNamedPaths` pin the `git restore` route, and RR32-RR34 replace both pins.
- `TestGateRunRetainsSubjectConstructionCause` pins the fixed write-access help row, and RR42-RR43 replace that pin.
- `TestReviewCheckpointRefusalRoute` gains the RR39 assertion that no `help[1]{cmd,why}` row prints.
- The commit help line `exit 3: published; the checkout did not reconcile — paste next= to repair` composes its label from the registry; no test pins that line today.
- The landing tests that compare a route with `landingRefusalFaceByName(...).route(...)` read the shared registry instead.
- The tests in `internal/landing`, `internal/commit`, and `internal/worktree` that read `run bench gate --fresh`, `fix the failures above`, or `run bench doctor` in a refusal sentence read the face route instead.

### Source-sentence-to-row table

| FT393 sentence | rows |
|---|---|
| A red source could not fold `main`, because the merge refused the inherited red and named no next action. | RR21-RR24, RR26, RR55-RR57, RR59 |
| A dirty fold also refused with no next action. | RR25 |
| The destructive-git guard denies a raw `git merge` after a conflict. | RR16, RR17, RR27 |
| The destructive-git guard denies `git restore --worktree` after a `bench commit` exit 3. | RR32-RR35 |
| One refusal-route registry serves every write verb. | RR01, RR14, RR29, RR36, RR44, RR48, RR49 |
| Each face declares a typed route and its authority, agent or reviewer. | RR02, RR03, RR05, RR06 |
| Every refusal goes through the registry. | RR49, RR50, RR61-RR65 |
| Every agent route passes the wired guards. | RR08-RR12, RR58, RR60 |
| Each producing fixture follows its route out of the refused face. | RR15, RR24, RR29, RR35, RR36, RR44, RR48 |
| The recovery matrix renders from the registry. | RR51-RR54 |
| A checkpoint refusal names the recovery route of its cause (FT330). | RR39-RR43, RR66, RR68 |

### Pre-review proof checklist

- `Cited symbols`: `landingRefusalFaces`, `landingFaceRefusalOf`, `conflictRepairPrefix`, `conflictContinuePrefix`, `restoreNext`, `publicationRemainder`, `refusalMessage`, `refusalGuidance`, `routedRefusal`, `subjectUnavailableHelp`, `refusalNext`, `checkoutClean`, `mergeReconcileNext`, `retryEmptyReasonInfrastructureFold`, `gitguard.Classify`, `benchguard.Classify`, and `benchguard.PoolReference` resolve in the tree at 120e2715.
- `Import edges`: each write-verb package to `internal/refusalroute` is new, and `internal/refusalroute` imports none of them; `internal/conformance` to `internal/gitguard` is new, and both packages are leaves.
- `Source-row clauses and occurrences`: two occurrences, the 2026-09-18 FT316 repro and the 2026-10-07 ft362 planning branch.
- `Promised field labels`: `recovery[N]{verb,face,authority,route}`, `next=`, and `reviewer: `.
- `Changed-function callers`: `conflictRepairPrefix` has `merge.go`, `land_refusal.go`, and `build.go`; `refusalMessage` has `landing.go`, `merge.go`, and `landing_reviewed_test.go`; `landingRefusalFaceByName` has eight worktree test files and `land_identity.go`.
- `Copy survival`: RR49 reds a hand-composed route that survives in a write verb.
- `Rendered-shape readers`: see the reader inventory.
- `Pin operators`: RR16, RR17, RR18, RR26, RR27, and RR45 use a prefix match; RR32 uses `==`; the other route rows use substring containment.
- `Entry reads`: each verb reads its own refusal state; the registry reads no ambient state.
- `Derived expectations`: the walks derive the face set from the registry; each route expectation comes from the fixture inputs.
- `Consolidated rules`: the per-verb route composition consolidates into the registry; each old route maps to its face in the face inventory table.
- `Quantified obligations`: every face (RR14, RR29, RR36, RR44, RR48), every agent route (RR08-RR11), and every write-verb literal (RR49).
- `Workflow-step writes`: none.

### Reader inventory

- `internal/worktree/land_refusal.go`: the landing registry, the conflict prefixes, and the resume route.
- `internal/worktree/merge.go`: the conflict, dirty, sibling, and reconcile routes.
- `internal/worktree/build.go`: a second caller of `conflictRepairPrefix`.
- `internal/worktree/reset.go` and `internal/worktree/reset_apply.go`: the reset routes.
- `internal/commit/commit.go`: `restoreNext` and the stderr refusal lines.
- `internal/landing/landing.go` and `internal/landing/merge.go`: `refusalMessage` and its guidance.
- `internal/gate/checkpoint.go`, `internal/gate/run_transaction.go`, and `internal/gate/gate.go`: the routed refusal, the funnel, and the fixed help row.
- `internal/commitment/commitcmd/command.go`: `refusal` and `refusalNext`.
- Tests that read the moved literals: `internal/commit/landing_test.go`, `internal/commit/dry_run_test.go`, `internal/landing/landing_reviewed_test.go`, `internal/worktree/merge_test.go`, `internal/worktree/merge_from_sha_test.go`, `internal/worktree/land_bench_home_test.go`, `internal/worktree/land_surface_test.go`, `internal/worktree/land_journey_test.go`, `internal/gate/run_outcomes_test.go`, and `internal/gate/review_checkpoint_test.go`.

### Primary sources

- `internal/worktree/land_refusal.go:187-340`: the landing registry and its one constructor.
- `internal/worktree/identity_component_test.go:200-360`: the producing-fixture walk.
- `internal/worktree/classifier.go:23-65`: the `refusal` record and its fields.
- `internal/worktree/merge.go:66-161`: the merge refusals and the conflict route.
- `internal/commit/commit.go:196-238`: the exit 3 record and `restoreNext`.
- `internal/landing/landing.go:283-315`: the authorization sentence and its inline action.
- `internal/gate/run_transaction.go:42-58` and `internal/gate/gate.go:338-343`: the funnel and the fixed help row.
- `internal/gate/complete_checkpoint.go:11-40`: `cleanCheckoutRefusal` and the three complete-checkpoint refusals.
- `internal/landing/landing.go:251-252`: the landing's own authorization red.
- `internal/gate/checkpoint.go:87-109`: the completion-evidence route.
- `internal/commitment/commitcmd/command.go:336-344`: the commitment refusal table.
- `internal/gitguard/gitguard.go:126-150` and `internal/benchguard/benchguard.go:101-120`, `internal/benchguard/pool.go:10-40`: the guard decisions.
- `internal/worktree/reset.go:28-240`: the reset plan keeps a dirty checkout in a recoverable envelope.

Cheap-tier research delegates read these sources and reported citations.
The author re-read each definition above before the rows locked.
