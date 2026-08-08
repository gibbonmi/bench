# Shared fixture staged-binary implementation evidence

This receipt records the independently green copy-based intermediate commit. The dependent `reduce-gate-fixture-materialization.md` ticket supersedes only its per-root copy after measuring the package residual; the process template, explicit re-seal boundary, sticky error, and teardown remain.

## Focused resource receipt

The repaired implementation ran:

`/usr/bin/time -v env GOMAXPROCS=2 go test -p 1 -parallel 1 -count=1 -run '^(TestKitShapedFixtureCarriesBuildAndCanary|TestKitShapedFixtureBinaryIsSealed|TestKitShapedFixturesPublishIndependentTemplateCopies)$' ./internal/gate`

It was green in 0.794 package seconds and 1.41 wall seconds, with 103,040 KiB
maximum RSS and 30,376 filesystem-output blocks. The same four-fixture command
after restoring a real link in every initial seal was green in 2.51 wall
seconds, with 214,972 KiB maximum RSS and 61,584 output blocks. The shared
implementation writes 49.3 percent of the mutation's output, inside SFB13's
60-percent ceiling.

The package-wide serial control did not meet the original provisional resource
targets: it was green in 127.602 package seconds and 128.00 wall seconds, with
245,500 KiB maximum RSS and 1,267,992 output blocks. The implementation therefore
does not claim those targets. They moved to the dependent
`reduce-gate-fixture-materialization.md` ticket.

## Compatibility receipts

- `GOMAXPROCS=2 go test -race -p 1 -parallel 2 -count=1 -timeout 600s ./internal/gate` was green in 100.679 package seconds and 104.82 wall seconds, with 240,280 KiB maximum RSS and 1,324,344 output blocks. This also exercised the two-core coordination path.
- `BENCH_KIT=/nonexistent GOMAXPROCS=2 go test -p 1 -parallel 2 -count=1 -timeout 600s ./internal/gate` was green in 94.303 package seconds and 95.83 wall seconds, with 244,900 KiB maximum RSS and 1,302,096 output blocks. This non-race full-package receipt also covers the two-core compatibility criterion without a redundant third package sweep.
- The final focused functional set covering template construction, publication, isolation, sticky failure, cleanup, authored builds, and planted builds was green in 1.420 seconds before the bounded mutations.

## Mutation receipts

- Directly publishing the process template made the second fixture construction fail because the first publication consumed the staging path.
- Replacing the ordinary copy with a hardlink made the two-fixture isolation witness fail freshness verification after the first executable was truncated.
- Routing a changed-source green build through the immutable template made `TestGreenBuildAttestsItsOwnBinary` fail because its executable digest did not change.
- Bypassing `sync.Once` made `TestKitShapedFixtureTemplateFailureIsSticky` fail in 0.120 seconds with two distinct retained-error values.
- Replacing template teardown with `os.RemoveAll("")` made the same parent-observed lifetime witness fail in 0.068 seconds because the child template directory remained.
- Restoring the exact subject made the lifetime witness green in 0.067 seconds.

## Debug receipt

The required second-loop `$bench-debug` pass observed four test processes, 101
root-owned template copies, and 24 build attempts in a serial package run. The
parent owned every copy and 23 successful behavior-owning builds; the only child
build was the intentional sticky-error witness. A throwaway exact helper census
wrote 580,520 blocks, while the instrumented full package wrote 1,302,464 and a
no-test control wrote 11,976. Roughly 722,000 blocks therefore remain in other
test materialization work. All diagnostic instrumentation and temporary census
code were removed before this evidence was recorded.

## Cross-harness review

Claude Opus reviewed Standards, Spec, and Coverage. Its first pass found an
implicit constructor/re-seal boundary, string-only error comparison, and no
in-tree teardown witness. The repair introduced separate initial-seal and
always-real re-seal APIs, error identity, and parent-observed process cleanup.
The second pass verified all three behavior repairs and left one stale helper
comment. The final comment states only the stable contract: the build helper
does not publish and its caller owns any matching publication. A caller audit
confirmed direct publication remains intentional at several test-specific
sites, so the comment does not enumerate routes.
