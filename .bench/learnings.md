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

## 2026-07-08 — write-delegate sharing a worktree edited the main tree via repo-root paths  [open]
- **What happened:** The FT43 fail-open fix was delegated into the *existing* FT43 worktree (a sanctioned shared-worktree case — the fix stacks on the prior commit). The delegate's charge named the worktree by `cd` and used relative paths for Bash, but its Write/Edit tool calls used repo-root absolute paths (`/home/mgibs/workspace/bench/tests/...`) instead of the worktree-absolute paths — so its first canary fixture landed in the *main tree* (pre-FT43, no Timeout field), which is why its Read/Edit looked stale and the red run reported "fixture not found." The delegate caught it, `rm -rf`'d the misplaced dir, verified the main tree clean, and redid the work in the worktree. Net deliverable effect: none (verified: main tree clean, worktree scope correct, string-timeout probe red on the real oracle, clean gate green).
- **Right behavior:** A write-delegate sharing an existing worktree must target every file tool (Write/Edit/Read) at worktree-absolute paths, not repo-root paths — `cd` only fixes Bash CWD, not the absolute paths passed to file tools. The charge should either state the worktree root once and require all tool paths be built under it, or the isolation mechanism should hand the delegate a checkout where repo-root *is* the worktree (as `Agent isolation: worktree` does — the miss is specific to the shared-existing-worktree pattern this fix needed).
- **Proposed rule change:** Add one clause to `craft-delegate`'s Isolation section: when a write-delegation *shares* an existing worktree (a stacked fix), the charge must pin all file-tool paths to the worktree root, since `cd` governs Bash only. Kit edit under `craft-synthesis`.
