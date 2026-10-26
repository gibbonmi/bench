# Record the grade key rule in an ADR and the current docs

Blocked by: 04-reuse-a-green-in-every-whole-gate-verb.md, 05-run-the-complete-checkpoint-under-the-broker.md
Writes: docs/adr/0031-one-grade-key-names-each-whole-gate-green.md (new), docs/adr/0016-the-installed-broker-owns-promotion.md, docs/adr/0002-accepted-enforcement-postures.md, projects/benchkit.md, .bench/BENCH-reference.md, CONTEXT.md, CHANGELOG.md
Covers: GR29, GR30

## What to build

ADR 31 records the decided state. One grade key names each whole-gate green by its verdict inputs, and the ADR names those members.
A green crosses checkouts and verbs under that key, the posture binds every reuse, and the kit's complete checkpoint runs under the landing broker.

ADR 16's evidence sentence names the grade key in place of the baseline runner identity. ADR 2's posture 5 names the grade key.
The profile and the reference state the reuse rule and the broker checkpoint.
`CONTEXT.md` adds the **grade key** term with its Avoid list: authorization key, evidence key, reuse key.
`CHANGELOG.md` records the user-facing changes under one Unreleased heading.

## Acceptance

- [ ] ADR 31 states the key members, reuse across checkouts and verbs, the posture binding, and the broker checkpoint, with no file path and no code.
- [ ] ADR 16, ADR 2, the profile, the reference, and the glossary state the grade key as the reuse rule, and no sentence says that a green stays in one checkout.
- [ ] The changelog entry names the reuse across checkouts, the broker checkpoint, and the new refusal face.
