# Landing refusal standard

Status: staged

Roadmap: FT169

Decision source: `specs/landing-refusal-standard/decisions/landing-authority.md`

Verification log: 1 iteration(s) to accept — the single round returned six blocking findings; all folded to the decision source, no second round by reviewer rule.

## Problem

The landing refuses often, and its refusals do not tell the operator what to do
next. One refusal formatter exists, and it makes the recovery route optional.
Four of its twenty-nine call sites set that route. The audit grades eight
recorded refusal faces as partial and two as open.

Three faces lose information at a boundary. The fence refusal folds its computed
paths into a sentence, so no path table prints. The dirty-destination refusal
and the dirty-source refusal both pass an empty route. The residue refusal names
the paths but no repair.

The conflict recovery cannot tell two states apart. After a hand merge, the
source worktree either holds a pending merge or holds the committed resolution.
The refusal prints one repair line for both, and it never names
`git merge --continue`.

Two faces outside the landing package have the same fault. The wrapper's
exec-route refusal names the environment variable and not the command. The
`bench spec retire` primary-checkout refusal names the worktree creation and not
the retire step that follows it.

## Solution

The recovery route becomes mandatory at the refusal seam. Every refusal this
spec reaches names its failed paths and one exact next command. A declared
registry holds the landing's refusal faces, and a registry test drives one
producing fixture for each face. A face without a route turns the gate red.

The three formatter families align on that shape. The `refusal{}` formatter
gains a required route for each landing face. The `toon.Errorf(kind, hint)`
family already requires its hint. The bare `fmt.Errorf` and `echo` sites carry
the route in their text.

The landing's conflict recovery reads `MERGE_HEAD` in the source worktree. A
pending merge names `git merge --continue`. A committed resolution names the
commit-and-review route, which is `bench commit`, then the review, then the
re-run with the new tip.

The standard binds the faces the FT169 occurrences name. These are the landing
path, its resume, the wrapper's exec route, and `bench spec retire`. No other
verb changes.

## User stories

### The landing states the route out of every refusal

Line: opus / medium. The work rewrites refusal construction across one package
and adds one registry, so it needs a competent implementer and not a top tier.

1. As a landing operator, I want every refusal the landing prints to name one
   exact next command, so that I do not guess the repair.
2. As a landing operator, I want a path-computing refusal to print a path table,
   so that I read the paths directly.
3. As a landing operator, I want every landing-preflight refusal to end its
   route with my own re-run command, so that my flags stay the same.
4. As a Bench maintainer, I want a declared registry of the landing's refusal
   faces, so that a face without a route reds the gate.
5. As a Bench maintainer, I want the three refusal formatter families to agree
   on one shape, so that no family leaves the route optional.
6. As a landing operator, I want the resume path to print the first run's
   refusal shape, so that the route survives a resume.

### The conflict recovery tells the two merge states apart

Line: opus / medium. The work adds one `MERGE_HEAD` read and branches the
existing route builder, so the same tier covers it.

7. As a landing operator with a pending merge, I want the refusal to name
   `git merge --continue`, so that I finish it.
8. As a landing operator who committed the resolution, I want the refusal to
   name the commit-and-review route, so that I merge once.

### The named partial faces gain their route

Line: opus / medium. Each face is a small local edit against a named audit row,
so the tier stays the same across the group.

9. As a landing operator, I want the fence refusal to print the unfenced paths
   as a table, so that I see them at once.
10. As a landing operator, I want a short-circuited proof group to say that its
    later proofs did not run, so that I expect another refusal.
11. As a landing operator, I want the residue refusal to name the declaration
    file and a removal command, so that I choose a route.
12. As a landing operator, I want the dirty-destination refusal to print a route
    beside its moved paths, so that I clean the destination.
13. As a landing operator, I want the dirty-source refusal to print a route on
    one line, so that the hostile-source surface stays bounded.
14. As a landing operator, I want the exec-route refusal to name the exact
    landing command and its venue, so that I land correctly.
15. As a phase closer, I want the retire verb's primary-checkout refusal to name
    the retire step that follows, so that I read one route.
16. As an operator, I want the source-tip mismatch refusal to name the re-run
    with the observed tip, so that a moved tip has a route.

### The standard stays inside its decided reach

Line: opus / medium. These are guard rows against scope creep, and they need no
higher tier.

17. As a Bench maintainer, I want `bench commit` and the idea verb to keep the
    shared refusal unchanged, so that the retire route stays local.
18. As a Bench maintainer, I want the landing to stay the only verb that writes
    the default branch, so that this spec moves no authority.
19. As a Bench maintainer, I want no new pool counter, so that the four existing
    one-active mechanisms stay the only limit.
20. As a reader, I want `.bench/BENCH-reference.md` to state the refusal shape
    the code enforces, so that the doc and the code agree.

## Implementation decisions

**One declared registry holds the landing's refusal faces.** The registry
follows the precedent `identityComponents` sets in
`internal/worktree/identity_component.go`. Each entry names one face and its
route builder. One constructor turns an entry into a `refusalError`, so no site
composes a route of its own. The registry is the authoritative inventory for the
universal claim, and the registry test is its enforcement seam.

**The route builder requires a route.** A landing face constructs its refusal
through the registry constructor, which takes the route as a required argument.
The compiler therefore refuses an omitted route. An empty route string still
compiles, so the registry walk grades the value the face prints. The `refusal`
struct keeps its optional `next` field, because the verbs outside this spec's
reach still use it.

**Each landing-preflight route ends with the caller's own re-run.** The
face-specific repair comes first. The re-run repeats the same `--request`,
`--base`, `--source-tip`, and `--spec` values the caller passed. The existing
`atSourceWorktree` pointer form still covers a path that is not line-safe.

**The unauthorized paths travel typed out of the preflight package.**
`preflight.AuthorizeReviewedSource` returns a typed error that carries the
unauthorized path list. `landingSourceRange` reads that list into the refusal's
paths rather than into its detail sentence. The joined detail string stays in
`pathsAuthorizedCheck`, because the preflight report reads it.

**The landing reads `MERGE_HEAD` in the source worktree.** The read decides
which of the two routes the conflict-related refusals name. A present
`MERGE_HEAD` names `git merge --continue`, and an absent one names the
commit-and-review route. An unreadable Git directory leaves the state undecided,
and the route falls back to the commit-and-review form. The `bench commit`
contract does not change, and FT258 still owns it.

**`bench spec retire` composes its own route.** The shared
`usage.PrimaryCheckoutRefusal` keeps its current text for `bench commit` and the
idea verb. The retire verb appends its own follow-on step, so the shared source
stays one source and the other two callers stay unchanged.

**The wrapper's exec-route refusal reads no repository.** The refusal text gains
the exact landing command and the venue sentence. The refusal still fires before
the first repository read, which the existing route test proves with an absent
marker file.

## Testing decisions

A good test drives the real command and reads the printed refusal record. The
record is the operator's whole evidence, so the assertion reads `refused{...}`,
its `next=` field, and its `refusal_paths` table rather than an internal value.

The engineering seams are the landing command in `internal/worktree`, the
wrapper route in `internal/systemtest`, the retire verb in `internal/spec`, and
the anchors registry in `internal/anchors`. The prior art is
`TestLandCommandReportsEveryRefusalInOnePreflight` for the preflight surface,
`TestWorktreeLandRouteRefusesInheritedRoutingBeforeRepositoryReads` for the
wrapper, and `TestRetireOnPrimaryCheckoutRefusesAndDeletesNothing` for the
retire verb.

The gate observes the feature through the ordinary Go test run and through the
canary fixture bite for the new documentation anchor.

### Seam diagram

    trigger: an operator runs `bench worktree land`, its resume, or `bench spec retire`
        │
        ▼
    failed proof  ──▶  [ registry face + route builder ]  ──▶  refused{detail,next} + refusal_paths
                          ◀ tests attach here: drive the command, read the printed record

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| LRS1 | 4 | The registry walk drives one producing fixture for each declared landing refusal face | `internal/worktree/identity_component_test.go` (`TestIdentityComponentRegistryHasAProducingFixture`) | A face declared without a fixture leaves the walk with no producer, and the walk reds |
| LRS2 | 1, 5 | The registry walk asserts that each face's fixture prints a non-empty `next=` field | `internal/worktree/identity_component_test.go` (`TestIdentityComponentRegistryHasAProducingFixture`) | A face that passes an empty route string still compiles, and the walk's non-empty assertion reds on it |
| LRS3 | 3 | Each landing-preflight refusal ends its `next=` field with the caller's own `bench worktree land` re-run | `internal/worktree/land_surface_test.go` (`TestLandCommandReportsEveryRefusalInOnePreflight`) | A face-specific repair without the re-run tail fails the suffix assertion |
| LRS4 | 2, 9 | The fence refusal prints `paths_total=` and a `refusal_paths` row naming the unfenced path | `internal/worktree/land_surface_test.go` (`TestLandCommandFenceRefusalNamesThePath`) | A detail sentence that carries the path prints no table, and the table assertion reds |
| LRS5 | 6 | Each resume refusal prints a non-empty `next=` field naming the `bench worktree land --resume` continuation | `internal/worktree/land_resume_refusal_test.go` (`TestResumeLandCommandPublicRefusesDestructiveDestinationState`) | A resume face that keeps a bare detail prints no resume continuation, and the assertion reds |
| LRS6 | 7 | A source worktree holding `MERGE_HEAD` produces a `next=` field naming `git merge --continue` | `internal/worktree/land_journey_test.go` (`TestLandCommandPublicConflictRepairRequiresNewReviewedTip`) | A route that names a second `git merge` omits the `--continue` token, and the row reds |
| LRS7 | 8 | A source worktree without `MERGE_HEAD` produces a `next=` field that names the commit-and-review route and omits `git merge --continue` | `internal/worktree/land_surface_test.go` (`TestLandCommandConflictRefusalNamesTheSourceRepair`) | An implementation that appends `git merge --continue` to every conflict route names it here too, and the absence assertion reds |
| LRS9 | 10 | A refusal from a short-circuited proof group names `later proofs in this group did not run` in its `next=` field | `internal/worktree/land_surface_test.go` (`TestLandCommandReportsIdentityAndDestinationInOnePreflight`) | A refusal that reports its own fault alone omits that sentence, and the row reds |
| LRS10 | 11 | The residue refusal's `next=` field names `.bench/build-outputs.json` | `internal/worktree/land_release_refusal_test.go` (`TestLandCommandRefusalListsIgnoredPaths`) | A route that names removal alone omits the declaration file, and the row reds |
| LRS11 | 11 | The residue refusal's `next=` field names an exact removal command for the undeclared path | `internal/worktree/land_release_refusal_test.go` (`TestLandCommandRefusalListsIgnoredPaths`) | A route that names the declaration alone omits the removal command, and the row reds |
| LRS12 | 12 | The dirty-destination refusal prints a `next=` field beside its moved-path table | `internal/worktree/land_release_refusal_test.go` (`TestLandCommandRefusalListsDestinationPaths`) | An empty route argument prints the table with no `next=` field, and the row reds |
| LRS13 | 13 | The dirty-source refusal prints a `next=` field and no `refusal_paths` table | `internal/worktree/land_identity_test.go` (`TestLandCommandInvalidatesAChangedSourceFingerprintBeforeTheGate`) | A route added through the shared path-carrying refusal prints a table, and the no-table assertion reds |
| LRS14 | 14 | The wrapper's exec-route refusal prints the exact `bench worktree land` command | `internal/systemtest/land_route_test.go` (`TestWorktreeLandRouteRefusesInheritedRoutingBeforeRepositoryReads`) | A refusal that names the variable alone omits the command token, and the row reds |
| LRS15 | 14 | The wrapper's exec-route refusal states that the landing runs outside `bench worktree exec` | `internal/systemtest/land_route_test.go` (`TestWorktreeLandRouteRefusesInheritedRoutingBeforeRepositoryReads`) | A refusal without the venue sentence leaves the operator inside the exec, and the row reds |
| LRS16 | 15 | The `bench spec retire` primary-checkout refusal names `bench spec retire` as the follow-on step | `internal/spec/spec_test.go` (`TestRetireOnPrimaryCheckoutRefusesAndDeletesNothing`) | A refusal equal to the shared line names only the worktree creation, and the equality assertion reds |
| LRS17 | 16 | The source-tip mismatch refusal's `next=` field names the observed tip | `internal/worktree/land_journey_test.go` (`TestLandCommandPublicConflictRepairRequiresNewReviewedTip`) | A route that names a placeholder tip omits the observed value, and the row reds |
| LRS18 | 17 | The idea verb's primary-checkout refusal names no `bench spec retire` step | `internal/roadmap/roadmap_test.go` (`TestIdeaRefusesPrimaryCheckout`) | A retire route added inside the shared function reaches the idea verb, and the absence assertion reds |
| LRS19 | 20 | The anchors registry requires the refusal-shape sentence in `.bench/BENCH-reference.md` | `internal/anchors/registry_data_test.go` (`TestWorktreeRuleAnchorsRedOnRemoval`) | A doc edit that drops the sentence leaves the needle unmatched, and the fixture bite reds |
| LRS20 | edge | A refusal raised before the assignment resolves names the operator's own worktree path in its re-run | `internal/worktree/land_surface_test.go` (`TestLandCommandReportsIdentityAndDestinationInOnePreflight`) | A route that addresses an empty assignment id prints an unusable command, and the row reds |
| LRS21 | edge | A source worktree path that is not line-safe produces a `next=` field in the `bench worktree exec` pointer form | `internal/worktree/land_reauthorization_test.go` (`TestLandCommandReauthorizeRecoveryPointsThroughUnsafePath`) | A route that quotes the unsafe path produces a command the operator cannot paste, and the row reds |
| LRS22 | edge | An unreadable source Git directory produces the commit-and-review route rather than the `--continue` route | `internal/worktree/land_journey_test.go` (`TestLandCommandPublicConflictRepairRequiresNewReviewedTip`) | A `MERGE_HEAD` read that treats an error as a present file names the wrong route, and the row reds |

Not covered: story 18 — this spec moves no authority, so the existing landing
suite is the only oracle, and a new row would assert an absent mutation.

Not covered: story 19 — this spec adds no counter, and a row that asserts an
absent counter cannot tell an unwritten counter from a removed one.

### Edge inventory

- A refusal raised before the assignment resolves. Covered by LRS20.
- A worktree path that is not line-safe. Covered by LRS21.
- An unreadable source Git directory during the `MERGE_HEAD` read. Covered by
  LRS22.
- **Won't handle** — an undeclared residue above `ignoredEntryLimit`: the
  bounded `paths_total=` behavior lives in the shared classifier. Another verb's
  test already pins it, and this spec changes no other verb.
- **Won't handle** — a third merge state beside the pending merge and the
  committed resolution: decision #9 decides two states, and a third is an
  addition beyond the decision source.
- **Won't handle** — a pending merge in the destination checkout: the
  dirty-destination refusal already fires first and owns that route.
- **Won't handle** — a second sanitizer for control characters in a route:
  `sanitize.Controls` already covers every printed field, and the landing
  surface tests pin it.
- **Won't handle** — a route for a verb no FT169 occurrence names: decision #12
  binds the standard to the named faces. Those verbs adopt it when someone
  touches their code.
- **Won't handle** — a machine-readable route field beside `next=`: the record
  is one line for an operator, and no consumer parses the route today.

## Ownership fences

- `internal/worktree/land_refusal.go`
- `internal/worktree/land_identity.go`
- `internal/worktree/land.go`
- `internal/worktree/land_resume.go`
- `internal/worktree/classifier.go`
- `internal/worktree/identity_component.go`
- `internal/worktree/land_surface_test.go`
- `internal/worktree/land_identity_test.go`
- `internal/worktree/land_journey_test.go`
- `internal/worktree/land_resume_refusal_test.go`
- `internal/worktree/land_release_refusal_test.go`
- `internal/worktree/land_reauthorization_test.go`
- `internal/worktree/land_tickets_only_test.go`
- `internal/worktree/identity_component_test.go`
- `internal/preflight/gather.go`
- `internal/preflight/decision.go`
- `internal/preflight/decision_test.go`
- `internal/spec/spec.go`
- `internal/spec/spec_test.go`
- `internal/roadmap/roadmap_test.go`
- `internal/systemtest/land_route_test.go`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry_data_test.go`
- `tests/canary/workflow-guidance-anchors/`
- `bin/bench.sh`
- `.bench/BENCH-reference.md`
- `reviews/landing-refusal-standard.md`

The bound packages `internal/anchors`, `internal/roadmap`, and
`internal/worktree` each carry the command-registry file set, so the fence adds
all five:

- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/main_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_test.go`

The canary fixtures that pin `.bench/BENCH-reference.md` outside
`tests/canary/workflow-guidance-anchors/` are:

- `tests/canary/skills-index-command-adapters/unindexed-skill`
- `tests/canary/skills-index-command-adapters/stale-index-wording`
- `tests/canary/skills-index-command-adapters/dangling-index`
- `tests/canary/skills-index-command-adapters/missing-index-field`
- `tests/canary/skills-index-command-adapters/debug-implicit-invocation-reverted`
- `tests/canary/skills-index-command-adapters/command-invocation-disabled-against-policy`
- `tests/canary/skills-index-command-adapters/adapter-inert-invocation-key`
- `tests/canary/docs-currency-token-diet/benchref-pointer-dropped`
- `tests/canary/docs-currency-token-diet/benchref-section-duplicated`
- `tests/canary/docs-currency-token-diet/benchref-imported`

The canary fixtures that pin `bin/bench.sh` are:

- `tests/canary/package-core-guard/unrouted-subcommand`
- `tests/canary/package-core-guard/reintroduced-bare-skip`
- `tests/canary/package-core-guard/bounds-duplicate-owner`
- `tests/canary/load-validity-metadata/extensionless-gate-ref`
- `tests/canary/docs-currency-token-diet/stale-skill-cli-reference`
- `tests/canary/docs-currency-token-diet/stale-cli-doc-reference`
- `tests/canary/docs-currency-token-diet/missing-cli-inventory`

The build may not edit this spec's acceptance rows, its ownership fences, or its
budget targets. A shortfall against a row returns to `/bench-write-spec`.

`internal/usage/worktree.go` is outside the fence, because story 18 requires the
shared refusal to stay unchanged. `internal/commit/commit.go` and
`internal/roadmap/learning.go` are outside the fence for the same reason.

The reference-doc edit keeps the anchored sentences verbatim. These are
"The spec is optional on the landing and on its resume" and the fast-lane
landing-shape marker. The doc edit and the new anchor move lines that canary
fixtures under `tests/canary/workflow-guidance-anchors/` compare. The fixture
bite test therefore belongs in the focused checks the build runs before the
gate.

## Out of scope

- **The landing interface spec.** It carries the symbolic source tip, the one
  `--spec` normalization seam, the staged-spec refusal, and the preflight base
  mirror from decisions #8 and #11. It also carries the accepted-forms clause on
  the source-tip mismatch refusal, because that clause states the interface and
  not the route. Estimate: 14 edits, 4 gate runs.
- **The `bench commit` `MERGE_HEAD` contract.** FT258 owns it, and FT254 is
  blocked on it. This spec reads `MERGE_HEAD` on the landing side only.
  Estimate: 6 edits, 2 gate runs.
- **The capture-storage decision.** FT249 owns it. The landing states its
  capture-related route after that decision lands. Estimate: 5 edits, 2 gate
  runs.
- **A repo-wide refusal reword.** Decision #12 binds the standard to the faces
  the FT169 occurrences name. Every other verb adopts the standard when someone
  touches its code. Estimate: 30 edits, 3 gate runs.
- **The `spec.Resolve` working-directory fault.** The audit found it beside the
  thirteen occurrences, and the one normalization seam in the interface spec
  removes it. Estimate: 4 edits, 1 gate run.

## Further notes

The audit grades occurrences 1, 6, and 7 as fixed by cause removal. This spec
does no work for them.

The strongest fixes in this area removed a cause rather than reworded a message.
The registry is the same move at the seam level: it removes the cause of an
absent route, because a face cannot exist without one.
