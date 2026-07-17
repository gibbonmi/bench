# Roadmap context gate attachment

## Result

The gate attaches at three existing seams. Package tests pin each typed owner's parser
and the roadmap aggregator's size policy. One AXI contract drives the real
`bench roadmap --context` command and pins the complete snapshot, errors, offline
posture, and unchanged repository state. One surface contract pins top-level-wrapper
forwarding to the repo-local wrapper. The gate already runs all contract packages as
one phase, so no new gate phase or registry is needed
([phase table](../internal/gate/phases.go#L62)).

Canaries do not repeat schemas or derivations. Each sabotages a real public surface and
expects the distinct failure message from its owning contract: an incomplete context
renderer expects `AXI roadmap context completeness contract failed`; an old wrapper
that skips local forwarding expects `repo-local wrapper forwarding contract failed`;
and a command file that drops the context invocation expects
`bench-what-next dropped the roadmap context query`. All three belong under
`behavior-owned` except the guidance anchor, which belongs to
`workflow-guidance-anchors`; those families already route nested canaries to only the
owning contract or conformance phase ([phase selection](../internal/canary/canary.go#L28),
[fixture registry](../internal/conformance/registry_test.go#L24)).

## Public snapshot contract

The default output uses these fixed flat-table blocks, in this order. Every block is
present even at zero rows; the shared TOON emitter already owns self-describing empty
tables, escaping, and refusal of unrepresentable control bytes
([TOON emitter](../internal/toon/toon.go#L20)).

| Block | Fields |
|---|---|
| `context[1]` | `schema,full` |
| `sources[N]` | `source,state,bytes` |
| `roadmap_rows[N]` | `id,title,spec,spec_status,external_trigger,body,body_bytes,truncated` |
| `roadmap_sequence[N]` | `rank,text,command` |
| `ideas[N]` | `date,text,text_bytes,truncated` |
| `learnings[N]` | `date,title,state,body,body_bytes,truncated` |
| `structure[N]` | `kind,path,actual,limit,state,detail` |
| `specs[N]` | `slug,status,roadmap_id` |
| `spec_history[N]` | `slug,hash,date,kind,subject` |
| `git[1]` | `branch,default_branch,dirty,ahead,behind` |
| `git_changes[N]` | `status,path` |
| `gate_cache[1]` | `present,status,cached_tree,work_tree,timestamp,stale` |
| `parse_failures[N]` | `source,reason,raw,raw_bytes,truncated` |

`schema` is integer `2`; `full` is a boolean. `sources.state` is one of
`absent`, `empty`, `parsed`, or `malformed`, preserving absent-versus-empty evidence.
The seven `sources.source` rows are fixed and ordered as `ROADMAP.md`, `IDEAS.md`,
`.bench/learnings.md`, `.bench/structure.budgets`, `.bench/structure-accept`, `specs/`,
and `.git/bench-last-gate` (the last is a logical label for the worktree's actual
git-dir cache path). Historical reconcile and promotion evidence stays in Git rather
than being copied into the snapshot.
An unreadable required source or failed git derivation emits one structured AXI error
and exit 1 rather than a partial snapshot. Unknown or conflicting arguments emit usage
on stdout and exit 2. `-h`/`--help` emit usage on stdout and exit 0. Successful output
uses exit 0 and empty stderr, following the query-surface contract
([AXI error and exit posture](../.agents/skills/bench-craft-cli/SKILL.md#L53)).

Body-like fields (`body`, `text`, and parse-failure `raw`) show at most 4096 source
bytes by default, cut back to a valid UTF-8 boundary. Their `*_bytes` field always
reports the full byte count and `truncated` is true exactly when bytes were withheld.
`--full` removes that ceiling and sets `context.full=true`; it does not change schemas
or ordering. Boundary probes use 4095, 4096, and 4097 bytes plus a multibyte rune across
the cut. Invalid or TOON-unrepresentable bytes fail closed instead of being replaced.

Blocks and rows are deterministic: document order for roadmap rows, sequence, ideas,
and learnings; path order for structure, specs, and git changes;
newest-first for each slug's history, matching the existing history query
([history ordering](../internal/spec/history.go#L50)). Exact stdout must repeat
byte-for-byte on a second unchanged invocation.

## Observable red signals

| Seam | Black-box assertion | Red signal |
|---|---|---|
| Typed owner APIs and roadmap aggregation | Parsed facts, source states, deterministic order, malformed raw preservation, 4095/4096/4097-byte policy, and propagated read/git errors | Focused package test fails with the named fact or boundary; no renderer parsing is involved. |
| AXI command | Real wrapper emits every fixed block and field; empty/malformed/control-byte/error/help/`--full` cases keep stdout, stderr, and 0/1/2 exits; two runs are byte-identical | `go test -count=1 ./internal/contract/axi -run '^TestAXIRoadmapContextContracts$'` fails, with the completeness subtest using `AXI roadmap context completeness contract failed`. |
| Read-only and offline posture | A before/after manifest of repository files, modes, bytes, git status, index hash, HEAD, and gate cache is unchanged; PATH sentinels for `bench`, `curl`, `wget`, `gh`, `glab`, `claude`, `codex`, and `opencode` are never called; a planted gate script leaves no marker | The AXI contract reports the changed artifact or sentinel name. This observes effects rather than scanning imports. |
| Wrapper routing | Source/global wrapper at root and deep CWD selects a distinct repo-local wrapper, preserves all args and exact stdout/exit/stderr, avoids self-recursion, reaches the main checkout binary from a linked worktree, and falls back to itself when no local wrapper exists | `go test -count=1 ./internal/contract/surface -run '^TestGoRoutingContracts$'` fails with `repo-local wrapper forwarding contract failed`. |
| Phase consumption | The canonical `/bench-what-next` command invokes `bench roadmap --context` once and states that a query error stops the phase without manual reconstruction | Root conformance reports `bench-what-next dropped the roadmap context query`. |
| Canary bite | The three sabotaged surfaces above still trigger their owning messages through the real inner gate | `bench canary` reports that the named fixture did not bite. |

The AXI fixture builder is the one source of the complete populated repository. Tests
mutate that fixture per hostile class rather than pasting multiple miniature context
fixtures. Context completeness is asserted at the public command; package tests cover
only facts that cannot be isolated cheaply at that higher seam. This follows the
gate's real-path rule and keeps the canary as a tripwire rather than a second oracle
([gate real-path rule](../.agents/skills/bench-craft-gate/SKILL.md#L38)).

## Hostile-input ownership

- The AXI seam owns spaces/globs and newlines in paths, missing trailing newlines,
  absent versus empty sources, malformed fragments, control bytes, 0/1/2 exits,
  4095/4096/4097-byte bodies, missing `git`, deep CWD, unchanged state, and repeat-run
  byte identity.
- The wrapper surface owns source, linked by-path, globally installed, symlinked,
  linked-worktree, deep-CWD, stale-PATH, missing-local, and self-resolution cases; it
  also proves `--context --full` remains two arguments rather than a joined string.
- Harness adapters own no query behavior. Conformance pins their single canonical
  phase source; the wrapper contract pins the bytes that all harnesses receive.
- SIGINT needs no separate state-recovery seam: the command creates no repository
  scratch, lease, worktree, cache, or mutation. Forced termination can only discard
  stdout; the before/after read-only assertion covers the relevant invariant.

The shell-CLI hostile classes come from the project profile
([hostile-input checklist](../projects/benchkit.md#L85)).
