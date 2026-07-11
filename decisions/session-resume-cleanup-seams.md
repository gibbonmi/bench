# Session resume cleanup: seam and envelope research

## Result

The durable ledger needs a new `internal/intent` package. It is neither worktree
cleanup nor rendering: shift, manual-worktree, Claude Agent, Stop, and status paths
all need the same persistence/reconciliation semantics. The deletion test therefore
puts storage, atomicity, objective safety, and proof-of-done filtering behind one
small interface; deleting that module would otherwise spread those rules across
five callers. `internal/worktree` retains cleanup and `internal/status` retains
composition ([seam standard](../.agents/skills/bench-craft-seams/SKILL.md#L15)).

The ledger address is the repository's **git common directory**, not the caller's
absolute git directory. Gate-cache code deliberately writes the latter
([writer](../internal/gate/gate.go#L187), [reader](../internal/status/status.go#L170));
a linked worktree therefore gets `.git/worktrees/<name>/bench-last-gate`. The live
probe below showed main, Claude-agent, and Bench-pool worktrees have distinct
absolute git dirs but the same common dir. `internal/worktree` already uses this
topology to recover the canonical checkout ([classifier](../internal/worktree/classifier.go#L39)).

## Hook envelopes and harness reach

Claude Code's Agent `PreToolUse` input carries common fields plus `tool_name`,
`tool_input`, and `tool_use_id`; the Agent input contains the description, prompt,
subagent type, requested model, and requested isolation. The official hook schema
does not attach the later `agent_id`, worktree path, or branch to this event
([Claude hooks](https://code.claude.com/docs/en/hooks#pretooluse-input)). A live
Claude Code 2.1.207 capture on 2026-07-11 matched that shape. Its `tool_use_id`
(`toolu_...`) was not a worktree/agent identifier. Later events split the missing
facts: `SubagentStart` carries `agent_id`/`agent_type` but no objective, while
`WorktreeCreate` carries only a generated `name`; neither shares the PreToolUse
correlation key ([Claude hooks](https://code.claude.com/docs/en/hooks#subagentstart-input),
[worktree hook](https://code.claude.com/docs/en/hooks#worktreecreate-input)). Under
parallel spawns, temporal joining would be nondeterministic. The agent-line writer
must therefore record **uncorrelated Claude intent**, then let later git evidence
resolve it; it must not invent a branch association.

Codex cannot reuse this writer. Its `PreToolUse` surface covers Bash, patch, and MCP
calls, while delegation is exposed as `SubagentStart`; that event carries
`agent_id`/`agent_type` but not the task objective and cannot stop the subagent
([Codex hooks](https://learn.chatgpt.com/docs/hooks#pretooluse),
[SubagentStart](https://learn.chatgpt.com/docs/hooks#subagentstart)). The repository
therefore wires `check-agent-line` only in Claude
([Claude config](../.claude/settings.json#L20), [Codex config](../.codex/hooks.json#L14),
[hook-layer decision](../.bench/BENCH-reference.md#L125)). OpenCode and headless
adapters have no project lifecycle-hook config; their harness-independent writer is
`bench shift`, whose objective is already available before its worktree starts
([shift](../internal/shift/loop.go#L70), [adapter contract](../.bench/BENCH-reference.md#L105)).

Both Claude and Codex Stop inputs carry `cwd`, session identity, and
`stop_hook_active`; Codex additionally documents `turn_id` and both document the
last assistant message ([Claude Stop](https://code.claude.com/docs/en/hooks#stop-input),
[Codex Stop](https://learn.chatgpt.com/docs/hooks#stop)). That is sufficient to
resolve the repository/common dir and refresh the ledger. It is not a recovery
source: Stop fires when a response finishes, not when the process is killed. The
existing shim already passes the full envelope to the Go core
([shim](../.bench/hooks/stop.sh#L47)), while the core currently reads only
`stop_hook_active` and returns before gate work when unarmed
([stophook](../internal/stophook/stophook.go#L34),
[run](../internal/stophook/stophook.go#L81)). Ledger refresh belongs before that
armed-gate branch and remains best-effort; it must never change the Stop guard's
fail-open rim or forge a gate verdict.

### Reproducible envelope/common-dir probes

The live Agent-envelope check can be repeated without editing a tracked hook by
passing Claude a temporary settings file (requires `claude` and `jq`):

```sh
tmp="$(mktemp -d)"
capture="$tmp/agent.json"
jq -n --arg cmd "tee '$capture' >/dev/null" \
  '{hooks:{PreToolUse:[{matcher:"Agent",hooks:[{type:"command",command:$cmd}]}]}}' \
  > "$tmp/settings.json"
claude --settings "$tmp/settings.json" -p \
  'Spawn one general-purpose agent with worktree isolation to run pwd only.'
jq '{hook_event_name,tool_name,tool_use_id,tool_input,
     has_agent_id:has("agent_id"),has_worktree_path:has("worktree_path")}' \
  "$capture"
```

The address edge is directly runnable against any registered linked worktree:

```sh
git worktree list --porcelain
git -C <worktree-path> rev-parse --absolute-git-dir
git -C <worktree-path> rev-parse --path-format=absolute --git-common-dir
```

On this repository the main, a `.claude/worktrees/agent-*` checkout, and a
`~/.bench/worktrees/*` checkout returned three different absolute git dirs and the
same `/home/mgibs/workspace/bench/.git` common dir.

## Package ownership

- `internal/intent` owns common-dir addressing, the ledger format, locked atomic
  replacement, entry upsert, safe one-line objective encoding, and a typed snapshot
  that filters entries only on worktree removal or the shared landed proof. The file
  contains intent only; git facts and gate/spec state stay derived elsewhere.
- `internal/worktree` already separates branch sweeping from out-of-pool removal
  ([clean entry](../internal/worktree/clean.go#L14),
  [sweep](../internal/worktree/clean.go#L39),
  [removal](../internal/worktree/clean.go#L83)). Add a sibling conservative primitive:
  reuse `LandedInDefault`, remove only clean out-of-pool registrations, and retain
  dirty, detached, locked, or leased entries. Do not thread an `auto` boolean through
  the salvage path. The classifier must retain porcelain `locked` state instead of
  discarding it ([classifier](../internal/worktree/classifier.go#L25)); Claude locks
  live agent worktrees specifically to prevent concurrent cleanup
  ([Claude worktrees](https://code.claude.com/docs/en/worktrees#clean-up-worktrees)).
- `internal/git` owns the typed dirty-path, upstream-ahead, and unique-local-branch
  facts. `internal/status` consumes those facts and `intent.Snapshot`; it owns only
  severity/order/rendering, consistent with its existing structured `Signals` seam
  ([status](../internal/status/status.go#L72)).
- `internal/shift` replaces its worktree-local `.bench-objective` as the durable
  source with an intent write before the first adapter run. The current scratch file
  is created only after acquire/branch creation and is removed during teardown
  ([write](../internal/shift/loop.go#L97), [teardown](../internal/shift/loop.go#L323));
  it is not a cross-session handoff.

## Cleanup and status contracts

Manual `bench worktree clean` is deliberately salvage-capable: it stages and commits
dirty branch-backed WIP before removal, and refuses dirty detached state
([manual removal](../internal/worktree/clean.go#L83)). Auto-clean may call the existing
orphan sweep because that sweep deletes only `worktree-*` branches absent from all
registered worktrees and proven landed by ancestry or patch containment; git failure
keeps rather than deletes ([branch inventory](../internal/worktree/branches.go#L12),
[landed proof](../internal/worktree/branches.go#L45)). The new worktree-removal half
must be separate and strictly clean/unlocked. It returns typed counts so the adapter
can render at most one plain line, suppressed when all counts are zero:

```text
bench resume: cleaned <worktrees> worktree(s), <branches> landed branch(es); kept <dirty> dirty, <locked> locked, <leased> leased; <intent> open intent(s)
```

This line is SessionStart developer context, not an AXI query. The hook already emits
plain CLI/status/guard context ([session-start](../.bench/hooks/session-start.sh#L29));
TOON would create a second contract on an operational adapter. The report prints
before `bench status`, so status observes post-clean state.

Status already owns one ascending ladder and a five-row default budget with `--all`
overflow ([signals](../internal/status/status.go#L72),
[renderer](../internal/status/status.go#L124)). Keep the existing ranks: severity 1
`git` becomes the consolidated landed-state verdict (dirty path count, unpushed
commit count, branches with unique commits); severity 2 `worktree` carries unresolved
or uncorrelated intent alongside kept worktree state. Aggregate each class into one
signal rather than spending a row per ledger entry. The runtime contract already
pins the five-row cutoff and `+N more` escape hatch
([budget test](../internal/contract/runtime/runtime_status_test.go#L242)) and pins
gate < git plus worktree < guards < drain ordering
([gate/git ladder](../internal/contract/runtime/runtime_status_test.go#L417),
[worktree ladder](../internal/contract/runtime/runtime_status_test.go#L33)).

## Gate attachment and red signals

| Seam | Real-path assertion | Distinct red signal |
|---|---|---|
| Intent package | Main/Claude/pool roots share one address; absent/empty/malformed/no-final-newline ledgers, hostile objectives, concurrent upserts, interrupted replacement, and proof-only removal are deterministic | Focused `internal/intent` test names the address, parse, atomicity, objective byte, or lifecycle case. |
| Conservative worktree cleanup | Real CLI fixture contains clean, dirty, locked, leased, detached, landed, unique-patch, and missing-default cases; only clean/unlocked worktrees and landed orphans disappear; HEAD/index/dirty bytes are unchanged; second run is a no-op | Runtime contract reports `session-start conservative cleanup contract failed` plus the wrongly removed or surviving path/ref. Existing manual-clean contracts remain unchanged ([worktree contracts](../internal/contract/runtime/runtime_worktree_test.go#L15)). |
| Claude intent writer | Real `check-agent-line.sh` receives allowed, denied, malformed, and hostile Agent envelopes from main and linked CWDs; only allowed calls append one safely encoded uncorrelated entry in the common ledger | Line-routing exec test reports `check-agent-line intent capture contract failed`; existing allow/deny exit assertions still own model enforcement ([line exec](../internal/conformance/line_routing_exec_test.go#L12)). |
| Stop refresh | Real `stop.sh` receives representative Claude and Codex envelopes while armed/unarmed/active; refresh is idempotent and existing 0/2 plus gate-cache behavior is unchanged | Runtime gate contract reports `stop hook intent refresh contract failed`; existing cache/active tests remain the guard regression signal ([Stop contracts](../internal/contract/runtime/runtime_gate_test.go#L12)). |
| SessionStart adapter | Real hook auto-cleans, emits zero or one report line before post-clean status, then guard brief; outside-repo and missing-core paths still exit 0 silently | AXI/runtime contract reports `session-start resume cleanup contract failed` ([current hook contract](../internal/contract/axi/axi_guards_test.go#L190)). |
| Status composition | Seed ledger plus dirty, ahead, unique, worktree, gate, and six total signals; assert exact aggregate counts, severity 0/1/2 ordering, five default rows, and full `--all` | Runtime status contract reports `status landed-state contract failed`; the behavior canary expects that substring rather than reimplementing counts. |
| Harness wiring | Claude keeps Agent/Stop/SessionStart; Codex keeps Stop and gains SessionStart; no false Codex Agent claim; OpenCode/headless remain explicitly unwired | Root conformance emits `claude settings.json ...` or `codex hooks.json SessionStart event does not run .bench/hooks/session-start.sh`; a targeted unwired canary proves each new diagnostic bites ([current Claude wiring test](../internal/conformance/line_routing_static_test.go#L199)). |

The gate should exercise these public paths, not scan implementations. Canaries
sabotage the SessionStart invocation, common-dir address, and status aggregation and
expect the owning contract's message; they do not carry their own ledger parser or
landedness rule ([gate discipline](../.agents/skills/bench-craft-gate/SKILL.md#L38)).
