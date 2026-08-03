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

## 2026-08-02 — mutated the tree while a gate was running  [open]

- **What happened:** Ran `bench idea` (a capture write) while a background
  `bench gate --fresh` was mid-run; the gate finished green on every phase but
  rejected its own verdict with "gate subject changed during execution", and
  the full run's cost was paid twice.
- **Right behavior:** The profile's cold-session note already states it: never
  mutate the repository while a gate is running — queue capture writes and
  roadmap edits until the verdict lands, or do them before starting the run.
- **Proposed rule change:** none — the rule exists; this is an execution miss
  worth counting, and the queue-until-verdict habit is the fix.

## 2026-08-03 — tip moved under an active spec-build run with zero checkpoints  [open]

- **What happened:** Committed capture files (handoff) while a spec-build run
  was active but before any checkpoint existed; the moved tip forced
  recomposition, which is unrecoverable at zero checkpoints (empty patch), and
  the run had to be abandoned and restarted with delegate diffs replayed.
- **Right behavior:** Freeze the tree for the life of an active spec-build
  run — no capture commits, no roadmap edits — until promote lands; every tip
  move mid-run costs at least a fresh full gate for recomposition, and before
  the first integration it costs the run.
- **Proposed rule change:** the phase-boundary handoff write in `--full` runs
  should happen before `bench spec build start` or after promote, never
  between; candidate for the bench-implement-spec command text.

## 2026-08-03 — a ticket row's named owner was never enumerated against the tree  [open]

- **What happened:** FT164's teach ticket claimed a `template heading depth`
  needle and mutation row (TT10); the delegate demonstrated the red as a
  self-probe but landed no artifact, and three review rounds' needle↔row
  enumerations never caught it because they keyed on registered needles, not
  on ticket-row owners. Found only when the fence-aware fix went looking for
  the artifact to retire.
- **Right behavior:** verifying a ticket's done-claim includes resolving every
  named owner in its Red-mutations table to a real artifact in the tree — the
  same resolve-the-identifier rule craft-delegate already states for absence
  claims.
- **Proposed rule change:** one sentence in craft-delegate's done-claim list
  or craft-tickets' evidence rules: a mutations-table owner that names a test
  or check is verified by finding it, not by the delegate's demonstration.
