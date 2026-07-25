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

## 2026-07-25 — a panicking contract test disarms unrelated canary fixtures  [open]

- **What happened:** FT122's first gate went red as
  `canary 'worktree-lifecycle-safety-bypassed' did not bite`. That fixture was
  untouched and correct. The real cause was a new runtime contract test that
  sliced a subject-reported hash to seven bytes without guarding it; against the
  fixture's stub subject the hash was empty, so the test panicked, the whole
  `internal/contract/runtime` binary aborted, and the worktree-safety contract
  never ran to log the substring the fixture greps for. The canary named the
  fixture that went quiet rather than the test that silenced it, and the
  diagnosis cost a full canary sweep at HEAD to establish the red was even ours.
- **Right behavior:** the canary should tell "the inner gate ran and the
  expected substring was absent" apart from "the inner test binary died before
  it could report", and name the panic in the second case. A panic is not a
  check that stopped biting; it is a check that never got to bite, and reporting
  the two identically points every future diagnosis at the wrong file.
- **Proposed rule change:** teach the canary to detect a panic or a non-test
  abort in the inner output and report it as its own failure class, naming the
  panicking test. Separately, consider a conformance check that no test in
  `internal/contract/runtime` slices subject-reported output without a length
  guard — the subject in a canary inner run is a stub, so any such slice is a
  latent tripwire-disarming panic.

## 2026-07-25 — write-delegate shared the main checkout because the work was uncommitted  [open]

- **What happened:** `/bench-implement-spec` requires a write subagent, and
  `craft-delegate` requires an isolated worktree for write-delegations. FT122's
  entire build (~1500 lines) sat uncommitted and largely untracked in the main
  checkout, and the gate was red, so it could not be committed first. A worktree
  branched from HEAD would not have contained the code under repair at all. I
  charged the delegate against the main checkout with an explicit file
  allowlist, no commit authority, and a `git status` check on return.
- **Right behavior:** unclear, which is why this is here. Isolation exists so
  stray edits cannot reach reviewer-owned files and so a mixed `git status`
  cannot make a done-claim unverifiable; with a single writer and a named
  allowlist both risks are bounded, but the skill states the worktree rule
  without an exception for this case.
- **Proposed rule change:** give `craft-delegate` an explicit clause for
  repairing uncommitted work — either sanctioning a shared main checkout when
  there is exactly one writer and the charge carries a file allowlist, or naming
  the route that gets the uncommitted work into a worktree first.

## 2026-07-25 — parked an idea while the gate was running, voiding a green run  [open]

- **What happened:** during FT122's gated commit I answered a reviewer question
  and ran `bench idea` to park the tangent it surfaced. That wrote `IDEAS.md`
  mid-run. Every phase came back green and the commit was still refused with
  `gate subject changed during execution`, costing a full ~15-minute re-run.
  `projects/benchkit.md` states the rule plainly — never mutate the repository
  while a gate is running — and the capture rule in `.bench/BENCH.md` says to
  park a tangent the moment it appears. I followed the second and broke the first.
- **Right behavior:** hold the capture until the gate returns. A parked idea is
  never so urgent that it outranks an in-flight verdict, and the two rules only
  collide because the gate is long enough that answering a question inside its
  window feels free.
- **Proposed rule change:** name the collision where the capture rule lives — the
  Capture section should say the parking is deferred while a gate or gated commit
  is in flight. Better, make it unnecessary: `bench idea` knows a gate is running
  (the subject lock exists), so it could queue the line and write it when the run
  finishes, or refuse with that reason rather than silently voiding the verdict.
