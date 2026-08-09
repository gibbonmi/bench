# FT173 AXI contract and migration sequence

Status: ready

## Destination

One shared AXI owner that states all ten published principles, provides one
derivation for each reusable output fact, and migrates Bench surfaces in
independently reviewable slices. The migration retains `bench diff` as the
coherent Git-inspection owner and extends full AXI behavior to every `bench spec
build` operation so lifecycle success, refusal, and recovery output is structured,
deterministic, and directly actionable. It supplies the AXI predecessor consumed
by the spec-build review-and-gate cadence without redefining gate authority,
gate-result evidence, or capture accounting.

## #1: What is the current ten-principle implementation census?

Blocked by: none
Type: Research

### Question

Re-read the current published AXI contract, `craft-cli`, command registry, output
helpers, and contract tests. Record every principle, its current documented and
implemented owner, all independent derivations, uncovered surfaces, and the exact
compatibility or gate assertion that constrains consolidation in
`decisions/byte-preserving-axi-foundation/assets/ft173-axi-surface-census.md`.

### Answer

The current implementation is a set of strong local owners rather than one
ten-principle contract. `internal/toon` owns flat-table bytes, `internal/usage`
owns the 23 shared-parser call sites, SessionStart owns ambient status, and spec
build already owns ordered operation grammar, structured success/refusal, explicit
empty status, and one lifecycle-next derivation. Four independent truncation
policies and command-local aggregate, empty-state, schema, and next-action
renderers remain. There are no production `help[]` emitters, no registry metadata
for AXI scope or help, and no standalone ten-principle gate matrix.

The gate's existing byte pins are fragmented command tests. Its intended
registry-wide `subcommand-routing` assertion still scans the former dispatcher
shape instead of `commandRegistry`/`Command.Run`, so it does not currently close
the help/grammar inventory. The complete census, current counts, exact owners,
and compatibility assertions are recorded in
`decisions/byte-preserving-axi-foundation/assets/ft173-axi-surface-census.md`.

## #2: What contextual help does real command usage require?

Blocked by: none
Type: Research

### Question

Build an exhaustive command-by-command inventory from the named Claude and Codex
usage logs plus the current CLI registry. For each command, record success, empty,
refusal, stale-state, and recovery cases; the state already known to the command;
the exact useful next action; and whether a contextual `help[]` row would have
prevented another turn. Name any unavailable log source rather than substituting
sampled commands. Store the evidence in
`decisions/byte-preserving-axi-foundation/assets/ft173-command-help-inventory.md`.

### Answer

The complete production surface is larger than either current advertisement:
the Go registry has 48 root names, the wrapper adds no-argument help and
`repair`, and the nested grammars expand worktree, spec, spec-build, release,
release-preflight, gate-go, and hook operations. `bench commands --brief` cannot
serve as the inventory owner because it hard-codes only `version`, itself, and
`status`.

Every declared call is inventoried in
`decisions/byte-preserving-axi-foundation/assets/ft173-command-help-inventory.md`, including unobserved public
and plumbing surfaces. The available Claude corpus supplies the decisive
wrong-remedy trace: `bench spec build start` knew exact-green evidence was absent
but suggested plain `bench gate` where only `bench gate --fresh` could satisfy
the precondition. The inventory therefore distinguishes exact state-derived
actions, plan/apply actions that still require a second call, terminal results
that should not invent busywork, and plumbing whose caller owns continuation.

The promised reviewer-named Codex usage-log set is unavailable: no manifest,
session-id list, or archive was supplied. The locally available Codex project
executions are recorded only as corroboration, never as a sampled substitute or
as evidence that an unmatched command is unused.

## #3: Which helper and byte-compatibility seams already exist?

Blocked by: none
Type: Research

### Question

Trace truncation, aggregate, structured success, structured error, empty-state,
exit-code, and contextual-help derivations to their current owners and consumers.
Identify which moves can preserve bytes, which necessarily change public output,
which consumers parse that output, and which exact mutation or paired-delta probe
proves each shared helper still bites. Store the result in
`decisions/byte-preserving-axi-foundation/assets/ft173-helper-compatibility-census.md`.

### Answer

Bench already has strong shared mechanics for exact flat-table TOON bytes and common
argv grammar, while semantic totals, completeness, empty-state meaning, and next
actions remain correctly owned by their command domains. The four truncation policies
are incompatible as policies: they differ in caps, units, escaping, metadata,
`--full`, and—in worktree cleanup—authority-bearing safety and fingerprint behavior.
They can route through one byte-preserving call only if it is a parameterized
projection mechanic that leaves those facts with their current owners.

No production consumer decodes Bench's TOON stdout. In-repo composition uses typed
values or forwards bytes: handoff reads status signals, dashboard gathers typed facts,
spec build renders service values, and SessionStart forwards status output. Tests and
external agents are still public byte consumers, so the foundation must pair baseline
and candidate stdout, stderr, exit, argv, and default/`--full` results exactly and
independently mutate every cap, unit, total, unknown, empty, sink, exit, and action
derivation.

Typed carriers, registry metadata, parameterized bounds, and owner-supplied aggregate
rendering can remain byte-preserving behind current renderers. `help[]`, changed
actions or refusals, new family homes/flags/schemas, coherent Git snapshots, and the
full spec-build AXI envelope necessarily change output and require separately approved
paired deltas. The complete owner/consumer ledger and mutation matrix are recorded in
`decisions/byte-preserving-axi-foundation/assets/ft173-helper-compatibility-census.md`.

## #4: Which CLI surfaces must satisfy the complete AXI contract?

Blocked by: none
Type: Grill

### Question

Does FT173 retain the old query-only scope, widen every operational command, or
make a bounded operational exception for the lifecycle family that now depends on
agent-actable refusals and recovery?

### Answer

FT173 remains the one AXI owner for the complete Bench CLI, with high-frequency
query surfaces migrated under its existing priority model. In addition, every
operation in the complete `bench spec build` family must satisfy all ten AXI
principles. That bounded operational expansion includes shared success and error
helpers, definitive empty states, honest exit codes, idempotent retries where the
operation is retryable, and contextual `help[]` rows populated with the known
slug, run, candidate, assignment, and exact next command.

This decision does not silently widen every other operational command. Each such
family still requires its own evidence and reviewer decision. The spec-build
cadence consumes FT173's owners and schemas rather than deriving lifecycle-only
formats or remediation prose.

## #5: Who owns gate result structure inside lifecycle output?

Blocked by: none
Type: Grill

### Question

Should FT173 define promotion's gate-result payload while making spec-build output
AXI-complete, or compose an independently owned result?

### Answer

FT185 owns the reusable gate-phase result payload shared by `bench gate`, `bench
commit`, and `bench spec build promote`. FT173 owns the surrounding AXI lifecycle
envelope and composes the FT185 payload without cloning its phase fields, verdict
meaning, or evidence derivation. Neither item changes which command may run a gate
or author project-green authority.

## #6: What is the Git-inspection command surface?

Blocked by: none
Type: Grill

### Question

Should FT173 extend `bench diff` or introduce another Git-oriented command family?

### Answer

Extend `bench diff`. It remains the single coherent owner for review-base
resolution, checkout facts, changed-file inventory, revision and aggregate facts,
whitespace status, landed-commit log, and exact patch bodies. It reuses the
existing Git-facts and diff-range owners, includes untracked regular-file bodies
under `--full`, and refuses or retries a snapshot whose HEAD, index, or worktree
moves mid-read. FT173 does not add `bench git` or a second porcelain parser.

## #7: How should the enlarged FT173 work be sliced and ordered?

Blocked by: #1, #2, #3, #4
Type: Grill

### Question

Does the old foundation → contextual disclosure → Git-inspection sequence remain
three independently reviewed specs, or does the full-AXI `bench spec build`
migration require its own independently useful spec? Choose the smallest ordering
that keeps byte-preserving consolidation separate from public-output migrations,
avoids rewriting pinned contracts twice, and makes the cadence predecessor
explicitly shippable.

### Answer

FT173 ships as eight independently reviewed specs. The byte-preserving foundation
is five ordered capabilities: an exact compatibility oracle; shared carriers and
registry declarations; outcome/action migration; bounded-projection migration;
and aggregate/empty migration. Each capability is independently promotable, emits
no `help[]`, and changes no existing command bytes, streams, exits, or accepted
argv.

After that foundation, two sibling public-output migrations may proceed
independently. One makes every `bench spec build` operation AXI-complete and emits
its contextual actions in the same atomic output migration; it composes FT185's
gate-result payload and is the explicitly shippable predecessor for the spec-build
review-and-gate cadence. The other makes `bench diff` the coherent Git snapshot
decided in #6 and emits its contextual actions in that same atomic output migration.
The spec-build slice is prioritized when FT185 is available; the Git slice does not
block it.

The eighth spec completes contextual disclosure across the remaining approved
surfaces after their output schemas are stable. It consumes the exhaustive command
inventory, leaves terminal and caller-owned plumbing without invented busywork, and
does not rewrite the already-migrated spec-build or diff contracts. This dependency
graph gives each of those two wide surfaces one public-output rewrite while retaining
separate compatibility review for every byte-changing slice.

## #8: What compatibility contract governs byte-changing output migrations?

Blocked by: #2, #3
Type: Grill

### Question

For contextual help, coherent diff snapshots, and full-AXI spec-build responses,
which existing stdout shapes remain supported, which receive a versioned or
additive transition, and which may change atomically? Base the ruling on observed
consumers and failure modes rather than preserving accidental prose by default.

### Answer

The byte-preserving foundation retains exact stdout, stderr, exit, accepted-argv,
default/`--full`, and empty-state behavior for every existing fixture. It introduces
no public schema version because it introduces no public output change.

Contextual disclosure on a surface whose primary response is otherwise stable is an
additive transition. The existing primary bytes, stream, exit, and argv behavior remain
exact, and one typed `help[]` block is appended after that response. An empty help set
is valid for a terminal result or caller-owned plumbing; compatibility does not justify
inventing a low-value command. Existing prose may remain visible during this slice, but
it is not a second action owner: the typed action derivation supplies `help[]` and any
later removal of redundant prose requires another explicit compatibility decision.

The coherent `bench diff` response and the complete spec-build family each receive one
atomic output migration under their existing command names. They do not ship a
`--legacy` flag, dual renderer, or parallel versioned schema: no production consumer
decodes current Bench TOON, and maintaining two public derivations would create
unwarranted compatibility sediment. `bench diff --commit` retains its post-landing
meaning. Every spec-build operation migrates together so the family never exposes a
mixed old/new contract.

Those atomic migrations preserve durable lifecycle and Git facts, command identity,
gate and landing authority, and the 0/1/2 taxonomy: success/help/no-op is 0, an
unsatisfied intent is 1, and invalid argv is 2. A slice may change schemas, block order,
error stream, accepted new argv, or the spec-build family home's former usage/2 result
only where its reviewed old-to-new fixture names that exact delta. Everything outside
the approved delta remains byte-equal. A future machine decoder must bind to an
explicit schema version when it is introduced; its hypothetical existence does not
force a legacy mode now.

## #9: What proportional proof makes each AXI slice authoritative?

Blocked by: #1, #2, #3, #7, #8
Type: Grill

### Question

Which behaviors require contract conformance, which are better proved by exact
unit or paired-delta fixtures, and which independent mutations demonstrate that
errors, empty states, help actions, truncation, aggregates, coherent Git snapshots,
and spec-build lifecycle responses fail red without expanding every command into
the full conformance matrix?

### Answer

Proof attaches at the smallest owner that can independently falsify the behavior.
Shared TOON, grammar, bounds, aggregate carriers, empty-state, error, and action
mechanics receive focused unit tests and one independent mutation per owned fact. The
byte-preserving foundation also runs exact baseline/candidate paired deltas over the
affected argv matrix, comparing stdout, stderr, exit, and accepted arguments. It does
not expand every command into a duplicate conformance suite.

One registry-derived AXI conformance check covers only the surfaces and reusable
principles declared by the registry. It proves inventory completeness, required output
and empty/error/help shapes, help spelling, deep-cwd routing where applicable, and that
unknown flags fail with usage/2. Command-domain semantics remain in their focused
owners rather than being re-derived by conformance.

Contextual actions are derived from the same typed success, refusal, or precondition
result that knows the next authorized transition. Every known fixed argument is carried
forward; only genuinely unknown future input remains a placeholder. A terminal result
or caller-owned plumbing may emit no action. A focused owner test must turn red when a
useful action is removed, a known value becomes a placeholder, an unknown value is
guessed, a fixed flag is dropped, a stale fingerprint is carried, or prose is advertised
as an executable command. The exhaustive usage inventory judges whether the derived
action is useful; it is evidence, not a second action derivation.

The coherent-diff slice is proved against raw Git over committed, staged, unstaged,
untracked, rename, deletion, binary, hostile-filename, clean, and mid-read-drift
fixtures. The spec-build slice uses an in-process operation/state matrix for lifecycle
semantics and exact old-to-new response deltas across success, no-op, empty, refusal,
stale, recovery, plan/apply, review, and promotion states. An exact Bench executable is
required only for an acceptance row that observes wrapper routing, executable identity,
environment, signals/teardown, or installed/stripped behavior.

The final remaining-contextual-disclosure spec begins with an open-ended review of the
Codex and Claude harness session logs accumulated across FT173 and representative recent
Bench work. The review looks for CLI leverage rather than a closed keyword checklist.
Its required examples are useful additions or corrections to a command's `help[]`,
repeated shell or tool-call sequences that one coherent Bench query could replace, and
output transformations agents repeatedly perform themselves—such as `head` or `tail`—
that belong in a bounded default, `--full` escape hatch, pre-computed aggregate, or other
CLI-owned projection. The resulting evidence asset gives every observed opportunity one
disposition: fold into the final spec, already owned, decline with a reason, or route to
a named roadmap item when it is a separate capability or would reopen a closed decision.
The asset informs the spec and its coverage rows; it does not become a second permanent
CLI inventory.

Ticket-local checks compile and run only their changed code slice and declared
integration surfaces. They run only conformance owners and canary fixture-owner
mutations affected by that slice. No ticket pays the whole-project gate; the completed
composition receives that sole prospective gate through spec-build promotion.

## #10: Is FT173 ready for spec authoring?

Blocked by: #1, #2, #3, #7, #8, #9
Type: Task

### Question

Verify the three evidence assets against the current tree, record the chosen spec
sequence and compatibility boundary, and confirm that each slice has independently
useful acceptance behavior, a named red signal, an ownership fence, and a
proportional mutation probe before handing the first slice to spec authoring.

### Answer

The map is ready for the first spec. The three evidence assets were rechecked at
`974020e4af8de5ed75098c4c5934a8907952bb2b` against an unchanged implementation
tree: 48 registry roots, 43 non-test TOON table calls, 23 shared usage-parser calls,
zero production `help[]` emitters, and four current truncation owners still match the
censuses. Each source is a regular file and `bench maps` reports no FT173 source or
schema diagnostic.

Each approved slice is independently useful and falsifiable:

- **Compatibility oracle:** acceptance is one pinned baseline/candidate matrix over
  every production member and argv class, comparing stdout, stderr, exit, and
  acceptance. Its fence is the manifest, case derivation, and paired executor.
- **Carriers and registry:** acceptance is one complete ten-principle owner plus typed
  outcome, action, projection, aggregate, empty, and registry mechanics with inert
  metadata. Its fence is the shared owner, production declarations, conformance, and
  exact guidance advertisements.
- **Outcome/action migration:** acceptance is production reachability for every declared
  outcome/action route with domain kinds, exits, actions, and bytes preserved. Its fence
  is the existing command families plus route conformance and contraction.
- **Bounded-projection migration:** acceptance is production reachability for all four
  owner-specific policies without cap, unit, total, full-mode, or authority movement.
  Its fence is sanitize, roadmap, worktree, outline, conformance, and contraction.
- **Aggregate/empty migration:** acceptance is production reachability for every ordered
  aggregate and exact empty class without inferred totals, unknown coercion, or empty
  normalization. Its fence is the named semantic owners, registry, conformance, and
  contraction.
- **Full-AXI spec build:** acceptance is one consistent ten-principle response contract
  across every operation and lifecycle state, including exact contextual actions. Its
  fence is the lifecycle envelope and action rendering; it composes but does not
  re-derive FT185's gate payload or change lifecycle/gate authority. Any operation/state
  matrix or approved old-to-new delta mismatch is red, including a missing known value,
  guessed placeholder, wrong exit, stale action, or mixed old/new operation.
- **Coherent diff:** acceptance is one compact orientation call and at most one full-body
  call over one stable Git snapshot. Its fence is the existing diff owner composed with
  existing Git-facts and range owners; it adds no second porcelain parser. A mismatch
  against raw Git, omitted untracked regular-file body, or accepted mid-read splice is
  red, with one independent mutation per hostile-state class.
- **Remaining contextual disclosure:** acceptance is registry-complete additive
  `help[]` on every approved stable surface where the command inventory names a useful
  next action, and an honest empty set for terminal or caller-owned results. Its fence is
  the typed action owner, renderer, registry declaration, and remaining approved
  surfaces; it does not rewrite spec-build or diff output or widen other operational
  families to full AXI. Removing a useful row, losing a fixed argument, guessing an
  unknown value, carrying stale authority, or advertising prose as a command is red. Its
  spec starts from the harness-log opportunity asset required by #9 and accounts for
  every finding through an acceptance row or an explicit disposition.

The five byte-preserving specs land in the order listed above. The full-AXI spec-build
and coherent diff atomic surface migrations follow them, with spec build prioritized once
FT185 is available. Remaining contextual disclosure is the eighth and final capstone
after those output schemas stabilize. Spec writers retain only the bounded discretion
already listed below.

## Not yet specified

## Spec-writer discretion

- Internal helper type and package names after ticket #3 identifies one durable
  owner; emitted-byte compatibility and public field meaning are not discretionary.
- Ticket decomposition within the reviewer-approved spec sequence, provided each
  ticket is an independently green tracer outcome rather than a thematic file
  bundle.
- Exact fixture names and focused test filters within ticket #9's decided proof
  boundary.

## Out of scope

- Implementing the spec-build review, checkpoint, candidate-pin, promotion, or
  post-landing capture cadence decided in the successor map.
- Changing which commands run the whole-project gate, changing gate verdict
  authority, or treating semantic review as deterministic evidence.
- Re-deriving FT185's gate-result payload or FT130's capture-accounting contract.
- Adding a `bench git` namespace or another Git porcelain parser.
- Reopening the terminal single-build serial-gate spec or lifecycle.
- Implementing, gating, committing, or landing FT173 during this shaping pass.

## Sources

- Path: `ROADMAP.md`
  Supports: #1 through #3 and #6 through #9 current FT173 evidence, closed constraints, and the prior three-spec sequence replaced by #7.
  Drift: re-read after FT173, FT175, or FT185 is reshaped, staged, or retired.
- Path: `.agents/skills/bench-craft-cli/SKILL.md`
  Supports: #1 current seven-principle project guidance and query-surface posture.
  Drift: re-read after any AXI guidance edit.
- Path: `decisions/byte-preserving-axi-foundation/assets/ft173-axi-surface-census.md`
  Supports: #1 and #7 through #10 current ten-principle implementation owners, independent derivations, uncovered surfaces, and exact compatibility or gate assertions.
  Drift: re-run after any command-registry, TOON, usage, truncation, output-helper, spec-build renderer, or AXI contract change.
- Path: `decisions/byte-preserving-axi-foundation/assets/ft173-command-help-inventory.md`
  Supports: #2 and #7 through #10 exhaustive current CLI and nested-operation contextual-help inventory, usage-log boundary, state-carrying next actions, and wrong-remedy evidence.
  Drift: re-run after any command-registry, wrapper routing, nested grammar, action/next renderer, or named usage-log corpus change.
- Path: `decisions/byte-preserving-axi-foundation/assets/ft173-helper-compatibility-census.md`
  Supports: #3 and #7 through #10 current helper owners and consumers, byte-preserving versus output-changing moves, runtime parser census, and exact paired-delta or mutation obligations.
  Drift: re-run after any TOON, usage, truncation, aggregate, empty-state, error/exit, contextual-action, renderer, or command-consumer change.
- Path: `decisions/spec-build-review-gate-cadence.md`
  Supports: #4 and #5 reviewer-approved full AXI expansion for the spec-build family and FT185 composition boundary.
  Drift: re-read if cadence tickets #6, #7, or #9 reopen.
