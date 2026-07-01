# git safety classifier

## Problem

The PreToolUse git guard enforces the rule that the agent has no destructive git
authority, but its implementation is a narrow denylist. It blocks `git reset --hard`,
force clean, push, rebase, and branch delete, yet allows path-level destructive
commands such as `git checkout -- README.md` and `git restore README.md`. Those
commands can discard user edits, so the guard's interface is broader than the
behavior behind it.

## Solution

Deepen destructive-git classification into one policy module used by the hook
adapter. The classifier blocks path-level checkout/restore in addition to the
existing destructive operations, while continuing to allow harmless read commands
such as `git status --short`. The hook remains a shell entry point because Claude and
Codex adapters call command paths, but the classification rule is one interface with
gate-tested examples.

## User stories

1. As a user with local edits, I want `git checkout -- <path>` blocked, so an agent
   cannot discard a file by path.
2. As a user with local edits, I want `git restore <path>` blocked, so an agent
   cannot revert tracked content by path.
3. As a user relying on the existing guard, I want current blocks for push, hard
   reset, force clean, branch delete, and rebase preserved.
4. As an agent doing read-only inspection, I want harmless commands such as
   `git status --short` to remain allowed.
5. As a kit maintainer, I want the runtime gate to assert both blocked and allowed
   examples, so the classifier cannot silently narrow again.

## Implementation decisions

- **Primary seam:** the `block-dangerous-git.sh` hook command. Tests feed it the same
  JSON shape a harness sends and assert exit codes plus stderr.
- Keep the shell script as the adapter, but make the Python classification block own
  destructive intent. Do not duplicate policy in Claude or Codex adapter JSON.
- Treat `git restore` with a pathspec as destructive. It mutates the worktree or
  index even when the source is implicit.
- Treat `git checkout -- <path>` and equivalent pathspec checkout forms as
  destructive. Branch creation or read-only git commands are not part of this spec.
- Keep existing global-option parsing (`git -C`, `--git-dir`, `--work-tree`) so the
  guard still finds git invocations embedded in shell commands.

## Testing decisions

- **Good tests here** exercise the hook command as a black box with representative
  commands. A passing implementation is one that blocks destructive examples with
  exit 2 and `BLOCKED:` output, and allows read-only examples with exit 0.
- **Seam:** `block-dangerous-git.sh`, because it is the interface harness adapters
  execute.
- **Prior art:** `.bench/gate-runtime-contracts.sh` already has a guard block for
  push, hard reset, clean, and harmless status. Extend that block.
- **Gate command:** `bench gate`.

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | `git checkout -- README.md` is blocked with exit 2 and actionable output. | `block-dangerous-git.sh` | Observed red before implementation: feeding `{"tool_input":{"command":"git checkout -- README.md"}}` exits 0. | The current classifier only blocks checkout when args contain literal `.`, so this path-level revert slips through. |
| 2 | `git restore README.md` is blocked with exit 2 and actionable output. | `block-dangerous-git.sh` | Observed red before implementation: feeding `{"tool_input":{"command":"git restore README.md"}}` exits 0. | The current classifier only blocks restore when args contain literal `.`, so this path-level revert slips through. |
| 3 | Existing destructive examples remain blocked. | `block-dangerous-git.sh` | Already covered by runtime contract examples for push, hard reset, force clean, and global-option forms. | This keeps the fix additive rather than trading one blocked class for another. |
| 4 | `git -C . status --short` remains allowed. | `block-dangerous-git.sh` | Already covered by the runtime guard contract's allowed status example. | The classifier must distinguish destructive mutations from normal read-only inspection. |
| 5 | The project gate fails if path-level checkout/restore are allowed again. | Project gate: runtime contracts | Not TDD-able before implementation beyond the two observed red hook probes. | The committed gate examples become the regression protection. |

## Out of scope

- Replacing shell command parsing with a full shell AST. That is a larger classifier
  product with quoting, pipelines, and aliases to decide, ~1-2 hours.
- Blocking every possible branch switch or checkout form. This spec targets confirmed
  destructive path restore cases; broader git workflow policy needs a separate
  reviewer decision, ~45 minutes.
- Adding a git pre-commit hook. The current enforcement layers are harness hooks, the
  pre-push guard, and the gate.
