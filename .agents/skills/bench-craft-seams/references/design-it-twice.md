# Design it twice — the briefs and the comparison

`craft-seams` routes here when a seam is genuinely uncertain and carries a
high-effort line. Spawn the sub-agents in parallel, one brief each. Paste the
brief verbatim, substituting `<module>`:

- Agent 1 — minimal: "Design the smallest possible interface for `<module>`:
  1–3 entry points, maximum leverage each. Return the interface (signatures,
  invariants, ordering, error modes), one usage example, what it hides, and
  where its leverage is thin."
- Agent 2 — flexible: "Design the most extensible interface for `<module>`:
  cover the known use cases and leave room for the unknown ones. Return the
  interface (signatures, invariants, ordering, error modes), one usage
  example, what it hides, and where its leverage is thin."
- Agent 3 — common-caller: "Design the interface for `<module>` that makes the
  most common call trivial: the default case should need near-zero
  configuration. Return the interface (signatures, invariants, ordering, error
  modes), one usage example, what it hides, and where its leverage is thin."
- Agent 4 — ports and adapters (only when the seam crosses a real dependency):
  "Design the interface for `<module>` as a port isolating `<dependency>`,
  with adapters behind it. Return the interface (signatures, invariants,
  ordering, error modes), one usage example, what it hides, and where its
  leverage is thin."

Then compare the returns on the three things that matter — **depth** (leverage
at the interface), **locality** (where change concentrates), and **seam
placement** — and give an opinionated recommendation, or a hybrid if elements
combine well. A menu isn't the deliverable; a strong read is.
