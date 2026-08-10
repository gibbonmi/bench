# Repair final Terra AXI findings

Blocked by: repair-fable-review-findings.md
Ownership fence: `internal/specbuild/testdata/axi-ledger.md`, `internal/conformance/`, `projects/benchkit.md`, `cmd/bench/specbuild_test.go`, `internal/specbuild/`
Integration surfaces: ledger summary→derived disclosure cell census; conformance registry input→profile advertisement; retained state enumerator→real family-root command renderer
Contracts: none crosses
Closure: TR1/derived-ledger-summary, TR2/fixture-aware-input, TR3/real-populated-home

## What to build

Close the three accepted findings from the Terra/xhigh review of candidate `210aa31c294d4e673d12820a2d8296546c331b43`: derive or correct the ledger summary to 84 applicable / 105 N/A; bind the fixture-reading AXI disclosure check and its profile row to the honest input class; and drive populated, hostile, and capped retained state through the real top-level spec-build family-home route without a fake home service or direct renderer call.

## Acceptance

- [ ] [TR1] (covers local) the checked ledger summary agrees with its 189 classified rows and the production-derived conformance census, with no stale literal count surviving.
- [ ] [TR2] (covers local) the registered AXI disclosure check and project profile advertise an input class that includes its checked-in JSONL fixture dependency.
- [ ] [TR3] (covers local) one real family-root command test reads bounded populated retained state containing healthy, hostile, and cap-pressure entries and renders the surviving rows, diagnostics, useful help, and exit 0 through the production route.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| TR1/derived-ledger-summary | restore 95/94 while keeping 84/105 classified rows | ledger census test | parse the checked ledger and require its summary counts to equal the derived row dispositions |
| TR2/fixture-aware-input | relabel the check `go-source` while it reads JSONL | registry/profile agreement test | resolve the executable check dependencies and require the registered and advertised input class to cover them |
| TR3/real-populated-home | discard `Runs()` output before rendering the family home | top-level command test | invoke the real command over bounded mixed retained state and require healthy rows and diagnostics rather than an empty projection |
