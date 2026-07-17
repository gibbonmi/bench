package preflight

import (
	_ "embed"
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
	Verify      []string                   `json:"verify"`
	PublishOnly []string                   `json:"publish_only"`
	Phases      map[string]PhaseDefinition `json:"phases"`
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
}

type PackageEvidence struct {
	Path   string `json:"path"`
	Schema string `json:"schema"`
	Mode   string `json:"mode"`
}

type ToolchainRequirement struct {
	Name        string   `json:"name"`
	VersionArgv []string `json:"version_argv"`
	Flags       []string `json:"flags"`
}

type ComponentManifestSchema struct {
	SchemaVersion   int      `json:"schema_version"`
	Path            string   `json:"path"`
	RootFields      []string `json:"root_fields"`
	ComponentFields []string `json:"component_fields"`
	TargetFields    []string `json:"target_fields"`
	FileFields      []string `json:"file_fields"`
}

type requirementRegistry struct {
	SchemaVersion     int                     `json:"schema_version"`
	PackageEvidence   []PackageEvidence       `json:"package_evidence"`
	Toolchains        []ToolchainRequirement  `json:"toolchains"`
	ComponentManifest ComponentManifestSchema `json:"component_manifest"`
	Records           []Requirement           `json:"records"`
}

func loadRegistry() phaseRegistry {
	var value phaseRegistry
	if err := decodeStrict(registryJSON, &value); err != nil {
		panic("invalid embedded preflight registry: " + err.Error())
	}
	for _, name := range append(append([]string{}, value.Verify...), value.PublishOnly...) {
		if value.Phases[name].Handler == "" {
			panic("missing preflight phase definition: " + name)
		}
	}
	return value
}

var requirements = loadRequirements()

func loadRequirements() requirementRegistry {
	var value requirementRegistry
	if err := decodeStrict(requirementsJSON, &value); err != nil {
		panic("invalid embedded requirement registry: " + err.Error())
	}
	if err := validateRequirementRegistry(value); err != nil {
		panic("invalid embedded requirement registry: " + err.Error())
	}
	return value
}

func packageEvidenceRegistry() []PackageEvidence {
	return append([]PackageEvidence(nil), requirements.PackageEvidence...)
}

func Requirements() []Requirement {
	return append([]Requirement(nil), requirements.Records...)
}

func PhaseNames(mode Mode) []string {
	names := append([]string{}, registry.Verify...)
	if mode == ModePublish {
		names = append(names, registry.PublishOnly...)
	}
	return names
}

func phaseDefinition(name string) (PhaseDefinition, bool) {
	definition, ok := registry.Phases[name]
	return definition, ok
}

func phaseSummaries(results []Result) []PhaseSummary {
	out := make([]PhaseSummary, 0, len(results))
	for _, result := range results {
		out = append(out, PhaseSummary{Name: result.Name, Status: result.Status, ExitCode: result.ExitCode})
	}
	return out
}
