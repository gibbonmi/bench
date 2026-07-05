# What-Next Roadmap Drain

<!-- command-currency: historical -->

Status: staged

## Problem

Bench has two maintenance paths for the same prioritization material. Parked ideas
live in `ROADMAP.md`, assessed priority lives in `ASSESSMENT.md`, and open
learnings leave through the old learnings-integration phase. A cold session has
to reconcile those surfaces by hand before it can know what to build next, and
the CLI can only print parked lines instead of the maintained plan.

## Solution

Make `ROADMAP.md` the working prioritization document, make `IDEAS.md` the
capture inbox, and add `/bench-what-next` as the single maintenance phase that
reconciles shipped work, drains ideas and open learnings into the roadmap, and
refreshes the recommended sequence. `bench roadmap` stays deterministic: it
prints the roadmap plus drain status when the inbox or journal needs attention,
or prints the roadmap's recommended sequence as the call to action when nothing
needs draining.

## User stories

1. As an agent parking an out-of-scope idea, I want `bench idea "<text>"` to
   append the existing dated two-space line format to `IDEAS.md`, so capture is
   an inbox instead of the working roadmap. Line: claude-sonnet-5 / low. This is
   an exact CLI retarget at a known file-effect seam with existing runtime
   contract coverage.
2. As an agent checking priorities, I want `bench roadmap` to print
   `ROADMAP.md` and a drain-status block when `IDEAS.md` has parked lines or the
   learnings journal has open entries, so I know to run `/bench-what-next` before
   trusting the sequence. Line: claude-opus-4-8 / medium. The state machine and
   wording are specified, but the status block shape has reviewer-approved
   latitude in the map.
3. As an agent checking priorities when there is nothing to drain, I want
   `bench roadmap` to print the roadmap and its `## Recommended sequence`
   section verbatim, so the CLI gives a deterministic next action without doing
   judgment. Line: claude-opus-4-8 / medium. This is covered at the CLI seam, but
   it introduces a new format contract and extraction path.
4. As an agent in a repo without a usable roadmap contract, I want
   `bench roadmap` to report the missing roadmap or missing recommended-sequence
   section explicitly and exit zero, so absence is a maintenance prompt rather
   than a crash or silent empty output. Line: claude-opus-4-8 / medium. The
   error posture is caller-visible and should be pinned before implementation.
5. As an agent starting a session, I want `bench status` to count parked ideas
   in `IDEAS.md` and open learnings entries as one maintenance signal pointing
   to `/bench-what-next`, so the dashboard has one drain action instead of two
   competing exits. Line: claude-opus-4-8 / medium. Status wording is
   gate-observed, and the story changes signal composition across two packages.
6. As a maintainer of this kit repo, I want the dogfood files migrated so
   current assessment content becomes `ROADMAP.md`, current parked lines become
   `IDEAS.md`, and `ASSESSMENT.md` is deleted, so the repo itself uses the new
   contract. Line: claude-opus-4-8 / medium. This is data migration plus
   stale-reference cleanup, and the gate observes the resulting files rather than
   a pure unit seam.
7. As a reviewer running maintenance, I want `/bench-what-next` to reconcile
   shipped or stale roadmap items, fully drain `IDEAS.md`, verdict every open
   learning, and propose one batch diff for approval, so the roadmap and journal
   have one reviewed exit. Line: claude-fable-5 / high. This is command guidance
   prose that steers future sessions, so the project profile's leverage override
   applies.
8. As a reviewer, I want the old learnings-integration command and its Codex adapter
   references retired, so the learnings journal cannot appear to have two
   authoritative drains. Line: claude-fable-5 / high. Retiring a workflow command
   touches shipped guidance and adapter surfaces, where ungated semantics can
   drift through future sessions.
9. As an agent shaping work from priorities, I want `/bench-shape-idea` to start
   cold from top `ROADMAP.md` items and leave roadmap rows in place until shipped
   retirement, so the roadmap is the working plan rather than a capture sink.
   Line: claude-fable-5 / high. This rewrites workflow guidance that future
   agents follow, so it uses the guidance-prose line.
10. As the gate owner, I want conformance and canary anchors to pin the new
    command, adapter, CLI inventory, stale-reference, and roadmap-promotion
    contracts, so a partial rename or weakened prose cannot pass. Line:
    claude-opus-4-8 / medium. Gate and conformance logic is the oracle surface,
    but the expected checks are concrete.
11. As a maintainer retiring shipped specs, I want spec retirement to remove the
    corresponding roadmap row in the same promote-then-delete pass, so presence
    in the roadmap remains the current status. Line: claude-fable-5 / high. This
    is workflow guidance that changes future cleanup behavior and is weakly
    graded except by prose anchors.

## Implementation decisions

- `internal/roadmap` remains the deep module for idea capture, roadmap rendering,
  drain counts, and recommended-sequence extraction. It hides file discovery and
  section parsing behind the `bench idea` and `bench roadmap` subcommands.
- `bench idea` retargets only the sink file. It keeps the current text joining,
  dated line format, missing-newline normalization, creation behavior, usage
  error, and write-failure structured error posture.
- `bench roadmap` becomes a small state machine over the roadmap document and
  drain sources. It never performs judgment and never mutates files.
- The `## Recommended sequence` format is the deterministic handoff contract:
  two or three numbered lines, each naming the item and the phase command to run.
  The CLI extracts that section verbatim after confirming both drain sources are
  empty.
- `internal/status` should not re-parse ideas or learnings independently. It
  consumes the same counters as `bench roadmap`, preserving one source per fact.
- `/bench-what-next` is a workflow command, not a CLI subcommand. It proposes an
  uncommitted batch diff that the reviewer approves or adjusts before commit.
- The old learnings-integration phase is retired rather than kept as an alias. Any
  durable synthesis discipline it still owns is moved into the what-next phase or
  referenced from the existing synthesis skill as the rule-shaped item path.
- The dogfood migration is part of the implementation, not a follow-up. The
  kit's own roadmap files must demonstrate the new contract.
- Command-adapter skills remain derived from command files. Adding
  `/bench-what-next` requires the matching Codex adapter metadata; retiring
  the old learnings-integration command removes its adapter and all guide
  references.
- The shell wrapper and Go dispatch both remain part of the CLI contract. Any
  new file-level behavior for `idea`, `roadmap`, or `status` must still route
  through the existing `bin/bench.sh` names and the compiled core command map.
- `bench init` scaffold text for `.bench/learnings.md` is part of the surface:
  it must name `/bench-what-next` as the journal exit after the retirement.
- Spec retirement removes the roadmap row because the row's presence is status.
  History stays in git through the retirement commit, not in checked-off roadmap
  lines.

## Testing decisions

- Good tests exercise caller-visible behavior at the runtime CLI seam and the
  conformance seam. Unit tests are acceptable when they pin pure parsing or
  counting helpers that the CLI and status both consume.
- Runtime contracts under `bench gate` should cover the observable CLI outputs,
  exit codes, and file effects for `bench idea`, `bench roadmap`, and
  `bench status`.
- Conformance and canary fixtures should cover command/adaptor guide wiring,
  stale command references, workflow anchors, and the new roadmap promotion and
  retirement guidance.
- The full gate is `.bench/gate.sh`.

### Seam diagram

Seam 1 - roadmap CLI and status contract:

    trigger: bench idea, bench roadmap, or bench status
        |
        v
    args + repo files + learnings journal -> [ internal/roadmap + internal/status ] -> stdout, exit code, file effects
                                               ^ tests attach here: Go/runtime contract fixtures drive the CLI and inspect stdout/files

Seam 2 - workflow guidance and adapter conformance:

    trigger: bench gate docs/conformance phase
        |
        v
    command files + adapter skills + guide prose -> [ conformance and canary checks ] -> green or targeted diagnostic
                                                     ^ tests attach here: conformance tests and deliberately broken canaries

Seam 3 - dogfood roadmap files:

    trigger: implementation migration and later /bench-what-next runs
        |
        v
    ASSESSMENT.md + old ROADMAP.md -> [ repo-level migration rules ] -> ROADMAP.md + IDEAS.md, no ASSESSMENT.md
                                      ^ tests attach here: gate stale-reference/docs checks and targeted file assertions if needed

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | `bench idea "<text>"` appends one dated line to `IDEAS.md`, creates it if absent, and does not create or append `ROADMAP.md` | roadmap CLI runtime contract | first red: add a runtime fixture expecting `IDEAS.md` after `bench idea "ship dark mode"`; current implementation writes `ROADMAP.md` | file-effect assertion proves capture was retargeted instead of only docs changing |
| 1 | empty or whitespace-only `bench idea` exits 2 with the usage line and leaves `IDEAS.md` absent | roadmap CLI runtime contract | existing empty-idea tests retargeted to `IDEAS.md` should fail if a blank inbox line is created | preserves the existing usage contract while changing the sink |
| 1 | unquoted multi-word args and no-trailing-newline append normalization still work against `IDEAS.md` | roadmap CLI runtime contract | retarget existing multi-word and newline-normalization cases from `ROADMAP.md` to `IDEAS.md` | catches a retarget that drops the hostile-input behavior already covered today |
| 2 | with roadmap present and either source non-empty, `bench roadmap` prints the roadmap plus a drain-status block naming counts and `/bench-what-next` | roadmap CLI runtime contract | first red: fixture with `ROADMAP.md`, one `IDEAS.md` line, and one open learning expects the drain block; current command prints only `ROADMAP.md` | exact stdout proves the CLI sees both drain sources without doing judgment |
| 3 | with both sources empty, `bench roadmap` prints the roadmap and the `## Recommended sequence` section verbatim as the call to action | roadmap CLI runtime contract | first red: fixture with a valid recommended section and empty sources expects the numbered lines in the output; current command has no extraction branch | a parser or summarizer cannot pass if it omits or rewrites the source section |
| 4 | absent `ROADMAP.md` exits 0 with a pointer to run `/bench-what-next` | roadmap CLI runtime contract | first red: replace the current absent-file `roadmap empty` expectation with the missing-roadmap pointer | distinguishes missing working document from empty capture inbox behavior |
| 4 | `ROADMAP.md` lacking `## Recommended sequence` exits 0 with an explicit missing-section message | roadmap CLI runtime contract | first red: fixture with a roadmap body but no heading expects the missing-section message | prevents silent empty output when the format contract is broken |
| 5 | `bench status` shows one combined maintenance row when `IDEAS.md` or open learnings require draining, and the action is `/bench-what-next` | status runtime contract | first red: status fixture with one idea and one open learning expects `/bench-what-next`; current status points learnings to the old learnings-integration phase and ideas to `bench roadmap` footer | proves the dashboard no longer advertises two exits from one drain |
| 5 | when no ideas or open learnings exist, status does not show a what-next row only because `ROADMAP.md` exists | status runtime contract | first red: clean fixture with a valid roadmap and empty sources expects no drain row | catches a noisy implementation that treats the working roadmap itself as pending capture |
| 6 | repo migration leaves `ROADMAP.md` in assessment format with `## Recommended sequence`, moves old parked lines into `IDEAS.md`, and deletes `ASSESSMENT.md` | dogfood roadmap files | first red: implementation check before migration sees `ASSESSMENT.md` present and no `IDEAS.md` | pins the dogfood state, not just library behavior |
| 7 | `/bench-what-next` command states reconcile-first, full inbox drain, open-learning verdicts, batch-propose approval, and commit-on-green | workflow guidance conformance | first red: new workflow-anchor check for these required phrases fails before the command exists | catches a phase file that omits one of the closed map decisions |
| 8 | old learnings-integration command and Codex adapter are removed, and no live docs reference the retired slash command or retired Codex adapter token | workflow guidance conformance | existing stale command and Codex adapter reference checks go red after deleting the command until every reference is updated | ensures the old drain cannot survive as a dangling or advertised surface |
| 8 | `bench init` scaffolds a learnings journal whose header says entries leave only through `/bench-what-next` | safe-link and init runtime contracts | first red: fixture after `bench init` expects the new journal header; current scaffold names the old learnings-integration phase | catches a stale consumer-facing scaffold even if live repo docs were swept |
| 9 | `/bench-shape-idea` cold start reads top roadmap items rather than parked ideas, and no longer removes a roadmap row on map promotion | workflow guidance conformance | first red: update roadmap-promotion anchor checks to expect working-roadmap language and row persistence | catches stale capture-sink wording that would make roadmap presence mean two different statuses |
| 10 | `.bench/BENCH.md`, `.bench/BENCH-reference.md`, command inventory, and command adapters list `/bench-what-next` and omit the retired command | workflow guidance conformance | existing command-guide, CLI-list, and adapter checks fail on any missing new command reference or leftover retired reference | pins the shipped guide and adapter surfaces together |
| 10 | canary fixtures prove the new what-next anchors and stale-reference checks bite | workflow guidance conformance | first red: add broken canaries that remove the what-next anchor or keep a retired-command reference; gate must fail with targeted substrings | proves the new enforcement is not an always-green structural check |
| 11 | spec-retire guidance says the corresponding roadmap row is removed in the same promote-then-delete commit | workflow guidance conformance | first red: workflow-anchor check for the roadmap-row retirement phrase fails before command prose is updated | prevents implemented work from staying listed as current roadmap status |
| Edge of 2 | `IDEAS.md` absent and present-but-empty are both zero drain count | roadmap CLI runtime contract | first red: table-driven fixtures cover absent and zero-byte inbox with a valid roadmap | protects the full-drain empty state from treating a missing inbox as an error |
| Edge of 2 | missing or unreadable learnings journal counts as zero for roadmap/status display, while `bench learnings` keeps its own structured behavior | roadmap CLI runtime contract | first red: fixture omits `.bench/learnings.md` and expects no drain count from the journal | keeps dashboard maintenance robust in consumer repos with old scaffolds |
| Edge of 3 | recommended-sequence extraction preserves numbering and phase-command text but does not require the CLI to validate item count beyond reporting malformed section explicitly | roadmap CLI runtime contract | first red: fixture with one numbered line or four numbered lines expects a malformed-section message | makes the format contract visible without moving prioritization judgment into the CLI |
| Edge of 5 | status uses the same counters as roadmap, including `IDEAS.md` lines and learnings open-heading parser | status runtime contract | first red: fixture with template learnings heading plus one real open heading expects count 1 | catches a second parser that counts examples or misses real open entries |
| Edge of 8 | stale command references inside shipped guidance surfaces - command files, skill files, guide docs, and adapter metadata - are all swept; the what-next spec and decision map opt out via the check's `command-currency: historical` marker as the record of the retirement | workflow guidance conformance | existing docs-currency stale-reference checks plus a retired-command canary | prevents the old command name from surviving in a shipped artifact while the documents that ordered the retirement stay green until spec-retire |
| 10 | `.bench/BENCH.md` Capture section names `IDEAS.md` as the capture sink, including the no-PATH fallback append instruction | workflow guidance conformance | first red: anchor check expecting IDEAS.md capture wording fails against the current ROADMAP.md-as-sink prose | stale sink wording would send agents appending raw capture into the working roadmap, the exact failure this build removes |

### Edge inventory

- Missing `IDEAS.md`: covered by Edge of 2.
- Present-but-empty `IDEAS.md`: covered by Edge of 2.
- Hand-edited `IDEAS.md` with no trailing newline: covered by story 1.
- Unquoted multi-word idea text: covered by story 1.
- `IDEAS.md` write failure: covered by story 1's existing write-error posture; add a
  focused unit case only if the current OpenFile error path changes.
- Missing `ROADMAP.md`: covered by story 4.
- Present roadmap without `## Recommended sequence`: covered by story 4.
- Malformed recommended sequence with fewer than two or more than three numbered
  items: covered by Edge of 3.
- Empty learnings journal, missing learnings journal, template heading, resolved
  heading, CRLF heading, and no-trailing-newline heading: covered by Edge of 5
  and existing learnings parser tests.
- Current old `ROADMAP.md` lines with parked items: covered by story 6 migration.
- Existing assessment content with shipped-work notes: covered by story 6 migration
  and the `/bench-what-next` reconcile rule in story 7.
- Stale retired learnings-integration slash-command or Codex-adapter references in
  shipped guidance surfaces (command files, skills, guide docs, adapter metadata):
  covered by story 8 and Edge of 8. The what-next map and this spec carry the
  `command-currency: historical` marker - the check's designed per-file opt-out -
  because they record the retirement and the not-yet-built command; they leave
  through spec-retire, not the sweep.
- Stale ROADMAP.md-as-capture-sink wording in shipped guidance, including the
  Capture section's no-PATH fallback append instruction: covered by the story 10
  capture-sink row.
- Stale retired learnings-integration references in scaffolded `.bench/learnings.md`
  text for newly initialized repos: covered by story 8.
- Command adapter missing for `/bench-what-next`: covered by story 10.
- Canary fixture copies under workflow-guidance anchors carrying old
  `/bench-shape-idea` or retired-command wording: covered by story 10 and Edge of
  8.
- A separate `bench what-next` CLI subcommand: won't handle - rejected by the
  closed map; the deterministic surface is `bench roadmap` and the judgment
  surface is `/bench-what-next`.
- Partial inbox drain: won't handle - rejected by the closed map; `IDEAS.md` is a
  pure inbox and empties to zero on each approved drain.
- Interactive per-item verdicts during `/bench-what-next`: won't handle - rejected
  by the closed map; the batch diff is the verdict sheet.
- Completion checkmarks in `ROADMAP.md`: won't handle - rejected by the closed
  map; roadmap row presence is status and history lives in git.
- Concurrent writer during a drain: won't handle in code - the existing worktree
  discipline owns this operational edge, and at least one in-scope caller still
  exercises the drain in a single checkout.

## Out of scope

No deliberate cuts. Every closed Handoff item is either part of this build or a
rejected alternative recorded in the edge inventory.
