# FT96 — the delegation discipline, batched

Status: staged

## Problem

Six things a session needs to know before it spawns a delegate are not written
down anywhere, and every one of them has already cost a real session.

A second concurrent `isolation: worktree` delegate is refused. The
`WorktreeCreate` hook derives its request ID from the harness session ID alone
(`internal/harness/worktree.go`), so `worktree.Create` finds an existing
assignment under a different label and returns `worktree create request
conflicts with its existing assignment`. The manual route
(`bench worktree create --request <opaque-id> --label <work-item>`) works and
respects the whole lifecycle, but no skill names it, so each session rediscovers
the refusal and improvises around it.

The whole-tree gate is a serialized resource, and nothing says so. Four
concurrent worktree `bench commit` gates flaked three load-sensitive contract
tests that all pass serially — a red that answers for machine load rather than
for any diff.

There is no sanctioned route for work no worktree can hold. FT122's ~1500-line
build sat uncommitted and largely untracked in the main checkout with the gate
red: it could not be committed first, and a worktree branched from HEAD would
not have contained the code under repair. That session improvised a
main-checkout delegate with a file allowlist, no commit authority, and a
`git status` check on return, and recorded the whole thing as unresolved.

Isolation's boundary is stated nowhere, and it is not where agents assume. A
worktree isolates the working tree, not repository-global git surfaces — the
stash stack above all. Two FT86 delegates in separate worktrees each reached for
`git stash push`/`pop`, the repository-global stash cross-applied their in-flight
edits, and neither diff could be attributed; one slice was discarded and re-run.
The destructive-git guard refuses `git checkout <path>` — the same hazard class —
but permits every working-tree-mutating stash verb.

Two charge-side defects round it out. A delegate's evidence that rests on an
absence is accepted without checking that the named identifiers exist: a
misspelled kit-only allowlist row once passed its contract by asserting the
absence of a path that never existed. And a charge that names prose conventions
instead of exemplar files degrades as the tree grows — "follow the repo's error
idiom" does not survive translation to a low-context delegate, while "mirror the
error shape in `internal/x/foo.go`" does. Relatedly, a write-delegate once
re-indented a workflow `run:` line and broke two canary fixtures that nothing in
its charged verification list could see, because the charge never named the gate
layer that owns workflow content.

All six live in one owner file, so they land as one diff.

## Solution

One pass over `.agents/skills/bench-craft-delegate/SKILL.md` closing all six
clauses, plus two enforcement arms so the clauses that assert a mechanical fact
are anchored to their source rather than left as prose that can rot.

The prose closes each clause in the section that already owns its subject:
Isolation gains the parallel-delegate route, the serialized-gate rule, the
shared-checkout exception, and isolation's true boundary; The charge gains the
gate-layer rule and the exemplar-file rule; Verifying the done-claim gains the
absence-evidence rule.

The first enforcement arm makes the stash ban real: the working-tree-mutating
stash verbs join `internal/gitguard`'s deny table, so the guard refuses them the
way it already refuses `git checkout <path>`, while `git stash list` and
`git stash show` stay runnable. The prose points at the guard rather than
re-listing the verbs, so the deny table stays the single source of what is
refused.

The second arm closes a hole the first clause would otherwise open: the
conformance sweep that catches a `bench <cmd>` reference naming no route today
covers `HANDOFF.md`, the two operating-guide files, and `.agents/commands/*.md`,
but not skills. Extending it to the `.agents` markdown tree means the CLI route
this spec writes into `craft-delegate` cannot silently become a dead pointer.

The third arm anchors the two clauses that assert a mechanical fact. The kit
already binds guidance prose to a literal through `checkWorkflowAnchors`, a
`require(file, needle)` table with its own canary fixtures. Adding rows for
`craft-delegate` means a later prune cannot quietly delete the parallel-delegate
route or the stash ban, and it closes the gap where a build could edit the skill
without ever writing the route in.

## Flagged for reviewer veto

No decision map backs this spec. It was compiled from the FT96 roadmap row under
`/bench-write-spec`'s reviewer-directed batch-drain override, so every decision
below was defaulted by the author rather than settled by you. Each names the
alternative not taken; veto any of them and the spec changes.

1. **Clause 1's fork resolved toward documentation.** The roadmap offered either
   re-keying hook assignments per delegate identity or documenting the
   `--request` route. This spec documents the route and defers the re-keying (see
   Out of scope), because FT96's own framing is one owner file and re-keying
   changes the create/remove contract on both ends.
2. **Clause 2's fork resolved toward sanctioning the shared checkout.** The
   roadmap offered either sanctioning it under exactly those conditions or naming
   the route that gets uncommitted work into a worktree first. This spec
   sanctions it (story 3). The unnamed alternative is dropped, not deferred: no
   such route exists today that would have held FT122's untracked build.
3. **Clause 3's "consider" promoted to a committed story.** The roadmap only
   suggested weighing `git stash` for the guard. Story 8 commits to it, because
   it is the one arm of this spec with a genuinely gate-observable red and
   because story 4's ban is otherwise advertisement with no enforcement.
4. **Story 9 is not on the roadmap row at all.** The sweep widening is the
   author's addition, admitted because story 1 writes a CLI route into a file the
   dead-pointer check does not currently read.
5. **Story 10 is not on the roadmap row at all.** The anchor rows are the
   author's addition, admitted because without them a build can satisfy every
   other row while never writing story 1's route or story 4's ban into the skill.
6. **The stash deny class splits in two rather than reusing the existing one.**
   Two labels, because destroying stash history and cross-applying working-tree
   state are different hazards and the agent-facing message names the hazard.

## User stories

1. As a coordinator spawning concurrent write-delegates, I want `craft-delegate`
   to name `bench worktree create --request <opaque-id> --label <work-item>` as
   the canonical parallel-delegate route and to say why the harness's own
   `isolation: worktree` refuses the second concurrent delegate, so that I stop
   rediscovering the refusal and improvising past it.
   Line: `gpt-5.6-sol` / high. `craft-line`'s leverage override routes any
   artifact that steers future generation to the top tier, because a defect in
   guidance prose is invisible to the gate and multiplies through every session
   that loads the skill.

2. As a coordinator running several worktree delegates at once, I want
   `craft-delegate` to state that the whole-tree gate is a serialized resource —
   delegates stop at "diff ready, focused tests green", and the coordinator runs
   `bench commit` per worktree one at a time — so that a red answers for a diff
   rather than for machine load.
   Line: `gpt-5.6-sol` / high. This clause is guidance prose in the same skill,
   so `craft-line`'s leverage override routes it to the top tier for the same
   reason: the gate cannot grade whether the rule is stated correctly.

3. As a coordinator facing a large uncommitted build the gate is red on, I want
   the shared-checkout exception written down with its exact conditions — one
   writer at a time, a named file allowlist, no commit authority, and a
   `git status` check verified on return — so that the exception is a rule I
   follow rather than an improvisation I record as unresolved.
   Line: `gpt-5.6-sol` / high. This clause loosens the isolation rule, so its
   wording is the only thing carrying the safety, and `craft-line`'s leverage
   override routes guidance prose of that weight to the top tier.

4. As a coordinator writing a charge, I want `craft-delegate` to state
   isolation's true boundary — a worktree isolates the working tree, not
   repository-global git surfaces, the stash stack above all — to ban `git stash`
   in a charge, and to name the per-worktree substitute the technique actually
   wants (copy the working file aside with `cp`, restore with
   `git show HEAD:<path> > <path>`, test, then copy back), so that two delegates
   cannot cross-apply each other's edits through a surface they both believed was
   isolated.
   Line: `gpt-5.6-sol` / high. The substitute has to be written precisely enough
   that a low-context delegate can follow it without inventing its own, which is
   exactly the prose judgment `craft-line`'s leverage override buys.

5. As a coordinator verifying a done-claim, I want the rule that a delegate's
   absence, exclusion, or withholding claim is accepted only after every
   identifier it names resolves to a real thing, so that a misspelled identifier
   cannot pass its contract by asserting the absence of something that never
   existed.
   Line: `gpt-5.6-sol` / high. This clause changes how a coordinator accepts
   evidence, and a rule stated loosely there weakens every future verification,
   which is the compounding cost `craft-line`'s leverage override prices in.

6. As a coordinator writing a charge, I want the rule that the charge names the
   gate layer owning each artifact class it touches — workflows and `.bench/`
   content to canary, gate output shape to canary, skills and commands to
   conformance — together with its converse for the delegate, that the named list
   is a floor rather than a ceiling, so that a delegate cannot break a layer
   nothing in its charge could see.
   Line: `gpt-5.6-sol` / high. The mapping and its converse have to be stated in
   one breath without either half reading as optional, and `craft-line`'s
   leverage override routes that kind of prose judgment to the top tier.

7. As a coordinator writing a charge, I want the rule that a charge names
   exemplar files to mirror, and says so explicitly when no exemplar exists, so
   that the charge does not degrade as the tree's unstated constraint set grows
   past its flat word count.
   Line: `gpt-5.6-sol` / high. `craft-skills` wants a contrastive good/bad pair
   for a rule like this, and writing a pair that teaches rather than decorates is
   prose judgment the gate cannot grade.

8. As an agent who reflexively reaches for `git stash`, I want the destructive-git
   guard to refuse the working-tree-mutating stash verbs while leaving
   `git stash list` and `git stash show` runnable, so that story 4's ban is
   enforcement rather than advertisement.
   Line: `gpt-5.6-terra` / medium. The profile's cached "Gate / conformance
   logic" row covers guard classification, and the seam is an existing pure
   verdict function with direct unit-test prior art beside it.

9. As a reviewer, I want the conformance dead-pointer sweep to cover the
   `.agents` markdown tree rather than only `.agents/commands`, so that a CLI
   route named in a skill cannot rot into a dead pointer the way story 1's route
   otherwise could.
   Line: `gpt-5.6-terra` / medium. The same cached "Gate / conformance logic" row
   applies, and the change is a file-set widening of an existing check plus the
   canary fixture that proves the widened path bites.

10. As a reviewer, I want the two clauses that assert a mechanical fact — story
    1's parallel-delegate route and story 4's `git stash` ban — anchored by the
    conformance anchor table, so that a build cannot satisfy this spec without
    writing them in and a later prune cannot silently delete them.
    Line: `gpt-5.6-terra` / medium. The same cached "Gate / conformance logic"
    row applies, and the work is adding rows to an existing `require` table plus
    the canary fixture that proves the new rows bite.

## Implementation decisions

**One owner file for the prose.** All seven prose clauses land in
`.agents/skills/bench-craft-delegate/SKILL.md`. `.claude/skills/bench-craft-delegate`
is a symlink to it, so there is no mirror to update; the conformance mirror check
asserts presence, not content.

**The prose points at the guard; it never re-lists the verbs.** Story 4's ban
says the destructive-git guard refuses `git stash` and names the substitute. It
does not enumerate which stash verbs are refused — `internal/gitguard`'s deny
table is the single source of that fact, and an enumeration in the skill would be
a second derivation that drifts.

**Legibility is a build constraint, not a story.** `craft-synthesis`'s legibility
loop and `craft-skills`' pruning rules bind this diff: the six clauses land inside
the sections that already own their subjects, adding at most one new heading, and
existing sentences that the new clauses subsume get deleted rather than left
beside them. A skill that grows past its legibility ceiling has failed even when
every clause is correct.

**Stash classification splits into two deny classes.** `internal/gitguard`'s
`stashVerdict` currently returns a verdict only for the free-arg subcommands
`drop` and `clear`, under the `stash` key labelled `git stash drop`. It gains a
second class for the working-tree-mutating verbs — a bare `git stash` (which is
`push`), and explicit `push`, `save`, `pop`, `apply`, and `branch` — under its own
key and label, because the two hazards differ: one destroys stash history, the
other cross-applies working-tree state between worktrees. `list` and `show` are
read-only and stay allowed. The deny table stays the ordered single source; the
guard script's static manifest header already advertises the class generically
and needs no edit.

**The bare-`git stash` case inverts the current default.** `stashVerdict` today
answers "allow" when there is no free argument. Since bare `git stash` is
`git stash push`, that default flips: absence of a free argument is now a block,
and the allow set is the explicit read-only verbs.

**The dead-pointer sweep derives its file set from the existing `.agents` walk.**
`checkColdPickupCLILists`'s reverse check currently globs `.agents/commands/*.md`
by hand. It switches to the same `walkConformanceDocs(.agents)` helper the slash
sweep in `checkStaleCommandReferences` already uses, filtered to `.md`, so the two
sweeps cannot disagree about which agent-facing documents are in scope. The tree
is green under the widened set today: every backticked `bench <cmd>` under
`.agents` resolves to a route in `bin/bench.sh`.

**The anchors go in the existing table, not a new check.**
`checkWorkflowAnchors` is already a `require(rel, needle)` table binding guidance
prose to literals across the command and skill files. Story 10 adds rows for
`.agents/skills/bench-craft-delegate/SKILL.md` — the literal
`bench worktree create --request` and the literal `git stash` — rather than
introducing a second anchor mechanism. What a presence anchor catches is exact
and worth stating: deletion of the clause, and a build that never wrote it. It
does not grade whether the surrounding sentence is correct; that stays the
reviewer's, which is why stories 1 and 4 still route top.

**Both new canary fixtures need registry entries.** Each fixture under
`tests/canary/` is registered in `canaryFixtureRegistry`
(`internal/conformance/registry_test.go`) with its owning gate entry and shell
sources, and listed in the bite test. Stories 9 and 10 each add one fixture under
`tests/canary/docs-currency-token-diet/`, registered exactly as its siblings in
that family are: story 9's plants a dead `bench <cmd>` reference in a skill file,
story 10's deletes an anchored literal from `craft-delegate`. Copy the sibling
registration rather than inventing one. The existing `stale-cli-doc-reference`
fixture keeps covering the original path.

## Testing decisions

A good test here exercises the seam's external behavior: for the guard, a command
string in and a deny label out, with the repository `Checker` injected; for the
conformance sweep, a tree in and a diagnostic list out, with the canary proving
the check bites rather than matching nothing. Neither arm is tested by reading
the diff.

The prose clauses have no behavioral seam. Two of them assert a mechanical fact
and get a presence anchor (story 10), which catches deletion and a build that
never wrote them; the other five are graded by the reviewer alone, and their rows
say so rather than dressing structural checks up as behavioral coverage. That
split is why stories 1–7 route to the top tier: an anchor proves a string is
present, never that the sentence around it is right.

Prior art: `internal/gitguard/verdict_test.go` for the guard seam;
`internal/conformance/docs_workflow_checks_test.go` plus
`tests/canary/docs-currency-token-diet/stale-cli-doc-reference` for the sweep;
`internal/conformance/docs_workflow_helpers_test.go` plus the
`acceptance-coverage-anchor` and `command-handoff-anchor` fixtures for the
anchors.

Gate command: `.bench/gate.sh` (the project gate).

### Seam diagram

Seam 1 — the guidance surface (structural checks only):

    trigger: a session loads the skill; the gate's conformance phase reads it
        │
        ▼
    SKILL.md frontmatter ───▶ [ conformance docs+skills ] ──▶ diagnostics
    backticked `bench <cmd>` ▶ [   checks, incl.        ]
    `/bench-*` references ───▶ [   checkWorkflowAnchors  ]
    anchored literals ───────▶ [                        ]
                  ◀ tests attach here: go test ./internal/conformance
                    -run '^TestRootConformance$' over the real tree, plus a
                    canary fixture per new bite path. Presence and liveness are
                    checkable; the clause's *meaning* is not — the reviewer
                    grades that

Seam 2 — destructive-git classification:

    trigger: the agent's Bash tool call → PreToolUse hook → `bench guard-git`
        │
        ▼
    command string ──▶ [ gitguard.Classify → scan → classify → ] ──▶ deny label
    Checker (repo) ──▶ [   stashVerdict                        ]     or "" (allow)
                  ◀ tests attach here: table-driven unit tests call Classify
                    with a literal command and a stub Checker, assert the label

Seam 3 — the docs dead-pointer sweep:

    trigger: the gate's conformance phase (and `bench canary` per fixture)
        │
        ▼
    bin/bench.sh routes ─▶ [ checkColdPickupCLILists ] ──▶ diagnostic strings
    .agents/**/*.md    ─▶ [   reverse dead-pointer   ]
    HANDOFF/BENCH docs ─▶ [   check                  ]
                  ◀ tests attach here: TestRootConformance over the real tree
                    for green; a canary fixture planting a dead ref for red

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | the skill contains the literal `bench worktree create --request` | conformance (`internal/conformance`) | `go test ./internal/conformance -run '^TestRootConformance$'` — passes today, which is the failure: no anchor requires the string | This is the row that forbids the cheapest wrong build of story 1 — edit the skill, never write the route in. Every other story-1 row is satisfied by that build. |
| 1 | whether the surrounding clause explains the refusal correctly | conformance (`internal/conformance`) | not TDD-able — a presence anchor proves the string is there, never that the sentence around it is right; the reviewer is the oracle | Recorded so the anchor above is not mistaken for semantic coverage. This is why story 1 routes top. |
| 1 | the route the clause names resolves to a real `bin/bench.sh` command | conformance (`internal/conformance`) | covered by story 9 — `go test ./internal/conformance -run '^TestRootConformance$'` after the sweep widens; green today because the reverse check does not read skills | Without story 9 this clause could name `bench worktee` and nothing would notice. The widened sweep is what turns a dead route into a red gate. |
| 2 | the skill states the serialized-gate rule and where a delegate stops | conformance (`internal/conformance`) | not TDD-able — guidance content, reviewer-graded | Same as story 1: no behavioral seam exists for prose. |
| 3 | the skill states the shared-checkout exception with all four conditions | conformance (`internal/conformance`) | not TDD-able — guidance content, reviewer-graded | The four conditions are a policy decision; the gate can see that the section exists, never that it is complete. Enumerated in the story so review has a checklist. |
| 4 | the skill contains the literal `git stash` | conformance (`internal/conformance`) | `go test ./internal/conformance -run '^TestRootConformance$'` — passes today, which is the failure: no anchor requires the string | Forbids the build that ships the boundary paragraph and drops the ban, and stops a later prune from deleting it silently. |
| 4 | the boundary, the ban's wording, and the substitute's four steps are right | conformance (`internal/conformance`) | not TDD-able — guidance content, reviewer-graded | The substitute is named concretely in story 4 so review has something exact to check; the ban's *enforcement* is story 8, which does have a red. |
| 5 | the skill requires absence-claim identifiers to resolve before acceptance | conformance (`internal/conformance`) | not TDD-able — guidance content, reviewer-graded | Same as above. |
| 6 | the skill maps artifact class to gate layer and states the floor-not-ceiling converse | conformance (`internal/conformance`) | not TDD-able — guidance content, reviewer-graded | Same as above. The converse is named separately here so a build that ships only the mapping is visibly short. |
| 7 | the skill requires exemplar files in a charge, and an explicit note when none exists | conformance (`internal/conformance`) | not TDD-able — guidance content, reviewer-graded | Same as above. The no-exemplar converse is named so it cannot be dropped silently. |
| 1–7 | the edited skill still parses: frontmatter, `index:` field, generated skills index | conformance (`internal/conformance`) | already covered — `go test ./internal/conformance -run '^TestRootConformance$'` is green today and must stay green | An edit that damages the frontmatter or drifts the `index:` line from the generated index turns this red without any new check. |
| 1–7 | the edited skill names no retired `/bench-*` or `$bench-*` command | conformance (`internal/conformance`) | already covered — the same `TestRootConformance` run; `.agents` is already in the slash sweep's file set | New prose that points a reader at a phase command must point at a live one; this bites today for slash references even though the CLI reverse check does not. |
| 8 | a bare `git stash` is refused | unit (`internal/gitguard`) | `go test ./internal/gitguard -run TestClassifyStashMutations` — fails, `stashVerdict` returns allow when there is no free argument | Bare `git stash` is `git stash push`, and it is the exact form the FT86 delegates ran. An implementation that only adds explicit verbs leaves the common case open, and this assertion names it. |
| 8 | `push`, `save`, `pop`, `apply`, and `branch` are each refused | unit (`internal/gitguard`) | `go test ./internal/gitguard -run TestClassifyStashMutations` — the five named subcases | The set is enumerated rather than left as "the mutating verbs" so a build cannot pick the cheapest reading and ship `push` alone. `apply` and `pop` are the cross-application half; `save` is the deprecated spelling of `push`; `branch` applies and branches. |
| 8 | `git stash list` and `git stash show` stay allowed | unit (`internal/gitguard`) | `go test ./internal/gitguard -run TestClassifyStashReadOnly` | The degenerate implementation is "block every `git stash`", which passes both rows above. This is the pair that forbids it, and it keeps the read-only inspection an agent legitimately needs. |
| 8 | `drop` and `clear` keep their existing label, not the new one | unit (`internal/gitguard`) | `go test ./internal/gitguard -run TestClassifyStashMutations` — the label assertion on the `drop` and `clear` cases | Collapsing both classes onto one label loses the distinction between destroying stash history and cross-applying working-tree state, and the agent-facing BLOCKED message names the wrong hazard. |
| 8 | the new class is a `denyTable` row, so its label reaches `BlockMessage` | unit (`internal/gitguard`) | `go test ./internal/gitguard -run TestClassifyStashMutations` — asserting the returned label equals the table's entry for the new key | A hardcoded literal identical to the table entry still passes, so this row pins drift rather than sourcing: it turns red the moment the table's label is reworded and the verdict's copy is not. |
| 9 | a dead `bench <cmd>` reference inside a skill file produces a diagnostic | conformance (`internal/conformance`) | `go test ./internal/conformance -run '^TestRootConformance$'` against a tree with `bench worktee` planted in a skill — passes today, which is the failure | The check reads only `.agents/commands`, so a skill can name a removed route and the gate stays green. The planted reference is the minimal demonstration. |
| 9 | the widened sweep is green on the real tree | conformance (`internal/conformance`) | already covered — `go test ./internal/conformance -run '^TestRootConformance$'` must stay green after the widening | Every backticked `bench <cmd>` under `.agents` resolves today, so a widening that turns the gate red has picked up files or a pattern it should not have. |
| 9 | the reverse sweep and the slash sweep agree on the `.agents` file set | conformance (`internal/conformance`) | not TDD-able — no persisted check forbids re-adding the glob beside the walk; review grades the single-sourcing against this repo's code standard | Stated as a row rather than left implicit because keeping the hand-rolled glob would satisfy every other story-9 row while leaving two derivations of "which agent-facing documents are in scope". |
| 9 | the gate goes red on a tree with a dead skill CLI reference | canary (`tests/canary/`) | `bench canary` against a new `docs-currency-token-diet` fixture whose `EXPECT` names the diagnostic — the fixture does not bite before the check widens | Proves the widened sweep actually bites rather than matching nothing, which is the failure mode a file-set change is most prone to. |
| edge of 8 | `git stash push -m "wip thing"` with a quoted multi-word message is refused | unit (`internal/gitguard`) | `go test ./internal/gitguard -run TestClassifyStashMutations` — the quoted-message case | The tokenizer splits on quoting; a verdict that scans raw text rather than tokens misreads the message as a subcommand and allows the command. |
| edge of 8 | `git stash push -- path/with space` is refused | unit (`internal/gitguard`) | `go test ./internal/gitguard -run TestClassifyStashMutations` — the pathspec case | Named in the profile's hostile-input checklist. A `--` separator and a spaced path must not change the verdict; a free-arg scan that reads past `--` could pick the path as the subcommand. |
| edge of 8 | `bash -c 'git stash pop'` is refused through one level of wrapper | unit (`internal/gitguard`) | `go test ./internal/gitguard -run TestClassifyStashMutations` — the wrapper case | The scanner recurses exactly one level by design; asserting it here confirms the new class rides the existing recursion instead of being wired only into the top-level path. |
| edge of 8 | `git stash` reached through `xargs` is refused | unit (`internal/gitguard`) | `go test ./internal/gitguard -run TestClassifyStashMutations` — the xargs case | `checkoutVerdict` and `restoreVerdict` take a `viaXargs` flag because arguments arrive from stdin and cannot be inspected. Stash needs the same posture, and without a row the build would silently allow it. |
| edge of 9 | a skill directory holding markdown outside `SKILL.md` is swept too | conformance (`internal/conformance`) | `go test ./internal/conformance -run '^TestRootConformance$'` against a tree with a dead reference planted in a `references/*.md` file | The walk returns every `.md`; a widening implemented as a `SKILL.md` glob would pass the main story row and leave the reference files unswept. |
| edge of 9 | a non-markdown file under `.agents` produces no CLI diagnostic | conformance (`internal/conformance`) | `go test ./internal/conformance -run '^TestRootConformance$'` must stay green | `walkConformanceDocs` also returns `.sh`, `.json`, and YAML. Feeding those to a regex written for prose invites a false positive, so the filter to `.md` is asserted by the tree staying green rather than by a new fixture. |
| 10 | both anchored literals are required, not one | conformance (`internal/conformance`) | `go test ./internal/conformance -run '^TestRootConformance$'` — the two new `require` rows, each asserted; passes today, which is the failure | The set is enumerated as two rows rather than "the mechanical literals" so a build cannot add one anchor and call the quantifier satisfied. |
| 10 | the gate goes red when an anchored literal is deleted from the skill | canary (`tests/canary/`) | `bench canary` against a new `docs-currency-token-diet` fixture that removes `bench worktree create --request` from the skill, `EXPECT` naming the anchor diagnostic | Proves the new `require` rows actually bite. An anchor added to the table but never exercised is the always-pass failure the canary layer exists to catch. |
| 10 | the anchor diagnostic names the file and the missing needle | conformance (`internal/conformance`) | covered by the canary row above — the fixture's `EXPECT` is a substring match against the diagnostic naming both | An anchor that fails with a bare "missing anchor" sends the next session to the wrong file; the fixture's `EXPECT` cannot match unless the diagnostic carries both. |
| edge of 10 | a tree with no `craft-delegate` skill yields the missing-file diagnostic, not a silent pass | conformance (`internal/conformance`) | covered by the existing helper — `checkWorkflowAnchors`'s `require` already emits an anchor-file-missing diagnostic when the path is absent | Absent versus present-but-lacking-the-needle are distinct states from the profile's checklist; reusing the existing table means the absent case is handled by construction rather than reimplemented. |

### Edge inventory

Walked per behavior against the canonical classes and the shell-CLI
hostile-input checklist in `projects/benchkit.md`. Each landed as a coverage row
above or as a **Won't handle** line here.

- **Won't handle** — anchoring the five clauses that assert no mechanical fact
  (stories 2, 3, 5, 6, 7): a presence anchor on an ordinary phrase pins wording
  the reviewer should stay free to improve, and it would read as semantic
  coverage the check cannot deliver. Those five stay reviewer-graded.
- **Won't handle** — a `bench <cmd>` reference not wrapped in backticks: the
  existing sweep's regex requires the backtick, and loosening it would flag
  ordinary prose that happens to say "bench status". Inherited constraint, not a
  new one.
- **Won't handle** — flag-level validation (`--request`, `--label`): the sweep
  resolves the subcommand token only, so a wrong flag in a skill stays
  reviewer-caught. Making flags checkable needs a usage-string parser, which is
  the separate capability named in Out of scope.
- **Won't handle** — `git stash` behind two or more levels of wrapper: the
  analyzer scans exactly one level deep by design, recorded in the package
  doc's threat model. This is an honest-mistake layer, not an evasion-resistant
  boundary.
- **Won't handle** — git reached other than through the agent's Bash tool (an
  MCP git tool, a delegate's own harness surface): the guard's boundary is
  `PreToolUse:Bash` and this spec does not move it.
- **Won't handle** — control bytes or non-UTF-8 in a skill file: the conformance
  reader already owns file-state classification for the tree it grades, and this
  change adds no new reader.
- **Won't handle** — a repository with no `.agents` directory: `walkConformanceDocs`
  returns nothing for a missing directory and the check already returns early on a
  missing `bin/bench.sh`, so a stripped fixture yields no diagnostics.
- Interrupted or partial state, re-run idempotency, and hostile environment: both
  code arms are pure reads over a tree with no writes and no external process, so
  a second run of either produces the same diagnostics. No row needed.

## Out of scope

- **Re-keying the `WorktreeCreate` hook request per delegate identity** — clause
  1's alternative fork. It is a separate capability because it changes the
  request-derivation contract on both ends: `claudeRemove` derives the same ID
  from `session_id` alone and would need a recoverable key, which is a lifecycle
  decision rather than a documentation one. Estimate: 6 edits, 4 gate runs.
- **A usage-string parser that validates documented `bench` flags against the
  CLI** — a distinct check with its own source of truth (the Go usage strings,
  not `bin/bench.sh`'s route list). Estimate: 5 edits, 3 gate runs.
- **FT131's guidance arm** — naming the `scripts/go-build.sh` rebuild where the
  phase names the `dist/bench`-driving seams. The roadmap admits it to this diff
  only if FT131's in-helper staleness check proves unreliable; FT131 is unbuilt,
  so the condition cannot be evaluated yet. Estimate: 1 edit, 1 gate run when
  triggered.
