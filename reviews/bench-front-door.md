## Standards

Finding count: 3

Worst issue: two production action derivations independently encode the same executable-action policy and can drift.

1. **ask-user** — `internal/status/route.go:36-76` independently parses the executable action allowlist also authored in `internal/status/status.go:142-159,271-298,366-385,409-486,576-691`. This violates the AGENTS.md one-source-per-fact policy; agreement in the status tests does not eliminate the drift risk.
2. **ask-user** — `cmd/bench/main.go:146-189` hand-maintains the CLI inventory independently of `commandRegistry` in `cmd/bench/main.go:70-124` and `.bench/BENCH.md:29-52`. This violates the same one-source-per-fact policy.
3. **auto-fix** — There is no concise typed `CHANGELOG.md` entry for the user-visible routing, root, adapters, and drain behavior, contrary to `.agents/skills/bench-craft-synthesis/SKILL.md:67-73`; the relevant changelog surface is `CHANGELOG.md:7-93`.

## Spec

Finding count: 5

Worst issue: an unreadable `.bench` path is falsely reported as needing setup.

1. **auto-fix** — `appendSetup` at `internal/status/status.go:142` emits setup for every `os.Stat` error, violating the quiet-unreadable behavior required by `specs/bench-front-door/spec.md:69`; `internal/status/status_test.go:34` covers only absent and present paths.
2. **auto-fix** — R14, R16, R20, and R21 at `specs/bench-front-door/spec.md:212,214,218-219` have only helper-level tests at `internal/status/status_test.go:95,136,179`, leaving their rendered routes unprotected.
3. **auto-fix** — The R24 all-producers test required by `specs/bench-front-door/spec.md:222` and implemented at `internal/status/status_test.go:242` omits the setup, guards, specs, drain, structure, maps, reviews, roadmap, leased, and out-of-pool producer families.
4. **auto-fix** — R27-R31 at `specs/bench-front-door/spec.md:225-229` lack exact fixtures for dirty-plus-unpushed, clean-unpushed, unique branch, leased/out-of-pool, and structure/map scan failure or orphan states; the current coverage is at `internal/status/status_test.go:766`.
5. **auto-fix** — R39 at `specs/bench-front-door/spec.md:237` requires a second-inventory-copy red test, but `cmd/bench/main_test.go:55` checks only runtime equality.

## Coverage

Finding count: 2

Worst issue: spec and decision-map paths containing spaces are misrouted.

1. **auto-fix** — A staged `specs/my spec/spec.md` or ready `decisions/my map.md` produces an action that `IsInvocable` rejects after whitespace splitting, so the candidate is skipped. The hostile-input checklist covers paths with spaces and glob characters, but R16, R20, and R21 do not exercise them.
2. **auto-fix** — `bench help extra` and direct-binary `help extra` ignore trailing arguments and exit 0 even though story 35 says help takes no arguments. The existing R37-R40 tests cover only valid forms.
