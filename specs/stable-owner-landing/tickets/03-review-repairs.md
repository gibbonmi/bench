# Review repairs

Blocked by: 02-gate-and-publish-under-the-stable-owner.md
Writes: bin/bench.sh, internal/adopt/broker.go, internal/adopt/broker_test.go, internal/worktree/land.go, internal/worktree/land_identity.go, internal/worktree/land_identity_test.go, internal/gate/prospective.go, internal/gate/prospective_owner_test.go, internal/systemtest, CHANGELOG.md, specs/stable-owner-landing/spec.md

## What to build

Close the accepted review findings in `reviews/stable-owner-landing.md`.
The reviews file names each target and its citation; it is the source.

## Acceptance

- [ ] S1: the CHANGELOG entry obeys the STE sentence limits.
- [ ] S2: one source supplies the broker platform fact; the wrapper stops
      deriving it for the broker check, and the digest comparison remains.
- [ ] P1: the landing binds `--base` to the assignment's recorded start,
      keeps the ancestry guard, and a test pins the exact refusal for a
      different valid ancestor.
- [ ] C1: each fail-closed manifest and broker refusal branch has a test.
- [ ] C2: each routing variable refuses when set but empty, under test.
- [ ] C5: a prospective build failure returns a clean error, under test.
- [ ] C6: the spec's Won't-handle list disposes of newline and control-byte
      roots.
