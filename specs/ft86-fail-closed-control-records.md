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

5. As an agent reading a *control record* that is a symlink to a regular file, I
   want the link followed and the target classified, so that a linked journal or
   roadmap is read rather than rejected. This rule binds control-record reads
   only; `bench outline` is deliberately not a consumer of it (story 16).
   Line: gpt-5.6-luna / medium. This is the companion assertion to story 2 and
   follows directly from the same Lstat-first rule.

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

16. As an agent running `bench outline`, I want its existing skip-row posture and
    its no-follow rule for tracked symlinks preserved exactly as they are, so
    that a symlink is still indexed once under its target's own path and the
    command still exits 0 on a skip. `bench outline` is therefore **out of** the
    classifier migration: it keeps its local Lstat and regular-file handling.
    Line: gpt-5.6-luna / medium. This is a deliberate non-change, and the guard
    is the contract test that already encodes the decision.

17. As an agent reading `bench status`, I want any signal whose underlying read
    failed to render an explicit `unknown` row naming the state and path, at its
    normal severity and never suppressed by a zero count, while the command still
    exits 0, so that the ambient board degrades visibly instead of reporting a
    clean repository. This covers the signals whose sources this spec already
    touches — the decision-map count, the capture-drain count, and the
    roadmap-reconcile read; `structure.ViolationCount` and the specs housekeeping
    counters are explicitly not migrated. Line: gpt-5.6-luna / medium. The
    failure-row pattern is already established by the worktree row; this applies
    it to the signals whose readers change here.

18. As a maintainer, I want `git.DefaultBranch` deleted and `git.ResolvedDefault`
    to be the sole owner of the default-branch fact, with every caller handling
    `ok=false` explicitly, so that no code path can fabricate `main` for a
    repository that has no resolvable default. Line: gpt-5.6-terra / medium.
    There are **eleven** call sites across `diff`, `status`, `adopt`, `worktree`,
    and `git` itself — enumerated in Implementation decisions — and each needs its
    own posture decision rather than a mechanical substitution.

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

**`unsupported-schema` is shape-based, not marker-based, and it is a parser
state.** No control record grows a schema or version marker — that was rejected
during shaping. The state means the file read cleanly and is valid UTF-8, but its
structure is not one the consuming parser recognizes. The classifier cannot
produce it, because the classifier has no parser: it is declared in the shared
state vocabulary and *returned by each parser*, which owns the predicate for its
own document (`maps`: a decisions file with no `## #<n>:` ticket heading;
`roadmap`: no recognizable rows; `learnings`: no dated heading). Consequently it
is asserted at the surfaces, not in the `bounds` unit — a `bounds` test for it
could only be satisfied by a constant.

**Migrations.** `roadmap`'s private `readSource` and `readDirSource` are deleted
and their callers move to the classifier. `learnings`, `maps`, and `status` drop
their local `os.ReadFile`/`os.ReadDir` handling and consume the classifier. No
migrated surface keeps reading logic of its own; each maps a state to its
posture.

**`bench outline` is out of the migration.** The map's ticket #1 grouped
`outline` with the fail-closed query commands, and this spec deviates. Two facts
force it. First, outline's walk classifies every non-regular entry as
`nonregular` and exits 0, and that is a *recorded* decision: following a tracked
symlink would index the target's symbols under the symlink's path, emitting
`file:line` anchors that do not hold, so the target must be indexed once under
its own path. Story 5's follow-the-link rule is therefore correct for control
records and wrong for outline. Second, outline's failure classes are per-file
within a listing over `git ls-files` output, where a tracked-but-deleted file is
routine in a dirty tree; failing closed there would turn an ordinary working
state into exit 1. Outline's skip rows already surface degradation honestly,
which is the property the spec is buying everywhere else. It keeps its local
Lstat and regular-file handling unchanged.

**Postures, by surface.** Query commands (`learnings`, `maps`, `roadmap`,
`coverage`) exit 1 with a `toon.Errorf`-shaped `error:` line on any non-absent
failure; absent renders the empty state with exit 0. `bench status` exits 0 and
renders explicit `unknown` rows. `roadmap --context` exits 0 and degrades the
affected source or block. `bench outline` is unchanged.

**Status signal counts carry state, narrowly.** `maps.UnresolvedCount` returns a
bare integer that cannot express failure, and status gates its row on `n > 0`, so
a failed scan renders as a clean repository. That signature changes to carry the
count plus its readability state, so the dashboard's `unknown` row and the
listing's rows derive from the same scan. The same applies to the capture-drain
counts and the roadmap-reconcile read, whose sources this spec already touches.
`structure.ViolationCount` and the specs housekeeping counters
(`retirementCount`, `orphanedPickupCount`) are **not** migrated: `internal/
structure` and `internal/spec` are outside the map's module boundary, and pulling
them in would widen the spec past what was decided.

**TOON state vocabulary is reused, not reinvented.** The state names are exactly
`roadmap --context`'s existing `absent` / `empty` / `parsed` / `malformed`
vocabulary extended with `unreadable`, `wrong-type`, and `unsupported-schema`,
and the dashboard's degradation literal stays `unknown`, matching the worktree
list's existing usage.

**`git.DefaultBranch` is deleted outright.** `ResolvedDefault` absorbs the
`origin/HEAD` lookup it currently delegates, and keeps its sole-local-branch
fallback — that fallback is the only thing that makes a `master`-only repository
resolve, so it is a behavior to preserve, not incidental code. There are
**eleven** call sites, enumerated so the migration has no residue:

| call site | posture on `ok=false` |
|---|---|
| `internal/diff/diff.go:178` | fail closed, exit 1, name the `benchBase` escape (story 19) |
| `internal/status/status.go:467` (`appendRetirement`) | tolerate — reads as "not the default branch", skips the advisory signal, matching the posture its existing comment records |
| `internal/status/status.go:498` (`appendRoadmapReconcile`) | tolerate, same as above |
| `internal/git/git.go:180` (`PruneLandedBranches`) | already errors; the message stops naming a fabricated branch |
| `internal/git/git.go:240` (`LandedState`) | already errors; same message fix |
| `internal/git/git.go:288` (`Facts`) | propagates the unresolved state (see below) |
| `internal/git/git.go:334` (`ResolvedDefault` itself) | the delegation disappears when the lookup is absorbed |
| `internal/adopt/link_hook.go:95` | candidate ref, already probe-guarded — skip the candidate |
| `internal/adopt/link_stage.go:97` | candidate ref — skip the candidate |
| `internal/worktree/refresh/refresh.go:48` | candidate ref, already falls through to `ResolvedDefault` then `HEAD` |
| `internal/worktree/lifecycle.go:194` | candidate ref, already falls through to the empty-ref `worktreeAdd` |

`internal/roadmap/context_parse.go:238` is **not** a call site — it reads the
`gf.DefaultBranch` struct field — but it is the rendering consumer of whatever
`Facts` returns, so it moves with story 20.

**`git.Facts` on an unresolvable default (flagged for veto).** The map did not
close this and the spec decides it: `RepoFacts` gains a `DefaultResolved bool`,
`DefaultBranch` holds the empty string when unresolved, and `Ahead`/`Behind`
stay zero rather than being computed against a branch that does not exist —
`Facts` currently derives them from `rev-list <default>...HEAD`, which would
error. `roadmap --context` renders `unknown` in the default-branch, ahead, and
behind cells of its five-column git row, keeping the row's arity. This changes a
public struct and an emitted row shape, so it is the third call in this spec you
should look at directly.

**Coverage validation.** `Check` gains a no-map rule: `State(p) == "no-map"` is a
violation unless the historical marker is present. Story references validate
against the exact set of declared story numbers rather than the maximum; a
reference to a number not in the set, a `0`, and a range whose end is below its
start each produce their own violation message. The existing violation phrasings
are matched downstream by substring, so new messages follow the same
`coverage map row %d …` shape and existing phrasings are not reworded.

**Flagged for veto — three calls the map did not close.**
(a) Story 16 removes `bench outline` from the classifier migration entirely,
against the map's ticket #1, which grouped it with the fail-closed query
commands. The reasoning is above; the short form is that outline's no-follow rule
is a recorded decision this spec would otherwise have overturned by accident.
(b) The `ROADMAP.md` row for FT86 names a `traversal` state the decision map's
enumeration does not include. This spec follows the map's six states and does not
add `traversal`; the caller-supplied-path surface it would guard (`bench
coverage <path>`) keeps its current behavior. If you want traversal classified,
it is a separate cut. (c) `git.Facts`'s unresolved-default shape — the new
`DefaultResolved` field and the zeroed `Ahead`/`Behind` — described above.

**Capability posture for the hostile-filesystem fixtures.** Several red signals
depend on filesystem features a host may lack or a root user may ignore: chmod
0o000 (`capability.Privilege`), FIFOs (`capability.Fifo`), and symlinks
(`capability.Symlink`). Those tests take the capability helper, exactly as
`axi_outline_test.go` already does for symlinks, so a host without the feature
emits an honest `bench-skip` line rather than a false green. Because a skip and a
deleted assertion both look green locally, the build is not done until the
affected suites have been run once with `BENCH_REQUIRE_CAPABILITIES=1` — the knob
both release workflows already wire on — and pass with zero skips on a host that
has all three.

## Testing decisions

A good test here drives a real `bench` subcommand in a throwaway fixture
repository and asserts the exit code and the stdout string — never a package
internal. The one exception is the classifier itself, which is a deep unit whose
hostile-filesystem semantics are cheaper and sharper to assert directly than
through five surfaces. `unsupported-schema` is the counter-case: it is a parser
state, so it is asserted only at the surfaces.

Seams, and their prior art:

- **`internal/bounds` unit** — the classifier's six states against real fixture
  files (dangling symlink, FIFO, chmod 0o000, symlink-to-regular, empty,
  oversized). Prior art: `internal/bounds/bounds_test.go`. The hostile fixtures
  route through `internal/capability` so a host lacking a feature skips honestly
  rather than passing vacuously.
- **`internal/contract/axi`** — every surface posture: exit codes and structured
  `error:` strings for `learnings`, `maps`, `roadmap`, `roadmap --context`,
  `coverage --check`, `diff`, and `status`, plus the two parser states.
  `bench outline` appears here only as a *regression guard* — its existing
  `TestAXIOutlineSymlinkSkipped` must stay green, which is what forbids migrating
  it onto the classifier. Prior art: `axi_fail_closed_test.go`,
  `axi_outline_test.go`, `axi_coverage_test.go`, `axi_roadmap_context_test.go`.
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
    fixture FS ─▶  [ diff      | status                   ]  ──▶  exit code
                   [ outline — unchanged, regression only ]
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
| 2 | dangling symlink classifies `unreadable`, not `absent` | unit (`internal/bounds`) | `go test ./internal/bounds -run TestClassifyDanglingSymlink` (`capability.Symlink`) | A `ReadFile`-first implementation sees ENOENT and answers `absent`, which this assertion forbids by name. |
| 3 | FIFO classifies `wrong-type` without opening | unit (`internal/bounds`) | `go test ./internal/bounds -run TestClassifyFIFOWithoutOpen` (`capability.Fifo`) | The test drives a FIFO with no writer under a deadline; an implementation that opens before type-checking blocks and fails the deadline rather than returning. |
| 3 | device and socket classify `wrong-type` | unit (`internal/bounds`) | `go test ./internal/bounds -run TestClassifyNonRegular` (`capability.Fifo`) | `/dev/null` is readable and would classify `empty` under a type-blind implementation. |
| 4 | mode 0o000 classifies `unreadable` with the reason preserved | unit (`internal/bounds`) | `go test ./internal/bounds -run TestClassifyPermissionDenied` (`capability.Privilege`) | Asserting a non-empty reason alongside the state forbids the implementation that discards the underlying error. |
| 5 | symlink to a regular file is followed and classified | unit (`internal/bounds`) | `go test ./internal/bounds -run TestClassifySymlinkFollowed` (`capability.Symlink`) | A blanket "any symlink is unreadable" implementation passes story 2 and fails here. |
| 6 | directory classification returns absent, empty, and unreadable | unit (`internal/bounds`) | `go test ./internal/bounds -run TestClassifyDir` (`capability.Privilege`) | An unreadable directory and an empty one both yield zero entries; only the state distinguishes them. |
| 7 | `bench maps` reports `unsupported-schema` for a decisions file with no ticket heading, and `malformed` for invalid UTF-8, as two distinct states | contract (`internal/contract/axi`) | `go test ./internal/contract/axi -run TestAXIMapsUnsupportedSchema` | Two fixtures asserted in one run; collapsing either state onto the other, or onto a generic error, fails one of them. Asserted at the surface because the classifier has no parser and a `bounds` test could only check a constant. |
| 7 | `bench roadmap` reports `unsupported-schema` for a document with no recognizable rows | contract (`internal/contract/axi`) | `go test ./internal/contract/axi -run TestAXIRoadmapUnsupportedSchema` | A second parser asserting the same state forbids implementing it for exactly one document and leaving the vocabulary otherwise dead. |
| 8 | `bench learnings` exits 1 with a structured error on an unreadable journal | contract (`internal/contract/axi`) | `go test ./internal/contract/axi -run TestAXILearningsFailClosed` (`capability.Privilege`) | A fixture with a mode-0o000 journal; the current implementation exits 0 with an empty table, which the exit assertion rejects. |
| 8 | `bench learnings` with no journal exits 0 and renders the empty state | contract (`internal/contract/axi`) | `go test ./internal/contract/axi -run TestAXILearningsAbsentIsEmpty` | Paired with the row above, this forbids the over-correction that fails closed on absence too. |
| 9 | malformed learning headings render as rows with line and reason, exit 1 | contract (`internal/contract/axi`) | `go test ./internal/contract/axi -run TestAXILearningsMalformedRows` | A fixture mixing one good and one malformed heading; asserting the good row is still present forbids a hard-fail that drops everything. |
| 10 | an unreadable decision file gets a row naming it, exit 1 | contract (`internal/contract/axi`) | `go test ./internal/contract/axi -run TestAXIMapsUnreadableFileRow` (`capability.Privilege`) | The current implementation `continue`s past the read error; asserting the filename appears in stdout catches the silent drop. |
| 10 | `bench maps` with no `decisions/` directory exits 0 and renders the empty state | contract (`internal/contract/axi`) | `go test ./internal/contract/axi -run TestAXIMapsAbsentIsEmpty` | Without this pair, an unconditional exit-1 stub satisfies the row above while destroying the absent-is-authoritative-empty contract the spec exists to establish. |
| 11 | the status map count and the `bench maps` listing agree on what was readable | contract (`internal/contract/axi`) | `go test ./internal/contract/axi -run TestAXIMapsCountMatchesListing` (`capability.Privilege`) | One fixture with an unreadable map file, asserted through both surfaces; a count that fabricates zero disagrees with a listing that shows the row. |
| 12 | `bench roadmap` exits 1 with a structured error on an unreadable ROADMAP.md | contract (`internal/contract/axi`) | `go test ./internal/contract/axi -run TestAXIRoadmapFailClosed` (`capability.Privilege`) | A mode-0o000 fixture; an empty-document print exits 0 and fails the assertion. |
| 12 | `bench roadmap` with no `ROADMAP.md` exits 0 and renders the empty state | contract (`internal/contract/axi`) | `go test ./internal/contract/axi -run TestAXIRoadmapAbsentIsEmpty` | The pair that forbids the always-exit-1 stub, as for `maps` above. |
| 13 | `roadmap --context` keeps the snapshot and marks one failed source | contract (`internal/contract/axi`) | `go test ./internal/contract/axi -run TestAXIRoadmapContextDegrades` | Asserting both exit 0 *and* the presence of unrelated sections forbids the whole-snapshot error return the code takes today. |
| 14 | a spec with no map and no historical marker fails `--check` | contract (`internal/contract/axi`) | `go test ./internal/contract/axi -run TestAXICoverageNoMapFails` | The current implementation returns nil violations for `no-map` and exits 0, which the exit-1 assertion rejects. |
| 14 | a spec carrying the historical marker still passes `--check` | contract (`internal/contract/axi`) | `go test ./internal/contract/axi -run TestAXICoverageHistoricalPasses` | Paired with the row above, this forbids the implementation that fails every unmapped spec including opted-out ones. |
| 15 | story `0`, a non-member story number, and a reversed range each fail | contract (`internal/contract/axi`) | `go test ./internal/contract/axi -run TestAXICoverageStoryMembership` | Three fixtures in one table, each currently passing; a max-only validator fails all three assertions. |
| 15 | a valid comma list and forward range still pass | contract (`internal/contract/axi`) | `go test ./internal/contract/axi -run TestAXICoverageValidStoryRefs` | Forbids the over-strict validator that rejects legitimate multi-story rows. |
| 16 | a tracked symlink is skipped, not followed, and outline still exits 0 | contract (`internal/contract/axi`) | already covered — `go test ./internal/contract/axi -run TestAXIOutlineSymlinkSkipped` is green today and must stay green | This is the regression guard for the deliberate non-change: if the build migrates outline onto the classifier, story 5's follow-the-link rule indexes the target under the link's path and this test fails by name. |
| 16 | a non-regular tracked entry keeps its `nonregular` skip row and exit 0 | contract (`internal/contract/axi`) | already covered — the same test asserts the `link.go,nonregular` row and exit 0 | An implementation that promotes non-regular to a failure class turns exit 0 into exit 1 here. |
| 17 | a signal whose read failed renders an `unknown` row and status still exits 0 | contract (`internal/contract/axi`) | `go test ./internal/contract/axi -run TestAXIStatusUnknownRow` | Asserting the row is present *and* the exit is 0 forbids both the fabricated zero (row absent) and the fail-closed over-correction. |
| 17 | the `unknown` row is not suppressed by a zero count | contract (`internal/contract/axi`) | `go test ./internal/contract/axi -run TestAXIStatusUnknownNotSuppressed` | Every signal today is gated on `count > 0`; an unknown state reaching that gate disappears, which this row catches. |
| 18 | no `DefaultBranch` symbol survives anywhere in the tree | conformance | `go test ./internal/conformance -run TestRootConformance` with a source-level sweep for the identifier | A partial migration that leaves one caller behind turns this red without needing a behavioral fixture per call site. |
| 18 | `RepoFacts` carries the unresolvable-default state rather than a guess | unit (`internal/git`) | `go test ./internal/git -run TestFactsUnresolvableDefault` | A fixture repo with no `origin/HEAD` and two local branches; an implementation returning `"main"` fails the state assertion. |
| 18 | a sole-`master` repository still resolves, via the single-local-branch fallback | unit (`internal/git`) | `go test ./internal/git -run TestResolvedDefaultSoleMaster` | Handoff item 6 names master-only as an owned hostile input. `ResolvedDefault` is rewritten to absorb the `origin/HEAD` lookup, and without this row a rewrite that drops the sole-branch fallback turns nothing red while breaking every `master` repository. |
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
  followed case, asserted separately, and story 16's guard that outline's
  opposite no-follow rule survives the classifier's arrival.
- **Permission-denied** — story 4, with the reason preserved.
- **Hand-edited file with no trailing newline** — Handoff item 6 assigns this to
  *each* parser. The learnings fixture in story 9 and the decisions fixture in
  story 7's `unsupported-schema` row are both written without a trailing newline,
  so the learnings and maps parsers each carry it; `roadmap` inherits it through
  story 7's second row.
- **Control bytes in rendered text** — the new `error:` lines and `unknown` rows
  carry a state name and a repo-relative path. Paths already route through
  `toon.Representable` at every emitter, and no new sink is introduced, so the
  existing posture covers it.
- **Unquoted multi-word arguments / argument grammar** — no new flags or
  positionals ship; every touched command keeps its existing `usage.Grammar`,
  which the grammar contract test already sweeps.
- **Oversized file** — the existing size bound survives as a distinguished
  outcome in story 1's state table.
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
- **Won't handle: a tracked file deleted from the worktree but not staged** —
  `bench outline` Lstats paths that come from `git ls-files`, so this is a
  routine dirty-tree state that would become a hard failure if outline were
  migrated. Story 16 keeps outline out of the migration, which is what removes
  the hazard; no new row is needed.
- **Won't handle: a host that lacks chmod, FIFO, or symlink semantics** — those
  fixtures take `internal/capability` and skip honestly. The posture, and the
  `BENCH_REQUIRE_CAPABILITIES=1` run that proves the skips are not hiding an
  unimplemented state, are stated in Implementation decisions.

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
