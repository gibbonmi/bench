# Make bench doctor recognize the kit source checkout

Blocked by: none
Writes: internal/adopt/doctor.go, internal/adopt/doctor_rows.go, internal/adopt/adopt_test.go

## What to build

`bench doctor` in the kit's own source checkout reports two reds and one
absent-hook row that `bench link` cannot fix. On 2026-08-30 a link of the kit
into itself landed and broke the shim route and the land route. The unlink
landed after it. The kit source checkout therefore never carries the
managed `AGENTS.md` block or the `.bench/bin/` launcher copy. The doctor
must say so instead of sending the reader to `bench link`.

The doctor recognizes the kit source checkout with one predicate: the
repository root from `git.Root()` equals the kit directory from `kitDir()`
in `internal/adopt/adopt.go`. The wrapper sets `BENCH_KIT` to the kit, and
in the kit repo that is the repository root. A consumer repo never satisfies
the predicate, because its kit lives under a package or a cache directory.
Put the predicate in one function in `internal/adopt/doctor.go` and let the
three rows below read it.

In the kit source checkout:

- The `AGENTS.md` row prints `ok: kit source checkout - AGENTS.md is the source agreement; no managed block applies` and is not red.
- The repo-local bench row prints `ok: kit source checkout - the launcher is bin/bench.sh; no .bench/bin copy applies` and is not red.
- The pre-push row keeps its current text for a managed hook. When the hook is absent, the row names `bench doctor --fix` instead of `bench link`, and `--fix` installs the managed hook through `installGitHook` in `internal/adopt/link_hook.go`. The row stays red until the hook is present.

Outside the kit source checkout, every row keeps its current text and its
current red posture; `--fix` keeps its current behavior for an absent hook.

`internal/adopt/adopt_test.go` holds the doctor tests. Mirror an existing
row test there for the fixture shape. A kit-checkout fixture sets `BENCH_KIT`
to the temporary repository root; a consumer fixture leaves it elsewhere.

## Acceptance

- [ ] With `BENCH_KIT` equal to the repository root and the pre-push hook present, `bench doctor` exits 0. The `AGENTS.md` and repo-local bench rows print the kit-source `ok` text.
- [ ] With `BENCH_KIT` equal to the repository root and no pre-push hook, the pre-push row names `bench doctor --fix`. That command installs the managed hook, and a second `bench doctor` then exits 0.
- [ ] With `BENCH_KIT` elsewhere, the `AGENTS.md` row and the repo-local bench row keep their current red text, and an absent pre-push row keeps `run bench link`.
- [ ] The delegate records the first two rows red before the fix.
- [ ] `go test ./internal/adopt/ -run 'Doctor' -parallel 2` passes.
