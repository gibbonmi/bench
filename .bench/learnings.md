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

## 2026-07-27 — I edited a reviewer-approved spec without asking  [open]

- **What happened:** The FT91 build's semantic review found four statements in
  `specs/ft91-artifact-build-tiering.md` that were factually wrong about the code —
  a coverage row describing a test seam that parks somewhere it does not, a story
  over-stating a deletion, and two edge-inventory exclusions resting on reasons the
  build invalidated. The reviewer was unreachable and had given a standing "fix a
  finding now instead of parking it". I corrected all four in the spec, each flagged
  in-line as `**Post-approval correction, flagged:**` for veto, and said so in the
  close.
- **Right behavior:** Unclear, which is why this is here. `.bench/BENCH.md` says
  spec content is the reviewer's call and a batch approval leaves each spec in
  `specs/` as post-hoc veto surface — which is what I relied on. But a
  reviewer-approved spec is a stronger artifact than an unreviewed one, and
  "correct a false citation" and "change what the spec asks for" are not obviously
  the same permission. The precedent commit `5a48cae` ("spec: correct the FT91
  out-of-scope residual with measured evidence") suggests correcting a staged spec
  against measured evidence is established here; whether that extends past reviewer
  approval is what I could not resolve.
- **Proposed rule change:** Let `/bench-review-implementation` state the disposition
  explicitly when a Spec-axis finding lands on the spec rather than the code. Two
  candidate rules, reviewer picks: either (a) a factual correction — a citation that
  resolves to nothing, or a described mechanism the tree contradicts — may be made
  post-approval under the existing in-line veto-flag convention, while anything that
  changes what gets built stops for sign-off; or (b) all post-approval spec edits
  stop, and the review persists them to `reviews/<spec-slug>.md` instead. The
  distinction matters because the Spec axis is charged to audit every coverage row,
  so it will keep producing this class of finding.
