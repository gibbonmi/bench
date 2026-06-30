# roadmap — capture an idea without committing to it

> Canonical term is **roadmap** (see `CONTEXT.md`). "Icebox" is a near-synonym we do
> not adopt in code or docs.

> Feature **B** of the ambient-feedback map (`decisions/ambient-feedback.md`, #7–#9).
> Decided: two specs, **B first**. Feature **A** (the surface) is a separate spec; its
> zero-severity footer that *displays* this roadmap is A's wire-up, not B's.

## Problem

Mid-build, I get an idea that's out of scope for where I am. Today the kit has nowhere
to put it that doesn't commit me to working it: a decision map implies intent to
resolve, a spec is committed work, an issue is a backlog I'll work. So the idea either
derails the current build or gets lost. I want to park it in one breath and move on.

## Solution

A capture-and-forget sink. `bench idea "<text>"` appends the idea to a plain
`ROADMAP.md` at the repo root and exits — no prompts, no workflow, no spec. `bench
roadmap` prints the parked list when I choose to look. The file is a dumb append-only
sink: no status, no lifecycle. When I later decide to act on a parked idea,
`/bench-ideate` (invoked cold) offers my parked items and pulls the chosen one into a
new decision map, removing its line.

## User stories

1. As a builder mid-build, I want to capture an out-of-scope idea with one command so I
   neither lose it nor derail my current work.
2. As a builder, I want capture to commit me to nothing — no prompt, no grill, no spec;
   `bench idea "<text>"` appends and exits.
3. As a builder, I want each entry timestamped (`- YYYY-MM-DD  <text>`) so I can see how
   long an idea has been parked.
4. As a builder, I want `bench roadmap` to print every parked idea on demand.
5. As a builder, `bench idea` with no text gives a usage error and appends nothing (no
   blank entry).
6. As a builder, the first `bench idea` creates `ROADMAP.md` if it doesn't exist.
7. As a builder, `bench roadmap` on an empty or absent roadmap tells me it's empty
   rather than erroring.
8. As a builder, idea text with spaces/punctuation is captured literally (all words
   after the subcommand are the idea).
9. As a builder running either subcommand outside a git repo, I get the same
   not-in-a-repo error as every other `bench` subcommand.
10. As a builder, `bench` help lists `bench idea` and `bench roadmap`.
11. As a builder running `/bench-ideate` cold (no idea already in hand), I'm shown my
    parked ideas and asked which, if any, to pull up.
12. As a builder, when `/bench-ideate` creates a map from a pulled idea, that idea's
    line is auto-removed from `ROADMAP.md`; an abandoned pull that writes no map keeps
    its line.
13. As a builder running `/bench-ideate` already carrying a fresh idea, I'm not
    interrupted with the roadmap prompt.
14. As a builder, the roadmap carries no status or lifecycle — removing an idea is a
    manual file edit or the promotion above; nothing in the file tracks state.

## Implementation decisions

- **Two new subcommands in `bin/bench.sh`**, dispatched from the `case` at the foot of
  the file and defined alongside `structure()`:
  - `idea` — text is all args after the subcommand, joined; append `- <date +%F>  <text>`
    to `<repo_root>/ROADMAP.md`, creating the file if absent, then print a
    `parked: <text>` confirmation to stdout (agent-facing CLI feedback, per the `bench-craft-cli`
    skill). No text → usage to stderr, exit 2 (matches `link`'s usage convention).
    `repo_root` already supplies the not-in-a-repo error.
  - `roadmap` — print `<repo_root>/ROADMAP.md`; absent or empty → `roadmap empty` to
    stdout, exit 0.
- **Storage:** `ROADMAP.md` at repo root, append-only, one entry per line. **Committed**
  (product-facing, visible — the point is "I can see it"), not gitignored. **Not** added
  to `package.json` `files[]`: it's per-consumer content generated in the user's repo,
  not part of the kit package.
- **No status/lifecycle** in the file (#9). The sink stays dumb; an idea's real state,
  once committed, lives in its map/spec/issue.
- **Promotion is agent-side**, via a `/bench-ideate` command-file edit, not shell:
  a cold invoke reads `ROADMAP.md`, offers the parked items, seeds the map from the
  chosen one, and deletes that line when it writes the map. No `bench` removal command.
- **Out of B:** the ambient surface and its footer count that *renders* this roadmap —
  feature A.

## Testing decisions

- **A good test here** exercises external CLI behavior at the seam — run the real
  subcommand in a throwaway git repo and assert `ROADMAP.md` content and stdout — never
  internals.
- **Seam:** the CLI (`bench idea` / `bench roadmap`). One seam, the highest that
  exercises real behavior. **Prior art:** the `bench init` and `bench link` contract
  blocks already in `.bench/gate.sh` use exactly this throwaway-repo, exercise-the-real-
  CLI pattern — extend it.
- **Gate contract block** (add to `.bench/gate.sh`) asserting: `idea` appends a
  timestamped line; `idea` creates the file when absent; no-arg `idea` errors and
  appends nothing; `roadmap` prints entries; `roadmap` on empty says empty.
- **Canary fixture** (`tests/canary/`) so the new contract can't rot into an always-pass
  — a broken `bench.sh` + an `EXPECT` substring, following the existing fixture shape.
- **`/bench-ideate` promotion** is agent-instruction, not shell, so it's verified by
  reading the command file; add a kit-conformance grep to the gate asserting the command
  references `ROADMAP.md` and the auto-remove-on-map-creation behavior, so the
  instruction can't silently disappear.
- **Gate command:** the project gate — `.bench/gate.sh` (`bench gate`).

## Out of scope

- **Feature A — the ambient-feedback surface** (`bench status` + SessionStart, six
  signals, severity ladder) and the zero-severity footer that *displays* this roadmap.
  A genuinely separate capability with its own spec; decided two-specs-B-first. Est:
  several hours (renderer, signals, ladder, SessionStart wiring, gate-result cache).
- **`bench roadmap promote <n>` / status fields.** #9 decided promotion is manual via
  `/bench-ideate`; a CLI promote helper or lifecycle tracking is a distinct future
  command, not the rest of this one. Est ~30 min if ever wanted.
- **CLI removal/editing of entries** (`bench idea --remove`, reorder). Manual file edit
  is the answer; a CLI editor is a separate capability. Est ~30 min.
