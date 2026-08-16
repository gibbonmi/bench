# Collapse the review loop and narrow --reviewer

Blocked by: move-the-spec-template-into-craft-spec.md
Writes: .agents/commands/bench-write-spec.md, internal/anchors/registry_data.go, projects/benchkit.md, tests/canary/workflow-guidance-anchors

## What to build

Spec authoring runs one review round instead of four. `bench-write-spec.md`
states a single round over the spec-and-tickets pair, with `craft-tickets`'
granularity / blocking-edges / merge-split quiz as that round's approval step,
replacing two rounds per artifact. `projects/benchkit.md`'s loop-1 and loop-2
routing rows collapse into one row naming the same mid model and high effort they
already name — the tier and venue are confirmed, not changed.

`--reviewer <tier> [effort]` resolves same-family through the invoking harness's
own `.bench/lines.env` column only. Both the cross-family branch and the bound
model-id form go, and the unbound-id refusal rules go with them — a model id is
now an invocation error, so naming the top tier is `--reviewer top high`.
`references/cross-harness-reviewers.md` itself stays for other phases and is not
edited.

This ticket closes the needle inventory the split opened. Every one of the twenty
needles anchored to `bench-write-spec.md` before the shrink must end in exactly
one state — retained in the thinned file, retargeted by the previous ticket, or
retired here with its canary — and this ticket records the full enumeration.

## Acceptance

- [ ] `(covers SR19)` `bench-write-spec.md` states one review round over the pair with the ticket
      quiz as its approval step, and no two-round sentence remains in the tree.
- [ ] `(covers SR19)` The profile carries one collapsed routing row naming mid model and high
      effort; no loop-1/loop-2 pair remains.
- [ ] `(covers SR20)` The `--reviewer` anchor names a tier-only grammar resolving same-family
      through the harness's own column, with no cross-family route and no
      model-id form.
- [ ] `(covers SR20)` `--reviewer <a bound model id>` is an invocation error.
- [ ] `references/cross-harness-reviewers.md` and its two canaries are unchanged.
- [ ] `(covers SR25)` All twenty pre-shrink needles are enumerated with their end state, and
      every retired needle's canary is retired with it.
- [ ] `(covers SR26)` `bench anchors` is green, and every touched canary still reproduces its own
      `EXPECT` line.
