# Safe Bench link

## Problem

Bench is supposed to be safe to dogfood in existing projects, but the current link
path can replace project-owned instructions and agent assets. That makes the kit
dangerous at the exact moment a user is trying to adopt it. The install path needs
to incorporate Bench without taking ownership of project decisions, while still
setting up the portable agent surfaces and the harness-specific git-safety adapters.

## Solution

Make `bench link` a safe incorporation step. It installs Bench-owned assets in
managed locations, adds only a small managed Bench block to project instructions,
fails on unowned filename conflicts, and records enough ownership to make reruns
idempotent. Full Bench operating guidance lives outside the project instruction
file, and Claude/Codex hook adapters point at shared hook scripts instead of
duplicating policy.

## User stories

1. As a reviewer linking Bench into a repo with existing instructions, I want my
   project-owned instructions preserved byte-for-byte outside a Bench-managed block,
   so the kit cannot silently rewrite the working agreement.
2. As a reviewer linking Bench into a repo without instructions, I want a minimal
   project instruction file created with the Bench block, so agents know Bench is
   installed without loading a full manual.
3. As a reviewer rerunning `bench link`, I want only the managed Bench block updated,
   so repeated setup is safe and predictable.
4. As an agent in a Bench-linked repo, I want the managed Bench block to point me at
   the gate, the profile, the portable commands/skills, and the full Bench operating
   docs, so I can find the right context without bloating `AGENTS.md`.
5. As a project that already has a skill or command with the same name as a Bench
   asset, I want `bench link` to fail with a clear conflict report, so I can decide
   whether to rename, remove, or replace it.
6. As a project that has already been linked by Bench, I want `bench link` to update
   Bench-owned assets when they have not been locally modified, so kit upgrades are
   straightforward.
7. As a project owner, I want locally modified Bench-owned assets detected rather
   than overwritten, so local changes are not lost on relink.
8. As a Codex user, I want Bench to install Codex hook config for git-safety events,
   so the same guardrails exist in Codex sessions after hook review/trust.
9. As a Claude Code user, I want Bench to install Claude hook config for the same
   shared scripts, so Claude and Codex enforce the same policy.
10. As a harness-agnostic user, I want skills and commands installed under the
    portable agent surface first, so non-Claude harnesses are not second-class.
11. As a Claude Code user, I want Claude-facing skill and command paths to point at
    the portable assets, so there is one Bench source for those files.
12. As a user browsing a linked repo, I want `.claude/README.md` to explain that the
    Claude-facing skill and command paths are adapters to `.agents/`, so the install
    shape is visible without knowing symlink conventions.
13. As a package consumer, I want `bench link` to use copy mode by default, so a
    target repo does not change when the source kit checkout changes.
14. As a kit maintainer, I want explicit symlink mode preserved, so local kit
    development can still dogfood live changes when that is the intended choice.
15. As a kit maintainer, I want package contents to include only the assets needed
    for safe setup, so local settings files cannot be published by accident.

## Implementation decisions

- **Primary seam:** the `bench link` CLI contract. Tests exercise the real command
  against throwaway git repos and assert observable filesystem results and exit
  codes.
- `bench link` defaults to copy mode. `bench link copy` remains accepted, and
  `bench link symlink` is the only mode that points installed assets back to the
  source kit.
- `AGENTS.md` is project-owned except for one managed Bench block delimited by
  `<!-- bench:start -->` and `<!-- bench:end -->`.
- If no project instruction file exists, `bench link` creates one containing the
  managed Bench block. If one exists without a Bench block, it appends the block. If
  one exists with a Bench block, it replaces only that block.
- The managed Bench block is intentionally small. It states that the gate is the
  oracle, points to the full Bench operating docs, names the portable command and
  skill locations, and reminds agents that the reviewer owns merge decisions.
- Full Bench operating docs live in a Bench-owned document installed under the
  `.bench` tree. They are not inlined into project instructions and are not named
  `CONTEXT-BENCH.md`.
- `.claude/README.md` explains that `.claude/skills/` and `.claude/commands/` are
  Claude Code adapters to `.agents/`, and that `.claude/settings.json` points at
  shared scripts in `.bench/hooks/`.
- Skills and commands have a canonical portable install surface under `.agents`.
  Claude-facing skill and command paths are adapters that point to those assets or
  are generated from them; they are not a second source of truth.
- Shared hook scripts live in a Bench-owned location. Claude and Codex adapter config
  point at those scripts. The git pre-push hook remains the harness-independent
  backstop.
- A link manifest records Bench-owned installed assets and their fingerprints. A
  later link may update an asset that is still manifest-owned and unmodified. If the
  destination exists without manifest ownership, or differs from the manifest
  fingerprint, link fails and reports the conflict.
- Conflict reporting is all-or-nothing: `bench link` discovers conflicts before
  writing and exits nonzero without a partial install.
- Package contents should be an allowlist for the installable kit surface. Local
  settings and repo-only development artifacts are excluded even if they live under
  an otherwise installed directory.

## Testing decisions

- **What a good test is here:** run the real `bench link` command in throwaway git
  repos and assert exit code plus filesystem state. Do not test helper internals.
- **Gate seam:** add a behavioral gate check for the safe-link contract. It should
  run fast, create temporary repos, and leave no workspace state behind.
- **Fresh repo case:** linking a repo with no project instructions creates a small
  instruction file with exactly one managed Bench block, installs the full Bench
  docs, portable commands/skills, Claude adapter config, Codex adapter config, shared
  hook scripts, and the git pre-push hook.
- **Existing instructions case:** linking a repo with existing project instructions
  preserves all existing text outside the managed block and appends exactly one Bench
  block.
- **Rerun case:** linking the same repo twice replaces only Bench-owned content and
  does not duplicate the managed block.
- **Conflict case:** a pre-existing same-named skill or command that is not
  manifest-owned makes `bench link` fail with a clear conflict list and no partial
  writes.
- **Modified-managed case:** if a manifest-owned installed file changed locally,
  relinking fails instead of overwriting it.
- **Mode case:** default link mode copies assets into the repo without references
  back to the source kit; explicit symlink mode may point managed assets at the
  source kit.
- **Hook adapter case:** Claude and Codex hook config both reference the shared hook
  scripts, and the git pre-push hook is still installed.
- **Discoverability case:** `.claude/README.md` is installed and explains the
  `.claude` adapter relationship to `.agents` and `.bench/hooks`.
- **Package case:** a dry-run package inspection proves required install assets are
  included and local-only settings are excluded.
- **Gate command:** `bench gate`.

## Out of scope

- **Shift worktree ownership** — separate operational-loop capability, not part of
  safe setup — ~60-90 minutes.
- **Non-Claude autonomous agent command adapter** — separate harness execution
  contract for `bench shift`, not required to safely install the kit — ~60-120
  minutes.
- **Full README command-name cleanup outside affected setup/install sections** —
  separate docs cleanup; deferred by reviewer choosing the install-safety slice —
  ~20 minutes.
- **Explicit conflict-replacement mode** — powerful migration command beyond the safe
  default; needs its own review because it can overwrite project-owned assets — ~45
  minutes.
