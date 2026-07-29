# The slice unit — tickets sized to one context window

## Destination

A slice unit ("ticket") that fits one fresh context, baked into
`/bench-implement-spec` as its breakdown step and reachable standalone so a
small unspecced change is the one-ticket degenerate case (FT154; the slicing
half of FT107's light-path boundary).

## #1: What is a ticket, observably?

Type: Grill

### Question

Bench has one sizing axis (the ownership fence — who writes where) and lacks
two: independent-green (can this land committed on a green gate by itself) and
context fit (can one session hold it). "Fits one context window" is
unfalsifiable as a rule. What defines the unit, and which axis is the
observable that grades it? Candidate: the smallest story group that can land
committed green on its own, with context fit as the sizing heuristic rather
than the rule. Everything else in this map keys off this definition.

### Answer

Independent-green defines the unit: a ticket is the smallest story group that
can land committed on a green gate by itself. The gate grades it — a commit
either lands green or the group wasn't a ticket. Context fit is the sizing
heuristic that says "split further"; it never grades anything. The ownership
fence stays orthogonal: fences say who writes where, tickets say what lands
green next. Verticality is part of the definition, not inherited from phase
prose (reviewer-confirmed 2026-07-28): a ticket is a tracer bullet — a small,
tightly related scope cutting a narrow but complete path through the layers,
demoable or verifiable on its own — so a horizontal green grouping (one
layer, or tests without behavior) does not qualify however green it lands.
"Fits one context" is the consequence of small-and-cohesive, never the rule.

## #2: Which file owns the breakdown procedure?

Blocked by: #1
Type: Grill

### Question

Reviewer-decided 2026-07-28: the breakdown lives inside
`/bench-implement-spec` and doubles as the light route. Open remainder:
`craft-spec` already owns "slicing a build for delegates" (the ownership
fence), so a new `craft-tickets` skill would be a second slicing source — the
knowledge-duplication defect. Does the ticket unit extend `craft-spec`'s
slicing section (one slicing owner, charged from the phase), or earn its own
skill because spec-time slicing (who writes where) and build-time slicing
(what lands green next) are genuinely different facts?

### Answer

A new `craft-tickets` skill. The audiences are different facts: `craft-spec`
fires at authoring time and owns who-writes-where; `craft-tickets` fires at
build entry and on unspecced changes and owns what-lands-green-next.
`/bench-implement-spec` charges it as its breakdown step; the light path
invokes it bare. Each skill points at the other's rule by name rather than
restating it — the fence stays single-sourced in `craft-spec`, the ticket
unit single-sourced in `craft-tickets`.

## #3: Where do tickets live, and what makes a reset safe?

Blocked by: #1
Type: Grill

### Question

`to-tickets` writes `tickets.md`; Bench already carries the spec, the coverage
map, and `session-handoff.md`. A third enumeration of the same work is the
knowledge-duplication defect, but a context reset is only safe if the next
ticket is resumable from its own text (invariant 3 at slice granularity).
Candidate: derive tickets at build entry from the spec's stories/coverage
rows, and record only the frontier — done / in-flight / next — in
`session-handoff.md`'s existing State section.

### Answer

Per-spec folders, modeled on how `to-spec` and `to-tickets` interact
(spec-as-parent, tickets-as-children): a spec becomes `specs/<slug>/spec.md`
and its breakdown lives beside it in `specs/<slug>/tickets/`. Each ticket
carries what-to-build as end-to-end behavior, its acceptance criteria, and
`Blocked by:` edges naming sibling tickets; the parent is the spec it sits
beside. The frontier is any ticket whose blockers are done, worked one at a
time with context cleared between tickets. Reviewer-decided 2026-07-28: the
duplication cost is accepted — the ticket is a working document, the spec the
contract.

## #4: Which phase writes the ticket files, and who approves them?

Blocked by: #1
Type: Grill

### Question

With tickets as files under `specs/<slug>/tickets/` (#3), someone must write
them and someone must approve the breakdown — `to-tickets` is its own
user-quizzed step between spec and implement. Reviewer-stated shape
(2026-07-28, pre-map) put the breakdown inside `/bench-implement-spec`.
Candidate: build entry writes `tickets/` as its first act, deriving edges from
the spec's stories, seams, and shared-primitives-first rule, and proceeds
under the session's existing approval (batch approval covers it AFK, sign-off
otherwise); no separate ticketing phase. The unspecced light path writes one
ticket the same way.

### Answer

Build entry writes them: `/bench-implement-spec`'s first act derives tickets
from the spec's stories and seams, writes `specs/<slug>/tickets/`, and
presents the breakdown as the build plan — approval rides the session's
existing surface (sign-off when present, batch approval AFK, tickets as
post-hoc veto surface like the spec itself). No new phase. The unspecced
light path invokes the same procedure and writes one ticket.

## #5: What is the reset, mechanically?

Blocked by: #1
Type: Grill

### Question

The phase already routes authorship through worktree-isolated write delegates
— each delegate is a fresh context, so a reset exists today wherever a ticket
becomes a delegate charge. Is the rule simply "one ticket = one delegate
charge" (sequential tickets = sequential fresh delegates; parallel only where
fences allow), or does the orchestrating session also need a reset story —
the `--full` run accumulating context across all fifteen stories is the
observed instance (FT154's second source), and a story-group reset inside
`--full` would be this decision applied there?

### Answer

One ticket = one write-delegate charge: the reset is the delegate boundary
the phase already has, and the orchestrator holds only the spec and the
frontier — which is also the `--full` answer, since story groups become
tickets. Done-marking is checkboxes: a landed ticket stays in
`specs/<slug>/tickets/` with its acceptance boxes checked until the spec
retires (the whole folder leaves together); the coverage map remains the
gate-facing done authority, the boxes are build-progress visibility.

## #6: Fold FT107's light-path table into this, or sequence behind it?

Blocked by: #1, #2
Type: Grill

### Question

FT107's first clause bounds the light path by blast radius; this map would
make "decomposes to one ticket" the observable approximating it. One boundary
decided from two sides must land consistent. Fold the table into this map's
build, or keep FT107 separate and sequence it behind the unit decision?

### Answer

Fold the table in: the light-path table lands in this map's build, with
"decomposes to one independently-green ticket that crosses no declared seam"
as its observable for blast radius. Only the table moves here — FT107's other
clauses (read-budget rerouting, fix-loop escalation, squash-merge commit
safety, self-contradicting-spec rule) stay on FT107.

## #7: How does the spec-folder layout land against the CLI?

Blocked by: #3
Type: Grill

### Question

`internal/spec` hard-codes the `specs/<slug>.md` shape — resolution
(`specs/<slug>.md` fallback for bare slugs), `Facts` reading `specs/*.md`,
retire validation, `history`'s `:(literal,top)specs/<slug>.md` pathspec — and
`bench roadmap`/`bench status` cross-check the same form, so #3's layout is a
Go change with a migration posture to pick: teach the CLI both forms and
migrate specs as they're touched, or migrate the existing staged spec(s) in
one commit and support only the folder form. Also to decide: whether
`bench spec` gains a tickets-aware surface (list the frontier) or the files
stay convention-only in v1.

### Answer

Folder-only: teach `internal/spec` the `specs/<slug>/spec.md` form and keep only
that shape going forward — no standing dual-form resolution. No migration is
needed: `specs/` holds no staged spec, so the folder form is the only live form
from the change forward. `bench spec history` stays able to see both forms,
because retired specs live at their old flat paths in git history and its
pathspec must keep resolving them. The roadmap/status
cross-check and the roadmap preamble's stated path convention move with the
layout in the same change. Ticket files are convention-only in v1 — the CLI
learns the folder, not the ticket format; a `bench` frontier surface
graduates later on demonstrated use, the FT6 posture.

## #8: What line does each stage of a ticketed build run?

Blocked by: #1
Type: Grill

### Question

With tickets small enough to charge cheap, does the kit bind a standing tier
per stage — orchestration, implementation, review — instead of per-build
deliberation, and how hard is the binding?

### Answer

Defaults with the ladder, recorded as a standing table in `craft-line`:
orchestration mid; ticket implementation cheap at low effort
(reviewer-amended 2026-07-28 at spec time from the original medium); review
mid at high. Each is the default charge, not a flat rule — a failed done-claim or
a red the delegate cannot clear escalates one tier, declared per the existing
ladder, never silently. The profile's leverage override stands: kit
always-loaded prose still routes top, for review and build alike. Stage is
not machine-visible to `check-agent-line` (the envelope carries no stage
field — FT128's gap class), so the binding is prose plus review in v1.
Reviewer-decided 2026-07-28: this replaces "builds route mid" for
ticket-sized charges, and the first builds under it are the cheap-tier
re-test evidence `decisions/cost-follows-project-size.md` #6 waits on.

## #9: Where does a closed decision map live after its spec compiles?

Blocked by: #3, #7
Type: Grill

### Question

The folder layout gives a compiled map a natural home beside its spec, but
`bench maps` currently treats top-level `decisions/` as the shaping frontier.
Should a closed source map remain top-level until final-check deletes it, be
copied beside the spec, or move into the spec folder with a defined asset,
reference, query, and retirement lifecycle?

### Answer

Pre-spec working maps stay top-level under `decisions/` through shaping and
until compilation. When
`/bench-write-spec` compiles a closed map, it moves rather than copies the
source map and any map-owned assets into `specs/<slug>/decisions/` and updates
every moved-path reference in the same green change. The spec-local files are
settled provenance, so `bench maps` continues to scan only top-level
`decisions/`. `bench spec retire` removes the whole spec folder, including the
compiled map, its owned assets, and tickets; final-check does not separately
delete a top-level shipped map.

- Gate anchors for the new prose (which conformance assertions pin the
  breakdown step and the light-path table), decided at spec time.
- The exact ticket-file template wording and the light-path table's wording —
  drafting, at spec time.
- The fail posture when one slug carries both the flat and folder form, and
  for a folder without `spec.md` — edge-inventory work at spec time.

## Out of scope

- The wide-refactor exception — expand–contract already lives in `craft-spec`
  and is mirrored in `to-tickets`; this map points at it, never copies it.
- Tracker publication (`to-tickets`' GitHub/Linear route) — Bench has no
  tracker; tickets never become external issues here.
- Retesting delegate tiers — fence and unit rules stay tier-independent.
- A `bench` ticket-parsing surface (frontier listing, edge validation) — held
  to the FT6 graduate-on-evidence posture, not built in v1.
- FT107's non-table clauses — they stay on FT107.

## Handoff

1. **Module boundaries.** `craft-tickets` skill (unit definition, breakdown
   procedure, ticket-file template, frontier rule) — new
   `.agents/skills/bench-craft-tickets/` plus its `.claude/skills/` symlink
   entry (the mirror is per-skill symlinks, not copied files).
   `/bench-implement-spec` gains the breakdown step charging it.
   `internal/spec` learns the `specs/<slug>/spec.md` layout, with no spec to
   migrate. `.bench/BENCH.md` gains the light-path table with the one-ticket
   observable. `craft-line` gains the per-stage
   default table (#8). `/bench-write-spec` owns the map move into the spec's
   `decisions/` folder (#9). Outside: no ticket parser, no new subcommand,
   no tracker.
2. **Contracts.** `internal/spec`: a bare slug resolves to
   `specs/<slug>/spec.md`; `Facts` enumerates folder specs; retire validates
   and removes the folder; `history` keeps resolving retired flat paths.
   Fail postures for flat/folder collisions are spec-time edge work. Ticket
   file (convention, unparsed): title, What to build (end-to-end behavior),
   `Blocked by:` naming sibling titles, acceptance checkboxes checked as work
   lands. Skill contract: one ticket = smallest independently-green story
   group = one write-delegate charge. Decision-map lifecycle: pre-spec working
   maps stay top-level; compilation moves the closed map and map-owned assets
   into `specs/<slug>/decisions/` with same-change reference updates;
   `bench maps` ignores the settled provenance and whole-folder retirement
   removes it.
3. **Deep vs thin.** `internal/spec`'s resolution is the deep unit — every
   consumer (coverage, commit --spec, retire, status, roadmap) sees the
   layout only through it. The phase's breakdown step is thin: it charges the
   skill, which owns the procedure.
4. **Black-box assertables.** `bench spec` resolution/retire/history and
   `bench coverage`/`bench status`/`bench roadmap` against a folder-form spec
   (exit codes, TOON rows naming the folder path). Conformance assertions
   that the breakdown step, the light-path table, and the skill's charge line
   are present in the owning prose files. `bench maps` reports no row for an
   open-looking map parked under a spec, and retirement leaves no spec-local
   decision provenance behind.
5. **Gate attachment.** The layout change is ordinary Go under existing
   test/contract phases; prose anchors ride the conformance phase like FT152's
   family. The later ticket-guidance conformance slice owns the decision-map
   lifecycle anchors and biting canary; existing AXI/runtime contracts cover
   the query and retirement behavior. Not gate-visible: ticket quality
   (sizing, edge correctness) — it is reviewed at breakdown approval, and
   that gap is accepted in v1.
6. **Hostile-input owners.** Flat/folder collision and folder-sans-spec.md →
   `internal/spec` resolution, fail closed naming the conflict. Malformed
   ticket file → n/a in v1, nothing parses it. Kit-versus-linked-repo (FT144
   prompt): linked repos own their own `specs/`; the layout and skill ship as
   kit surface, so the edge inventory must walk both audiences at spec time.
7. **Uncertainty flags.** None blocking. The ticket-file template is
   deliberately unparsed so it can evolve; FT107-table wording is drafting.
8. **Rejected alternatives.** Extending `craft-spec` instead of a new skill
   (#2). `tickets.md` single file, handoff-frontier-only, and in-session-only
   artifacts (#3). A separate ticketing phase, and spec-phase ticketing (#4).
   Delete-on-land done-marking (#5 — checkboxes chosen). Context fit as the
   grading rule (#1). Dual-form resolution and convention-only layout (#7).
   A frontier CLI surface in v1 (#7). A flat no-escalation stage rule, and
   dropping the review leverage override (#8). Leaving compiled maps
   top-level, copying them into specs, scanning spec-local provenance with
   `bench maps`, and separate final-check deletion (#9).
9. **Domain watch-outs.** `bench spec history` resolves retired specs by
   literal git pathspec, so retired flat paths must stay reachable after the
   layout moves. Any `specs/*.md` glob (in `Facts`, gate prose, or
   conformance) silently misses folder specs — an invisible-skip class, so
   each glob is enumerated and moved, not patched as found. A map move is one
   green change with every map-owned asset and exact path reference; splitting
   those edits creates either duplicate authority or a dangling reference.

Dependency order within this map's build: the CLI layout change lands before
the compiled-map move; the skill/phase/table prose follows the lifecycle it
names. Slicing within that is the reviewer's call.

## Sources

Read in full for this map, 2026-07-28. The reference repo is vendored at
`~/workspace/reference-skill-repos/skills`; paths under it may drift with the
vendor, tree paths with the tree — re-verify before citing onward.

- `skills/engineering/to-tickets/SKILL.md` (vendored) — the ticket template,
  frontier rule, and wide-refactor expand–contract; the model for #1, #3, #4.
- `skills/engineering/to-spec/SKILL.md` (vendored) — the spec-as-parent /
  tickets-as-children interaction #3's folder layout mirrors.
- `skills/engineering/implement/SKILL.md` and `skills/engineering/tdd/SKILL.md`
  (vendored) — the loop rules; confirmed they mirror what `craft-tdd` and the
  phase already carry, so nothing was borrowed from them.
- `.agents/commands/bench-implement-spec.md` — the venue routing whose
  write-delegate boundary is #5's reset, and the file the breakdown step edits.
- `.agents/skills/bench-craft-spec/SKILL.md` — the ownership fence #1 keeps
  orthogonal and the expand–contract rule Out of scope points at.
- `internal/spec/spec.go`, `internal/spec/history.go` — the flat-path
  consumers #7 enumerates (resolution fallback, `Facts` glob, retire
  validation, history's literal pathspec).
