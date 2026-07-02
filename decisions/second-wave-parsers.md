# Second-wave check parsers

Wave 2 of the structured state surface (`decisions/state-surface.md`): subcommands
that replace named ad-hoc derivations in the workflow phases. Every wave-2 surface
inherits wave-1's decided contract pattern unchanged — dedicated subcommand, flat
TOON via the shared helper, AXI conformance (structured stdout errors, definitive
empty state, honest exit codes), and a gate conformance layer with canary fixtures.
That inheritance is decided here once; per-command schemas are spec detail.

## #1: What is the admission rule for a wave-2 parser?

Blocked by: —
Type: Grill

### Question
The roadmap parked five candidates (`diff`, `refs`, `coverage`, `doctor`,
`detect`). Do all five ship, or must each earn its place?

### Answer
Consumer-anchored: a parser earns a ticket only by naming the existing phase step
it replaces. `diff` (review-implementation's pin-the-diff step), `refs`
(update-kit's stale-reference sweeps), `coverage` (review-implementation's
coverage axis), and `detect` (setup-repo's stack interview) qualify. `doctor` has
no named consumer yet — #5 decides whether it gets one or stays parked.

## #2: Does `bench diff` own review-base semantics?

Blocked by: —
Type: Grill

### Question
review-implementation resolves its base as the default branch, but a shift
branch's true base is the pre-shift HEAD (roadmap workflow-exits item). Is that
fixed here or there?

### Answer
Folded in here. `bench diff` is the single source of review-base truth: pre-shift
HEAD for shift branches, merge-base with the default branch otherwise.
review-implementation re-points its pin-the-diff step to it. The rest of the
workflow-exits roadmap item (capped-shift routing, superseded-spec retirement)
stays parked.

## #3: How is the pre-shift base recorded, and what is `diff`'s contract?

Blocked by: —
Type: Grill

### Question
For `diff` to resolve a shift branch's base, the shift loop must persist the
pre-shift HEAD somewhere `diff` can read (a ref, a file under `.bench/`, worktree
metadata) — and rollback/cleanup semantics must not orphan it. Also the row
schema: merge-base plus changed-file rows, with what per-file fields.

Codebase facts: `shift_loop` already computes the pre-shift HEAD at branch
creation and discards it when the loop exits. It is not derivable afterward — a
shift stacked on unmerged work makes merge-base land on the feature's fork
point, and reflog recovery dies with expiry or a fresh clone. Recording is
required; the open question is where.

### Answer
Branch-scoped git config: the shift loop writes `branch.<name>.benchBase` at
branch creation. `bench diff` resolves the base as recorded key first,
merge-base with the default branch when the key is absent (non-shift branches,
clones). The key is local-only — review happens where the shift ran — and an
orphaned key after branch deletion is harmless residue the fallback ignores.
Row schema and per-file fields are spec detail.

## #4: What does `refs` sweep, and does the gate consume it?

Blocked by: —
Type: Grill

### Question
Which corpora (kit files, whole repo, docs only) and which reference forms
(`/name`, `$name`, basename) does `bench refs <stem>` cover? And do the gate's
existing stale-reference checks re-point to it — one parser per signal — or stay
independent so the oracle doesn't depend on the surface under review?

### Answer
— (open)

## #5: Who consumes `bench doctor`?

Blocked by: —
Type: Grill

### Question
"Link/wiring drift" names a signal but no phase step that reads it. Under the #1
admission rule it needs a named consumer (update-kit's post-link verification?
session-start?) or it stays parked. Overlap with what the gate already asserts
must also be resolved — doctor must not become a second oracle.

### Answer
— (open)

## #6: What does `detect` discover, into what schema?

Blocked by: —
Type: Grill

### Question
Read-only stack discovery feeding the setup-repo interview: which facts
(language, test runner, linter, CI, package manager), how deep the probing goes,
and how the interview consumes the rows (pre-fill vs. suggestion).

### Answer
— (open)

## #7: What is the coverage-row convention `coverage` parses?

Blocked by: —
Type: Research

### Question
Specs already carry acceptance-coverage maps in prose tables. `bench coverage
<spec>` needs a parseable convention: read the existing maps and the gate's
coverage checks, then formalize row anchors and the status vocabulary the command
reports. Output: a short summary asset proposing the convention.

### Answer
— (open)

## #8: Does the build slice, and in what order?

Blocked by: #3, #4, #5, #6, #7
Type: Grill

### Question
One spec for all admitted parsers (the wave-1 precedent) or staged per parser?
Reviewer's call, made against the finished shapes — sequencing matters because
`diff` also carries the review-base bug fix.

### Answer
— (open)
