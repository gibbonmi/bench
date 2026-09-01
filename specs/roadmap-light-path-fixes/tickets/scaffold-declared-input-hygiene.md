# Scaffold declared-input hygiene for consumers

Blocked by: none
Writes: internal/adopt/setup.go, internal/adopt/setup_test.go, internal/conformance/validity_checks_test.go, internal/conformance/registry/scope.go, internal/testreport/command.go, internal/testreport/check_test.go
Covers: LF2

## What to build

Ship the gitignored-declared-input check in the consumer gate scaffold. Reuse
the declared-input grammar and keep ordinary ignored files outside the check.

## Acceptance

- [ ] A linked consumer rejects a declared input that Git ignores.
- [ ] An undeclared ignored file remains allowed.
- [ ] Paths with spaces or glob characters are treated literally.

## Repair evidence

The detected-project removal mutation deleted the
`consumerGateHygieneCheck()` composition from the detected branch of
`setupGateScript`. The focused command built the mutated tree as the selected
executable and authenticated it through the system-test owner:

```sh
GOTOOLCHAIN=local GOMODCACHE=/home/mgibs/go/pkg/mod GOPATH=/home/mgibs/go GOCACHE=/tmp/bench-lf2-gocache /home/mgibs/.local/opt/go-v1.25.0/bin/go build -o /tmp/bench-lf2-no-detected-hygiene ./cmd/bench && GOTOOLCHAIN=local GOMODCACHE=/home/mgibs/go/pkg/mod GOPATH=/home/mgibs/go GOCACHE=/tmp/bench-lf2-gocache BENCH_KIT=/home/mgibs/.bench/worktrees/bench-2826441890/0c256629d0b059bab43a28a57c5e6312-b9ef50a510a022917fa9e3e2415d11ac BENCH_RUN_BINARY=/tmp/bench-lf2-no-detected-hygiene /home/mgibs/.local/opt/go-v1.25.0/bin/go test -tags=system ./internal/systemtest -run '^(TestDetectedProjectGateRejectsIgnoredDeclaredInput|TestSelectedExecutableComposition|TestRedTimeoutAndDescendantTeardown)$' -count=1 -v
```

The journey turned red at its consumer-visible assertion:

```text
=== RUN   TestDetectedProjectGateRejectsIgnoredDeclaredInput
    adoption_test.go:55: detected-project gate = (1, "", "go: warning: \"./...\" matched no packages\nno packages to test\n")
--- FAIL: TestDetectedProjectGateRejectsIgnoredDeclaredInput (0.86s)
FAIL
FAIL github.com/gibbonmi/bench/internal/systemtest 1.075s
```

The seeded-routing-input removal mutation deleted `BENCH_RUN_BINARY` from
`scaffoldGateInputs`. This focused command turned the independent seed
expectation red:

```sh
GOTOOLCHAIN=local GOMODCACHE=/home/mgibs/go/pkg/mod GOPATH=/home/mgibs/go GOCACHE=/tmp/bench-lf2-gocache /home/mgibs/.local/opt/go-v1.25.0/bin/go test ./internal/adopt -run '^TestSetupSeedsGateInputs$' -count=1 -v
```

```text
=== RUN   TestSetupSeedsGateInputs
    setup_test.go:57: seeded gate-inputs.json =
        {
          "schema": 1,
          "closure": "local",
          "environment": ["BENCH_HOME", "BENCH_KIT", "HOME"],
          "paths": [],
          "tools": ["bash", "basename", "dirname", "git", "readlink", "uname"]
        }

        want:
        {
          "schema": 1,
          "closure": "local",
          "environment": ["BENCH_HOME", "BENCH_KIT", "BENCH_RUN_BINARY", "HOME"],
          "paths": [],
          "tools": ["bash", "basename", "dirname", "git", "readlink", "uname"]
        }
--- FAIL: TestSetupSeedsGateInputs (0.38s)
FAIL
FAIL github.com/gibbonmi/bench/internal/adopt 0.378s
```
