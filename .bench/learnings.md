# Learnings — usage journal

Append one entry when you deviate from the workflow, make a judgment call you're
unsure about, or catch a should-have-asked in hindsight. You capture; the reviewer
decides. `/resynthesize` reviews the open entries, promotes the generalizable ones
into the kit with sign-off, and marks them resolved. Never rewrite a kit rule
yourself — that is the whole point of capturing here instead.

Format per entry:

## <date> — <short title>  [open]
- **What happened:** …
- **Right behavior:** …
- **Proposed rule change:** … (or "none")

An entry only becomes `[resolved]` via /resynthesize.

<!-- entries below -->

## 2026-06-27 — gate was red at HEAD; a commit claimed green without running it  [open]
- **What happened:** Fixing the missing `.bench/learnings.md` scaffold, I ran the
  full gate and found it red at clean HEAD — commit `cea2f42` renamed the commands
  and its message asserted "Gate green: the index<->disk conformance check confirms
  no reference was missed," but AGENTS.md still carried the old names
  (`/map`, `/diagnose`, `/review`, `/verify`). The one file the conformance check
  reads was the one the rename missed, which means the gate was never actually run
  after that commit. Same failure class as `d77063c` (a rename that skipped a
  reference).
- **Right behavior:** Never write "gate green" in a commit message from belief; run
  `bench gate` and paste/observe its exit. The oracle is the gate, not the diff.
- **Proposed rule change:** Add a one-line reminder to the build/verify-gate path (or
  a pre-commit nudge) that "gate green" in a message must come from an actual run,
  not inspection. Possibly a Stop-hook check that the gate was run since the last
  edit.

## 2026-06-27 — scaffolded files must be created by init, not just referenced  [open]
- **What happened:** `bench init` scaffolded `.bench/gate.sh` but not
  `.bench/learnings.md`, while AGENTS.md (write side) and `/resynthesize` (read side)
  both depend on that file existing. The self-learning contract pointed at a path
  nothing created.
- **Right behavior:** Any file the kit's prose instructs an agent to read or append
  to must be produced by `init` (guarded, idempotent) and locked by a gate check that
  exercises the real init path.
- **Proposed rule change:** When adding a contract that names a `.bench/*` file,
  the same change must (1) scaffold it in `init()` and (2) add a behavioral gate check.
  Consider generalizing gate check 1d to assert every kit-referenced `.bench/*` file
  is scaffolded.
