# Recovery-aware ref inventory

Status: staged

Roadmap: FT199

Decision source: the compiled map `specs/ref-inventory/decisions/ref-inventory.md`, ready on 2026-09-27 with eleven resolved tickets.

Verification log: 2 iteration(s) to accept — the opus/high round returned 12 blocking findings, then 6, and both folds are in. Iteration 1 found the import cycle, the unreachable path operand, the wide holder set, and the unbound class fingerprint. It also found the unfenced test readers and the absent apply action. Iteration 2 found the content-landed holder and the status precedence, which the reviewer decided on 2026-09-27. It also found the checked-out candidate refusals, the slug export, and the shift-path route. The cap of 2 is exhausted, so reviewer sign-off accepts the second fold.

## Problem

`bench worktree clean --discard-branch --unclaimed` discards every unrecorded branch in the two Bench namespaces.
It never calls the landed classifier that the landing prune and `clean --landed` share.
On 2026-09-27 the repository held 43 refs in the assignment namespace, 42 of them unrecorded.
Every one of them carried content that `main` lacks under all four landed proofs.
Seven of the ten root tips are evidence that roadmap rows and research documents cite.
`bench status` named the destructive sweep as its top action on the count alone, with no recovery ref behind it.

Only an assignment record reaches the explicit discard forms, so an unrecorded ref has no per-ref route.
The assignment-scoped recovery namespace empties at each session start.
So nothing can hold a discarded ref's work across sessions.
No record carries a spec field, so no route lists the assignments that a spec retirement superseded.

## Solution

The unclaimed plan classifies each unrecorded Bench-namespace branch as landed, subsumed, or unique.
A landed ref passes one of the four landed proofs.
A subsumed ref has a tip that equals, or is a strict ancestor of, a holder's tip, and the row names that holder.
Every other ref is unique.
The bulk sweep, with a fingerprint or with `--apply-current`, deletes landed and subsumed rows only.
A unique row is retained, and its detail names the explicit discard command.

`bench status` counts the three classes and routes to the plan-only unclaimed command.
The plan output alone names the apply command, and only when a row removes.

`--target` resolves an unrecorded Bench-namespace branch by its branch path or by its assignment id segment.
The apply of a unique target writes `refs/bench/discarded/<yyyymmdd>/<branch-path>` at the exact tip.
Then it deletes the branch at that tip.
The session-start sweep removes a discarded ref 30 or more days after its date and reports the count.
`bench spec retire` lists the active or cleanup-pending assignments whose label or request token contains the slug.
Each listed line carries the explicit discard command, and one more line carries the unique ref count.

## User stories

Line: opus / high.
Implementation-line reason: RI-C2b is the hardest chunk. Its transaction writes a ref before an exact delete and must survive a fault between them. The spec is exact, and the seams are the existing planner, fingerprint, and sweep. The transaction rows have no gate coverage until their fault fixtures exist, so the weak-gate row bumps mid plus medium to mid plus high.
Harder chunks: RI-C2b.

Classification:

1. As a coordinator, I want the inventory to cover every branch under both Bench namespaces, so that no Bench-created ref escapes a class.
2. As a coordinator, I want a recorded branch to stay out of every plan row, so that live work is never offered.
3. As a coordinator, I want a checked-out branch and the default branch to stay protected, so that the sweep cannot delete a live checkout.
4. As a coordinator, I want a branch whose tip is an ancestor of `main` to classify as landed, so that a merged ref can retire.
5. As a coordinator, I want a squash-folded branch to classify as landed through the reverse-apply proof, so that ancestry alone does not decide.
6. As a coordinator, I want an equal-tip branch to classify as subsumed by the lexically first ref, so that one member survives.
7. As a coordinator, I want a strict ancestor of a holder's tip subsumed by that holder, so that the holder is fixed.
8. As a coordinator, I want an active recorded branch to hold the refs at or beneath its tip, so that review-axis branches retire.
9. As a coordinator, I want landed refs excluded from holders, so that a content-landed ref cannot authorize another ref's deletion. A ref beneath an ancestry-landed ref classifies as landed.
10. As a coordinator, I want every other unrecorded branch to classify as unique, so that content `main` lacks is named as such.
11. As a coordinator, I want a ref with no commit tip to print an error row, so that a damaged ref never rounds up.
12. As a coordinator, I want the plan to read refs only and write nothing, so that a plan is safe to repeat.

Plan rows and the bulk sweep:

13. As a reviewer, I want each classified row's detail to start with `class=<class>`, so that I read the class without a second command.
14. As a reviewer, I want a subsumed row to continue with `holder=<ref>`, so that I can verify reachability myself.
15. As a reviewer, I want a landed or subsumed row to carry the action `discard-remove`, so that the sweep's reach is visible.
16. As a reviewer, I want unique rows to carry `retain` and a branch-path discard command, so that each route selects its ref.
17. As a reviewer, I want the plan fingerprint to bind each row's ref, tip, class, and holder, so that a class change refuses at apply.
18. As a reviewer, I want `--apply <fingerprint>` to delete landed and subsumed rows only, at their exact tips, so that a unique ref survives.
19. As a reviewer, I want `--apply-current` to plan and apply in one call with the same reach, so that the shortcut stays available.
20. As a reviewer, I want a stale apply to refuse and print the plan-only re-plan command, so that I re-read before I act.
21. As a reviewer, I want every deleted tip to stay reachable from `main` or from a surviving ref, so that nothing is lost.
22. As a reviewer, I want the plan to name the apply command only when a row removes, so that a unique-only plan stays safe.
23. As a reviewer, I want a plan over an empty namespace to exit 0 with an empty table, so that a clean repository stays quiet.

The status signal:

24. As an operator, I want the status git row to count landed, subsumed, and unique refs, so that the counts match the plan.
25. As an operator, I want the status git row to route to the plan-only unclaimed command, so that status never names a destructive command.
26. As an operator, I want any Bench-namespace ref to route the git row to the plan command with every detail, so that nothing is lost.
27. As an operator, I want a repository with no Bench-namespace refs to keep today's git row, so that the change is invisible there.
28. As an operator, I want the routed status command to remove nothing when I run it, so that a copied action is safe.

Explicit discard of a unique ref:

29. As a reviewer, I want `--target <branch path>` to resolve an unrecorded Bench-namespace branch, so that I can name one ref.
30. As a reviewer, I want `--target <assignment id>` to resolve the one unrecorded branch whose id segment matches, so that the id from the row suffices.
31. As a reviewer, I want an id segment that matches two unrecorded branches to refuse and name both refs, so that ambiguity never picks one.
32. As a reviewer, I want `--target archive/x` to keep today's refusal, so that a foreign branch is never selected.
33. As a reviewer, I want a unique target's plan row to show `class=unique` and the planned discarded ref, so that I see the preservation.
34. As a reviewer, I want the explicit apply to need the set fingerprint, so that no one-call discard exists for a unique ref.
35. As a reviewer, I want the transaction to write the discarded ref before the delete, so that a fault leaves a handle.
36. As a reviewer, I want a fault before the write to leave the branch and write no ref, so that the failure is clean.
37. As a reviewer, I want a fault after the write to leave both refs, so that a second apply can finish.
38. As a reviewer, I want an existing discarded ref at another tip to refuse the apply, so that the write never overwrites a handle.
39. As a reviewer, I want the outcome row to name the discarded ref in its recovery cell, so that the handle is in the output.
40. As a reviewer, I want a landed or subsumed explicit target to discard without a discarded ref, so that both routes agree on reach.
41. As a reviewer, I want an explicit target that names a recorded assignment to keep today's route, so that the record path is unchanged.

Discarded ref lifetime:

42. As an operator, I want a discarded ref to live under `refs/bench/discarded/<yyyymmdd>/<branch-path>`, so that its date is in its name.
43. As an operator, I want the session-start sweep to remove a discarded ref 30 or more days old, so that the namespace stays bounded.
44. As an operator, I want a discarded ref 29 days old to survive the sweep, so that the window holds.
45. As an operator, I want the lifecycle-namespace emptying rule to leave the discarded namespace alone, so that a session start keeps the handle.
46. As an operator, I want a discarded ref with an unparseable date segment to survive the sweep, so that a malformed name fails closed.
47. As an operator, I want the sweep to delete each discarded ref at its listed object, so that a concurrent move refuses.

Retire listing:

48. As a reviewer, I want retire to list each live assignment whose label or request names the slug, so that superseded work has a route.
49. As a reviewer, I want each listed assignment line to carry the exact `--target` discard command, so that I copy it.
50. As a reviewer, I want retire to add one count line with the plan command, so that I see the unique inventory.
51. As a reviewer, I want the retire listing to discard nothing and change no exit code, so that retire stays a spec operation.
52. As a reviewer, I want a retire with no matching assignment to print the count line alone, so that the shape is stable.

Reviewed exclusions:

53. As a reviewer, I want no hold marker, keep list, or hold verb, so that the unique class alone is the keep.
54. As a reviewer, I want no automatic supersession proof from a spec retirement, so that a human claim never becomes a discard proof.
55. As a reviewer, I want `clean --landed` and the landing prune unchanged, so that recorded assignments keep their route.
56. As a reviewer, I want the live refs of this repository untouched by the build, so that their disposition stays with FT346, FT347, and FT348.
57. As a maintainer, I want the glossary to name the shift namespace as `refs/heads/bench/shift-`, so that the term matches the code.

Review-round additions:

58. As a reviewer, I want `--target <unrecorded>` without `--discard-branch` to plan a `retain` row that names the flag, so that no discard runs unasked.
59. As a reviewer, I want the explicit fingerprint to bind the class, so that a stale explicit apply writes no discarded ref.
60. As a reviewer, I want the explicit delete to use the exact tip, so that a branch moved after the write survives.
61. As a reviewer, I want retire to print `unique refs: unavailable` with the error when the count fails, so that the exit code holds.
62. As a reviewer, I want the fallback to select from the unclaimed set only, so that a checkout or a recorded branch never enters it.
63. As a maintainer, I want the glossary `holder` term to name only a recorded branch or a unique root, so that it matches the rule.

## Implementation decisions

The class computation is one function beside the unclaimed planner, and it reads refs only.
It takes the recorded assignments, the protected set, the default branch, and the sorted unrecorded refs.
It returns one class and one holder per ref.
It calls the shared landed proof for each unrecorded ref, so the sweep and the landing prune cannot disagree.
A ref whose tip does not resolve to a commit produces a row with action `error`, and the set fingerprint stays empty.
The error row's detail names the ref and the object type.

The landed class takes precedence over the subsumed class.
A holder is an active or cleanup-pending recorded assignment branch, or a unique root.
No landed ref is a holder.
A recorded assignment branch that `LandedInDefault` proves landed by content only is not a holder, including active and cleanup-pending assignments.
A ref beneath an ancestry-landed ref is itself landed by ancestry, because ancestry is transitive.
A resolving unrecorded symbolic ref in either Bench namespace produces an error row naming the ref and `symref`, never a silent exclusion.

A ref that is landed by content only, the squash fold through the reverse-apply proof, is deleted by the bulk sweep but holds nothing.
So a ref under a content-landed ref that is not itself landed classifies as unique.
A checked-out foreign branch, a complete record's branch, and a subsumed ref are never holders.
A ref is subsumed when its tip equals, or is a strict ancestor of, a holder's tip.

Among unrecorded refs with an equal tip and no other holder, the lexically first full ref name is the unique root.
The holder a row names is the lexically first holder whose tip reaches the ref, with recorded branches first, then unique roots.

The unclaimed plan row keeps the seven cleanup columns.
The `tracked` cell stays `unclaimed`.
The `detail` cell starts with `class=<class>`, continues with ` holder=<ref>` for a subsumed row, and ends with the existing removal text for a removing row.
A unique row carries the action `retain`.
Under ticket 1 its detail ends with `retained: content main lacks`, and ticket 4 replaces that suffix with `bench worktree clean --discard-branch --target <branch path>`.

The branch path is the row's ref without `refs/heads/`, for an assignment branch and for a shift branch alike.
The printed route and the explicit apply action derive their selector from one function.
An id operand stays an input form under stories 30, 31, and 41, and it is never a printed route.

The unclaimed fingerprint version moves to `bench-unclaimed-assignment-branches/v2` and binds each row's ref, tip, class, and holder.
The apply loop, for a fingerprint and for `--apply-current`, skips a row whose action does not remove.
A plan that holds an error row exits 1 and prints no apply help action, and `--apply-current` beside an error row refuses before any delete.
The apply help action is the exact token `bench worktree clean --discard-branch --unclaimed --apply <fingerprint>`, and today's plan prints no apply action at all.

Ticket 1 changes the status action for the git row to `bench worktree clean --discard-branch --unclaimed` and removes the `--apply-current` entry from the action table.
That change lands with the classifier.
After ticket 1 the routed `--apply-current` call keeps a unique branch, so the system route test would red otherwise.

Ticket 2 then makes `bench status` read the class counts through one exported function of the worktree package that wraps the same planner.
The landing prune test also reads that function, so no exported function serves a test alone.
The git row details read `<n> landed ref`, `<n> subsumed ref`, and `<n> unique ref`, plural as the existing helper renders, and the row omits a zero class.
Any Bench-namespace ref count above zero routes the git row to the plan command, and the dirty-path and unpushed-commit details stay in the row text.

Status shows a nonzero `<n> faulted ref` count through the plural helper, or `unclaimed refs unavailable` on a planner failure, and routes both to the plan command.
When Git state fails in a repository, unclaimed rows or planner errors still route status to the plan with their details beside `git state unavailable`.
A non-repository reports only `git state unavailable` with the `git status` action.
The worktree package exports the plan command spelling, and status and the retire listing read it.

The plan-only route costs one landed proof per unrecorded ref plus one ancestry check per root per ref.
It also costs one landed proof per active or cleanup-pending record whose branch resolves to a commit.
RI-C1b records the plan time over a 43-ref fixture as a number, not a bound.

The discarded namespace is one ledger constant, `refs/bench/discarded/`, beside the recovery and reset namespaces.
It joins neither the lifecycle emptying list nor the reset rule.
The sweep parses the first path segment under the namespace as a UTC date `yyyymmdd`.
It keeps a ref whose segment does not parse.
It deletes a ref at its listed object when the resume instant is at or past the date plus 30 days at 00:00:00Z.

The delete never follows a symbolic ref, so a discarded symref leaves its target in place.
The swept count joins the existing swept-refs total.

The explicit set planner routes an operand to the unrecorded-branch fallback before the relative-path check.
The fallback takes an operand that starts with `refs/heads/bench/assign/`, `bench/assign/`, `refs/heads/bench/shift-`, or `bench/shift-`.
It also takes a 32-character hexadecimal operand that resolves no record, and matches it to the last segment of exactly one unclaimed assignment branch.
Two matches refuse with both refs, and no match keeps the unassigned refusal.
Every other operand keeps today's path check, so `archive/x` still refuses with `relative path targets are unsupported`.

The fallback's candidate set is the unclaimed set that the bulk planner selects, not every unrecorded branch.
A prefix operand that names a checked-out branch refuses with an error row that names the checkout.
A prefix operand that names the default branch refuses with an error row.
A prefix operand that names a recorded branch plans through the record and prints no `class=` prefix.
The set collapses `--target <id>` and `--target bench/assign/<owner>/<id>` into one row by branch ref.
Recorded rows sort by assignment id first, then unrecorded rows sort by branch ref, and the explicit set fingerprint version moves to `bench-explicit-set/v2`.

An unrecorded target row classifies through the same class function.
Without `--discard-branch` the row carries `retain` and its detail names `--discard-branch`, and the apply removes nothing.
A unique row plans the discarded ref `refs/bench/discarded/<yyyymmdd>/<branch-path>` in its recovery cell.
The date comes from a clock seam on the package's join set, so a test fixes the day.
The explicit set fingerprint binds each unrecorded row's ref, tip, class, holder, and planned discarded ref.

The apply of an unrecorded row runs inside the set's own transaction order.
When the planned ref exists at the row's tip as a direct ref, the write is skipped.
When it exists at another tip, or when the planned path is a symbolic ref, the apply refuses and writes nothing.
Otherwise the write uses the zero old value.
The delete uses the exact tip, so a branch moved after the write survives with an error row.

A fault boundary step precedes the read of the planned path, one sits between that read and the write, and one follows the write.
A target that the class function faults prints an error row with no fingerprint, and both apply forms refuse.

The retire listing lives at the `cmd/bench` dispatch of `bench spec retire`, because `internal/worktree` already imports `internal/spec`.
The retire command returns as today.
The dispatcher appends the candidate lines and the count line only after a code-0 retire that printed a `next:` line.
It inserts them before that line, and it appends nothing after `--help` or after a refusal.
`bench spec retire` takes `<spec.md | slug>`, so `internal/spec` exports one slug derivation that the retire command and the dispatcher both call.

A candidate line reads `superseded candidate: <assignment id> <label> — bench worktree clean --discard-branch --target <assignment id>`.
The count line reads `unique refs: <n> — bench worktree clean --discard-branch --unclaimed`.

When the ledger read or the planner fails, the count line reads `unique refs: unavailable — <error>`, and the exit code holds.
When the faulted count is nonzero, the count line reads `unique refs: <n>, <faulted> faulted — bench worktree clean --discard-branch --unclaimed`, and a zero faulted count omits the suffix.
The wrapper takes the slug from the operand the retire parsed, not from the last raw argument.
The listing includes the calling worktree's own assignment, and the explicit discard retains that assignment on its live lease as today.

The glossary term for the unclaimed ref names the shift namespace as `refs/heads/bench/shift-`, which is what the ledger declares.
The glossary term for the holder names an active or cleanup-pending recorded assignment branch, or a unique root, and no landed ref.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| RI-C1a / 1-classify-unclaimed-refs.md, 6-repair-glossary-shift-namespace.md | The unclaimed plan classes each ref, names holders, retains unique rows, applies landed and subsumed rows only, and status routes to the plan | RI1 to RI15, RI17 to RI23, RI25, RI28, RI55, RI57, RI59 to RI63, RI66, RI72, RI73, RI81 to RI87 | `internal/worktree` clean classes and unclaimed tests, the system route test | no |
| RI-C1b / 2-route-status-to-the-plan.md | Status counts the three classes | RI24, RI26, RI27, RI88 to RI92 | `internal/status` producible signals and the landed second case of the system route test | no |
| RI-C2a / 3-sweep-discarded-refs.md | The discarded namespace exists, survives the lifecycle emptying, and expires at 30 days | RI42 to RI47, RI93 to RI95 | `internal/worktree` reconcile tests | no |
| RI-C2b / 4-discard-a-unique-ref-by-target.md, 5-list-retire-candidates.md | An unrecorded unique ref discards by target with a discarded ref first, and retire lists candidates | RI16, RI29 to RI41, RI48 to RI52, RI58, RI64, RI65, RI67 to RI71, RI74 to RI80, RI96 to RI101 | `internal/worktree` discard tests and `cmd/bench` retire dispatch tests | yes |

## Testing decisions

- A good test drives `cleanCommandWith` over a real Git repository with real branches, reads the rendered rows, and then reads the refs back with `show-ref`.
- The class function receives its own table test over the ref shapes of the edge inventory.
- The transaction rows use the `cleanupBoundary` step seam that the existing apply tests use, with two new step tokens around the discarded ref write.
- The sweep rows drive `reconcileLifecycleDebris` with a fixed instant and refs planted under the discarded namespace.
- The status rows extend the producible-signal table, and the system test drives the routed command through the sealed binary.
- The retire rows drive the `cmd/bench` dispatcher over a fixture repository with a merged-implemented spec and planted records.
- The gate observes the feature through `bench test` over `internal/worktree`, `internal/status`, `cmd/bench`, and the system suite.

### Seam diagram

    trigger: bench worktree clean --discard-branch --unclaimed [--apply <fp> | --apply-current]
        │
        ▼
    refs, ledger, default branch  ──▶  [ planUnclaimedAssignmentSet + classifyUnclaimedRefs ]  ──▶  rows with class, holder, action
                                          ◀ tests attach here: cleanCommandWith over a fixture repo; show-ref after apply

    trigger: bench worktree clean --discard-branch --target <id | branch path> [--apply <fp>]
        │
        ▼
    operand  ──▶  [ planExplicitSet: namespace-prefix fallback over the unclaimed set, else the record path ]  ──▶  explicit row, planned discarded ref
                        │
                        ▼
                  [ discard transaction: skip, refuse, or write refs/bench/discarded/..., then DeleteBranchExact ]
                        ◀ tests attach here: cleanupBoundary steps before and after the write; the clock join

    trigger: session start (bench resume)
        │
        ▼
    refs/bench/discarded/<date>/...  ──▶  [ sweepLifecycleRefs with the resume instant ]  ──▶  swept count
                                            ◀ tests attach here: planted refs, fixed instant

    trigger: bench status / bench spec retire <slug>
        │
        ▼
    class counts, ledger  ──▶  [ appendGit / the spec retire dispatch in cmd/bench ]  ──▶  git row, candidate lines
                                  ◀ tests attach here: producible-signal table; dispatcher output table

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| RI1 | 1 | A plan over one assignment branch and one shift branch, both unrecorded, prints two rows | `internal/worktree/clean_unclaimed_test.go` (`TestPlanUnclaimedShiftResidueBranch`) | A namespace the selector drops leaves one row |
| RI2 | 2 | A recorded active branch with one unique commit prints no row | `internal/worktree/clean_classes_test.go` (`TestClassifyUnclaimedRefsOverTheEdgeInventory`) | A dropped record filter offers live work |
| RI3 | 3 | A checked-out unrecorded branch and the default branch print no row | `internal/worktree/clean_unclaimed_test.go` (`TestPlanUnclaimedAssignmentSetExcludesClaimedCheckedOutAndForeignRefs`) and `internal/worktree/clean_unclaimed_test.go` (`TestPlanUnclaimedAssignmentSetExcludesDefaultBranchInAssignmentNamespace`) | A dropped checkout filter offers a live checkout |
| RI4 | 4 | A branch at the `main` tip prints `class=landed` and action `discard-remove` | `internal/worktree/clean_classes_test.go` (`TestClassifyUnclaimedRefsOverTheEdgeInventory`) | A classifier that reads content alone misses the ancestor proof |
| RI5 | 5 | A branch whose two commits were squash-folded into `main` prints `class=landed` | `internal/worktree/clean_classes_test.go` (`TestClassifyUnclaimedRefsOverTheEdgeInventory`) | A classifier that stops at ancestry or cherry prints `unique` |
| RI6 | 6 | Two unrecorded branches at one unique tip print the lexically first as `class=unique` and the second as `class=subsumed holder=<first>` | `internal/worktree/clean_classes_test.go` (`TestClassifyUnclaimedRefsOverTheEdgeInventory`) | A rule that marks both subsumed deletes the tip |
| RI7 | 7 | A branch one commit under a unique ref prints `class=subsumed holder=<that ref>` | `internal/worktree/clean_classes_test.go` (`TestClassifyUnclaimedRefsOverTheEdgeInventory`) | A rule that reads equality alone prints `unique` |
| RI8 | 8 | A branch one commit under an active recorded branch prints `class=subsumed holder=<active branch>` | `internal/worktree/clean_classes_test.go` (`TestClassifyUnclaimedRefsOverTheEdgeInventory`) | A rule that ignores recorded branches prints `unique` |
| RI9 | 9 | A branch one commit under an ancestry-landed ref prints `class=landed` without a `holder=` field | `internal/worktree/clean_classes_test.go` (`TestClassifyUnclaimedRefsOverTheEdgeInventory`) | A rule that prints `class=subsumed holder=<landed ref>` fails the assertion |
| RI10 | 10 | A branch with one commit `main` lacks and no holder prints `class=unique` and action `retain` | `internal/worktree/clean_classes_test.go` (`TestClassifyUnclaimedRefsOverTheEdgeInventory`) | The old sweep prints `discard-remove` for this row |
| RI11 | 11 | A loose ref file under `.git/refs/heads/bench/assign/` that names a blob prints action `error` with the ref and `blob` in its detail, and no fingerprint | `internal/worktree/clean_unclaimed_test.go` (`TestCleanUnclaimedErrorRowRefusesTheSet`) | A classifier that skips resolution offers the set |
| RI12 | 12 | Two consecutive plans over one repository print identical rows and change no ref | `internal/worktree/clean_unclaimed_test.go` (`TestCleanUnclaimedPlanIsReadOnlyAndQuietWhenEmpty`) | A plan that writes a ref changes the second output |
| RI13 | 13 | Every classified row's detail cell begins with `class=` | `internal/worktree/clean_classes_test.go` (`TestClassifyUnclaimedRefsOverTheEdgeInventory`) | A row without the prefix hides the class |
| RI14 | 14 | A subsumed row's detail contains ` holder=refs/heads/` | `internal/worktree/clean_classes_test.go` (`TestClassifyUnclaimedRefsOverTheEdgeInventory`) | A row without the holder cannot be verified |
| RI15 | 15 | A landed row and a subsumed row both carry action `discard-remove` | `internal/worktree/clean_classes_test.go` (`TestClassifyUnclaimedRefsOverTheEdgeInventory`) | A subsumed row marked retain never retires |
| RI16 | 16 | A unique row's detail ends with `bench worktree clean --discard-branch --target <branch path>`, using the row's ref without `refs/heads/` | `planned` | An id selector can select a recorded label or refuse branches that share an id segment |
| RI17 | 17 | A plan, then one new commit on a subsumed ref, then an apply with the old fingerprint refuses as stale | `internal/worktree/clean_unclaimed_test.go` (`TestCleanUnclaimedStaleClassRefusesTheOldPlan`) | A fingerprint without the tip accepts the moved ref |
| RI18 | 18 | An apply with the fingerprint over a landed, a subsumed, and a unique row removes two refs and keeps the unique ref | `internal/worktree/clean_unclaimed_test.go` (`TestCleanUnclaimedBulkSweepKeepsUniqueRefs`) | The old apply loop deletes every row |
| RI19 | 19 | `--apply-current` over the same three rows prints the plan, removes two refs, and keeps the unique ref | `internal/worktree/clean_unclaimed_test.go` (`TestCleanUnclaimedBulkSweepKeepsUniqueRefs`) and `internal/worktree/clean_unclaimed_test.go` (`TestCleanUnclaimedApplyCurrent`) | The old shortcut deletes the unique ref |
| RI20 | 20 | A stale apply prints the re-plan command `bench worktree clean --discard-branch --unclaimed` with no apply flag | `internal/worktree/clean_set_outcomes_test.go` (`TestCleanSetUnclaimedStaleReplanAction`) | A re-plan that carries an apply flag skips the read |
| RI21 | 21 | After an apply over an equal-tip pair, a chain, and a ref under an ancestry-landed ref, every deleted tip stays reachable from `main` or from a surviving ref | `internal/worktree/clean_unclaimed_test.go` (`TestCleanUnclaimedBulkSweepKeepsUniqueRefs`) | A holder rule omission deletes both members of a pair |
| RI22 | 22 | A plan over unique rows only prints no apply help action | `internal/worktree/clean_unclaimed_test.go` (`TestCleanUnclaimedPlanNamesApplyOnlyWhenARowRemoves`) | An apply action beside retained rows invites a no-op or a mistake |
| RI23 | 23 | A plan over a repository with no Bench-namespace branch prints the empty table and exits 0 | `internal/worktree/clean_unclaimed_test.go` (`TestCleanUnclaimedPlanIsReadOnlyAndQuietWhenEmpty`) | A refusal on an empty set breaks the status route |
| RI24 | 24 | Status over one landed, one subsumed, and one unique ref prints `1 landed ref, 1 subsumed ref, 1 unique ref` | `internal/status/status_producible_test.go` (`TestAllProducibleBoardActionsAreInvocableOrEmpty`) | The old detail counts one number |
| RI25 | 25 | The git row action equals `bench worktree clean --discard-branch --unclaimed` | `internal/status/status_producible_test.go` (`TestAllProducibleBoardActionsAreInvocableOrEmpty`) | The old action carries `--apply-current` |
| RI26 | 26 | Status over one dirty path and one unique ref prints both details and routes to `bench worktree clean --discard-branch --unclaimed` | `internal/status/status_producible_test.go` (`TestAllProducibleBoardActionsAreInvocableOrEmpty`) | A precedence that routes the dirty path first hides the inventory |
| RI27 | 27 | A repository with one unique feature branch and no Bench-namespace ref prints `1 unique branch` and `git push` | `internal/status/status_producible_test.go` (`TestAllProducibleBoardActionsAreInvocableOrEmpty`) | A class count that fires on zero changes today's row |
| RI28 | 28 | The routed status command over a unique ref exits 0 and leaves the ref | `internal/systemtest/status_route_converge_test.go` (`TestStatusRouteExecutesUnclaimedCleanup`) | The old route deletes the ref |
| RI29 | 29 | `--target refs/heads/bench/assign/<owner>/<id>`, `--target bench/assign/<owner>/<id>`, and `--target bench/shift-<stamp>` each plan one row for the unrecorded branch | `planned` | The path check refuses the operand before the fallback |
| RI30 | 30 | `--target <id>` with one unrecorded branch whose last segment is `<id>` plans that row | `planned` | A resolver that reads the ledger alone refuses |
| RI31 | 31 | `--target <id>` with two unrecorded branches whose last segment is `<id>` prints an error row that names both refs and no fingerprint | `planned` | A first-match resolver discards the wrong branch |
| RI32 | 32 | `--target archive/x` for a branch outside both namespaces prints the refusal `relative path targets are unsupported` | `planned` | A fallback over every operand selects foreign work |
| RI33 | 33 | The explicit plan row of a unique unrecorded target shows `class=unique` and recovery `refs/bench/discarded/<yyyymmdd>/bench/assign/<owner>/<id>` | `planned` | A row without the planned ref hides the preservation |
| RI34 | 34 | `--target <id> --apply-current` exits 2 with the usage line | `internal/worktree/clean_set_command_test.go` (`TestCleanSetGrammar`) | A one-call discard of a unique ref skips the read |
| RI35 | 35 | After the apply, the discarded ref resolves to the old tip and the branch is gone | `planned` | A delete-first order leaves no handle |
| RI36 | 36 | A fault at the step before the write leaves the branch, and `for-each-ref refs/bench/discarded/` prints nothing | `planned` | A write before the boundary leaves residue on refusal |
| RI37 | 37 | A fault at the step after the write leaves the branch and the discarded ref, and a re-plan plus second apply removes the branch | `planned` | A second apply that refuses its own ref strands the retry |
| RI38 | 38 | A planted discarded ref at the planned path at another commit makes the apply refuse and keep the branch | `planned` | An unconditional write overwrites a handle |
| RI39 | 39 | The outcome row's recovery cell equals the discarded ref | `planned` | A row with recovery `none` hides the handle |
| RI40 | 40 | `--target <id>` for a landed unrecorded branch applies with recovery `none` and writes no discarded ref | `planned` | A ref for a landed row fills the namespace with noise |
| RI41 | 41 | `--target <label>` for a recorded active assignment plans through the record and prints no `class=` prefix | `planned` | A fallback that shadows the record breaks the release route |
| RI42 | 42 | The discarded ref path equals `refs/bench/discarded/` plus the UTC date plus `/` plus the branch path without `refs/heads/` | `internal/worktree/reconcile_test.go` (`TestDiscardedRefNamesTheDateAndTheBranchPath`) | A path without the date defeats the sweep |
| RI43 | 43 | A ref dated D is deleted at the instant D plus 30 days 00:00:00Z, and the swept count is 1 | `internal/worktree/reconcile_test.go` (`TestSweepDeletesADiscardedRefAtThirtyDays`) | A sweep that skips the namespace leaves the ref |
| RI44 | 44 | A ref dated D survives at the instant D plus 29 days 23:59:59Z, and the swept count is 0 | `internal/worktree/reconcile_test.go` (`TestSweepDeletesADiscardedRefAtThirtyDays`) | An off-by-one sweep removes a live handle |
| RI45 | 45 | A planted ref dated today survives the lifecycle emptying pass while a planted recovery ref is deleted | `internal/worktree/reconcile_test.go` (`TestSweepKeepsATodayDiscardedRefWhileItEmptiesRecovery`) | A namespace in the emptying list loses the handle |
| RI46 | 46 | A planted ref whose date segment is `latest` survives and the swept count is 0 | `internal/worktree/reconcile_test.go` (`TestSweepKeepsADiscardedRefWithAnUnparseableDate`) | A parser that treats a bad date as old deletes it |
| RI47 | 47 | A ref moved between the listing and the delete stays, and the sweep reports an error | `internal/worktree/reconcile_test.go` (`TestSweepRefusesADiscardedRefMovedAfterListing`) | A delete without the listed object removes the moved ref |
| RI93 | 47 | A discarded symref dated past the window leaves its target ref in place after the sweep | `internal/worktree/reconcile_test.go` (`TestSweepDeletesADiscardedSymrefAndNotItsTarget`) | A delete that follows the symref removes a branch outside the namespace |
| RI94 | 46 | A planted ref whose date segment is `20200101x` survives and the swept count is 0 | `internal/worktree/reconcile_test.go` (`TestSweepKeepsADiscardedRefWithAnUnparseableDate`) | A parser that reads the first eight bytes deletes the ref |
| RI95 | 42 | The path function renders the UTC date for an instant whose local date differs from its UTC date | `internal/worktree/reconcile_test.go` (`TestDiscardedRefNamesTheDateAndTheBranchPath`) | A renderer that formats the local date names the wrong day |
| RI48 | 48 | Retire of slug `s` with one active assignment labelled `s-build` prints one `superseded candidate:` line with its id | `planned` | A retire without the listing leaves the assignment with no route |
| RI49 | 49 | The candidate line ends with `bench worktree clean --discard-branch --target <assignment id>` | `planned` | A line without the command needs a lookup |
| RI50 | 50 | Retire with two unique unrecorded refs prints `unique refs: 2 — bench worktree clean --discard-branch --unclaimed` | `planned` | A retire without the count hides the inventory |
| RI51 | 51 | Retire with a candidate exits 0 and every ref survives | `planned` | A listing that discards turns retire destructive |
| RI52 | 52 | Retire with no matching assignment prints `unique refs: 0` and no candidate line | `planned` | A shape that omits the count line on zero is unstable |
| RI55 | 55 | A squash-folded sibling still prunes at the landing | `internal/worktree/land_prunes_landed_siblings_test.go` (`TestLandCommandPrunesSquashFoldedSiblingBranch`) | A shared-proof change that breaks the prune shows here |
| RI57 | 57 | The glossary term `unclaimed ref` names `refs/heads/bench/shift-` | `review-owned` | A term with the wrong namespace misleads a reader |
| RI58 | 58 | `--target <unrecorded id>` without `--discard-branch` plans action `retain` with `--discard-branch` in its detail, and the apply leaves the branch | `planned` | A discard without the flag skips the operator's assertion |
| RI59 | 8 | An unrecorded ref at exactly the tip of an active recorded branch prints `class=subsumed holder=<active branch>` | `internal/worktree/clean_classes_test.go` (`TestClassifyUnclaimedRefsOverTheEdgeInventory`) | An equality rule limited to unrecorded refs prints `unique` |
| RI60 | 10 | An unrecorded ref one commit under a checked-out foreign branch prints `class=unique` | `internal/worktree/clean_classes_test.go` (`TestClassifyUnclaimedRefsOverTheEdgeInventory`) | A holder set that admits every protected branch prints `subsumed` |
| RI61 | 7 | A chain A under B under C with C unique prints both A and B with `holder=<C>` | `internal/worktree/clean_classes_test.go` (`TestClassifyUnclaimedRefsOverTheEdgeInventory`) | A nearest-reaching rule names B for A |
| RI62 | 4 | Two unrecorded refs at the `main` tip both print `class=landed` | `internal/worktree/clean_classes_test.go` (`TestClassifyUnclaimedRefsOverTheEdgeInventory`) | A subsumed-first precedence prints `subsumed` for the second |
| RI63 | 17 | A plan, then a recorded active branch created at a descendant of the unique ref, then the old fingerprint refuses as stale with the tip unchanged | `internal/worktree/clean_unclaimed_test.go` (`TestCleanUnclaimedStaleClassRefusesTheOldPlan`) | A fingerprint without the class accepts the changed row |
| RI64 | 59 | An explicit plan, the same descendant record, then the old fingerprint refuses as stale | `planned` | An explicit fingerprint without the class accepts the changed row |
| RI65 | 59 | After the stale explicit refusal, `for-each-ref refs/bench/discarded/` prints nothing | `planned` | A write before the oracle leaves residue on refusal |
| RI66 | 22 | A plan with one landed row prints the help action `bench worktree clean --discard-branch --unclaimed --apply <fingerprint>` with the set fingerprint | `internal/worktree/clean_unclaimed_test.go` (`TestCleanUnclaimedPlanNamesApplyOnlyWhenARowRemoves`) | A plan without the action has no route to the apply |
| RI67 | 60 | A branch moved at the after-write fault step survives, the discarded ref stays at the old tip, and the row reads action `error` | `planned` | A delete without the exact tip removes the moved branch |
| RI68 | 40 | `--target <id>` for a subsumed unrecorded branch applies with recovery `none` and writes no discarded ref | `planned` | A ref for a subsumed row fills the namespace with noise |
| RI69 | 61 | Retire with an unreadable ledger prints `unique refs: unavailable — <error>` and exits 0 | `planned` | A listing failure that changes the exit code breaks the retire route |
| RI70 | 48 | Retire of slug `s` with one active assignment whose request token is `s-run` and whose label is `other` prints one candidate line | `planned` | A label-only match misses the token |
| RI71 | 48 | Retire of slug `s` with one complete assignment labelled `s-build` prints no candidate line | `planned` | A state-blind match lists retired work |
| RI72 | 11 | A plan with an error row exits 1 and prints no apply help action | `internal/worktree/clean_unclaimed_test.go` (`TestCleanUnclaimedErrorRowRefusesTheSet`) | An apply action beside an error row invites a partial sweep |
| RI73 | 19 | `--apply-current` over an error row and a landed row refuses before any delete, and the landed ref survives | `internal/worktree/clean_unclaimed_test.go` (`TestCleanUnclaimedErrorRowRefusesTheSet`) | A shortcut that skips the error row deletes around it |
| RI74 | 62 | `--target bench/shift-<stamp>` for a checked-out unrecorded shift branch prints an error row that names the checkout path and no fingerprint | `planned` | A fallback over every unrecorded branch offers a live checkout |
| RI75 | 62 | `--target bench/assign/<owner>/<id>` for a Bench-namespace default branch prints an error row and no fingerprint | `planned` | A fallback without the default filter offers the default branch |
| RI76 | 62 | `--target bench/assign/<owner>/<id>` for a recorded active branch plans through the record and prints no `class=` prefix | `planned` | A prefix match that shadows the record breaks the release route |
| RI77 | 48 | Retire with the operand `specs/s/spec.md` prints the same candidate lines as the operand `s` | `planned` | A wrapper that reads the operand as the slug lists nothing for a path |
| RI78 | 16 | A unique `bench/shift-<stamp>` row's detail ends with `bench worktree clean --discard-branch --target bench/shift-<stamp>`, and that command plans one row for it | `planned` | A shift row with the id form names a target that cannot resolve |
| RI79 | 50 | The candidate lines and the count line print after the `retired:` lines and before the `next:` line | `planned` | A listing after `next:` hides the remainder step |
| RI80 | 51 | `bench spec retire --help` prints no candidate line and no count line | `planned` | A wrapper that appends on every call pollutes the help |
| RI81 | 9 | A adds `x=1`, B changes `x` to `2` and is squash-folded into `main`: the plan prints B as `class=landed` and A as `class=unique` | `internal/worktree/clean_classes_test.go` (`TestClassifyUnclaimedRefsOverTheEdgeInventory`) | A holder set that admits a content-landed ref deletes A with no surviving handle |
| RI82 | 63 | The glossary term `holder` names a recorded assignment branch or a unique root, and no landed ref | `review-owned` | A term that admits a landed ref contradicts the rule |
| RI83 | 3, 21 | A resolving unrecorded symref in either Bench namespace produces action `error` with its name and `symref` in the detail, no class prefix, an empty set fingerprint, no apply action, and both apply forms refuse before any delete while the target and a removable sibling survive | `internal/worktree/clean_classes_test.go` (`TestCleanUnclaimedSymrefFailsClosed`) | Silent exclusion omits the required error row, and a dereferenced classification lets the apply delete the target |
| RI84 | 8, 9, 21 | Recorded R adds `x=1` then changes it to `x=2`, `main` receives a squash of R, and unrecorded A at R's first commit with no other holder prints `class=unique` and `retain`, for an active R and for a cleanup-pending R | `internal/worktree/clean_classes_test.go` (`TestClassifyUnclaimedRefsOverTheEdgeInventory`) | A holder set that admits a content-landed recorded R prints A as subsumed and authorizes its deletion before R's later retirement |
| RI85 | 17 | A plan, then `main` fast-forwarded onto the unique ref so that only its class changes, then the old fingerprint refuses as stale | `internal/worktree/clean_unclaimed_test.go` (`TestCleanUnclaimedStaleClassRefusesTheOldPlan`) | A fingerprint without the class accepts a row whose tip and holder are unchanged |
| RI86 | 10 | A branch one commit under a complete record's branch prints `class=unique` | `internal/worktree/clean_classes_test.go` (`TestClassifyUnclaimedRefsOverTheEdgeInventory`) | A holder filter that admits every record state prints `subsumed` |
| RI88 | 24 | Status over one unique ref and one faulted ref includes `1 unique ref, 1 faulted ref` and routes to `bench worktree clean --discard-branch --unclaimed` | `internal/status/status_producible_test.go` (`TestAllProducibleBoardActionsAreInvocableOrEmpty`) | The exact signal assertion fails if status hides the fault, folds it into another class, or changes the plan action |
| RI89 | 24 | Status over one unique Bench ref, a resolving Bench symref to that ref, and one unique feature branch prints exactly `1 unique ref, 1 faulted ref, 1 unique branch` | `internal/status/status_producible_test.go` (`TestAllProducibleBoardActionsAreInvocableOrEmpty`) | Counting the faulted symref again as an ordinary unique branch produces `2 unique branches` |
| RI90 | 26 | Status with an unreadable ledger prints `unclaimed refs unavailable` with action `bench worktree clean --discard-branch --unclaimed` | `internal/status/status_producible_test.go` (`TestAllProducibleBoardActionsAreInvocableOrEmpty`) | The ledger read forces a planner failure, so an omitted failure detail or a wrong action fails the exact signal assertion |
| RI91 | 26 | Status with only one faulted unclaimed symref to `main` prints `1 faulted ref` with action `bench worktree clean --discard-branch --unclaimed` | `internal/status/status_producible_test.go` (`TestAllProducibleBoardActionsAreInvocableOrEmpty`) | No other ref class or planner error can trigger the route, so a faulted count left out of the row total fails the action assertion |
| RI92 | 26 | Status over a blob-tip Bench ref beside an unreadable ledger prints `git state unavailable, unclaimed refs unavailable` with action `bench worktree clean --discard-branch --unclaimed` | `internal/status/status_producible_test.go` (`TestAllProducibleBoardActionsAreInvocableOrEmpty`) | A route that drops the planner error whenever Git state fails prints `git status` for a repository that still holds refs |
| RI87 | 8 | A branch one commit under a cleanup-pending recorded branch prints `class=subsumed holder=<that branch>` | `internal/worktree/clean_classes_test.go` (`TestClassifyUnclaimedRefsOverTheEdgeInventory`) | A holder filter limited to active records prints `unique` |
| RI96 | 35, 38 | A discarded-path symref pointing to the target branch makes the apply print an `error` row, keep both refs unchanged, and write nothing | `planned` | Accepting the dereferenced tip skips preservation and deletes the branch, leaving the recovery ref dangling |
| RI97 | 51 | A retire refusal whose operand contains a newline followed by `next: ` preserves its nonzero exit code and prints neither candidate nor count lines | `planned` | Removing the exit-code guard lets the operand's `next: ` line trigger the listing |
| RI98 | 29, 62 | An explicit target naming a resolving Bench symref prints an `error` row naming `symref`, no fingerprint, and both apply forms refuse without ref changes | `planned` | Ignoring the candidate's fault admits a dereferenced branch into the explicit discard plan |
| RI99 | 21, 59 | After an explicit set removes a recorded holder, its formerly subsumed unrecorded member reports a stale error and survives at its original tip | `planned` | Omitting the per-member requalification deletes the unrecorded branch after its holder disappears |
| RI100 | 50 | Retire with a resolving Bench symref to `main` and a Bench blob-tip ref prints exactly one count line, `unique refs: 0, 2 faulted — bench worktree clean --discard-branch --unclaimed` | `planned` | Ignoring faulted rows or counting them as unique changes the exact line |
| RI101 | 35, 38 | A direct ref planted at the planned path between the absent-ref read and the write makes the apply print an `error` row, keep both refs unchanged, and write nothing | `planned` | Removing the zero old value lets the write overwrite a handle that appeared after the read |

Not covered: story 53 — reviewed exclusion, and the review round confirms that no hold surface is added.
Not covered: story 54 — reviewed exclusion, and RI48 to RI52 show the listing discards nothing.
Not covered: story 56 — the build runs on fixtures only, and the reviewer runs the live plan after the landing.

### Edge inventory

- An unrecorded branch at the default branch tip: landed by ancestry, RI4, and two of them, RI62.
- A squash fold: landed by the reverse-apply proof, RI5.
- Equal tips: RI6, RI21, and RI59.
- A strict-ancestor chain of three refs: the root is the holder, RI61 and RI21.
- A ref under an active recorded branch: RI8.
- A ref under an ancestry-landed ref, landed by transitivity: RI9, and under a content-landed ref: RI81.
- A ref under a checked-out foreign branch: RI60.
- A ref under a complete record's branch, and under a cleanup-pending recorded branch: RI86 and RI87.
- A ref under a recorded branch that is landed by content only: RI84, unique.
- A resolving unrecorded symref in either Bench namespace: an error row, never a silent exclusion, RI83. Its detail names the ref and `symref`. The set has no fingerprint or apply action, and both apply forms refuse before any delete.
- A stale fingerprint after a class change alone: RI85.
- A ref that points at a blob: RI11, RI72, and RI73.
- A plan over an empty namespace: RI23.
- A stale fingerprint after a tip move: RI17, and after a class change: RI63, RI64, and RI65.
- An ambiguous id segment: RI31.
- A foreign branch as a target: RI32.
- A checkout, the default branch, and a recorded branch as a prefix target: RI74, RI75, RI76.
- A target without the discard flag: RI58.
- A unique shift row's explicit route: RI78.
- A retire path operand, the listing position, and the help form: RI77, RI79, RI80.
- A fault before the write, a fault after the write, a branch moved after the write, and a conflicting planted ref: RI36, RI37, RI67, RI38.
- The 30-day boundary on each side: RI43 and RI44.
- A malformed date segment: RI46.
- A concurrent ref move under the sweep: RI47.
- A retire listing failure, a token match, and a complete record: RI69, RI70, RI71.
- **Won't handle** a branch outside the two Bench namespaces — the landing prune still retires a landed foreign branch, and a unique one is the operator's.
- **Won't handle** a restore of a discarded ref into a new assignment — `bench worktree create --from` reads a commit, and a later spec decides the seam.
- **Won't handle** a remote-tracking ref — the landing prune reads local heads only, and Bench creates no remote ref.
- **Won't handle** an assignment-scoped recovery ref — its owner is the release path, and its sweep rule stays as it is.
- **Won't handle** a discarded ref as a holder — a discarded ref is never a holder, so the unique class and the explicit route stay.
- **Won't handle** a plan and an apply on different UTC days — the planned ref carries the plan's date, so the apply refuses as stale.
- **Won't handle** the calling worktree's own assignment — the retire listing includes it, and the explicit discard retains it on a live lease.
- **Won't handle** an unrecorded dangling symref in either Bench namespace, whose target ref does not exist. Git's ref enumeration skips it, so the unclaimed plan omits it and cannot delete through it. If its target appears before the apply recheck, the changed row set triggers the stale refusal. The unclaimed planner still rejects a resolving symref under RI83.

## Ownership fences

- `internal/worktree/clean_classes.go`
- `internal/worktree/clean_classes_test.go`
- `internal/worktree/clean_unclaimed.go`
- `internal/worktree/clean_unclaimed_test.go`
- `internal/worktree/clean_landed_apply_test.go`
- `internal/worktree/clean_set_apply.go`
- `internal/worktree/clean_set_apply_test.go`
- `internal/worktree/clean_set_command_test.go`
- `internal/worktree/clean_set_outcomes_test.go`
- `internal/worktree/clean_set_wiring_test.go`
- `internal/worktree/clean_set.go`
- `internal/worktree/clean_discard.go`
- `internal/worktree/clean_discard_test.go`
- `internal/worktree/joins.go`
- `internal/worktree/land_prunes_landed_siblings_test.go`
- `internal/worktree/path.go`
- `internal/worktree/worktree.go`
- `internal/worktree/reconcile.go`
- `internal/worktree/reconcile_test.go`
- `internal/worktree/resume_reconcile_test.go`
- `internal/intent/ledger/ledger.go`
- `internal/intent/ledger_aliases.go`
- `internal/status/status.go`
- `internal/status/status_producible_test.go`
- `internal/systemtest/status_route_converge_test.go`
- `internal/spec/spec.go`
- `internal/spec/spec_test.go`
- `internal/spec/history.go`
- `cmd/bench/main.go`
- `cmd/bench/spec_retire_listing.go`
- `cmd/bench/spec_retire_listing_test.go`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry_decision_maps.go`
- `internal/anchors/registry_decision_maps_test.go`
- `tests/canary/package-core-guard/unrouted-subcommand`
- `tests/canary/docs-currency-token-diet/signal-vocabulary-drift`
- `tests/canary/workflow-guidance-anchors/context-acceptance-row-vocabulary`
- `tests/canary/workflow-guidance-anchors/context-coverage-map-term`
- `tests/canary/workflow-guidance-anchors/context-coverage-row-parts`
- `tests/canary/workflow-guidance-anchors/context-coverage-row-vocabulary`
- `tests/canary/workflow-guidance-anchors/context-decision-map-term`
- `tests/canary/workflow-guidance-anchors/context-reader-sweep-term`
- `tests/canary/workflow-guidance-anchors/context-ticket-vocabulary`
- `CONTEXT.md`
- `reviews/ref-inventory.md`

## Out of scope

- A restore verb that turns a discarded ref into a new assignment: 4 edits, 2 gate runs.
- A hold marker, keep list, or hold verb: reviewed exclusion, no estimate.
- A class column in the shared cleanup table for every mode: 6 edits, 2 gate runs, and it moves every cleanup fixture.
- A cache for the landed proofs under `bench status`: 3 edits, 1 gate run. It waits for the RI-C1b measurement.
- A paragraph on the three classes in the operating reference: 1 edit, 1 gate run, and 20 fixture pins join the fence.
- The disposition of the seven cited tips and the three FT336 drafts: FT346, FT347, and FT348 own it.

## Further notes

Flagged additions beyond the decision source:

- The `class=` and `holder=` detail spelling, the `retain` action for a unique row, and the fingerprint version bump are spec-writer discretion under the map.
- The candidate and count line spellings of `bench spec retire`, and the `unique refs: unavailable` line, are spec-writer discretion under the map.
- The unresolvable-ref error row, RI11, and its two apply rows, RI72 and RI73, are fail-closed edges the map did not name.
- The fallback grammar for `--target`, and the `retain` row without `--discard-branch`, RI58, are spec-writer discretion under ticket 6 of the map.
- The conflicting planted ref refusal, RI38, is the fail-closed side of the write rule the map's ticket 3 fixed.
- The malformed date segment rule, RI46, is the fail-closed side of the sweep rule the map's ticket 3 fixed.
- The glossary namespace repair, RI57, corrects a landed defect in the term the map introduced.

Map ticket 1 of the decision source spells the shift namespace as `refs/heads/bench/shift/`.
The ledger declares `refs/heads/bench/shift-`, and this spec uses the ledger's spelling.
The compiled map stays as the reviewer confirmed it, and this note flags the typo for reviewer veto.

Decision ticket 7 of the map names "a landed ref" as a holder.
By reviewer decision on 2026-09-27 that phrase first narrowed to an ancestry-landed ref, because a content-landed ref leaves no handle for a ref beneath it.
The ticket 1 author then found that the narrowed kind is unreachable.
Ancestry is transitive, so every ref beneath an ancestry-landed ref is itself landed, and the landed-first precedence never names that holder.
By a second reviewer decision on 2026-09-27, through the Codex Astra route, no landed ref is a holder.
The compiled map stays untouched, RI9 pins the transitive landed class, and RI81 pins the content-landed case.

The RI-C1a chunk review added five rows by reviewer decision on 2026-09-27, through the Codex Astra route.
RI83 fails a symbolic ref closed, because Git dereferences a symref on delete.
The Coverage axis observed the deletion of a unique root and of `main` through such a ref.
RI84 excludes a content-landed recorded branch from the holders, because its later retirement would leave a subsumed ref with no handle.

RI85, RI86, and RI87 pin the class-only fingerprint binding and the two recorded states, because the review probes at those sites stayed silent.
RI5 reads two commits, because a one-commit squash lands by patch containment.
RI13 and story 13 read classified rows, because an error row carries no class prefix.

The confirming round narrowed the symref rule to a resolving symref, by reviewer decision on 2026-09-27 through the same route.
Git's ref enumeration skips a dangling symref, so the plan cannot list it or delete through it, and the Won't-handle list names that case.
The same round corrected the plan-only cost sentence, because the classifier now proves each active or cleanup-pending recorded branch before it can hold.

The RI-C1b build added row RI88 by reviewer decision on 2026-09-27, through the same route.
The status row shows a faulted ref count and a planner failure.
A hidden fault would leave the operator with a green row over a set that the plan refuses.
The same decision defers the `main` composition to the landing preparation, with its own review round.
The review chain refuses a `main` merge after the first chunk.

The RI-C1b review added rows RI89 to RI91 and two sentences by reviewer decision on 2026-09-27, through the same route.
A blob-tip ref made the Git state read fail before the class counts.
The row now keeps its details and its plan route beside `git state unavailable`.
A symref to a unique ref was counted twice, so RI89 pins the single count.
The plan command spelling moves to one exported source in the worktree package, because status and the retire listing print it too.

The repair then found that a non-repository fails both reads, so the reviewer route separated that case from a repository.
RI92 pins a blob-tip ref beside an unreadable ledger in a repository.

The RI-C2a review added rows RI93 to RI95 and one sentence, by orchestrator plan expansion inside the approved behavior on 2026-09-27.
A discarded symref made the sweep's delete follow it to a branch outside the namespace, so the sweep's own delete never dereferences.
That change touches only the sweep in `reconcile.go`, and the earlier decision on `DeleteBranchExact` stands; the reviewer can veto it.
RI94 pins a date segment with extra bytes, and RI95 pins the UTC date for a local instant.

The RI-C2b review added rows RI96 to RI101 and moved the printed discard route to the branch path.
The reviewer decided both on 2026-09-27 through the Codex Astra route.
A symref at the planned discarded path made the apply skip the write and delete the branch, so the apply refuses a symbolic planned path.
An id route could select a recorded label or refuse two branches that share an id segment.
Story 16 and RI16 therefore name the branch path.

The retire count line gains a faulted suffix, because a symref or blob-tip ref was invisible in the unique count.
A third boundary step sits between the read of the planned path and the write, so a competing direct ref has a reachable window.

### Completion plan

```bench-completion-plan
{
  "version": 2,
  "execution": {
    "mode": "delegate",
    "run_id": "ft199-ref-inventory-full-20260927",
    "orchestrator_session": "claude:session_0156tkEZcRSowaafegWfFZJP",
    "author_limit": 1,
    "assignments": {
      "6-repair-glossary-shift-namespace.md": [
        {
          "session": "claude:bench-writer/ri-t6-author",
          "assignment": "ri-t6-author",
          "model": "opus",
          "effort": "high",
          "source": "72a749a35dc4c37de87f56534959ac9c137099e1",
          "native_ref": "claude:agent/ri-t6-author-20260927@72a749a35dc4c37de87f56534959ac9c137099e1"
        }
      ],
      "1-classify-unclaimed-refs.md": [
        {
          "session": "claude:bench-writer/ri-t1-author",
          "assignment": "ri-t1-author",
          "model": "opus",
          "effort": "high",
          "source": "f3ff2543aa7e40b07fd481fe4286161402d96d2b",
          "native_ref": "claude:agent/ri-t1-author-20260927@f3ff2543aa7e40b07fd481fe4286161402d96d2b"
        },
        {
          "session": "claude:bench-writer/ri-t1-repair-1",
          "assignment": "ri-t1-repair-1",
          "model": "opus",
          "effort": "high",
          "source": "c6b2a74effa57cbafdd8bc449c9fbbb24d1468eb",
          "native_ref": "claude:agent/ri-t1-repair-1-20260927@c6b2a74effa57cbafdd8bc449c9fbbb24d1468eb",
          "predecessor": "claude:bench-writer/ri-t1-author",
          "trigger": "user-directed",
          "stopped": "ri-t1-author returned its final report and idle notification after the chunk-tip verification rerun at c6b2a74e; the worktree was clean and no further write came from it",
          "preserved": "2ace6d83cbe1430ea443cdee7743d081e6a6aa54 on bench/assign/9cd9510fff4093f7f9f4456f6a029560/3a0fa26e2c3c38179f908f3636fb07ed, chunk tip c6b2a74effa57cbafdd8bc449c9fbbb24d1468eb"
        },
        {
          "session": "claude:bench-writer/ri-t1-repair-2",
          "assignment": "ri-t1-repair-2",
          "model": "opus",
          "effort": "high",
          "source": "4b80686a12ea46afa0e6c0ba4fe9aab2e00935f6",
          "native_ref": "claude:agent/ri-t1-repair-2-20260927@4b80686a12ea46afa0e6c0ba4fe9aab2e00935f6",
          "predecessor": "claude:bench-writer/ri-t1-repair-1",
          "trigger": "user-directed",
          "stopped": "ri-t1-repair-1 returned its final report and idle notification after the chunk-tip verification rerun at 4b80686a; the worktree was clean and no further write came from it",
          "preserved": "4b80686a12ea46afa0e6c0ba4fe9aab2e00935f6 on bench/assign/9cd9510fff4093f7f9f4456f6a029560/3a0fa26e2c3c38179f908f3636fb07ed, the RI-C1a repair 1 tip"
        }
      ],
      "2-route-status-to-the-plan.md": [
        {
          "session": "claude:bench-writer/ri-t2-author",
          "assignment": "ri-t2-author",
          "model": "opus",
          "effort": "high",
          "source": "09611663f8a49bd5f037b24ba9385eaa2e42be41",
          "native_ref": "claude:agent/ri-t2-author-20260927@09611663f8a49bd5f037b24ba9385eaa2e42be41"
        },
        {
          "session": "claude:bench-writer/ri-t2-repair-1",
          "assignment": "ri-t2-repair-1",
          "model": "opus",
          "effort": "high",
          "source": "f12046941a0ad1f2812e82de0bf4e9c61303079b",
          "native_ref": "claude:agent/ri-t2-repair-1-20260927@f12046941a0ad1f2812e82de0bf4e9c61303079b",
          "predecessor": "claude:bench-writer/ri-t2-author",
          "trigger": "user-directed",
          "stopped": "ri-t2-author returned its final report and idle notification after the chunk-tip verification rerun at f1204694; the worktree was clean and no further write came from it",
          "preserved": "0e8f29cbb4e0fb4159baa23c34cddea9ffbd1cf1 on bench/assign/9cd9510fff4093f7f9f4456f6a029560/3a0fa26e2c3c38179f908f3636fb07ed, chunk tip f12046941a0ad1f2812e82de0bf4e9c61303079b"
        },
        {
          "session": "claude:bench-writer/ri-t2-repair-2",
          "assignment": "ri-t2-repair-2",
          "model": "opus",
          "effort": "high",
          "source": "f0a8a3c6f3b0e4e70edacbda5d0a7ac9d80fc246",
          "native_ref": "claude:agent/ri-t2-repair-2-20260927@f0a8a3c6f3b0e4e70edacbda5d0a7ac9d80fc246",
          "predecessor": "claude:bench-writer/ri-t2-repair-1",
          "trigger": "user-directed",
          "stopped": "ri-t2-repair-1 returned its final report and idle notification after its commit f0a8a3c6; the worktree was clean and no further write came from it",
          "preserved": "f0a8a3c6f3b0e4e70edacbda5d0a7ac9d80fc246 on bench/assign/9cd9510fff4093f7f9f4456f6a029560/3a0fa26e2c3c38179f908f3636fb07ed, the RI-C1b repair 1 tip"
        }
      ],
      "3-sweep-discarded-refs.md": [
        {
          "session": "claude:bench-writer/ri-t3-author",
          "assignment": "ri-t3-author",
          "model": "opus",
          "effort": "high",
          "source": "c6d2cfbf66d4b82f284e008ab063266bf61c4a23",
          "native_ref": "claude:agent/ri-t3-author-20260927@c6d2cfbf66d4b82f284e008ab063266bf61c4a23"
        },
        {
          "session": "claude:bench-writer/ri-t3-repair-1",
          "assignment": "ri-t3-repair-1",
          "model": "opus",
          "effort": "high",
          "source": "5c54f668f2d79291e447dda034c78a2046d9d999",
          "native_ref": "claude:agent/ri-t3-repair-1-20260927@5c54f668f2d79291e447dda034c78a2046d9d999",
          "predecessor": "claude:bench-writer/ri-t3-author",
          "trigger": "user-directed",
          "stopped": "ri-t3-author returned its final report and idle notification after the chunk-tip verification reruns at 5c54f668; the worktree was clean and no further write came from it",
          "preserved": "5c54f668f2d79291e447dda034c78a2046d9d999 on bench/assign/9cd9510fff4093f7f9f4456f6a029560/3a0fa26e2c3c38179f908f3636fb07ed, the RI-C2a chunk tip"
        }
      ],
      "4-discard-a-unique-ref-by-target.md": [
        {
          "session": "claude:bench-writer/ri-t4-author",
          "assignment": "ri-t4-author",
          "model": "opus",
          "effort": "high",
          "source": "114f94a2d6541d11833af640e5a886cbe8d01966",
          "native_ref": "claude:agent/ri-t4-author-20260927@114f94a2d6541d11833af640e5a886cbe8d01966"
        },
        {
          "session": "claude:bench-writer/ri-t4-repair-1",
          "assignment": "ri-t4-repair-1",
          "model": "opus",
          "effort": "high",
          "source": "7ccd9aaab3a0f4be51d8b2bab0041f8e8ac28835",
          "native_ref": "claude:agent/ri-t4-repair-1-20260927@7ccd9aaab3a0f4be51d8b2bab0041f8e8ac28835",
          "predecessor": "claude:bench-writer/ri-t4-author",
          "trigger": "user-directed",
          "stopped": "ri-t4-author returned its final report and idle notification after the chunk-tip verification rerun at 7ccd9aaa; the worktree was clean and no further write came from it",
          "preserved": "505a4c659f8000ae0a088f42038bc69fbfc782ef on bench/assign/9cd9510fff4093f7f9f4456f6a029560/3a0fa26e2c3c38179f908f3636fb07ed, chunk tip 7ccd9aaab3a0f4be51d8b2bab0041f8e8ac28835"
        }
      ],
      "5-list-retire-candidates.md": [
        {
          "session": "claude:bench-writer/ri-t5-author",
          "assignment": "ri-t5-author",
          "model": "opus",
          "effort": "high",
          "source": "505a4c659f8000ae0a088f42038bc69fbfc782ef",
          "native_ref": "claude:agent/ri-t5-author-20260927@505a4c659f8000ae0a088f42038bc69fbfc782ef"
        },
        {
          "session": "claude:bench-writer/ri-t5-repair-1",
          "assignment": "ri-t5-repair-1",
          "model": "opus",
          "effort": "high",
          "source": "7ccd9aaab3a0f4be51d8b2bab0041f8e8ac28835",
          "native_ref": "claude:agent/ri-t5-repair-1-20260927@7ccd9aaab3a0f4be51d8b2bab0041f8e8ac28835",
          "predecessor": "claude:bench-writer/ri-t5-author",
          "trigger": "user-directed",
          "stopped": "ri-t5-author returned its final report and idle notification after the chunk-tip verification rerun at 7ccd9aaa; the worktree was clean and no further write came from it",
          "preserved": "7ccd9aaab3a0f4be51d8b2bab0041f8e8ac28835 on bench/assign/9cd9510fff4093f7f9f4456f6a029560/3a0fa26e2c3c38179f908f3636fb07ed, the RI-C2b chunk tip"
        }
      ]
    }
  },
  "chunks": [
    {
      "id": "RI-C1a",
      "tickets": ["1-classify-unclaimed-refs.md", "6-repair-glossary-shift-namespace.md"],
      "verification": [
        {"id": "6-anchors", "command": "bench test --package ./internal/anchors", "ticket": "6-repair-glossary-shift-namespace.md"},
        {"id": "6-conformance", "command": "bench test --package ./internal/conformance", "ticket": "6-repair-glossary-shift-namespace.md"},
        {"id": "1-worktree", "command": "bench test --package ./internal/worktree", "ticket": "1-classify-unclaimed-refs.md"},
        {"id": "1-status", "command": "bench test --package ./internal/status", "ticket": "1-classify-unclaimed-refs.md"},
        {"id": "1-cmd", "command": "bench test --package ./cmd/bench", "ticket": "1-classify-unclaimed-refs.md"},
        {"id": "1-conformance", "command": "bench test --package ./internal/conformance", "ticket": "1-classify-unclaimed-refs.md"},
        {"id": "1-system", "command": "bench test --check system", "ticket": "1-classify-unclaimed-refs.md"}
      ]
    },
    {
      "id": "RI-C1b",
      "tickets": ["2-route-status-to-the-plan.md"],
      "verification": [
        {"id": "2-status", "command": "bench test --package ./internal/status", "ticket": "2-route-status-to-the-plan.md"},
        {"id": "2-worktree", "command": "bench test --package ./internal/worktree", "ticket": "2-route-status-to-the-plan.md"},
        {"id": "2-system", "command": "bench test --check system", "ticket": "2-route-status-to-the-plan.md"}
      ]
    },
    {
      "id": "RI-C2a",
      "tickets": ["3-sweep-discarded-refs.md"],
      "verification": [
        {"id": "3-worktree", "command": "bench test --package ./internal/worktree", "ticket": "3-sweep-discarded-refs.md"},
        {"id": "3-ledger", "command": "bench test --package ./internal/intent/ledger", "ticket": "3-sweep-discarded-refs.md"}
      ]
    },
    {
      "id": "RI-C2b",
      "tickets": ["4-discard-a-unique-ref-by-target.md", "5-list-retire-candidates.md"],
      "verification": [
        {"id": "4-worktree", "command": "bench test --package ./internal/worktree", "ticket": "4-discard-a-unique-ref-by-target.md"},
        {"id": "4-cmd", "command": "bench test --package ./cmd/bench", "ticket": "4-discard-a-unique-ref-by-target.md"},
        {"id": "5-spec", "command": "bench test --package ./internal/spec", "ticket": "5-list-retire-candidates.md"},
        {"id": "5-cmd", "command": "bench test --package ./cmd/bench", "ticket": "5-list-retire-candidates.md"}
      ]
    }
  ],
  "final_verification": [
    {"id": "coverage", "command": "bench coverage --check specs/ref-inventory/spec.md"},
    {"id": "worktree", "command": "bench test --package ./internal/worktree"},
    {"id": "status", "command": "bench test --package ./internal/status"},
    {"id": "spec", "command": "bench test --package ./internal/spec"},
    {"id": "cmd", "command": "bench test --package ./cmd/bench"},
    {"id": "conformance", "command": "bench test --package ./internal/conformance"},
    {"id": "anchors", "command": "bench test --package ./internal/anchors"},
    {"id": "system", "command": "bench test --check system"}
  ]
}
```

Source-sentence-to-row table:

| map ticket sentence | rows |
| --- | --- |
| The inventory covers every branch under the two namespaces | RI1 |
| A recorded branch classifies as active and stays protected | RI2, RI3, RI83 |
| An unrecorded branch classifies as landed under the four proofs | RI4, RI5, RI55, RI62 |
| An unrecorded branch classifies as subsumed under a holder | RI6, RI7, RI8, RI9, RI59, RI61 |
| Every other unrecorded branch classifies as unique | RI10, RI60 |
| The bulk sweep discards landed and subsumed rows only | RI15, RI18, RI19, RI21 |
| A unique row never enters a bulk discard | RI10, RI16, RI18, RI22 |
| No marker, the class is the keep | Not covered: story 53 |
| The discard writes the dated ref at the exact tip before the delete | RI35, RI42, RI67, RI95 |
| The outcome row names the discarded ref | RI39 |
| The sweep removes a ref 30 days after its date and reports the count | RI43, RI44, RI47, RI93, RI94 |
| The emptying rule does not apply to the discarded namespace | RI45 |
| A superseded ref classifies as unique and needs an explicit discard | RI10, RI33 |
| Status counts the three classes and routes to the plan-only form, above the dirty-path route | RI24, RI25, RI26, RI27, RI88, RI89, RI90, RI91, RI92 |
| Status never names a destructive command | RI25, RI28 |
| The plan output alone names the apply command, and only when a landed or subsumed row exists | RI22, RI66 |
| `--target` resolves an unrecorded branch by id segment or path | RI29, RI30, RI31, RI32, RI58, RI74, RI75, RI76, RI78, RI98 |
| The explicit plan shows class unique and the planned ref, and the apply needs the fingerprint | RI33, RI34 |
| A class change between plan and apply refuses as stale | RI17, RI63, RI64, RI65, RI85 |
| The preserve step and the discard step never split | RI36, RI37, RI38, RI96, RI99, RI101 |
| Equal tips make the lexically first ref the root | RI6 |
| A holder is an active branch, a unique root, or a landed ref (narrowed: no landed ref holds) | RI7, RI8, RI9, RI59, RI60, RI61, RI81, RI82, RI84, RI86, RI87 |
| A subsumed row names its holder and writes no discarded ref | RI14, RI40, RI68 |
| `--apply-current` stays, narrows, and leaves the status table | RI19, RI25, RI73 |
| Retire lists recorded rows by slug with the discard command | RI48, RI49, RI52, RI70, RI71, RI77 |
| Retire adds the unique count line and discards nothing | RI50, RI51, RI69, RI79, RI80, RI97, RI100 |
| One spec, two chunks, A before B | the chunk table |
| Fixtures prove the classes, the live refs stay | Not covered: story 56 |

Pre-review proof checklist:

- `Cited symbols`: each symbol below resolves in the tree at `534a8c69`.
  - `planUnclaimedAssignmentSet`, `applyUnclaimedAssignmentSet`, `UnclaimedAssignmentBranchRefs`, `cleanCommandWith`, `parseCleanSelection`
  - `planExplicitSet`, `resolveAssignmentIn`, `selectAssignment`, `targetPath`, `errTargetUnassigned`, `StepUnlockedReplan`
  - `git.LandedInDefault`, `git.DeleteBranchExact`, `sweepLifecycleRefs`, `reconcileLifecycleDebris`, `lifecycleRefNamespaces`
  - `intent.RecoveryRefNamespace`, `intent.ResetRefNamespace`, `appendGit`, `cleanUnclaimedWorktreeAction`, `spec.Command`
  - `cleanupFields`, `cleanupRow`, `joins`, `slugOf`, `explicitSetFingerprintVersion`
- `Import edges`: `internal/status` imports `internal/worktree` today. `internal/worktree` imports `internal/spec` in three landing files, so `internal/spec` gains no import. `cmd/bench` already imports both packages, and ticket 5 exports the slug derivation from `internal/spec` for the dispatcher.
- `Source-row clauses and occurrences`: the table above.
- `Promised field labels`: `class=`, `holder=`, `superseded candidate:`, `unique refs:`, `unique refs: unavailable`, `<n> landed ref`, `<n> subsumed ref`, `<n> unique ref`, `retained: content main lacks`.
- `Changed-function callers`: `UnclaimedAssignmentBranchRefs` has two callers, `appendGit` and the landing prune test. `applyUnclaimedAssignmentSet` has two callers, `cleanCommandWith` and the outcomes test at its line 347. All four are in the fences.
- `Copy survival`: none.
- `Rendered-shape readers`: the status action string has two test readers, `internal/status/status_producible_test.go` and `internal/systemtest/status_route_converge_test.go`, and both join ticket 1. The re-plan prefix has two more, `internal/worktree/clean_set_wiring_test.go` and `internal/worktree/clean_set_outcomes_test.go`, and each joins the ticket that changes it.

Changed fixtures: `TestCleanUnclaimedApplyCurrent` and `TestPlanUnclaimedShiftResidueBranch` plant a branch with a unique commit and expect its removal.
Ticket 1 rewrites them to plant a landed branch for the removal case and to keep the unique case retained.
The unique-commit fixture of `internal/worktree/clean_landed_apply_test.go` reds under the class rule, and ticket 1 repairs it.

`TestStatusRouteExecutesUnclaimedCleanup` expects the routed command to delete the branch, and it reds at ticket 1's checkpoint once the unique branch survives.
Ticket 1 rewrites it to expect survival, and ticket 2 adds a second case for a landed branch through the printed apply command.
The `retained: content main lacks` pin lives only in `clean_unclaimed_test.go`, so ticket 4's suffix change touches one test file.
`TestGitSignalIgnoresLandedUnclaimedAssignmentBranch` expects no git row for a landed unclaimed branch, and ticket 2 changes that posture to `1 landed ref` with the plan route.

Reader sweep: the string `bench worktree clean --discard-branch --unclaimed --apply-current` occurs in `internal/status/status.go`, `internal/status/status_producible_test.go`, and the roadmap row FT199.
The roadmap row stays as history.
No workflow file, `.mjs` script, or public doc names the string.
`bench anchors CONTEXT.md` lists eight anchors, none on the worktree terms.

Claim verification: each current-code claim in the problem and solution traces to a cited symbol, and the probe record is in decision ticket 1.
`DATA_HANDLING.md` records no value this spec changes, because a discarded ref carries a branch path and a date and no objective text.
