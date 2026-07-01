# Human-facing Bench phase names

## Problem

Bench now exposes its workflow phases as Codex skills, but the current names still
come from internal Bench vocabulary: `$bench-ideate`, `$bench-qa`, `$bench-update`,
and similar. Those names are consistent, but they do not tell a human scanning a
skill list why to run each phase or where it sits in the workflow. The phase surface
needs names that are self-explanatory across Claude, Codex, and other harnesses.

## Solution

Rename the portable `.agents` phase surface to verb/object names that describe the
job the reviewer is starting. Rename the canonical command files and the matching
Codex adapter skill directories together, keeping the same basename for each pair so
there is no mapping layer. Harness adapters continue to point at `.agents`, so Claude
and Codex share one human-facing vocabulary.

| Current phase | New phase |
| --- | --- |
| `bench-setup` | `bench-setup-repo` |
| `bench-ideate` | `bench-shape-idea` |
| `bench-spec` | `bench-write-spec` |
| `bench-diagnose` | `bench-debug` |
| `bench-build` | `bench-implement-spec` |
| `bench-review` | `bench-review-implementation` |
| `bench-qa` | `bench-final-check` |
| `bench-update` | `bench-update-kit` |
| `bench-learn` | `bench-integrate-learnings` |

## User stories

1. As a reviewer in Codex, I want the explicit `$bench-*` phase skills to describe
   why I would run them, so I can pick the right phase from the skill surface without
   already knowing Bench vocabulary.
2. As a Claude Code reviewer, I want the slash command names to match the Codex skill
   names, so moving between harnesses does not require translating phase names.
3. As an agent, I want each `.agents/commands/<name>.md` file to have the matching
   `.agents/skills/<name>/SKILL.md` adapter, so the command and skill surfaces cannot
   drift.
4. As a Bench maintainer, I want the adapter skill frontmatter, descriptions, headings,
   and command-file references updated to the new names, so Codex exposes the new
   names and still delegates to the canonical command file.
5. As a reviewer, I want `.bench/BENCH.md` to explain the new phase sequence, so the
   operating guide remains the canonical source for how Bench is used.
6. As a maintainer, I want `AGENTS.md`, `README.md`, and `projects/benchkit.md` to point
   at the new current names without duplicating the full operating guide.
7. As a maintainer, I want `.bench/gate.sh` to assert the renamed command/adapter
   contract, so a missing rename fails at the oracle rather than in a future session.
8. As a user, I want no surviving old phase identifiers on current-state surfaces, so
   stale names do not appear in command pickers, docs, or gate output.

## Implementation decisions

Rename the portable command files in `.agents/commands/` and the matching Codex
adapter skill directories in `.agents/skills/`. Use `git mv` for tracked files and
directories so history survives.

Keep the invariant simple: after the rename,
`.agents/commands/<name>.md` pairs with `.agents/skills/<name>/SKILL.md`. There is no
manifest or mapping table for old-to-new names in the implementation; the only mapping
is this spec and the decision map.

Update each adapter skill:

- frontmatter `name:` to the new basename
- description to the new invocation form and matching command file
- heading to `$<new-name>`
- body reference to `.agents/commands/<new-name>.md`
- `agents/openai.yaml` remains explicit-only

Update each command file where the phase name appears in frontmatter, H1, prose, or
handoffs. Keep concept words such as "spec", "review", "gate", and "idea" when they are
plain English, not identifiers.

Update current-state docs and contracts:

- `.bench/BENCH.md` workflow list, harness invocation list, and command references
- `AGENTS.md` command list and any short pointers
- `README.md` onboarding and harness-switching references
- `projects/benchkit.md` content-surface and gate descriptions
- `.bench/gate.sh` hardcoded file paths, package/link checks, status action strings,
  and command-adapter conformance comments
- `bin/bench.sh` user-facing status/action strings if they name old phases

Do not rewrite historical `decisions/`, older `specs/`, `CHANGELOG.md`, or `ROADMAP.md`
only to make old identifiers disappear. The build may update this spec and
`decisions/human-facing-skill-names.md` if needed because they are the active artifacts
for this rename.

## Testing decisions

- **Primary seam:** `bench gate`. This rename spans files, docs, package contents, and
  link behavior; the project oracle is the highest existing seam that exercises the
  shipped kit surface.
- **Gate contract:** `.bench/gate.sh` should continue to prove every command file is
  referenced, every command has a same-basename Codex adapter skill, adapter
  frontmatter names match, adapter bodies reference the exact command file, and
  adapter metadata disables implicit invocation.
- **Install/package contract:** the existing link and npm dry-run checks should name
  at least one renamed command and one renamed adapter skill so packaging and link
  fixtures move with the rename.
- **Completeness sweep:** after edits, run a grep sweep over current-state surfaces for
  old identifiers:
  `bench-setup`, `bench-ideate`, `bench-spec`, `bench-diagnose`, `bench-build`,
  `bench-review`, `bench-qa`, `bench-update`, `bench-learn`.
  Exclude historical `decisions/`, older `specs/`, `CHANGELOG.md`, `ROADMAP.md`, and
  git metadata. Eyeball any remaining hits before calling the build done.
- **Gate command:** `bench gate`.

## Out of scope

- **Renaming `bench-craft-*` guidance skills** — separate naming problem for
  auto-invoked know-how, not part of the human phase surface. Est. 30-45 min.
- **Shortening `bench` to `bn`** — rejected for this feature because it saves a few
  characters while losing namespace clarity. Est. 0 min.
- **Adding a mapping/alias layer from old names to new names** — intentionally avoided
  so harnesses share one vocabulary. A compatibility alias system would be a separate
  support feature. Est. 45-90 min.
- **Rewriting historical decisions/specs/changelog entries** — history may contain old
  names. Updating it would add noise and risks corrupting records. Est. 30-60 min and
  not desired.
