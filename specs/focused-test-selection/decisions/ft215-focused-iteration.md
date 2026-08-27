# FT215 focused test runs and the fast lane

Status: ready

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

Produce
`specs/focused-test-selection/decisions/assets/ft215-focused-test-inputs.md`.
Derive the input
partitions from these producers:

- The current `bench test` grammar, run binary, and result renderer.
- The coherent diff subject and its base and source-tip forms.
- The Go package and reverse-import graph.
- The conformance registry, selection transport, and input derivations.

Include staged, unstaged, committed, deleted, non-Go, and outside-module
changes when the producers emit them. State each unresolved behavior choice.
Do not choose the public grammar.

### Answer

Resolved 2026-08-27. The cited summary is
`specs/focused-test-selection/decisions/assets/ft215-focused-test-inputs.md`.
The current command accepts
one package expression and no test filter. The diff owner supplies coherent
live and frozen subjects.

A live Go 1.25 query verified production, test, and
external-test import edges. No current producer maps diff paths to packages or
conformance checks. Tickets #7, #8, and #9 own the remaining choices.

## #6: Which paths can narrow the fast lane?

Blocked by: #4
Type: Research

### Question

Produce
`specs/focused-test-selection/decisions/assets/ft215-fast-lane-inputs.md`.
Derive the path
partitions from these producers:

- The paths that `bench commit` attributes to one commit.
- The current fast-lane phase table and path placeholders.
- Go source and `go:embed` inputs.
- Decision, spec, ticket, coverage, preflight, and capture inputs.

For each partition, name the checks that can observe it. Identify every
unknown partition and the safe fallback. State each unresolved behavior
choice.

### Answer

Resolved 2026-08-27. The cited summary is
`specs/focused-test-selection/decisions/assets/ft215-fast-lane-inputs.md`.
The attributed-path producer
accepts safe files, directories, deletions, and symlinks.

The current lane
always runs gofmt, prose, vet, and build. Positive path classes can omit a
check only when its producer cannot observe that class. Unknown paths have no
safe narrow selection. Ticket #10 owns the class and fallback policy.

## #7: Which diff subject does changed selection use?

Blocked by: #5
Type: Grill

### Question

Should changed selection inspect the worktree against `HEAD`, a named
base-to-tip subject, or both? What happens when it derives no runnable Go
package?

### Answer

Resolved 2026-08-27. Changed mode accepts both coherent subject families. Its
default and explicit base produce live subjects. An explicit base and source
tip produce an immutable subject.

A proven non-Go subject emits an explicit
empty result. Any Go-relevant path that cannot resolve safely refuses. The
command never drops or widens that path silently.

## #8: What is the focused test grammar?

Blocked by: #5
Type: Grill

### Question

Which flags express a package, a test pattern, changed selection, and a named
conformance check? What compatibility does the current positional package form
retain?

### Answer

Resolved 2026-08-27. The public forms are:

- `bench test [--full] [--package <expr> | <legacy-package> | --changed]
  [--run <go-regex>]`
- `bench test [--full] --check <name>`

The positional package form remains compatible. `--check` is exclusive.
`--run` can combine with the default, package, or changed subject. A run that
matches zero tests refuses. Ticket #11 adds the changed-subject flags.

## #9: Does changed selection include conformance checks?

Blocked by: #5
Type: Grill

### Question

Should changed selection run only Go packages and reverse dependents? Should
it also select conformance checks through their registered input derivations?

### Answer

Resolved 2026-08-27. Changed mode selects Go packages and their reverse
dependents only. It does not derive conformance checks from changed paths. The
`Inputs` labels are not path selectors, and catch-all rows would make the
result unbounded. `--check <name>` remains the explicit conformance form.

## #10: How does the fast lane narrow safely?

Blocked by: #6
Type: Grill

### Question

Which path classes omit Go checks? Does an unknown or mixed input select the
current complete fast lane, or refuse the commit?

### Answer

Resolved 2026-08-27. The lane classifies the composed changed paths and takes
the union of their checks:

- Go source selects gofmt, vet, and build.
- Go metadata or a named embed input selects vet and build.
- Markdown or prose policy selects prose.
- A known document class also selects its existing focused validator.
- Mixed known paths select the union.
- Any unknown or symlink path selects the current complete lane.

The attribution owner continues to refuse special and unreadable paths before
the lane. The composed-path input expands named directories and represents a
rename as its deletion and addition.

## #11: How does changed mode name a frozen subject?

Blocked by: #7, #8
Type: Grill

### Question

Should changed mode reuse the coherent diff flags for live and frozen
subjects?

### Answer

Resolved 2026-08-27. Changed mode accepts
`--base <commit> [--source-tip <commit>]`. A base without a source tip
selects the live subject. A base with a source tip selects the immutable pair.
A source tip without a base refuses as a grammar error.

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
- Path: `specs/focused-test-selection/decisions/assets/ft215-focused-test-inputs.md`
  Supports: #5's producer-derived package, diff, and conformance input partitions.
  Drift: Re-run the research if test, diff, Go graph, or conformance producers change.
- Path: `specs/focused-test-selection/decisions/assets/ft215-fast-lane-inputs.md`
  Supports: #6's producer-derived attributed-path and fast-lane partitions.
  Drift: Re-run the research if attribution, lane, embed, or document producers change.
