# Repair kit prose and advertisements

Blocked by: remove-preservation-and-recovery.md

## What to build

Author the removed-verb sweep first — literal `spec build` and
`worktree recovery` tokens over `.bench/`, `.agents/`, `projects/`,
`README.md`, `CHANGELOG.md`, `ROADMAP.md`, and staged specs — observe it red
on the current prose, then repair every reference until it is green and
leave it as a standing conformance check. Repairs are reference-scrubs only
(the doctrine rewrite is Spec C's): `.bench/BENCH.md` CLI inventory and
workflow paragraph (reviewed spec-backed builds land serially
commit-on-green, review before final check), `.bench/BENCH-reference.md`'s
lifecycle lookup table, the three phase commands' lifecycle sections,
`craft-tickets`/`craft-delegate` lifecycle mentions, `projects/benchkit.md`,
`ROADMAP.md`, `README.md` if it matches, one `CHANGELOG.md` removal entry
naming the deleted family, close the learnings `--spec` entry, and update
`capture/session-handoff.md`. Covers RM8, RM12.

## Acceptance

- [ ] The removed-verb sweep exists as a standing check, was observed red
      before the prose repair, and is green after.
- [ ] `CHANGELOG.md` carries the one removal entry.
- [ ] The learnings `--spec` entry is closed; no kit surface recommends a
      removed verb.
- [ ] `go test ./internal/conformance/...` green.
