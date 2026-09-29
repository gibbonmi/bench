# Explicit tree targets for tree-scoped Bench verbs (FT341)

Status: ready

## Destination

Each public Bench verb declares its tree scope in the command registry. A
tree-scoped verb takes `--in <label|primary>`, reads its tracked content from
that tree target, and names the tree target in its response. The directory
inference and the exec route for Bench commands then end, so a caller never
grades one tree while it means another.

## Notes

Domain: the command registry, the wrapper's kit selection, the worktree
lifecycle verbs, and the hooks and skills that call Bench. Consult
`bench-craft-domain` and `bench-craft-cli`.

Canonical terms for this map:

- **tree target** — the checkout that a tree-scoped verb reads: a worktree label or `primary`. Not "root", not "kit dir", not "target" alone.
- **tree-scoped verb** — a public registry leaf that reads or grades tracked content of one checkout.
- **repository-scoped verb** — a public registry leaf that reads only Git refs, the worktree ledger, or git-ignored capture files.

The lifecycle verbs keep `<target>` for their worktree operand; that operand is not a tree target.

A map-owned asset stays in the map's assets folder,
specs/tree-targets/decisions/tree-targets/assets/.

## Decisions so far

- [How does a caller name the tree target?](tree-targets/tickets/1.md): `--in <label|primary>` after the verb leaf; no `--in` means the primary checkout; no environment form.
- [Which leaves are tree-scoped, and where is that declared?](tree-targets/tickets/2.md): a required per-leaf registry field; tracked-content readers are tree-scoped; lifecycle and plumbing leaves stay outside `--in`.
- [Which running executable serves a named tree target?](tree-targets/tickets/3.md): in the kit repository, the target's worktree build; a stale build refuses and names `bench worktree build`.
- [How does a response name its tree target?](tree-targets/tickets/4.md): one `tree{target,head,dirty}` row and no pool path.
- [Which path derivation resolves a tree target?](tree-targets/tickets/5.md): `internal/canonicalpath`; the gate's `resolvedPath` becomes its caller.
- [How does the exec route for Bench commands end?](tree-targets/tickets/6.md): exec refuses a Bench child; spec B removes directory inference in the same landing.
- [How is FT341 delivered?](tree-targets/tickets/7.md): two specs from this map; spec A moves every tree-scoped leaf; the list filter is out.

## Not yet specified

## Spec-writer discretion

- Whether a refusal that comes after target resolution also carries the `tree{target,head,dirty}` row.
- How exec detects a Bench child in its command, provided that the detection reuses the follow-on guard's parser.
- The exact refusal text for an unknown, released, or ambiguous label, provided that it names the live labels or the listing command.

## Out of scope

- A `bench worktree list` filter by spec; FT125 owns precise readers.
- The comfort of `bench worktree exec` for commands that are not Bench commands; FT254 owns it.
- A change to the lifecycle verbs' `<target>` operand.
- Staleness of the primary checkout's own binary; the current freshness route keeps it.

## Sources

- Path: `roadmap/FT341.md`
  Supports: the destination and the occurrences that every ticket answers.
  Drift: a body edit or a new occurrence on the row.
- Path: `bin/bench.sh`
  Supports: tickets 1, 3, and 6. `kit_dir` selects a worktree's kit from the current directory, and the binary resolves from the primary checkout.
  Drift: a change to `kit_dir`, `main_tree_kit`, `bench_binary_path`, or `route_binary`.
- Path: `internal/worktree/exec.go`
  Supports: ticket 6. `execEnv` removes `BENCH_KIT` and `BENCH_RUN_BINARY` so that a child resolves its own kit.
  Drift: a change to `execEnv` or to the exec child's environment.
- Path: `internal/freshness/freshness.go`
  Supports: ticket 3. `Digest` covers every build input in the working tree, so an uncommitted Go edit makes a build stale.
  Drift: a change to `Digest` or `buildInputs`.
- Path: `.bench/gate.sh`
  Supports: ticket 2. The gate passes its root as an operand and its binary through `BENCH_RUN_BINARY`.
  Drift: a change to how the gate script starts `gate-phases`.
- Path: `internal/git/git.go`
  Supports: ticket 1. `Root` derives the root from the current directory, and 55 non-test call sites use it.
  Drift: a change to `Root` or `RootAt`.
