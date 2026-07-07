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
— (open)
