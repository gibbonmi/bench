# 1. Name every seeded gate input in the setup preview

Blocked by: none
Writes: internal/adopt/setup.go, internal/adopt/setup_prompt_test.go
Covers: none

## What to build

The seeded `.bench/gate-inputs.json` declares four environment names: `BENCH_HOME`, `BENCH_KIT`, `BENCH_RUN_BINARY`, and `HOME`.
The `bench setup` preview line says that the seed declares only `BENCH_HOME` and `HOME`.
The preview test pins this incorrect text.
A test helper also decodes the names from the seed JSON again.

Put the declared environment names in one Go slice in `setup.go`.
The seed, the preview line, and the `runInteractiveSetup` test helper read that slice.
The preview test reads the slice and finds each name in the preview line.

Keep `wantSeededGateInputs` as an independent expectation.
Do not derive it from the slice, because it is the one test that finds an incorrect seed.

## Acceptance

- [ ] The preview line for an absent `.bench/gate-inputs.json` names each environment name that the seed declares.
- [ ] The seed bytes are equal to `wantSeededGateInputs`.
- [ ] When a probe removes one name from the slice, the `wantSeededGateInputs` test fails.
- [ ] No code decodes the environment names from the seed JSON again.
