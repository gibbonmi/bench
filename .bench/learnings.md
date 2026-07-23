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

## 2026-07-23 — a committed diagnosis report outlived its own fix  [open]

- **What happened:** The session briefing told me to reproduce a flake, read the
  goroutine dump, and — pre-authorized — fix at the fsync seam. It was built from
  `GATE-REPORT.md`, which reads `Status: diagnosis partial, no fix landed`. But
  the fix had landed two commits earlier (`380fd00`), the spec was already
  `Status: implemented`, and `session-handoff.md` correctly said so. The report
  was committed in `96ddc9f` *after* the fix, carrying a body written at an
  older baseline. I checked `git log` before running anything and reframed the
  task from "hunt a repro" to "run the acceptance load window", which is what
  the handoff actually asked for. Had I followed the briefing literally I would
  have spent three gate runs hunting a failure the fix had already removed, and
  then possibly re-fixed a fixed seam.
- **Right behavior:** What I did — reconcile a report's stated status against
  `git log` and the spec status before acting on it, and surface the conflict to
  the reviewer rather than silently following either source. Invariant 3 says
  docs describe the current decided state; a report that says "nothing is
  committed" while sitting in a commit violates that.
- **Proposed rule change:** A capture-style report committed to the tree must
  carry its status at the top and be re-read at the moment of commit, not only
  at the moment of writing — or, better, `bench` should refuse to commit a doc
  whose own body claims uncommitted state. Candidate codification: extend the
  gate's doc conformance check to flag phrases like "nothing is committed" /
  "uncommitted work in the tree" in a tracked file. Weak signal, cheap check;
  reviewer's call whether the false-positive rate is worth it.
