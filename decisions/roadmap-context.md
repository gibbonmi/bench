## Destination

Define a fast, read-only `bench roadmap --context` evidence snapshot that gives
`/bench-what-next` the same complete local inputs in every harness.

## #1: What authority does the context command have?

Blocked by: none
Type: Grill

### Question

Should the command gather evidence only, or also reconcile and mutate the roadmap?

### Answer

It is read-only. It emits complete source evidence and mechanical derivations;
`/bench-what-next` retains every removal, rewording, prioritization, and drain verdict.

## #2: Where does harness portability live?

Blocked by: none
Type: Grill

### Question

Should output adapt to each harness or remain one portable contract?

### Answer

The CLI emits byte-identical AXI/TOON across Claude Code, Codex, OpenCode, and
headless shifts. Harness adapters only invoke the command; they never change its
schema or evidence.

## #3: What evidence and output shape does the command own?

Blocked by: #1, #2
Type: Grill

### Question

Should the snapshot embed source documents or normalize their evidence, and may it
research external triggers?

### Answer

Emit normalized, lossless typed sections for roadmap rows and sequence, capture
entries, structure findings, spec/history state, git state, cached gate state, and
promotion-record state. Preserve raw text when parsing fails. Stay offline and mark
rows whose graduation requires external evidence instead of researching them.

## #4: What are the latency, size, and failure postures?

Blocked by: #3
Type: Grill

### Question

May the command run expensive verification, how does it handle pathological bodies,
and may the phase reconstruct evidence after an error?

### Answer

Use cheap local reads and derivations only; never run a cold gate or canary sweep.
Keep the common case complete in one call, but truncate oversized individual bodies
with an explicit byte count and `--full` escape hatch. Emit structured AXI errors and
fail closed; `/bench-what-next` does not reconstruct the snapshot manually.

## #5: Which existing modules own every snapshot fact?

Blocked by: #1, #3, #4
Type: Research

### Question

Trace every proposed field to its current Go owner, distinguish reusable APIs from
facts trapped behind renderers, and identify the smallest refactors that avoid a
second parser or derivation. Record the inventory and citations in a short research
asset.

### Answer

The context snapshot belongs in `internal/roadmap`, extending its existing roadmap,
sequence, ideas, and drain seam. It composes typed facts from `internal/learnings`,
`internal/structure`, `internal/spec`, `internal/git`, and `status.GateVerdict`; each
current renderer must project from those facts rather than have context parse its
stdout. Promotion linkage stays local to the context aggregation until a second
consumer earns a separate seam. The cited ownership inventory and minimal refactors
are recorded in [roadmap-context-ownership.md](roadmap-context-ownership.md).

## #6: How does the command reach every shipped invocation surface?

Blocked by: #2, #5
Type: Research

### Question

Trace real-kit, linked-repo by-path, global-wrapper, deep-CWD, hook, and headless
adapter routing for Claude Code, Codex, and OpenCode. Identify any surface where the
same command would resolve a different binary, root, stdout, or exit code, with a
runnable probe for each claimed equivalence.

### Answer

Phase discovery is already single-sourced through `.agents/commands`: Claude Code's
command adapter, Codex's explicit skill, and other AGENTS.md harnesses all reach that
file. Query execution is harness-neutral only when it uses the linked/source wrapper.
Real-kit, linked by-path, deep-CWD, shared-resolver, and linked-worktree probes produced
identical stdout and exit codes; the worktree wrapper correctly used the main checkout's
binary. Bare `bench` remains PATH-owned and can select different bytes. The cited surface
trace and runnable probes are recorded in
[roadmap-context-routing.md](roadmap-context-routing.md).

## #7: Where does the gate attach without duplicating the contract?

Blocked by: #5, #6, #8
Type: Research

### Question

Map the AXI contract checks, linked-surface behavior tests, hostile-input classes,
and canary needed to prove the command stays complete, harness-neutral, offline, and
read-only. Name the observable red signal for each seam.

### Answer

Attach the gate at three existing seams: typed package tests for fact ownership and
truncation, one real-command AXI contract for the complete/read-only/offline snapshot,
and the linked-surface routing contract for repo-local wrapper forwarding. Add three
tripwires that sabotage those public checks without reimplementing them: context
completeness and wrapper forwarding as behavior-owned fixtures, plus the
`/bench-what-next` invocation anchor as a conformance fixture. Exact schemas, the
4096-byte per-body ceiling, hostile-input ownership, and every observable red signal
are recorded in [roadmap-context-gate.md](roadmap-context-gate.md).

## #8: Which layer pins the repo-local CLI version?

Blocked by: #6
Type: Grill

### Question

Should `/bench-what-next` resolve the linked/source wrapper before invoking context, or
should the top-level `bench` wrapper forward every command to a repo-local wrapper when
one exists?

### Answer

The top-level `bench` wrapper forwards to the repo-local wrapper whenever one is
present, with self-resolution preventing recursion. This makes bare
`bench roadmap --context` select the linked repo's CLI version for harnesses and
humans alike; `/bench-what-next` needs no resolver-specific shell ceremony.

## Not yet specified

- None.

## Out of scope

- Mutating or prioritizing `ROADMAP.md`, `IDEAS.md`, or the learnings journal.
- Network or provider probes for upstream graduation triggers.
- Harness-specific schemas, rendering branches, or manual fallback gathering.
- Running a fresh gate, canary sweep, or other cold verification command.

## Handoff

1. **Module boundaries.** `internal/roadmap` owns the typed roadmap document,
   `ContextSnapshot`, aggregation, source-state/raw fallback, size policy, and AXI
   rendering. `internal/learnings`, `internal/structure`, `internal/spec`, and
   `internal/git` expose their typed facts; `status.GateVerdict` remains the sole gate
   cache reader. The top-level shell wrapper owns repo-local forwarding. The canonical
   `/bench-what-next` command consumes the query and retains all judgment.
2. **Contracts.** `bench roadmap --context [--full]` is read-only and offline, emits
   the fixed schema-1 TOON blocks in the gate asset, uses exit 0 with empty stderr on
   success, structured stdout/exit 1 on unsatisfied reads or derivations, and stdout
   usage with exit 2 for bad arguments. `--full` removes only the 4096-byte body ceiling.
   The top-level wrapper forwards the entire argv to a distinct repo-local wrapper and
   otherwise uses its own normal binary resolution without recursion.
3. **Deep vs thin.** Roadmap context is deep: it hides cross-owner gathering,
   completeness, ordering, truncation, and rendering. Typed owner APIs are deep
   extensions of existing engines. Harness adapters and CLI dispatch stay thin and
   gain no independent schema, parser, resolver, or test seam.
4. **Black-box assertables.** Fixed TOON blocks/fields/order, definitive empties,
   malformed raw evidence, truncation counts and `--full`, structured errors and
   0/1/2 exits, empty stderr, no repository/git/cache changes, no network/provider/gate
   sentinel calls, repeated byte identity, exact wrapper bytes/exits across invocation
   surfaces, and one canonical phase invocation are all externally observable.
5. **Gate attachment.** Owner package tests cover typed facts; the AXI contract covers
   the real context command; the surface contract covers wrapper forwarding; root
   conformance covers phase consumption. Three targeted canaries keep the new checks
   live. No new phase, duplicate registry, or manual-only seam is required.
6. **Hostile-input owners.** AXI owns spaces/globs/control bytes in evidence, newline-less
   and absent/empty/malformed files, missing git, boundary-size bodies, deep CWD,
   read-only state, and idempotent bytes. Wrapper routing owns missing global tools,
   symlinks, source/linked/global/worktree surfaces, multiword argv, PATH shadowing,
   missing local wrappers, and self-recursion. SIGINT has no persistent state to clean.
7. **Uncertainty flags.** None; exact schemas, ordering, error posture, ceiling, routing,
   and oracle seams are fixed in the three research assets.
8. **Rejected alternatives.** Do not mutate or verdict the roadmap, parse human CLI
   output, compose lossy dashboard/status snapshots, create harness-specific schemas,
   inject context through hooks, probe providers, run cold verification, teach the phase
   resolver ceremony, or add a second completeness registry.
9. **Domain watch-outs.** TOON refuses several control bytes; fail closed without
   partial output. Git-derived text and paths are hostile bytes. Linked worktrees lack
   ignored binaries and must re-anchor to the main checkout. A wrapper must compare its
   resolved path before forwarding or it can recurse forever. Read-only includes the
   index, gate cache, and untracked files, not only tracked Markdown.

Dependency order: expand typed owner APIs first; build the context aggregator and AXI
contract second; add wrapper forwarding and its surface contract third; wire the
canonical phase anchor and all canaries last. These are recommended green build slices,
not separate product capabilities.
