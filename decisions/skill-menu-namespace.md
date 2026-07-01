# Skill menu namespace — separate human phases from model guidance

> **GRILLED & RESOLVED (2026-06-30).** Keep human-run phase adapters in the
> `$bench-*` namespace. Move model-invoked craft guidance to visible `craft-*`
> skill names, while leaving the existing `.agents/skills/bench-craft-*`
> directories in place as stable file paths.

## #1: What namespace should model-invoked craft skills use?

Type: Grill

### Question
Human-invoked Bench phases and model-invoked craft guidance currently share the
`bench-*` skill namespace. That makes `$` / `/` menus show internal guidance skills
when the reviewer is looking for a phase to run. Should craft skills move to
`bench-z*`, or should they leave the `bench-*` user namespace?

### Answer
Resolved from grill. Keep human-run phase skills under `$bench-*`; rename the
model-invoked craft skills' visible names to `craft-*`:

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

Do not use `bench-z*`: it sorts later, but it still appears when a reviewer filters
for `bench`, and the `z` prefix does not explain what the skill does.

## #2: Is changing frontmatter names sufficient, or must directories also rename?

Blocked by: #1
Type: Research

### Question
The visible menu appears to use skill `name:` frontmatter, but linked installs,
Claude symlinks, gate checks, and cold-session docs may still depend on directory
basenames. What is the smallest rename surface that keeps menus clean without
breaking link/package/gate contracts?

### Answer
Changing craft skill frontmatter `name:` is sufficient for the visible menu, and
it is the smallest safe build. Keep `.agents/skills/bench-craft-*` directory
basenames stable because AGENTS.md indexes those paths, `bench link` copies the
whole `.agents/skills` tree, package contents include `.agents/`, and existing
Claude adapter paths point through `.agents/skills`.

The gate requires command adapter skill frontmatter names to match command
basenames, but it does not impose that rule on craft skills. Craft skills only need
frontmatter and AGENTS.md path coverage.

## #3: What current-state surfaces must reflect the namespace split?

Blocked by: #1, #2
Type: Research

### Question
If craft skills rename, which docs and contracts must change so users see only
human-run phases under `$bench-*` while agents can still discover craft guidance
reliably?

### Answer
Update current-state prose that names craft skills as visible invocations so the
agent sees and reaches the new `craft-*` names consistently. Keep file-path indexes
pointing at `.agents/skills/bench-craft-*/SKILL.md`.

Current-state surfaces to update:

- craft skill `name:` frontmatter
- craft skill descriptions and bodies when they refer to another craft skill by name
- `.agents/commands/*.md` references such as `bench-craft-grill`,
  `bench-craft-seams`, `bench-craft-tdd`, `bench-craft-design-system`, and
  `bench-craft-synthesis`
- `bin/bench.sh` status/action strings and prompts that name craft guidance
- `.bench/gate.sh` status contract assertions that check those strings
- `README.md`, `projects/benchkit.md`, `HANDOFF.md`, and any other current-state
  docs that name the craft skills as skills rather than as paths

Leave `AGENTS.md` path references unchanged unless the directory names change.
