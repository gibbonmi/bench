# Introduce the typed action carrier

Blocked by: none
Ownership fence: `internal/axi`
Integration surfaces: action API→`internal/axi`; registry declarations→declare-production-axi-registry.md
Contracts: ordered token vector, fixed argument values, open placeholder names, and executable disposition cross caller→`internal/axi`, domain is command template or non-invokable prose, order is literal argv order, and absence is zero actions, asserted by AC1
Closure: AC1/tokens, AC1/fixed, AC1/open, AC1/prose, AC1/absence

## What to build

typed actions preserve tokens, fixed arguments, open placeholders, prose disposition, and zero-action absence without rendering help.

## Acceptance

- [ ] [AC1] (covers CR2) typed actions preserve tokens, fixed arguments, open placeholders, prose disposition, and zero-action absence without rendering help.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| AC1/tokens | flatten tokens into one prose string | action test | construct and require ordered tokens |
| AC1/fixed | turn a fixed argument into an open placeholder | action test | construct and require the fixed value |
| AC1/open | replace an open placeholder with an empty fixed value | action test | construct and require it remain open |
| AC1/prose | mark orchestration prose executable | action validation test | construct and require refusal |
| AC1/absence | invent a default action for an empty list | action test | construct zero actions and require zero |

