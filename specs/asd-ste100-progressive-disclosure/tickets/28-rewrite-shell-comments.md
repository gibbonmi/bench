# Rewrite the shell comments and the shell scaffold strings

Blocked by: none
Writes: bin/bench.sh, bin/bench-postinstall.sh, .bench/gate.sh, .bench/gate-prospective.sh, .bench/hooks/, .bench/adapters/, .bench/lib/, scripts/aggregate-native-proofs.sh, scripts/build-artifacts.sh, scripts/build-offline-archives.sh, scripts/compare-artifacts.sh, scripts/gen-platform-packages.sh, scripts/go-build.sh, scripts/gremlins-diff.sh, scripts/install-govulncheck.sh, scripts/native-proof.sh, scripts/release-preflight.sh, scripts/smoke-artifacts.sh, scripts/smoke-offline.sh, scripts/lib/search.sh, internal/adopt/prepush.sh, internal/adopt/init.go, internal/adopt/setup.go, internal/adopt/
Line: `sonnet` / medium under Claude Code, `gpt-5.6-luna` / medium under Codex; the orchestrator validates at `fable` / high or `gpt-5.6-sol` / high.

## What to build

Every explanatory `#` comment in the shell surfaces reads in ASD-STE100 inside the `craft-comments` register. The shell surfaces are `bin/*.sh`, the gate scripts, the hooks, the adapters, the lib files, `scripts/*.sh`, `scripts/lib/*.sh`, and the embedded `prepush.sh`. The shell scaffold strings the kit writes into a tree join them: the gate scaffold from `scaffoldGate` and `setupGateScript` in `internal/adopt`. Tests that quote those strings update with them.

The guard manifest headers keep their `# key: value` lines byte-identical, and `session-start.sh` keeps `denies: nothing (informational)`. Shellcheck directives and shebangs stay byte-identical. A comment that restates its line is deleted.

The delegate changes comment lines in shell files and the scaffold string literals in Go only. The orchestrator verifies the comment-only diff on the shell files and reads each rewritten comment and scaffold sentence against its code.

## Acceptance

- [ ] `shellcheck`, `gofmt`, `vet`, and `test` stay green, and every shell file's comments and the gate scaffold read in STE (covers PD33).
- [ ] The guard manifest headers and the `session-start.sh` `denies:` value are byte-identical (covers PD34).
- [ ] Every changed line in a shell file starts with `#` after whitespace (covers PD35).
