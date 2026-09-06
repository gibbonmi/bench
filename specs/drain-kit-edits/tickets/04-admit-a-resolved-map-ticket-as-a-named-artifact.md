# Admit a resolved map ticket answer as a named reviewed artifact

Blocked by: none
Writes: .agents/commands/bench-write-spec.md, .claude/commands/bench-write-spec.md
Covers: none

## What to build

A shaping map that holds one ticket per factory item cannot reach `ready`
before its last item decides. FT303 resolved through map tickets #7 to #12 with
a research asset while #8 and #9 stay open. The FT302 specs cited the same
kind of source. The entry contract's `Named reviewed artifact` bullet already
admits settled decisions. It does not say that a resolved ticket answer inside
a shaping map qualifies, or what the phase copies.

Add one or two ASD-STE100 sentences to the `Named reviewed artifact` bullet in
`## Entry contract`. A resolved ticket answer and its asset in a shaping map
count as a named reviewed artifact when the map is a multi-item sequence. The
phase copies only that ticket's asset, and the map stays in place. Mirror the
edit byte-exact into `.claude/commands/bench-write-spec.md`, because the
`entry-point-parity` check compares the two copies.

## Acceptance

- [ ] The `Named reviewed artifact` bullet admits a resolved ticket answer and its asset from a multi-item shaping map.
- [ ] The bullet says the phase copies only that ticket's asset.
- [ ] `.agents/commands/bench-write-spec.md` and `.claude/commands/bench-write-spec.md` are byte-identical.
- [ ] `bench gate-prose . -- .agents/commands/bench-write-spec.md` exits 0.
- [ ] `bench test --check entry-point-parity`, `bench test --check guidance-prose-budgets`, and `bench test --check kit-compliance` exit 0 over the worktree.
