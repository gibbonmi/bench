# FT215 focused test runs and the fast lane

Status: shaping

## Destination

Preserve `bench gate` as the whole-project oracle. Add non-authoritative
focused test runs for explicit and changed inputs. Make the fast lane select
only checks that its named paths can affect. Split the two outcomes into
separate specs after this map resolves their shared input rules.

## #1: Does the gate admit a scoped verdict?

Blocked by: none
Type: Grill

### Question

Should FT215 reopen the earlier decision that `bench gate` always grades the
whole project?

### Answer

Resolved 2026-08-27. `bench gate` always grades the whole project. A focused
test run writes no gate verdict and moves no green marker. Diff-scoped gating
remains out of scope.

## #2: Does focused iteration derive changed packages?

Blocked by: #1
Type: Grill

### Question

Should focused iteration derive touched Go packages and their reverse
dependents from a diff?

### Answer

Resolved 2026-08-27. `bench test` includes a changed-input mode. It derives
touched Go packages and their reverse dependents. The result is
non-authoritative and cannot satisfy a gate precondition.

## #3: Which public surface owns focused execution?

Blocked by: #1
Type: Grill

### Question

Should one non-oracle command own explicit test selection, changed-package
selection, and named conformance checks?

### Answer

Resolved 2026-08-27. `bench test` owns explicit package and test selection,
changed-package selection, and named conformance checks. FT215 rejects the
proposed `bench gate --check` form. Ticket #8 decides the exact public
grammar.

## #4: Does one spec own both outcomes?

Blocked by: #2, #3
Type: Grill

### Question

Should one spec bundle the focused-test surface and the path-aware fast lane?

### Answer

Resolved 2026-08-27. Use two specs. One spec owns the focused `bench test`
surface. One spec owns fast-lane selection from the paths that `bench commit`
names. This map resolves their shared boundaries before either spec starts.

## #5: Which inputs can focused test selection derive?

Blocked by: #2, #3
Type: Research

### Question

Produce `decisions/assets/ft215-focused-test-inputs.md`. Derive the input
partitions from these producers:

- The current `bench test` grammar, run binary, and result renderer.
- The coherent diff subject and its base and source-tip forms.
- The Go package and reverse-import graph.
- The conformance registry, selection transport, and input derivations.

Include staged, unstaged, committed, deleted, non-Go, and outside-module
changes when the producers emit them. State each unresolved behavior choice.
Do not choose the public grammar.

### Answer

— (open)

## #6: Which paths can narrow the fast lane?

Blocked by: #4
Type: Research

### Question

Produce `decisions/assets/ft215-fast-lane-inputs.md`. Derive the path
partitions from these producers:

- The paths that `bench commit` attributes to one commit.
- The current fast-lane phase table and path placeholders.
- Go source and `go:embed` inputs.
- Decision, spec, ticket, coverage, preflight, and capture inputs.

For each partition, name the checks that can observe it. Identify every
unknown partition and the safe fallback. State each unresolved behavior
choice.

### Answer

— (open)

## #7: Which diff subject does changed selection use?

Blocked by: #5
Type: Grill

### Question

Should changed selection inspect the worktree against `HEAD`, a named
base-to-tip subject, or both? What happens when it derives no runnable Go
package?

### Answer

— (open)

## #8: What is the focused test grammar?

Blocked by: #5
Type: Grill

### Question

Which flags express a package, a test pattern, changed selection, and a named
conformance check? What compatibility does the current positional package form
retain?

### Answer

— (open)

## #9: Does changed selection include conformance checks?

Blocked by: #5
Type: Grill

### Question

Should changed selection run only Go packages and reverse dependents? Should
it also select conformance checks through their registered input derivations?

### Answer

— (open)

## #10: How does the fast lane narrow safely?

Blocked by: #6
Type: Grill

### Question

Which path classes omit Go checks? Does an unknown or mixed input select the
current complete fast lane, or refuse the commit?

### Answer

— (open)

## Not yet specified

## Spec-writer discretion

- The internal type and function placement for a resolved selection.
- The cache strategy for package-graph queries within one command invocation.

## Out of scope

- A partial gate verdict, a changed-input green marker, or scoped evidence
  reuse.
- A change to the whole-project gate at landing.
- The spec-build checkpoint lifecycle in
  `decisions/spec-build-review-gate-cadence.md`.
- Release-tier selection or a general file-to-test build system.

## Sources

- Path: `roadmap/FT215.md`
  Supports: The destination, the two command projections, and the measured fixed-cost problem.
  Drift: Re-read if FT215 changes before this map becomes ready.
- Path: `decisions/gate-pipeline.md`
  Supports: #1's closed ruling against diff-scoped gating and #3's existing named-check transport.
  Drift: Re-read if the gate pipeline or conformance entry point changes.
- Path: `decisions/gate-budget.md`
  Supports: #1's whole-project oracle boundary.
  Drift: Re-read if the project reopens scoped gate authority.
- Path: `decisions/spec-build-review-gate-cadence.md`
  Supports: #4 and #6's focused ticket evidence and no-Go compile boundary.
  Drift: Re-read if that map changes its ticket-local evidence rules.
- Path: `internal/testreport/testreport.go`
  Supports: #3 and #5's current focused test producer.
  Drift: Re-read before #5 if the test command grammar or execution changes.
- Path: `internal/gate/lane.go`
  Supports: #4 and #6's current fast-lane phase and record boundaries.
  Drift: Re-read before #6 if the lane declaration or record changes.
