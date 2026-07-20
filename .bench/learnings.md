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

## 2026-07-20 — FT88 spec premise was stale against a landed feature  [open]

- **What happened:** `specs/minimal-subprocess-data-exposure.md` states as its problem
  that "a project gate gets everything except `BENCH_KIT` and `BENCH_WRAPPER`", and
  builds story 2 plus four coverage rows on that premise. It is false: FT78's closed
  gate subject already launches the project gate with `PATH` plus only the names
  declared in `.bench/gate-inputs.json`. The build reached stage 1b before the seam-A
  contract delegate hit it, having already wired `internal/env` into a `gateEnv` whose
  only live caller is the kit's own phase runner.
- **Right behavior:** a spec compiled from a closed decision map should verify its
  problem-statement premises against the current tree before locking coverage rows,
  not inherit them from the map's snapshot. Here one `rg` for the gate's env
  construction at spec time would have caught it. The map closed before FT78 landed.
- **Proposed rule change:** `/bench-write-spec` (or `craft-spec`) should require that
  each "today the code does X" claim in a Problem section be checked against the tree
  at spec time, with the check named in the spec — the same standard the coverage map
  already applies to its red signals.

## 2026-07-20 — reproducibility probe reds on uncommitted tracked-file deletions  [open]

- **What happened:** FT88's stage 3 deleted `internal/intent/preview.go` (moved to
  `internal/sanitize`). The gate's contract phase then failed
  `TestGoBuildIgnoresCheckoutTopology`: its clone-and-overlay fixture copies the
  working tree over a HEAD clone with tar, which propagates adds and modifies but
  not deletions, so the clone built with the resurrected deleted file and the
  binaries legitimately differed. Under batch approval I had a delegate fix the
  fixture to a true mirror (deletions included) rather than stall — a
  gate-attached contract test edited without a live reviewer sign-off, flagged
  here and in the exit report for veto.
- **Right behavior:** the fixture fix is the faithful one (it strengthens the
  probe rather than weakening it), but gate-check edits mid-build should always
  land with an explicit flag, and the probe should have mirrored deletions from
  the start so a refactor that moves a file cannot red the gate for a fixture
  reason.
- **Proposed rule change:** none — the fixture fix closes the gap; this entry is
  the veto surface for it.
