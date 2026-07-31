# FT128 agent-line binding

Status: implemented

Decision source: the compiled map
`specs/ft128-agent-line-binding/decisions/multi-harness-line-binding.md` (the schema
decisions), scoped by `ROADMAP.md`'s FT128 entry with FT97 merged in 2026-07-29 (the fork
verdict, the denial bite, and the static guidance-token sweep). The map carries no
`## Sources` entries, so every claim it makes was re-verified against the current tree on
2026-07-31; two re-verification findings changed this draft and are called out below.

## Problem

Three defects share one enforcement surface.

A fork delegation runs on its parent's model and the harness ignores any `model` it
declares, but `check-agent-line` decides from the declared model alone. A fork that
declares `sonnet` while the session runs `opus` passes the guard today — verified against
the pinned tree on 2026-07-31 — which is the silent escalation invariant 2 exists to
block, and the one delegation shape the guard grades backwards.

The denial that fires for every other unbound model leads with Codex model ids and trails
the Claude aliases. Inside a Claude Code session the aliases are the only tokens the Agent
tool can pass, so the recovery instruction names ids that harness cannot use.

Underneath both, `BENCH_TIER_*` is named for the abstract tier but holds one family's ids,
so that family is structurally "the real model" and every other harness reads through a
translation.

## Solution

Tier becomes the only shared identity, resolved through a closed per-harness binding
matrix in which no family is canonical. Every runtime caller names its own harness, so the
guard denies in tokens the asking harness can actually pass. The guard learns the
delegation type and refuses a fork that declares a model it will not run on. Conformance
holds the matrix, the profile's rendering of it, and the kit's own guidance prose to the
same single binding.

### Two re-verification findings that changed this draft

**The fork discriminator exists; the comparison the source asked for does not.** The map
and the roadmap both specify denying a fork "whose declared alias is not the session's own
tier". `tool_input.subagent_type` is a documented Agent input and was captured from a real
`fork` delegation in this repo on 2026-07-31, so the discriminator is pinned. The other
half is not obtainable: the Claude Code hooks reference states that only `SessionStart`
hooks can receive a `model` field, that it is not guaranteed to be present there, and that
no `$CLAUDE_MODEL` environment variable exists — so a `PreToolUse` hook cannot read the
invoking session's tier. Caching a SessionStart `model` is not a substitute: the same
reference records that the field may be omitted and that a mid-session `/model` switch is
re-reported by no hook event, so the cache goes stale silently and a stale cache produces
wrong denials. Story 4 therefore denies on a fact the envelope does carry — see its
implementation decision. This is a reviewer decision, flagged in the approval table.

**The map's SessionStart line report describes a surface that does not exist.** Decision #7
assumes `check-agent-line --describe` prints the binding at session start. No `--describe`
exists anywhere in the runtime; it was retired because SessionStart executed hook script
bodies, and it was replaced by static leading-comment manifests read as data. Session start
runs `bench session-inspect`, whose resume/status/guards output names no tier value at all.
The observed FT97 bite is the deny message, which story 5 fixes. No story re-adds a session
report; a data-only reporting surface is priced in Out of scope.

## User stories

Each `Line:` names the tier's currently bound Codex id, because that is the family
`projects/benchkit.md` declares today. After story 9 lands, a Claude Code session declares
the same tiers as `fable` / `opus` / `sonnet`. Three stories route mid + high, which is not
a row of `craft-line`'s starting table but is the pairing the profile's own cached routings
already use for a falsification pass and for oracle logic; the table sets the start and the
cached routings refine it.

1. As a reviewer, I want Codex and Claude to hold equal, explicitly named columns in the
   binding, so that neither family is structurally the real model. Line: gpt-5.6-terra /
   high. The schema is exact, but it is the parser every verdict reads, so the gate must
   observe it precisely.

2. As an operator, I want every adapter, hook, and CLI entry point to name its own harness
   and resolve through one core, so that the kit CLI, a linked repo's CLI, the guard, and
   the three adapters cannot drift apart. Line: gpt-5.6-luna / medium. Thin wiring over the
   specified core, with focused conformance coverage at a known seam.

3. As an operator, I want an unbound OpenCode column to refuse to launch and to be named by
   `bench doctor`, so that an unadopted harness fails closed instead of guessing. Line:
   gpt-5.6-terra / medium. The launch refusal is mechanical, but the doctor advice is
   operator-facing prose the gate can only grade structurally, so it takes the up-bias.

4. As a Claude Code user, I want the guard to refuse a fork delegation that declares a model
   the harness will ignore, so that a cheap-looking request cannot silently spend the
   session's tier. Line: gpt-5.6-terra / high. The verdict is correctness-critical and its
   residual postures are the guard's whole value.

5. As a Claude Code user, I want a denial to lead with tokens I can actually pass, so that
   the recovery instruction is executable where it is shown. Line: gpt-5.6-terra / high.
   Enforcement stays permissive across the matrix while the advice becomes harness-native,
   and keeping those two apart is the subtle part.

6. As a reviewer, I want the gate to reject an incomplete, malformed, or prose-drifted
   matrix, so that the reviewer-owned binding stays complete and human-readable. Line:
   gpt-5.6-terra / medium. Conformance is the oracle for this static contract and needs a
   mutation-specific bite proof.

7. As an agent author, I want the gate to reject an unbound model literal or a retired
   schema key in kit guidance prose, so that guidance cannot silently resurrect a
   hard-coded binding. Line: gpt-5.6-terra / medium. A new static oracle with an
   all-files quantifier that must be enumerated and proven to bite.

8. As a maintainer migrating a binding, I want `bench doctor` to name the exact hard-cut
   rewrite without touching my file, so that migration is explicit and reviewer-owned.
   Line: gpt-5.6-luna / medium. Settled report-and-offer posture with a black-box CLI
   result, over six retired keys that must all be named.

9. As an agent author, I want `craft-line`, `/bench-setup-repo`, and the profile to declare
   and recommend in the current harness's family, so that the rule that produced this bite
   cannot regenerate it. Line: gpt-5.6-sol / high. Guidance prose compounds through every
   session that loads it, so `craft-line`'s leverage override applies — and per the
   profile's escalation policy this top-tier line pauses for reviewer approval.

## Implementation decisions

- `internal/lines` stays the only deep module. It parses the closed
  `BENCH_<HARNESS>_<TIER>` matrix over the fixed harness set `codex`, `claude`, `opencode`;
  resolves a tier token for a named harness; renders the named harness first; and owns every
  allow, deny, and degraded verdict. The `BENCH_TIER_*` / `BENCH_ALIAS_*` schema is a hard
  cut with no dual read. The alias concept dissolves: the Claude column *is* what the alias
  keys used to hold, so alias projection and its "no bound alias" failure disappear rather
  than being ported.
- Cell grammar is decided per harness rather than inherited from the retired split. Today's
  tier values validate as `modelid.SafeToken` while alias values validate as a bare
  `^[a-z0-9-]+$`; under one matrix every cell validates as `modelid.SafeToken`, which
  already accepts the bare alias shape. OpenCode's provider-qualified requirement stays a
  column rule on its own cells, not a post-filter on a resolved value.
- `--harness codex|claude|opencode` replaces `--alias` and `--provider-model` on
  `resolve-model`, and `check-agent-line` gains the same flag. Both are value-taking flags,
  so the hand-rolled arity switch in the CLI becomes real parsing. `BENCH_MODEL` names
  `top`, `mid`, or `cheap` and never a concrete id, so the reviewer picks the tier and the
  harness picks the family.
- Only *declared* harnesses require all three cells. An absent OpenCode column is not a gate
  failure; its adapter refuses to launch and `bench doctor` names the missing column. No
  OpenCode ids are invented.
- **Story 4's verdict, fixed against what the envelope carries.** The guard reads
  `tool_input.subagent_type`. When it is exactly `fork`, the harness runs the delegation on
  the invoking session's model and ignores any declared model, so: a fork carrying a model
  denies with exit 2, because the declaration is a claim the harness will not honor and the
  guard cannot verify; a fork carrying no model allows with a warning naming that it
  inherits the session's model, because omission is the honest signal and inheritance is the
  documented, unavoidable behavior. Every non-fork path is unchanged, including the existing
  routed-complete deny for an omitted model. Comparison is exact-string on the discriminator
  only: no provider lookup, no model discovery, no session-model cache.
- The fork branch degrades safely by construction. `subagent_type: "fork"` is documented as
  experimental and feature-gated, so in a deployment without fork mode no envelope carries
  the value and the branch never fires; if the value is renamed upstream the exact-string
  match simply stops matching and the guard returns to today's behavior. A missing,
  malformed, or unexpected `subagent_type` never impersonates a fork.
- `ModelFromEnvelope` keeps reading `tool_input.resolvedModel` before `tool_input.model`,
  but its comment is corrected: `resolvedModel` is a `PostToolUse` `tool_response` field,
  not a `PreToolUse` input, so the first read is a defensive fallback and not the documented
  contract it currently claims to be. No behavior changes.
- The denial renderer takes the active harness and uses that harness's three bound tokens as
  its actionable listing. It does not print the rest of the matrix: enforcement stays
  permissive across every bound cell, but advice names only what the asking harness can
  pass, and the full matrix is human-readable in the profile. That removes the map's
  verbose-flag discretion item rather than answering it.
- Command parsing, adapters, and the hook stay thin pass-throughs, each supplying its own
  harness name. No shim parses bindings or reconstructs a verdict. The hook keeps its
  existing fail-open rims; an intentional core deny stays exit 2.
- `checkLineBinding` validates the matrix, declared-harness completeness, and the profile's
  rendering of it from the same parsed binding the runtime uses. Its prose check gains a
  section anchor: today it does a bare whole-file substring search, so any mention anywhere
  in the profile satisfies it. A six-cell matrix rendered as a table needs the check
  anchored to the `Lines` section, or drift hides behind an unrelated paragraph. Every cell
  is cross-checked, not a sample.
- **Named exception to one-source-per-fact.** Keeping the six ids in both `lines.env` and the
  profile is a second authoring of the same knowledge, which the code standard normally
  rejects. It is a reviewer decision closed in the map's #5 — the binding must stay
  human-readable in the profile — and the per-cell cross-check is what makes the duplication
  safe rather than drifting. No other fact in this feature gets a second author.
- `linesEnv()` currently folds an unreadable binding into "absent", so a matrix that fails
  to read takes the unrouted fail-open branch silently. The new parser distinguishes the two
  and the absent-versus-empty contract already pinned in conformance survives the migration.
- One new conformance owner sweeps kit guidance prose. It discovers entries under
  `.agents/commands` and `.agents/skills` at run time, sorts them, traverses every regular
  file, and rejects a non-regular discovered entry before any read. It rejects two token
  classes: a literal matching the tier-model token grammar that names no cell in the parsed
  matrix, and a retired `BENCH_TIER_*` or `BENCH_ALIAS_*` schema key. The diagnostic names
  the file, the token, and the binding source.
- **The swept inventory, stated once.** At spec time the sweep's directories hold ten command
  files (`bench-assess.md`, `bench-debug.md`, `bench-final-check.md`,
  `bench-implement-spec.md`, `bench-review-implementation.md`, `bench-setup-repo.md`,
  `bench-shape-idea.md`, `bench-update-kit.md`, `bench-what-next.md`, `bench-write-spec.md`)
  and twenty-five skill directories. This paragraph is the spec's only statement of that
  inventory; coverage rows refer to it rather than restating the counts. It is a
  specification snapshot, not a production or executable-test registry — run-time directory
  discovery owns the universal set, which is exactly what the newly-discovered-file row
  proves.
- `bench doctor` gains two binding rows: it detects the six retired keys and reports their
  exact rewrites, and it names any known harness whose column is unbound. It never writes
  `lines.env`.
- **No dual read.** A `lines.env` carrying only retired keys resolves as *no binding at all*,
  not as a legacy binding. That is what makes the cut hard rather than a silent
  compatibility shim, and it is the behavior the migration row pins.
- The same green change migrates this repo's own `.bench/lines.env` to `BENCH_CODEX_*` and
  `BENCH_CLAUDE_*` and rewrites the profile's `Lines` section as a matrix table, because
  `checkLineBinding` grades the kit against itself and would otherwise turn the gate red on
  the tree that implements the feature. No concrete model choice changes: the same six
  tokens move to the new keys.

## Testing decisions

- The primary seam is the `internal/lines` exported verdict and resolver surface. Tests pass
  binding bytes and real-shaped envelopes and observe exit code, stdout, and stderr rather
  than private parser state. The captured fork envelope shape — `tool_input` carrying
  `description`, `prompt`, `subagent_type`, and an optional `model` — is the fixture shape,
  taken from a real delegation rather than assumed.
- Two existing hook execution cases drive `tool_input.resolvedModel`, a field a real
  `PreToolUse` envelope never carries. They are proving the guard on a shape the harness
  does not send, so they move to `tool_input.model` as part of story 4's fixture work; the
  defensive `resolvedModel` read stays, but no test may depend on it as the live contract.
- The conformance seam grades a temporary repository through `checkLineBinding` and the new
  guidance-token sweep. Each mutation proof starts from an otherwise clean copy and adds
  exactly one defect, so the target diagnostic cannot be collateral.
- Runtime wiring stays at the existing adapter and hook subprocess seam, where the kit
  wrapper, a linked repo's wrapper, the Claude hook, and all three adapters drive the same
  core. The doctor report is graded as a black-box CLI result in a fixture repo.
- Prior art: `internal/lines/lines_test.go`'s parser, verdict, and resolver tables;
  `checkAgentHookBehavior` and `checkAdapterLineGuards` in
  `internal/conformance/line_routing_exec_test.go`; `checkLineBinding` in
  `line_routing_static_test.go`; and the five `tests/canary/line-routing/` fixtures, whose
  registry entries name `.bench/gate-line-contracts.sh` as a *retired* fragment for
  provenance only — the Go conformance phase owns them, and the gate is asserted never to
  reference that filename again, so there is no shell twin to edit.
- Two existing full-string deny-message assertions and the `line-binding-prose-drift` canary
  are pinned to the retired two-row tier/alias rendering. They move to the matrix rendering
  in the same change as the formatter; the canary carries its binding in a `files/` overlay
  rather than a `MUTATE.json`, so its `lines.env` and profile stub migrate together or its
  EXPECT stops matching.
- The feature gate is `bench gate`; focused unit, conformance, and runtime checks precede it.

### Seam diagram: line verdict and harness resolution

```text
trigger: resolve-model, check-agent-line, or a shift adapter
    |
    v
lines.env + --harness + BENCH_MODEL/envelope --> [ internal/lines ] --> model, message, or exit 0/1/2
                                              [ matrix + verdicts ]
                                                 ^ tests attach here: bytes in, result bytes out
```

### Seam diagram: static matrix and guidance-token conformance

```text
trigger: conformance phase
    |
    v
lines.env + profile Lines section     --> [ checkLineBinding ]      --> targeted diagnostics
.agents/commands + .agents/skills     --> [ guidance-token sweep ]  --> targeted diagnostics
                                            [ sorted regular-file discovery ]
                                               ^ tests attach here: clean temp root, then one mutation
```

### Seam diagram: runtime surfaces and the migration report

```text
trigger: bench shift, an Agent delegation, or bench doctor
    |
    v
adapter/hook/CLI + --harness --> [ subprocess -> same core ] --> launch, deny, or rewrite report
                                    ^ tests attach here: run the real shim, observe exit + stderr
```

### Acceptance coverage map

Red signals are classified honestly. Every seam this feature touches already exists —
`internal/lines`, `checkLineBinding`, the hook and adapter execution checks, and the doctor
CLI — so almost every row is TDD-able today: write the failing case first at the named seam.
Only the guidance-token sweep introduces a new owner, and even there the test is written
before the check. A row that says "already covered" names the existing test.

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | each of the six declared cells resolves for its own harness: `--harness codex` and `--harness claude` each map `BENCH_MODEL` of `top`, `mid`, and `cheap` to that harness's cell | line verdict | TDD-able now: a six-case table in `internal/lines` fails today because a tier token is not a bound model id and no `--harness` exists. | A resolver handling only one tier or only one harness fails the other five cases, so neither a cheap-only nor a Codex-only implementation survives. |
| 1 | a `BENCH_<HARNESS>_<TIER>`-shaped key naming a harness outside the closed set is rejected rather than silently accepted | line verdict | TDD-able now: the case fails today because no key validation exists. | An open-set parser would let a typo'd harness create a phantom column that nothing grades. |
| 1 | a `lines.env` carrying only retired `BENCH_TIER_*`/`BENCH_ALIAS_*` keys resolves as no binding, not as a legacy binding | line verdict | TDD-able now: the case fails today because those keys *are* the live schema. | It pins the hard cut; a dual-read shim would pass every other row while quietly keeping the retired schema alive. |
| 1 | an unreadable `lines.env` is distinguished from an absent one instead of failing open as unrouted | line verdict | TDD-able now: the case fails today because the locator folds both states into "absent". | A corrupt matrix would otherwise silently disable enforcement rather than announce itself. |
| 2 | all six runtime surfaces pass their own `--harness` to the core: the kit CLI, a linked-repo CLI, the Claude hook, and the Codex, Claude, and OpenCode adapters | adapter/hook subprocess | TDD-able now: a per-surface case at the existing execution seam fails today because each surface still uses a retired shape flag or none. | An omitted surface either selects the wrong family or keeps a retired flag, and per-surface cases name which one. |
| 2 | the retired `--alias` and `--provider-model` flags are rejected rather than quietly accepted | line verdict | TDD-able now: both currently exit 0 and resolve a model. | Retiring a flag in prose while the parser still honours it leaves two ways to ask the same question, which is the drift this story removes. |
| 2 | each adapter's resolved model equals what the core returns for the same harness and tier, rather than being recomputed in the shim | adapter/hook subprocess | TDD-able now: the comparison fails today because the Claude adapter's alias projection is computed by a core mode the matrix removes. | A shim that reimplements resolution passes a static "calls resolve-model" check while drifting from the core, so the values must be compared. |
| 2 | an unknown `--harness` value and an unknown or unset tier token each exit 1 | line verdict | TDD-able now: no `--harness` parsing exists, so both cases fail. | Argument validation must not accept an arbitrary family or a concrete model id smuggled in through `BENCH_MODEL`. |
| 3 | the OpenCode adapter refuses to launch while its column is unbound | adapter/hook subprocess | TDD-able now: the refusal is currently a provider-shape rejection rather than an unbound-column one, so a column-scoped case fails. | A fallback to another column would launch the wrong family; only a fail-closed refusal is safe for an unadopted harness. |
| 3 | `bench doctor` names the unbound OpenCode column and the action that would bind it | doctor CLI | TDD-able now: doctor reads no binding at all today, so the case fails on absent output. | Failing closed without saying why leaves an operator with a refusal and no route forward. |
| 3 | an absent OpenCode column produces no conformance diagnostic while Codex and Claude are complete | matrix conformance | TDD-able now at `checkLineBinding`: the case fails until declared-versus-known is distinguished. | Without the distinction an unadopted harness would red the gate, which would make fail-closed unusable. |
| 4 | a fork envelope declaring a bound model denies with exit 2 | line verdict | **Observed red 2026-07-31:** `{"tool_input":{"subagent_type":"fork","model":"sonnet"}}` through the current `check-agent-line` returns exit 0. | This is the backwards verdict the story exists to fix, and an always-allow branch reproduces it exactly. |
| 4 | a fork envelope declaring no model allows with a warning stating that the delegation inherits this session's model | line verdict | **Observed red 2026-07-31:** the same envelope without `model` returns exit 2 under the routed-complete omitted-model deny. | A deny-all-forks implementation fails this positive control, so the two fork branches cannot collapse into one policy. The warning states the inheritance; it cannot name the model, which is unknowable at this hook event. |
| 4 | every bound cell in the matrix is still allowed through the Claude guard, including the Codex column | line verdict | TDD-able now: a six-cell case fails today because only three ids and three aliases are bound. | It pins the map's permissive enforcement across the whole matrix rather than narrowing to the active column. |
| 4 | a non-fork envelope with no model keeps its routed-complete exit 2, and an unrouted or incomplete binding keeps its fail-open rim | line verdict | Already covered by `TestAgentLineVerdict`'s `missing-model-routed-complete-denies`, `missing-model-unrouted-fails-open`, and `missing-model-incomplete-fails-open` cases; the fork field only changes fixture shape. | The fork branch cannot silently widen or replace the residual postures the guard already settled. |
| 4 | a missing, malformed, or unexpected `subagent_type` never impersonates a fork | line verdict | TDD-able now: the case fails once the branch exists but treats arbitrary text as a trusted type. | Hostile envelope text must not reach a trusted comparison; only the exact literal may select the fork branch. |
| 5 | a Claude denial lists the Claude column as parsed from the binding, proven with a fixture whose Claude cells differ from this repo's | line verdict | TDD-able now: the case fails against today's `bound: top=gpt-5.6-sol …` message. | A renderer that hard-codes `fable`/`opus`/`sonnet` passes a live-binding assertion but fails a fixture that binds different cells. |
| 5 | a Codex denial lists the Codex column and names no Claude token | line verdict | TDD-able now: today's message carries both families in every session. | It is the mirror of the Claude row, so an implementation hard-coded to either family fails one of the two. |
| 5 | an unbound token denies with the active harness's recovery instruction | line verdict | TDD-able now: today's message is family-blind. | A plain allow, the wrong family, or a generic message each fails the black-box denial contract. |
| 6 | each of the six declared cells is cross-checked against the profile, one mutation per cell | matrix conformance | TDD-able now at `checkLineBinding`: six per-cell cases fail until the matrix is parsed. | A checker that validates one sampled cell passes a single-cell test while five cells rot; per-cell mutation is the only proof of the quantifier. |
| 6 | a matrix cell named only outside the profile's `Lines` section does not satisfy the prose check | matrix conformance | TDD-able now: the case fails because today's whole-file substring search accepts it. | Without an anchor a six-cell table can rot while an unrelated paragraph keeps the check green. |
| 6 | a declared Claude harness missing its mid cell emits a diagnostic | matrix conformance | TDD-able now: the case fails until completeness is graded per declared harness. | It proves declared-harness completeness rather than mere token syntax. |
| 6 | a malformed matrix token and a missing `lines.env` each emit their own named diagnostic | matrix conformance | Already covered in shape by `TestLineBindingRejectsUnsafeModelTokens` and the absent-input case in `TestRunConformanceDistinguishesAbsentAndEmptyInputs`; both migrate to matrix keys and fail on the retired names until they do. | Safe-token validation and the required-file check are separate ways an unsafe binding could otherwise pass. |
| 7 | every regular file discovered under the sweep's two directories is scanned, including one newly added after the check was written | guidance-token conformance | TDD-able now: a temporary-root case adding a newly named file carrying one unbound token fails until run-time discovery exists. | A hard-coded list or a partial traversal misses the new file, which is why production discovery and not a spec-time list owns the universal set. |
| 7 | one unbound model literal inserted into an otherwise clean copy emits its file-and-token diagnostic, with no other diagnostic present | guidance-token conformance | TDD-able now: the mutation case fails until the sweep exists. | The clean baseline excludes collateral failure and proves the check bites for the intended omission rather than for noise. |
| 7 | each of the six retired schema keys in guidance prose emits its own diagnostic | guidance-token conformance | TDD-able now: the case fails against the two live occurrences in `craft-line`'s skill and `/bench-setup-repo`, which no current check guards. | The hard cut rots back into prose exactly where nothing grades it, which is where the drift already is. |
| 7 | every bound matrix token and ordinary non-model prose token stays accepted | guidance-token conformance | TDD-able now: the clean-fixture case fails while the allowlist is anything other than the parsed binding. | It stops a broad text matcher from turning all guidance red, which is the failure mode a new prose oracle invites. |
| 7 | each of four non-regular entry kinds — FIFO, character device, socket, and symlink — is rejected before it is read | guidance-token conformance | TDD-able now: a case installing a FIFO with no writer must return the discovery diagnostic without blocking, and fails until classification precedes reading. | It proves classification precedes reads per kind, so static inspection neither hangs nor follows bytes outside the swept directories. |
| 8 | all six retired keys are named with their exact rewrites and `lines.env` is left byte-identical | doctor CLI | TDD-able now: doctor emits no binding output, so the case fails on absent text. | Naming five of six leaves a silent survivor, and a byte comparison is what separates report-and-offer from silent mutation. |
| 8 | a second `bench doctor` run reports the same rewrites and still changes nothing | doctor CLI | TDD-able now: the case fails until the report exists. | Re-run idempotency is the property that makes a report safe to leave in a routed repo. |
| 9 | `craft-line`, `/bench-setup-repo`, and the profile carry no retired schema key and name the matrix | guidance-token conformance | Covered by story 7's retired-key row, which fails against the two live occurrences until this story rewrites them. | A rewrite that misses one of the three files leaves its retired key behind and stays red. |
| 9 | whether the guidance actually *recommends* in the current harness's family | guidance-token conformance | **No red signal — stated exception.** The gate can grade that the retired keys are gone and the matrix is named; it cannot grade whether prose gives good advice. | Recorded as residual risk rather than pretended coverage: this is the story's reason for the top-tier line and for reviewer reading, not something the oracle can settle. |
| edge of 1 | quoted, CRLF, last-wins, and no-final-newline matrix entries parse consistently | line verdict | Already covered by `TestTierValue` and `TestTierValueEdgeCases`; their carrier keys migrate to the matrix schema. | A schema migration that drops shell-compatible parsing fails a focused case rather than surfacing at runtime. |
| edge of 1 | an absent `lines.env` and a present-but-empty one stay distinguishable | matrix conformance | Already covered by `TestRunConformanceDistinguishesAbsentAndEmptyInputs`; its pinned diagnostic migrates to the first matrix cell. | A reader that collapses the two states changes either the fail-open rim or the diagnostic operators rely on. |
| edge of 2 | the Claude hook fails open with its warning when the core cannot be resolved | adapter/hook subprocess | TDD-able now, and **not covered today**: no test drives the hook's missing-core rim, so the case is new rather than migrated. | A guard that dies with its core would block every delegation, which is worse than the gap it guards. |
| edge of 2 | each of the three adapters refuses to launch with exit 1 when the core cannot be resolved | adapter/hook subprocess | TDD-able now, and **not covered today**: no test drives the adapters' missing-core rims. | The adapters must fail closed where the hook fails open, and nothing currently proves either direction. |
| edge of 8 | the migration report is the only repeatable in-scope command and performs no binding write | doctor CLI | Covered by the story 8 idempotency row. | A command whose own write changes the fact it reports would falsify its output on the next run. |

Degenerate implementations considered, one per story. A cheap-only or Codex-only resolver
(1) satisfies whichever cell it implements and fails the remaining cases of the six-cell
row. A shim that passes `--harness` but recomputes resolution, or keeps the retired flags
(2), fails the core-equality row and the retired-flag row. A refusal with no operator advice
(3) fails the doctor OpenCode row. An always-allow fork branch (4) reproduces the observed
exit 0 and fails the fork deny row, while a deny-all-forks branch fails the fork allow row.
A renderer that hard-codes `fable`/`opus`/`sonnet` (5) passes a live-binding assertion and
fails both the differing-fixture row and the Codex-denial row. A checker that validates one
sampled cell (6) fails the per-cell mutation row. A sweep driven by a hard-coded file list
(7) fails the newly-discovered-file row, and one emitting a generic diagnostic fails the
clean-baseline row. A generic "your binding is outdated" warning (8) fails the six-key row,
and a doctor that rewrites the file fails the byte-identical row. A rename that leaves the
recommendation rule itself unchanged (9) passes every row the gate owns — that is the stated
exception above, and it is why the story takes the top tier and reviewer reading rather than
a claimed oracle.

### Edge inventory

- Error path — resolved by the unbound OpenCode, unknown harness, unknown tier, and doctor
  rows.
- Empty and absent input — resolved by the absent-versus-empty, unreadable-binding, and
  missing-model rows.
- Boundary values — resolved by the newly-discovered-regular-file row against the swept
  inventory named once in the implementation decisions.
- Malformed input — resolved by the malformed-token, malformed-envelope, and
  hostile-discovery rows.
- Interrupted or partial state — resolved by the declared-harness partial-matrix rows.
- Re-run idempotency — resolved by the doctor no-write and second-run rows.
- Hostile environment — resolved by the quoted/CRLF/no-final-newline parsing, special-file
  discovery, and missing-core rim rows.
- Paths containing spaces or glob characters — **Won't handle**: discovery is rooted in
  tree-controlled kit directories and accepts no caller-supplied path, so no external path
  crosses this seam.
- Control bytes in model text — **Won't handle**: safe-token validation rejects them before
  the renderer or the token sweep sees them, and the malformed-token row owns that rejection.
- Tabs, newlines, and returns a sink permits — **Won't handle**: the model token grammar
  rejects them before a line-oriented diagnostic could receive a value that splits its own
  field.
- A command whose write changes a fact it reports — **Won't handle**: the doctor rows pin
  report-only migration, so this feature performs no binding write whose state it must then
  describe.
- Unquoted multi-word arguments, or a flag value read as a positional — **Won't handle**:
  the fixed `--harness` grammar owns only its named flag, its unknown and missing values
  have their own row, and this feature adds no shell word parser.
- A required tool absent from PATH — resolved by the two missing-core rim rows, which are new
  coverage rather than a claim on existing tests: nothing tests those rims today.
- Invocation through a symlinked wrapper path — **Won't handle**: path resolution belongs to
  the shared wrapper resolver, which this build does not touch; the rim rows above cover the
  outcome that matters, which is what happens when resolution fails.
- A dangling symlink where a swept file is expected — resolved by the special-file discovery
  row, which stats each entry before reading, so a broken link is rejected rather than read
  as an authoritative empty file.
- Destructive worktree state, or SIGINT mid-loop — **Won't handle**: no FT128 behavior
  allocates, leases, or cleans a worktree, so the lifecycle seam keeps those owners.
- Invocation from a nested working directory — **Won't handle**: existing repository-root
  discovery already owns it and this feature adds no cwd-sensitive caller.
- Non-TTY stdin — **Won't handle**: the guard's JSON stdin is deliberately non-interactive,
  so a TTY contract would amputate its only in-scope caller.
- Host-backed filesystem I/O pressure — **Won't handle**: this static parser and conformance
  sweep add no durability or fsync protocol, so the gate-run owner keeps that contract.

## Out of scope

- Binding live OpenCode model identifiers — a harness-adoption capability, not the remainder
  of this one, since the fail-closed behavior is complete without it — 7 edits, 2 gate runs.
- Harness-scoped `bench models` discovery output — advisory discovery is a separate
  capability from binding enforcement and rendering — 5 edits, 2 gate runs.
- A session-start report of the live binding — the map's decision #7 surface, which does not
  exist; re-adding an executed `--describe` would reintroduce a closed red finding, so this
  is a new data-only reporting subcommand rather than a restoration — 6 edits, 2 gate runs.
- Detecting a fork's actual model after the fact from `PostToolUse` `tool_response`
  `resolvedModel` — an observability capability at a different hook event that cannot deny,
  so it neither completes nor substitutes for story 4 — 6 edits, 2 gate runs.
- Codex parity for agent-line denial — blocked by the current non-deny-capable spawn event —
  0 sound repository edits, 0 gate runs until Codex exposes a deny-capable spawn event.
