# Demonstrate the BENCH_KIT seed red

Blocked by: cover-detected-consumer-hygiene.md
Writes: specs/roadmap-light-path-fixes/tickets/scaffold-declared-input-hygiene.md
Covers: none

## What to build

Remove `BENCH_KIT` temporarily from the scaffolded routing-input seed. Run the
focused independent seed expectation, restore production byte-for-byte, and
record the exact command and red output beside the existing routing-input
mutation evidence.

## Acceptance

- [ ] The focused seed expectation turns red when `BENCH_KIT` is omitted.
- [ ] Production is restored byte-for-byte after the mutation.
- [ ] The original LF2 ticket durably records the named mutation, command, and
      red output.
