# Carry per-test bite scoping through the canary runner

Blocked by: Guard fixture marker reads against special files

## What to build

Stories 1 and 5 of `specs/ft91-gate-fastpath/spec.md`, plus the argv half of
story 4: `RunCall` gains a `Test` field and a `RunList` kind; `runnerCommand`
gains real-argv cases for both (bite appends `-test.run ^<QuoteMeta(owner)>$`;
list invokes the compiled binary with `-test.list '.*'`); `subjectCall` sets
`Test` for fixtures only, reading an *optional* `TEST` marker through the
shared guarded reader — enforcement and `-test.list` validation come in a
later ticket, so a fixture without `TEST` still runs package-wide and the live
sweep is unchanged. Contract-group baselines never carry `Test`.

## Acceptance

- [ ] Two fixtures with distinct `TEST` owners produce bite RunCalls each
      carrying its own anchored `-test.run` value, no clash (injected-Runner
      sweep test, two-fixture shape).
- [ ] An owner name containing regexp metacharacters matches only itself
      (argv-builder unit test; unquoted-interpolation mutation observed red
      and reverted, recorded in the return log).
- [ ] `TestDefaultRunnerDispatchesOnCallKind` extended: real argv asserted for
      both new shapes (bite with `-test.run`, list with `-test.list '.*'`).
- [ ] Contract-group baselines carry no `-test.run` (asserted against a
      scoped-baseline variant).
- [ ] A fixture with no `TEST` file behaves exactly as today (package-wide
      run, no refusal).
