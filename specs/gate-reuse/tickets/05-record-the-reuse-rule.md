# Record the reuse rule and the merge plan in the ADRs and the docs

Blocked by: 04-plan-a-merge-grade.md
Writes: docs/adr/0031-gate-evidence-crosses-runs-under-one-authorization-key.md (new), docs/adr/0002-accepted-enforcement-postures.md, docs/adr/0016-the-installed-broker-owns-promotion.md, projects/benchkit.md, .bench/BENCH-reference.md, CONTEXT.md, CHANGELOG.md
Covers: GR29, GR30, GR43

## What to build

ADR 31 records the decision: one authorization key names each prospective green, and a green crosses runs and checkouts under that key. It names the key members, the posture binding, and the trust posture of the spec's "Trust posture" section.
It holds no file path and no code snippet, as invariant 3 requires.

ADR 2's posture 5 names the authorization key as the subject binding of reuse. It states that a reuse crosses checkouts, that the posture binds, and that a store record cannot authenticate its writer.
ADR 16's evidence sentence names the authorization key, which binds the prospective tree and the baseline runner identity among its members.

`projects/benchkit.md` updates three paragraphs to the current rule:

- the gate contract sentence on prospective evidence keys
- the worktree merge paragraph, which also states the plan form and its four grades
- the verdict paragraphs on reuse and record classes

`.bench/BENCH-reference.md` updates the gate output sentence on a reused green, and its Command Notes state the merge plan form, its record, and its grades.

`CONTEXT.md` adds the terms `authorization key`, `oracle identity`, `grade`, and `merge plan` with their Avoid lists, and its `worktree merge` entry names the plan form.
`CHANGELOG.md` adds Unreleased entries for the authorization key, the gate-cache write on a reuse, and the merge plan.

Write every sentence in ASD-STE100, so that `bench test --check prose` stays green. No anchor pins a sentence of these paragraphs.
Keep each pinned line byte-identical: `projects/benchkit.md` lines 637 and 715, `.bench/BENCH-reference.md` lines 45, 47, 98, and 752, and `CHANGELOG.md` lines 638 and 990 at the base.

## Acceptance

- [ ] ADR 31 states the key members, the reuse across checkouts, the posture binding, and the trust posture (GR29).
- [ ] ADR 2, ADR 16, the profile, the reference, and the glossary state the authorization key as the reuse rule (GR30).
- [ ] The reference's Command Notes, the profile's merge paragraph, and the glossary state the plan form and its four grades (GR43).
- [ ] `bench test --check prose` and `bench test --check root-conformance` pass on the ticket commit.
