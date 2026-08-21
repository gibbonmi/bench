# Retro — learnings-dated-line-visibility

## Outcome

`capture/learnings.md` no longer reports a false zero for content its parser
cannot interpret. Landed at `954c34de` from the reviewed integration source
`a1cb0ed5..3c8fd1ea`, four gated commits composed into one published commit, and
the spec flipped to `Status: implemented` by that landing.

`learnings.Parse` gained two rules. A line that leads with a `YYYY-MM-DD` date
but is not a `## ` heading becomes a malformed record naming its line. Content
below the scaffold's `<!-- entries below -->` marker that belongs to no entry
becomes one record per contiguous run, naming the run's first line. The marker
moved into `internal/learnings` as `JournalEntriesMarker`, and
`internal/adopt`'s scaffold and `internal/conformance`'s docs-reference check
both read it from there, so the boundary the parser enforces is the boundary a
fresh repo receives.

No consumer needed a production edit. `bench learnings`, `bench roadmap
--context`, and `roadmap.learningCount` already read the malformed list, so all
three surfaces flipped on their own. Production changed in exactly two files:
`internal/learnings/learnings.go` and `internal/adopt/init.go`.

The observed red is flipped, demonstrated through the built binary rather than
through a test: the pre-drain journal that printed `learnings[0]{date,title}:`
at exit 0 now prints two `line <n>` rows and exits 1, and a scaffolded journal
with one undated note appended below its marker prints its row and exits 1 while
`bench status` renders `unknown (capture/learnings.md is malformed)` instead of
`0 open learning(s)`.

36 acceptance rows, all covered, none excused except story 26, whose control-byte
exposure is pre-existing and identical on the malformed-heading path.

## Gate-stage timings

Measured on the landing's own gate run.

| phase | elapsed |
| --- | --- |
| gofmt | 0.07 s |
| vet | 1.3 s |
| test | 63.9 s |
| race | 4.9 s |
| system | 4.0 s |
| shellcheck | 0.5 s |

The `test` phase is 85% of the wall clock. Seven full gate runs were paid this
session — two capture commits, three ticket commits, one spec commit, one
landing — so roughly seven and a half minutes went to the same test phase
re-running over diffs that touched five packages at most.

## Ticket-versus-spec-slice and delegate performance

Three ticket-sized charges, all Opus/medium, all first-pass accepted at the diff
level. None needed a returned write pass.

The two build tickets were genuine tracer bullets and behaved like it. Ticket 01
carried the dated rule from journal bytes through all three consumer surfaces
and proved each without editing one. Ticket 02 did its prefactoring first — the
marker export across three files — then added the second rule on top, and its
delegate reported the anchor's residues (marker below a real heading, no real
heading at all, second marker) as falling out of a single forward scan rather
than as special cases, which is what made the Spec axis able to verify them
structurally.

Ticket-sized charges outperformed what a spec-slice charge would have bought
here, because the coverage rows did the compressing. Each charge carried its
rows with the why-it-catches clause attached, and both delegates used those
clauses to write red-capable tests rather than presence assertions.

Two delegates returned something the charge did not ask for. Ticket 01's
delegate flagged that the DL15 fixture it had to build was a second copy of the
scaffold's bytes and named ticket 02's DL29 as the place to reconcile it —
correctly, and it was reconciled there as a byte-equality mirror test. Ticket
02's delegate found its own mutation probe weak: dropping the placeholder
exclusion red only one command test, so it added a unit case and re-probed
rather than reporting a single-row guard for a load-bearing exclusion.

One delegate corrected an expectation against the tree instead of the other way
round: ticket 02's charge implied a `[open]` record at line 8, the tree emits
line 7, and the delegate pinned 7 after reading why. That is the right
direction.

## Coordinator catches

Every done-claim was probed independently, each at a site and kind different
from the delegate's own.

- Ticket 01, omission at `opensWithDate`: 14 cases went red across three
  packages, including all three quiet-posture rows. The delegate had probed a
  swap at the prefix walk.
- Ticket 02, swap at the run-tracking site (`runStart` set on every line rather
  than only when no run is open): DL23 and DL34 red, nothing else. The delegate
  had probed an omission at the anchor.
- Ticket 03, omission of the `marker < 0` guard: DL34 red alone. Then both of
  the ticket's own target mutations were re-run by hand — each reds its own row
  and only its own row.

The review's one accepted finding was reproduced before it was accepted: the
post-loop `flushRun()` call was deleted whole and all four packages stayed
green.

One finding came from the coordinator rather than any axis. `unaccountedRegion`
matches the marker after `strings.TrimSpace`, so a marker carrying trailing
whitespace still opens the rule — real behavior, no test, and no spec sentence.
It was pinned as DL36 rather than removed, because the leniency fails in the
loud direction and an exact match would let one invisible space disable the
diagnostic silently, which is the failure class this spec exists to close.

Two self-inflicted errors were caught by looking rather than by assuming. An
unquoted heredoc let the shell expand backticks inside a spec edit; the
assertion fired before the write, and `git diff` confirmed the file was
untouched. Earlier, `bench worktree path`'s `~`-prefixed output was consumed as
a shell path and `cp` built a stray `./~/` tree inside the repository — found by
inspecting the failure rather than by re-running the command.

## Repair attribution

| ticket | repair rounds | cause per round |
| --- | --- | --- |
| 01-report-dated-lines-that-miss-heading-shape | 0 | none |
| 02-report-unaccounted-content-below-the-entries-marker | 1 | spec-row |
| 03-pin-the-end-of-input-flush-and-the-marker-tolerance | 0 | none |

Ticket 02's round is charged to `spec-row`, not to its delegate. The coverage
map had no row for the walk's end-of-input flush, so no charge could have asked
for it; the map's own DL20 covered the no-trailing-newline edge for the dated
rule only, and the edge inventory claimed the class as covered on that basis.

The spec itself took two review iterations before any build started. Round 1's
blocking finding was that the rule's anchor contradicted the shipped scaffold's
byte order, which made two coverage rows mutually unsatisfiable and left the
unaccounted rule with no reachable red.

## Agent-experience improvements

### Bench CLI

- `bench worktree path` prints a `~`-prefixed portable path. A shell does not
  expand it inside quotes, so the obvious consumption — `cp x "$(bench worktree
  path L)/"` — silently writes into a literal `~` directory inside the repo. One
  safe consumer, several unsafe ones. Either emit an absolute path or name the
  expansion in the command's own help.
- `bench worktree land --source-tip` refused a short SHA with
  `refused{detail=worktree source tip mismatch}`, the same message a genuinely
  moved tip produces. Distinguishing "not a full object id" from "tip moved"
  would have saved a diagnosis step.
- `bench worktree land` names `bench worktree release` as the next action for a
  retained worktree, but `release` has no `--discard-ignored` flag; the residue
  is cleared by `bench worktree clean`. The refusal should name the command that
  actually accepts the flag it demands.
- `bench spec retire` closes with `next: ... remove the ROADMAP row`, but
  `/bench-final-check`'s post-merge tail says the opposite in the same breath —
  "Leave the roadmap and capture rows to `/bench-drain`". One of the two is
  wrong; the phase text was followed and FT243's row was left standing.

### Skills

- `craft-spec` should require that a rule anchoring on a literal is walked over
  the real output of whatever generator emits that literal, before its rows are
  locked. This spec's round-1 blocking finding was exactly that omission: the
  author read the scaffold for the literal and the parser for the placeholder
  helper, but never traced the anchor across the scaffold's actual line order.

### Process

- `/bench-write-spec` should run `bench preflight build <slug>` before it exits.
  This spec reached the build phase with `rows-owned` red — its ticket cited no
  row IDs — so the red landed in the phase that does not own the ticket file.
- The `test` phase dominates every gate run at ~64 s, and this landing paid it
  seven times over diffs touching at most five packages. Worth measuring whether
  a per-commit gate can narrow to affected packages while the landing keeps the
  full run.
