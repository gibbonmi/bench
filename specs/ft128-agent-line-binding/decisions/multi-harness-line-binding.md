# Multi-harness line binding

## Destination

Make the tier→model binding symmetric across harnesses so each harness reads,
reports, and recommends its own model family. Today `BENCH_TIER_*` holds Codex
gpt ids as the privileged canonical and `BENCH_ALIAS_*` is a Claude-only
translation, so every advisory surface — the SessionStart line report, the
profile "Lines" section, and the agent's own recommendations — leads with gpt
ids even inside a Claude Code session. After this change no family is canonical:
the tier (top/mid/cheap) is the only abstract identity, and every reader resolves
it to the current harness's family. Enforcement already runs the right family per
harness (adapters + the Agent-line hook); this is the *advisory* layer catching up.

Scope confirmed with reviewer: the bite is presentation only — what the agent
recommends, the profile framing, and the SessionStart report. No run is
mis-launching today.

## #1: Data model — one canonical family, or symmetric peers?

Type: Grill

### Question
`BENCH_TIER_*` is named for the abstract tier but holds Codex ids, structurally
privileging Codex as "the real model." Fix the presentation only, or restructure
the data model so no family is canonical?

### Answer
**Symmetric peers.** The tier (top/mid/cheap) becomes the only canonical identity;
every harness family is a peer binding, none privileged. Bigger blast radius
(lines.env schema, Go core, conformance check, all three adapters) but it's the
only shape that stops a family being "the real id." Presentation-only was rejected
because prose would still call gpt "the tier."

## #2: lines.env schema — key shape and harness set

Type: Grill

### Question
How are the per-harness-per-tier cells keyed, and is the harness set open or fixed?

### Answer
**`BENCH_<HARNESS>_<TIER>`, fixed validated set.** e.g. `BENCH_CODEX_TOP`,
`BENCH_CLAUDE_MID`. Groups by harness. The core validates against a known harness
list (codex/claude/opencode); adding a harness is a deliberate code change, not a
free-form key. Chosen over `BENCH_TIER_<TIER>_<HARNESS>` open-set for tighter
validation.

## #3: How does the resolver/hook learn the current harness?

Type: Grill

### Question
Today the harness is implied by output-shape flags (`--alias`→claude,
`--provider-model`→opencode, bare→codex). Keep that, or make harness explicit?

### Answer
**Explicit `--harness <name>`.** Each adapter and the Agent-line hook passes
`--harness codex|claude|opencode`; output shape follows from the harness's column.
The old shape-flags (`--alias`, `--provider-model`) are retired — they conflated
"output shape" with "which harness," the exact conflation being removed.

## #4: Agent-line enforcement scope — own-family-only or any bound tier?

Type: Grill

### Question
In Claude Code, a delegation names e.g. `opus`. Should the hook (`--harness claude`)
accept only `BENCH_CLAUDE_*` values, or any bound cell across the matrix?

### Answer
**Any bound tier value.** Enforcement stays permissive — the honest-mistake guard
keeps today's "is it on a bound tier" semantics, now matrix-wide. A model matching
no cell still denies (exit 2), and the deny message uses `--harness` to recommend
the *current* harness's family. **Steering to the native family is advisory, not
enforced** — realized in recommendation/report/skill (#8), not the verdict.

## #5: Docs + conformance for a 3×3 matrix

Type: Grill

### Question
Nine concrete ids. Keep listing them in the profile prose (cross-checked), or defer
to lines.env as the one source?

### Answer
**Full matrix in prose, cross-checked.** The profile "Lines" renders the matrix as
a table; `checkLineBinding` cross-checks every cell against lines.env, extending
today's prose↔env check. Completeness rule: for any *declared* harness, all three
tiers must be present and be safe tokens. Deferring to lines.env was rejected —
reviewer wants the binding human-readable in the profile.

## #6: Which harnesses does this repo bind now?

Type: Grill

### Question
The schema knows three harnesses. Populate all three, or only those in use?

### Answer
**codex + claude only.** OpenCode stays unbound until adopted; its adapter refuses
to launch (fail-closed) until `BENCH_OPENCODE_*` is added, and `bench doctor` says
so. The completeness check requires all tiers only for *declared* harnesses, so an
absent opencode column is not a gate failure. Don't invent opencode ids not being run.

## #7: SessionStart line report shape

Type: Grill

### Question
`check-agent-line --describe` is where "claude code looks to use gpt models" bit.
With `--harness claude`, what does it print?

### Answer
**Current-harness column primary.** The report leads with the asking harness's
family (Claude: `top=fable mid=opus cheap=sonnet`), full matrix behind a verbose
flag. Directly fixes the bite — a Claude session sees Claude models first, no gpt
row in the default report.

## #8: BENCH_MODEL contract for a shift

Type: Grill

### Question
Adapters now pass `--harness`. A concrete-id `BENCH_MODEL` + `--harness claude`
is incoherent (which family wins?). What does BENCH_MODEL name?

### Answer
**A tier token (top/mid/cheap).** `BENCH_MODEL=cheap bench shift …` → codex adapter
resolves `gpt-5.6-luna`, claude adapter resolves `sonnet`. The reviewer picks a
tier; the harness picks the family. Fully realizes "tier is the abstract identity"
and removes any family/harness mismatch at launch. The opencode `--provider-model`
projection is retired — opencode gets explicit `BENCH_OPENCODE_*` ids (full
provider/model strings) when bound, not a projection off a canonical id.

## Migration (resolved under #1)

Hard cut, no dual-read shim. `bench doctor` detects the retired
`BENCH_TIER_*`/`BENCH_ALIAS_*` keys and reports the exact rewrite to
`BENCH_CODEX_*`/`BENCH_CLAUDE_*`. The core reads only the new keys. (Auto-rewrite
vs report-and-offer for a reviewer-owned file: spec-writer's call — lean
report-and-offer, see Handoff item 9.)

## The recommendation rule (resolved under #4/#7)

The actual ask — "recommend claude models in claude code." Since enforcement is
permissive, this lives in the advisory layer: `craft-line` gains a rule to declare
the line / recommend in the *current harness's* family; the SessionStart report
(#7) reinforces it; the tier memory is updated. No family is named in the abstract
rule — it resolves per harness.

## Not yet specified

- Exact deny-message wording and the verbose-flag name (`--all`?) for the full-matrix
  report — spec-writer picks.
- Whether `bench models` discovery output should also become harness-scoped.

## Out of scope

- Changing which concrete models are bound to each tier (a reviewer binding
  decision, not this restructure).
- Adding OpenCode's live binding (deferred to adoption, #6).
- Any change to the shift loop, gate, or worktree machinery beyond the BENCH_MODEL
  contract (#8).
- The agent's own auto-memory update (not a repo artifact).

## Handoff

1. **Module boundaries.**
   - `internal/lines` (deep, Go core) — owns binding parse into a harness×tier
     matrix, the Agent-line verdict, `ResolveModel*` verdicts, and `DescribeBinding`.
     The heart of the change: new schema, `--harness`, tier-token BENCH_MODEL,
     permissive matrix verdict, harness-primary describe.
   - `cmd/bench/main.go` (thin) — `resolveModel`/`check-agent-line` CLI arg parsing:
     add `--harness`, retire `--alias`/`--provider-model`, accept tier-token BENCH_MODEL.
   - `.bench/adapters/{codex,claude,opencode}` (thin shims) — pass `--harness <self>`;
     BENCH_MODEL is now a tier; opencode fail-closed until bound.
   - `.bench/hooks/check-agent-line.sh` (thin shim) — pass `--harness claude`; its
     `--describe` shows the claude column primary.
   - `.bench/lines.env` — new `BENCH_<HARNESS>_<TIER>` schema, codex+claude cells.
   - `projects/benchkit.md` "Lines" — matrix table replacing the gpt-canonical prose.
   - `internal/conformance` `checkLineBinding` — matrix cross-check + per-declared-harness
     completeness.
   - `bench doctor` — old-schema detection + rewrite report.
   - `craft-line` skill — harness-native recommendation rule.

2. **Contracts.**
   - `bench resolve-model --harness <h>` with `BENCH_MODEL=<tier>` → stdout the cell
     `BENCH_<H>_<TIER>`, exit 0; unbound harness/cell → exit 1 with a stderr naming
     the missing key; unset/unknown tier token → exit 1.
   - `bench check-agent-line --harness <h>` (stdin: Agent envelope) → exit 0 if the
     model matches any bound cell (permissive), exit 2 if it matches none, with a
     deny message recommending harness `<h>`'s family. Fail-open rims unchanged.
   - `bench check-agent-line --describe --harness <h>` → manifest with the `<h>`
     column as the primary line; verbose flag prints the full matrix.
   - `checkLineBinding` → diagnostics (gate red) on: missing lines.env, a declared
     harness missing a tier, a malformed/unsafe token, prose↔env cell mismatch.

3. **Deep vs thin.** `internal/lines` is the one deep module — all binding and
   verdict logic. Adapters, the hook, and `cmd/bench` are pass-throughs with no
   logic of their own; their seam is the core's exported funcs, not their own body.

4. **Black-box assertables.**
   - `resolve-model --harness codex`, `BENCH_MODEL=cheap` → stdout `gpt-5.6-luna`, exit 0.
   - `resolve-model --harness claude`, `BENCH_MODEL=cheap` → stdout `sonnet`, exit 0.
   - `resolve-model --harness opencode` (unbound) → exit 1, stderr names `BENCH_OPENCODE_*`.
   - `check-agent-line --harness claude`, model=`opus` → exit 0; model=`gpt-5.6-sol`
     → exit 0 (permissive); model=`bogus` → exit 2, message names the claude family.
   - `check-agent-line --describe --harness claude` → stdout shows `top=fable mid=opus
     cheap=sonnet` as the primary line.
   - `checkLineBinding`: cell mismatch prose↔env → non-empty diags; claude bound but
     `BENCH_CLAUDE_MID` unset → non-empty diags; opencode absent → no diag.
   - `bench doctor` on old-schema lines.env → reports the `BENCH_TIER_*`→`BENCH_CODEX_*`
     / `BENCH_ALIAS_*`→`BENCH_CLAUDE_*` rewrite.

5. **Gate attachment.** `internal/lines` unit tests (`go test`) observe the verdict
   and resolver seams — the primary gate surface. `checkLineBinding` runs in the
   conformance phase against lines.env + the profile. The shell shims (adapters,
   hook) carry no logic; their fail-open/closed rims have existing coverage. Two
   seams the gate cannot see: the live Claude Code SessionStart report rendering and
   the actual per-harness shift launch — both need manual verify.

6. **Hostile-input owners** (from the profile checklist):
   - malformed/quoted/CRLF/last-wins lines.env values → `internal/lines` ParseBinding
     (existing coverage; extend to new keys).
   - hand-edited file without trailing newline / absent vs empty lines.env →
     `checkLineBinding` + ParseBinding.
   - unknown `--harness` value → `resolveModel`/`check-agent-line` arg validation.
   - unknown/unset BENCH_MODEL tier token → resolve verdict.
   - retired old-schema keys present → `bench doctor` + core ignores them (hard cut).
   - Agent envelope missing the model field → agent-line fail-open (existing).
   - partial matrix (harness with 2 of 3 tiers) → completeness check.
   - invocation through every surface (kit CLI, by-path CLI, hooks, adapters) reaches
     the same routed core → the `--harness` plumbing must be uniform across all four.

7. **Uncertainty flags.** None require escalation — every decision was grilled and
   closed. Deny-message wording and the verbose-flag name are free spec-writer choices,
   not tier-escalating uncertainty.

8. **Rejected alternatives.** Presentation-only (kept Codex canonical); open-set
   `BENCH_TIER_<TIER>_<HARNESS>` keys; shape-flags as harness proxy; own-family-only
   enforcement; profile-defers-to-lines.env; dual-read compat shim; all-three-harnesses
   now; full-matrix (non-primary) report; concrete-id BENCH_MODEL. Do not reopen.

9. **Domain watch-outs.**
   - Claude Code's Agent tool only accepts model *aliases* (fable/opus/sonnet), not
     arbitrary ids — so `BENCH_CLAUDE_*` values are aliases, and a gpt id, though it
     passes the permissive hook, is not launchable there. This is why permissive
     enforcement is acceptable but native recommendation is the real fix.
   - OpenCode ids are provider/model-qualified; with the projection retired,
     `BENCH_OPENCODE_*` must hold full provider/model strings.
   - lines.env is read by hook + adapters *and* cross-checked by the gate; editing a
     cell in one place without the other turns the gate red — by design.
   - `bench doctor` rewriting lines.env touches a reviewer-owned config file; prefer
     report-and-offer over silent auto-rewrite. Confirm the posture in the spec.

Dependency order: recommended as one spec, but if sliced — (A) `internal/lines`
core: schema + `--harness` + tier-token BENCH_MODEL + permissive verdict +
harness-primary describe, with unit tests; (B) wiring: lines.env migration,
adapters, hook shim, `bench doctor`; (C) advisory: profile matrix table +
`checkLineBinding` + `craft-line` rule. A depends on nothing; B and C depend on A.
Slicing stays the reviewer's call.
