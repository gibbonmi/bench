package roadmap

import (
	"strings"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

func contextUsage() string { return roadmapUsage() }

// contextGrammar is the declared argument shape usage.Parse enforces for the
// `bench roadmap --context` form — arity, flag recognition, `--`, repeated flags, and
// help all come from there rather than a local switch. Help is contextUsage without
// its trailing newline, because the caller appends one.
var contextGrammar = usage.Grammar{
	Cmd:   "bench roadmap --context",
	Help:  strings.TrimSuffix(contextUsage(), "\n"),
	Flags: []usage.Flag{{Name: "--context"}, {Name: "--full"}},
}

// ContextCommand implements the read-only schema-3 AXI roadmap snapshot.
func ContextCommand(args []string, gate func(string) GateCacheFact) (string, int) {
	parsed, line, code := usage.Parse(contextGrammar, args)
	if line != "" {
		return line + "\n", code
	}
	// --context is the mode selector, not an optional flag: --full alone names no mode
	// and is the misuse this reports. Arity, unknown flags, and a repeated --context are
	// already answered by the grammar above.
	if _, context := parsed.Flags["--context"]; !context {
		arg := "arguments"
		if len(args) > 0 {
			arg = args[0]
		}
		return toon.Usage("bench roadmap --context", arg) + "\n", 2
	}
	_, full := parsed.Flags["--full"]
	root, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}
	s, err := BuildContext(root, full, gate(root))
	if err != nil {
		return toon.Errorf("roadmap context failed", err.Error()) + "\n", 1
	}
	out, err := renderContext(s)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	return out, 0
}

func renderContext(s ContextSnapshot) (string, error) {
	type block struct {
		name   string
		fields []string
		rows   [][]any
	}
	var bs []block
	bs = append(bs, block{"context", []string{"schema", "full", "sequence_trusted"}, [][]any{{3, s.Full, s.SequenceTrusted}}})
	rows := make([][]any, len(s.Sources))
	for i, r := range s.Sources {
		rows[i] = []any{r.Source, r.State, r.Bytes}
	}
	bs = append(bs, block{"sources", []string{"source", "state", "bytes"}, rows})
	rows = nil
	for _, r := range s.Roadmap.Rows {
		rows = append(rows, []any{r.ID, r.Title, r.Spec, r.SpecStatus, r.ExternalTrigger, r.Body, r.BodyBytes, r.Truncated, r.OccurrenceCount, r.OccurrenceKeys})
	}
	bs = append(bs, block{"roadmap_rows", []string{"id", "title", "spec", "spec_status", "external_trigger", "body", "body_bytes", "truncated", "occurrence_count", "occurrence_keys"}, rows})
	rows = nil
	for _, r := range s.Roadmap.Sequence {
		rows = append(rows, []any{r.Rank, r.Text, r.Command})
	}
	bs = append(bs, block{"roadmap_sequence", []string{"rank", "text", "command"}, rows})
	rows = nil
	for _, r := range s.Ideas {
		rows = append(rows, []any{r.Date, r.Text, r.TextBytes, r.Truncated})
	}
	bs = append(bs, block{"ideas", []string{"date", "text", "text_bytes", "truncated"}, rows})
	rows = nil
	for _, r := range s.Learnings {
		rows = append(rows, []any{r.Date, r.Title, r.State, r.Body, r.BodyBytes, r.Truncated})
	}
	bs = append(bs, block{"learnings", []string{"date", "title", "state", "body", "body_bytes", "truncated"}, rows})
	rows = nil
	for _, r := range s.Retros {
		rows = append(rows, []any{r.Path, r.State, r.Body, r.BodyBytes, r.Truncated})
	}
	bs = append(bs, block{"retros", []string{"path", "state", "body", "body_bytes", "truncated"}, rows})
	rows = nil
	for _, r := range s.CaptureOccurrences {
		rows = append(rows, []any{r.Owner, r.Incident, r.Source, r.CaptureUnit, r.State})
	}
	bs = append(bs, block{"capture_occurrences", []string{"owner", "incident", "source", "capture_unit", "state"}, rows})
	rows = nil
	for _, r := range s.Roadmap.OccurrenceDiscrepancies {
		rows = append(rows, []any{r.Source, r.CaptureUnit, r.Kind, r.Owner, r.Incident, r.Structural})
	}
	bs = append(bs, block{"occurrence_discrepancies", []string{"source", "capture_unit", "kind", "owner", "incident", "structural"}, rows})
	bs = append(bs, block{"structure", []string{"kind", "path", "actual", "limit", "state", "detail"}, s.Structure})
	rows = stringRows(s.Specs)
	bs = append(bs, block{"specs", []string{"slug", "status", "roadmap_id"}, rows})
	rows = stringRows(s.SpecHistory)
	bs = append(bs, block{"spec_history", []string{"slug", "hash", "date", "kind", "subject"}, rows})
	bs = append(bs, block{"git", []string{"branch", "default_branch", "dirty", "ahead", "behind"}, s.Git})
	bs = append(bs, block{"git_changes", []string{"status", "path"}, stringRows(s.GitChanges)})
	bs = append(bs, block{"gate_cache", []string{"present", "state", "pending_status", "status", "cached_tree", "work_tree", "timestamp", "stale"}, s.GateCache})
	rows = nil
	for _, r := range s.Failures {
		rows = append(rows, []any{r.Source, r.Reason, r.Raw, r.RawBytes, r.Truncated})
	}
	bs = append(bs, block{"parse_failures", []string{"source", "reason", "raw", "raw_bytes", "truncated"}, rows})
	bs = append(bs, block{"help", []string{"cmd", "why"}, nil})
	var out strings.Builder
	for _, b := range bs {
		x, e := toon.TableTyped(b.name, b.fields, b.rows)
		if e != nil {
			return "", e
		}
		out.WriteString(x)
	}
	return out.String(), nil
}

func stringRows(in [][]string) [][]any {
	out := make([][]any, len(in))
	for i, r := range in {
		out[i] = make([]any, len(r))
		for j, v := range r {
			out[i][j] = v
		}
	}
	return out
}
