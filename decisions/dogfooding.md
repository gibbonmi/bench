# Dogfooding benchkit-on-benchkit

How to develop the kit using the kit without breakage in the kit-under-development
corrupting the development workflow. The egg bites in three places: the **gate**
(break it, lose the oracle), the **skills/commands** (break a trigger, next session
misfires), the **CLI** (break `bench shift`, lose the loop). The current session
loads skills at start, so mid-session edits hit the *next* session, not this one.

Frontier resolved — path to build is clear.

## #1: What dogfooding posture — self-host or pinned?

Type: Grill

### Question
How much does building the kit depend on the in-progress kit? Self-host (in-tree ==
running kit, truest signal, can hand breakage forward) vs pinned (develop against a
stable installed copy, egg dissolved, weak signal until promote).

### Answer
**True self-host + guard.** The in-tree kit is the running kit; changes are eaten
immediately for the best breakage signal. The cost (handing a broken workflow
forward) is bought down by #2. Settles that there is no second pinned copy to reason
about — the dev repo is the source of truth, and the gate must be trustworthy enough
to commit against.

## #2: How does the gate guard itself?

Blocked by: #1
Type: Grill

### Question
A self-hosted gate is both oracle and code-under-development. It `bash -n`s the CLI
and hooks but not itself, and never tests its own checks against bad input — a check
that silently always-passes would go unnoticed and the oracle would lose its teeth.
How does the gate prove it still bites?

### Answer
**Fixture canary.** Keep a small set of known-broken kit fixtures (bad frontmatter,
a dangling index reference, a no-op `init`, an extensionless gate ref). The gate runs
its checks against each and asserts RED; a fixture going green means a check rotted.
Run the fixtures through the gate as self-contained mini-repos (the check-1d pattern:
exercise the real path in a throwaway dir, assert the expected exit), rather than
refactoring every check to take a target dir. Fixtures mirror the failure modes the
gate already claims to catch — including the three real regressions seen so far
(dangling index, referenced-but-uncreated file, extensionless gate ref).

## #3: Does the canary cover behavioral skill firing?

Blocked by: #2
Type: Grill

### Question
Structural checks + the canary catch "does the gate bite" and "is the kit
structurally valid," but not "does a skill actually trigger and produce sane output"
— that needs a running agent.

### Answer
**Out of scope for the gate — deliberately.** Behavioral skill firing is
non-deterministic and needs an agent, so it stays out of the deterministic oracle.
Recommend handling it out of band (a periodic manual or `/resynthesize`-time smoke
check), never as a gate check. Recorded here so a future session doesn't try to build
agent-behavior testing into the gate.

## Build shape (for /spec)

A `tests/canary/` (or similar) tree of broken mini-repo fixtures + one gate check
that runs each through the gate and asserts the expected RED. Small, deterministic,
in-tree. No slicing decision — single focused build.
