# Learnings — usage journal

Append one entry when you deviate from the workflow, make a judgment call you're
unsure about, catch a should-have-asked in hindsight, or catch yourself assembling
the same ad-hoc check a second time (a codification candidate — name the `bench`
subcommand it wants to be). You capture; the reviewer
decides. `/bench-what-next` verdicts every open entry in its reviewed batch
diff: work-shaped and rule-shaped entries become roadmap items (rule-shaped ones
built later under the synthesis discipline), the rest are dismissed with one line
of why. A resolved entry leaves this file, and its verdict is recorded in the
drain commit. The journal holds open entries only; history lives in git. Never
rewrite a kit rule yourself — that is the whole point of capturing here instead.

Format per entry:

## <date> — <short title>  [open]
- **What happened:** …
- **Right behavior:** …
- **Proposed rule change:** … (or "none")

An entry leaves this file only via /bench-what-next.

<!-- entries below -->

## 2026-07-07 — transient gate red recurred under back-to-back full-gate load  [open]
- **What happened:** During the ten-commit spec-retire pass, one `bench commit`
  gate run went red on a docs-only deletion; the byte-identical tree gated
  green on the next run. This matches the ROADMAP Watch note (worktree
  concurrent-acquire failed once under full-gate load on 2026-07-06), but the
  failing check is unattributed because the batch loop piped gate output
  through `tail -2`, discarding the failing phase.
- **Right behavior:** Batch loops over `bench commit` capture full output (no
  pipe-masking of exit codes or failing checks), and a recurred flake gets a
  deflake row: make the concurrent-acquire contract test deterministic under
  gate phase concurrency or serialize its acquire window.
- **Proposed rule change:** none — work-shaped; drain to a roadmap row for the
  deflake.

## 2026-07-07 — bench commit cannot stage a deleted or renamed path (git 2.43)  [open]
- **What happened:** structure-splits deletes axi_wave2_test.go and renames
  line_routing_checks_test.go. `bench commit` gated green and flipped the spec,
  then died at staging: its per-path `git add -A -- :(literal)<p>` errors
  `fatal: pathspec ... did not match any files` for every path absent from the
  worktree — each deletion and the old side of each rename. I finished with a
  direct `git commit` (index already fully staged by an earlier `git add -A`,
  oracle already green), then captured this.
- **Right behavior:** bench commit's staging must handle a named deletion/rename
  (e.g. detect an absent tracked path and stage its removal via `git rm --cached`,
  or `git add -A -- <dir>`), so any split/rename spec can use the sanctioned
  commit path instead of falling back to raw git.
- **Proposed rule change:** none — a kit defect; drain to a roadmap row to fix
  bench commit's deletion staging (its own doc comment claims "a named deletion
  included", which git 2.43 contradicts for literal-pathspec `git add`).

## 2026-07-07 — structure-splits landed on a drifted tree; "structure → 0" needed an out-of-spec call  [open]
- **What happened:** The spec assumed 4 over-budget items (2 split, 2 seed). By
  build time other merges had pushed the tree to 6: a new FILE TOO LONG
  (internal/contract/surface/link_test.go, 425) and a new DIR CROWDED
  (internal/contract/runtime/, 14 files). After the specced splits + seed,
  `bench structure` reads 2, not 0. I did not expand the seed to hit 0 —
  accepting a DIR CROWDED "group into modules" signal is a reviewer split-vs-accept
  call outside this spec's seed — and flagged it plus parked an idea instead.
- **Right behavior:** correct — implement the spec's scope, keep the gate green,
  surface fresh drift as a reviewer decision rather than suppress a new signal to
  hit a number.
- **Proposed rule change:** none — parked as an idea for the residual two.
