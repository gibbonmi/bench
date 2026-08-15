# Learnings — usage journal

## 2026-08-14 — write-spec falsification loop took 33 rounds on worktree-enumeration-hang [open]

What happened: /bench-write-spec's loop 1 needed 33 author-fix/re-review rounds
before ACCEPT (tickets loop: 2). The author drafted from the compiled map plus
targeted code reading; the reviewer's tree-verification kept finding real
holes — the linked-worktree producer input, callers that swallow or relabel
the new typed error, per-child bound scoping, process-status partition arms,
and fixture reachability at command surfaces.
Right behavior: most catches trace to one root: the draft reasoned from the
producing seam outward and enumerated consumers late. Enumerating every
caller/renderer of the changed contract *before* writing coverage rows —
call-graph first, rows second — would have removed roughly half the rounds
(list.go, resume.go, dashboard, PruneLandedBranches were each late finds).
Proposed rule change: none — a proposal was made and empirically rejected.
The candidate craft-spec sentence ("enumerate every production consumer of a
changed contract before writing rows") was A/B tested in Codex (2026-08-14,
3 runs per arm, first-draft grading against the FT189 answer key): treatment
medians tied at 3/5, controls did not improve, and no run in either arm
found the intra-package PruneLandedBranches caller or all eight
--git-common-dir sites. Reviewer verdict: dropped. Residual observation for
the drain: the deep consumer surfaces were found only by adversarial
tree-verifying review, not by authoring guidance — the verification loop,
not the skill, is the load-bearing control.

## 2026-08-14 — CLI changes can leave the canonical command inventory stale [open]

What happened: recent Bench CLI work added public worktree commands without
updating `.bench/BENCH.md`'s always-loaded CLI Inventory, so cold sessions were
given an incomplete command surface.
Right behavior: every change to the Bench CLI's public command surface must keep
the canonical markdown command definition current in the same change. A
codified validator must compare the production command surface and the markdown
inventory in both directions, fail with an attributable diagnostic when a
command or definition is missing or stale, and carry a retained omission
mutation proving that failure bites.
Proposed rule change: add that validator to the conformance layer and make its
green result an acceptance requirement for every public Bench CLI change.

## 2026-08-15 — interrupted semantic review leaves ambiguous pickup state [open]

What happened: the worktree-enumeration-hang semantic-review session was
interrupted after it wrote and committed `reviews/worktree-enumeration-hang.md`.
The retained source therefore looks clean but carries an actionable review
artifact, with no terminal review verdict or handoff that distinguishes it from
an intentionally paused repair pass.
Right behavior: phase recovery must classify a persisted review artifact as
unfinished pickup state, reacquire the exact candidate base and tip, and rerun
all three review axes before either repair or landing. It must never infer a
clean review, acceptance, or landing authority from the artifact or its commit.
Proposed rule change: add an interruption/recovery contract that records the
phase state and frozen pair atomically, then makes resume surface the required
next phase and refuse landing until a fresh terminal review verdict exists.

## 2026-08-15 — write-spec loops accepted by reviewer cap, not by clean re-review

**What happened.** The reviewer set `--reviewer` to Codex `gpt-5.6-sol` and capped both
verification loops of `specs/gate-run-transaction` at one round. Loop 1 returned eight
blocking findings (two named mutations were not red-capable as written; a false
"already covered" citation; missing persistence-failure rows; the deep-module story
was satisfiable by symbol deletion alone). All were folded, but no reviewer round
re-checked the folds; loop 2's per-row split finding was declined by the author as a
cost call and surfaced in the approval table.

**Right behavior.** Under a round cap, the author still folds every blocking finding
and states in the verification log that acceptance was by cap; a declined blocking
finding is a reviewer decision surfaced at sign-off, never silently dropped.

**Proposed rule change.** `/bench-write-spec` step 9: when `--reviewer` carries a
round cap below "until clean", the verification-log line must say "accepted by cap"
and name any blocking finding the author declined, so the sign-off table carries it.
