# Command-first usability

## Problem

Bench is meant to feel lightweight from the reviewer's seat, but the public docs
currently start with provenance, invariants, layout, and CLI setup before they show
the command workflow the reviewer actually uses. That leaks the worker-facing CLI
substrate into onboarding and makes the kit feel heavier than its operating model.

## Solution

Make the reviewer-facing path command-first: the README starts with the phase
commands for each harness, setup is framed as asking the agent to run the setup
phase, and every Bench command has an entry orientation and exit handoff. Keep the
CLI mechanics, but move them behind a clearly labeled worker and maintainer surface
and preserve `.bench/BENCH.md` as the operational reference.

## User stories

1. As a reviewer opening the README cold, I want the first workflow section to show
   the Bench commands I invoke, so I can start without learning the CLI substrate.
2. As a Claude Code reviewer, I want the docs to show `/bench-*` phase commands, so
   I can use the native slash-command surface.
3. As a Codex reviewer, I want the docs to show `$bench-*` adapter skills, so I can
   use the installed Codex command surface without guessing at file paths.
4. As a reviewer setting up a repo, I want to ask the agent to run the setup phase,
   so I do not have to learn `bench link` and `bench init` first.
5. As the worker, I want `/bench-setup-repo` to check whether link/init already ran
   and handle the worker-facing CLI steps when needed, so setup works from the
   command workflow instead of assuming the repo is already wired.
6. As a reviewer entering any Bench command, I want to be told what phase I am in,
   what artifact or state this phase produces, and where it sits in the workflow.
7. As a reviewer leaving any Bench command, I want a concise handoff: what changed,
   current gate or artifact state, and the single recommended next command.
8. As a worker or maintainer, I want exact `bench link`, `bench init`,
   `bench status`, `bench shift`, and `bench gate` mechanics preserved in a labeled
   CLI section, so the operational details remain discoverable.
9. As a kit maintainer, I want the gate to catch structural regressions in the
   command-first docs, so the reviewer-facing surface does not silently drift back
   to CLI-first.

## Implementation decisions

- **README order:** the first H2 is `## Reviewer quick start`. It shows the
  reviewer command path before the invariants, provenance, layout, or CLI install
  details. Keep the philosophy and origin story, but move them below the quick
  start.
- **Harness-aware quick start:** the README names Claude Code as `/bench-*`, Codex
  as `$bench-*`, and other AGENTS.md harnesses as "read the matching command file"
  when they lack a native command surface.
- **Setup path:** reviewer-facing setup says to ask the agent to run
  `/bench-setup-repo` or `$bench-setup-repo`. The setup command preflights whether
  Bench-owned files and `.bench/gate.sh` exist, then runs or reports the required
  worker-facing `bench link` / `bench init` step before continuing the repo-specific
  interview.
- **Worker and maintainer CLI:** exact CLI mechanics live under a labeled README
  section and in `.bench/BENCH.md`. The CLI remains fully documented, but it is not
  presented as the first thing a reviewer operates.
- **Command handoff sections:** every `.agents/commands/*.md` file has literal
  `## Entry orientation` and `## Exit handoff` sections near the top. The sections
  are command-specific, not boilerplate: each states the phase purpose, produced
  artifact or state, workflow position, final state to report, and next command to
  recommend.
- **Conformance shape:** add minimal gate checks for structure only: README's first
  H2 must be `## Reviewer quick start`, and every command file must carry both
  handoff headings. Do not gate exact prose beyond those anchors; prose quality is
  caught by review.

## Testing decisions

- A good deterministic test here checks the command-first structure, not the prose
  style. The gate should prove the reviewer-first section exists in the right
  position and that every command file carries the entry/exit contract.
- The highest seam is `bench gate`, because this is a shipped kit-conformance
  contract across README and command files.
- The setup-command content and CLI placement still need reviewer review: exact
  wording and usefulness are semantic docs quality, not something a shell gate can
  judge without brittle text policing.
- Gate command: `.bench/gate.sh`.

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1, 2, 3 | README starts with a reviewer command quick start before philosophy or CLI mechanics. | `bench gate` README conformance check | `node` heading check failed: expected first README H2 `## Reviewer quick start`, got `## The four invariants`. | Fails if onboarding drifts back to invariant/provenance-first before the command path. |
| 4, 5 | Setup is presented as a setup phase, and `/bench-setup-repo` preflights link/init instead of assuming they already happened. | `.agents/commands/bench-setup-repo.md` plus README reviewer setup text | Not TDD-able without brittle prose assertions; current setup command begins from "`bench link` wired..." which assumes the CLI bootstrap already happened. | Human review catches whether the setup path truly hides CLI mechanics from the reviewer while keeping worker actions explicit. |
| 6, 7 | Every Bench command has command-specific entry orientation and exit handoff sections. | `bench gate` command-file conformance check | Shell section check failed for every `.agents/commands/*.md`: missing `Entry orientation` and `Exit handoff`. | Fails if any command lacks the structural handoff contract the reviewer depends on. |
| 8 | CLI mechanics remain available under a worker/maintainer surface and in `.bench/BENCH.md`. | README + `.bench/BENCH.md` docs review, with existing command-currency gate coverage | Already partially covered by existing gate checks that `.bench/BENCH.md` lists real CLI subcommands; README placement needs semantic review. | Keeps CLI details discoverable without making them the reviewer-first path. |
| 9 | Structural regressions are caught by the oracle. | `bench gate` | The two proposed conformance probes currently fail before implementation. | Proves the new command-first structure has deterministic teeth while avoiding brittle prose checks. |

## Out of scope

- **A new one-shot `bench setup` CLI that combines link/init/interview** — this is a
  separate product surface with its own invocation and harness questions; estimate
  60-120 minutes.
- **Changing `bench shift`, `bench status`, or the autonomous loop behavior** — this
  spec changes how the surfaces are presented, not the loop contract; estimate
  60-120 minutes for any loop behavior change.
- **The visual Bench dashboard idea parked on the roadmap** — separate UI product,
  not part of command-first docs; estimate 2-4 hours for a useful prototype.
