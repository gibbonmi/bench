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

## 2026-08-03 — Started a code fix in the main checkout, then could not undo it  [open]

- **What happened:** A `/bench-debug` session began editing five files directly in the
  main working tree. A second session's `/bench-what-next` drain was already dirty in
  that same tree, and `bench commit` refuses on any dirt outside its named paths, so
  both changesets deadlocked: neither could commit past the other. Moving to a
  `bench worktree` afterwards did not resolve it, because clearing the duplicate from
  main needs a `git restore`, which `block-dangerous-git` correctly blocks. The reviewer
  had to run the restore by hand, twice — the first invocation wrapped across two lines
  and silently restored only three of the five paths.
- **Right behavior:** Take the worktree *before* the first edit, not after the collision.
  Invariant 1 already says to use one when `git status` shows another writer, but the
  status check has to happen at the moment work starts; here the tree was clean at
  session start and dirtied by the other session mid-run, so a single entry check was
  not enough. A debug session that expects to write code should open its worktree at
  entry regardless of what the tree looks like then.
- **Proposed rule change:** `/bench-debug`'s Phase 5 says to route code authorship
  through `craft-delegate`, which owns worktree isolation — but Phases 1 and 2 routinely
  write throwaway harnesses and repro tests before Phase 5 is reached, with no isolation
  guidance. Either the entry orientation names the worktree as the debug session's
  workspace, or the fix-at-a-seam phase's isolation rule moves earlier. Reviewer decides
  whether this is `/bench-debug`'s or a general session-entry rule.

## 2026-08-03 — `bench worktree release` grew the orphan count while retiring a worktree  [open]

- **What happened:** A safety worktree was released after its content had already been
  committed to main. Release still marked the row `recovered` and preserved a recovery
  ref, taking open assignments 49→51 and recovery refs 58→59, because the checkout was
  dirty and proof-of-landing is OID identity — a dirty tree's payload never matches a
  commit, however equivalent its content. `bench worktree recovery --apply` then refuses
  that row permanently ("every recovery payload must be proven landed").
- **Right behavior:** Unclear, and that is the finding — the retire half of the
  "deliberate recover-or-retire" that `internal/worktree/resume.go` names in prose does
  not exist as a command, so the count is monotonic by construction.
- **Proposed rule change:** none from me; parked to `capture/IDEAS.md` with the
  measurements, because discarding preserved work is a reviewer decision and a new seam.

## 2026-08-03 — A ticket's `Assumptions:` field mostly holds verified preconditions, not assumptions  [open]

- **What happened:** Writing eight tickets for `specs/recovery-discard`, every
  `Assumptions:` line I authored turned out to hold one of two things, and neither is an
  assumption. Most clauses were inherited preconditions that are checkable against the
  tree ("`enrich-recovery-plan.md` has landed, so orphaned and absent are distinct plan
  verdicts") — a delegate can verify that in one command, and the ticket's own
  `Blocked by:` line already asserts it. The rest were restatements of a standing kit
  rule: all eight tickets ended with a clause meaning "claims re-derived from the tree at
  pickup", which is `craft-tickets`' instruction to every ticket, copied eight times into
  the artifacts that instruction governs. The reviewer noticed the tension with
  `.bench/BENCH.md`'s "NEVER assume, always verify" and asked whether it was the cause of
  a separate parsing defect. It was not — that defect was `listValue` splitting on commas
  — but the naming question stands on its own.
- **Right behavior:** Unclear, and that is the finding. A field named `Assumptions` in a
  kit whose operating guide forbids assuming invites authors to write verifiable facts
  into a slot whose name says they were not verified. Two clause kinds are being pooled:
  inherited preconditions (verifiable, and partly redundant with `Blocked by:`) and
  standing-rule restatements (a one-source-per-fact violation, since the rule is
  canonical in the skill). Only a genuine third kind — something the ticket takes on
  faith because it *cannot* be checked at authoring time — matches the field's name.
- **Proposed rule change:** Reviewer decides between three shapes, and I have no
  preference strong enough to recommend one: rename the field to something like
  `Preconditions:` and let it hold checkable facts honestly; keep the name and have
  `craft-tickets` say explicitly that only unverifiable-at-authoring-time claims belong
  there, with inherited preconditions left to `Blocked by:`; or split it into two fields.
  Any of the three is a change to the ticket grammar, so it touches `craft-tickets`, the
  taught example the `example-agreement` conformance check grades, and `ParseTicket` in
  `internal/specbuild/assign.go` — kit surface, under the synthesis discipline, not
  something I should act on from inside a build.

## 2026-08-04 — An authoritative review's findings survived only in conversation, and were lost  [open]

- **What happened:** The `recovery-discard` build's second review ran in the Codex CLI
  under the reviewer's top binding, returned eight findings including two criticals, and
  was never written to the tree. The findings reached the next session only as a summary
  the reviewer pasted by hand — finding tags, severities, a fence table, and prose for the
  two the summarizing session happened to describe in full. When that context was cleared
  the underlying output was gone, and the reviewer confirmed they no longer had it. Four
  of the eight findings (C4, C5, C6, SP5) exist now as tags with no text. The repair round
  proceeded by reproducing the two criticals independently and reading the other two
  defects out of the tree, so every repair ticket rests on evidence gathered fresh — but
  no one can say whether those tickets close the findings that were actually raised.
- **Right behavior:** A review that authorizes or blocks a promotion is durable evidence,
  and evidence that lives only in a chat transcript is not evidence the next session can
  reach. The review's output should have landed in the tree — the phase already has the
  shape for it, since `/bench-implement-spec` tells the build to read `reviews/<spec-slug>.md`
  and the promoted composition deletes it, so resolved findings cannot resurface as pickup
  work. That file was never written here. The delegated-review path writes a bounded
  receipt for `bench spec build review`, which records that a review happened and its
  verdict; it does not preserve the findings themselves, and those are what a repair round
  actually consumes.
- **Proposed rule change:** Make writing `reviews/<spec-slug>.md` a required step of
  `/bench-review-implementation` rather than an artifact the implement phase merely reads
  if present, and say so in both commands — including for a review delegated to another
  harness, where the coordinator owns capturing the returned findings to that path before
  acting on them. Reviewer decides; this touches two canonical command surfaces and the
  spec-build lifecycle's review receipt, so it is kit surface under the synthesis
  discipline, not something to act on from inside a build.

## 2026-08-04 — Coordinator checkpoint receipts are hand-assembled, and the ticket fence is only checked at checkpoint  [open]

- **What happened:** Landing four repair tickets through `bench spec build` cost far more
  coordinator effort than the repairs themselves. Two causes, both mechanical. First, the
  checkpoint receipt has no producer: the coordinator hand-writes JSON carrying the run and
  assignment identities, the assignment base, a tree hash computed the way `git.TreeHash`
  does it, the ticket digest, one row per acceptance ID, the changed-path set diffed against
  the assignment base, and a sha256 of every assumption string — each cross-checked by
  `validateReceipt`, so any one of them wrong is an opaque `invalid spec build receipt`.
  Assembling two receipts took more coordinator turns than verifying the two critical
  repairs they attested. Second, a ticket's ownership fence is validated only when a
  receipt arrives: one ticket asked for a change whose only possible site lay outside its
  own fence, and that contradiction surfaced after a delegate had already done the work,
  not at slicing time. There is no way to release a single open assignment to re-fence it —
  `abandon` is whole-run — so the correction had to be a sibling ticket and a further
  delegate cycle.
- **Right behavior:** The receipt is derived data. Every field is computable from the run
  record, the assignment worktree, and the coordinator's own probe result, so the
  coordinator should supply only what it alone knows — the row outcomes, the check list,
  and the probe command, output, and exit — and have the tool derive the rest. And a fence
  that cannot hold every path its own acceptance rows require is a defect the ticket parser
  could catch when the ticket is read, long before a delegate is charged against it.
- **Proposed rule change:** Two, both kit surface for the reviewer. Add a receipt producer
  — `bench spec build receipt <slug> --assignment <id>` emitting the derived skeleton for
  the coordinator to complete — so the hand-assembly and its opaque refusal disappear.
  Separately, consider whether `assign` should refuse a ticket whose acceptance rows name a
  path outside its declared fence, or whether that stays an authoring discipline
  `craft-tickets` teaches; the second is cheaper but did not prevent this.
