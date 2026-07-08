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

## 2026-07-08 — splitting a long file can trip the dir file-count budget  [open]
- **What happened:** To clear a `FILE TOO LONG` structure flag on `internal/conformance/package_core_checks_test.go` (430 > 400 lines, from the new identity sweep), I split the sweep into its own file. The dir was already at 12/12 source files, so the split immediately tripped `DIR CROWDED` (13 > 12) — one violation traded for another.
- **Right behavior:** Before splitting a file to fix a length flag, check the dir's file-count budget (`bench structure`). If the dir is at its file limit, a split can't help; propose a reviewer grant for the modest overage or a module-grouping refactor instead.
- **Proposed rule change:** Add a line to craft-seams / the implement-spec structure-housekeeping note: "a file-length split is only free when the dir has file-count headroom; check both budgets before choosing split-vs-grant."
