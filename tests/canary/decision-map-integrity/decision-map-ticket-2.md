# lines.env schema — key shape and harness set

Blocked by: none
Type: Grill

### Question

How does the file key each per-harness, per-tier cell, and is the harness set
open or fixed?

### Answer

One `BENCH_<HARNESS>_<TIER>` key per cell, over a fixed validated set. The core
validates the harness name, so a new harness is a deliberate code change.
