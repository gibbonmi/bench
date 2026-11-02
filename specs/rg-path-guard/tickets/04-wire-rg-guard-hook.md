# Wire the rg guard into both harnesses and retire the prose rule

Blocked by: 03-add-guard-rg-verb.md
Writes: .claude/settings.json, .codex/hooks.json, internal/conformance/line_routing_static_test.go, internal/conformance/line_routing_static_cases_test.go (new), internal/conformance/validity_checks_test.go, internal/conformance/harness_record_test.go, internal/conformance/harness_record_fixture_test.go, internal/conformance/refusal_route_guard_test.go, cmd/bench/isolated_command_fixtures_test.go, tests/canary/line-routing/stop-hook-unwired, tests/canary/line-routing/agent-hook-unwired, tests/canary/load-validity-metadata/codex-hooks-timeout, tests/canary/load-validity-metadata/codex-hooks-timeout-typed, internal/harnesses/command_test.go, .bench/BENCH-reference.md, AGENTS.md
Covers: R39, R40, R47, R48, R49

## What to build

Claude Code and Codex both run `.bench/hooks/block-unsafe-rg.sh` from `03-add-guard-rg-verb.md` in their existing PreToolUse Bash group, so an agent in either harness gets the refusal before the call runs.

Add the shim to the Bash group of `.claude/settings.json` (after `block-bench-follow-on.sh`) and of `.codex/hooks.json`, in the shape of the existing entries. The PreToolUse Bash group already exists, so `internal/harnesses` takes no change (`internal/conformance/harness_record_test.go:160-181`).

Add one needle for the shim to the Claude wiring check (`internal/conformance/line_routing_static_test.go:503-512`) and to the Codex wiring check (`internal/conformance/validity_checks_test.go:211-217`). Add the shim to each intact fixture that those checks grade (`line_routing_static_test.go:575`, `validity_checks_test.go:229`), or the new needle reds them.
The mutation case "follow-on hook dropped" (`line_routing_static_test.go:609-611`) slices `Hooks[:1]`, so it also drops the shim, and `:640-645` refuses a sibling diagnostic. Change that case so it drops only the follow-on hook, and add one mutation case that drops only the shim (R39). Add the matching Codex mutation case (R40).
`internal/conformance/line_routing_static_test.go` has 653 lines over its budget with no grant, and the gate lane runs `structure --growth`. Move at least as many lines as you add into `internal/conformance/line_routing_static_cases_test.go`, as a cohesive block such as the mutation case table, and change no behavior.

Update each fixture or canary that mirrors the real Bash group, so that no mirror goes stale: the harness-record fixtures (`harness_record_test.go:311`, `:321`, `:351`; `harness_record_fixture_test.go:53`, `:60`), the `bench guards` fixture (`cmd/bench/isolated_command_fixtures_test.go:44-72`), and the canary settings files under `tests/canary/line-routing/` and `tests/canary/load-validity-metadata/` that list the Bash hooks. Change a mirror only when its check compares the Bash script set or its purpose is to copy the live config. Leave each canary's expected red unchanged.

Add the rg verdict to `stepGuardFaults` (`internal/conformance/refusal_route_guard_test.go:53-68`) through `rgguard.Classify`, and add a case where a route step `rg foo` gets an rg fault (R47).

Add one bullet for the rg guard to the Hook Layers section of `.bench/BENCH-reference.md`, after the shape of the `guard-file-write` bullet at `:883`. Name `guard-rg` and both harnesses. Keep the `bench harnesses codex` anchor (`internal/anchors/registry_data.go:157`) and the agent-line bullet name that `internal/harnesses/command_test.go:59` and `:183` pin.
Retire the bullet at `AGENTS.md:61-63` that states the explicit-path rule. No anchor pins it (`bench anchors AGENTS.md`). The shim's `# denies:` line and the refusal text are then the one source of the rule. Keep the other shell-convention bullets unchanged.

## Acceptance

- [ ] The Claude wiring check reds when the Bash group does not run `block-unsafe-rg.sh` (R39).
- [ ] The Codex wiring check reds when the Bash group does not run `block-unsafe-rg.sh` (R40).
- [ ] `stepGuardFaults` names an rg refusal for a route step `rg foo` (R47).
- [ ] The Hook Layers section names `guard-rg` (R48, review-owned).
- [ ] `AGENTS.md` holds no bullet that states the explicit-path rule (R49, review-owned).
- [ ] `bench test --check line-routing`, `bench test --check harness-record`, and `bench test --check root-conformance` pass.
- [ ] `bench guards` lists `block-unsafe-rg` as wired in `claude,codex`.
- [ ] `bench structure --growth` reports no grown file.
