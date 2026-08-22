# Spec and ticket prose: ASD-STE100

Charged from `craft-spec` and `craft-tickets`. Write every sentence in a
`specs/<slug>/spec.md` and in a ticket file in ASD-STE100 Simplified Technical
English. The rules govern the sentences; the template's headings, labels,
identifiers, paths, commands, and table cells stay as the template fixes them.

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

Why: a coverage row or an acceptance row written as one short active sentence
is one predicate, and a ticket lifts it unchanged. A 2026-08-22 two-arm
comparison of the same decision map found the STE spec shorter, more exact
in its rows, and the easier one to slice.
