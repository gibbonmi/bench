# FT311 recoverable reset

Status: staged

Roadmap: FT311

Decision source: specs/ft311-landing-completion/decisions/ft311-coordinator-work.md

Verification log: 2 iteration(s) to accept — opus/high reviewed through the native agent surface. Iteration one returned 23 findings with 9 blocking, and the author folded 22 and recorded 1 as a reviewer-visible exception. Iteration two accepted the nine folds and returned 8 residuals with 1 blocking, which the author folded after the round. On 2026-09-11 the reviewer decided the ignored-path collision policy, and the repair added story 53 and rows RR70 through RR72.

## Problem

A delegate applies a migration script halfway and stops, and the worktree holds a mix of the old tree and the new tree.
A shift leaves an assignment worktree on its shift branch, and the exec, path, and clean verbs refuse the branch mismatch.
A merge publishes its tip and then fails the checkout catch-up, and the exit-3 record names a raw Git command.
In each case the coordinator runs raw Git by hand, and the guard refuses the destructive forms.
The reviewer then runs the checkout by hand, and the uncommitted work is preserved by memory or not at all.

## Solution

`bench worktree reset --to <commit> <target>` plans a return to one named recovery checkpoint, and `--apply <fingerprint>` executes that plan.
The checkpoint is an explicit commit in the assignment history, and the verb never infers it from the default branch.
Before the checkout moves, the apply preserves every staged, unstaged, and untracked change under a reset ref, and it verifies that envelope.
The record prints the preserved ref and the exact restore command, `bench worktree reset --restore <ref> <target>`.
The restore is the same plan-and-apply verb, and it returns HEAD, the index, and the working tree to the preserved state byte-exact.

This is the fourth of five FT311 specs.
It serves this kit and every repository that links it.
It changes neither the shift loop, nor the cleanup verdict, nor the landing.
The reset verb repairs the shift-branch state, and the shift end stays as it is.

## User stories

### Plan the reset

Line: gpt-5.6-terra / medium.
The plan composes existing readers behind an exact grammar under a covering test seam, which the scorecard routes at medium.

1. As a coordinator, I want `bench worktree reset --to <commit> <target>` to print a plan and change nothing, so that I read the effect before I apply it.
2. As a coordinator, I want the plan to name every fact the apply acts on, so that I read them before the apply.
3. As a coordinator, I want the plan to list every affected staged, unstaged, untracked, and renamed path, so that I know what the preservation carries.
4. As a coordinator, I want the plan to print the exact apply command with its fingerprint, so that the apply costs one paste.
5. As a coordinator, I want a checkout that already matches the checkpoint reported as `action=none`, so that an idle reset writes nothing.
6. As a reviewer, I want the checkpoint operand required, so that the verb never infers it from the default branch.
7. As a reviewer, I want a checkpoint outside the start-to-tip range refused with the range named, so that the range alone decides.
8. As a reviewer, I want a shift-branch commit outside that range refused as a checkpoint, so that the branch keeps its own history.
9. As a coordinator, I want a spelling that is not a commit refused by its own fault, so that a typo never reads as prose.
10. As a coordinator, I want a target that names no active assignment refused, so that the primary checkout and a retired assignment stay untouched.
11. As a coordinator, I want a conflicted index refused with the cleanup route named, so that conflict stages are never flattened.
12. As a coordinator, I want a dirty nested repository, an embedded repository, and an unknown nested state refused, so that the envelope carries no gitlink.
13. As a coordinator, I want a checkout with a live lease held by another process refused, so that two sessions never move one checkout.
14. As a coordinator, I want the fingerprint to follow the status, the content, the head, the ref, and the checkpoint, so that a change refuses.
15. As a coordinator, I want an ignored file to leave the fingerprint unchanged, so that a build output does not stale the plan.
16. As a coordinator, I want a control byte in a changed path rendered sanitized, so that the plan never splits its own line.
17. As a coordinator, I want a checkpoint spelling with a control byte refused before any lookup, so that it addresses nothing.
18. As a coordinator, I want a shift-branch checkout and a drifted lock planned rather than refused, so that the plan reaches the states it repairs.

### Preserve and apply

Line: gpt-5.6-terra / high.
The apply runs destructive Git under a preservation transaction, which the scorecard routes at high.

19. As a coordinator, I want `--apply <fingerprint>` to preserve every affected layer under a reset ref before the checkout moves, so that no work is lost.
20. As a coordinator, I want the preservation verified against the ref, the manifest, and the root's parents, so that a broken envelope stops the reset.
21. As a coordinator, I want the apply to leave a clean checkout on the assignment branch at the checkpoint, so that work restarts there.
22. As a coordinator, I want untracked files removed and ignored files kept, so that a half-applied migration leaves no residue and the build output survives.
23. As a coordinator, I want the record to print the preserved ref and the exact restore command, so that I can go back.
24. As a coordinator, I want a rewind or an off-branch checkout to keep the branch tip reachable, so that no commit is lost.
25. As a coordinator, I want a stale fingerprint refused with nothing written, so that a changed checkout never resets.
26. As a coordinator, I want a pure ref or lock repair to write no envelope, so that a no-loss repair leaves no ref behind.
27. As a coordinator, I want a move fault after preservation to exit 3 with the restore command, so that the envelope is never orphaned silently.
28. As a coordinator, I want a post-move verification failure to exit 3, so that a reset that did not land never reads as complete.
29. As a maintainer, I want the apply to take the cleanup apply's per-target lock, so that a reset and a cleanup never interleave.
30. As a reviewer, I want the assignment record left unchanged by a reset, so that the ledger's active-state rule survives.
31. As a coordinator, I want a second apply with a consumed fingerprint refused as stale, so that a reset applies once.

### Repair the known broken states

Line: gpt-5.6-terra / high.
These outcomes ride the apply ticket, so they take its line.

32. As a coordinator, I want a worktree left on its shift branch returned to the assignment branch, so that the lifecycle verbs accept it again.
33. As a coordinator, I want the shift branch and its commits retained, so that the shift's work stays reachable.
34. As a coordinator, I want a detached checkout re-attached to the assignment branch, so that the identity bundle passes again.
35. As a coordinator, I want a registration lock whose reason drifted restored to the exact Bench lock, so that the identity bundle passes again.
36. As a coordinator, I want an unreconciled merge checkout brought to the published tip, so that the merge's exit-3 repair has a Bench form.
37. As a reviewer, I want the shift end unchanged, so that this fold adds no behavior to the shift loop.

### Restore the preserved state

Line: gpt-5.6-terra / high.
The restore writes the index and the working tree from an envelope, which is the same destructive class as the apply.

38. As a coordinator, I want the restore to return HEAD, the index, and the working tree byte-exact, so that a regret costs one apply.
39. As a coordinator, I want the restore to preserve the current state first, so that a restore is itself recoverable.
40. As a coordinator, I want a ref outside this assignment's reset namespace refused, so that another assignment's envelope never lands here.
41. As a coordinator, I want an envelope that does not verify refused, so that a partial envelope never restores.
42. As a reviewer, I want the restore to return the assignment branch to the recorded tip, so that the branch never moves onto foreign history.
43. As a coordinator, I want the restored layers verified by a production recapture, so that the record reports a restore only when the trees match.
44. As a coordinator, I want a restore with a stale fingerprint refused, so that a changed checkout never restores over new work.
45. As a coordinator, I want a restore fault after preservation to exit 3 with a restore command, so that the interruption is recoverable.

### Keep the envelope durable

Line: gpt-5.6-terra / medium.
The namespace and the sweep rule compose existing owners under existing tests.

46. As a maintainer, I want the reset namespace outside the session-start sweep, so that an envelope survives the next session.
47. As a maintainer, I want the reconcile to delete a reset ref whose assignment record is gone, so that envelopes do not leak forever.
48. As a maintainer, I want an unreadable ledger to delete no reset ref, so that a broken ledger never costs an envelope.
49. As a maintainer, I want the cleanup's recovery capture unchanged after the move to one owner, so that the cleanup keeps every layer it preserves.
50. As a maintainer, I want the reset apply recorded under its own seam, so that the seam record names it like every other mutating verb.

### Keep the guidance current

Line: gpt-5.6-terra / high.
Guidance prose compounds through every session, so the leverage override applies.

51. As a teammate, I want the reference, the changelog, and the glossary to state the verb, so that a cold session finds the plan-and-apply form.
52. As a reviewer, I want the trials, the residuals, and the landing selector kept out of this build, so that this spec has one outcome.

### Refuse the ignored collision

Line: gpt-5.6-terra / high.
The refusal guards the same destructive move as the apply, so it takes the apply's line.

53. As a coordinator, I want a move over an ignored path refused with the paths named, so that no build output is lost.

## Implementation decisions

### The grammar and the target

The grammar is `bench worktree reset (--to <commit> | --restore <ref>) <target> [--apply <fingerprint>]`.
Exactly one of the two modes is required, and neither or both is a usage error at exit 2 before any ledger read.
The usage grammar has no exactly-one-of construct, so the verb checks the pair by hand after the parse.
The verb joins the worktree leaf family beside the merge verb.
The apply opens one verb span named `worktree.reset`, and the plan opens none, as the plan-and-apply cleanup verb opens none.

The target resolves through the assignment selector alone, and not through the shared resolver.
The shared resolver runs the creation bundle, which refuses a checked-out ref other than the assignment branch and a drifted lock.
Those two are the states the verb exists to repair.
So the verb composes the selector, the state check, the tree check, and the owner-marker step.

The owner-marker step proves a registered worktree of this repository, inside the pool, with a marker that names this owner and this path.
The verb reads the registration's branch and lock as facts to repair, and it refuses neither.
The primary checkout, a retired assignment, and a pooled shift worktree that is no assignment all refuse as unassigned targets.

A `--to` or `--restore` value with a control byte refuses before any lookup, through the shared line-safety predicate.
The verb owns its own two refusal sentences, because the merge verb's guard hard-codes its own flag name.
The verb's ledger read comes after that check, which is the verb's own ordering promise.

### The checkpoint predicate

A checkpoint is a commit that the assignment start reaches and that reaches the assignment branch tip, where reach includes equality.
Both ancestry questions go through the gate authorization's ancestor query, which the merge verb already runs.
A commit that fails either question refuses `checkpoint is outside the assignment history`, and the refusal names the start and the tip as the wanted range.
A default-branch commit the merge verb folded in satisfies both legs and is an accepted checkpoint.
A spelling that Git cannot peel to a commit refuses `checkpoint is not a commit`.
An ambiguous abbreviated sha is such a spelling, because the peel fails.

### The plan

The plan reads the head, the checked-out ref, the branch tip, the tracked state, and the affected paths.
It also reads the lock reason, the lease state, the nested state, and the content identity.
The tracked state comes from the same status argv the explicit cleanup planner runs.
The affected paths come from its NUL-framed records through the shared strict parser.
A rename's path-only second record carries no status, and the table drops it, so a rename is one row.

The plan prints one `reset_plan{...}` line on stdout and then one `reset_paths[N]{path,status}` table when N is not zero.
The line carries `worktree`, `mode`, `action`, `checkpoint`, `head`, `ref`, `tip`, `tracked`, `lock`, `preserve`, `fingerprint`, and `next`.
The `action` cell reads `reset` or `none`, and a `none` plan carries `fingerprint=none` and no `next` cell.
A `reset` plan's `next` cell is the exact apply command with the id operand.

A plan is `none` when the checkout is clean, HEAD is on the assignment branch, the head equals the checkpoint, and the lock is exact.
The `preserve` cell reads `envelope` when the checkout is dirty, when the head differs from the checkpoint, or when the tip differs from the checkpoint.
It reads `none` otherwise, which is the pure ref or lock repair.
A shift-branch checkout and a detached checkout each plan with their observed ref, and a drifted lock plans `lock=repair`.

A conflicted index, a dirty submodule, an embedded repository, an unknown nested state, and a live lease of another process each refuse at exit 1.
The conflicted refusal names the explicit cleanup verb as its next command, because that verb's envelope keeps conflict stages.
The unknown nested state is the walk error the nested classifier returns, and the verb fails closed on it.
Every refusal renders through the worktree refusal record at exit 1, with a sanitized path table where the refusal reads paths.

An ignored path that the move would materialize refuses `ignored content would be overwritten` at exit 1, with every colliding path in the table.
The reset materializes the checkpoint tree, and the restore materializes the envelope's tip, base, and working layer.
An ignored path collides when one of those trees tracks the path itself or a directory above it.
The refusal writes nothing, because the reset envelope carries no ignored bytes and the overwrite would be unrecoverable.
The coordinator moves the ignored content aside and plans again.

### The fingerprint

The fingerprint binds the version tag, the common directory, the owner id, the assignment id, the mode, and the checkpoint.
It also binds the head, the checked-out ref, the tip, the lock reason, the raw status bytes, the content identity, and the restore ref.
The content identity is the one the explicit cleanup planner computes, so an edit inside an already-dirty file stales the plan.
The status bytes exclude ignored entries, so a build output written between the plan and the apply does not stale the plan.
The landing's checkout fingerprint includes ignored entries except the runtime-log and local-capture paths, because the landing must see them.
The digest composes through the shared fingerprint-parts helper, so the byte layout has one owner.

The fingerprint also binds the raw index entries from the staged listing.
So a staged blob that changes under equal status and working bytes stales the plan too.

### The reset envelope

A reset envelope is the recovery envelope shape the cleanup already writes: one root commit whose tree holds `manifest.json`, with one payload commit per layer.
The layer capture moves out of the cleanup's recovery function into one shared function, and the cleanup calls it with its existing behavior.
The reset calls the same function with one difference: the working layer is always recorded, even when it equals the head tree.
Each payload's parent is the head at capture, so the head stays reachable through the envelope.

The reset manifest carries one more field, `tip`, which the cleanup leaves empty.
The field records the assignment branch tip at capture.
When the tip is not the head, the root commit takes the tip as one more parent.
So an off-branch capture keeps the tip reachable too.
The shared manifest reader accepts the field, and the cleanup's own envelopes never set it.
The field is additive, and the schema token stays `bench-recovery/v1`, because only the reset writes the field and the reader ignores an absent one.

The envelope ref is `refs/bench/reset/<owner>/<assignment>/<n>`, and the ordinal comes from the same next-ref walk the recovery prefix uses.
The namespace is a sibling of the recovery namespace in the ledger package, with its own prefix function.
The reset writes nothing to the assignment record, because the ledger refuses an active assignment that carries recovery metadata.
The ref is written with a zero old value, so a concurrent writer of the same ordinal refuses rather than overwrites.

Verification is the reset's own composition of three checks, because the cleanup's verifier requires an assignment record that names the envelope.
First, the ref resolves to the recorded root.
Second, the manifest parses through the shared manifest reader, with a non-empty `tip`.
Third, every layer payload the manifest names is a parent of the root commit, so the payload set has the root as its independent source.
A ref that does not exist under the prefix fails the first check and refuses `reset envelope does not verify`.

The reset namespace is Bench-owned.
A verified envelope under it is trusted as the cleanup trusts its own refs, and that is the stated trust assumption.

### The apply

The apply takes the per-target cleanup registration lock, re-plans, and compares the fingerprint.
A mismatch refuses `reset plan is stale` with the fresh fingerprint as the wanted value and the fresh plan command as the next command.
The apply then writes and verifies the envelope when the plan says `preserve=envelope`, and a verification failure refuses before any move.
The move attaches HEAD to the assignment branch and resets the branch and the checkout to the checkpoint.
It removes untracked files and keeps ignored files, through a clean that names no ignored entry.
A lock repair unlocks the registration and locks it again with the exact Bench reason, because Git refuses a second lock on a locked worktree.

The move runs through one seam-set field, so a test injects a fault or a no-op move without a Git stand-in.
The restore's layer write runs through a second seam-set field, which ticket 4 adds beside the first.
After the move the apply reads the head, the ref, and the status again, and it requires the checkpoint, the branch, and a clean tree.
A move fault and a failed post-move read both exit 3 with a record that carries `preserved=<ref>` and a `next` cell.
The `next` cell names the restore of that envelope, and a pure ref or lock repair names the plan command instead.
The success record is `reset{worktree,mode,checkpoint,previous,ref,preserved,restore}`, and `restore` is the exact restore plan command or `none`.

### The restore mode

`--restore <ref>` names an envelope under this assignment's reset prefix, and any other ref refuses `restore ref is not this assignment's`.
The envelope must verify, and its recorded tip must be a commit the assignment start reaches.
A tip the start does not reach refuses `envelope tip is not reached by the assignment start`, with the start as the wanted value.

The restore returns the assignment branch to the recorded tip and HEAD to the base.
HEAD is attached when the base equals the tip, and it is detached at the base otherwise.
So an off-branch capture restores without moving the branch onto that history.
A restore that leaves HEAD detached prints a `next` cell, `bench worktree reset --to <tip> <id>`, because the lifecycle verbs refuse a detached checkout.

The plan carries one more cell, `envelope=<ref>`, and its checkpoint cell is the base.
The restore preserves the current state under a new envelope first, whenever the checkout is dirty or its head or tip differs from the envelope's.

After the move, the restore writes the working layer into the index and the working tree with one reset read-tree.
It then writes the staged layer into the index alone, or the base tree when the envelope has no staged layer.
The restore then recaptures both layers in production and requires them equal to the envelope's.
So the record reports a restore only on a byte-exact match.
A layer fault and a recapture mismatch both exit 3, and the `next` cell names the restore of the new envelope when one exists.
When the restore wrote no new envelope, the `next` cell names the same restore plan again.

### The namespace sweep

The session-start reconcile keeps its two lifecycle namespaces and adds no third one to that list.
It gains one rule.
A ref under the reset namespace whose owner and assignment match no ledger record is deleted against the object the listing read.
A reset ref of any recorded assignment survives every reconcile.
The rule reads the ledger once, and an unreadable ledger deletes no reset ref.
An assignment record leaves the ledger on the reconcile after its cleanup completes, so its envelopes leave on the reconcile after that.

### Guidance

The reference's command notes gain one paragraph on the reset verb.
It states the two modes, the plan and apply, the envelope, the restore, and the states the verb repairs.
The changelog records the verb under the unreleased heading.
The glossary gains `reset envelope`, because the cleanup's `recovery` envelope and the reset's envelope must not share one name.
The `recovery checkpoint` term stays as written, because it already names what `--to` restores.
No anchored sentence changes, and no canary fixture moves.

## Testing decisions

The primary seam is the package-internal reset entry point with its seam set replaced.
That seam is how the merge verb proves its exit-3 boundary, and it lets a fixture fault the move without a Git stand-in.
Every fixture is a real repository with a real assignment, created through the existing creation helper, because the verb's facts are Git facts.
The shift-branch state is reached through the shift branch prefix, a raw switch, and the shift's retention lock reason inside the fixture.
The unreconciled merge state is reached through the merge verb with its reconcile seam faulted, as the merge suite already does.

The reconcile rows attach at the resume-clean command over a seeded repository, as the reconcile suite already does.
The differential row for the moved capture is the existing layer-preservation test, unchanged.
Each of the three verb tickets splits its tests across two or three files, with the row-to-file assignment named in the ticket.
A new file over the budget reds the lane on its own commit.
Every new test runs in the ordinary Go phase of the project gate, and the full landing gate remains the code oracle.

Two Git flags have no observed run from the authoring session, because the guard denies them there: `git reset --hard` and `git clean -fd`.
The dry run `git clean -nd` is not denied, and it ran on a scratch repository on 2026-09-10 over four shapes.
The shapes were an untracked file inside an untracked directory, an ignored file inside that same directory, an ignored directory, and an untracked symlink.
The listing was `Would remove link.txt` and `Would remove untracked-dir/keep.txt`, and the two ignored entries and the ignored directory stayed.
So `-d` without `-x` removes the untracked file, the untracked symlink, and the untracked directory's untracked files, and it keeps every ignored entry.
Ticket 3's fixture is that nested shape, and its row is the observed run of the destructive form.

A clean submodule under `git clean -fd` and `git reset --hard` is documented behavior only.
Git does not descend into a nested repository without a second force, and a reset does not recurse into a submodule by default.
The `read-tree -u --reset` round-trip and the double-lock refusal were observed on a scratch repository on 2026-09-10.
The round-trip recaptured both layers equal to the captured trees.

### Seam diagram

    trigger: coordinator command from any checkout of the repository
        |
        v
    argv --> [ resetWith: grammar, target, checkpoint, plan ] --> reset_plan line, reset_paths table
                  ^ tests attach here: real assignment fixture, status read-back, fingerprint pairs
        |
        v
    plan + fingerprint --> [ apply: lock, re-plan, envelope, move, verify ] --> reset record
                  ^ tests attach here: seam-set move fault, envelope fault, lock attempt hook
        |
        v
    envelope ref --> [ restore: verify, move to tip and base, write layers, recapture ] --> reset record
                  ^ tests attach here: layer trees before and after, second envelope, seam corruption
        |
        v
    session start --> [ reconcileLifecycleDebris ] --> swept count
                  ^ tests attach here: seeded reset refs with and without a record

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| RR1 | 1, 2 | `bench worktree reset --to <tip> <id>` over a dirty checkout prints one `reset_plan{...}` line with `action=reset`, the checkpoint, the head, `ref=<branch>`, `tracked=dirty`, and `lock=ok` at exit 0, and the status bytes after the plan equal the bytes before. | New TestResetPlanReportsTheDirtyCheckoutAndWritesNothing through resetWith | A plan that resets as it plans changes the status bytes. |
| RR2 | 3 | The plan prints one `reset_paths[N]{path,status}` table with the staged path, the unstaged path, the untracked path, and one row for a renamed file, and no ignored path. | New TestResetPlanListsEveryAffectedPath through resetWith | A missing path hides work the preservation carries, and an unfiltered rename record renders an empty status. |
| RR3 | 16 | A changed path with an ESC byte renders through the control sanitizer, and the plan line stays one line. | New TestResetPlanSanitizesAHostilePath through resetWith | A raw byte splits the record for every reader. |
| RR4 | 4 | The plan's `next` cell is `bench worktree reset --to <checkpoint> <id> --apply <fingerprint>` with the id operand and the plan's own fingerprint. | New TestResetPlanNamesTheApplyCommand through resetWith | A path operand or a foreign fingerprint makes the paste refuse. |
| RR5 | 5 | A clean checkout on its branch at the checkpoint with the exact lock prints `action=none`, `fingerprint=none`, and no `next` cell. | New TestResetPlanReportsNothingToReset through resetWith | An idle plan that prints an apply command invites a needless envelope. |
| RR6 | 6 | `bench worktree reset <id>` with neither mode exits 2 with the grammar, and an unreadable ledger does not change that answer. | New TestResetRefusesNeitherMode through resetWith | A default checkpoint reads the default branch tip. |
| RR7 | 7 | A default-branch commit past the start that the branch never folded refuses `checkpoint is outside the assignment history` with `wanted=<start>..<tip>` at exit 1. | New TestResetRefusesACheckpointOutsideTheHistory through resetWith | A predicate that checks only the tip's ancestry accepts the merged-in commit. |
| RR8 | 8 | A shift-branch commit that descends from the start refuses with the same detail at exit 1. | New TestResetRefusesAShiftBranchCheckpoint through resetWith | A predicate that checks only the start's ancestry accepts it. |
| RR9 | 9 | A spelling that names no commit refuses `checkpoint is not a commit` at exit 1. | New TestResetRefusesANonCommitCheckpoint through resetWith | A shared headline hides which repair the operator owes. |
| RR10 | 17 | A `--to` value with an ESC byte refuses `--to contains control characters` at exit 1 over a ledger the fixture made unreadable. | New TestResetRefusesAControlByteCheckpoint through resetWith | A lookup that runs first reports the ledger fault instead. |
| RR11 | 10 | The primary checkout's path refuses at exit 1 with the checkout unchanged. | New TestResetRefusesThePrimaryCheckout through resetWith | A selector that accepts any path resets the primary. |
| RR12 | 11 | A conflicted index refuses `checkout is conflicted` at exit 1, and the `next` cell names the explicit cleanup verb for the target. | New TestResetRefusesAConflictedIndex through resetWith | A capture that flattens stages loses the theirs side. |
| RR13 | 12 | A dirty submodule refuses `nested repository is dirty` at exit 1. | New TestResetRefusesADirtyNestedRepository through resetWith | A capture over a dirty submodule records a stale gitlink. |
| RR14 | 12 | A clean embedded repository refuses `embedded repository is retained` at exit 1. | New TestResetRefusesAnEmbeddedRepository through resetWith | A capture records the embedded tree as a gitlink and loses its files. |
| RR15 | 13 | A lease line naming a live pid of another process refuses `assignment has a live lease` at exit 1. | New TestResetRefusesALiveLease through resetWith | Two sessions move one checkout. |
| RR16 | 14 | The fingerprint differs after a tracked edit to a clean file. | New TestResetFingerprintTracksTheCheckout through resetWith | A fingerprint over the operand alone lets a changed checkout apply. |
| RR17 | 15 | An ignored file written between two plans leaves the fingerprint equal. | New TestResetFingerprintIgnoresIgnoredFiles through resetWith | A fingerprint over the ignored inventory stales on every build. |
| RR18 | 19 | `--apply` over a dirty checkout writes `refs/bench/reset/<owner>/<id>/1`, and the envelope's working and staged layers equal the trees captured before the apply. | New TestResetApplyPreservesTheLayersAndMovesTheCheckout through resetWith | A move before the capture preserves the checkpoint's tree. |
| RR19 | 21 | After the apply, HEAD is attached to the assignment branch, the branch and the head equal the checkpoint, and the status is clean. | New TestResetApplyPreservesTheLayersAndMovesTheCheckout through resetWith | A move of the ref alone leaves the tree dirty. |
| RR20 | 22 | Over the nested shape, the apply removes the untracked directory's untracked file and the untracked symlink, and it keeps the ignored file inside that directory, the ignored directory, and their bytes. | New TestResetApplyPreservesTheLayersAndMovesTheCheckout through resetWith | A clean with the ignored flag removes the build output, and a clean without the directory flag leaves the untracked file. |
| RR21 | 23 | The record reads `preserved=<ref>` and `restore=bench worktree reset --restore <ref> <id>` at exit 0. | New TestResetRecordNamesTheRestoreCommand through resetWith | A record without the ref leaves the operator no way back. |
| RR22 | 20 | An envelope whose root the seam builds without one named payload as a parent refuses `reset envelope failed verification` at exit 1, and the head and the status are unchanged. | New TestResetApplyRefusesAnUnverifiedEnvelope through resetWith with the envelope seam | A verification after the move discovers the loss too late. |
| RR23 | 24 | A `--to` of the commit before the tip moves the branch to that commit, and the previous tip is the parent of the envelope's working payload. | New TestResetApplyKeepsTheRewoundTipReachable through resetWith | A capture that skips an unchanged working layer leaves the tip unreachable. |
| RR24 | 25 | A file edited between the plan and the apply refuses `reset plan is stale` at exit 1 with no new ref and no move. | New TestResetApplyRefusesAStalePlan through resetWith | An apply that trusts the operand resets a checkout it never planned. |
| RR25 | 26, 34 | A clean detached checkout at the checkpoint with the branch at the checkpoint applies with HEAD attached, `preserved=none`, `restore=none`, and no ref under the namespace. | New TestResetApplyReattachesWithoutAnEnvelope through resetWith | An unconditional envelope leaves a ref for every repair. |
| RR26 | 27 | A move fault after preservation exits 3, the ref survives, and the record carries `preserved=<ref>`. | New TestResetApplyExitsThreeOnAMoveFault through resetWith with the move seam | A refusal-shaped exit hides a written envelope. |
| RR27 | 28 | A move that returns without moving the checkout exits 3 with `preserved=<ref>`. | New TestResetApplyExitsThreeWhenTheMoveDidNotLand through resetWith with the move seam | A record that trusts the seam reports a reset that did not happen. |
| RR28 | 29 | The apply records one cleanup lock attempt for the target, and the plan records none. | New TestResetApplyTakesTheCleanupLock through resetWith with the lock attempt hook | An unlocked apply interleaves with an explicit cleanup. |
| RR29 | 30 | After the apply, the assignment record has state `active` and an empty recovery set, and the ledger validation accepts it. | New TestResetApplyLeavesTheAssignmentRecordUnchanged through resetWith | A record write trips the active-with-recovery rule at the next read. |
| RR30 | 32, 33 | A checkout on a shift-namespace branch locked with the shift's retention reason applies with HEAD attached to the assignment branch and the exact lock, and the shift branch keeps its tip. | New TestResetApplyRepairsAShiftBranchCheckout through resetWith | A branch deletion loses the shift's commits, and a repair of the branch alone leaves the bundle refusing. |
| RR31 | 35 | A registration locked with a foreign reason on the assignment branch applies with the exact Bench lock back, after which the target resolves through resolveAssignment. | New TestResetApplyRestoresTheExactLock through resetWith | A second lock without an unlock fails, and the bundle keeps refusing. |
| RR32 | 36 | A merge left unreconciled by a reconcile fault, then `--to <tip>` applied, leaves the checkout clean at the tip with the incoming file present. | New TestResetApplyReconcilesAnUnfinishedMerge through resetWith and mergeWith | A verb that refuses a dirty checkout leaves the exit-3 state raw. |
| RR33 | 37 | The shift loop's branch creation, teardown, and retention bytes are unchanged by the build. | Review-owned scope audit over the shift package, with the existing TestRetainAndLockLocksDropsLeaseAndPreservesDirt green | A shift-end repair reopens the decided seam. |
| RR34 | 38 | `--restore <ref>` after RR18 returns HEAD to the base on the branch, the recaptured working and staged trees equal the envelope's layers, the untracked file is back untracked, the staged file is back staged, and the record reads `mode=restore`. | New TestResetRestoreReturnsThePreservedStateByteExact through resetWith | A restore that writes the working layer alone loses the staged bytes. |
| RR35 | 39 | A restore over a dirty checkout writes a second envelope first and prints `preserved=<ref 2>`, and a restore of ref 2 returns that dirty state. | New TestResetRestorePreservesTheCurrentStateFirst through resetWith | A restore without its own envelope is the one unrecoverable step. |
| RR36 | 40 | A ref under another assignment's reset prefix refuses `restore ref is not this assignment's` at exit 1. | New TestResetRestoreRefusesAForeignRef through resetWith | An envelope of another assignment lands its tree here. |
| RR37 | 41 | A ref whose root has no manifest refuses `reset envelope does not verify` at exit 1. | New TestResetRestoreRefusesAMissingManifest through resetWith | A partial envelope restores a partial tree. |
| RR38 | 42 | An envelope whose recorded tip is not reached by the start refuses `envelope tip is not reached by the assignment start` with the start as `wanted` at exit 1. | New TestResetRestoreRefusesAForeignTip through resetWith | A restore moves the branch onto history the assignment never had. |
| RR39 | 44 | A restore apply after an edit refuses `reset plan is stale` at exit 1 with no move and no new ref. | New TestResetRestoreRefusesAStalePlan through resetWith | A restore mode that skips the recheck writes over new work. |
| RR40 | 45 | A layer-apply fault after preservation exits 3, and the record's `next` cell names the restore of the new envelope. | New TestResetRestoreExitsThreeOnALayerFault through resetWith with the layer-write seam | An exit that names the old envelope loses the state the restore just preserved. |
| RR41 | 46 | A session-start reconcile leaves a reset ref of an active assignment in place and reports it in no swept count. | New TestResumeReconcileKeepsAnActiveAssignmentsResetRefs through ResumeCleanCommand | A namespace added to the sweep list deletes the envelope at the next session. |
| RR42 | 47 | A reset ref whose owner and assignment match no ledger record is deleted by the reconcile against the object the listing read. | New TestResumeReconcileSweepsOrphanedResetRefs through ResumeCleanCommand | A sweep with no record rule leaks one ref per reset forever. |
| RR43 | 49 | The explicit cleanup of a dirty assignment preserves the same layers as before the extraction. | Existing `internal/worktree/lifecycle_acquire_test.go` (`TestRecoveryPreservesEveryGitVisibleLayerWithoutMovingBranchOrIndex`) | A changed capture drops a layer the cleanup promised. |
| RR44 | 50 | The `worktree.reset` seam is in the seam registry, an apply records one span whose subject is the assignment id, and a plan records none. | Existing `internal/worktree/otel_seams_test.go` (`TestWorktreeSeamsMatchTheRegistry`) and new TestResetApplyRecordsOneVerbSpan through resetWith | An unregistered seam reds the registry test, and a plan span diverges from the cleanup precedent. |
| RR45 | 1 | `bench worktree reset --help` answers its grammar at exit 0, `bench worktree --help` lists the reset grammar, and `bench help` carries the reset row. | Existing `cmd/bench/command_registry_test.go` (`TestKeptRoutesAnswerTheirOwnHelp`) and `cmd/bench/help_inventory_test.go` (`TestHelpInventoryIsComplete`), with the reset entries added | A leaf outside the family help is unreachable from the inventory. |
| RR46 | 51 | The reference, the changelog, and the glossary state the verb, its two modes, and the envelope term. | Review-owned Spec axis over the three guidance files | A missing paragraph leaves a cold session running the hand checkout. |
| RR47 | 52 | No trial, probe residual, landing selector, or shift-loop change enters the build. | Review-owned scope audit | An added grammar elsewhere belongs to another spec's fence. |
| RR48 | 41 | An envelope whose manifest names a payload that is not a parent of the root refuses `reset envelope does not verify` at exit 1. | New TestResetRestoreRefusesAnUnparentedPayload through resetWith | A verification that derives the payload set from the manifest alone is a tautology. |
| RR49 | 42 | An envelope captured on a shift branch restores the assignment branch to the recorded tip and HEAD detached at the base, with the layers back. | New TestResetRestoreReturnsAnOffBranchCapture through resetWith | A restore to the base on the branch moves the assignment branch onto shift history. |
| RR50 | 18 | A checkout on a shift-namespace branch plans `action=reset` with `ref=<shift branch>` at exit 0. | New TestResetPlanReachesAShiftBranchCheckout through resetWith | A target resolved through the shared resolver refuses the state the verb repairs. |
| RR51 | 18 | A registration locked with a foreign reason plans `lock=repair` at exit 0. | New TestResetPlanReachesADriftedLock through resetWith | A target resolved through the creation bundle refuses the drifted lock. |
| RR52 | 43 | A restore whose seam writes a wrong working layer exits 3 after the production recapture, and the record's `next` cell names a restore. | New TestResetRestoreExitsThreeOnARecaptureMismatch through resetWith with the layer-write seam | A restore that skips its own recapture reports a match it never checked. |
| RR53 | 27 | After a move fault with an envelope, the record's `next` cell is `bench worktree reset --restore <ref> <id>`. | New TestResetRecordNamesTheRestoreCommand through resetWith with the move seam | An exit-3 record without the command leaves the envelope orphaned. |
| RR54 | 6 | `bench worktree reset --to <tip> --restore <ref> <id>` exits 2 with the grammar. | New TestResetRefusesBothModes through resetWith | A verb that prefers one mode silently discards the other. |
| RR55 | 9 | An abbreviated sha that names two objects refuses `checkpoint is not a commit` at exit 1. | New TestResetRefusesAnAmbiguousCheckpoint through resetWith | A peel that takes the first match resets to the wrong commit. |
| RR56 | 10 | The label of a retired assignment refuses the assignment-state component at exit 1. | New TestResetRefusesARetiredAssignment through resetWith | A selector over every state resets a completed worktree. |
| RR57 | 12 | A dirty embedded repository refuses `embedded repository is retained` at exit 1. | New TestResetRefusesADirtyEmbeddedRepository through resetWith | A verdict keyed on dirtiness alone captures the gitlink. |
| RR58 | 14 | The fingerprint differs after a new commit on the branch. | New TestResetFingerprintTracksTheTip through resetWith | A fingerprint without the tip lets a rewind apply against a moved branch. |
| RR59 | 14 | The fingerprint differs after a detach of HEAD. | New TestResetFingerprintTracksTheRef through resetWith | A fingerprint without the ref lets a ref repair apply against a changed ref. |
| RR60 | 14 | The fingerprint differs between two `--to` values over one checkout. | New TestResetFingerprintTracksTheCheckpoint through resetWith | A fingerprint without the checkpoint applies one plan's fingerprint to another checkpoint. |
| RR61 | 31 | A second apply with a consumed fingerprint refuses `reset plan is stale` at exit 1 with no new ref. | New TestResetApplyRefusesAConsumedFingerprint through resetWith | A fingerprint that survives its own apply resets twice. |
| RR62 | 40 | A ref outside the reset namespace refuses `restore ref is not this assignment's` at exit 1. | New TestResetRestoreRefusesARefOutsideTheNamespace through resetWith | A namespace check alone accepts a green marker ref. |
| RR63 | 12 | A nested state the classifier cannot walk refuses `nested repository state is unknown` at exit 1. | New TestResetRefusesAnUnknownNestedState through resetWith | An unknown state read as clean captures a tree the walk never saw. |
| RR64 | 24 | A clean detached checkout at the checkpoint with the branch tip ahead plans `preserve=envelope`, and after the apply the previous tip is a parent of the envelope root. | New TestResetApplyKeepsTheTipOfADetachedCheckout through resetWith | A preserve rule keyed on the head alone rewinds the branch with no envelope. |
| RR65 | 24 | A clean shift-branch checkout at the checkpoint with the assignment tip ahead applies with the previous tip a parent of the envelope root. | New TestResetApplyKeepsTheTipOfAShiftBranchCheckout through resetWith | The same omission loses the tip through the shift-branch route. |
| RR66 | 48 | A reconcile over an unreadable ledger deletes no reset ref. | New TestResumeReconcileKeepsResetRefsOverAnUnreadableLedger through ResumeCleanCommand | A rule that reads absence as no record deletes every envelope. |
| RR67 | 14 | The fingerprint differs after an edit inside an already-dirty file. | New TestResetFingerprintTracksTheContent through resetWith | A fingerprint over the status bytes alone applies the plan to content it never showed. |
| RR68 | 42 | A restore that leaves HEAD detached prints `next=bench worktree reset --to <tip> <id>`, and a restore that leaves HEAD attached prints no `next` cell. | New TestResetRestoreNamesTheReattachWhenDetached through resetWith | A detached checkout with no named way back refuses every lifecycle verb. |
| RR69 | 41 | A ref under this assignment's prefix that does not exist refuses `reset envelope does not verify` at exit 1. | New TestResetRestoreRefusesAMissingRef through resetWith | A missing ref read as empty restores nothing and reports a restore. |
| RR70 | 53 | Over a checkpoint that tracks a file and a directory the tip removed and ignored, with ignored bytes at the file and inside the directory, `--to <checkpoint>` refuses `ignored content would be overwritten` at exit 1 with both paths in the table and a non-colliding ignored path absent, the ignored bytes unchanged, the head unchanged, and no ref written. | New TestResetRefusesAnIgnoredCollision through resetWith | A move that keeps ignored files by omission still overwrites the one the checkpoint tracks, and a check over the file alone misses the directory. |
| RR71 | 53 | After a tracked path in an envelope's working layer becomes ignored with new bytes, `--restore <ref>` refuses the same detail with that path in the table, the bytes unchanged, and the head and the branch unchanged. | New TestResetRestoreRefusesAnIgnoredCollision through resetWith | A restore that checks the checkpoint tree alone overwrites through the layer write. |
| RR72 | 14 | A staged blob changed under equal status and working bytes changes the fingerprint, and an apply of the old fingerprint refuses `reset plan is stale` in both modes with the new staged bytes intact. | New TestResetFingerprintTracksTheIndex and TestResetRestoreRefusesAStaleIndex through resetWith | A fingerprint over the status and the working diff alone accepts a plan the staged layer outgrew. |

### Edge inventory

The canonical walk covers empty input, boundaries, errors, repetition, process boundaries, and hostile environments.
The attached profile is the shell CLI checklist in projects/benchkit.md.
Every behavior serves this repository and every repository that links the kit.

| Class | Concrete disposition | Rows |
|---|---|---|
| Absent versus empty | A ref that does not exist under this prefix refuses, and a clean checkout at the checkpoint plans `none`. | RR5, RR69 |
| Already satisfied | A `none` plan carries no apply command and writes no envelope. | RR5, RR25 |
| Errors | A non-commit, a foreign checkpoint, a conflicted index, a nested repository, and a live lease each refuse before any write. | RR7, RR9, RR12, RR13, RR15 |
| Ordering | The envelope is written and verified before the move, and the lock is taken before the re-plan. | RR18, RR22, RR28 |
| Repetition | A consumed fingerprint refuses, and a restore of a restore's envelope returns the intermediate state. | RR61, RR35 |
| Interruption | A move fault and a silent move both exit 3 with the envelope named. | RR26, RR27, RR53 |
| Process boundary | The reconcile in a fresh process keeps a live assignment's refs and drops an orphaned one. | RR41, RR42 |
| Output shape | The plan line precedes the paths table, and the `none` plan has no `next` cell. | RR2, RR5 |
| Paths | A control byte in a changed path renders sanitized, and every command names the id operand. | RR3, RR4 |
| Identity | A control byte in the checkpoint, an ambiguous sha, and a foreign restore ref each refuse by their own fault. | RR10, RR55, RR36 |
| Ignored files | Ignored entries leave the fingerprint alone and survive the apply, and a collision with a materialized path refuses. | RR17, RR20, RR70, RR71 |
| Staged layer | A staged blob that changes under equal status and working bytes stales the plan. | RR72 |
| Ledger | The assignment record is never written, and an unreadable ledger deletes no envelope. | RR29, RR66 |
| Grammar | Neither mode and both modes are usage errors at exit 2. | RR6, RR54 |
| Off-branch capture | A detached or shift-branch capture keeps the tip reachable, and its restore returns the branch to that tip and names the re-attach. | RR64, RR65, RR49, RR68 |

**Won't handle:** A restore of conflict stages is refused at the plan, because the reset envelope carries no stage layers. `bench worktree clean` keeps its conflict-stage capture for a conflicted checkout.

**Won't handle:** The name of the previous HEAD ref is not restored. The shift branch keeps its ref and its commits, and `bench worktree clean --discard-branch --unclaimed` retires it.

**Won't handle:** Ignored files are neither preserved nor removed. `bench worktree clean --discard-ignored` owns the ignored inventory. A move that would overwrite an ignored path refuses before any write, and the coordinator moves that content aside by hand.

**Won't handle:** A pooled shift worktree that is no assignment is refused as unassigned. `bench worktree clean --discard-branch <path>` retires it, which is the route the light-path landing of 2026-09-09 used.

**Won't handle:** A reset of the primary checkout is refused as unassigned. The primary receives writes only through landings.

**Won't handle:** The cleanup's recovery namespace keeps its session-start sweep as it is. A change to that lifetime is a reviewer decision outside this spec.

**Won't handle:** A reset after a review froze the source tip makes the landing refuse the tip mismatch. The coordinator re-runs the review over the new range, which is the landing's existing rule.

**Won't handle:** A reset while an exec child runs in the target is guarded by the plan-time lease read alone. The same holds for a reset from a working directory inside the target. The coordinator runs the reset from the primary checkout with no child in the target, as every lifecycle verb expects.

No package-variable substitution and no deletion of a tree file is required by this design.
If implementation needs one, its ticket must attach the corresponding hostile case before it claims its row.

## Ownership fences

These entries are the union of the ticket Writes.
They authorize implementation only after the reviewer approves this spec and its ticket graph.

- `internal/intent/ledger/ledger.go`
- `internal/intent/ledger/validate.go`
- `internal/intent/ledger_aliases.go`
- `internal/worktree/clean.go`
- `internal/worktree/layers.go`
- `internal/worktree/reconcile.go`
- `internal/worktree/resume_reconcile_test.go`
- `internal/worktree/reset.go`
- `internal/worktree/reset_plan_test.go`
- `internal/worktree/reset_refusal_test.go`
- `internal/worktree/reset_fingerprint_test.go`
- `internal/worktree/reset_apply.go`
- `internal/worktree/reset_apply_test.go`
- `internal/worktree/reset_repair_test.go`
- `internal/worktree/reset_envelope.go`
- `internal/worktree/reset_restore.go`
- `internal/worktree/reset_restore_test.go`
- `internal/worktree/reset_restore_refusal_test.go`
- `internal/worktree/verb_span.go`
- `internal/worktree/joins.go`
- `internal/worktree/worktree.go`
- `internal/worktree/land.go`
- `internal/otelrecord/registry.go`
- `internal/usage/worktree.go`
- `cmd/bench/worktree_leaves.go`
- `cmd/bench/worktree_leaves_test.go`
- `cmd/bench/main.go`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `.bench/BENCH-reference.md`
- `CHANGELOG.md`
- `CONTEXT.md`
- `tests/canary/package-core-guard/unrouted-subcommand`
- `tests/canary/docs-currency-token-diet`
- `tests/canary/skills-index-command-adapters`
- `tests/canary/workflow-guidance-anchors`
- `reviews/ft311-recoverable-reset.md`
- `specs/ft311-recoverable-reset`

The last entry is the spec phase's own change, which the landed range carries.
The two conformance registries and the four canary entries are closure declarations from the preflight proposal, and the expected edit inside them is none.
The coordinator owns the conditional review pickup, so no ticket names it.

The build cannot edit this spec, its acceptance rows, or its ticket graph without the existing spec-change authority.
Four fenced files are at or over the line budget: `cmd/bench/main.go`, `cmd/bench/command_registry_test.go`, `internal/worktree/worktree.go`, and `internal/worktree/land.go`.
`internal/worktree/worktree.go` sheds the verb-span block into `verb_span.go` in the ticket that adds the reset seam.
`internal/worktree/land.go` sheds the seam set into `joins.go` in the ticket that adds the reset seam fields.
`cmd/bench/main.go` and `cmd/bench/command_registry_test.go` shed at least the lines they gain in the same commit.
Each verb ticket splits its tests across two or three new files, so no new file crosses the budget on its own commit.

## Out of scope

These are planning estimates of file edits, not measured implementation costs.
Each future capability gets its own specification and one whole-project landing gate.

| Capability | Estimated price | Derivation |
|---|---|---|
| Lower-tier trials | 6 edits, 1 gate runs | Task definitions, twelve matched runs' evidence projection, comparison, and reviewer decision record |
| Restore of conflict stages | 4 edits, 1 gate runs | A stage-layer writer, its index-info apply, a verification, and two tests |
| Reset envelope inventory in `bench worktree list` | 3 edits, 1 gate runs | One column, its reader, and one test |
| Recovery namespace lifetime | 3 edits, 1 gate runs | A sweep predicate, its reconcile row, and one test |
| Reset lock against a running exec child | 4 edits, 1 gate runs | A child lease in the exec verb, its probe in the plan, and two tests |

The remaining order is the trials spec.
FT120, FT254, FT258, and FT290 retain their neighboring subjects.
The existing merge lane, the cleanup verdict, the landing, and the shift loop remain unchanged.

## Further notes

### Review round and dogfood runs

The review round ran over the uncommitted draft at opus/high through the native agent surface, read-only, with a one-iteration cap.
Iteration one returned 23 findings across the six rubric questions, the map discipline, and the ticket quiz, with 9 blocking.
The author folded 22 into this draft and recorded the copy-survival finding as the reviewer-visible exception below.
Iteration two was a scoped pass over the nine blocking folds, declared under the 2026-09-09 precedent for a blocking first round.
It accepted all nine and returned 8 residuals with 1 blocking, a missing seam-set fence on ticket 4.
The author folded all 8 after the round.

No dogfood run applies, because the staged spec changes no executable.

### Decision provenance

The decision source is consumed in place under the landing-completion spec's folder, and this spec copies nothing.
The three structured map sources were reread on 2026-09-10.
The FT311 detail file's recoverable-reset group is the fold this spec settles, and its residual list is not.
The research report names no reset fact, and its Q1 census evidence is the origin of the verb request.
The kit profile's tier binding, hostile-input checklist, and lane table were reread, and none changes here.

The FT311 detail file leaves one decision to this spec: whether the shift end or the reset verb restores the assignment branch.
This spec decides that the reset verb restores it, and the shift end stays as it is.
The shift's retention path locks the dirty worktree in place on purpose, so the repair belongs to the explicit verb that plans and preserves.
The 2026-09-09 occurrence was a registered assignment with both a branch drift and a lock drift, and RR30 reproduces both together.

Every line names the Codex mid column, as the third FT311 spec did.
A Claude Code build resolves the same tier to `opus` at the same effort through the harness binding.

### Flagged additions and engineering choices

No new product promise beyond the reviewed map is intended.
The plan and record cell names, the `none` action, the fingerprint parts, the manifest `tip` field, and the ref ordinal are engineering choices.

Four calls are flagged for the reviewer.
First, the restore is a mode of the reset verb rather than a raw Git command.
The two-layer restore needs two plumbing commands, and the guard denies the destructive raw forms.
Second, the reset envelope always records the working layer and records the branch tip, so every capture keeps the head and the tip reachable.

Third, the reset namespace is `refs/bench/reset/` with a record-bound sweep rule.
The recovery namespace is emptied at every session start, and the ledger refuses recovery metadata on an active assignment.
Fourth, the fingerprint binds the cleanup planner's content identity, so an edit inside an already-dirty file stales the plan.

One contestable call is the branch move.
A rewind moves the assignment branch to the checkpoint, because the glossary defines the checkpoint as the commit a reset restores.
The alternative keeps the branch tip and stages the reversion, which never moves a ref but leaves HEAD off the checkpoint.
The reviewer may reopen this before ticket 3 dispatches.

One reviewer-visible exception is recorded.
The copy-survival rule wants a red-capable row that fails if a copy of the layer capture survives inside the cleanup.
No black-box row can tell an identical surviving copy from the shared function.
So the coordinator probes the shared function with an omission before ticket 1 commits.
The Standards axis grades the one-source rule over the diff.

### Source-sentence-to-row accounting

References below identify resolved tickets in the single decision source and the FT311 detail file's recoverable-reset group.
Each table entry paraphrases its assigned clause without creating another decision source.

| Source clause | Disposition |
|---|---|
| FT311: `bench worktree reset <label>` repairs a half-applied migration, and it plans before it applies | RR1, RR4, RR18 through RR21 |
| FT311: a worktree that `bench shift` left on its shift branch is one more state the reset plan repairs | RR30, RR50, RR65 |
| FT311: the spec decides whether the shift end or the reset verb restores the assignment branch | RR30, RR33, and the provenance note |
| Ticket 1: worktree reset provides planned, recoverable checkpoint restoration | RR1, RR18, RR34 |
| Ticket 3: the CLI derives revisions, file lists, and evidence handles from canonical sources | RR2, RR16 |
| Ticket 11: reset restores an explicitly named checkpoint from the assignment history | RR6, RR7, RR8 |
| Ticket 11: it preserves every affected staged, unstaged, and untracked change | RR18, RR20 |
| Ticket 11: it provides an exact restore command | RR21, RR34 |
| Ticket 11: apply requires the approved plan fingerprint | RR24, RR39 |
| Ticket 11: a changed checkout causes refusal | RR16, RR24, RR67 |
| Ticket 11: reset never infers the checkpoint from the latest main | RR6, RR7 |
| Ticket 12: a missing or conflicting revision identity blocks a dependent continuation instruction | RR9, RR10, RR55 |
| Ticket 16: a derived fence addition needs approval only when it grants new authority | The fence list above |
| Tickets 10 and 18: five specs, one source, fixed order | RR47 and Out of scope |

### Reader sweep and enforcement reads

The repository sweep used hidden-file coverage over Go, Markdown, fixtures, scripts, and workflow files.
No tree file names `bench worktree reset` today, and the decision map and the FT311 detail file name it in prose alone.
The recovery prefix's Go readers were enumerated with a repository-wide search for the prefix function, and every reader is listed below.

| Reader or owner | Read evidence and disposition |
|---|---|
| worktree mergeWith, mergeTarget, mergeTargetTip, and mergeReconcileNext | The verb shape, the target proof, the checkout proofs, and the exit-3 next command. The reset composes the same shape and refuses nothing the merge refuses by accident. |
| worktree fromRepresentable | The merge's control-byte guard, which hard-codes its flag name and runs after the ledger read. The reset composes the line-safety predicate and owns its two sentences. |
| worktree selectAssignment and resolveAssignment | The selector and the shared resolver. The reset composes the selector alone, because the resolver's creation bundle refuses the two repaired states. |
| worktree validateCreationBundle, identityBundleRefusal, ownerMarkerRefusal, and registrationRefusal | The bundle components and their order. The reset runs the marker step and reads the registration step's facts, and RR31 proves the bundle passes after the lock repair. |
| worktree recoverAssignmentWithFault, worktreeTree, realIndexTree, readIndexEntries, and commitTree | The layer capture and the envelope writer. Ticket 1 moves the capture into one shared function, and RR43 is the differential. |
| worktree readRecoveryManifest and nextRecoveryRef | The manifest reader and the ordinal walk. The reset composes both with its own prefix, and the reader accepts the `tip` field. |
| The cleanup's recoveryEnvelopeValid and verifyRecovery in worktree | The cleanup's validity reader and its record-naming verifier. Neither proves the reset's envelope, so the reset composes its own three checks. |
| The cleanup's recoveryMetadataMatches in worktree | The reader that compares recorded recovery refs against the refs present. A sibling namespace does not change it. |
| worktree predictedForeignRef | The reader that predicts an unowned checkout's recovery ref from digests. A sibling namespace does not change it. |
| worktree retireCheckout and releaseLeftover | The cleanup's preservation callers. Neither changes, because the capture keeps its signature behind the moved function. |
| worktree reconcile sweepLifecycleRefs and reconcileLifecycleDebris | The two swept namespaces and the ledger purge. RR41, RR42, and RR66 add the record-bound reset rule beside them. |
| ledger RecoveryRefNamespace, RecoveryRefPrefix, and ValidateAssignment | The namespace, the prefix, and the active-with-recovery refusal. The reset namespace joins as a sibling, and RR29 proves the refusal stays. |
| worktree lockCleanupRegistration and lockCleanupFile | The per-target flock and its attempt hook. RR28 proves the apply takes it, and its only other taker is the cleanup apply. |
| worktree relock and lockReason | The exact Bench lock and its restore. The lock repair unlocks first, because the scratch probe showed a second lock refuses. |
| worktree ProbeLease and LeaseFile | The lease probe. RR15 composes it. |
| worktree classifyNestedState | The nested-repository probe and its five kinds. RR13, RR14, RR57, and RR63 dispose of every kind but clean. |
| worktree fingerprintParts, explicitRetainFingerprint, and explicitContentIdentity | The digest helper, its precedent, and the content identity. The reset fingerprint composes the helper and the content identity. |
| landing CheckoutFingerprint and fingerprintStatus | The landing's checkout digest, which keeps ignored entries except two path classes. The reset states its own status filter and does not reuse this one. |
| worktree beginVerbSpan and otelVerbSeams, and otelrecord registry | The verb span and the seam registry. RR44 adds the reset seam to both, on the apply alone. |
| worktree landingSource | The landing's proof that the head is the frozen source tip. A reset after the freeze makes the landing refuse, which the edge inventory names. |
| shift createShiftBranch, teardown, and preserveAndRecover | The shift-branch creation and the retention path. RR33 keeps all three unchanged. |
| lifecyclepolicy DecideExplicit RegistrationShiftBranch | The one exception to the branch-mismatch retain. It stays the cleanup's rule, and the reset repairs the state it retains. |
| gitguard denyTable | The destructive forms the agent cannot run. The restore command is a Bench verb, so an agent can run it. |
| git ParsePorcelainZStrict | The strict status parser and its path-only rename record. RR2 drops that record. |
| cmd/bench worktreeLeaves, helpRow inventory, keptRoutes, and keptWorktreeGrammars | The leaf family and its pins. Ticket 2 adds the reset leaf, the help row, and the two pins. |
| usage worktreeCommands and WorktreeUsage | The family help. Ticket 2 adds the reset grammar, and ticket 4 widens it. |
| structure Growth | The lane's growth rule, which grades a new file from zero. Every verb ticket names two test files. |
| CONTEXT.md `recovery checkpoint` | The term already names the `--to` operand. Ticket 5 adds `reset envelope` beside it. |
| projects/benchkit.md lane table and tier binding | The lane argv and the model binding. Neither changes. |

### Pre-review proof checklist

- Cited symbols: every symbol in the reader sweep table above was read in the current tree on 2026-09-10, and the review round confirmed each one.
- Import edges: worktree already imports intent, git, landing, sanitize, usage, otelrecord, and the gate authorization package. The ledger package gains no import. The reconcile file already imports intent and git.
- Source-row clauses and occurrences: the accounting table enumerates the three fold sentences of the FT311 detail file and the six binding decision tickets. No source clause is rewritten by this spec.
- Promised field labels: the records name `reset_plan`, `reset_paths`, `reset`, `worktree`, `mode`, `action`, `checkpoint`, `head`, `ref`, `tip`, `tracked`, `lock`, `preserve`, `fingerprint`, and `next`. The cells and flags are `envelope`, `previous`, `preserved`, `restore`, `to`, `none`, `repair`, `ok`, `--to`, `--restore`, and `--apply`.
- Changed-function callers: recoverAssignmentWithFault is called by retireCheckout only. reconcileLifecycleDebris is called by the resume-clean command only. The seam set and the verb-span block move without a signature change, and their callers stay inside the worktree package. worktreeCommands is read by WorktreeUsage only.
- Copy survival: the layer capture gets one owner, and RR43 is the behavior differential. No black-box row can red an identical surviving copy, which the reviewer-visible exception above records.
- Git flags: `read-tree -u --reset`, `read-tree`, `worktree lock`, `worktree unlock`, and `clean -nd` were observed on scratch repositories on 2026-09-10. `reset --hard` and `clean -fd` are stated assumptions until ticket 3's nested fixture runs.

### Implementation ticket approval table

| Ticket | Blocked by | Delivered outcome |
|---|---|---|
| 1. Extract the layer capture and add the reset namespace | none | One shared capture with the `tip` field, the reset prefix, and the record-bound sweep rule |
| 2. Plan the reset | none | `bench worktree reset --to <commit> <target>` with its plan, its fingerprint, its refusals, and its help, and no `next` cell yet |
| 3. Apply the reset | 1.md, 2.md | `--apply <fingerprint>` with the envelope, the move, the lock repair, the exit-3 boundary, and the plan's `next` cell |
| 4. Restore the preserved state | 3.md | `--restore <ref>` with its verification, its layer write, its recapture proof, and the `restore` and `next` command cells |
| 5. Fold the guidance and the changelog | 4.md | The reference, the changelog, and the glossary state the verb |

Tickets 1 and 2 form the first frontier and run in parallel, one worktree each.
Their overlapping registry and canary entries are closure declarations that neither ticket edits, so the frontier needs no serialization.
Ticket 3 follows both, because it writes the envelope through ticket 1's function and the seam set beside ticket 2's plan.
Ticket 4 follows ticket 3, because it widens the grammar and the apply, and it carries the worktree package's invariant as its last writer.
Ticket 5 is the last writer of the shared guidance files.
The review pickup is created only for actionable findings.

### Approval surface

| Subject | Disposition requested |
|---|---|
| Stories and lines | Approve the six outcome groups and their bound model efforts. |
| Seams | Approve the package-internal reset entry point, the seam-set move field, the shared capture, and the reconcile rule. |
| Acceptance and edges | Approve RR1 through RR69 and the eight explicit exclusions. |
| Ownership fences | Approve the exact union above for implementation. |
| Scope and tickets | Approve the five-ticket graph within the fourth FT311 capability. |
| Branch move | Decide before ticket 3 dispatches whether a rewind moves the assignment branch, as recommended. |
| Copy-survival exception | Accept the coordinator's omission probe in place of a black-box row, or name the row you want. |
