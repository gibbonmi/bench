# Exercise hostile process and grammar cases

Blocked by: compare-four-observations.md
Ownership fence: `internal/axi/compatibility`, `cmd/bench/axi_compatibility_test.go`, `specs/axi-compatibility-oracle/testdata`
Integration surfaces: four-observation comparator→compare-four-observations.md; command classes→`decisions/byte-preserving-axi-foundation/assets/ft173-command-help-inventory.md`; exact executable→`scripts/go-build.sh` exercised unchanged
Contracts: argv vector, cwd path, PATH value, invocation path, deadline, and teardown state cross hostile fixtures→`cmd/bench/axi_compatibility_test.go`, membership is every inventory class, order is setup-run-teardown, and missing teardown refuses, asserted by HP1
Closure: HP1/malformed-argv, HP1/outside-repo, HP1/deep-cwd, HP1/stripped-path, HP1/symlink-invocation, HP1/bounded-interrupt, HP1/teardown

## What to build

every command-inventory hostile argv and environment class retains exact observations under bounded fresh processes.

## Acceptance

- [ ] [HP1] (covers CO6) every command-inventory hostile argv and environment class retains exact observations under bounded fresh processes.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| HP1/malformed-argv | accept one repeated flag | hostile grammar test | run the derived argv and require acceptance mismatch |
| HP1/outside-repo | reuse repository cwd for the outside-repo case | hostile environment test | run in a fresh non-repository root and require the exact refusal |
| HP1/deep-cwd | force root cwd before invoking the candidate | nested cwd test | invoke from a deep subdirectory and require baseline equality |
| HP1/stripped-path | resolve an optional tool from ambient PATH | stripped PATH test | invoke with the named minimal PATH and require baseline equality |
| HP1/symlink-invocation | canonicalize away the public symlink path before dispatch | symlink invocation test | invoke through the fixture symlink and require baseline equality |
| HP1/bounded-interrupt | remove the child deadline | interrupt test | signal the fresh process under the named timeout and require terminal observation |
| HP1/teardown | leave one descendant alive after comparison | process teardown test | after the deadline require no descendant and no reusable process state |

