# Non-Claude Model Tiers

## #1: Who chooses the cheap/mid/top tier binding?

Blocked by: —
Type: Grill

### Question
The kit already has abstract cheap/mid/top roles, but the current live binding is
Claude-shaped. The build changes depending on whether Bench lets each harness pick
its own tier mapping, or keeps tier ownership with the reviewer and only makes the
accepted model id space wider.

### Answer
The reviewer owns the tier binding. `.bench/lines.env` remains the authoritative
machine-readable source for `BENCH_TIER_TOP`, `BENCH_TIER_MID`, and
`BENCH_TIER_CHEAP`; harnesses may expose candidate model ids, but they do not
choose tiers. Discovery is advisory input to the reviewer-owned binding.

## #2: What model-id grammar should `.bench/lines.env` accept?

Blocked by: #1
Type: Grill

### Question
Conformance currently treats a tier value as valid only when it matches
`claude-*`. Codex and other harnesses need owner-defined ids in the same binding
surface, so the spec needs a portable model-id grammar and malformed-value
posture.

### Answer
Tier values are opaque reviewer-owned model-id tokens, not provider-prefixed
Claude ids. The gate validates only that each `BENCH_TIER_*` value is a safe
non-empty printable token for a harness to receive, rejecting whitespace,
control bytes, and shell-dangerous characters with the key and value named.
Provider existence is not enforced by the gate: discovery is advisory, so an
unknown-to-discovery but syntactically safe model id remains a valid binding.

## #3: What should `bench models` discover?

Blocked by: #1, #2
Type: Research

### Question
`bench models` currently queries the Anthropic Models API when
`ANTHROPIC_API_KEY` is set, then otherwise points the reviewer at each harness's
manual model list. The build needs to decide whether this command becomes a
multi-provider discovery surface, a harness-local advisory surface, or remains
Anthropic-only with corrected copy.

### Answer
`bench models` becomes a multi-source advisory inventory of candidate model ids,
not a validator for `.bench/lines.env`. It should list every source it can reach
and print a clear unavailable row for each source it cannot reach, while exiting
0 so discovery never blocks tier ownership.

Sources for the spec:

- Codex local catalog when `codex` is on `PATH`: run `codex debug models`,
  parse JSON `models[].slug`, and list visible slugs. Use `--bundled` as the
  no-network fallback when refresh fails.
- OpenAI API when `OPENAI_API_KEY` is set: call `GET /v1/models` and list
  `data[].id`; OpenAI documents those ids as model identifiers usable in API
  endpoints.
- Anthropic API when `ANTHROPIC_API_KEY` is set: keep the existing
  `/v1/models` behavior and list `data[].id`.
- Manual fallback when no source is reachable: print exact commands/env vars
  the reviewer can run or set, including `codex debug models --bundled`,
  `OPENAI_API_KEY`, and `ANTHROPIC_API_KEY`.

The output groups rows by source and labels freshness (`live`, `bundled`, or
`unavailable`). It does not infer cheap/mid/top, does not reject safe ids absent
from discovery, and does not require every provider to be installed.

## #4: Which docs and gate checks must stop saying Claude is the tier namespace?

Blocked by: #2, #3
Type: Research

### Question
The line guidance, project profile, conformance checks, and adapter tests all
name Claude ids in some places. The spec needs a concrete edit list so one source
still owns tier facts and the gate catches drift without rejecting non-Claude
owner-defined bindings.

### Answer
Spec edit list:

- `internal/conformance/line_routing_checks_test.go`: replace the
  `^claude-*` tier-value regex with the opaque safe-token grammar from #2.
  Add bite tests for accepted non-Claude ids (`gpt-5.4`,
  `gpt-5.3-codex-spark`, and a provider-qualified id such as
  `openai/gpt-5`) and rejected unsafe ids (empty, whitespace, control bytes,
  shell metacharacters). Keep profile-drift checks against the literal values
  from `.bench/lines.env`.
- `internal/conformance/line_routing_checks_test.go`: update hook/adapter
  fixture bindings so at least one routed repo uses Codex-shaped ids. The
  adapter and Agent hook behavior should prove exact membership against
  `lines.env`, not accidental compatibility with Claude-shaped examples.
- `internal/models/models.go` and `internal/models/models_test.go`: replace the
  Anthropic-only `bench models` implementation and no-key text with the #3
  multi-source advisory inventory. Tests should cover Codex catalog parsing,
  OpenAI/Anthropic `data[].id` parsing, source unavailable rows, and exit-0
  failure posture.
- `.agents/commands/bench-setup-repo.md`: update setup guidance so `bench
  models` is described as multi-source advisory discovery. The reviewer still
  binds cheap/mid/top; the harness or provider never assigns tiers.
- `.agents/skills/bench-craft-line/SKILL.md`: clarify that tier values are
  opaque model-id tokens and `bench models` is candidate discovery, not
  validation. Keep the Claude alias caveat as harness-specific, not as the
  general tier model.
- `projects/benchkit.md` and `.bench/lines.env`: keep the current Claude
  binding if it remains this repo's chosen line, but reword surrounding prose so
  those ids are examples/current binding, not the allowed namespace. The
  conformance check should still require the profile to name the literal bound
  values.
- `.bench/BENCH.md` / `.bench/BENCH-reference.md` only need edits if their
  line or `bench models` prose still implies provider validation after the
  changes above.
- Canary fixtures that embed any changed command/guide text need the same
  fixture updates as usual so stale-reference canaries keep testing the current
  source text.

Do not change the rule that `.bench/lines.env` is the owner-controlled binding,
and do not make unknown-to-discovery model ids red.

## Handoff

1. **Module boundaries.** `internal/lines` remains the pure parser/verdict engine
   for membership in `.bench/lines.env`; `internal/conformance` owns structural
   validation of the binding and profile drift; `internal/models` owns advisory
   discovery; `.agents`/`.bench`/`projects` own the human guidance.
2. **Contracts.** `.bench/lines.env` carries `BENCH_TIER_TOP`,
   `BENCH_TIER_MID`, and `BENCH_TIER_CHEAP` as reviewer-owned opaque safe model
   tokens; optional `BENCH_ALIAS_*` values stay bare aliases for harnesses that
   need them. `bench models` exits 0 and prints grouped candidate ids with source
   freshness or unavailable rows.
3. **Deep vs thin.** `internal/lines` is deep for parsing and exact-membership
   verdicts; `internal/models` is deep for provider/catalog probing and output
   shaping; adapters and hooks stay thin callers of the existing verdict
   surfaces.
4. **Black-box assertables.** Tests can assert `bench models` stdout rows and
   exit code, `checkLineBinding` diagnostics for valid/invalid tier tokens, hook
   allow/deny behavior for non-Claude bound ids, and adapter behavior for
   `BENCH_MODEL` exact membership.
5. **Gate attachment.** The gate observes the conformance package, model command
   tests, adapter/hook behavior tests, docs drift checks, and canary fixtures.
   No manual-only seam is needed.
6. **Hostile-input owners.** `internal/lines` owns spaces, quotes, CR/no-newline,
   control bytes, and shell-dangerous tier values; `internal/models` owns absent
   commands, missing API keys, malformed JSON, live endpoint failure, and
   unavailable providers; adapters own dash-leading prompts and exact
   multi-word argument preservation.
7. **Uncertainty flags.** n/a — the reviewer resolved ownership, token grammar,
   and discovery posture.
8. **Rejected alternatives.** Harness-chosen tier mappings are rejected;
   provider-prefix-only validation is rejected; discovery-backed enforcement is
   rejected; `bench models` failing nonzero when a provider is unavailable is
   rejected.
9. **Domain watch-outs.** The current Bench repo may still bind to Claude ids in
   `.bench/lines.env`; that is a chosen binding, not the namespace. Discovery is
   information for the reviewer and must never become the oracle for line
   membership.

Dependency order: n/a — single spec.
