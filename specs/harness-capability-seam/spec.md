# Harness capability seam

Status: staged

Roadmap: FT239

Decision source: named reviewed artifact `roadmap/FT239.md`, plus the reviewer's 2026-08-26 closure that the record is a Go table in the core and FT239 no longer waits on FT222

Verification log: 1 iteration to accept — the reviewer-named session reviewed the spec and the opus slices in one round and folded three partials: the HC10 split, the prefix-reader shape in ticket 03, and the roadmap closure out of ticket 09

## Problem

Five readers enumerate harnesses on their own today. `lines.Harnesses` closes
the binding matrix. `status.harnessPrefix` translates phase actions.
`guards.wiredHarnesses` scans two config files it names itself. The
conformance `harnessOf` map and its literal loops grade the adapters. The Hook
Layers prose holds the Codex agent-line verdict as one dated sentence.

A new harness therefore lands as five edits. A skill that needs a mechanic
reads prose or guesses.

The route rejects `opencode` and has no model-free form. So the Claude-outage
path and the shell-only path exist only as prose. The parity property that
proves one core covers one verb, `resolve-model`, across six surfaces. The
other plumbing verbs the shims call have no such proof, and the CI entry has
none.

## Solution

One package, `internal/harnesses`, owns a versioned record with one row per
harness: `codex`, `claude`, `opencode`, and `none`. Each row carries the
providers it binds, the phase invocation form, and the hook config with its
wired events. It also carries the deny-capable delegation verdict, the
headless adapter, and the roadmap's mechanics as cells. Each `yes` or `no`
cell names its source and its check date. Every cell the tree does not record starts as `unknown`.

`lines`, `status`, `handoff`, `guards`, and the conformance package derive
their harness lists from the record. `bench harnesses` projects the record as
TOON. `bench status --route --harness none` routes past every phase action, so
a shell operator receives a runnable command or an honest empty cell.

Two conformance checks join the gate. `harness-record` grades the record
against the tree: every shipped adapter and hook config has a row, and every
row's claims match the files. `entry-point-parity` holds every shim, adapter,
front door, and the CI script to one registry name, and it compares each
shim's result to the direct call.

## User stories

### The record is one queryable source

Line: opus / high.

Four packages derive from this seam, and the oracle's checks read it. The
scorecard routes a foundational Go-seam rewrite to Opus at high.

1. As an agent, I want one Go table to name every harness Bench knows, so that no package keeps its own list.
2. As an agent, I want each row to carry the providers the harness binds, so that the parser knows which cells need a qualified id.
3. As an agent, I want each row to carry the harness's phase invocation form, so that the route translates from one table.
4. As an agent, I want each row to carry the hook config path and its wired events, so that the guard report reads one list.
5. As an agent, I want each row to state whether a delegation event can deny, so that the FT24 verdict is data.
6. As an agent, I want each row to name the headless adapter, so that the shift loop's entry is recorded.
7. As an agent, I want the roadmap's mechanics as cells that hold `yes`, `no`, or `unknown`, so that no skill infers a mechanic.
8. As a reviewer, I want each `yes` or `no` cell to carry a source and a date, so that a claim names what was read.
9. As a reviewer, I want a `none` row for the model-free path, so that the degraded operation is explicit.
10. As a reviewer, I want the record to carry a schema version, so that a consumer names the shape it read.
11. As an agent, I want `lines.Harnesses` to be the record's binding harnesses, so that the matrix stays closed from one source.
12. As an agent, I want the provider-qualified cell rule to derive from the row's providers, so that no literal names a harness.
13. As an agent, I want the route's prefix table to derive from the record, so that a new harness lands as one row.
14. As an agent, I want the guard wiring reader to derive its config list from the record, so that `wired` never misses a harness.
15. As an agent, I want the conformance harness loops to derive from the record, so that an adapter check covers a new row.
16. As a reviewer, I want a unit test to red an empty cell or an out-of-enum value, so that the record cannot rot silently.

### The route reads the record

Line: opus / medium.

The route composes the record through the existing selection owner, so the
cached mid routing for status work applies.

17. As a shell operator, I want `bench status --route --harness none` to print the first invocable non-phase command, so that I get a runnable next step without a harness.
18. As an agent on OpenCode, I want `--harness opencode` to route past every phase action, so that a dead key is never recommended.
19. As a shell operator, I want a phase-only board under `--harness none` to print an empty command cell, so that the answer is honest.
20. As a reviewer, I want one board to lead with the same state and why under every harness, so that only the command form differs.
21. As an agent, I want the usage line to advertise the record's harnesses, so that the grammar and the record agree.
22. As a shell operator, I want `bench handoff --harness none` to write a shell command under `## Next command`, so that a cold session without a harness resumes.
23. As an agent, I want an unknown `--harness` value to stay a usage error, so that a typo never routes.

### The record is queryable

Line: opus / medium.

The verb composes the record and the existing TOON emitter, so the cached mid
routing applies.

24. As an agent, I want `bench harnesses` to print one TOON row per harness, so that the matrix is one call away.
25. As an agent, I want `bench harnesses <harness>` to print every cell with its source and date, so that one row's detail is one call away.
26. As an agent, I want an unknown harness argument to be a usage error, so that a typo never reads as an empty state.
27. As a reviewer, I want the verb in the AXI registry, the profile seam, and the craft-cli table, so that the registry check passes.
28. As an agent, I want the wrapper to route `harnesses` to the core, so that every shipped CLI reaches one implementation.

### The matrix check grades the record against the tree

Line: opus / high.

This is oracle logic, and a wrong check is the worst class of bug in the kit.

29. As a reviewer, I want a shipped adapter with no matching row to red the gate, so that every supported runtime has an explicit row.
30. As a reviewer, I want a shipped hook config with no matching row to red the gate, so that a wired harness is never unrecorded.
31. As a reviewer, I want a declared hook event that its config does not wire to red the gate, so that the record cannot overclaim.
32. As a reviewer, I want a wired `.bench/hooks/` script the row omits to red the gate, so that the record cannot underclaim.
33. As a reviewer, I want a row that names an absent headless adapter to red the gate, so that the entry is real.
34. As a reviewer, I want a `delegation_guard` cell that contradicts the `check-agent-line` wiring to red the gate, so that verdict and config agree.
35. As a reviewer, I want a special file or symlink at a config path refused and named, so that the check never follows a link.
36. As a reviewer, I want the check to own one canary family with a planted red, so that the tripwire proves it bites.

### Every entry point reaches one core

Line: opus / high.

This is oracle logic that runs subprocesses, and the existing line-routing
exec checks are its prior art.

37. As a reviewer, I want one parity table to map every shim and adapter to a registry name, so that an unlisted entry is red.
38. As a reviewer, I want each shim, run through a stub wrapper under `BENCH_COMMAND_OBSERVE=1`, to print `command-registry:<name>`, so that the shim reaches the registry.
39. As a reviewer, I want each shim's exit code and stdout tail to match the direct call, so that no shim holds a second opinion.
40. As a reviewer, I want the wrapper's bare invocation and `bench status --route` to print the same output, so that both front doors agree.
41. As a reviewer, I want `scripts/release-preflight.sh` to exec `release-preflight` verbatim and the workflow to run that script, so that CI reaches the registry.
42. As a reviewer, I want the front-door phase file to name `bench status --route` exactly, so that both phase adapters route through the core.
43. As a reviewer, I want every plumbing verb reached by a parity row or exempt with a reason, so that none escapes the table.
44. As a reviewer, I want the parity check to own one canary family with a planted red, so that the tripwire proves it bites.

### The docs name the record

Line: opus / high.

The reference and the profile are guidance prose, so the leverage override
routes them mid and high.

45. As a cold session, I want Hook Layers to name `bench harnesses` as the verdict's source, so that the prose does not restate a cell.
46. As a cold session, I want the profile to list the verb and the two checks, so that the profile advertises the gate's shape.
47. As a reviewer, I want the changelog to record the verb and the `none` route, so that the release notes are complete.

## Implementation decisions

**One package owns the record.** `internal/harnesses` exports the schema
version, the ordered rows, and a lookup by name. The row order is `codex`,
`claude`, `opencode`, `none`, so every existing diagnostic keeps its wording.

A row holds the providers as a closed value: one provider name, `any`, or
`none`. It holds the phase form as a string that may be empty. It holds the
hook config path and the wired events as one list. It holds the deny-capable
delegation verdict, the headless adapter path, and the mechanics cells.

The mechanics are the roadmap's twelve:

- steering during an active turn
- structured user questions
- tool-permission controls
- hooks
- MCP support
- subagent support
- subagent isolation
- effort selection
- persistent tasks
- resume and recovery
- structured output and exit status
- headless execution

A cell holds `yes`, `no`, or `unknown`. A `yes` or `no` cell carries a source
and an ISO date. An `unknown` cell carries neither. The initial cells come
only from facts the tree records. Those facts are the hook configs, the
adapters, the Hook Layers verdict with its 2026-07-11 date, and the
reference's effort rule. Every other cell starts as `unknown`, and a later
reviewed edit fills it.

**Four derivations collapse.** `lines.Harnesses` becomes the record's rows
whose providers are not `none`. `lines.CellFault`'s provider-qualified rule
reads the row's providers instead of a harness literal. `status` builds its
prefix table from the record's phase forms at init, and `HarnessChoices` and
`ValidHarness` read the record. `guards.wiredHarnesses` reads each row's hook
config path. The conformance `harnessOf` map and its literal loops read the
record.

The `harness-prefix-single-source` check moves its owner constants to the
record's file, so the handoff package still holds no literal.

**The route treats a missing form as a dead key.** A phase action is
invocable for a harness only when that harness has a phase form. `Route`
skips a phase action for a formless harness the way it skips prose today. A
board with only phase actions therefore renders an empty command cell for
`none` and `opencode`. The lead's state and why never depend on the harness.

**`bench harnesses` is an AXI query.** Bare, it prints
`schema: <n>` and `harnesses[N]{harness,provider,phase_form,hooks,delegation_guard,headless,checked}`.
With one harness argument, it prints `cells[N]{field,value,source,checked}`
for that row. Any other argument is a usage error at exit 2. The verb joins
`approvedAXIQueries`, the profile's AXI seam bullet, and the craft-cli
disclosure table with a terminal `help[0]` disposition. It also joins the
wrapper's case labels with a `subcommand-routing` entry.

**Two conformance checks join the registry.** `harness-record` runs at the
dev tier over the root. It enumerates `.bench/adapters/*` and the shipped
hook configs from disk, and each must map to a row. For each row it grades
the hook config against the declared events in both directions. It grades
the adapter path against disk and the `delegation_guard` cell against the
`check-agent-line` wiring. It classifies every path with the no-follow
classifier before it reads, and it names a refused path.

`entry-point-parity` runs at the dev tier over the root. Its table maps each
shim and adapter basename to a registry name and a canned input. It
enumerates the shims and the adapters from disk, and an entry outside the
table is a diagnostic.

For each runtime row it runs the entry through the stub wrapper with
`BENCH_COMMAND_OBSERVE=1`, and it runs the direct verb with the same input.
It requires the observed `command-registry:<name>` line, equal exit codes,
and the direct stdout as a suffix of the shim's stdout.

Three rows are static. The `worktree-lifecycle.sh` row names `worktree-hook`
and runs nothing, because a create needs a live pool. The CI row grades the
script's exec line and the workflow's run line. The front-door row grades the
phase file's exact verb.

Every internal-inventory registry command is reached by a row or carries an
exemption reason in the check's own table. Both checks register in
`conformanceChecks`, the registry package, the profile's conformance table,
and one canary family each under `tests/canary/`.

**The docs point at the record.** The Hook Layers bullet on the agent-line
guard keeps its one-clause why and names `bench harnesses codex` as the
verdict's source. The Files section's adapter bullet names the record. The
profile's AXI seam bullet, its conformance table, and the changelog gain
their rows.

## Testing decisions

The highest seam that shows each failure is the record package's data and
the conformance check functions over a synthetic root. The record's unit test
walks every row and every cell. The derivation tests extend the existing unit
tests in `lines`, `status`, `handoff`, and `guards`. The check tests follow
`TestHarnessPrefixSingleSourcedBites` and `checkAdapterLineGuards`: a
synthetic root, one mutation per diagnostic, and a subprocess through the
stub wrapper. The gate's `test` phase runs the registered checks over the
live root.

### Seam diagram

    trigger: bench status --route / bench harnesses / the gate's test phase
        │
        ▼
    harness name ──▶ [ internal/harnesses: rows, cells, lookup ] ──▶ row
                          │                                   ◀ tests attach here: walk every row and cell
                          ├──▶ lines.Harnesses, CellFault
                          ├──▶ status prefix table, HarnessChoices, Route
                          ├──▶ guards.wiredHarnesses
                          └──▶ conformance: harness-record, entry-point-parity
                                    ◀ tests attach here: synthetic root, one mutation per diagnostic
    shim + canned input ──▶ [ stub wrapper, BENCH_COMMAND_OBSERVE=1 ] ──▶ stderr id, stdout, exit
                                    ◀ tests attach here: compare with the direct verb

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| HC01 | 1 | The record's rows are exactly `codex`, `claude`, `opencode`, and `none` in that order. | record unit | A missing or reordered row changes every diagnostic. |
| HC02 | 2 | The `codex` row binds `openai`, the `claude` row binds `anthropic`, the `opencode` row binds `any`, and the `none` row binds `none`. | record unit | A wrong provider mis-rules a cell. |
| HC03 | 3 | The `claude` phase form is `/bench-`, the `codex` form is `$bench-`, and the other two are empty. | record unit | A wrong form renders a dead key. |
| HC04 | 4 | The `claude` row names `.claude/settings.json` with six events, and the `codex` row names `.codex/hooks.json` with three. | record unit | A dropped event hides a wired hook. |
| HC05 | 5 | Only the `claude` row holds `delegation_guard: yes`. | record unit | A `yes` on codex contradicts FT24. |
| HC06 | 6 | Each of the three harness rows names `.bench/adapters/<harness>`, and `none` names no adapter. | record unit | A wrong path fails the matrix check. |
| HC07 | 7 | Every row holds every one of the twelve mechanics with a value in the closed enum. | record unit | A missing cell reads as a zero value. |
| HC08 | 8 | Every `yes` or `no` cell carries a non-empty source and a date that parses as ISO. | record unit | An unsourced claim is inference. |
| HC09 | 9 | The `none` row has an empty phase form, no hook config, no adapter, and `headless execution: no`. | record unit | A `none` row with a form is not model-free. |
| HC10 | 10 | `harnesses.Schema` equals 1. | record unit | A consumer cannot name the shape. |
| HC11 | 11 | `lines.Harnesses` equals `codex`, `claude`, `opencode`. | lines unit | A second list drifts. |
| HC12 | 12 | `lines.CellFault("opencode", "gpt-5")` reports the provider-qualified fault and `CellFault("codex", "gpt-5")` reports none. | lines unit | A dropped rule accepts a bare opencode id. |
| HC13 | 13 | `status.HarnessChoices()` returns the four record names with `claude` first and the rest sorted. | status unit | A stale choice list rejects a record harness. |
| HC14 | 14 | `bench guards` reports `claude,codex` for a script both configs name and `none` for a script neither names. | guards unit | A dropped config path hides a wiring. |
| HC15 | 15 | `checkAdapterLineGuards` on a root missing `.bench/adapters/opencode` names that adapter. | conformance unit | A literal loop dropped by the derivation skips it. |
| HC16 | 16 | The record unit test fails on a row with an empty cell or a value outside the enum. | record unit with a mutated copy | A permissive walk lets the record rot. |
| HC17 | 17 | A board with a phase signal above a `git push` signal routes `--harness none` to `git push`. | route unit | A translation that keeps the phase leads with a dead key. |
| HC18 | 18 | The same board routes `--harness opencode` to `git push`. | route unit | An opencode form invented from another column recommends a phase. |
| HC19 | 19 | A board with only phase signals routes `--harness none` to the first signal with an empty command and `NoCommand`. | route unit | A fallback command claims a phase the operator cannot run. |
| HC20 | 20 | The lead's `state` and `why` are equal across all four harnesses for one board. | route unit | A re-ranked board per harness disagrees. |
| HC21 | 21 | `bench status -h` advertises `--harness` with the four record names, `claude` first. | status command unit | A stale help pin names the old grammar. |
| HC22 | 22 | `bench handoff --harness none` on a board led by `git push` writes `git push` under `## Next command`. | handoff unit | A handoff that rejects `none` cannot pin a shell resume. |
| HC23 | 23 | `bench status --route --harness cursor` prints the grammar and exits 2. | status command unit | A permissive parse routes a typo. |
| HC24 | 24 | `bench harnesses` prints `schema: 1` and then `harnesses[4]{harness,provider,phase_form,hooks,delegation_guard,headless,checked}` with four rows. | verb unit | A projection that skips `none` hides the degraded path. |
| HC25 | 25 | `bench harnesses codex` prints `cells[13]{field,value,source,checked}` with the `delegation_guard` source naming the Codex hooks docs. | verb unit | A detail view without sources restates prose. |
| HC26 | 26 | `bench harnesses cursor` prints the usage line and exits 2. | verb unit | An unknown name rendered as an empty table reads as a definitive empty state. |
| HC27 | 27 | `checkAXIQueryRegistry` over the live root reports no diagnostic. | conformance over the live root | A verb missing from any of the three seams reds the registry check. |
| HC28 | 28 | `bin/bench.sh harnesses` prints the same output as the direct verb. | entry-point-parity row | A wrapper without the label reaches the default case with no routing entry. |
| HC29 | 29 | A root with `.bench/adapters/cursor` and no `cursor` row yields a diagnostic naming the adapter. | conformance unit | A check that walks only the record never sees a new adapter. |
| HC30 | 30 | A root with `.cursor/hooks.json` naming a `.bench/hooks/` script and no `cursor` row yields a diagnostic naming the config. | conformance unit | A check that walks only the record never sees a new config. |
| HC31 | 31 | A `claude` config missing the `Stop` wiring yields a diagnostic naming `Stop`. | conformance unit | The record overclaims a hook. |
| HC32 | 32 | A `codex` config that wires `check-agent-line.sh` yields a diagnostic naming the script. | conformance unit | The record underclaims a hook. |
| HC33 | 33 | A root without `.bench/adapters/codex` yields a diagnostic naming the row's adapter path. | conformance unit | A recorded entry that does not exist. |
| HC34 | 34 | A `claude` config without `check-agent-line.sh` yields a diagnostic naming `delegation_guard`. | conformance unit | The verdict and the wiring disagree. |
| HC35 | 35 | A FIFO at `.codex/hooks.json` yields a diagnostic naming that path and no other diagnostic for that row. | conformance unit | A plain read blocks or reads through. |
| HC36 | 36 | The `harness-record` canary fixture turns the check red, and the restored fixture turns it green. | canary mutation test | A family without a bite is rot. |
| HC37 | 37 | A root with `.bench/hooks/extra.sh` outside the table yields a diagnostic naming `extra.sh`. | conformance unit | A shim outside the table has no parity proof. |
| HC38 | 38 | `block-dangerous-git.sh` with a benign envelope through the stub wrapper prints `command-registry:guard-git` on stderr. | conformance unit with subprocess | A shim that decides alone prints no id. |
| HC39 | 39 | `stop.sh` with a `stop_hook_active` envelope exits like the direct `stop-verdict` and its stdout ends with the direct stdout. | conformance unit with subprocess | A shim with a second opinion differs. |
| HC40 | 40 | `bin/bench.sh` with no argument prints the same bytes as `bin/bench.sh status --route`. | conformance unit with subprocess | The two front doors diverge. |
| HC41 | 41 | A `scripts/release-preflight.sh` whose exec line names `release-preflight-2` yields a diagnostic. | conformance unit | CI reaches a verb the registry does not own. |
| HC42 | 42 | A `.agents/commands/bench.md` that names `bench status --routes` yields a diagnostic. | conformance unit | The front door routes through a dead verb. |
| HC43 | 43 | A registry command with `internalInventory` that no row reaches and no exemption names yields a diagnostic naming the command. | conformance unit | A new plumbing verb escapes the table. |
| HC44 | 44 | The `entry-point-parity` canary fixture turns the check red, and the restored fixture turns it green. | canary mutation test | A family without a bite is rot. |
| HC45 | 45 | The Hook Layers agent-line bullet contains `bench harnesses codex`. | anchors registry test | A dropped pointer leaves the verdict as prose. |
| HC46 | 46 | The profile's conformance table lists `harness-record` and `entry-point-parity`. | conformance-meta over the live root | The advertised table and the registry disagree. |
| HC47 | 47 | `CHANGELOG.md` names `bench harnesses` under Unreleased. | anchors registry test | The release notes omit the verb. |
| HC48 | 38 | `session-start.sh` in a temp repo prints `command-registry:session-inspect` and its stdout ends with the direct stdout. | conformance unit with subprocess | A shim that renders its own dashboard prints no id. |
| HC49 | 38 | Each of the three adapters under `BENCH_MODEL=mid` prints `command-registry:resolve-model`. | conformance unit with subprocess | An adapter that recomputes the line prints no id. |
| HC50 | 20 | `bench status --route --harness codex` on a board led by a phase prints `$bench-` in the command cell. | status command unit | A dropped translation prints the canonical form. |

### Edge inventory

- `--harness` without `--route` stays a usage error at exit 2.
- `--route --all` stays a usage error at exit 2.
- A `lines.env` key `BENCH_NONE_MID=` is a foreign key, because `none` binds no provider.
- A `lines.env` that binds only the opencode column is complete for opencode and unadopted for the rest.
- An absent `.codex/hooks.json` with a `codex` row that declares three events is a diagnostic, and an empty file is a second, distinct diagnostic.
- A `.claude/settings.json` that references a hook through `$CLAUDE_PROJECT_DIR` and one through `${CLAUDE_PROJECT_DIR}` both count as wired, as today's wiring check accepts.
- A dangling symlink at an adapter path is absent, and a live symlink is refused and named.
- A shim run with no wrapper on PATH takes its own rim, and the parity check does not grade that state.
- A shim's stderr carries the observed id and its own warnings, and the check reads only the id line.
- The `worktree-lifecycle.sh` row is static: it names `worktree-hook` and runs nothing.
- A record cell whose `checked` date is in the future is a unit-test failure.
- The `harnesses` verb takes at most one positional argument, and two arguments are a usage error.

**Won't handle** a runtime CI row through a staged binary — the script builds its own executable, and the static exec line proves the route.

**Won't handle** an OpenCode phase adapter — the record shows the empty form, and a `$bench` equivalent is a separate capability.

**Won't handle** Codex agent-line guard parity — FT24 stays parked, and the record's `no` cell names the upstream reason.

**Won't handle** `bench models`' provider source list — discovery sources are not harnesses, and the list stays where it is.

**Won't handle** filling `unknown` cells from upstream docs — the build records only what the tree already records, and a reviewed edit fills the rest.

## Ownership fences

- `internal/harnesses/`
- `internal/lines/`
- `internal/status/`
- `internal/handoff/`
- `internal/guards/`
- `internal/conformance/`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry_data_test.go`
- `cmd/bench/main.go`
- `cmd/bench/main_test.go`
- `cmd/bench/command_registry_test.go`
- `bin/bench.sh`
- `tests/canary/harness-record/`
- `tests/canary/entry-point-parity/`
- `.bench/BENCH-reference.md`
- `.agents/skills/bench-craft-cli/SKILL.md`
- `projects/benchkit.md`
- `CHANGELOG.md`
- `ROADMAP.md`
- `roadmap/FT222.md`
- `roadmap/FT239.md`
- `specs/harness-capability-seam/`

The record package lands first. The four derivations follow, then the route
and the verb, then the two checks, then the docs.

## Out of scope

- A runtime CI parity row through a staged `dist/bench-preflight`: 3 edits, 1 gate run.
- An OpenCode front door that routes with `--harness opencode`: 6 edits, 2 gate runs.
- A reviewed edit that fills the `unknown` cells from named upstream docs: 1 edit per cell, 1 gate run per batch.
- `bench models` sources derived from the record's providers: 2 edits, 1 gate run.
- A `bench harnesses --check` that grades a linked repo's configs outside the kit gate: 4 edits, 1 gate run.

## Further notes

The reviewer closed one fork on 2026-08-26: the record is a Go table in the
core, and FT239 no longer waits on FT222. Three choices are the author's and
stand for veto. The record keys on harness with a providers value, so a
provider × harness pair is a row plus a value rather than a second row. The
`none` row is the model-free path, and `opencode` routes like `none` until a
front door exists. The parity comparison uses stdout as a suffix, because
`session-start.sh` prints one advisory line before the core's output.

The word "capability" names a host facility in the gate (`capability-skips`,
`BENCH_REQUIRE_CAPABILITIES`). This spec names the harness data "the record"
and its rows "harnesses", so the two never share a term.
