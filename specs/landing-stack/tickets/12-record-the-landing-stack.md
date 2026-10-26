# Record the landing stack decision

Blocked by: 11-resume-a-stack.md
Writes: docs/adr/0032-a-landing-stack-grades-several-sources-once.md (new), CONTEXT.md, .bench/BENCH-reference.md, CHANGELOG.md, .agents/commands/bench-drain.md
Covers: LS83, LS84, LS85

## What to build

This ticket authors guidance, so its author raises the effort to high under the leverage override of `craft-line`. Write each artifact in ASD-STE100 prose for a teammate who walks in with no history. Put no file path and no code snippet in the ADR.

- ADR 0032 records the decided state of the stack:
  - the members that a stack admits and the members that land alone;
  - the stack proof, and the `capture/` paths that the settle policy owns;
  - the whole gate of every stack grade;
  - the bisect and the attribution to the first member that is red on top of a green prefix;
  - the build window, which keeps every rule;
  - the absent persistent queue.
- `CONTEXT.md` adds the terms **landing stack**, **stack member**, **stack prefix**, and **attributed member**, each with an Avoid list. **Landing queue** goes on the Avoid list of **landing stack**, because the stack keeps no state across runs. The **landing lease** entry states that a stack's lease names its first member.
- `.bench/BENCH-reference.md` Command Notes state the stack forms, the stack file, the overrides, the preflight, the members table and the `stack{...}` record, the exits, and the resume. The list of whole-gate cases of the landing plan adds `the landing is a landing stack`.
- `CHANGELOG.md` gains one entry under `[Unreleased]`. The author appends near the top and does not read the whole file.
- `.agents/commands/bench-drain.md` replaces the sentence "Then the coordinator lands the ticket as the light-path row states." When the run holds two or more verified light-path tickets at one time, the coordinator lands them as one landing stack. Each member keeps its tickets-only `--spec`. One ticket still lands as the light-path row states. The edit keeps every anchored needle of that file, and it writes no "commit on green".

## Acceptance

- [ ] An ADR states the members of a stack, its proof, its schedule, its attribution, its window rule, and the absent persistent queue (LS83).
- [ ] The glossary, the reference, and the changelog state the landing stack, and the reference lists the `landing stack` reader-set case (LS84).
- [ ] `/bench-drain` lands two or more verified implement-now light-path tickets as one landing stack (LS85).

## Checkpoint verification

- `bench test --check root-conformance`
- `bench test --check prose`
- `bench test --package ./internal/conformance`
