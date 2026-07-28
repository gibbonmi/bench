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

## 2026-07-28 — FT91 stage-1 floor estimate keyed on the wrong tail

What happened: the gate-critical-path map predicted post-scoping solo canary at 60–80 s by keying on the 45 s surface straggler; the measured result was 151.5 s because the five surface/artifact-bound fixtures each pay the ~134 s artifact suite at inner width 2. The spec's ≤100 s acceptance then failed and needed a reviewer veto to ship.

Right behavior: a scoping floor estimate should key on the largest bound package's suite time inflated to the inner width, not on the observed straggler of the pre-change mix.

Proposed rule change: none — map-authoring judgment, not a kit rule; fold into the next gate-critical-path re-measure (map #5) when stage 2 lands.
