package assessment

import (
	"encoding/json"
	"fmt"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
	"os"
	"sort"
	"strings"
)

const help = "usage: bench assessment list | show <run-id> | record --input <file>\n"

func Command(s Store, args []string) (string, int) {
	if len(args) == 0 {
		return help, 2
	}
	if len(args) == 1 && (args[0] == "help" || args[0] == "--help" || args[0] == "-h") {
		return help, 0
	}
	g := usage.Grammar{Cmd: "bench assessment " + args[0], Help: help, MaxArgs: 0}
	switch args[0] {
	case "list":
	case "show":
		g.MinArgs = 1
		g.MaxArgs = 1
	case "record":
		g.Flags = []usage.Flag{{Name: "--input", HasValue: true, NoEmptyValue: true, Required: true}}
	default:
		return toon.Usage("bench assessment", args[0]) + "\n", 2
	}
	parsed, line, code := usage.Parse(g, args[1:])
	if line != "" {
		return line, code
	}
	if s.Root == "" {
		return fail(fmt.Errorf("not in a git repository"))
	}
	if s.Home == "" {
		return fail(fmt.Errorf("BENCH_HOME is not set"))
	}
	switch args[0] {
	case "record":
		var r Run
		if err := readJSON(parsed.Flags["--input"], &r); err != nil {
			return fail(err)
		}
		if err := s.Record(r); err != nil {
			return fail(err)
		}
		path, _ := s.path(r.RunID)
		return table("record", []string{"run_id", "path"}, [][]string{{r.RunID, path}}, false)
	case "show":
		r, err := s.Read(parsed.Positionals[0])
		if err != nil {
			return fail(err)
		}
		return detail(r)
	default:
		if err := noLinks(s.Dir()); err != nil {
			return fail(err)
		}
		entries, err := os.ReadDir(s.Dir())
		if err != nil && !os.IsNotExist(err) {
			return fail(err)
		}
		var rows [][]string
		for _, entry := range entries {
			if !strings.HasSuffix(entry.Name(), ".json") {
				continue
			}
			r, err := s.Read(strings.TrimSuffix(entry.Name(), ".json"))
			if err != nil {
				return fail(err)
			}
			rows = append(rows, []string{r.RunID, r.State, fmt.Sprint(len(r.Attempts)), r.TaskID})
		}
		return table("runs", []string{"run_id", "state", "attempts", "task"}, rows, true)
	}
}

func fail(err error) (string, int) {
	return toon.Errorf("assessment refused", err.Error()) + "\nhelp[0]{cmd,why}:\n", 1
}
func table(name string, fields []string, rows [][]string, query bool) (string, int) {
	out, err := toon.Table(name, fields, rows)
	if err != nil {
		return fail(err)
	}
	if query {
		out += "help[0]{cmd,why}:\n"
	}
	return out, 0
}
func detail(r Run) (string, int) {
	summary, err := Summarize(r)
	if err != nil {
		return fail(err)
	}
	usages, err := RunUsage(r)
	if err != nil {
		return fail(err)
	}
	var rows [][]string
	for _, a := range r.Attempts {
		u := usages[a.AttemptID]
		cost, _ := estimateUsage(a, u)
		rows = append(rows, []string{a.AttemptID, a.ChunkID, a.Role, a.State, encoded(u), encoded(cost)})
	}
	out, code := table("attempts", []string{"attempt_id", "chunk", "role", "state", "usage", "cost"}, rows, false)
	if code != 0 {
		return out, code
	}
	more, code := table("summary", []string{"run_id", "metrics"}, [][]string{{r.RunID, encoded(summary)}}, false)
	if code != 0 {
		return more, code
	}
	out += more
	raw, err := json.Marshal(r)
	if err != nil {
		return fail(err)
	}
	object := jsonObject(raw)
	rows = nil
	flatten(&rows, "", object)
	more, code = table("record", []string{"field", "value"}, rows, true)
	return out + more, code
}
func encoded(v any) string { data, _ := json.Marshal(v); return string(data) }
func flatten(rows *[][]string, path string, v any) {
	switch item := v.(type) {
	case map[string]any:
		keys := make([]string, 0, len(item))
		for k := range item {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		if len(keys) == 0 {
			*rows = append(*rows, []string{path, "{}"})
		}
		for _, k := range keys {
			flatten(rows, path+"/"+k, item[k])
		}
	case []any:
		if len(item) == 0 {
			*rows = append(*rows, []string{path, "[]"})
		}
		for i, x := range item {
			flatten(rows, fmt.Sprintf("%s/%d", path, i), x)
		}
	default:
		*rows = append(*rows, []string{path, encoded(v)})
	}
}
