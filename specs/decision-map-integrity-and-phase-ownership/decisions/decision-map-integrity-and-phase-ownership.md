# Decision-map integrity and phase ownership

Status: ready

## Destination

Make decision maps an honest, situational planning surface with machine-checked
readiness, an operational dependency graph, and a clean authority boundary:
shaping resolves reviewer decisions and research context; spec authoring owns
engineering seams, tests, coverage, and gate attachment.

## #1: Which Wayfinder-derived improvements belong to this effort?

Blocked by: none
Type: Grill

### Question

Decide whether this effort takes every recommendation from the Wayfinder review
or only the changes that define decision-map integrity and phase ownership.

### Answer

This effort owns four connected changes: graph and readiness integrity,
situational map use, map-versus-spec authority, and decision-ticket
terminology. Prototype discipline remains separate future work. FT165 remains
the roadmap owner for domain-model maintenance. Tracker-backed maps,
assignment claims, and one-ticket-per-session remain rejected defaults.

## #2: When is a decision map required?

Blocked by: #1
Type: Grill

### Question

Decide whether every spec still requires a map or maps return to the
multi-session uncertainty case named by `/bench-shape-idea`'s trigger.

### Answer

Decision maps are situational. Use one when unresolved decisions form a
multi-session dependency tree. Do not manufacture a zero-fog map for a clear
idea merely to authorize `/bench-write-spec`.

## #3: Which phase owns engineering seams and test attachment?

Blocked by: #2
Type: Grill

### Question

Choose whether shaping or spec authoring owns module seams, deep-versus-thin
design, black-box test attachment, acceptance coverage, and gate attachment.

### Answer

`/bench-write-spec` owns those engineering derivations from the decided
constraints and current tree. A decision map owns the destination,
reviewer-owned behavioral or architectural decisions, rejected alternatives,
out-of-scope rulings, bounded spec-writer discretion, and required research
objects. A seam appears as a decision-ticket answer only when the reviewer
explicitly chose that seam.

The future decision-map format removes the nine-item `## Handoff`. The map
itself is the handoff: ticket answers own decisions, `## Sources` owns the
read-set, `## Out of scope` owns exclusions, and `## Spec-writer discretion`
owns intentionally delegated choices. This compiled map retains a transitional Handoff
because the currently installed phase contract still requires one.

## #4: Where does the spec writer's required read-set live?

Blocked by: #3
Type: Grill

### Question

Decide whether to add a research-object manifest to the Handoff or strengthen
the existing `## Sources` section.

### Answer

`## Sources` remains the single owner. It records evidence already used by
decision tickets and every code file, document, asset, or external source the
spec writer must read. Each entry names what it supports and its drift posture.
`/bench-write-spec` must consume and re-verify the section before deriving
seams or drafting the spec. No second research-object list appears in the
Handoff or spec.

For a ready map, the Go validator requires structured source entries.
Repository paths must resolve to regular files. External URLs are
syntax-checked but never fetched by the gate.

## #5: What makes a map honestly ready?

Blocked by: #2, #3, #4
Type: Grill

### Question

Define the state transition and the treatment of unresolved fog without making
an in-progress map gate-red.

### Answer

A map carries `Status: shaping | ready`. The Go gate always validates its
schema and dependency graph. `Status: ready` additionally requires every
decision ticket resolved, `## Not yet specified` empty, required sources valid,
and every remaining choice classified under `## Spec-writer discretion` or
`## Out of scope`. `/bench-write-spec` accepts only a ready map.

`## Spec-writer discretion` may contain only reversible technical choices that
do not change observable behavior, scope, an architectural boundary,
compatibility, or what the gate proves. Anything outside that boundary remains
a reviewer-owned decision ticket.

The shaping phase changes the tracked status line directly, then runs
`bench maps` and the focused Go validation. `bench maps` remains a read-only
AXI query; no status-mutating CLI is added.

## #6: What graph integrity does the oracle enforce?

Blocked by: #5
Type: Grill

### Question

Define whether `Blocked by` stays advisory prose or becomes a complete,
machine-checked graph.

### Answer

The existing `internal/maps` package remains the one parser used by CLI,
status, and gate consumers. It enforces unique decision-ticket numbers, valid
ticket types, blocker references to existing tickets, no self-edge, no cycle,
and no resolved ticket currently depending on an unresolved blocker. It derives
the actionable frontier from the same graph.

The Go gate validates every decision map that exists: active top-level maps and
compiled maps under `specs/*/decisions/`. A spec may legally have no map.
Existing maps migrate to the new schema in the same green implementation
change.

## #7: How does `bench maps` present the remaining route?

Blocked by: #6
Type: Grill

### Question

Choose the terminology and default query view once the CLI understands the
graph.

### Answer

`Decision ticket` is the canonical shaping term, distinct from an
independently-green implementation ticket. Human-facing output refers to a
decision ticket by its question or title; numeric IDs remain graph handles.

By default, `bench maps` emits every unresolved decision ticket with its title,
type, actionable state (`frontier`, `blocked`, or `deferred`), and blocker
titles. It does not hide blocked work merely because only the frontier is
currently takeable.

## #8: What authorizes spec writing without a map?

Blocked by: #2, #3, #4
Type: Grill

### Question

Define the no-map entry contract without allowing the spec writer to invent
reviewer-owned decisions.

### Answer

`/bench-write-spec` may proceed without a map from either a
reviewer-confirmed current conversation or a named reviewed source artifact,
such as a roadmap row, issue, or assessment, that fixes the destination and
constraints. The spec records that source.

When spec exploration encounters late uncertainty, ask one or two clarifying
questions one at a time with recommendations. If those answers expose a
dependency tree or multi-session fog, stop and route to
`$bench-shape-idea`; ordinary late uncertainty remains inside spec authoring.

## Not yet specified

## Spec-writer discretion

- Exact Go API and TOON column names, provided `internal/maps` remains the
  single semantic owner and the default output exposes the decided information.
- Exact conformance, contract, and canary fixture split proving the parser,
  gate, CLI, status, active-map, and compiled-map behaviors.
- The mechanical migration shape for existing maps, provided it lands in the
  same green change as the new validator.
- The exact spec field that cites a no-map source, provided the source is named
  and reviewer-verifiable.

## Out of scope

- Prototype-ticket lifecycle and prototype isolation; separate future work.
- Domain-model maintenance during shaping; FT165 already owns it.
- Tracker-backed maps, issue assignment as a claim, and one decision ticket per
  session; the prior rejection stands.
- Splitting maps into low-resolution indexes plus ticket files; revisit only
  with measured context pressure.
- Any change to implementation-ticket slicing or `craft-tickets`.

## Sources

- Path: `.agents/commands/bench-shape-idea.md`
  Supports: current map schema, Handoff requirement, lifecycle, and phase trigger.
  Drift: repository source; re-read by the spec writer.
- Path: `.agents/commands/bench-write-spec.md`
  Supports: current mandatory-map entry, seam ownership, source inputs, and clarification posture.
  Drift: repository source; re-read by the spec writer.
- Path: `.agents/skills/bench-craft-grill/SKILL.md`
  Supports: reviewer authority, one-question cadence, and shared-understanding gate.
  Drift: repository source; re-read by the spec writer.
- Path: `.agents/skills/bench-craft-tickets/SKILL.md`
  Supports: the distinct implementation-ticket meaning.
  Drift: repository source; re-read when changing terminology.
- Path: `internal/maps/maps.go`
  Supports: the existing single parser, unresolved-marker scan, and close-readiness projection.
  Drift: repository source; re-read before choosing the implementation seam.
- Path: `internal/status/status.go`
  Supports: the existing `maps.UnresolvedCount` consumer and dashboard behavior.
  Drift: repository source; re-read before changing parser output.
- Path: `projects/benchkit.md`
  Supports: the AXI query contract, gate shape, and one-source-per-fact standard.
  Drift: repository source; re-read before writing the spec.
- Path: `ROADMAP.md`
  Supports: FT165 as the existing domain-modeling owner and the absence of direct owners for the other decisions.
  Drift: repository working document; re-run `bench roadmap` before spec authoring.
- Path: `skills-assessment.md`
  Supports: the prior rejection of tracker-backed maps and one-ticket-per-session.
  Drift: historical assessment; verify closed decisions against the current tree before citing onward.
- URL: `https://github.com/mattpocock/skills/blob/main/skills/engineering/wayfinder/SKILL.md`
  Supports: decision-ticket terminology, graph frontier, fog, and situational-map comparison.
  Drift: mutable upstream `main`; re-fetch before citing onward.

## Handoff

1. **Module boundaries.** Exact build boundaries are deliberately not chosen
   here: `/bench-write-spec` owns them. The binding constraint is that
   `internal/maps` remains the single semantic parser consumed by CLI, status,
   and the Go gate check; command and skill prose contain no second parser.
2. **Contracts.** Maps are situational; existing maps carry
   `Status: shaping | ready`; ready means resolved decision tickets, empty fog,
   valid sources, and only bounded spec-writer discretion remaining. Specs may
   exist without maps. `bench maps` exposes the whole unresolved route and its
   actionable frontier.
3. **Deep vs thin.** `internal/maps` is the required deep owner of schema,
   graph, readiness, and source validation. Exact adapters and function
   boundaries are spec-owned.
4. **Black-box assertables.** A shaping map with honest fog remains valid; the
   same map marked ready is rejected. Duplicate IDs, invalid types, dangling,
   self, and cyclic blockers, or a resolved ticket blocked by open work are
   rejected. Valid active and compiled maps pass. A spec with no map is
   accepted. CLI output names every unresolved decision ticket and its
   frontier/blocked/deferred state.
5. **Gate attachment.** A Go gate check consumes the `internal/maps` parser and
   validates active and compiled maps. Focused package tests and deliberately
   broken canary fixtures must prove each refusal; the spec chooses their exact
   placement.
6. **Hostile-input owners.** The map parser owns absent, empty, malformed,
   unreadable, unsupported-schema, duplicate-ID, bad-type, bad-edge, cyclic,
   non-empty-ready-fog, invalid-discretion, and invalid-source states. The spec
   walks the project profile's filesystem and hostile-input checklist before
   locking fixtures.
7. **Uncertainty flags.** None. Exact APIs, output column names, fixture split,
   migration mechanics, and no-map source-field syntax are explicitly bounded
   spec-writer discretion above.
8. **Rejected alternatives.** Mandatory zero-fog maps; map-owned engineering
   seams; a second research-object manifest; an unvalidated prose graph; a
   second gate parser; a mutating maps CLI; retaining the architectural
   Handoff; tracker-backed default planning; one decision ticket per session.
9. **Domain watch-outs.** Keep decision tickets distinct from implementation
   tickets. Map location distinguishes active top-level work from compiled
   provenance. `## Sources` is both evidence and the required spec-writer
   read-set, never a second decision store.

Dependency order: n/a — single spec.
