# Written-artifact prose: ASD-STE100

The one source for the prose standard `.bench/BENCH.md` names. It is charged
from `craft-spec`, `craft-tickets`, `craft-adr`, `craft-comments`, and
`craft-skills`. It governs every written artifact: a doc, an ADR, a roadmap
row, a handoff, a retro. It also governs a journal entry, a spec, a ticket, a
skill, and a comment. Write every sentence in ASD-STE100 Simplified
Technical English.

The rules govern the sentences. A template's headings, labels, identifiers,
paths, commands, and table cells stay as the template fixes them. A code
comment obeys these rules too, inside the register `craft-comments` owns.

- Use the active voice and name the agent: "The gate reds the row", not "The
  row is rejected".
- Keep one topic per paragraph and six sentences or fewer per paragraph.
- Keep a procedural sentence to 20 words or fewer and a descriptive sentence
  to 25 words or fewer.
- Give one instruction per sentence, in the imperative, with the condition
  before the instruction.
- Use the present tense for a description and the past tense only for an
  event that occurred before.
- Write the articles ("the", "a"); do not use telegraphic style.
- Do not use a gerund as a noun, and do not use a noun as a verb.
- Use one word for one thing; do not use a synonym for variety.
- Do not make a noun cluster of more than three nouns.
- Use a vertical list for three or more parallel items.

A short label with a colon and no sentence terminator is a field line; the
prose mechanics check does not grade it.
A terminated label is a field line only when it names `Blocked by`, `Covers`, `Drift`, `Occurrence`, `Occurrences`, `Source`, `Sources`, `Supports`, or `Writes`.
Keep a grammar field, a marker phrase, and an anchor needle on one physical line.

Why: a coverage row or an acceptance row written as one short active sentence
is one predicate, and a ticket lifts it unchanged. A 2026-08-22 two-arm
comparison of the same decision map found the STE spec shorter, more exact
in its rows, and the easier one to slice.
