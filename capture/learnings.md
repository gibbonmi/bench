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

Format per entry. Heading: `## YYYY-MM-DD — short title  [open]`

- **What happened:** …
- **Right behavior:** …
- **Proposed rule change:** … (or "none")

## 2026-08-01 — a tree-wide sweep skipped the kit's own canonical files  [open]

- **What happened:** The capture co-location migration rewrote about 278 path
  references across the tree, driven by `rg -l` over the old paths. Ripgrep skips
  hidden paths by default, so `.bench/` and `.agents/` — where `BENCH.md`,
  `BENCH-reference.md`, and every phase command live — were never in the file list,
  while the canary *fixture* copies of those same files were, because the fixtures
  store them under literal `dot-bench/` and `dot-agents/` directory names. The sweep
  therefore updated each fixture and left its real counterpart untouched, which is
  the exact desync the conformance anchors exist to detect. The gate caught all
  seven, named each anchor, and the second sweep with `--hidden` closed them. Cost
  was one red gate cycle; nothing shipped wrong.
- **Right behavior:** Two things, and the first is the generalizable one. A
  repository-wide sweep in this repo must include hidden paths, because the kit's
  canonical guidance lives in dot-directories while its test fixtures live in
  visible ones — so the default tool posture updates the copy and misses the source,
  which is worse than missing both. Second, `AGENTS.md` already tells agents to
  prefer `rg` over `grep` in this repo, so it is the surface that induced the
  default and the natural place to qualify it.
- **Proposed rule change:** Extend the "Shell conventions for agents in this repo"
  bullet in `AGENTS.md`: when sweeping or auditing repository-wide rather than
  searching for one symbol, pass `--hidden` (and exclude `.git/**`), because
  `.bench/` and `.agents/` hold the canonical files and their canary fixtures shadow
  them under visible `dot-*` names. Worth weighing against a mechanical alternative:
  a conformance check asserting that each `dot-*` fixture copy still matches the
  real file it shadows would catch this class regardless of tooling, and would
  remove the instruction rather than duplicating it.

## 2026-08-01 — no sanctioned way to name "everything currently changed"  [open]

- **What happened:** Landing the migration took three `bench commit` attempts, all
  refused before the gate at no real cost, but all spent on assembling a path list
  rather than on the diff. `bench commit -m … .` does not expand `.` to the changed
  set, so the renamed-away paths fell outside the named set and were refused. The
  obvious fix, `git diff --name-only HEAD`, silently omits them too: rename
  detection collapses a move into a single entry naming only the destination, so the
  deletions are invisible until `--no-renames` is passed. The working incantation
  was `{ git diff --name-only --no-renames HEAD; git ls-files --others
  --exclude-standard; } | sort -u` piped through `xargs`.
- **Right behavior:** The path-scoped refusal is correct and I would not weaken it —
  it is what makes a gate verdict answer for exactly one diff. What is missing is a
  sanctioned way to say "the whole current working set, deliberately", so that a
  reviewer-approved whole-tree change does not require reconstructing a
  rename-aware path list by hand. Reaching for `git` to compute the argument to
  `bench commit` is the smell: the porcelain is making the session assemble a fact
  the tool already holds.
- **Proposed rule change:** A codification candidate rather than a rule — `bench
  commit --all -m <msg>`, an explicit opt-in that names the full working set as the
  attribution target, keeping the default path-scoped and the refusal intact. It is
  adjacent to FT166's `bench capture commit` porcelain and to FT98's set-aside
  primitive, so it should be weighed with them rather than built alone; all three
  are faces of the same missing "name this working set" vocabulary.
