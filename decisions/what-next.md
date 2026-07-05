# what-next — working-document restructure and the drain workflow

<!-- command-currency: historical -->

Restructure the capture/priority docs and add one maintenance phase: ROADMAP.md
becomes the working document (assessment format), IDEAS.md becomes the capture
sink, `/bench-what-next` drains both learnings and ideas into the roadmap, and
`bench roadmap` (CLI) pulls the document with drain status.

## #1: Where does the drain judgment live — CLI or workflow?

Type: Grill

### Question
Triaging a parked line or learnings entry into a prioritized roadmap item needs
judgment a deterministic CLI can't supply. Where does it live?

### Answer
Hybrid. `bench roadmap` (CLI) is deterministic: prints ROADMAP.md plus a drain
status (N parked ideas, N open learnings entries → run `/bench-what-next`); when
both counts are zero it prints the roadmap's `## Recommended sequence` section
verbatim as the call to action. `/bench-what-next` (workflow command) does the
judgment: reconcile, drain, re-sequence, recommend.

## #2: Does `/bench-what-next` replace `/bench-integrate-learnings`?

Type: Grill

### Question
Learnings already drains through `/bench-integrate-learnings` (verdict with
reviewer sign-off, kit edits under craft-synthesis). Two drains of one source
would conflict.

### Answer
Yes — `/bench-integrate-learnings` is retired; the journal has exactly one exit.
`/bench-what-next` inherits the verdict-with-sign-off: every open entry gets a
verdict in the batch diff — work-shaped → roadmap item; rule-shaped → roadmap
item whose next action is the kit edit (built later under craft-synthesis,
gated as usual); or dismissed with reviewer OK. The retirement ripples through
every file naming the old command (BENCH.md, learnings.md header,
craft-synthesis skill, command files, canary fixtures).

## #3: What are the two files called?

Type: Grill

### Question
The assessment format becomes the working document; the old roadmap becomes the
capture sink. Names?

### Answer
ROADMAP.md keeps its name and takes the assessment format — the working
document. IDEAS.md is the capture sink (`bench idea` pairs with it). Current
ASSESSMENT.md content migrates into the new ROADMAP.md and ASSESSMENT.md is
deleted; current ROADMAP.md content becomes IDEAS.md.

## #4: Full drain or partial?

Type: Grill

### Question
Do all ideas leave IDEAS.md on a drain, or can some stay parked in the inbox?

### Answer
Full drain. IDEAS.md is a pure inbox and empties to zero each run; ROADMAP.md
holds all triaged state including a parked-pending-evidence tier (the FT6
pattern). The empty inbox is the normal post-drain state, so the top-2–3
recommendation fires whenever capture is clear.

## #5: How does the CLI know the top items?

Type: Grill

### Question
`bench roadmap`'s empty-state output needs a deterministic source for "the top
two or three items."

### Answer
Format contract: ROADMAP.md ends with `## Recommended sequence` — 2–3 numbered
lines, each naming the item and the phase command to run. `/bench-what-next`
refreshes the section every run; the CLI extracts and prints it verbatim, no
judgment.

## #6: Does `/bench-what-next` reconcile against the tree?

Type: Grill

### Question
Assessment discipline included verifying items against the tree (shipped work
removed). Drain-only, or reconcile too?

### Answer
Reconcile every run, before draining: shipped → removed, stale → reworded. The
empty-state recommendation is only trustworthy if the roadmap is current.
`/bench-what-next` is the single roadmap-maintenance phase.

## #7: Per-item sign-off or batch?

Type: Grill

### Question
Does the reviewer verdict each item interactively during a run, or approve the
pass as a whole?

### Answer
Batch-propose. The agent drafts the full pass — reconciled roadmap, drained
inbox, journal verdicts including dismissals — as one uncommitted diff; the
reviewer approves or adjusts once; commit on green. The diff is the verdict
sheet; nothing leaves the journal without that one approval.

## #8: How is completed work marked in the roadmap?

Type: Grill

### Question
After an item ships, is it checked off, and how does spec retirement relate?

### Answer
No completion markers — presence is status (invariant 3; history lives in git).
spec-retire removes the roadmap item in the same commit that deletes the spec;
`/bench-what-next`'s reconcile is the backstop for anything that slipped.
Archaeology stays the git recipes (`git log --grep=spec-retire`,
`git log --diff-filter=D -- specs/`); `bench specs --retired` stays parked
pending evidence.

## #9: What happens to the existing CLI surface?

Type: Grill

### Question
`bench roadmap` prints the old file verbatim; `bench idea` appends to it. What
changes?

### Answer
No new subcommand. `bench roadmap` is upgraded per #1 (a new `bench what-next`
subcommand was considered and reverted). `bench idea` retargets IDEAS.md,
behavior otherwise unchanged (created on first append). `bench status`'s footer
counts IDEAS.md lines plus open learnings entries and points at
`/bench-what-next`. `/bench-shape-idea`'s cold-start section retargets: it
offers top ROADMAP.md items instead of parked lines, and the
remove-line-on-promotion rule moves to retirement (#8) — a roadmap row persists
until the work ships.

## Handoff

1. **Module boundaries.** `internal/roadmap` (Go): `bench idea` append →
   IDEAS.md; `bench roadmap` print + drain counts + Recommended-sequence
   extraction. `internal/status`: footer counts and pointer.
   `.agents/commands/bench-what-next.md` (new phase; replaces
   `bench-integrate-learnings.md` and its per-harness copies). Prose surfaces:
   BENCH.md (Capture + Workflow + CLI inventory), `/bench-shape-idea` roadmap
   section, spec-retire wiring, learnings.md header. Dogfood migration:
   ASSESSMENT.md → new-format ROADMAP.md, old ROADMAP.md → IDEAS.md.
2. **Contracts.** `bench idea "<text>"`: appends `- YYYY-MM-DD  <text>` to
   IDEAS.md, creating it if absent; exit 0, TOON error + exit 1 on write
   failure. `bench roadmap`: ROADMAP.md present + drainable items → doc plus
   drain-status block naming `/bench-what-next`; both sources empty → doc plus
   `## Recommended sequence` verbatim; ROADMAP.md absent → pointer to run
   `/bench-what-next`; exit 0 in all non-error states. Format contract:
   `## Recommended sequence`, 2–3 numbered lines, each `<item> — <phase
   command>`.
3. **Deep vs thin.** `internal/roadmap` is the deep module (file discovery,
   counting, section extraction hidden behind the two subcommands).
   `internal/status` footer is a thin pass-through over the same counts — reuse
   the roadmap package's counter, don't re-derive (one source per fact).
   `/bench-what-next` is prose, no code seam.
4. **Black-box assertables.** CLI stdout + exit code per state in #2's contract
   (undrained, empty-state, missing-file); IDEAS.md content after `bench idea`;
   status footer line; canary anchors on the new command prose and on the
   removal of integrate-learnings referents.
5. **Gate attachment.** Go tests and conformance run under `bench gate` as
   usual; workflow-guidance-anchor canaries pin `/bench-what-next` prose. The
   drain's judgment itself is not gate-visible — it lands as a reviewed batch
   diff (advisory review + reviewer approval), which is the designed check.
6. **Hostile-input owners.** Missing/empty ROADMAP.md or IDEAS.md → `bench
   roadmap` state machine (#2). ROADMAP.md lacking the `## Recommended
   sequence` heading → CLI reports the missing contract section explicitly,
   never silent-empty output. Multiline/odd idea text → existing `bench idea`
   append semantics carry over. Concurrent writer during a drain → worktree
   discipline (existing learnings rule), not new code.
7. **Uncertainty flags.** None blocking. Exact wording of the drain-status
   block and whether `bench status` shows one combined row or two is mid-tier
   spec-writer latitude.
8. **Rejected alternatives.** A separate `bench what-next` subcommand (reverted
   in-grill; `bench roadmap` supersedes). Keeping `/bench-integrate-learnings`
   alongside the drain (two exits from one journal). Interactive per-item
   verdicts. Completion checkmarks in the roadmap. Letting ideas stay parked in
   IDEAS.md (kills the empty-state trigger).
9. **Domain watch-outs.** A roadmap row leaves only at spec-retire (same
   commit) or a reconcile pass; a journal entry leaves only via a verdict in an
   approved `/bench-what-next` diff. The renames ripple into canary fixtures
   and conformance checks — updating fixtures to the new prose is part of the
   diff, never a gate weakening.

Dependency order: (1) doc restructure + CLI retarget — IDEAS.md, new-format
ROADMAP.md with migration, `bench idea`/`bench roadmap`/`bench status`; (2)
`/bench-what-next` phase + `/bench-integrate-learnings` retirement; (3)
spec-retire wiring for roadmap-row removal. Recommendation only — slicing is
the reviewer's call.
