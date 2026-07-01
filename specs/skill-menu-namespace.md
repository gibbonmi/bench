# Skill menu namespace

## Problem

Bench now exposes human-run workflow phases as explicit `$bench-*` skills, but the
model-invoked craft guidance skills still use visible `bench-craft-*` names. In a
skill or slash menu, filtering for `bench` shows both the phases a reviewer should
choose and the internal craft guidance the model should reach on its own. That
pollutes the user surface and makes the menu harder to scan.

## Solution

Keep human-run phase adapters in the `$bench-*` namespace. Rename only the visible
`name:` frontmatter of model-invoked craft skills to `craft-*`, while leaving their
`.agents/skills/bench-craft-*` directories stable as installed file paths.

The menu-facing names become:

| Current visible name | New visible name |
| --- | --- |
| `bench-craft-adr` | `craft-adr` |
| `bench-craft-cli` | `craft-cli` |
| `bench-craft-design-system` | `craft-design-system` |
| `bench-craft-grill` | `craft-grill` |
| `bench-craft-seams` | `craft-seams` |
| `bench-craft-skills` | `craft-skills` |
| `bench-craft-synthesis` | `craft-synthesis` |
| `bench-craft-tdd` | `craft-tdd` |

## User stories

1. As a reviewer opening the `$` or `/` menu and typing `bench`, I want to see only
   human-run Bench phases, so I can choose the next workflow step without internal
   guidance skills competing for attention.
2. As a reviewer, I want craft guidance skills to have clear `craft-*` names, so if
   I do invoke one directly, the name explains that it is generation guidance rather
   than a workflow phase.
3. As an agent, I want model-invoked craft skills to keep their descriptions, so I
   can still reach them autonomously when their trigger applies.
4. As an agent, I want current-state command and skill prose to refer to the new
   `craft-*` names, so cross-skill handoffs point at names visible in the harness.
5. As a Bench maintainer, I want existing `.agents/skills/bench-craft-*` directories
   to remain stable, so AGENTS.md path indexes, package contents, link manifests, and
   Claude adapter symlinks do not churn for a menu-only rename.
6. As a Bench maintainer, I want the gate to assert the visible-name split, so a
   future craft skill cannot accidentally re-enter the `$bench-*` user namespace.
7. As a Bench maintainer, I want current-state docs to explain the namespace split
   without rewriting historical specs and decisions, so users see the current rule
   and git preserves the naming history.

## Implementation decisions

Only craft skill frontmatter `name:` values change. Directory basenames stay
`bench-craft-*`, and AGENTS.md path entries keep pointing at those directories.

Use these exact frontmatter names:

- `.agents/skills/bench-craft-adr/SKILL.md` -> `name: craft-adr`
- `.agents/skills/bench-craft-cli/SKILL.md` -> `name: craft-cli`
- `.agents/skills/bench-craft-design-system/SKILL.md` -> `name: craft-design-system`
- `.agents/skills/bench-craft-grill/SKILL.md` -> `name: craft-grill`
- `.agents/skills/bench-craft-seams/SKILL.md` -> `name: craft-seams`
- `.agents/skills/bench-craft-skills/SKILL.md` -> `name: craft-skills`
- `.agents/skills/bench-craft-synthesis/SKILL.md` -> `name: craft-synthesis`
- `.agents/skills/bench-craft-tdd/SKILL.md` -> `name: craft-tdd`

Keep command adapter skills unchanged: their visible names, directories, and
command files must continue to share one basename because the gate enforces that
one-to-one command adapter contract.

Update current-state prose that names craft skills as invocations or handoffs:

- craft skill descriptions and bodies when they name another craft skill
- `.agents/commands/*.md` references to craft skills
- `bin/bench.sh` prompts and status action strings
- `.bench/gate.sh` status assertions and any conformance literals
- `README.md`, `projects/benchkit.md`, `HANDOFF.md`, and `CONTEXT.md` where they
  name craft skills as skills

Keep file-path references such as `.agents/skills/bench-craft-seams/SKILL.md`
unchanged unless they are intentionally documenting the on-disk path.

## Testing decisions

- **Primary seam:** `bench gate`. This change spans skill frontmatter, command prose,
  package/link behavior, and docs. The existing kit gate is the highest seam that
  exercises the shippable content surface.
- **Visible-name contract:** add a gate assertion that every `bench-craft-*` skill
  directory has a `name: craft-*` frontmatter value, and that no non-command-adapter
  skill uses `name: bench-*`.
- **Command-adapter contract:** keep the existing gate assertion that every
  `.agents/commands/<name>.md` has a matching `.agents/skills/<name>/SKILL.md` with
  `name: <name>`, because human-run phases remain `$bench-*`.
- **Current-state sweep:** after edits, search current-state files for
  `bench-craft-` references. Remaining hits must be either on-disk file paths,
  historical specs/decisions, or explicit text explaining old-to-new naming. Prose
  that says "use `bench-craft-*`" should become `craft-*`.
- **Menu proxy check:** list skill frontmatter names and verify the sorted visible
  names put all human-run phases under `bench-*` and all craft guidance under
  `craft-*`.
- **Gate command:** `bin/bench.sh gate`.

## Out of scope

- **Renaming `.agents/skills/bench-craft-*` directories** — this is path churn for a
  menu-only problem and would touch AGENTS.md indexes, link manifests, package
  assertions, and Claude adapter paths. Est. 30-45 min, not needed.
- **Adding alias skills for old `bench-craft-*` visible names** — aliases would keep
  the unwanted menu pollution and double the surface. Est. 20-30 min, rejected.
- **Changing human-run `$bench-*` phase names again** — recently resolved and already
  fit the user menu goal. Est. 60-90 min, separate naming problem.
- **Rewriting historical decisions/specs/changelog entries** — history can contain
  old names. Est. 30-60 min, not desired.
