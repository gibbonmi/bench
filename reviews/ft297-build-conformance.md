# FT297 author verification

Status: author verified; native reviews pending
Base: 216ced3fb95ff10314dc72ab237609f14a6e4bab
Assignment: batch-ft297
Ticket: specs/ft297-build-conformance/tickets/pin-go-build-contracts.md
Author: ft297_author
Line: gpt-6-astra / ultra
Pre-review coherent attempts consumed: 1 of 3
Post-review repair cycles consumed: 0 of 2
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
