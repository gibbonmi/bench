# Document anchors use a declarative substring registry

Document guidance anchors are ordered declarative data, evaluated by one shared matcher, with section-scoped placement and HTML-comment exclusion where required. The agent-facing query reads that registry rather than re-parsing enforcement code. Substring matching is intentionally retained, because it gives deterministic, byte-stable diagnostics and cheap omission or relocation oracles. Paraphrased semantic restatements remain review-owned until evidence justifies a stronger mechanism.

## Admission

An anchor earns a place in the registry only in one of five classes:

- FORBID: the anchor forbids a stale or wrong form.
- LABEL: the anchor names a label or a heading that a parser or a template reads.
- ROUTE: the anchor names a phase, a skill, or a verb as a routing reference.
- INVARIANT: the anchor carries a platform invariant from the core sections of the operating guide.
- COMMAND: the anchor pins a Bench command line with an argument or a flag.

COMMAND is an interim class. A grammar check of command spans against the executable help replaces it.

A Forbid anchor takes the FORBID class from its kind, and each Require anchor declares one of the other classes. The registry test refuses a Require anchor with no class or with a shape outside its class.

The wording of a rule sentence is review-owned, and no anchor pins it. Review grades that wording directly. So a reworded rule changes one file, and the gate stays green.

Canary coverage stays at one planted needle for each check family.
