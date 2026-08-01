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

## 2026-08-01 — Parking capture blocks the spec-build lifecycle  [open]

**What happened.** Mid-way through the `reduced-gate-phase-set` build, the reviewer
asked me to park a recommendation. I ran `bench idea`, which appended one line to
`capture/IDEAS.md`. The next `bench spec build checkpoint` refused:

    error: spec build start requires a clean working checkout:  M capture/IDEAS.md

Three verified tickets were blocked behind a one-line capture edit. Clearing it
myself was correctly denied by the `block-dangerous-git` hook (`git checkout <path>`
is not the agent's authority), so the build stopped for a reviewer decision. Writing
*this* learning has the same effect: `capture/learnings.md` is a tracked file, so the
act of recording the problem reproduces it.

**Why this is a real conflict, not a papercut.** Two rules in the kit disagree:

- `.bench/BENCH.md`'s Capture section makes parking deliberately frictionless —
  *"Parking an idea is conversational — never a CLI chore for the reviewer"* — and
  assigns the running agent the job of doing it the moment a tangent appears.
- The spec-build lifecycle requires a clean working checkout for its mutations.

So during any spec build — the workflow's own heaviest phase, and the one most likely
to surface tangents worth parking — following the capture rule breaks the lifecycle.
The agent must either refuse to park (violating the capture rule), pay a full ~10
minute gate to commit one line, or stop and ask. All three are wrong answers to
"write this down before we lose it."

**What the right behavior was.** Park to a staging area that is not the graded tree,
and land it as part of the build's own capture commit. I improvised this after the
fact by holding the text in the session scratchpad; that works but depends on the
agent remembering, and a session that ends first loses the idea entirely — which is
precisely the failure `bench idea` exists to prevent.

**Proposed rule / tooling change.** Either:

1. **Exempt the capture surfaces from the lifecycle's clean-checkout precondition.**
   The declared allowlist from this very spec (`capture/`, `specs/`, `ROADMAP.md`,
   `.bench-notes.md`) already names the paths no gate phase grades behaviorally. A
   dirty path inside it cannot invalidate a checkpoint's evidence, because the
   checkpoint's subject is the assignment worktree, not the main checkout. This is
   the smaller change and it composes with the feature being built here.

2. Or **give `bench idea` a staging mode** that writes outside the graded tree and a
   drain step that folds staged entries into `capture/IDEAS.md` at the next capture
   commit.

I'd take (1): the precondition exists to stop *uncommitted implementation* from
riding into a checkpoint, and capture files are exactly the class that cannot.

**Two smaller defects found alongside, both worth fixing with it.**

- The refusal message hardcodes `start` regardless of which mutation raised it
  (`internal/specbuild/precondition.go:105`). It fired on `checkpoint` and said
  "spec build start requires...". I initially read the shared precondition as
  start-only *because of that message*, told the reviewer parking was safe, and was
  wrong.
- `bench spec build checkpoint`'s receipt `rows[].outcome` enum is
  `passed | already-covered | not-tdd-able`, and nothing advertises it. An invalid
  value returns only `invalid spec build receipt`, with no indication which of the
  ~15 validated fields failed. Parked separately as the receipt-assembly item.

## 2026-08-01 — Ran an interactive porcelain in automation and leaked two worktrees  [open]

**What happened.** Looking for `bench worktree`'s usage string, I ran the bare verb in a
non-interactive Bash call. It hung for two minutes and was killed with SIGTERM (exit
143). I repeated it twice more while diagnosing. Each invocation created a real
worktree; the two that were signal-killed leaked, and I retired them afterwards with
`bench worktree clean` (plan then apply, `count=0 bytes=0 recovery=none`, nothing lost).

**Why it happened.** `bench worktree [objective...]` is not a usage-error verb — it is a
human porcelain (`internal/worktree/worktree.go:412`, `Subshell`). It creates a worktree,
prints `🪵 worktree: <path>  (exit to release)`, then runs `$SHELL` inside it with stdin
inherited and blocks on `cmd.Run()` until the shell exits. With an inherited open stdin
in automation, that wait never ends. The non-interactive surface is the subcommands:
`list`, `path`, `exec`, `release`, `clean`, `recovery`.

**What the right behavior was.** Read the usage rather than probe for it — `bench
commands --brief`, or the help block in `bin/bench.sh` (which documents exactly the four
worktree subcommands and not the bare form). Never invoke an unrecognized `bench` verb
bare in a non-interactive call to discover what it does: in this CLI a bare verb can be a
porcelain that acts, not a parser that refuses. Where a probe is genuinely needed,
`</dev/null` bounds it — an interactive command then sees EOF and exits immediately,
which is also what distinguished the hang from a slow command during diagnosis.

**Proposed rule change.** Add to `AGENTS.md`'s shell conventions for agents in this repo:
discover a subcommand's shape from `bench commands --brief` or `bin/bench.sh`, never by
running the bare verb; and redirect stdin from `/dev/null` for any `bench` invocation
whose interactivity is not already known. This is a one-line convention that would have
prevented all three invocations.

**Separable tool defect, parked as its own idea.** `Subshell` releases the worktree only
on the line after `cmd.Run()` returns, so the release is not signal-safe: SIGTERM or
SIGINT to the `bench` process skips it and leaks a registered worktree. The clean-exit
path works correctly. See the staged idea for the proposed fix.
