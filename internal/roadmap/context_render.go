package roadmap

import (
	"strings"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
)

func contextUsage() string { return "usage: bench roadmap --context [--full]\n" }

// ContextCommand implements the read-only schema-2 AXI roadmap snapshot.
func ContextCommand(args []string, gate func(string) GateCacheFact) (string, int) {
	if len(args) == 1 && (args[0] == "-h" || args[0] == "--help") {
		return contextUsage(), 0
	}
	full := len(args) == 2 && args[0] == "--context" && args[1] == "--full"
	if !(len(args) == 1 && args[0] == "--context") && !full {
		arg := "arguments"
		if len(args) > 0 {
			arg = args[0]
		}
		return toon.Usage("bench roadmap --context", arg) + "\n", 2
	}
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
	bs = append(bs, block{"context", []string{"schema", "full"}, [][]any{{2, s.Full}}})
	rows := make([][]any, len(s.Sources))
	for i, r := range s.Sources {
		rows[i] = []any{r.Source, r.State, r.Bytes}
	}
	bs = append(bs, block{"sources", []string{"source", "state", "bytes"}, rows})
	rows = nil
	for _, r := range s.Roadmap.Rows {
		rows = append(rows, []any{r.ID, r.Title, r.Spec, r.SpecStatus, r.ExternalTrigger, r.Body, r.BodyBytes, r.Truncated})
	}
	bs = append(bs, block{"roadmap_rows", []string{"id", "title", "spec", "spec_status", "external_trigger", "body", "body_bytes", "truncated"}, rows})
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
