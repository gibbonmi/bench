# Google model scorecard

Last incorporated landing: `map-only-route` build (landing `38224aac1e719dd973daed40fd4d0751550de269`, 2026-10-14).

The `map-only-route` build ran `gemini-4-argon-high` as the author of one guidance ticket, by reviewer direction, as the second trial of the model. The orchestrator ran it through the `agy` CLI in print mode inside the build worktree, with a self-contained charge. The Claude Code permission classifier refused the launch until the reviewer added a Bash rule for `agy`. No Bench hook guarded the session. Opus/high ran the review axes, and Opus/low ran the two repairs.

Provider usage and charges remain unknown, because `agy` print mode reports no token count. The session took about 40 minutes of wall time, and its ticket commit came at about 25 minutes. About 15 minutes went to the floor runs and the eight record calls after the commit.

The first trial on `probe-mutant-sets` ran the same model on one guidance ticket. Its session took 13 minutes, and its ticket commit came at about 10 minutes.

## Current routing

| model / effort | role and sample | observed quality | current use | calibration |
| --- | --- | --- | --- | --- |
| `gemini-4-argon-high` / high | implementer, 2 guidance tickets | On both trials it obeyed every safety rule of the charge, stayed inside its fence, showed the one required red, and passed every check on its first attempt. Each first review found prose findings only: seven on `probe-mutant-sets` and six on `map-only-route`, on terms, clause order, paragraph topic, a copied false claim, and a duplicated literal. | Trial only, by reviewer direction; an exact-spec guidance ticket under a self-contained charge | implementation 0.0 over 10 pairs, 0 abstained |

## Representative evidence

| landing | model / effort / role | what it shows |
| --- | --- | --- |
| `probe-mutant-sets` PM-E | `gemini-4-argon-high` / high / implementer | The author changed 10 paths for 9 rows and named a reason for each of 29 untouched `Writes:` entries. Fresh Opus/high axes met 8 of 9 rows and found no test, anchor, fence, or gate defect. Two Opus/low prose repairs then closed the findings. |
| `map-only-route` MR-A | `gemini-4-argon-high` / high / implementer | The author changed 7 paths for 33 rows and named a reason for each of 20 untouched `Writes:` entries. Fresh Opus/high axes met 32 of 33 rows and found no test, anchor, fence, or gate defect. One Opus/low repair closed the six prose findings and the one anchor gap. |

## Current decisions

- Two comparable guidance trials now exist, and both passed every check first-pass with one prose repair round each. A routing change for guidance tickets is the reviewer's decision.
- Run the model through `agy --model <id> --dangerously-skip-permissions -p "<charge>"`; `-p` takes the prompt as its value and ignores a stdin prompt.
- Add a Bash permission rule for `agy` through `/permissions` before the dispatch, because the auto-mode classifier refuses the launch.
- Give each `agy` author a self-contained charge that spells every safety rule, because no Bench hook guards the session.
- Verify each `agy` return with the fence, `git status`, and an independent coordinator probe before the review, as for every author.
- Do not read an `agy` author's wall time as a model measure, because the floor runs and the record calls take most of it.
