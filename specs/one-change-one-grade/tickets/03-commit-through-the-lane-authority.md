# Commit through the lane authority

Line: opus / medium.

Blocked by: 02-run-a-declared-lane-in-the-gate.md
Writes: internal/commit/commit.go, internal/commit/lane_test.go (new), internal/commit/landing_test.go, internal/landing/landing.go, internal/gate/authorization/authorization.go, internal/gate/authorization/lane_test.go (new), internal/worktree/land_journey_test.go, internal/status/status_test.go

## What to build

A ticket commit in a worktree costs seconds. `bench commit` builds the landing
owner with the lane authority when the root declares a lane, and with the gate
authority otherwise. A lane pass publishes onto the worktree branch. A lane fail
refuses the commit, exits 1, and names the failing check with its first
diagnostic line. `--dry-run` runs the same authority and stops before it
publishes anything.

The authorization result gains two lane kinds, a lane pass and a lane fail. Each
owner construction names the kinds it accepts. The commit owner accepts a green
kind or a lane pass kind. The reviewed landing owner accepts a green kind alone,
so a lane pass never publishes onto `main`. This ticket consumes the lane run,
the lane record, and the built-in lane table that the gate ticket landed.

The lane grades the composed snapshot and not the working tree, so an unnamed
dirty file never fails it. The prose check receives the commit's named Markdown
paths, so a violation elsewhere does not refuse the commit. The existing gofmt
rewrite of the named Go files still runs before the lane, so a rewritten file
passes gofmt. The lane's output states `outcome=pass` or `outcome=fail`, and it
never carries the gate's `green` token.

Reuse the existing commit landing fixture. Give it a phase manifest with a
declared lane and a controllable check, and keep a second fixture that declares
no lane.

## Acceptance

- [ ] OG01 shows `bench commit` with a declared lane prints `lane{outcome=pass,checks=...}` and no `phase test:` line.
- [ ] OG02 shows an unnamed dirty file with a compile error does not fail the lane's build check.
- [ ] OG03 shows the worktree branch ref points at a commit whose tree equals the composed snapshot after a lane pass.
- [ ] OG04 shows a lane check that exits nonzero makes `bench commit` exit 1 and print `lane{outcome=fail,check=<name>}` with the check's first diagnostic line.
- [ ] OG05 shows a lane fail leaves the worktree branch ref unchanged.
- [ ] OG06 shows a named Markdown file with a 27-word sentence fails the lane and names the file and the line.
- [ ] OG07 shows a 27-word sentence in an unnamed Markdown file does not fail the lane.
- [ ] OG09 shows the lane's gofmt check passes on a named Go file that was misformatted before the commit.
- [ ] OG10 shows `bench commit --dry-run` with a declared lane prints the lane record line and no `phase` line.
- [ ] OG11 shows `bench commit --dry-run` after a lane pass leaves the branch ref unchanged.
- [ ] OG12 shows the lane's stdout carries `outcome=pass` or `outcome=fail` and never the token `green`.
- [ ] OG14 shows `bench worktree land` on a source that a lane pass committed prints `gate: green` before it publishes.
- [ ] OG15 shows the reviewed landing owner refuses an authority result of the lane pass kind.
- [ ] OG20 shows `bench status` reports the gate row from the last gate verdict after a lane commit.
- [ ] OG23 shows a fixture repo with a gate script and no lane prints `gate: green` at `bench commit`.
- [ ] OG24 shows a `lane` entry with an empty argv makes `bench commit` exit 1 and name the defect.
