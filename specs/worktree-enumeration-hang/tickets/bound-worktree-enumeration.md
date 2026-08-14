# Bound the worktree enumeration's git children

Blocked by: resolve-git-common-dir.md, refuse-malformed-admin-entries.md
Writes: internal/git, internal/gittest, internal/bounds, internal/conformance/bounds_policy_test.go, internal/status/status_test.go, internal/worktree/list_actions_test.go

## What to build

With clean-shaped admin entries and a git that blocks anyway, `git.Worktrees`
returns a typed timeout refusal within the named `internal/bounds` deadline
per launched child, per the spec's "The bound" and "Conformance binding"
decisions: the `WorktreeListTimeout = 15 * time.Second` registry constant;
the package var consumed from it, with the exported set-and-restore test
hook (the one cross-package knob, never called by production); the
stream-preserving five-status `internal/bounds` variant (caller stdout
writer, `Setpgid` + group kill retained); both children driven through it
via the predecessors' argv helper, with start and canceled mapped into the
typed class and the expiry message rendering the var's live value, never
the constant; and the `bounds-policy` `required` and `owners` tripwire rows
proven to bite before the consumption lands. One comment at
`refCheckTimeout`'s declaration distinguishes the hook-scoped fail-safe
from this policy bound, and the registry constant carries the git 2.43.0
retirement comment. The stub gains the block-rev-parse, noisy-list, and
vanish-after-rev-parse modes. The named mutation probe is two single-site
`bounds.Run` substitutions, one per child. Tests reuse `git_test.go`'s
`runGit`/`newRepo` — the census caps `internal/git` test constructors at
one. "Within the bounded wait" here derives from the test-overridden
package var (goroutine plus deadline), not the `TestDeadline(0)` floor the
scanner rows use.

## Acceptance

- [ ] The `bounds-policy` check reds with the owner row present and the consumption absent, and greens once `internal/git/git.go` consumes the token (covers WE8)
- [ ] A blocking porcelain child with a pgid-shared descendant holding stdout yields, within the bounded wait, a typed error naming the invocation and the enforced (overridden) bound value, not the constant's 15s (covers WE9)
- [ ] The package var's default equals `bounds.WorktreeListTimeout` (covers WE14)
- [ ] Stderr noise on both children leaves the parsed worktree list correct; each single-site `bounds.Run` substitution turns its probe row red (covers WE15)
- [ ] A blocking rev-parse child with a stdout-holding descendant yields a typed timeout refusal naming the rev-parse invocation and the enforced bound value (covers WE17)
- [ ] A bound expiry renders at `appendWorktree` with the investigate action and no re-run advice (covers WE21)
- [ ] A bound expiry renders at `bench worktree list` naming the timed-out invocation, action `investigate the git failure`, not `inspect and remove it` or the retry advice (covers WE22)
- [ ] A git that cannot start yields the typed error with the exec failure text, the failing invocation's name per fixture, and the investigate action — never a silent empty list or a borrowed timeout message (covers WE23)
