# Bench command & skill renaming

Rename every Bench-owned command and skill into a `bench-*` namespace whose name
states its **intention** and, for commands, its **stage in the workflow**. Motive:
the bare names (`adr`, `axi`, `prep-shift`, `seams`…) don't say what they do, and
they collide visually with the many non-Bench skills a harness loads (`grill-me`,
`tdd`, `prototype`…). The `bench-` prefix disambiguates; the rename makes purpose
legible.

## #1: What gets renamed? — RESOLVED

Type: Grill

Both the 7 skills and the 8 commands, under one `bench-` namespace. A half-rename
(skills only, or commands only) leaves the surface you actually type as the bare
names that confuse — so both.

## #2: What's the naming grammar? — RESOLVED

Blocked by: #1
Type: Grill

- **Commands → `bench-<verb>`.** An action verb, ordered by the canonical workflow.
  **No numbers** — the workflow branches (the bug path is an alternate entry;
  setup/update/learn are lifecycle, not line positions), so a numeric `1..N` would
  misrepresent the shape and renumber on every insert.
- **Skills → `bench-craft-<noun>`.** One `craft` bucket for all of them — the
  `craft` infix is the visual signal for "auto-firing know-how" as opposed to "a
  phase I invoke." No sub-buckets (technique/standard/meta) — too fuzzy to name
  crisply over 7 items.
- **Stem rule.** Keep the stem when it's a recognizable term — industry or
  Bench-branded (`spec`, `build`, `seams`, `grill`, `adr`). Replace it only when
  it's opaque internal jargon (`axi`, `prep-shift`, `resynthesize`).

The distinction the scheme rests on: a **command is a place you go in the
workflow** (sequenced, you invoke it); a **skill is know-how you carry**
(cross-cutting, auto-fires). This is itself codified in `bench-craft-skills` as
user-invoked vs model-invoked.

## #3: The full name slate — RESOLVED

Blocked by: #2
Type: Grill

Commands:

| from | to |
|---|---|
| start-ideation | `bench-ideate` |
| spec | `bench-spec` |
| build | `bench-build` |
| fix-bug | `bench-diagnose` |
| prep-shift | `bench-review` |
| verify-gate | `bench-qa` |
| setup | `bench-setup` |
| resynthesize | `bench-update` + `bench-learn` (see #4) |

Skills:

| from | to |
|---|---|
| adr | `bench-craft-adr` |
| axi | `bench-craft-cli` |
| design-system | `bench-craft-design-system` |
| grill | `bench-craft-grill` |
| seams | `bench-craft-seams` |
| tdd-at-seams | `bench-craft-tdd` |
| writing-great-skills | `bench-craft-skills` |

Reasoning on the non-obvious renames:

- **`prep-shift → bench-review`** — "shift" is CLI jargon; the command is a
  pre-merge **semantic** review (does the diff follow conventions + match the
  spec). Distinct from `bench-qa`, the **mechanical** oracle: the review judges
  whether the code is *right*; the gate proves it's *correct*.
- **`verify-gate → bench-qa`** — QA covers all four checks (types, tests, lint,
  conformance) without foregrounding lint, which is the narrowest. The *concept*
  stays "the gate"; only the command surface is `bench-qa`. Consequence: `bench-qa`
  (slash) and `bench gate` (CLI) diverge in name — intentional.
- **`fix-bug → bench-diagnose`** — foregrounds the repro-first core of the bug path
  (the hard half); the `bench-` prefix keeps it clear of any loose `diagnose` skill.
- **`axi → bench-craft-cli`** — "axi" is opaque; it is CLI-design standards.

## #4: Should `resynthesize` stay one command? — RESOLVED

Blocked by: #3
Type: Grill

No — split into two, because they read different inputs:

- **`bench-update`** — compares this repo against the **upstream kit**, pulls
  improvements.
- **`bench-learn`** — promotes the **local `.bench/learnings.md`** journal into the
  kit.

Both are commands (deliberate actions you invoke), not skills. **This is a behavior
change, not a mechanical rename — it needs `/spec` → `/build`.**

## #5: Does the user use the CLI? — RESOLVED

Type: Grill

No. The reviewer's surface is **slash-commands + conversation**; the agent (the
worker) drives the `bench` CLI underneath. So:

- **Capture** stays the `bench idea` CLI verb, but the user never types it — they
  say "park this," the agent runs it. No `bench-capture` slash-command.
- **Out of scope:** `ROADMAP.md` (a file, not a command) and the `bench` CLI verbs
  (`link`, `init`, `gate`, `shift`, `idea`, `models`) — already `bench`-prefixed by
  construction.

## #6: Known follow-on edits — RESOLVED

Type: Research

Beyond the directory + frontmatter renames, the rename has a blast radius the spec
must cover:

- The skills index and command references in `AGENTS.md`, `.bench/BENCH.md`, and
  `CLAUDE.md`.
- Cross-references between command/skill files (e.g. the build phase points at
  `seams`/`tdd-at-seams`; ideation points at `spec`/`grill`).
- **`bench-craft-skills`' own example list** names the old commands and skills as
  its user-invoked / model-invoked examples — it must be rewritten to the new slate.
- The `.claude/` mirror of skills and commands.

## Next

`/spec`. Slice decision for the spec (yours): the **mechanical renames** (most of
the slate + reference fixes) are large-churn but simple; the **behavior changes**
(#4 split, #5 capture wiring) need real design. They can ship as one spec or two
slices — your call, against this finished map.

## Parked (not in this map)

"Should `bench-craft-skills` require good/bad output examples in new skills
(Pocock-style contrastive examples)?" — a real possible improvement to the skill's
*content*, orthogonal to this renaming. Parked to `ROADMAP.md`, decided separately.
