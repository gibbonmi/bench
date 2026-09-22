# FT297 author verification

Status: author acceptance satisfied; all native axes pass; landing pending
Base: 216ced3fb95ff10314dc72ab237609f14a6e4bab
Assignment: batch-ft297
Ticket: specs/ft297-build-conformance/tickets/pin-go-build-contracts.md
Author: ft297_author
Line: gpt-6-astra / ultra
Pre-review coherent attempts consumed: 1 of 3
Post-review repair cycles consumed: 1 of 2
Hardening cycles consumed: 0 of 1

## Current-source readiness

The Go owner derives the published executable path from `dist` and `bench`.
The resolver and wrapper each derive the same path.
The tree has four direct Go build or list calls through `exec.Command`.
The calls belong to freshness, testreport, and the command brief test.

The freshness diagnostic lacked the VCS-disable flag.
The other three calls already carried the flag.
The coordinator accepted that diagnostic repair within the approved invariant.

## Acceptance evidence

| Acceptance | Check | Result |
|---|---|---|
| Either shell path or the Go owner can drift | TestPublishedExecutablePathSourceFamily | Passed all three changed-source cases and each restore |
| Matching paths pass without a fixed destination | TestPublishedExecutablePathSourceFamily | Passed the shared out/tool destination |
| A missing flag fails | TestGoBuildVCSRejectsMissingFlag | Behavioral red before implementation; green after implementation |
| Literal and gate constant forms pass | TestGoBuildVCSCallFamily | Passed both accepted forms |
| Test files and import aliases remain covered | TestGoBuildVCSCallFamily | Passed alias, context, dot-import, and test-file cases |
| Malformed ambient Git metadata remains harmless | TestDigestIgnoresMalformedAmbientVCSMetadata | Passed; 0.093 seconds |

The first path regression returned an empty diagnostic list against its compiled stub.
The implemented check then returned the expected resolver mismatch.
The first VCS regression also failed against its compiled stub and passed after implementation.
The registered VCS check found the real missing flag before the diagnostic repair.
It named `internal/freshness/freshness_digest_test.go:166`.

## Commands and results

The focused command used the repository resource limits and executed every selected test without skips.

```sh
go test -count=1 -parallel 2 ./internal/conformance ./internal/freshness -run 'Test(PublishedExecutablePath|GoBuildVCS|ConformanceMetaBites|DigestIgnoresMalformedAmbientVCSMetadata)' -v
```

The conformance package took 0.037 seconds.
The freshness package took 0.093 seconds.
The registry metadata test passed all eleven mutation cases.
The structure growth check passed against the recorded base.

The assignment binary came from `bench worktree build batch-ft297`.
Both named checks passed through that binary with zero skips.
The path check took 5 milliseconds.
The VCS check took 196 milliseconds.
The VCS check uses catch-all input selection because it also reads test and fixture source files.

## Mutation evidence

Each probe used `GOFLAGS=-p=4 -parallel=2` and the assignment binary.
Each probe reported `bit`, one executed test, zero skips, and `restored=yes`.

```sh
./dist/bench probe .bench/lib/resolve-bench.sh --swap '$1/dist/bench' --with '$1/wrong/bench' --check published-executable-path
./dist/bench probe internal/freshness/freshness_buildinputs.go --swap '-buildvcs=false' --with '-buildvcs=true' --check go-build-vcs
./dist/bench probe internal/testreport/selection.go --omit '"-buildvcs=false", ' --check go-build-vcs
```

The first two probes are author-selected swaps at separate production sites.
The coordinator selected the last probe as the independent omission.
The retained author executed it at the coordinator's request.
Its kind and site differ from both author probes.
The path failure named the resolver and both path values.
The VCS failures named the mutated files and their missing flag requirement.

## Scope and headroom

The ordered check table moved to `internal/conformance/registry/checks.go`.
The binding type and methods moved to `internal/conformance/check_bindings_test.go`.
The executable map stayed in its existing fixture-advertised file.
The moves preserve one owner per registry fact and keep existing file budgets unchanged.
The coordinator approved the fence expansion, and the author captured it with `bench learning`.

The tickets-only preflight refused because this light-path ticket has no `spec.md`.
The approved light path uses the explicit base and this ticket as its source.

## Remaining work

The coordinator owns native Standards, Spec, and Coverage reviews.
The coordinator also owns the serialized commit, full gate, and local landing.
The author retains all production, test, probe, and repair work.
No full gate or commit ran during this author pass.

## Initial native reviews

Source pair: 216ced3fb95ff10314dc72ab237609f14a6e4bab..0d3f61f986a96c233872025290b78351fecd877b

| Axis | Findings | Worst | Status | Confidence |
|---|---:|---|---|---:|
| Standards | 1 | ST1 | claimed | 9 |
| Spec | 0 | none | claimed | 9 |
| Coverage | 1 | COV1 | claimed | 10 |

ST1: auto-fix; accepted.
The import resolver duplicates the architecture scanner's import-alias parser.
AGENTS.md requires one source per fact, including parsers.
Both scanners must use one import resolver.

COV1: auto-fix; accepted.
Both checks can read a FIFO before they reject its special-file mode.
The project profile requires rejection of special files before a read.
Both checks must reject each nonregular source before they read its bytes.

These accepted findings share repair cycle one.
The retained author remains gpt-6-astra at ultra effort.
The repair keeps the existing two-cycle allowance.

## Repair cycle one

Initial review record commit: 9721cf8ea39e4619447224755a5685c9c3876a9c

ST1 author result: verified; confidence 9.
Both scanners now call `importedPackage` for import resolution.
The architecture scanner no longer keeps its own alias parser.
The existing architecture census and the complete build-contract fixture family pass.

COV1 author result: verified; confidence 9.
Every new source reader uses `bounds.ClassifyNoFollow` before it receives source bytes.
The classifier owns the special-file refusal and the bounded read.
The shared fixture planter owns FIFO capability handling.
Four FIFO fixtures cover the Go owner, both shell sources, and the Go source sweep.
Each fixture returns the specific wrong-type refusal without a writer or a wait.

The focused repair suite passed without skips in 0.296 seconds.
The structure growth check passed against the frozen review tip.
The skip-ownership check also passed without skips.

```sh
go test -count=1 -parallel 2 ./internal/conformance -run 'Test(PublishedExecutablePath|GoBuildVCS|BuildContractChecksRefuseSpecialFiles|BranchNativeArchitectureCensus|ConformanceMetaBites)' -v
./dist/bench test --check skip-ownership
./dist/bench probe internal/conformance/build_contracts_test.go --swap 'if classified.State != bounds.StateParsed && classified.State != bounds.StateEmpty {' --with 'if false && classified.State != bounds.StateParsed && classified.State != bounds.StateEmpty {' --package ./internal/conformance --run '^TestBuildContractChecksRefuseSpecialFiles$'
```

The probe bypassed the classifier-refusal branch.
It reported `bit`, four failed fixture cases, zero skips, and `restored=yes`.
The specific wrong-type expectations independently catch that refusal omission.
The retained author consumed one of the two permitted repair cycles.
All three native axes must reaffirm the repaired source before the full gate.

## Final native review results

Reviewed base: 216ced3fb95ff10314dc72ab237609f14a6e4bab
Reviewed source: 6987b214fec96b48c69f5a297831bfaeace66cf6

The coordinator supplied the native terminal excerpts below.
Each SHA-256 digest covers the exact UTF-8 excerpt bytes.
All three axes used separate native sessions at gpt-5.6-sol with high effort.
The review claims remain claimed; they do not replace the executed author checks.

| Axis | Final findings | Outcome | Confidence |
|---|---:|---|---:|
| Standards | 0 | pass | 10 |
| Spec | 0 | pass | 9 |
| Coverage | 0 | pass | 10 |

The final Standards result supersedes ST1 after the shared resolver repair.
The final Coverage result supersedes COV1 after all four reader repairs.
The initial findings and their accepted dispositions remain in this record.

This light-path ticket has no spec-backed preflight source digest or plan digest.
Those unavailable fields are omitted rather than invented.
The typed terminal metadata records the frozen commit pairs.
It does not claim spec-backed record-schema verification.

```json
[
  {
    "id": "ft297-standards-initial",
    "performer": "/root/ft326_standards",
    "role": "independent-review",
    "model": "gpt-5.6-sol",
    "effort": "high",
    "state": "completed",
    "outcome": "findings",
    "native_ref": {
      "ref": "/root/ft326_standards",
      "digest": "sha256:022dc9f80c387b5db9a6d90f4d03f8d1691ed5ac1a46e6e7c965dea2c80c7548",
      "excerpt": "ST1(auto-fix,c9) new importedPackage at build_contracts_test95-109 duplicates exec alias derivation scanArchitectureGo ordinary_build_census_test206-217; AGENTS34-47 parsers single-source."
    },
    "axis": "Standards",
    "base": "216ced3fb95ff10314dc72ab237609f14a6e4bab",
    "tip": "0d3f61f986a96c233872025290b78351fecd877b",
    "finding_ids": [
      "ST1"
    ],
    "supersedes": []
  },
  {
    "id": "ft297-standards-final",
    "performer": "/root/ft326_standards",
    "role": "independent-review",
    "model": "gpt-5.6-sol",
    "effort": "high",
    "state": "completed",
    "outcome": "pass",
    "native_ref": {
      "ref": "/root/ft326_standards",
      "digest": "sha256:acda3e81976635df416fb10871a5534bf9d1985c564d60ca026fc5022f87577e",
      "excerpt": "ST1sharedimportedPackage soleowner, existingboundsClassifyNoFollow/hostileSkillPlantersreused no duplicates."
    },
    "axis": "Standards",
    "base": "216ced3fb95ff10314dc72ab237609f14a6e4bab",
    "tip": "6987b214fec96b48c69f5a297831bfaeace66cf6",
    "finding_ids": [],
    "supersedes": [
      "ft297-standards-initial"
    ]
  },
  {
    "id": "ft297-spec-initial",
    "performer": "/root/ft326_spec",
    "role": "independent-review",
    "model": "gpt-5.6-sol",
    "effort": "high",
    "state": "completed",
    "outcome": "pass",
    "native_ref": {
      "ref": "/root/ft326_spec",
      "digest": "sha256:526852fbe63cbf2eabba53b6291fe6ebbbfc32c9093d1fbbdd4589a764a116c4",
      "excerpt": "Spec0 findingsc9."
    },
    "axis": "Spec",
    "base": "216ced3fb95ff10314dc72ab237609f14a6e4bab",
    "tip": "0d3f61f986a96c233872025290b78351fecd877b",
    "finding_ids": [],
    "supersedes": []
  },
  {
    "id": "ft297-spec-final",
    "performer": "/root/ft326_spec",
    "role": "independent-review",
    "model": "gpt-5.6-sol",
    "effort": "high",
    "state": "completed",
    "outcome": "pass",
    "native_ref": {
      "ref": "/root/ft326_spec",
      "digest": "sha256:956a290025539e48656f07c309d34eb8a8d0f941f9e500a2ecbcdd94f269e3ef",
      "excerpt": "all6acceptancesunchanged; exact4readsiteswrongtypefixtures,approvedfence."
    },
    "axis": "Spec",
    "base": "216ced3fb95ff10314dc72ab237609f14a6e4bab",
    "tip": "6987b214fec96b48c69f5a297831bfaeace66cf6",
    "finding_ids": [],
    "supersedes": [
      "ft297-spec-initial"
    ]
  },
  {
    "id": "ft297-coverage-initial",
    "performer": "/root/ft326_coverage",
    "role": "independent-review",
    "model": "gpt-5.6-sol",
    "effort": "high",
    "state": "completed",
    "outcome": "findings",
    "native_ref": {
      "ref": "/root/ft326_coverage",
      "digest": "sha256:952943c820041ec5e066329da19476fcae7d0e346b05a78c21dfc92e1fdb869a",
      "excerpt": "COV1(auto-fix,c10): pathchecker parser.ParseFile(owner)/os.ReadFile(shell), Go walk os.ReadFile each .go can block on FIFO; profile226-228 requires reject specialfiles before read."
    },
    "axis": "Coverage",
    "base": "216ced3fb95ff10314dc72ab237609f14a6e4bab",
    "tip": "0d3f61f986a96c233872025290b78351fecd877b",
    "finding_ids": [
      "COV1"
    ],
    "supersedes": []
  },
  {
    "id": "ft297-coverage-final",
    "performer": "/root/ft326_coverage",
    "role": "independent-review",
    "model": "gpt-5.6-sol",
    "effort": "high",
    "state": "completed",
    "outcome": "pass",
    "native_ref": {
      "ref": "/root/ft326_coverage",
      "digest": "sha256:e4b3eff960a7288f8b7cb797f0411f7e9c5c833bee2ef08b7696d09cc512c40b",
      "excerpt": "COV1all4FIFOreadsites classifybeforebytes, bypassswapbitall4/restorednoskip."
    },
    "axis": "Coverage",
    "base": "216ced3fb95ff10314dc72ab237609f14a6e4bab",
    "tip": "6987b214fec96b48c69f5a297831bfaeace66cf6",
    "finding_ids": [],
    "supersedes": [
      "ft297-coverage-initial"
    ]
  }
]
```

## Final author acceptance

The retained author confirms all six ticket acceptance criteria as satisfied.
The original acceptance claims retain confidence 9.
The repair evidence above supplies the executed tests and mutation results.
The final native results close both accepted findings.
One of two repair cycles is consumed, and no hardening cycle is consumed.

This record-only update changes no production or test bytes.
The coordinator owns the remaining whole-project gate and local landing.
