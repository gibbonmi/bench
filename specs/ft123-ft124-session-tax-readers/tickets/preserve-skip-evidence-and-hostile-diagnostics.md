# Preserve skip evidence and hostile diagnostics

Blocked by: Render Go failures and result postures

## What to build

The structured test reader keeps generic and Bench-structured skip reasons,
contains hostile diagnostics in bounded TOON cells, cancels the whole Go process
group without partial output, and is reachable and advertised through every
shipped surface.

## Acceptance

- [x] Generic and structured capability or environment skips appear at default
  verbosity with reasons, while a no-skip run prints `skips[0]`.
- [x] An inherited `BENCH_SKIP_LOG` cannot divert a structured skip reason or
  modify the sentinel side channel.
- [x] SIGINT kills the fake Go process group, emits one structured interrupted
  stdout error, exits 1, and emits no partial TOON; a stream with no terminal
  package event is also a structured error.
- [x] Diagnostics at 120 and 121 code points and containing ESC, BEL, newline,
  tab, backslash, or invalid UTF-8 remain one representable bounded cell, while
  `--full` restores the complete escaped selected line.
- [x] Compiled routing, real-kit and linked-repository wrappers, help, runtime
  contract registration, CLI inventory, and the subcommand registry agree on
  `bench test`.
