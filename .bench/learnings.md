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

- 2026-07-09 — Staging a spec for a new phase command turns the gate red: the
  stale-command-reference sweep flags `/bench-<new>` in specs/ before the command
  file exists (hit on specs/assess-owner.md naming the planned assessment phase).
  What we did: held the spec out of the tree and land it in the same diff as the
  command, so the sweep sees the file. Right behavior unclear — options are (a)
  accept spec-lands-with-implementation for new-command specs, or (b) teach the
  sweep that a `Status: staged` spec may name its own deliverable. Proposed rule:
  (a), recorded in /bench-write-spec, unless the reviewer prefers (b).
