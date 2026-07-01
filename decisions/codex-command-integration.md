# Codex command integration — make Bench phases invokable in Codex

> **BOOTSTRAPPED (2026-06-30).** Bench commands live in `.agents/commands/` as
> portable phase files. Claude Code can expose those through `.claude/commands/`;
> Codex currently does not scan `.agents/commands` as slash commands. This map
> decides the Codex adapter layer that makes the same Bench phases usable without
> pretending the `bench` CLI is the reviewer-facing surface.

## Grounding

- Bench's command files are the canonical workflow phases:
  `/bench-setup`, `/bench-ideate`, `/bench-spec`, `/bench-diagnose`,
  `/bench-build`, `/bench-review`, `/bench-qa`, and `/resynthesize`.
- Codex scans `.agents/skills/**/SKILL.md` for skills. Skills can be invoked
  explicitly with `$skill-name` or through `/skills`, and can also be selected
  implicitly from their descriptions.
- Codex custom prompts can appear in the `/` menu as `/prompts:<name>`, but the
  Codex manual marks custom prompts deprecated, user-local, and stored under
  `~/.codex/prompts/*.md`, not shared through the repository.
- Therefore, `.agents/commands/*.md` is currently a Bench convention that Codex
  agents can read and follow, but not a Codex-native slash-command surface.
- `decisions/command-first-usability.md` decides that the reviewer-facing surface
  is commands plus conversation; the `bench` CLI is worker-facing infrastructure.

## #1: What does "work with Codex" mean for Bench commands?

Type: Grill

### Question
Should "work with Codex" mean native slash-menu entries, explicit `$bench-*` skills,
or simply agent-readable phase files that Codex is instructed to follow?

### Answer
Resolved from grill. For Bench commands in Codex, "work with Codex" means the
phases are first-class Codex-invokable as explicit `$bench-*` skills. Literal
top-level `/bench-*` slash-menu entries are not the target because the Codex
surface that provides user-defined slash entries is custom prompts, and custom
prompts are deprecated. The acceptable Codex user surface is `$bench-ideate` or
selection through `/skills`, backed by the same canonical phase files in
`.agents/commands/`.

## #2: Should Bench generate Codex custom-prompt shims?

Blocked by: #1
Type: Grill

### Question
If the desired experience is typing `/` and seeing Bench phases, should Bench create
`~/.codex/prompts/bench-*.md` shims that each read and follow the matching
`.agents/commands/*.md` file, despite custom prompts being deprecated and user-local?

### Answer
Resolved from grill. No. Bench should not build its Codex command integration on
custom prompts because the Codex manual marks them deprecated and user-local. A
deprecated compatibility surface is acceptable for a user's personal workaround,
but not as the kit's supported integration path.

## #3: Should commands be mirrored as explicit-only Codex skills instead?

Blocked by: #1
Type: Grill

### Question
If Codex's recommended reusable workflow surface is skills, should Bench expose each
phase as an explicit-only skill, or keep commands and skills distinct and avoid
duplicating phase instructions into `SKILL.md` files?

### Answer
Resolved from grill. Yes, but only as thin explicit-only adapter skills. Each
`$bench-*` skill should delegate to the matching `.agents/commands/*.md` phase file
rather than duplicate the phase body. This keeps commands and skills distinct:
`.agents/commands` remains the canonical workflow phase surface, while
`.agents/skills` provides Codex's supported invocation and discovery mechanism.

## #4: Where should Codex adapter files live?

Blocked by: #2, #3
Type: Research

### Question
If an adapter layer is needed, determine whether it belongs in the repository
(`.codex/`, `.agents/`, or plugin packaging), the user's Codex home
(`~/.codex/prompts`), or as a generated install artifact created by `bench link`.

### Answer
Resolved from grill. The adapter skills live in the repo's existing
`.agents/skills` surface. Codex already scans that tree, and Bench already ships it
through `bench link`, so no user-local `~/.codex` files or separate plugin-only
location are needed for repo use.

## #5: What should `bench link` install for Codex?

Blocked by: #4
Type: Grill

### Question
Once the adapter location is chosen, should `bench link` install Codex command
integration automatically, offer it as an explicit mode, or only document the manual
setup?

### Answer
Resolved from grill. `bench link` installs the `$bench-*` adapter skills
automatically as part of the normal `.agents/skills` install. They are not an
optional mode because Codex invocation is part of the harness-neutral promise, and a
missing adapter would make the command workflow silently harder to discover.

## #6: What is the gateable contract?

Blocked by: #2, #3, #4, #5
Type: Research

### Question
Define the tests or conformance checks that prove every `.agents/commands/*.md`
phase has a working Codex invocation path, without relying on an interactive Codex
TUI session in the gate.

### Answer
Resolved from grill. Gate the repository file contract, not the interactive Codex
TUI. For every `.agents/commands/<name>.md` phase, the gate should assert there is
a matching explicit adapter skill in `.agents/skills` using the chosen naming
pattern, that its frontmatter disables implicit invocation, and that its body
references the exact command file. This proves Bench ships the Codex invocation path
without requiring an interactive slash-menu test.
