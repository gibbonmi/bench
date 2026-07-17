package preflight

import (
	_ "embed"
	"encoding/json"
)

type Mode string
type Scope string
type Status string

const (
	ModeVerify        Mode   = "verify"
	ModePublish       Mode   = "publish"
	ScopePreflight    Scope  = "preflight"
	ScopeFocused      Scope  = "focused"
	StatusGreen       Status = "green"
	StatusRed         Status = "red"
	StatusNotRun      Status = "not_run"
	StatusInterrupted Status = "interrupted"
)

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

var registry = loadRegistry()

type phaseRegistry struct {
	Verify      []string `json:"verify"`
	PublishOnly []string `json:"publish_only"`
}

func loadRegistry() phaseRegistry {
	var value phaseRegistry
	if err := json.Unmarshal(registryJSON, &value); err != nil {
		panic("invalid embedded preflight registry: " + err.Error())
	}
	return value
}

func PhaseNames(mode Mode) []string {
	names := append([]string{}, registry.Verify...)
	if mode == ModePublish {
		names = append(names, registry.PublishOnly...)
	}
	return names
}

func phaseSummaries(results []Result) []PhaseSummary {
	out := make([]PhaseSummary, 0, len(results))
	for _, result := range results {
		out = append(out, PhaseSummary{Name: result.Name, Status: result.Status, ExitCode: result.ExitCode})
	}
	return out
}
