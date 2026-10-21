# Run the exec and charge tests in parallel

Blocked by: 11-run-the-land-and-telemetry-tests-in-parallel.md
Writes: internal/systemtest/exec_bound_test.go, internal/systemtest/exec_nested_bound_test.go, internal/systemtest/charge_cleanup_test.go, internal/systemtest/charge_evidence_test.go, internal/systemtest/candidate_install_test.go, internal/systemtest/parallel_census_test.go (new)
Covers: LGP34

## What to build

The exec-bound, charge, and candidate-install tests run in parallel. None of them reaches `owner.repos`, so each one only gains `t.Parallel`.

Each test in these five files calls `t.Parallel`. Remove each test's "pending migration" row from the serial table. A test that fails under parallel load keeps a serial row, and its reason quotes the failing run.

`charge_evidence_test.go` has exactly 400 lines, and its eight new `t.Parallel` lines would put it over the limit. Split it along responsibility.

Keep `TestEvidenceInterruptedPublication` in `charge_evidence_test.go`, because `internal/conformance/injected_ports_registry_test.go:51` names that file and test.

The candidate-install tests run `go-build.sh` and share the Go build cache. If parallel builds conflict, keep those tests serial, and quote the failing run.

BENCH_KIT is supplied by bench test --check system through the existing sealed system owner.

## Acceptance

- [ ] Each test in these five files passes in parallel with its assertions unchanged, or holds a serial row that quotes its failing run.
- [ ] The injected-port registry still finds `TestEvidenceInterruptedPublication` in `charge_evidence_test.go`.
- [ ] The probe "remove t.Parallel from one exec-bound test" reds the census.
- [ ] The verification record holds the system phase time before and after the change (LGP34).
