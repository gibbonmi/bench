package releaseevidence

import (
	_ "embed"
	"sort"
)

type Mode string
type Scope string
type Status string

const (
	ModeVerify        Mode    = "verify"
	ModePublish       Mode    = "publish"
	ScopePreflight    Scope   = "preflight"
	ScopeFocused      Scope   = "focused"
	StatusGreen       Status  = "green"
	StatusRed         Status  = "red"
	StatusNotRun      Status  = "not_run"
	StatusInterrupted Status  = "interrupted"
	ProfilePublic     Profile = "public"
	ProfileBank       Profile = "bank"
)

type Profile string

type RunEvidence struct {
	Mode     Mode
	Scope    Scope
	Identity Identity
	Profile  Profile
	Phases   []Result
}

type Failure struct {
	Kind    string `json:"kind"`
	Message string `json:"message"`
}

type Record struct {
	SchemaVersion int      `json:"schema_version"`
	Phase         string   `json:"phase"`
	Mode          Mode     `json:"mode"`
	Status        Status   `json:"status"`
	ExitCode      *int     `json:"exit_code"`
	Error         *Failure `json:"error"`
}

type Identity struct {
	Tag              *string `json:"tag"`
	PackageVersion   *string `json:"package_version"`
	SourceCommit     *string `json:"source_commit"`
	BinaryVersion    *string `json:"binary_version"`
	ChangelogHeading *string `json:"changelog_heading"`
	Toolchain        *string `json:"toolchain"`
}

type PhaseSummary struct {
	Name     string `json:"name"`
	Status   Status `json:"status"`
	ExitCode *int   `json:"exit_code"`
}

type Manifest struct {
	SchemaVersion int            `json:"schema_version"`
	Mode          Mode           `json:"mode"`
	Scope         Scope          `json:"scope"`
	Status        Status         `json:"status"`
	Identity      Identity       `json:"identity"`
	Phases        []PhaseSummary `json:"phases"`
}

type Result struct {
	Name     string
	Status   Status
	ExitCode *int
	Failure  *Failure
}

//go:embed registry.json
var registryJSON []byte

//go:embed requirements.json
var requirementsJSON []byte

var registry = loadRegistry()

type phaseRegistry struct {
	Verify         []string                   `json:"verify"`
	PublishBefore  []string                   `json:"publish_before"`
	PublishOnly    []string                   `json:"publish_only"`
	EvidenceInputs []string                   `json:"evidence_inputs"`
	Phases         map[string]PhaseDefinition `json:"phases"`
}

type PhaseDefinition struct {
	Handler        string   `json:"handler"`
	Argv           []string `json:"argv"`
	Tools          []string `json:"tools"`
	Inputs         []string `json:"inputs"`
	Requires       []string `json:"requires"`
	ExactToolchain bool     `json:"exact_toolchain"`
}

type Requirement struct {
	Key          string    `json:"key"`
	Owner        string    `json:"owner"`
	Schema       string    `json:"schema"`
	Profiles     []Profile `json:"profiles"`
	Requiredness string    `json:"requiredness"`
	Path         string    `json:"path"`
	Producer     bool      `json:"producer,omitempty"`
	PackageMode  string    `json:"package_mode,omitempty"`
}

type PackageEvidence struct {
	Path   string `json:"path"`
	Schema string `json:"schema"`
	Mode   string `json:"mode"`
}

type ToolchainRequirement struct {
	Name        string              `json:"name"`
	VersionArgv []string            `json:"version_argv"`
	Operations  map[string][]string `json:"operations"`
}

type ComponentManifestSchema struct {
	SchemaVersion   int                     `json:"schema_version"`
	Path            string                  `json:"path"`
	RootFields      ComponentRootFields     `json:"root_fields"`
	ComponentFields ComponentIdentityFields `json:"component_fields"`
	TargetFields    ComponentTargetFields   `json:"target_fields"`
	FileFields      ComponentFileFields     `json:"file_fields"`
}

type ComponentRootFields struct {
	SchemaVersion string `json:"schema_version"`
	Component     string `json:"component"`
	Files         string `json:"files"`
}
type ComponentIdentityFields struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Target  string `json:"target"`
}
type ComponentTargetFields struct {
	OS   string `json:"os"`
	Arch string `json:"arch"`
}
type ComponentFileFields struct {
	Path   string `json:"path"`
	Mode   string `json:"mode"`
	Size   string `json:"size"`
	SHA256 string `json:"sha256"`
}

type Registry struct {
	SchemaVersion     int                     `json:"schema_version"`
	Toolchains        []ToolchainRequirement  `json:"toolchains"`
	ComponentManifest ComponentManifestSchema `json:"component_manifest"`
	Records           []Requirement           `json:"records"`
}

func loadRegistry() phaseRegistry {
	var value phaseRegistry
	if err := decodeStrict(registryJSON, &value); err != nil {
		panic("invalid embedded preflight registry: " + err.Error())
	}
	allPhases := append(append(append([]string{}, value.PublishBefore...), value.Verify...), value.PublishOnly...)
	for _, name := range allPhases {
		if value.Phases[name].Handler == "" {
			panic("missing preflight phase definition: " + name)
		}
	}
	return value
}

var requirements = loadRequirements()

func SetRequirementsForTesting(value Registry) func() {
	previous := requirements
	requirements = value
	return func() { requirements = previous }
}

func loadRequirements() Registry {
	var value Registry
	if err := decodeStrict(requirementsJSON, &value); err != nil {
		panic("invalid embedded requirement registry: " + err.Error())
	}
	if err := validateRequirementRegistry(value); err != nil {
		panic("invalid embedded requirement registry: " + err.Error())
	}
	return value
}

func RequirementsRegistry() Registry { return requirements }

func PackageEvidenceRegistry() []PackageEvidence {
	var evidence []PackageEvidence
	for _, record := range requirements.Records {
		if record.PackageMode != "" {
			evidence = append(evidence, PackageEvidence{Path: record.Path, Schema: record.Schema, Mode: record.PackageMode})
		}
	}
	return evidence
}

func Requirements() []Requirement {
	return append([]Requirement(nil), requirements.Records...)
}

func PhaseNames(mode Mode) []string {
	names := append([]string{}, registry.Verify...)
	if mode == ModePublish {
		names = append(append(append([]string{}, registry.PublishBefore...), names...), registry.PublishOnly...)
	}
	return names
}

func PhaseDefinitionFor(name string) (PhaseDefinition, bool) {
	definition, ok := registry.Phases[name]
	return definition, ok
}

func releaseInputPaths() []string {
	seen := map[string]bool{}
	paths := make([]string, 0, len(registry.EvidenceInputs))
	add := func(values []string) {
		for _, value := range values {
			if !seen[value] {
				seen[value] = true
				paths = append(paths, value)
			}
		}
	}
	add(registry.EvidenceInputs)
	allPhases := append(append(append([]string{}, registry.PublishBefore...), registry.Verify...), registry.PublishOnly...)
	for _, name := range allPhases {
		add(registry.Phases[name].Inputs)
	}
	sort.Strings(paths)
	return paths
}

func PhaseSummaries(results []Result) []PhaseSummary {
	out := make([]PhaseSummary, 0, len(results))
	for _, result := range results {
		out = append(out, PhaseSummary{Name: result.Name, Status: result.Status, ExitCode: result.ExitCode})
	}
	return out
}
