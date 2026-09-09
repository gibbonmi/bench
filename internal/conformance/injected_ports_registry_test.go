package conformance

// auditedPortPackages are the packages whose injected ports the registry covers, the
// audit's inventory, no wider. A package outside it that declares a port is out of scope
// until a reviewer binds it in. A package inside it that derives nothing fails closed.
var auditedPortPackages = []string{
	"internal/gitguard",
	"internal/preflight",
	"internal/publication",
}

// This block names the five failure modes, each with its own message, so a red names its
// cause without archaeology. They are constants because the canary fixture's EXPECT and
// the bite proof both quote them. A message typed a second time drifts from the one that
// fires.
const (
	unregisteredPortMessage   = "injected port has no registry row"
	missingPortTestMessage    = "injected port registry row names a test the tree does not declare"
	emptyPortExemptionMessage = "injected port registry row is exempt with an empty reason"
	zeroPortInventoryMessage  = "injected port derivation found no ports"
	orphanPortRowMessage      = "injected port registry row names a port the derivation no longer reports"
)

// injectedPortRow is one derived port's disposition. The disposition is either a real-
// producer test or an exemption whose reason someone had to write down. A real-producer
// test drives the actual producer through the consuming surface. A row carrying neither
// is the empty-exemption red.
type injectedPortRow struct {
	pkg, port          string
	testFile, testName string
	exempt             string
}

// injectedPortRegistry is the advertisement half: one row per port the derivation finds
// in auditedPortPackages. Rows are grouped by package and ordered as the derivation
// reports them. A diff against a re-derived inventory therefore reads straight down.
var injectedPortRegistry = []injectedPortRow{
	{
		pkg: "internal/gitguard", port: "Checker",
		testFile: "internal/gitguard/checker_junction_test.go", testName: "TestClassifyRealCheckerResolvedComposition",
	},
	{
		pkg: "internal/publication", port: "Registry",
		exempt: "the only adapter without gate coverage is NPMCLIRegistry, which is runbook-only: the gate drives FixtureRegistry against the hermetic offline registry, and no NPMCLIRegistry path performs gate egress or touches a credential (internal/publication/registry.go:5-11)",
	},
	{
		pkg: "internal/preflight", port: "reviewEvidenceObserver",
		testFile: "internal/preflight/review_charge_test.go", testName: "TestReviewChargeSharedEvidence",
	},
}
