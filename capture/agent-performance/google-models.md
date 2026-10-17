# Google model scorecard

Last incorporated landing: `preflight-ownership-closure` build (landing `53c67bb7792f9732c5c04b95bdbcb3d4c852c19c`, 2026-10-17).

The `preflight-ownership-closure` build ran `gemini-4-argon-medium` at medium effort as the author of two code tickets in OC-C3, by reviewer direction. The orchestrator ran it through `agy` with `--dangerously-skip-permissions` inside the build worktree. Fresh Opus/medium repair sessions then took over both tickets: ticket 06 took one repair round and ticket 07 took three. A planned third ticket moved to a Claude author by reviewer direction before its session started.

Provider usage and charges remain unknown, because `agy` print mode reports no token count. This session holds no claims table and no first-review finding count for the two tickets, so their calibration stays unknown.

The `map-only-route` build ran `gemini-4-argon-high` as the author of one guidance ticket, by reviewer direction, as the second trial of the model. Its session took about 40 minutes of wall time, and about 15 minutes went to the floor runs and the record calls after the commit. The first trial on `probe-mutant-sets` ran the same model on one guidance ticket in 13 minutes.

## Current routing

| model / effort | role and sample | observed quality | current use | calibration |
| --- | --- | --- | --- | --- |
| `gemini-4-argon-high` / high | implementer, 2 guidance tickets | On both trials it obeyed every safety rule of the charge, stayed inside its fence, showed the one required red, and passed every check on its first attempt. Each first review found prose findings only, on terms, clause order, paragraph topic, a copied false claim, and a duplicated literal. | Trial only, by reviewer direction; an exact-spec guidance ticket under a self-contained charge | implementation 0.0 over 10 pairs, 0 abstained |
| `gemini-4-argon-medium` / medium | implementer, 2 code tickets | On `preflight-ownership-closure` both code tickets needed Claude repair sessions after the first review: one round for ticket 06 and three for ticket 07. The ticket 07 rounds include a system-suite red that its floor did not run. | Trial only, by reviewer direction; no code ticket without a reviewer decision | unknown |

## Representative evidence

| landing | model / effort / role | what it shows |
| --- | --- | --- |
| `probe-mutant-sets` PM-E | `gemini-4-argon-high` / high / implementer | The author changed 10 paths for 9 rows and named a reason for each of 29 untouched `Writes:` entries. Fresh Opus/high axes met 8 of 9 rows and found no test, anchor, fence, or gate defect. Two Opus/low prose repairs then closed the findings. |
| `map-only-route` MR-A | `gemini-4-argon-high` / high / implementer | The author changed 7 paths for 33 rows and named a reason for each of 20 untouched `Writes:` entries. Fresh Opus/high axes met 32 of 33 rows and found no test, anchor, fence, or gate defect. One Opus/low repair closed the six prose findings and the one anchor gap. |
| `preflight-ownership-closure` OC-C3 | `gemini-4-argon-medium` / medium / implementer | Two code tickets took four Claude repair rounds between them, and the chunk went past the two-cycle allowance by reviewer decision. |

## Current decisions

- Two comparable guidance trials passed every check first-pass with one prose repair round each. A routing change for guidance tickets is the reviewer's decision.
- Two code tickets at medium effort took four repair rounds between them. Do not route a code ticket to the model without a reviewer decision.
- Run the model through `agy --model <id> --dangerously-skip-permissions -p "<charge>"`; `-p` takes the prompt as its value and ignores a stdin prompt.
- Add a Bash permission rule for `agy` through `/permissions` before the dispatch, because the auto-mode classifier refuses the launch.
- Give each `agy` author a self-contained charge that spells every safety rule and asks for a claims table, because no Bench hook guards the session.
- Verify each `agy` return with the fence, `git status`, and an independent coordinator probe before the review, as for every author.
- Do not read an `agy` author's wall time as a model measure, because the floor runs and the record calls take most of it.
