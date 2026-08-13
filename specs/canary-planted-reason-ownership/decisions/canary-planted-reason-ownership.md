# FT153 canary planted-reason ownership

Status: ready

## Destination

One branch-native canary contract for the Bench kit and linked repos. The map
decides what the inventory command proves, who owns a retained fixture's
planted-reason proof, what a linked repo may claim, which stale canary
statements must leave, and whether the repair is one spec or multiple specs.
It does not reinstate the retired nested sweep or authorize the public command
to execute a full fixture journey.

## #1: What contract does the current tree actually implement?

Blocked by: none
Type: Research

### Question

After the branch-native test refactor, which surface inventories fixtures,
which surface proves a planted reason, and where does the current proof stop?

### Answer

Resolved 2026-08-12 from the current tree. `bench canary` is an inventory
command: it discovers fixtures, resolves one binding for each, and aggregates
the selection without materializing a fixture or invoking its check
(`internal/canary/inventory.go:90-140`). The current profile assigns the
planted-red and restored-green proof to ordinary mutation tests
(`projects/benchkit.md:198-216`), and `runFixtureBite` implements that direct
shape for the fixtures it receives
(`internal/conformance/fixture_bite_test.go:545-567`).

That proof is incomplete. The `package-core-guard` bite test names five of its
33 fixtures and delegates the remainder to a canary sweep that no longer
exists (`internal/conformance/fixture_bite_test.go:222-239`). A sentinel placed
in one excluded fixture's `EXPECT` left both `bench canary` and the complete
`internal/conformance` package green. The linked-repo scaffold likewise calls
only the inventory command while claiming that its seed fixture proves the
example check bites (`internal/adopt/init.go:117-133`). The current FT153 work
is therefore missing planted-reason ownership, not a mandate to turn the
inventory command into an executor.

## #2: What planted-reason guarantee does a linked repo receive?

Blocked by: #1
Type: Grill

### Question

Should Bench preserve a generic linked-repo promise by introducing a
language-agnostic owner-execution protocol, or should the kit state the
boundary honestly: Bench inventories fixture bindings while each linked
repo's native gate tests own any planted-reason proof?

### Answer

Resolved 2026-08-12: state the boundary honestly. A linked repo receives
fixture-inventory validation from Bench, not a generic guarantee that Bench
executed its project checks against fixture mutations. Its project-native gate
tests own any planted-reason proof, and the kit must not advertise that proof
unless a later reviewed decision introduces an executable owner protocol.

`bench init` must stop seeding an example fixture whose `EXPECT` no retained
test consumes and must stop claiming that `bench canary` proves the example
check bites. A newly initialized repo remains intentionally red until the
project supplies real checks, real fixture bindings where it wants them, and
its own check-specific tests. The alternative—making the public command own a
generic fixture journey—is rejected because no current owner protocol exists
and the branch-native architecture deliberately removed nested gate and Go
subprocess canaries.

## #3: What qualifies as a retained kit fixture?

Blocked by: #1
Type: Grill

### Question

May a kit fixture remain as inventory-only evidence when its old sweep owner
has vanished, or must every retained fixture carry direct observable
planted-red and restored-green proof through its exact registered owner?

### Answer

Resolved 2026-08-12: every retained kit fixture must carry the direct proof.
Its `EXPECT` is mechanically compared with the diagnostic produced after the
fixture mutation is materialized, and the same diagnostic is absent after the
subject is restored. No fixture may rely on the removed sweep, an empty-tree
collision screen, a hard-coded sibling diagnostic, or inventory membership as
a substitute.

If a current fixture cannot run through an in-process registered owner, it
does not remain as ungraded canary evidence. The implementation may migrate
the check to such an owner, or retire the fixture and replace it with an
owning-package mutation test that preserves the same semantic tripwire. In
either case, weakening the behavior the fixture was retained to guard must
make the ordinary gate red for the named reason. Reintroducing a nested gate,
wrapper, `go test`, or `go run` invocation per fixture is rejected.

## #4: What does the public canary command say and prove?

Blocked by: #1, #2
Type: Grill

### Question

Keep the current “owners dispatched” vocabulary, make the command execute
owners, or retain inventory-only behavior and rename every public claim to
match it?

### Answer

Resolved 2026-08-12: retain inventory-only behavior and make its vocabulary
honest. A successful public run reports exactly
`canary inventory ok (<n> fixture bindings)`. It proves a non-empty,
well-formed inventory with one accepted binding per fixture; it does not claim
that a binding was invoked or that the fixture bit.

The command keeps its current grammar and exit classes. Invalid, duplicate,
empty, or unbound inventory remains red. Production comments, result fields,
tests, the ship inventory check, and system-journey expectations must stop
using “dispatch”, “invoke”, or “sweep” for this inventory-only path. The exact
internal factoring remains a spec decision, but no vestigial callback may be
presented as a production owner invocation.

## #5: Which current-state documents move with the contract?

Blocked by: #2, #3, #4
Type: Grill

### Question

Limit the repair to code and tests, or reconcile the linked-repo scaffold,
profile, and ADRs that still describe the retired sweep and the false consumer
promise?

### Answer

Resolved 2026-08-12: reconcile every current-state claim in the same delivery.
ADR 0001 must distinguish the kit's direct ordinary planted-reason proofs from
the linked repo's inventory-only Bench guarantee. ADRs 0003 and 0009 must stop
describing the removed full inner-gate sweep and record their superseded or
deprecated status under `craft-adr`. The project profile remains the authority
for the branch-native split and is corrected only where the fixture census or
wording is incomplete.

The init scaffold, compatibility shim comments, production canary comments,
system journey, and any conformance rationale that says an absent sweep proves
a fixture must all state the resulting contract. History stays in Git; the
documents describe only the resulting state.

## #6: Is this one spec or two independently useful specs?

Blocked by: #2, #3, #4, #5
Type: Grill

### Question

Split kit fixture coverage from linked-repo and inventory-vocabulary honesty,
or deliver one contract repair with independently-green tracer tickets?

### Answer

Resolved 2026-08-12: one spec, with independently-green tracer tickets. The
kit proof, linked-repo boundary, CLI vocabulary, and ADR/scaffold claims are
different consumers of one canary contract. Splitting them would permit a
landed state in which the code and current-state documentation intentionally
describe different guarantees.

The spec may sequence an inventory-vocabulary tracer, retained-fixture proof
tracers, and a current-state documentation tracer, but its final accepted
candidate closes the complete contract above.

## #7: What happens to FT168 focused canary selection?

Blocked by: #4, #6
Type: Grill

### Question

Does this map also authorize fixture or family execution through
`bench canary`, or does FT168 remain blocked for revalidation against the
branch-native contract?

### Answer

Resolved 2026-08-12: FT168 remains separate and blocked on FT153. This map
does not add execution selectors to an inventory-only command. Before FT168
can be specified, its removed full-sweep premise must be revalidated and any
focused iteration surface must select the direct ordinary planted-reason
proofs without becoming a second oracle or receiving gate credit.

`/bench-shape-idea` never edits the roadmap. FT153 and the FT153-to-FT168 edge
remain there until the normal reviewed roadmap drain reconciles their text.

## Not yet specified

## Spec-writer discretion

- Per fixture, migrate to an in-process registered owner or retire the fixture
  in favor of an equivalent owning-package mutation test, provided #3's
  observable red/restore predicate and existing semantic protection both hold.
- Internal names and factoring for the inventory decision, provided #4's exact
  public output and no-execution contract hold.
- Independently-green ticket order inside the single spec, provided no landed
  ticket makes a stronger canary claim than the code at that ticket proves.

## Out of scope

- A production fixture-journey owner in `bench canary`.
- Nested gates, per-fixture Go subprocesses, or restoration of the retired
  canary sweep.
- FT168's focused fixture/family execution surface.
- Editing or retiring FT153 or FT168 in `ROADMAP.md` during shaping.
- Gate performance pricing; this map fixes correctness and ownership only.

## Sources

- Path: `specs/canary-planted-reason-ownership/decisions/assets/ft153-branch-native-canary-contract.md`
  Supports: #1's current-tree derivation and the factual premises of #2 through #7, including the two executed sentinel probes and the 31-member reader census.
  Drift: re-verify if the gate phase table, canary inventory path, conformance registry or fixture-bite tests, init scaffold, project profile, or retained fixture inventory changes before `/bench-write-spec` reads this map.
