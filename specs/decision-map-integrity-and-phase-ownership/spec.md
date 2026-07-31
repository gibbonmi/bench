# Decision-map integrity and phase ownership

Status: implemented

Decision source: `specs/decision-map-integrity-and-phase-ownership/decisions/decision-map-integrity-and-phase-ownership.md`

Same-session map: reviewer-directed on 2026-07-30 after every load-bearing fork
was closed. The reviewer may veto the author-derived seams and exact interface
choices below before implementation.

## Problem

Decision maps currently claim more authority than they can prove. Their
dependency edges are prose, readiness is inferred from unresolved markers and a
handwritten `## Handoff`, `bench maps` does not derive the actual frontier, and
the gate does not validate active or compiled maps. At the same time, every spec
is forced through a map and shaping pre-decides engineering seams that should be
derived later from the current tree.

The result is ceremony without integrity: a zero-fog idea gets a manufactured
map, a malformed or cyclic map can look ready, compiled provenance can rot
outside the query surface, and the spec writer can inherit stale engineering
choices from the wrong phase.

## Solution

Make decision maps situational and machine-checked. A map is used only for a
multi-session unresolved decision tree, carries an explicit `shaping` or `ready`
status, and is parsed into one graph model that feeds validation, frontier
rendering, status counts, and the gate. Ready maps have no unresolved decision
tickets or fog, use structured and locally verifiable sources, and leave only
bounded reversible choices to the spec writer.

Shaping owns reviewer decisions, constraints, exclusions, research objects, and
bounded discretion. Spec authoring re-reads those sources and owns engineering
seams, tests, coverage, and gate attachment. A clear idea may enter spec writing
from a reviewer-confirmed conversation or named reviewed artifact without a
manufactured map.

### Defaulted decisions and veto points

These choices were derived while writing the spec, not explicitly approved in
the same-session map. They are part of this approval surface.

| decision | default | veto consequence |
|---|---|---|
| canonical skeleton | Add read-only `bench maps --template`, derived from the parser’s schema data, so exact Markdown grammar has one source. | Keeping the skeleton only in phase prose requires another single-source mechanism; dropping a rendered skeleton leaves the author without a paste-ready format. |
| candidate discovery | Preserve current active behavior: non-hidden direct Markdown children, case-insensitive README exemption, nested asset directories excluded; apply the same direct-child rule to compiled map directories. | Recursive or hidden discovery reclassifies map-owned assets/editor files and needs a new asset convention and migration. |
| required shape | Require title, status, destination, at least one decision ticket, the four terminal sections, ticket question/answer fields, and fence-aware parsing. | A looser schema must name which omissions stay valid and how the validator distinguishes a map from an asset or malformed note. |
| source cardinality | Require the `## Sources` section, but allow it to be empty when the map has no external or repository research object; every present entry remains structured and valid. | Requiring at least one source would force a synthetic citation on an all-grilled map; allowing unstructured `n/a` would add a second entry form. |
| compiled status | Require every compiled map to be `ready`; only active maps may remain `shaping`. | Allowing compiled shaping maps needs a defined post-spec owner and query surface for unfinished decisions. |
| active query and count | Count each shaping or invalid active map once, omit compiled maps from query/status, and emit a fog row when a shaping map has no unresolved ticket. | A different projection must still prevent an empty detailed query from disagreeing with a positive dashboard count and keep compiled provenance out of the shaping frontier. |

## User stories

1. As a map author, I want one canonical map schema and skeleton, so that the
   document I write is the document every consumer parses. Line:
   `gpt-5.6-terra / medium.` The parser is a known Go seam, but schema and
   filesystem correctness need the project’s mid-tier gate posture.

2. As a reviewer, I want map dependencies validated and the whole unresolved
   route rendered with actionable states, so that frontier, blocked, and
   deferred work cannot be guessed differently by different sessions. Line:
   `gpt-5.6-terra / medium.` The graph rules are exact and testable, but cycle
   and resolution-state interactions are not cheap plumbing.

3. As a spec writer, I want readiness and research inputs verified before I
   derive engineering seams, so that I start from settled reviewer decisions and
   current evidence rather than stale or incomplete context. Line:
   `gpt-5.6-terra / medium.` The behavior is precise and gate-observable, while
   hostile path handling warrants the mid tier.

4. As the gate operator, I want one Go conformance check to validate every
   active and compiled decision map, so that malformed provenance cannot remain
   green merely because `bench maps` does not list it. Line:
   `gpt-5.6-terra / high.` This changes the oracle and must prove every refusal
   with its own targeted bite.

5. As a cold-session user, I want `bench maps`, `bench maps --count`, and
   `bench status` to project the same active-map model, so that the dashboard
   never disagrees with the detailed query or hides shaping fog. Line:
   `gpt-5.6-luna / low.` The existing CLI, AXI, and status seams fully observe
   this mechanical projection once the graph model exists.

6. As a reviewer with a clear idea, I want spec writing to proceed from my
   confirmed conversation or a named reviewed artifact without a placeholder
   map, so that shaping is paid for only when there is a real multi-session
   decision tree. Line: `gpt-5.6-sol / high.` This changes high-leverage phase
   guidance whose semantic failures are only partly machine-observable.

7. As a Bench maintainer, I want existing active and compiled maps migrated
   without changing their decisions, so that the new validator lands green with
   no compatibility limbo or lost provenance. Line:
   `gpt-5.6-terra / high.` The mechanics are bounded, but classifying existing
   state and preserving decision meaning are not fully gate-gradeable.

8. As a workflow user, I want shaping and spec authoring to use distinct
   ownership and ticket vocabulary everywhere, so that decision tickets are not
   confused with independently-green implementation tickets and engineering
   seams are chosen in the phase that can inspect the current tree. Line:
   `gpt-5.6-sol / high.` Command, skill, glossary, and project guidance compound
   through future sessions, so the leverage override applies.

## Implementation decisions

- `internal/maps` becomes the single deep owner of map discovery, schema,
  ticket resolution, graph validation, readiness, source validation, active
  query rows, and the distinct active-map count. CLI, status, and conformance
  consume that model; none carries a second parser or re-derives a count.
- The same schema data used by the parser renders a read-only
  `bench maps --template` Markdown skeleton. `/bench-shape-idea` points to that
  skeleton instead of embedding a second exact schema. This is an explicit
  paste-ready Markdown detail mode rather than a list query: it writes the
  skeleton to stdout, writes no file, exits 0, and retains AXI usage/error
  postures for invalid argument combinations.
- A map candidate is a non-hidden direct `*.md` child of top-level
  `decisions/` or of `specs/<slug>/decisions/`; case-insensitive `README.md` is
  the standing directory-document exemption. Nested directories remain the
  place for map-owned assets and are not reclassified as maps. A missing
  candidate directory is authoritative absence. Present candidate files retain
  ADR 0010’s fail-closed empty, malformed, unreadable, wrong-type, and
  unsupported-schema states.
- The schema requires one title, one `Status: shaping` or `Status: ready`, one
  non-empty `## Destination`, at least one decision ticket, and exactly one each
  of `## Not yet specified`, `## Spec-writer discretion`, `## Out of scope`,
  and `## Sources`. `## Handoff` is unsupported in the new format.
- A decision ticket uses `## #<positive integer>: <non-empty question or title>`,
  exactly one `Blocked by: none` or comma-separated `#<id>` list, one of
  `Research`, `Prototype`, `Grill`, or `Task`, one `### Question`, and one
  `### Answer`. Headings and field markers inside fenced examples are inert.
  CRLF and a missing final newline are accepted without changing meaning.
- A ticket is unresolved when its answer is empty, begins with the existing
  open or deferred em-dash marker, or carries the existing `GRILL DEFERRED`
  banner. A deferred marker yields `deferred`; another unresolved ticket is
  `blocked` when any blocker is unresolved and `frontier` otherwise.
- The graph rejects duplicate IDs, unknown types, duplicate blocker IDs on one
  ticket, missing blocker targets, self-edges, cycles, and a resolved ticket
  whose blocker is unresolved. Diagnostics name the map, ticket title and graph
  handle, and the specific bad edge; validation records all independent
  failures it can safely derive.
- A `ready` map additionally requires every ticket resolved, literal emptiness
  under `## Not yet specified` apart from whitespace, valid structure for every
  present source, and a compiled-safe discretion/out-of-scope shape. A compiled
  map must be `ready`; an active map may be either status. A shaping map remains
  gate-valid while it carries honest unresolved tickets or fog.
- `## Sources` entries have one `Path` or `URL`, followed by non-empty
  `Supports` and `Drift` fields. Paths are repository-root-relative, cannot
  escape the root, and must resolve through any symlink to a regular file still
  inside the root. URLs must parse as absolute HTTP(S) URLs with a host. An empty
  section means the map has no required research object; it is not replaced by
  an unstructured sentinel. The validator performs no DNS lookup, HTTP request,
  or other network I/O.
- `## Spec-writer discretion`, `## Out of scope`, and shaping fog are Markdown
  bullet lists when non-empty. The validator can prove classification and
  placement, not whether a discretion item is semantically reversible.
  `/bench-shape-idea` remains the owner of the reviewer-facing rule: discretion
  is limited to technical choices that do not change observable behavior,
  scope, an architectural seam, compatibility, or what the gate proves.
- Default `bench maps` output changes to
  `maps[N]{map,title,type,state,blockers}`. These five fields are the minimal
  default because rows from multiple maps need their map identity and the
  reviewer required all four ticket facts. It emits every unresolved decision
  ticket using its title and the titles of its currently unresolved blockers;
  numeric IDs remain graph handles inside the map rather than the human-facing
  query. A shaping map with no unresolved ticket but non-empty fog emits one
  map-level row (`title=Not yet specified`, `type=fog`, `state=shaping`) so the
  detailed query cannot look empty while status counts the map. Invalid active
  candidates render explicit error rows in the same five-column schema
  alongside valid rows and make the command exit 1.
- `bench maps --count` and `maps.UnresolvedCount` count distinct active maps
  whose status is `shaping`, plus invalid active candidates. A valid ready map
  contributes zero. Compiled maps never enter the AXI listing or ambient count;
  the gate validates them separately.
- A new dev-tier conformance check named `decision-map-integrity` is a thin
  adapter over `internal/maps` tree validation. It scans active and compiled
  candidates, records every diagnostic, and is bound in the conformance
  registry. A dedicated canary family proves its failure messages still bite.
- `/bench-write-spec` records one `Decision source:` line. For a map-backed
  spec it names the compiled map path. For a no-map spec it names either the
  reviewer-confirmed current conversation with date or the reviewed source
  artifact. The line is provenance, not a second research manifest; a compiled
  map’s `## Sources` remains beside the spec.
- `/bench-write-spec` re-reads and re-verifies every map source before choosing
  seams. It asks late clarifications one at a time, recommends an answer, and
  stops after at most two if the uncertainty has not closed. A dependency tree
  or multi-session fog routes to `$bench-shape-idea`; ordinary late uncertainty
  stays in spec authoring.
- `/bench-shape-idea` owns reviewer decisions, constraints, rejected
  alternatives, exclusions, research objects, and bounded discretion.
  `/bench-write-spec` owns module seams, deep-versus-thin design, tests,
  acceptance coverage, hostile-input attachment, and the gate seam. A seam
  remains in a decision-ticket answer only when the reviewer explicitly chose
  it.
- The workflow uses **decision ticket** only for shaping questions.
  `craft-tickets` keeps **implementation ticket** as the independently-green
  build unit and its behavior does not change.
- This is a wide format migration. At build entry, `craft-tickets` owns the
  expand–migrate–contract sequence: introduce the new model without enabling a
  tree-wide refusal, migrate every current candidate while the project stays
  green, then switch CLI/status/gate consumers and remove legacy Handoff
  acceptance. No terminal state retains dual schema authority.
- Ownership fences for future build slices are:
  `internal/maps/**`; `internal/status/**` plus AXI runtime contracts;
  `internal/conformance/**` plus the new canary family;
  workflow skills/commands plus their existing guidance-anchor fixtures and
  generated indexes; and the decision-map corpus. A slice owns every file in
  its fence, and consumers block on the `internal/maps` primitive.
- The guidance pass replaces stale mandatory-map, Handoff, and map-owned-seam
  rules at their owners rather than adding exceptions around them. It updates
  README, CONTEXT, the Bench profile, relevant craft skills and command
  adapters, generated/anchored guidance, and adds one concise typed
  `CHANGELOG.md` entry.

## Testing decisions

Tests attach at three existing seams. Parser and graph cases drive the
`internal/maps` model directly. AXI and status cases drive the built `bench`
binary in throwaway repositories. Oracle cases drive the registered conformance
check and the real canary path. Guidance anchors are deletion tripwires, not
behavioral proof; the required fresh-session dogfood run is the semantic check.

The implementation gate is `.bench/gate.sh`. Focused work starts with
`go test -count=1 ./internal/maps`, the relevant AXI/status package, or
`BENCH_CONFORMANCE_CHECK=decision-map-integrity go test -count=1 ./internal/conformance -run '^TestRootConformance$'`.
The completed build must observe each new canary’s targeted red, restore it,
finish `bench gate` green, and run a fresh-session map-to-spec dogfood flow
because command and skill triggers changed.

### Seam diagram: decision-map model

    trigger: CLI, status, or conformance asks for map state
        │
        ▼
    repo root + candidate bytes ──▶ [ internal/maps model ] ──▶ maps, graph, rows, count, diagnostics
                                     ◀ tests attach here: package fixtures drive schema, graph, and sources

### Seam diagram: AXI query and ambient status

    trigger: user or SessionStart invokes the built binary
        │
        ▼
    fixture repository ──▶ [ bench maps / bench status ] ──▶ TOON stdout + exit code
                            ◀ tests attach here: runtime fixture invokes the shipped surface

### Seam diagram: decision-map gate check

    trigger: conformance phase or canary runs the oracle
        │
        ▼
    active + compiled map tree ──▶ [ decision-map-integrity check ] ──▶ diagnostics / gate verdict
                                   ◀ tests attach here: clean and mutation fixture trees

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | The parser and `bench maps --template` derive headings, status values, ticket types, and field names from one schema owner. | decision-map model | Future TDD row: introduce the schema/template equality test first; it must fail while no template flag or shared schema exists. | A handwritten template or parser-only registry cannot satisfy the same-owner equality assertion. |
| 1 | A well-formed shaping map parses while missing status, duplicate required sections, missing ticket fields, Handoff, and an unsupported type each produce distinct diagnostics. | decision-map model | Future `go test -count=1 ./internal/maps -run '^TestMapSchema'`; each malformed table case must be observed red before its parser rule. | It closes the cheapest implementation that recognizes only `Status:` and otherwise accepts prose. |
| 1 | Research, Prototype, Grill, and Task each parse successfully as the ticket’s type. | decision-map model | Future four-case positive table; a Grill-only parser must first fail the other three cases. | It closes the inverse type bug: rejecting `Bogus` does not prove every allowed value is accepted. |
| 1 | Direct active and compiled candidates are discovered, while hidden Markdown, case variants of README, and nested asset Markdown are ignored. | decision-map model | Future discovery matrix; recursive or no-exemption scans must first fail on the controlled ignored files. | It pins both positive locations and every compatibility exemption instead of testing only ordinary candidates. |
| 2 | Duplicate IDs, duplicate blocker entries, dangling blockers, self-edges, cycles, and resolved-on-unresolved edges are each rejected. | decision-map model | Future `go test -count=1 ./internal/maps -run '^TestMapGraph'`; add and observe each failing case before graph behavior. | A parser that stores blocker text without constructing a complete graph passes none of these mutations. |
| 2 | Every unresolved ticket is emitted as frontier, blocked, or deferred; blocked rows carry unresolved blocker titles, while numeric IDs remain graph handles only inside the map. | AXI query and ambient status | Future AXI graph-route fixture; it must first fail against the current `map,ticket,type,state` output. | A frontier-only query, an ID-only query, or a state guessed from file order has observably different rows. |
| 2 | A shaping map with resolved tickets but non-empty fog remains visible as one map-level shaping row. | AXI query and ambient status | Future zero-open/non-empty-fog AXI fixture; it must first fail with today’s empty listing. | It prevents the degenerate implementation that counts unresolved fog but renders no explanation. |
| 3 | Marking the same honest-fog map `shaping` passes and `ready` fails with a non-empty-fog diagnostic. | decision-map model | Future paired readiness test; the ready half must be observed red before the status predicate exists. | The paired body isolates status as the only changed fact and proves readiness is not inferred from ticket markers alone. |
| 3 | A ready map rejects every unresolved/deferred ticket and accepts only resolved tickets with literal empty fog. | decision-map model | Future readiness table test, first red on each marker class and non-whitespace fog body. | It catches a validator that checks only the status token or only one legacy marker spelling. |
| 3 | Structured sources accept an in-root regular path and absolute HTTP(S) URL, reject empty locators, absent/special/escaping paths, and malformed URLs, and never access the network. | decision-map model | Future source-validation test with a dead hostname and hostile filesystem cases; each invalid case must first red while the current parser ignores sources. | It proves local evidence exists, repository scope holds, and URL validation remains syntactic/offline. |
| 3 | A ready all-grilled map with an empty Sources body passes, while an unstructured `n/a` sentinel fails. | decision-map model | Future paired source-cardinality test; first red if the parser either requires an entry or accepts the sentinel form. | It pins the defaulted zero-source posture without weakening the one structured entry grammar. |
| 3 | The validator requires discretion, exclusions, and fog to use their declared list shape while leaving discretion semantics to reviewer sign-off. | decision-map model plus guidance conformance | Future structural table test and a guidance-anchor mutation; both must first fail before the list and policy checks exist. | It prevents unclassified prose while avoiding the false claim that Go can decide architectural reversibility. |
| 4 | The registered check accepts valid active and compiled maps and a spec with no decisions directory. | decision-map gate check | Future focused conformance fixture, first red until the new check is registered and bound. | It proves both allowed map locations and the situational no-map path, rather than a blanket “map required” check. |
| 4 | An invalid compiled map makes the focused check and full gate red with its compiled path and targeted diagnostic. | decision-map gate check | Future compiled-map canary; observe its own diagnostic red, then restore the valid fixture. | A check that scans only top-level `decisions/` cannot see this mutation. |
| 4 | Each schema, graph, readiness, and source refusal has a canary expectation matching that refusal’s own diagnostic. | decision-map gate check | Future dedicated canary family; deleting or weakening any validator arm must make `bench canary` red. | Per-refusal bites prevent one generic failing fixture from laundering untested branches. |
| 5 | Listing and `--count` derive from one active scan: two unresolved tickets in one shaping map count once, one invalid candidate counts once, and a valid ready map counts zero. | AXI query and ambient status | Future count/listing contract fixture; it must first fail against current marker/Handoff counting. | It catches ticket-counting, ready-map counting, and silent invalid-file omission in one controlled tree. |
| 5 | `bench status` reports the same distinct active-map count and degrades to unknown when the active directory scan fails. | AXI query and ambient status | Future runtime status fixture paired with `bench maps --count`; first red if either surface disagrees. | A second status derivation or fabricated zero cannot satisfy both outputs on the same root. |
| 5 | Compiled-map failures never enter `bench maps` rows or the ambient count even though they red the gate. | both high seams | Future paired fixture invoking maps/status and the focused conformance check on one bad compiled map. | It pins the deliberate active-query versus all-map-oracle split. |
| 1, 5 | Repeated default, count, template, and focused-validator invocations leave tracked and untracked repository state unchanged. | AXI query and ambient status | Future before/after fixture snapshot; first red when a deliberately mutating template/cache stub changes the tree. | Formatting and validation rows cannot detect a query that writes correct output and also mutates state. |
| 6 | `/bench-write-spec` accepts a ready map, a reviewer-confirmed current conversation, or a named reviewed artifact and records exactly one `Decision source:` line. | guidance conformance | Future workflow-anchor clean/mutation tests; removing any one authorization or the provenance line must red its targeted fixture. | It rejects both old mandatory-map behavior and an unbounded “conversation alone” bypass. |
| 6 | Late uncertainty is asked one question at a time with a recommendation, is bounded to two, and routes multi-session dependency fog to `$bench-shape-idea`. | guidance conformance plus fresh-session dogfood | Future anchor mutations cover the structural rule; a fresh session must exercise one clarification and inspect the produced source line. | An unbounded interview, bundled questions, or silent spec-writer decision violates an observable phase step. |
| 7 | Every current active candidate migrates to shaping or ready according to its actual state, every compiled candidate is ready, and no migrated map retains Handoff or loses a ticket answer, exclusion, or source fact. | decision-map gate check plus artifact review | Future migration inventory test and full-tree focused check; current legacy corpus must first fail once the new validator is enabled. | It catches migrating only files listed by `bench maps`, blanket-ready status, and lossy section deletion. |
| 7 | The feature reaches a terminal state with no legacy-schema acceptance path or second parser. | decision-map model plus conformance | Future source-level ownership test and old-format rejection case must first red while compatibility remains. | A permanent dual parser would otherwise make the migrated tree green while preserving drift. |
| 8 | Shaping guidance uses decision-ticket vocabulary, points to `bench maps --template`, removes Handoff, and makes maps situational. | guidance conformance | Future targeted workflow canaries; each deleted or reverted rule must red with its own message. | It catches the cheap prose edit that changes the command headline but leaves old mandatory-map and Handoff rules below. |
| 8 | Spec guidance owns seams/tests/coverage/gate attachment and consumes map Sources without copying a research manifest. | guidance conformance plus fresh-session dogfood | Future anchor mutations plus one fresh-session map-backed spec; first red until stale Handoff-derived-seam text is removed. | It catches phase-ownership drift across command and craft-skill surfaces. |
| 8 | `craft-tickets` retains independently-green implementation-ticket semantics while CONTEXT, README, profile, and changelog consistently describe decision tickets. | guidance conformance | Future vocabulary sweep and consistency canary; first red on the current mandatory-map/Handoff vocabulary. | It prevents a global “ticket” rename from changing the build unit or leaving contradictory cold-start guidance. |

### Degenerate-implementation check

- A status-only parser fails stories 1–3’s malformed schema, graph, and paired
  readiness rows.
- A graph parser that emits only the frontier fails the blocked/deferred route
  and fog-visibility rows.
- A gate adapter that re-parses Markdown or scans only top-level maps fails the
  ownership and compiled-map rows.
- A source validator that calls `url.Parse` but never checks local targets fails
  absent, special, escape, and dead-host cases.
- A count derived from rendered ticket rows fails the distinct-map and
  zero-ticket-fog cases.
- A prose-only headline change fails the stale-rule canaries, vocabulary sweep,
  and fresh-session dogfood.
- A migration that merely adds `Status: ready` fails ready-map validation and
  the preserved-fact inventory.

### Quantifier enumeration

- Map locations are exactly top-level `decisions/*.md` and direct
  `specs/<slug>/decisions/*.md` candidates, with hidden names and README exempt
  and nested asset directories excluded.
- Valid ticket types are exactly Research, Prototype, Grill, and Task.
- Graph refusals are duplicate ticket ID, duplicate blocker in one ticket,
  dangling blocker, self-edge, cycle, and resolved-on-unresolved blocker.
- Unresolved states are open, deferred, and GRILL DEFERRED; query states are
  frontier, blocked, and deferred.
- Source locator classes are repository Path and absolute HTTP(S) URL; source
  fields are locator, Supports, and Drift.
- Guidance owners swept are `/bench-shape-idea`, `/bench-write-spec`,
  `craft-grill`, `craft-spec`, `craft-delegate`, README, CONTEXT, the Bench
  profile, applicable workflow canaries/indexes, and CHANGELOG.
- The migration inventory is every non-exempt active and compiled candidate
  present at implementation time, not the filenames listed in this draft.

### Edge inventory

| edge class | resolution |
|---|---|
| error path | The schema, graph, readiness, source, and compiled-path rows require distinct diagnostics while continuing across independent candidates. |
| empty/absent input | The schema, readiness, source-cardinality, and no-map rows distinguish absent directories, zero-byte candidates, empty answers, empty fog, empty source sections, and a legal spec with no map. |
| boundary values | Schema and graph rows cover the smallest positive ID, duplicate IDs, no blockers, multiple blockers, one-ticket maps, and a shaping map with only fog. |
| malformed input | Schema, graph, and source rows cover malformed headings, types, edges, bytes, paths, and URLs. |
| interrupted or partial state | A half-migrated tree remains green only before the contract step; once the check is enabled, any legacy or partial candidate reds. The build uses `craft-tickets`’ wide-refactor sequence. |
| re-run idempotency | Re-running the read-only query/template and validator changes no file; a second migration pass finds no legacy candidate and changes no decision content. |
| hostile environment | Source, status, and read-only rows plus the checklist below cover spaces/globs, special files, symlinks, nested cwd, malformed bytes, and offline execution. |

### Hostile-input checklist disposition

| checklist class | resolution |
|---|---|
| spaces or glob characters in paths | Source and candidate fixtures use repository roots and filenames containing both; discovery and diagnostics preserve the path. |
| refused control bytes in git-sourced text | Candidate titles and source fields pass through TOON; existing `toon.TableTyped` refusal remains the renderer owner and the AXI fixture asserts exit 1 rather than raw control output. |
| permitted tab/newline/return splitting a line sink | Structured fields are line-scoped and reject embedded control line breaks; TOON escaping owns rendered ticket titles. |
| command self-write changes reported fact | `bench maps`, `--count`, and `--template` are read-only; a tracked fixture compares status before and after repeated invocations. |
| hand-edited file lacks trailing newline | Parser and AXI fixtures cover a valid map and one malformed field without a final newline. |
| absent versus present-but-empty | Missing candidate directories are valid absence; zero-byte candidate files and empty required sections are distinct diagnostics. |
| special discovered files | FIFO, device, and socket candidates and local source targets are rejected before open; capability-aware tests prove no blocking read. |
| dangling symlink | Candidate and source-path dangling links diagnose unreadable; live source links are accepted only when their resolved regular target remains inside the repository. |
| unquoted multi-word arguments | Runtime fixtures invoke the query from a repository path containing spaces and glob characters. |
| flag value mistaken for positional | `usage.Parse` owns `--count`, `--template`, `--`, unknown flags, and arity; existing grammar contracts extend to the new flag. |
| required tool missing from PATH | The built Go query and conformance check require no external parser or network tool; missing `bench` is outside the invoked binary seam. |
| invocation through a symlink | Existing CLI routing remains the owner; the AXI query contract exercises the shipped built surface rather than adding map-specific path resolution. |
| every shipped invocation surface | Direct kit and linked-repository `bench maps` reach the same Go command through existing routing; hooks consume status, which consumes the same model. |
| destructive worktree state | **Won’t handle:** map validation is read-only and owns no worktree mutation; the existing worktree subsystem remains the callable owner. |
| SIGINT mid-loop | **Won’t handle:** validation performs bounded local reads with no loop, subprocess, or writes; gate process-group teardown remains a separate capability. |
| re-run idempotency | Repeated default, count, template, focused-check, and full-gate invocations preserve tracked state and stable semantic output. |
| cwd deeper than repo root | Runtime fixtures invoke maps and status below the root and require identical discovery. |
| non-TTY stdin | Query, template, status, and conformance are non-prompting and never read stdin. |
| host-backed filesystem pressure | **Won’t handle:** deterministic latency guarantees are a separate bounds capability; bounded record size and classification-before-open prevent unbounded content and FIFO waits. |

## Out of scope

- **Prototype-ticket lifecycle and prototype isolation** — a separate
  artifact-lifecycle capability with reviewer interaction and cleanup rules —
  7 edits, 3 gate runs.
- **Domain-model maintenance during shaping (FT165)** — a separate ubiquitous
  language and ADR-maintenance capability already owned by the roadmap —
  5 edits, 2 gate runs.
- **Tracker-backed maps, assignment claims, and one decision ticket per
  session** — a separate collaborative-planning substrate that reopens rejected
  single-reviewer workflow decisions — 12 edits, 4 gate runs.
- **Low-resolution map indexes with one file per decision ticket** — a separate
  storage/context-pressure capability with migration and concurrency semantics —
  9 edits, 3 gate runs.
- **Changes to implementation-ticket slicing or `craft-tickets` semantics** — a
  separate build-orchestration capability; this spec only invokes its existing
  wide-refactor rule — 4 edits, 2 gate runs.
- **Network reachability, freshness, or trust verification for external
  sources** — a separate bounded-network and provenance capability; this gate
  remains hermetic and syntax-only — 6 edits, 3 gate runs.
