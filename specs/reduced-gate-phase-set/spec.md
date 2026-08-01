# Reduced gate phase set for a declared path allowlist

Status: implemented

Decision source: reviewer-confirmed current conversation, 2026-08-01. Three rulings. First, a changeset confined to an exact declared path allowlist runs a reduced phase set, with the allowlist and the phase set both single-sourced and gate-asserted, rather than reusing a whole-tree green verdict across the subject change. Second, enforcement stays by construction — the excludable phases run against a tree the allowlisted paths are absent from — accepting that every full gate materializes a stripped worktree. Third, the capture surfaces co-locate under `capture/` — landed as a separate light-path migration before this spec is built — so the allowlist the profile pins expands to that one directory; the profile names expansion a new decision, and this is that decision.

## Problem

Every commit pays the whole gate, including commits that cannot affect most of
what the gate grades. The 2026-08-01 roadmap drain touched exactly two files —
`ROADMAP.md` and the learnings journal — and paid a full run: measured that day,
the contract and artifact phases alone took about 574 seconds against 113 for the
conformance suite. A documentation-only changeset waits roughly ten minutes for
evidence about Go behavior it did not touch, and the drain phase pays it again
whenever a red or a concurrent writer forces a second attempt.

The existing machinery does not reach this case. `bench commit` already reuses a
fresh green through `ExecuteReusingFreshGreen`, but reuse requires a matching
subject and any content change computes a new one, so a capture-only edit always
falls through to a real run. `internal/status` already holds a capture-only path
set, but it governs only the gate-staleness signal on the ambient board; it has
never scoped an oracle run.

The naive fix is unsound and the tree proves it. FT166 proposed exempting the
capture set from the gate outright, conditioned on "no gate check reads those
paths" — and three of the four paths it names are graded from the real root
today: `checkOccurrenceLedgerMigration` reads `ROADMAP.md`,
`TestHandoffShapeSingleSourced` reads the handoff from the graded root,
and the conformance doc sweep includes `capture/learnings.md`. An exemption would
stop running exactly the checks that exist to grade capture commits.

The second trap is subtler and killed this spec's first draft. Making a file
absent does not reliably make a dependent check fail: the kit's own idiom for a
missing subject file is `skipIfSubjectFileMissing`, which converts absence into a
capability skip, and the profile makes capability skips informational in the dev
tier unless capabilities are explicitly required. A design that hides the
allowlisted paths from a phase and trusts a red to appear would, for any check
written in that idiom, produce a permanent silent green instead — worse than the
exemption it replaced, because it degrades every run rather than only the
capture-confined ones.

## Solution

The gate learns one declaration holding two lists: the allowlisted paths, and the
phases that are excludable because they cannot observe them. A changeset confined
to the allowlist runs only the phases that can see it, and records a verdict
marked as reduced that names the full-green run whose evidence it inherited for
everything else.

Soundness comes from construction plus a posture, not from an assertion. On a
full run, the excludable phases execute against a materialized worktree of the
current tree with the allowlisted paths removed — and they run there with
capabilities *required*, so a check that quietly degrades to a capability skip
because its subject file vanished is a red rather than an informational line. A
phase that depends on an allowlisted path therefore cannot pass in the dark by
either route: it fails loudly, on an ordinary full gate, at the moment the
dependency is introduced, and the fix is to move that phase into the included
set.

Reuse then keys on a second identity computed over the tree with those same paths
excluded. When a capture-only commit leaves that stripped identity unchanged, the
excludable phases' most recent full green still answers for them exactly, and only
the included phases run. Consecutive capture commits all inherit from that same
full green rather than from each other, so a run of drain commits stays cheap
without evidence quietly ageing through a chain.

The operator sees the reduction. A reduced run announces which phases it skipped
and which full-green run supplied their evidence, in the same spirit as the
existing reused-verdict announcement, because a skipped run that says nothing
reads as a gate that never ran.

## User stories

1. As a maintainer, I want one declaration to hold the path allowlist and the
   excludable phase set, with the ambient staleness signal, the commit path, and
   the gate all routing through it, so that the set of files on the fast path
   cannot differ between the board's advice and the oracle's behavior. Line:
   `opus` / medium. Consolidating an existing map across three consumers is a
   known seam, and the risk is import-direction rather than logic.

2. As a maintainer, I want the profile's pinned allowlist prose updated in the
   same change that widens the list, so that the document naming expansion "a new
   decision" records the decision that was made rather than contradicting the
   code. Line: `fable` / high. This is kit guidance prose, which `craft-line`'s
   leverage override routes to the top tier.

3. As a maintainer, I want the gate to compute a stripped subject identity that
   excludes the allowlisted paths alongside its existing whole-tree identity, so
   that a capture-only change is recognizable as one that leaves every excludable
   phase's observable input untouched. Line: `opus` / high. This is new identity
   arithmetic on the oracle's reuse path, where a permissive mistake reuses
   evidence that no longer answers for the tree.

4. As a maintainer, I want every excludable phase on a full run to execute against
   a materialized worktree with the allowlisted paths absent and capabilities
   required, so that an undeclared dependency — whether it hard-fails or
   soft-skips — becomes a red on an ordinary gate instead of a silent hole in
   every later run. Line: `fable` / high. This story converts the reviewer's
   "gate-asserted" constraint from a claim into an enforced property, changes what
   a full gate grades, and is where the first draft was wrong; the top-tier
   escalation is reviewer-approved for this story.

5. As an agent landing a capture-only commit, I want the gate to run just the
   included phases and record a verdict marked reduced that names the full-green
   run its skipped evidence came from, so that a later reader can tell exactly
   what was and was not graded. Line: `opus` / high. Verdict shape is what every
   downstream consumer trusts, and an ambiguous record is worse than a slow gate.

6. As an agent landing several capture commits in a row, I want each to inherit
   from the same full-green ancestor rather than from its reduced predecessor, so
   that a drain that lands twice stays cheap without excludable evidence ageing
   through a chain of records. Line: `opus` / high. This is the case the motivating
   scenario actually hits, and the naive single-slot implementation silently loses
   the ancestor on the first reduced write.

7. As a reviewer, I want a reduced verdict to satisfy nothing that requires a
   whole-tree green — `bench prep-release` and any release-evidence precondition
   refuse it with that reason — so that the fast path can never leak into release
   authority. Line: `opus` / medium. The refusal is a fail-closed default at an
   existing precondition seam.

8. As an agent, I want `bench commit` to take the reduced path automatically when
   its staged set is allowlist-confined, and to say which phases it skipped and
   why, so that the saving needs no new flag and no session has to decide when the
   fast path applies. Line: `opus` / medium. The staged-set intersection is
   mechanical once story 1 lands, and the announcement follows the existing
   reused-verdict precedent.

9. As a maintainer, I want a conformance check binding the allowlist, the
   excludable phase set, and the prose that documents them to their single source,
   so that a future edit cannot widen the fast path in one place only. Line:
   `opus` / high. Correctness of the oracle's own scoping matters more than speed,
   per the profile's gate-logic routing.

## Implementation decisions

**One declaration owns both lists.** The path allowlist and the excludable phase
set live together, because they are two halves of a single claim — *these paths
are invisible to these phases* — and splitting them lets one half move without
the other. `internal/status`'s existing `captureOnlyStalePaths` is the seed; the
declaration moves to a package both `internal/status` and `internal/gate` may
import without a cycle, and each consumer routes through it rather than keeping a
private copy beside it.

**The allowlist, enumerated.** The `capture/` and `specs/` directories,
`ROADMAP.md`, and `.bench-notes.md`. The co-location migration collapsed the inbox,
the journal, the handoff, and the retros into `capture/`, so what was six scattered
entries is now two directories plus two files: `ROADMAP.md`, which stays at the root
because it is a working document a reader opens directly, and `.bench-notes.md`,
which is not repository capture at all but per-worktree shift scratch, carried here
only because the staleness signal already carries it.

`specs/` joins on the reviewer's 2026-08-01 ruling, by the same argument: it is
entirely formatted documents — specs, decision maps, tickets, retros' siblings — and
the phases that grade them are the included ones. Conformance owns acceptance-coverage
map validation, decision-map integrity, and the docs sweep, so a spec edit is still
graded on a reduced run; what it stops paying for is evidence about Go behavior it
did not touch. Nothing about the excludable set is asserted here beyond what story 4
continuously tests: if some excludable phase does read the real `specs/` rather than a
fixture tree, the stripped construction reds the next full gate and names it, and that
phase moves to the included set.

**Membership is location, and the design polices mis-filing itself.** The
profile's pin excluded directory matching because a bare prefix delegates future
membership to whoever drops a file inside — the list stops being a decision and
becomes a pattern. Co-location answers that differently from an admission
condition: `capture/` *is* the definition of the capture surface, and `specs/` *is*
the definition of the spec surface, so a new file in either is that surface by
construction rather than by anyone's judgment. The residual risk is a file mis-filed
into a declared directory that the gate really does grade, and the
design absorbs it in both directions — a grader in the included set keeps running
normally, and a grader in the excludable set reds under story 4's stripped
construction. Mis-filing is therefore a loud failure or a harmless one, never a
silent hole, which is what makes location-based membership stronger here than the
enumerated list it replaces.

**The excludable phase set, enumerated:** `gofmt`, `vet`, `test`, `race`,
`contract`, `shellcheck`, `canary`. The included set is `conformance` and
`conformance-suite`, which are the phases that grade the allowlisted files. The
`build` phase is unconditional in both modes, because it produces the binary the
other phases exec and its output is content-derived from sources the allowlist
never touches. This enumeration is a declaration that story 4's construction
continuously tests: a phase declared excludable that in fact reads an allowlisted
path reds the next full gate and moves to the included set.

**Excludability is proven by construction and by posture.** The stripped worktree
is materialized through git so that it is a real repository — a contract test
stages the subject with `git ls-files` and fails outright against a non-repo — and
the build phase runs inside it so its `dist/` output is produced there exactly as
on the primary root. The excludable phases then run with capabilities required, so
a check whose subject file vanished cannot report an informational skip. This is a
new mechanism: the gate owns no working-tree materializer today.
`.bench/gate-prospective.sh` is a build-then-exec hook and `internal/gate/prospective.go`
only hashes it, so there is nothing existing to compose here and the spec adds one
rather than claiming reuse.

**The contract phase has two roots, and construction only reaches one.** That
phase runs the suite from the kit checkout against a separate subject root, so
stripping the subject leaves any test that resolves paths from the kit checkout
reading the real tree. Materializing a stripped kit as well would mean compiling
the suite from a copy, which is a larger change than the saving justifies. The
narrower closure is a static check over that package: a contract test may not
read an allowlisted path relative to the kit checkout, only relative to the
subject root. This is the one place the design is an assertion rather than a
construction, and it is scoped to a single package and a single resolution
helper so the assertion stays checkable.

**Two identities, one subject.** The existing whole-tree identity continues to
govern included phases and the overall subject record. A second, stripped identity
governs whether excludable-phase evidence carries forward. `TreeHash` already
builds a throwaway index, so the stripped variant is the same construction with
the allowlisted paths dropped from that index; `buildSubjectForPolicy`'s policy
parameter distinguishes the two hash domains but does not by itself vary the tree
contribution, which the stripped construction supplies.

**A reduced verdict is a distinct record class, and the loader learns its shape.**
The verdict cache is a single slot with strict field validation, so a reduced
record is a new shape rather than the existing one — the loader gains the marker,
the phases actually run, and the ancestor's identity and recorded time, and
rejects a record mixing the two shapes. Because the slot is single, the reduced
record carries the ancestor's identity forward rather than pointing at a record
that no longer exists: the next capture commit validates against that same
full-green ancestor, never against its reduced predecessor. Inherited evidence
does not re-stamp its recorded time; the existing freshness window applies to the
ancestor's, so a stale ancestor falls back to a full run rather than being
refreshed by inheritance.

**No new flag.** The reduced path is selected by the changeset, not by the
operator, because an opt-in fast path is one an agent under time pressure will
reach for in cases the rule does not cover. `bench gate --fresh` remains the
escape to a real whole-tree run.

**This does not reopen diff-scoped gating.** FT91 ruled diff-scoped gating unsound
here because contract and canary are behavior contracts with no file-to-test map,
and that ruling stands. The carve-out is exactly the one FT101 names as
legitimate: a boundary the reviewer declared, not a boundary inferred from a diff.
Speed is the motivation but never the justification — the justification is that a
phase which provably cannot observe a change has nothing to say about it.

## Testing decisions

A good test here exercises the oracle's own scoping through behavior rather than
through the declaration, because the failure this feature can cause is a verdict
that credits work nobody graded. The strongest evidence is adversarial: plant
exactly the dependency the design forbids and require the gate to notice, in both
the hard-failing and the soft-skipping form.

Seams receiving tests: the gate's execution and identity units in `internal/gate`,
driven through the existing `gateEngine` interface, whose test double in
`verdict_reuse_test.go` is the prior art for running the real decision logic
without a real gate; the runtime contract surface for command-level behavior over
a fixture root that declares its own phase table through `.bench/phases.json`, so
there are real phases to include and exclude; the preprelease contract surface for
the release refusal, which is where those contracts live; and the conformance
surface for the single-source binding, following `checkLineBinding`, which
cross-checks a prose table against its machine-readable source and reds on drift.

The gate seam that observes the feature is the conformance phase for the binding,
the contract phase for command behavior, and — for story 4 — every full gate run
in this repository, since the stripped construction is exercised continuously by
the gate's own operation.

### Seam diagram

    trigger: bench commit with an allowlist-confined staged set (or bench gate)
        │
        ▼
    staged paths  ──▶  [ one declaration: paths + excludable phases ]  ──▶ confined? yes/no
        │                  ◀ tests attach here: unit table drives membership;
        │                    conformance binds declaration to profile prose
        ▼
    tree + prior verdict ──▶ [ stripped subject identity ]  ──▶ ancestor reusable?
        │                  ◀ tests attach here: unit builds both identities over a
        │                    fixture tree; stable across allowlisted edits only
        ▼
    included phases only ──▶ [ gate run via gateEngine ]  ──▶ verdict{reduced, phases, ancestor}
                             ◀ tests attach here: engine double asserts which phases
                               ran; runtime contract reads the verdict record

    trigger: any full gate run
        │
        ▼
    materialized stripped worktree ──▶ [ excludable phases, capabilities required ]
                             ◀ tests attach here: contract plants a hard-reading and
                               a soft-skipping phase; both must red

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | Every consumer resolves the allowlist from the declaration rather than a private copy | unit (`internal/status`) | `go test ./internal/status -run TestStaleSofteningRoutesThroughDeclaration` — observed red: a private literal left beside the declaration keeps softening a path the declaration no longer carries | Pins behavior rather than equality of two exported values, so exporting the declaration for the test while leaving the behavior path on its own literal still fails |
| 1 | Every co-located capture surface is a member through the directory, not through its own entry | unit (`internal/gate`) | `go test ./internal/gate -run TestCaptureDirectoryCoversItsSurfaces` — observed red: a declaration listing the four files individually passes membership but fails the assertion that a newly added `capture/` file is covered | The migration's whole point is that membership follows location; re-enumerating the members inside the directory would reintroduce the exhaustiveness burden it removed |
| 1 | A path outside the declaration is never confined, including a near-miss sibling | unit (`internal/gate`) | `go test ./internal/gate -run TestConfinementRejectsUndeclaredPath` — observed red: prefix matching that admits `ROADMAP.md.bak` passes without the exact-match rule | The cheapest wrong implementation is a substring or prefix test; the near-miss fixture is what distinguishes it from exact-path membership |
| 1 | A file directly under a declared directory is a member; a nested or sibling-prefixed path is not | unit (`internal/gate`) | `go test ./internal/gate -run TestDeclaredDirectoryMembership` — observed red: a `strings.HasPrefix` implementation admits `capture-old/x.md` | Directory membership is the one non-exact rule in the declaration, so its boundary is where an over-broad match would hide |
| 2 | The profile's pinned allowlist prose matches the declaration | conformance | `BENCH_CONFORMANCE_ROOT="$PWD" go test ./internal/conformance -run TestRootConformance` — observed red: widening the declaration without the profile edit produces the binding diagnostic | The profile calls expansion a new decision; a code-only widening would leave the document asserting a list the oracle no longer uses |
| 3 | The stripped identity is unchanged by an edit confined to allowlisted paths | unit (`internal/gate`) | `go test ./internal/gate -run TestStrippedIdentityIgnoresAllowlistedEdit` — observed red: an identity hashing the whole tree moves and fails the stability assertion | This is the arithmetic the entire reuse decision rests on; if it moves, no reduction ever happens and the feature is inert |
| 3 | The stripped identity moves for any edit outside the allowlist | unit (`internal/gate`) | `go test ./internal/gate -run TestStrippedIdentityMovesOnUnlistedEdit` — observed red: an over-broad strip excluding a parent directory holds the identity stable and fails | The dangerous failure is an identity too *insensitive*, which reuses evidence for a tree that really changed; this is the degenerate-implementation check for story 3 |
| 4 | An excludable phase that hard-reads an allowlisted path fails a full gate | contract (runtime) | `go test ./internal/contract/runtime -run TestExcludablePhaseCannotReadAllowlistedPath` — observed red: without the stripped construction the planted phase reads the file and passes | The enforcement itself: plants exactly the dependency the design forbids and requires the gate to notice |
| 4 | An excludable phase that soft-skips on the missing path fails a full gate | contract (runtime) | `go test ./internal/contract/runtime -run TestExcludablePhaseCannotCapabilitySkip` — observed red: without capabilities required the planted phase reports an informational skip and the gate stays green | The failure that killed the first draft: absence degrading to a green skip is the kit's dominant idiom, and construction alone does not catch it |
| 4 | Included phases still see the allowlisted paths on the same run | contract (runtime) | `go test ./internal/contract/runtime -run TestIncludedPhasesSeeAllowlistedPaths` — observed red: stripping for every phase makes the conformance-shaped probe fail to find `ROADMAP.md` | Over-stripping breaks the checks that must keep grading these files — the regression that makes the design unsound in the other direction |
| 4 | The materialized stripped tree is a git repository | contract (runtime) | `go test ./internal/contract/runtime -run TestStrippedSubjectIsGitRepository` — observed red: a plain directory copy makes the `git ls-files` staging path fail outright | A tracked-file copy is the obvious implementation and it breaks a contract test that stages the subject through git |
| 4 | No contract test reads an allowlisted path relative to the kit checkout | conformance | `go test ./internal/conformance -run TestContractTestsReadCapturePathsFromSubjectOnly` — observed red: a planted kit-relative read of `ROADMAP.md` in the contract package passes today because nothing looks | Construction strips only the subject root, so this is the one hole it cannot close; without the check the phase stays observably coupled to a capture path |
| 4 | The excludable set is not degenerate: `contract` is in it | unit (`internal/gate`) | `go test ./internal/gate -run TestExcludableSetCoversContractPhase` — observed red: a declaration excluding only `shellcheck` fails the assertion | Without this the whole feature is satisfiable by excluding one trivial phase, passing every other row while saving nothing |
| 5 | A capture-only changeset runs only included phases | unit (`internal/gate`, engine double) | `go test ./internal/gate -run TestReducedRunExecutesIncludedPhasesOnly` — observed red: today every change runs all phases, so the recorded phase list is complete | Asserting the executed phase list through the engine double makes the reduction observable without a real gate run |
| 5 | The verdict records the reduced marker, the phases run, and the ancestor | unit (`internal/gate`) | `go test ./internal/gate -run TestReducedVerdictRecordsAncestor` — observed red: the existing record shape has no such fields, so the round-trip drops them | Pins that the reduction happened and that what it inherited is identified rather than implied |
| 5 | The loader rejects a record mixing the full and reduced shapes | unit (`internal/gate`) | `go test ./internal/gate -run TestVerdictLoaderRejectsMixedShape` — observed red: a permissive loader accepts a record carrying an ancestor but no reduced marker | The existing loader validates fields strictly; a new class must fail closed on a malformed hybrid rather than guessing which it is |
| 6 | A second consecutive capture commit inherits from the same full green | unit (`internal/gate`) | `go test ./internal/gate -run TestConsecutiveReducedRunsShareAncestor` — observed red: an implementation reading the ancestor from the previous record's own identity finds a reduced verdict and falls back to a full run | The motivating scenario: without this the drain's second commit pays full and the feature does not deliver what it was approved for |
| 6 | A reduced verdict is never itself a valid ancestor | unit (`internal/gate`) | `go test ./internal/gate -run TestReducedVerdictIsNotAnAncestor` — observed red: chaining permitted, so a fixture chain produces evidence attributed to a run that graded nothing | The symmetric hazard to the row above: sharing an ancestor is sound, inheriting from a reduction is not |
| 6 | An ancestor older than the freshness window forces a full run | unit (`internal/gate`) | `go test ./internal/gate -run TestStaleAncestorFallsBackToFullRun` — observed red: re-stamping the inherited evidence with the current time keeps a day-old ancestor reusable forever | Inheritance that refreshes its own timestamp defeats the window the existing cache relies on |
| 6 | An allowlist-confined change with no ancestor runs the full gate | unit (`internal/gate`) | `go test ./internal/gate -run TestNoAncestorFallsBackToFullRun` — observed red: a missing ancestor treated as reusable emits a reduced verdict with an empty ancestor field | Fail-closed at the one state where there is nothing to inherit — a first commit, a fresh clone, or a pruned cache |
| 7 | `bench prep-release` refuses a reduced verdict and names that reason | contract (preprelease) | `go test ./internal/contract/surface/preprelease -run TestPrepReleaseRefusesReducedVerdict` — observed red: the precondition checks only for a green status and the reduced record satisfies it | The release path is the first consumer that must not accept a partial verdict; the refusal message is what tells the maintainer why |
| 8 | `bench commit` takes the reduced path for an allowlist-confined staged set | contract (runtime) | `go test ./internal/contract/runtime -run TestCommitReducesForConfinedStagedSet` — observed red: the commit path always calls whole-tree execution and the announced phase list is complete | The end-to-end behavior the reviewer asked for; asserting on the phase list rather than wall-clock keeps it deterministic |
| 8 | A staged set mixing an allowlisted and an unlisted path runs the full gate | contract (runtime) | `go test ./internal/contract/runtime -run TestCommitMixedStagedSetRunsFullGate` — observed red: confinement computed with "any" rather than "all" reduces a mixed commit | The most likely coding error in story 8, and the one whose consequence is an ungraded code change |
| 8 | The reduced run announces the skipped phases and the ancestor | contract (runtime) | `go test ./internal/contract/runtime -run TestReducedRunAnnouncesSkippedPhases` — observed red: silent reduction, matching the failure the existing reused-verdict announcement exists to prevent | A silent skip reads as a gate that never ran; the announcement is the operator's only signal that this verdict is narrower |
| 9 | The declaration, the phase set, and the prose cannot drift | conformance | `BENCH_CONFORMANCE_ROOT="$PWD" go test ./internal/conformance -run TestRootConformance` — observed red: mutating either side alone produces the binding diagnostic | Follows `checkLineBinding`'s shape, where a prose table and its machine-readable source red on divergence |
| 9 | The check bites in both directions | conformance | `go test ./internal/conformance -run TestDeclaredAllowlistBindingBites` — observed red: the bite test mutates the prose alone, then the declaration alone, and requires a diagnostic for each | A subset-or-substring binding survives a prose-only addition; running both directions is what rules that implementation out |

### Edge inventory

- **Empty / absent:** a missing verdict cache, a first commit with no ancestor, and an absent allowlisted file are covered by the no-ancestor row (story 6) and exact membership (story 1); an absent file is simply not in the staged set.
- **Malformed:** a record mixing the full and reduced shapes has its own row (story 5); a truncated or unreadable record is already refused by the existing loader, whose fail-closed posture this class inherits unchanged.
- **Boundary:** near-miss path names and the declared-directory boundary are the story 1 rows; the all-versus-any confinement boundary is the story 8 mixed-set row; the freshness-window boundary is the stale-ancestor row.
- **Adversarial content:** an allowlisted file whose *content* a check grades is the story 4 trio — the excludable phase cannot read it by hard read or by soft skip, and the included phase still must.
- **Capability degradation:** a phase that reports an informational skip rather than failing is a first-class class here rather than a footnote, because it is the kit's dominant idiom for a missing subject file and the first draft's central error. Covered by the soft-skip row.
- **Concurrency:** a full gate running elsewhere while a reduced run starts enters the existing execution-lock and fresh-green reuse path unchanged; no new lock is introduced. Two runs materializing stripped worktrees concurrently use distinct temporary roots, as the worktree owner already requires.
- **Interruption / crash:** an interrupted reduced run leaves no verdict, exactly as an interrupted full run does; a stripped worktree orphaned by an interrupt is retired by the existing worktree recovery path rather than by new cleanup.
- **Permission:** an unreadable allowlisted path during materialization opens the subject through the existing `open(reason)` path, which degrades to a non-closed subject and forces a real run.
- **Symlink / traversal:** allowlist entries are repository-relative and refused if they escape the root by the same lexical containment the manifest loader already applies; the one directory entry is a declared constant, not a pattern.
- **Unicode / encoding:** membership is byte-exact against a fixed list with no normalization, so a homoglyph or differently-normalized filename is not a member — the fail-closed direction.
- **Resource limits:** materialization copies a tracked tree measured at 13 MB and 1,316 files, and the build cache is content-keyed, so the marginal cost is seconds against a contract phase measured at 574; no new bound is introduced.
- **Interaction with the pre-push pin:** the pre-push hook pins the whole `.bench` tree, and the co-location migration moved the journal and the retros out of it, so a capture-only commit no longer moves the pinned tree and no longer needs `bench gate pin` before a push. That was true of the journal and retros before the migration and is named here because it is the one behavior change the move makes to this feature's surroundings — an improvement, not a regression, and not something this spec has to design for.

## Out of scope

- **The focused canary invocation (FT168's first face).** Running one named fixture or family as iteration evidence is a separate capability with its own selection surface; it shares this row's motivation but no seam. Estimate: 3 edits, 1 gate run.
- **`bench capture commit` porcelain (FT166's first half).** A conventional-message commit command for the capture set composes over this allowlist and needs it to exist first. Estimate: 2 edits, 1 gate run.
- **Extending the reduced path to any non-capture boundary** — a monorepo profile boundary, a declared package scope. That is FT101's grill, and reopening it here would smuggle a whole-tree ruling into a capture-set change. Estimate: not priced; starts as a decision, not a build.
- **Making the excludable phase set reviewer-configurable per repository.** The set is declared and asserted here for this kit; a linked repo declaring its own needs the two-audience analysis FT144 owns, and a manifest-declared excludability grammar for arbitrary roots is that work. Estimate: 4 edits, 2 gate runs.
- **Retiring `checkOccurrenceLedgerMigration`'s pinned count map** (FT172's third face). It stays graded by an included phase, so this spec neither needs nor makes that change. Estimate: 2 edits, 1 gate run.
