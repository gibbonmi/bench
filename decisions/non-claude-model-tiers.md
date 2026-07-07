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
— (open)

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
— (open)

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
