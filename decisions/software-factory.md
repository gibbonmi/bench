# Software factory

Status: shaping

## Destination

Bench will develop a traceable, recoverable local factory. It starts with a
structural refactor, a module-deepening refactor, and one CLI improvement. A
shared projection then supports a useful roadmap and execution view. Durable
execution and a qualified Regroup pilot follow through existing release
requirements.

## #1: What direction and authority govern the factory?

Blocked by: none
Type: Grill

### Question

What work does the software factory organize?

### Answer

One factory roadmap organizes Bench and one later Regroup pilot. Bench retains
its gate, release, and project-owned policy. Regroup retains its domain and
design authorities. The factory does not select new product direction or permit
implementation without a target-specific approval.

## #2: What recurring maintenance cadence applies?

Blocked by: none
Type: Grill

### Question

How does recurring maintenance select work?

### Answer

FT243 checks introduced structure violations on each run and makes existing
violations ranked proposals. Module-deepening assessment runs weekly. CLI and
other deep assessments start biweekly and become monthly only after successive
runs find little actionable friction. The policy favors seams and tracer-bullet
slices, not count splitting.

## #3: Which durable-execution principles apply?

Blocked by: none
Type: Grill

### Question

Which principles guide the future execution and view work?

### Answer

Use model-free routine supervision, durable pending decisions, unknown state
with provenance, monitor continuity, typed lifecycle controls, and one shared
projection. Existing authoritative readers remain state owners. These principles
create no Firstmate dependency and make no runtime verification claim.

## #4: What qualifies the external pilot?

Blocked by: none
Type: Grill

### Question

What must precede a Regroup factory pilot?

### Answer

FT306 requires durable execution, FT71, and existing release qualification.
The pilot is one user-visible Regroup change with browser behavior and visual
evidence. It begins only after the local workflow is inspected and a change is
selected.

## #5: Which structural target starts the sequence?

Blocked by: #1
Type: Research

### Question

Which selected structural pass establishes the baseline?

### Answer

The structural-refactor-pass landing at `e86d344d` establishes this baseline.
Its retained retro records the survey, accepted scope, gate results, and cost.
The next target requires a separate module-deepening survey.

## #6: Which module-deepening target should FT302 select?

Blocked by: #5
Type: Research

### Question

Which post-structural survey target will hide behavior behind a useful seam?

### Answer

The reviewer approved the four-candidate deepening batch on 2026-09-05 from the
`/bench-deepen` survey. All four candidates shipped. The named Git admin readers
landed at `5589e73a`, and the landing-destination fact at `a928bebc`. The
ledger transaction with the settle policy landed at `3ed9140f`. The two residuals, the diff
package policy extraction and the Git-reader promotion, stay on the FT302 row.
They need no map ticket.

## #7: Which CLI improvement should FT303 select?

Blocked by: #6
Type: Research

### Question

Which grouped calls, existing verbs, and attributed accessible sources justify
one CLI improvement?

### Answer

The assessment of 2026-09-06 is recorded in
`decisions/assets/ft303-cli-assessment.md`. The Claude Code transcripts and the
Codex rollouts are the accessible sources. The census is absent, because each
release deletes its record. The window holds 6373 calls in 36 sessions, 4642
with a Bench head.

Five candidates carry grouped evidence. A mutation-probe verb ranks first. 102
copy-aside probe sequences ran in 10 sessions at a median of 3 calls each. The
census counts a probe that names the pool path as raw shell. FT168 owns
`bench probe`, and FT98 owns the preserve-and-restore face as a new
`bench worktree` subcommand.

A production-or-test projection on
`bench consumers` and `bench outline` ranks second, with 80 calls in 9 sessions
and one 91.7 KB overflow. A `bench structure --path` filter ranks third, with
49 whole-census reads and three contract deviations. A one-line note on
`bench worktree path` ranks fourth as a guidance repair. An inbox-emptying verb
ranks last on one event.

Two asks are served already: the shell `time` prefix wrapped 99 exec calls, and
a heredoc fed 362. The assessment recommends the probe verb as the one
improvement. The reviewer selects in #10.

## #8: What does the shared view project?

Blocked by: #10
Type: Research

### Question

Which existing readers provide the useful observation view without a second
work-state owner?

### Answer

— (open)

## #9: What execution target fulfills the settled principles?

Blocked by: #8
Type: Research

### Question

Which existing execution seams support durable continuity and typed lifecycle
control?

### Answer

— (open)

## #10: Which improvement does FT303 select?

Blocked by: #7
Type: Grill

### Question

Which one evidence-backed improvement does FT303 deliver, by which route, and
with which limits? The route is the light path or `/bench-write-spec`. The
limits state whether the `bench worktree path` note rides with it, and whether
the two served asks close as out of scope.

### Answer

The reviewer selected the mutation-probe verb on 2026-09-06. One spec delivers
FT168's probe verdict and FT98's preserve-and-restore face together. One
guidance ticket in that spec adds the one-line note to `bench worktree path`.
The route is `/bench-write-spec`, because the verb crosses the declared CLI
seam. The `--time` face and the heredoc route are out of scope, and the `PWD`
ask stays with FT254. Ticket #11 bounds the verb.

## #11: Where does the probe verb live, and what bounds it?

Blocked by: #10
Type: Grill

### Question

FT168 names a root verb `bench probe`, and FT98 names a `bench worktree` face.
Which placement holds, which mutation grammar does the verb take, which exit
posture does each verdict carry, and in which checkouts does it run?

### Answer

The verb is the root verb `bench probe`, and FT98's probe face folds into it.
A worktree probe runs as `bench worktree exec <target> -- bench probe ...`, the
only worktree form. The mutation is an exact string pair or an omission on one
named file. The forms are `--swap <old> --with <new>` and `--omit <old>`, and
the old string matches exactly once.

The verdicts are `bit` at exit 0, `silent` and `invalid` at exit 1, and
`restore-failed` at exit 2 with the preserved copy named. The verb runs in any
checkout and refuses while a gate run holds the tree. The restore is proven
byte-exact before exit. The focused run composes the `bench test` selection
grammar, because FT168 requires the registered check owner. Ticket #12 bounds
the focused-run forms and the verdict record.

## #12: Which focused-run forms and records does the probe verb carry?

Blocked by: #11
Type: Grill

### Question

Does the focused run accept both `bench test` selection forms, the package and
run regex and the `--check <name>`? Does the verb record its verdicts beyond
the terminal answer?

### Answer

The focused run accepts both `bench test` selection forms: a package expression
with `--run <go-regex>`, and `--check <name>` for a registered check. The verb
writes no record beyond the terminal verdict row, and the caller cites that
row. A probe ledger or a census entry is a later row if a retro asks for it.

## Not yet specified

- The targets selected by FT304 and FT305.
- The Regroup pilot change and its browser evidence.

## Spec-writer discretion

- Choose reversible internal placement after the selected target and source
  readers are verified.
- Choose a view layout that consumes the shared projection without becoming a
  state owner.
- Choose the preserved copy's location for the probe verb; it is reversible and
  changes no verdict.

## Out of scope

- A Firstmate runtime dependency, copied implementation, or runtime-capability
  claim.
- A release, deployment, or Regroup mutation before qualification.
- Automatic implementation from `/bench-deepen`.
- A `--time` face on `bench worktree exec` and a heredoc route to the exec
  child: the shell `time` prefix and the exec child's stdin serve both
  (2026-09-06).
- A probe verdict record in the census or a probe ledger: the terminal verdict
  row is the complete evidence (2026-09-06).

## Sources

- Path: `ROADMAP.md`
  Supports: owner rows, dependencies, release qualification, and holds.
  Drift: re-read before a factory target is specified or retired.
- URL: https://github.com/kunchenguid/firstmate/blob/main/docs/architecture.md
  Supports: shared projection and durable-pending-decision design input.
  Drift: use only as inspiration; it establishes no Bench capability.
- URL: https://github.com/kunchenguid/firstmate/blob/main/docs/watcher-continuity.md
  Supports: monitor-continuity design input.
  Drift: use only as inspiration; it establishes no Bench capability.
- URL: https://github.com/kunchenguid/firstmate/blob/main/docs/agent-control.md
  Supports: typed lifecycle-control design input.
  Drift: use only as inspiration; it establishes no Bench capability.
- Path: `decisions/assets/ft303-cli-assessment.md`
  Supports: #7 and the #10 selection; the git-ignored inventory at `research/ft303-cli-assessment/` holds the extraction.
  Drift: re-run the extraction when a probe verb, a production filter, or a structure path filter lands, or at the next CLI assessment.
