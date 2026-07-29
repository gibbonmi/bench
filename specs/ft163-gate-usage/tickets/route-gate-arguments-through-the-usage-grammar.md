# Route gate arguments through the usage grammar

Blocked by: none

## What to build

`bench gate` treats help and invalid invocations as usage requests without
starting the oracle, while preserving the existing `bench gate pin` route.

## Acceptance

- [x] `bench gate --help`, `bench gate -h`, and `bench gate help` print
  `usage: bench gate` on stdout and exit 0 without starting a gate run.
- [x] Unknown flags and excess arguments print a usage refusal on stderr and
  exit 2 without starting a gate run or replacing gate evidence.
- [x] `bench gate pin` continues to route its remaining arguments to the
  existing pin command.
