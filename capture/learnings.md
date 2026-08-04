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

## 2026-08-04 — A spec edit during an active build wedges the run, and the wedge is invisible until checkpoint  [open]

- **What happened:** Mid-build, a read-only research pass found that the `authoring-hardening`
  spec's story-1 implementation decision claimed `ParseTicket` "already owns both fields" when
  it owns only the fence. The reviewer approved correcting the sentence, so the correction
  landed on `main` as its own gated commit while three assignments were out. The first
  `bench spec build checkpoint` then refused with `staged spec no longer matches recorded
  subject`: `precondition.go` derives `specTip` as `git rev-parse HEAD:<spec path>`, the spec
  blob, and the run pins it at `start`. Nothing warned at edit time or at commit time — the
  refusal arrived one delegate round later, against work that was already finished. The exit
  was to restore the exact recorded blob and re-land the correction after promotion.
- **Right behavior:** Unclear which end should move, and that is the finding. Either a spec is
  explicitly frozen for a run's duration — in which case `bench commit` should refuse a write
  to the pinned spec of an active run, the way it already refuses dirt outside its named paths
  — or a run should accept a spec revision the reviewer approved, with an explicit
  re-pin operation. What is certainly wrong is the current shape, where the edit succeeds, the
  gate goes green, and the lifecycle discovers the drift at the next checkpoint.
- **Proposed rule change:** Reviewer decides between the freeze and the re-pin. Worth noting
  that `abandon`'s apply path refuses on the same identity drift (`TestApplyAbandonStillRefuses
  IdentityDriftOnMovedTip`), so a run wedged this way cannot be cleaned up either until the
  blob is restored — the cheap fix is at least a refusal at write time, where the operator is
  still holding the decision.
