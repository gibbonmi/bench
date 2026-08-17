# Collapse the review loop and narrow --reviewer

Blocked by: move-the-spec-template-into-craft-spec.md
Writes: .agents/commands/bench-write-spec.md, .agents/skills/bench-craft-spec/SKILL.md, internal/anchors/registry_data.go, projects/benchkit.md, tests/canary/workflow-guidance-anchors

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

Three sentences are two-loop residue and are reworded to one round, not merely
left standing: the slicing step ("After loop 1 accepts … then run loop 2 on that
breakdown"), the learnings hook ("When *either loop* takes more than one
iteration"), and the `Verification log: spec <n> + tickets <m> iteration(s)`
two-count schema, whose two counts *are* the two loops encoded into the artifact.
Each needle is reworded together with its registry row and its canary, so the
anchor keeps guarding a true claim rather than a sentence describing machinery
this change deletes.

The review rubric moves to `craft-spec` rather than being retired: the materiality
exit ("A finding blocks only when it changes observable behavior, an ownership
fence, or the ticket graph…"), the cheapest-plausible degenerate standard, and the
falsification questions are rubric, not phase choreography, and `craft-spec`
already owns the process. Retarget their registry rows with the receiving section
named; do not delete them.

This ticket closes the needle inventory the split opened. The file carries **47**
anchored registry rows before the shrink — 34 `Require`, 8 `RequireInSection`, 4
`Forbid`, 1 `ForbidInSection`, all live and green; the "twenty" the spec carried at
staging was a miscount. Every one must end in exactly one state — retained in the
thinned file, retargeted by the previous ticket or by this one, or reworded here
with its canary — and this ticket records the full enumeration. A `BASE`/`MUTATE.json`
fixture that mutates a moved needle follows it to the receiving file; payload-only
fixtures are hand-written stand-ins, not copies of `bench-write-spec.md`.

This ticket also carries the budget row the previous one deferred. Only after the
rubric moves and the loop prose collapses does `bench-write-spec.md` fit its budget,
so `projects/benchkit.md` gains its `.agents/commands/bench-write-spec.md | 73` row
(reviewer-accepted at build time in place of 60, which the retained needles cannot meet)
here, in the same green commit that makes the file fit.

## Acceptance

- [ ] `(covers SR18)` The profile's budget table carries the
      `.agents/commands/bench-write-spec.md | 73` row, the file is at most 73 lines, and
      raising it above 73 turns the budget check red.
- [ ] `(covers SR19)` `bench-write-spec.md` states one review round over the pair with the ticket
      quiz as its approval step, and no two-round sentence remains in the tree.
- [ ] `(covers SR19)` The profile carries one collapsed routing row naming mid model and high
      effort; no loop-1/loop-2 pair remains.
- [ ] `(covers SR20)` The `--reviewer` anchor names a tier-only grammar resolving same-family
      through the harness's own column, with no cross-family route and no
      model-id form.
- [ ] `(covers SR20)` `--reviewer <a bound model id>` is an invocation error.
- [ ] `references/cross-harness-reviewers.md` and its two canaries are unchanged.
- [ ] `(covers SR25)` All 47 pre-shrink registry rows are enumerated with their end
      state, and every retired or reworded needle's canary moves with it.
- [ ] `(covers SR19)` No sentence in the tree names `loop 1`, `loop 2`, or `either loop`,
      and the verification log carries one count for the single round over the pair.
- [ ] `(covers SR25)` The materiality exit, the degenerate standard, and the falsification
      questions read in `craft-spec`, with their registry rows naming its section.
- [ ] `(covers SR26)` `bench anchors` is green, and every touched canary still reproduces its own
      `EXPECT` line.
