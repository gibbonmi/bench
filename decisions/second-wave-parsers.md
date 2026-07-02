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
Evidence-anchored: a parser earns a ticket by naming the phase step it replaces
**and** showing evidence — an observed wrong answer, or recurrence through the
learnings funnel (`decisions/state-surface.md` #2). Consumer-naming alone is not
enough; it admits pattern-completion. Under this bar the wave is two parsers:
`diff` (live review-base bug in review-implementation) and `coverage` (reviewer-
admitted as part of the same review story — the coverage axis needs attribution).
`refs`, `doctor`, and `detect` returned to the roadmap until the funnel names
them; `refs` additionally duplicates a signal the gate already enforces.

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

## #4–#6: retired

`refs`, `doctor`, and `detect` failed the #1 evidence bar and returned to the
roadmap. Their tickets are removed; numbering stays stable.

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

Blocked by: #7
Type: Grill

### Question
One spec for both admitted parsers (`diff`, `coverage`) or staged? Reviewer's
call, made against the finished shapes — sequencing matters because `diff`
carries the review-base bug fix.

### Answer
— (open)
