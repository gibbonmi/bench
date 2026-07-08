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

## 2026-07-08 — session guessed the git idiom because bench commit's signature isn't in context  [open]
- **What happened:** During the what-next drain commit, the session ran `git add` then a pathless `bench commit -m <msg>` (exit 2). No always-loaded file states the contract — `BENCH.md`'s CLI inventory lists the bare name `bench commit`, `BENCH-reference.md` doesn't cover it, and the drain command says only "commit on green" — so the git stage-then-commit prior filled the gap. Repro: the pathless invocation exits 2 with usage; `rg 'commit -m' .bench/BENCH.md .bench/BENCH-reference.md` finds no signature.
- **Right behavior:** The inventory line should carry the one-clause contract (`bench commit -m <msg> <path>...`, path-scoped, stages its own paths) so a cold session doesn't reach for `git add`. Interim: check `--help` before first use of any bench subcommand with arguments.
- **Proposed rule change:** Amend the `BENCH.md` CLI inventory's work-execution line with the commit signature clause (kit edit under `craft-synthesis`).
