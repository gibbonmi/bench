package preflight

import (
	"context"

	"github.com/gibbonmi/bench/internal/releaseevidence"
)

type Mode = releaseevidence.Mode
type Scope = releaseevidence.Scope
type Status = releaseevidence.Status
type Profile = releaseevidence.Profile
type RunEvidence = releaseevidence.RunEvidence
type Failure = releaseevidence.Failure
type Record = releaseevidence.Record
type Identity = releaseevidence.Identity
type PhaseSummary = releaseevidence.PhaseSummary
type Manifest = releaseevidence.Manifest
type Result = releaseevidence.Result
type PhaseDefinition = releaseevidence.PhaseDefinition
type Requirement = releaseevidence.Requirement
type PackageEvidence = releaseevidence.PackageEvidence
type requirementRegistry = releaseevidence.Registry
type releaseIndex = releaseevidence.Index
type requirementStatus = releaseevidence.RequirementStatus
type releaseIntentError = releaseevidence.ReleaseIntentError

const (
	ModeVerify        = releaseevidence.ModeVerify
	ModePublish       = releaseevidence.ModePublish
	ScopePreflight    = releaseevidence.ScopePreflight
	ScopeFocused      = releaseevidence.ScopeFocused
	StatusGreen       = releaseevidence.StatusGreen
	StatusRed         = releaseevidence.StatusRed
	StatusNotRun      = releaseevidence.StatusNotRun
	StatusInterrupted = releaseevidence.StatusInterrupted
	ProfilePublic     = releaseevidence.ProfilePublic
	ProfileBank       = releaseevidence.ProfileBank
)

var requirements = releaseevidence.RequirementsRegistry()

func Requirements() []Requirement   { return releaseevidence.Requirements() }
func PhaseNames(mode Mode) []string { return releaseevidence.PhaseNames(mode) }
func phaseDefinition(name string) (PhaseDefinition, bool) {
	return releaseevidence.PhaseDefinitionFor(name)
}
func phaseSummaries(results []Result) []PhaseSummary { return releaseevidence.PhaseSummaries(results) }
func packageEvidenceRegistry() []PackageEvidence     { return releaseevidence.PackageEvidenceRegistry() }
func terminalStatus(results []Result) Status         { return releaseevidence.TerminalStatus(results) }
func contains(items []string, want string) bool      { return releaseevidence.Contains(items, want) }
func FinalizeEvidence(ctx context.Context, root string, run RunEvidence) error {
	return releaseevidence.FinalizeEvidence(ctx, root, run)
}
func PromoteEvidence(root string, mode Mode, results []Result, manifest Manifest) error {
	return releaseevidence.PromoteEvidence(root, mode, results, manifest)
}
func readPackageVersion(root string) (string, error) { return releaseevidence.ReadPackageVersion(root) }
func readRegular(path string) ([]byte, error)        { return releaseevidence.ReadRegular(path) }
func setArchiveMemberLimitForTesting(limit int64) func() {
	return releaseevidence.SetArchiveMemberLimitForTesting(limit)
}
func validateTarballForTesting(data []byte) error {
	return releaseevidence.ValidateTarballForTesting(data)
}
func setExchangeForTesting(exchange func(string, string) error) func() {
	return releaseevidence.SetExchangeForTesting(exchange)
}
func atomicExchangeForTesting(left, right string) error {
	return releaseevidence.AtomicExchangeForTesting(left, right)
}
