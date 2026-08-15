## Outcome

FT189 landed as `c0e62b4a5b4eb485c2095668e8cd3c2123bb3bf3` from reviewed source pair `ea8b3aa248910618647453e1f5456e3ac47158a4` → `c4df72c76dcfb25e045c4edf488cd3d9350be9e3`. It fails closed before `git worktree list --porcelain` can block on malformed private admin entries, bounds the enumeration as a backstop, and carries the typed refusal through list, status, resume, dashboard, and doctor. The retained worktree was removed after its 39 ignored gate logs were discarded through the fingerprinted cleanup plan. Retirement is `00290cf2`.

## Gate-stage timings

Retained landing gate `20260815T152136.407646804Z-1503429` was green in 121.385s: binary selection 2.562s; gofmt 0.122s; vet 3.468s; test 100.040s; race 7.770s; system 4.406s; shellcheck 0.515s.

## Ticket-versus-spec-slice and delegate performance

The four serial tickets kept the implementation independently gateable: resolve common dir, refuse malformed entries, bound enumeration, then report the refusal in doctor. The staged spec records 34 spec and 3 ticket verification iterations. Per-delegate elapsed timings were not retained; the landing gate is the available terminal performance evidence.

## Coordinator catches

- Exact-tip review caught the missing `Writes:` ownership for the malformed-entry ticket and required the approved fence amendment before landing.
- Landing recovery caught two lifecycle contracts: callers must supply the original request token rather than its persisted digest, and source and destination must carry identical staged-spec bytes.
- The cleanup receipt prevented deletion of retained ignored logs until its target and fingerprint were reviewed.

## Repair attribution

| ticket | repair rounds | cause per round |
| --- | ---: | --- |
| resolve-git-common-dir | 3 | delegate-error; delegate-error; delegate-error |
| refuse-malformed-admin-entries | 5 | delegate-error; delegate-error; delegate-error; delegate-error; spec-row |
| bound-worktree-enumeration | 0 | none |
| report-admin-entry-in-doctor | 1 | delegate-error |

## Agent-experience improvements

### Bench CLI

`bench worktree land --request` advertises an opaque identifier while lifecycle records expose only its digest; reusing that digest causes a second hash and refusal. Accept the listed assignment/digest selector or emit a resumable request token so crash recovery can replay landing without historical context.

### Skills

The implementation/review handoff should name the original landing request token beside the frozen base and source tip. That makes the exact post-review landing invocation recoverable without inspecting lock metadata.

### Process

Keep authored maps and specs in assigned worktrees so in-progress authoring cannot dirty `main` during an unrelated landing. The resulting coordination rule is parked in `capture/IDEAS.md` for the next drain.
