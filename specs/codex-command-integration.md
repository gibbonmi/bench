# Codex command integration

## Problem

Bench's workflow phases are command files under the portable `.agents/commands`
surface, but Codex does not treat that directory as an invocation surface. In Codex
today, the reviewer can see repo skills but cannot invoke Bench phases as first-class
Codex workflows unless they know to ask the agent to read a command file by path.
That makes Bench feel less command-first in Codex than it does in Claude Code.

## Solution

Expose each Bench phase as an explicit Codex skill whose only job is to load and
follow the matching command file. The canonical phase body stays in
`.agents/commands`; the new skill is an adapter for Codex invocation and discovery.
Custom prompts are not used because Codex marks them deprecated and user-local. The
adapter skills install by default through `bench link`, because Codex support is part
of Bench's harness-neutral promise.

## User stories

1. As a Codex reviewer, I want to invoke each Bench phase as `$bench-setup`,
   `$bench-ideate`, `$bench-spec`, `$bench-diagnose`, `$bench-build`,
   `$bench-review`, and `$bench-qa`, so I do not need to know the command file path.
2. As a Codex reviewer, I want maintenance phases such as the synthesis commands to
   have the same adapter treatment, so every installed Bench command has one
   discoverable Codex invocation path.
3. As a Bench maintainer, I want `.agents/commands` to remain the canonical phase
   body, so Claude Code adapters, Codex adapters, and other AGENTS.md harnesses do
   not drift into separate instructions.
4. As a Bench maintainer, I want adapter skills to be explicit-only, so Codex does
   not fire workflow phases implicitly from a vague task description.
5. As a Codex user, I want Bench to avoid deprecated custom prompts, so the supported
   integration follows Codex's current reusable-workflow surface.
6. As a project using `bench link`, I want the adapter skills installed by default
   with the rest of `.agents/skills`, so Codex invocation works after the normal link
   path.
7. As a kit maintainer, I want the gate to prove every command file has a matching
   adapter skill, so a new phase cannot be added without Codex invocation support.
8. As a kit maintainer, I want adapter skills to reference the exact command file
   they wrap, so renames and missing files fail deterministically.

## Implementation decisions

Use `.agents/skills/<command-name>/SKILL.md` for each command adapter. The skill
frontmatter `name:` matches the command basename, which gives the Codex invocation
form `$<command-name>`; for example, the `bench-spec` command gets a `bench-spec`
skill.

Adapter skill bodies stay tiny. Each one tells the agent to read the matching
`.agents/commands/<command-name>.md` file completely and follow it as the active
Bench phase. Do not copy the command body into the skill; copying would create a
second command source and reintroduce the drift this feature is meant to remove.

Mark each adapter skill explicit-only using Codex skill metadata. The documented
Codex mechanism is `agents/openai.yaml` with `policy.allow_implicit_invocation:
false`; if Bench also keeps a frontmatter convention for other harnesses, that
metadata may be present too, but the Codex behavior depends on the Codex metadata.

Install the adapters through the existing `.agents/skills` surface. `bench link`
already copies that tree into linked repos and fails on project-owned name conflicts,
so no user-local `~/.codex` files, custom prompts, or separate plugin-only install
path are part of this feature.

The adapter naming intentionally means there are two Bench skill families:
`bench-craft-*` skills for reusable generation guidance, and `bench-*` skills for
explicit workflow phases. The difference is invocation policy, not directory root:
craft skills may be model-invoked; phase adapters are explicit-only.

Because new skill directories are installed into `.agents/skills`, update the
working agreement's skills index so the existing kit conformance rule remains true:
every skill directory on disk is referenced, and every indexed skill exists.

## Testing decisions

- **Primary seam:** the kit content surface plus the project gate. The observable
  contract is files on disk and gate conformance, not an interactive Codex TUI menu.
- **Good test:** run `bench gate` and assert the repository contract: every
  `.agents/commands/*.md` file has a matching `.agents/skills/<command-name>/SKILL.md`,
  each adapter skill has frontmatter `name: <command-name>`, each adapter body names
  `.agents/commands/<command-name>.md`, and each adapter carries explicit-only Codex
  metadata.
- **Install contract:** extend the existing `bench link` throwaway-repo checks enough
  to prove a linked repo receives at least one command adapter skill, with the generic
  command-adapter conformance proving the full set.
- **Index contract:** rely on the existing AGENTS.md skills-index checks after adding
  the adapter skills to the index; do not weaken the index check to make the new
  directories pass.
- **No custom-prompt contract:** no `~/.codex/prompts` writes, no packaged custom
  prompt files, and no docs that present custom prompts as the supported Bench path.
- **Gate command:** `bench gate`.

## Out of scope

- **Literal top-level `/bench-*` slash-menu entries in Codex** — that path relies on
  Codex custom prompts, which are deprecated and user-local. Personal prompt shims are
  a workaround, not a kit feature. Est. ~20 min if someone wants a local-only helper.
- **Plugin packaging for Bench** — useful for broader distribution, but not required
  for repo-local Codex invocation because Codex already scans `.agents/skills`.
  Separate packaging capability. Est. 60-120 min.
- **Changing the canonical command format** — `.agents/commands/*.md` remains the
  phase source. Rewriting phases as full skills would duplicate content instead of
  adapting it. Est. 45-60 min and intentionally rejected for this slice.
- **Autonomous `bench shift` command-contract portability** — this feature is about
  user-invoked phases in Codex, not the headless agent command used inside the shift
  loop. Separate operational-loop capability. Est. 60-120 min.
