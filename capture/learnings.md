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

## 2026-08-04 — A promotion squash silently reverted a commit made to `main` during the run  [open]

- **What happened:** While auditing ticket quality across two builds, `git diff 950e354 fafb049`
  showed that the `recovery-discard` promotion does not contain `950e354`'s ownership-fence
  corrections: that commit landed directly on `main` while the run was active, and the promoted
  composition — built from the candidate refs, based at the run's own base — reinstated the
  pre-correction text over it. Nothing warned, and the loss surfaced only because a later
  analysis diffed the two commits for an unrelated reason. The affected files have since been
  deleted by spec retirement, so nothing durable is lost here, but the mechanism is live.
- **Right behavior:** A run's composition should not be able to silently overwrite a commit that
  landed on the branch after its base. Either promote refuses when the prospective tree would
  revert a path changed on the branch since the run's base, or it reports those paths and makes
  the reviewer decide. Recomposition onto a moved tip already exists as an operation, so the
  detection has a home; what is missing is noticing the conflict at all.
- **Proposed rule change:** Pair this with the same day's mid-run spec-edit wedge — they are two
  faces of one gap. The run pins a subject and refuses reads that disagree with it (the spec blob),
  while writing through disagreement without a word (a landed commit). One rule covers both: an
  active run states, at write time, which paths it has pinned, and the CLI that would violate the
  pin refuses while the operator is still holding the decision.

## 2026-08-04 — The project-green marker desyncs from a run's base, wedging promotion repo-wide  [open]

- **What happened:** `bench spec build promote` refuses every run on `main` with
  `project-green marker conflicts with another tip`. `refs/bench/green/main` points at
  `572b9046`, the tip the `authoring-hardening` run started from. When a commit landed on
  `main` while that run was still empty, `fastForwardEmptyRun` advanced `run.Base` and — by
  documented design — ran no gate and left the marker alone. `recomposePromotion` then calls
  `authorization.Bootstrap` with `run.Base` as the expected prior marker, finds the older
  value, and refuses. The marker only ever advances through that same call, so no amount of
  `bench gate --fresh` clears it, and `checkpoint`, `assign`, and `review` all refuse once the
  tip has moved (`abandon` is the sole exemption). The run could neither finish nor take
  another ticket. This is not local to one run: the marker last moved under an earlier run,
  `recovery-discard` landed through `bench commit` rather than promotion, so any spec
  attempting promotion on `main` hits the same refusal.
- **Right behavior:** The empty-run fast-forward and the green marker have to move together,
  or `Bootstrap` has to accept a marker that is a recognized ancestor of both the expected base
  and the tip rather than requiring exact equality. The current pair is a state machine with one
  transition that silently invalidates the next one. Either fix restores promotion; the second
  is smaller and matches how `recognizedAdvance` already reasons about a moved tip elsewhere in
  the same file.
- **Proposed rule change:** None for the operator — this is a code defect, not a discipline
  problem. Worth naming that the recovery required abandoning a reviewed spec-backed build and
  landing five completed tickets through the light path with reviewer approval, because the
  sanctioned exit from the wedge discards the composition. A wedge whose only escape hatch is
  destructive should be a release-blocking bug, not a papercut.
