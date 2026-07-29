# Surface pending implementation retros

Blocked by: none

## What to build

Add one bounded scanner for regular `.bench/retros/*.md` capture artifacts.
Include its pending count in `bench status` and `bench roadmap`, and include
stable path/body facts in `bench roadmap --context`, with explicit degraded
evidence for special, dangling, unreadable, or oversized entries.

## Acceptance

- [ ] Story 22's CLI acceptance-coverage row is green.
- [ ] Absent, empty, one-file, and multiple-file directories have deterministic counts and context order.
- [ ] FIFO and oversized retro entries fail closed without blocking or disappearing.
