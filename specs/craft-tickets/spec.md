# craft-tickets — the slice unit, the folder spec layout, and the light path

Status: staged
Roadmap: FT154

Compiled from `decisions/slice-unit.md` (closed 2026-07-28, all eight decisions
resolved, no blocking uncertainty flags). Every seam below is map-sourced;
deviations from the map are flagged inline as **[flagged]**.

## Problem

Bench has one sizing axis (the ownership fence — who writes where) and no unit
for what lands green next. Builds either run whole-spec in one accumulating
context or get sliced ad hoc mid-loop; small unspecced changes have no defined
light route, so every few-line change either over-runs the full pipeline or
skips it silently. Separately, spec artifacts are single flat files with
nowhere for a build breakdown to live beside them.

## Solution

A ticket: the smallest story group that can land committed on a green gate by
itself — the gate grades it; context fit is the sizing heuristic that says
"split further", never the rule; verticality (a tracer bullet through the
layers, demoable or verifiable on its own) is part of the definition. A new
`craft-tickets` skill owns the unit; `/bench-implement-spec`'s first act writes
the breakdown to `specs/<slug>/tickets/` beside the spec, which moves to
`specs/<slug>/spec.md` (folder-only — no dual-form resolution); the unspecced
light path is the one-ticket degenerate case, bounded by a standing table in
`.bench/BENCH.md`; one ticket = one write-delegate charge; `craft-line` gains
the per-stage default line table. `internal/spec` is the deep unit every
consumer sees the layout through.

## User stories

Go — the folder layout (the deep unit lands first; map dependency order):

1. As an agent, I want a bare spec slug to resolve to `specs/<slug>/spec.md`
   anchored at the repo root, so that `bench coverage`, `bench commit --spec`,
   `bench spec implemented|retire|history` all address folder specs from any
   cwd. Line: gpt-5.6-luna / low. This matches the stage table's cheap-at-low
   default for ticket implementation (reviewer-amended 2026-07-28), the ladder
   corrects it if the gate disagrees, and this story is precise, well-gated Go
   at a known seam.
2. As an agent, I want resolution to fail closed when a flat `specs/<slug>.md`
   exists — alone ("flat spec layout: move to `specs/<slug>/spec.md`") or
   beside a folder (a collision naming both paths) — so that a stray or
   pre-upgrade flat spec is a named refusal, never an invisible skip. Line:
   gpt-5.6-luna / low. The refusal postures are enumerated in this spec and
   contract-tested, so story 1's reviewer-set row applies here too.
3. As an agent, I want a `specs/<slug>/` folder without `spec.md` to fail
   closed at resolution naming the missing file, so that a half-deleted or
   half-created spec folder cannot read as an authoritative empty state. Line:
   gpt-5.6-luna / low. Same seam, same enumerated posture, and therefore
   the same line as story 2.
4. As an agent, I want `spec.Facts` (feeding `bench status` and `bench
   handoff`) to enumerate `specs/*/spec.md` with slug = folder name, retaining
   evidence rows for malformed content, so that ambient surfaces see folder
   specs everywhere they saw flat ones. Line: gpt-5.6-luna / low. This is a
   mechanical consumer move behind the deep unit, gate-covered by existing
   suites re-pointed to folder fixtures.
5. As a reviewer, I want `bench spec retire` to validate merged-implemented at
   HEAD, remove the review pickup and the whole `specs/<slug>/` folder
   (tickets and all), and give every interrupted intermediate state the
   defined re-run outcome enumerated under Implementation decisions, so that a
   retired spec leaves nothing behind and no partial state is undefined. Line:
   gpt-5.6-luna / low. This is the riskiest Go in the diff, but every
   posture is enumerated below and contract-tested, which is what makes the
   reviewer-set cheap-at-low row workable; a red the delegate cannot clear
   escalates per the ladder.
6. As a reviewer, I want `bench spec history <slug>` to keep resolving specs
   retired at the old flat path *and* find folder-form deletions, so that git
   stays the archive across the layout change. Line: gpt-5.6-luna / low.
   This is two literal pathspecs through the existing one-parser design.
7. As a reviewer, I want `bench status` to count folder specs awaiting
   retirement, pair `reviews/<slug>.md` orphans against folder slugs, and
   cross-check ROADMAP rows for `specs/<slug>/spec.md` tokens, so that the
   ambient board stays truthful under the new layout. Line: gpt-5.6-luna /
   low. These are consumer moves with existing fail-closed status tests
   re-pointed to folder fixtures.
8. As a reviewer, I want `bench roadmap`'s context parse and `bench handoff`'s
   live-spec rendering to speak the folder path, so that no ambient surface
   prints a path form the CLI no longer resolves. Line: gpt-5.6-luna / low.
   These are small mechanical consumer edits behind the same deep unit.
9. As a maintainer, I want the conformance coverage sweep to validate
   `specs/*/spec.md` maps and to turn red on any stray flat `specs/*.md` in
   this repo, so that the kit repo cannot silently carry the dead form the CLI
   refuses. Line: gpt-5.6-terra / medium. This is gate/conformance logic, and
   the profile's cached routing keeps the oracle's own checks on the mid tier
   — a deliberate deviation from stories 1–8's cheap default. **[flagged]**

Prose and skill — lands after the Go slice (map dependency order):

10. As a build session, I want a `craft-tickets` skill
    (`.agents/skills/bench-craft-tickets/SKILL.md` + its `.claude/skills/`
    symlink + skills-index regeneration) owning the unit definition, the
    breakdown procedure, the ticket-file template, the frontier rule, and the
    one-ticket-one-delegate reset, so that build-time slicing has one source.
    Line: gpt-5.6-sol / high. Skill authoring is the leverage override —
    guidance prose compounds through every session that loads it.
11. As a build session, I want `/bench-implement-spec`'s first act to derive
    tickets from the spec's stories and seams, write `specs/<slug>/tickets/`,
    and present the breakdown as the build plan under the session's existing
    approval surface, so that every build starts from an approved frontier.
    Line: gpt-5.6-sol / high. Phase-command prose is the same leverage
    override.
12. As a reviewer, I want the light-path table in `.bench/BENCH.md` with
    "decomposes to one independently-green ticket and crosses no declared
    seam" as its observable, so that small changes have the standing rule the
    workflow section has always licensed but never stated (FT107's table
    clause folds in here; its other clauses stay on FT107). Line: gpt-5.6-sol
    / high. Always-loaded platform prose is the leverage override at its
    strongest.
13. As a coordinator, I want `craft-line` to carry the per-stage default table
    — orchestration mid/medium, ticket implementation cheap/low, review
    mid/high — with the ladder and the kit-prose leverage override intact, so
    that ticketed builds inherit routing instead of deliberating per build.
    Line: gpt-5.6-sol / high. Same leverage override.
14. As a spec author, I want `craft-spec` and `craft-tickets` to point at each
    other's rule by name (fence: spec-time who-writes-where; ticket: build-time
    what-lands-green-next) without restating it, so that slicing keeps one
    owner per fact. Line: gpt-5.6-sol / high. Cross-pointer prose in two
    skills.
15. As a cold session, I want every kit-prose statement of the spec path
    convention moved to the folder form in the same change — the
    `/bench-write-spec` template and retire step, `/bench-final-check`'s
    retirement-detector line, `craft-delegate`'s charge example, the ROADMAP
    preamble's stated convention, the README walkthrough, and the field
    guide — so that no document teaches a path the CLI refuses
    (enumerate-every-glob applies to prose too). Line: gpt-5.6-sol / high.
    Sweep of always-loaded and shipped prose; cheap edits, compounding reach.
16. As a maintainer, I want conformance anchors plus canary fixtures pinning
    the new prose — the breakdown step, the light-path table, the stage
    table, and the skill's charge line — so that the guidance cannot silently
    regress out of the files that carry it. Line: gpt-5.6-terra / medium.
    Anchor plumbing follows the existing workflow-guidance-anchors family
    mechanically.

## Implementation decisions

- **Folder-only resolution.** `internal/spec.Resolve` keeps path-first
  behavior for explicit path arguments; the separator-free fallback becomes
  `specs/<slug>/spec.md`. No standing dual-form support (map #7, closed).
  There is no migration: `specs/` was empty when this spec was staged, and
  this spec is the first folder occupant.
- **Fail postures (map left these to spec time; drafted here, veto surface):**
  bare-slug resolution refuses flat-only (`specs/<slug>.md` exists → exit 1
  naming the file and the folder form), collision (both forms → exit 1 naming
  both), and folder-sans-`spec.md` (exit 1 naming the missing file).
  Precedence is deterministic: any flat `specs/<slug>.md` fires the flat
  refusal first — naming the folder too when one exists, whether or not that
  folder contains `spec.md` — and the missing-`spec.md` refusal fires only
  when no flat file exists. The refusal text is the linked-repo migration
  prompt: a repo upgrading the kit with staged flat specs gets a named
  instruction, not silence. Enumeration (`Facts`) keeps the
  retain-with-empty-status evidence convention for malformed content; in the
  kit repo the conformance sweep additionally turns any flat `specs/*.md` red
  (story 9), so both audiences fail closed at their own surface.
- **Retire removes the folder, and every partial state has a defined re-run
  outcome.** Validation unchanged (merged-implemented, marker at HEAD via the
  folder path). Deletion order — review pickup, then `tickets/`, then
  `spec.md`, then the folder — is an implementation decision chosen to make
  the following black-box outcomes true; the outcomes, not the order, are
  what tests observe: (a) pickup removed, folder intact → re-run retires
  cleanly (an absent pickup is not an error); (b) `tickets/` partially or
  fully removed, `spec.md` present → re-run retires cleanly; (c) `spec.md`
  removed but folder residue remains → re-run refuses at resolution naming
  the residue, and the refusal text names the hand-clean step. `slugOf`
  derives the slug from the folder name for folder paths.
- **History keeps both pathspecs.** `historyDeleteLog` queries
  `:(literal,top)specs/<slug>.md` and `:(literal,top)specs/<slug>/spec.md`;
  the existing parse/merge/dedupe path is unchanged. The flat pathspec is the
  only place the flat form survives, deliberately — retired specs live there.
- **Tickets are convention-only in v1.** No parser, no subcommand, no
  frontier surface; the CLI learns the folder, not the ticket format. Ticket
  files carry: a verb-first title, `Blocked by:` naming sibling ticket titles
  (or `none`), `## What to build` as end-to-end behavior, and `## Acceptance`
  checkboxes checked as work lands; a landed ticket stays in place until the
  spec retires (whole folder leaves together). The exact template lives in the
  `craft-tickets` skill; the skill is its one source.
- **The light-path table** (drafted; wording is veto surface) lands beside
  `.bench/BENCH.md`'s "Right-size the process" rule as the standing rule that
  paragraph licenses: a change that decomposes to one independently-green
  ticket and crosses no declared seam takes the light path — charge
  `craft-tickets` bare, write the one ticket, fix and gate, no spec. The
  table is the standing OK for taking that route, so no per-change ask to
  skip the spec phase; the one ticket itself still rides the session's
  existing approval surface exactly as map #4 decided (sign-off when the
  reviewer is present, batch approval AFK, the ticket as post-hoc veto
  surface). Anything larger, a new seam, or an open decision takes the full
  workflow; the existing reviewer-requested-fix route is unchanged. Only
  FT107's table clause moves; its other clauses stay on FT107.
- **The stage table** in `craft-line`: orchestration mid/medium; ticket
  implementation cheap/low (reviewer-amended 2026-07-28 at spec time from map
  #8's original cheap/medium); review (axis or falsification delegate)
  mid/high. Defaults under the existing ladder — a failed done-claim or an
  uncleared red escalates one declared tier — and the leverage override is
  untouched: kit always-loaded prose still routes top/high for build and
  review alike. This table supersedes the standing "builds route mid" note for
  ticket-sized charges; the first builds under it are the cheap-tier re-test
  evidence `decisions/cost-follows-project-size.md` #6 waits on.
- **Skill wiring.** New skills need three artifacts in one diff: the
  `.agents/skills/bench-craft-tickets/` directory (frontmatter with `name`,
  `description`, `index:`), the `.claude/skills/bench-craft-tickets` symlink
  (the mirror is per-skill symlinks — creating the entry is part of adding a
  skill), and a skills-index regeneration (`.bench/skills-index.sh --write`).
  The conformance phase enforces all three.
- **Gate anchors** (map left these to spec time): new `require()` lines in
  `checkWorkflowAnchors`, one per load-bearing clause rather than one per
  file — `bench-implement-spec.md` charges `bench-craft-tickets`, writes
  `tickets/`, derives tickets "from the spec's stories and seams", and
  presents the breakdown under the existing approval surface;
  `.bench/BENCH.md` carries "independently-green ticket" and "crosses no
  declared seam"; `craft-line` carries the three stage-table rows;
  `bench-craft-tickets/SKILL.md` carries "smallest independently-green",
  "one write-delegate charge", and the ticket template's structural headings
  (`Blocked by:`, "What to build", acceptance checkboxes); `craft-spec` and
  `craft-tickets` carry each other's names. Each new anchor family gets a
  canary fixture in `tests/canary/workflow-guidance-anchors/` proving it
  bites, landed in the same commit as its anchor; the reviewer's completeness
  check is fixture count = new anchor-family count. Anchors pin presence of
  the load-bearing clauses only — prose *semantics* stay review-owned, the
  gap the map accepts for v1.
- **Self-hosting window, accepted:** between this spec's staging commit and
  story 9 landing, the conformance sweep (today globbing `specs/*.md`) does
  not validate this spec's coverage map, and `bench status` does not
  cross-check its ROADMAP path token. `bench coverage --check
  specs/craft-tickets/spec.md` was run by path at author time; the Go slice
  lands first in the build, closing the window.

## Testing decisions

- A good test here drives the built binary (or the contract/conformance test
  binaries) against throwaway fixture trees and asserts exit codes, TOON rows,
  and error text — never internals. Runtime/AXI contract suites run against a
  rebuilt `dist/bench`.
- Prior art: `internal/contract/runtime/runtime_spec_test.go`,
  `runtime_spec_history_test.go`, `runtime_status_test.go`,
  `internal/contract/axi/axi_coverage_test.go`, `internal/coverage` and
  `internal/spec` unit tests — all currently exercising the flat form; they
  move to the folder form (and gain the fail-posture cases) rather than being
  duplicated.
- Gate: `.bench/gate.sh` (project gate, dev tier).
- **Codex operational note (reviewer-directed 2026-07-28):** any `codex exec`
  call made from this environment must run with the yolo posture
  (approvals/sandbox bypass) or the call blocks. This spec's own falsification
  pass ran that way at `gpt-5.6-sol` / high / fast service tier on the
  reviewer's explicit instruction — a one-time top-binding approval carried in
  the phase invocation, not a new standing route; the profile's cached
  mid-tier route for the pass is unchanged. `codex exec --help` stays
  authoritative for flag spellings. Yolo delegates are unsandboxed, so
  read-only charges are verified with `git status` after every call.

### Seam diagram

Seam 1 — the CLI contract seam (`internal/spec` is the deep unit behind it;
every consumer — coverage, commit --spec, retire, history, status, roadmap,
handoff — sees the layout only through it):

    trigger: contract test (or agent) invokes the built binary in a fixture repo
        │
        ▼
    slug or path      ──▶  [ bench <coverage|spec|status|          ──▶  TOON rows / error text
    fixture tree odd       roadmap|handoff|commit --spec>               + exit code 0/1/2
    states (flat file,      → internal/spec Resolve/Facts/
    empty folder, …)        retire/History ]
                      ◀ tests attach here: run the subcommand against the
                        fixture tree; assert exit code, stdout shape, and the
                        named refusal text

Seam 2 — the conformance/canary seam (prose anchors and the sweep's bite):

    trigger: gate conformance phase / canary sweep
        │
        ▼
    kit tree (or a     ──▶  [ TestRootConformance:                ──▶  diagnostics; canary
    canary fixture          checkWorkflowAnchors +                     EXPECT match = red
    with one fact           checkCoverageMaps over                     proven
    mutated)                specs/*/spec.md ]
                      ◀ tests attach here: fixture mutates one anchored fact;
                        the phase must go red with the targeted substring

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | bare slug resolves `specs/<slug>/spec.md` from repo root and deeper cwd | seam 1 | observed 2026-07-28: `./dist/bench coverage craft-tickets` exits 1 `spec not found: craft-tickets, …/specs/craft-tickets.md` with this file present | the current binary tries only the flat fallback; the red clears only when folder resolution exists |
| 1 | `bench spec implemented` and `bench spec retire` resolve bare slugs to the folder form | seam 1 | observed 2026-07-28: `./dist/bench spec implemented craft-tickets` and `./dist/bench spec retire craft-tickets` both exit 1 trying only `specs/craft-tickets.md` | the story quantifies over every consumer; a coverage-only special case would leave these two red |
| 1 | `bench commit --spec` validates and flips a folder spec on its green commit | seam 1 | red at build: `runtime_commit_test.go` `--spec` cases moved to folder fixtures | CheckStaged/Flip is the remaining resolution consumer; against a coverage-only fix the moved cases stay red |
| 1 | explicit path args stay path-first (a cwd file shadows the fallback) | seam 1 | already covered — `internal/coverage/coverage_test.go` shadow case, re-pointed to the folder fallback | regression guard: the fallback change must not break path-first behavior |
| 2 | bare slug with flat-only `specs/<slug>.md` exits 1 naming the flat file and the folder form | seam 1 | red at build: new `internal/spec` unit case + runtime contract case, written first against the fixture | today the flat file resolves successfully, so the test asserting a refusal starts red |
| 2 | bare slug with both forms exits 1 naming both paths | seam 1 | red at build: same new cases, collision fixture | today the flat file wins silently; the collision refusal cannot pass until written |
| 3 | folder-sans-`spec.md` exits 1 naming the missing file | seam 1 | red at build: new resolution case (today a bare slug with only an empty folder reports flat-form not-found, the wrong message) | pins the distinct refusal so a half-created folder is not a generic not-found |
| 4 | `Facts` enumerates folder specs, slug = folder name, malformed retained with empty status | seam 1 | red at build: `internal/spec/spec_test.go` Facts cases moved to folder fixtures | the current glob `specs/*.md` returns nothing for folder fixtures, so moved cases start red |
| 5 | retire on a merged-implemented folder spec removes the pickup and the whole folder, `tickets/` included | seam 1 | red at build: `runtime_spec_test.go` retire cases moved to folder fixtures with a `tickets/` file present | current retire removes a single file; a folder with residue fails the moved assertions (deletion *order* is not black-box observable — the interrupt rows below pin its consequences) |
| 5 | retire refuses staged, not-at-HEAD, unknown slug, and orphaned pickup — unchanged under the layout | seam 1 | already covered — existing retire refusal contract cases, re-pointed to folder fixtures | proves the refusal set survives the layout move |
| 5 | retire re-run after interrupt states (a) pickup removed, folder intact, and (b) `tickets/` removed, `spec.md` present, completes the retire cleanly | seam 1 | red at build: new contract cases hand-building each intermediate state | pins that recoverable partial states recover; without these rows only the terminal refusal is defined |
| 5 | retire re-run after the terminal interrupt state (`spec.md` removed, folder residue remains) refuses naming the residue and the hand-clean step | seam 1 | red at build: new contract case with a hand-built residue fixture | pins the partial-state posture instead of leaving it to whatever `RemoveAll` does |
| 2 | a flat `specs/<slug>.md` beside a folder without `spec.md` fires the flat refusal naming both paths | seam 1 | red at build: new resolution case with that combined fixture | the state satisfies two refusal predicates at once; the row pins the precedence so it cannot be implementation-defined |
| 6 | `spec history` finds a spec retired at the old flat path | seam 1 | already covered — `runtime_spec_history_test.go` flat-deletion cases stay green unchanged; plus observed 2026-07-28: `./dist/bench spec history craft-tickets` renders the definitive empty state | the flat pathspec must remain; these cases are the guard against dropping it |
| 6 | `spec history` finds a folder-form deletion | seam 1 | red at build: new history contract case retiring a folder spec in a fixture repo | today only the flat pathspec is queried, so the folder deletion is invisible |
| 7 | `bench status` counts folder specs awaiting retirement; orphan pairing keys folder slugs | seam 1 | red at build: `runtime_status_test.go` retirement/orphan cases moved to folder fixtures | the current `specs/*.md` scan yields 0 on folder fixtures, so moved cases start red |
| 7 | ROADMAP cross-check classifies `specs/<slug>/spec.md` tokens; literal `<slug>` template text stays excluded | seam 1 | red at build: status roadmap-reconcile cases with folder tokens (today the regex cannot match them) | an unmatched token silently drops the row from the ambient check — the invisible-skip class |
| 8 | `bench handoff` renders live specs at the folder path; `bench roadmap` context parse recognizes it | seam 1 | red at build: `runtime_handoff_facts_test.go` + roadmap context cases moved to folder paths | ambient surfaces printing the dead flat form would send sessions to a path the CLI refuses |
| 9 | conformance validates `specs/*/spec.md` coverage maps | seam 2 | red at build: canary fixture (conformance family) carrying a folder spec with a malformed map, EXPECT its `coverage --check` diagnostic | a sweep still globbing `specs/*.md` skips the folder spec and the fixture stays green — exactly the miss this row exists to catch |
| 9 | conformance reds on a stray flat `specs/*.md` in the kit repo | seam 2 | red at build: canary fixture with one flat spec file, EXPECT the stray-flat diagnostic | without it a dead-form file survives invisibly beside the folder convention |
| 10 | skill exists, indexed, mirrored | seam 2 | already covered — `checkSkillsIndexGenerateVerify` + Claude skill-mirroring checks red on a missing symlink or unregenerated index the moment the skill directory lands | existing conformance owns the three-artifact contract; no new check needed |
| 10, 13, 14 | skill carries the unit-definition and charge clauses; `craft-line` carries all three stage-table rows; cross-pointers present | seam 2 | red at build: new `require()` anchors (one per load-bearing clause, enumerated under Implementation decisions) + one canary fixture per family mutating the anchored text | anchors pin the presence of each load-bearing clause; the canary proves each bites; wording semantics beyond the pinned clauses are review-owned, the map's accepted v1 gap |
| 10 | skill carries the ticket template's structural headings (`Blocked by:`, What to build, acceptance checkboxes) | seam 2 | red at build: template-heading anchors + a canary fixture dropping a heading | the template is the Handoff item-2 ticket-file contract; heading anchors are its gate-visible half, since nothing parses tickets in v1 |
| 11 | `/bench-implement-spec` carries the breakdown step: charges `bench-craft-tickets`, writes `tickets/`, derives "from the spec's stories and seams", presents under the existing approval surface | seam 2 | red at build: one anchor per clause + workflow-guidance-anchors canary fixture dropping the step | a name-drop edit that only mentions the skill fails the derivation and approval-surface anchors |
| 12 | `.bench/BENCH.md` carries the light-path table with both observables ("independently-green ticket", "crosses no declared seam") | seam 2 | red at build: two anchors + canary fixture removing the table | a table missing the seam clause would license the light path for cross-seam changes — the anchor pair pins both halves of the observable |
| 15 | no kit prose states the flat convention | seam 2 | not TDD-able as a class — the sweep is a hand-enumerated edit list (write-spec, final-check, craft-delegate, ROADMAP preamble, README, field guide); review verifies it | a "no stale path form" checker would need semantic judgment; the enumerated list plus review is the honest posture, and story 9's stray-file red guards the repo itself |
| 16 | every new anchor family bites | seam 2 | red at build: the canary sweep itself — each new fixture's EXPECT must match a real red; fixture *existence* is not gate-checkable, so the build lands each fixture in the same commit as its anchor and the reviewer verifies fixture count = new anchor-family count | the canary baseline rejects vacuous EXPECTs, so a rotted anchor fails here; the same-commit rule plus the count check covers the omission the sweep cannot see |

### Edge inventory

Walked per the profile's hostile-input checklist and `craft-spec`'s canonical
classes; each lands as a row above or a **Won't handle** line here.

- Flat/folder collision; flat-only; folder-sans-spec.md; the combined
  flat-beside-folder-sans-spec.md state (precedence pinned) → rows (stories
  2, 3).
- Absent `specs/` vs present-but-empty → already covered by existing
  definitive-empty-state contract cases, re-pointed (story 4/7 rows).
- Special file (FIFO/device) parked at `specs/<slug>/spec.md` → already
  covered: resolution reads through `bounds.Classify`; the folder fallback
  keeps the same `readCandidate` path (story 1 rows exercise it).
- Dangling symlink at `specs/<slug>` or `specs/<slug>/spec.md` → row via the
  classifier's stat-first posture, exercised in the moved resolution cases.
- Slug containing glob metacharacters or spaces → history pathspecs carry
  `literal` magic (both forms); resolution uses `filepath.Join`, not globs;
  moved history cases keep the existing metacharacter fixture.
- Missing trailing newline in `spec.md` → already covered: `Flip` preserves
  every byte; existing case re-pointed.
- Interrupt mid-retire → rows (story 5: two recoverable states re-run
  cleanly, the terminal residue state refuses).
- Re-run idempotency of `bench spec implemented` (double flip refuses) →
  already covered, re-pointed.
- cwd deeper than repo root → already covered by anchoring cases, re-pointed.
- Control bytes in commit subjects (history TOON) → already covered, unchanged.
- Kit vs linked repo: the kit repo gets the conformance red (story 9); a
  linked repo upgrading with staged flat specs gets the named CLI refusal
  whose text is the migration instruction (story 2). Both audiences walked.
- **Won't handle:** malformed ticket files — nothing parses tickets in v1;
  quality is caught at breakdown approval (accepted v1 gap, map #5/#7).
- **Won't handle:** ticket sizing/edge quality as a gate check — reviewed at
  breakdown approval; not machine-checkable without a parser v1 excludes.
- **Won't handle:** `check-agent-line` enforcement of the stage table — the
  delegation envelope carries no stage field (FT128's gap class, map #8);
  binding is prose plus review in v1.

## Out of scope

- **Frontier CLI surface** (`bench spec tickets` listing/validation) — a
  separate read-surface capability held to the FT6 graduate-on-evidence
  posture; ~10 edits, 3 gate runs when demonstrated use pulls it.
- **Tracker publication** (tickets → GitHub/Linear) — Bench has no tracker; a
  distinct integration capability; ~15 edits, 4 gate runs if ever wanted.
- **FT107's non-table clauses** (read-budget reroute, fix-loop shrink
  observable, squash-merge commit safety, self-contradicting-spec rule) —
  they remain FT107's batch; only the table clause moved here (map #6).
- **Retesting existing delegate-tier routings** beyond the stage table's own
  defaults — `decisions/cost-follows-project-size.md` #6 consumes this
  build's evidence; it is not built here.
