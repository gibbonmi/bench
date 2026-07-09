# Docs batch (FT52)

Status: implemented

> **Source and authority.** This spec is not compiled from a `decisions/<topic>.md`
> map. It is a reviewer-directed batch drain: the six scope items are pre-decided in
> `ASSESSMENT.md` §6 (Docs, skills, vocabulary) plus ranked backlog item 11, pushed
> into `specs/` on the reviewer's instruction. Per the `/bench-write-spec` batch-drain
> override, every defaulted decision is transcribed as-decided and flagged for
> **post-hoc veto** below — none is reopened here.

## Problem

Six low-severity doc and packaging drifts, each individually below the pipeline
threshold, cluster into one coherent batch:

1. `CONTEXT.md`'s canonical **signal** definition names 6 of the 10 signals
   `bench status` actually emits, and no gate check protects the list — so the file
   that pins the ubiquitous language teaches half the board's vocabulary, and it has
   drifted silently before.
2. "Dashboard" names two things: `CONTEXT.md` defines *ambient dashboard* =
   `bench status`, but the `bench dashboard` HTML artifact has no canonical term —
   and FT38 (the artifact's visual-identity pass) is next in the roadmap and needs
   the word.
3. `README.md`'s curated `internal/` layout omits `dashboard/` and `outline/`, the
   packages behind two commands its own CLI list names.
4. The lighter-path threshold is worded three ways across `.bench/BENCH.md` and the
   two command files — a knowledge duplication that has stayed open across three
   assessments.
5. `projects/benchkit.md` — this repo's internal dogfood profile, naming its Go
   seams — ships to every consumer via the `projects/` glob, alongside the two
   intended example templates.
6. `research/unit_testing.pdf` is referenced nowhere and isn't shipped — unowned
   scratch.

The through-line (assessment theme 2): **records rot unless a gate check pins
them.** The one item with real drift history — the signal vocabulary — gets a
conformance cross-check so it can't drift again; the rest are one-time corrections.

## Solution

- Correct `CONTEXT.md`'s signal enumeration to the full live set, and add a
  conformance cross-check that ties that list to `internal/status/status.go`'s
  emitted signal names — prior art `checkLineBinding`, which cross-checks the
  profile's Lines prose against `.bench/lines.env`.
- Define **"the dashboard page"** as the canonical term for the `bench dashboard`
  HTML artifact, alongside the existing *ambient dashboard* = `bench status`.
- Add `dashboard/` and `outline/` to `README.md`'s `internal/` layout, sweeping for
  any other package behind a named CLI command that the curated list omits.
- Keep the lighter-path threshold's one canonical wording in `.bench/BENCH.md`'s
  "Right-size the process" paragraph; have `/bench-write-spec` and
  `/bench-implement-spec` **reference** it rather than carry their own wording.
- Stop shipping `projects/benchkit.md` in the npm tarball while keeping the two
  example templates (`gl-axi`, `regroup`); update the package-shape conformance
  expectation to match.
- Delete `research/unit_testing.pdf` (git history preserves it).

## Defaulted decisions flagged for post-hoc veto

| # | Decision (as handed down) | Why it is a veto surface |
|---|---|---|
| D1 | Canonical term for the `bench dashboard` artifact is **"the dashboard page."** | Coins vocabulary that feeds the FT38 grill; the reviewer may prefer another term. |
| D2 | `projects/benchkit.md` **stops shipping**; `gl-axi` and `regroup` keep shipping. | Deliberately reverses part of a recent "projects/ ships" fix — the *examples* ship, the internal dogfood profile does not. |
| D3 | `research/unit_testing.pdf` is **deleted**, not relocated. | Irreversible-ish (recoverable only from git history); the reviewer may want it kept in-tree. |
| D4 | Batch line cap: **no story routes `claude-fable-5`.** Prose stories the leverage override would send top route `claude-opus-4-8 / high`, with the top-tier orchestrator reviewing. | A routing cap the reviewer set for this batch; noted so it is visible, not silent. |

## User stories

1. **As a cold session, I want `CONTEXT.md`'s signal definition to name every
   signal `bench status` can print, so that I learn the whole board's vocabulary
   from the file that claims to pin it.** The corrected enumeration is the ten live
   signal names: `gate`, `git`, `worktree`, `guards`, `drain`, `structure`,
   `decisions`, `specs`, `reviews`, `roadmap`.
   Line: claude-opus-4-8 / high. This is ubiquitous-language guidance prose that the
   leverage override would route top, capped to the mid tier for this batch with the
   top-tier orchestrator reviewing the wording.

2. **As a maintainer, I want a conformance check that fails when the `CONTEXT.md`
   signal list drifts from `status.go`, so that the vocabulary can't silently rot
   again.** The check reads the graded tree's `status.go`, extracts the emitted
   signal-name set, and asserts each name appears in `CONTEXT.md`'s signal
   enumeration; it lands in `TestRootConformance` beside the existing docs-drift
   checks.
   Line: claude-opus-4-8 / medium. Gate/conformance logic is the profile's cached
   mid row, and the source-extraction seam carries a real false-pass hole (common
   words like "git" and "specs" appear all over `CONTEXT.md`) that wants care.

3. **As a session reading or writing about the `bench dashboard` artifact, I want a
   canonical name for it, so that "dashboard" stops ambiguously meaning both the
   `bench status` board and the HTML page.** Add **"the dashboard page"** to
   `CONTEXT.md`'s core terms, distinguished from *ambient dashboard* = `bench
   status`.
   Line: claude-opus-4-8 / high. Vocabulary guidance prose under the leverage
   override, capped to mid for this batch with the orchestrator reviewing.

4. **As a reader of `README.md`, I want the `internal/` layout to name the packages
   behind the commands the same README's CLI list advertises, so that the map
   matches the territory.** Add `dashboard/` and `outline/`; while there, verify the
   curated list against the actual `internal/` directory and add any other package
   that sits behind a named CLI command and is currently omitted.
   Line: claude-opus-4-8 / high. Documentation authoring routes top under the
   leverage override, capped to mid for this batch with the orchestrator reviewing;
   the edit itself is mechanical.

5. **As a reader of the three lighter-path wordings, I want one canonical statement
   and two pointers, so that the threshold can't drift three ways.** Keep
   `.bench/BENCH.md`'s "Right-size the process" paragraph as the single source;
   change `/bench-write-spec`'s "more than a trivial change" and
   `/bench-implement-spec`'s "the seam is obvious" to reference that paragraph
   instead of restating a threshold. The pointers must not re-state the threshold,
   or the duplication reopens.
   Line: claude-opus-4-8 / high. Kit guidance prose under the leverage override,
   capped to mid with the orchestrator reviewing.

6. **As a consumer of the npm package, I want the tarball to carry the two example
   profiles but not this repo's internal dogfood profile, so that I get templates,
   not the kit's own seams.** Exclude `projects/benchkit.md` from `package.json`
   `files[]` while keeping `projects/gl-axi.md` and `projects/regroup.md`; update the
   package-shape conformance expectation (`packagesurface` assets) so the exclusion
   is enforced.
   Line: claude-opus-4-8 / medium. Package-shape conformance is the profile's cached
   mid row, and the files[]-vs-`checkPackageFiles` interaction is a genuine seam
   decision (see Implementation decisions).

7. **As a maintainer, I want the unreferenced scratch PDF gone, so that the tree
   carries no unowned artifacts.** Delete `research/unit_testing.pdf`.
   Line: claude-sonnet-5 / low. Pure mechanical deletion of an unreferenced,
   unshipped file with no seam.

## Implementation decisions

**Modules touched.**
- `CONTEXT.md` — signal enumeration (S1), the dashboard-page term (S3).
- `internal/conformance/` — a new docs-drift cross-check for the signal list (S2,
  new check in the `TestRootConformance` aggregation); the package-shape expectation
  (S6, via `internal/packagesurface/assets.go`).
- `tests/canary/` — a fixture that proves the S2 check bites.
- `README.md` — the `internal/` layout list (S4).
- `.bench/BENCH.md` (canonical wording, unchanged text), `.agents/commands/bench-write-spec.md`,
  `.agents/commands/bench-implement-spec.md` (pointers) — S5. Note `.claude/`
  mirrors are symlink adapters, not a second source.
- `package.json` `files[]` (S6); `research/unit_testing.pdf` deleted (S7).

**S2 cross-check design (prior art `checkLineBinding`).** `checkLineBinding` reads
the graded tree's machine source (`.bench/lines.env`) and asserts the prose file
(`projects/benchkit.md`) names each bound value. The signal check mirrors this:
read the graded tree's `internal/status/status.go` as text, extract the distinct
signal-name literals (the `row{<sev>, "<name>", …}` first string), and assert each
appears in `CONTEXT.md`'s signal enumeration.
- **Direction.** Forward only — every emitted signal name must be present in
  `CONTEXT.md` — matching `checkLineBinding`'s one-directional prose check. This
  closes the drift the assessment actually observed (signals added to `status.go`,
  `CONTEXT.md` not updated).
- **Scoped match, not whole-file `Contains`.** `checkLineBinding` gets away with
  `strings.Contains` because model ids are distinctive tokens. Signal names are
  common words (`git`, `specs`, `reviews`) that appear throughout `CONTEXT.md`, so a
  whole-file substring test false-passes. The check must match each name as a
  delimited token *within the signal-definition enumeration* (the parenthesized list
  on the `signal` term), not anywhere in the file.
- **Absent source.** If `status.go` is unreadable in the graded tree, the check
  no-ops (as `checkLineBinding` skips an empty profile); safe because the
  compiled-core build/test conformance check already fails a tree missing
  `status.go`.

**S6 package-shape design.** `RequiredPackAssets`/`ForbiddenPackAssets`
(`internal/packagesurface/assets.go`) is the one list both the contract surface test
and the conformance `npm pack --dry-run` check iterate.
- Add `projects/benchkit.md` to `ForbiddenPackAssets`; add `projects/regroup.md` to
  `RequiredPackAssets` (pinning the second example, since only `gl-axi.md` is pinned
  today).
- **`files[]` mechanism.** Prefer replacing the `"projects/"` glob with the two
  explicit paths `"projects/gl-axi.md"` and `"projects/regroup.md"` over a
  `"!projects/benchkit.md"` negation. `checkPackageFiles` validates each `files[]`
  entry with `exists()`; a literal `!`-prefixed entry would fail that as a missing
  path. Two explicit includes are the clean "or equivalent" the decision allows and
  need no change to `checkPackageFiles`.

## Testing decisions

- **A good test here** exercises the conformance oracle at its seam: feed it a tree
  and read the diagnostics, never a reading of the diff. The prose-only stories
  (S1 partially, S3, S4, S5, S7) have no gate-observable behavior and are classified
  honestly below, not faked into coverage.
- **Seams tested.** Both live in `internal/conformance`'s `TestRootConformance`
  aggregation: the new signal-list cross-check (S2), and the existing package-shape
  check (S6, via the `packagesurface` list). The gate command is the project gate,
  `.bench/gate.sh`, whose conformance phase runs `TestRootConformance` against the
  tree under grade.
- **Prior art.** `checkLineBinding` (`line_routing_static_test.go`) for prose-vs-source
  cross-checking; `checkNpmPackAssets` (`package_core_checks_test.go`) for the
  Required/Forbidden pack-asset shape; `tests/canary/` for the bite proof.

### Seam diagram

Signal-list cross-check (S2):

    trigger: gate → conformance phase → TestRootConformance
        │
        ▼
    status.go text ─┐
                    ├─▶ [ checkSignalVocabulary ] ──▶ diags: "CONTEXT.md
    CONTEXT.md text ┘        (extract names,             signal enumeration
                             scoped match)               missing signal '<name>'"
                      ◀ tests attach here: TestRootConformance over a fixture
                        tree; canary asserts the red substring

Package-shape check (S6):

    trigger: gate → conformance phase → TestRootConformance
        │
        ▼
    package.json files[] ─┐
                          ├─▶ [ checkPackageFiles →     ──▶ diags: "npm package
    npm pack --dry-run    │      checkNpmPackAssets ]        includes local-only
    (Required/Forbidden)  ┘      over packagesurface        file projects/benchkit.md"
                      ◀ tests attach here: add benchkit.md to ForbiddenPackAssets,
                        run conformance; green only after files[] excludes it

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1, 2 | `CONTEXT.md` names all 10 live signals; conformance fails when the set drifts | `internal/conformance` `TestRootConformance` (new signal cross-check) | `go test ./internal/conformance -run '^TestRootConformance$'` red on the current tree: `status.go` emits `gate,git,worktree,guards,drain,structure,decisions,specs,reviews,roadmap`, `CONTEXT.md` names only 6 — diagnostics name each missing signal | a signal added to `status.go` without a `CONTEXT.md` edit is named by the forward check; the degenerate whole-file `Contains` impl false-passes on common words, so the row demands scoped extraction |
| 2 | the cross-check bites (not an always-pass) | `tests/canary/` fixture: `CONTEXT.md` with one signal name removed | the canary meta-run reports the fixture red with the targeted `missing signal '<name>'` substring | proves the check still fails after S1 corrects the real list — a check rotted to always-pass fails the canary |
| 6 | `projects/benchkit.md` excluded from the tarball; `gl-axi` + `regroup` still ship | `internal/conformance` package-shape (`npm pack --dry-run`, `packagesurface`) | adding `projects/benchkit.md` to `ForbiddenPackAssets` turns conformance red (`npm package includes local-only file projects/benchkit.md`) because the `projects/` glob still ships it; green after `files[]` lists the two examples explicitly | catches failure-to-exclude (benchkit Forbidden) directly |
| 6 | the `regroup` example keeps shipping | same | already covered — `projects/regroup.md` added to `RequiredPackAssets` is green now, and goes red only if a later `files[]` edit over-excludes it | guards the exclusion from amputating the second example |
| 3 | "the dashboard page" is canonical for the `bench dashboard` artifact | `CONTEXT.md` prose | not TDD-able — vocabulary prose with no gate surface; decision D1 adds no conformance guard | — |
| 4 | `README.md` `internal/` layout names `dashboard/` and `outline/` (and any other omitted package behind a named CLI command) | `README.md` prose | not TDD-able — no conformance check ties the curated `internal/` list to the directory; a one-time doc correction | — |
| 5 | lighter-path threshold single-sourced in `.bench/BENCH.md`; the two commands reference it | kit prose | not TDD-able — the "reference, don't restate" property is semantic; the shared-rule single-sourcing check covers only the four invariants and communication rules, not this threshold | — |
| 7 | `research/unit_testing.pdf` removed | filesystem | not TDD-able — an unreferenced, unshipped orphan with no calling behavior; git preserves history | — |

### Edge inventory

Walked per behavior; each lands as a coverage row above or a **Won't handle** line
here.

- **Package-shape drift (over/under exclusion)** — covered: the benchkit-Forbidden
  row catches failure-to-exclude, the regroup-Required row catches over-exclusion.
- **Stale / dangling references (docs-drift class)** — covered by the existing
  stale-command-reference conformance check plus the gate: after S5 edits the two
  command files, a dangling `/bench-*` phase or path reference goes red there.
- **Token false-match (malformed input for S2)** — covered: folded into the S2
  degenerate check — a whole-file `Contains` impl false-passes on `git`/`specs`, so
  the check must match names within the signal enumeration specifically.
- **Bidirectional drift** — **Won't handle**: `CONTEXT.md` naming a signal
  `status.go` does not emit (a removed signal) is a lower-frequency, lower-risk drift
  than the additions the assessment observed; the forward direction matches the
  `checkLineBinding` precedent and closes the real failure.
- **Absent / empty source for S2** — **Won't handle**: an unreadable `status.go` in
  the graded tree makes the check no-op (mirroring `checkLineBinding` on an empty
  profile) because the compiled-core build/test check already fails such a tree.
- **Re-run idempotency** — **Won't handle**: PDF deletion is idempotent and
  conformance is a pure read; there is no state to corrupt on re-run.
- **Shell-CLI hostile-input classes** (paths with spaces/glob chars, control bytes
  in git text, missing-trailing-newline files, absent-vs-empty file, unquoted
  multi-word args, required tool missing from PATH, symlink invocation, SIGINT
  mid-loop, cwd deeper than root) — **Won't handle**: this batch adds no shell CLI
  surface; it edits docs and two Go conformance lists, so the profile's shell-CLI
  checklist classes do not apply.

## Out of scope

- **A conformance guard tying `README.md`'s `internal/` list to the directory** —
  a separate protection capability (its own drift check), and decision D-scope
  protects only the signal vocabulary this pass. ~3 edits, 2 gate runs.
- **A wording-drift conformance check for the lighter-path threshold** — a separate
  capability; decision 4 is single-source-only, not gate-enforced. ~2 edits, 2 gate
  runs.
- **A conformance guard for the "dashboard page" term** — separate capability; the
  term is being coined, not yet worth pinning. ~2 edits, 2 gate runs.
- **`README.md` build-command / clone-path refresh** — named in backlog item 11
  ("+ build command") but not in the orchestrator's six decisions; a separate doc
  pass. ~1 edit, 1 gate run.
- **FT54 — assessment owner (the assessment-drill phase command)** — its own
  roadmap row; the recurring §6 finding, not part of this batch.
- **FT38 — dashboard visual-identity pass** — tabled roadmap row; this batch only
  supplies the term FT38's grill will use.
