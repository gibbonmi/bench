# Retro: kit-guidance-fold

## Outcome

The landing commit `6e51b1ec` published `specs/kit-guidance-fold/spec.md` as
`Status: implemented` and closed FT158, FT236, FT259, FT269, and FT273. The reviewed
pair was base `dab7ca38` and source tip `2030391d`. Five tickets and one repair ticket
landed through the retained integration source `kit-guidance-fold`. The diff added 28
anchor tuples, five registry test functions, and one extended absence test. It also
added 27 live-mirror canary fixtures, two guard table rows, and one new reference file.

Five reviewer-approved light-path fixes landed on `main` during the build. They are
the `bench test` `elapsed_ms` column (`c7c8588a`), the `gate-prose` help example
(`c8413ded`), and the fixture-bite diagnostic projection (`fb1498fa`). The other two
are the named preview cap (`471699d2`) and `bench repair --help` at exit 0
(`dab7ca38`).

## Gate-stage timings

The landing gate ran once over the composed pair. Each merge from `main` paid one
prospective gate of the same shape.

| phase | elapsed |
| --- | --- |
| gofmt | 117 ms |
| vet | 934 ms |
| test | 56.2 s |
| race | 2.6 s |
| system | 21.7 s |
| shellcheck | 511 ms |

## Ticket-versus-spec-slice and delegate performance

Every charge was ticket-sized. Five write charges ran at Opus medium, and all five
landed first-pass on behavior. Each returned a red-before and green-after log per row
and a self-probe that bit. Two charges reported an out-of-fence demand and reworded
their sentence instead of editing the fixture.

Five light-path fix charges ran at Opus medium or low. Four landed first-pass, and one
took one repair round for a test in a package outside its search list. The repair
charge folded nine findings in one pass and took one more round for a row citation.
Three review axes, one Codex falsification pass, and two scoped re-reviews ran at the
mid tier. The axes returned 19 raw findings that collapsed to nine repair targets.

## Coordinator catches

- A coordinator probe ran as a heading rename at each ticket, kept constant across the
  batch, and bit through the owning check every time.
- The preview-cap delegate searched three packages for the old literal, and a fourth
  package pinned it. The landing gate found it, and the resumed delegate repaired it.
- One delegate reported the follow-on hook refused a trailing `; echo` after an exec
  segment. Two coordinator runs of the same shape passed, so the report did not
  reproduce.
- The merge gate went red on a shipped-surface sweep that read "includes" beside
  `tests/` as a package claim. One word changed, and the merge went green.
- A worktree test reported a release refusal once during a merge gate and never again.
  It was classified as environmental.

## Repair attribution

| ticket | repair rounds | cause per round |
| --- | --- | --- |
| add-the-both-ends-rule-to-craft-gate | 0 | none |
| name-the-repair-ticket-owner-before-re-review | 0 | none |
| make-the-falsification-pass-standing-for-guidance-diffs | 1 | spec-row |
| fold-the-probe-oracle-and-fence-rules-into-delegation-discipline | 1 | tree-drift |
| add-the-finding-discipline-reference-to-craft-review | 1 | spec-row |
| repair-the-review-findings | 1 | delegate-error |

## Agent-experience improvements

### Bench CLI

- Give `bench worktree exec` a probe form that copies a file aside, applies one swap,
  runs the named check, and restores the file with a byte confirmation. The
  census entry `kit-guidance-fold landed with census 48` records the heads it removes.
  Feeds: new
- Add `bench test --check <name> --fixtures` so a charge author reads the canary
  family a check owns from the CLI instead of from the fixture-bite test.
  Feeds: new
- Keep the whole body of a multi-line failure message under `bench test --full`, so
  a fixture-bite red is readable inside the boundary.
  Feeds: new
- Render `bench preflight` closure diagnostics as a `closure[N]{ticket,path,requires}`
  table instead of one prose cell.
  Feeds: new
- Attribute a compile error directly when `bench test --run` matches no test.
  Feeds: new
- Default the `bench gate-prose` root to `.` when every operand is a file.
  Feeds: FT270

### Skills

- Add the shipped-surface claim words to the `craft-spec` reader sweep, so a guidance
  sentence that names `tests/` beside "includes" is caught before the merge gate.
  Feeds: new
- Name the closest pinning package in a cap-change charge's search list, because a
  literal can be pinned outside the packages that consume it.
  Feeds: none

### Process

- Land a light-path fix before a spec's final merge only when its `CHANGELOG.md` entry
  sits under a heading no sibling touches. Two entries under one heading conflict at
  composition.
  Feeds: none
- Set the review base to the `main` tip merged into the source before a spec landing,
  so the reviewed range holds the spec diff alone.
  Feeds: none
