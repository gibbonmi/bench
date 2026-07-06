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

## 2026-07-06 — `bench commit` cannot stage a file deletion  [open]
- **What happened:** Retiring the `commit-wrappers` spec (a promote-then-delete
  pass) required committing a deleted `specs/*.md` alongside a modified profile.
  `bench commit -m … projects/benchkit.md specs/commit-wrappers.md` ran the gate
  green, then failed at staging: `git add :(literal)specs/commit-wrappers.md`
  exits 128 on a path whose file no longer exists. The commit never landed; I
  completed it with plain `git commit` on the already-green, already-staged tree.
- **Right behavior:** `bench commit` should stage a deletion in the named set,
  since spec-retire — a core workflow — always deletes a file. `git add` of a
  removed path fails; the fix is to stage named paths with a mode that records
  deletions (e.g. `git add -A -- :(literal)<path>` or `git rm --cached` on absent
  paths), so the sanctioned commit path can complete the very workflow it exists
  to serve. Falling back to raw `git commit` defeats the block-check + gate-order
  guarantees `bench commit` is meant to enforce.
- **Proposed rule change:** none to the working agreement; this is a `bench
  commit` defect. Roadmap item: teach `bench commit` staging to record deletions
  for named paths, with a gate row driving a deleted path through the command.
