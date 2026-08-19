# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `d4a46f82` pre-commit, clean but for this drain's own batch, 37 unpushed commits
Spec: none staged — `specs/` is empty
Gate: green at `d4a46f82`

## State

**The opt-in `--restructure` board pass ran.** Capture was already empty — no ideas, no
learnings, no retros, no staged specs — so this pass was reconcile plus restructure only.
The board went from 73 rows to 71 while gaining seven new ones, because ten merged away.

The structural change worth knowing: the 2026-08 capability audit's action items drove the
recommended sequence while existing only inside the audit directory, so the board could not
be read on its own. They are rows now — FT226 (A6, `BENCH_HOME` leak), FT227 (A4, adoption
smoke), FT228 (A5, `/bench-debug`), FT229 (A9, hygiene batch), FT230 (A7, release
workflow), FT231 (A11, measurement harness), FT232 (A10, repair-loop tripwire) — and the
sequence names board rows rather than external audit IDs.

Merges, each over one owner surface: FT205 + FT221 into FT213 (one `craft-delegate` visit),
FT192 + FT209 into FT214 (one `craft-spec` visit), FT111 into FT179 (one comments visit),
FT206 into the new FT233, FT112 into FT228, FT138 + FT170 into FT231, FT180 and FT177's
landing-guard clause into FT229. FT199 split: the coordinator stayed, its six landing
diagnostics became FT233.

Thirteen accreted rows were pruned to what is genuinely open, each verified against the
tree by a read-only delegate and spot-checked here — not taken on the audit's word. Two of
the audit's own dispositions were stale and were not followed: FT89's registry half had
already shipped with the front door (only `commands --brief` and the guide's CLI inventory
remain), and FT197's premise is inverted, since the Go runner now owns the process group and
calls the gate script rather than sitting behind it, so FT197 dropped to LOW and parked.

One finding landed as evidence rather than prose: this file's previous pins, `HEAD 0ee2106`
and `green at 761a839`, both resolve to *trees*, not commits — they never named anything.
That is now FT162's second occurrence, and it is why the pins above were generated from the
tree.

## Next command

`/bench-write-spec` for FT226 — kit tests stop writing into the operator's real
`BENCH_HOME`. Rank 1 on the refreshed sequence, no dependencies, and the only actionable
row whose defect damages the operator's own machine (759 orphaned pool directories, 41 MB,
roughly ten added per gate run).

## Shape

Rewritten in full at every phase close, pruned rather than accreted: a fresh
session pays for every line it reads cold, so drop anything it would not act on.
Operational gotchas are placed by lifetime, not copied here: one that recurs across
phases belongs in `projects/benchkit.md`'s cold-session notes, and one scoped to a
build belongs in that spec's coverage rows. This file names at most when you'll hit
one, never the command — a second copy drifts from the source.
Keep the three sections above — **State** (what is true now, including anything
uncommitted), **Next command** (the exact harness-native invocation, not a
description of it), and this one.

The handoff carries no date of its own. `bench status` computes its age from the
commit that last wrote this file and reports a `handoff` row once anything has
landed since. Where this document and the tree disagree, the tree wins.
