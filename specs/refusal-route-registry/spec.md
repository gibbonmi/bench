# Every write-verb refusal names a typed recovery route

Status: staged

Roadmap: FT393

Decision source: named reviewed artifact `roadmap/FT393.md`, opened by drain d-6259cc0d8421, which folded FT330 into it, and placed second in quality-1 by the reviewer on 2026-10-07

Audience: every repository that links the kit

Verification log: pending

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
Harder chunks: RR-C1, RR-C3.

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

15. As an orchestrator, I want each landing face to keep its current route, except where its authority changes, so that the move adds no behavior.
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

The one constructor takes the face name and the facts.
It returns a typed refusal that carries the face, the sentence, the paths, and the rendered route.
An unknown face name returns the refusal `refusal face <name> is unregistered`, with a reviewer route.
This keeps the landing's fail-soft rule for a bookkeeping fault.

### Authority

A face has agent authority when an agent can clear its cause with Bench verbs and its own file edits, inside its own worktree.
A face has reviewer authority when the clear needs one of these:

- a raw merge
- a change to the primary checkout
- a commitment change
- a decision that FT342 owns

The destructive-git guard already states that the merge and any history rewrite are the reviewer's, so the conflict faces take reviewer authority.

### Output shape

Each verb keeps its output shape.
The landing, the merge, and the reset print `next=<route>` in the `refused{...}` record.
The commit prints `next=<route>` on its stderr refusal line, and its exit 3 record keeps `committed{published_commit=…,path=…,next=…}`.
The gate checkpoint prints `next=<route>` on stderr after its reason, and it no longer prints the fixed `help[1]{cmd,why}` row.
The commitment verb keeps its `next[1]{command}` table, and the cell holds the face's route.

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
| land | `source-not-clean` | agent | `bench commit --in <label> -m <msg> -- <path>...`; then the re-run |
| land | `source-not-fenced` | agent | the current fence instruction; then the re-run |
| land | `composition-conflict` | reviewer | the hand merge of the destination commit; then `bench commit`; then `/bench-review-implementation`; then the re-run |
| land | `composition-conflict-pending` | reviewer | finish the merge in progress; then `/bench-review-implementation`; then the re-run |
| land | `resume-destination-residue` | reviewer | commit or discard the destination's uncommitted work; then the resume |
| land | `resume-marker` | agent | `bench gate`; then the resume |
| merge | `merge-target-red` | agent | repair each failing check in `<label>`; then `bench commit --in <label> -m <msg> -- <path>...`; then the re-run |
| merge | `merge-fold-red` | reviewer | the fold of `<from>` reds `<label>`; decide the moved-`main` route under FT342 |
| merge | `merge-target-not-clean` | agent | `bench commit --in <label> -m <msg> -- <path>...`; then the re-run |
| merge | `merge-sibling-not-clean` | agent | `bench commit --in <sibling-label> -m <msg> -- <path>...`; then the re-run |
| merge | `merge-conflict` | reviewer | the hand merge of the incoming commit; then `bench commit`; then the re-run |
| merge | `merge-infrastructure` | agent | `bench doctor`; then the re-run |
| reset | `reset-checkout-conflicted` | agent | `bench worktree clean <id>` |
| reset | `reset-plan-stale` | agent | the current plan command |
| reset | `reset-tree-missing` | agent | the current missing-tree route |
| commit | `commit-published-unreconciled` | agent | `bench worktree reset --to <published-commit> <checkout>` |
| commit | `commit-primary-checkout` | agent | `bench worktree create --request <opaque-id> --label <work-item>` |
| commit | `commit-red` | agent | repair each failure that the run reports; then the re-run |
| commit | `commit-infrastructure` | agent | `bench doctor`; then the re-run |
| gate | `checkpoint-completion-evidence` | agent | `bench preflight review <slug>` |
| gate | `checkpoint-dirty-checkout` | agent | `bench commit -m <msg> -- <path>...`; then the re-run |
| gate | `checkpoint-composition` | reviewer | the delivery closure of the spec does not compose; hand back |
| gate | `checkpoint-tip-moved` | agent | the re-run |
| gate | `checkpoint-subject-unavailable` | agent | `bench doctor`; then the re-run with `--fresh` |
| commitment | `commitment-needs-assignment` | agent | `bench worktree create --request <request> --label <label>` |
| commitment | `commitment-plan-input` | agent | correct the input file; then the re-run |
| commitment | `commitment-verify-evidence` | agent | `bench commitment verify --milestone <id> --evidence <file>` |
| commitment | `commitment-decision` | reviewer | `bench commitment plan --input <file>` |

The `<verb>-handback` reviewer faces join this list for each verb.
The `commitment-decision` face covers each start, block, plan, and approve refusal whose clear changes the active commitment.
The reset table keeps each current route of the reset verb unchanged.

### The red-source fold (collision 5a)

The decision source closes the merge policy: the gate stays the oracle, and a fold publishes only on a green or a lane pass.
So a fold red picks its face by its cause.
When the target tip alone is red, the `merge-target-red` face routes the agent to repair the target first.
The repair commit runs the lane on the repaired tree, and the fold then composes onto a green target.
When the target tip alone is green, the fold adds the red.
Then the `merge-fold-red` face routes to the reviewer, because FT342 owns the moved-`main` decision.

The gate kinds already carry this attribution: `inherited` takes `merge-target-red`, and `candidate` takes `merge-fold-red`.
The `lane fail` kind carries none, and the collision 5a repro refuses with it.
So on a `lane fail` fold, the merge grades the target tip alone under the same lane, before it picks the face.
That run is one extra lane run on the refusal path only, and it changes no publish rule.

### The commit exit 3 (collision 8b)

The commit publishes, and then the checkout fails to reconcile.
The route names the reset plan at the published commit: `bench worktree reset --to <published-commit> <checkout>`.
The reset verb keeps the dirty layer in a recoverable envelope and then prints its own `--apply` command.
`<checkout>` is the commit's root path, shell-quoted, or the placeholder `<checkout>` when the path is not line-safe.

### The conflict repair (collision 8a)

The conflict faces of the landing and the merge take reviewer authority.
Their route keeps the hand merge text that `conflictRepairPrefix` and `conflictContinuePrefix` print today, behind the `reviewer: ` prefix.
The guard check grades agent routes only, so a reviewer step that names `git merge` passes the check.

### The gate checkpoint (FT330)

Each cause in the pre-oracle funnel maps to its own face.
The completion-evidence face keeps the route that `routedRefusal` prints today.
The FT392 dirty-checkout refusal takes the `checkpoint-dirty-checkout` face.
A subject-capture fault that no other face claims takes `checkpoint-subject-unavailable`.
The checkpoint no longer prints the fixed write-access help row.

### The authorization sentence

`landing.refusalMessage` keeps its prefix, its kind, and its explanation.
It drops its inline action, because the face now carries the route.
The merge's single infrastructure retry matches the typed kind, not the old sentence.

### The FT392 dependency

This spec depends on FT392 landing first.
FT392 ticket 02 writes `internal/gate/gate.go`, `internal/gate/checkpoint.go`, `internal/gate/complete_checkpoint.go`, `.bench/BENCH-reference.md`, `CHANGELOG.md`, and the anchor registry files.
The gate chunk of this spec writes the same gate files, so the two fences overlap.
The build starts only on a base that contains the FT392 landing.

### The guard check

A conformance test renders each agent route with sample facts.
A path fact takes a sample under the pool prefix, so a `cd` or a `git -C` form turns red.
The test classifies each command step through `gitguard.Classify`, `benchguard.PoolReference`, and `benchguard.Classify`.
It also refuses a step that runs `bench worktree exec <target> -- bench ...`, under FT341's closed decision.
The production guards do not change.

### The static bypass check

A conformance test scans the string literals in the production Go files of the write-verb packages.
It refuses a literal that holds `next=`, a `next[` table header, or the route joiner `; then ` outside `internal/refusalroute`.
The test holds one reviewed allowlist of files whose routes serve a non-write verb, such as `bench worktree list`.

### The recovery matrix

`bench recovery` is a repository-scoped read verb.
It prints `recovery[N]{verb,face,authority,route}` with one row per registered face, in registry order.
It reads only the compiled registry.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| RR-C1 / to be sliced | The shared registry, the guard check, the landing faces on the registry, and `bench recovery` | RR01-RR20, RR51-RR54 | `internal/refusalroute`, `internal/conformance`, `internal/worktree` landing tests | yes |
| RR-C2 / to be sliced | The merge and the reset print registry routes, and the red-source fold has an exit | RR21-RR31, RR55, RR56 | `internal/worktree` merge and reset tests, `internal/landing` | no |
| RR-C3 / to be sliced | The commit prints registry routes, and the exit 3 route is the reset plan | RR32-RR38 | `internal/commit`, `internal/worktree` | yes |
| RR-C4 / to be sliced | The gate checkpoint prints the route of its cause | RR39-RR44 | `internal/gate` | no |
| RR-C5 / to be sliced | The commitment verb prints routes of the right authority | RR45-RR48 | `internal/commitment/commitcmd`, `cmd/bench` | no |
| RR-C6 / to be sliced | No write-verb source composes a route outside the registry | RR49-RR50 | `internal/conformance` | no |

RR-C1 comes first, because every later chunk raises its faces through the shared constructor.
RR-C2 changes the authorization sentence that RR-C3 also prints.
RR-C4 waits for the FT392 landing.
RR-C6 comes last, because it turns red on any route that an earlier chunk has not moved.

### Completion plan

The version 1 plan records future implementation evidence.
It claims no current implementation pass, red, or probe result.
This plan is provisional, and the slice replaces it.

```bench-completion-plan
{
  "version": 1,
  "chunks": [
    {
      "id": "RR-C1",
      "tickets": [
        "00-draft-closure.md"
      ],
      "verification": [
        {
          "id": "registry",
          "command": "bench test --package ./internal/refusalroute"
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
  ]
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
| RR09 | 9 | Every agent route command step gets `Blocked == false` from `benchguard.Classify` | planned TestAgentRoutesPassTheWiredGuards in internal/conformance/refusal_route_guard_test.go | A step that chains a second Bench call passes unseen |
| RR10 | 9 | Every agent route command step gets the empty string from `benchguard.PoolReference` with a pool-path sample | planned TestAgentRoutesPassTheWiredGuards in internal/conformance/refusal_route_guard_test.go | A `cd` or `git -C` into the pool passes unseen |
| RR11 | 10 | The guard check refuses an agent step that runs `bench worktree exec <target> -- bench <verb>` | planned TestAgentRoutesPassTheWiredGuards in internal/conformance/refusal_route_guard_test.go | An exec route survives today and breaks when FT341 refuses it |
| RR12 | 11 | The guard check reds an injected agent face whose route step is `git merge <commit>` | planned TestAgentRouteGuardCheckBites in internal/conformance/refusal_route_guard_test.go | A check that never classifies passes every route |
| RR13 | 12 | A route step that runs a tree-scoped Bench verb at a worktree renders `--in <label>` | planned TestRouteRendering in internal/refusalroute/route_test.go | A renderer that keeps the exec form fails RR11 at the first real face |
| RR14 | 13, 15 | Each land face in the registry has exactly one producing fixture, and each fixture produces a registered land face | `internal/worktree/identity_component_test.go` (`TestLandingRefusalRegistryHasAProducingFixture`) | A face moved without its fixture reaches an operator unproven |
| RR15 | 14 | Each land fixture follows the printed route, reruns the landing, and the face's sentence no longer prints | planned TestLandingFacesFollowTheirRoutes in internal/worktree/refusal_route_test.go | A route that names the wrong repair leaves the face in place |
| RR16 | 16 | The landing composition-conflict refusal prints a `next=` value that starts with `reviewer: ` | planned TestConflictRepairIsAReviewerRoute in internal/worktree/refusal_route_test.go | An agent route that prints `git merge` is the collision 8 denial |
| RR17 | 16 | The landing composition-conflict-pending refusal prints a `next=` value that starts with `reviewer: ` | planned TestConflictRepairIsAReviewerRoute in internal/worktree/refusal_route_test.go | The pending-merge arm prints `git merge --continue`, which the guard denies |
| RR18 | 17 | The landing destination-not-clean refusal prints a `next=` value that starts with `reviewer: ` | planned TestLandingFacesFollowTheirRoutes in internal/worktree/refusal_route_test.go | An agent route lets the agent discard the reviewer's primary-checkout work |
| RR19 | 18 | The landing source-not-clean refusal prints a route that contains `bench commit --in ` and the source label | planned TestLandingFacesFollowTheirRoutes in internal/worktree/refusal_route_test.go | A route that names no commit command leaves the agent at a raw commit, which the guard denies |
| RR20 | 19 | Each landing preflight route ends with the caller's own re-run | `internal/worktree/land_surface_test.go` (`TestLandCommandReportsEveryRefusalInOnePreflight`) | A face that drops the re-run leaves a second lookup |
| RR21 | 20 | A fold of `main` into a target whose committed tip fails its lane prints a `refused{` record that contains `next=` | planned TestRedSourceFoldNamesAnExit in internal/worktree/merge_route_test.go | This is the collision 5a repro: the `lane fail` refusal prints no `next=` today |
| RR22 | 21 | The target-red fold route contains `bench commit --in ` and the target label | planned TestRedSourceFoldNamesAnExit in internal/worktree/merge_route_test.go | A route that reruns the fold alone loops on the same red |
| RR23 | 21 | The target-red fold route ends with `bench worktree merge --from ` and the incoming spelling and the target | planned TestRedSourceFoldNamesAnExit in internal/worktree/merge_route_test.go | A route that stops at the commit leaves the fold undone |
| RR24 | 21 | After the fixture removes the red file and runs the printed commit step, the printed fold step exits 0 | planned TestRedSourceFoldNamesAnExit in internal/worktree/merge_route_test.go | A route whose steps do not clear the red keeps the deadlock |
| RR25 | 22 | A fold into a dirty target prints a `refused{` record whose `next=` contains `bench commit --in ` | planned TestRedSourceFoldNamesAnExit in internal/worktree/merge_route_test.go | This is the second refusal of the collision 5a repro, with no `next=` today |
| RR26 | 23 | A `lane fail` fold whose incoming commit adds the red to a lane-green target prints a `next=` value that starts with `reviewer: ` | planned TestMergeFacesFollowTheirRoutes in internal/worktree/merge_route_test.go | A merge that never grades the target alone gives every lane red the agent route, which loops |
| RR55 | 23 | A `lane fail` fold whose target tip alone fails the lane prints a `next=` value that does not start with `reviewer: ` | planned TestRedSourceFoldNamesAnExit in internal/worktree/merge_route_test.go | A merge that sends every lane red to the reviewer hands back a repair that is the agent's |
| RR56 | 23 | A `candidate` gate-kind fold red prints a `next=` value that starts with `reviewer: ` | planned TestMergeFacesFollowTheirRoutes in internal/worktree/merge_route_test.go | A face map that ignores the gate kind's attribution gives a fold red the agent route |
| RR27 | 24 | A merge composition conflict prints a `next=` value that starts with `reviewer: ` | planned TestConflictRepairIsAReviewerRoute in internal/worktree/refusal_route_test.go | The merge keeps printing the guard-denied `git merge` as an agent step |
| RR28 | 25 | A fold with a dirty sibling prints a route that contains `bench commit --in ` and no `bench worktree exec` | planned TestMergeFacesFollowTheirRoutes in internal/worktree/merge_route_test.go | The current route prints the exec form that FT341 refuses |
| RR29 | 13, 27 | Each merge face and each reset face has exactly one producing fixture that follows its route out of the face | planned TestMergeFacesFollowTheirRoutes in internal/worktree/merge_route_test.go | A merge or reset face with no fixture reaches an operator unproven |
| RR30 | 26 | The `inherited` authorization refusal sentence equals `prospective authorization refused: inherited (the gate ran red on the composed tree and no green baseline attributes the red to this diff)` | `internal/landing/landing_reviewed_test.go` (`TestRefusalMessageNamesTheOperatorActionAndTheOpenReason`) | An inline action prints a second route beside the face's route |
| RR31 | 26 | The merge retries an empty-reason infrastructure refusal exactly once | `internal/worktree/merge_test.go` (`TestMergeRetriesOnlyAVerifiedEmptyReasonInfrastructureRefusal`) | A retry that matches the old sentence never fires after the sentence changes |
| RR32 | 28 | A commit exit 3 prints `next=` with the value `bench worktree reset --to <published-commit> <checkout>` for its published commit and its root | planned TestPublishedUnreconciledRouteIsTheResetPlan in internal/commit/refusal_route_test.go | This is the collision 8 restore repro: the current route is `git restore`, which the guard denies |
| RR33 | 28 | A commit exit 3 `next=` value does not contain `git restore` | planned TestPublishedUnreconciledRouteIsTheResetPlan in internal/commit/refusal_route_test.go | A route that keeps the restore beside the reset still prints a denied step |
| RR34 | 28 | A commit exit 3 on a root path that is not line-safe prints the placeholder `<checkout>` | planned TestPublishedUnreconciledRouteIsTheResetPlan in internal/commit/refusal_route_test.go | A raw control byte reaches the line-structured record |
| RR35 | 29 | After a commit exit 3 in an assignment worktree, the printed reset plan and its `--apply` leave `git status --porcelain` empty at the published commit | planned TestCommitExitThreeRouteReconcilesTheCheckout in internal/worktree/commit_route_test.go | A reset route that does not reconcile leaves the exit 3 state |
| RR36 | 30 | Each commit face has exactly one producing fixture that follows its route out of the face | planned TestCommitFacesFollowTheirRoutes in internal/commit/refusal_route_test.go | A commit face with no fixture reaches an operator unproven |
| RR37 | 30 | A red commit refusal prints a `next=` line on stderr whose route ends with the caller's commit command | planned TestCommitFacesFollowTheirRoutes in internal/commit/refusal_route_test.go | The current `inherited` refusal prints only `run bench gate --fresh`, which re-grades the same red |
| RR38 | 30 | A commit in the primary checkout prints a `next=` route that contains `bench worktree create --request` | planned TestCommitFacesFollowTheirRoutes in internal/commit/refusal_route_test.go | A refusal that keeps its route only inside the sentence bypasses the registry |
| RR39 | 31 | A checkpoint refusal for a missing completion record prints no `help[1]{cmd,why}` row | `internal/gate/review_checkpoint_test.go` (`TestReviewCheckpointRefusalRoute`) | This is the FT330 defect: the fixed write-access row prints for every cause |
| RR40 | 32 | A checkpoint refusal for a missing completion record prints `next=bench preflight review example` | `internal/gate/review_checkpoint_test.go` (`TestReviewCheckpointRefusalRoute`) | The move to the registry drops the current evidence route |
| RR41 | 33 | A complete checkpoint on a dirty checkout prints a `next=` route that contains `bench commit -m <msg> -- <path>...` | planned TestCheckpointFacesFollowTheirRoutes in internal/gate/refusal_route_test.go | The FT392 refusal names no typed route |
| RR42 | 34 | A subject-capture fault prints `next=` with a route that contains `bench doctor` and `--fresh` | `internal/gate/run_outcomes_test.go` (`TestGateRunRetainsSubjectConstructionCause`) | The fallback keeps the write-access text for a fault that is not a write-access fault |
| RR43 | 35 | A subject-capture fault prints no `help[1]{cmd,why}` row | `internal/gate/run_outcomes_test.go` (`TestGateRunRetainsSubjectConstructionCause`) | A route printed beside the old row gives two answers |
| RR44 | 13, 14 | Each gate face has exactly one producing fixture that follows its route out of the face | planned TestCheckpointFacesFollowTheirRoutes in internal/gate/refusal_route_test.go | A gate face with no fixture reaches an operator unproven |
| RR45 | 36 | A `bench commitment start` refusal for an outcome outside the active milestone prints a `next` cell that starts with `reviewer: ` | planned TestCommitmentFacesFollowTheirRoutes in internal/commitment/commitcmd/refusal_route_test.go | The current cell tells the agent to re-plan a commitment that only the reviewer changes |
| RR46 | 37 | A `bench commitment start` refusal without an owned assignment prints a `next` cell that contains `bench worktree create --request` | planned TestCommitmentFacesFollowTheirRoutes in internal/commitment/commitcmd/refusal_route_test.go | The current cell prints the plan command for a cause that the agent clears |
| RR47 | 37 | A `bench commitment verify` refusal keeps its `bench commitment verify --milestone` cell | `internal/commitment/verification_test.go` (`TestCommitmentUnmetCriterion`) | The move to the registry drops the current agent route |
| RR48 | 13, 14 | Each commitment face has exactly one producing fixture that follows its route out of the face | planned TestCommitmentFacesFollowTheirRoutes in internal/commitment/commitcmd/refusal_route_test.go | A commitment face with no fixture reaches an operator unproven |
| RR49 | 38 | The bypass check passes on the tree after every chunk lands | planned TestNoWriteVerbComposesARouteOutsideTheRegistry in internal/conformance/refusal_route_bypass_test.go | A hand-composed route left in a verb shows as a red |
| RR50 | 39 | The bypass check reds a planted `next=` literal in a write-verb production file | planned TestRouteBypassCheckBites in internal/conformance/refusal_route_bypass_test.go | A check that scans no file passes every tree |
| RR51 | 40 | `bench recovery` prints `recovery[N]{verb,face,authority,route}` and exits 0 | planned TestRecoveryListsEveryFace in internal/refusalroute/command_test.go | A missing header breaks the TOON contract |
| RR52 | 41 | `bench recovery` prints exactly one row for each registered face, in registry order | planned TestRecoveryListsEveryFace in internal/refusalroute/command_test.go | A hand-kept matrix drifts from the registry |
| RR53 | 40 | `bench help` lists `bench recovery` | `cmd/bench/help_inventory_test.go` (`TestHelpInventoryIsComplete`) | An unlisted verb is a dead key |
| RR54 | 42 | The reference guide names `bench recovery` as the recovery matrix | review-owned | No gate check reads this prose |

Not covered: story 43 — FT342 owns the moved-`main` decision, and RR26 routes to it.
Not covered: story 44 — the exit-2 grammar refusals keep their behavior, and the usage package owns their usage line.

### Edge inventory

The hostile-input checklist applies to the facts that a route renders.
A path fact can hold a space, a glob byte, or a control byte.
The renderer shell-quotes a line-safe value and prints the slot placeholder for a value that is not line-safe; RR34 grades the placeholder.
A label fact comes from the assignment ledger, and the `--in` resolver refuses an ambiguous label under FT341.

The guard check runs over sample facts, and the producing fixtures run over real values.
A route can pass the check with samples and fail it with a real path.
So RR15 and the other follow rows run the real printed steps.

The fail-closed refusals are the unregistered face (RR07) and each verb's `<verb>-handback` face.

- **Won't handle**: the moved-`main` decision — FT342 owns it, and `merge-candidate-red` routes to it.
- **Won't handle**: a fold that admits an inherited red — the decision source closes the merge policy, and `merge-inherited-red` routes the repair before the fold.
- **Won't handle**: the exit-2 grammar refusals — the usage package owns their usage line, and `usage.Parse` keeps it.
- **Won't handle**: the gate run-state refusals outside the checkpoint funnel — they serve every gate run, and `operational` keeps their reason line.
- **Won't handle**: the inline route tails in commitment repository errors — `bench status` reads them, and the face prints the typed route beside them.
- **Won't handle**: the non-write worktree verbs (`list`, `path`, `clean`, `release`, `build`) — the decision source names six write verbs, and `recoveryRoute` keeps their routes.
- **Won't handle**: a change to either wired guard — the closed decision forbids it, and the guard check reads the guards unchanged.

## Ownership fences

The prospective build owns these exact paths:

- `internal/refusalroute/registry.go`
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
- `internal/worktree/build.go`
- `internal/worktree/reset.go`
- `internal/worktree/reset_apply.go`
- `internal/worktree/reset_restore.go`
- `internal/worktree/path.go`
- `internal/worktree/classifier.go`
- `internal/worktree/identity_component.go`
- `internal/worktree/refusal_route_test.go`
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

This fence is provisional. The slice closes it with the closure that build preflight proposes.

## Out of scope

- The inline route tails in commitment repository error sentences: 9 edits, 2 gate runs.
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
- The target-alone lane run on a `lane fail` fold (RR26, RR55). The `lane fail` kind carries no attribution, so without this run a fold red cannot pick its authority.

### Posture change

- `TestPublicationRemainderNextRestoresEveryNamedPathShellQuoted` and `TestPublicationRemainderSanitizesTheFailedPathAndPointsAtNamedPaths` pin the `git restore` route, and RR32-RR34 replace both pins.
- `TestGateRunRetainsSubjectConstructionCause` pins the fixed write-access help row, and RR42-RR43 replace that pin.
- The landing tests that compare a route with `landingRefusalFaceByName(...).route(...)` read the shared registry instead.
- The tests in `internal/landing`, `internal/commit`, and `internal/worktree` that read `run bench gate --fresh`, `fix the failures above`, or `run bench doctor` in a refusal sentence read the face route instead.

### Source-sentence-to-row table

| FT393 sentence | rows |
|---|---|
| A red source could not fold `main`, because the merge refused the inherited red and named no next action. | RR21-RR24, RR26, RR55, RR56 |
| A dirty fold also refused with no next action. | RR25 |
| The destructive-git guard denies a raw `git merge` after a conflict. | RR16, RR17, RR27 |
| The destructive-git guard denies `git restore --worktree` after a `bench commit` exit 3. | RR32-RR35 |
| One refusal-route registry serves every write verb. | RR01, RR14, RR29, RR36, RR44, RR48, RR49 |
| Each face declares a typed route and its authority, agent or reviewer. | RR02, RR03, RR05, RR06 |
| Every refusal goes through the registry. | RR49, RR50 |
| Every agent route passes the wired guards. | RR08-RR12 |
| Each producing fixture follows its route out of the refused face. | RR15, RR24, RR29, RR35, RR36, RR44, RR48 |
| The recovery matrix renders from the registry. | RR51-RR54 |
| A checkpoint refusal names the recovery route of its cause (FT330). | RR39-RR43 |

### Pre-review proof checklist

- `Cited symbols`: `landingRefusalFaces`, `landingFaceRefusalOf`, `conflictRepairPrefix`, `conflictContinuePrefix`, `restoreNext`, `publicationRemainder`, `refusalMessage`, `refusalGuidance`, `routedRefusal`, `subjectUnavailableHelp`, `refusalNext`, `checkoutClean`, `mergeReconcileNext`, `retryEmptyReasonInfrastructureFold`, `gitguard.Classify`, `benchguard.Classify`, and `benchguard.PoolReference` resolve in the tree at a9506ce5.
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
- Tests that read the moved literals: `internal/commit/landing_test.go`, `internal/commit/dry_run_test.go`, `internal/landing/landing_reviewed_test.go`, `internal/worktree/merge_test.go`, `internal/worktree/land_surface_test.go`, `internal/worktree/land_journey_test.go`, `internal/gate/run_outcomes_test.go`, and `internal/gate/review_checkpoint_test.go`.

### Primary sources

- `internal/worktree/land_refusal.go:187-340`: the landing registry and its one constructor.
- `internal/worktree/identity_component_test.go:200-360`: the producing-fixture walk.
- `internal/worktree/classifier.go:23-65`: the `refusal` record and its fields.
- `internal/worktree/merge.go:66-161`: the merge refusals and the conflict route.
- `internal/commit/commit.go:196-238`: the exit 3 record and `restoreNext`.
- `internal/landing/landing.go:283-315`: the authorization sentence and its inline action.
- `internal/gate/run_transaction.go:43-58` and `internal/gate/gate.go:331-336`: the funnel and the fixed help row.
- `internal/gate/checkpoint.go:87-109`: the completion-evidence route.
- `internal/commitment/commitcmd/command.go:336-344`: the commitment refusal table.
- `internal/gitguard/gitguard.go:126-150` and `internal/benchguard/benchguard.go:101-120`, `internal/benchguard/pool.go:10-40`: the guard decisions.
- `internal/worktree/reset.go:28-240`: the reset plan keeps a dirty checkout in a recoverable envelope.

Cheap-tier research delegates read these sources and reported citations.
The author re-read each definition above before the rows locked.
