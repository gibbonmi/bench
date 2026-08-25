package racetests

// Test is one regression test the race phase must execute.
type Test struct {
	PackagePath string
	Name        string
}

// Tests names the regression tests that need race instrumentation because the ordinary
// suite cannot observe their concurrent failure modes.
var Tests = []Test{
	{PackagePath: "./internal/worktree", Name: "TestConcurrentCleanupRecordsOneTransaction"},
	{PackagePath: "./internal/worktree", Name: "TestParallelJourneysShareTheHarnessSafely"},
	{PackagePath: "./internal/worktree", Name: "TestParallelJourneysRecordEverySelection"},
	{PackagePath: "./internal/guards", Name: "TestScanTimeoutPreservesPartialRowsAndHonestCounts"},
	{PackagePath: "./internal/guards", Name: "TestScanEnumerationTimeoutUsesUnknownCounts"},
}
