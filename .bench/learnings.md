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

## 2026-07-11 — grill parked an answerable question as fog  [open]
- **What happened:** During the session-resume-cleanup shape-idea grill, the ledger
  row-lifecycle question was parked under "Not yet specified" as dependent on the
  open Research ticket, and the grill stopped — but its policy half (expire on what
  proof, silent vs surfaced) was answerable by the present reviewer immediately.
- **Right behavior:** Before stopping a grill, split each fog line into
  reviewer-answerable policy vs evidence-dependent detail; grill the policy half now
  and leave only the evidence half as fog.
- **Proposed rule change:** none — craft-grill already covers this ("stop when the
  remaining questions no longer change what gets built"); this is an application miss.
