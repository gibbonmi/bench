# parallel-session-landings

Status: staged

Decision source: `specs/parallel-session-landings/decisions/parallel-session-landings.md` (ready compiled map, reviewer-approved 2026-08-13)

## Problem

A reviewed spec-backed build currently has no identity separate from the destination
branch's mutable review base. Its tickets land as independently green commits,
but a later phase-owned handoff or unrelated commit becomes part of the same
contiguous `benchBase..HEAD` range. Moving that base either admits paths the
spec never authorized or hides implementation that still needs semantic review.

FT198 demonstrates the dead end. Three implementation commits are already on
`main`, its remaining work is preserved on an owned assignment branch, and a
later roadmap/handoff commit is also on `main`. Every base that exposes the
complete FT198 patch makes ownership preflight red; the only green base makes
the review diff empty. No scalar branch review base represents that authorship.

## Solution

Treat a reviewed build as a **landing source**: one owned integration branch,
one frozen base commit, and its current source tip. Tickets continue to land
serially and independently green on that branch. Review and ownership checks
grade exactly `frozen base..source tip`, independent of the current **landing
destination**.

After semantic review, `bench worktree land` composes the reviewed source tip
onto the invoking checkout's destination tip with ordinary Git merge semantics.
A conflict refuses without mutation. A clean **prospective landing tree**,
including the staged-to-implemented spec transition, receives the authoritative
whole-project gate. Publication compare-and-swaps the destination from the tip
captured for that composition, records the source tip as the second parent,
advances project-green only after publication, and then releases the owned
source worktree.

Destination movement invalidates only the prospective tree and its gate
verdict: rerunning recomposes the unchanged reviewed source and gates again.
Source movement invalidates semantic review and requires a new reviewed tip.

## User stories

1. **A reviewer can inspect the complete source without destination history.**
   From the integration worktree, an explicit frozen-base mode for `bench diff`
   and `bench preflight build|review` resolves one ancestor base and the current
   source tip. The live diff and changed-path ownership check root their
   committed, index, tracked-worktree, and untracked inventory at that base
   rather than `branch.<name>.benchBase` or merge-base with the destination.
   Review mode refuses uncommitted source state, so its non-empty predicate and
   reported `(base, source tip)` pair describe a committed review subject.
   `Line: gpt-5.6-terra / high.` The query owners already exist, but their source
   identity and drift checks are correctness-bearing and cross process state.

2. **A clean source composes like an ordinary controlled merge.** An owned
   integration branch and reviewed source tip compose onto the destination tip
   captured at landing entry. Destination-only commits survive, already-landed
   source commits are not replayed, source-only commits appear once, and a
   textual or structural merge conflict refuses before the gate without
   creating merge state in either checkout.
   `Line: gpt-5.6-terra / high.` The existing exact landing owner is the seam,
   but three-way Git composition and conflict classification need adversarial
   graph coverage.

3. **The exact prospective landing tree is the oracle subject.** The composed
   tree includes the final staged-to-implemented spec transition and receives
   the existing whole-project prospective gate. Only reusable evidence for the
   identical tree and oracle may satisfy authorization. Distinct destination
   tips cannot inherit an opaque whole-tree verdict, and no partial or
   diff-scoped result may publish the landing.
   `Line: gpt-5.6-terra / high.` The gate seam is strong, but attaching the new
   composition to it determines whether green still means shippable.

4. **Publication is fail-closed against both identities.** The landing command
   requires the reviewed source tip, verifies that the integration branch still
   names it, and publishes a commit whose first parent is the captured
   destination and whose second parent is that source tip. Source movement
   refuses as stale review. Destination movement loses the expected-old update,
   preserves the winner and the source, and returns a rerun action. Only a
   successful publication advances the destination's project-green marker.
   `Line: gpt-5.6-terra / high.` Two moving refs plus post-gate publication and
   marker order require deterministic fault seams rather than timing tests.

5. **The owned worktree closes only after the landing is durable.** The command
   authenticates the request, assignment record, owner marker, worktree path,
   branch, and reviewed tip; it validates and discloses the reviewer-authorized
   frozen base under the explicit trust assumption below. It accepts a clean
   source and a clean invoking destination checkout. A pre-publication refusal
   changes neither.
   Success reconciles the destination checkout to the published tree and
   releases the now-landed source worktree. Any post-publication reconciliation,
   marker, or release failure reports the landed commit and retains enough state
   to finish safely; it never rolls the destination backward.
   `Line: gpt-5.6-terra / high.` Existing ownership and recovery owners reduce
   the surface, but their ordering around an irreversible ref update is subtle.

6. **The workflow carries one source from implementation through review.** A
   full spec-backed build creates or retains one owned integration worktree after the
   tickets are sliced, commits each reviewer-approved ticket there in blocker
   order, reviews the explicit source range, and hands the reviewed pair to the
   landing command. Phase-owned handoff writes remain outside implementation
   fences. The already preserved FT198 doctrine worktree resumes with base
   `0924e02e`, modifies only the named `context_render.go`/`context_types.go`
   single-owner collapse while keeping the other seven retained paths byte-
   identical, and never reruns or rewrites the three commits already on `main`.
   `Line: gpt-5.6-sol / high.` This changes kit guidance that steers every later
   spec-backed build, so the leverage override applies even though the Git seam is set.

## Implementation decisions

**The source identity is explicit, not another lifecycle record.** The complete
review identity is `(owned integration branch, frozen base commit, reviewed
source tip)`. The command line and phase handoff carry the two commits; Git and
the existing assignment ledger carry the branch and worktree identity. Do not
add a commit ledger, run revision, receipt schema, prepared-operation journal,
or mutable replacement for `benchBase`.

**One integration worktree carries the serial ticket chain.** At implementation
entry, `bench worktree create` starts from the then-current destination as it
does today. Each dependent ticket receives a fresh write-delegate charge in
that same isolated worktree and lands through ordinary path-scoped
`bench commit`, advancing only the integration branch. Reuse is limited to the
dependent chain; unrelated builds and parallel writers retain separate
worktrees. The frozen review base is the reviewed ticket-graph commit on which
that worktree started, not a value inferred later from destination ancestry. The
legacy FT198 source is the explicit exception described by the trust assumption
below: its retained assignment started after three ticket commits had already
landed.

**Explicit source review extends the existing query and preflight seams.**
`bench diff [--full] --base <commit>` is mutually exclusive with `--commit` and
renders the same live snapshot as bare `bench diff`, but rooted at the explicit
base: committed `base..HEAD` changes plus current index, tracked-worktree, and
untracked state. `bench preflight build|review <slug> --base <commit>` uses the
same live changed-path inventory. It requires a real commit that is an ancestor
of source `HEAD`; review mode additionally requires a clean source checkout, so
the accepted review subject collapses to the committed `base..HEAD` range. Its
source-validity row replaces the default-branch-ancestor predicate for this
mode, while ownership and row checks keep their current meanings. Both surfaces
resolve the base and source tip once per movement-checked attempt and print the
full commit identities; they never write Git config.

**The landing porcelain belongs to the worktree family.** The public operation
is:

    bench worktree land --request <opaque-id> --base <commit> --source-tip <commit> --spec <slug> -m <message> <path>
    bench worktree land --resume <published-commit> --request <opaque-id> --base <commit> --source-tip <commit> --spec <slug> <path>
    bench worktree reauthorize --assignment <assignment-id> --request <new-opaque-id> --base <commit> --source-tip <commit> <path>

It runs from the destination checkout. `--request` plus `<path>` authenticate
the existing assignment; `--base` and `--source-tip` reproduce the reviewed
pair; `--spec` is required because this scope owns reviewed spec-backed builds. The
destination is the resolved default branch, which must be checked out in the
clean invoking checkout; its expected tip is captured once at entry. The
`--resume` form accepts only the exact commit reported by a prior
landed-but-incomplete result and performs no composition, gate, commit, or
destination update. The mutating command stays operational rather than
joining the approved AXI query set, but its success output names the source
base, source tip, destination base, published commit, tree, and terminal
worktree state. Usage exits 2; unsatisfied intent and landed-but-incomplete
results exit 1; complete success and authenticated `already-complete` results
exit 0. Progress or Git diagnostics never replace the structured terminal
result.

**The frozen base and legacy token replacement are reviewer-held authority.**
Bench has no durable semantic-
review receipt, so a commit cannot authenticate its own role as the beginning
of the approved patch. The reviewer is the trust root: before implementation,
the approved ticket-graph commit is recorded in the phase handoff and repeated
verbatim through source review and landing. Ordinarily it equals the
assignment's recorded start because the integration worktree is created after
ticket slicing. FT198 is the reviewer-visible recovery exception: its approved
review base is `0924e02e`, while retained assignment
`f70c5f8dc98f203fac19bdd6e07df1d3` records start `c46b135a` because that
worktree was created after three implementation tickets landed. The command
proves that the supplied base is a full ancestor of the source tip and exposes
both it and the recorded assignment start; it does not misstate that equality
with the start authenticates review. Substituting another ancestor changes the
review subject and requires reviewer approval.

The explicitly reviewer-authorized `bench worktree reauthorize` operation is
the trust root for replacing a lost legacy request preimage. It takes the
assignment ID, new opaque request token, reviewer-authorized base, current
source tip, and owned path shown above. It compare-and-swap updates only that
named active assignment's request digest after verifying its owner marker,
worktree path, branch, source tip, recorded start, and authorized base. The
recorded start must resolve to a full commit that is an ancestor of, or equal to,
the source tip and is disclosed; it is never required to equal the reviewer-
authorized base. The prior request digest is read under the ledger lock and used
as the expected-old value in the same operation's single compare-and-swap;
concurrent assignment movement refuses. This introduces no assignment field or
schema version. It never discovers or accepts an assignment by path alone, and
ordinary callers without explicit reviewer approval may not invoke this recovery
operation. FT198 uses it because the retained ledger preserves only the old
request digest, not its preimage.

**Preflight and landing share one source-range derivation.** Extend the current
review-range owner rather than writing a second ancestor/path parser inside the
worktree command. It resolves full commit IDs, proves base ancestry, enumerates
the committed `base..source tip` path set, and rejects source-tip movement.
`bench diff`, preflight, and final ownership authorization consume that typed
fact. Live diff/preflight consumers add checkout state to that committed range;
final landing consumes only the committed base and reviewed tip. The spec's
fence matcher remains the one owner of path authorization.

**Final landing accepts committed source state only.** The assignment branch,
its worktree `HEAD`, and `--source-tip` must name the same commit. The source
index, tracked worktree, untracked set, and nested repository state must be
clean. Ignored declared build outputs do not change the source commit and are
handled by existing release policy; any other ignored residue retains the
worktree after publication rather than being discarded. A source base that is
not an ancestor of the reviewed tip, a detached/foreign worktree, a mismatched
request or owner marker, or an already moved branch refuses before gating.

**The destination checkout is a clean publication surface.** The invoking
checkout must be attached to the resolved default branch and clean across index,
tracked worktree, untracked files, nested repositories, and undeclared ignored
residue before composition. Revalidate the same fingerprint immediately before
publication. This scope supports other committed destination landings, not an
uncommitted second writer in the destination checkout. A changed fingerprint
refuses without moving the ref; recoverable dirty-destination set-aside remains
separate work.

**Composition is an in-memory Git merge over committed identities.** Extend the
deep landing owner with a source-merge request. It derives the real Git merge
base from the destination and reviewed source tip, performs Git's ordinary
three-way tree merge without changing an index or worktree, and returns either
one conflict classification or one immutable tree. The frozen review base is
for authorship and semantic review; it is deliberately not forced to be the
merge algorithm's base. This distinction is what lets FT198 keep `0924e02e` as
its review base while Git recognizes the first three source commits already
ancestral to `main`.

**Conflict is a source repair, not an ambient merge session.** A conflict emits
a bounded refusal, leaves both refs and worktrees unchanged, and names that the
source needs a repair commit plus fresh review. Bench does not write index
stages, conflict markers, or `MERGE_HEAD`. Resolving a conflict changes the
source tip, so the prior semantic review cannot survive even when the changed
path remains within the same fence.

**The spec transition is part of the composed tree.** The destination and source
must both resolve the same staged spec bytes. The landing owner applies the
existing staged-to-implemented transition to the prospective tree before the
gate. Invalid, absent, divergent, or already-implemented spec state refuses
before gate execution. The source branch remains an inspectable staged parent;
the published merge commit and destination checkout contain the identical
implemented bytes observed by the gate.

**The existing prospective gate remains the sole oracle.** Authorization runs
on the immutable prospective landing tree through `gate.ExecuteTree`. Reuse is
legal only when the tree and derived oracle identity match exactly. A changed
destination that produces a distinct tree therefore runs the whole missing
closure again; the landing owner does not synthesize component inheritance or
accept source-branch green as publication authority.

**Publication has two immutable parents and one expected-old destination.**
After green, create a commit with the authorized tree, first parent equal to the
captured destination tip, and second parent equal to the reviewed source tip.
Before the expected-old destination update, verify the source branch and both
checkout fingerprints again. A destination compare-and-swap loss leaves the
new object unreachable and returns a recomposition-required result; it does not
merge, rebase, retry, or invalidate the source review by itself.

**Project-green accepts an absent starting marker.** Capture the destination's
project-green marker before composition. If absent, the expected prior marker
is the zero ref and successful publication creates it. If present, it must be a
full commit ancestral to the captured destination and is the expected old value
for marker advancement. A marker that moves before advancement produces a
landed-but-incomplete result; it never authorizes rollback of the published
destination.

**Post-publication steps never deny that publication happened.** Once the
destination ref advances, record project-green for the exact published commit,
reconcile the clean destination checkout, then release the owned source
worktree. The default-branch requirement makes the existing release owner's
landedness predicate authoritative. A failure reports `landed` with the commit,
captured prior marker, and incomplete step, preserves the source assignment/
recovery evidence, and never rolls back the ref or tells the caller to repeat
the first-run form.

The `--resume <published-commit>` form is the replay-safe continuation. It
proves that the named commit is reachable from the current destination, has the
reviewed source tip as its second parent, has the expected implemented spec
bytes, and corresponds to the retained assignment. It then resumes only the
observable incomplete marker, destination-reconciliation, and release steps.
If the published commit is still the destination tip, reconcile to its tree and
advance a missing or unchanged prior marker to it. If later destination commits
exist, reconcile only to the current destination tree; a current project-green
marker at the published commit or any descendant of it already satisfies the
marker step, while an absent, behind, or non-descendant marker refuses for
reviewer repair rather than being moved backward. The form accepts the source
parent's staged spec and never re-enters the first-run staged-spec, composition,
gate, commit, or destination-CAS predicates. A normal success removes the
assignment worktree only after its source tip is an ancestor of the published
default branch. Repeating a completed resume returns a definitive
`already-complete` result after the published commit and existing terminal
release receipt authenticate the same request, path, source tip, and completed
cleanup; it does not require the removed assignment or recreate any state. If
that bounded receipt has been evicted, repeat resume returns a definitive
fail-closed `missing-terminal-receipt` refusal rather than claiming completion
or entering the first-run path, and changes no state.

**Workflow prose changes as one current-state contract.** Update the platform
workflow, spec-authoring, implementation, semantic-review, final-check,
delegation, ticketing, profile, reference, field-guide, greenfield-sequence, and
help surfaces to name the integration source, explicit reviewed pair, and
`bench worktree land` handoff. The suggested fresh session begins after reviewed
ticket slicing, not between map, spec, and ticket planning. Do not duplicate
command grammar outside its usage owner; prose names the operation and intent,
while executable help owns flags. Every edited subject that the profile's
guidance-prose table budgets stays within its current limit through compensating
cuts; unbudgeted command, reference, profile, and documentation surfaces are
outside that check. This scope does not raise a budget to make new wording fit.

## Testing decisions

- **Primary seam.** Drive the real `bench worktree land` command in disposable
  Git repositories with an owned assignment worktree, a destination checkout,
  a staged spec, and a gate script that inspects its materialized tree. Assert
  refs, parents, tree bytes, project-green, both checkout states, assignment
  state, output, and exit code.
- **Composition seam.** Extend `internal/landing`'s real-Git table for converged,
  diverged, partly ancestral, textual-conflict, structural-conflict, rename,
  delete, mode, symlink, and gitlink graphs. This lower seam is justified because
  the public command must not create an ambient conflict state merely to expose
  the merge result.
- **Race seam.** Inject source movement, destination movement, marker failure,
  reconciliation failure, and release failure at their exact operation joins.
  Public two-process coverage proves one destination winner and a successful
  rerun; deterministic lower tests own every interleaving.
- **Review seam.** Extend `bench diff` and preflight contract tests with an
  explicit frozen base, a moved destination, a dirty source, an unreachable
  base, and the actual FT198 graph shape. Both consumers must cite the same
  source-range fact.
- **Gate seam.** The existing prospective authorization owner proves exact
  tree/oracle reuse. Integration coverage asserts the tree inspected by the
  gate equals the published commit tree and that project-green changes only
  after destination publication.
- **Prior art.** Reuse `internal/landing` for immutable-tree authorization and
  expected-old publication, `internal/diff` for movement-checked review ranges,
  `internal/preflight` for ownership fences and ticket rows, and
  `internal/worktree` for assignment authentication and safe release. No test
  fixture reimplements those owners.
- **Central mutation probe.** Make source-tip verification compare only tree
  equality, allowing a different source commit with identical bytes. The held
  review test must go red because semantic review binds the commit identity,
  not merely an equivalent tree.

### Seam diagrams

    trigger: review the finished integration source
        │
        ▼
    frozen base + source HEAD ──▶ [ shared source-range owner ]
                                       │
                         ┌─────────────┴─────────────┐
                         ▼                           ▼
                  bench diff --base       preflight ownership/rows
                         │                           │
                         └─────────────┬─────────────┘
                                       ▼
                              reviewed (base, tip)
                  ◀ tests attach here: moved destination changes neither range

    trigger: bench worktree land --base B --source-tip S --spec X <path>
        │
        ▼
    owned clean source S + clean destination D
        │
        ▼
    [ exact landing owner: Git three-way composition ]
        │ conflict: refuse unchanged
        ▼ clean
    prospective tree + implemented spec ──▶ [ whole-project prospective gate ]
        │                                           │ exact green
        └───────────────────────────────────────────▼
        commit(tree, parents=D,S) ──CAS D────────▶ destination
                                                  │
                                                  ▼
                                   project-green → reconcile → release
        ◀ tests attach here: exact tree/parents, ref races, terminal recovery

### Acceptance coverage map

| row | story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|---|
| PL1 | 1 | Explicit-base diff renders the complete live snapshot rooted at the frozen base, including `base..HEAD`, index, tracked-worktree, and untracked state plus full base/tip IDs, while excluding destination-only commits; accepted explicit-base diff and preflight runs leave `branch.<name>.benchBase` and every Git config byte unchanged | `bench diff --base` and preflight contract tests with diverged source/destination plus config-byte snapshots | to be observed at build time; the flag is currently a usage error and bare diff falls back to merge-base/`benchBase` | A repointed scalar review base can make ownership green only by hiding source commits; the expected committed and checkout inventory plus config snapshot catch both hidden paths and a write-through reuse of the mutable scalar |
| PL2 | 1, 6 | FT198 source preflight uses `0924e02e..source-tip`, keeps all implementation paths, and excludes the later phase-owned handoff | real retained graph fixture at preflight seam | observed 2026-08-13: destination preflight is red on phase-owned history, while the retained source preflight is red only because default `main` is not its ancestor | The exact public symptom stays red until source authorship replaces destination ancestry rather than widening FT198's fence |
| PL3 | 1 | Missing, non-commit, non-ancestor, and ambiguous frozen bases each refuse with no config or ref write; `bench diff --base <commit> --commit <commit>` refuses as usage before repository reads | diff/preflight source-range and cross-flag grammar table | to be observed at build time | A fallback to merge-base silently reviews a different patch than the caller named, while accepting both subject selectors makes the reviewed range ambiguous |
| PL4 | 1, 4 | Review mode requires a clean source and prints the exact source tip; any later branch-tip movement makes final landing refuse stale review before gate | preflight plus injected source-ref movement | to be observed at build time | Tree-equivalent or appended source commits must not inherit semantic approval for another commit identity |
| PL5 | 1 | Explicit-base query retries once on source HEAD, index, worktree, or base-resolution drift and then returns the existing structured snapshot-drift refusal with the exact retry action | `bench diff` movement test | existing bare-diff drift coverage; explicit mode is to be observed at build time | Reading the base and tip in separate moments can produce a range that never existed |
| PL6 | 2 | A clean diverged graph preserves destination-only content and applies source-only content once | landing real-Git composition table | to be observed at build time | A replay/cherry-pick approximation can duplicate already-integrated changes or drop destination work |
| PL7 | 2, 6 | When the first source commits are already ancestors of the destination, composition applies only later source commits while the merge commit still names the reviewed source tip | FT198-shaped graph at landing seam | to be observed at build time | This is the composition degenerate in the motivating case: replaying `base..tip` duplicates the first three FT198 commits |
| PL8 | 2 | Textual, modify/delete, rename/rename, file/directory, mode, symlink, and gitlink conflicts refuse before gate and leave refs, indexes, worktrees, and merge-state files unchanged | landing conflict table plus gate tally | to be observed at build time per enumerated class | A happy-path-only tree overlay can silently choose one side while claiming ordinary merge semantics |
| PL9 | 2, 4 | Repairing a conflict with a new source commit invalidates the old reviewed tip; after a new review the same destination composes cleanly | public source-repair journey | to be observed at build time | Refusal alone can strand progress; accepting the old review would instead bless unreviewed resolution bytes |
| PL10 | 3 | The gate inspects the exact prospective tree and the published commit has the identical tree | public landing fixture with gate-side tree assertions | to be observed at build time | Gating the source tree or destination checkout can stay green while their merge is wrong |
| PL11 | 3 | The gate and published tree both contain identical `Status: implemented` bytes while both parents retain the staged spec | public spec landing fixture | to be observed at build time | A post-gate transition recreates the known tree mismatch; mutating the source destroys the reviewed parent |
| PL12 | 3 | `authorization.Green` alone permits publication; `Candidate`, `Inherited`, and `Infrastructure` each leave destination, project-green, both checkouts, and assignment unchanged | four-kind authorization table through landing command | existing producer classification is covered; source-landing state assertions are to be observed at build time for every kind | Enumerating the producer's exact partition prevents one red attribution or an infrastructure outcome from being treated as green |
| PL13 | 3 | Exact tree-and-oracle evidence may be reused for the identical composition, but a destination change producing another tree executes the missing whole-project closure | gate tally across identical and changed compositions | existing exact evidence-key tests plus new integration tally | Reusing source green or the previous destination's opaque verdict weakens the final oracle |
| PL14 | 4 | Published commit has first parent equal to captured destination, second parent equal to reviewed source tip, and no additional parent | real Git object assertion | to be observed at build time | A squash loses source ancestry; a normal ambient merge may select another parent or tree |
| PL15 | 4 | Destination movement during the gate loses the expected-old update, preserves the winner and source review, and does not advance project-green | injected ref updater plus two-process winner/loser journey | to be observed at build time | An unconditional update overwrites the winner; automatic recomposition publishes a tree that did not receive the completed verdict |
| PL16 | 4 | Rerunning after destination-only movement recomposes and regates the unchanged reviewed source, then publishes without a second semantic review | two-process public journey | to be observed at build time | The opposite degenerates either retain stale gate evidence or erase valid source review on every unrelated landing |
| PL17 | 4 | Project-green advances to the published commit only after destination CAS; conflict, gate red, source drift, and CAS loss leave the prior marker unchanged | marker fault/order table | to be observed at build time | Advancing before publication lets a non-destination tree authorize later reuse or release work prematurely |
| PL18 | 5 | Wrong request, mismatched path, foreign or detached branch, invalid owner marker, moved assignment record, source-tip mismatch, a base that is not a full ancestor, and any committed source path outside the staged spec's ownership fences refuse before gate; a valid reviewer-supplied base and the assignment start are both disclosed even when FT198 makes them differ | worktree-land authentication, source-fence, and trust-assumption table | existing worktree ownership refusal fixtures plus new landing operation | A path alone is not proof that the caller owns the branch, while skipping final path authorization permits publication outside reviewed scope and pretending the assignment start authenticates FT198's review base rejects the motivating recovery |
| PL19 | 5 | Dirty source state and dirty destination state each refuse before gate; a destination fingerprint change during the gate refuses before CAS | public state matrix plus injected fingerprint mutation | to be observed at build time | Uncommitted bytes are absent from the reviewed Git identities and must neither disappear nor leak into publication |
| PL20 | 5 | Complete success and authenticated `already-complete` exit 0; unsatisfied intent and `landed` with an incomplete step exit 1; usage exits 2. Complete success updates the clean destination checkout to the published tree, advances project-green, releases the landed assignment with no recovery ref, and emits one structured result naming source base, source tip, destination base, published commit, tree, and terminal worktree state; progress or Git diagnostics never replace that terminal result | public success journey, complete exit-partition table, and `bench worktree list` | to be observed at build time | A ref-only implementation leaves the destination checkout reverse-dirty or leaks the active assignment forever, while a fail-open incomplete exit or ungraded terminal envelope can send callers down the wrong recovery path or hide the immutable subject |
| PL21 | 5 | Marker, checkout reconciliation, or release failure after CAS reports `landed` with commit and incomplete step, keeps recoverable state, and never rolls back or republishes | ordered fault injection after publication | to be observed at build time | Reporting generic failure invites a duplicate first-run invocation; rollback can overwrite a later destination winner |
| PL22 | 5 | Control-bearing source paths, hostile worktree paths, and unsafe Git diagnostics produce bounded structured refusal without forged lines; special files are rejected before read | public hostile-input table | existing TOON/path refusal tests cover the old queries; new source and landing sinks are to be observed at build time | Git-valid hostile names can corrupt the terminal envelope or block a reader even when merge logic is correct |
| PL23 | 6 | Platform workflow, spec-authoring, implementation, semantic-review, final-check, delegation, ticketing, profile, reference, field-guide, greenfield-sequence, and help surfaces all route one integration source through explicit review and worktree landing, with the spec-authoring fresh-session recommendation moved after ticket slicing; the platform's current claim that path-scoped `bench commit` is the sole landing path is removed | workflow conformance and stale-reference sweep over the enumerated subjects | to be observed at build time; current platform rules call path-scoped `bench commit` the sole landing path, spec authoring recommends a fresh session before ticket slicing, guidance lands tickets directly on the destination, and review resolves from `benchBase` | The code can work while canonical platform rules or fresh-agent guidance continue invoking or advertising the broken scalar-base flow |
| PL24 | 6 | The recovery journey starts from the preserved FT198 assignment, modifies only the named `internal/roadmap/context_render.go`/`context_types.go` single-owner collapse while keeping the other seven retained paths byte-identical, reviews from `0924e02e`, and lands without rewriting `main` history or rerunning its first three tickets | bounded FT198 dogfood journey after generic tests are green | current public preflight stops red before that work can resume, and the retained handoff requires the one named knowledge-deduplication repair | A synthetic fixture may miss the exact interleaving, bounded retained repair, and release semantics this feature exists to recover |
| PL25 | 5 | `--resume <published-commit>` accepts the source parent's staged spec, proves the published commit/source/destination relationship, skips gate and publication, and completes each marker, reconciliation, or release failure without another commit; repeating it after completed release returns `already-complete` from the published commit and terminal receipt, while an evicted terminal receipt returns a fail-closed `missing-terminal-receipt` refusal without changing state or entering the first-run path | one replay case after each injected post-CAS failure, a repeated completed resume with its receipt present, and the same repeat after FIFO receipt eviction | to be observed at build time | Without a distinct and idempotent resume predicate, the implemented destination spec or removed assignment makes the promised recovery reject; treating an evicted receipt as completion invents authority, while retrying first-run can republish |
| PL26 | 5 | Landing from a checkout not attached to the resolved default branch refuses before gate, while a successful default-branch landing satisfies the existing release landedness predicate | public destination-branch table plus worktree list | to be observed at build time | Existing release classifies landedness only against the default branch; accepting another destination strands a supposedly closed source worktree |
| PL27 | 5 | Each flag value is parsed only as that flag's value: land and reauthorize invocations where a would-be `<path>` appears only as `--request`, `--base`, `--source-tip`, `--spec`, `-m`, `--resume`, or `--assignment` data refuse for missing positional input | public grammar table with one flag-value-only case per value flag on each applicable invocation and `--` path separator controls | to be observed at build time | A first-non-dash parser can consume a flag value as the source path and mutate the wrong worktree or ref |
| PL28 | 4 | An absent project-green marker is created at the published commit; a present ancestral marker advances from its captured value; concurrent marker movement yields landed-but-incomplete and resumes without republishing | marker absent/present/moved table | current producer has an absent-marker branch but no landing caller; to be observed at build time | Assuming a prior marker makes the first landing fail after publication, while an unconditional marker write overwrites concurrent green evidence |
| PL29 | 5, 6 | Explicitly reviewer-authorized `bench worktree reauthorize` replaces FT198's lost request preimage with a new token only after exact assignment-ID, owner, path, branch, tip, recorded-start, and approved-base verification; the recorded start is a full commit that is an ancestor of, or equal to, the source tip and is disclosed but need not equal the approved base; one compare-and-swap uses the prior request digest read under the ledger lock in the same operation, without a new assignment field or schema version; the new token then authenticates the unchanged retained worktree | re-authorization grammar/CAS table plus real FT198 assignment dogfood | current ledger exposes only the old request digest, has no CAS update, and no public operation can supply its preimage | Requiring an unrecoverable secret makes PL24 impossible, while path-only discovery or a blind ledger put lets another caller seize or overwrite an assignment; strict start ancestry or start/base equality rejects the motivating retained graph |
| PL30 | 5 | Resume after later destination commits reconciles to the current destination tree, treats a descendant project-green marker as satisfied, and neither moves the marker nor destination backward; absent, behind, or divergent marker state refuses without changing either | resume-after-destination-movement table | to be observed at build time | Reachability alone otherwise authorizes reconciliation to an obsolete tree or marker regression |
| PL31 | 5 | A source with only declared build-output residue releases normally, while any undeclared ignored residue publishes successfully but retains the assignment and names the residue/release action without deleting it | public ignored-residue terminal-state table | existing release policy classifies ignored residue; source landing integration is to be observed at build time | A single success assertion can either leak ordinary worktrees or discard ignored bytes it never reviewed |
| PL32 | 6 | Each edited subject named by the profile's guidance-prose budget table remains at or below that table's current limit through compensating cuts; unbudgeted surfaces are not counted by this row | guidance-prose budget conformance over the enumerated table subjects | current implementation, ticketing, and delegation guidance are already exactly at their limits | Correct workflow prose that simply appends lines to a budgeted subject makes the gate red |
| PL33 | 3 | Invalid, absent, source/destination-divergent, and already-implemented spec state each refuses before gate execution with destination, project-green marker, source and destination checkouts, and assignment unchanged | public spec-state refusal table with gate tally and full state assertions | existing `spec.CheckStaged` validates one checkout only; source/destination byte comparison and landing state preservation are to be observed at build time | A green-only transition row cannot catch landing the wrong spec bytes, gating an invalid transition, or mutating state before the precondition is known |

The cheapest composition degenerate replays every commit after the frozen review
base onto the destination. PL6 may pass on a fully disjoint graph, but PL7 goes
red because commits already ancestral to the destination are applied twice. The
cheapest authority degenerate gates the source tip rather than the composed
tree; PL10 goes red because destination-only bytes visible in the published
tree never appear in the gate checkout.

### Edge inventory

- **Error path** — PL3-PL5, PL8-PL9, PL12, PL15, PL17-PL19, PL21, PL22, and
  PL33 cover
  invalid identities, drift, conflicts, gate outcomes, publication loss,
  authentication, dirty state, partial completion, and hostile diagnostics.
- **Empty or absent input** — missing flags/path remain usage errors; PL27
  separates absent positional input from flag values; PL33 covers absent spec
  state before the gate; an empty source diff makes review preflight red; PL28
  covers an absent project-green marker. A valid first-run landing always changes
  the tree through the required staged-to-implemented spec transition.
- **Boundary values** — one source commit, multiple commits, a source whose early
  commits are destination ancestors, a destination ancestor of source, source
  ancestor of destination, and a true divergence are enumerated across PL6 and
  PL7. A root commit cannot be a frozen base without being an ancestor.
- **Malformed input** — abbreviated/ambiguous/non-commit IDs, invalid refs,
  wrong request IDs, paths outside the repository, symlinked owner records,
  special files, control-bearing values, and flag values mistaken for the
  positional path are covered by PL3, PL18, PL22, and PL27.
- **Interrupted or partial state** — before CAS, interruption leaves refs and
  checkouts unchanged apart from unreachable Git objects. After CAS, the commit
  is authoritative; PL21 preserves and reports incomplete marker/reconcile/
  release work instead of denying or rolling back publication.
- **Re-run idempotency** — an unchanged pre-publication refusal is repeatable;
  a CAS loser follows PL16; PL25's explicit resume form finishes a
  landed-but-incomplete result without creating another merge commit, refuses
  fail-closed when its bounded terminal receipt has been evicted, and PL30 covers
  that continuation after later destination movement.
- **Process-boundary lifecycle** — PL15, PL16, and PL24 use fresh processes and
  re-read both refs, assignment state, gate evidence, and project-green rather
  than retaining in-memory identity.
- **Hostile environment** — held Git index/ref locks, gate lock contention,
  missing Git, cwd below repository root, unreadable common-dir state, and a
  destination checkout changed during the gate fail closed at their current
  owners. Ref-lock loss is classified as CAS refusal only when the destination
  actually moved; other lock failures remain infrastructure errors.
- **Prose-budget headroom** — PL32 requires every guidance addition to be paid
  for by a compensating cut; raising the budget is not an implementation escape.
- **Ignored residue** — PL31 distinguishes declared disposable build output
  from undeclared ignored bytes, which retain the assignment after publication
  and are never silently discarded.
- **Dirty destination set-aside** — **Won't handle:** this scope requires the
  destination checkout clean and supports concurrency through committed
  destination movement. Preserving an uncommitted destination writer needs
  FT98's recoverable set-aside and retains a surviving in-scope caller in
  ordinary path-scoped `bench commit`.
- **Non-spec and multi-source landing** — **Won't handle:** the command requires
  one staged spec and one owned source. Light-path landing and octopus/multiple
  sources need separate authority and have surviving current manual/serial
  routes.

## Ownership fences

- `internal/diff/`
- `internal/preflight/`
- `internal/landing/`
- `internal/worktree/`
- `internal/intent/`
- `internal/gate/authorization/`
- `internal/axi/`
- `internal/usage/`
- `cmd/bench/main.go`
- `cmd/bench/command_registry_test.go`
- `bin/bench.sh`
- `internal/systemtest/`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/docs_workflow_helpers_test.go`
- `internal/conformance/package_shipped_surface_test.go`
- `internal/anchors/registry_data.go`
- `tests/canary/`
- `.bench/BENCH.md`
- `.bench/BENCH-reference.md`
- `.agents/commands/bench-implement-spec.md`
- `.agents/commands/bench-review-implementation.md`
- `.agents/commands/bench-final-check.md`
- `.agents/commands/bench-write-spec.md`
- `.agents/skills/bench-craft-delegate/SKILL.md`
- `.agents/skills/bench-craft-tickets/SKILL.md`
- `projects/benchkit.md`
- `docs/field-guide.html`
- `docs/greenfield-build-sequence.md`
- `CHANGELOG.md`
- `CONTEXT.md`
- `specs/parallel-session-landings/`

`capture/session-handoff.md` is deliberately not an implementation fence. The
phase owns that ambient handoff write; the landing source does not absorb it to
make ownership preflight green.

The `restore-ft198-doctrine-commit.md` and `finish-and-land-ft198-source.md`
tickets execute in the retained FT198 source under
`specs/roadmap-progressive-index/spec.md` and its ownership fences, not under the
list above. Their delegates run
`bench preflight build roadmap-progressive-index --base 0924e02e`; those recovery
paths do not widen parallel-session-landings ownership.

## Out of scope

- **A new multi-coordinator spec-build lifecycle.** Run records, revisions,
  receipt schemas, planning subjects, checkpoint state, status queries, and
  recovery journals are rejected by decision #6, not deferred. Reintroducing
  them would be a new reviewed capability (`30+ edits, 7+ gate runs`).
- **Dirty-destination set-aside.** Capturing, clearing, and restoring another
  writer's uncommitted destination state belongs to FT98's recoverable payload
  owner (`12 edits, 3 gate runs`). This build requires a clean destination and
  supports committed destination movement.
- **Generic no-spec/light-path worktree landing.** Without a staged spec there
  is no approved ownership-fence or coverage-row subject to authorize the
  source. A separate capability must choose that authority (`8 edits, 3 gate
  runs`).
- **A non-default landing destination.** Existing assignment release proves
  landedness against the resolved default branch, and reviewed spec-backed builds land
  there. Destination-relative release for arbitrary branches is separate
  lifecycle work (`6 edits, 2 gate runs`).
- **Parallel ticket writers inside one source.** This workflow keeps tickets
  serial and independently green on one integration branch. Multiple frontier
  branches plus an integration scheduler are separate concurrency semantics
  (`15+ edits, 4+ gate runs`).
- **Automatic conflict resolution, source rebasing, or history rewriting.** A
  conflict becomes an explicit source repair and new review. Automated repair
  policy is a separate reviewer decision (`8 edits, 3 gate runs`).
- **New component/check inheritance policy.** The gate may use only its existing
  authoritative exact-tree/oracle evidence rules; this spec does not add
  diff-scoped or cross-tree green transfer (`10 edits, 3 gate runs`).
- **Weakening `bench prep-release` or changing release publication.** The ship
  tier retains its current exact-green prerequisites (`0 edits, 0 gate runs`).
- **Rewriting or deleting FT198 or roadmap-maintenance commits.** PL24 consumes
  the current graph and preserved worktree as-is (`0 edits, 0 gate runs`).
