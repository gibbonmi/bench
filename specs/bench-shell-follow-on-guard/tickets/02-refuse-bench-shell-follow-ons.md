# Refuse Bench shell follow-ons

Blocked by: 01-extract-shared-shell-token-stream.md
Writes: .bench/hooks/block-bench-follow-on.sh (new), .claude/settings.json, .codex/hooks.json, bin/bench.sh, cmd/bench, internal/benchguard (new), internal/systemtest, internal/conformance, internal/guards, internal/packagesurface, tests/canary/load-validity-metadata, tests/canary/line-routing, tests/canary/package-core-guard, projects/benchkit.md, CHANGELOG.md

## What to build

Consume the shared shell token-stream interface without changing it.
Add the Bench-follow-on classifier, its executable identity resolver, and its internal guard command.
Recognize direct, repository-relative, absolute, and live symlink spellings without running the target.

Add the thin fail-open hook shim and wire it into both Bash hook groups.
Require the shipped script, its manifest, and both harness wires.

The refusal must explain that the Bench response is bounded, complete, and self-contained.
Do not add Bench argument validation or a public projection option.

## Acceptance

- [ ] FOG01 allows one bare Bench call.
- [ ] FOG02 and FOG03 refuse both pipeline operator partitions before the marker runs.
- [ ] FOG04 and FOG05 refuse descriptor duplication and input redirection before Bench runs.
- [ ] FOG06 refuses an `&&` chain before the marker runs.
- [ ] FOG07 and FOG08 refuse semicolon and newline follow-ons before the marker runs.
- [ ] FOG09 and FOG10 scan the complete outer `worktree exec` command.
- [ ] FOG11 returns the bounded, complete, and self-contained explanation.
- [ ] FOG12, FOG13, FOG14, and FOG15 preserve supported arguments, quoted children, and incidental Bench text.
- [ ] FOG16 and FOG17 classify JSON-escaped shell operators.
- [ ] FOG18 inspects one routine shell wrapper.
- [ ] FOG19 and FOG20 warn and allow when the degraded rim fires.
- [ ] FOG21 and FOG22 require both harness wiring points.
- [ ] FOG23 advertises the complete guard row.
- [ ] FOG26 refuses an `||` chain before the marker runs.
- [ ] FOG27 recognizes a Bench invocation behind an environment prefix.
- [ ] FOG28, FOG29, FOG30, FOG31, and FOG32 warn and allow each unreadable command-field partition.
- [ ] FOG33 warns and allows when the platform binary is missing.
- [ ] FOG34, FOG35, and FOG36 refuse repository-relative, absolute, and live symlink spellings.
- [ ] FOG37 refuses plain stdout redirection before Bench runs.
- [ ] FOG38 refuses a Bench heredoc input redirection.
- [ ] FOG39 allows Bench-looking heredoc body text on a non-Bench command.
- [ ] FOG40 recognizes a Bench invocation behind a shell assignment prefix.
