# Compiled base-relative git-context (FT9)

Grow `bench diff` with a `--full` flag that bundles the diff body and the
base-relative log alongside the changed-file list it already emits, replacing
the two-call git prose in `/bench-review-implementation` with one command. The
original rationale — partitioning self-vs-other-writer git state — died when
invariant 1 made one-writer-per-tree the rule; what remains is call-count
reduction at the one genuinely repeated base-relative call site.

## #1: Does FT9 survive at all, and as what?

Type: Grill

### Answer
It survives **narrowed to the base-relative bundle**. The generic
`status + diff + log + staged` framing dies: recon found no repeated call site
for it — the git-context agents actually repeat is base-relative (review diffs)
or specialized (retirement greps). The only genuinely repeated multi-call
base-relative pattern is the review phase's `git diff <base>...HEAD` +
`git log <base>..HEAD`, so FT9 is scoped to exactly that. Rejected: build the
generic bundle (no call site); drop FT9 to FT6's parked tier (the review pair is
a real repeated pattern worth collapsing).

## #2: New command or grow `bench diff`?

Type: Grill

### Answer
**Grow `bench diff` behind a flag.** `bench diff` already owns base resolution
(recorded key or merge-base) and the changed-file list; a new command would
duplicate that base machinery — a one-source-per-fact violation. The flag adds
the log and diff-body sections; bare `bench diff` keeps its current output so
existing callers are untouched. Rejected: a new command (duplicate base
resolution).

## #3: What does the flag cover — base-relative only, or working-tree too?

Type: Grill

### Answer
**Purely base-relative.** The flag adds the diff body (`base...HEAD`, matching
the existing file list) and `git log <base>..HEAD`, and nothing else. Working-
tree `git status` (uncommitted state) stays out: the anchor use is the review
phase, which reads committed branch content, and folding in working-tree state
reintroduces the generic bundle killed in #1 and mixes two unrelated contexts on
one flag. Rejected: fold in `git status` (a mid-shift working-tree view is a
different context; park it if it ever earns evidence).

## #4: How do the log and diff body reconcile with the TOON-stdout gate?

Type: Grill

### Answer
Split by nature:
- **log → TOON table** `log[N]{sha,subject}` — naturally tabular, cheap, fully
  conformant.
- **diff body → a single delimited raw block, appended last, declared a
  documented output contract** (the craft-cli exemption for a surface with its
  own contract). TOON-escaping a unified diff serves nobody and would mangle
  `@@`/`+`/`-` framing; the agent reads it verbatim.
- **Bare `bench diff` stays 100% TOON** — only the `--full` path carries the one
  documented raw exception.

Rejected: TOON-encode the diff body (mangles the diff, no reader benefit).

## #5: Does the payoff land in a phase doc?

Type: Grill

### Answer
**Yes — the spec includes replacing the review-phase prose with
`bench diff --full`.** `/bench-review-implementation` (the `git diff <base>...HEAD`
+ `git log <base>..HEAD` instruction) is the single base-relative call site and
the only edit; without it FT9 is unused surface. `bench-debug`'s
`--grep=spec-retire` / `--diff-filter=D` are specialized lookups, not the
base-relative bundle, and stay untouched.

## Handoff

1. **Module boundaries.** `internal/diff/diff.go` `Command` grows a `--full`
   branch that appends two sections (log, diff body). The phase-doc edit lands
   in `.claude/commands/bench-review-implementation.md` and its
   `.agents/commands/` mirror. No new module.
2. **Contracts.** `bench diff --full`: input the `--full` flag; output the
   existing `branch/base/method` preamble + `files[N]{status,path}` table, then a
   `log[N]{sha,subject}` TOON table (`git log <base>..HEAD`), then a delimited raw
   diff-body block (`git diff <base>...HEAD`). Exit codes unchanged: 0 ok, 1
   not-in-repo / base-unresolvable, 2 usage (unknown arg). Bare `bench diff`
   output is byte-unchanged.
3. **Deep vs thin.** `diff.Command` is the deep unit — it already hides base
   resolution and NUL-safe name-status parsing. The `--full` branch is a thin
   addition (two `git.Raw` calls + rendering); the seam attaches at
   `diff.Command`'s stdout, as it does today.
4. **Black-box assertables.** Exit code; bare `bench diff` stdout carries no
   `log[`/diff sections (regression assert); `--full` stdout carries the
   `log[` TOON header when commits exist and raw `@@` hunk markers for a changed
   file; empty-since-base yields an empty log table and empty diff block.
5. **Gate attachment.** AXI surface/contract tests observe `bench diff` stdout.
   The conformance assertion that stdout is TOON-shaped must be scoped so it
   still bites bare `bench diff` but exempts the documented raw diff-body block on
   `--full` (assert the structured prefix, not the raw tail). This is the one
   seam the gate can't observe generically — see uncertainty flag.
6. **Hostile-input owners.** Paths with spaces / non-ASCII / quotes: the existing
   `-z` NUL framing owns the files table; log subjects and the diff body are raw
   passthrough (no TOON layer, so no escaping hazard). Empty diff / no commits
   since base → empty sections, not error. Base unresolvable → existing error
   path (exit 1). Detached HEAD → existing `(detached)` label.
7. **Uncertainty flags.** How the conformance TOON-stdout assertion exempts the
   `--full` raw block (item 5) is the seam the grill could not fully settle —
   the spec-writer should consult `craft-gate` on whether to assert the
   structured prefix only, gate `--full` under a separate documented-contract
   check, or split the test. Escalate per the `craft-line` ladder rather than
   guessing.
8. **Rejected alternatives.** New command (duplicate base machinery); the generic
   `status+diff+log+staged` bundle (no repeated call site); working-tree `git
   status` in the flag (different context, reintroduces the killed bundle);
   TOON-encoding the diff body (mangles the diff).
9. **Domain watch-outs.** The log uses two-dot `base..HEAD` while the diff and
   file list use three-dot `base...HEAD`; this asymmetry is intentional and
   matches the existing review-phase prose and git convention (three-dot diff =
   changes on the HEAD side since merge-base; two-dot log = commits on HEAD since
   base). Do not "normalize" the two to match. Unified-diff bodies can be large;
   no truncation is specified — the flag is opt-in and the reader wants it whole.

Dependency order: n/a — single spec.
