# Roadmap context routing

## Result

Harnesses discover the phase differently, but none owns query execution. Claude Code
loads the canonical command through its `.claude/commands` adapter, Codex's explicit
skill reads the same command, and other AGENTS.md harnesses read that command directly
([harness invocation](../.bench/BENCH-reference.md#harness-invocation)). Once the phase
runs a shell command, wrapper selection—not the harness—determines the binary and bytes.

The portable execution path is the linked repo's local wrapper. `bench link` installs
the wrapper plus a matching binary snapshot ([link plan](../internal/adopt/link.go#L85));
the wrapper executes that local binary and, from a linked worktree where the ignored
binary is absent, re-anchors resolution to the main checkout
([binary resolution](../bin/bench.sh#L119)). Repo root discovery belongs to the Go core,
so an absolute local-wrapper invocation is invariant to deep CWD.

Hooks and headless shift adapters already use one resolver whose order is linked local
wrapper, source-kit wrapper, then PATH ([resolver](../.bench/lib/resolve-bench.sh#L16)).
The adapters only launch Claude, Codex, or OpenCode with a prompt; they do not proxy
later CLI stdout ([Claude adapter](../.bench/adapters/claude), [Codex adapter](../.bench/adapters/codex),
[OpenCode adapter](../.bench/adapters/opencode)). The context query does not belong in a
hook: SessionStart intentionally emits only the bounded ambient status and guard brief
([session-start](../.bench/hooks/session-start.sh#L29)).

## Verified bytes and exits

A linked throwaway repo used the existing AXI `bench maps` query as a routing surrogate.
The real-kit wrapper, linked by-path wrapper, current global wrapper, deep-CWD linked
wrapper, shared resolver, linked-worktree wrapper, and deep linked-worktree wrapper all
returned exit 0, empty stderr, and SHA-256
`2fc75b5f47357ab1e0c535291f221a6faba3216dda77e849fa0872832d9d970a` for:

```text
maps[0]{map,ticket,type,state}:
```

The linked-worktree copy had no local binary; it reached the main checkout's binary.
The usage edge was byte-identical at exit 2, and the outside-repo structured error was
byte-identical at exit 1 with empty stderr. The repository's executable routing probes
remain runnable with:

```sh
go test -count=1 ./internal/contract/surface \
  -run '^(TestGoRoutingContracts|TestLinkContracts)$'
```

To reproduce the unresolved PATH edge from any linked repo, shadow `bench` in the shell
while comparing it with the shared resolver:

```sh
bench() { printf 'stale-global\n'; return 7; }
export -f bench
bench maps
. .bench/lib/resolve-bench.sh
"$(bench_resolve_wrapper)" maps
```

The first call returns the shadow's bytes and exit 7; the resolver still selects
`.bench/bin/bench.sh` and returns the canonical TOON at exit 0. A globally installed
wrapper is explicitly convenience rather than hook infrastructure
([hook layers](../.bench/BENCH-reference.md#hook-layers)), so bare `bench` cannot prove
repo-version identity.

## Consequence

The top-level wrapper forwards to the linked/source wrapper before dispatch, with
self-resolution preventing recursion. This keeps bare `bench roadmap --context` as the
portable public surface while making the repo-local version authoritative for every
harness and human caller. `/bench-what-next` does not embed resolver-specific shell
ceremony. Gate research must prove forwarding across root, deep-CWD, linked-worktree,
PATH-shadow, and missing-local-wrapper states.
