# ft228-debug-restoration

Status: implemented

Decision source: `roadmap/FT228.md` — the board row the reviewer opened from the 2026-08 capability audit (ledger L-10, L-11, L-30, action item A5), carrying the 2026-07-23 trustworthy-gate-under-load occurrence. The row's one open fork — the Codex trigger — was closed by the reviewer in this session on 2026-08-19: `$bench-debug` becomes implicitly invocable on Codex.

Verification log: 2 iterations to accept — the round (`opus` / high, read-only) returned two blocking findings and eleven partials on iteration 1: the Codex flip shipped with an anti-trigger adapter description (folded as story 24, row IP10, and the description rewrite in ticket 03), and `internal/conformance/registry_test.go` sat outside every fence though the new fixture family cannot land without it (folded into the fences and ticket 03). The partials folded: the anchor-file-missing diagnostic replaces the reads-as-empty claim, the budget row tightened to 170 with the build deriving the final number, the gate seam renamed to the ordinary test phase, the DR rows renumbered and narrowed to their needles, IP4 rebuilt on ghost-command and lookup-ahead semantics, IP7 restated on the hardened producer path, the adapter fixtures' `BASE` mechanics stated with the ghost carve-out, the CHANGELOG split across tickets 01 and 04, group C's line naming the doc-paragraph exception, and the promise-inflation flag on story 12 recorded in Further notes. Iteration 2 accepted with five acceptance partials, all folded. The round's quiz split ticket 03 into a Codex half and a Claude half. Mutation probes (a)–(d) in Testing decisions were named, not executed — the build's tickets execute them.

## Problem

The debug phase lost half its upstream discipline, and its trigger is unreliable
on both harnesses.

`.agents/commands/bench-debug.md` descends from the pinned upstream
`diagnosing-bugs` skill (`d574778f`, 134 lines). The audit's clause-by-clause
diff (ledger L-11) found the four load-bearing constraints preserved verbatim —
no hypothesis before a red loop, the exact symptom, one command already run,
the original-loop rerun — plus five real Bench additions. It also found three
compressions nothing ever decided:

1. **The loop-construction menu dropped from ten entries to five.** Upstream
   Phase 1 lists ten ways to build a feedback loop, tried in order — failing
   test, curl, CLI invocation, headless browser, trace replay, throwaway
   harness, property/fuzz loop, bisection harness, differential loop, and a
   structured human-in-the-loop last resort. The current file names five in one
   sentence.
2. **Both hard stop-gates and the checkbox completion form are gone.** Upstream
   ends Phase 1 with "No red-capable command, no Phase 2" over a four-item
   `- [ ]` checklist, and ends Phase 2 with "Do not proceed until you have
   reproduced and minimised" over a three-item confirmation checklist. The
   current file softens both to "Don't proceed on a theory" and plain bullets.
3. **"Tighten the loop" vanished as a named step** — treat the loop as a
   product: faster, sharper, more deterministic — folded into three adjectives.

The compressions happened silently because nothing pins this text: the anchors
registry guards the Bench additions (accused command, expected-failure form,
write isolation, spec archaeology) but none of the upstream constraints. The
same mechanism that lost them once will lose them again.

Separately, the phase has no reliable trigger. On Claude the command file
carries no `disable-model-invocation` key, so the model can route a bug to the
phase — but no check grades that fact, on this or any phase. On Codex every
phase adapter's `agents/openai.yaml` pins `allow_implicit_invocation: false`,
so `$bench-debug` fires only when the reviewer types it: the
strongest-evidence behavior has the weakest trigger (ledger L-10). And all
thirteen Codex adapter `SKILL.md` files carry `disable-model-invocation: true`
— a Claude key that no harness reads there (Claude's surface is
`.claude/commands`, Codex's policy file is `openai.yaml`), pure dead
configuration (ledger L-30).

Finally, the 2026-07-23 trustworthy-gate-under-load diagnosis exposed a
reproduction-economics gap in the phase's own text: synthetic load shapes
stayed green while real `bench gate` host load reproduced the failure. Nothing
in the phase says a green stand-in only narrows a hypothesis.

## Solution

Three moves, each durable.

- **Restore the constraints at the current phase-file contract.** The ten-entry
  menu returns as one shipped reference file,
  `.agents/skills/bench-debug/references/loop-constructions.md`, pointed at
  from Phase 1 — both harnesses read the same canonical file, so one pointer
  serves both. The two hard stop-gates return verbatim. The Phase 1 completion
  criterion, the Phase 2 confirmation, and the Phase 6 close-out return as
  `- [ ]` checkbox forms. "Tighten the loop" returns as a named step. The
  reproduction-economics rule lands in Phase 2: a green proxy only narrows a
  hypothesis — a load- or environment-sensitive failure is reproduced through
  the accused command under exposing conditions before another stand-in is
  trusted. Every preserved constraint and every Bench addition stays.
- **Pin what was restored.** New anchor rows in the anchors registry require
  the restored needles in the phase file and the menu content in the reference
  file, and two new `workflow-guidance-anchors` canary fixtures prove the pins
  bite through the registered owner. The phase file joins the reviewer-owned
  prose-budget table with an exact row, so regrowth is graded too.
- **Settle the trigger as declared per-phase policy.** A per-phase
  invocation-policy table becomes the reviewed expectation the conformance
  check grades both harness surfaces against: the command file's
  `disable-model-invocation` frontmatter on the Claude side (the check that was
  missing), and `openai.yaml` on the Codex side (which today hardcodes
  `false` for everyone). `bench-debug`'s row says implicitly invocable on both
  — the reviewer's 2026-08-19 decision — and its `openai.yaml` flips to
  `true`. Its adapter description also drops the explicit-only phrasing for a
  symptom-bearing trigger: once implicit invocation is on, the description is
  the only text Codex matches against, and the current one instructs the model
  not to fire. Every other phase's row records today's posture unchanged. The inert
  Claude key leaves all thirteen adapter `SKILL.md` files, and its
  reintroduction turns the gate red. `.bench/BENCH-reference.md`'s invocation
  section states the per-phase policy and names the exception.

FT24 stays parked as the upstream re-check trigger for a deny-capable Codex
hook surface; nothing here graduates it.

## User stories

### Restoring the debug discipline

Line: `fable` / high. Guidance prose compounds through every session that loads
it while the edit costs few tokens — the leverage override in `craft-line`.

1. As a debugging agent, I want Phase 1 to point at a shipped loop-construction
   menu of ten constructions plus the human-in-the-loop last resort, so that
   loop building starts from the full option space instead of the five
   survivors.
2. As a debugging agent, I want the hard stop-gate "No red-capable command, no
   Phase 2" back in Phase 1, so that hypothesis work cannot start before a
   red-capable loop exists.
3. As a debugging agent, I want the hard stop-gate "Do not proceed until you
   have reproduced and minimised" back in Phase 2, so that diagnosis cannot
   start on an unconfirmed or unminimised repro.
4. As a debugging agent, I want Phase 1's completion criterion as a `- [ ]`
   checkbox form (red-capable, deterministic, fast, agent-runnable), so that
   Phase 1 completion is a checked state rather than an adjective.
5. As a debugging agent, I want Phase 2's confirmation as a `- [ ]` checkbox
   form (the user's failure mode, a debuggable reproduction rate, the captured
   symptom), so that the wrong-bug trap is a checked step.
6. As a debugging agent, I want "Tighten the loop" restored as a named step —
   faster, sharper, more deterministic — so that a loop is treated as a product
   rather than accepted as found.
7. As a debugging agent, I want the reproduction-economics rule in Phase 2 — a
   green proxy only narrows a hypothesis, and a load- or environment-sensitive
   failure is reproduced through the accused command under exposing conditions
   before another stand-in is trusted — so that a synthetic stand-in's green
   cannot close a diagnosis the way it nearly did on 2026-07-23.
8. As a reviewer, I want every Bench-specific addition — write isolation before
   the first repro artifact, the accused-command rule, the expected-failure
   quarantine form, the seam-versus-edit-owner rule, spec archaeology, the
   delegate blocked-report shape — to survive the restoration, so that
   recovering upstream text does not trade away what Bench added.
9. As a teammate on any harness, I want the menu in one reference file reached
   from the canonical phase file, so that Claude, Codex, and every AGENTS.md
   harness read one source.

### Pinning the restoration

Line: `opus` / medium. Registry data rows and canary fixtures at an established
conformance seam; the oracle's correctness matters more than speed.

10. As a maintainer, I want anchor rows that turn the gate red when a restored
    constraint is deleted or softened, so that the silent compression FT228
    repairs cannot recur.
11. As a maintainer, I want canary fixtures demonstrating the new anchors bite
    through their registered owner, so that each pin is proved red-capable
    rather than assumed.
12. As a maintainer, I want `.agents/commands/bench-debug.md` under an exact
    row in the reviewer-owned prose-budget table, so that the restored file has
    a graded ceiling instead of an open invitation to regrow.

### Settling the trigger on both harnesses

Line: `opus` / medium. Conformance logic at the adapters seam plus one-line
config edits. The one `.bench/BENCH-reference.md` paragraph rides here as a
deliberate exception to the doc-authoring override: a paragraph inside a
conformance slice keeps the slice vertical, and review reads the prose.

13. As a Codex operator, I want `$bench-debug` implicitly invocable, so that a
    reported bug routes to the debug path without me remembering the phase
    name (reviewer decision, 2026-08-19).
14. As a maintainer, I want one per-phase invocation-policy table graded
    against both harness surfaces, so that each phase's trigger is a declared
    reviewed decision instead of ambient frontmatter.
15. As a maintainer, I want a Claude phase whose command-file frontmatter
    disagrees with the table to red the gate, so that the Claude side gains the
    invocation-policy check it never had.
16. As a maintainer, I want a Codex adapter whose `openai.yaml` disagrees with
    the table to red the gate, so that the `bench-debug` flip cannot silently
    revert.
17. As a maintainer, I want a command file with no policy row to red the gate,
    so that a new phase arrives with a declared trigger or not at all.
18. As a maintainer, I want a policy row naming no command file on disk to red
    the gate, so that the table cannot go stale in the other direction.
19. As a maintainer, I want the inert `disable-model-invocation` key removed
    from every Codex adapter `SKILL.md` and its reintroduction graded, so that
    dead configuration cannot mislead the next reader into thinking it is the
    control.
20. As a teammate who just walked in, I want `.bench/BENCH-reference.md`'s
    invocation section to state the per-phase policy and name the `bench-debug`
    exception, so that the docs match the tree.
21. As a maintainer, I want the modified check to keep refusing hostile
    `SKILL.md` files — symlink, FIFO, oversized, invalid UTF-8 — with a
    diagnostic naming the refused path, so that the hardening composition
    holds through the rewrite.
22. As a maintainer, I want an `openai.yaml` that spells neither policy value to
    red the gate as undeclared, so that a file declaring nothing cannot pass by
    lacking the wrong value.
23. As a maintainer, I want the Claude-side parse anchored to the frontmatter
    block, so that a command body quoting the key in prose does not flip that
    phase's graded policy.
24. As a Codex operator, I want the `bench-debug` adapter's description to
    carry the symptom trigger — broken, throwing, failing, or slow — instead of
    "Use only when the reviewer invokes", so that the implicit invocation the
    flip enables has matchable text and an anti-trigger cannot silently defeat
    the settle.

## Implementation decisions

- **The menu is a reference file under the phase's Codex adapter skill.**
  `.agents/skills/bench-debug/references/loop-constructions.md` follows the
  established references location (`craft-seams`, `craft-review`,
  `craft-delegate` precedent), ships automatically via the package's `.agents/`
  glob, and sits outside the prose-budget universe by that table's own rule.
  Its content adapts upstream's ten constructions to this kit's voice; entry
  ten describes the structured human-in-the-loop pattern — drive the human
  with a script so the loop stays structured — without shipping an executable
  template.
- **Restoration edits the current file in place.** The four preserved upstream
  constraints and every Bench addition stay word-for-word; the existing anchor
  rows over the file are the regression net for that claim. The restored
  pieces land at their upstream positions: the pointer and stop-gate in
  Phase 1, the confirmation checklist, stop-gate, and reproduction-economics
  rule in Phase 2, "Tighten the loop" as a named step inside Phase 1, and the
  checkbox close-out in Phase 6. Upstream's dead
  `/improve-codebase-architecture` pointer is not imported.
- **Anchors pin exact sentences, not topics.** New `Require` rows in
  `internal/anchors/registry_data.go`: the two stop-gate sentences, the
  `Tighten the loop` step name, the reference-file pointer, the
  reproduction-economics sentence, and the `- [ ]` spellings of the first
  Phase 1 criterion and the first Phase 2 confirmation in the phase file; two
  menu needles (`Bisection harness`, `Property / fuzz loop`) in the reference
  file. A deleted reference reds with the registry's own
  anchor-file-missing diagnostic, and an emptied one reds both needles.
  Existing debug rows are untouched.
- **Two new `workflow-guidance-anchors` fixtures** prove the pins fire: one
  softens the Phase 1 stop-gate sentence, one deletes the
  reproduction-economics sentence, each with `BASE` copying the phase file and
  an `EXPECT` naming the anchor's diagnostic. Mutation anchors must occur
  exactly once in the file; the fixture materializer enforces that. The
  `workflow-guidance-anchors` family already has its registration in
  `canaryFixtureFamilyRegistry`, so these two fixtures need no
  `registry_test.go` edit.
- **The invocation-policy table is the independently authored expectation.**
  A `phaseInvocationPolicy` table beside the check in
  `internal/conformance/skills_index_checks_test.go` maps each phase to two
  booleans: Claude-model-invocable and Codex-implicit. The check (the existing
  `skills-index-command-adapters` registration and its parent implementation
  keep their names) grades three facts per command file: the command
  frontmatter's `disable-model-invocation` key equals the negation of the
  Claude boolean, `openai.yaml` spells exactly the value the Codex boolean
  demands (a file spelling neither, or an empty file, is red as undeclared),
  the adapter `SKILL.md` frontmatter carries no `disable-model-invocation`
  key, and a Codex-implicit phase's adapter description carries no
  explicit-only phrasing (the marker `Use only when the reviewer invokes`). A
  command file with no table row is red, and a table row with no command file
  is red — the policy lookup runs ahead of the adapter-existence checks, so an
  undeclared phase reds as undeclared even when its adapter is also missing.
  The Claude parse reads only the frontmatter block, so a body mention of the
  key is inert. This table duplicates the frontmatter facts deliberately: its
  independence is what lets a silent flip of either surface turn the gate red,
  the named exception to the one-source rule, and the canary fixtures record
  and demonstrate those reds.
- **Policy values.** `bench-debug`: Claude-invocable, Codex-implicit (the
  settle). `bench`, `bench-final-check`, `bench-implement-spec`,
  `bench-review-implementation`, `bench-shape-idea`, `bench-write-spec`:
  Claude-invocable, Codex-explicit (today's posture). `bench-assess`,
  `bench-deepen`, `bench-drain`, `bench-setup-repo`, `bench-update-kit`,
  `bench-what-next`: disabled on Claude, explicit on Codex (today's posture).
  The `bench-what-next` alias keeps its existing special-case body handling;
  its Claude fact lives in its own command file's frontmatter like every other
  row's.
- **Four new `skills-index-command-adapters` fixtures**: `openai.yaml` flipped
  against the table, a command file's frontmatter gaining the disable key
  against the table, a command file present with no table row (the fixture
  supplies only the ghost command file, and the row reds on the missing policy
  before the missing adapter), and the inert key reintroduced into an adapter
  `SKILL.md`. Each fixture's `BASE` names the command file, the adapter's own
  `SKILL.md`, its `agents/openai.yaml`, and the two guide files (`.bench/BENCH.md`,
  `.bench/BENCH-reference.md`) — the check's universe is the command-file glob
  in the materialized root, and without the guide files the guide-mention
  diagnostic would dominate the intended red — except the unlisted-phase
  fixture, which overlays only the ghost command file. The family gains one
  `canaryFixtureFamilyRegistry` registration in
  `internal/conformance/registry_test.go` naming
  `skills_index_checks_test.go` and `checks_test.go`; the five existing exact
  rows stay and override. The stale-table direction and
  the parse edges (body-prose mention, undeclared value, explicit-only
  description on an implicit row) are in-package cases over the check function
  with constructed roots, the pattern the file already uses.
- **`openai.yaml` flips for `bench-debug` only** —
  `allow_implicit_invocation: true`. The other twelve stay `false`. In the same
  slice, `.agents/skills/bench-debug/SKILL.md`'s description is rewritten to a
  symptom-bearing trigger mirroring the command file's own — something is
  broken, throwing, failing, or slow — dropping `Explicit` and `Use only when
  the reviewer invokes`; the other twelve adapters keep theirs.
- **The inert key leaves all thirteen adapter `SKILL.md` files.** It is dead on
  Codex (`openai.yaml` is the policy surface) and invisible on Claude
  (`checkClaudeSkillMirror` keeps phase adapters out of `.claude/skills`).
- **Docs.** `.bench/BENCH-reference.md`'s Harness Invocation section drops the
  blanket "explicit-only" claim in favor of the per-phase policy statement and
  names `$bench-debug` as the implicitly invocable exception with its
  one-clause why. The prose-budget table in `projects/benchkit.md` gains an
  exact `.agents/commands/bench-debug.md` row at 170 — the build derives the
  final number from the file ticket 01 landed, tight to the table's 0–5-line
  convention (restored size is ~165 lines); the section's "every other
  `.agents/commands/*.md` file" sentence remains true as written. Two
  CHANGELOG entries under Unreleased, each landing with the ticket it
  documents: the restoration entry with ticket 01, the invocation-policy entry
  with the final trigger ticket.

## Testing decisions

- The external behavior a good test exercises: the shipped guidance carries the
  restored constraints and cannot silently lose them again, and each harness's
  invocation trigger for every phase matches one reviewed table — with the
  gate, not review, catching drift on both.
- Seams. Prose facts attach to the anchors registry
  (`workflow-guidance-anchors` → `docs-currency-workflow`), whose bite
  discipline (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) proves
  each fixture reds through its registered owner and goes green on restore.
  Policy facts attach to the `skills-index-command-adapters` check, graded by
  its canary family and by in-package cases over constructed roots. The budget
  fact attaches to the `guidance-prose-budgets` check parsing the profile's
  table. All three are existing gate-run conformance seams; no new check
  registration is needed.
- Gate seam: the gate's ordinary `test` phase, which carries the graded root
  and the dev tier to the conformance entry point, observes everything here;
  no shell or gate-script edit.
- The hostile-input composition holds by keeping the modified check in
  `hostileSkillReaders`; `TestRegisteredSkillReadersRefuseHostileSkillFiles`
  re-proves symlink, FIFO, oversized, and invalid-UTF-8 refusal against the
  rewritten reader.
- Mutation probes recorded in the verification log, not retained: (a) restore
  the stop-gates as paraphrases → the exact-needle anchors red; (b) inline the
  menu in the phase file without the reference → the registry's
  anchor-file-missing diagnostic reds on the absent reference; (c) derive the
  policy expectation from disk instead
  of the table → the flipped-yaml and flipped-frontmatter fixtures red; (d)
  parse the disable key with a whole-file Contains → the body-prose in-package
  case reds.

### Seam diagram

    kit gate (ordinary test phase → conformance entry point)
        │
        ├─▶ [ workflow-guidance-anchors: anchors registry ]
        │       .agents/commands/bench-debug.md ──▶ restored needles present? ──▶ diag on loss
        │       .agents/skills/bench-debug/references/loop-constructions.md ──▶ menu needles present?
        │           ◀ tests attach here: 2 canary fixtures mutate the phase file and
        │             prove the exact diagnostic fires through the registered owner
        │
        ├─▶ [ skills-index-command-adapters: phaseInvocationPolicy table ]
        │       .agents/commands/<phase>.md frontmatter ──▶ matches Claude column?
        │       .agents/skills/<phase>/agents/openai.yaml ──▶ matches Codex column?
        │       .agents/skills/<phase>/SKILL.md ──▶ carries no inert Claude key?
        │       .agents/skills/<phase>/SKILL.md description ──▶ no explicit-only phrasing on an implicit row?
        │           ◀ tests attach here: 4 canary fixtures + in-package constructed roots
        │
        └─▶ [ guidance-prose-budgets: profile table ]
                .agents/commands/bench-debug.md ──▶ within its exact budget row?

### Acceptance coverage map
| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| DR1 | 1, 9 | the phase file's Phase 1 carries the pointer naming `.agents/skills/bench-debug/references/loop-constructions.md` | anchor row over the phase file | a restoration that re-inlines a short menu never names the reference, and the pointer needle is absent |
| DR2 | 1, 9 | the reference file carries the menu needles `Bisection harness` and `Property / fuzz loop` | anchor rows over the reference file | a dead pointer, a deleted file, or a five-entry paraphrase satisfies DR1 while the menu stays lost — a deleted file reds the registry's anchor-file-missing diagnostic and an emptied one reds both needles |
| DR3 | 2 | the phase file carries the exact sentence `No red-capable command, no Phase 2` | anchor row over the phase file | softening to "don't proceed on a theory" is the precise compression the audit found, and the exact needle reds on it |
| DR4 | 3 | the phase file carries the exact sentence `Do not proceed until you have reproduced and minimised` | anchor row over the phase file | the second stop-gate was dropped entirely, and a topic-level check would pass on the surviving minimise prose |
| DR5 | 4 | the first Phase 1 completion criterion is a `- [ ]` checkbox line | anchor row on that criterion's checkbox spelling | plain bullets satisfy a word-presence check while dropping the checkable form — the needle includes the `- [ ]` marker |
| DR6 | 5 | the first Phase 2 confirmation is a `- [ ]` checkbox line | anchor row on that confirmation's checkbox spelling | same class as DR5 on the second surface, pinned separately so losing either form reds alone |
| DR7 | 6 | the phase file carries a step named `Tighten the loop` | anchor row on the step name | the step was folded into three unnamed adjectives once, and the named-step needle reds on that fold — the adjectives themselves are ticket 01 prose review reads |
| DR8 | 7 | the phase file carries the sentence beginning `A green proxy only narrows a hypothesis` | anchor row over the phase file | the 2026-07-23 diagnosis had no textual home, and this needle is that rule's one source in the phase |
| DR9 | 8 | every pre-existing bench-debug `Require` anchor still matches the restored file and the `Forbid` needle stays absent | existing registry rows, run by the kit gate over the real tree | the cheapest restoration is a wholesale paste of upstream `d574778f`, which drops the Bench additions and reds these rows |
| DP1 | 10, 11 | a fixture softening the Phase 1 stop-gate produces that anchor's exact diagnostic through the registered owner, and its restore runs green | fixture bite proof over `workflow-guidance-anchors` | an anchor can be registered with an unreachable diagnostic, and the bite-and-restore pair is what proves the pin fires on the mutation rather than on ambient state |
| DP2 | 10, 11 | a fixture deleting the reproduction-economics sentence produces that anchor's exact diagnostic through the registered owner, and its restore runs green | fixture bite proof over `workflow-guidance-anchors` | same class as DP1 for the rule this spec adds rather than restores |
| DP3 | 12 | the budget table carries an exact `.agents/commands/bench-debug.md` row and the restored file passes under it | `guidance-prose-budgets` check over the profile table | without the row the file is outside the reviewed universe, and the gate says nothing when it doubles |
| IP1 | 13 | `.agents/skills/bench-debug/agents/openai.yaml` spells `allow_implicit_invocation: true` | the policy check over the real tree | with `false` still on disk, Codex cannot self-invoke the phase and the check now reds because the table's row says implicit |
| IP2 | 14, 16 | a fixture flipping `bench-debug`'s `openai.yaml` back to `false` reds with a mismatch diagnostic naming the adapter | `skills-index-command-adapters` fixture | a check that derives its expectation from disk is green on any drift — the independent table is what makes a silent revert observable |
| IP3 | 15 | a fixture adding `disable-model-invocation: true` to a command file the table marks Claude-invocable reds naming the phase | `skills-index-command-adapters` fixture | the Claude side previously had no invocation-policy check at all, and this is the parity check biting on a Claude surface |
| IP4 | 17 | a fixture supplying only a ghost command file the table does not name reds on the missing policy row | `skills-index-command-adapters` fixture | the policy lookup runs ahead of the adapter-existence checks, so the red names the undeclared trigger — a lookup that skips misses lets the next phase arrive triggerless and ungraded |
| IP5 | 18 | a policy row whose command file is absent from the materialized root reds naming the stale row | in-package case over a constructed root omitting one command file from `BASE` | the table is compiled in while fixtures control only the tree, so omitting the file is the one way to make a row stale — a check that only walks disk never reads the table's own keys |
| IP6 | 19 | a fixture reintroducing `disable-model-invocation` into a phase adapter's `SKILL.md` frontmatter reds naming the file | `skills-index-command-adapters` fixture | the key is dead on Codex and invisible on Claude, so ungraded it re-accretes exactly the way it did |
| IP7 | 21 | over each hostile `SKILL.md` planting, the check completes and emits a diagnostic naming the refused path, with its skill reads staying on the hardened producer path | existing `TestRegisteredSkillReadersRefuseHostileSkillFiles` composition test | a rewrite that reads adapter files off the hardened path opens hostile bytes before classifying them, hanging on a FIFO or reporting a clean skill it never read |
| IP8 | 23 | a command file quoting `disable-model-invocation` only in body prose keeps its table policy green | in-package case over a constructed root | a whole-file Contains parse flips a phase's graded policy on documentation, the profile checklist's quoted-grammar-token class |
| IP9 | 22 | an `openai.yaml` that is empty or spells neither policy value reds as undeclared | in-package case over a constructed root | a parse that only rejects the wrong value reports green for a file that declares nothing |
| IP10 | 24 | `.agents/skills/bench-debug/SKILL.md`'s description carries the symptom trigger and no `Use only when the reviewer invokes` clause | the policy check over the real tree, plus an in-package case planting the explicit-only phrase on a Codex-implicit row | the yaml flip alone leaves the anti-trigger description as the only text Codex matches against, so implicit invocation never fires and the settle silently fails |

Not covered: story 20 — one paragraph of reference prose; review reads it.

### Edge inventory

- **Absent vs present-but-empty:** an absent `openai.yaml` keeps its existing
  missing-metadata diagnostic; a present-but-empty one is red as undeclared
  (IP9). An absent
  reference file reds the registry's anchor-file-missing diagnostic, and a
  present-but-empty one reds both DR2 needles.
- **A grammar token quoted in surrounding prose:** the profile checklist's
  class, walked at the Claude parse — a body mention of the key is inert
  (IP8). `bench-craft-skills/SKILL.md` mentions the key in prose and is not a
  phase adapter, so it is outside the check's universe.
- **Special files and symlinks in a discovered path:** the check stays in
  `hostileSkillReaders` and refuses before reading (IP7).
- **The alias:** `bench-what-next` keeps its thin-alias body handling in the
  check; its policy row says disabled on Claude and explicit on Codex, matching
  the `disable-model-invocation: true` its own command file carries today.
- **Mutation-anchor uniqueness:** each fixture's `old` string must occur
  exactly once in the phase file; the materializer refuses otherwise, so the
  build picks needles that are unique sentences.
- **Won't handle:** moving the Bench-mechanics tail of the phase file behind a
  pointer (A5's second half) — the tail stays in place, where seven existing
  anchor rows and two fixtures already grade it; the in-scope caller is the
  phase file carrying it unchanged (retargeting is Out of scope).
- **Won't handle:** shipping upstream's executable HITL script template — the
  reference file's tenth entry describes the structured pattern, and a project
  that needs the script writes its own.
- **Won't handle:** a harness that reads `.agents/skills/` natively with
  Claude-key semantics — after the inert-key removal such a harness could
  implicitly match a phase adapter, but no adopted harness does:
  `openai.yaml` is Codex's operative control, `checkClaudeSkillMirror` keeps
  phase adapters off Claude's surface, and OpenCode is unadopted with an
  unbound column.
- **Won't handle:** re-deciding the other twelve phases' Codex posture — the
  table records today's explicit-only values for them, and each future flip is
  a one-line reviewer edit plus its yaml (Out of scope).
- **Won't handle:** `openai.yaml` files without a trailing newline or with
  reordered keys — the value read is substring-based over a two-line file the
  kit itself writes, and the undeclared-value red (IP9) covers the degenerate
  end of that laxity.

## Ownership fences

Tickets are serial on one retained integration source. Reviewer disposition:
approve, merge, or split at sign-off.

- `.agents/commands/bench-debug.md`
- `.agents/skills/bench-debug/`
- `.agents/skills/bench/SKILL.md`
- `.agents/skills/bench-assess/SKILL.md`
- `.agents/skills/bench-deepen/SKILL.md`
- `.agents/skills/bench-drain/SKILL.md`
- `.agents/skills/bench-final-check/SKILL.md`
- `.agents/skills/bench-implement-spec/SKILL.md`
- `.agents/skills/bench-review-implementation/SKILL.md`
- `.agents/skills/bench-setup-repo/SKILL.md`
- `.agents/skills/bench-shape-idea/SKILL.md`
- `.agents/skills/bench-update-kit/SKILL.md`
- `.agents/skills/bench-what-next/SKILL.md`
- `.agents/skills/bench-write-spec/SKILL.md`
- `internal/anchors/registry_data.go`
- `internal/conformance/skills_index_checks_test.go`
- `internal/conformance/registry_test.go`
- `tests/canary/workflow-guidance-anchors/`
- `tests/canary/skills-index-command-adapters/`
- `.bench/BENCH-reference.md`
- `projects/benchkit.md`
- `CHANGELOG.md`
- `specs/ft228-debug-restoration/`
- `capture/session-handoff.md`

## Out of scope

- **Moving the Bench-mechanics tail behind a pointer** (the remaining half of
  audit action A5): 4 edits (the phase file, a new reference, ~7 anchor
  retargets in `registry_data.go`, 2 fixture rewrites), 2 gate runs. Its own
  capability: relocating graded text without weakening the grading.
- **Flipping any other phase's Codex invocation posture:** 2 edits (the table
  row, that phase's `openai.yaml`), 1 gate run each — reviewer decisions the
  table now makes one-line.
- **The repair-loop tripwire** (audit A10, the alternative trigger the
  explicit-only path would have leaned on): 3 edits (gate-record projection,
  status row, test), 1 gate run.
- **Shipping a HITL loop script template:** 2 edits (the script, its manifest
  row), 1 gate run.
- **The measurement harness for debug-discipline arms D/E/F** (audit A11):
  its own spec; nothing here blocks it.

## Further notes

FT24 stays parked pending upstream: the Codex changelog adding a spawn tool
name or a deny-capable SubagentStart is its graduation trigger, and this spec's
Codex settle (implicit invocation for one phase) neither needs nor provides
that surface.

The restored text is adapted from the pinned upstream `diagnosing-bugs` at
`reference-skill-repos/skills@d574778f`, read in full during authoring, side by
side with the audit's clause-by-clause ledger entry L-11. The upstream
`scripts/hitl-loop.template.sh` was read and deliberately not imported.

Open reviewer call, flagged per the promise-inflation guard: story 12 / DP3 —
the prose-budget row — is a promise beyond `roadmap/FT228.md` (the audit defers
volume work to the measurement arms). Keep or cut at sign-off; the spec carries
it as staged.

`decisions/assets/gate-pipeline-fixture-inventory.md` records per-family
fixture counts this build makes stale. It stays a dated snapshot; it is not
gate-graded and no ticket refreshes it.
