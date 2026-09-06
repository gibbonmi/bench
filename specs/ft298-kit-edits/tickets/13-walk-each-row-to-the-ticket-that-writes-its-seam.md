# Walk each row to the ticket that writes its seam

Blocked by: 12-enumerate-the-exact-match-tests-on-a-changed-message.md
Writes: .agents/skills/bench-craft-spec/references/map-discipline.md
Covers: none

## What to build

Add four rules to `.agents/skills/bench-craft-spec/references/map-discipline.md` under `At ticket slicing`. Write them in ASD-STE100, in the shape of the neighbouring items. Source: roadmap row FT298.

Before the first review charge, the author walks each coverage row's seam path against the ticket `Writes:` lines. The author moves the row to the ticket that writes that path. A golden-file row states its exact predicate, and the old-versus-new diff stays review-owned. A build-time rewrite scope excludes `specs/*/spec.md` and the tickets by name. The spec quotes every anchor needle verbatim.

## Acceptance

- [ ] the rules sit as list items under `## At ticket slicing`.
- [ ] the items name the `Writes:` walk, the golden-file predicate, the rewrite-scope exclusion, and the verbatim needle.
- [ ] `bench gate-prose` passes on the edited file.
