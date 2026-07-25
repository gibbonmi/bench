# FT86 — fail-closed control records and single-sourced repository facts

Status: staged

## Problem

Every control-record read in the kit collapses failure to success. `bench
learnings` on an unreadable `.bench/learnings.md` prints an empty table and
exits 0 — byte-identical to a repo that genuinely has no journal. `bench maps`
silently drops a decision file it cannot read, from both the listing and the
count the dashboard shows. `bench learnings` parses malformed `## ` headings and
throws them away. `bench status` reports a fabricated zero for every one of
those. Only `roadmap --context` distinguishes absent from malformed, and it does
so with its own private reader.

The same fabrication runs through the default-branch fact: `git.DefaultBranch`
returns the literal string `"main"` when `origin/HEAD` is unset, so a repo whose
default is `master`, or one with no resolvable default at all, gets a confident
wrong answer that `bench diff` then uses to compute a review base.

An agent reading these surfaces cannot tell "nothing to report" from "I could not
read it." That is the whole failure: a silent empty state is indistinguishable
from an authoritative one, and every consumer downstream — the dashboard, the
drain, the review base — inherits the lie.

Separately, `bench coverage --check` validates almost nothing: a spec with no
acceptance-coverage map at all passes green, and story references validate only
against the highest story number, so story `0`, a reversed range `3-1`, and a
reference to a story number the spec skips all pass.

## Solution

One classified reader owns file state, and every surface takes an explicit,
documented posture on each state.

`internal/bounds` becomes the single classified reader: it Lstats first, rejects
special files before reading, keeps the existing size bound, and returns one of
six typed states — `absent`, `empty`, `malformed`, `unreadable`, `wrong-type`,
`unsupported-schema`. Absence is the only authoritative empty state.

Query commands fail closed. `bench learnings`, `bench maps`, `bench roadmap`,
and `bench coverage` exit 1 with a structured `error:` line naming the state when
a read is anything but absent. Where a failure is per-entry or per-file rather
than whole-document — a malformed learning heading, an unreadable decision file —
the command shows every well-formed row *plus* an explicit row naming the broken
one, and still exits 1.

`bench status` keeps exit 0, because the dashboard's job is to render whatever it
can, but it renders an explicit `unknown` row for a signal it could not compute
instead of a fabricated zero — the same surfaced-degradation pattern its
worktree-list failure row already uses. `roadmap --context` likewise keeps its
snapshot and degrades the affected block to `unknown`.

`git.DefaultBranch` is deleted. `git.ResolvedDefault` becomes the sole owner of
the default-branch fact, and every caller handles `ok=false` explicitly. `bench
diff` fails closed naming reality — no resolvable default — and names the
`branch.<name>.benchBase` git-config escape hatch.

`bench coverage --check` requires a map or an explicit
`<!-- coverage-map: historical -->` marker, and validates story references as
exact membership in the declared story set, positive, and forward-ranged.

## User stories

1. As an agent reading a control record, I want one classified reader in
   `internal/bounds` that returns a typed state (`absent`, `empty`, `malformed`,
   `unreadable`, `wrong-type`, `unsupported-schema`) alongside the bytes, so that
   every surface classifies file state by one rule instead of five partial ones.
   Line: gpt-5.6-terra / medium. This is the deep unit the whole spec rests on and
   a wrong classification silently poisons all six consuming surfaces, so it gets
   the mid tier even though the seam is known.

2. As an agent reading a control record behind a broken symlink, I want the
   classifier to Lstat before reading so a dangling link classifies as
   `unreadable` and never as `absent`, so that a broken link cannot masquerade as
   an authoritative empty state. Line: gpt-5.6-terra / medium. `os.ReadFile`
   reports a dangling symlink as ENOENT, so this is precisely the trap the
   classifier exists to close and it belongs with the classifier's own tier.

3. As an agent whose repository contains a FIFO, device node, or socket where a
   control record should be, I want the classifier to return `wrong-type` without
   opening the path, so that static inspection cannot block on a file that never
   yields EOF. Line: gpt-5.6-terra / medium. Rejecting before the open is the
   correctness property here and it is decided inside the classifier, so it stays
   on the classifier's tier.

4. As an agent reading a control record whose permissions deny me, I want
   `unreadable` with the underlying reason preserved, so that a permission
   problem reads as a permission problem rather than as an empty file.
   Line: gpt-5.6-luna / medium. The classification rule is already fixed by story
   1 and the fixture pattern (chmod 0o000 plus cleanup) already exists in the
   status worktree test.

5. As an agent reading a control record that is a symlink to a regular file, I
   want the link followed and the target classified, so that a linked journal or
   roadmap is read rather than rejected. Line: gpt-5.6-luna / medium. This is the
   companion assertion to story 2 and follows directly from the same Lstat-first
   rule.

6. As an agent reading a control *directory* — `specs/`, `decisions/` — I want
   the same typed classification for directory reads, so that an unreadable
   directory is distinguishable from an empty one. Line: gpt-5.6-luna / medium.
   The state vocabulary is settled by story 1; this is the directory arm of the
   same contract.

7. As an agent reading a control record whose bytes are fine but whose structure
   no parser recognizes — a decisions file with no ticket headings, a roadmap
   with no recognizable rows — I want `unsupported-schema` with a reason,
   distinct from byte-level `malformed`, so that "I read it but it is not the
   document you think" is its own state. Line: gpt-5.6-terra / medium. Which
   shapes count as unrecognized is a judgment per parser rather than a mechanical
   mapping, so it takes the mid tier.

8. As an agent running `bench learnings` against an unreadable, wrong-type, or
   malformed journal, I want exit 1 and a structured `error:` line naming the
   state, so that I never read an empty table as "no open learnings."
   Line: gpt-5.6-luna / medium. Thin consumer mapping a classifier state to an
   exit code at a seam the AXI contract tests already cover.

9. As an agent running `bench learnings` against a journal containing malformed
   `## ` headings, I want every well-formed entry rendered plus an explicit row
   per malformed heading carrying its line and reason, and exit 1, so that a
   broken entry is visible instead of dropped. Line: gpt-5.6-luna / medium. The
   parser already collects malformed fragments; this is rendering them and
   setting the exit code.

10. As an agent running `bench maps` where a decision file cannot be read, I want
    a row naming that file and its state rather than a silent omission, and exit
    1, so that an unreadable map is not invisibly missing from the listing.
    Line: gpt-5.6-luna / medium. Same shape as story 9 at a sibling seam.

11. As an agent reading `bench status`'s decision-map signal, I want the
    unresolved-map count to carry its own readability state rather than
    collapsing a failed scan to zero, so that the dashboard and the listing
    cannot disagree about what was readable. Line: gpt-5.6-luna / medium. The
    count and the listing already share one engine; this threads the state
    through that single source.

12. As an agent running `bench roadmap` against an unreadable or malformed
    `ROADMAP.md`, I want exit 1 with a structured error naming the state, so that
    a read failure never prints as an empty working document.
    Line: gpt-5.6-luna / medium. Same posture as story 8, applied at the roadmap
    command.

13. As an agent running `bench roadmap --context` where one source fails, I want
    the snapshot to survive with that source's state rendered explicitly rather
    than the whole context aborting, so that one bad file does not cost me the
    entire repository picture. Line: gpt-5.6-luna / medium. `--context` already
    carries an absent/empty/malformed state vocabulary per source; this extends
    it and stops the whole-snapshot error return.

14. As an agent running `bench coverage --check` against a spec with no
    acceptance-coverage map and no `<!-- coverage-map: historical -->` marker, I
    want exit 1, so that an unmapped spec cannot pass the gate's docs layer by
    having nothing to validate. Line: gpt-5.6-luna / medium. Effort is medium
    rather than low because this is gate logic, per the profile's cached routing.

15. As an agent running `bench coverage --check` against a map whose story
    references are invalid, I want story `0`, a reference to a number outside the
    declared story set, and a reversed range such as `3-1` each to fail, so that
    membership is validated exactly rather than against a maximum.
    Line: gpt-5.6-luna / medium. Gate logic again, and the parser's row-walking
    structure already exists.

16. As an agent running `bench outline` where a listed file cannot be read or is
    not a regular file, I want that classification to come from the shared
    classifier and the command to exit 1, while a binary or oversized file keeps
    its visible skip row and exit 0, so that a genuine read failure is a failure
    and an un-outlineable file is not. Line: gpt-5.6-terra / medium. This splits
    outline's existing skip classes into failure and content classes, which is a
    posture judgment the map did not spell out (flagged below for veto).

17. As an agent reading `bench status`, I want any signal whose underlying read
    failed to render an explicit `unknown` row naming the state and path, at its
    normal severity and never suppressed by a zero count, while the command still
    exits 0, so that the ambient board degrades visibly instead of reporting a
    clean repository. Line: gpt-5.6-luna / medium. The failure-row pattern is
    already established by the worktree row; this applies it to the remaining
    signals.

18. As a maintainer, I want `git.DefaultBranch` deleted and `git.ResolvedDefault`
    to be the sole owner of the default-branch fact, with every caller handling
    `ok=false` explicitly, so that no code path can fabricate `main` for a
    repository that has no resolvable default. Line: gpt-5.6-terra / medium.
    There are eight call sites across `diff`, `status`, `adopt`, `worktree`, and
    `git` itself, and each needs its own posture decision rather than a
    mechanical substitution.

19. As an agent running `bench diff` in a repository with no resolvable default
    branch, I want exit 1 with a structured error naming that reality and the
    `git config branch.<name>.benchBase <sha>` escape hatch, so that I am told
    how to proceed instead of receiving a review base computed against a branch
    that does not exist. Line: gpt-5.6-luna / medium. The error shape and the
    hint text are decided; this is wiring them to the `ok=false` path.

20. As an agent running `bench roadmap --context` in a repository with no
    resolvable default branch, I want the git facts block to render `unknown`
    while the rest of the snapshot stays intact and the command exits 0, so that
    an unresolvable default costs me one block rather than the whole context.
    Line: gpt-5.6-luna / medium. Same degradation posture as story 13 at the git
    block.

## Implementation decisions

**The classifier lives in `internal/bounds`.** It is an extension of the
existing size-bounded `Read`, not a new package: `bounds` already owns the
kit's read policy, and a second package would split that ownership. It exposes a
typed result carrying the state, the bytes when a read completed, and a reason
string for the diagnostic states. `bounds.Read`'s existing `ReadStatus`
(`complete`/`oversized`/`failed`) stays as the stream-level primitive the new
path-level classifier composes; oversized remains distinguishable.

**Classification order is Lstat, then type check, then read.** Lstat first
because `os.ReadFile` reports a dangling symlink as ENOENT, which would
classify a broken link as authoritative absence. A symlink to a regular file is
followed and its target classified. A FIFO, device, or socket is `wrong-type`
and is never opened. A permission failure at any step is `unreadable` with the
underlying error preserved as the reason.

**`unsupported-schema` is shape-based, not marker-based.** No control record
grows a schema or version marker — that was rejected during shaping. The state
means the file read cleanly and is valid UTF-8, but its structure is not one the
consuming parser recognizes. Each parser owns the predicate for its own document;
the classifier carries the state so parsers can return it uniformly.

**Migrations.** `roadmap`'s private `readSource` and `readDirSource` are deleted
and their callers move to the classifier. `learnings`, `maps`, `status`, and
`outline` drop their local `os.ReadFile`/`os.ReadDir`/`Lstat` handling and
consume the classifier. No surface keeps reading logic of its own; each maps a
state to its posture.

**Postures, by surface.** Query commands (`learnings`, `maps`, `roadmap`,
`coverage`) exit 1 with a `toon.Errorf`-shaped `error:` line on any non-absent
failure; absent renders the empty state with exit 0. `bench status` exits 0 and
renders explicit `unknown` rows. `roadmap --context` exits 0 and degrades the
affected source or block. `bench outline` splits its skip classes: `unreadable`
and `wrong-type` are failures and exit 1; `binary` and `oversized` are content
classes that keep their visible skip row and exit 0.

**Status signal counts carry state.** `maps.UnresolvedCount` and the sibling
counters status consumes stop returning a bare integer that cannot express
failure. They return the count plus its readability state, so the dashboard's
`unknown` row and the listing's rows derive from the same scan — the count and
the listing must not be two derivations of one fact.

**TOON state vocabulary is reused, not reinvented.** The state names are exactly
`roadmap --context`'s existing `absent` / `empty` / `parsed` / `malformed`
vocabulary extended with `unreadable`, `wrong-type`, and `unsupported-schema`,
and the dashboard's degradation literal stays `unknown`, matching the worktree
list's existing usage.

**`git.DefaultBranch` is deleted outright.** `ResolvedDefault` absorbs the
`origin/HEAD` lookup it currently delegates. The eight call sites resolve as
follows. `diff` fails closed (story 19). `roadmap --context`'s git facts block
degrades to `unknown` (story 20). `status`'s two default-branch comparisons
(`appendRetirement`, `appendRoadmapReconcile`) already treat an unreadable
branch as "not the default branch" and skip their advisory signal; an
unresolvable default takes the same tolerant path, which is a posture the
existing code comments already record. `git.Facts` propagates the unresolvable
state into `RepoFacts` rather than filling `DefaultBranch` with a guess.
`PruneLandedBranches` and `LandedState` already error on `ok=false` and only
need their message to stop naming a fabricated branch. `adopt`'s hook and
staging paths, and `worktree`'s refresh and lifecycle paths, use the default
branch as a *candidate* ref that is already probe-guarded, so on `ok=false` they
skip the candidate rather than probing a fabricated name.

**Coverage validation.** `Check` gains a no-map rule: `State(p) == "no-map"` is a
violation unless the historical marker is present. Story references validate
against the exact set of declared story numbers rather than the maximum; a
reference to a number not in the set, a `0`, and a range whose end is below its
start each produce their own violation message. The existing violation phrasings
are matched downstream by substring, so new messages follow the same
`coverage map row %d …` shape and existing phrasings are not reworded.

**Flagged for veto — two calls the map did not close.**
(a) Story 16 splits `outline`'s skip classes into failure classes (exit 1) and
content classes (exit 0). The map listed `outline` in the fail-closed group, but
exiting 1 on every binary file in a repository would make the command unusable;
the split keeps the fail-closed property for genuine read failures. (b) The
`ROADMAP.md` row for FT86 names a `traversal` state the decision map's
enumeration does not include. This spec follows the map's six states and does not
add `traversal`; the caller-supplied-path surface it would guard (`bench
coverage <path>`) keeps its current behavior. If you want traversal classified,
it is a separate cut.

## Testing decisions

A good test here drives a real `bench` subcommand in a throwaway fixture
repository and asserts the exit code and the stdout string — never a package
internal. The one exception is the classifier itself, which is a deep unit whose
hostile-filesystem semantics are cheaper and sharper to assert directly than
through six surfaces.

Seams, and their prior art:

- **`internal/bounds` unit** — the classifier's six states against real fixture
  files (dangling symlink, FIFO, chmod 0o000, symlink-to-regular, empty,
  oversized). Prior art: `internal/bounds/bounds_test.go`.
- **`internal/contract/axi`** — every surface posture: exit codes and structured
  `error:` strings for `learnings`, `maps`, `roadmap`, `roadmap --context`,
  `coverage --check`, `outline`, `diff`, and `status`. Prior art:
  `axi_fail_closed_test.go`, `axi_outline_test.go`, `axi_coverage_test.go`,
  `axi_roadmap_context_test.go`.
- **`internal/conformance`** — the repository-wide facts: no surviving
  `DefaultBranch` symbol, and the `specs/*.md` sweep applying the new no-map
  rule. Prior art: `docs_workflow_checks_test.go`'s `checkCoverageMaps`.
- **`tests/canary/`** — proof the gate still bites. Prior art:
  `coverage-map-validation/broken-coverage-map` with its `EXPECT` file.

The gate command is the project gate: `.bench/gate.sh`.

### Seam diagram

**Seam 1 — the classifier (deep unit).**

    trigger: any surface needing a control record's bytes
        │
        ▼
    path ────▶  [ bounds: Lstat ▸ type check ▸ bounded read ]  ──▶  state
    limit ───▶  [                                           ]  ──▶  bytes
                                                               ──▶  reason
                  ◀ tests attach here: build a hostile fixture file
                    (dangling symlink, FIFO, mode 0o000, oversized)
                    and assert the returned state and reason

**Seam 2 — the AXI query surface (contract).**

    trigger: `bench <cmd>` in a fixture repo
        │
        ▼
    argv ───────▶  [ learnings | maps | roadmap | coverage ]  ──▶  stdout (TOON or error:)
    fixture FS ─▶  [ outline   | diff | status            ]  ──▶  exit code
                     ◀ tests attach here: run the built binary in a
                       throwaway repo whose control records are hostile,
                       assert exit code and stdout substring

**Seam 3 — conformance and canary (the gate biting).**

    trigger: `.bench/gate.sh`
        │
        ▼
    repo tree ──▶  [ conformance sweep + canary fixtures ]  ──▶  diagnostics
                                                            ──▶  gate verdict
                     ◀ tests attach here: a fixture spec with no map and no
                       historical marker, and a source tree still calling
                       DefaultBranch, must each turn the gate red

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | classifier returns six typed states with bytes and reason | unit (`internal/bounds`) | `go test ./internal/bounds -run TestClassifyStates` — fails, no such symbol | Asserting all six states in one table means an implementation covering only absent and complete cannot pass. |
| 1 | absent and present-but-empty are distinct states | unit (`internal/bounds`) | `go test ./internal/bounds -run TestClassifyAbsentVsEmpty` | Two fixtures whose only difference is existence; collapsing either onto the other fails one assertion. |
| 2 | dangling symlink classifies `unreadable`, not `absent` | unit (`internal/bounds`) | `go test ./internal/bounds -run TestClassifyDanglingSymlink` | A `ReadFile`-first implementation sees ENOENT and answers `absent`, which this assertion forbids by name. |
| 3 | FIFO classifies `wrong-type` without opening | unit (`internal/bounds`) | `go test ./internal/bounds -run TestClassifyFIFOWithoutOpen` | The test drives a FIFO with no writer under a deadline; an implementation that opens before type-checking blocks and fails the deadline rather than returning. |
| 3 | device and socket classify `wrong-type` | unit (`internal/bounds`) | `go test ./internal/bounds -run TestClassifyNonRegular` | `/dev/null` is readable and would classify `empty` under a type-blind implementation. |
| 4 | mode 0o000 classifies `unreadable` with the reason preserved | unit (`internal/bounds`) | `go test ./internal/bounds -run TestClassifyPermissionDenied` | Asserting a non-empty reason alongside the state forbids the implementation that discards the underlying error. |
| 5 | symlink to a regular file is followed and classified | unit (`internal/bounds`) | `go test ./internal/bounds -run TestClassifySymlinkFollowed` | A blanket "any symlink is unreadable" implementation passes story 2 and fails here. |
| 6 | directory classification returns absent, empty, and unreadable | unit (`internal/bounds`) | `go test ./internal/bounds -run TestClassifyDir` | An unreadable directory and an empty one both yield zero entries; only the state distinguishes them. |
| 7 | valid UTF-8 with unrecognized structure classifies `unsupported-schema` | unit (`internal/bounds`) | `go test ./internal/bounds -run TestClassifyUnsupportedSchema` | Asserting it against a byte-level `malformed` fixture in the same table forbids collapsing the two. |
| 8 | `bench learnings` exits 1 with a structured error on an unreadable journal | contract (`internal/contract/axi`) | `go test ./internal/contract/axi -run TestAXILearningsFailClosed` | A fixture with a mode-0o000 journal; the current implementation exits 0 with an empty table, which the exit assertion rejects. |
| 8 | `bench learnings` with no journal exits 0 and renders the empty state | contract (`internal/contract/axi`) | `go test ./internal/contract/axi -run TestAXILearningsAbsentIsEmpty` | Paired with the row above, this forbids the over-correction that fails closed on absence too. |
| 9 | malformed learning headings render as rows with line and reason, exit 1 | contract (`internal/contract/axi`) | `go test ./internal/contract/axi -run TestAXILearningsMalformedRows` | A fixture mixing one good and one malformed heading; asserting the good row is still present forbids a hard-fail that drops everything. |
| 10 | an unreadable decision file gets a row naming it, exit 1 | contract (`internal/contract/axi`) | `go test ./internal/contract/axi -run TestAXIMapsUnreadableFileRow` | The current implementation `continue`s past the read error; asserting the filename appears in stdout catches the silent drop. |
| 11 | the status map count and the `bench maps` listing agree on what was readable | contract (`internal/contract/axi`) | `go test ./internal/contract/axi -run TestAXIMapsCountMatchesListing` | One fixture with an unreadable map file, asserted through both surfaces; a count that fabricates zero disagrees with a listing that shows the row. |
| 12 | `bench roadmap` exits 1 with a structured error on an unreadable ROADMAP.md | contract (`internal/contract/axi`) | `go test ./internal/contract/axi -run TestAXIRoadmapFailClosed` | A mode-0o000 fixture; an empty-document print exits 0 and fails the assertion. |
| 13 | `roadmap --context` keeps the snapshot and marks one failed source | contract (`internal/contract/axi`) | `go test ./internal/contract/axi -run TestAXIRoadmapContextDegrades` | Asserting both exit 0 *and* the presence of unrelated sections forbids the whole-snapshot error return the code takes today. |
| 14 | a spec with no map and no historical marker fails `--check` | contract (`internal/contract/axi`) | `go test ./internal/contract/axi -run TestAXICoverageNoMapFails` | The current implementation returns nil violations for `no-map` and exits 0, which the exit-1 assertion rejects. |
| 14 | a spec carrying the historical marker still passes `--check` | contract (`internal/contract/axi`) | `go test ./internal/contract/axi -run TestAXICoverageHistoricalPasses` | Paired with the row above, this forbids the implementation that fails every unmapped spec including opted-out ones. |
| 15 | story `0`, a non-member story number, and a reversed range each fail | contract (`internal/contract/axi`) | `go test ./internal/contract/axi -run TestAXICoverageStoryMembership` | Three fixtures in one table, each currently passing; a max-only validator fails all three assertions. |
| 15 | a valid comma list and forward range still pass | contract (`internal/contract/axi`) | `go test ./internal/contract/axi -run TestAXICoverageValidStoryRefs` | Forbids the over-strict validator that rejects legitimate multi-story rows. |
| 16 | `bench outline` exits 1 when a listed file is unreadable or non-regular | contract (`internal/contract/axi`) | `go test ./internal/contract/axi -run TestAXIOutlineFailureClasses` | A fixture with a mode-0o000 tracked file; the current implementation emits a skip row and exits 0. |
| 16 | a binary or oversized file keeps its skip row and exit 0 | contract (`internal/contract/axi`) | `go test ./internal/contract/axi -run TestAXIOutlineContentClasses` | Forbids the blanket fail-closed that would make `bench outline` exit 1 in any repo containing an image. |
| 17 | a signal whose read failed renders an `unknown` row and status still exits 0 | contract (`internal/contract/axi`) | `go test ./internal/contract/axi -run TestAXIStatusUnknownRow` | Asserting the row is present *and* the exit is 0 forbids both the fabricated zero (row absent) and the fail-closed over-correction. |
| 17 | the `unknown` row is not suppressed by a zero count | contract (`internal/contract/axi`) | `go test ./internal/contract/axi -run TestAXIStatusUnknownNotSuppressed` | Every signal today is gated on `count > 0`; an unknown state reaching that gate disappears, which this row catches. |
| 18 | no `DefaultBranch` symbol survives anywhere in the tree | conformance | `go test ./internal/conformance -run TestRootConformance` with a source-level sweep for the identifier | A partial migration that leaves one caller behind turns this red without needing a behavioral fixture per call site. |
| 18 | `RepoFacts` carries the unresolvable-default state rather than a guess | unit (`internal/git`) | `go test ./internal/git -run TestFactsUnresolvableDefault` | A fixture repo with no `origin/HEAD` and two local branches; an implementation returning `"main"` fails the state assertion. |
| 19 | `bench diff` exits 1 naming the unresolvable default and the `benchBase` escape | contract (`internal/contract/axi`) | `go test ./internal/contract/axi -run TestAXIDiffUnresolvableDefault` | Asserting the literal `benchBase` substring forbids a bare failure that leaves the agent with no next action. |
| 20 | `roadmap --context` renders the git block `unknown` with the snapshot intact | contract (`internal/contract/axi`) | `go test ./internal/contract/axi -run TestAXIRoadmapContextGitUnknown` | Asserting unrelated sections are still present forbids dropping the whole snapshot on an unresolvable default. |
| 14 | the gate goes red on a spec with no coverage map | canary (`tests/canary/`) | `bench gate` against a new fixture whose `EXPECT` names the no-map violation | A check rotted into an always-pass leaves the fixture green, which the canary sweep reports as a failed expectation. |
| 18 | the gate goes red on a tree still calling `DefaultBranch` | canary (`tests/canary/`) | `bench gate` against a fixture reintroducing the identifier, `EXPECT` naming the conformance diagnostic | Proves the conformance sweep in the row above actually bites rather than matching nothing. |

### Edge inventory

Walked per behavior; each landed as a coverage row above or as a **Won't handle**
line here. The classes come from the shell-CLI hostile-input checklist in
`projects/benchkit.md`.

- **Absent vs present-but-empty** — story 1 and story 8's paired rows: distinct
  states at the classifier and distinct postures at the surface.
- **Special files in a discovery path** — story 3: FIFO under a deadline, device,
  and socket, all rejected before the open.
- **Invocation through a symlink** — stories 2 and 5: the dangling case and the
  followed case, asserted separately.
- **Permission-denied** — story 4, with the reason preserved.
- **Hand-edited file with no trailing newline** — the learnings and maps parsers
  already split on `\n` and trim `\r`; the malformed-heading fixture in story 9
  is written without a trailing newline so the row exercises it.
- **Control bytes in rendered text** — the new `error:` lines and `unknown` rows
  carry a state name and a repo-relative path. Paths already route through
  `toon.Representable` at every emitter, and no new sink is introduced, so the
  existing posture covers it.
- **Unquoted multi-word arguments / argument grammar** — no new flags or
  positionals ship; every touched command keeps its existing `usage.Grammar`,
  which the grammar contract test already sweeps.
- **Oversized file** — the existing size bound survives as a distinguished
  outcome (story 1's state table, story 16's content-class row).
- **Re-run idempotency** — every touched command is read-only; none writes.
- **cwd deeper than the repo root** — all touched commands already resolve
  through `git.Root()`, unchanged by this spec.
- **Won't handle: path traversal via a caller-supplied spec path** — `bench
  coverage <path>` accepts an arbitrary path today and keeps that behavior; the
  decision map's state enumeration does not include `traversal`, and adding it is
  a separate cut (flagged above for veto).
- **Won't handle: concurrent mutation of a control record mid-read** — the
  classifier reads once and reports what it saw; a file rewritten between the
  Lstat and the read is a torn read no single-process classifier can close, and
  no consumer treats these reads as transactional.
- **Won't handle: a repository with no commits at all** — `ResolvedDefault`
  already returns `ok=false` there, which is exactly the path stories 18–20
  assert; no separate fixture adds signal.

## Out of scope

- **Outline output bounding and truncation metadata (`RR:C-06`).** A separate
  capability: it governs how much of a *successful* outline is emitted, not
  whether a read failed. Roughly 4 edits, 2 gate runs.
- **Tamper-proofing or signing control records.** A separate capability owned by
  FT71's evidence posture — integrity against a hostile writer, not honesty about
  a failed read. Not estimated here; FT71 carries it.
- **Path-traversal classification for caller-supplied paths.** A separate
  capability: it guards the argument surface rather than the control-record
  reads this spec owns. Roughly 3 edits, 2 gate runs. Flagged above for veto if
  you would rather fold it in.
