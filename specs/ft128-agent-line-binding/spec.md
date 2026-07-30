# FT128 agent-line binding

Status: staged

## Problem

The agent-line guard accepts the requested model on a context-forked delegation even
though that delegation runs on its parent's model. Its denial also leads with tokens
that the active harness cannot pass, and command guidance can silently hard-code a
retired model binding.

## Solution

Make tier the only shared identity and resolve it through a closed per-harness binding
matrix. Subject to a reproducible real-envelope evidence prerequisite, the Claude guard
will include the delegation discriminator and authoritative inherited session tier in
its verdict, present Claude tokens first, and reject a context fork whose declared model
does not match that session tier. The conformance phase will reject every unbound
tier-model token in the portable command set.

Compiled decision map: `specs/ft128-agent-line-binding/decisions/multi-harness-line-binding.md`.

## User stories

1. As a reviewer, I want a tier binding that gives Codex and Claude equal, explicit
   columns, so that neither family is treated as the real model. Line:
   gpt-5.6-terra / high. The closed schema is exact, but the core verdict and migration
   are oracle behavior that the gate must observe precisely.

2. As a Claude Code user, I want the delegation guard to reject a context fork whose
   declared model differs from its inherited session tier, so that a cheap-looking
   request cannot spend the parent’s top-tier line. Line: gpt-5.6-terra / high. The
   real hook-envelope contract and the deny path are correctness-critical; their JSON
   spelling remains blocked on captured evidence.

3. As an operator, I want every adapter and hook to resolve the current harness’s tier
   through one core, so that the kit CLI, linked CLI, hook, and adapters cannot drift.
   Line: gpt-5.6-luna / medium. This is thin wiring over the specified core contract
   with focused conformance coverage.

4. As a reviewer, I want the gate to reject an incomplete, malformed, or prose-drifted
   harness matrix, so that the reviewer-owned binding remains complete and readable.
   Line: gpt-5.6-terra / medium. Conformance is the oracle for this static contract and
   needs a mutation-specific bite proof.

5. As an agent author, I want each portable command’s literal tier-model token checked
   against the binding, so that guidance cannot silently resurrect a hard-coded model.
   Line: gpt-5.6-terra / medium. This is a new static oracle with an all-command
   quantifier that must be enumerated and proven to bite.

6. As a Claude Code user, I want denial and SessionStart advice to lead with the
   current harness’s tokens, so that the recovery instruction is executable where it is
   shown. Line: gpt-5.6-sol / high. The harness-native recommendation rule changes
   reusable guidance and therefore takes the kit’s leverage line.

7. As a maintainer migrating a binding, I want `bench doctor` to name the hard-cut
   rewrite without changing my file, so that migration is explicit and reviewer-owned.
   Line: gpt-5.6-terra / medium. The report-and-offer posture is settled and has a
   black-box CLI result the gate can observe.

## Implementation decisions

- `internal/lines` remains the only deep module. It parses the closed
  `BENCH_<HARNESS>_<TIER>` matrix for the fixed harness set `codex`, `claude`, and
  `opencode`; resolves a tier token for the named harness; renders the named harness
  first; and owns every allow, deny, and degraded verdict. The old
  `BENCH_TIER_*`/`BENCH_ALIAS_*` schema is a hard cut, never a dual read.
- `--harness codex|claude|opencode` replaces shape flags on both plumbing commands.
  `BENCH_MODEL` names `top`, `mid`, or `cheap`, never a concrete id. Only declared
  harnesses require all three safe tokens; OpenCode stays intentionally unbound and
  fails closed at launch.
- **Evidence prerequisite and approval block:** the pinned tree has no reproducible
  real Claude context-fork envelope that proves both the discriminator’s JSON path and
  value and an authoritative inherited-session-tier field. Those are open evidence,
  not decisions: this spec assumes neither a field name nor a non-fork value. Before a
  parser or verdict change, capture a reproducible real context-fork fixture that
  establishes all three facts. Until then, the default is no parser change and story 2
  is **blocked for approval and implementation**. If a real envelope has no
  authoritative inherited tier, fail closed at this feature boundary: surface the
  implementation blocker and do not replace the missing fact with requested-model
  inference or a broader deny-all-forks policy.
- The semantic verdict is fixed independently of JSON spelling: when the captured
  delegation type denotes a context fork and its declared model differs from the
  captured authoritative inherited session tier, deny with exit 2. Equality preserves
  the normal bound-model path; an omitted model keeps the existing routed-complete
  deny rim. Comparisons are exact strings only: no provider lookup, model discovery,
  or inference is added.
- The denial renderer takes the active harness, uses that harness’s three bound tokens
  as its actionable primary listing, and may show the remaining matrix only through a
  deliberately named verbose describe mode. The spec writer chooses the exact prose
  and verbose flag; the ordering and single rendering owner are fixed.
- Command parsing, adapters, and the hook are thin pass-throughs. Each supplies its
  own harness name; no shim parses bindings or reconstructs a verdict. The hook keeps
  its existing fail-open rims, while an intentional core deny remains exit 2.
- `checkLineBinding` validates the matrix, the profile table, and declared-harness
  completeness from the same parsed binding the runtime uses. `bench doctor` detects
  retired keys and reports the one-way rewrite; it never rewrites the file.
- Add one conformance extraction/check owner for command prose. It discovers entries
  from `.agents/commands` at run time, sorts them, and traverses every regular file;
  a non-regular discovered entry is rejected before any read. Today’s quantified
  directory inventory is these ten files: `bench-assess.md`, `bench-debug.md`,
  `bench-final-check.md`, `bench-implement-spec.md`,
  `bench-review-implementation.md`, `bench-setup-repo.md`, `bench-shape-idea.md`,
  `bench-update-kit.md`, `bench-what-next.md`, and `bench-write-spec.md`. That list is
  a specification snapshot, not a production or executable-test registry. For every
  literal matching the command-model token grammar (a safe `gpt-` or
  provider-qualified model token, or a Claude alias token), the owner accepts only a
  value in the corresponding parsed matrix; the diagnostic names the file, token, and
  binding source.
- The profile presents the full declared matrix and `craft-line` requires a line
  declaration or recommendation in the current harness’s family. No concrete model
  choice changes in this build.

## Testing decisions

- The primary seam is the `internal/lines` exported verdict/resolver surface. Tests
  pass binding bytes and real-shaped envelopes, then observe exit code, stdout, and
  stderr rather than private parser state.
- The conformance seam grades a temporary repository through `checkLineBinding` and
  the command-token sweep. Its mutation proof must start from an otherwise clean copy
  and add only one unbound token, so the target diagnostic cannot be collateral.
- Runtime wiring stays at the existing adapter/hook subprocess seam; the kit wrapper,
  linked-repo wrapper, Claude hook, and all three adapters drive the same core.
- Prior art: the existing `internal/lines` table tests, line-routing execution tests,
  and line-routing canary fixtures. The feature gate is `bench gate`; focused unit,
  conformance, and runtime checks precede it.
- Story 2 has no green implementation path yet. Its first build step is to capture a
  reproducible real Claude context-fork envelope fixture that identifies the actual
  discriminator and authoritative inherited tier; no claimed version, field path, or
  value is evidence in the current tree. If that capture cannot supply the inherited
  tier, the build reports the blocker and stops rather than broadening the verdict. A
  fixture test may then make the semantic mismatch case red without any live provider
  query.

### Seam diagram: line verdict and harness resolution

```text
trigger: resolve-model, check-agent-line, or a shift adapter
    |
    v
lines.env + --harness + BENCH_MODEL/envelope --> [ internal/lines ] --> model, report, or exit 0/1/2
                                             [ matrix + verdicts ]
                                                ^ tests attach here: bytes in, result bytes out
```

### Seam diagram: static command-token conformance

```text
trigger: conformance phase
    |
    v
.agents/commands directory + lines.env --> [ command-token sweep ] --> targeted diagnostics
                                         [ sorted regular-file discovery ]
                                            ^ tests attach here: add a newly named file, then one-token mutation
```

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | `resolve-model --harness codex` maps `BENCH_MODEL=cheap` to the Codex cheap cell | line verdict | Not TDD-able before the new resolver exists; future `internal/lines` unit test must first fail with the wrong cell. | A canonical-family fallback or concrete BENCH_MODEL cannot produce the required cell. |
| 1 | `resolve-model --harness claude` maps `BENCH_MODEL=cheap` to the Claude cheap cell | line verdict | Not TDD-able before the new resolver exists; future `internal/lines` unit test must first fail with the wrong cell. | It independently proves the peer matrix rather than only the Codex column. |
| 1 | an unbound OpenCode harness exits 1 and names its missing binding | line verdict | Not TDD-able before the new resolver exists; future unit test must first fail on an allow or vague error. | An accidental OpenCode fallback makes this black-box result green only if the fail-closed rule is missing. |
| 2 | a captured real envelope whose type denotes a context fork and whose declared cheap alias differs from its captured inherited top tier denies exit 2 before requested-model allow logic | line verdict | **Blocked by evidence prerequisite:** no authoritative real-envelope fixture exists in the tree. After capture, its mismatch case must first fail with the current allow. | The current implementation accepts the alias; the independent inherited-tier fact distinguishes the reviewed mismatch verdict from a broader fork policy. |
| 2 | a captured real envelope with equal declared and inherited bound models preserves normal allow behavior | line verdict | **Blocked by evidence prerequisite:** add the equality positive control only after the same authoritative fixture establishes the actual envelope facts. | A deny-all-forks implementation fails this positive control, while a requested-model-only check cannot establish the independent equality fact. |
| 2 | Claude `opus` remains an allowed bound token | line verdict | Already covered in part by current Agent-line execution tests; migrate its alias case to the harness matrix. | It proves the permissive matrix verdict did not become own-family-only or fork-only denial. |
| 2 | Codex `gpt-5.6-sol` remains allowed through the Claude guard | line verdict | Not TDD-able before the peer matrix exists; future table test must first fail on the cross-family allow. | It preserves the map’s closed permissive-enforcement decision independently of native advice. |
| 2 | an unbound `bogus` token denies with a Claude-family recovery instruction | line verdict | Not TDD-able before the native renderer exists; future table test must first fail with a non-Claude or generic deny. | A plain allow, wrong family, or unhelpful message fails the black-box denial contract. |
| 2 | malformed envelope and missing model preserve the settled fail-open and routed-complete deny rims | line verdict | Already covered; update existing table cases only where the new field changes fixture shape. | The fork branch cannot silently widen or replace accepted residual postures. |
| 3 | kit CLI, linked-repo CLI, Claude hook, and Codex/Claude/OpenCode adapters each pass their explicit harness to the same core | adapter/hook subprocess | Not TDD-able before the new flags exist; future line-routing execution test must first fail per surface. | An omitted surface either selects the wrong family or takes a retired flag path. |
| 4 | a profile matrix cell mismatching lines.env emits a diagnostic | matrix conformance | Not TDD-able before the matrix checker exists; future temporary-root test must first fail on one changed profile cell. | It proves prose and machine binding cannot drift independently. |
| 4 | a declared Claude harness missing its mid cell emits a diagnostic | matrix conformance | Not TDD-able before the matrix checker exists; future temporary-root test must first fail on that one empty cell. | It proves declared-harness completeness rather than merely token syntax. |
| 4 | a malformed matrix token and missing lines.env each emit their named diagnostics | matrix conformance | Not TDD-able before the matrix checker exists; future table test must first fail on each input. | Safe-token and required-file checks are separate ways an unsafe binding could otherwise pass. |
| 4 | absent OpenCode cells do not diagnose while Codex and Claude are complete | matrix conformance | Not TDD-able before the matrix checker exists; future temporary-root test must first fail if absence is treated as a declared partial matrix. | It proves the closed declared-versus-known harness distinction. |
| 5 | every regular command file discovered from `.agents/commands` is scanned for tier-model tokens; today the quantified inventory is the ten named files in the implementation decision | command-token conformance | Not TDD-able before the sweep exists; future temporary-root test must add a newly named command file containing one unbound token and first fail with its file-and-token diagnostic. | A stale named list or partial traversal misses the newly discovered file; production discovery, not an executable test registry, owns the universal set. |
| 5 | a single unbound literal inserted into an otherwise clean command copy emits its file-and-token diagnostic | command-token conformance | Not TDD-able before the sweep exists; future mutation test must observe the targeted diagnostic with no pre-existing diagnostics. | The clean baseline excludes collateral failure and proves this check bites for the intended omission. |
| 5 | every bound matrix token and non-model prose token remains accepted | command-token conformance | Not TDD-able before the sweep exists; future clean-fixture test must fail before the allowlist is derived from the parsed binding. | It prevents a broad text matcher from making all command guidance red. |
| 6 | Claude denial leads with `fable`, `opus`, and `sonnet`, while the verbose describe output retains matrix visibility | line verdict | Not TDD-able before the renderer exists; future exact ordering test must first fail with the current model-id-first message. | A sorted/global listing is observably not executable in Claude Code. |
| 6 | `--describe --harness claude` leads with the Claude column | line verdict | Not TDD-able before describe is rewritten; future stdout test must first fail with the current Codex-first report. | It covers the SessionStart-facing assertable separately from denial text. |
| 7 | `bench doctor` reports retired-key rewrites without modifying lines.env | doctor CLI | Not TDD-able before doctor owns this migration; future fixture CLI test must first fail on missing rewrite text or a changed file. | It distinguishes report-and-offer from silent mutation and from an unhelpful generic warning. |
| edge of 1 | quoted, CRLF, last-wins, and no-final-newline matrix entries parse consistently | line verdict | Already covered by existing parser tables; extend their key set to the new schema. | A schema migration that drops shell-compatible parsing fails a focused case. |
| edge of 1 | absent and present-but-empty lines.env remain distinct | matrix conformance | Already covered in current line-routing tests; migrate fixtures to matrix keys. | A reader that collapses the states changes either the fail-open rim or the diagnostic. |
| edge of 1 | unknown `--harness` and unknown/unset tier token exit 1 | line verdict | Not TDD-able before `--harness` parsing exists; future CLI table test must first fail on each invalid input. | Argument validation cannot accept an arbitrary family or model. |
| edge of 2 | a missing, malformed, or unexpected captured discriminator never impersonates a context fork | line verdict | **Blocked by evidence prerequisite:** the fixture must establish the actual discriminator before a degraded-posture test can be specified. | It keeps hostile envelope handling explicit instead of treating arbitrary text as a trusted type. |
| edge of 2 | if the real context-fork capture has no authoritative inherited-tier field, story 2 stops with an implementation blocker and no substitute fork verdict lands | evidence prerequisite | **Blocked now:** inspect the required capture before parser work; absence of the tier fact is the red approval signal. | It makes lack of evidence fail closed for the feature, preventing requested-model inference or deny-all-forks from masquerading as the reviewed mismatch rule. |
| edge of 3 | hook/adapters retain their existing missing-core fail-open rims | adapter/hook subprocess | Already covered by current line-routing execution tests; migrate arguments without changing rim assertions. | The new plumbing cannot brick delegation when its core is unavailable. |
| edge of 5 | a FIFO, device, socket, or symlink in command discovery is rejected before reading | command-token conformance | Not TDD-able before the sweep exists; future temporary-root test must install a FIFO with no writer and require the targeted discovery diagnostic without blocking, then cover the other special entry kinds. | It proves classification precedes reads, so static inspection neither blocks nor follows bytes outside the command directory. |

Degenerate implementations considered: an always-allow mismatch branch fails the
captured mismatch row once evidence unblocks it; a deny-all-forks branch fails the
captured equality row; a Codex-only matrix fails the Claude resolver and ordering rows;
a stale named-list sweep fails the newly discovered-file row; and a sweep that emits a
generic failure passes neither the clean-baseline mutation requirement nor the
targeted-token row.

### Edge inventory

- Error path — resolved by the unbound OpenCode, unknown harness, unknown tier, and
  doctor rows.
- Empty/absent input — resolved by the absent-versus-empty and missing-model rows.
- Boundary values — resolved by the current ten-file inventory and newly-discovered
  regular-file row; ten is the spec-time count, while directory discovery is the
  production universal set.
- Malformed input — resolved by malformed token, malformed envelope, and hostile
  command-discovery rows.
- Interrupted or partial state — resolved by declared-harness partial-matrix rows.
- Re-run idempotency — resolved by the doctor no-write row and migration hard-cut;
  reporting the same retired file twice changes no repository state.
- Hostile environment — resolved by quoted/CRLF/no-final-newline parsing, special-file
  discovery, and missing-core rim rows.
- Paths containing spaces or glob characters — **Won't handle**: discovery is rooted
  in the tree-controlled command directory and accepts no caller-supplied path, so no
  external path crosses this seam.
- Control bytes in model text — **Won't handle**: safe-token validation rejects them
  before renderer or command-token matching, and the matrix-malformed row owns that
  rejection.
- Tabs, newlines, and returns permitted by another sink — **Won't handle**: model and
  alias token grammars reject them before a line-oriented diagnostic can receive them.
- A command changing the fact it reports — **Won't handle**: the doctor row asserts
  report-only migration, so this feature performs no binding write whose state it must
  then describe.
- Unquoted multi-word arguments or a flag value read as a positional — **Won't
  handle**: the fixed `--harness` grammar owns only its named flag and the future
  argument-validation row covers its unknown and missing values; this feature adds no
  shell word parser.
- Required tool absent from PATH or invocation through a symlink — **Won't handle**:
  the existing wrapper-resolver and hook-rim tests own those transport conditions, and
  this build keeps the shim thin.
- Destructive worktree state or SIGINT mid-loop — **Won't handle**: no FT128 behavior
  allocates, leases, or cleans a worktree, so those owners remain the lifecycle seam.
- Re-run lifecycle commands — **Won't handle**: the only in-scope repeatable command is
  the doctor report, whose no-write row proves idempotency without borrowing link/setup
  behavior.
- Invocation from a nested cwd — **Won't handle**: existing repository-root discovery
  already owns it and this feature does not add a cwd-sensitive caller.
- Non-TTY stdin — **Won't handle**: the Agent hook’s JSON stdin is deliberately
  non-interactive, so a TTY contract would amputate its sole in-scope caller.
- Host-backed filesystem I/O pressure — **Won't handle**: this static parser and
  conformance sweep add no durability or fsync protocol; the gate-run owner retains
  that environment-specific contract.

## Out of scope

- Binding live OpenCode model identifiers — a separate harness-adoption capability,
  not completion of the unbound fail-closed behavior — 7 edits, 2 gate runs.
- Harness-scoped `bench models` discovery output — an advisory discovery capability
  separate from binding enforcement and rendering — 5 edits, 2 gate runs.
- Codex hook parity for agent-line denial — blocked by the current non-deny-capable
  spawn event: external-capability blocker — 0 sound repository edits, 0 gate runs
  until Codex exposes a deny-capable spawn event.
