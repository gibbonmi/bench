# commit-wrappers — `bench spec implemented` + `bench commit`

Source: ROADMAP.md FT3. Turns the commit discipline currently carried as prose,
duplicated across `/bench-implement-spec` (close-on-green) and `/bench-final-check`,
into two CLI wrappers over existing logic, then collapses that prose to pointers.
The commit-behavior wrappers were parked out of `implement-spec-lean` (#1 there)
as "changing what a build commit does — needs its own shaping"; this map is that
shaping.

## #1: Does `bench commit` fold in the status-flip, or stay decoupled?

Type: Grill

### Question
The roadmap row says "commit could fold in the spec status flip." One command
that does both, or two composable pieces?

### Answer
**Decoupled but composable.** `bench spec implemented <slug>` is a standalone
primitive that flips the one status line — the single source of the flip logic,
independently testable. `bench commit` does the gated commit; `bench commit
--spec <slug>` additionally calls the flip so the finishing build-commit is one
command. Plain `bench commit` (no flag) never touches any spec. Rejected:
auto-folding the flip into every commit with active-spec inference — that makes
every mid-build checkpoint a footgun that must remember *not* to flip, and the
"which spec is active" guess is exactly the ambiguity that bites.

## #2: Does `bench commit` run the gate, or trust the caller?

Type: Grill

### Question
The `bench shift` loop already gates-then-commits internally. Should the
interactive `bench commit` re-run the gate, or assume the caller (final-check)
already did?

### Answer
**Gate-inside: `bench commit` runs the gate and commits only on green.** It is
the mechanization of invariant 1 ("commit on green, never on red"), so the
oracle belongs in the command, not in prose the agent must remember. On red it
reports the first failure and exits nonzero — commit-on-red is impossible by
construction. This collapses `/bench-final-check`'s "gate → report → commit on
green" into one call and `/bench-implement-spec`'s finishing commit into `bench
commit --spec <slug>`. Cost accepted: a second gate run at the finish, which is
the moment you most want the freshest verdict (files can change after the last
gate). Rejected: pure-plumbing commit that trusts the caller already gated — it
re-creates the footgun, since files can drift between the gate and the commit.

## #3: What does `bench commit` do with working-tree files it wasn't named?

Type: Grill

### Question
The gate runs against the whole working tree, but commit lands only the named
files. Changes outside the named set make the green verdict describe files that
won't be committed. Current prose: "an unexplained working-tree file blocks the
commit." Mechanize *how* it blocks.

### Answer
**Block.** Refuse — before gating — if any tracked-modified or untracked file
exists outside the named set (plus the `--spec` file). List the offending files,
exit nonzero; the reviewer resolves them (their call: "surface it, don't commit
or revert on your own"), then re-runs. Order: clean-except-named check first
(fail fast before the ~40s gate), then gate, then commit. This keeps the
invariant honest — when `bench commit` succeeds, the tree *equals* the committed
diff, so the green gate is a verdict on exactly what landed. Rejected:
stash-isolate (git-stash-in-a-command is a partial-state / interrupted-mid-stash
footgun for a case the reviewer should just look at); commit-named-and-warn (a
warning the agent scrolls past defeats the point).

## #4: Does `bench commit` carry a commit-time branch guard?

Type: Grill

### Question
final-check prose says "lands on the working branch, never the default branch."
But this repo's working branch *is* `main` (its default). A naive "refuse on
default" guard would refuse every commit here. Port the guard into code, or not?

### Answer
**No commit-time branch guard — `bench commit` is branch-agnostic.** The
harness-independent branch protection is the pre-push hook (guards direct push to
the default branch); that is the irreversible, outward step. Duplicating the
branch rule at commit-time is a second derivation of a fact the hook already
owns (violating one-source-per-fact), it breaks main-is-working-branch repos like
this one, and it would need new machinery to parse the working branch out of
`projects/<name>.md` (nothing does today). This *removes* the final-check prose
guard rather than porting it. Rejected: refuse-on-default (wrong here); parse the
profile's working branch (new seam for a rule the hook enforces).

## #5: How does `bench spec implemented` behave? (resolved inline — veto at spec)

Type: Proposal

### Question
Resolution, validation, and staging posture of the flip.

### Answer
**Single responsibility: rewrite the one status line, nothing else.** Resolve the
spec path-first, then `specs/<slug>.md` for a separator-free arg (reuse
`coverage.Command`'s exact convention). Require exactly one line-start `Status:
staged`; rewrite it to the exact retirement-detector form
`^Status:[ \t]+implemented[ \t]*$`. Error (nonzero, naming file + reason) on:
not-found, no `Status: staged` line (already implemented, or missing), or more
than one. Edit the file only — never stage; staging is `bench commit`'s job, and
when composed via `--spec`, commit adds the spec path to its named set. Re-run on
an already-implemented spec errors (no staged line), non-destructive. One
mechanic for the spec to record (not a fork): `bench commit` refuses an empty
commit (named paths that produce no staged change) — nonzero, "nothing to commit."

## Handoff

1. **Module boundaries.**
   - `internal/spec` (new): owns `bench spec implemented <slug>` — resolve,
     validate, exact-format flip. The single source of the flip logic.
   - `internal/commit` (new): owns `bench commit [-m <msg>] [--spec <slug>]
     <path>...` — orchestrates block-check → gate → stage → commit. Reuses
     `internal/gate`'s runner; does **not** touch `internal/shift`'s loop.
     When `--spec` is set, composes `internal/spec` for the flip.
   - `bin/bench.sh`: two new routes — top-level `commit`, and a `spec` namespace
     with `implemented` — each routing into the Go core so every surface (kit
     CLI, linked by-path CLI, hooks) reaches the same implementation.
   - `.agents/commands/bench-implement-spec.md` and `bench-final-check.md`: the
     duplicated commit-discipline prose collapses to pointers at the two wrappers.
     This is the "replaces footgun prose" half of FT3, in the same feature.
2. **Contracts.**
   - `bench spec implemented <slug>`: exit 0 flip written (one confirmation line);
     exit 1 not-found / no-single-staged-line (message names file + reason); exit
     2 usage. No stdin.
   - `bench commit [-m <msg>] [--spec <slug>] <path>...`: requires ≥1 path and
     `-m`. Exit 0 committed (commit summary on stdout). Exit 1 on unexplained-file
     block (offending files on stderr), gate-red (first failing phase +
     reproduction), flip failure, or empty commit. Exit 2 usage (no paths, no
     `-m`, unknown flag). Gate command is the project gate; commit forms no
     opinion of its own.
   - Phase pointers name their target by the token the stale-reference sweep
     recognizes (`bench commit`, `bench spec implemented`).
3. **Deep vs thin.** `internal/gate` is the existing deep unit hiding gate
   resolution/run — commit composes it, adds no gate logic. `internal/spec` is
   deep over the flip (resolution + validation + exact-format rewrite), the
   single source composed by commit. `internal/commit` is a thin orchestrator;
   its seam is the CLI boundary, not internal helpers.
4. **Black-box assertables.** Go table tests on both Commands, asserting git
   state / file content / exit code / stderr:
   - `spec implemented`: staged→implemented produces the exact retire-regex line;
     path arg and bare-slug arg both resolve; slug with a separator gets no
     fallback; not-found names the tried forms; no-`staged`-line errors;
     already-`implemented` errors; only that line changes (rest byte-identical,
     trailing-newline preserved).
   - `commit`: clean-tree-except-named commits only the named paths; an
     unexplained tracked-modified *or* untracked file → block, nonzero, no
     commit, file listed; gate-red → no commit, nonzero, first failure reported;
     gate-green → commit lands on the current branch (including `main`); `--spec`
     → spec flipped and included in a single commit; never `git add -A` (a second
     untracked file never rides in); empty commit refused.
   - Prose: the phase pointers and the stale-reference sweep are the observable
     surface for the doc half; cold-session legibility of the collapsed prose is
     review-phase judgment, named in testing decisions as not TDD-able.
5. **Gate attachment.** `bench gate`'s `go test` phase runs both packages' table
   tests — fully gate-observable (git state + exit codes are black-box), no
   manual-verify seam. The conformance phase gains a route-anchor asserting
   `bin/bench.sh` routes `commit` and `spec implemented` (invocation-through-
   every-surface), so a later edit cannot silently drop the route. Per
   `craft-gate`, the spec records proving that anchor bites (remove a route once,
   red, restore).
6. **Hostile-input owners** (profile checklist, `projects/benchkit.md`):
   - paths with spaces/glob chars → explicit `:(literal)` pathspecs (prior art:
     `internal/shift/shift.go`); table-tested.
   - spec file with no trailing newline (hand-edited) → flip preserves byte
     structure, only the status line changes; asserted.
   - absent vs empty spec → distinct: absent → not-found; empty/no-staged-line →
     no-staged-line error.
   - git missing from PATH → error, nonzero (git is the core dependency).
   - invocation through every surface → both commands route through the core from
     kit CLI, by-path CLI, and hooks; conformance route-anchor asserts it.
   - re-run idempotency → second `spec implemented` on the same spec errors
     (non-destructive); `bench commit` re-run after a clean commit → empty-commit
     refusal.
   - cwd deeper than repo root → resolve `git.Root` before path/spec resolution;
     slug fallback is repo-root-relative `specs/<slug>.md`, path args are
     cwd-relative — asserted.
7. **Uncertainty flags.** None open — the four forks are grilled closed and #5
   follows from the retirement-detector format. The two verifications the spec
   records: the route-anchor bite-proof (item 5), and that the empty-commit
   refusal (#5) has a red-capable row.
8. **Rejected alternatives.** Auto-fold the flip into every commit / active-spec
   inference (#1); pure-plumbing commit trusting the caller gated (#2);
   stash-isolate or commit-named-and-warn for unexplained files (#3); commit-time
   branch guard / parsing the profile working branch (#4); sharing
   `internal/shift`'s iterate loop for commit (reuse `internal/gate` directly).
9. **Domain watch-outs.** The retirement detector (`internal/status`) fires on
   any unfenced line-start `Status: implemented` in `specs/*.md`; the flip must
   write exactly that form and touch no other spec, or `bench status` mislabels
   an unfinished spec as done. The gate runs against the whole working tree —
   that is why the unexplained-file block precedes the gate; without it a green
   verdict describes files that won't be committed. Commit is branch-agnostic by
   design; the machine-independent branch protection is the pre-push hook — a
   re-added commit-time branch check is a second source of a fact the hook owns
   and breaks main-is-working-branch repos. The phase-prose collapse is a
   leverage artifact — command files compound through every session that loads
   them — so its story rides the `craft-line` leverage override (top + high)
   while the Go wrappers, being exact-spec and gate-covered, route cheap.

Dependency order: single feature, natural two-part build order within it — the Go
wrappers (`internal/spec`, `internal/commit`, the routes) land first, then the
phase-prose collapse that points at commands which must already exist. Whether to
slice that into two specs is the reviewer's call.
