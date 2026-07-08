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

## 2026-07-08 — sanctioned commands lose to raw-git priors at the moment of use  [open]
- **What happened:** Two instances in one drain session. (1) The drain commit ran `git add` then a pathless `bench commit -m <msg>` (exit 2) — no always-loaded file states the contract: `BENCH.md`'s CLI inventory lists the bare name, `BENCH-reference.md` doesn't cover it, the drain command says only "commit on green", so the git stage-then-commit prior filled the gap. Repro: the pathless invocation exits 2 with usage; `rg 'commit -m' .bench/BENCH.md .bench/BENCH-reference.md` finds no signature. (2) Verifying FT42's retirement, the session hand-ran `git log` + `git show --stat` instead of the sanctioned `bench spec history npm-identity` — the command existed but wasn't recalled at the moment of use (~250 extra tokens of unfiltered output; the named shortcut lives only in debug-skill prose).
- **Right behavior:** The inventory line should carry the one-clause contract (`bench commit -m <msg> <path>...`, path-scoped, stages its own paths) so a cold session doesn't reach for `git add`, and recall cues for sanctioned shortcuts need to live where the work happens (the what-next reconcile step naming `bench spec history` for shipped-row checks), not only in another skill's prose. Interim: check `--help` before first use of any bench subcommand with arguments.
- **Proposed rule change:** Amend the `BENCH.md` CLI inventory's work-execution line with the commit signature clause, and add `bench spec history` as the named tool in the what-next reconcile step (kit edits under `craft-synthesis`). Don't codify new subcommands for the drain's bespoke spot-checks — the second-assembly bar stands.

## 2026-07-08 — bench commit --spec flips Status on the wrong commit when misused  [open]
- **What happened:** Third contract-guess in one day: committing the freshly *staged* spec, the session passed `--spec defect-batch-ft43-49` assuming it merely associated the commit with the spec. The flag's real behavior is `bench spec implemented`'s flip — it rewrote `Status: staged` to `implemented` inside the authoring commit, so a spec whose build never started shipped claiming implemented. Caught immediately; corrected in the next commit.
- **Right behavior:** `--spec` belongs only on the implementation's green commit. The usage string (`[--spec <slug>]`) states no semantics, so the same fix as the signature entry applies: the one-clause contract must be visible at the point of use.
- **Proposed rule change:** Extend the usage string to `[--spec <slug: mark implemented>]` (or similar) and add the clause to the BENCH.md inventory line alongside the signature fix already proposed above — one kit edit covering both.

## 2026-07-08 — main-tree vs shared-worktree path/CWD confusion bit both delegate and orchestrator  [open]
- **What happened:** Two mirrored instances of the same confusion during the FT43 fix, which was built in the *existing* FT43 worktree (a sanctioned shared-worktree case — the fix stacks on the prior commit). (1) **Delegate side:** the charge named the worktree by `cd` and used relative Bash paths, but the delegate's Write/Edit tool calls used repo-root absolute paths (`/home/mgibs/workspace/bench/tests/...`) instead of worktree-absolute paths — so its first canary fixture landed in the *main tree*, its Read/Edit looked stale, and the red run reported "fixture not found." It caught it, `rm -rf`'d the misplaced dir, verified main clean, redid the work in the worktree. (2) **Orchestrator side:** the coordinating session edited the main tree's `.bench/learnings.md` via an absolute-path Edit while its own Bash CWD had persisted in the worktree (from an earlier verification `cd`). `bench commit .bench/learnings.md` then reported "nothing to commit — the named paths produced no staged change," because it correctly resolved its root from CWD (the worktree), where that file is unmodified. This was **misreported to the reviewer as a `bench commit` wrong-root bug**; a repro showed `bench commit` was correct — running it from the main tree committed cleanly. Net deliverable effect on FT43: none.
- **Right behavior:** Worktrees have independent working trees; a file edited in one is invisible in the other even though `.git` is shared. Both a delegate's file-tool paths AND the orchestrator's Bash CWD must be pinned to the *same* tree the edit targets. `bench commit` resolving root from CWD is correct, not a defect — the failure mode is a stale CWD, so a "nothing to commit" on a file `git status` shows modified elsewhere is a CWD/tree mismatch to check first, never a tool bug to file. `cd` fixes Bash CWD only, not absolute paths passed to file tools.
- **Proposed rule change:** (a) Add one clause to `craft-delegate`'s Isolation section: when a write-delegation *shares* an existing worktree (a stacked fix), the charge must pin all file-tool paths to the worktree root, since `cd` governs Bash only. (b) Note in the same place that an orchestrator coordinating a shared worktree must `cd` back to the main tree before main-tree commits, and read a `bench commit` "nothing to commit" against a visibly-modified file as a CWD mismatch. Kit edits under `craft-synthesis`. No `bench commit` code change — it behaved correctly.
