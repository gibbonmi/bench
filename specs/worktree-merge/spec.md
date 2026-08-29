# Worktree merge

Status: staged

Roadmap: FT252

Decision source: named reviewed artifact `roadmap/FT252.md` (drained 2026-08-28, `Next: spec`; five occurrences read in full).

Verification log: 1 iteration(s) to accept — round 1 (opus / medium) returned 4 blocking and 10 non-blocking findings; all 14 folded, see Further notes.

## Problem

No Bench verb moves a retained worktree onto a new base. An integration
source that must absorb a moved `main` reaches its base only through raw Git
inside the worktree. A delegate worktree that must start from the integration
tip has the same route. A coordinator folds a sibling's diff into the
integration source with `git apply`, `git apply --3way`, or a whole-file `cp`,
and resolves the conflicts by hand. Both routes break the rule
that the lifecycle runs through Bench verbs. The landing's own conflict repair
advertises the raw route, because no verb exists.

## Solution

`bench worktree merge --from <commit|target> <target>` composes one commit into
one active owned worktree by merge. The incoming commit is a commit in the
default branch's history or the branch tip of a sibling assignment. The verb
proves the target's identity bundle and refuses a dirty checkout. It composes
the pair through the landing owner's merge-tree composition and runs the
root's fast lane on the composed tree. A lane pass publishes a two-parent
merge commit onto the worktree branch by compare-and-swap.

When the target tip is an ancestor of the incoming commit, the branch moves by fast-forward. A
conflict outside the capture rule table refuses and names every conflicting
path. The guidance and the landing's repair text name the verb. The raw route
survives only as the hand resolution of a named conflict.

## User stories

### Group A — composition and publication

Line: opus / low. The seams exist, the spec is exact, and the ordinary test
phase covers the package.

1. As a coordinator, I want one verb to merge a commit into an owned worktree,
   so that raw Git never moves the integration source.
2. As a coordinator, I want a target behind the incoming commit moved by
   fast-forward, so that a delegate worktree starts exactly at the integration
   tip.
3. As a coordinator, I want a target that already contains the commit reported
   as `kind=current` with nothing changed, so that a re-run is idempotent.
4. As a reviewer, I want the previous tip first and the incoming commit second
   among the merge's parents, so that first-parent history stays the
   worktree's.
5. As a coordinator, I want the verb to derive the merge commit subject, so
   that no `-m` flag exists and the log reads one way.
6. As a coordinator, I want `--from` to accept a sibling assignment and take its
   branch tip, so that committed delegate work folds through one verb.
7. As a coordinator, I want the record to name the kind, both tips, and the
   tree, so that the next phase reads identities there.
8. As an agent, I want the target operand to take a label, an id, a prefix, or
   a path, so that the address matches `exec`.

### Group B — refusals before publication

Line: opus / low. Each refusal is one predicate at the command seam, and the
fixture harness exists.

9. As a coordinator, I want a conflict outside the capture table refused with
   every path named, so that a hand resolution starts from a list.
10. As a coordinator, I want a conflicted capture path settled by the landing's
    rule table and disclosed, so that the merge and the landing agree.
11. As a coordinator, I want a dirty target checkout refused before composition,
    so that no uncommitted work is reset away.
12. As a coordinator, I want a dirty or detached sibling refused with the repair
    named, so that a fold never drops uncommitted delegate work.
13. As a coordinator, I want a `--from` naming both an assignment and a commit
    refused as ambiguous, so that a ref-shaped label merges nothing.
14. As a coordinator, I want a `--from` that resolves to the target itself
    refused, so that a self-merge cannot hide an operand error.
15. As a coordinator, I want a `--from` that names nothing refused with the
    value named, so that a typo is never a silent no-op.
16. As a coordinator, I want a failed identity component refused by name, so
    that the verb never writes to a tree Bench does not own.
17. As a coordinator, I want a branch tip unequal to the checkout HEAD refused
    with both named, so that the swap's old value is visible.
31. As an agent, I want `--from` required with no default, so that a bare call
    never merges the default branch silently.
32. As a reviewer, I want a `--from` commit outside the default branch and every
    sibling refused, so that the lane never runs an unowned tree.

### Group C — the lane and the publication boundary

Line: opus / medium. The group edits gate authorization logic, which the
profile routes to medium effort.

18. As a reviewer, I want the composed tree to pass the root's fast lane before
    the branch moves, so that a broken merge never publishes.
19. As a reviewer, I want a lane failure to refuse, name the check, and
    change nothing, so that a red lane reads as a refusal.
20. As a reviewer, I want the branch ref updated by compare-and-swap on the
    previous tip, so that a tip moved during the lane refuses.
21. As a coordinator, I want a failed reconcile reported at exit 3 naming
    the published commit, so that the boundary reads apart from a refusal.
22. As a reviewer, I want a checkout changed during the lane refused before the
    ref moves, so that the reconcile never resets a fresh edit.
23. As a coordinator, I want a fast-forward to run the lane too, so that an
    arbitrary commit never places an ungraded tree on the branch.
24. As a reviewer, I want the lane's prose check to grade the Markdown the
    merge changed, so that incoming prose is graded like named
    Markdown.

### Group D — advertisement and guidance

Line: opus / medium. The group edits guidance prose, which the leverage rule
routes to the mid tier, and the reviewer's 2026-08-26 direction caps every
subagent at medium.

25. As a cold agent, I want `bench help` and `bench worktree --help` to list the
    verb and its grammar, so that discovery needs no guesswork.
26. As a cold agent, I want the profile to describe the verb beside the other
    worktree verbs, so that the seam paragraph stays the
    advertisement.
27. As a coordinator, I want the landing's conflict repair to say the verb
    refuses that conflict, so that `next=` never names an impossible repair.
28. As a reviewer, I want the `craft-delegate` stale-base rule, anchor
    needle, and fixture-bite case to name the verb, so that guidance and
    oracle agree.
29. As a reader, I want the reference's repair list and ADR 0014's consequence
    to name the verb, so that no document denies the verb.
30. As a reader, I want `CONTEXT.md` to define `sibling assignment` and
    `worktree merge`, so that each term has one meaning.

## Implementation decisions

- The verb is `bench worktree merge --from <commit|target> <target>`. `--from`
  is required and has no default. The record names the exact commit it
  resolved. There is no `-m` flag. The subject is
  `merge: compose <from-spelling> <8-char commit> into <label>`, the shape
  hand landings already use.
- One verb serves both faces of FT252. A move onto a new base is a merge from
  the default branch tip or from any commit. A fold is a merge from a sibling
  assignment's branch tip. The direction is always into the target.
- `--from` resolves in two lookups. The assignment lookup uses the target
  grammar over active, bundle-valid assignments: label, id, 8–12 character
  prefix, or absolute path. The commit lookup peels `<from>^{commit}` at the
  repository root and then requires the commit to be an ancestor of the
  resolved default branch tip. When both resolve, the verb refuses as
  ambiguous and names the assignment id and the full commit.
- When neither `--from` lookup resolves, the verb refuses and names the value. A commit outside the default
  branch's history that is no sibling tip refuses and names the commit. A
  sibling that is the target refuses.
- Bootstrap authority: the lane builds and runs an executable from the
  composed tree, so the incoming bytes must come from a trusted root. The two
  roots are the default branch, whose every commit passed a whole-project
  gate, and a Bench-owned assignment branch, whose every commit passed the
  lane. The identity bundle authenticates the sibling, and the default branch
  ancestry authenticates the commit. No other object reaches the lane.
- The verb needs the assignment record, not only its path: the id, the label
  for the subject, and the branch for the compare-and-swap. It selects the
  assignment the way the target-taking verbs do, checks the active state, and
  validates the creation bundle. It does not call the path-only resolver.
- A sibling contributes its branch tip only. The sibling's checkout must have
  its branch checked out at that tip, and it must be clean. The refusal names
  `bench commit` at the sibling. `bench commit` stays the one snapshot
  composer, so the merge verb never composes a working tree.
- The kind is decided by ancestry before composition. A target tip that
  contains the incoming commit is `current`: exit 0, no lane, no change. A
  target tip that is an ancestor of the incoming commit is `fast-forward`.
  Then the lane grades the incoming tree, and the branch moves to the incoming
  commit with no new object. Every other pair is `merge`.
- Ancestry is reflexive, so an equal pair satisfies both tests; the `current` test runs
  first and decides it. The ancestry query separates exit 1 from a Git
  failure, so a failed query refuses instead of classifying as `merge`.
- Composition and publication live in the landing owner as one merge operation
  beside the reviewed landing. It resolves exact commits and composes through
  the existing merge-tree composition with the capture rule table. It
  authorizes the graded tree through the owner's worktree-commit policy, which
  accepts a graded green or a lane pass. It rechecks the branch tip and the
  checkout fingerprint, creates the two-parent commit, and updates the branch
  ref by compare-and-swap on the previous tip. A conflict outside the table
  returns the landing's conflict error with its paths.
- The verb has the worktree commit's authority class, not the landing's. The
  identity bundle is the proof: state, path, owner marker, registration, and
  lock, plus the checked-out branch. The registry gains no branch component;
  the branch refusal names the assignment branch ref in its detail, as the
  landing's does. The verb takes no `--request`. It resolves the
  ledger from the repository's common directory, so it runs from the primary
  checkout or from any worktree.
- The lane is the root's declared fast lane, resolved for the target worktree
  the way `bench commit` resolves it. A root with no declared lane keeps the
  whole-project gate, as the commit does. The lane's prose subject is the set
  of Markdown paths that differ between the previous tip's tree and the
  composed tree. The lane authority gains one input, the previous tip. When
  that input is set, it derives the subject by a tree diff against the graded
  tree, so composition runs once.
- The target checkout must be clean before composition; any status entry,
  untracked included, is dirty. The verb takes its fingerprint before the lane
  and rechecks it before the ref update. After the update, the verb reconciles
  the checkout with `git reset --merge <published>` at the target alone. The
  landing's residue bracket does not apply. The pre-check and the fingerprint
  recheck already hold the precondition it stands in for. Ignored residue in
  a worktree is not a refusal.
- Exit codes follow the publication boundary. `0` is published or current,
  `1` is a refusal before publication, `2` is invalid usage, and `3` is
  published with the checkout not reconciled. The exit-3 record is
  `merged{...,next=git -C <path> reset --merge <published-commit>}`, and a path
  that is not line-safe takes the `bench worktree exec <id> --` pointer form,
  as `bench commit`'s remainder record does.
- The stdout record is
  `merged{worktree=<assignment id>,from=<commit>,kind=<current|fast-forward|merge>,previous_tip=<commit>,tip=<commit>,tree=<tree>}`.
  The assignment id addresses the worktree, because it is line-safe by
  construction. Every refusal renders on stdout as the landing's
  `refused{detail=...}` record, with the path table when paths exist, at exit
  1. A target or `--from` resolution error wraps into that record; the verb
  never uses the stderr line the `exec` and `path` verbs print. Usage errors
  print usage on stderr at exit 2. A settled capture path is disclosed on
  stderr as `merge composition{resolved=<path>:<side>,...}`, the landing's
  line under the merge prefix.
- Registration follows `reauthorize`: the grammar constant, the worktree usage
  block, the registry help rows, and the nested `worktree` dispatch. The
  `bench worktree --help` description row adds `merge` to its enumeration. The
  wrapper already routes every non-`land` worktree child as porcelain. The
  child is a mutating verb, so the AXI query set and the wrapper routes do not
  change.
- The landing's conflict `next=` keeps the raw merge as the hand resolution.
  It drops the claim that no verb exists and states that `bench worktree
  merge` refuses the same conflict. The `craft-delegate` stale-base sentence
  names `bench worktree merge --from main <target>` for the coordinator. Its
  anchor needle and its fixture-bite case move with the sentence in the same
  commit. That skill file sits exactly at its line budget, so the edit replaces
  the sentence line-neutrally. If the sentence needs a line, the same commit
  raises the profile's budget row by exactly that count.
- `internal/worktree` and `internal/landing` are over their file budgets
  already. The new files follow the `reauthorize` and `composition` precedent,
  and the structure verb stays advisory.

## Testing decisions

- A good test drives the public verb in process against a real fixture
  repository. It passes argv and reads the record and the exit code. Then it
  reads the branch ref, the checkout HEAD, the status, and the lane record
  through Git.
- The command seam is `MergeCommand` in the worktree package, on the package's
  fixture repositories and assignments. The lane fixture is the worktree
  package's own inline phase manifest, extended with controllable check
  scripts; the commit package's lane harness is not pasted. WM24's prose check
  is a script that records its argv.
- The owner seam is the landing owner's merge operation, tested in the landing
  package on fixture repositories as the composition tests are. It grades the
  parents, the kind decision, and compare-and-swap loss through the owner's
  update-ref field.
- The fault seams are the worktree package's joins value and the owner's
  update-ref field. The joins value carries the reconcile failure and a
  checkout edit between the lane and the ref update.
- The registry help and usage tests pin the advertised rows. The fixture-bite
  test pins the moved anchor. The landing surface tests pin the conflict
  repair text.
- The landing gate observes the feature through the ordinary `test` phase and
  the `bench-sh-routes` and `subcommand-routing` conformance checks. The
  fixture-bite check grades the moved anchor, the `guidance-prose-budgets`
  check grades the skill file's line count, and the prose check grades the
  edited Markdown.

### Seam diagram

    trigger: a coordinator, from the primary checkout or any worktree
        │
        ▼
    argv ──▶ [ usage.Parse ] ──▶ [ target + --from resolution, identity bundle, dirty checks ]
                                        │
                                        ▼
                     [ landing owner: kind, merge-tree composition, lane authorization,
                       fingerprint recheck, commit-tree, compare-and-swap ]
                                        │
                                        ▼
                     [ checkout reconcile ] ──▶ merged{...} record, exit 0 / 1 / 2 / 3
                      ◀ tests attach at MergeCommand with fixture repositories and a lane manifest;
                        owner tests attach at the merge operation with its update-ref field

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| WM1 | 1 | on a diverged target, `merge --from <commit>` publishes a commit whose tree equals the merge-tree of the pair, and the checkout HEAD equals that commit | command seam, lane fixture | a verb that moves only the ref leaves the checkout at the previous tip |
| WM2 | 2 | a target whose tip is an ancestor of `--from` prints `kind=fast-forward`, and the branch tip equals the incoming commit with no new commit object | command seam | a merge-always implementation mints a redundant merge commit |
| WM3 | 3 | a target that already contains `--from` prints `kind=current` at exit 0, and the tip, the checkout, and the lane record are unchanged | command seam | a verb that reruns the lane or mints an empty merge changes state on a re-run |
| WM4 | 4 | the published commit's first parent is the previous tip and its second parent is the incoming commit | owner seam | swapped parents move the first-parent history the review reads |
| WM5 | 5 | the published commit subject is `merge: compose <from-spelling> <8-char commit> into <label>` | command seam | a hand-typed message varies by session |
| WM6 | 6 | `--from <sibling label>` publishes a commit whose second parent is the sibling's branch tip | command seam, two-assignment fixture | a from-resolver that takes only commits refuses every sibling |
| WM7 | 7 | stdout is one `merged{worktree=,from=,kind=,previous_tip=,tip=,tree=}` record with the exact values | command seam | an incomplete record forces a `git log` to find the new tip |
| WM8 | 8 | the target operand accepts the label, the id, an 8–12 character prefix, and the absolute path, and an ambiguous prefix refuses naming both ids | command seam | a path-only operand breaks the `exec`-style address |
| WM9 | 9 | a conflicting non-capture path prints `refused{...}` and the path table on stdout at exit 1, names every conflicting path, and leaves the tip, the checkout, and the lane record unchanged | command seam | a verb that writes conflict markers leaves the worktree mid-merge |
| WM10 | 10 | a conflicted `capture/learnings.md` composes as the union and stderr carries `merge composition{resolved=capture/learnings.md:union}` | command seam, tracked-capture fixture | a verb that bypasses the table refuses a conflict the landing would settle |
| WM11 | 11 | a target with one untracked path or one modified path refuses before composition, names the path, and runs no lane | command seam | a reconcile over a dirty checkout discards the edit |
| WM12 | 12 | a sibling with a modified path refuses naming `bench commit` at the sibling, and a detached sibling refuses naming its assignment branch ref | command seam | a fold of the branch tip silently omits uncommitted delegate work |
| WM13 | 13 | a `--from` equal to an assignment label that is also a branch name refuses as ambiguous, naming the assignment id and the full commit | command seam, label-equals-branch fixture | a first-match resolver merges whichever lookup runs first |
| WM14 | 14 | a `--from` that resolves to the target's own assignment refuses | command seam | a self-merge prints `kind=current` and hides the operand error |
| WM15 | 15 | a `--from` that names nothing refuses at exit 1 naming the value | command seam | a silent `kind=current` on a typo |
| WM16 | 16 | a target whose owner marker or lock fails refuses naming that component, a non-active state refuses naming the state, and a checkout off the assignment branch refuses naming the branch ref | command seam, identity component registry | a verb that skips the bundle writes to a foreign tree |
| WM17 | 17 | a target whose branch tip is not its checkout HEAD refuses naming both commits | command seam | the compare-and-swap old value comes from the checkout while the ref differs |
| WM18 | 18 | a declared lane runs on the composed tree before the ref moves, and stdout carries `lane{outcome=pass,checks=...}` | command seam, lane manifest fixture | a verb that skips the lane publishes a broken build |
| WM19 | 19 | a failing lane check prints `lane{outcome=fail,check=<name>}` at exit 1, and the tip and the checkout stay unchanged | command seam | a verb that publishes on a lane fail |
| WM20 | 20 | a branch tip moved after the lane makes the ref update refuse, and the moved tip survives | owner seam, update-ref field | an unconditional update-ref overwrites a concurrent commit |
| WM21 | 21 | a reconcile failure after the ref update exits 3 with `merged{...,next=git -C <path> reset --merge <published-commit>}` naming the published commit | command seam, joins fault | a refusal-shaped exit hides that the ref moved |
| WM22 | 22 | a checkout edited between the pre-check and the ref update refuses before the ref moves, and the edit survives | command seam, joins hook | a verb that trusts the first read resets over a fresh edit |
| WM23 | 23 | a fast-forward from a commit whose tree fails the lane refuses, and the tip stays | command seam | a lane skipped on fast-forward publishes an ungraded tree |
| WM24 | 24 | the lane's prose placeholder resolves to exactly the Markdown paths that differ between the previous tip's tree and the composed tree | command seam, lane manifest whose prose check records its argv | an empty placeholder grades no incoming prose, and a whole-tree placeholder grades unchanged files |
| WM25 | 25 | `bench help` lists the merge row, and `bench worktree --help` prints the merge grammar | registry help test, usage test | an unadvertised verb is a dead key |
| WM27 | 27 | the landing conflict `next=` names `git -C <path> merge <destination>` as the hand resolution, states that `bench worktree merge` refuses the same conflict, and no longer claims that no verb exists | landing surface test | a stale claim after the verb exists |
| WM28 | 28 | the anchor registry needle and the fixture-bite case carry the merge-verb sentence, and the mutation that hands the merge to the delegate goes red | conformance fixture-bite test | a moved sentence with an unmoved needle reds the gate, and an unmoved sentence keeps the raw route advertised |
| WM31 | 9 | a conflicting path that holds a control byte renders escaped in the stdout refusal table at exit 1 | command seam | a raw path splits the refusal record |
| WM32 | 31 | a missing `--from`, an empty `--from`, or a second positional exits 2 with usage on stderr | usage grammar | a defaulted `--from` merges the default branch silently |
| WM33 | 9 | an add/add `capture/learnings.md` whose two sides differ in file mode refuses naming the path | command seam, add/add fixture | an edit-shaped fixture auto-resolves the mode and never reds |
| WM34 | 32 | a `--from` commit that is not an ancestor of the default branch tip and is no sibling tip refuses at exit 1 naming the commit, and no lane runs | command seam, off-branch commit fixture | a verb that lanes any commit executes an unowned tree's checks |

Not covered: story 26 — the profile paragraph is prose; the review round and
the gate's prose check grade it, and no behavior seam exists.

Not covered: story 29 — the reference list and the ADR consequence are prose;
the review round and the gate's prose check grade them.

Not covered: story 30 — the glossary entries are prose; the review round grades
them.

### Edge inventory

The walk covers the profile's hostile-input classes that reach this surface:

- A `--from` or target that holds a control byte refuses through the target
  grammar's line-safety check. A commit spelling with a control byte fails to
  peel, so WM15 names the value escaped.
- A pool path under a symbolic-link component resolves through the canonical
  path the ledger stores, as every worktree verb does.
- A cwd deeper than the repository root resolves through the Git root, as
  `reauthorize` does.
- An annotated tag as `--from` peels to its commit; the commit must still sit
  in the default branch's history, and the record names the commit.
- Ignored residue in the target (`dist/`, `node_modules/`) does not block the
  merge. The status predicate excludes ignored paths, and the fingerprint
  carries them unchanged across the run.
- A lane interrupted by a signal returns the lane's own error; the verb
  refuses, and the ref never moved.
- The assignment's recorded start stays an ancestor of every published tip,
  because the verb only adds descendants. The landing's reviewed-range check
  is unaffected.
- A tracked capture file exists only in a linked repository, because this
  repository git-ignores all three. WM10 and WM33 drive a fixture that tracks
  the file.
- A fold of a sibling already merged into the target is `current`. A sibling
  merged into the integration source reads as landed once the integration
  source lands, so `bench worktree clean --landed` retires it.

**Won't handle** lines:

- Conflict resolution inside the verb — the refusal names the paths. The hand
  resolution is raw Git in the worktree, which the platform rules leave
  outside the command boundary.
- A fold of an uncommitted sibling checkout — the sibling commits with
  `bench commit` first, so one snapshot composer exists.
- A default `--from` of the default branch — the operand is explicit, and the
  record names the exact commit.
- A `--request` proof on the merge — the verb carries the worktree commit's
  authority, and the identity bundle plus the default branch ancestry are its
  proof.
- An arbitrary `--from` commit — the decision source names the default branch
  and a sibling only, and the lane would execute an unowned tree.
- The `reauthorize` claim correction from the source — no tree text claims
  that `reauthorize` moves a base, so no edit is owed.
- A rewrite of the assignment's recorded start — every published tip descends
  from it.
- A per-direction capture rule table — one table serves the landing and the
  merge, and the disclosure names the side taken.
- A tagged system journey for the verb — the system suite proves the nested
  `worktree` dispatch for `reauthorize`, and this child routes the same way.
  The in-process command tests grade the verb.

## Ownership fences

- `internal/worktree/`
- `internal/landing/`
- `internal/gate/authorization/`
- `internal/usage/worktree.go`
- `cmd/bench/`
- `internal/anchors/registry_data.go`
- `internal/conformance/fixture_bite_test.go`
- `.agents/skills/bench-craft-delegate/SKILL.md`
- `.bench/BENCH-reference.md`
- `docs/adr/0014-main-receives-writes-only-through-landings.md`
- `projects/benchkit.md`
- `CONTEXT.md`
- `specs/worktree-merge/`

## Out of scope

- A `next=` remedy on the bare `bench preflight` `base-current` red that names
  the merge verb — 2 edits, 1 gate run.
- A `bench worktree create --from <target>` that starts a new assignment at a
  sibling's tip in one step — 3 edits, 2 gate runs.
- A fold of an uncommitted sibling snapshot through the merge verb — 4 edits,
  2 gate runs.
- An interactive conflict-resolution mode — 6 edits, 3 gate runs.
- A tagged system journey for the merge verb — 2 edits, 2 gate runs.

## Further notes

Round 1 folded eight changes. `--from` narrowed to the default branch's
history and sibling tips. One refusal renderer prints on stdout. The branch
refusal names the ref without a new identity component. The reconcile is the
bare reset, and the exit-3 `next=` is spelled exactly. Group C rose to medium,
the canary fence left, and the skill file's line budget is named.

The `--request` omission and the committed-siblings-only rule are the two
calls the reviewer is most likely to veto. The approval table names both. The
landing's conflict repair still ends in raw Git, because that conflict is one
the composition cannot settle. The verb's value is the common non-conflicting
case, which the five occurrences describe.
