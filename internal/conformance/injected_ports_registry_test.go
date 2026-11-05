package conformance

import "github.com/gibbonmi/bench/internal/conformance/injectedports"

// auditedPortPackages are the packages whose injected ports the registry covers, the
// audit's inventory, no wider. A package outside it that declares a port is out of scope
// until a reviewer binds it in. A package inside it that derives nothing fails closed.
var auditedPortPackages = []string{
	"internal/assessment",
	"internal/chargeevidence",
	"internal/gitguard",
	"internal/preflight",
	"internal/publication",
	"internal/subprocess/protection",
}

// injectedPortRegistry is the advertisement half: one row per port the derivation finds
// in auditedPortPackages. Rows are grouped by package and ordered as the derivation
// reports them. A diff against a re-derived inventory therefore reads straight down.
var injectedPortRegistry = []injectedports.Row{
	{Package: "internal/assessment", Port: "FileOps", TestFile: "internal/assessment/replacement_test.go", TestName: "TestAssessmentReplacement"},
	{
		Package: "internal/chargeevidence", Port: "ExportWrite",
		TestFile: "internal/preflight/evidencecmd/evidence_export_test.go", TestName: "TestEvidenceExportWritesSources",
	},
	{
		Package: "internal/chargeevidence", Port: "SourceAt",
		TestFile: "internal/preflight/evidencecmd/evidence_export_test.go", TestName: "TestEvidenceExportVerifiesBeforeWrite",
	},
	{
		Package: "internal/chargeevidence", Port: "StoreOptions",
		TestFile: "internal/systemtest/charge_evidence_test.go", TestName: "TestEvidenceInterruptedPublication",
	},
	{
		Package: "internal/gitguard", Port: "Checker",
		TestFile: "internal/gitguard/checker_junction_test.go", TestName: "TestClassifyRealCheckerResolvedComposition",
	},
	{
		Package: "internal/publication", Port: "Registry",
		Exempt: "a real-producer test needs a live npm registry, and NPMCLIRegistry is runbook-only: it never performs gate egress or touches a credential (the Registry port comment in internal/publication/registry.go states each adapter's gate coverage)",
	},
	{
		Package: "internal/preflight", Port: "reviewEvidenceObserver",
		TestFile: "internal/preflight/review_charge_test.go", TestName: "TestEvidenceReviewCollectors",
	},
	{
		Package: "internal/subprocess/protection", Port: "Probe",
		TestFile: "internal/subprocess/protection/observe_test.go", TestName: "TestObserveGroupPresent",
	},
}
