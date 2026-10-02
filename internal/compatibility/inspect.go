// Package compatibility owns the facts and verdicts shared by compatibility consumers.
package compatibility

import (
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/toon"
)

// Interface names the harness interface selected by the caller.
type Interface string

const (
	CodexCLI     Interface = "codex-cli"
	CodexDesktop Interface = "codex-desktop"
)

func interfaces() []Interface {
	return []Interface{CodexCLI, CodexDesktop}
}

// InterfaceOperand renders the accepted interface vocabulary for command usage.
func InterfaceOperand() string {
	names := []string{}
	for _, selected := range interfaces() {
		names = append(names, string(selected))
	}
	return "<" + strings.Join(names, "|") + ">"
}

// ParseInterface accepts the complete compatibility interface vocabulary.
func ParseInterface(value string) (Interface, bool) {
	for _, selected := range interfaces() {
		if string(selected) == value {
			return selected, true
		}
	}
	return "", false
}

// Fact is one observed value and the source that supplied it.
type Fact struct {
	Value  string
	Source string
}

// Context holds the observation context that determines whether evidence can be reused.
type Context struct {
	Interface         Interface
	Repository        Fact
	Environment       Fact
	ConfigurationHome Fact
	ActiveRuntime     Fact
	LauncherVersion   Fact
	PolicyProvenance  Fact
}

// FileFact retains every bounded file disposition, including absence and emptiness.
type FileFact struct {
	Path   string
	Source string
	State  bounds.FileState
	Reason string
	digest string
}

// ReadFile classifies one compatibility input without opening a special file.
func ReadFile(path, source string) FileFact {
	classified := bounds.Classify(path, bounds.ControlRecordLimit)
	fact := FileFact{Path: path, Source: source, State: classified.State, Reason: classified.Reason}
	if classified.State == bounds.StateParsed || classified.State == bounds.StateEmpty {
		fact.digest = fmt.Sprintf("%x", sha256.Sum256(classified.Data))
	}
	return fact
}

// Asset is one required repository asset and its supported restoration action.
type Asset struct {
	Name          string
	File          FileFact
	RestoreAction string
}

// Input is the complete typed input to a local compatibility inspection.
type Input struct {
	Context       Context
	Configuration FileFact
	Assets        []Asset
	GlobalBench   bool
	ByPathAction  string
	HookDeclared  bool
	HookAction    string
	Live          []LiveRow
}

// CheckState is the local diagnostic vocabulary.
type CheckState string

const (
	StateOK          CheckState = "ok"
	StateFailed      CheckState = "failed"
	StateUnknown     CheckState = "unknown"
	StateNotRequired CheckState = "not-required"
)

// CheckRow is one local compatibility verdict.
type CheckRow struct {
	Check  string
	State  CheckState
	Action string
}

// LiveRow is one capability that only the selected interface can verify.
type LiveRow struct {
	Capability string
	Action     string
}

// ContextRow is one named compatibility observation.
type ContextRow struct {
	Field  string
	Value  string
	Source string
}

// Report separates observed context, local checks, and live obligations.
type Report struct {
	Context []ContextRow
	Checks  []CheckRow
	Live    []LiveRow
}

// Fingerprint identifies the compatibility observation context.
func Fingerprint(context Context) string {
	return fingerprint(
		string(context.Interface), context.Repository.Value, context.Repository.Source,
		context.Environment.Value, context.Environment.Source,
		context.ConfigurationHome.Value, context.ConfigurationHome.Source,
		context.ActiveRuntime.Value, context.ActiveRuntime.Source,
		context.PolicyProvenance.Value, context.PolicyProvenance.Source,
	)
}

// PolicyProvenance identifies inspected policy files without retaining their contents.
func PolicyProvenance(files ...FileFact) Fact {
	values, sources := []string{}, []string{}
	for _, file := range files {
		values = append(values, file.Path, string(file.State), file.digest)
		sources = append(sources, file.Source+": "+file.Path+" ("+string(file.State)+")")
	}
	return Fact{Value: fingerprint(values...), Source: strings.Join(sources, "; ")}
}

func fingerprint(values ...string) string {
	hash := sha256.New()
	for _, value := range values {
		fmt.Fprintf(hash, "%d:%s\n", len(value), value)
	}
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// Inspect derives the report from typed observations without reading or writing external state.
func Inspect(input Input) Report {
	report := Report{}
	appendContext := func(field string, fact Fact) {
		value := fact.Value
		if value == "" {
			value = "unknown"
		}
		report.Context = append(report.Context, ContextRow{Field: field, Value: value, Source: fact.Source})
	}
	appendContext("interface", Fact{Value: string(input.Context.Interface), Source: "argument"})
	appendContext("repository", input.Context.Repository)
	appendContext("execution-environment", input.Context.Environment)
	appendContext("configuration-home", input.Context.ConfigurationHome)
	appendContext("active-runtime", input.Context.ActiveRuntime)
	appendContext("launcher-version", input.Context.LauncherVersion)
	appendContext("policy-provenance", input.Context.PolicyProvenance)
	appendContext("compatibility-fingerprint", Fact{Value: Fingerprint(input.Context), Source: "derived"})

	report.Checks = append(report.Checks, configurationCheck(input.Configuration))
	for _, asset := range input.Assets {
		state, action := StateOK, "none"
		if asset.File.State != bounds.StateParsed {
			state, action = StateFailed, asset.RestoreAction
		}
		report.Checks = append(report.Checks, CheckRow{Check: "asset:" + asset.Name, State: state, Action: action})
	}
	if input.GlobalBench {
		report.Checks = append(report.Checks, CheckRow{Check: "global-bench", State: StateOK, Action: "none"})
	} else {
		report.Checks = append(report.Checks, CheckRow{Check: "global-bench", State: StateFailed, Action: input.ByPathAction})
	}
	if input.HookDeclared {
		report.Live = append(report.Live, LiveRow{Capability: "hook-behavior", Action: input.HookAction})
	}
	report.Live = append(report.Live, input.Live...)
	return report
}

func configurationCheck(file FileFact) CheckRow {
	row := CheckRow{Check: "effective-configuration", State: StateOK, Action: "none"}
	switch file.State {
	case bounds.StateParsed:
		row.State = StateUnknown
		row.Action = "observe effective settings in the selected interface; declaration: " + file.Source
		return row
	case bounds.StateAbsent:
		row.State, row.Action = StateUnknown, "configuration absent: "+file.Path
	case bounds.StateEmpty:
		row.State, row.Action = StateUnknown, "configuration empty: "+file.Path
	default:
		row.State = StateUnknown
		row.Action = fmt.Sprintf("inspect %s: %s", file.Source, file.Reason)
	}
	return row
}

// ExitCode reports whether every required local check passed.
func (r Report) ExitCode() int {
	for _, row := range r.Checks {
		if row.State == StateFailed || row.State == StateUnknown {
			return 1
		}
	}
	return 0
}

// Render emits the three compatibility tables through the shared TOON encoder.
func (r Report) Render() (string, error) {
	contextRows := make([][]string, 0, len(r.Context))
	for _, row := range r.Context {
		contextRows = append(contextRows, []string{row.Field, row.Value, row.Source})
	}
	checkRows := make([][]string, 0, len(r.Checks))
	for _, row := range r.Checks {
		checkRows = append(checkRows, []string{row.Check, string(row.State), row.Action})
	}
	liveRows := make([][]string, 0, len(r.Live))
	for _, row := range r.Live {
		liveRows = append(liveRows, []string{row.Capability, row.Action})
	}
	var output strings.Builder
	for _, table := range []struct {
		name   string
		fields []string
		rows   [][]string
	}{
		{"context", []string{"field", "value", "source"}, contextRows},
		{"checks", []string{"check", "state", "action"}, checkRows},
		{"live", []string{"capability", "action"}, liveRows},
	} {
		block, err := toon.Table(table.name, table.fields, table.rows)
		if err != nil {
			return "", err
		}
		output.WriteString(block)
	}
	return output.String(), nil
}
