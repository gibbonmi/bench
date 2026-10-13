# Google model scorecard

Last incorporated landing: `probe-mutant-sets` build (landing `062f9e41cbab5e2a0f51a8a80d5f01ca790680e5`, 2026-10-13).

The `probe-mutant-sets` build ran `gemini-4-argon-high` as the author of one guidance ticket, by reviewer direction, as the first trial of the model. The orchestrator ran it through the `agy` CLI in print mode inside the build worktree, with a self-contained charge. By reviewer decision, the session skipped its permission prompts. No Bench hook guarded the session. Opus/high ran the review axes, and Opus/low ran the two repairs.

Provider usage and charges remain unknown, because `agy` print mode reports no token count. The session took 13.0 minutes of wall time, and its ticket commit came at about 10 minutes.

## Current routing

| model / effort | role and sample | observed quality | current use | calibration |
| --- | --- | --- | --- | --- |
| `gemini-4-argon-high` / high | implementer, 1 guidance ticket | On `probe-mutant-sets` it obeyed every safety rule of the charge, stayed inside its fence, showed the one required red, and passed every check on its first attempt. The first review found seven prose findings: four on terms and duplicated facts, and three on accuracy gaps in new reference text. | Trial only, by reviewer direction; an exact-spec guidance ticket under a self-contained charge | unknown |

## Representative evidence

| landing | model / effort / role | what it shows |
| --- | --- | --- |
| `probe-mutant-sets` PM-E | `gemini-4-argon-high` / high / implementer | The author changed 10 paths for 9 rows and named a reason for each of 29 untouched `Writes:` entries. Fresh Opus/high axes met 8 of 9 rows and found no test, anchor, fence, or gate defect. Two Opus/low prose repairs then closed the findings. |

## Current decisions

- Route a second comparable guidance ticket to `gemini-4-argon-high` before any routing change, because one trial gives no comparison.
- Run the model through `agy --model <id> -p "<charge>"`; `-p` takes the prompt as its value and ignores a stdin prompt.
- Give each `agy` author a self-contained charge that spells every safety rule, because no Bench hook guards the session.
- Verify each `agy` return with the fence, `git status`, and an independent coordinator probe before the review, as for every author.
