//go:build system

package systemtest

// ownerRepositories is the serial reason of a test whose subject is the owner's own
// repositories.
const ownerRepositories = "the owner's repositories are the test's subject"

// systemSerialTests is the one list of the package's serial tests. Each row names a
// top-level test that calls no t.Parallel() and the reason it stays serial. Go runs these
// tests one at a time before every parallel test, so a serial test may name a ledger
// field or bind process state. TestSystemParallelCensus and TestSystemSerialTableIsExact
// own the limits on each row and on the process state that a test binds.
var systemSerialTests = map[string]string{
	"TestOwnerSelectionChild":           ownerRepositories,
	"TestSelectedExecutableComposition": ownerRepositories,
}
