# Publish one sealed local subject

Blocked by: none
Ownership fence: `scripts/go-build.sh`, `cmd/bench`, `internal/freshness`, `internal/contract/freshness_subject_test.go`, `internal/conformance/subcommand_routing_test.go`, `internal/gate/prospective_test.go`
Integration surfaces: default builder grammar and command trace→`scripts/go-build.sh`; compiled plumbing dispatch→`cmd/bench`; sealed publication owner and standalone-publisher removal→`internal/freshness`; fresh-process subject journey→`internal/contract/freshness_subject_test.go`; plumbing classification→`internal/conformance/subcommand_routing_test.go`; prospective fixture copy inventory→`internal/gate/prospective_test.go`; unchanged default callers in `package.json`, `scripts/release-preflight.sh`, and `internal/gate/phases.go`→existing paths + LS1
Contracts: the absolute source-root string and literal destination path cross `scripts/go-build.sh`→the staged Bench plumbing operation→`internal/freshness`; LS1 and LS3 assert that the operation accepts exactly root then destination, binds the staged input to its own invoked executable, preserves argument order, and treats an absent operand as refusal before publication

## What to build

The default two-argument builder produces one host-runnable Bench subject with one `go build`. The staged Bench executable invokes an internal plumbing operation that binds publication to its own executable path and delegates atomic sealed publication to the freshness owner. The standalone publisher command disappears. The fence is deliberately broad because the real builder journey, compiled dispatch, publication owner, routing registry, and prospective fixture form one observable outcome; splitting any one away leaves a horizontal layer rather than an independently green subject.

## Acceptance

- [ ] [LS1] `go-build.sh <root> <output>` invokes exactly one `go build`, invokes no `go run` or second Go-built publisher, and leaves an executable whose seal matches its bytes and current source closure when verified in a fresh process.
- [ ] [LS2] Ambient non-host `GOOS` and `GOARCH` cannot change default mode: the builder resolves the host target, and the resulting sealed output executes `version` on the build host.
- [ ] [LS3] The internal publication operation accepts only source root and destination, binds the staged input to its own invoked executable, and reaches the existing freshness publication owner through the compiled command dispatch.
- [ ] [LS4] `internal/freshness/cmd`, `go run ./internal/freshness/cmd`, and every equivalent standalone publisher production path are absent, while plumbing classification and prospective fixture inventories remain complete.
- [ ] [LS5] Compile, validation, and sealed-publication failures plus bounded SIGINT preserve any prior output and prior seal byte-for-byte and leave no staged executable behind.
- [ ] [LS6] Missing operands and literal spaced, globbed, dash-led, symlinked, special, or unwritable output paths either publish the exact requested regular file safely or refuse before damaging an existing subject.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| LS1 | restore the `go run` publisher after the staged build | the real-builder command recorder | build a local subject, enumerate every Go and helper execution, expect the second build or helper invocation to fail the tally |
| LS2 | let ambient `GOOS` select the default build target | the non-host default journey | set a non-host target, build in default mode, execute `version`, and expect the non-runnable subject or missing valid seal failure |
| LS3 | accept a caller-selected staged executable path | the plumbing operation test | supply an extra staged path or substitute another executable, invoke the operation, and expect usage or executable-binding refusal |
| LS4 | retain the standalone publisher package or its prospective copy entry | the topology and routing contracts | enumerate the exact package, invocation stem, dispatch registry, and fixture copy list, and expect the retained production path or stale inventory to fail |
| LS5 | remove the cleanup trap or remove the old seal before the replacement pair is ready | the failure and interrupt table | fingerprint an existing output and seal, block or fail each stage, signal the builder where applicable, and expect both fingerprints plus an empty staging set |
| LS6 | parse a dash-led output as a mode or follow a symlinked destination | the hostile-output table | invoke each hostile path through the real builder and expect exact literal publication or a bounded refusal with the prior pair unchanged |
