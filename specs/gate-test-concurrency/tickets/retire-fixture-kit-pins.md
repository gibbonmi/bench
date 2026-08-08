# Retire the fixture constructors' kit pins

Blocked by: inject-kit-root-below-entries.md
Ownership fence: `internal/gate`
Integration surfaces: kit-taking execution boundary→inject-kit-root-below-entries.md; fixture constructor signatures→`internal/gate`; representative pinned entry tests→`internal/gate`
Contracts: the kit-taking execution boundary crosses inject-kit-root-below-entries.md→`internal/gate`, with root and kit as path strings, each migrated call preserving its existing operation order, and absence represented by no alternate entry call, asserted by RP1 against the real blocker-produced boundary
Closure: RP1/kitshaped-constructor, RP1/routed-constructor, RP1/entry-migration, RP1/hostile-run

## What to build

The kit-shaped and routed fixture constructors claim kit identity through the
injection seam instead of `t.Setenv("BENCH_KIT", root)`, so the roughly 52
construction sites stop pinning the process environment.

Constructor-side injection alone cannot carry the fixture tests that drive the
exported entries directly (`RunCommand --fresh`, the reuse execute, the
composed-green query — the build-skip, verdict-reuse, evaluation, and
component-decision families), because entry-time resolution reads the ambient
environment. Those call sites migrate to the kit-taking boundary the blocker
ticket exposed, except a small representative set that keeps driving each
exported entry as its subject: those representatives pin `BENCH_KIT` to the
fixture root explicitly (env is the input at that seam) and are thereby
serial. The two deliberate env-behavior tests — subject-env stripping and the
kit≠root manifest resolution — likewise keep their pins untouched.

The proof the retirement is real is one hostile-environment package run
recorded as build evidence: `BENCH_KIT=<foreign path> go test -count=1
./internal/gate` green. Before this ticket that command's green depends on the
per-test pins; after it, injection carries every migrated construction while
the exported variable stays hostile for the entire run, and only the
enumerated representatives and env-behavior tests still pin.

## Acceptance

- [ ] [RP1] (covers KC3) no fixture construction pins `BENCH_KIT`; entry-driving fixture tests run through the kit-taking boundary except an enumerated representative set that pins explicitly; the whole package is green with a foreign `BENCH_KIT` exported for the entire run.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RP1/kitshaped-constructor | drop the kit injection from the kit-shaped constructor (construct without claiming kit identity) | the hostile-environment package run | apply, run `BENCH_KIT=/nonexistent go test -count=1 ./internal/gate`, expect kit-shaped fixture tests red |
| RP1/routed-constructor | drop the kit injection from the routed constructor | the hostile-environment package run | apply, run the same command, expect routed-fixture tests red |
| RP1/entry-migration | move one migrated fixture test back to the exported entry without a pin | the hostile-environment package run | apply, run the same command, expect that test red against the foreign kit |
| RP1/hostile-run | make one constructor prefer ambient `BENCH_KIT` when set | the hostile-environment package run | apply, run the same command, expect that constructor's tests red against the foreign kit |
