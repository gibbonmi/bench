# Close the system serial table

Blocked by: 12-run-the-exec-and-charge-tests-in-parallel.md
Writes: internal/systemtest/parallel_census_test.go (new)
Covers: LGP31, LGP32, LGP33, LGP34, LGP35, LGP36, LGP38

## What to build

The serial table of the system census closes. The temporary "pending migration" reason is gone, and the census forbids it.

`TestSystemSerialTableIsExact` refuses any row with the reason "pending migration". The owner-repository reason names exactly `TestSelectedExecutableComposition` and `TestOwnerSelectionChild`, and the census proves that both reach `owner.repos`.

Each other serial row quotes a failing parallel run that an earlier ticket recorded.

The gate argv, the phase table, the `projects/benchkit.md` system sentences, and ADR 0024 stay unchanged.

Record the system phase time after the change. The orchestrator then pairs the landing-gate timings with each retained red proof (LGP35, LGP36).

BENCH_KIT is supplied by bench test --check system through the existing sealed system owner.

## Acceptance

- [ ] A serial row with the reason "pending migration" reds by name (LGP38).
- [ ] Each serial row names an existing test and a reason (LGP31).
- [ ] A test kept serial on an owner repository beyond the two subject tests reds by name (LGP32).
- [ ] `TestBenchkitPhasesUseBranchNativeDrivers` passes with the system argv unchanged (LGP33).
- [ ] The verification record holds the system phase time before and after the change (LGP34).
- [ ] The build log names each reverted attempt, and the final report pairs timings with retained reds (LGP35, LGP36).
