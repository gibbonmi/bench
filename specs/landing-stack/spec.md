# A landing stack grades several reviewed sources once and bisects a red

Status: staged

Roadmap: FT450

Decision source: reviewer-confirmed current conversation, 2026-10-26, decisions D1 to D6 under Further notes, with the closed decision of 2026-10-26 that FT450 owns several light-path sources landed in order under one gate.

Verification log: 5 iteration(s) to accept — the Opus/high spec round took 3 iterations after the reviewer raised its cap from 2 to 3, and the Opus/high ticket round took 2.

## Problem

Each landing holds the landing lease for its own whole gate. When several reviewed sources are ready at one time, they land one after another, and each one pays a whole gate.

A tickets-only source also pays a fold. Its completion proof needs the destination tip, so after each landing the next source folds `main` again.

The landing spans of 2026-10-19 to 2026-10-26 hold 81 landings that ran a gate. The whole gate took a median of 277 s, with a range of 139 s to 431 s. Those gates took 378 minutes in total.

Sixteen sources were ready before the landing in front of them started. All sixteen were light-path sources of drain dispatches, and each one folded `main` once.

A stack model saves about 62 minutes at most in those 8 days, which is about 54 minutes each week. The spread is wide. One drain dispatch on 2026-10-21 gives 42 minutes, and the median day gives about 1 minute.

Since FT449, a landing with no Go change runs the reader set in about 140 s. So the real saving for a stack of documentation-only sources is lower than this bound.

## Solution

`bench worktree land --stack <file>` lands two to sixteen reviewed sources with one command. The file lists the members in landing order, and each member carries the fields of one landing.

Every member folds the same destination tip first. The stack composes the members in order, and each member gets its own landing merge commit and its own tickets-only folder close. The whole gate grades the top tree once, and one ref update publishes every member.

The gate proves that the top tree holds only reviewed bytes, except a `capture/` path that the settle policy resolves. Each other path keeps the destination's bytes or the bytes of the one member that changes it.

On a red, the stack bisects its ordered prefixes with whole gates. It names the first member whose prefix is red on top of a green prefix. It publishes that green prefix and leaves each later member for a rerun.

Each member passes the build window on its own, and each preflight refusal stops the whole stack before composition. A `--source-tip <request>=<commit>` override pins the repaired tip of one member, so each refusal route reruns with an exact command. `bench worktree land --resume <published-commit> --stack <file>` finishes a published stack whose later step stopped.

## User stories

Line: opus / medium.
Implementation-line reason: LS-C2 changes the oracle's completion proof, LS-C3 owns the publication, and LS-C6 owns the bisect. The profile routes gate logic at mid effort. The rows are exact, and the seams exist: the completion proof, the authorization field of the landing owner, and the route walk. Each wrong answer fails a grade-sequence row or a tally row. Ticket 12 authors guidance, so its author raises the effort to high under the leverage override of `craft-line`.
Harder chunks: LS-C2, LS-C3, LS-C6.

### Compose and grade once (D2, D4)

1. As a drain coordinator, I want to land two or more verified light-path sources with one command, so that the batch pays one whole gate.
2. As a landing operator, I want the stack to compose its members in file order, so that the history follows the order that I chose.
3. As a reviewer, I want each member to keep its own landing merge commit and message, so that each reviewed source stays visible in the history.
4. As a reviewer, I want each tickets-only member to close its own folder in its own merge commit, so that each source keeps its own deliverable transition.
5. As a landing operator, I want a green stack to grade only its top tree, so that the stack costs one whole gate.
6. As a landing operator, I want a green stack to publish every member with one ref update, so that the destination never holds half of a green stack.
7. As a gate maintainer, I want a stack to run the whole gate and never the reader set, so that a green moves the green marker to the stack tip.
8. As a release operator, I want the stack to run only under the installed landing broker, so that no member supplies the executable that grades it.
9. As a landing operator, I want one landing lease to hold the destination for the whole stack, so that no other landing composes during the gates of the stack.
10. As a landing operator, I want the census count of each published member printed and its assignment released, so that each worktree retires as one landing retires it.

### One verdict covers only reviewed bytes (D3)

11. As a gate maintainer, I want the gate to prove that each path of a stack tree keeps the destination's bytes or the bytes of the one member that changes it, so that one verdict covers only reviewed changes.
12. As a gate maintainer, I want the proof to refuse a stack tree that changes the bytes of a member, so that the composition cannot rewrite a reviewed change.
13. As a gate maintainer, I want the proof to refuse a stack tree that changes a path that no member changes, so that no unreviewed destination delta lands.
14. As a gate maintainer, I want the proof to refuse a stack tree that keeps an entry under the closed tickets folder of a member, so that each folder close stays exact.
15. As a drain coordinator, I want the capture rule table to settle a `capture/` path that several members change, so that an inbox file never blocks a stack.
16. As a landing operator, I want a member whose tip lacks the destination tip refused before composition with the fold route, so that I fold each member once before the stack runs.
17. As a landing operator, I want two members that change one path outside `capture/` refused before composition, so that I learn of the overlap before any gate runs.
18. As a reviewer, I want a staged-spec member refused with a route to land it alone, so that a spec completion keeps its own completion evidence.
19. As a reviewer, I want a member whose assignment holds a delivery binding refused with a route to land it alone, so that a deliverable closure and the chunk-gap landings of a build keep their own landing.
20. As a gate maintainer, I want the stack preflight and the gate proof to read one owner rule, so that the early refusal and the oracle cannot drift apart.
21. As a reviewer, I want the commitment admission of each member to grade that member's own landing tree, so that a light-path member answers only for its own change.

### Red attribution and the bisect (D5)

22. As a landing operator, I want a red stack to bisect its ordered prefixes with whole gates, so that the red names one member.
23. As a landing operator, I want the stack to attribute the red to the first member whose prefix is red on top of a green prefix, so that the attribution has one exact meaning.
24. As a landing operator, I want the stack to publish the green prefix below the attributed member, so that the green work lands without a second run.
25. As a landing operator, I want each member after the attributed member left unlanded with a rerun route, so that I land it after the repair.
26. As a cost owner, I want the bisect to make at most one plus the ceiling of log2 of the member count whole gates, so that a red stack still costs less than serial landings.
27. As a landing operator, I want a red that the stack attributes to the first member to publish nothing, so that a red never lands.
28. As a landing operator, I want an infrastructure outcome or an inherited red to stop the stack with no further grade and no publication, so that a host fault or a red destination never reads as the red of a member.
29. As a landing operator, I want the refusal of the attributed member to route the repair to that member's worktree, so that the next action is exact.

### The build window and the preflight (D6)

30. As a build orchestrator, I want the build window to admit each member on its own, so that a stack never lands a member that a single landing would refuse.
31. As a landing operator, I want each preflight refusal to stop the whole stack before composition and to name its member, so that one run lists every repair.
32. As a reviewer, I want a landing override to admit only its own member, so that the stack consumes the override of each published member once.
33. As a landing operator, I want a malformed stack file refused as invalid usage, so that a bad file never reaches the lease.
34. As a landing operator, I want two members that name one worktree refused as invalid usage, so that one source never lands twice.
35. As a landing operator, I want the fields of each member to pass the single landing's grammar, so that a member accepts exactly what one landing accepts.
36. As a landing operator, I want a stack path that is not a regular file refused before any read, so that a FIFO never blocks the landing.

### Interruption and resume

37. As a landing operator, I want a stack whose post-publication step stopped to exit 3 with a stack resume command, so that one command finishes the published stack.
38. As a landing operator, I want the stack resume to prove each published member merge against the stack file before any step, so that a resume never finishes a stack that it did not publish.
39. As a landing operator, I want the stack resume to release each active published member and to change nothing on a second run, so that a repeated resume does no harm.
40. As a landing operator, I want the stack resume to leave each unlanded member active, so that a partial publication keeps the rest of the work.

### Records and adoption

41. As a teammate who walks in, I want an ADR, the glossary, the reference, and the changelog to state the landing stack, so that I read the decided state without its history.
42. As a drain coordinator, I want `/bench-drain` to land its verified light-path tickets as one stack, so that the drain pays one gate for its implement-now items.
43. As an agent, I want `bench help` to print the stack grammar and the stack resume grammar, so that I find the form without a trial call.
44. As a reviewer, I want a persistent queue, the background grades, the chunk-gap landings, and the merge plan disclosure left out, so that this spec ships only the stack.
45. As a cost owner, I want the serial gate time that a stack saves measured before the spec, so that the premise of the build stays checkable.

### Route recovery

46. As a landing operator, I want each refusal route that commits in a member to rerun with that member's repaired tip as an override, so that I follow the route to its end without an edit of the stack file.
47. As a landing operator, I want a route that leaves one member to rerun that member's single landing, so that the rerun never meets the floor of two members.
48. As a landing operator, I want a rerun to keep each override of the call that refused, so that two repairs in a row converge on one command.

## Implementation decisions

### Surface

- The verb gets two forms: `bench worktree land --stack <file> [--source-tip <request>=<commit>]...` and `bench worktree land --resume <published-commit> --stack <file> [--source-tip <request>=<commit>]...`. The usage owner declares both grammars beside the single forms, and `bench help` prints both rows.
- A `--source-tip <request>=<commit>` override replaces the `source_tip` of the member whose `request` matches. An override whose request names no member, a second override for one request, and a value with no `=` or an empty side exit 2.
- The land leaf dispatches on the `--stack` flag, as it dispatches on `--resume` today. The wrapper already sends every `worktree land` call to the installed broker, so the stack adds no executable hop.
- The stack runs inside the landing span at the verb boundary. The span names the assignment of the first member.

### The stack file

- The file is one strict JSON object that the exact-document decoder reads: `{"members":[{"path":"…","request":"…","base":"…","source_tip":"…","spec":"…","message":"…"}]}`. The `spec` field is optional, and every other member field is required.
- The file holds two to sixteen members. The bound keeps the worst bisect at five whole gates.
- The verb reads the file only when the path resolves to a regular file. A relative path resolves against the working directory of the caller.
- Each member maps to the argument vector of one landing and parses through the single landing grammar. An invalid member exits 2 and names its 1-based index.
- Two members that resolve to one canonical worktree path exit 2 and name both indexes.
- One stack reader owns the file read, the decoder, the member bound, the member grammar, and the override grammar. The verb maps each reader refusal to exit 2 before the lease.

### The preflight

- The destination proofs run once. Each member then runs the assignment proofs and the source proofs of one landing. One preflight prints every refusal of these groups, and each refusal names its member index and label.
- Three stack proofs run for each member after its source proofs:
  - A staged-spec member refuses with `stack-member-alone`.
  - A member whose assignment holds a delivery binding refuses with `stack-member-alone`. A binding of any deliverable kind counts, a started tickets-only deliverable included, because that member's landing closes the deliverable. An open build binds only a staged-spec deliverable, so an open-build test alone misses this member.
  - A member whose tip does not contain the destination tip refuses with `land-destination-moved` and its fold route.
- After the member proofs, the owner rule finds each path outside `capture/` that two members change. Each such pair refuses with `stack-members-overlap`, which names both labels and the paths.
- The build window then decides each member on its own, and the lease follows, as for one landing. The lease record names the assignment of the first member.

### The owner rule

- One pure policy owner maps each changed path to the members that change it, relative to the destination tree. It compares the mode and the object of each tree entry, so an attribute change counts, and a rename owns its old path and its new path.
- The owner rule leaves each `capture/` path out of the overlap set. The settle policy decides that path in the composition.
- The stack preflight and the gate proof both read this owner. It follows the precedent of the pure settle policy, and the gate already imports a landing subpackage.

### The composition and the grade

- The landing owner composes each member onto the previous member merge with the existing three-way composition. The first member composes onto the destination.
- After each composition, the owner applies the tickets-only folder close of the member with the published-tree transform. It then writes one merge commit with three facts:
  - The first parent is the previous merge, or the destination for the first member.
  - The second parent is the member tip.
  - The message is the member's message.
- These merge commits stay unreachable until publication.
- A composition conflict refuses the stack with the `composition-conflict` face for that member, before any grade.
- The commitment admission of each member grades that member's own landing tree: the destination composed with that member alone, after its transform.
- The stack grades a tree with the landing value and a stack context. The context names each member of the graded prefix by its source tip and its tickets folder.

### The gate

- The reader plan refuses the reader set for a stack context with the case `landing stack`. The case comes right after `no landing value`, so a kit stack prints `gate: whole gate (reader set refused: landing stack)`.
- When a stack context is present, the stack proof replaces the single completion proof. It runs for every stack context, a context of spec-less members included, so the early return of a run with no spec and no tickets folder does not skip it.
- The stack proof applies these rules to each path, in this order:
  - A path under the tickets folder of a member is absent. This rule wins over the next rules for the paths of that folder.
  - A path that no member changes keeps the destination's entry.
  - A path that one member changes keeps that member's entry.
  - A `capture/` path that two or more members change passes, because the settle policy owns it.
  - A path outside `capture/` that two or more members change refuses.
- The stack proof also applies the verified closure of each member that has a tickets folder, with that member's source tip, as the single proof does. A spec-less member has no deliverable, so it has no closure. The preflight refuses each member that holds a delivery binding, so each closure that a stack proves is empty.
- The proof runs on every grade, a reused green included, as the evidence checks of the complete checkpoint run on a hit.
- A proof fault refuses the stack with `land-handback`. The preflight refuses each known cause first, so this refusal is a defensive stop with no reachable input after the preflight.

### The bisect

- A green top tree publishes the top merge.
- The bisect continues only on the authorization kinds `green` and `candidate`. A `candidate` red means that the destination held a green before the run, so the destination counts as a green prefix.
- A `candidate` top tree starts the bisect over the prefixes 0 to n. The destination is prefix 0. The search keeps a green bound `lo` at 0 and a red bound `hi` at n.
- While `hi - lo > 1`, the stack grades the prefix `mid = (lo + hi) / 2`, rounded down. A green moves `lo` to `mid`, and a red moves `hi` to `mid`.
- The stack attributes the red to member `hi` and publishes the merge of prefix `lo`. When `lo` is 0, the stack publishes nothing.
- The stack makes at most one plus the ceiling of log2 n grades. Each grade is a whole gate.
- Any other kind on any grade stops the stack, and the stack publishes nothing. An `infrastructure` result refuses with `land-infrastructure`. An `inherited` red refuses with `land-handback`, because the destination held no green and no member can carry the repair.

### The publication

- Before publication, the owner rechecks the branch tip and the checkout fingerprint of each published member, and the destination.
- One commitment decision under the intent lock decides the build window and the admission of each published member again. In the same decision, it runs one compare-and-swap of the destination to the published merge. It also consumes the override of each published member and drops the lease.
- After publication, the stack runs these steps once for the published merge: the green marker advance, the commitment reconcile, the destination checkout reconcile, and the branch prune.
- It then reads the census count of each published member and releases that member's assignment, in file order. It runs the two landing effects once.

### The output and the exits

- Each member prints its `landing source{...}` line on stderr before composition.
- After the lease, stdout prints one `members[N]{member,label,source_tip,published_commit,outcome,census}` table in file order. The `outcome` cell reads `landed`, `red`, or `unlanded`. The `published_commit` cell reads `none` for a member that did not land.
- After a publication, stdout prints the effects table of one landing. A refusal after the lease prints its record next.
- After the lease, the last stdout line is `stack{destination_base=…,published_commit=…,landed=<k>,worktree=…}`. Its `published_commit` field reads `none` when nothing published. A refusal before the lease prints no members table and no `stack{...}` record.
- Exit 0 means that every member landed and every step finished. Exit 1 means that the stack refused at least one member and no step after publication stopped. Exit 2 means invalid usage. Exit 3 means that a step after publication stopped, and the record carries the stack resume.

### The refusal faces

- Three new land faces join the registry:
  - `stack-member-alone` has one instruction step that removes the member from the stack file and routes it to its own landing, then the rerun.
  - `stack-members-overlap` has one instruction step that removes the later member of the pair from the stack file, then the rerun. That member lands later on its own route.
  - `stack-red` routes the repair to the worktree of the attributed member, a commit there, and its review. One instruction step follows: when the members table shows a landed member, remove each landed member from the stack file, fold the destination into each remaining member, and set the new `source_tip` of each folded member in the stack file. The rerun comes last. The operator fills its `<repaired-tip>` slot after the instruction step, because the fold moves the tip of the attributed member a second time.
- Every other refusal keeps its single-landing face.
- The raising site composes the rerun of each stack-raised refusal by these rules, in this order:
  1. The rerun starts from the overrides of the refused call.
  2. It drops the override of each member that the route removes from the stack file.
  3. It drops the override of each member whose tip the route sets in the stack file.
  4. A face repairs the source when the single landing's repaired-tip rule names it, and `stack-red` repairs the source too. For such a face, the rerun adds or replaces `--source-tip <request>=<repaired-tip>` for the refused member, which is the attributed member of `stack-red`.
  5. For `source-tip-mismatch`, the rerun adds or replaces the override of the refused member with the observed tip.
  6. When the route leaves two or more members, the rerun is `bench worktree land --stack <file>` with the overrides that remain.
  7. When the route leaves one member, the rerun is the single `bench worktree land` of that member. Its `--source-tip` is the override that remains for that member, or else its file tip.
- A member whose branch tip differs from its `source_tip` or its override refuses with `source-tip-mismatch`, which names the observed tip.

### The resume

- The stack resume reads the members and the overrides of the call. It proves that the published commit resolves exactly and that the destination reaches it.
- The second parent of the published commit must be the tip of exactly one member, and that member's index is the published count k.
- The resume walks k first parents back from the published commit. Each merge must have the tip of the next member, in reverse file order, as its second parent. The last first parent of the walk is the destination base of the stack.
- The reviewed range of each published member must authenticate, as in one resume.
- The exit-3 record of a stack names the stack resume with each override of the call.
- After the proofs, it finishes the post-publication steps once for the published merge. It releases each published member that is still active, and it skips each released member that has a terminal receipt. It leaves each later member active.

### Persisted state

The stack adds no intent-ledger field and no review-record field. The lease record keeps its one assignment, which is the assignment of the first member.

## Implementation chunks

Each ticket that creates a seam that a later ticket consumes is its own chunk, so its review closes before its consumer starts. LS-C1 and LS-C2 give the oracle its owner rule and its stack proof. LS-C3 builds the stack owner and its commitment decision, with a green publication and no bisect. LS-C4 moves the help rows and adds the stack reader. LS-C5 ships the verb for a green stack. LS-C6 adds the bisect, LS-C7 the red attribution, and LS-C8 the preflight refusals. LS-C9 adds the resume and the help rows, and LS-C10 records the decided state.

Each ticket targets an author context of about 150k tokens, and the Further notes size each one.

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| LS-C1 / 01-share-the-stack-owner-rule.md | One owner rule names the members of each changed path | LS1, LS2, LS3, LS4, LS5, LS6 | TestStackOwnersNamesEachChangingMember, TestStackOverlapsLeaveOutCapture, TestStackOverlapsKeepAnOutsidePath, TestStackOwnersCountsAModeChange, TestStackOwnersCountsBothRenamePaths, TestStackOwnersKeepsExactPaths | no |
| LS-C2 / 02-prove-a-stack-in-the-gate.md | The gate proves a stack tree and grades it whole | LS7, LS8, LS9, LS10, LS11, LS12, LS13, LS14, LS15, LS16, LS17 | TestStackProofAcceptsTheExactStack, TestStackProofRefusesAChangedMemberPath, TestStackProofRefusesAnUnreviewedDelta, TestStackProofRefusesAKeptFolder, TestStackProofLeavesASharedCapturePathToTheSettlePolicy, TestStackProofRefusesAnOverlap, TestStackProofHonorsAMemberDeletion, TestStackProofAcceptsAOneMemberPrefix, TestStackProofRunsOnAReusedGreen, TestReaderPlanRunsWholeForALandingStack, TestReaderPlanNamesTheStackCaseFirst | yes |
| LS-C3 / 03-admit-a-stack-publication.md, 04-compose-and-publish-a-stack.md | The commitment decides every published member under one lock, and the landing owner composes, grades once, and publishes a green stack | LS18, LS19, LS20, LS21, LS22, LS23, LS24, LS25, LS26, LS27, LS28, LS36, LS38, LS39 | TestStackPublicationRefusesARefusedMember, TestStackPublicationConsumesEachOverride, TestStackPublicationKeepsAnUnpublishedOverride, TestStackPublicationDropsTheLease, TestBindingQueryNamesATicketsOnlyAssignment, TestBindingQueryLeavesALightPathAssignment, TestLandStackGradesOnce, TestLandStackChainsMemberMerges, TestLandStackPublishesWithOneUpdate, TestLandStackClosesEachFolderInItsMerge, TestLandStackAdmitsEachMemberTree, TestLandStackStopsOnInfrastructure, TestLandStackRefusesAMemberConflict, TestLandStackRechecksEachMember | yes |
| LS-C4 / 05-move-the-worktree-help-rows.md, 06-read-a-stack-file.md | The worktree help rows sit beside the worktree leaves, and one stack reader reads a strict stack file and its overrides | LS54, LS55, LS56, LS57, LS58, LS59, LS60, LS61, LS62, LS64, LS65, LS96, LS97, LS98 | TestHelpInventoryIsComplete, TestStackReaderRefusesAMalformedFile, TestStackReaderAcceptsSixteenMembers, TestStackReaderParsesEachMemberWithTheLandingGrammar, TestStackReaderRefusesASpecialFile, TestStackReaderRefusesABadOverride | no |
| LS-C5 / 07-land-a-green-stack.md | `bench worktree land --stack` lands a green stack, with each override pinned | LS40, LS41, LS42, LS43, LS44, LS45, LS50, LS51, LS52, LS53, LS63, LS86, LS87, LS94, LS95, LS101, LS104, LS105 | TestLandStackJourney, TestLandStackRefusesAnUnreadableStack, TestLandStackRefusesARepeatedMember, TestLandStackPrintsEveryMemberRefusal, TestLandStackRefusesAHeldLease, TestLandStackHoldsTheLease, TestLandStackRoutesEndWithTheStackRerun, TestLandStackRefusesAStaleMemberTip, TestLandStackRerunKeepsEachOverride, TestLandStackOverrideReplacesTheFileTip, TestLandCommandDoesNotTreatMessageValueAsStackFlag | no |
| LS-C6 / 08-bisect-a-red-stack.md | The landing owner bisects a red stack to one member | LS29, LS30, LS31, LS32, LS33, LS34, LS35, LS37, LS99 | TestLandStackGradesEachPrefixWithItsMembers, TestLandStackBisectsToTheRedMember, TestLandStackPublishesNothingForAFirstMemberRed, TestLandStackAttributesARedTransition, TestLandStackBoundsTheBisect, TestLandStackStopsOnInfrastructureInTheBisect, TestLandStackStopsOnAnInheritedRed | yes |
| LS-C7 / 09-attribute-a-red-stack.md | A red stack refuses its attributed member with `stack-red`, and the route walk proves each stack rerun | LS46, LS47, LS48, LS49, LS68, LS69, LS89, LS102, LS103, LS106 | TestLandStackRedJourney, TestLandStackInheritedRedHandsBack, TestLandingFacesFollowTheirRoutes | no |
| LS-C8 / 10-refuse-a-stack-member.md | Each stack preflight refusal stops the stack before composition and routes to an exact recovery | LS70, LS71, LS72, LS73, LS74, LS75, LS76, LS77, LS88, LS90, LS91, LS92, LS93 | TestLandStackRefusesABehindMember, TestLandStackRefusesAnOverlap, TestLandStackRefusesAStagedSpecMember, TestLandStackRefusesABoundMember, TestLandStackAdmitsEachMemberThroughTheWindow, TestLandingFacesFollowTheirRoutes | no |
| LS-C9 / 11-resume-a-stack.md | `bench worktree land --resume <published-commit> --stack <file>` finishes a published stack, and help lists both stack forms | LS66, LS67, LS78, LS79, LS80, LS81, LS82, LS100 | TestLandStackIncompleteNamesTheResume, TestLandStackResumeProvesEachMember, TestLandStackResumeFinishes, TestLandStackResumeKeepsUnlandedMembers, TestHelpInventoryIsComplete, TestWorktreeLandRouteForwardsTheStackForm | no |
| LS-C10 / 12-record-the-landing-stack.md | An ADR, the glossary, the reference, the changelog, and the drain state the landing stack | LS83, LS84, LS85 | review-owned rows, root conformance, prose | no |

### Completion plan

```bench-completion-plan
{"version":2,"execution":{"mode":"delegate","run_id":"ft450-landing-stack-full-20261026","orchestrator_session":"claude:ls-orchestrator-20261026","author_limit":3,"assignments":{"01-share-the-stack-owner-rule.md":[{"session":"claude:ls-t01-author-20261026","assignment":"1cc76f5a694166ff98662b56b893ac23","model":"opus","effort":"medium","source":"4b03797988ac14cf6e30d6a51a5afb36eaac6649","native_ref":"claude-agent:bench-writer-medium:ls-t01-author-20261026"},{"session":"claude:ls-t01-repair1-20261026","assignment":"1cc76f5a694166ff98662b56b893ac23","model":"opus","effort":"medium","source":"d866a582578bcdc1b68909f9488ecfe872b2adb9","native_ref":"claude-agent:bench-writer-medium:ls-t01-repair1-20261026","predecessor":"claude:ls-t01-author-20261026","trigger":"user-directed","stopped":"claude-agent:bench-writer-medium:ls-t01-author-20261026","preserved":"d172cd9967dcd2b4dff64de7732b5c959f388336"}],"02-prove-a-stack-in-the-gate.md":[{"session":"claude:ls-t02-author-20261026","assignment":"1cc76f5a694166ff98662b56b893ac23","model":"opus","effort":"medium","source":"554fa8647e60183b7fb8621a47a08b80561a339f","native_ref":"claude-agent:bench-writer-medium:ls-t02-author-20261026"},{"session":"claude:ls-t02-repair1-20261026","assignment":"1cc76f5a694166ff98662b56b893ac23","model":"opus","effort":"medium","source":"1d2ad57d872a71989b5af63f3af5da94c220e398","native_ref":"claude-agent:bench-writer-medium:ls-t02-repair1-20261026","predecessor":"claude:ls-t02-author-20261026","trigger":"user-directed","stopped":"claude-agent:bench-writer-medium:ls-t02-author-20261026","preserved":"a663e75135c53473df75b4da34c3328773d47e11"}],"03-admit-a-stack-publication.md":[{"session":"claude:ls-t03-author-20261026","assignment":"1cc76f5a694166ff98662b56b893ac23","model":"opus","effort":"medium","source":"c84e7c35c0003d04a1edb76c6bc997c392c03c23","native_ref":"claude-agent:bench-writer-medium:ls-t03-author-20261026"},{"session":"claude:ls-t03-repair1-20261026","assignment":"1cc76f5a694166ff98662b56b893ac23","model":"opus","effort":"medium","source":"548df27601399dd2ee9e28bcd94895fd417ccb07","native_ref":"claude-agent:bench-writer-medium:ls-t03-repair1-20261026","predecessor":"claude:ls-t03-author-20261026","trigger":"user-directed","stopped":"claude-agent:bench-writer-medium:ls-t03-author-20261026","preserved":"9648f3cce4ac8a6a9e39d75c61ca2bbcf310e66e"}],"04-compose-and-publish-a-stack.md":[{"session":"claude:ls-t04-author-20261026","assignment":"1cc76f5a694166ff98662b56b893ac23","model":"opus","effort":"medium","source":"9648f3cce4ac8a6a9e39d75c61ca2bbcf310e66e","native_ref":"claude-agent:bench-writer-medium:ls-t04-author-20261026"},{"session":"claude:ls-t04-repair1-20261026","assignment":"1cc76f5a694166ff98662b56b893ac23","model":"opus","effort":"medium","source":"c4a007e7b87478a342b96a8f031e050b305edf25","native_ref":"claude-agent:bench-writer-medium:ls-t04-repair1-20261026","predecessor":"claude:ls-t04-author-20261026","trigger":"user-directed","stopped":"claude-agent:bench-writer-medium:ls-t04-author-20261026","preserved":"d966b20288a8a026882d87a6acb38c9e85c18baf"}],"05-move-the-worktree-help-rows.md":[{"session":"claude:ls-t05-author-20261026","assignment":"1cc76f5a694166ff98662b56b893ac23","model":"opus","effort":"medium","source":"5d027145cfc46cda3286284b1fc867440846ddbb","native_ref":"claude-agent:bench-writer-medium:ls-t05-author-20261026"}],"06-read-a-stack-file.md":[{"session":"claude:ls-t06-author-20261026","assignment":"1cc76f5a694166ff98662b56b893ac23","model":"opus","effort":"medium","source":"fde72e0423d8e8ad9ebf841e3cb4cfac7308db77","native_ref":"claude-agent:bench-writer-medium:ls-t06-author-20261026"},{"session":"claude:ls-t06-repair1-20261026","assignment":"1cc76f5a694166ff98662b56b893ac23","model":"opus","effort":"medium","source":"2d6d6dd49900fddca5d0d2e49a446ef3a790bf66","native_ref":"claude-agent:bench-writer-medium:ls-t06-repair1-20261026","predecessor":"claude:ls-t06-author-20261026","trigger":"user-directed","stopped":"claude-agent:bench-writer-medium:ls-t06-author-20261026","preserved":"f2b81fc002c18b450b3c7a32b8a4505506ff21d4"}],"07-land-a-green-stack.md":[{"session":"claude:ls-t07-author-20261026","assignment":"1cc76f5a694166ff98662b56b893ac23","model":"opus","effort":"medium","source":"b0d6723440ce076a50867606934ad1f44ad55786","native_ref":"claude-agent:bench-writer-medium:ls-t07-author-20261026"},{"session":"claude:ls-t07-repair1-20261026","assignment":"1cc76f5a694166ff98662b56b893ac23","model":"opus","effort":"medium","source":"5501ecb815316974f6c9464ed201aa65b55477da","native_ref":"claude-agent:bench-writer-medium:ls-t07-repair1-20261026","predecessor":"claude:ls-t07-author-20261026","trigger":"user-directed","stopped":"claude-agent:bench-writer-medium:ls-t07-author-20261026","preserved":"5eb51e6cf66e0b7ed4d142f9888f7b0e13eabfd1"}],"08-bisect-a-red-stack.md":[{"session":"claude:ls-t08-author-20261026","assignment":"1cc76f5a694166ff98662b56b893ac23","model":"opus","effort":"medium","source":"b850660b2da2c037d71c2da74a16447e3f0da220","native_ref":"claude-agent:bench-writer-medium:ls-t08-author-20261026"},{"session":"claude:ls-t08-repair1-20261026","assignment":"1cc76f5a694166ff98662b56b893ac23","model":"opus","effort":"medium","source":"d7e895e8d431632bc0833fa0d1a3ddd22806fc35","native_ref":"claude-agent:bench-writer-medium:ls-t08-repair1-20261026","predecessor":"claude:ls-t08-author-20261026","trigger":"user-directed","stopped":"claude-agent:bench-writer-medium:ls-t08-author-20261026","preserved":"f9ea6a781c5bc02387879512e7df7f0efb2e40c1"}],"09-attribute-a-red-stack.md":[{"session":"claude:ls-t09-author-20261026","assignment":"1cc76f5a694166ff98662b56b893ac23","model":"opus","effort":"medium","source":"40cb573292fdb2ae8e62de46939bd471c8a1bbfa","native_ref":"claude-agent:bench-writer-medium:ls-t09-author-20261026"},{"session":"claude:ls-t09-repair1-20261026","assignment":"1cc76f5a694166ff98662b56b893ac23","model":"opus","effort":"medium","source":"8c4704bdef9295356cf47a06aa706ed9ab55f3fc","native_ref":"claude-agent:bench-writer-medium:ls-t09-repair1-20261026","predecessor":"claude:ls-t09-author-20261026","trigger":"user-directed","stopped":"claude-agent:bench-writer-medium:ls-t09-author-20261026","preserved":"a3eae91e35f6ca0d9f28c4d5e5bd1798d3b838d9"}],"10-refuse-a-stack-member.md":[{"session":"claude:ls-t10-author-20261026","assignment":"1cc76f5a694166ff98662b56b893ac23","model":"opus","effort":"medium","source":"e6e5f6c226826320c3e933fd59db97ea65a82793","native_ref":"claude-agent:bench-writer-medium:ls-t10-author-20261026"},{"session":"claude:ls-t10-repair1-20261026","assignment":"1cc76f5a694166ff98662b56b893ac23","model":"opus","effort":"medium","source":"f6b3039a76bf9bc65680dcf2bcb77872ea8ca394","native_ref":"claude-agent:bench-writer-medium:ls-t10-repair1-20261026","predecessor":"claude:ls-t10-author-20261026","trigger":"user-directed","stopped":"claude-agent:bench-writer-medium:ls-t10-author-20261026","preserved":"7aad4fd4297c90e316dc03b4500d4155c46d1f9a"}],"11-resume-a-stack.md":[{"session":"claude:ls-t11-author-20261026","assignment":"1cc76f5a694166ff98662b56b893ac23","model":"opus","effort":"medium","source":"3d8605b9b34c179d824f6ec68114a3e29bf9c3d9","native_ref":"claude-agent:bench-writer-medium:ls-t11-author-20261026"},{"session":"claude:ls-t11-repair1-20261026","assignment":"1cc76f5a694166ff98662b56b893ac23","model":"opus","effort":"medium","source":"7df4689348b9485aeaf3dd155ed5af601d941687","native_ref":"claude-agent:bench-writer-medium:ls-t11-repair1-20261026","predecessor":"claude:ls-t11-author-20261026","trigger":"user-directed","stopped":"claude-agent:bench-writer-medium:ls-t11-author-20261026","preserved":"3ed09a7196ce70e96f305cae6b80d6adc90965d8"}],"12-record-the-landing-stack.md":[{"session":"claude:ls-t12-author-20261027","assignment":"1cc76f5a694166ff98662b56b893ac23","model":"opus","effort":"high","source":"c7087287f38556230468153df3bfe60d73bfdf16","native_ref":"claude-agent:bench-writer-high:ls-t12-author-20261027"},{"session":"claude:ls-t12-repair1-20261027","assignment":"1cc76f5a694166ff98662b56b893ac23","model":"opus","effort":"high","source":"30ddb6312f681e4eb199ea6c929045509c042ab9","native_ref":"claude-agent:bench-writer-high:ls-t12-repair1-20261027","predecessor":"claude:ls-t12-author-20261027","trigger":"user-directed","stopped":"claude-agent:bench-writer-high:ls-t12-author-20261027","preserved":"c06c7ec58b64da0b96a316a9beb1224432d0ba76"}]}},"chunks":[{"id":"LS-C1","tickets":["01-share-the-stack-owner-rule.md"],"verification":[{"id":"stackpolicy","command":"bench test --package ./internal/landing/stackpolicy","ticket":"01-share-the-stack-owner-rule.md"},{"id":"root-conformance","command":"bench test --check root-conformance","ticket":"01-share-the-stack-owner-rule.md"}]},{"id":"LS-C2","tickets":["02-prove-a-stack-in-the-gate.md"],"verification":[{"id":"gate","command":"bench test --package ./internal/gate","ticket":"02-prove-a-stack-in-the-gate.md"},{"id":"root-conformance","command":"bench test --check root-conformance","ticket":"02-prove-a-stack-in-the-gate.md"}]},{"id":"LS-C3","tickets":["03-admit-a-stack-publication.md","04-compose-and-publish-a-stack.md"],"verification":[{"id":"commitment-repository","command":"bench test --package ./internal/commitment/repository","ticket":"03-admit-a-stack-publication.md"},{"id":"landing","command":"bench test --package ./internal/landing","ticket":"04-compose-and-publish-a-stack.md"},{"id":"commands","command":"bench test --package ./cmd/bench","ticket":"04-compose-and-publish-a-stack.md"},{"id":"root-conformance","command":"bench test --check root-conformance","ticket":"03-admit-a-stack-publication.md"},{"id":"root-conformance-04","command":"bench test --check root-conformance","ticket":"04-compose-and-publish-a-stack.md"},{"id":"commitmenttest","command":"bench test --package ./internal/commitment/commitmenttest","ticket":"03-admit-a-stack-publication.md"},{"id":"gate-04","command":"bench test --package ./internal/gate","ticket":"04-compose-and-publish-a-stack.md"}]},{"id":"LS-C4","tickets":["05-move-the-worktree-help-rows.md","06-read-a-stack-file.md"],"verification":[{"id":"commands","command":"bench test --package ./cmd/bench","ticket":"05-move-the-worktree-help-rows.md"},{"id":"worktree","command":"bench test --package ./internal/worktree","ticket":"06-read-a-stack-file.md"},{"id":"root-conformance","command":"bench test --check root-conformance","ticket":"05-move-the-worktree-help-rows.md"},{"id":"root-conformance-06","command":"bench test --check root-conformance","ticket":"06-read-a-stack-file.md"},{"id":"conformance-05","command":"bench test --package ./internal/conformance","ticket":"05-move-the-worktree-help-rows.md"}]},{"id":"LS-C5","tickets":["07-land-a-green-stack.md"],"verification":[{"id":"worktree","command":"bench test --package ./internal/worktree","ticket":"07-land-a-green-stack.md"},{"id":"usage","command":"bench test --package ./internal/usage","ticket":"07-land-a-green-stack.md"},{"id":"commands","command":"bench test --package ./cmd/bench","ticket":"07-land-a-green-stack.md"},{"id":"root-conformance","command":"bench test --check root-conformance","ticket":"07-land-a-green-stack.md"}]},{"id":"LS-C6","tickets":["08-bisect-a-red-stack.md"],"verification":[{"id":"landing","command":"bench test --package ./internal/landing","ticket":"08-bisect-a-red-stack.md"},{"id":"commands","command":"bench test --package ./cmd/bench","ticket":"08-bisect-a-red-stack.md"},{"id":"root-conformance","command":"bench test --check root-conformance","ticket":"08-bisect-a-red-stack.md"}]},{"id":"LS-C7","tickets":["09-attribute-a-red-stack.md"],"verification":[{"id":"worktree","command":"bench test --package ./internal/worktree","ticket":"09-attribute-a-red-stack.md"},{"id":"refusalroute","command":"bench test --package ./internal/refusalroute","ticket":"09-attribute-a-red-stack.md"},{"id":"commands","command":"bench test --package ./cmd/bench","ticket":"09-attribute-a-red-stack.md"},{"id":"root-conformance","command":"bench test --check root-conformance","ticket":"09-attribute-a-red-stack.md"}]},{"id":"LS-C8","tickets":["10-refuse-a-stack-member.md"],"verification":[{"id":"worktree","command":"bench test --package ./internal/worktree","ticket":"10-refuse-a-stack-member.md"},{"id":"refusalroute","command":"bench test --package ./internal/refusalroute","ticket":"10-refuse-a-stack-member.md"},{"id":"commands","command":"bench test --package ./cmd/bench","ticket":"10-refuse-a-stack-member.md"},{"id":"root-conformance","command":"bench test --check root-conformance","ticket":"10-refuse-a-stack-member.md"},{"id":"gate-10","command":"bench test --package ./internal/gate","ticket":"10-refuse-a-stack-member.md"}]},{"id":"LS-C9","tickets":["11-resume-a-stack.md"],"verification":[{"id":"worktree","command":"bench test --package ./internal/worktree","ticket":"11-resume-a-stack.md"},{"id":"commands","command":"bench test --package ./cmd/bench","ticket":"11-resume-a-stack.md"},{"id":"usage","command":"bench test --package ./internal/usage","ticket":"11-resume-a-stack.md"},{"id":"system","command":"bench test --check system","ticket":"11-resume-a-stack.md"},{"id":"root-conformance","command":"bench test --check root-conformance","ticket":"11-resume-a-stack.md"}]},{"id":"LS-C10","tickets":["12-record-the-landing-stack.md"],"verification":[{"id":"prose","command":"bench test --check prose","ticket":"12-record-the-landing-stack.md"},{"id":"conformance-package","command":"bench test --package ./internal/conformance","ticket":"12-record-the-landing-stack.md"},{"id":"root-conformance","command":"bench test --check root-conformance","ticket":"12-record-the-landing-stack.md"}]}],"final_verification":[{"id":"integration","command":"bench test --changed --base 34ac699156f965752a13fc3c2f1c703c842348fc"},{"id":"system","command":"bench test --check system"},{"id":"whole-project","command":"bench gate"}]}
```

## Testing decisions

- A good test drives the stack through the seam that owns each decision, and it counts grades. The authorization field of the landing owner records each graded tree, so a bisect row compares the exact sequence of graded prefixes. The worktree fixture gate appends to a tally file, so a journey row counts whole-gate runs.
- The prior art is `greenOwner` in `internal/landing/landing_helpers_test.go` for the authorization stub. For the proof, it is `gradeTicketsClose` in `internal/gate/commitment_completion_test.go`, which calls `completionTree` on a crafted tree. For the verb, it is `landingFixtureWithGateStep` and `landingGateTally` in `internal/worktree/land_fixtures_test.go`.
- A stack journey builds two or more assignment worktrees from one review base, so each member contains the destination tip.
- A red journey first records a green of the destination, so its top red reads `candidate` and the bisect runs. The inherited row (LS106) uses a destination with no green.
- Each bisect row states its arithmetic. The grade bound for n members is 1 + ⌈log2 n⌉.
- Each expected `members` row comes from the same `toon.Table` call that the producer makes, because a hex tip can read as a number to the TOON encoder.
- The route walk adds one producing fixture for each new face, because `routetest.Fixtures` reds a walked face that has no producing fixture.
- A stack-raised single face, such as `land-destination-moved` in LS92, takes its own cause key, because `routetest.Fixtures` refuses two producing fixtures for one face and cause.
- The walk's success check reads a `landed{` record today. A stack rerun ends with the `stack{...}` record, so the stack fixtures check that record, and a single rerun keeps the `landed{` check.
- The reader-plan rows reuse the `landingFixture` of `internal/gate/reader_plan_test.go` from a new file, because that file is at its 400-line limit.

### Seam diagram

    trigger: an operator or a drain coordinator runs bench worktree land --stack <file>
        │
        ▼
    stack file ──▶ [ member proofs + owner rule + window + lease ] ──▶ refusals, or the lease
        │                                    ◀ tests attach here: land_stack_test.go drives the real binary and counts the tally
        ▼
    members ──▶ [ landing owner: compose, grade, bisect, publish ] ──▶ the published prefix and the attributed member
        │                                    ◀ tests attach here: stack_test.go stubs the authorization and records each graded tree
        ▼
    prefix tree + stack context ──▶ [ gate: reader plan + stack proof ] ──▶ green, red, or a proof fault
                                             ◀ tests attach here: stack_completion_test.go grades crafted trees

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| LS1 | 20 | The owner rule maps a path that two members change to both member indexes, and a path that one member changes to that member alone | `internal/landing/stackpolicy/stackpolicy_test.go` (`TestStackOwnersNamesEachChangingMember`) | A rule that keeps only the last owner hides the overlap |
| LS2 | 15 | The overlap set leaves out a `capture/` path that two members change | `internal/landing/stackpolicy/stackpolicy_test.go` (`TestStackOverlapsLeaveOutCapture`) | A rule with no exemption refuses every stack that touches one inbox file |
| LS3 | 17 | The overlap set keeps a path outside `capture/` that two members change | `internal/landing/stackpolicy/stackpolicy_test.go` (`TestStackOverlapsKeepAnOutsidePath`) | A rule that exempts every path accepts an overlap |
| LS4 | 20 | A member that changes only the mode of one path owns that path | `internal/landing/stackpolicy/stackpolicy_test.go` (`TestStackOwnersCountsAModeChange`) | A rule that compares only object ids misses an attribute change |
| LS5 | 20 | A member that renames one path owns the old path and the new path | `internal/landing/stackpolicy/stackpolicy_test.go` (`TestStackOwnersCountsBothRenamePaths`) | A rule with rename detection hides the old path from an overlap |
| LS6 | 20 | A changed path that holds a space and a `*` maps to its exact spelling | `internal/landing/stackpolicy/stackpolicy_test.go` (`TestStackOwnersKeepsExactPaths`) | A rule that splits or globs a path names the wrong path |
| LS7 | 11 | The stack proof accepts the exact tree of a two-member stack | `internal/gate/stack_completion_test.go` (`TestStackProofAcceptsTheExactStack`) | A proof that refuses the exact tree refuses every stack |
| LS8 | 12 | A tree of a stack of two spec-less members whose bytes at a member's changed path differ from that member's bytes refuses with a message that names the path | `internal/gate/stack_completion_test.go` (`TestStackProofRefusesAChangedMemberPath`) | A proof that reads only the path set accepts rewritten bytes, and a proof placed beside the tickets-only branch never runs for spec-less members |
| LS9 | 13 | A tree of a stack of two spec-less members that changes a path that no member changes refuses with `stack composition changes <path>; review the destination delta` | `internal/gate/stack_completion_test.go` (`TestStackProofRefusesAnUnreviewedDelta`) | A proof that skips the unowned paths, or that skips a context with no tickets folder, lands an unreviewed change |
| LS10 | 14 | A stack tree that keeps an entry under the tickets folder of a member refuses with `completion keeps closed <path>` | `internal/gate/stack_completion_test.go` (`TestStackProofRefusesAKeptFolder`) | A proof that skips the folder close lands the open tickets |
| LS11 | 15 | A stack tree whose bytes at a `capture/` path that two members change match neither member passes the proof | `internal/gate/stack_completion_test.go` (`TestStackProofLeavesASharedCapturePathToTheSettlePolicy`) | A proof that demands one member's bytes refuses the union that the settle policy writes |
| LS12 | 11 | A stack tree with a path outside `capture/` that two members change refuses with a message that names both member indexes | `internal/gate/stack_completion_test.go` (`TestStackProofRefusesAnOverlap`) | A proof that takes the bytes of the last owner accepts a silent overwrite |
| LS13 | 11 | A stack tree passes only when it lacks a destination path that one member deletes | `internal/gate/stack_completion_test.go` (`TestStackProofHonorsAMemberDeletion`) | A proof that reads only the present entries keeps a deleted path |
| LS14 | 11, 22 | A one-member stack context over a member that holds the destination accepts the exact prefix tree | `internal/gate/stack_completion_test.go` (`TestStackProofAcceptsAOneMemberPrefix`) | A proof that demands two members refuses every bisect prefix of one member |
| LS15 | 11 | After a green of the top tree under no stack context, a grade of that tree under a stack context whose proof fails refuses | `internal/gate/stack_completion_test.go` (`TestStackProofRunsOnAReusedGreen`) | A reuse that skips the proof lands a tree that the proof refuses |
| LS16 | 7 | A landing with a stack context prints exactly `gate: whole gate (reader set refused: landing stack)` and runs the whole gate | `internal/gate/stack_completion_test.go` (`TestReaderPlanRunsWholeForALandingStack`) | A plan that ignores the stack context runs the reader set |
| LS17 | 7 | A stack landing over an uncovered path prints the `landing stack` case and not the `uncovered path` case | `internal/gate/stack_completion_test.go` (`TestReaderPlanNamesTheStackCaseFirst`) | A case placed after the change cases names the wrong cause |
| LS18 | 30 | The stack publication refuses, and runs no publish, when the build window refuses one member | `internal/commitment/repository/stack_publication_test.go` (`TestStackPublicationRefusesARefusedMember`) | A decision that reads only the first member lands through an open build |
| LS19 | 32 | A stack publication consumes the override of each published member | `internal/commitment/repository/stack_publication_test.go` (`TestStackPublicationConsumesEachOverride`) | A stack that consumes no override admits the member a second time |
| LS20 | 32 | A stack publication keeps the override of a member that it does not publish | `internal/commitment/repository/stack_publication_test.go` (`TestStackPublicationKeepsAnUnpublishedOverride`) | A stack that consumes every listed override strands the rerun of an unlanded member |
| LS21 | 9 | The stack publication drops the one landing lease in its own decision | `internal/commitment/repository/stack_publication_test.go` (`TestStackPublicationDropsTheLease`) | A lease that survives the publication refuses the next landing |
| LS22 | 19 | The binding query answers true for an assignment that holds the delivery binding of a started tickets-only deliverable | `internal/commitment/repository/stack_publication_test.go` (`TestBindingQueryNamesATicketsOnlyAssignment`) | A query over the open builds alone misses this assignment, because a tickets-only deliverable opens no build |
| LS23 | 19 | The binding query answers false for an unbound light-path assignment during an open build | `internal/commitment/repository/stack_publication_test.go` (`TestBindingQueryLeavesALightPathAssignment`) | A query that answers true for every assignment during a build refuses every stack in a chunk gap |
| LS24 | 1, 5 | A green three-member stack calls the authorization once, with the top tree | `internal/landing/stack_test.go` (`TestLandStackGradesOnce`) | A stack that grades each member pays one gate for each member |
| LS25 | 2, 3 | A green stack publishes a chain whose first parents run from the destination base through each member merge in file order, with each member tip as second parent and each member's message | `internal/landing/stack_test.go` (`TestLandStackChainsMemberMerges`) | A composition that reorders or squashes the members loses the reviewed history |
| LS26 | 6 | A green stack moves the destination ref once, from the destination base to the top merge | `internal/landing/stack_test.go` (`TestLandStackPublishesWithOneUpdate`) | A publication of one ref update for each member exposes half of a stack |
| LS27 | 4 | The tree of the first merge lacks the tickets folder of the first member, and the tree of the second merge lacks both folders | `internal/landing/stack_test.go` (`TestLandStackClosesEachFolderInItsMerge`) | A close applied only to the top tree leaves the folder in the first merge |
| LS28 | 21 | The admission receives the own landing tree of each member and never the top tree | `internal/landing/stack_test.go` (`TestLandStackAdmitsEachMemberTree`) | An admission over the top tree refuses a light-path member for the path of another member |
| LS29 | 11, 22 | Each grade carries a stack context that names exactly the members of its prefix | `internal/landing/stack_bisect_test.go` (`TestLandStackGradesEachPrefixWithItsMembers`) | A prefix grade with the full member list fails the proof for each member above the prefix |
| LS30 | 22, 26 | With four members and member 3 red, the stack grades the trees of prefixes 4, 2, and 3 in that order, which is 1 + ⌈log2 4⌉ = 3 grades | `internal/landing/stack_bisect_test.go` (`TestLandStackBisectsToTheRedMember`) | A linear search or a wrong midpoint grades another sequence |
| LS31 | 23 | With four members and member 3 red, the stack attributes the red to member 3 | `internal/landing/stack_bisect_test.go` (`TestLandStackBisectsToTheRedMember`) | An attribution to the last graded red prefix names member 4 or member 2 |
| LS32 | 24 | With four members and member 3 red, the stack publishes the merge of member 2 and no later merge | `internal/landing/stack_bisect_test.go` (`TestLandStackBisectsToTheRedMember`) | A stack that publishes nothing on a red discards the green prefix |
| LS33 | 27 | With four members and member 1 red, the stack grades prefixes 4, 2, and 1, attributes member 1, and moves no ref | `internal/landing/stack_bisect_test.go` (`TestLandStackPublishesNothingForAFirstMemberRed`) | A stack that publishes prefix 0 writes an empty merge or moves the ref |
| LS34 | 23 | When prefixes 1 to 4 grade green, red, green, and red, the stack grades prefixes 4, 2, and 1, attributes member 2, and publishes the merge of member 1 | `internal/landing/stack_bisect_test.go` (`TestLandStackAttributesARedTransition`) | An attribution that assumes a monotone red names member 4 |
| LS35 | 26 | With sixteen members and member 16 red, the stack grades prefixes 16, 8, 12, 14, and 15, which is 1 + ⌈log2 16⌉ = 5 grades | `internal/landing/stack_bisect_test.go` (`TestLandStackBoundsTheBisect`) | A search past the bound costs more gates than serial landings save |
| LS36 | 28 | An infrastructure result on the top tree stops the stack after one grade with no ref update | `internal/landing/stack_test.go` (`TestLandStackStopsOnInfrastructure`) | A stack that bisects an infrastructure result attributes a host fault to a member |
| LS37 | 28 | An infrastructure result on a bisect prefix stops the stack with no further grade and no ref update | `internal/landing/stack_bisect_test.go` (`TestLandStackStopsOnInfrastructureInTheBisect`) | A stack that publishes the last green prefix after a host fault lands on an unknown verdict |
| LS99 | 28 | An `inherited` red on the top tree stops the stack after one grade, moves no ref, and attributes no member | `internal/landing/stack_bisect_test.go` (`TestLandStackStopsOnAnInheritedRed`) | A stack that bisects an inherited red counts a red destination as green and attributes the red to member 1 |
| LS38 | 31 | A conflict of member 2 with member 1 refuses with a conflict error that names member 2, and the stack makes no grade | `internal/landing/stack_test.go` (`TestLandStackRefusesAMemberConflict`) | A stack that grades a conflicted tree runs a gate for nothing |
| LS39 | 6 | A member whose branch tip moves during the grade refuses the publication with no ref update | `internal/landing/stack_test.go` (`TestLandStackRechecksEachMember`) | A stack that rechecks only the top member publishes a moved source |
| LS40 | 1, 5 | Two tickets-only members from one destination land with one stack command, and the fixture gate tally reads one run | `internal/worktree/land_stack_test.go` (`TestLandStackJourney`) | A stack that grades each member leaves a tally of two |
| LS41 | 10 | After a green two-member stack, both member assignments are released | `internal/worktree/land_stack_test.go` (`TestLandStackJourney`) | A stack that releases only the lease owner leaves the second worktree active |
| LS42 | 4 | After a green two-member stack, the destination holds neither tickets folder | `internal/worktree/land_stack_test.go` (`TestLandStackJourney`) | A stack that closes only the first folder keeps the second |
| LS43 | 7 | After a green two-member stack, the green marker names the stack tip | `internal/worktree/land_stack_test.go` (`TestLandStackJourney`) | A stack that advances the marker to the first merge leaves the marker behind |
| LS44 | 10 | A green two-member stack prints one `members[2]{member,label,source_tip,published_commit,outcome,census}` table whose rows read `landed` and carry each member's census count | `internal/worktree/land_stack_test.go` (`TestLandStackJourney`) | A stack that prints one landing record hides the second member |
| LS45 | 10 | The last stdout line of a green two-member stack is the `stack{...}` record whose `published_commit` is the stack tip and whose `landed` is 2 | `internal/worktree/land_stack_test.go` (`TestLandStackJourney`) | A record of the first merge names the wrong tip |
| LS101 | 10 | A green two-member stack prints one `landing source{review_base=…,assignment_start=…}` line on stderr for each member, in file order | `internal/worktree/land_stack_test.go` (`TestLandStackJourney`) | A stack that prints the disclosure of the first member alone hides the reviewed range of the second |
| LS46 | 24 | With three members and member 2 red, the stack publishes member 1 and exits 1 | `internal/worktree/land_stack_red_test.go` (`TestLandStackRedJourney`) | A stack that stops at the first red publishes nothing |
| LS47 | 26 | With three members and member 2 red, the gate tally reads 3, which is 1 + ⌈log2 3⌉ | `internal/worktree/land_stack_red_test.go` (`TestLandStackRedJourney`) | A verb that grades the published prefix again before its publication leaves a tally of 4 |
| LS48 | 25 | With three members and member 2 red, member 3 stays active and its `outcome` reads `unlanded` | `internal/worktree/land_stack_red_test.go` (`TestLandStackRedJourney`) | A stack that releases every member drops the unlanded work |
| LS49 | 29 | The `next=` route of the `stack-red` refusal names the label of member 2 in its repair step | `internal/worktree/land_stack_red_test.go` (`TestLandStackRedJourney`) | A route that names the first member sends the repair to the wrong worktree |
| LS50 | 31 | A stack whose members 1 and 2 both fail a source proof prints two refusals, each with its member index | `internal/worktree/land_stack_test.go` (`TestLandStackPrintsEveryMemberRefusal`) | A preflight that stops at the first member hides the second repair |
| LS51 | 31 | After a refusal of two members, the intent ledger holds no landing lease | `internal/worktree/land_stack_test.go` (`TestLandStackPrintsEveryMemberRefusal`) | A preflight that takes the lease first strands the destination |
| LS52 | 9 | While another landing holds the lease, the stack refuses with `land-lease-held`, and the tally reads 0 | `internal/worktree/land_stack_test.go` (`TestLandStackRefusesAHeldLease`) | A stack that composes under a held lease races the publication of the owner |
| LS53 | 9 | A single landing that starts during the grade of a stack refuses with `land-lease-held` and names the assignment of member 1 | `internal/worktree/land_stack_test.go` (`TestLandStackHoldsTheLease`) | A stack that takes no lease lets a second landing publish under it |
| LS54 | 33 | The stack reader refuses an absent stack file | `internal/worktree/land_stack_input_test.go` (`TestStackReaderRefusesAMalformedFile`) | A reader that reads absence as an empty stack reaches the lease |
| LS55 | 33 | The stack reader refuses an empty stack file | `internal/worktree/land_stack_input_test.go` (`TestStackReaderRefusesAMalformedFile`) | A reader that decodes empty bytes as zero members reaches the lease |
| LS56 | 33 | The stack reader refuses a stack file with a duplicate key | `internal/worktree/land_stack_input_test.go` (`TestStackReaderRefusesAMalformedFile`) | A lenient decoder keeps the last value and lands the wrong member |
| LS57 | 33 | The stack reader refuses a stack file with an unknown field | `internal/worktree/land_stack_input_test.go` (`TestStackReaderRefusesAMalformedFile`) | A lenient decoder drops a misspelled field and lands without it |
| LS58 | 33 | The stack reader refuses a stack file with content after the object | `internal/worktree/land_stack_input_test.go` (`TestStackReaderRefusesAMalformedFile`) | A decoder that stops at the first value ignores a second stack |
| LS59 | 33 | The stack reader refuses a stack file with one member | `internal/worktree/land_stack_input_test.go` (`TestStackReaderRefusesAMalformedFile`) | A stack of one pays the stack proof for a plain landing |
| LS60 | 33 | The stack reader refuses a stack file with seventeen members | `internal/worktree/land_stack_input_test.go` (`TestStackReaderRefusesAMalformedFile`) | An unbounded stack holds the lease past the bisect bound |
| LS61 | 33 | The stack reader refuses a member that has no `message` field and names the member index | `internal/worktree/land_stack_input_test.go` (`TestStackReaderRefusesAMalformedFile`) | A reader that defaults an absent field writes a merge with no message |
| LS62 | 33 | The stack reader accepts a file with sixteen members and no final newline | `internal/worktree/land_stack_input_test.go` (`TestStackReaderAcceptsSixteenMembers`) | A bound of fifteen or a demand for a final newline refuses a valid stack |
| LS63 | 34 | Two members that name one worktree exit 2 and name both indexes, and the ledger holds no lease | `internal/worktree/land_stack_test.go` (`TestLandStackRefusesARepeatedMember`) | A stack that lands one source twice publishes two merges of one tip |
| LS64 | 35 | The stack reader refuses a member with an empty `message` with its index and the grammar line of the single landing | `internal/worktree/land_stack_input_test.go` (`TestStackReaderParsesEachMemberWithTheLandingGrammar`) | A reader that parses members with its own rules accepts a value that a landing refuses |
| LS65 | 36 | The stack reader refuses a FIFO at the stack path as not a regular file, and the call returns with no writer on the FIFO | `internal/worktree/land_stack_input_test.go` (`TestStackReaderRefusesASpecialFile`) | A reader that opens the path first blocks the landing |
| LS105 | 35 | A single landing whose `-m` value is `--stack` stays a single landing and never selects the stack form | `internal/worktree/land_flags_test.go` (`TestLandCommandDoesNotTreatMessageValueAsStackFlag`) | A dispatch that scans the arguments for the bare token reads a flag value as the stack flag |
| LS104 | 33 | A stack file that the stack reader refuses exits 2, and the ledger holds no lease | `internal/worktree/land_stack_test.go` (`TestLandStackRefusesAnUnreadableStack`) | A verb that maps a reader refusal to exit 1 or takes the lease first strands the destination |
| LS66 | 8 | The wrapper forwards `worktree land --stack <file>` to the installed broker with the arguments unchanged | `internal/systemtest/land_route_test.go` (`TestWorktreeLandRouteForwardsTheStackForm`) | A later wrapper change that matches only the single grammar routes the stack form to a repository executable |
| LS67 | 43 | `bench help` prints the stack grammar `bench worktree land --stack <file> [--source-tip <request>=<commit>]...` as one row | `cmd/bench/help_inventory_test.go` (`TestHelpInventoryIsComplete`) | A help that omits the form leaves an agent to find it by trial |
| LS68 | 29 | In a three-member stack with member 2 red, the printed route of `stack-red`, carried out step by step, ends in a stack landing that exits 0 | `internal/worktree/refusal_route_follow_test.go` (`TestLandingFacesFollowTheirRoutes`) | A route that does not remove the landed member or fold the remaining members repeats a refusal |
| LS89 | 47 | In a two-member stack with member 2 red, the rerun of the `stack-red` route is the single landing of member 2, and the carried-out route exits 0 | `internal/worktree/refusal_route_follow_test.go` (`TestLandingFacesFollowTheirRoutes`) | A stack rerun over one remaining member exits 2 at the floor of two members |
| LS69 | 31 | A stack file whose path holds a space reruns from the printed `next=` route | `internal/worktree/refusal_route_follow_test.go` (`TestLandingFacesFollowTheirRoutes`) | A route that prints the path unquoted splits it into two arguments |
| LS86 | 46 | In a three-member stack, the last step of a `source-not-clean` route for member 2 is `bench worktree land --stack <file> --source-tip <request>=<repaired-tip>` with the request of member 2 | `internal/worktree/land_stack_rerun_test.go` (`TestLandStackRoutesEndWithTheStackRerun`) | A rerun with no override for the repaired member meets `source-tip-mismatch` |
| LS87 | 35 | A member whose branch tip differs from its `source_tip` refuses with `source-tip-mismatch`, which names the observed tip | `internal/worktree/land_stack_rerun_test.go` (`TestLandStackRefusesAStaleMemberTip`) | A stack that lands the branch tip publishes a commit that nobody reviewed |
| LS94 | 48 | A call with an override for member 1 that refuses member 2 prints a rerun that keeps the member 1 override and adds the slot of member 2 | `internal/worktree/land_stack_rerun_test.go` (`TestLandStackRerunKeepsEachOverride`) | A rerun that drops the earlier override loops between the two members |
| LS95 | 46 | An override replaces the file `source_tip` of its member, and the stack publishes that member's merge with the override tip as second parent | `internal/worktree/land_stack_rerun_test.go` (`TestLandStackOverrideReplacesTheFileTip`) | A stack that reads the file tip alone refuses the repaired member or lands its old tip |
| LS96 | 33 | The stack reader refuses an override whose request names no member | `internal/worktree/land_stack_input_test.go` (`TestStackReaderRefusesABadOverride`) | A stack that ignores an unknown override lands a tip that the operator did not pin |
| LS97 | 33 | The stack reader refuses two overrides for one request | `internal/worktree/land_stack_input_test.go` (`TestStackReaderRefusesABadOverride`) | A stack that keeps the last override lands a tip that the operator did not choose |
| LS98 | 33 | The stack reader refuses an override with no `=`, an empty request, or an empty commit | `internal/worktree/land_stack_input_test.go` (`TestStackReaderRefusesABadOverride`) | A lenient split pins a member to an empty tip |
| LS70 | 16 | A member whose tip lacks the destination tip refuses with `land-destination-moved` and the label of that member, and the tally reads 0 | `internal/worktree/land_stack_refusal_test.go` (`TestLandStackRefusesABehindMember`) | A stack that leaves the missing fold to the gate pays a whole gate for a handback |
| LS71 | 17 | Two members that change one path outside `capture/` refuse with `stack-members-overlap`, which names both labels and the path, and the tally reads 0 | `internal/worktree/land_stack_refusal_test.go` (`TestLandStackRefusesAnOverlap`) | A stack that leaves the overlap to the gate proof pays a whole gate for a handback |
| LS72 | 18 | A member whose `spec` names a staged spec refuses with `stack-member-alone`, and the tally reads 0 | `internal/worktree/land_stack_refusal_test.go` (`TestLandStackRefusesAStagedSpecMember`) | A stack that admits a staged spec grades a completion with no completion record proof |
| LS73 | 19 | A member whose assignment holds the delivery binding of a started tickets-only deliverable refuses with `stack-member-alone`, and the tally reads 0 | `internal/worktree/land_stack_refusal_test.go` (`TestLandStackRefusesABoundMember`) | A stack that admits the member writes a deliverable closure that no member reviewed |
| LS88 | 19 | A spec-less member whose assignment an open build binds refuses with `stack-member-alone`, and the tally reads 0 | `internal/worktree/land_stack_refusal_test.go` (`TestLandStackRefusesABoundMember`) | A stack that tests only the `spec` field lands the chunk-gap work of a build outside FT424 |
| LS74 | 30 | With an open build and no chunk gap, a light-path member refuses with `land-build-window` and its index, and the ledger holds no lease | `internal/worktree/land_stack_refusal_test.go` (`TestLandStackAdmitsEachMemberThroughTheWindow`) | A window decision for the first member alone lands the second through the build |
| LS75 | 32 | With an open build and an override for member 1 alone, the stack refuses member 2 with `land-build-window` | `internal/worktree/land_stack_refusal_test.go` (`TestLandStackAdmitsEachMemberThroughTheWindow`) | An override applied to the whole stack admits a member that the reviewer did not name |
| LS76 | 18, 19 | In a three-member stack, the printed route of `stack-member-alone`, carried out step by step, ends in a stack landing that exits 0 | `internal/worktree/refusal_route_follow_test.go` (`TestLandingFacesFollowTheirRoutes`) | A route that does not take the member out of the stack repeats the refusal |
| LS90 | 47 | In a two-member stack, the rerun of the `stack-member-alone` route is the single landing of the other member, and the carried-out route exits 0 | `internal/worktree/refusal_route_follow_test.go` (`TestLandingFacesFollowTheirRoutes`) | A stack rerun over one remaining member exits 2 at the floor of two members |
| LS77 | 17 | In a three-member stack, the printed route of `stack-members-overlap`, carried out step by step, ends in a stack landing that exits 0 | `internal/worktree/refusal_route_follow_test.go` (`TestLandingFacesFollowTheirRoutes`) | A route that keeps both members repeats the refusal |
| LS91 | 47 | In a two-member stack, the rerun of the `stack-members-overlap` route is the single landing of the earlier member, and the carried-out route exits 0 | `internal/worktree/refusal_route_follow_test.go` (`TestLandingFacesFollowTheirRoutes`) | A stack rerun over one remaining member exits 2 at the floor of two members |
| LS92 | 46 | A `land-destination-moved` refusal that a stack raises, carried out step by step, ends in a stack landing that exits 0 | `internal/worktree/refusal_route_follow_test.go` (`TestLandingFacesFollowTheirRoutes`) | A rerun that keeps the stale `source_tip` meets `source-tip-mismatch` after the fold |
| LS93 | 46 | A `source-tip-mismatch` refusal that a stack raises, carried out step by step, ends in a stack landing that exits 0 | `internal/worktree/refusal_route_follow_test.go` (`TestLandingFacesFollowTheirRoutes`) | A rerun that keeps the file tip loops on the same refusal |
| LS102 | 46 | In a two-member stack with member 1 red, the rerun of the `stack-red` route adds `--source-tip <request>=<repaired-tip>` for member 1, and the carried-out route exits 0 | `internal/worktree/refusal_route_follow_test.go` (`TestLandingFacesFollowTheirRoutes`) | A `stack-red` rerun with no slot for the attributed member meets `source-tip-mismatch`, because nothing published and no step wrote its tip |
| LS103 | 48 | A three-member stack called with an override for member 1, with member 2 red, publishes member 1, and its carried-out `stack-red` route exits 0 | `internal/worktree/refusal_route_follow_test.go` (`TestLandingFacesFollowTheirRoutes`) | A rerun that keeps the override of the removed member 1 exits 2, and one that keeps an override over a folded tip meets `source-tip-mismatch` |
| LS106 | 28 | A stack whose top tree reads `inherited` refuses with `land-handback` after one grade, and the gate tally reads 1 | `internal/worktree/land_stack_red_test.go` (`TestLandStackInheritedRedHandsBack`) | A verb that maps an unattributed stop to `stack-red` sends the repair to a member that did not cause the red |
| LS78 | 37 | A stack called with one override whose release of member 2 fails exits 3, and its `next=` names `bench worktree land --resume <tip> --stack <file>` with that override | `internal/worktree/land_stack_resume_test.go` (`TestLandStackIncompleteNamesTheResume`) | A route that names a single resume reconciles from the wrong base, and one that drops the override fails the member proof |
| LS100 | 43 | `bench help` prints the stack resume grammar `bench worktree land --resume <published-commit> --stack <file> [--source-tip <request>=<commit>]...` as one row | `cmd/bench/help_inventory_test.go` (`TestHelpInventoryIsComplete`) | A help that omits the form leaves an agent to find it by trial |
| LS79 | 38 | The stack resume refuses a published commit whose second parent is not the tip of the listed member, before any step | `internal/worktree/land_stack_resume_test.go` (`TestLandStackResumeProvesEachMember`) | A resume that trusts the file finishes a stack that it did not publish |
| LS80 | 39 | The stack resume releases the active published member and exits 0 | `internal/worktree/land_stack_resume_test.go` (`TestLandStackResumeFinishes`) | A resume that skips the release leaves the member worktree active |
| LS81 | 39 | A second stack resume of the same stack exits 0 and changes no ref, no ledger record, and no checkout | `internal/worktree/land_stack_resume_test.go` (`TestLandStackResumeFinishes`) | A resume that repeats a release refuses or writes twice |
| LS82 | 40 | The resume of a stack that published members 1 and 2 of four leaves members 3 and 4 active | `internal/worktree/land_stack_resume_test.go` (`TestLandStackResumeKeepsUnlandedMembers`) | A resume that releases every listed member drops unlanded work |
| LS83 | 41 | An ADR states the members of a stack, its proof, its schedule, its attribution, its window rule, and the absent persistent queue | review-owned at the LS-C10 checkpoint | A missing ADR leaves the decision only in this spec |
| LS84 | 41 | The glossary, the reference, and the changelog state the landing stack, and the reference lists the `landing stack` reader-set case | review-owned at the LS-C10 checkpoint | A stale reference tells a teammate that each landing grades alone |
| LS85 | 42 | `/bench-drain` lands two or more verified implement-now light-path tickets as one landing stack | review-owned at the LS-C10 checkpoint | A drain that keeps serial landings pays one gate for each fix |

Not covered: story 44 — FT451, FT424, the parked merge plan idea, and a future queue spec own the excluded scope.
Not covered: story 45 — the Problem section and Further notes hold the measurement, and no build behavior proves it.

### Edge inventory

The walk covers each canonical edge class and the profile's hostile-input checklist:

- Empty, one, and boundary inputs: LS54, LS55, LS59, LS60, and LS62.
- Duplicate inputs: LS56 and LS63.
- Order: LS17 and LS25.
- Concurrency: LS39, LS52, and LS53.
- Interruption: LS78 to LS82.
- Oracle failure: LS36, LS37, and LS99.
- A route that leaves one member, and a route that commits in a member: LS86, LS89 to LS94.
- A malformed override: LS96 to LS98.
- A flag value read as a flag: LS105, which slicing added from the profile's hostile-input checklist.
- A red that does not grow with the prefix: LS34.
- An attribute-only change, a rename, and a deletion: LS4, LS5, and LS13.
- Paths with a space or a glob character: LS6 and LS69.
- A special file at the input path: LS65.
- A hand-edited file with no final newline: LS62.
- A two-sided change of a `capture/` path that the settle policy rewrites: LS2 and LS11.
- A numeric-looking hex cell in the members table: the Testing decisions derive each expected row through `toon.Table`.

**Won't handle** lines:

- A flaky red that the bisect attributes to a member — the single `bench worktree land` of that member grades it again, so a flake costs one rerun.
- A queue that lasts across sessions — a second session meets `land-lease-held` and reruns its own `bench worktree land` when the lease ends.
- An automatic fold of the members — each member folds with `bench worktree merge --from`, which the `land-destination-moved` route names.
- A cheaper reader-set grade of a Go-free stack — D4 runs the whole gate, and a single `bench worktree land` still selects the reader set.

## Out of scope

- A landing queue that lasts across sessions and that several sessions add to: about 15 edits, 4 gate runs. It needs an intent-ledger field and its own recovery.
- A staged-spec member in a stack, with a review-record proof for each member: about 8 edits, 2 gate runs.
- An automatic fold of each stack member before the stack runs: about 6 edits, 2 gate runs.
- The background grade of each lane-passed tip: FT451 owns it.
- The chunk-gap landings of an open build: FT424 owns them.
- The `bench worktree merge --plan` disclosure: its parked idea owns it.

## Further notes

### Decisions D1 to D6 (closed 2026-10-26)

- D1: The build goes ahead with the narrow scope of D2 to D6. The measured bound is real but small, and most of it comes from drain dispatches.
- D2: `bench worktree land --stack <file>` runs once over an ordered strict-JSON list of two or more members, each with the fields of one landing. One lease and one compare-and-swap publish the own merge commit of each member. The stack keeps no persistent queue and adds no ledger field.
- D3: A member is spec-less or tickets-only. A staged-spec member, or a member that an open build binds, lands alone. Each member tip contains the destination tip, and the deltas of two members overlap only on `capture/` paths that the settle policy resolves.
- D4: A stack always runs the whole gate and never the reader set. A green moves the green marker to the stack tip.
- D5: A red stack bisects its ordered prefixes with whole gates. It attributes the red to the first member whose prefix is red on top of a green prefix. It publishes that green prefix, refuses the attributed member, and leaves each later member unlanded with its rerun. The worst case is 1 + ⌈log2 n⌉ whole gates, and an infrastructure outcome stops the stack with no publication.
- D6: The build window keeps every rule. Each member passes the window admission on its own, and each preflight refusal stops the whole stack before composition and names each refused member.

### Measurement method

- The source is the landing spans (`worktree.land`) and their `gate.prospective` children in the span store of this repository, from 2026-10-19 to 2026-10-26, with the source commits in Git.
- A landing ran a gate when its span has a gate child longer than 60 s.
- A source was ready when its last non-merge commit on the source side predates the start of the landing in front of it.
- The model batches, at the start of each landing, every later landing whose source was ready by then. It prices a green batch at its longest gate. It prices a red batch at that gate times 2 + ⌈log2 n⌉, which counts the bisect and one regrade of the rest.
- Not verified: the Codex session records, the lane time of each fold, and the gate time of a Go-free stack after FT449.

### Validation plan

- After the build lands, count the stack landings of the next two drain dispatches, with their tallies and their wall time.
- Compare that wall time with the serial landings of the 2026-10-21 and 2026-10-23 dispatches.

### Research findings for each seam beyond the decision source

- The owner rule package: the settle policy is a pure policy subpackage of the landing (`internal/landing/composition.go:13`), and the gate already imports a landing subpackage (`internal/gate/complete_checkpoint.go:12`).
- The stack proof: the single completion proof compares each graded entry with the reviewed source (`internal/gate/completion.go:63`), and `checkpointEvaluation` reads its context for both the run path and the reuse path (`internal/gate/engine.go:65`, `internal/gate/engine.go:99`).
- The reader-plan case: the plan names its first refusing case in one line (`internal/gate/reader_plan.go:23`).
- The commitment decision: `PublishAdmitted` decides the window, the admission, the publish, the override, and the lease in one transaction (`internal/commitment/repository/publication.go:35`).
- The binding query: `openBuilds` counts only a binding whose deliverable reads as a staged spec (`internal/commitment/repository/window.go:246`), and `TestOpenBuildsDeriveFromActiveOwners` proves that a started tickets-only deliverable opens no build (`internal/commitment/repository/window_test.go:84`). So the query reads the bindings of the ledger, not the open builds.
- The bisect kinds: the landing routes a red by the kinds `inherited`, `candidate`, and `lane fail` (`internal/landing/merge.go:127`), and `inherited` means that the subject held no green before the run (`internal/gate/authorization/authorization.go:184`).
- The stack resume: the single resume proves a published merge by its two parents and its reach from the destination (`internal/worktree/land_resume.go:256`).
- The refusal faces: the route walk reds a walked face with no producing fixture (`internal/refusalroute/routetest/routetest.go:64`).
- The tip override: the single landing points its rerun at the repaired tip for each face that `repairsSource` names (`internal/worktree/land_refusal.go:269`), and the route walk fills that slot in `operatorFill` (`internal/worktree/refusal_route_follow_test.go:82`). A preface cannot carry the file edit, because `routetest.Steps` strips the preface before the walk (`internal/refusalroute/routetest/routetest.go:134`).

### Pre-review proof checklist

- `Cited symbols`: `landWith`, `landAttributed`, `landingAssignment`, `landingSourceProofs`, `landingSource`, `Owner.landReviewed`, `Owner.Compose`, `published.Tree`, `gateEvaluation.completionTree`, `closedFolder`, `gateEvaluation.completionClosure`, `gateEvaluation.applyCheckpoint`, `checkpointEvaluation`, `readerPlanner.decide`, `WithLanding`, `WithCompletion`, `holdWholeGate`, `Store.AdmitWindow`, `Store.PublishAdmitted`, `Store.AcquireLease`, `Store.refusingBuilds`, `Store.authorizeCandidate`, `jsonfile.DecodeExactDocument`, `publicationFacts`, `landingpolicy.ResumeMarker`, `advanceLandedMarker`, `reconcileLandingDestination`, `landedAfterEffects`, `routetest.Fixtures`, `routetest.Steps`, `repairsSource`, `retargetSourceTip`, `landing.RedKind`, `authorization.attribute`, `Store.openBuilds`, `Store.validateDeliverableAt`, and `usage.WorktreeUsage`. This session read each definition.
- `Import edges`: the new owner rule package imports only the standard library and the settle policy, which owns the `capture/` rule. The gate, the landing package, and the worktree package import it. `go list` shows that the gate imports neither the landing package nor the worktree package, so the new edges make no cycle.
- `Source-row clauses and occurrences`: the source-sentence table below quotes each clause of the two criteria and of D1 to D6.
- `Promised field labels`: the stack file fields `members`, `path`, `request`, `base`, `source_tip`, `spec`, and `message`. The override `--source-tip <request>=<commit>`. The table `members[N]{member,label,source_tip,published_commit,outcome,census}`. The record `stack{destination_base,published_commit,landed,worktree}`. The faces `stack-member-alone`, `stack-members-overlap`, and `stack-red`. The reader-plan case `landing stack`.
- `Changed-function callers`: `bench consumers --production` at `34ac699` names these callers. `gateEvaluation.applyCheckpoint` has one caller, `gateEvaluation.build`. `readerPlanner.decide` has one caller, `readerPlanner.plan`. `checkpointEvaluation` has four callers: `executeTreeWithOwner`, `inspectProspective`, `execute`, and `executeAfterAcquire`. `landAttributed` has one caller, `landWith`. `usage.WorktreeUsage` has one reference, `commandRegistry`.
- `Copy survival`: none.
- `Rendered-shape readers`: the reader-set refusal line keeps its format, and only a new case value joins it. `bench help` gains two rows, and `TestHelpInventoryIsComplete` pins them. The single `landed{...}` record does not change, and no production code parses it.
- `Pin operators`: the tally rows compare with `==`. The grade-sequence rows compare the exact ordered list of graded prefixes. The reader-plan row compares the exact line.
- `Entry reads`: the stack verb reads the stack file below its entry. The input rows drive the real verb over real files, a FIFO included.
- `Derived expectations`: the bisect rows derive their grade lists from the algorithm stated in Implementation decisions, not from the code. The members rows derive through `toon.Table`.
- `Consolidated rules`: none.
- `Quantified obligations`: each member runs the proofs of one landing (LS50, LS87), the stack proofs (LS70 to LS73, LS88), its override (LS95), the window (LS74), the admission over its own tree (LS28), the folder close (LS27, LS42), the override rule (LS19, LS20), the census and release (LS41, LS44), and the resume proof (LS79). Each prefix grade carries its own members (LS29).
- `Workflow-step writes`: the new drain step lands the verified light-path tickets as one stack. It writes the destination ref, the intent-ledger lease and overrides, the release receipts, the green marker, the reader hold of a red whole run, and the gate evidence. It reads the HEAD of each member worktree after the coordinator verifies each done-claim, and the ledger and the destination tip at the preflight. The capture inboxes are ignored here, so the step reads no capture source.

### Flagged additions

- The bound of sixteen members (LS60, LS62).
- The stack resume form (LS78 to LS82).
- The `landing stack` reader-plan case line (LS16, LS17).
- The `members` table and the `stack` record (LS44, LS45).
- The three new refusal faces (LS49, LS68, LS71, LS72, LS73, LS76, LS77).
- The `--source-tip <request>=<commit>` override and the rerun rules of a stack-raised refusal (LS86, LS89 to LS98).
- The owner rule package that the preflight and the proof share (LS1 to LS6).
- The binding query of the commitment store (LS22, LS23).
- For reviewer veto: an `inherited` red stops the stack with `land-handback` and no bisect (LS99). This reads D5's "red on top of a green prefix" to exclude a destination that held no green.
- For reviewer veto: the alone predicate reads "holds a delivery binding", which is broader than the words of D3 (LS73, LS88). D3 names a member bound to an open build, but a started tickets-only deliverable opens no build, and its closure must not join a stack.

### Source sentence to row table

| source clause | rows |
|---|---|
| Ready landing sources compose in order | LS25, LS40 |
| and grade once as one stack | LS24, LS40 |
| and each source keeps its own spec transition | LS27, LS42 |
| A red stack bisects with whole gates | LS16, LS30, LS35 |
| and the red is attributed to one source | LS31, LS33, LS34, LS49 |
| D2: one lease | LS21, LS52, LS53 |
| D2: one compare-and-swap publishes the own merge commit of each member | LS25, LS26 |
| D2: no persistent queue and no ledger field | Won't handle line, Persisted state |
| D3: a staged-spec member lands alone | LS72 |
| D3: a member that an open build binds lands alone | LS22, LS23, LS73, LS88 |
| D3: each member tip contains the destination tip | LS70 |
| D3: deltas overlap only on settled `capture/` paths | LS2, LS3, LS11, LS12, LS71 |
| D4: the whole gate, never the reader set | LS16, LS17 |
| D4: a green moves the green marker to the stack tip | LS43 |
| D5: publish the green prefix | LS32, LS46 |
| D5: leave each later member unlanded with its rerun | LS48, LS68 |
| D5: the worst case is 1 + ⌈log2 n⌉ whole gates | LS30, LS35, LS47 |
| D5: an infrastructure outcome publishes nothing | LS36, LS37 |
| D5: red on top of a green prefix (an `inherited` red has no green prefix) | LS99 |
| D6: each member passes the window admission on its own | LS18, LS74, LS75 |
| D6: each preflight refusal stops the whole stack and names each member | LS50, LS51 |

### Ticket sizing

Each ticket targets an author context of about 150k tokens. Each estimate uses the corrected model of the open learning of 2026-10-26: the `Writes:` bytes and the read set outside `Writes:` at about four bytes per token, 50k tokens of fixed context, 25k for the verification runs, about 5k for each witness probe, and 10k for one lane retry.

| ticket | `Writes:` bytes | read set | probes | estimate |
|---|---|---|---|---|
| 01-share-the-stack-owner-rule.md | about 10 KB of new files; the author does not read `.bench/commitment.json` | about 10 KB | 2 | about 100k |
| 02-prove-a-stack-in-the-gate.md | about 64 KB | about 45 KB | 3 | about 127k |
| 03-admit-a-stack-publication.md | about 26 KB | about 30 KB | 2 | about 110k |
| 04-compose-and-publish-a-stack.md | about 46 KB | about 50 KB | 2 | about 119k |
| 05-move-the-worktree-help-rows.md | about 38 KB | about 20 KB | 0 | about 100k |
| 06-read-a-stack-file.md | about 26 KB | about 30 KB | 2 | about 110k |
| 07-land-a-green-stack.md | about 136 KB | about 50 KB | 2 | about 142k |
| 08-bisect-a-red-stack.md | about 35 KB | about 40 KB | 2 | about 114k |
| 09-attribute-a-red-stack.md | about 101 KB | about 40 KB | 2 | about 130k |
| 10-refuse-a-stack-member.md | about 114 KB | about 45 KB | 2 | about 134k |
| 11-resume-a-stack.md | about 121 KB | about 36 KB | 1 | about 129k |
| 12-record-the-landing-stack.md | about 123 KB; `CHANGELOG.md` takes an append and counts as a 10 KB read | about 45 KB | 0 | about 127k |

Ticket 07 is the closest to the target, so its help rows and its wrapper row moved to ticket 11. A continuation of any ticket goes to a fresh narrow session, not a resume of the full author.

### Research route

By reviewer direction, this session ran the research and the slicing inline, with no research delegate, no authoring fork, and no paired slicing readers. An independent Opus reviewer at high effort graded the spec. These records replace the paired readers:

- The measurement method above, over the span store and Git.
- The `path:line` citations of the research findings above, each read in this session.
- The `bench consumers --production` answer at `34ac699` for each changed function.
- The fence-expansion classes of the slicing checks, walked for each ticket, with the sizes above.

### Gap log

Each entry takes one class: (a) spec-writer discretion, (b) a seam or test choice, (c) a missing fact, or (d) a reviewer decision.

- (a) The bound of sixteen members keeps the worst bisect at five whole gates. The measured largest ready set was nine sources.
- (a) A stack after a partial publication needs a fold of each unlanded member, because the destination moved to the published prefix. D5 accepts that cost.
- (a) The lease record names the assignment of the first member, so no ledger field changes.
- (a) The override names a member by its `request`, because a label can repeat and an index moves when a route removes a member.
- (a) A proof fault after the preflight refuses with `land-handback` as a defensive stop. The preflight refuses each known cause first, so no row reaches it.
- (a) Ticket 04 composes, grades once, and publishes a green stack. It stops a red top tree with no publication until ticket 08 adds the bisect.
- (a) Ticket 07 handles each result of the landing owner: k published members and an optional refusal. Until ticket 09 adds `stack-red`, it refuses a red member with `land-handback`. Until ticket 11 ships the stack resume, it ends a stack whose post-publication step stopped with `land-handback`. So no commit advertises a form before it works.
- (b) Ticket 05 moves the worktree help rows out of `cmd/bench/main.go` with no change of output, so ticket 11 adds both stack rows with no growth of that file.
- (a) Ticket 07 ships the stack form, and ticket 11 lists both stack forms in `bench help` and pins the wrapper route. This keeps ticket 07 under the size target, and no ticket advertises a form before it works.
- (a) From ticket 07 until ticket 10, the verb runs without the stack proofs of the preflight and without the window decision of each member before the lease. The gate proof of LS-C1 still refuses an overlap and a member that lacks the destination, and the commitment decision of ticket 03 still decides the window under the lock. The build lands as a whole, so this interim state never reaches `main`.
- (d) The alone predicate "holds a delivery binding" is broader than the words of D3. The flagged additions name it for reviewer veto.
- (b) The owner rule compares tree entries by path with no rename detection, so a rename owns two paths.
- (b) `internal/worktree/land.go` has 397 lines, so the stack dispatch stays at three lines there, and the stack flow takes its own files.
- (b) `cmd/bench/main.go` is over its line limit, so the growth ratchet reds any growth. The two help rows land with no net growth of that file.
- (b) The reader-plan rows take a new test file, because `internal/gate/reader_plan_test.go` is at its 400-line limit.
- (c) The staged spec `specs/shared-test-fixtures/` plans to change the landing test fixtures. The later build of that spec reads the stack fixtures as one more caller.
