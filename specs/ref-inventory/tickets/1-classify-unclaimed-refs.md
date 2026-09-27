# 1. Classify each unclaimed ref and narrow the bulk sweep

Blocked by: 6-repair-glossary-shift-namespace.md
Writes: internal/worktree/clean_classes.go (new), internal/worktree/clean_classes_test.go (new), internal/worktree/clean_unclaimed.go, internal/worktree/clean_unclaimed_test.go, internal/worktree/clean_landed_apply_test.go, internal/worktree/clean_set_apply_test.go, internal/worktree/clean_set_outcomes_test.go, internal/worktree/clean_set_wiring_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/status/status.go, tests/canary/docs-currency-token-diet/signal-vocabulary-drift, internal/status/status_producible_test.go, internal/systemtest/status_route_converge_test.go
Covers: RI1, RI2, RI3, RI4, RI5, RI6, RI7, RI8, RI9, RI10, RI11, RI12, RI13, RI14, RI15, RI17, RI18, RI19, RI20, RI21, RI22, RI23, RI25, RI28, RI55, RI59, RI60, RI61, RI62, RI63, RI66, RI72, RI73, RI81, RI83, RI84, RI85, RI86, RI87

## What to build

Chunk: RI-C1a.
This ticket edits a system-tagged test, so the author runs the system suite under `BENCH_KIT` through `bench test --check system`.

Add the class function beside the unclaimed planner.
It takes the recorded assignments, the protected set, the default branch, and the sorted unrecorded refs.
It returns one class and one holder per ref.
It calls `git.LandedInDefault` for each unrecorded ref, and the landed class takes precedence over the subsumed class.
A ref whose tip does not resolve to a commit produces a row with action `error`.
That row names the ref and the object type, and the set fingerprint stays empty.

A holder is an active or cleanup-pending recorded assignment branch, or a unique root.
No landed ref is a holder.
A recorded assignment branch that `LandedInDefault` proves landed by content only is not a holder, including active and cleanup-pending assignments.
A ref beneath an ancestry-landed ref is itself landed by ancestry, because ancestry is transitive.
A resolving unrecorded symbolic ref in either Bench namespace produces an error row naming the ref and `symref`, never a silent exclusion.
Git's ref enumeration skips a dangling symref, and the spec does not handle it.

A ref that is landed by content only is deleted by the bulk sweep but holds nothing.
So a ref under a content-landed ref that is not itself landed classifies as unique.

A checked-out foreign branch, a complete record's branch, and a subsumed ref are never holders.
A ref is subsumed when its tip equals, or is a strict ancestor of, a holder's tip.
Among unrecorded refs with an equal tip and no other holder, the lexically first full ref name is the unique root.
The named holder is the lexically first holder whose tip reaches the ref, with recorded branches first, then unique roots.

The plan row keeps the seven cleanup columns and the `tracked` cell `unclaimed`.
The detail cell starts with `class=<class>`, continues with ` holder=<ref>` for a subsumed row, and ends with the existing removal text for a removing row.
A unique row carries the action `retain`, and under this ticket its detail ends with `retained: content main lacks`.
That suffix is pinned only in `clean_unclaimed_test.go`, so ticket 4's suffix change touches one test file.

The fingerprint version becomes `bench-unclaimed-assignment-branches/v2` and binds each row's ref, tip, class, and holder.
The apply loop, under a fingerprint and under `--apply-current`, skips a row whose action does not remove.
The plan prints the help action `bench worktree clean --discard-branch --unclaimed --apply <fingerprint>` only when a row removes.
A plan with an error row exits 1 and prints no apply action, and `--apply-current` beside an error row refuses before any delete.

Change the status action for the git row to `bench worktree clean --discard-branch --unclaimed` and remove the `--apply-current` entry from the action table.
The routed `--apply-current` call keeps a unique branch after this ticket, so the system route test reds at this checkpoint otherwise.
Rewrite that test so that the routed command over a unique ref exits 0 and leaves the ref.
The git row details stay as they are until ticket 2.

Two existing tests plant a unique branch and expect its removal.
Rewrite them so that a landed branch is the removal case and the unique case is retained.
Repair the unique-commit fixture of the landed apply test the same way.
The error-row fixture writes the loose ref file under `.git/refs/heads/` directly, because `git update-ref` may refuse a non-commit.

## Acceptance

- [ ] A plan over the ref shapes of the edge inventory prints the class and holder of the spec for each row.
- [ ] Two unrecorded refs at the `main` tip both print `class=landed`.
- [ ] A chain A under B under C with C unique prints `holder=<C>` for A and B.
- [ ] An unrecorded ref at the tip of an active recorded branch prints `class=subsumed`, and one under a checked-out foreign branch prints `class=unique`.
- [ ] A loose ref that names a blob prints action `error` with the ref and `blob` in its detail.
- [ ] A plan with an error row exits 1 and prints no apply action.
- [ ] `--apply-current` over an error row and a landed row refuses before any delete, and the landed ref survives.
- [ ] An apply with the fingerprint, and an `--apply-current` run, each remove the landed and subsumed refs and keep the unique ref.
- [ ] After each apply, every deleted tip stays reachable from `main` or from a surviving ref, with a ref under an ancestry-landed ref in the fixture.
- [ ] A branch one commit under an ancestry-landed ref prints `class=landed` without a `holder=` field.
- [ ] A resolving unrecorded symref in either namespace prints action `error` with its name and `symref`, and the set has no fingerprint and no apply action.
- [ ] Both apply forms refuse a set with a symref before any delete, and the target and a removable sibling survive.
- [ ] An unrecorded ref at the first commit of a squash-folded recorded branch prints `class=unique` and `retain`, for an active and for a cleanup-pending record.
- [ ] A plan, then `main` fast-forwarded onto the unique ref, then the old fingerprint refuses as stale.
- [ ] A branch one commit under a complete record's branch prints `class=unique`, and one under a cleanup-pending recorded branch prints `class=subsumed holder=<that branch>`.
- [ ] A adds `x=1`, B changes `x` to `2` and is squash-folded into `main`: the plan prints B as `class=landed` and A as `class=unique`.
- [ ] A plan, one new commit on a subsumed ref, then the old fingerprint refuses as stale and prints `bench worktree clean --discard-branch --unclaimed`.
- [ ] A plan, then a recorded active branch at a descendant of the unique ref, then the old fingerprint refuses as stale with the tip unchanged.
- [ ] A plan with one landed row prints the apply help action with the set fingerprint, and a plan over unique rows only prints none.
- [ ] A plan over no Bench-namespace branch prints the empty table and exits 0.
- [ ] Two consecutive plans print identical rows and `for-each-ref refs/bench/` is unchanged between them.
- [ ] The git row action equals `bench worktree clean --discard-branch --unclaimed`, and the `--apply-current` action is absent from the action table.
- [ ] The system route test runs the routed command over a unique ref, exits 0, and the ref survives.
- [ ] `TestLandCommandPrunesSquashFoldedSiblingBranch` still passes.
