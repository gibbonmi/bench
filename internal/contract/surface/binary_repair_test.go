package surface

import (
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

func TestBinaryRepairContracts(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	contract.RunParallel(t, "binary repair download-and-run contract failed", testRepairDownloadsAndRuns)
	contract.RunParallel(t, "binary repair bad-integrity contract failed", testRepairRefusesBadHash)
	contract.RunParallel(t, "binary repair malformed-tar contract failed", testRepairRefusesMalformedTar)
	contract.RunParallel(t, "binary repair offline fallback contract failed", testRepairOffline)
	contract.RunParallel(t, "binary repair plumbing exclusion contract failed", testRepairSkipsPlumbing)
	contract.RunParallel(t, "binary repair announcement contract failed", testRepairAnnounces)
	contract.RunParallel(t, "binary repair idempotency contract failed", testRepairIdempotent)
	contract.RunParallel(t, "binary repair version-keyed cache contract failed", testRepairVersionKeyed)
	contract.RunParallel(t, "binary repair no-node contract failed", testRepairNoNode)
	contract.RunParallel(t, "binary repair disabled contract failed", testRepairDisabled)
	contract.RunParallel(t, "binary repair BENCH_OFFLINE exact-value contract failed", testRepairBenchOfflineExact)
	contract.RunParallel(t, "binary repair opt-in exact-value contract failed", testRepairOptInExact)
	contract.RunParallel(t, "binary repair suppression precedence contract failed", testRepairSuppressionPrecedence)
	contract.RunParallel(t, "binary repair torn-cache contract failed", testRepairReplacesTornCache)
	contract.RunParallel(t, "linked manifest repair contract failed", testRepairReadsLinkedManifestWithoutNewline)
	contract.RunParallel(t, "malformed manifest version escaped repair cache", testRepairRejectsMalformedVersion)
	contract.RunParallel(t, "interrupted repair promoted partial cache", testRepairInterruptedPromotion)
	contract.RunParallel(t, "fresh clone repair required ambient tooling", testRepairMinimalPortablePath)
	contract.RunParallel(t, "repair explicit-default contract failed", testRepairExplicitDefault)
	contract.RunParallel(t, "repair explicit subcommand contract failed", testRepairSubcommand)
	contract.RunParallel(t, "repair argument contract failed", testRepairArguments)
	contract.RunParallel(t, "repair prune contract failed", testRepairPrune)
	contract.RunParallel(t, "repair pin manifest fail-closed contract failed", testRepairPinManifestFailures)
	contract.RunParallel(t, "repair resource bounds contract failed", testRepairResourceBounds)
	contract.RunParallel(t, "repair losing-racer cleanup contract failed", testRepairLosingRacerPreservesWinner)
	contract.RunParallel(t, "repair earliest interrupt cleanup contract failed", testRepairEarliestInterrupt)
}
